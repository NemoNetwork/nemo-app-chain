package v_2_0_0

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	ccvconsumerkeeper "github.com/cosmos/interchain-security/v5/x/ccv/consumer/keeper"

	"github.com/nemo-network/v4-chain/protocol/lib"
)

// CreateUpgradeHandler returns an upgrade handler that only runs module migrations.
// Used when the upgrade is registered without consumer init (e.g. for store loader compatibility).
func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		return mm.RunMigrations(ctx, configurator, vm)
	}
}

// CreateUpgradeHandlerWithConsumerInit returns an upgrade handler that:
// 1. Runs module migrations.
// 2. Initializes the CCV consumer module state for a standalone->consumer changeover (PreCCV=true, initial val set from staking).
func CreateUpgradeHandlerWithConsumerInit(
	mm *module.Manager,
	configurator module.Configurator,
	consumerKeeper ccvconsumerkeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades/v2.0.0")

		// 1. Run migrations for all modules (including registering the new consumer module in the version map).
		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		sdkCtx.Logger().Info("v2.0.0 upgrade: CCV consumer module addition")

		return vm, nil
	}
}
