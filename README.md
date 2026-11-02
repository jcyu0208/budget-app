# budget-app

Personal budgeting app — a Go backend (`backend/`) and a frontend (`frontend/`, not built yet).

## Running locally

Prerequisites:

- **Go 1.27+** (backend)
- One-time per clone — enable the shared git hooks: `git config core.hooksPath .githooks`

Start the app:

```sh
./scripts/dev.sh
```

This starts the backend on <http://localhost:8080> (override with `PORT=9000 ./scripts/dev.sh`).
It also starts the frontend automatically once `frontend/package.json` exists; until then it
runs the backend only. Press **Ctrl+C** to stop everything.

Verify the backend is up:

```sh
curl localhost:8080/health   # {"status":"ok"}
```

### Running one part on its own

- **Backend only:** see [backend/README.md](backend/README.md).

## Project conventions

- Operating rules (feature branches, PR-only merges to `main`, test-first, no committed
  secrets): [AGENTS.md](AGENTS.md).
- Decision log: [docs/decisions/](docs/decisions/).
