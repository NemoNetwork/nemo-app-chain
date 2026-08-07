package keeper_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	"github.com/nemo-network/v4-chain/protocol/dtypes"
	testapp "github.com/nemo-network/v4-chain/protocol/testutil/app"
	"github.com/nemo-network/v4-chain/protocol/testutil/constants"
	testutil "github.com/nemo-network/v4-chain/protocol/testutil/util"
	satypes "github.com/nemo-network/v4-chain/protocol/x/subaccounts/types"
	vaulttypes "github.com/nemo-network/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

// TestValidateDepositAgainstCap covers the cap arithmetic in isolation, including
// the "zero means uncapped" rule that a chain upgraded from pre-cap state relies on.
func TestValidateDepositAgainstCap(t *testing.T) {
	tests := map[string]struct {
		// Megavault equity before the deposit.
		equity *big.Int
		// Deposit cap to set, or nil to leave megavault params untouched.
		depositCap *big.Int
		// Quote quantums to deposit.
		quoteQuantums *big.Int

		// Whether the deposit is expected to be rejected.
		expectedErr error
	}{
		"No cap set - deposit allowed": {
			equity:        big.NewInt(1_000),
			quoteQuantums: big.NewInt(999_999_999),
		},
		"Cap of zero is uncapped - deposit allowed": {
			equity:        big.NewInt(1_000),
			depositCap:    big.NewInt(0),
			quoteQuantums: big.NewInt(999_999_999),
		},
		"Deposit lands below cap": {
			equity:        big.NewInt(1_000),
			depositCap:    big.NewInt(5_000),
			quoteQuantums: big.NewInt(3_999),
		},
		"Deposit lands exactly on cap": {
			equity:        big.NewInt(1_000),
			depositCap:    big.NewInt(5_000),
			quoteQuantums: big.NewInt(4_000),
		},
		"Deposit exceeds cap by one quantum": {
			equity:        big.NewInt(1_000),
			depositCap:    big.NewInt(5_000),
			quoteQuantums: big.NewInt(4_001),
			expectedErr:   vaulttypes.ErrDepositCapExceeded,
		},
		"Equity already above cap - any deposit rejected": {
			equity:        big.NewInt(10_000),
			depositCap:    big.NewInt(5_000),
			quoteQuantums: big.NewInt(1),
			expectedErr:   vaulttypes.ErrDepositCapExceeded,
		},
		"Cap counts existing equity, not just the deposit": {
			equity: big.NewInt(4_999),
			// A deposit of 2 alone is far below the cap, but equity after the
			// deposit is 5_001, which is not.
			depositCap:    big.NewInt(5_000),
			quoteQuantums: big.NewInt(2),
			expectedErr:   vaulttypes.ErrDepositCapExceeded,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				// Fund the megavault main subaccount so that megavault equity
				// equals `tc.equity`.
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *satypes.GenesisState) {
						genesisState.Subaccounts = []satypes.Subaccount{
							{
								Id: &vaulttypes.MegavaultMainSubaccount,
								AssetPositions: []*satypes.AssetPosition{
									testutil.CreateSingleAssetPosition(0, tc.equity),
								},
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper

			if tc.depositCap != nil {
				require.NoError(t, k.SetMegavaultParams(ctx, vaulttypes.MegavaultParams{
					DepositCapQuoteQuantums: dtypes.NewIntFromBigInt(tc.depositCap),
				}))
			}

			err := k.ValidateDepositAgainstCap(ctx, tc.quoteQuantums)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestDepositToMegavaultRespectsCap checks that the cap is enforced on the deposit
// path itself, and that a rejected deposit mints no shares and moves no funds.
func TestDepositToMegavaultRespectsCap(t *testing.T) {
	depositor := constants.Alice_Num0
	depositorBalance := big.NewInt(10_000)

	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *satypes.GenesisState) {
				genesisState.Subaccounts = []satypes.Subaccount{
					{
						Id: &depositor,
						AssetPositions: []*satypes.AssetPosition{
							testutil.CreateSingleAssetPosition(0, depositorBalance),
						},
					},
				}
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Cap megavault at 1_000 quote quantums.
	require.NoError(t, k.SetMegavaultParams(ctx, vaulttypes.MegavaultParams{
		DepositCapQuoteQuantums: dtypes.NewInt(1_000),
	}))

	// A deposit that fits under the cap succeeds.
	mintedShares, err := k.DepositToMegavault(ctx, depositor, big.NewInt(600))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(600), mintedShares)
	require.Equal(t, big.NewInt(600), k.GetTotalShares(ctx).NumShares.BigInt())

	// A second deposit that would breach the cap is rejected, and nothing moves.
	_, err = k.DepositToMegavault(ctx, depositor, big.NewInt(401))
	require.ErrorIs(t, err, vaulttypes.ErrDepositCapExceeded)
	require.Equal(t, big.NewInt(600), k.GetTotalShares(ctx).NumShares.BigInt())
	ownerShares, exists := k.GetOwnerShares(ctx, depositor.Owner)
	require.True(t, exists)
	require.Equal(t, big.NewInt(600), ownerShares.NumShares.BigInt())

	// A deposit that exactly reaches the cap succeeds.
	_, err = k.DepositToMegavault(ctx, depositor, big.NewInt(400))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1_000), k.GetTotalShares(ctx).NumShares.BigInt())

	// Raising the cap re-opens deposits.
	require.NoError(t, k.SetMegavaultParams(ctx, vaulttypes.MegavaultParams{
		DepositCapQuoteQuantums: dtypes.NewInt(2_000),
	}))
	_, err = k.DepositToMegavault(ctx, depositor, big.NewInt(500))
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1_500), k.GetTotalShares(ctx).NumShares.BigInt())
}
