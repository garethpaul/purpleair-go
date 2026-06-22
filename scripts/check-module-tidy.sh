#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)
TEMP_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/purpleair-go-module-tidy-XXXXXX")
OUTPUT="$TEMP_ROOT/tidy.out"
trap 'chmod -R u+w "$TEMP_ROOT" 2>/dev/null || :; rm -rf "$TEMP_ROOT"' EXIT HUP INT TERM

run_offline_tidy() (
  if [ "$#" -gt 0 ]; then
    GOMODCACHE=$1
    export GOMODCACHE
  fi
  HTTP_PROXY=http://127.0.0.1:1
  HTTPS_PROXY=http://127.0.0.1:1
  ALL_PROXY=http://127.0.0.1:1
  NO_PROXY=localhost,127.0.0.1,::1
  GOPROXY=off
  GOSUMDB=off
  GOWORK=off
  GOFLAGS=''
  export HTTP_PROXY HTTPS_PROXY ALL_PROXY NO_PROXY GOPROXY GOSUMDB GOWORK GOFLAGS
  exec go mod tidy -diff
)

failed=0
last_byte=$(tail -c 1 "$ROOT_DIR/go.mod" | od -An -tu1 | tr -d '[:space:]')
if [ "$last_byte" != "10" ]; then
  printf '%s\n' "go.mod must end with exactly one newline." >&2
  failed=1
fi

if ! (
  cd "$ROOT_DIR"
  run_offline_tidy
) >"$OUTPUT" 2>&1; then
  if grep -Fq "disabled by GOPROXY=off" "$OUTPUT"; then
    MODULE_CACHE="$TEMP_ROOT/module-cache"
    STAGED_MOD="$TEMP_ROOT/staged.mod"
    mkdir -p "$MODULE_CACHE"
    cp "$ROOT_DIR/go.mod" "$STAGED_MOD"
    cp "$ROOT_DIR/go.sum" "$TEMP_ROOT/staged.sum"
    if ! (
      cd "$ROOT_DIR"
      GOMODCACHE="$MODULE_CACHE" GOWORK=off GOFLAGS='' \
        go mod download -modfile="$STAGED_MOD"
    ) >"$TEMP_ROOT/download.out" 2>&1; then
      cat "$TEMP_ROOT/download.out" >&2
      failed=1
    elif ! (
      cd "$ROOT_DIR"
      run_offline_tidy "$MODULE_CACHE"
    ) >"$OUTPUT" 2>&1; then
      cat "$OUTPUT" >&2
      failed=1
    fi
  else
    cat "$OUTPUT" >&2
    failed=1
  fi
fi

if [ "$failed" -ne 0 ]; then
  printf '%s\n' "Go module metadata must be tidy with the complete local module cache and network access disabled." >&2
  exit 1
fi
