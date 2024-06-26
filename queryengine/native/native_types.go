package queryengine

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	csr "github.com/Canto-Network/Canto/v6/x/csr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	query "github.com/cosmos/cosmos-sdk/types/query"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types"
	inflation "github.com/cosmos/cosmos-sdk/x/mint/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/rs/zerolog/log"
)

// STAKING
type Validator struct {
	// operator_address defines the address of the validator's operator; bech encoded in JSON.
	OperatorAddress string `json:"operator_address"`
	// jailed defined whether the validator has been jailed from bonded status or not.
	Jailed bool `json:"jailed"`
	// status defines the validator's status (bonded(3)/unbonding(2)/unbonded(1)).
	Status string `json:"status"`
	// tokens defines the amount of staking tokens delegated to the validator.
	Tokens string `json:"tokens"`
	// description of validator includes moniker, identity, website, security contact, and details.
	Description staking.Description `json:"description"`
	// commission defines the commission rate.
	Commission string `json:"commission"`
}

// get all Validators for staking
// will return full response string and mapping of operator address to response string
func GetValidators(ctx context.Context, queryClient staking.QueryClient) ([]Validator, map[string]string, error) {
	respValidators, err := queryClient.Validators(ctx, &staking.QueryValidatorsRequest{
		Pagination: &query.PageRequest{
			Limit: 1000,
		},
	})
	if err != nil {
		return nil, nil, err
	}
	allValidators := new([]Validator)
	validatorMap := make(map[string]string)
	for _, validator := range respValidators.Validators {
		valResponse := Validator{
			OperatorAddress: validator.OperatorAddress,
			Jailed:          validator.Jailed,
			Status:          validator.Status.String(),
			Tokens:          validator.Tokens.String(),
			Description:     validator.Description,
			Commission:      validator.Commission.CommissionRates.Rate.String(),
		}
		*allValidators = append(*allValidators, valResponse)
		validatorMap[validator.OperatorAddress] = GeneralResultToString(valResponse)
	}
	return *allValidators, validatorMap, nil
}
func GetStakingAPR(ctx context.Context, stakingQueryClient staking.QueryClient, inflationQueryClient inflation.QueryClient) (sdk.Dec, error) {
    // Fetch pool information
    poolResp, err := stakingQueryClient.Pool(ctx, &staking.QueryPoolRequest{})
    if err != nil {
        return sdk.Dec{}, err // Return error if fetching pool fails
    }
    
    // Fetch mint provision information
    mintProvisionResp, err := inflationQueryClient.AnnualProvisions(ctx, &inflation.QueryAnnualProvisionsRequest{})
    if err != nil {
        return sdk.Dec{}, err // Return error if fetching mint provision fails
    }
    
    bondedTokensDec := sdk.NewDecFromInt(poolResp.Pool.BondedTokens)
    if bondedTokensDec.IsZero() {
        return sdk.NewDec(0), nil // Avoid division by zero, not necessarily an error
    }
    
    mintProvisionAmount := mintProvisionResp.AnnualProvisions
    apr := mintProvisionAmount.Quo(bondedTokensDec).MulInt64(100) // Adjusted to multiply by 100 for percentage
    
    return apr, nil // Return calculated APR and no error
}

// USER DELEGATIONS

type DelegationResponse struct {
	Delegations          []DelegationInfo      `json:"delegations"`
	UnbondingDelegations []UnbondingDelegation `json:"unbondingDelegations"`
	Rewards              RewardsInfo           `json:"rewards"`
}


// DelegationInfo holds information about a single delegation.
type DelegationInfo struct {
	Delegation Delegation `json:"delegation"`
	Balance    Balance    `json:"balance"`
}

