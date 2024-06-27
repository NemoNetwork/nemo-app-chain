package lib

import (
	"math/big"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/lib/margin"
	"github.com/dydxprotocol/v4-chain/protocol/x/perpetuals/types"
	pricestypes "github.com/dydxprotocol/v4-chain/protocol/x/prices/types"
)

// GetSettlementPpm returns the net settlement amount ppm (in quote quantums) given
// the perpetual and position size (in base quantums).
func GetSettlementPpmWithPerpetual(
	perpetual types.Perpetual,
	quantums *big.Int,
	index *big.Int,
) (
	bigNetSettlementPpm *big.Int,
	newFundingIndex *big.Int,
) {
	fundingIndex := perpetual.FundingIndex.BigInt()

	// If no change in funding, return 0.
	if fundingIndex.Cmp(index) == 0 {
		return big.NewInt(0), fundingIndex
	}

	// The settlement is a signed value.
	// If the index delta is positive and the quantums is positive (long), then settlement is negative.
	// Thus, always negate the value of the multiplication of the index delta and the quantums.
	result := new(big.Int).Sub(fundingIndex, index)
	result = result.Mul(result, quantums)
	result = result.Neg(result)

	return result, fundingIndex
}

// GetNetCollateralAndMarginRequirements returns the net collateral, initial margin requirement,
// and maintenance margin requirement in quote quantums, given the position size in base quantums.
func GetNetCollateralAndMarginRequirements(
	perpetual types.Perpetual,
	marketPrice pricestypes.MarketPrice,
	liquidityTier types.LiquidityTier,
	quantums *big.Int,
) (
	risk margin.Risk,
) {
	nc := GetNetNotionalInQuoteQuantums(
		perpetual,
		marketPrice,
		quantums,
	)
	imr, mmr := GetMarginRequirementsInQuoteQuantums(
		perpetual,
		marketPrice,
		liquidityTier,
		quantums,
	)
	return margin.Risk{
		NC:  nc,
		IMR: imr,
		MMR: mmr,
	}
}

// GetNetNotionalInQuoteQuantums returns the net notional in quote quantums, which can be
// represented by the following equation:
//
// `quantums / 10^baseAtomicResolution * marketPrice * 10^marketExponent * 10^quoteAtomicResolution`.
// Note that longs are positive, and shorts are negative.
func GetNetNotionalInQuoteQuantums(
	perpetual types.Perpetual,
	marketPrice pricestypes.MarketPrice,
	bigQuantums *big.Int,
) (
	bigNetNotionalQuoteQuantums *big.Int,
) {
	bigQuoteQuantums := lib.BaseToQuoteQuantums(
		bigQuantums,
		perpetual.Params.AtomicResolution,
		marketPrice.Price,
		marketPrice.Exponent,
	)

	return bigQuoteQuantums
}

// GetMarginRequirementsInQuoteQuantums returns initial and maintenance margin requirements
// in quote quantums, given the position size in base quantums.
func GetMarginRequirementsInQuoteQuantums(
	perpetual types.Perpetual,
	marketPrice pricestypes.MarketPrice,
	liquidityTier types.LiquidityTier,
	bigQuantums *big.Int,
) (
	bigInitialMarginQuoteQuantums *big.Int,
	bigMaintenanceMarginQuoteQuantums *big.Int,
) {
	// Calculate the notional value of the position in quote quantums.
	// Always consider the magnitude of the position regardless of whether it is long/short.
	bigQuoteQuantums := new(big.Int).Abs(bigQuantums)
	bigQuoteQuantums = lib.BaseToQuoteQuantums(
		bigQuoteQuantums,
		perpetual.Params.AtomicResolution,
		marketPrice.Price,
		marketPrice.Exponent,
	)
	// Calculate the perpetual's open interest in quote quantums.
	openInterestQuoteQuantums := lib.BaseToQuoteQuantums(
		perpetual.OpenInterest.BigInt(), // OpenInterest is represented as base quantums.
		perpetual.Params.AtomicResolution,
		marketPrice.Price,
		marketPrice.Exponent,
	)

	// Initial margin requirement quote quantums = size in quote quantums * initial margin PPM.
	bigBaseInitialMarginQuoteQuantums := lib.BigMulPpm(
		bigQuoteQuantums,
		lib.BigU(liquidityTier.InitialMarginPpm),
		true, // Round up initial margin.
	)

	// Maintenance margin requirement quote quantums = IM in quote quantums * maintenance fraction PPM.
	bigMaintenanceMarginQuoteQuantums = lib.BigMulPpm(
		bigBaseInitialMarginQuoteQuantums,
		lib.BigU(liquidityTier.MaintenanceFractionPpm),
		true,
	)

	bigInitialMarginQuoteQuantums = liquidityTier.GetInitialMarginQuoteQuantums(
		bigQuoteQuantums,
		openInterestQuoteQuantums, // pass in current OI to get scaled IMR.
	)
	return bigInitialMarginQuoteQuantums, bigMaintenanceMarginQuoteQuantums
}
