package types

import "math/big"

// Constants for the megavault fee engine.
// See x/vault/spec/adr-001-megavault-fees.md.
const (
	// SecondsPerYear is the denominator used to convert an annualized
	// `operator_fee_ppm` into a per-accrual amount. 365 days; leap seconds and
	// leap days are immaterial at ppm precision.
	SecondsPerYear = int64(365 * 24 * 60 * 60)

	// FeeAccrualIntervalSeconds is the minimum time between fee accruals driven
	// by the EndBlocker. Withdrawals and deposits force an accrual regardless of
	// this interval.
	FeeAccrualIntervalSeconds = int64(3600)

	// MaxFeeFractionPerAccrualPpm caps how much of megavault equity a single
	// accrual may take. It exists so that `equity - fee` stays comfortably
	// positive in the share-minting formula, and as a backstop against a
	// mis-set parameter or an absurdly long gap between accruals.
	MaxFeeFractionPerAccrualPpm = uint32(500_000) // 50%

	// MaxAccrualElapsedSeconds caps the elapsed time any single accrual may
	// charge for.
	//
	// This bounds the damage if `last_accrual_time` is ever left uninitialized
	// (a chain that starts at a real timestamp with a stored zero would
	// otherwise be charged decades of fees on its first accrual), and stops a
	// long chain halt from producing a huge charge on resume. It is set to a
	// year so that a fee too small to be representable over one interval still
	// accumulates enough to be charged within a reasonable period.
	MaxAccrualElapsedSeconds = SecondsPerYear
)

// NavPerShareScale is the fixed-point scale used to store NAV per share.
// NAV per share is a ratio of two big integers; storing it scaled by 1e18
// preserves far more precision than quote-quantum arithmetic requires.
var NavPerShareScale = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
