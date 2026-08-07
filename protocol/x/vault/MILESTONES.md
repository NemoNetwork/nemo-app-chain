# Vault Module — Delivery Milestones

**Companion to:** [spec.md](./spec.md) (status & estimation, 2026-07-24)
**Scope:** `protocol/x/vault` → indexer, end to end.

**Scoping decisions:**

1. **A real performance-fee engine is required at launch** — not fields-only at 0%.
   This is spec.md's Scenario B. M2 gains a mechanism-design phase and grows to 3–4 weeks.
2. **Abacus is deferred to a later phase.** It lives in a separate repository; it will be
   scoped once the M5 API contract is settled. M6 is out of the current plan.
3. **The chain has live state that must survive.** An upgrade handler is mandatory, and
   this fork has never shipped one. This is now a first-class milestone, not a buffer item.

Decisions 1 and 3 interact, and that interaction is the single largest design risk in this
plan — see [Risk: charging fees to existing depositors](#risk-charging-fees-to-existing-depositors).

---

## 0. Ground truth: this fork's vault module has an exact upstream ancestor

The single most important planning fact, established by diffing this repo against a
blobless mirror of `dydxprotocol/v4-chain`:

> **`protocol/x/vault` in this repo is byte-identical to upstream `dydxprotocol/v4-chain`
> at commit `e0a20164` (2024-08-28, *"implement megavault shares locking" #2160*),
> modulo the `dydxprotocol` → `nemo_network` proto package rename.**

Verification: after normalizing import paths, every hand-written `.go` file in
`x/vault` diffs to **zero lines** against `e0a20164` except three trivial spots:

| File | Delta vs `e0a20164` |
|---|---|
| `genesis.go` | 3 lines — `genesis.Vaults = k.GetAllVaults(ctx)` is duplicated in `ExportGenesis` (copy-paste bug, harmless but should be cleaned) |
| `keeper/orders.go:47` | `vault.PerpetualPositions == nil \|\| len(...) == 0` → `len(...) == 0` (lint cleanup) |
| `keeper/msg_server_deposit_to_megavault_test.go` | 1 line |
| `types/*.pb.go` | rename-only (gzipped file descriptors + registration paths) |

Two upstream test files were dropped in the fork: `keeper/shares_test.go`, `types/params_test.go`.

**Consequence:** Workstreams 1 and 2 of [spec.md](./spec.md) are not a design-and-build
exercise — they are a **replay of a known, bounded upstream commit range**:

- **41 commits** touch `x/vault` / `proto/…/vault` between `e0a20164` and `protocol/v9.6.3`,
  ~90% of them clustered in **September 2024** (the megavault withdrawal/operator push).
- **22 commits** touch the upstream indexer vault stack over the same range.
- The vault module is essentially **frozen upstream after v7.0.0** — `v7.0.0 → v9.6.3` is
  only ~540 lines across 14 files, mostly generated code and one quoting-param addition.

This lowers risk substantially and changes the *shape* of the work from "implement" to
"port, rename, re-test". It does **not** shrink the indexer API work, which is genuinely new.

### What upstream gives us for free

| Feature | Upstream file(s) | Status here |
|---|---|---|
| `MsgWithdrawFromMegavault` + redemption/slippage math | `keeper/withdraw.go` (426 ln) | missing |
| Withdrawal preview query | `keeper/grpc_query_withdrawal_info.go` | missing |
| `MsgUnlockShares` + unlock processing (closes `TODO (TRA-565)`) | `keeper/msg_server_unlock_shares.go` | missing |
| Operator params + `MsgUpdateOperatorParams` | `keeper/msg_server_update_operator_params.go` | missing |
| `MsgAllocateToVault` / `MsgRetrieveFromVault` | `keeper/msg_server_{allocate,retrieve}*.go` | missing |
| Bank→subaccount fund sweep | `keeper/sweep_funds.go` | missing |
| Indexer vault events (`UpsertVault`, transfers) | `types/events.go` | missing |
| `vaults` table, hourly/daily PnL materialized views | `postgres/…/vault-table.ts`, `vault-pnl-ticks-view.ts` | missing |
| Ender `upsert-vault` validator + handler + SQL | `ender/…/upsert-vault-*`, `dydx_vault_upsert_handler.sql` | missing |
| Roundtable `refresh-vault-pnl` | `roundtable/…/refresh-vault-pnl.ts` | missing |
| Vault PnL/positions endpoints w/ resolution + caching | `comlink/…/vault-controller.ts` (895 ln) | 3 stub endpoints, `EXPERIMENT_VAULTS` config hack |

### What has NO upstream equivalent (must be designed here)

Upstream `OperatorParams` is **only** `{ operator address, {name, description} }`. There is:

- **no performance fee, no profit share, no minimum operator share** anywhere upstream
- **no max-deposit cap**
- **no** summary/metrics/APR/MDD/NAV/charts/my-activity endpoints — upstream's vault
  controller exposes exactly three routes (`megavault/historicalPnl`,
  `vaults/historicalPnl`, `megavault/positions`)

Everything the NLP page needs beyond PnL + positions is net-new indexer work.

---

## Progress

| Milestone | Status |
|---|---|
| M0 — baseline, port harness | ✅ done |
| M1a — keeper plumbing + indexer events | ✅ done (indexer track unblocked) |
| M1b — share unlocking | ✅ done (closes `TODO (TRA-565)`) |
| M1c — megavault withdrawal | ✅ done |
| M2a — operator params (pulled forward, see below) | ✅ done |
| M1d — allocate / retrieve, quoting modes | ✅ done — **M1 complete** |
| M3a — indexer TS proto codegen | ✅ done |
| M3b — `vaults` table (migration, model, store, types) | ✅ done |
| M3c — ender upsert-vault path (+ restored ender test harness) | ✅ done |
| M3d — vault PnL views, roundtable refresh, Athena DDL, controller off env vars | ✅ done — **M3 complete** |
| M2a remainder — max-deposit cap (net-new) | ✅ done — **M2a complete** |
| M2b — fee mechanism design | ✅ done — [ADR-001](./spec/adr-001-megavault-fees.md) |
| M2c — fee engine implementation | ✅ done — **M2 complete** |
| M7 — upgrade handler | ✅ done (state migration + backfill; rehearsal on a live export still owed) |
| M4 — port the vault controller | ✅ done — **M4 complete** |
| M5 — NLP page API surface | next — **blocked on `NewSpec.md`** |

**Verification status:** the DB-backed indexer tests **have now been run**. There is no
Postgres on this machine and no root to install one, so a userspace PostgreSQL 18 was stood
up from the `embedded-postgres` npm package (binaries only, nothing installed system-wide).
Against it:

- `packages/postgres`: **362 tests pass, 35 suites**, including the new `vault-table` and
  `vault-pnl-ticks-view` suites — so the migration, the store, and both materialized views
  are verified against a real database, not just typechecked.
- `services/ender`: both vault suites pass, exercising the real SQL handler, the status-enum
  helper, the block-processor dispatch arm and `dydx_to_jsonb`.

Running them surfaced **two pre-existing bugs that block every fresh environment** — see
[Two blocking bugs found by actually running the tests](#two-blocking-bugs-found-by-actually-running-the-tests).

Still unverified: `services/comlink` (the fork deleted its 27 test files, so there is nothing
to run), Redis- and Kafka-backed suites, and the M7 upgrade rehearsal.

### M1b notes

Ported from upstream `6d86d0b9` (#2169) — the commit immediately after this fork's
ancestor, so it applied almost verbatim. `MsgUnlockShares` + `UnlockShares` keeper method;
`LockShares` now schedules the unlock via `delayMsgKeeper.DelayMessageByBlocks`.

Two pre-existing bugs surfaced while registering the new message:

1. **20 message type URLs were wrong.** `app/msgs/internal_msgs.go` and `normal_msgs.go`
   keyed entries as `/nemo-network.<module>.<Msg>` (hyphen) while the real proto package is
   `nemo_network` — verified with `sdk.MsgTypeURL`. Affected affiliates, listing, revshare
   and vault, including `MsgDepositToMegavault`. These maps are consumed **only by tests**,
   so this was not rejecting live transactions — but the tests that enforce the catalogue
   against the real registry were failing because of it.
2. **The expected-message lists were never re-sorted after the rename.** They are compared
   with `require.Equal` against `lib.GetSortedKeys(...)`, so order is significant, and
   `dydxprotocol` sorted before `/ibc.` where `nemo_network` sorts after it. Fixed by
   regrouping each list into sorted order, preserving the per-module comments — including
   the bare placeholder comments (`// consensus`, `// perpetuals`, `// prices`) that mark
   modules with no messages, which get a synthetic sort key from the following group's
   namespace so they stay in their alphabetical slot.

**And one genuine production bug:** `lib/ante/unsupported_msgs.go` listed
`*crisis.MsgVerifyInvariant` as an unsupported message, so `IsUnsupportedMsg` returned true
for it — and unlike the registries above, that function *is* on the live ante path
(`app/ante/msg_type.go`, `lib/ante/disallow_msg.go`, `lib/ante/nested_msg.go`). Meanwhile
`app/msgs/normal_msgs.go` explicitly allows `/cosmos.crisis.v1beta1.MsgVerifyInvariant` and
`UnsupportedMsgSamples` does not list it. Net effect: crisis invariant-verification
transactions would be rejected on chain despite being allow-listed.

Upstream has two deprecated vault messages in that `case`
(`MsgSetVaultQuotingParams`, `MsgUpdateParams`); neither exists in this fork, and it looks
like `crisis.MsgVerifyInvariant` was substituted to keep the `case` syntactically valid —
the orphaned `// vault` comment above it is the tell. Removed, leaving ICA + gov.

### M1c notes — megavault withdrawal

`withdraw.go` (slippage, `RedeemFromMainAndSubVaults`, `WithdrawFromMegavault`),
`MsgWithdrawFromMegavault`, `QueryMegavaultWithdrawalInfo`, `sweep_funds.go` wired into the
EndBlocker, the withdrawal CLI command, and the deposit/withdraw/sweep typed events.

The final `withdraw.go` needed supporting pieces first, ported in this order:

- `GetVaultQuotingParams` → `GetVaultAndQuotingParams` (also returns vault params); 3 call
  sites in `orders.go`
- `GetVaultLeverageAndEquity` and `GetVaultClobPerpAndMarket` in `vault.go`
- deactivated vaults excluded from megavault equity (upstream #2308)
- `lib/vault` helper package (`SkewAntiderivative`, `SpreadPpm`) — `orders.go` had the
  spread formula inline; upstream extracted it, so this de-duplicates rather than adds
- `lib.BigRatMin`
- `PerpetualsKeeper.GetLiquidityTier`

**`GetLiquidityTier` needed an interface change.** The vault expects it, but
`x/perpetuals/types.PerpetualsKeeper` didn't declare it, so `mocks.PerpetualsKeeper` didn't
have it either. Upstream has the same gap — its *committed mock* has the method while its
interface does not, i.e. upstream's mock is out of sync with its own interface. Rather than
copy that inconsistency, the method was added to the perpetuals interface (the concrete
keeper already implements it) and the mock regenerated.

**Error codes.** Upstream repurposed vault error code 4 from `ErrInvalidDepositAmount` to
`ErrInvalidQuoteQuantums`, broadening it to cover deposits, allocations, retrievals and
withdrawals. Renamed to match — same code, so no error-code break. New errors were given
upstream's exact codes (25, 27, 28, 30) with 24, 26, 29 and 31 left as documented gaps for
the operator-params and vault-status work still to come, so those ports stay mechanical.

### M2a pulled ahead of M1d — operator params first

M1d's `MsgAllocateToVault` is gated on *"module authority **or operator**"*, so it needs
`GetOperatorParams`. Upstream's own ordering confirms this (#2259 operator params landed
before #2277 allocate), so the operator-params half of M2a was ported first rather than
writing throwaway gov-only authority logic and revisiting it.

Landed: `OperatorParams` / `OperatorMetadata` protos, `MsgUpdateOperatorParams`, keeper
accessors, genesis round-trip, operator params in the params query, and error codes 24
(`ErrEmptyOperator`) and 26 (`ErrInvalidAuthority`) filling two of the gaps reserved in M1c.
`MsgSetVaultParams` is now operator-gated too, matching upstream.

**One proto divergence, deliberate.** Upstream dropped the deprecated `params` field from
`QueryParamsResponse` and renumbered, making `operator_params` field 2. On a chain with live
clients that renumbering is wire-breaking, so `operator_params` was appended at field 3 and
the deprecated field left in place. Noted in the proto so the divergence is not mistaken for
an oversight during the next port.

### M1d notes — allocate / retrieve

`MsgAllocateToVault`, `MsgRetrieveFromVault`, and the `AllocateToVault` keeper method.

**Message classification was the trap.** These were initially registered as *internal*
messages alongside the other authority-gated ones, and every test failed with
*"internal msg cannot be submitted externally"*. Internal means gov-proposal-only, but the
operator is not gov — it submits these as ordinary transactions. Upstream puts
`MsgAllocateToVault`, `MsgRetrieveFromVault` and `MsgSetVaultParams` in **normal** msgs and
only `MsgUpdateOperatorParams` in internal. Corrected to match.

The rule of thumb: a vault message gated on *"authority or operator"* belongs in
`NormalMsgs`; one gated on module authority alone belongs in `InternalMsgs` and
`lib/ante/internal_msg.go`.

**Quoting modes** (#2299, #2310, #2329, #2301, plus the skew-factor check #2249) landed in
the same milestone. `orders.go` was replaced wholesale with upstream's version — by then the
only differences were forward changes plus one fork-local lint cleanup, which was re-applied
on top. That brought close-only reduce-only sizing, no-quoting for deactivated/stand-by,
post-only orders on inventory-increasing sides, `CancelVaultClobOrder` /
`TryToCancelVaultClobOrder`, and cancellation of no-longer-needed orders.
`SetVaultParams` now refuses to deactivate a vault with positive equity and cancels
outstanding orders when moving a vault to deactivated or stand-by. Error codes 29 and 31 are
now used, so the vault error range 1–31 matches upstream exactly.

**A genesis-ordering consequence worth knowing:** `InitGenesis` calls `SetVaultParams`, so
the new "cannot deactivate a vault with positive equity" rule applies at genesis too — a
genesis state with a deactivated vault that has a funded subaccount will panic on
`InitChain`. Upstream has the identical constraint; its tests were updated in the same
commit, which is why `vault_test.go`, `orders_test.go` and `params_test.go` were taken from
upstream rather than patched locally.

The ported `params_test.go` is strictly better coverage than the hand-written version it
replaced (it exercises deactivation-blocking and order cancellation), and it injects a test
message sender via `WithAppOptions` so its indexer-event assertions are real rather than
silently passing against a nil block. The one assertion carried over from the local version:
a *rejected* `SetVaultParams` must emit no `upsert_vault` event.

**One non-obvious wiring step:** `MsgWithdrawFromMegavault` sets
`(cosmos.msg.v1.signer) = "subaccount_id"`, and because the signer lives in a nested
message the SDK cannot resolve it from the annotation alone. It needs an entry in
`CustomGetSigners` in [app/module/interface_registry.go](../../app/module/interface_registry.go)
alongside the existing `MsgDepositToMegavault` one. Without it every withdrawal fails
`CheckTx` with *"no cosmos.msg.v1.signer option found for message
nemo_network.subaccounts.SubaccountId"*. Any future vault message that signs via a
subaccount needs the same entry.

### M2a notes — the max-deposit cap (net-new)

No upstream equivalent. Shipped as a new `MegavaultParams` message holding the parameters
that apply to megavault as a whole rather than to one vault or to the operator's identity,
plus `MsgUpdateMegavaultParams`, keeper accessors, genesis round-trip, and the cap check at
the top of `DepositToMegavault`.

**Two decisions worth knowing:**

1. **The cap is on megavault equity after the deposit, not on cumulative net deposits.**
   Cumulative-deposit accounting would need new state that tracks deposits minus withdrawals
   forever; an equity cap reuses `GetMegavaultEquity` and needs no state at all. The
   consequence is ordinary TVL-cap semantics: trading gains can fill the cap on their own,
   and losses re-open it.
2. **A cap of zero means uncapped.** New state reads as zero, so a chain upgraded from a
   state that predates the field keeps accepting deposits rather than silently halting them.

`MsgUpdateMegavaultParams` is **operator-gated** (`NormalMsgs`), unlike the fee parameters
which are gov-only. The dividing line is stated in [ADR-001 D7](./spec/adr-001-megavault-fees.md):
risk knobs the operator must be able to pull quickly are operator-gated; anything that moves
value from depositors to the operator is gov-gated.

Fork-local error codes start at **32**, since upstream's vault range ends at 31. Future
upstream ports therefore stay mechanical.

### M2b notes — the fee design

[ADR-001](./spec/adr-001-megavault-fees.md). The open question this plan recorded — *"fee vs.
profit share: one charge or two?"* — is **answered by the evidence in the repo**: the
frontend renders three distinct zero-valued fields (*operator fee*, *minimum operator share*,
*profit share*), and a UI does not show one charge twice. So: two charges, which is also the
conventional management-fee-plus-performance-fee structure those names imply.

The decision is cheap to reverse — the two charges are independent parameters computed in
separate steps, so if product says they are one charge, setting either to zero collapses the
design with no code change. That is why this was resolved by design rather than by waiting.

### M2c notes — the fee engine

`keeper/fees.go`: `AccrueFees`, `MaybeAccrueFees` (the EndBlocker interval entry point),
`GetNavPerShare`, `updateHighWaterMark`, `mintSharesTo`, `ValidateOperatorShareFloor`. Fee
params appended to `OperatorParams` at fields 3–5; new `FeeState` under its own store key;
`MegavaultFeeState` query; `accrue_megavault_fees` event; genesis round-trip.

**Three things that turned out to matter more than the arithmetic:**

1. **The sentinel problem.** `last_accrual_time == 0` reads as "never accrued", but zero is
   also a real Unix timestamp — and the test app's genesis block time is 1 ns after the
   epoch, so `Unix()` is exactly 0. The first implementation guarded on
   `LastAccrualTime > 0` and silently never charged an operator fee. The fix is not a bigger
   guard: it is that the clock must **always** be initialized (`InitGenesis` seeds it, and so
   does the M7 upgrade handler), plus a `MaxAccrualElapsedSeconds` clamp of one year so that
   a missed initialization costs at most a year of fees instead of five decades.
2. **Rounding to zero must not advance the clock.** If it did, a fee rate too small to be
   representable over one interval would accrue *nothing ever*, however much time passed.
   Leaving the timestamp unadvanced lets the elapsed time grow until the fee is
   representable; the one-year clamp bounds how long that can take.
3. **Accrue before deposits, not just before withdrawals.** Accruing before a withdrawal is
   the obvious half — it is what makes the exiting depositor pay their share. Accruing before
   a *deposit* is the less obvious half and is a correctness requirement: without it, a
   deposit landing just before a large accrual dilutes the new depositor for fees earned
   before they arrived.

Dilution is what makes the whole thing work: because minting shares to the operator reduces
every existing share equally, "force an accrual immediately before a withdrawal" *is*
per-depositor settlement, with no per-owner lot tracking anywhere in state.

### M7 notes — the upgrade handler

[`app/upgrades/v1.0`](../../app/upgrades/v1.0/) — the first upgrade handler this fork ships.
`Upgrades` in [app/upgrades.go](../../app/upgrades.go) is no longer an empty slice.

Three steps, all idempotent and all conditional, so the upgrade is safe to re-run against
state that already carries the values (a genesis export/import round-trip, for instance):

1. **Operator params** default to the gov module account if unset. Fee params are left at
   zero — the engine ships inert and is turned on by a separate governance action, which is
   mitigation 3 of ADR-001 §7 and buys an announcement window.
2. **Fee state**: the high-water mark is seeded at the NAV per share observed *at upgrade
   height*, and the accrual clock at the upgrade block time. This is the single most
   consequential line in the upgrade — **a high-water mark of zero would charge the
   operator's profit share against the vault's entire trading history, retroactively, on
   depositors who deposited under a 0%-fee regime.** An existing mark is only ever raised,
   never lowered.
3. **Indexer backfill**: every existing vault has its params re-set, which is a no-op on
   chain state but emits the `UpsertVault` event each vault has never emitted. Without it the
   indexer's `vaults` table starts empty and every vault endpoint reports nothing.

**No `StoreUpgrades`.** All new state (`OperatorParams`, `MegavaultParams`, `FeeState`) lives
under the existing `vault` store key as new prefixes — confirmed, not assumed.

**Two of upstream's three v7 vault migrations were correctly skipped**:
`migrateVaultSharesToMegavaultShares` and `migrateVaultQuotingParamsToVaultParams` are both
already applied in this fork.

`App.GetKey(storeKey)` was added — a standard accessor most Cosmos apps expose, needed to
model genuinely-unset pre-upgrade state in the test rather than approximating it.

**Still owed for M7:** a rehearsal against a real state export from the live chain, and an
upgrade-container test. Neither can be done on this machine (no Docker, no state export).

### Spec-compliance audit (2026-08-07): closing every unintended delta vs upstream

Prompted by the requirement that the vault module be *complete against the specification*,
the whole of `x/vault` was re-diffed against upstream `protocol/v9.6.3` (rename-normalized),
with the rule that **every remaining delta must be a documented, deliberate divergence**.
The diff found real gaps that the milestone-by-milestone port had missed:

| Gap | Why it matters | Status |
|---|---|---|
| `MegavaultOwnerShares(address)` returned raw shares only | This was spec.md's own complaint ("no per-owner equity") and the §1 *user equity* data path. Upstream's by-address query returns shares, unlocks, **equity and withdrawable equity** | ✅ ported; old paginated query lives on as `MegavaultAllOwnerShares` |
| No `VaultParams` query | Per-vault params were readable only via the full `Vault` query | ✅ ported |
| `MsgUpdateDefaultQuotingParams` was gov-only | Upstream gates it operator-or-authority; M2a's notes *claimed* this was done — it was not | ✅ ported (with upstream's tests) |
| `QueryVaultResponse` never populated `most_recent_client_ids` | Field silently absent from every response; the proto field itself was missing too | ✅ fixed |
| `types.VaultKeeper` interface still had pre-megavault per-vault signatures | Stale enough that nothing could implement it; zero consumers | ✅ replaced with upstream's |
| CLI had only deposit/withdraw + 5 queries | Upstream ships set-vault-params, allocate, retrieve, update-default-quoting-params, withdrawal-info, owner-shares | ✅ ported (tx.go, query.go, util.go) |
| `MegavaultWithdrawalInfo` was GET | Upstream uses POST with body | ✅ aligned |
| `DefaultOperatorParams` had empty metadata | Upstream defaults to `Governance` / `Governance Module Account`; the operator name is user-visible on the NLP page | ✅ aligned (fixtures updated) |
| Missing upstream tests | `grpc_query_shares_test`, `grpc_query_vault_params_test`, `msg_server_update_operator_params_test`, skew-factor validation cases | ✅ ported, keeping the fork-local assertions (rejected-SetVaultParams-emits-no-event, `TestValidateMegavaultParams`) |

**Deltas verified as deliberate, left in place:** the fee engine and deposit cap (fork-local
by design), the deprecated `params` field kept in `QueryParamsResponse` (wire compat),
`types/codec.go` not registering upstream's two deprecated messages (they never existed in
this fork), `keeper/deprecated_state.go` not ported (exists solely for upstream's v7.x
migration, which this fork's state predates), `vault_test.go` fee assertions using this
fork's 3-argument `GetPerpetualFeePpm` (upstream's takes affiliate parameters this fork does
not have), and `withdraw.go` carrying *more* error logging than upstream.

The result: `x/vault` now equals upstream `v9.6.3` plus exactly the documented fork-local
additions, and every requirement in spec.md's protocol scope has an implementation and a
test.

### Two blocking bugs found by actually running the tests

Neither is vault work. Both were invisible to `tsc`, `eslint` and the build, and both break
any environment brought up from scratch.

**1. `migrate()` fails on every fresh database.**
`20260507000000_create_user_complaints_table` creates `user_complaints` *with* a
`walletAddress` column, and `20260507000001_add_wallet_address_to_user_complaints` then adds
the same column again:

```
error: alter table "user_complaints" add column "walletAddress" varchar(255) not null
default '' - column "walletAddress" of relation "user_complaints" already exists
```

The cause is that the create-table migration was edited *after* it had already been applied.
The live database is fine — it ran the original version and then the add-column migration —
but any new deployment, CI job, or local dev database cannot migrate at all.

Fixed by guarding the second migration on `hasColumn`, which makes both paths work without
rewriting migration history a second time.

**2. Ender's build puts its SQL files where ender does not look for them.**
`postgres-functions.ts` resolves `path.resolve(__dirname, '../../scripts/…')`, which from
`build/src/helpers/postgres/` means `build/src/scripts/`. Commit `f1b18fde` ("debug upddate
ender package", Sep 2024) changed the build script to copy them to `build/scripts` instead:

```diff
-"build": "... && mkdir build/src && cp -r src/scripts build/src/scripts",
+"build": "... && cp -r src/scripts build/scripts",
```

`createPostgresFunctions()` runs at ender startup ([index.ts:65](../../../indexer/services/ender/src/index.ts)),
so **a built ender image cannot install any of its `dydx_*` SQL functions** — and every
handler is one of those functions. Restored to upstream's `build/src/scripts`.

The original line's redundant `mkdir build/src` was presumably what prompted the change;
`tsc` already creates that directory, so removing the `mkdir` alone would have been correct.

### M4 notes — the vault controller port

`vault-controller.ts` replaced wholesale with upstream's version at `92733cad` (#2598, the
materialized-view commit): 362 lines → 803. That target was chosen deliberately over
`v9.6.3`: everything after it is the 2025 request-level Redis caching work, which depends on
the `vault-cache` and the `CachedMegavaultPnl` types that only exist to serve it. Porting to
v9.6.3 would have meant dragging that in; stopping at `92733cad` gives a complete, coherent
controller with nothing dangling.

What the port brings that the fork did not have: a `resolution` query parameter (hour/day)
validated by `checkSchema`, materialized-view-backed PnL reads, main-subaccount equity folded
into megavault PnL, per-interval tick filtering, pnl rebasing onto a common start date,
chunked funding-index lookups, and a current-equity tick appended to every series.

**Six pieces had to be ported first**, none of which existed in the fork:

| Piece | Where |
|---|---|
| `MEGAVAULT_MODULE_ADDRESS` / `MEGAVAULT_SUBACCOUNT_ID` | `packages/postgres/src/constants.ts` |
| `PnlTicksTable.getLatestPnlTick` | `packages/postgres/src/stores/pnl-ticks-table.ts` |
| `FundingIndexUpdatesTable.findFundingIndexMaps` (plural) | `packages/postgres/src/stores/funding-index-updates-table.ts` |
| `aggregateHourlyPnlTicks` | `comlink/src/lib/helpers.ts` |
| `AggregatedPnlTick`, `MegavaultHistoricalPnlRequest`, `VaultsHistoricalPnlRequest` | `comlink/src/types.ts` |
| `binary-searching` npm dependency, five `VAULT_*` config keys | `comlink/package.json`, `comlink/src/config.ts` |

**The megavault subaccount UUID had to be recomputed, not copied.** Upstream hardcodes
`c7169f81-0c80-54c5-a41f-9cbb6a538fdf`, which is a deterministic UUID over the string
`"<address>-0"` — and the address differs between chains. The nemo value is
`535b9ea5-af1f-51b9-b5ab-004d982a0343`, computed with the repo's own
`SubaccountTable.uuid`, with the dydx value reproduced first as a check that the derivation
was right. Copying upstream's constant would have produced a controller that silently
reported zero main-vault equity forever.

`swagger.json` and `api-documentation.md` were regenerated; the `resolution` parameter now
appears on both historicalPnl endpoints. The regeneration also picked up a
`UserComplaintsResponse` section that had never been regenerated after the fork added that
endpoint.

**Known remaining debt, not introduced here:** the generated API docs carry
*"For the deployment by DYDX token holders, use baseURL = 'https://indexer.dydx.trade/v4'"*
101 times at `HEAD` (103 now, from the two new endpoint sections). The string does not appear
anywhere in repo source — it comes from the widdershins toolchain — so fixing it is a
separate, self-contained task.

### M3 notes — indexer vault infrastructure

**The indexer forked at the same point as the protocol.** `packages/postgres/src/stores/
pnl-ticks-table.ts` is byte-identical to upstream at `fbda06f8` (2024-08-28), one day and
a handful of commits from the protocol's ancestor `e0a20164`. The last migration in this
fork is `20240717171246`, which is exactly upstream's last migration at that commit. So the
same "replay a bounded upstream range" strategy applies to the indexer, and the vault stack
lands as four upstream commits: `0330ae7d` (#2263, table), `f2fb2b90` (#2274, ender),
`35a70aad` (#2364, controller off env vars), `92733cad` (#2598, materialized views).

**M3a — proto codegen.** Regenerated `packages/v4-protos/src/codegen` with telescope. Before
trusting the toolchain, the whole tree was regenerated into a scratch directory and diffed
against the committed output: **10 files differed and all 10 were files my proto changes
touch**, everything else byte-identical. The new `indexer/protocol/v1/vault.ts` is
byte-identical to upstream's. `bundle.ts` churns 200 lines purely from import renumbering.

**M3b — `vaults` table.** Migration, `VaultModel`, `VaultTable` store, `vault-types`,
`VaultFromDatabase`, `VaultQueryConfig`, `layer1Tables`, `SQL_TO_JSON_DEFINED_MODELS`, plus
the store test and `defaultVault` seeded from `mock-generators`.

**M3c — ender.** `UpsertVaultValidator`, `UpsertVaultHandler`, `dydx_vault_upsert_handler.sql`,
`dydx_protocol_vault_status_to_vault_status.sql`, the `upsert_vault` dispatch arm in
`dydx_block_processor_ordered_handlers.sql`, subtype enum, decode case, and both validator
mappings.

**M3d — PnL views and downstream.** Hourly/daily materialized views, `VaultPnlTicksView`,
`PnlTickInterval`, the roundtable `refresh-vault-pnl` loop with its three config keys, and
the `vaults` Athena DDL. `vault-controller` now derives its vault set from the `vaults`
table via `getVaultMapping()`; `EXPERIMENT_VAULTS` / `EXPERIMENT_VAULT_MARKETS` and both
`TODO(TRA-570)` placeholders are gone. **M3's exit criterion is met: vaults are discovered
from chain events, not env vars.**

#### The fork deleted the indexer's service test suites

Discovered when the ported ender tests would not compile: `services/ender/__tests__` did not
exist. Upstream has 44 ender test files at the fork point; this fork has none. Same story
across the service layer:

| Package | Upstream at fork point | This fork |
|---|---|---|
| `services/ender` | 44 | 0 |
| `services/comlink` | 27 | 0 |
| `services/roundtable` | 17 | 0 |
| `services/socks` | 5 | 0 |
| `services/vulcan` | 5 | 0 |
| `packages/postgres` / `redis` / `kafka` | kept | kept |

The `packages/*` suites survived; the `services/*` suites did not. Restored the eight ender
test helpers (`constants`, `indexer-proto-helpers`, `validator-helpers`, `kafka-helpers`,
`kafka-publisher-helpers`, `postgres-helpers`, `redis-helpers`, `conversion-helpers`) from
the fork point so ender tests can exist again — this unblocks all future ender work, not
just vault. Two upstream defects surfaced in the process: `indexer-proto-helpers.ts` uses
`Long` as a type without importing it (upstream relies on an ambient declaration this fork
does not get), and upstream's own files violate this fork's `import/order` lint rule —
`services/ender/src` fails lint with **41 pre-existing errors**, unrelated to vault.

#### Rename debt found in the indexer

Same class of bug as M0's, in three places, all live-path:

1. **`comlink/src/lib/helpers.ts`** — `checkIfValidDydxAddress` tested
   `/^dydx[0-9a-z]{39}$/`. On a `nemo` chain this matches nothing, so the trader-search
   endpoint (`social-trading-controller`) could never resolve an address. Fixed and renamed
   to `checkIfValidNemoAddress`.
2. **`packages/postgres/src/lib/vault-addresses.json`** — 1,000 `dydx1…` vault addresses used
   by `pnl-ticks-table` to *exclude vaults from the leaderboard*. None of them match a real
   address on this chain, so vault subaccounts would have ranked on the leaderboard. All
   1,000 re-encoded to `nemo1…` preserving bytes.
3. **`__tests__/helpers/constants.ts`** — `vaultAddress` re-encoded to match, otherwise the
   leaderboard-exclusion tests would silently stop testing exclusion.

The re-encode was independently confirmed: entry 0 of the JSON re-encodes to
`nemo1c0m5x87llaunl5sgv3q5vd7j5uha26d2z2yp69`, which is exactly what
`VaultId{CLOB, 0}.ToModuleAccountAddress()` returns from this fork's protocol code.

#### Deliberate divergences from upstream in M3

| Divergence | Reason |
|---|---|
| Migration timestamps `20260804…` instead of upstream's `20240912` / `20241119` | A migration that sorts *before* already-applied ones on a live database is a footgun, and the fork already has `20260507…` migrations. Names and bodies otherwise identical. |
| Materialized views hardcode `nemo18tkxrnrkqc2t0lr3zxr5g6a4hdvqksylyq4j0f` | Upstream hardcodes its own megavault main address; a view definition cannot take a parameter. Value derived from this fork's `types.MegavaultMainAddress`. |
| `defaultVault` reuses the existing `vaultAddress` constant rather than adding upstream's `defaultVaultAddress` | The fork already had a constant with that exact value; two names for one address invites drift. |
| SQL functions keep the `dydx_` prefix | The fork has ~40 `dydx_*.sql` functions and `dydx_to_jsonb` is called from generated code. A half-renamed SQL namespace is worse than a consistently-legacy one; rename it wholesale or not at all. |
| Redis `vault-cache` **not** ported | It exists only to serve the 895-line vault controller's request-level PnL cache, and its types (`CachedMegavaultPnl`, `RedisVaultsArray`) come with that controller. Dead code until M4. |

### Unplanned work absorbed into M0: the `dydx` → `nemo` rename was never finished

Discovered while establishing the baseline — `x/vault/keeper` tests panicked in `InitChain`,
and the cause had nothing to do with vault:

1. **Bech32 addresses** in `testutil/constants/genesis.go` were still `dydx1…`, so bank
   `InitGenesis` panicked on HRP mismatch. That fixture backs **136 test files across ~20
   modules** — the entire testapp-based integration suite was dead.
2. **Gentx signatures were never regenerated.** `chain_id` was renamed to
   `localnemo-network` and the denom to `unemo`, but a signature covers the chain-id and
   signer address, so all four gentxs failed verification. Emptying `gen_txs` is not an
   option — `InitGenesis` then panics with `validator set is empty`.

Fixed by re-encoding every address (preserving the underlying bytes, so `nemo1…` maps to
the same account) and re-signing the gentxs from the mnemonics in
`testutil/constants/mnemonics.go`. The derived pubkeys matched the originals exactly,
confirming the derivation rather than merely producing something that verifies.

The sweep was then extended repo-wide: 34 files, ~130 addresses, plus the localnet/testnet
shell scripts, the containertest pre-upgrade fixture (converted **and** re-signed), and the
bech32-prefix assertions in `app/config/config_test.go`.

Also repaired: gRPC gateway route assertions in 12 `module_test.go` files requested
`/nemo-network/<module>/…` while the protos declare `/nemo_network/[v4/]<module>/…` —
wrong separator *and*, for six modules, a missing `/v4/` segment.

**Deliberately not touched:** `testing/testnet/genesis.json` (`dydx-testnet-4`, denom
`adv4tnt`) and `testing/mainnet/genesis.json` (`dydx-mainnet-1`, denom `adydx`) are dYdX's
own network genesis files, carried over from upstream. Their `dydx` addresses are correct
for those chains, 87 gentx signatures depend on them, and no script in this repo consumes
them. Renaming them would corrupt a historical record, not complete a migration.

### Result, measured against a pristine baseline

Full `go test ./...` run before and after, comparing failing-package sets:

| | Baseline | After |
|---|---|---|
| Packages passing | 91 | **135** |
| Packages failing | 59 | **15** |
| Regressions | — | **0** |

**44 packages fixed, zero regressions.** The 15 remaining failures all fail identically on
the pristine baseline and are unrelated to addresses: msg-registry and gas-metering logic
(`app/ante`, `app/msgs`, `app/process`, `lib/ante`, all five `x/clob` packages,
`testing/e2e/trading_rewards`), address-derived hash fixtures
(`indexer/off_chain_updates`), and numeric balance mismatches (`x/rewards/keeper`,
`x/ratelimit/keeper`). Several only became *visible* after this fix, because packages now
run past the first panic instead of aborting at it.

Worth noting: the five `x/clob` packages were a genuine suspicion — this change converted
addresses in two clob test files, and some clob tests assert hashes derived from addresses
(the exact failure mode present in `indexer/off_chain_updates`). The baseline comparison
ruled it out; they were already broken. That class of fixture — an expected hash computed
over an address — is the one thing a bech32 re-encode cannot fix mechanically, and is worth
checking first if similar failures appear later.

None of the remaining 15 are vault work; they are not addressed here.

---

## Milestone plan

Two tracks run in parallel after M1a lands (the indexer needs the event protos).
Sizing is engineer-weeks; ranges assume the port strategy above.

```
                        M1a gates the indexer track
                        │
Protocol:  M0 ──> M1 ──┬─> M2a ──> M2c ──> M7 ──> M8
                       │    ↑                      ↑
                  M2b ─┘    │  (design, runs early and in parallel)
                       │
Indexer:           M3 ─┴──> M4 ──> M5 ────────────┘
```

---

### M0 — Baseline, tooling, and port harness · **0.5 wk**

Make the port mechanical and repeatable before writing feature code.

- Add upstream as a read-only remote (or a pinned vendored snapshot) so
  `e0a20164..protocol/v9.6.3` is cherry-pickable and re-diffable on demand.
- Build a rename-aware port script (`dydxprotocol` ↔ `nemo_network` / `nemo-network`)
  so each upstream commit can be applied and diffed with low friction.
- Restore the two dropped upstream test files; fix the duplicated `ExportGenesis` line.
- Green baseline: `x/vault` unit tests + `app/e2e` pass; a `make test-vault` target exists.
- Decide and record the **port target tag** (recommendation: `protocol/v9.6.3` — the
  post-v7 delta is tiny and includes the skew-factor validation and withdrawal fixes).

**Exit:** clean baseline, reproducible "diff vs upstream" report, port target pinned.

---

### M1 — Protocol: withdrawal, unlock, and the keeper plumbing they need · **2–3 wks**

The core of the port. Sub-steps are ordered by dependency, each independently mergeable.

**M1a — Keeper surface + indexer events (blocks the indexer track)**
- Extend `NewKeeper` with `assetsKeeper`, `bankKeeper`, `delayMsgKeeper`,
  `indexerEventManager`; update `app.go` wiring.
- Port `types/events.go`; add vault event protos to
  `proto/nemo_network/indexer/events/events.proto`; regenerate Go + TS protos.
- Emit `UpsertVault` from `SetVaultParams`.
- **Exit:** indexer team is unblocked on M3.

**M1b — Share unlocking**
- `MsgUnlockShares`, delayed-message dispatch from `LockShares` (closes the literal
  `TODO (TRA-565)` at `keeper/shares.go:200`), unlock processing, genesis round-trip.

**M1c — Megavault withdrawal**
- `keeper/withdraw.go` (redemption value, slippage, per-vault pro-rata withdrawal),
  `MsgWithdrawFromMegavault`, `QueryMegavaultWithdrawalInfo` preview endpoint,
  `keeper/sweep_funds.go`, CLI commands.

**M1d — Vault fund management**
- `MsgAllocateToVault` / `MsgRetrieveFromVault`, `TransferToVault`,
  deactivated-vault equity exclusion, order cancellation on deactivate/stand-by,
  `GetVaultAndQuotingParams` refactor, close-only order sizing.

**Exit:** deposit → lock → unlock → withdraw round-trips on a local chain;
upstream's `x/vault` test suite ports over and passes.

---

### M2 — Protocol: operator params and the performance-fee engine · **3–4 wks**

Split into a portable part and a net-new part. The portable part should land first so it
does not sit behind the design phase.

**M2a — Operator params + deposit cap (port) · 1–1.5 wks**
- `OperatorParams` (address + name/description metadata), `MsgUpdateOperatorParams`,
  operator-gated `SetVaultParams` / `UpdateDefaultQuotingParams`,
  `grpc_query_vault_params`, operator params surfaced in the params query.
- **New:** max-deposit-cap validation on `MsgDepositToMegavault`.

**M2b — Fee mechanism design · 0.5–1 wk, gates M2c**

There is nothing upstream to port: upstream `OperatorParams` is only
`{ operator address, { name, description } }`. This needs a written design (ADR in
`x/vault/spec/`) resolving at minimum:

| Question | Options | Recommendation |
|---|---|---|
| Fee basis | Global NAV-per-share high-water mark vs. per-depositor cost basis | **Global per-share HWM** — megavault shares are fungible and pooled; per-depositor basis requires per-owner lot tracking in state |
| Crystallization trigger | On withdrawal only, vs. periodic (epoch/block) accrual | **Both**: accrue on an epoch, crystallize on withdrawal, so exiting depositors pay their share |
| Payment mechanism | Mint shares to operator (dilution) vs. transfer quote quantums | **Mint shares** — no subaccount transfer on every accrual, no settlement failure mode |
| Fee vs. profit share | Are these one field or two distinct charges? | Needs a product answer — spec.md lists both, the frontend shows both at 0% |
| Minimum operator share | Enforced where — blocking deposits, or blocking operator withdrawal? | **Block operator withdrawal** below the floor; blocking user deposits punishes the wrong party |
| Existing depositors | See risk section below | HWM seeded at upgrade-height NAV/share |

**M2c — Fee engine implementation · 1.5–2 wks**
- HWM state + accrual + crystallization + distribution, params (`fee_ppm`,
  `profit_share_ppm`, `min_operator_share_ppm`), events, queries, CLI.
- Fee interaction with the M1c withdrawal path (redemption value is net of crystallized fee).
- Invariant tests: fees never exceed profit, no double-charge across accrual + withdrawal,
  HWM monotonic, total shares reconcile after operator minting.

**Exit:** operator can manage vaults; fees accrue, crystallize, and are distributed;
cap and minimum-operator-share enforced; invariants covered by tests.

---

### M3 — Indexer: vault data infrastructure · **1.5–2.5 wks**

Mostly portable. Depends on M1a.

- `vaults` table + migration + `vault-model` + `vault-table` store + `vault-types`.
- Ender: `upsert-vault-validator`, `upsert-vault-handler`,
  `nemo_vault_upsert_handler.sql`, status-enum helper SQL, block-processor registration.
- `vault_pnl_ticks` hourly + daily materialized views; roundtable `refresh-vault-pnl`.
- Redis `vault-cache`; Athena DDL for `vaults`.
- **Delete the `EXPERIMENT_VAULTS` / `EXPERIMENT_VAULT_MARKETS` config hack** and the
  `TODO(TRA-570/571)` placeholders in `vault-controller.ts`.

**Exit:** vaults are discovered from chain events, not env vars; vault PnL views populate.

---

### M4 — Indexer: port + harden the existing vault endpoints · **1–1.5 wks**

- Port upstream's 895-line `vault-controller`: resolution argument, main-subaccount
  equity inclusion, tick de-duplication/filtering, materialized-view backing,
  request-level PnL caching, chunked funding-index lookups.
- Point the three existing endpoints at the `vaults` table.

**Exit:** `megavault/historicalPnl`, `vaults/historicalPnl`, `megavault/positions`
are production-grade and vault-table-driven.

---

### M5 — Indexer: the NLP page API surface (net-new) · **3.5–5 wks**

No upstream equivalent; this is where the real design work is. Reuses `pnl_ticks`,
the existing controllers, and the max-drawdown / equity work on `feat/update-indexer`.

| Deliverable | Size |
|---|---|
| Summary & metrics: TVL, APR (30d annualized), NAV, MDD, volume, user shares/equity | 1.5–2 wks |
| Charts: return % / TVL / NAV × 6 periods, with resolution downsampling | ~1 wk |
| Bottom tables (8 tabs) scoped to vault subaccounts — thin wrappers over existing controllers | 1–1.5 wks |
| My Activity: deposits/withdrawals with share counts, NAV-at-time, full tx hash | 1–1.5 wks |
| Deposit/withdraw status by tx hash (`requestId = txHash`, wallet-signed broadcast) | 0.5–1 wk |

**Contract note:** the frontend has zero backend integration today, so the API shape is
ours to define. Recommended: frontend broadcasts wallet-signed txs and polls a thin
status endpoint. A backend transaction-submission gateway would add 2–3 weeks plus
key-custody surface, and nothing in the spec or frontend requires it.

**Exit:** every panel on the `/vault` page is served by a real endpoint.

---

### M6 — Abacus / client integration · **deferred**

Out of scope for this phase by decision. The abacus lives in a separate repository;
upstream `v4-abacus` already implements the client side of the megavault deposit/withdraw
flow and is a strong porting candidate. Scope it once M5 fixes the API contract — doing it
earlier means building against an interface that is still moving.

---

### M7 — Upgrade handler · **1.5–2 wks**

Mandatory, because the chain carries live state. This is the **first upgrade handler this
fork will ship**: `Upgrades` in [app/upgrades.go:14](../../app/upgrades.go#L14) is an empty
slice and `app/upgrades/` contains only `types.go`.

- Build the upgrade package scaffold (`app/upgrades/<name>/{constants,upgrade}.go`),
  register in `Upgrades`, wire `setupUpgradeHandlers`.
- **No `StoreUpgrades` needed for vault** — all new state (operator params, fee state,
  unlock queue) lives under the existing `vault` store key as new prefixes. Confirm before
  relying on this.
- Initialize new state that `InitGenesis` would otherwise set and a running chain will not
  have: `OperatorParams`, fee params, and the initial high-water mark.
- **Two of upstream's three v7 vault migrations are already unnecessary here** —
  `migrateVaultSharesToMegavaultShares` and `migrateVaultQuotingParamsToVaultParams` are
  both already applied in this fork (it sits after upstream #2109 and the vault-params
  refactor). Do not port them.
- **Indexer backfill:** existing vaults never emit `UpsertVault`, so the indexer's `vaults`
  table would start empty. Upstream's handler happens to solve this by calling
  `SetVaultParams` per vault during migration. Do the same deliberately: iterate all vaults
  in the handler and re-set params, so M1a's event backfills the table at upgrade height.
- Upgrade-container test (upstream has `upgrade_container_test.go` as a template).

**Exit:** upgrade rehearsed on a state export from the live chain, with vault balances,
shares, and the indexer's `vaults` table all correct on the far side.

---

### M8 — E2E and testnet · **1–1.5 wks**

- End-to-end deposit → lock → unlock → withdraw across protocol + indexer.
- Fee accrual and crystallization observed end-to-end over a multi-epoch soak.
- Testnet soak, load check on the new endpoints.

---

## Risk: charging fees to existing depositors

Decisions 1 and 3 together create a problem that neither creates alone.

The chain already holds real deposits made under a **0%-fee regime**. Introducing a
performance fee at an upgrade height means existing depositors begin paying a fee they did
not agree to when they deposited. Worse, a naively initialized high-water mark charges them
retroactively: if HWM starts at zero (or at their original deposit NAV), the operator
immediately crystallizes a fee on **all gains earned before the fee existed**.

Mitigations, in order of preference:

1. **Seed the HWM at the NAV-per-share observed at upgrade height.** Only post-upgrade
   gains are ever charged. Cheap, and it is the obviously correct default.
2. Announce the parameter change ahead of the upgrade with a withdrawal window, so
   depositors who object can exit fee-free. Requires M1c (withdrawals) to ship *before*
   M2c (fees) — which is already the plan's ordering, and is now a hard constraint.
3. Launch fee params at 0 and raise them via a separate governance action after the
   upgrade, decoupling "the mechanism exists" from "the mechanism charges".

This is a product and governance decision as much as a technical one, and it should be
settled during M2b rather than discovered during M7.

---

## Totals

| Milestone | Track | Effort | Nature |
|---|---|---|---|
| M0 baseline & port harness | protocol | 0.5 | port |
| M1 withdrawal + unlock + plumbing | protocol | 2–3 | port |
| M2a operator params + deposit cap | protocol | 1–1.5 | port + small new |
| M2b fee mechanism design | protocol | 0.5–1 | **new** |
| M2c fee engine implementation | protocol | 1.5–2 | **new** |
| M3 indexer vault infra | indexer | 1.5–2.5 | port |
| M4 endpoint port/harden | indexer | 1–1.5 | port |
| M5 NLP page API surface | indexer | 3.5–5 | **new** |
| M7 upgrade handler | protocol | 1.5–2 | new (upstream template) |
| M8 e2e + testnet | both | 1–1.5 | — |
| **Total (M6 deferred)** | | **~14.5–20 eng-weeks** | |
| **Calendar, 2 engineers (protocol ∥ indexer)** | | **~8–10 weeks** | |

Against [spec.md](./spec.md)'s Scenario B range (~13–20 eng-weeks, ~8–10 weeks calendar):
the totals land in the same place, but the composition shifted. The upstream-ancestor
finding removed roughly 1–1.5 weeks of protocol uncertainty; the mandatory upgrade handler
put it straight back. Roughly **half the remaining effort is genuinely new work** — the fee
engine (M2b/M2c) and the NLP page API (M5) — and that half carries nearly all the risk.

### Sequencing constraints

- **M1a gates the entire indexer track.** Land the keeper/event plumbing early, even
  though it delivers no user-visible behavior; everything in M3 waits on it.
- **M1c must ship before M2c.** Depositors need a working exit before fees turn on
  (see the risk section). This is now a hard ordering constraint, not a preference.
- **M2b gates M2c** and should start early — a design phase parallelizes cleanly against
  M1's port work, and it is the only milestone that can block on decisions outside the team.
- **M5 is the critical path.** It is the longest single milestone, it is entirely new, and
  its contract depends on `NewSpec.md` (see below).

---

## Open questions

1. **`NewSpec.md` is referenced by [spec.md](./spec.md) but is not in this repository.**
   M5's endpoint list is reconstructed from spec.md's summary. The source document is
   needed to finalize the API contract, and M5 is the critical path.
2. **Fee vs. profit share:** are these one charge or two? Needed for M2b.
3. **Fee-introduction governance:** which mitigation from the risk section applies —
   HWM seeded at upgrade height, a pre-upgrade withdrawal window, or launch-at-zero
   followed by a separate governance raise? Needed before M7.
