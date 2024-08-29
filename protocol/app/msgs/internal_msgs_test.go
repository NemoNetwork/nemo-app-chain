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
		"/nemo-network.blocktime.MsgUpdateDowntimeParams",
		"/nemo-network.blocktime.MsgUpdateDowntimeParamsResponse",

		// bridge
		"/nemo-network.bridge.MsgCompleteBridge",
		"/nemo-network.bridge.MsgCompleteBridgeResponse",
		"/nemo-network.bridge.MsgUpdateEventParams",
		"/nemo-network.bridge.MsgUpdateEventParamsResponse",
		"/nemo-network.bridge.MsgUpdateProposeParams",
		"/nemo-network.bridge.MsgUpdateProposeParamsResponse",
		"/nemo-network.bridge.MsgUpdateSafetyParams",
		"/nemo-network.bridge.MsgUpdateSafetyParamsResponse",

		// clob
		"/nemo-network.clob.MsgCreateClobPair",
		"/nemo-network.clob.MsgCreateClobPairResponse",
		"/nemo-network.clob.MsgUpdateBlockRateLimitConfiguration",
		"/nemo-network.clob.MsgUpdateBlockRateLimitConfigurationResponse",
		"/nemo-network.clob.MsgUpdateClobPair",
		"/nemo-network.clob.MsgUpdateClobPairResponse",
		"/nemo-network.clob.MsgUpdateEquityTierLimitConfiguration",
		"/nemo-network.clob.MsgUpdateEquityTierLimitConfigurationResponse",
		"/nemo-network.clob.MsgUpdateLiquidationsConfig",
		"/nemo-network.clob.MsgUpdateLiquidationsConfigResponse",

		// delaymsg
		"/nemo-network.delaymsg.MsgDelayMessage",
		"/nemo-network.delaymsg.MsgDelayMessageResponse",

		// feetiers
		"/nemo-network.feetiers.MsgUpdatePerpetualFeeParams",
		"/nemo-network.feetiers.MsgUpdatePerpetualFeeParamsResponse",

		// govplus
		"/nemo-network.govplus.MsgSlashValidator",
		"/nemo-network.govplus.MsgSlashValidatorResponse",

		// perpeutals
		"/nemo-network.perpetuals.MsgCreatePerpetual",
		"/nemo-network.perpetuals.MsgCreatePerpetualResponse",
		"/nemo-network.perpetuals.MsgSetLiquidityTier",
		"/nemo-network.perpetuals.MsgSetLiquidityTierResponse",
		"/nemo-network.perpetuals.MsgUpdateParams",
		"/nemo-network.perpetuals.MsgUpdateParamsResponse",
		"/nemo-network.perpetuals.MsgUpdatePerpetualParams",
		"/nemo-network.perpetuals.MsgUpdatePerpetualParamsResponse",

		// prices
		"/nemo-network.prices.MsgCreateOracleMarket",
		"/nemo-network.prices.MsgCreateOracleMarketResponse",
		"/nemo-network.prices.MsgUpdateMarketParam",
		"/nemo-network.prices.MsgUpdateMarketParamResponse",

		// ratelimit
		"/nemo-network.ratelimit.MsgSetLimitParams",
		"/nemo-network.ratelimit.MsgSetLimitParamsResponse",

		// rewards
		"/nemo-network.rewards.MsgUpdateParams",
		"/nemo-network.rewards.MsgUpdateParamsResponse",

		// sending
		"/nemo-network.sending.MsgSendFromModuleToAccount",
		"/nemo-network.sending.MsgSendFromModuleToAccountResponse",

		// stats
		"/nemo-network.stats.MsgUpdateParams",
		"/nemo-network.stats.MsgUpdateParamsResponse",

		// vault
		"/nemo-network.vault.MsgUpdateParams",
		"/nemo-network.vault.MsgUpdateParamsResponse",

		// vest
		"/nemo-network.vest.MsgDeleteVestEntry",
		"/nemo-network.vest.MsgDeleteVestEntryResponse",
		"/nemo-network.vest.MsgSetVestEntry",
		"/nemo-network.vest.MsgSetVestEntryResponse",

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
