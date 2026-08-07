package keeper_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/nemo-network/v4-chain/protocol/dtypes"
	testapp "github.com/nemo-network/v4-chain/protocol/testutil/app"
	"github.com/nemo-network/v4-chain/protocol/testutil/constants"
	testutil "github.com/nemo-network/v4-chain/protocol/testutil/util"
	satypes "github.com/nemo-network/v4-chain/protocol/x/subaccounts/types"
	"github.com/nemo-network/v4-chain/protocol/x/vault/keeper"
	vaulttypes "github.com/nemo-network/v4-chain/protocol/x/vault/types"
	"github.com/stretchr/testify/require"
)

// setupMegavault builds a test app whose megavault main subaccount holds `equity`
// quote quantums and whose total shares are `totalShares`, all owned by
// `shareOwner`.
func setupMegavault(
	t *testing.T,
	equity *big.Int,
	totalShares *big.Int,
	shareOwner string,
) (*testapp.TestApp, sdk.Context, keeper.Keeper) {
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
						Owner:  shareOwner,
						Shares: vaulttypes.BigIntToNumShares(totalShares),
					},
				}
			},
		)
		return genesis
	}).Build()
	ctx := tApp.InitChain()

	return tApp, ctx, tApp.App.VaultKeeper
}

// fundMegavault changes megavault equity by `delta` quote quantums by rewriting
// the main subaccount's USDC position directly. Deposits and withdrawals go
// through the fee engine, so they cannot be used to set up a fee test.
func fundMegavault(
	t *testing.T,
	tApp *testapp.TestApp,
	ctx sdk.Context,
	delta *big.Int,
) {
	t.Helper()

	sa := tApp.App.SubaccountsKeeper.GetSubaccount(ctx, vaulttypes.MegavaultMainSubaccount)
	quantums := big.NewInt(0)
	if len(sa.AssetPositions) > 0 {
		quantums = sa.AssetPositions[0].Quantums.BigInt()
	}
	quantums.Add(quantums, delta)

	sa.Id = &vaulttypes.MegavaultMainSubaccount
	sa.AssetPositions = []*satypes.AssetPosition{
		testutil.CreateSingleAssetPosition(0, quantums),
	}
	tApp.App.SubaccountsKeeper.SetSubaccount(ctx, sa)
}

// setFeeParams appoints `operator` and sets the three fee parameters.
func setFeeParams(
	t *testing.T,
	k keeper.Keeper,
	ctx sdk.Context,
	operator string,
	operatorFeePpm uint32,
	profitSharePpm uint32,
	minOperatorSharePpm uint32,
) {
	t.Helper()
	require.NoError(t, k.SetOperatorParams(ctx, vaulttypes.OperatorParams{
		Operator:            operator,
		OperatorFeePpm:      operatorFeePpm,
		ProfitSharePpm:      profitSharePpm,
		MinOperatorSharePpm: minOperatorSharePpm,
	}))
}

// TestAccrueFees_Inert verifies that with both fee rates at zero, no shares are
// ever minted — but the high-water mark still tracks NAV, so that turning fees on
// later does not retroactively charge for gains earned while they were off.
func TestAccrueFees_Inert(t *testing.T) {
	owner := constants.AliceAccAddress.String()
	_, ctx, k := setupMegavault(t, big.NewInt(1_000), big.NewInt(1_000), owner)

	require.NoError(t, k.AccrueFees(ctx))

	require.Equal(t, big.NewInt(1_000), k.GetTotalShares(ctx).NumShares.BigInt())
	// NAV per share is 1.0, scaled by 1e18.
	require.Equal(
		t,
		vaulttypes.NavPerShareScale,
		k.GetFeeState(ctx).HighWaterMarkNavPerShare.BigInt(),
	)
}

