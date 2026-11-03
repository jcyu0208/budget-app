# budget-app — Product Spec (YNAB-style)

A like-for-like model of how YNAB works at its core, told through user journeys,
with the goal of identifying the **essential components for an MVP**. The focus is
the **money flow** — accounts, assigning money, transactions, and how balances
move — not convenience features (bank sync, recurring transactions, reports).

---

## 1. The idea in one paragraph

You can only budget money you **actually have** (not money you expect). When money
arrives it lands in a single pool called **Ready to Assign**. You then give every
dollar a job by **assigning** it to **categories** (Rent, Groceries, …). When you
spend, you record a **transaction** against an account and a category; the
category's remaining balance goes down, and so does the account. The whole app is a
loop: *money in → assign it → spend it → watch the balances → adjust.*

## 2. Mental model & glossary

- **Account** — a real place money sits (Checking, Savings, Cash). For the MVP, all
  accounts are **on-budget** (their money is budgetable). _Tracking accounts_
  (investments, loans) are out of scope (see §9).
- **Ready to Assign (RTA)** — the pool of money that has arrived but hasn't been
  given a job yet. The app constantly nudges this toward **0**.
- **Category** — a job for money (Rent, Groceries). Grouped into **Category Groups**
  (e.g. "Monthly Bills") for organisation.
- **Assigned / Activity / Available** — the three numbers every category has **for a
  given month**:
  - **Assigned** — money you put into the category this month.
  - **Activity** — money that moved out (spending) this month (negative).
  - **Available** — the running balance you can still spend. **It rolls over to next
    month.** `Available = carryover_from_prev_month + Assigned + Activity`.
- **Transaction** — money entering or leaving an account: an **inflow** (+) or
  **outflow** (−), tied to a **payee**, a **category**, a date, and an optional memo.
- **Transfer** — moving money between two of your own accounts (no category).

### The one invariant that must always hold

> **Σ(account balances) = Ready to Assign + Σ(category Available)**

Every operation below is designed to keep this true. If a change ever breaks it,
there's a bug. This identity is the single most important thing to protect in code
and in tests.

## 3. Money representation

Store money as **integer minor units** (cents), never floats. All arithmetic is
integer. (YNAB itself uses "milliunits"; cents are enough for us.) This avoids
rounding drift and makes the invariant in §2 exact.

---

## 4. Primary user journey (fresh account → first spend)

This is the flow the MVP must support, end to end.

### Step 0 — Empty state
A new user has a budget with **no accounts**, **no categories**, and **Ready to
Assign = 0**. The app should guide them to "Add your first account."

### Step 1 — Add the first account
The user creates an account (name + type + **starting balance**, e.g. Checking,
`$1,000.00`).
- The account now has a balance of `$1,000`.
- That starting balance is **income into the budget**, so **Ready to Assign becomes
  `$1,000`**.
- Nothing is assigned yet. The app shows: "You have $1,000 to assign."

### Step 2 — Create categories
The user sets up a few categories (optionally grouped), e.g.
`Monthly Bills → Rent`, `Everyday → Groceries`, `Everyday → Fun Money`.
- New categories start with Assigned `0`, Activity `0`, Available `0`.
- **Ready to Assign is unchanged** — creating a category doesn't move money.

### Step 3 — Give every dollar a job (assign)
The user assigns money from Ready to Assign into categories:
Rent `$600`, Groceries `$300`, Fun Money `$100`.
- Each assignment: `RTA −= amount`, `category.Assigned += amount`,
  `category.Available += amount`.
- After assigning all $1,000, **Ready to Assign = $0**. Goal reached: every dollar
  has a job.
- If they assign more than RTA, RTA goes **negative** — a clear warning state
  ("You've assigned more than you have"). Allowed, but surfaced loudly.

### Step 4 — Spend money (record a transaction)
The user buys groceries for `$80` on the Checking account, category Groceries.
- `account.balance −= $80` → Checking `$920`.
- `Groceries.Activity −= $80`, `Groceries.Available −= $80` → Groceries Available
  `$220`.
- **Ready to Assign is unchanged.**
- Invariant check: accounts `$920` = RTA `$0` + categories (`Rent 600` + `Groceries
  220` + `Fun 100`) `$920`. ✓

### Step 5 — More income arrives
A paycheck of `$2,000` is recorded as an **inflow** to Checking, categorised as
**Ready to Assign** (i.e. "income", not a spending category).
- `account.balance += $2,000` → Checking `$2,920`.
- `RTA += $2,000` → Ready to Assign `$2,000`.
- The user is back to Step 3: assign the new money.

### Step 6 — Overspending (Roll With the Punches)
The user spends `$120` on Fun Money, which only had `$100` available.
- Fun Money Available = `−$20` (allowed, shown in red).
- The app prompts: **cover it** by moving `$20` from another category
  (e.g. Groceries → Fun Money). Moving money between categories:
  `from.Available −= 20`, `to.Available += 20`. Accounts and RTA untouched.
- This is Rule 3: the budget is meant to be edited, not blown.

### Step 7 — Transfer between accounts (no budget impact)
The user moves `$500` Checking → Savings.
- `Checking −= 500`, `Savings += 500`. No category, no Activity, **RTA unchanged**.
- Total account money is unchanged, so the invariant holds automatically.

### Step 8 — A new month begins (rollover)
At the 1st of the next month, each category's positive **Available carries forward**
as next month's starting balance; Assigned and Activity reset to `0` for the new
month. (Handling of *negative* carryover from overspending is a refinement — see §9.)

