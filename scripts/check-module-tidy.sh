#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)
OUTPUT=$(mktemp "${TMPDIR:-/tmp}/purpleair-go-module-tidy-XXXXXX")
trap 'rm -f "$OUTPUT"' EXIT HUP INT TERM

failed=0
last_byte=$(tail -c 1 "$ROOT_DIR/go.mod" | od -An -tu1 | tr -d '[:space:]')
if [ "$last_byte" != "10" ]; then
  printf '%s\n' "go.mod must end with exactly one newline." >&2
  failed=1
fi

if (
  cd "$ROOT_DIR"
  HTTP_PROXY=http://127.0.0.1:1 \
    HTTPS_PROXY=http://127.0.0.1:1 \
    ALL_PROXY=http://127.0.0.1:1 \
    NO_PROXY=localhost,127.0.0.1,::1 \
    GOPROXY=off \
    GOSUMDB=off \
    GOWORK=off \
    GOFLAGS='' \
    go mod tidy -diff
) >"$OUTPUT" 2>&1; then
  :
else
  cat "$OUTPUT" >&2
  failed=1
fi

if [ "$failed" -ne 0 ]; then
  printf '%s\n' "Go module metadata must be tidy with the complete local module cache and network access disabled." >&2
  exit 1
fi
