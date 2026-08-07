package keeper

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/nemo-network/v4-chain/protocol/dtypes"
	"github.com/nemo-network/v4-chain/protocol/lib"
	"github.com/nemo-network/v4-chain/protocol/lib/log"
	"github.com/nemo-network/v4-chain/protocol/x/vault/types"
)

// GetFeeState returns the megavault `FeeState` in state.
func (k Keeper) GetFeeState(
	ctx sdk.Context,
) (
	feeState types.FeeState,
) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get([]byte(types.MegavaultFeeStateKey))
	k.cdc.MustUnmarshal(b, &feeState)
	return feeState
}

// SetFeeState sets the megavault `FeeState` in state.
func (k Keeper) SetFeeState(
	ctx sdk.Context,
	feeState types.FeeState,
) error {
	if err := feeState.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&feeState)
	store.Set([]byte(types.MegavaultFeeStateKey), b)

	return nil
}

// GetNavPerShare returns megavault NAV per share scaled by
// `types.NavPerShareScale`, along with the equity and total shares it was
// derived from.
//
// `exists` is false when NAV per share is undefined, which is the case when
// there are no shares or when equity is non-positive. Callers must treat that
// as "skip", never as "NAV is zero".
func (k Keeper) GetNavPerShare(
	ctx sdk.Context,
) (
	navPerShare *big.Int,
	equity *big.Int,
	totalShares *big.Int,
	exists bool,
	err error,
) {
	totalShares = k.GetTotalShares(ctx).NumShares.BigInt()
	if totalShares.Sign() <= 0 {
		return nil, nil, totalShares, false, nil
	}

	equity, err = k.GetMegavaultEquity(ctx)
	if err != nil {
		return nil, nil, totalShares, false, err
	}
	if equity.Sign() <= 0 {
		return nil, equity, totalShares, false, nil
	}

	navPerShare = new(big.Int).Mul(equity, types.NavPerShareScale)
	navPerShare.Quo(navPerShare, totalShares)

	return navPerShare, equity, totalShares, true, nil
}

