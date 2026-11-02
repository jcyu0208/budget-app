#!/usr/bin/env bash
#
# Run budget-app tests.
#
# Usage:
#   ./scripts/test.sh                 # run all test types
#   ./scripts/test.sh unit            # unit tests only
#   ./scripts/test.sh integration     # integration tests only
#   ./scripts/test.sh e2e             # end-to-end tests only
#   ./scripts/test.sh unit -run TestHealthHandler -v   # extra args pass through
#
# Test types (more can be added as the project grows):
#   unit         fast, no external dependencies (Go: the default build)
#   integration  needs external deps; Go test files tagged `//go:build integration`
#   e2e          end-to-end against a running app (not set up yet)
#
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

TYPE="${1:-all}"
shift || true # remaining args pass through to the underlying test runner

usage() {
	sed -n '2,19p' "$0" | sed 's/^#\{1,\} \{0,1\}//'
}

have_go() { command -v go >/dev/null 2>&1; }

run_unit() {
	echo "==> unit tests"
	if have_go; then
		(cd backend && go test ./... "$@")
	else
		echo "   skipped: Go not installed (need Go 1.27+)" >&2
	fi
	if [ -f frontend/package.json ] && command -v npm >/dev/null 2>&1; then
		(cd frontend && npm test)
	fi
}

run_integration() {
	echo "==> integration tests"
	local tagged
	tagged="$(grep -rl '//go:build integration' backend --include='*_test.go' 2>/dev/null || true)"
	if [ -n "$tagged" ] && have_go; then
		(cd backend && go test -tags=integration ./... "$@")
	else
		echo "   none yet (add Go test files tagged \`//go:build integration\`)"
	fi
}

run_e2e() {
	echo "==> e2e tests"
	echo "   none yet"
}

case "$TYPE" in
unit) run_unit "$@" ;;
integration) run_integration "$@" ;;
e2e) run_e2e "$@" ;;
all) run_unit "$@" && run_integration "$@" && run_e2e "$@" ;;
-h | --help | help) usage ;;
*)
	echo "unknown test type: $TYPE" >&2
	echo "usage: ./scripts/test.sh [all|unit|integration|e2e] [extra args]" >&2
	exit 2
	;;
esac
