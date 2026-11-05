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

## Testing

```sh
./scripts/test.sh                # run all test types
./scripts/test.sh unit           # unit tests only
./scripts/test.sh integration    # integration tests only (none yet)
./scripts/test.sh e2e            # end-to-end tests only (none yet)
```

Extra arguments pass straight through to the underlying runner, e.g.
`./scripts/test.sh unit -run TestHealthHandler -v`.

Test types (more can be added as the project grows):

- **unit** — fast, no external dependencies. (Go: the default build.)
- **integration** — needs external deps; Go test files tagged `//go:build integration`.
- **e2e** — end-to-end against a running app (not set up yet).

## Merging changes

`main` is protected: changes land via pull request only (no direct/force pushes, no
bypass). A PR merges when **CI is green** and the **branch is up to date with `main`**.

CI ([`.github/workflows/ci.yml`](.github/workflows/ci.yml)) runs `gofmt`, `go vet`,
`go build`, and `go test -race` on every PR; the `backend` check is required. Local hooks
also block direct pushes and scan for secrets (`git config core.hooksPath .githooks`).

## Project conventions

- Operating rules (feature branches, PR-only merges to `main`, test-first, no committed
  secrets): [AGENTS.md](AGENTS.md).
- Decision log: [docs/decisions/](docs/decisions/).