// AccrueFees charges megavault's operator fee and profit share, if any are due.
//
// It is called from the EndBlocker on an interval, and forced at the top of
// every deposit and withdrawal so that the depositor entering or exiting sees a
// NAV that carries no unaccrued fee liability.
//
// Fees are paid by minting shares to the operator, which dilutes every existing
// share equally — that is what makes "force an accrual before a withdrawal"
// equivalent to "the exiting depositor pays their pro-rata share".
//
// See x/vault/spec/adr-001-megavault-fees.md.
func (k Keeper) AccrueFees(ctx sdk.Context) error {
	operatorParams := k.GetOperatorParams(ctx)
	if operatorParams.OperatorFeePpm == 0 && operatorParams.ProfitSharePpm == 0 {
		// The fee engine is inert. Still maintain the high-water mark, so that
		// enabling fees later does not retroactively charge for the gains made
		// while they were off.
		return k.updateHighWaterMark(ctx)
	}

	_, equity, totalShares, exists, err := k.GetNavPerShare(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	feeState := k.GetFeeState(ctx)
	blockTime := ctx.BlockTime().Unix()

	// Step 1: the operator fee, charged on NAV pro-rata over elapsed time.
	operatorFee := big.NewInt(0)
	if operatorParams.OperatorFeePpm > 0 {
		elapsedSeconds := blockTime - feeState.LastAccrualTime
		// Clamp so that a stale or uninitialized `LastAccrualTime` cannot
		// produce an enormous one-off charge. See ADR-001 §4.
		if elapsedSeconds > types.MaxAccrualElapsedSeconds {
			elapsedSeconds = types.MaxAccrualElapsedSeconds
		}
		if elapsedSeconds > 0 {
			// equity * operator_fee_ppm / 1e6 * elapsed / secondsPerYear
			operatorFee = new(big.Int).Mul(equity, big.NewInt(int64(operatorParams.OperatorFeePpm)))
			operatorFee.Mul(operatorFee, big.NewInt(elapsedSeconds))
			operatorFee.Quo(operatorFee, big.NewInt(int64(lib.OneMillion)))
			operatorFee.Quo(operatorFee, big.NewInt(types.SecondsPerYear))
		}
	}

	// Step 2: the profit share, charged on the excess over the high-water mark.
	//
	// Computed against NAV per share after the operator fee has diluted, so the
	// profit share is levied on gains net of the operator fee.
	equityAfterOperatorFee := new(big.Int).Sub(equity, operatorFee)
	profitShare := big.NewInt(0)
	if operatorParams.ProfitSharePpm > 0 && equityAfterOperatorFee.Sign() > 0 {
		hwm := feeState.HighWaterMarkNavPerShare.BigInt()
		if hwm == nil {
			hwm = big.NewInt(0)
		}
		// NAV per share does not change when the operator fee is paid in shares
		// in aggregate terms, but the profit base does: value already promised
		// to the operator is not depositor profit. Compute the post-operator-fee
		// NAV per share directly from the reduced equity.
		navAfterOperatorFee := new(big.Int).Mul(equityAfterOperatorFee, types.NavPerShareScale)
		navAfterOperatorFee.Quo(navAfterOperatorFee, totalShares)

		if navAfterOperatorFee.Cmp(hwm) > 0 {
			// profit = (nav - hwm) * totalShares / scale
			profit := new(big.Int).Sub(navAfterOperatorFee, hwm)
			profit.Mul(profit, totalShares)
			profit.Quo(profit, types.NavPerShareScale)

			profitShare = new(big.Int).Mul(profit, big.NewInt(int64(operatorParams.ProfitSharePpm)))
			profitShare.Quo(profitShare, big.NewInt(int64(lib.OneMillion)))
		}
	}

	totalFee := new(big.Int).Add(operatorFee, profitShare)
	if totalFee.Sign() <= 0 {
		// Nothing was owed, or the amount owed rounded to zero.
		//
		// `LastAccrualTime` is deliberately left untouched when the fee rounds
		// to zero: advancing it would mean a rate too small to be representable
		// over one interval never accrues at all, however much time passes.
		// `MaxAccrualElapsedSeconds` bounds how long that can go on.
		return k.updateHighWaterMark(ctx)
	}

	// Clamp so that `equity - totalFee` stays comfortably positive.
	maxFee := new(big.Int).Mul(equity, big.NewInt(int64(types.MaxFeeFractionPerAccrualPpm)))
	maxFee.Quo(maxFee, big.NewInt(int64(lib.OneMillion)))
	if totalFee.Cmp(maxFee) > 0 {
		totalFee = maxFee
	}

	// Mint shares worth `totalFee` quote quantums to the operator:
	//   sharesToMint = totalShares * totalFee / (equity - totalFee)
	// so that the minted shares are worth exactly `totalFee` afterwards.
	denominator := new(big.Int).Sub(equity, totalFee)
	if denominator.Sign() <= 0 {
		return nil
	}
	sharesToMint := new(big.Int).Mul(totalShares, totalFee)
	sharesToMint.Quo(sharesToMint, denominator)
	if sharesToMint.Sign() <= 0 {
		// The fee is real but too small to be expressed in whole shares. Leave
		// `LastAccrualTime` unadvanced so it accumulates.
		return k.updateHighWaterMark(ctx)
	}

	if err := k.mintSharesTo(ctx, operatorParams.Operator, sharesToMint, totalShares); err != nil {
		return err
	}

	feeState.LastAccrualTime = blockTime
	if err := k.SetFeeState(ctx, feeState); err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(
		types.NewAccrueMegavaultFeesEvent(
			operatorParams.Operator,
			operatorFee.Uint64(),
			profitShare.Uint64(),
			sharesToMint.Uint64(),
		),
	)

	// The high-water mark is set from NAV per share *after* dilution, so the
	// same gains are never charged twice.
	return k.updateHighWaterMark(ctx)
}

// updateHighWaterMark raises the stored high-water mark to the current NAV per
// share if the latter is higher. The mark never decreases.
func (k Keeper) updateHighWaterMark(ctx sdk.Context) error {
	navPerShare, _, _, exists, err := k.GetNavPerShare(ctx)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}

	feeState := k.GetFeeState(ctx)
	hwm := feeState.HighWaterMarkNavPerShare.BigInt()
	if hwm != nil && hwm.Cmp(navPerShare) >= 0 {
		return nil
	}

	feeState.HighWaterMarkNavPerShare = dtypes.NewIntFromBigInt(navPerShare)
	return k.SetFeeState(ctx, feeState)
}