---

## 5. What happens after each operation (the state-transition table)

| Operation | Account | Category Activity | Category Available | Ready to Assign |
|---|---|---|---|---|
| Add account (starting balance +X) | +X | — | — | **+X** |
| Assign X to a category | — | — | **+X** | **−X** |
| Move X between categories | — | — | −X (from), +X (to) | — |
| Outflow (spend X, categorised) | −X | **−X** | **−X** | — |
| Inflow categorised as Ready to Assign (income +X) | +X | — | — | **+X** |
| Transfer X between two on-budget accounts | −X / +X | — | — | — |
| Delete / edit a transaction | reverse of the above | reverse | reverse | reverse |

**Every row preserves `Σ(accounts) = RTA + Σ(Available)`.** This table is effectively
the spec for the backend and the source of our test cases.

---

## 6. What the user needs on screen (MVP)

Three views cover the whole flow:

1. **Accounts / Register**
   - List of accounts with balances; total of all accounts.
   - Open an account to see its **register**: a list of transactions (date, payee,
     category, memo, amount) with add / edit / delete.
2. **Budget** (the heart of the app)
   - **Ready to Assign** shown prominently at the top (green at 0, red if negative).
   - Category groups → categories, each row showing **Assigned | Activity |
     Available**, with Available colour-coded (positive / zero / negative).
   - Inline edit of **Assigned**; a "move money between categories" action.
   - Month selector (even if we start with the current month only).
3. **Add transaction**
   - Account, date, payee, inflow/outflow amount, category (or "Ready to Assign" for
     income), memo.

## 7. Worked example (doubles as acceptance tests)

Start: no accounts. All values in dollars.

| # | Action | Checking | Savings | RTA | Rent | Groceries | Fun |
|---|---|---|---|---|---|---|---|
| 1 | Add Checking, start $1,000 | 1000 | — | 1000 | 0 | 0 | 0 |
| 2 | Assign 600 / 300 / 100 | 1000 | — | 0 | 600 | 300 | 100 |
| 3 | Spend $80 Groceries | 920 | — | 0 | 600 | 220 | 100 |
| 4 | Income +$2,000 (to RTA) | 2920 | — | 2000 | 600 | 220 | 100 |
| 5 | Spend $120 Fun | 2800 | — | 2000 | 600 | 220 | −20 |
| 6 | Move $20 Groceries→Fun | 2800 | — | 2000 | 600 | 200 | 0 |
| 7 | Transfer $500 →Savings | 2300 | 500 | 2000 | 600 | 200 | 0 |

At every row: `Checking + Savings == RTA + Rent + Groceries + Fun`. (Row 7:
`2300 + 500 = 2000 + 600 + 200 + 0 = 2800`. ✓)

---

## 8. Essential MVP components

**Data entities**
- `Account { id, name, type, balance }` — all on-budget for MVP.
- `CategoryGroup { id, name }`
- `Category { id, group_id, name }`
- `CategoryMonth { category_id, month, assigned, activity, available }` — the
  per-month numbers (or derive `activity`/`available` from transactions + rollover).
- `Transaction { id, account_id, date, payee, category_id(nullable), memo, amount,
  transfer_account_id(nullable) }` — `category_id` is null for income-to-RTA and for
  transfers; amount is signed (− outflow, + inflow).
- **Ready to Assign** — a derived value: `Σ(account balances) − Σ(category Available)`
  (keeping it derived guarantees the invariant can't drift).

**Operations (map directly to backend endpoints)**
- Create / list accounts (with starting balance → RTA).
- Create / list category groups & categories.
- **Assign** money to a category for a month (and un-assign / change).
- **Move** money between categories.
- Create / edit / delete a **transaction** (outflow, income-to-RTA, transfer).
- Read the **budget view** for a month (RTA + each category's three numbers).
- Read an **account register**.

**Invariant to enforce & test everywhere:** §2's identity, after every mutation.

## 9. Explicitly out of scope for the MVP (and why)

Deferring these keeps the core flow clean; none of them are required to demonstrate
the YNAB loop:
- **Recurring / scheduled transactions** — convenience, not flow.
- **Bank syncing / auto-import** — integration; MVP is manual entry.
- **Tracking accounts & net worth** — not budgetable; separate concern.
- **Targets / goals** — guide *how much* to assign, but assigning works without them.
- **Reconciliation** — matching to a real bank balance; add once the core is solid.
- **Reports, Age of Money** — reporting on top of the same data.
- **Credit cards (YNAB's auto payment-category mechanic)** — genuinely complex;
  MVP treats every account as cash-like (spending reduces category Available
  directly). Revisit as its own spec.
- **Negative-balance rollover rules** — in YNAB, cash overspending is pulled from
  next month's Ready to Assign. MVP can start single-month; add the rollover rule
  when we implement months fully.

## 10. Open decisions (to confirm before building)

1. **Single budget** for MVP (one user, one budget), or multi-budget from day one?
   _Recommended: single._
2. **Months**: implement real month rollover now, or start with a single "current
   month" and add rollover next? _Recommended: current month first._
3. **Payees**: free-text string on the transaction (simple) vs. a managed Payee
   entity? _Recommended: free-text string for MVP._
4. **Credit cards**: confirm MVP treats them as cash-like (no auto payment category).

---

*This spec intentionally describes the **flow**. Once we agree on §10, the §5 table
and §7 example become our first test cases, and §8 becomes the first set of backend
models and endpoints.*
