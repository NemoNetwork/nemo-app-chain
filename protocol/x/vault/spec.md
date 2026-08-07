# Vault Module — Development Status & Work Estimation

**Date:** 2026-07-24
**Scope:** Work estimation for the Vault / NLP page backend requirements in [NewSpec.md](./NewSpec.md), based on three sources of research:

1. **Protocol** — code review of `protocol/x/vault` (Cosmos SDK, dYdX v4 fork)
2. **Indexer** — code review of `indexer/` (comlink, ender, roundtable, postgres)
3. **Frontend** — bundle analysis of the `ui-vault` preview deployment on Vercel (`nemo-network-interface`, `/vault` page)

---

## TL;DR

The chain can accept deposits and mint vault shares, but **withdrawals, the operator/fee model, and nearly all indexer-side vault APIs do not exist yet**. The frontend `/vault` page is a finished-looking but **100% static mock** with no backend integration, which means the API contract is still fully negotiable.

Two findings from the frontend materially reduce the estimate: the UI hard-codes **0% operator fee / profit share** (suggesting the fee *mechanism* can be deferred) and its only transaction flow is **wallet-signed broadcast** (eliminating the need for a backend transaction gateway).

| Scenario | Effort | Calendar (2 engineers) |
|---|---|---|
| **A — 0% fee at launch, wallet-signed txs** (likely, per frontend) | **~10–15 engineer-weeks** | **~6–8 weeks** |
| B — full performance-fee mechanism required at launch | ~13–20 engineer-weeks | ~8–10 weeks |

---

## Current Development Status

### Protocol (`protocol/x/vault`) — deposit-only megavault

Corresponds to dYdX megavault at a **pre-withdrawal** stage.

**Implemented:**

- `MsgDepositToMegavault` with proportional share minting (1:1 on first deposit, `quantums × totalShares / equity` after) — `keeper/deposit.go`
- Megavault equity/TVL computation from subaccount net collateral — `keeper/vault.go:17`
- Vault quoting params, per-vault params (gov-gated), per-block CLOB order refresh in EndBlocker
- Queries: `Params`, `Vault`, `AllVaults`, `MegavaultTotalShares`, `MegavaultOwnerShares` (raw shares only — no per-owner equity)
- Share **locking** scaffold (`keeper/shares.go:146`) — locks are written but nothing ever unlocks or consumes them (literal `TODO (TRA-565)` at `keeper/shares.go:200`)

**Missing entirely (zero code, not partial):**

- Any withdrawal path — no `MsgWithdrawFromMegavault`, no share burning, no payout, no slippage handling, no withdrawal-preview query
- Operator concept — no operator address, name, metadata, or params
- Any fee — no performance fee, no profit share, no fee accrual anywhere
- Unlock processing in ABCI, `MsgUnlockShares`
- Max-deposit-cap validation

### Indexer — placeholder-grade vault support

- Three early-generation endpoints exist (`/vault/v1/megavault/historicalPnl`, `/vaults/historicalPnl`, `/megavault/positions`), but vault subaccounts are **hard-coded via env config** (`EXPERIMENT_VAULTS`), with explicit `TODO(TRA-570/571)` markers.
- **No** `vaults` table, no vault migrations/models, no vault event handling in ender, no vault roundtable jobs, no shares/deposit/withdraw/activity endpoints.
- Reusable building blocks all exist: positions, orders, fills, transfers, funding controllers; `pnl_ticks` + the roundtable job computing them; and the portfolio controller built out on the `feat/update-indexer` branch (equity, live-pnl, realized-pnl, **max-drawdown**, profit-factor).

---

## Gap Analysis vs. NewSpec.md

