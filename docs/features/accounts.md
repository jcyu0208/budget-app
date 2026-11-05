# Feature: Accounts + Ready to Assign

First MVP slice (top of Journey 1 in `docs/specs.md`): create a cash account with a
starting balance, and read **Ready to Assign**. Also establishes the core / ports /
adapters skeleton.

## Behaviour

- **Create account** (`POST /accounts`): `{ name, kind, starting_balance }` → the new
  account `{ id, name, kind, balance }`.
- **Read budget** (`GET /budget`): `{ ready_to_assign, accounts: [{id, name, kind, balance}] }`.
- Starting with no accounts: `ready_to_assign = 0`, `accounts = []`.

## Rules

- `name` required (non-empty after trim). Duplicate names allowed.
- `kind` must be `cash` (covers checking/savings/cash/debit). `credit`/`loan` rejected → roadmap.
- `starting_balance`: whole cents (`int64`); may be **zero** or **negative**; very large OK.
- **Current balance is derived**, never stored: `starting_balance + Σ(transactions)` (no
  transactions exist yet → balance = `starting_balance`).
- **Ready to Assign is derived**, read-only: `Σ(current balances) − Σ(category Available)`
  → `Σ(current balances)` here (no categories yet).
- Invariant (always): `Σ(account balances) = Ready to Assign + Σ(category Available)`.

## Edge cases

- **zero** starting balance → account at 0, no transaction created, RTA unchanged.
- **negative** starting balance → allowed; RTA can go negative (shown loudly).
- **empty** → no accounts: RTA 0, empty list.
- **duplicate** names → both accounts exist, distinguished by id.
- **very large** → int64 cents; totals stay correct.
- invalid `kind` / empty `name` → `400`.

## Decisions (from the interview)

- Single implicit budget.
- `kind = cash` only for MVP; starting balance is a **stored input** (account attribute),
  not a transaction; a $0 start creates no transaction (no clash with "$0 txns rejected").
- Numbers derived from records, never stored; whole cents.
- Records are dated, but **month aggregation is deferred** to the assigning ticket
  (resolves open item #1).
- Storage: one `Account` concept + repository **port**; in-memory adapter now, physical
  SQL table layout deferred (base + extension tables later).

## Architecture

- **Core** `internal/budget` (no HTTP/SQL): `Account`, `Kind`, `Money` (int64 cents),
  ports `AccountRepository` + `IDGenerator`, and `Service` (`CreateAccount`, `Budget`).
- **Adapter (storage)** `internal/budget/memory`: in-memory repo + sequential id generator.
- **Adapter (HTTP)** `internal/server`: `NewRouter(svc)` with `POST /accounts`, `GET /budget`,
  `GET /health`; wired in `cmd/server/main.go`.

## Test plan

- **Unit (core):** create-account validation; balance derivation; `ReadyToAssign` = Σ balances.
- **Contract (port):** shared `AccountRepository` suite run against every adapter (in-memory now).
- **Functional (API → in-memory):** drive HTTP via `httptest`, assert through the read side
  (`GET /budget`), check the invariant after every step; cover all edge cases above.

## Links

- Core: `backend/internal/budget/` (`account.go`, `ports.go`, `service.go`).
- Adapters: `backend/internal/budget/memory/` (storage + id gen); `backend/internal/server/server.go` (HTTP); wired in `backend/cmd/server/main.go`.
- Tests: `backend/internal/budget/service_test.go` (unit), `backend/internal/budget/budgettest/contract.go` + `backend/internal/budget/memory/repository_test.go` (contract), `backend/internal/server/accounts_functional_test.go` (functional + invariant).

## API

- `POST /accounts` — `{ name, kind: "cash", starting_balance }` (cents) → `201 { id, name, kind, balance }`; `400` on empty name / non-cash kind / bad JSON.
- `GET /budget` → `200 { ready_to_assign, accounts: [{ id, name, kind, balance }] }`.
