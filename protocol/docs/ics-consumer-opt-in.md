## Nemo App Chain – Opt-In Consumer of Cosmos Hub (Local Setup)

This document explains, step by step, how to turn `nemo-app-chain` into an **ICS 2.0 Partial-Set Security (PSS) Opt-In consumer** of a local Cosmos Hub (`gaia`) using a local Hermes IBC relayer.

Assumptions:

- **Chains**
  - Provider: `gaia-local-1` (Cosmos Hub with `x/ccv/provider` enabled)
  - Consumer: `nemo-local-1` (Nemo app chain, initially a standalone chain)
- **Binaries**
  - `gaiad` – Cosmos Hub (provider)
  - `nemod` – Nemo app chain (consumer)
  - `hermes` – IBC relayer
- All commands are written for **local single-validator setups**; adapt paths, keys, and chain-ids as needed.

---

## 1. Code-Level Changes on Nemo: Wiring the ICS Consumer Module

This section summarizes what changes are required in the Nemo codebase to support ICS 2.0 consumer mode and PSS.

### 1.1. Align `go.mod` for ICS 2.0 Consumer

In `protocol/go.mod`:

- **Add Interchain Security dependency**:

```go
require (
  github.com/cosmos/cosmos-sdk v0.50.9
  github.com/cosmos/interchain-security/v5 v5.1.1
  // ...
)
```

- **Ensure IBC-Go is compatible with ICS**:
  - Remove any `replace github.com/cosmos/ibc-go/v8 => ...` overrides that pin a conflicting version.
  - Let ICS pull in the required `ibc-go/v8.3.x` dependency.

Then run:

```bash
cd /root/dapp/nemo-app-chain/protocol
go mod tidy
```

This resolves all ICS and IBC dependencies to versions compatible with PSS.

### 1.2. Add the CCV Consumer Module to `app.go`

In `protocol/app/app.go`:

- **Imports**:

```go
import (
  // ...
  ccvconsumermodule "github.com/cosmos/interchain-security/v5/x/ccv/consumer"
  ccvconsumerkeeper "github.com/cosmos/interchain-security/v5/x/ccv/consumer/keeper"
  ccvconsumertypes "github.com/cosmos/interchain-security/v5/x/ccv/consumer/types"
  ccvtypes "github.com/cosmos/interchain-security/v5/x/ccv/types"
)
```

- **Extend the `App` struct**:

```go
type App struct {
  // ...
  StakingKeeper    *stakingkeeper.Keeper
  // ...
  CCVConsumerKeeper ccvconsumerkeeper.Keeper
}
```

- **Store key**:

```go
keys := storetypes.NewKVStoreKeys(
  // ...
  ccvconsumertypes.StoreKey,
)
```

- **Scope capabilities and create the consumer keeper**:

```go
scopedCCVConsumerKeeper := app.CapabilityKeeper.ScopeToModule(ccvconsumertypes.ModuleName)

validatorAddrCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix())
consAddrCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix())

app.CCVConsumerKeeper = ccvconsumerkeeper.NewKeeper(
  appCodec,
  keys[ccvconsumertypes.StoreKey],
  app.getSubspace(ccvconsumertypes.ModuleName),
  scopedCCVConsumerKeeper,
  app.IBCKeeper.ChannelKeeper,
  app.IBCKeeper.PortKeeper,
  app.IBCKeeper.ConnectionKeeper,
  app.IBCKeeper.ClientKeeper,
  app.SlashingKeeper,
  app.BankKeeper,
  app.AccountKeeper,
  app.TransferKeeper,
  app.IBCKeeper,
  authtypes.FeeCollectorName,
  lib.GovModuleAddress.String(),
  validatorAddrCodec,
  consAddrCodec,
)
```

### 1.3. IBC Routing and Module Manager Integration

- **IBC router**: register the consumer IBC module on the **`consumer`** port:

```go
ccvConsumerModule := ccvconsumermodule.NewAppModule(
  app.CCVConsumerKeeper,
  app.getSubspace(ccvconsumertypes.ModuleName),
)
var ccvConsumerIBCModule ibcporttypes.IBCModule = ccvConsumerModule

ibcRouter := ibcporttypes.NewRouter()
ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferStack)
ibcRouter.AddRoute(ccvtypes.ConsumerPortID, ccvConsumerIBCModule) // port "consumer"
ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostIBCModule)

app.IBCKeeper.SetRouter(ibcRouter)
```

