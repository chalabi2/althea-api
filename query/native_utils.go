package query

import (
	inflation "github.com/Canto-Network/Canto/v6/x/inflation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// GetStakingAPR returns the APR for a all bonded tokens and mint provision for current epoch
func GetStakingAPR(pool staking.QueryPoolResponse, mintProvision inflation.QueryEpochMintProvisionResponse) sdk.Dec {
	//get bonded tokens from pool
	bondedTokens := pool.GetPool().BondedTokens
	//get mint provision amount from epoch (in acanto)
	mintProvisionAmount := mintProvision.GetEpochMintProvision().Amount

	//check if bonded tokens are zero so we don't divide by zero
	if bondedTokens.IsZero() {
		return sdk.NewDec(0)
	}

	//calculate apr (mint provision / bonded tokens) * 365 (days) * 100%
	return mintProvisionAmount.Mul(sdk.NewDec(36500)).QuoInt(bondedTokens)
}
