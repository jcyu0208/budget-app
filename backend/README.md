# budget-app backend

Go HTTP backend for budget-app (standard library `net/http`).

## Requirements

- Go 1.27+

## Run

```sh
cd backend
go run ./cmd/server
# listens on :8080 (override with PORT, e.g. PORT=9000 go run ./cmd/server)
```

Check it:

```sh
curl localhost:8080/health
# {"status":"ok"}
```

## Test

```sh
cd backend
go test ./...
```

## Layout

- `cmd/server` — entrypoint (`main`).
- `internal/server` — router and handlers.
