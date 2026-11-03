# How YNAB Works — User Journeys

A like-for-like description of YNAB's core flow, in plain language and domain words
only. No architecture or technical choices here — that's yours to decide.

**The core idea:** you can only budget money you actually have. When money arrives it
sits in one pool called **Ready to Assign**. You give every dollar a job by moving it
into **categories**. When you spend, a transaction lowers an account and lowers the
category it came from. One rule always holds:

> **Everything in your accounts = Ready to Assign + everything Available across your categories.**

Each category, each month, has three numbers: **Assigned** (what you put in this
month), **Activity** (what moved out this month), and **Available** (what's left to
spend — this is the number that carries forward).

---

## Journey 1 — A fresh account

**What the user sees first:** an empty budget. No accounts, no categories, Ready to
Assign is $0. The only sensible next step is "Add your first account."

**What they do:**
1. **Add an account** — name, type, and a starting balance (e.g. Checking, $0, with
   income to come). A positive starting balance immediately appears in **Ready to
   Assign** as money waiting for a job.
2. **Create categories** (optionally in groups) — the jobs money can have: Rent,
   Groceries, Fun Money. New categories start empty: Available $0.
3. **Start budgeting** — move money out of Ready to Assign into categories. Each
   assignment lowers Ready to Assign and raises that category's Available. The goal
   is to reach **Ready to Assign = $0** — every dollar has a job.

**Edge cases:** *zero* — an account added with $0 balance leaves nothing to assign
yet. *Negative* — assigning more than Ready to Assign holds pushes it below $0; show
it loudly, don't block. *Empty* — clear empty states for "no accounts" and "no
categories." *Duplicate* — two accounts or two categories with the same name should
be allowed but easy to tell apart. *Very large* — a huge starting balance must still
display and total correctly.

## Journey 2 — After each transaction (the core flow)

There are four kinds of transaction. Each one keeps the rule above true.

| The user records… | Account balance | Category Activity | Category Available | Ready to Assign |
|---|---|---|---|---|
| **Spending** (account + category) | goes down | goes down | goes down | unchanged |
| **Refund** (account + category) | goes up | goes up | goes up | unchanged |
| **Income** (marked "to Ready to Assign") | goes up | — | — | goes up |
| **Transfer** (account → account) | one down, one up | — | — | unchanged |

**What the user sees:** after a spend, the account's balance and the category's
Available both drop by the amount, and the register shows the new line. A **Refund** is
the exact reverse of a spend: money returns to the account and back **into the category
it was spent from**, so that category's Available goes up. (A refund is *not* income —
income has no category and goes to Ready to Assign.) After income, Ready to Assign grows
and they're invited to assign it.

**Editing and deleting** a transaction are **not built for the MVP** (they're on the
roadmap) — in the MVP a transaction is entered once, and a mistake is corrected with a
balancing entry (e.g. a Refund). We still design the roll-up numbers now so edits are
correct later; see *Planning note — edit-safe roll-ups* below.

**Edge cases:** *zero* — **$0 transactions are rejected** (they change nothing and only
add noise). *Negative* — money coming back is the **Refund** action above, not a
negative spend. *Empty* — a spend or refund with no category is flagged; income and
transfers legitimately have no category. *Duplicate* — the same amount/date/payee
entered twice is valid but worth a gentle "looks like a duplicate." *Very large* — a
spend bigger than the account balance takes the account negative; allow it, but flag it.

## Journey 3 — Month end: overspending and rollover

**During the month:** a category can be **overspent** — its Available goes negative
(red). The app's suggestion is to **cover it** by moving money from another category.
A move lowers one category's Available and raises another's; accounts and Ready to
Assign don't change.

**When the month rolls over:**
- Each category's **positive Available carries forward** as its starting Available
  next month.
- **Assigned and Activity reset to $0** for the new month; Available is the running
  number that carries.
- **Leftover Ready to Assign carries forward** to the new month.
- **Cash overspending does not carry as a negative.** Instead, the overspent amount is
  **subtracted from next month's Ready to Assign**, and the category resets to $0. The
  budget as a whole pays for the overspend.

**Edge cases:** *zero* — a category ending at exactly $0 carries nothing, no special
case. *Negative* — the overspend case above. *Empty* — a never-used category just
stays at $0. *Duplicate* — rolling the same month over twice must not double-count;
rollover should be safe to re-run. *Very large* — a large carried balance, or a large
overspend pulled from Ready to Assign, must stay correct.

---

## Planning note — edit-safe roll-ups (not built for MVP)

We won't build transaction **editing or deleting** in the MVP, but we design the numbers
now so it can be added later without a rewrite. Every rolled-up number — a category's
**Activity** and **Available**, **Ready to Assign**, and the amounts **carried into the
next month** — should be **calculated from the underlying transactions and assignments**,
so that changing or removing one transaction produces the correct new totals rather than
leaving a stale running tally.

The hardest case to keep correct is an edit to a transaction in an **already-rolled-over
month**: it changes that month's Available, which changes what carried forward, which
cascades through every later month's Ready to Assign. Design the roll-up so that
re-running it after such a change always lands on the right numbers and is safe to
repeat. Getting this shape right now is the whole reason we flag edits as out-of-scope
rather than ignoring them.

