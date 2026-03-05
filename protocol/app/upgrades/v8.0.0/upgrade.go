package v_8_0_0

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	ccvconsumerkeeper "github.com/cosmos/interchain-security/v5/x/ccv/consumer/keeper"
	ccvconsumertypes "github.com/cosmos/interchain-security/v5/x/ccv/consumer/types"
	ccv "github.com/cosmos/interchain-security/v5/x/ccv/types"

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
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades/v8.0.0")

		// 1. Run migrations for all modules (including registering the new consumer module in the version map).
		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		// 2. Initialize the CCV consumer module state for a chain that was previously standalone.
		// PreCCV=true: consumer module will not drive validator set until the provider CCV channel is established.
		// InitialValSet: current staking validator set, used when changeover completes.
		initialValSet, err := getInitialValidatorSet(ctx, stakingKeeper)
		if err != nil {
			return vm, err
		}

		consumerGenesis := ccvconsumertypes.DefaultGenesisState()
		consumerGenesis.PreCCV = true
		consumerGenesis.Params = ccv.DefaultParams()
		consumerGenesis.Provider.InitialValSet = initialValSet

		if err := consumerGenesis.Validate(); err != nil {
			return vm, err
		}

		consumerKeeper.InitGenesis(sdkCtx, consumerGenesis)
		sdkCtx.Logger().Info("v8.0.0 upgrade: CCV consumer module initialized (PreCCV=true)")

		return vm, nil
	}
}

// getInitialValidatorSet returns the current bonded validator set as abci.ValidatorUpdate slice.
func getInitialValidatorSet(ctx context.Context, k *stakingkeeper.Keeper) ([]abci.ValidatorUpdate, error) {
	validators, err := k.GetLastValidators(ctx)
	if err != nil {
		return nil, err
	}
	powerReduction := k.PowerReduction(ctx)
	updates := make([]abci.ValidatorUpdate, 0, len(validators))
	for _, val := range validators {
		updates = append(updates, val.ABCIValidatorUpdate(powerReduction))
	}
	return updates, nil
}