// mintSharesTo credits `sharesToMint` to `owner` and increases total shares.
//
// Unlike `MintShares`, this does not derive the share count from a deposit — the
// caller has already computed it — and it does not move any funds.
func (k Keeper) mintSharesTo(
	ctx sdk.Context,
	owner string,
	sharesToMint *big.Int,
	existingTotalShares *big.Int,
) error {
	newTotalShares := new(big.Int).Add(existingTotalShares, sharesToMint)
	if err := k.SetTotalShares(ctx, types.BigIntToNumShares(newTotalShares)); err != nil {
		return err
	}

	ownerShares, exists := k.GetOwnerShares(ctx, owner)
	newOwnerShares := new(big.Int).Set(sharesToMint)
	if exists {
		newOwnerShares.Add(newOwnerShares, ownerShares.NumShares.BigInt())
	}

	return k.SetOwnerShares(ctx, owner, types.BigIntToNumShares(newOwnerShares))
}

// ValidateOperatorShareFloor returns an error if `sharesToWithdraw` by `owner`
// would take the operator's share of megavault below `min_operator_share_ppm`.
//
// The floor binds the operator's own actions only. A user deposit that dilutes
// the operator below the floor is allowed — rejecting it would punish depositors
// for the operator's failure to co-invest, and would let the operator halt
// deposits by withdrawing. See ADR-001 D6.
func (k Keeper) ValidateOperatorShareFloor(
	ctx sdk.Context,
	owner string,
	sharesToWithdraw *big.Int,
) error {
	operatorParams := k.GetOperatorParams(ctx)
	if operatorParams.MinOperatorSharePpm == 0 || owner != operatorParams.Operator {
		return nil
	}

	ownerShares, exists := k.GetOwnerShares(ctx, owner)
	if !exists {
		return nil
	}

	totalShares := k.GetTotalShares(ctx).NumShares.BigInt()
	operatorSharesAfter := new(big.Int).Sub(ownerShares.NumShares.BigInt(), sharesToWithdraw)
	totalSharesAfter := new(big.Int).Sub(totalShares, sharesToWithdraw)
	if totalSharesAfter.Sign() <= 0 {
		// The operator is withdrawing everything there is; there is no ratio to
		// enforce and the ordinary share checks cover it.
		return nil
	}

	// operatorSharesAfter / totalSharesAfter >= minPpm / 1e6
	lhs := new(big.Int).Mul(operatorSharesAfter, big.NewInt(int64(lib.OneMillion)))
	rhs := new(big.Int).Mul(totalSharesAfter, big.NewInt(int64(operatorParams.MinOperatorSharePpm)))
	if lhs.Cmp(rhs) < 0 {
		return types.ErrOperatorShareBelowMinimum
	}

	return nil
}

// MaybeAccrueFees runs a fee accrual if at least `FeeAccrualIntervalSeconds`
// have elapsed since the last one. It is the EndBlocker entry point.
//
// Accruing every block would be wasteful and, at realistic rates, would round to
// zero every time. Deposits and withdrawals call `AccrueFees` directly and are
// not subject to this interval.
func (k Keeper) MaybeAccrueFees(ctx sdk.Context) {
	feeState := k.GetFeeState(ctx)
	blockTime := ctx.BlockTime().Unix()

	if feeState.LastAccrualTime > 0 &&
		blockTime-feeState.LastAccrualTime < types.FeeAccrualIntervalSeconds {
		return
	}

	if err := k.AccrueFees(ctx); err != nil {
		log.ErrorLogWithError(ctx, "failed to accrue megavault fees", err)
	}
}
