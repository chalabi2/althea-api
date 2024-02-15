package queryengine

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	inflation "github.com/cosmos/cosmos-sdk/x/mint/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// CalculateStakingAPR calculates the APR based on mint provision and total bonded tokens
func CalculateStakingAPR(pool *staking.QueryPoolResponse, mintProvision *inflation.QueryAnnualProvisionsResponse) sdk.Dec {
    // Convert bondedTokens to sdk.Dec
    bondedTokensDec := sdk.NewDecFromInt(pool.Pool.BondedTokens)

    // Ensure bondedTokensDec is not zero to avoid division by zero
    if bondedTokensDec.IsZero() {
        return sdk.NewDec(0)
    }

    // mintProvisionAmount is already of type sdk.Dec
    mintProvisionAmount := mintProvision.AnnualProvisions

    // Perform the calculation: (mintProvision / totalStake) * 100
    apr := mintProvisionAmount.Quo(bondedTokensDec)

    return apr
}

func GeneralResultToString(results interface{}) string {
	ret, err := json.Marshal(results)
	if err != nil {
		return "GeneralResultToString:" + err.Error()
	}
	return string(ret)
}