---

## Worked example

Fresh Checking account (starts at $0). Three categories: Rent, Groceries, Fun Money.

| Step | Action | Checking | Ready to Assign | Rent | Groceries | Fun Money |
|---|---|---|---|---|---|---|
| 1 | Income $3,000 (to Ready to Assign) | 3,000 | 3,000 | 0 | 0 | 0 |
| 2 | Assign 1,500 / 1,000 / 500 | 3,000 | **0** | 1,500 | 1,000 | 500 |
| 3 | Spend Rent $1,500 | 1,500 | 0 | 0 | 1,000 | 500 |
| 4 | Spend Groceries $850 | 650 | 0 | 0 | 150 | 500 |
| 5 | Spend Fun $600 (overspent $100) | 50 | 0 | 0 | 150 | **−100** |

End of month: Checking $50 = Ready to Assign $0 + (Rent 0 + Groceries 150 + Fun −100). ✓

**Next month rolls over:**
- Groceries carries its **+$150**; Rent and Fun carry $0 (Fun's −$100 does *not* carry).
- Fun's **$100 cash overspend is deducted from Ready to Assign**.
- New month start: Checking **$50**, Ready to Assign **−$100**, Rent 0, Groceries 150,
  Fun 0. (Still balances: 50 = −100 + 150.) The user must bring in or free up $100 to
  get Ready to Assign back to $0.

---

## What I'm sure of vs. unsure of

**Sure (checked against YNAB's help pages):** the core loop; the "accounts = Ready to
Assign + Available" rule; positive Available carries forward; leftover Ready to Assign
carries forward; **cash overspending is deducted from next month's Ready to Assign and
does not carry as negative.**

**Less sure / deliberately simplified:** the precise wording of in-app prompts, and the
exact rules for the out-of-scope items (editing, credit cards, loans) for when we build
them. (Decided: **$0 transactions are rejected.**) These are open choices, not settled
facts.

**Where this came from:** YNAB's official support pages for month rollover and
overspending, plus general knowledge of the app and third-party reviews for the
feature overview.

## Out of scope for MVP — the roadmap

The MVP is **cash-only**: every account behaves like cash or a bank account you spend
from directly. The capabilities below are deliberately **not** in the MVP. Each is
designed to be added later as its own feature **without changing the cash-only core**,
so deferring them now costs us nothing. The order is a proposed priority — you decide.

1. **Editing & deleting transactions** — correcting or removing a past entry and having
   every roll-up update correctly (see the planning note). Deferred from the MVP, but the
   numbers are designed for it now.
2. **Repeating transactions** (scheduled / recurring) — saved transactions that appear
   automatically on a date (rent, salary, subscriptions). The MVP is manual entry; this
   is pure convenience layered on the same flow.
3. **Targets / goals** — a per-category "assign $X this month / by this date" that
   guides and tracks how much to assign. Assigning already works without it.
4. **Reconciliation** — matching an account to its real-world balance and recording an
   adjustment when they differ. Add once the core numbers are trusted.
5. **Reports & Age of Money** — net worth, income vs. expense, spending by category, and
   the "how old is the money you spend" metric. Reporting sits on top of existing data;
   it adds nothing to the flow.
6. **Credit cards** — the full YNAB mechanic: spending on a credit account moves the
   budgeted money into a dedicated "Credit Card Payment" category, so the cash to pay
   the bill is always set aside. This is genuinely its own feature. The cash-only MVP is
   built so it can be added later without reworking the core.
7. **Loans & tracking accounts** — non-budget accounts (mortgages, car loans,
   investments) that affect net worth but aren't spent from day to day, plus loan payoff
   tracking. A separate concern from the budgeting loop.
8. **Bank syncing / import** — pulling transactions from a bank automatically instead of
   typing them. An integration layer, not core behaviour.

Also later, smaller: multiple budgets, a managed payee list, multiple currencies.

---

## Proposed MVP *(a proposal — you decide what's in)*

The smallest behaviors that complete the three journeys, in build order. Each is meant
to be small enough to build in ~45 minutes or less. **All accounts are cash-only** —
credit cards, loans, and repeating transactions are on the roadmap above, not here.

1. See an **empty budget** and **add an account** with a starting balance; Ready to
   Assign reflects it.
2. **Create category groups and categories** (start empty).
3. **Assign** money from Ready to Assign to a category; both numbers update.
4. Record a **spending** transaction; account balance, category Activity and Available
   update.
5. Record a **Refund** into a category; account, category Activity and Available go back
   up.
6. Record **income** to Ready to Assign; Ready to Assign updates.
7. **Move money** between categories (to cover overspending).
8. **Transfer** between two accounts.
9. **View the budget** for the month: Ready to Assign plus each category's Assigned /
   Activity / Available.
10. **View an account's transactions** (the register).
11. **Roll into a new month:** carry positive Available, carry leftover Ready to Assign,
    and deduct cash overspending from Ready to Assign.

Steps 1–4 already demonstrate the core loop (add account → assign → spend); the rest
complete the journeys. **$0 transactions are rejected; editing/deleting is not in the
MVP** (a mistake is fixed with a balancing entry such as a Refund).
