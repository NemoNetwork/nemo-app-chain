package v_1_0

import (
	store "cosmossdk.io/store/types"
	"github.com/nemo-network/v4-chain/protocol/app/upgrades"
)

const (
	// UpgradeName is the name of the first software upgrade this chain ships.
	// It carries the megavault work: operator params, the deposit cap, the fee
	// engine, and the indexer `vaults` table backfill.
	UpgradeName = "v1.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName: UpgradeName,
	// No store upgrades. Every piece of state this upgrade introduces
	// (`OperatorParams`, `MegavaultParams`, `FeeState`) lives under the existing
	// `vault` store key as a new key prefix, so no store is added, removed or
	// renamed.
	StoreUpgrades: store.StoreUpgrades{},
}