// TestAccrueFees_ProfitShare covers the core profit-share path and the
// no-double-charge invariant.
func TestAccrueFees_ProfitShare(t *testing.T) {
	owner := constants.AliceAccAddress.String()
	operator := constants.BobAccAddress.String()

	// Start at NAV per share of 1.0 with a 20% profit share.
	tApp, ctx, k := setupMegavault(t, big.NewInt(1_000), big.NewInt(1_000), owner)
	setFeeParams(t, k, ctx, operator, 0, 200_000, 0)

	// Establish the high-water mark at the starting NAV.
	require.NoError(t, k.AccrueFees(ctx))
	hwmBefore := k.GetFeeState(ctx).HighWaterMarkNavPerShare.BigInt()
	require.Equal(t, vaulttypes.NavPerShareScale, hwmBefore)
	require.Equal(t, big.NewInt(1_000), k.GetTotalShares(ctx).NumShares.BigInt())

	// Megavault doubles: equity 1_000 -> 2_000, NAV per share 1.0 -> 2.0.
	fundMegavault(t, tApp, ctx, big.NewInt(1_000))

	require.NoError(t, k.AccrueFees(ctx))

	// Profit is 1_000 quote quantums; 20% of it, 200, accrues to the operator.
	// sharesToMint = 1_000 * 200 / (2_000 - 200) = 111.
	operatorShares, exists := k.GetOwnerShares(ctx, operator)
	require.True(t, exists)
	require.Equal(t, big.NewInt(111), operatorShares.NumShares.BigInt())
	require.Equal(t, big.NewInt(1_111), k.GetTotalShares(ctx).NumShares.BigInt())

	// The operator's shares are worth approximately the fee that was charged.
	operatorValue := new(big.Int).Mul(operatorShares.NumShares.BigInt(), big.NewInt(2_000))
	operatorValue.Quo(operatorValue, big.NewInt(1_111))
	require.InDelta(t, 200, operatorValue.Int64(), 2)

	// Invariant: no double-charge. Accruing again with no NAV change mints nothing.
	sharesAfterFirst := k.GetTotalShares(ctx).NumShares.BigInt()
	require.NoError(t, k.AccrueFees(ctx))
	require.Equal(t, sharesAfterFirst, k.GetTotalShares(ctx).NumShares.BigInt())

	// Invariant: the high-water mark is monotonic and rose.
	hwmAfter := k.GetFeeState(ctx).HighWaterMarkNavPerShare.BigInt()
	require.Equal(t, 1, hwmAfter.Cmp(hwmBefore))
}

// TestAccrueFees_NoProfitShareBelowHighWaterMark checks that a loss followed by a
// partial recovery is not charged, which is the whole point of a high-water mark.
func TestAccrueFees_NoProfitShareBelowHighWaterMark(t *testing.T) {
	owner := constants.AliceAccAddress.String()
	operator := constants.BobAccAddress.String()

	tApp, ctx, k := setupMegavault(t, big.NewInt(2_000), big.NewInt(1_000), owner)
	setFeeParams(t, k, ctx, operator, 0, 200_000, 0)

	// High-water mark is set at NAV per share 2.0.
	require.NoError(t, k.AccrueFees(ctx))
	hwm := k.GetFeeState(ctx).HighWaterMarkNavPerShare.BigInt()
	sharesBefore := k.GetTotalShares(ctx).NumShares.BigInt()

	// Megavault loses half, then recovers part of the way — still under the mark.
	fundMegavault(t, tApp, ctx, big.NewInt(-1_000))
	require.NoError(t, k.AccrueFees(ctx))
	require.Equal(t, sharesBefore, k.GetTotalShares(ctx).NumShares.BigInt())

	fundMegavault(t, tApp, ctx, big.NewInt(500))
	require.NoError(t, k.AccrueFees(ctx))

	// No shares minted, and the mark did not fall.
	require.Equal(t, sharesBefore, k.GetTotalShares(ctx).NumShares.BigInt())
	require.Equal(t, hwm, k.GetFeeState(ctx).HighWaterMarkNavPerShare.BigInt())
	_, exists := k.GetOwnerShares(ctx, operator)
	require.False(t, exists)
}

