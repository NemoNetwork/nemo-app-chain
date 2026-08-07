# Vault Module — Work Summary

**Period:** 2026-08-03 → 2026-08-04
**Branch:** `feat/vault`
**Scope covered:** M0 (baseline) and all of M1 (protocol: withdrawal path), plus the
operator-params half of M2a pulled forward, plus an unplanned repo-wide rename repair.

Companion documents: [spec.md](./spec.md) (status & estimation) ·
[MILESTONES.md](./MILESTONES.md) (delivery plan and per-milestone detail)

---

## 1. Headline

**The megavault protocol path is complete end to end** — deposit → lock → unlock →
withdraw, plus operator fund management (allocate / retrieve) and quoting behavior for every
vault status.

| | Before | After |
|---|---|---|
| Packages passing (`go test ./...`) | 91 | **138** |
| Packages failing | 59 | **13** |
| Regressions introduced | — | **0** |
| Vault `Msg` RPCs | 3 | **8** |
| Vault `Query` RPCs | 5 | **6** |
| Vault module `.go` files | 57 | **75** |

Total diff: **101 files changed, ~7,200 insertions**, plus 33 new files.

Every milestone was verified by running the full suite against a **stashed pristine
baseline** and diffing the failing-package sets, rather than reasoning about impact. All 13
remaining failures fail identically on the untouched baseline.

---

## 2. The finding that shaped everything

