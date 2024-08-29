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
	blocktime "github.com/nemo-network/v4-chain/protocol/x/blocktime/types"
	bridge "github.com/nemo-network/v4-chain/protocol/x/bridge/types"
	clob "github.com/nemo-network/v4-chain/protocol/x/clob/types"
	delaymsg "github.com/nemo-network/v4-chain/protocol/x/delaymsg/types"
	feetiers "github.com/nemo-network/v4-chain/protocol/x/feetiers/types"
	govplus "github.com/nemo-network/v4-chain/protocol/x/govplus/types"
	perpetuals "github.com/nemo-network/v4-chain/protocol/x/perpetuals/types"
	prices "github.com/nemo-network/v4-chain/protocol/x/prices/types"
	ratelimit "github.com/nemo-network/v4-chain/protocol/x/ratelimit/types"
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
		// blocktime
		"/nemo-network.blocktime.MsgUpdateDowntimeParams":         &blocktime.MsgUpdateDowntimeParams{},
		"/nemo-network.blocktime.MsgUpdateDowntimeParamsResponse": nil,

		// bridge
		"/nemo-network.bridge.MsgCompleteBridge":              &bridge.MsgCompleteBridge{},
		"/nemo-network.bridge.MsgCompleteBridgeResponse":      nil,
		"/nemo-network.bridge.MsgUpdateEventParams":           &bridge.MsgUpdateEventParams{},
		"/nemo-network.bridge.MsgUpdateEventParamsResponse":   nil,
		"/nemo-network.bridge.MsgUpdateProposeParams":         &bridge.MsgUpdateProposeParams{},
		"/nemo-network.bridge.MsgUpdateProposeParamsResponse": nil,
		"/nemo-network.bridge.MsgUpdateSafetyParams":          &bridge.MsgUpdateSafetyParams{},
		"/nemo-network.bridge.MsgUpdateSafetyParamsResponse":  nil,

		// clob
		"/nemo-network.clob.MsgCreateClobPair":                             &clob.MsgCreateClobPair{},
		"/nemo-network.clob.MsgCreateClobPairResponse":                     nil,
		"/nemo-network.clob.MsgUpdateBlockRateLimitConfiguration":          &clob.MsgUpdateBlockRateLimitConfiguration{},
		"/nemo-network.clob.MsgUpdateBlockRateLimitConfigurationResponse":  nil,
		"/nemo-network.clob.MsgUpdateClobPair":                             &clob.MsgUpdateClobPair{},
		"/nemo-network.clob.MsgUpdateClobPairResponse":                     nil,
		"/nemo-network.clob.MsgUpdateEquityTierLimitConfiguration":         &clob.MsgUpdateEquityTierLimitConfiguration{},
		"/nemo-network.clob.MsgUpdateEquityTierLimitConfigurationResponse": nil,
		"/nemo-network.clob.MsgUpdateLiquidationsConfig":                   &clob.MsgUpdateLiquidationsConfig{},
		"/nemo-network.clob.MsgUpdateLiquidationsConfigResponse":           nil,

		// delaymsg
		"/nemo-network.delaymsg.MsgDelayMessage":         &delaymsg.MsgDelayMessage{},
		"/nemo-network.delaymsg.MsgDelayMessageResponse": nil,

		// feetiers
		"/nemo-network.feetiers.MsgUpdatePerpetualFeeParams":         &feetiers.MsgUpdatePerpetualFeeParams{},
		"/nemo-network.feetiers.MsgUpdatePerpetualFeeParamsResponse": nil,

		// govplus
		"/nemo-network.govplus.MsgSlashValidator":         &govplus.MsgSlashValidator{},
		"/nemo-network.govplus.MsgSlashValidatorResponse": nil,

		// perpetuals
		"/nemo-network.perpetuals.MsgCreatePerpetual":               &perpetuals.MsgCreatePerpetual{},
		"/nemo-network.perpetuals.MsgCreatePerpetualResponse":       nil,
		"/nemo-network.perpetuals.MsgSetLiquidityTier":              &perpetuals.MsgSetLiquidityTier{},
		"/nemo-network.perpetuals.MsgSetLiquidityTierResponse":      nil,
		"/nemo-network.perpetuals.MsgUpdateParams":                  &perpetuals.MsgUpdateParams{},
		"/nemo-network.perpetuals.MsgUpdateParamsResponse":          nil,
		"/nemo-network.perpetuals.MsgUpdatePerpetualParams":         &perpetuals.MsgUpdatePerpetualParams{},
		"/nemo-network.perpetuals.MsgUpdatePerpetualParamsResponse": nil,

		// prices
		"/nemo-network.prices.MsgCreateOracleMarket":         &prices.MsgCreateOracleMarket{},
		"/nemo-network.prices.MsgCreateOracleMarketResponse": nil,
		"/nemo-network.prices.MsgUpdateMarketParam":          &prices.MsgUpdateMarketParam{},
		"/nemo-network.prices.MsgUpdateMarketParamResponse":  nil,

		// ratelimit
		"/nemo-network.ratelimit.MsgSetLimitParams":         &ratelimit.MsgSetLimitParams{},
		"/nemo-network.ratelimit.MsgSetLimitParamsResponse": nil,

		// rewards
		"/nemo-network.rewards.MsgUpdateParams":         &rewards.MsgUpdateParams{},
		"/nemo-network.rewards.MsgUpdateParamsResponse": nil,

		// sending
		"/nemo-network.sending.MsgSendFromModuleToAccount":         &sending.MsgSendFromModuleToAccount{},
		"/nemo-network.sending.MsgSendFromModuleToAccountResponse": nil,

		// stats
		"/nemo-network.stats.MsgUpdateParams":         &stats.MsgUpdateParams{},
		"/nemo-network.stats.MsgUpdateParamsResponse": nil,

		// vault
		"/nemo-network.vault.MsgUpdateParams":         &vault.MsgUpdateParams{},
		"/nemo-network.vault.MsgUpdateParamsResponse": nil,

		// vest
		"/nemo-network.vest.MsgSetVestEntry":            &vest.MsgSetVestEntry{},
		"/nemo-network.vest.MsgSetVestEntryResponse":    nil,
		"/nemo-network.vest.MsgDeleteVestEntry":         &vest.MsgDeleteVestEntry{},
		"/nemo-network.vest.MsgDeleteVestEntryResponse": nil,
	}
)