// TestAccrueFees_OperatorFeeIsTimeBased checks that the operator fee scales with
// elapsed time and is charged regardless of performance.
func TestAccrueFees_OperatorFeeIsTimeBased(t *testing.T) {
	owner := constants.AliceAccAddress.String()
	operator := constants.BobAccAddress.String()

	// 10% annualized operator fee, no profit share.
	_, ctx, k := setupMegavault(t, big.NewInt(1_000_000_000), big.NewInt(1_000_000_000), owner)
	setFeeParams(t, k, ctx, operator, 100_000, 0, 0)

	// The accrual clock is seeded at genesis, so an accrual at the genesis block
	// time has no elapsed time and charges nothing.
	require.NoError(t, k.AccrueFees(ctx))
	require.Equal(t, big.NewInt(1_000_000_000), k.GetTotalShares(ctx).NumShares.BigInt())
	require.Equal(t, ctx.BlockTime().Unix(), k.GetFeeState(ctx).LastAccrualTime)

	// Half a year later, roughly 5% of equity should have accrued.
	halfYear := time.Duration(vaulttypes.SecondsPerYear/2) * time.Second
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(halfYear))
	require.NoError(t, k.AccrueFees(ctx))

	operatorShares, exists := k.GetOwnerShares(ctx, operator)
	require.True(t, exists)

	totalShares := k.GetTotalShares(ctx).NumShares.BigInt()
	operatorValue := new(big.Int).Mul(operatorShares.NumShares.BigInt(), big.NewInt(1_000_000_000))
	operatorValue.Quo(operatorValue, totalShares)
	// 10% per year for half a year on 1e9 equity is 5e7 quote quantums.
	require.InEpsilon(t, 50_000_000, operatorValue.Int64(), 0.001)
}

// TestAccrueFees_RoundToZeroDoesNotAdvanceClock is the subtle case called out in
// ADR-001 §4: if the accrual timestamp advanced on a fee that rounded to zero, a
// low rate would accrue nothing forever.
func TestAccrueFees_RoundToZeroDoesNotAdvanceClock(t *testing.T) {
	owner := constants.AliceAccAddress.String()
	operator := constants.BobAccAddress.String()

	// A tiny equity and a small rate, so one second of accrual rounds to zero.
	_, ctx, k := setupMegavault(t, big.NewInt(1_000), big.NewInt(1_000), owner)
	setFeeParams(t, k, ctx, operator, 1_000, 0, 0) // 0.1%/yr

	require.NoError(t, k.AccrueFees(ctx))
	startTime := k.GetFeeState(ctx).LastAccrualTime
	require.Equal(t, ctx.BlockTime().Unix(), startTime)

	// One second later the fee rounds to zero...
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(time.Second))
	require.NoError(t, k.AccrueFees(ctx))
	require.Equal(t, big.NewInt(1_000), k.GetTotalShares(ctx).NumShares.BigInt())
	// ...and the clock must not have advanced, or the fee would be lost forever.
	require.Equal(t, startTime, k.GetFeeState(ctx).LastAccrualTime)

	// Given enough time, the accumulated fee becomes representable and is
	// charged. `MaxAccrualElapsedSeconds` bounds how long that can take.
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(
		time.Duration(vaulttypes.MaxAccrualElapsedSeconds) * time.Second,
	))
	require.NoError(t, k.AccrueFees(ctx))
	require.Equal(t, 1, k.GetTotalShares(ctx).NumShares.BigInt().Cmp(big.NewInt(1_000)))
	require.Equal(t, ctx.BlockTime().Unix(), k.GetFeeState(ctx).LastAccrualTime)
}

// TestAccrueFees_ElapsedTimeIsClamped verifies the ADR-001 §4 backstop: a stale
// or uninitialized accrual clock cannot produce an unbounded one-off charge.
func TestAccrueFees_ElapsedTimeIsClamped(t *testing.T) {
	owner := constants.AliceAccAddress.String()
	operator := constants.BobAccAddress.String()

	// 10% annualized operator fee.
	_, ctx, k := setupMegavault(t, big.NewInt(1_000_000_000), big.NewInt(1_000_000_000), owner)
	setFeeParams(t, k, ctx, operator, 100_000, 0, 0)

	// Jump forward a decade in one step.
	ctx = ctx.WithBlockTime(ctx.BlockTime().Add(
		time.Duration(vaulttypes.SecondsPerYear*10) * time.Second,
	))
	require.NoError(t, k.AccrueFees(ctx))

	operatorShares, exists := k.GetOwnerShares(ctx, operator)
	require.True(t, exists)

	totalShares := k.GetTotalShares(ctx).NumShares.BigInt()
	operatorValue := new(big.Int).Mul(operatorShares.NumShares.BigInt(), big.NewInt(1_000_000_000))
	operatorValue.Quo(operatorValue, totalShares)
	// One year of fees at 10% on 1e9, not ten.
	require.InEpsilon(t, 100_000_000, operatorValue.Int64(), 0.001)
}