`protocol/x/vault` is **byte-identical to upstream `dydxprotocol/v4-chain` at commit
`e0a20164`** (2026-08-28 upstream date, *"implement megavault shares locking" #2160*),
modulo the `dydxprotocol` → `nemo_network` proto rename.

Verified by normalizing import paths and diffing: every hand-written `.go` file matched at
**zero lines**, except a duplicated `ExportGenesis` line, one lint cleanup in `orders.go`,
and two upstream test files the fork had dropped.

This converted the protocol work from "design and build" into **replaying a known upstream
commit range**, and it is why the M1 milestones landed quickly and with low risk. The
milestone plan, effort estimates, and sequencing in [MILESTONES.md](./MILESTONES.md) are all
built on this.

---

## 3. Unplanned: the `dydx` → `nemo` rename was never finished

Discovered while establishing an M0 baseline — `x/vault/keeper` tests panicked in
`InitChain` for reasons unrelated to vault.

**Two distinct defects, both repo-wide:**

1. **Bech32 addresses** in `testutil/constants/genesis.go` were still `dydx1…`, so bank
   `InitGenesis` panicked on HRP mismatch. That fixture backs **136 test files across ~20
   modules** — the entire testapp-based integration suite was dead.
2. **Gentx signatures were never regenerated.** The fork renamed `chain_id` to
   `localnemo-network` and the denom to `unemo`, but a signature covers the chain-id and
   signer address, so all four gentxs failed verification. Emptying `gen_txs` is not an
   escape — `InitGenesis` then panics with `validator set is empty`.

**Fixed** by re-encoding every address (preserving the underlying bytes, so `nemo1…` maps to
the same account) and re-signing the gentxs from the mnemonics in
`testutil/constants/mnemonics.go`. The derived pubkeys matched the originals exactly, which
confirmed the key derivation rather than merely producing something that verifies.

Then swept repo-wide: **34 files, ~130 addresses**, plus localnet/testnet/staging shell
scripts, the containertest pre-upgrade fixture (converted **and** re-signed), bech32-prefix
assertions, and gRPC gateway route paths in 12 `module_test.go` files that requested
`/nemo-network/…` where the protos declare `/nemo_network/[v4/]…`.

**Deliberately not touched:** `testing/testnet/genesis.json` (`dydx-testnet-4`, denom
`adv4tnt`) and `testing/mainnet/genesis.json` (`dydx-mainnet-1`, denom `adydx`) are dYdX's
own network genesis files carried over from upstream. Their `dydx` addresses are correct for
those chains, 87 gentx signatures depend on them, and no script in this repo consumes them.
Renaming them would corrupt a historical record rather than complete a migration. If they
are unwanted, deleting is the sensible action — not rewriting.

**Result: 44 packages fixed, zero regressions.** Several failures only became *visible*
afterwards, because packages now run past the first panic instead of aborting at it.

---

## 4. Bugs found and fixed along the way

| Bug | Where | Impact |
|---|---|---|
| `crisis.MsgVerifyInvariant` listed as an unsupported message | `lib/ante/unsupported_msgs.go` | **Production.** `IsUnsupportedMsg` is on the live ante path, so crisis invariant-verification txs were rejected on chain despite being allow-listed. Upstream has two deprecated vault msgs in that `case`; neither exists here, and one was evidently substituted to keep the `case` syntactically valid — the orphaned `// vault` comment above it was the tell. |
| 20 message type URLs keyed `/nemo-network.…` (hyphen) instead of `/nemo_network.…` | `app/msgs/{internal,normal}_msgs.go` | Test-only. Verified with `sdk.MsgTypeURL`. The maps are consumed only by tests, so this was not rejecting live transactions — but the tests enforcing the catalogue against the real registry were failing. |
| Expected-message lists never re-sorted after the rename | `app/msgs/*_test.go` | Test-only. Compared with `require.Equal` against sorted keys, and `dydxprotocol` sorted before `/ibc.` where `nemo_network` sorts after. |
| Duplicated `genesis.Vaults = k.GetAllVaults(ctx)` | `x/vault/genesis.go` | Harmless copy-paste in `ExportGenesis`. |

---

## 5. What was built, by milestone

### M0 — baseline and port harness
Green baseline established (via the rename repair above), the duplicated `ExportGenesis`
line removed, and the two dropped upstream test files restored
(`keeper/shares_test.go`, `types/params_test.go`) and confirmed executing.

Proto toolchain stood up locally (no docker/buf present): `buf` v1.34.0 plus
`protoc-gen-gocosmos` and `protoc-gen-grpc-gateway`. **Validated by regenerating an
unchanged proto and diffing byte-for-byte against the committed output** — a mismatched
generator version would have silently rewritten every `.pb.go` in the repo.

### M1a — keeper plumbing and indexer events *(unblocked the indexer track)*
`VaultStatus` enum + `UpsertVaultEventV1` protos, the `VaultStatusToIndexerVaultStatus`
mapper, `NewUpsertVaultEvent`, subtype/version constants, and `SetVaultParams` emitting
`UpsertVault`. `NewKeeper` extended with `assetsKeeper`, `bankKeeper`, `delayMsgKeeper` and
`indexerEventManager` — all four at once so `app.go` and the fixtures churn once rather than
three times.

### M1b — share unlocking *(closes `TODO (TRA-565)`)*
`MsgUnlockShares`, the `UnlockShares` keeper method, and `LockShares` scheduling its own
unlock via `delayMsgKeeper.DelayMessageByBlocks`.

### M1c — megavault withdrawal
`withdraw.go` (slippage, `RedeemFromMainAndSubVaults`, `WithdrawFromMegavault`),
`MsgWithdrawFromMegavault`, the `MegavaultWithdrawalInfo` preview query, `sweep_funds.go`
wired into the EndBlocker, the withdrawal CLI command, and deposit/withdraw/sweep typed
events. Required porting `GetVaultAndQuotingParams`, `GetVaultLeverageAndEquity`,
`GetVaultClobPerpAndMarket`, deactivated-vault equity exclusion, the `lib/vault` helper
package, `lib.BigRatMin`, and `PerpetualsKeeper.GetLiquidityTier` first.

### M2a (operator half, pulled forward)
`OperatorParams` / `OperatorMetadata`, `MsgUpdateOperatorParams`, keeper accessors, genesis
round-trip, operator params in the params query, operator-gated `MsgSetVaultParams`.
Pulled ahead of M1d because `MsgAllocateToVault` is gated on *"authority **or operator**"* —
upstream's own ordering confirms the dependency.

### M1d — allocate / retrieve and quoting modes *(completes M1)*
`MsgAllocateToVault`, `MsgRetrieveFromVault`, the `AllocateToVault` keeper method, and
`orders.go` replaced wholesale with upstream's: close-only reduce-only sizing, no quoting
for deactivated/stand-by, post-only orders on inventory-increasing sides,
`CancelVaultClobOrder` / `TryToCancelVaultClobOrder`, and cancellation of no-longer-needed
orders. `SetVaultParams` now refuses to deactivate a positive-equity vault and cancels
outstanding orders on the move to deactivated or stand-by.

**Vault error codes 1–31 now match upstream exactly**, so future ports slot in without
collisions.

---

## 6. Traps worth remembering

These cost real debugging time and will recur:

1. **Nested-message signers need explicit registration.** `MsgWithdrawFromMegavault` signs
   via `subaccount_id`; because the signer sits in a nested message, the SDK cannot resolve
   it from the annotation. It needs an entry in `CustomGetSigners` in
   [app/module/interface_registry.go](../../app/module/interface_registry.go). Without it,
   every withdrawal fails `CheckTx` with *"no cosmos.msg.v1.signer option found for message
   nemo_network.subaccounts.SubaccountId"*. **Any future vault message signing via a
   subaccount needs the same entry.**

2. **Message classification.** A vault message gated on *"authority or operator"* belongs in
   `NormalMsgs`; one gated on module authority alone belongs in `InternalMsgs` and
   `lib/ante/internal_msg.go`. Getting this wrong yields
   *"internal msg cannot be submitted externally"* on every test.

3. **Indexer-event assertions can pass vacuously.** `ProduceBlock` returns `nil` when the
   indexer is disabled, and the plain test-app builder disables it. Event tests must inject
   `msgsender.NewIndexerMessageSenderInMemoryCollector()` via `WithAppOptions`.

4. **Genesis applies the deactivation rule.** `InitGenesis` calls `SetVaultParams`, so a
   genesis state with a deactivated vault whose subaccount is funded now panics on
   `InitChain`.

5. **Address-derived hash fixtures cannot be fixed by re-encoding.** A bech32 re-encode
   preserves the address bytes but changes the encoded string that fed any stored hash. This
   is the failure mode already present in `indexer/off_chain_updates`; check it first if
   similar failures appear.

---

## 7. Deliberate divergences from upstream

Each is intentional and commented at the site, so it is not mistaken for an oversight:

- **`QueryParamsResponse` field numbering.** Upstream dropped the deprecated `params` field
  and renumbered, making `operator_params` field 2. On a chain with live clients that is
  wire-breaking, so `operator_params` was appended at **field 3** and the deprecated field
  left in place.
- **`PerpetualsKeeper.GetLiquidityTier` added to the interface.** Upstream's committed mock
  has the method while its interface does not — upstream's mock is out of sync with itself.
  Rather than copy that, the method was declared properly (the concrete keeper already
  implements it) and the mock regenerated.
- **`ErrInvalidDepositAmount` → `ErrInvalidQuoteQuantums`.** Matches upstream's broadening of
  error code 4 to cover deposits, allocations, retrievals and withdrawals. Same code, so no
  error-code break.
- **`orders.go` lint cleanup retained** (`len(x) == 0` over
  `x == nil || len(x) == 0`) when the file was replaced with upstream's.

---

## 8. Remaining failures (13, all pre-existing)

None are vault-related; all fail identically on the pristine baseline:

- `x/clob`, `x/clob/{ante,e2e,keeper,memclob,types}` — clob short-term-order semantics
- `app/ante`, `app/e2e`, `app/process` — msg-registry and gas-metering logic
- `indexer/off_chain_updates` — address-derived hash fixtures
- `testing/e2e/trading_rewards`
- `x/rewards/keeper`, `x/ratelimit/keeper` — numeric balance mismatches

The clob cluster is the largest and should go to whoever owns clob.

---

## 9. What is next, and what is blocked

**Unblocked and portable:**
- **M3 — indexer vault infrastructure.** Unblocked since M1a and entirely portable from
  upstream (`vaults` table, hourly/daily PnL views, ender `upsert-vault` handler, roundtable
  `refresh-vault-pnl`). Feeds **M5, the critical path**. *Recommended next.*
- **M2a remainder — max-deposit cap.** Small, net-new, no upstream equivalent.

**Blocked on decisions:**

| Blocked item | Needs |
|---|---|
| **M2b/M2c — fee engine** | Are **fee and profit share one charge or two?** Nothing on the fee engine can start without the mechanism design. |
| **M7 — upgrade handler** | The **high-water-mark seeding decision** — see the risk section in [MILESTONES.md](./MILESTONES.md). Introducing a fee on a chain with live deposits charges existing depositors retroactively unless the HWM is seeded at upgrade-height NAV/share. |
| **M5 — NLP page API** | **`NewSpec.md`** is referenced by [spec.md](./spec.md) but absent from the repo. M5's endpoint contract is currently reconstructed from summary. |

The ordering constraint from the plan is now satisfied: **withdrawals ship before fees**, so
depositors have an exit before any fee mechanism turns on.

---

## 10. Repo hygiene notes

- Working tree is clean of stray artifacts; a temporary stash created during baseline
  checking was inspected and dropped (it was strictly superseded).
- The bech32 re-encoder and gentx re-signer used for the rename repair were scaffolding and
  were **not** committed. They are preserved in the session scratchpad. Worth landing
  properly if genesis regeneration should be repeatable — note that
  `testing/testnet-local/local.sh` still hard-codes `dydx1` test accounts, so **regenerating
  the localnet genesis today would reintroduce the bug**.
- Six files flagged by `gofmt` were already unformatted at `HEAD` and were left alone rather
  than adding unrelated churn.
- `mocks/Makefile` updated so `make mock-gen` reproduces the new `AssetsKeeper` mock.

---

## 11. Addendum — M3, M2 and M7 (2026-08-04)

Everything in sections 1–10 above describes the protocol port through M1d. This addendum
covers the work that followed. Full detail lives in [MILESTONES.md](./MILESTONES.md); this is
the short version and the things that would otherwise be lost.

### What shipped

| Milestone | Outcome |
|---|---|
| M3 — indexer vault infrastructure | ✅ complete. Vaults are discovered from chain events, not env vars. |
| M2a — max-deposit cap | ✅ complete. `MegavaultParams`, `MsgUpdateMegavaultParams`, cap check on deposit. |
| M2b — fee mechanism design | ✅ [ADR-001](./spec/adr-001-megavault-fees.md). |
| M2c — fee engine | ✅ complete. HWM, dual accrual, crystallization, operator share floor. |
| M7 — upgrade handler | ✅ state migration + indexer backfill. **Rehearsal on a live state export still owed.** |
| M4 — vault controller port | ✅ complete. 362 → 803 lines; resolution, view-backed PnL, main-subaccount equity. |
| Spec-compliance audit | ✅ x/vault re-diffed against upstream v9.6.3; every unintended delta closed (per-owner equity query, VaultParams query, operator gating on quoting params, CLI, client-ids field, metadata defaults). Details in MILESTONES. |

### The two blocking questions from section 9 are resolved

1. **"Fee vs. profit share — one charge or two?"** Answered by evidence rather than by
   waiting: the frontend renders three distinct zero-valued fields (*operator fee*, *minimum
   operator share*, *profit share*), and a UI does not show one charge twice. Implemented as
   two independent parameters, so if product disagrees, setting either to zero collapses the
   design with **no code change**.
2. **"HWM seeding."** The upgrade handler seeds the high-water mark at the NAV per share
   observed at upgrade height, and ships the fee parameters at zero so governance turns the
   engine on separately. That applies two of the three mitigations from the risk section at
   once, and the third (withdrawals before fees) was already satisfied by the M1c/M2c
   ordering.

Only **`NewSpec.md`** remains outstanding, and it gates M5 alone — which is now the only
milestone left with unblocked predecessors and no way to start.

### Findings worth carrying forward

- **The indexer forked at the same commit as the protocol** (`fbda06f8`, 2024-08-28). The
  same "replay a bounded upstream range" strategy applies to both halves of the repo.
- **The fork deleted the indexer's service test suites** — ender 44→0, comlink 27→0,
  roundtable 17→0, socks and vulcan 5→0. The `packages/*` suites survived. The ender test
  helpers have been restored, which unblocks all future ender work, not just vault.
- **More rename debt, on live paths**: `checkIfValidDydxAddress` matched `/^dydx…/` so
  address search could never resolve; `vault-addresses.json` held 1,000 `dydx1…` addresses
  used to exclude vaults from the leaderboard, so it excluded nothing. Both fixed.
- **A sentinel bug worth remembering**: `last_accrual_time == 0` means "never accrued", but
  zero is also a real Unix timestamp — and the test app's genesis block time is 1 ns after
  the epoch, so `Unix()` is exactly 0. The fix was not a better guard but always initializing
  the clock, plus a one-year clamp bounding the damage if that is ever missed.
- **A constant that had to be recomputed rather than copied**: upstream hardcodes the
  megavault subaccount UUID, which is derived from the chain's address string. Copying it
  would have produced a controller that silently reported zero main-vault equity forever. The
  nemo value was computed with the repo's own `SubaccountTable.uuid`, verifying the
  derivation by reproducing upstream's value from upstream's address first. This is the same
  class of trap as the address-derived hash fixtures in §6.

### Verification status, stated plainly

- **Protocol**: `go build ./...` clean; `x/vault`, `app`, `app/msgs`, `app/upgrades` all
  pass. The remaining failures (`app/ante`, `app/e2e`, `app/process`,
  `indexer/off_chain_updates`, the `x/clob` packages, `x/rewards/keeper`,
  `x/ratelimit/keeper`) are the same pre-existing set documented in section 8.
- **Indexer**: `tsc` and `eslint` clean on every touched package; `pnpm -r build:prod` green
  across all 19 workspace projects; regenerated protos byte-compared against upstream.
- **DB-backed indexer tests now run.** With no Postgres on the machine and no root to install
  one, a userspace PostgreSQL 18 was stood up from the `embedded-postgres` npm package
  (binaries only; nothing installed system-wide, nothing written outside the scratchpad).
  `packages/postgres` passes **362 tests across 35 suites**, including the new `vault-table`
  and `vault-pnl-ticks-view` suites; both `services/ender` vault suites pass, exercising the
  real SQL handler and block-processor dispatch.
- **Still not verified**: `services/comlink` (the fork deleted its 27 test files — nothing to
  run), Redis- and Kafka-backed suites, and the M7 upgrade rehearsal, which needs a real
  state export.

### Two blocking bugs that only running the tests could find

Neither is vault work; both survived `tsc`, `eslint` and a clean build, and both break any
environment brought up from scratch. Full detail in [MILESTONES.md](./MILESTONES.md).

1. **`migrate()` fails on every fresh database.** The `user_complaints` create-table migration
   was edited after it had already been applied, so it now creates a column that the next
   migration also adds. Production is unaffected; every new deployment, CI run and local dev
   database is not. Fixed by guarding the second migration on `hasColumn`.
2. **Ender's build copies its SQL files to the wrong directory.** `createPostgresFunctions()`
   runs at ender startup and reads from `build/src/scripts/`; commit `f1b18fde` changed the
   build to copy them to `build/scripts/`. A built ender image therefore cannot install any
   of its `dydx_*` SQL functions — and every block handler is one of those functions.
   Restored to upstream's path.