- **Module manager**: register the consumer module and set begin/end/init order:

```go
app.ModuleManager = module.NewManager(
  // ...
  ibc.NewAppModule(app.IBCKeeper),
  ccvconsumermodule.NewAppModule(
    app.CCVConsumerKeeper,
    app.getSubspace(ccvconsumertypes.ModuleName),
  ),
  // ...
)

app.ModuleManager.SetOrderBeginBlockers(
  // ...
  slashingtypes.ModuleName,
  ccvconsumertypes.ModuleName,
  evidencetypes.ModuleName,
  // ...
)

app.ModuleManager.SetOrderEndBlockers(
  // ...
  ibcexported.ModuleName,
  ccvconsumertypes.ModuleName,
  ibctransfertypes.ModuleName,
  // ...
)

app.ModuleManager.SetOrderInitGenesis(
  // ...
  ibcexported.ModuleName,
  ccvconsumertypes.ModuleName,
  genutiltypes.ModuleName,
  // ...
)
```

- **Params subspace**:

```go
paramsKeeper.Subspace(ccvconsumertypes.ModuleName)
```

These changes ensure the CCV consumer module participates correctly in IBC routing and in the ABCI lifecycle.

### 1.4. Add a CCV Consumer Upgrade Handler (`v8.0.0`)

Create `protocol/app/upgrades/v8.0.0/constants.go`:

```go
package v_8_0_0

import (
  store "cosmossdk.io/store/types"
  "github.com/nemo-network/v4-chain/protocol/app/upgrades"
  ccvconsumertypes "github.com/cosmos/interchain-security/v5/x/ccv/consumer/types"
)

const UpgradeName = "v8.0.0"

var Upgrade = upgrades.Upgrade{
  UpgradeName: UpgradeName,
  StoreUpgrades: store.StoreUpgrades{
    Added: []string{
      ccvconsumertypes.StoreKey,
    },
  },
}
```

Create `protocol/app/upgrades/v8.0.0/upgrade.go` with:

- **Step 1**: run module migrations.
- **Step 2**: construct a CCV consumer genesis with:
  - `PreCCV = true` – enabling standalone → consumer changeover mode.
  - `Params = ccv.DefaultParams()` – baseline consumer parameters.
  - `Provider.InitialValSet` – set to the current Nemo staking validator set.

Sketch:

```go
func CreateUpgradeHandlerWithConsumerInit(
  mm module.Manager,
  configurator module.Configurator,
  consumerKeeper ccvconsumerkeeper.Keeper,
  stakingKeeper *stakingkeeper.Keeper,
) upgradetypes.UpgradeHandler {
  return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
    sdkCtx := lib.UnwrapSDKContext(ctx, "app/upgrades/v8.0.0")

    // 1. Run migrations
    vm, err := mm.RunMigrations(ctx, configurator, vm)
    if err != nil {
      return vm, err
    }

    // 2. Collect current validator set
    initialValSet, err := getInitialValidatorSet(ctx, stakingKeeper)
    if err != nil {
      return vm, err
    }

    // 3. Build consumer genesis
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
```

Helper to build the initial validator set:

```go
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
```

Register the upgrade in `protocol/app/upgrades.go`:

```go
import (
  // ...
  v7_0_0 "github.com/nemo-network/v4-chain/protocol/app/upgrades/v7.0.0"
  v8_0_0 "github.com/nemo-network/v4-chain/protocol/app/upgrades/v8.0.0"
)

var Upgrades = []upgrades.Upgrade{
  v7_0_0.Upgrade,
  v8_0_0.Upgrade,
}

func (app *App) setupUpgradeHandlers() {
  // v7
  app.UpgradeKeeper.SetUpgradeHandler(
    v7_0_0.UpgradeName,
    v7_0_0.CreateUpgradeHandler(app.ModuleManager, app.configurator),
  )

  // v8 – CCV consumer changeover
  app.UpgradeKeeper.SetUpgradeHandler(
    v8_0_0.UpgradeName,
    v8_0_0.CreateUpgradeHandlerWithConsumerInit(
      app.ModuleManager,
      app.configurator,
      app.CCVConsumerKeeper,
      app.StakingKeeper,
    ),
  )
}
```

