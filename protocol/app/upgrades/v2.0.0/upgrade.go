package v_2_0_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	ccvconsumerkeeper "github.com/cosmos/interchain-security/v5/x/ccv/consumer/keeper"
	ccvconsumertypes "github.com/cosmos/interchain-security/v5/x/ccv/consumer/types"
	ccvtypes "github.com/cosmos/interchain-security/v5/x/ccv/types"

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
// 2. Builds a PreCCV consumer genesis from the current staking validator set.
// 3. Initializes the CCV consumer module with the PreCCV genesis state.
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

		// 2. Build the initial validator set from the current staking module state.
		initialValSet, err := getInitialValidatorSet(ctx, stakingKeeper)
		if err != nil {
			return vm, fmt.Errorf("getting initial validator set: %w", err)
		}

		// 3. Construct a PreCCV consumer genesis state for standalone-to-consumer changeover.
		consumerGenesis := ccvconsumertypes.DefaultGenesisState()
		consumerGenesis.PreCCV = true
		consumerGenesis.Params = ccvtypes.DefaultParams()
		consumerGenesis.Params.Enabled = true
		consumerGenesis.Provider.InitialValSet = initialValSet

		if err := consumerGenesis.Validate(); err != nil {
			return vm, fmt.Errorf("validating consumer genesis: %w", err)
		}

		// 4. Initialize the CCV consumer module with the PreCCV genesis.
		consumerKeeper.InitGenesis(sdkCtx, consumerGenesis)

		sdkCtx.Logger().Info("v2.0.0 upgrade: CCV consumer module initialized (PreCCV=true)",
			"initial_valset_size", len(initialValSet),
		)

		return vm, nil
	}
}

// getInitialValidatorSet builds the initial validator set from the current staking state.
func getInitialValidatorSet(ctx context.Context, k *stakingkeeper.Keeper) ([]abci.ValidatorUpdate, error) {
	vals, err := k.GetLastValidators(ctx)
	if err != nil {
		return nil, err
	}
	powerReduction := k.PowerReduction(ctx)

	updates := make([]abci.ValidatorUpdate, 0, len(vals))
	for _, val := range vals {
		updates = append(updates, val.ABCIValidatorUpdate(powerReduction))
	}
	return updates, nil
}
