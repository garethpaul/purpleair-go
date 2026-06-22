#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)
TEMP_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/purpleair-go-module-tidy-test-XXXXXX")
trap 'rm -rf "$TEMP_ROOT"' EXIT HUP INT TERM

copy_checkout() {
  destination=$1
  mkdir -p "$destination"
  cp -R "$ROOT_DIR"/. "$destination"
  rm -rf "$destination/.git"
}

expect_failure() {
  scenario=$1
  checkout=$2
  if "$checkout/scripts/check-module-tidy.sh" >"$TEMP_ROOT/$scenario.out" 2>&1; then
    printf '%s\n' "Module tidy checker accepted $scenario." >&2
    exit 1
  fi
  grep -Fq "Go module metadata must be tidy" "$TEMP_ROOT/$scenario.out"
}

canonical="$TEMP_ROOT/canonical"
copy_checkout "$canonical"
(
  cd "$canonical"
  HTTP_PROXY=http://127.0.0.1:1 \
    HTTPS_PROXY=http://127.0.0.1:1 \
    ALL_PROXY=http://127.0.0.1:1 \
    NO_PROXY=localhost,127.0.0.1,::1 \
    GOPROXY=off \
    GOSUMDB=off \
    GOWORK=off \
    GOFLAGS='' \
    go mod tidy
)
"$canonical/scripts/check-module-tidy.sh"

missing_checksum="$TEMP_ROOT/missing-checksum"
cp -R "$canonical" "$missing_checksum"
awk '
  $0 != "gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405 h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM="
' "$missing_checksum/go.sum" >"$missing_checksum/go.sum.new"
mv "$missing_checksum/go.sum.new" "$missing_checksum/go.sum"
expect_failure missing-checksum "$missing_checksum"

missing_newline="$TEMP_ROOT/missing-newline"
cp -R "$canonical" "$missing_newline"
go_mod_contents=$(cat "$missing_newline/go.mod")
printf '%s' "$go_mod_contents" >"$missing_newline/go.mod"
expect_failure missing-final-newline "$missing_newline"

printf '%s\n' "Module tidy tests passed: canonical metadata accepted, missing checksum rejected, and final-newline drift rejected."