Finally, build the new Nemo binary:

```bash
cd /root/dapp/nemo-app-chain/protocol
go build -o build/nemod ./cmd/nemod
```

---

## 2. Running the Nemo Upgrade via Governance

This section covers how to use the governance module to perform the **Pre-CCV software upgrade** on Nemo.

### 2.1. Export the Pre-CCV Genesis (Optional but Recommended)

Before running the upgrade, export the current (standalone) genesis for backup and inspection:

```bash
nemod export > genesis_pre_ccv.json
```

This file has **no** `ccvconsumer` section yet.

### 2.2. Submit a Software Upgrade Proposal on Nemo

Choose:

- `UPGRADE_NAME = v8.0.0`
- `UPGRADE_HEIGHT` = a block height comfortably in the future.

Submit the proposal:

```bash
UPGRADE_NAME=v8.0.0
UPGRADE_HEIGHT=2000     # example; ensure > current height

nemod tx gov submit-proposal software-upgrade "$UPGRADE_NAME" \
  --title "Nemo: Enable ICS Consumer (Pre-CCV)" \
  --description "Upgrade to v8.0.0 to add CCV consumer module in Pre-CCV mode." \
  --upgrade-height "$UPGRADE_HEIGHT" \
  --from validator \
  --deposit 1000000unemo \
  --chain-id nemo-local-1 \
  --gas auto --gas-adjustment 1.3 -y
```

Vote:

```bash
nemod tx gov vote <proposal-id> yes \
  --from validator \
  --chain-id nemo-local-1 \
  --gas auto --gas-adjustment 1.3 -y
```

Monitor:

```bash
nemod q gov proposal <proposal-id>
```

### 2.3. Switch Binaries at the Upgrade Height

Watch height:

```bash
nemod status | jq '.SyncInfo.latest_block_height'
```

When the chain reaches `UPGRADE_HEIGHT`, the old binary will halt with an upgrade error. At that point:

```bash
sudo systemctl stop nemod    # or pkill nemod

cp /root/dapp/nemo-app-chain/protocol/build/nemod /usr/local/bin/nemod

nemod start --home ~/.nemo
```

On startup:

- `UpgradeKeeper` applies the `v8.0.0` **store upgrades** (adds `ccvconsumer` store).
- `CreateUpgradeHandlerWithConsumerInit`:
  - Runs module migrations.
  - Sets `PreCCV = true`.
  - Sets `Provider.InitialValSet` from the current staking validator set.

Nemo is now in **Pre-CCV** mode: the consumer module is initialized, but validator set updates still come from local staking until CCV connection to Gaia is established.

### 2.4. Export the Post-CCV Genesis Snapshot (Optional)

After the upgrade:

```bash
nemod export > genesis_post_ccv.json
```

Inspect the new `ccvconsumer` section:

```bash
jq '.app_state.ccvconsumer' genesis_post_ccv.json
```

You should see fields like:

- `preCCV: true`
- `params` – consumer parameters
- `provider.initial_val_set` – initial validator set used when CCV fully activates.

You can also query on-chain:

```bash
nemod q ccvconsumer params --chain-id nemo-local-1 -o json
nemod q ccvconsumer provider --chain-id nemo-local-1 -o json
```

---

## 3. Setting Up Hermes Between Gaia and Nemo

This section assumes both `gaiad` (provider) and `nemod` (consumer) nodes are already running.

### 3.1. Hermes Configuration

Create or edit `~/.hermes/config.toml`:

```toml
[global]
log_level = "info"
rpc_timeout = "10s"
clear_on_start = true

[[chains]]
id = "gaia-local-1"
rpc_addr = "http://127.0.0.1:26657"
grpc_addr = "http://127.0.0.1:9090"
websocket_addr = "ws://127.0.0.1:26657/websocket"
rpc_timeout = "10s"
account_prefix = "cosmos"
key_name = "relayer"
store_prefix = "ibc"
max_gas = 2000000
gas_price = { price = 0.01, denom = "uatom" }
clock_drift = "5s"
trusting_period = "336h"
memo_prefix = "hermes-gaia-nemo"

[[chains]]
id = "nemo-local-1"
rpc_addr = "http://127.0.0.1:36657"
grpc_addr = "http://127.0.0.1:9091"
websocket_addr = "ws://127.0.0.1:36657/websocket"
rpc_timeout = "10s"
account_prefix = "nemo"
key_name = "relayer"
store_prefix = "ibc"
max_gas = 2000000
gas_price = { price = 0.01, denom = "unemo" }
clock_drift = "5s"
trusting_period = "336h"
memo_prefix = "hermes-gaia-nemo"
```

