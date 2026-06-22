#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)
TEMP_ROOT=$(mktemp -d "${TMPDIR:-/tmp}/purpleair-go-module-tidy-test-XXXXXX")
trap 'chmod -R u+w "$TEMP_ROOT" 2>/dev/null || :; rm -rf "$TEMP_ROOT"' EXIT HUP INT TERM

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
"$canonical/scripts/check-module-tidy.sh"

ci_cache="$TEMP_ROOT/ci-cache"
mkdir -p "$ci_cache"
(
  cd "$canonical"
  GOMODCACHE="$ci_cache" go test ./...
)
GOMODCACHE="$ci_cache" "$canonical/scripts/check-module-tidy.sh"

complete_cache="$TEMP_ROOT/complete-cache"
mkdir -p "$complete_cache"
(
  cd "$canonical"
  GOMODCACHE="$complete_cache" go mod download
)
GOMODCACHE="$complete_cache" GOPROXY=off "$canonical/scripts/check-module-tidy.sh"

missing_checksum="$TEMP_ROOT/missing-checksum"
cp -R "$canonical" "$missing_checksum"
awk '
  $0 != "gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405 h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM="
' "$missing_checksum/go.sum" >"$missing_checksum/go.sum.new"
mv "$missing_checksum/go.sum.new" "$missing_checksum/go.sum"
expect_failure missing-checksum "$missing_checksum"
if GOMODCACHE="$ci_cache" "$missing_checksum/scripts/check-module-tidy.sh" >"$TEMP_ROOT/missing-checksum-partial-cache.out" 2>&1; then
  printf '%s\n' "Module tidy checker repaired a missing checksum while staging a partial cache." >&2
  exit 1
fi
grep -Fq "Go module metadata must be tidy" "$TEMP_ROOT/missing-checksum-partial-cache.out"

missing_newline="$TEMP_ROOT/missing-newline"
cp -R "$canonical" "$missing_newline"
go_mod_contents=$(cat "$missing_newline/go.mod")
printf '%s' "$go_mod_contents" >"$missing_newline/go.mod"
expect_failure missing-final-newline "$missing_newline"

printf '%s\n' "Module tidy tests passed: canonical metadata accepted with partial, complete, and proxy-disabled caches; missing checksum and final-newline drift rejected."
