package queryengine

import (
	"reflect"
	"testing"

	inflation "github.com/cosmos/cosmos-sdk/x/mint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func TestGetStakingAPR(t *testing.T) {
	type args struct {
		pool          staking.QueryPoolResponse
		mintProvision inflation.QueryAnnualProvisionsResponse
	}
	tests := []struct {
		name string
		args args
		want sdk.Dec
	}{
		{
			name: "test bonded tokens are zero",
			args: args{
				pool: staking.QueryPoolResponse{
					Pool: staking.Pool{
						BondedTokens: sdk.ZeroInt(),
					},
				},
				mintProvision: inflation.QueryAnnualProvisionsResponse{
					EpochMintProvision: sdk.NewDecCoin("acanto", sdk.NewInt(100)),
				},
			},
			want: sdk.NewDec(0),
		},
		{
			name: "mint provision is zero",
			args: args{
				pool: staking.QueryPoolResponse{
					Pool: staking.Pool{
						BondedTokens: sdk.NewInt(100),
					},
				},
				mintProvision: inflation.QueryAnnualProvisionsResponse{
					EpochMintProvision: sdk.NewDecCoin("acanto", sdk.ZeroInt()),
				},
			},
			want: sdk.NewDec(0),
		},
		{
			name: "bonded tokens is less than mint provision",
			args: args{
				pool: staking.QueryPoolResponse{
					Pool: staking.Pool{
						BondedTokens: sdk.NewInt(100),
					},
				},
				mintProvision: inflation.QueryAnnualProvisionsResponse{
					AnnualProvisions: types.NewDecWithPrec(100, 0)
				},
			},
			want: sdk.NewDec(36500000000000),
		},
		{
			name: "mint provision is less than bonded tokens",
			args: args{
				pool: staking.QueryPoolResponse{
					Pool: staking.Pool{
						BondedTokens: sdk.NewInt(100000000000),
					},
				},
				mintProvision: inflation.QueryAnnualProvisionsResponse{
					EpochMintProvision: sdk.NewDecCoin("acanto", sdk.NewInt(100)),
				},
			},
			want: sdk.MustNewDecFromStr("0.0000365"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateStakingAPR(tt.args.pool, tt.args.mintProvision); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStakingAPR() = %v, want %v", got, tt.want)
			}
		})
	}
}
