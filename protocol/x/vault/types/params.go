package types

import (
	"math"
	"math/big"

	"github.com/nemo-network/v4-chain/protocol/dtypes"
	"github.com/nemo-network/v4-chain/protocol/lib"
)

// DefaultQuotingParams returns a default set of `x/vault` parameters.
func DefaultQuotingParams() QuotingParams {
	return QuotingParams{
		Layers:                           2,                            // 2 layers
		SpreadMinPpm:                     10_000,                       // 100 bps
		SpreadBufferPpm:                  1_500,                        // 15 bps
		SkewFactorPpm:                    2_000_000,                    // 2
		OrderSizePctPpm:                  100_000,                      // 10%
		OrderExpirationSeconds:           60,                           // 60 seconds
		ActivationThresholdQuoteQuantums: dtypes.NewInt(1_000_000_000), // 1_000 USDC
	}
}

// Validate validates `x/vault` parameters.
func (p QuotingParams) Validate() error {
	// Layers must be less than or equal to MaxUint8.
	if p.Layers > math.MaxUint8 {
		return ErrInvalidLayers
	}
	// Spread min ppm must be positive.
	if p.SpreadMinPpm == 0 {
		return ErrInvalidSpreadMinPpm
	}
	// Order size must be positive.
	if p.OrderSizePctPpm == 0 {
		return ErrInvalidOrderSizePctPpm
	}
	// Order expiration seconds must be positive.
	if p.OrderExpirationSeconds == 0 {
		return ErrInvalidOrderExpirationSeconds
	}
	// Activation threshold quote quantums must be non-negative.
	if p.ActivationThresholdQuoteQuantums.Sign() < 0 {
		return ErrInvalidActivationThresholdQuoteQuantums
	}
	// Skew factor times order_size_pct must be less than 2 to avoid skewing over the spread
	skewFactor := new(big.Int).SetUint64(uint64(p.SkewFactorPpm))
	orderSizePct := new(big.Int).SetUint64(uint64(p.OrderSizePctPpm))
	skewFactorOrderSizePctProduct := new(big.Int).Mul(skewFactor, orderSizePct)
	skewFactorOrderSizePctProductThreshold := big.NewInt(2_000_000 * 1_000_000)
	if skewFactorOrderSizePctProduct.Cmp(skewFactorOrderSizePctProductThreshold) >= 0 {
		return ErrInvalidSkewFactor
	}

	return nil
}

// DefaultOperatorParams returns a default set of `x/vault` operator parameters.
func DefaultOperatorParams() OperatorParams {
	return OperatorParams{
		Operator: lib.GovModuleAddress.String(),
	}
}

// Validate validates OperatorParams.
func (o OperatorParams) Validate() error {
	// Validate that operator is non-empty.
	if o.Operator == "" {
		return ErrEmptyOperator
	}

	// Fee parameters must be expressible as a fraction of one.
	// See x/vault/spec/adr-001-megavault-fees.md.
	for _, ppm := range []uint32{
		o.OperatorFeePpm,
		o.ProfitSharePpm,
		o.MinOperatorSharePpm,
	} {
		if ppm > lib.OneMillion {
			return ErrInvalidFeePpm
		}
	}

	return nil
}

// DefaultFeeState returns the default megavault fee state: no high-water mark
// and no accrual yet.
func DefaultFeeState() FeeState {
	return FeeState{
		HighWaterMarkNavPerShare: dtypes.NewInt(0),
	}
}

// Validate validates FeeState.
func (f FeeState) Validate() error {
	if f.HighWaterMarkNavPerShare.Sign() < 0 {
		return ErrNegativeHighWaterMark
	}

	return nil
}

// DefaultMegavaultParams returns a default set of `x/vault` megavault
// parameters. The zero deposit cap means deposits are uncapped.
func DefaultMegavaultParams() MegavaultParams {
	return MegavaultParams{
		DepositCapQuoteQuantums: dtypes.NewInt(0),
	}
}

// Validate validates MegavaultParams.
func (m MegavaultParams) Validate() error {
	// Deposit cap must be non-negative. Zero means no cap.
	if m.DepositCapQuoteQuantums.Sign() < 0 {
		return ErrNegativeDepositCap
	}

	return nil
}

// Validate validates individual vault parameters.
func (v VaultParams) Validate() error {
	// Validate status.
	if v.Status == VaultStatus_VAULT_STATUS_UNSPECIFIED {
		return ErrUnspecifiedVaultStatus
	}
	// Validate quoting params.
	if v.QuotingParams != nil {
		if err := v.QuotingParams.Validate(); err != nil {
			return err
		}
	}

	return nil
}
