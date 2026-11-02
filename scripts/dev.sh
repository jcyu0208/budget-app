#!/usr/bin/env bash
#
# Run budget-app locally.
#
# Starts the Go backend, and the frontend too once it exists
# (detected by frontend/package.json). Press Ctrl+C to stop everything.
#
# Usage:
#   ./scripts/dev.sh
#   PORT=9000 ./scripts/dev.sh   # change the backend port (default 8080)
#
set -euo pipefail

# Resolve repo root regardless of where the script is invoked from.
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

BACKEND_PORT="${PORT:-8080}"

pids=()
cleaned=0
cleanup() {
	[ "$cleaned" -eq 1 ] && return
	cleaned=1
	echo ""
	echo "Shutting down..."
	for pid in "${pids[@]:-}"; do
		kill "$pid" 2>/dev/null || true
	done
	wait 2>/dev/null || true
}
trap cleanup INT TERM EXIT

# --- Backend (Go) -----------------------------------------------------------
if ! command -v go >/dev/null 2>&1; then
	echo "error: Go is not installed (need Go 1.27+). See backend/README.md." >&2
	exit 1
fi
# Build first, then run the compiled binary directly. (Running via `go run`
# leaves an orphaned child process on Ctrl+C that keeps holding the port.)
echo "Building backend..."
( cd backend && go build -o bin/dev-server ./cmd/server )
echo "Starting backend on http://localhost:${BACKEND_PORT} ..."
( cd backend && exec env PORT="$BACKEND_PORT" ./bin/dev-server ) &
pids+=($!)

# --- Frontend (optional, auto-enabled once it exists) -----------------------
if [ -f frontend/package.json ]; then
	if command -v npm >/dev/null 2>&1; then
		echo "Starting frontend ..."
		( cd frontend && exec npm run dev ) &
		pids+=($!)
	else
		echo "warning: frontend/package.json found but npm is not installed; skipping frontend." >&2
	fi
else
	echo "No frontend/package.json yet — running backend only."
fi

echo ""
echo "budget-app is running. Press Ctrl+C to stop."
# Wait on all services; if any exits, cleanup() (EXIT trap) stops the rest.
wait
