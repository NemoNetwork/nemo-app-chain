package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeDepositToMegavault = "deposit_to_megavault"
	AttributeKeyDepositor       = "depositor"
	AttributeKeyQuoteQuantums   = "quote_quantums"
	AttributeKeyMintedShares    = "minted_shares"

	EventTypeWithdrawFromMegavault    = "withdraw_from_megavault"
	AttributeKeyWithdrawer            = "withdrawer"
	AttributeKeySharesToWithdraw      = "shares_to_withdraw"
	AttributeKeyTotalShares           = "total_shares"
	AttributeKeyMegavaultEquity       = "megavault_equity"
	AttributeKeyRedeemedQuoteQuantums = "redeemed_quote_quantums"

	EventTypeSweepToMegavault      = "sweep_to_megavault"
	AttributeKeySweptQuoteQuantums = "swept_quote_quantums"

	// Fork-local; see x/vault/spec/adr-001-megavault-fees.md.
	EventTypeAccrueMegavaultFees = "accrue_megavault_fees"
	AttributeKeyOperator         = "operator"
	AttributeKeyOperatorFee      = "operator_fee_quote_quantums"
	AttributeKeyProfitShare      = "profit_share_quote_quantums"
)

// NewDepositToMegavaultEvent constructs a new deposit_to_megavault sdk.Event.
func NewDepositToMegavaultEvent(
	depositorAddress string,
	quoteQuantums uint64,
	mintedShares uint64,
) sdk.Event {
	return sdk.NewEvent(
		EventTypeDepositToMegavault,
		sdk.NewAttribute(AttributeKeyDepositor, depositorAddress),
		sdk.NewAttribute(AttributeKeyQuoteQuantums, fmt.Sprintf("%d", quoteQuantums)),
		sdk.NewAttribute(AttributeKeyMintedShares, fmt.Sprintf("%d", mintedShares)),
	)
}

// NewWithdrawFromMegavaultEvent constructs a new withdraw_from_megavault sdk.Event.
func NewWithdrawFromMegavaultEvent(
	withdrawerAddress string,
	sharesToWithdraw uint64,
	totalShares uint64,
	megavaultEquity uint64,
	redeemedQuoteQuantums uint64,
) sdk.Event {
	return sdk.NewEvent(
		EventTypeWithdrawFromMegavault,
		sdk.NewAttribute(AttributeKeyWithdrawer, withdrawerAddress),
		sdk.NewAttribute(AttributeKeySharesToWithdraw, fmt.Sprintf("%d", sharesToWithdraw)),
		sdk.NewAttribute(AttributeKeyTotalShares, fmt.Sprintf("%d", totalShares)),
		sdk.NewAttribute(AttributeKeyMegavaultEquity, fmt.Sprintf("%d", megavaultEquity)),
		sdk.NewAttribute(AttributeKeyRedeemedQuoteQuantums, fmt.Sprintf("%d", redeemedQuoteQuantums)),
	)
}

// NewAccrueMegavaultFeesEvent constructs a new accrue_megavault_fees sdk.Event.
//
// Note: fork-local; upstream dydxprotocol has no fee mechanism.
func NewAccrueMegavaultFeesEvent(
	operatorAddress string,
	operatorFeeQuoteQuantums uint64,
	profitShareQuoteQuantums uint64,
	mintedShares uint64,
) sdk.Event {
	return sdk.NewEvent(
		EventTypeAccrueMegavaultFees,
		sdk.NewAttribute(AttributeKeyOperator, operatorAddress),
		sdk.NewAttribute(AttributeKeyOperatorFee, fmt.Sprintf("%d", operatorFeeQuoteQuantums)),
		sdk.NewAttribute(AttributeKeyProfitShare, fmt.Sprintf("%d", profitShareQuoteQuantums)),
		sdk.NewAttribute(AttributeKeyMintedShares, fmt.Sprintf("%d", mintedShares)),
	)
}

func NewSweepToMegavaultEvent(
	quoteQuantums uint64,
) sdk.Event {
	return sdk.NewEvent(
		EventTypeSweepToMegavault,
		sdk.NewAttribute(AttributeKeySweptQuoteQuantums, fmt.Sprintf("%d", quoteQuantums)),
	)
}
