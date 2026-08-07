package v_1_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/nemo-network/v4-chain/protocol/dtypes"
	"github.com/nemo-network/v4-chain/protocol/lib"
	vaultkeeper "github.com/nemo-network/v4-chain/protocol/x/vault/keeper"
	vaulttypes "github.com/nemo-network/v4-chain/protocol/x/vault/types"
)

// initializeVaultParams writes the megavault state that `InitGenesis` would have
// set but that a chain running from an older genesis does not have.
//
// Every write here is conditional: an upgrade must be safe to run against state
// that already carries the value (for instance, after a genesis export/import
// round-trip), and must never overwrite an operator or a cap that governance has
// already chosen.
func initializeVaultParams(ctx sdk.Context, k vaultkeeper.Keeper) {
	// 1. Operator params.
	//
	// A chain whose genesis predates operator params reads back an empty
	// operator. Default it to the gov module account, matching
	// `DefaultOperatorParams`, so that `MsgSetVaultParams` and friends have a
	// valid authority from the first block after the upgrade.
	//
	// Fee parameters are deliberately left at zero. The fee engine ships inert
	// and is turned on by a separate governance action, which is mitigation 3 of
	// ADR-001 §7 — it decouples "the mechanism exists" from "the mechanism
	// charges" and leaves room to announce the change.
	operatorParams := k.GetOperatorParams(ctx)
	if operatorParams.Operator == "" {
		operatorParams = vaulttypes.DefaultOperatorParams()
		if err := k.SetOperatorParams(ctx, operatorParams); err != nil {
			panic(fmt.Sprintf("failed to initialize vault operator params: %s", err))
		}
		ctx.Logger().Info(fmt.Sprintf(
			"v1.0 upgrade: initialized megavault operator to %s",
			operatorParams.Operator,
		))
	}

	// 2. Megavault params.
	//
	// A zero deposit cap already means "uncapped", so reading unset state is
	// correct on its own. Write it explicitly anyway so that an exported genesis
	// after the upgrade round-trips identically to one produced by `InitGenesis`.
	if err := k.SetMegavaultParams(ctx, k.GetMegavaultParams(ctx)); err != nil {
		panic(fmt.Sprintf("failed to initialize megavault params: %s", err))
	}
}

// initializeFeeState seeds the fee accounting state.
//
// This is the single most consequential step in the upgrade. The high-water mark
// is seeded at the NAV per share observed *at upgrade height*, so that only gains
// made after the upgrade can ever be charged a profit share. A high-water mark of
// zero would charge the operator's profit share against the entire trading
// history of the vault, retroactively, on depositors who deposited under a
// 0%-fee regime. See ADR-001 §7 and D8.
func initializeFeeState(ctx sdk.Context, k vaultkeeper.Keeper) {
	feeState := k.GetFeeState(ctx)

	navPerShare, _, _, exists, err := k.GetNavPerShare(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to read megavault NAV per share: %s", err))
	}

	if exists {
		// Only raise the mark; never lower one that is already set.
		currentHwm := feeState.HighWaterMarkNavPerShare.BigInt()
		if currentHwm == nil || currentHwm.Cmp(navPerShare) < 0 {
			feeState.HighWaterMarkNavPerShare = dtypes.NewIntFromBigInt(navPerShare)
		}
	} else if feeState.HighWaterMarkNavPerShare.IsNil() {
		feeState.HighWaterMarkNavPerShare = dtypes.NewInt(0)
	}

	// Start the accrual clock at the upgrade block time. Leaving it at zero
	// would make the first accrual charge for every second since the Unix epoch;
	// `MaxAccrualElapsedSeconds` caps that at a year, which is still wrong.
	if feeState.LastAccrualTime == 0 {
		feeState.LastAccrualTime = ctx.BlockTime().Unix()
	}

	if err := k.SetFeeState(ctx, feeState); err != nil {
		panic(fmt.Sprintf("failed to initialize megavault fee state: %s", err))
	}

	ctx.Logger().Info(fmt.Sprintf(
		"v1.0 upgrade: seeded megavault high-water mark at %s (scaled by 1e18), accrual clock at %d",
		feeState.HighWaterMarkNavPerShare.BigInt(),
		feeState.LastAccrualTime,
	))
}

// backfillIndexerVaults re-sets every existing vault's params so that
// `SetVaultParams` emits an `UpsertVault` indexer event for each one.
//
// Vaults that already exist have never emitted that event — it was added in the
// same release as this upgrade — so the indexer's `vaults` table would otherwise
// start empty and the vault endpoints would report nothing. Re-setting the params
// is a no-op on chain state and produces exactly the events the indexer needs.
func backfillIndexerVaults(ctx sdk.Context, k vaultkeeper.Keeper) {
	vaults := k.GetAllVaults(ctx)
	for _, vault := range vaults {
		if err := k.SetVaultParams(ctx, vault.VaultId, vault.VaultParams); err != nil {
			// A vault already in state must by definition have valid params, so
			// this can only mean a validation rule changed under it. Log loudly
			// rather than halt the chain: the cost is one vault missing from the
			// indexer, not a failed upgrade.
			ctx.Logger().Error(fmt.Sprintf(
				"v1.0 upgrade: failed to backfill indexer event for vault %s: %s",
				vault.VaultId.ToString(),
				err,
			))
			continue
		}
	}

	ctx.Logger().Info(fmt.Sprintf(
		"v1.0 upgrade: emitted upsert_vault events for %d vaults",
		len(vaults),
	))
}

// MigrateVaultState performs this upgrade's vault state migration.
//
// It is split out of the handler so that tests can exercise the real migration
// without standing up a module manager, and so the ordering is stated in one
// place: params before fee state (the fee state seeds from NAV, which is
// independent), and the indexer backfill last so it emits events against the
// final state.
func MigrateVaultState(ctx sdk.Context, vaultKeeper vaultkeeper.Keeper) {
	initializeVaultParams(ctx, vaultKeeper)
	initializeFeeState(ctx, vaultKeeper)
	backfillIndexerVaults(ctx, vaultKeeper)
}

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	vaultKeeper vaultkeeper.Keeper,
) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades")
		sdkCtx.Logger().Info(fmt.Sprintf("Running %s Upgrade...", UpgradeName))

		MigrateVaultState(sdkCtx, vaultKeeper)

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
