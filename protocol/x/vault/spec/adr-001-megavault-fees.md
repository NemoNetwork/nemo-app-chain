# ADR-001 — Megavault operator fees

**Status:** Accepted (M2b deliverable)
**Date:** 2026-08-04
**Milestone:** [M2b](../MILESTONES.md#m2--protocol-operator-params-and-the-performance-fee-engine), gates M2c and M7
**Supersedes:** nothing. **Upstream equivalent:** none — `dydxprotocol/v4-chain` has no fee mechanism at any version.

---

## 1. Context

The megavault must charge the operator's economics on-chain. Nothing in upstream can be
ported: upstream `OperatorParams` is exactly `{ operator address, { name, description } }`.
Every decision below is ours.

Two facts constrain the design:

1. **The chain is live and already holds deposits made under a 0%-fee regime.** Introducing
   a fee at an upgrade height means existing depositors begin paying a fee they did not
   agree to. A naive high-water mark charges them *retroactively* on gains earned before
   the fee existed. See [§7](#7-introducing-fees-to-a-live-chain).
2. **Megavault shares are fungible and pooled.** There is no per-owner lot tracking in
   state — `OwnerShares` is one number per address. Any design requiring per-depositor cost
   basis would need new state proportional to deposits-per-owner, and would make deposits
   and withdrawals unboundedly expensive.

### 1.1 Resolving "fee vs. profit share: one charge or two?"

[MILESTONES.md](../MILESTONES.md) lists this as an open product question. The evidence in
the repository answers it: the frontend renders **three distinct zero-valued fields** —
*operator fee*, *minimum operator share*, and *profit share* ([spec.md](../spec.md) §1/§2).
A UI would not show one charge twice.

So: **two distinct charges**, which is also the conventional hedge-fund structure the field
names imply.

| Frontend field | This ADR | Charged on | Conventional name |
|---|---|---|---|
| Operator fee | `operator_fee_ppm` | NAV, pro-rata over time | management fee |
| Profit share | `profit_share_ppm` | gains above the high-water mark | performance fee |
| Minimum operator share | `min_operator_share_ppm` | — (a constraint, not a charge) | operator skin-in-the-game |

**This decision is cheap to reverse.** The two charges are independent parameters computed
in separate steps. If product later says they are one charge, set one of them to zero — no
code change, no migration. That is the reason for designing two rather than asking and
waiting.

---

## 2. Decisions

| # | Question | Decision |
|---|---|---|
| D1 | Fee basis | **Global NAV-per-share high-water mark**, not per-depositor cost basis |
| D2 | Charge structure | **Two charges**: a time-based operator fee and a HWM-gated profit share |
| D3 | Crystallization trigger | **Both** — accrue on an epoch, and force an accrual immediately before any withdrawal |
| D4 | Payment mechanism | **Mint shares to the operator** (dilution), never a quote-quantum transfer |
| D5 | Ordering | Operator fee first, then profit share, so the profit share is charged net of the operator fee |
| D6 | Minimum operator share | Enforced by **blocking operator withdrawals**, never by blocking user deposits |
| D7 | Parameter authority | Fee params are **gov-only**; the operator cannot raise its own fee |
| D8 | HWM at upgrade | **Seeded at the NAV per share observed at upgrade height** |

### D1 — Global NAV-per-share high-water mark

The HWM is a single number: the highest NAV per share the megavault has ever reached, *after*
any fee dilution. A profit share is charged only on the excess of current NAV per share over
that mark.

Rejected alternative — per-depositor cost basis: it is fairer to a depositor who buys in at a
peak (they pay no performance fee until they are personally in profit), but it requires a
lot-tracking table keyed by (owner, deposit) and turns withdrawal into a loop over lots.
Given fungible pooled shares, the global HWM is the only design that keeps deposit and
withdrawal O(1).

**Known consequence, accepted:** a depositor who enters above the HWM pays a profit share on
the recovery back to their own entry price. This is standard for pooled vaults and is why
D8 matters — see [§7](#7-introducing-fees-to-a-live-chain).

### D2/D5 — Two charges, ordered

Let `E` = megavault equity in quote quantums, `S` = total shares, `NAV = E / S`.

**Step 1 — operator fee (time-based).** For elapsed time `Δt` seconds since the last accrual:

```
operatorFeeQuantums = E × operator_fee_ppm / 1e6 × Δt / SECONDS_PER_YEAR
```

`operator_fee_ppm` is therefore an **annualized** rate. Charged whether or not the vault is
in profit — that is what distinguishes it from the profit share.

**Step 2 — profit share (HWM-gated).** After step 1 has diluted, recompute `NAV`. If
`NAV > HWM`:

```
profitQuantums      = (NAV − HWM) × S
profitShareQuantums = profitQuantums × profit_share_ppm / 1e6
```

Charging in this order means the profit share is levied on gains *net of* the operator fee,
which is the conventional treatment and strictly favours depositors over the reverse order.

### D4 — Payment by minting shares

To transfer `F` quote quantums of value to the operator without moving funds, mint

```
sharesToMint = S × F / (E − F)
```

so that after minting, the operator's new shares are worth exactly `F`:
`E × sharesToMint / (S + sharesToMint) = F`.

Rejected alternative — transferring quote quantums to the operator's subaccount: every
accrual becomes a subaccount transfer that can fail (insufficient free collateral, the
funds are deployed in sub-vaults), so a fee accrual could fail in the EndBlocker. Minting
cannot fail. It also keeps the operator's economics automatically invested in the vault,
which interacts well with D6.

`E − F` is guaranteed positive because `F` is clamped to a fraction of `E` (see
[§4](#4-edge-cases-that-must-be-handled)).

### D3 — Accrue on an epoch, crystallize on withdrawal

Accruing every block is wasteful and rounds to zero at realistic rates. Accruing *only* on
withdrawal lets a depositor who exits just before an accrual escape their share.

So: accrue on a fixed epoch (recommended: hourly, reusing `x/epochs`), **and** force an
accrual as the first step of `WithdrawFromMegavault`, before redemption value is computed.
Because dilution applies to all shares equally, an accrual immediately before a withdrawal
*is* the exiting depositor paying their pro-rata share — no separate per-depositor
settlement is needed.

### D6 — Minimum operator share

`min_operator_share_ppm` is a floor on `operatorShares / totalShares`. Enforce it by
**rejecting an operator withdrawal that would breach the floor**.

Rejected alternative — blocking user deposits when the ratio is too low: a user deposit
dilutes the operator's ratio, so this would punish depositors for the operator's failure to
co-invest, and would let the operator halt deposits by withdrawing. Blocking the operator's
own exit puts the constraint on the party it is meant to bind.

**Deliberately not enforced against dilution from user deposits.** If the floor were a hard
invariant, a large deposit would have to be rejected. The floor binds the operator's
*actions* only.

### D7 — Fee parameters are gov-only

Fee fields are appended to `OperatorParams`, which is set by `MsgUpdateOperatorParams` — an
**internal** message, gov-proposal-only. The operator therefore cannot raise its own fee.

Contrast with the deposit cap (M2a), which lives in `MegavaultParams` and *is*
operator-settable. The dividing line: **risk knobs the operator must be able to pull
quickly are operator-gated; anything that moves value from depositors to the operator is
gov-gated.**

---

## 3. State and parameters

### 3.1 Parameters — appended to `OperatorParams` (gov-only)

```proto
message OperatorParams {
  string operator = 1;
  OperatorMetadata metadata = 2;

  // Annualized management fee on megavault NAV, in parts-per-million.
  uint32 operator_fee_ppm = 3;
  // Share of gains above the high-water mark, in parts-per-million.
  uint32 profit_share_ppm = 4;
  // Minimum fraction of total shares the operator must retain, in ppm.
  uint32 min_operator_share_ppm = 5;
}
```

Appending at fields 3–5 is wire-compatible; existing clients decoding `OperatorParams` are
unaffected, and a chain upgraded from pre-fee state reads all three as zero — i.e. **the fee
engine is inert until governance turns it on**. That is deliberate and is mitigation 3 of
[§7](#7-introducing-fees-to-a-live-chain), available for free.

Validation:

- each ppm value ≤ 1,000,000
- `profit_share_ppm` ≤ 1,000,000 — a profit share above 100% would take more than the profit
- `operator_fee_ppm` — recommend a governance-level sanity ceiling, but no protocol cap
  beyond 1,000,000 (100%/yr), which is already absurd

### 3.2 New state

| Key | Type | Meaning |
|---|---|---|
| `MegavaultFeeState` | `FeeState` | HWM and last-accrual timestamp |

```proto
message FeeState {
  // Highest NAV per share ever reached after fee dilution, scaled by 1e18.
  bytes high_water_mark_nav_per_share = 1 [(gogoproto.customtype) = "...SerializableInt"];
  // Unix seconds of the last completed accrual.
  int64 last_accrual_time = 2;
}
```

NAV per share is a ratio of two big integers and must not be stored lossily. Scaling by
`1e18` keeps ~18 significant decimal digits, far beyond quote-quantum precision.

All new state lives under the existing `vault` store key as a new prefix, so **no
`StoreUpgrades` entry is needed** — consistent with the M7 note.

---

## 4. Edge cases that must be handled

Each of these must have a test in M2c.

| Case | Required behaviour |
|---|---|
| `S == 0` (no shares) | Skip accrual entirely; NAV per share is undefined |
| `E <= 0` | Skip accrual; never mint against non-positive equity |
| `Δt <= 0` (same block, or clock skew) | Skip; never accrue negative time |
| Computed fee rounds to 0 | Skip, but **do not** advance `last_accrual_time` — otherwise a low rate accrues nothing forever |
| Computed fee ≥ `E` | Clamp to a fraction of `E` (recommend ≤ 50% in any one accrual); `E − F` must stay positive |
| `sharesToMint` rounds to 0 | Skip the mint, leave the HWM unchanged |
| `NAV <= HWM` | No profit share. HWM unchanged (it is a *high*-water mark) |
| Operator address changes | Already-minted shares stay with the old address — they are ordinary shares. Governance must be aware |
| Fee params set to 0 mid-life | Accrual becomes a no-op; the HWM is still maintained so re-enabling does not retroactively charge |
| Deactivated vaults | Excluded from megavault equity already (M1c), so they are excluded from the fee base for free |

**The rounding-to-zero case is the subtle one.** If `last_accrual_time` advances on an
accrual that minted nothing, a fee rate low enough to round to zero per epoch accrues
*nothing ever*, no matter how much time passes. Leaving the timestamp unadvanced lets `Δt`
grow until the fee is representable.

---

## 5. Invariants (M2c test obligations)

1. **HWM is monotonic non-decreasing.** No code path lowers it.
2. **Fees never exceed profit.** `profitShareQuantums ≤ profitQuantums` for all inputs.
3. **No double-charge.** Accrual followed immediately by a second accrual with `Δt = 0`
   mints nothing.
4. **Shares reconcile.** `totalShares == Σ ownerShares` after every accrual.
5. **Withdrawal is net of fees.** Redemption value computed after the forced accrual is
   strictly less than or equal to the value computed before it, for any positive fee.
6. **A depositor who enters and exits with no NAV change and zero operator fee pays
   nothing.** Guards against a profit share leaking on flat performance.
7. **Operator floor holds.** After any successful operator withdrawal,
   `operatorShares / totalShares ≥ min_operator_share_ppm`.

---

## 6. Interaction with the M1c withdrawal path

`WithdrawFromMegavault` gains one step at the top:

```
1. AccrueFees(ctx)                     <-- new; crystallizes for the exiting depositor
2. compute redemption value            <-- unchanged, now operates on post-dilution shares
3. apply withdrawal slippage           <-- unchanged
4. burn shares, transfer quantums      <-- unchanged
```

Nothing downstream needs to know fees exist: the accrual has already moved value by changing
the share count, and every later step reads `totalShares` fresh.

The same applies to `MsgDepositToMegavault` — accrue first, so an incoming depositor does not
buy shares at a NAV that still contains an unaccrued fee liability. **This is a correctness
requirement, not an optimization:** without it, a deposit landing just before a large accrual
dilutes the new depositor for fees earned before they arrived.

---

## 7. Introducing fees to a live chain

This is the largest risk in M2, created by the combination of "a real fee engine at launch"
and "the chain has live state". Mitigations, in the order they should be applied:

1. **Seed the HWM at the NAV per share observed at upgrade height (D8).** Only post-upgrade
   gains are ever charged. Cheap, and the obviously correct default. **The upgrade handler
   must do this — a HWM of zero would charge a profit share on the entire history of the
   vault.**
2. **Withdrawals ship before fees.** Depositors who object must be able to exit fee-free.
   M1c (withdrawals) is already complete and M2c is not, so this ordering holds — it is now
   a hard constraint, not a preference.
3. **Launch the parameters at zero and raise them by a separate governance action.** This
   comes for free from §3.1: a chain upgraded from pre-fee state reads all fee params as
   zero. It decouples "the mechanism exists" from "the mechanism charges", and gives an
   announcement window between the two.

Applying all three is recommended, and costs nothing beyond the upgrade-handler work in M7
that is already planned.

---

## 8. What this ADR deliberately does not decide

- **The actual parameter values.** A product and governance decision. The protocol default
  is zero for all three.
- **The accrual epoch length.** Recommended hourly; it is a constant, trivially changed, and
  does not affect the design. Anything from per-block to daily works.
- **Whether the operator's minted shares are subject to the M1b share lock.** Argument for:
  it aligns the operator with depositors. Argument against: it complicates the D6 floor
  check, which would then need to distinguish locked from unlocked operator shares.
  Recommend **not** locking them in v1 and revisiting if the floor proves insufficient.
- **Fee display in the indexer.** M5's surface. The chain exposes the params and the HWM via
  the params query; how the page renders "fees paid to date" is an indexer question, and the
  natural data source is the delta in operator shares over time.
