package v_2_0_0

import (
	store "cosmossdk.io/store/types"
	ccvconsumertypes "github.com/cosmos/interchain-security/v5/x/ccv/consumer/types"
	"github.com/nemo-network/v4-chain/protocol/app/upgrades"
)

const (
	// UpgradeName is the name of the upgrade that adds the ICS CCV consumer module
	// (standalone chain -> PSS opt-in consumer changeover).
	UpgradeName = "v2.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			ccvconsumertypes.StoreKey,
		},
	},
}
