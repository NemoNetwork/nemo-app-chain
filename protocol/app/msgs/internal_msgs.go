package msgs

import (
	upgrade "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensus "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisis "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distribution "github.com/cosmos/cosmos-sdk/x/distribution/types"
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	slashing "github.com/cosmos/cosmos-sdk/x/slashing/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v8/modules/core/02-client/types" //nolint:staticcheck
	ibcconn "github.com/cosmos/ibc-go/v8/modules/core/03-connection/types"
	"github.com/nemo-network/v4-chain/protocol/lib"
	affiliates "github.com/nemo-network/v4-chain/protocol/x/affiliates/types"
	blocktime "github.com/nemo-network/v4-chain/protocol/x/blocktime/types"
	bridge "github.com/nemo-network/v4-chain/protocol/x/bridge/types"
	clob "github.com/nemo-network/v4-chain/protocol/x/clob/types"
	delaymsg "github.com/nemo-network/v4-chain/protocol/x/delaymsg/types"
	feetiers "github.com/nemo-network/v4-chain/protocol/x/feetiers/types"
	govplus "github.com/nemo-network/v4-chain/protocol/x/govplus/types"
	listing "github.com/nemo-network/v4-chain/protocol/x/listing/types"
	perpetuals "github.com/nemo-network/v4-chain/protocol/x/perpetuals/types"
	prices "github.com/nemo-network/v4-chain/protocol/x/prices/types"
	ratelimit "github.com/nemo-network/v4-chain/protocol/x/ratelimit/types"
	revshare "github.com/nemo-network/v4-chain/protocol/x/revshare/types"
	rewards "github.com/nemo-network/v4-chain/protocol/x/rewards/types"
	sending "github.com/nemo-network/v4-chain/protocol/x/sending/types"
	stats "github.com/nemo-network/v4-chain/protocol/x/stats/types"
	vault "github.com/nemo-network/v4-chain/protocol/x/vault/types"
	vest "github.com/nemo-network/v4-chain/protocol/x/vest/types"
)

