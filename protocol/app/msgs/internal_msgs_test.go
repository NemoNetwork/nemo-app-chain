package msgs_test

import (
	"sort"
	"testing"

	"github.com/nemo-network/v4-chain/protocol/app/msgs"
	"github.com/nemo-network/v4-chain/protocol/lib"
	"github.com/stretchr/testify/require"
)

func TestInternalMsgSamples_All_Key(t *testing.T) {
	expectedAllInternalMsgs := lib.MergeAllMapsMustHaveDistinctKeys(msgs.InternalMsgSamplesGovAuth)
	require.Equal(t, expectedAllInternalMsgs, msgs.InternalMsgSamplesAll)
}

func TestInternalMsgSamples_All_Value(t *testing.T) {
	validateMsgValue(t, msgs.InternalMsgSamplesAll)
}

func TestInternalMsgSamples_Gov_Key(t *testing.T) {
	expectedMsgs := []string{
		// auth
		"/cosmos.auth.v1beta1.MsgUpdateParams",

		// bank
		"/cosmos.bank.v1beta1.MsgSetSendEnabled",
		"/cosmos.bank.v1beta1.MsgSetSendEnabledResponse",
		"/cosmos.bank.v1beta1.MsgUpdateParams",
		"/cosmos.bank.v1beta1.MsgUpdateParamsResponse",

		// consensus
		"/cosmos.consensus.v1.MsgUpdateParams",
		"/cosmos.consensus.v1.MsgUpdateParamsResponse",

		// crisis
		"/cosmos.crisis.v1beta1.MsgUpdateParams",
		"/cosmos.crisis.v1beta1.MsgUpdateParamsResponse",

		// distribution
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpend",
		"/cosmos.distribution.v1beta1.MsgCommunityPoolSpendResponse",
		"/cosmos.distribution.v1beta1.MsgUpdateParams",
		"/cosmos.distribution.v1beta1.MsgUpdateParamsResponse",

		// gov
		"/cosmos.gov.v1.MsgExecLegacyContent",
		"/cosmos.gov.v1.MsgExecLegacyContentResponse",
		"/cosmos.gov.v1.MsgUpdateParams",
		"/cosmos.gov.v1.MsgUpdateParamsResponse",

		// slashing
		"/cosmos.slashing.v1beta1.MsgUpdateParams",
		"/cosmos.slashing.v1beta1.MsgUpdateParamsResponse",

		// staking
		"/cosmos.staking.v1beta1.MsgUpdateParams",
		"/cosmos.staking.v1beta1.MsgUpdateParamsResponse",

		// upgrade
		"/cosmos.upgrade.v1beta1.MsgCancelUpgrade",
		"/cosmos.upgrade.v1beta1.MsgCancelUpgradeResponse",
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgrade",
		"/cosmos.upgrade.v1beta1.MsgSoftwareUpgradeResponse",

		// blocktime
		"/nemo_network.blocktime.MsgUpdateDowntimeParams",
		"/nemo_network.blocktime.MsgUpdateDowntimeParamsResponse",

		// bridge
		"/nemo_network.bridge.MsgCompleteBridge",
		"/nemo_network.bridge.MsgCompleteBridgeResponse",
		"/nemo_network.bridge.MsgUpdateEventParams",
		"/nemo_network.bridge.MsgUpdateEventParamsResponse",
		"/nemo_network.bridge.MsgUpdateProposeParams",
		"/nemo_network.bridge.MsgUpdateProposeParamsResponse",
		"/nemo_network.bridge.MsgUpdateSafetyParams",
		"/nemo_network.bridge.MsgUpdateSafetyParamsResponse",

		// clob
		"/nemo_network.clob.MsgCreateClobPair",
		"/nemo_network.clob.MsgCreateClobPairResponse",
		"/nemo_network.clob.MsgUpdateBlockRateLimitConfiguration",
		"/nemo_network.clob.MsgUpdateBlockRateLimitConfigurationResponse",
		"/nemo_network.clob.MsgUpdateClobPair",
		"/nemo_network.clob.MsgUpdateClobPairResponse",
		"/nemo_network.clob.MsgUpdateEquityTierLimitConfiguration",
		"/nemo_network.clob.MsgUpdateEquityTierLimitConfigurationResponse",
		"/nemo_network.clob.MsgUpdateLiquidationsConfig",
		"/nemo_network.clob.MsgUpdateLiquidationsConfigResponse",

		// delaymsg
		"/nemo_network.delaymsg.MsgDelayMessage",
		"/nemo_network.delaymsg.MsgDelayMessageResponse",

		// feetiers
		"/nemo_network.feetiers.MsgUpdatePerpetualFeeParams",
		"/nemo_network.feetiers.MsgUpdatePerpetualFeeParamsResponse",

		// govplus
		"/nemo_network.govplus.MsgSlashValidator",
		"/nemo_network.govplus.MsgSlashValidatorResponse",

		// perpeutals
		"/nemo_network.perpetuals.MsgCreatePerpetual",
		"/nemo_network.perpetuals.MsgCreatePerpetualResponse",
		"/nemo_network.perpetuals.MsgSetLiquidityTier",
		"/nemo_network.perpetuals.MsgSetLiquidityTierResponse",
		"/nemo_network.perpetuals.MsgUpdateParams",
		"/nemo_network.perpetuals.MsgUpdateParamsResponse",
		"/nemo_network.perpetuals.MsgUpdatePerpetualParams",
		"/nemo_network.perpetuals.MsgUpdatePerpetualParamsResponse",

		// prices
		"/nemo_network.prices.MsgCreateOracleMarket",
		"/nemo_network.prices.MsgCreateOracleMarketResponse",
		"/nemo_network.prices.MsgUpdateMarketParam",
		"/nemo_network.prices.MsgUpdateMarketParamResponse",

		// ratelimit
		"/nemo_network.ratelimit.MsgSetLimitParams",
		"/nemo_network.ratelimit.MsgSetLimitParamsResponse",

		// rewards
		"/nemo_network.rewards.MsgUpdateParams",
		"/nemo_network.rewards.MsgUpdateParamsResponse",

		// sending
		"/nemo_network.sending.MsgSendFromModuleToAccount",
		"/nemo_network.sending.MsgSendFromModuleToAccountResponse",

		// stats
		"/nemo_network.stats.MsgUpdateParams",
		"/nemo_network.stats.MsgUpdateParamsResponse",

		// vault
		"/nemo_network.vault.MsgUpdateParams",
		"/nemo_network.vault.MsgUpdateParamsResponse",

		// vest
		"/nemo_network.vest.MsgDeleteVestEntry",
		"/nemo_network.vest.MsgDeleteVestEntryResponse",
		"/nemo_network.vest.MsgSetVestEntry",
		"/nemo_network.vest.MsgSetVestEntryResponse",

		// ibc
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParams",
		"/ibc.applications.interchain_accounts.host.v1.MsgUpdateParamsResponse",
		"/ibc.applications.transfer.v1.MsgUpdateParams",
		"/ibc.applications.transfer.v1.MsgUpdateParamsResponse",
		"/ibc.core.client.v1.MsgUpdateParams",
		"/ibc.core.client.v1.MsgUpdateParamsResponse",
		"/ibc.core.connection.v1.MsgUpdateParams",
		"/ibc.core.connection.v1.MsgUpdateParamsResponse",
	}

	require.Equal(t, expectedMsgs, lib.GetSortedKeys[sort.StringSlice](msgs.InternalMsgSamplesGovAuth))
}

func TestInternalMsgSamples_Gov_Value(t *testing.T) {
	validateMsgValue(t, msgs.InternalMsgSamplesGovAuth)
}