Adjust ports, prefixes, and denoms to match your local setup.

### 3.2. Add and Fund Relayer Keys

Add keys:

```bash
hermes keys add --chain gaia-local-1 --mnemonic-file gaia-relayer.mnemonic
hermes keys add --chain nemo-local-1 --mnemonic-file nemo-relayer.mnemonic
```

Fund the relayer addresses on both chains (from your validator or faucet accounts), then validate:

```bash
hermes keys list --chain gaia-local-1
hermes keys list --chain nemo-local-1
```

### 3.3. Run Hermes

Start the relayer:

```bash
hermes start
```

For CCV, you usually do **not** manually create transfer channels. The **provider module** on Gaia will:

- Create the **consumer client / connection**.
- Initiate the **CCV channel handshake** on the special CCV ports.

Hermes relays all IBC handshakes and subsequent CCV packets.

You can inspect channels:

```bash
hermes query channels --chain gaia-local-1
hermes query channels --chain nemo-local-1
```

Later, you should see:

- A channel on Gaia associated with the provider port.
- A channel on Nemo with `port_id = "consumer"`.

---

## 4. Creating the Consumer Chain on Gaia (MsgCreateConsumer)

On the Gaia (provider) side, you use `MsgCreateConsumer` (wrapped in a governance proposal) to register Nemo as a consumer chain.

### 4.1. Understanding `spawn_time` and Other Parameters

Key fields for `MsgCreateConsumer`:

- **`chain_id`**: the consumer chain ID (`"nemo-local-1"`).
- **`spawn_time`**:
  - RFC3339 UTC timestamp.
  - The time when the provider completes setting up the consumer chain.
  - Must be strictly in the **future** relative to current Gaia time.
- **`initial_height`**:
  - Height at which the consumer chain starts under provider security.
  - For a brand-new chain, `{revision_number: 0, revision_height: 1}`.
  - For a changeover chain, this can be aligned with the height after the Pre-CCV upgrade.
- **`unbonding_period`**, **`ccv_timeout_period`**, **`transfer_timeout_period`**:
  - Go duration strings (e.g. `"1209600s"`).
  - Control unbonding length and CCV/transfer packet timeouts.
- **`consumer_redistribution_fraction`**:
  - Fraction of staking rewards routed back to the provider or to a dedicated redistribution address.
- **`blocks_per_distribution_transmission`**, **`historical_entries`**, etc.:
  - Fine-grained CCV economics and history parameters.

In **local testing**, you can choose a `spawn_time` 5–10 minutes in the future:

```bash
GAIA_TIME=$(gaiad status | jq -r '.SyncInfo.latest_block_time')
echo "$GAIA_TIME"

# Example: if GAIA_TIME = 2026-03-03T10:00:00Z,
# choose spawn_time = 2026-03-03T10:10:00Z manually.

SPAWN_TIME="2026-03-03T10:10:00Z"
```

### 4.2. Create the Consumer Proposal JSON

Example `consumer-proposal.json`:

```json
{
  "title": "Create Nemo Consumer Chain",
  "description": "Make nemo-local-1 a partial-set security opt-in consumer of gaia-local-1",
  "chain_id": "nemo-local-1",
  "spawn_time": "2026-03-03T10:10:00Z",
  "initial_height": {
    "revision_number": "0",
    "revision_height": "1"
  },
  "genesis_hash": "",
  "binary_hash": "",
  "blocks_per_distribution_transmission": "1000",
  "historical_entries": "10000",
  "unbonding_period": "1209600s",
  "ccv_timeout_period": "2419200s",
  "transfer_timeout_period": "3600s",
  "consumer_redistribution_fraction": "0.75",
  "provider_fee_pool_addr": "",
  "allowlist": [],
  "denylist": [],
  "distribution_transmission_channel": "",
  "top_n": "0",
  "validators_power_cap": "0",
  "exponentials": []
}
```

Notes:

- In production, `genesis_hash` and `binary_hash` should be filled with SHA256 hashes of the consumer genesis and binary.
- For local testing, leaving them empty is acceptable.