| Spec section | Status | What's needed |
|---|---|---|
| §1 Summary — TVL, user equity | 🟡 partial | Equity exists on-chain; needs indexer exposure + user-shares data path |
| §1/§2 — APR | 🟡 partial | Abacus already computes 30-day annualized APR from PNL ticks; serve it server-side |
| §1/§2 — Operator name/fee, min operator share, profit share | 🔴 missing | Fields don't exist anywhere on-chain; frontend shows all-zero values → likely fields-only for v1 |
| §2 — Vault PNL, MDD, volume | 🟡 partial | `pnl_ticks` + branch's max-drawdown work reusable; needs vault-scoped aggregation |
| §3 Charts — return %, TVL, NAV × 6 periods | 🟡 partial | historicalPnl endpoint close for return %; NAV/TVL series + period/resolution handling is new |
| §4 Deposit flow | 🟡 partial | Chain side works; max-cap validation and status-by-txHash endpoint are new |
| §5 Withdraw flow | 🔴 missing | Entire withdrawal path, protocol + indexer |
| §6 Bottom tables (8 tabs) | 🟢 mostly exists | Reuse existing controllers scoped to vault subaccounts; "My Activity" (shares + NAV-at-time + tx hash) is new |
| §7 Pagination | 🟢 exists | Standard in existing controllers |
| §8 Wallet scoping/auth | 🟡 partial | Follows existing patterns; explicit unauthenticated responses need adding |
| §9 Endpoint shapes | 🟢 negotiable | Frontend has zero integration — backend defines the contract |

---

## Work Estimate

### Workstream 1 — Protocol: withdrawals & share unlock (**2.5–4 weeks**)

`MsgWithdrawFromMegavault` (burn shares, redemption value with slippage), withdrawal-preview query, unlock processing for the existing lock scaffold, upgrade handler, proto regen, tests. Low end assumes **porting from upstream dydxprotocol v6/v7**, which already has all of this — strongly recommended over rebuilding, since this repo is a fork. The frontend bundle's Abacus code already implements the client side of exactly this flow.

### Workstream 2 — Protocol: operator & fee fields (**1–1.5 weeks**; **3–4 weeks** if a real fee engine is required)

The frontend hard-codes 0% for operator fee, minimum operator share, and profit share. If launch parameters are confirmed zero, v1 only needs: operator params (name, address, metadata), the fee/min-share fields stored and served, and max-deposit cap — **no accrual/crystallization/distribution engine**. If a real performance-fee mechanism is required at launch, this grows to 3–4 weeks and needs mechanism design first (upstream dYdX has no performance fee to port).

### Workstream 3 — Indexer: vault infrastructure (**2–3 weeks**)

`vaults` table + migration + model + store, ender handling for vault events, megavault PNL/NAV tick materialization in roundtable, replace the `EXPERIMENT_VAULTS` config hack. Largely portable from upstream's indexer.

### Workstream 4 — Indexer: vault-page API endpoints (**3.5–5 weeks**)

- Summary/metrics/performance endpoints (APR, TVL, NAV, MDD, volume, user deposits/equity): ~1.5–2 wks — reuses `pnl_ticks` and the portfolio controller's max-drawdown/equity work
- Charts endpoint, 3 metrics × 6 periods with resolution downsampling: ~1 wk
- Bottom-table endpoints scoped to the vault (thin wrappers over existing controllers): ~1–1.5 wks
- My Activity / deposits & withdrawals with shares, NAV-at-time, full tx hash: ~1–1.5 wks

### Workstream 5 — Deposit/withdraw status lifecycle (**0.5–1 week**)

The frontend's only transaction flow is wallet-signed broadcast (Abacus), and its activity tables show tx hashes directly — nothing suggests a backend-owned lifecycle. Recommended contract: **frontend broadcasts, `requestId = txHash`**, and a thin status endpoint serves Pending/Completed/Failed from indexed transactions (spec §4's "Processing" state barely exists on-chain). The alternative — a backend transaction-submission gateway owning the full lifecycle — would cost 2–3 weeks plus key-custody/security surface, and nothing in the spec or frontend requires it.

### Totals

| | Scenario A (0% fee, wallet-signed) | Scenario B (full fee engine) |
|---|---|---|
| Workstream 1 — withdrawals | 2.5–4 wks | 2.5–4 wks |
| Workstream 2 — operator & fees | 1–1.5 wks | 3–4 wks |
| Workstream 3 — indexer infra | 2–3 wks | 2–3 wks |
| Workstream 4 — API endpoints | 3.5–5 wks | 3.5–5 wks |
| Workstream 5 — status lifecycle | 0.5–1 wk | 0.5–3 wks |
| Integration / e2e / testnet buffer | ~1–1.5 wks | ~1–1.5 wks |
| **Total** | **~10–15 engineer-weeks** | **~13–20 engineer-weeks** |
| **Calendar, 2 engineers (protocol ∥ indexer)** | **~6–8 weeks** | **~8–10 weeks** |

---