// Delegation holds the delegator and validator addresses and the amount of shares.
type Delegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`
	Shares           string `json:"shares"`
}

// Balance holds the denomination and amount of tokens.
type Balance struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

// RewardsInfo holds information about rewards.
type RewardsInfo struct {
    Rewards []ValidatorReward `json:"rewards"`
    Total   []Balance         `json:"total"`
}

type UnbondingDelegation struct {
    DelegatorAddress   string `json:"delegator_address"`
    ValidatorAddress   string `json:"validator_address"`
    CreationHeight     int64  `json:"creation_height"`
    CompletionTime     time.Time `json:"completion_time"`
    InitialBalance     string `json:"initial_balance"`
    Balance            string `json:"balance"`
}

func (r RewardsInfo) MarshalJSON() ([]byte, error) {
    type Alias RewardsInfo
    if r.Rewards == nil {
        r.Rewards = []ValidatorReward{} // Ensure an empty slice, not nil
    }
    if r.Total == nil {
        r.Total = []Balance{} // Ensure an empty slice, not nil
    }
    return json.Marshal((Alias)(r))
}

// ValidatorReward holds information about rewards from a specific validator.
type ValidatorReward struct {
	ValidatorAddress string    `json:"validator_address"`
	Reward           []Balance `json:"reward"`
}



// FetchUserDelegations fetches delegations for a specific user.
func FetchUserDelegations(ctx context.Context, stakingQueryClient staking.QueryClient, distributionQueryClient distrtypes.QueryClient, delegatorAddress string) (*DelegationResponse, error) {
    response := &DelegationResponse{}

    // Fetch delegations
    delegationResp, err := stakingQueryClient.DelegatorDelegations(ctx, &staking.QueryDelegatorDelegationsRequest{
        DelegatorAddr: delegatorAddress,
        Pagination: &query.PageRequest{Limit: 100}, // Adjust pagination as needed
    })
    if err != nil {
        return nil, fmt.Errorf("failed to fetch delegations: %w", err)
    }

    // Handle delegations response
    for _, del := range delegationResp.DelegationResponses {
        response.Delegations = append(response.Delegations, DelegationInfo{
            Delegation: Delegation{
                DelegatorAddress: del.Delegation.DelegatorAddress,
                ValidatorAddress: del.Delegation.ValidatorAddress,
                Shares:           del.Delegation.Shares.String(),
            },
            Balance: Balance{
                Denom:  del.Balance.Denom,
                Amount: del.Balance.Amount.String(),
            },
        })
    }

    // Fetch unbonding delegations
    unbondingResp, err := stakingQueryClient.DelegatorUnbondingDelegations(ctx, &staking.QueryDelegatorUnbondingDelegationsRequest{
        DelegatorAddr: delegatorAddress,
        Pagination: &query.PageRequest{Limit: 100}, // Adjust pagination as needed
    })
    if err != nil {
        return nil, fmt.Errorf("failed to fetch unbonding delegations: %w", err)
    }

    // Handle unbonding delegations response
    for _, unbond := range unbondingResp.UnbondingResponses {
        for _, entry := range unbond.Entries {
            response.UnbondingDelegations = append(response.UnbondingDelegations, UnbondingDelegation{
                DelegatorAddress: unbond.DelegatorAddress,
                ValidatorAddress: unbond.ValidatorAddress,
                CreationHeight:   entry.CreationHeight,
                CompletionTime:   entry.CompletionTime,
                InitialBalance:   entry.InitialBalance.String(),
                Balance:          entry.Balance.String(),
            })
        }
    }

    // If no unbonding delegations, set to nil explicitly
    if len(response.UnbondingDelegations) == 0 {
        response.UnbondingDelegations = nil
    }

    // Fetch rewards
    rewardsResp, err := distributionQueryClient.DelegationTotalRewards(ctx, &distrtypes.QueryDelegationTotalRewardsRequest{
        DelegatorAddress: delegatorAddress,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to fetch rewards: %w", err)
    }

    // Handle rewards response
    for _, reward := range rewardsResp.Rewards {
        var validatorRewards []Balance
        for _, valReward := range reward.Reward {
            validatorRewards = append(validatorRewards, Balance{
                Denom:  valReward.Denom,
                Amount: valReward.Amount.String(),
            })
        }
        response.Rewards.Rewards = append(response.Rewards.Rewards, ValidatorReward{
            ValidatorAddress: reward.ValidatorAddress,
            Reward:           validatorRewards,
        })
    }

    // Calculate total rewards
    var total []Balance
    for _, reward := range response.Rewards.Rewards {
        for _, valReward := range reward.Reward {
            found := false
            for i, totalReward := range total {
                if totalReward.Denom == valReward.Denom {
                    amount, _ := strconv.ParseFloat(totalReward.Amount, 64)
                    rewardAmount, _ := strconv.ParseFloat(valReward.Amount, 64)
                    total[i].Amount = fmt.Sprintf("%f", amount+rewardAmount)
                    found = true
                    break
                }
            }
            if !found {
                total = append(total, valReward)
            }
        }
    }
    response.Rewards.Total = total

    return response, nil
}

// GOVSHUTTLE

type Proposal struct {
	// proposalId defines the unique id of the proposal.
	ProposalId uint64 `json:"proposal_id"`
	// typeUrl indentifies the type of the proposal by a serialized protocol buffer message
	TypeUrl string `json:"type_url"`
	// title of the proposal
	Title string `json:"title"`
	// description of the proposal
	Description string `json:"description"`
	// status defines the current status of the proposal.
	Status string `json:"status"`
	// finalVote defined the result of the proposal
	FinalVote gov.TallyResult `json:"final_vote"`
	// submitTime defines the block time the proposal was submitted.
	SubmitTime time.Time `json:"submit_time"`
	// depositEndTime defines the time when the proposal deposit period will end.
	DepositEndTime time.Time `json:"deposit_end_time"`
	// totalDeposit defines the total amount of coins deposited on this proposal
	TotalDeposit sdk.Coins `json:"total_deposit"`
	// votingStartTime defines the time when the proposal voting period will start
	VotingStartTime time.Time `json:"voting_start_time"`
	// votingEndTime defines the time when the proposal voting period will end
	VotingEndTime time.Time `json:"voting_end_time"`
}

// get all proposals from gov shuttle
// will return full response string and mapping of proposal id to response string
func GetAllProposals(ctx context.Context, queryClient gov.QueryClient) ([]Proposal, map[string]string, error) {
	resp, err := queryClient.Proposals(ctx, &gov.QueryProposalsRequest{
		Pagination: &query.PageRequest{
			Limit: 1000,
		},
	})
	if err != nil {
		return nil, nil, err
	}
	allProposals := new([]Proposal)
	proposalMap := make(map[string]string)
	for _, proposal := range resp.GetProposals() {
		// deal with votes
		var votes gov.TallyResult
		// if vote is still ongoing, query the current tally
		if proposal.Status == 2 {
			resp, err := queryClient.TallyResult(ctx, &gov.QueryTallyResultRequest{
				ProposalId: proposal.ProposalId,
			})
			if err == nil {
				votes = resp.Tally
			}
		} else {
			votes = proposal.FinalTallyResult
		}

		// get proposal metadata
		title := ""
		description := ""
		metadata, err := GetProposalMetadata(proposal.Content)
		if err != nil {
			log.Log().Msgf("Error getting proposal metadata: %v", err)
		} else {
			title = metadata.Title
			description = metadata.Description
		}

		proposalResponse := Proposal{
			ProposalId:      proposal.ProposalId,
			TypeUrl:         proposal.Content.TypeUrl,
			Title:           title,
			Description:     description,
			Status:          proposal.Status.String(),
			FinalVote:       votes,
			SubmitTime:      proposal.SubmitTime,
			DepositEndTime:  proposal.DepositEndTime,
			TotalDeposit:    proposal.TotalDeposit,
			VotingStartTime: proposal.VotingStartTime,
			VotingEndTime:   proposal.VotingEndTime,
		}
		*allProposals = append(*allProposals, proposalResponse)
		proposalMap[strconv.Itoa(int(proposal.ProposalId))] = GeneralResultToString(proposalResponse)
	}
	return *allProposals, proposalMap, nil
}

// CSR
type CSR struct {
	// ID of the CSR
	Id uint64 `json:"id"`
	// all contracts under this csr id
	Contracts []string `json:"contracts"`
	// total number of transactions under this csr id
	Txs uint64 `json:"txs"`
	// The cumulative revenue for this CSR NFT -> represented as a big.Int
	Revenue string `json:"revenue"`
}

// get all CSRS
// will return full response string and mapping of nft id to response string
func GetCSRS(ctx context.Context, queryClient csr.QueryClient) ([]CSR, map[string]string, error) {
	resp, err := queryClient.CSRs(ctx, &csr.QueryCSRsRequest{Pagination: &query.PageRequest{
		Limit: 1000,
	}})
	if err != nil {
		return nil, nil, err
	}
	allCsrs := new([]CSR)
	csrMap := make(map[string]string)
	for _, csr := range resp.GetCsrs() {
		csrResponse := CSR{
			Id:        csr.GetId(),
			Contracts: csr.GetContracts(),
			Txs:       csr.GetTxs(),
			Revenue:   csr.Revenue.String(),
		}
		*allCsrs = append(*allCsrs, csrResponse)
		csrMap[strconv.Itoa(int(csr.GetId()))] = GeneralResultToString(csrResponse)
	}
	return *allCsrs, csrMap, nil
}