### 4.3. Submit `MsgCreateConsumer` via Governance on Gaia

Submit proposal:

```bash
gaiad tx gov submit-proposal create-consumer consumer-proposal.json \
  --from gaia-validator \
  --chain-id gaia-local-1 \
  --deposit 1000000uatom \
  --gas auto --gas-adjustment 1.3 -y
```

Vote:

```bash
gaiad tx gov vote <proposal-id> yes \
  --from gaia-validator \
  --chain-id gaia-local-1 \
  --gas auto --gas-adjustment 1.3 -y
```

Monitor:

```bash
gaiad q gov proposal <proposal-id> -o json
gaiad q provider consumer-chains
```

Once the proposal passes and `spawn_time` is reached:

- Gaia registers `nemo-local-1` as a consumer chain.
- The provider module begins establishing the CCV client/connection/channel toward Nemo.

Hermes relays all IBC handshakes and CCV packets.

### 4.4. Consumer Genesis on the Provider Side

After `spawn_time`, you can query the **consumer genesis data** generated by Gaia:

```bash
gaiad q provider consumer-genesis nemo-local-1 -o json > ccvconsumer_genesis.json
```

This file represents the CCV-related genesis state from the provider’s perspective. For a changeover chain like Nemo:

- It is combined with Nemo’s own app-state (which already has `ccvconsumer` initialized in Pre-CCV mode).
- In practice, for changeover, the main source of truth is Nemo’s upgraded chain state; the provider’s consumer genesis query is more critical when bootstrapping a new standalone consumer chain.

---

## 5. Verifying Nemo as an Opt-In Consumer of Gaia

### 5.1. On Gaia (Provider)

Check consumer chain registration:

```bash
gaiad q provider consumer-chains -o json
gaiad q provider consumer-chain nemo-local-1 -o json
```

You should see:

- `chain_id: "nemo-local-1"`
- Associated `client_id`, `channel_id`, and consumer status.

### 5.2. On Nemo (Consumer)

Check CCV consumer params and provider info:

```bash
nemod q ccvconsumer params -o json
nemod q ccvconsumer provider -o json
```

Expected:

- `params.enabled` – true once the CCV channel is fully established and consumer is active.
- `provider_client_id` and `provider_channel_id` – set to the values corresponding to the Gaia connection/channel.

Inspect IBC channels on Nemo:

```bash
nemod q ibc channel channels -o json \
  | jq '.channels[] | select(.port_id == "consumer")'
```

You should see:

- A channel bound to the `consumer` port (`ccvtypes.ConsumerPortID`).

### 5.3. Validator-Set Alignment and Slashing Path

To verify Nemo is secured by Gaia’s validator set:

```bash
# Gaia validator set
gaiad q staking validators -o json > gaia_validators.json

# Nemo (consumer) effective validator set (command may vary with ICS version)
nemod q ccvconsumer validator-set -o json > nemo_ccv_validators.json
```

The active validator set on Nemo should match the provider’s validator set modulo PSS/top-N logic.

For **slashing and alerts**:

- Downtime / double-sign evidence on Nemo is relayed back to Gaia over the CCV channel.
- Gaia’s provider module applies the appropriate slashing and jailing.

You can confirm on Gaia:

```bash
gaiad q slashing signing-info <val-cons-address>
gaiad q staking validator <val-op-address>
```

---

## 6. Summary

High-level steps to make `nemo-app-chain` an Opt-In consumer of Cosmos Hub locally:

- **Code**: add `interchain-security/v5`, wire CCV consumer keeper and module, register IBC routes, and implement a `v8.0.0` Pre-CCV upgrade handler.
- **Nemo upgrade**: use governance to schedule `v8.0.0`, switch to the new binary at the upgrade height, and let the handler initialize the CCV consumer state (Pre-CCV).
- **Relayer**: configure Hermes between `gaia-local-1` and `nemo-local-1`, ensure it is running and relaying.
- **Provider**: submit `MsgCreateConsumer` (via gov) on Gaia with a future `spawn_time`, let the provider create the CCV client/connection/channel.
- **Verification**: check CCV consumer params and channels on Nemo, check consumer chains on Gaia, and confirm that Nemo’s validator set matches Gaia’s (subject to PSS settings).

