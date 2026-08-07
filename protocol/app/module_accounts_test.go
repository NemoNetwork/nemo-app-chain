package app_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/nemo-network/v4-chain/protocol/app"
	bridgemoduletypes "github.com/nemo-network/v4-chain/protocol/x/bridge/types"
	perpetualsmoduletypes "github.com/nemo-network/v4-chain/protocol/x/perpetuals/types"
	rewardsmoduletypes "github.com/nemo-network/v4-chain/protocol/x/rewards/types"
	satypes "github.com/nemo-network/v4-chain/protocol/x/subaccounts/types"
	vaultmoduletypes "github.com/nemo-network/v4-chain/protocol/x/vault/types"
	vestmoduletypes "github.com/nemo-network/v4-chain/protocol/x/vest/types"
	marketmapmoduletypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestModuleAccountsToAddresses(t *testing.T) {
	expectedModuleAccToAddresses := map[string]string{
		authtypes.FeeCollectorName:                   "nemo17xpfvakm2amg962yls6f84z3kell8c5lmq203q",
		bridgemoduletypes.ModuleName:                 "nemo1zlefkpe3g0vvm9a4h0jf9000lmqutlh9sw4c2x",
		distrtypes.ModuleName:                        "nemo1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vxdnnz",
		stakingtypes.BondedPoolName:                  "nemo1fl48vsnmsdzcv85q5d2q4z5ajdha8yu37zqqr2",
		stakingtypes.NotBondedPoolName:               "nemo1tygms3xhhs3yv487phx3dw4a95jn7t7l2zu347",
		govtypes.ModuleName:                          "nemo10d07y265gmmuvt4z0w9aw880jnsr700j3m62vw",
		ibctransfertypes.ModuleName:                  "nemo1yl6hdjhmkf37639730gffanpzndzdpmh9xlxd7",
		satypes.ModuleName:                           "nemo1v88c3xv9xyv3eetdx0tvcmq7ung3dywpkux9zs",
		perpetualsmoduletypes.InsuranceFundName:      "nemo1c7ptc87hkd54e3r7zjy92q29xkq7t79wc4h5e2",
		rewardsmoduletypes.TreasuryAccountName:       "nemo16wrau2x4tsg033xfrrdpae6kxfn9kyueprnegt",
		rewardsmoduletypes.VesterAccountName:         "nemo1ltyc6y4skclzafvpznpt2qjwmfwgsndph5qgpt",
		vestmoduletypes.CommunityTreasuryAccountName: "nemo15ztc7xy42tn2ukkc0qjthkucw9ac63pgr7ghee",
		vestmoduletypes.CommunityVesterAccountName:   "nemo1wxje320an3karyc6mjw4zghs300dmrjkvned3u",
		icatypes.ModuleName:                          "nemo1vlthgax23ca9syk7xgaz347xmf4nunefv3lckd",
		marketmapmoduletypes.ModuleName:              "nemo16j3d86dww8p2rzdlqsv7wle98cxzjxw62j40ce",
		vaultmoduletypes.MegavaultAccountName:        "nemo18tkxrnrkqc2t0lr3zxr5g6a4hdvqksylyq4j0f",
	}

	require.True(t, len(expectedModuleAccToAddresses) == len(app.GetMaccPerms()),
		"expected %d, got %d", len(expectedModuleAccToAddresses), len(app.GetMaccPerms()))
	for acc, address := range expectedModuleAccToAddresses {
		expectedAddr := authtypes.NewModuleAddress(acc).String()
		require.Equal(t, address, expectedAddr, "module (%v) should have address (%s)", acc, expectedAddr)
	}
}

func TestBlockedAddresses(t *testing.T) {
	expectedBlockedAddresses := map[string]bool{
		"nemo17xpfvakm2amg962yls6f84z3kell8c5lmq203q": true,
		"nemo1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vxdnnz": true,
		"nemo1tygms3xhhs3yv487phx3dw4a95jn7t7l2zu347": true,
		"nemo1fl48vsnmsdzcv85q5d2q4z5ajdha8yu37zqqr2": true,
		"nemo1yl6hdjhmkf37639730gffanpzndzdpmh9xlxd7": true,
		"nemo1vlthgax23ca9syk7xgaz347xmf4nunefv3lckd": true,
	}
	require.Equal(t, expectedBlockedAddresses, app.BlockedAddresses())
}

func TestMaccPerms(t *testing.T) {
	maccPerms := app.GetMaccPerms()
	expectedMaccPerms := map[string][]string{
		"bonded_tokens_pool":     {"burner", "staking"},
		"bridge":                 {"minter"},
		"distribution":           nil,
		"fee_collector":          nil,
		"gov":                    {"burner"},
		"insurance_fund":         nil,
		"not_bonded_tokens_pool": {"burner", "staking"},
		"subaccounts":            nil,
		"transfer":               {"minter", "burner"},
		"interchainaccounts":     nil,
		"rewards_treasury":       nil,
		"rewards_vester":         nil,
		"community_treasury":     nil,
		"community_vester":       nil,
		"marketmap":              nil,
		"megavault":              nil,
	}
	require.Equal(t, expectedMaccPerms, maccPerms, "default macc perms list does not match expected")
}

func TestModuleAccountAddrs(t *testing.T) {
	expectedModuleAccAddresses := map[string]bool{
		"nemo17xpfvakm2amg962yls6f84z3kell8c5lmq203q": true, // x/auth.FeeCollector
		"nemo1zlefkpe3g0vvm9a4h0jf9000lmqutlh9sw4c2x": true, // x/bridge
		"nemo1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8vxdnnz": true, // x/distribution
		"nemo1fl48vsnmsdzcv85q5d2q4z5ajdha8yu37zqqr2": true, // x/staking.bondedPool
		"nemo1tygms3xhhs3yv487phx3dw4a95jn7t7l2zu347": true, // x/staking.notBondedPool
		"nemo10d07y265gmmuvt4z0w9aw880jnsr700j3m62vw": true, // x/ gov
		"nemo1yl6hdjhmkf37639730gffanpzndzdpmh9xlxd7": true, // ibc transfer
		"nemo1vlthgax23ca9syk7xgaz347xmf4nunefv3lckd": true, // interchainaccounts
		"nemo1v88c3xv9xyv3eetdx0tvcmq7ung3dywpkux9zs": true, // x/subaccount
		"nemo1c7ptc87hkd54e3r7zjy92q29xkq7t79wc4h5e2": true, // x/clob.insuranceFund
		"nemo16wrau2x4tsg033xfrrdpae6kxfn9kyueprnegt": true, // x/rewards.treasury
		"nemo1ltyc6y4skclzafvpznpt2qjwmfwgsndph5qgpt": true, // x/rewards.vester
		"nemo15ztc7xy42tn2ukkc0qjthkucw9ac63pgr7ghee": true, // x/vest.communityTreasury
		"nemo1wxje320an3karyc6mjw4zghs300dmrjkvned3u": true, // x/vest.communityVester
		"nemo16j3d86dww8p2rzdlqsv7wle98cxzjxw62j40ce": true, // x/marketmap
		"nemo18tkxrnrkqc2t0lr3zxr5g6a4hdvqksylyq4j0f": true, // x/vault.megavault
	}

	require.Equal(t, expectedModuleAccAddresses, app.ModuleAccountAddrs())
}