// TestAccrueFees_SkipsWhenUndefined covers the "no shares" and "non-positive
// equity" cases, which must be skips rather than divisions by zero.
func TestAccrueFees_SkipsWhenUndefined(t *testing.T) {
	operator := constants.BobAccAddress.String()

	t.Run("no shares", func(t *testing.T) {
		tApp := testapp.NewTestAppBuilder(t).Build()
		ctx := tApp.InitChain()
		k := tApp.App.VaultKeeper
		setFeeParams(t, k, ctx, operator, 100_000, 200_000, 0)

		require.NoError(t, k.AccrueFees(ctx))
		require.Equal(t, big.NewInt(0), k.GetTotalShares(ctx).NumShares.BigInt())
	})

	t.Run("zero equity", func(t *testing.T) {
		owner := constants.AliceAccAddress.String()
		_, ctx, k := setupMegavault(t, big.NewInt(0), big.NewInt(1_000), owner)
		setFeeParams(t, k, ctx, operator, 100_000, 200_000, 0)

		require.NoError(t, k.AccrueFees(ctx))
		require.Equal(t, big.NewInt(1_000), k.GetTotalShares(ctx).NumShares.BigInt())
	})
}

// TestValidateOperatorShareFloor covers ADR-001 D6.
func TestValidateOperatorShareFloor(t *testing.T) {
	operator := constants.BobAccAddress.String()
	other := constants.AliceAccAddress.String()

	tests := map[string]struct {
		withdrawer          string
		sharesToWithdraw    *big.Int
		minOperatorSharePpm uint32
		expectedErr         error
	}{
		"No floor configured": {
			withdrawer:       operator,
			sharesToWithdraw: big.NewInt(400),
		},
		"Non-operator is unaffected by the floor": {
			withdrawer:          other,
			sharesToWithdraw:    big.NewInt(400),
			minOperatorSharePpm: 500_000,
		},
		"Operator withdrawal that keeps the floor": {
			withdrawer: operator,
			// 400 of 1_000 operator shares out of 2_000 total leaves
			// 600/1_600 = 37.5%, above a 30% floor.
			sharesToWithdraw:    big.NewInt(400),
			minOperatorSharePpm: 300_000,
		},
		"Operator withdrawal that breaches the floor": {
			withdrawer: operator,
			// Leaves 600/1_600 = 37.5%, below a 40% floor.
			sharesToWithdraw:    big.NewInt(400),
			minOperatorSharePpm: 400_000,
			expectedErr:         vaulttypes.ErrOperatorShareBelowMinimum,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Total 2_000 shares: 1_000 to the operator, 1_000 to another owner.
			tApp := testapp.NewTestAppBuilder(t).WithGenesisDocFn(func() (genesis types.GenesisDoc) {
				genesis = testapp.DefaultGenesis()
				testapp.UpdateGenesisDocWithAppStateForModule(
					&genesis,
					func(genesisState *vaulttypes.GenesisState) {
						genesisState.TotalShares = vaulttypes.BigIntToNumShares(big.NewInt(2_000))
						genesisState.OwnerShares = []vaulttypes.OwnerShare{
							{
								Owner:  operator,
								Shares: vaulttypes.BigIntToNumShares(big.NewInt(1_000)),
							},
							{
								Owner:  other,
								Shares: vaulttypes.BigIntToNumShares(big.NewInt(1_000)),
							},
						}
					},
				)
				return genesis
			}).Build()
			ctx := tApp.InitChain()
			k := tApp.App.VaultKeeper
			setFeeParams(t, k, ctx, operator, 0, 0, tc.minOperatorSharePpm)

			err := k.ValidateOperatorShareFloor(ctx, tc.withdrawer, tc.sharesToWithdraw)
			if tc.expectedErr != nil {
				require.ErrorIs(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestFeeStateValidation checks that state written before the fee engine existed
// (a nil high-water mark) reads back as valid rather than panicking.
func TestFeeStateValidation(t *testing.T) {
	tApp := testapp.NewTestAppBuilder(t).Build()
	ctx := tApp.InitChain()
	k := tApp.App.VaultKeeper

	require.Equal(t, vaulttypes.DefaultFeeState(), k.GetFeeState(ctx))
	require.NoError(t, vaulttypes.FeeState{}.Validate())

	require.ErrorIs(
		t,
		k.SetFeeState(ctx, vaulttypes.FeeState{
			HighWaterMarkNavPerShare: dtypes.NewInt(-1),
		}),
		vaulttypes.ErrNegativeHighWaterMark,
	)
}

// TestWithdrawCrystallizesFees checks ADR-001 §6: a withdrawal forces an accrual
// first, so the exiting depositor pays their pro-rata share of fees earned up to
// that moment rather than escaping them.
func TestWithdrawCrystallizesFees(t *testing.T) {
	owner := constants.Alice_Num0
	operator := constants.BobAccAddress.String()

	// Megavault holds 2_000 quote quantums against 1_000 shares, all Alice's,
	// so NAV per share is 2.0.
	tApp, ctx, k := setupMegavault(t, big.NewInt(2_000), big.NewInt(1_000), owner.Owner)
	setFeeParams(t, k, ctx, operator, 0, 200_000, 0)

	// The high-water mark already tracks the starting NAV — it is maintained
	// even while fees are off, which is what stops enabling fees from charging
	// retroactively. Gains have to come after that mark to be chargeable.
	require.NoError(t, k.AccrueFees(ctx))
	_, exists := k.GetOwnerShares(ctx, operator)
	require.False(t, exists)

	// Megavault gains 1_000: NAV per share 2.0 -> 3.0.
	fundMegavault(t, tApp, ctx, big.NewInt(1_000))

	totalSharesBefore := k.GetTotalShares(ctx).NumShares.BigInt()
	_, err := k.WithdrawFromMegavault(ctx, owner, big.NewInt(100), big.NewInt(0))
	require.NoError(t, err)

	// The operator now holds shares that did not exist before the withdrawal,
	// which is the accrual having been forced by the withdrawal itself.
	operatorShares, exists := k.GetOwnerShares(ctx, operator)
	require.True(t, exists)
	require.Equal(t, 1, operatorShares.NumShares.BigInt().Sign())

	// Total shares moved by the mint minus the burn, not by the burn alone.
	totalSharesAfter := k.GetTotalShares(ctx).NumShares.BigInt()
	burnedOnly := new(big.Int).Sub(totalSharesBefore, big.NewInt(100))
	require.Equal(t, 1, totalSharesAfter.Cmp(burnedOnly))
}

// TestWithdrawRespectsOperatorShareFloor checks that the floor is enforced on the
// real withdrawal path, not just in the helper.
func TestWithdrawRespectsOperatorShareFloor(t *testing.T) {
	operator := constants.Alice_Num0

	// The operator owns every share, and must retain at least 50%.
	_, ctx, k := setupMegavault(t, big.NewInt(2_000), big.NewInt(1_000), operator.Owner)
	setFeeParams(t, k, ctx, operator.Owner, 0, 0, 500_000)

	// Withdrawing 600 of 1_000 leaves 400/400 = 100%, which is above the floor
	// because the burn removes the shares from both sides of the ratio.
	// Give another owner shares so the ratio can actually fall.
	require.NoError(t, k.SetOwnerShares(
		ctx,
		constants.BobAccAddress.String(),
		vaulttypes.BigIntToNumShares(big.NewInt(1_000)),
	))
	require.NoError(t, k.SetTotalShares(ctx, vaulttypes.BigIntToNumShares(big.NewInt(2_000))))

	// Operator holds 1_000 of 2_000 = 50%. Withdrawing 400 leaves 600/1_600 =
	// 37.5%, below the 50% floor.
	_, err := k.WithdrawFromMegavault(ctx, operator, big.NewInt(400), big.NewInt(0))
	require.ErrorIs(t, err, vaulttypes.ErrOperatorShareBelowMinimum)

	// A non-operator withdrawing the same amount is unaffected.
	bob := satypes.SubaccountId{Owner: constants.BobAccAddress.String(), Number: 0}
	_, err = k.WithdrawFromMegavault(ctx, bob, big.NewInt(400), big.NewInt(0))
	require.NoError(t, err)
}
