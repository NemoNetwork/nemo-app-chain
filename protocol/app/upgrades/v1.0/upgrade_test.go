package v_1_0_test

import (
	"math/big"
	"testing"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/nemo-network/v4-chain/protocol/dtypes"
	"github.com/nemo-network/v4-chain/protocol/lib"
	testapp "github.com/nemo-network/v4-chain/protocol/testutil/app"
	testutil "github.com/nemo-network/v4-chain/protocol/testutil/util"
	satypes "github.com/nemo-network/v4-chain/protocol/x/subaccounts/types"
	vaultkeeper "github.com/nemo-network/v4-chain/protocol/x/vault/keeper"
	vaulttypes "github.com/nemo-network/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"

	v_1_0 "github.com/nemo-network/v4-chain/protocol/app/upgrades/v1.0"
)

// setupPreUpgradeChain builds a chain that looks like one running from a genesis
// that predates the megavault work: megavault holds funds and has shares
// outstanding, but no operator params and no fee state.
func setupPreUpgradeChain(
	t *testing.T,
	equity *big.Int,
	totalShares *big.Int,
) (sdk.Context, vaultkeeper.Keeper) {
	t.Helper()

	tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
		genesis = testapp.DefaultGenesis()
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *satypes.GenesisState) {
				genesisState.Subaccounts = []satypes.Subaccount{
					{
						Id: &vaulttypes.MegavaultMainSubaccount,
						AssetPositions: []*satypes.AssetPosition{
							testutil.CreateSingleAssetPosition(0, equity),
						},
					},
				}
			},
		)
		testapp.UpdateGenesisDocWithAppStateForModule(
			&genesis,
			func(genesisState *vaulttypes.GenesisState) {
				genesisState.TotalShares = vaulttypes.BigIntToNumShares(totalShares)
				genesisState.OwnerShares = []vaulttypes.OwnerShare{
					{
						Owner:  "nemo1c0m5x87llaunl5sgv3q5vd7j5uha26d2z2yp69",
						Shares: vaulttypes.BigIntToNumShares(totalShares),
					},
				}
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	// Rewind the state to what a pre-upgrade chain would hold: no operator, no
	// fee state. `InitGenesis` sets both, so they have to be cleared explicitly
	// to model the chain this upgrade actually runs against.
	store := ctx.KVStore(tApp.App.GetKey(vaulttypes.StoreKey))
	store.Delete([]byte(vaulttypes.OperatorParamsKey))
	store.Delete([]byte(vaulttypes.MegavaultParamsKey))
	store.Delete([]byte(vaulttypes.MegavaultFeeStateKey))

	return ctx, k
}

// TestUpgrade_SeedsHighWaterMarkAtUpgradeHeight is the test that matters most:
// a zero high-water mark would charge a profit share on the vault's entire
// history, retroactively, against depositors who deposited under a 0%-fee regime.
func TestUpgrade_SeedsHighWaterMarkAtUpgradeHeight(t *testing.T) {
	// NAV per share of 3.0 at upgrade height.
	ctx, k := setupPreUpgradeChain(t, big.NewInt(3_000), big.NewInt(1_000))

	require.Equal(t, vaulttypes.FeeState{}, k.GetFeeState(ctx))

	v_1_0.MigrateVaultState(ctx, k)

	feeState := k.GetFeeState(ctx)
	expectedHwm := new(big.Int).Mul(big.NewInt(3), vaulttypes.NavPerShareScale)
	require.Equal(t, expectedHwm, feeState.HighWaterMarkNavPerShare.BigInt())
	require.Equal(t, ctx.BlockTime().Unix(), feeState.LastAccrualTime)

	// With the mark seeded at the pre-upgrade NAV, turning the profit share on
	// immediately charges nothing until the vault makes new gains.
	require.NoError(t, k.SetOperatorParams(ctx, vaulttypes.OperatorParams{
		Operator:       lib.GovModuleAddress.String(),
		ProfitSharePpm: 200_000,
	}))
	sharesBefore := k.GetTotalShares(ctx).NumShares.BigInt()
	require.NoError(t, k.AccrueFees(ctx))
	require.Equal(t, sharesBefore, k.GetTotalShares(ctx).NumShares.BigInt())
}

// TestUpgrade_InitializesOperatorAndMegavaultParams covers the state that
// `InitGenesis` would have set and a running chain will not have.
func TestUpgrade_InitializesOperatorAndMegavaultParams(t *testing.T) {
	ctx, k := setupPreUpgradeChain(t, big.NewInt(1_000), big.NewInt(1_000))

	require.Equal(t, "", k.GetOperatorParams(ctx).Operator)

	v_1_0.MigrateVaultState(ctx, k)

	operatorParams := k.GetOperatorParams(ctx)
	require.Equal(t, lib.GovModuleAddress.String(), operatorParams.Operator)

	// The fee engine ships inert; governance turns it on separately.
	require.Equal(t, uint32(0), operatorParams.OperatorFeePpm)
	require.Equal(t, uint32(0), operatorParams.ProfitSharePpm)
	require.Equal(t, uint32(0), operatorParams.MinOperatorSharePpm)

	// Deposits stay uncapped until a cap is chosen.
	require.Equal(t, 0, k.GetMegavaultParams(ctx).DepositCapQuoteQuantums.Sign())
}

// TestUpgrade_DoesNotOverwriteExistingState checks that the upgrade is safe to
// run against state that already carries these values — for instance after a
// genesis export and re-import.
func TestUpgrade_DoesNotOverwriteExistingState(t *testing.T) {
	ctx, k := setupPreUpgradeChain(t, big.NewInt(3_000), big.NewInt(1_000))

	chosenOperator := "nemo1c0m5x87llaunl5sgv3q5vd7j5uha26d2z2yp69"
	require.NoError(t, k.SetOperatorParams(ctx, vaulttypes.OperatorParams{
		Operator:       chosenOperator,
		OperatorFeePpm: 20_000,
	}))
	require.NoError(t, k.SetMegavaultParams(ctx, vaulttypes.MegavaultParams{
		DepositCapQuoteQuantums: dtypes.NewInt(123_456),
	}))
	// A high-water mark above the current NAV must survive: lowering it would
	// let the operator charge a profit share on a recovery.
	highMark := new(big.Int).Mul(big.NewInt(10), vaulttypes.NavPerShareScale)
	require.NoError(t, k.SetFeeState(ctx, vaulttypes.FeeState{
		HighWaterMarkNavPerShare: dtypes.NewIntFromBigInt(highMark),
		LastAccrualTime:          99,
	}))

	v_1_0.MigrateVaultState(ctx, k)

	require.Equal(t, chosenOperator, k.GetOperatorParams(ctx).Operator)
	require.Equal(t, uint32(20_000), k.GetOperatorParams(ctx).OperatorFeePpm)
	require.Equal(t, big.NewInt(123_456), k.GetMegavaultParams(ctx).DepositCapQuoteQuantums.BigInt())
	require.Equal(t, highMark, k.GetFeeState(ctx).HighWaterMarkNavPerShare.BigInt())
	require.Equal(t, int64(99), k.GetFeeState(ctx).LastAccrualTime)
}

// TestUpgrade_HandlesEmptyMegavault checks the no-shares case, where NAV per
// share is undefined and must not be treated as zero.
func TestUpgrade_HandlesEmptyMegavault(t *testing.T) {
	ctx, k := setupPreUpgradeChain(t, big.NewInt(0), big.NewInt(0))

	require.NotPanics(t, func() {
		v_1_0.MigrateVaultState(ctx, k)
	})

	feeState := k.GetFeeState(ctx)
	require.Equal(t, 0, feeState.HighWaterMarkNavPerShare.Sign())
	require.Equal(t, ctx.BlockTime().Unix(), feeState.LastAccrualTime)
}