var (
	// InternalMsgSamplesAll are msgs that are used only used internally.
	InternalMsgSamplesAll = lib.MergeAllMapsMustHaveDistinctKeys(InternalMsgSamplesGovAuth)

	// InternalMsgSamplesGovAuth are msgs that are used only used internally.
	// GovAuth means that these messages must originate from the gov module and
	// signed by gov module account.
	// InternalMsgSamplesAll are msgs that are used only used internally.
	InternalMsgSamplesGovAuth = lib.MergeAllMapsMustHaveDistinctKeys(
		InternalMsgSamplesDefault,
		InternalMsgSamplesDydxCustom,
	)

	// CosmosSDK default modules
	InternalMsgSamplesDefault = map[string]sdk.Msg{
		// auth
		"/cosmos.auth.v1beta1.MsgUpdateParams": &auth.MsgUpdateParams{},

		// bank
		"/cosmos.bank.v1beta1.MsgSetSendEnabled":         &bank.MsgSetSendEnabled{},
		"/cosmos.bank.v1beta1.MsgSetSendEnabledResponse": nil,
		"/cosmos.bank.v1beta1.MsgUpdateParams":           &bank.MsgUpdateParams{},
		"/cosmos.bank.v1beta1.MsgUpdateParamsResponse":   nil,

		// consensus
		"/cosmos.consensus.v1.MsgUpdateParams":         &consensus.MsgUpdateParams{},
		"/cosmos.consensus.v1.MsgUpdateParamsResponse": nil,

		// crisis
		"/cosmos.crisis.v1beta1.MsgUpdateParams":         &crisis.MsgUpdateParams{},
		"/cosmos.crisis.v1beta1.MsgUpdateParamsResponse": nil,

		// distribution
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpend":         &distribution.MsgCommunityPoolSpend{},
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpendResponse": nil,
		"/cosmos.distribution.v1beta1.MsgUpdateParams":               &distribution.MsgUpdateParams{},
		"/cosmos.distribution.v1beta1.MsgUpdateParamsResponse":       nil,

		// gov
		"/cosmos.gov.v1.MsgExecLegacyContent":         &gov.MsgExecLegacyContent{},
		"/cosmos.gov.v1.MsgExecLegacyContentResponse": nil,
		"/cosmos.gov.v1.MsgUpdateParams":              &gov.MsgUpdateParams{},
		"/cosmos.gov.v1.MsgUpdateParamsResponse":      nil,

		// slashing
		"/cosmos.slashing.v1beta1.MsgUpdateParams":         &slashing.MsgUpdateParams{},
		"/cosmos.slashing.v1beta1.MsgUpdateParamsResponse": nil,

		// staking
		"/cosmos.staking.v1beta1.MsgUpdateParams":         &staking.MsgUpdateParams{},
		"/cosmos.staking.v1beta1.MsgUpdateParamsResponse": nil,

		// upgrade
		"/cosmos.upgrade.v1beta1.MsgCancelUpgrade":           &upgrade.MsgCancelUpgrade{},
		"/cosmos.upgrade.v1beta1.MsgCancelUpgradeResponse":   nil,
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade":         &upgrade.MsgSoftwareUpgrade{},
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgradeResponse": nil,

		// ibc
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParams":         &icahosttypes.MsgUpdateParams{},
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParamsResponse": nil,
		"/ibc.applications.transfer.v1.MsgUpdateParams":                         &ibctransfer.MsgUpdateParams{},
		"/ibc.applications.transfer.v1.MsgUpdateParamsResponse":                 nil,
		"/ibc.core.client.v1.MsgUpdateParams":                                   &ibcclient.MsgUpdateParams{},
		"/ibc.core.client.v1.MsgUpdateParamsResponse":                           nil,
		"/ibc.core.connection.v1.MsgUpdateParams":                               &ibcconn.MsgUpdateParams{},
		"/ibc.core.connection.v1.MsgUpdateParamsResponse":                       nil,
	}

	// Custom modules
	InternalMsgSamplesDydxCustom = map[string]sdk.Msg{

		// affiliates
		"/nemo-network.affiliates.MsgUpdateAffiliateTiers":         &affiliates.MsgUpdateAffiliateTiers{},
		"/nemo-network.affiliates.MsgUpdateAffiliateTiersResponse": nil,

		// blocktime
		"/nemo_network.blocktime.MsgUpdateDowntimeParams":         &blocktime.MsgUpdateDowntimeParams{},
		"/nemo_network.blocktime.MsgUpdateDowntimeParamsResponse": nil,

		// bridge
		"/nemo_network.bridge.MsgCompleteBridge":              &bridge.MsgCompleteBridge{},
		"/nemo_network.bridge.MsgCompleteBridgeResponse":      nil,
		"/nemo_network.bridge.MsgUpdateEventParams":           &bridge.MsgUpdateEventParams{},
		"/nemo_network.bridge.MsgUpdateEventParamsResponse":   nil,
		"/nemo_network.bridge.MsgUpdateProposeParams":         &bridge.MsgUpdateProposeParams{},
		"/nemo_network.bridge.MsgUpdateProposeParamsResponse": nil,
		"/nemo_network.bridge.MsgUpdateSafetyParams":          &bridge.MsgUpdateSafetyParams{},
		"/nemo_network.bridge.MsgUpdateSafetyParamsResponse":  nil,

		// clob
		"/nemo_network.clob.MsgCreateClobPair":                             &clob.MsgCreateClobPair{},
		"/nemo_network.clob.MsgCreateClobPairResponse":                     nil,
		"/nemo_network.clob.MsgUpdateBlockRateLimitConfiguration":          &clob.MsgUpdateBlockRateLimitConfiguration{},
		"/nemo_network.clob.MsgUpdateBlockRateLimitConfigurationResponse":  nil,
		"/nemo_network.clob.MsgUpdateClobPair":                             &clob.MsgUpdateClobPair{},
		"/nemo_network.clob.MsgUpdateClobPairResponse":                     nil,
		"/nemo_network.clob.MsgUpdateEquityTierLimitConfiguration":         &clob.MsgUpdateEquityTierLimitConfiguration{},
		"/nemo_network.clob.MsgUpdateEquityTierLimitConfigurationResponse": nil,
		"/nemo_network.clob.MsgUpdateLiquidationsConfig":                   &clob.MsgUpdateLiquidationsConfig{},
		"/nemo_network.clob.MsgUpdateLiquidationsConfigResponse":           nil,

		// delaymsg
		"/nemo_network.delaymsg.MsgDelayMessage":         &delaymsg.MsgDelayMessage{},
		"/nemo_network.delaymsg.MsgDelayMessageResponse": nil,

		// feetiers
		"/nemo_network.feetiers.MsgUpdatePerpetualFeeParams":         &feetiers.MsgUpdatePerpetualFeeParams{},
		"/nemo_network.feetiers.MsgUpdatePerpetualFeeParamsResponse": nil,

		// govplus
		"/nemo_network.govplus.MsgSlashValidator":         &govplus.MsgSlashValidator{},
		"/nemo_network.govplus.MsgSlashValidatorResponse": nil,

		// listing
		"/nemo-network.listing.MsgSetMarketsHardCap":                    &listing.MsgSetMarketsHardCap{},
		"/nemo-network.listing.MsgSetMarketsHardCapResponse":            nil,
		"/nemo-network.listing.MsgSetListingVaultDepositParams":         &listing.MsgSetListingVaultDepositParams{},
		"/nemo-network.listing.MsgSetListingVaultDepositParamsResponse": nil,

		// perpetuals
		"/nemo_network.perpetuals.MsgCreatePerpetual":               &perpetuals.MsgCreatePerpetual{},
		"/nemo_network.perpetuals.MsgCreatePerpetualResponse":       nil,
		"/nemo_network.perpetuals.MsgSetLiquidityTier":              &perpetuals.MsgSetLiquidityTier{},
		"/nemo_network.perpetuals.MsgSetLiquidityTierResponse":      nil,
		"/nemo_network.perpetuals.MsgUpdateParams":                  &perpetuals.MsgUpdateParams{},
		"/nemo_network.perpetuals.MsgUpdateParamsResponse":          nil,
		"/nemo_network.perpetuals.MsgUpdatePerpetualParams":         &perpetuals.MsgUpdatePerpetualParams{},
		"/nemo_network.perpetuals.MsgUpdatePerpetualParamsResponse": nil,

		// prices
		"/nemo_network.prices.MsgCreateOracleMarket":         &prices.MsgCreateOracleMarket{},
		"/nemo_network.prices.MsgCreateOracleMarketResponse": nil,
		"/nemo_network.prices.MsgUpdateMarketParam":          &prices.MsgUpdateMarketParam{},
		"/nemo_network.prices.MsgUpdateMarketParamResponse":  nil,

		// ratelimit
		"/nemo_network.ratelimit.MsgSetLimitParams":         &ratelimit.MsgSetLimitParams{},
		"/nemo_network.ratelimit.MsgSetLimitParamsResponse": nil,

		// revshare
		"/nemo-network.revshare.MsgSetMarketMapperRevShareDetailsForMarket":         &revshare.MsgSetMarketMapperRevShareDetailsForMarket{}, //nolint:lll
		"/nemo-network.revshare.MsgSetMarketMapperRevShareDetailsForMarketResponse": nil,
		"/nemo-network.revshare.MsgSetMarketMapperRevenueShare":                     &revshare.MsgSetMarketMapperRevenueShare{}, //nolint:lll
		"/nemo-network.revshare.MsgSetMarketMapperRevenueShareResponse":             nil,

		// rewards
		"/nemo_network.rewards.MsgUpdateParams":         &rewards.MsgUpdateParams{},
		"/nemo_network.rewards.MsgUpdateParamsResponse": nil,

		// sending
		"/nemo_network.sending.MsgSendFromModuleToAccount":         &sending.MsgSendFromModuleToAccount{},
		"/nemo_network.sending.MsgSendFromModuleToAccountResponse": nil,

		// stats
		"/nemo_network.stats.MsgUpdateParams":         &stats.MsgUpdateParams{},
		"/nemo_network.stats.MsgUpdateParamsResponse": nil,

		// vault
		"/nemo-network.vault.MsgSetVaultParams":                     &vault.MsgSetVaultParams{},
		"/nemo-network.vault.MsgSetVaultParamsResponse":             nil,
		"/nemo-network.vault.MsgUpdateDefaultQuotingParams":         &vault.MsgUpdateDefaultQuotingParams{},
		"/nemo-network.vault.MsgUpdateDefaultQuotingParamsResponse": nil,

		// vest
		"/nemo_network.vest.MsgSetVestEntry":            &vest.MsgSetVestEntry{},
		"/nemo_network.vest.MsgSetVestEntryResponse":    nil,
		"/nemo_network.vest.MsgDeleteVestEntry":         &vest.MsgDeleteVestEntry{},
		"/nemo_network.vest.MsgDeleteVestEntryResponse": nil,
	}
)
