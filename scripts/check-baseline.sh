#!/usr/bin/env sh
set -eu

ROOT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
README="$ROOT_DIR/README.md"
MAKEFILE="$ROOT_DIR/Makefile"
GITIGNORE="$ROOT_DIR/.gitignore"
DOCS_PLANS="$ROOT_DIR/docs/plans"
WORKFLOW="$ROOT_DIR/.github/workflows/check.yml"

require_file() {
  path=$1
  if [ ! -f "$ROOT_DIR/$path" ]; then
    printf '%s\n' "Required file is missing: $path" >&2
    exit 1
  fi
}

for path in \
  ".github/workflows/check.yml" \
  ".gitignore" \
  "CHANGES.md" \
  "Makefile" \
  "README.md" \
  "SECURITY.md" \
  "VISION.md" \
  "AGENTS.md" \
  "client.go" \
  "client_test.go" \
  "go.mod" \
  "go.sum" \
  "results.go" \
  "sensor.go" \
  "sensor_test.go" \
  "plans/2026-06-12-001-fix-sensor-result-id-validation-plan.md" \
  "docs/plans/2026-06-08-purpleair-go-baseline.md" \
  "docs/plans/2026-06-09-scripted-baseline-check.md" \
  "docs/plans/2026-06-10-ci-baseline.md" \
  "docs/plans/2026-06-10-hosted-go-validation.md" \
  "docs/plans/2026-06-12-default-http-timeout-boundary.md" \
  "docs/plans/2026-06-12-sensor-response-identity.md" \
  "docs/plans/2026-06-13-response-error-context-and-body-close.md" \
  "docs/plans/2026-06-13-response-content-length-preflight.md" \
  "docs/plans/2026-06-14-location-independent-make.md" \
  "docs/plans/2026-06-16-sensor-process-exit-boundary.md" \
  "docs/plans/2026-06-17-active-stack-nil-context-guard.md" \
  "docs/plans/2026-06-21-safe-make-root.md" \
  "docs/plans/2026-06-25-default-redirect-policy.md" \
  "docs/plans/2026-06-25-sensor-caller-migration.md" \
  "scripts/check-module-tidy.sh" \
  "scripts/test-module-tidy.sh" \
  "scripts/test-makefile-root.sh" \
  "scripts/check-baseline.sh"; do
  require_file "$path"
done

if ! grep -Fq "// Deprecated: Use SensorWithError" "$ROOT_DIR/sensor.go"; then
  printf '%s\n' "Sensor compatibility wrapper must carry standard Go deprecation guidance." >&2
  exit 1
fi

unexpected_sensor_callers=$(grep -RIn --include='*.go' '\.Sensor(' "$ROOT_DIR" |
  grep -v '/sensor_test.go:' || true)
if [ -n "$unexpected_sensor_callers" ]; then
  printf '%s\n' "Repository Go callers must use SensorWithError or SensorWithContext:" >&2
  printf '%s\n' "$unexpected_sensor_callers" >&2
  exit 1
fi

for document in "$README" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md" "$ROOT_DIR/AGENTS.md"; do
  if ! grep -Fq '`SensorWithError` is the preferred default' "$document"; then
    printf '%s\n' "$document must preserve Sensor caller migration guidance." >&2
    exit 1
  fi
done

if grep -Fq -- '- Migrate callers from `Sensor` to `SensorWithError`' "$ROOT_DIR/VISION.md"; then
  printf '%s\n' "VISION.md must remove the completed Sensor caller migration priority." >&2
  exit 1
fi

SENSOR_MIGRATION_PLAN="$ROOT_DIR/docs/plans/2026-06-25-sensor-caller-migration.md"
if ! grep -Fq "Status: Completed" "$SENSOR_MIGRATION_PLAN" ||
  ! grep -Fq "make check" "$SENSOR_MIGRATION_PLAN"; then
  printf '%s\n' "Sensor caller migration plan must record completed status and verification." >&2
  exit 1
fi

if ! grep -Fq "if ctx == nil" "$ROOT_DIR/sensor.go" ||
  ! grep -Fq "purpleair: context is required" "$ROOT_DIR/sensor.go"; then
  printf '%s\n' "SensorWithContext must reject nil context before request construction." >&2
  exit 1
fi

sensor_id_line=$(grep -nF "requestedSensorID, parseErr := strconv.Atoi(sensorId)" "$ROOT_DIR/sensor.go" | cut -d: -f1)
nil_context_line=$(grep -nF "if ctx == nil" "$ROOT_DIR/sensor.go" | cut -d: -f1)
request_line=$(grep -nF "http.NewRequestWithContext" "$ROOT_DIR/sensor.go" | cut -d: -f1)
if [ -z "$sensor_id_line" ] || [ -z "$nil_context_line" ] || [ -z "$request_line" ] ||
  [ "$sensor_id_line" -ge "$nil_context_line" ] || [ "$nil_context_line" -ge "$request_line" ]; then
  printf '%s\n' "Sensor ID and nil context validation must precede request construction." >&2
  exit 1
fi

for test_contract in \
  "TestSensorWithContextRejectsNilContext" \
  'assert.EqualError(t, err, "purpleair: context is required")' \
  "nil context must fail before HTTP requests" \
  "sensor id validation must remain before nil context validation"; do
  if ! grep -Fq "$test_contract" "$ROOT_DIR/sensor_test.go"; then
    printf '%s\n' "Sensor tests must preserve nil-context contract: $test_contract" >&2
    exit 1
  fi
done

for document in "$README" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md"; do
  if ! grep -Fq "active-stack nil context guard" "$document"; then
    printf '%s\n' "$document must document the active-stack nil context guard." >&2
    exit 1
  fi
done

NIL_CONTEXT_PLAN="$ROOT_DIR/docs/plans/2026-06-17-active-stack-nil-context-guard.md"
if ! grep -Fq "Status: Completed" "$NIL_CONTEXT_PLAN" ||
  ! grep -Fq "make check" "$NIL_CONTEXT_PLAN"; then
  printf '%s\n' "Active-stack nil context plan must record completed status and verification." >&2
  exit 1
fi

SENSOR_WRAPPER=$(sed -n '/^func (c \*Client) Sensor(/,/^}/p' "$ROOT_DIR/sensor.go")
if printf '%s\n' "$SENSOR_WRAPPER" | grep -Eq 'log\.Fatal|panic\(' ||
  ! printf '%s\n' "$SENSOR_WRAPPER" | grep -Fq 'pa, err := c.SensorWithError(sensorId)' ||
  ! printf '%s\n' "$SENSOR_WRAPPER" | grep -Fq 'if err != nil {' ||
  ! printf '%s\n' "$SENSOR_WRAPPER" | grep -Fq 'return nil' ||
  ! printf '%s\n' "$SENSOR_WRAPPER" | grep -Fq 'return pa'; then
  printf '%s\n' "Sensor compatibility wrapper must return nil without exiting on lookup errors." >&2
  exit 1
fi

for test_contract in \
  "TestSensorReturnsNilInsteadOfExitingOnError" \
  'assert.Nil(t, client.Sensor(" "))' \
  "TestSensorReturnsDataOnSuccess" \
  'assert.NotNil(t, sensor)' \
  'assert.Equal(t, 17937, sensor.Results[0].ID)'; do
  if ! grep -Fq "$test_contract" "$ROOT_DIR/sensor_test.go"; then
    printf '%s\n' "Sensor compatibility tests must preserve: $test_contract" >&2
    exit 1
  fi
done

for document in "$README" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md"; do
  if ! grep -Fq "Sensor process exit boundary" "$document"; then
    printf '%s\n' "$document must document the Sensor process exit boundary." >&2
    exit 1
  fi
done

SENSOR_EXIT_PLAN="$ROOT_DIR/docs/plans/2026-06-16-sensor-process-exit-boundary.md"
SENSOR_EXIT_VERIFICATION=$(sed -n '/^## Verification Completed$/,$p' "$SENSOR_EXIT_PLAN")
if ! grep -Fq "Status: Completed" "$SENSOR_EXIT_PLAN" ||
  ! printf '%s\n' "$SENSOR_EXIT_VERIFICATION" | grep -Fq "All repository and external-directory Make gates passed" ||
  ! printf '%s\n' "$SENSOR_EXIT_VERIFICATION" | grep -Fq "Seven isolated hostile mutations were rejected" ||
  ! printf '%s\n' "$SENSOR_EXIT_VERIFICATION" | grep -Fq "go test -race" ||
  printf '%s\n' "$SENSOR_EXIT_VERIFICATION" | grep -Eiq '\b(pending|todo|tbd|not run)\b'; then
  printf '%s\n' "Sensor process exit boundary plan must record completed verification." >&2
  exit 1
fi

for evidence in \
  "Go 1.25.11" \
  "Go 1.26.4" \
  "unrelated directory" \
  "hostile mutations rejected"; do
  if ! grep -Fq "$evidence" "$ROOT_DIR/docs/plans/2026-06-14-location-independent-make.md"; then
    printf '%s\n' "Location-independent Make plan must preserve verification evidence: $evidence" >&2
    exit 1
  fi
done

if ! grep -Fq "res.ContentLength > maxSensorResponseBytes" "$ROOT_DIR/sensor.go"; then
  printf '%s\n' "Sensor responses must reject oversized declared Content-Length values before reading." >&2
  exit 1
fi

for test_contract in \
  "TestSensorWithErrorRejectsDeclaredOversizedBodiesBeforeReading" \
  "ContentLength: maxSensorResponseBytes + 1" \
  "assert.Equal(t, 0, reader.reads)" \
  '"declared oversized body must close without reading"' \
  "TestSensorWithErrorReadsBodiesDeclaredAtLimit" \
  "ContentLength: maxSensorResponseBytes" \
  "assert.Greater(t, reader.reads, 0)" \
  '"declared exact-limit body must close after reading"' \
  "ContentLength: -1"; do
  if ! grep -Fq "$test_contract" "$ROOT_DIR/sensor_test.go"; then
    printf '%s\n' "Declared response length tests must preserve: $test_contract" >&2
    exit 1
  fi
done

for document in "$README" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md"; do
  if ! grep -Fq "declared Content-Length" "$document"; then
    printf '%s\n' "$document must document declared Content-Length preflight rejection." >&2
    exit 1
  fi
done

for evidence in \
  "Go 1.25.11" \
  "Go 1.26.4" \
  'Canonical `make check` passed' \
  "Eight hostile mutations were rejected" \
  '`go mod verify`' \
  '`git diff --check`' \
  "secret, captured-prompt, and generated-artifact scans"; do
  if ! grep -Fq "$evidence" "$ROOT_DIR/docs/plans/2026-06-13-response-content-length-preflight.md"; then
    printf '%s\n' "Response Content-Length preflight plan must preserve verification evidence: $evidence" >&2
    exit 1
  fi
done

if ! grep -Fq 'purpleair: read response body: %w' "$ROOT_DIR/sensor.go" ||
  ! grep -Fq 'purpleair: decode response body: %w' "$ROOT_DIR/sensor.go"; then
  printf '%s\n' "Sensor response read and decode failures must preserve wrapped PurpleAir context." >&2
  exit 1
fi

for test_contract in \
  "TestSensorWithErrorWrapsResponseReadErrors" \
  "TestSensorWithErrorWrapsJSONDecodeErrors" \
  "TestSensorWithErrorClosesResponseBodies" \
  "errors.Is(err, readErr)" \
  "errors.As(err, &syntaxErr)" \
  "assert.True(t, body.closed)"; do
  if ! grep -Fq "$test_contract" "$ROOT_DIR/sensor_test.go"; then
    printf '%s\n' "Sensor response lifecycle tests must preserve: $test_contract" >&2
    exit 1
  fi
done

if ! grep -Fq "defaultHTTPTimeout = 30 * time.Second" "$ROOT_DIR/client.go" ||
  ! grep -Fq "TestZeroValueClientUsesDefaultTimeout" "$ROOT_DIR/client_test.go" ||
  ! grep -Fq "TestNilClientUsesDefaultTimeout" "$ROOT_DIR/client_test.go" ||
  ! grep -Fq "TestClientPreservesCallerProvidedHTTPTimeout" "$ROOT_DIR/client_test.go"; then
  printf '%s\n' "Client timeout tests must preserve the 30-second default and caller overrides." >&2
  exit 1
fi

for redirect_contract in \
  "CheckRedirect: rejectRedirect" \
  "return http.ErrUseLastResponse"; do
  if ! grep -Fq "$redirect_contract" "$ROOT_DIR/client.go"; then
    printf '%s\n' "Default HTTP client must preserve redirect rejection: $redirect_contract" >&2
    exit 1
  fi
done

for redirect_test in \
  "TestDefaultHTTPClientRejectsRedirects" \
  "TestSensorWithErrorRejectsRedirectsBeforeFollowing" \
  'assert.EqualError(t, err, "purpleair: unexpected status 302")' \
  "assert.Equal(t, 0, destinationRequests)"; do
  if ! grep -Fq "$redirect_test" "$ROOT_DIR/client_test.go" "$ROOT_DIR/sensor_test.go"; then
    printf '%s\n' "Redirect tests must preserve: $redirect_test" >&2
    exit 1
  fi
done

for document in "$README" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md"; do
  if ! grep -Fiq "legacy" "$document" || ! grep -Fiq "redirect" "$document"; then
    printf '%s\n' "$document must document the legacy endpoint and redirect boundary." >&2
    exit 1
  fi
done

for document in "$README" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md"; do
  if ! grep -Fq "30-second" "$document"; then
    printf '%s\n' "$document must document the 30-second default HTTP timeout." >&2
    exit 1
  fi
done

for document in "$README" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md"; do
  if ! grep -Fq "response bodies are closed" "$document"; then
    printf '%s\n' "$document must document that response bodies are closed." >&2
    exit 1
  fi
done

if ! grep -Fq "result.ID <= 0" "$ROOT_DIR/sensor.go" ||
  ! grep -Fq "result %d has invalid sensor id %d" "$ROOT_DIR/sensor.go"; then
  printf '%s\n' "Sensor responses must reject non-positive result IDs." >&2
  exit 1
fi

if ! grep -Fq "TestSensorWithErrorRejectsInvalidResultIDs" "$ROOT_DIR/sensor_test.go" ||
  ! grep -Fq "TestSensorWithErrorAcceptsMultipleValidResultIDs" "$ROOT_DIR/sensor_test.go"; then
  printf '%s\n' "Sensor tests must cover invalid and multiple valid result IDs." >&2
  exit 1
fi

if ! grep -Fq "for _, digit := range sensorId" "$ROOT_DIR/sensor.go" ||
  ! grep -Fq "digit < '0' || digit > '9'" "$ROOT_DIR/sensor.go" ||
  ! grep -Fq "strconv.Atoi(sensorId)" "$ROOT_DIR/sensor.go" ||
  ! grep -Fq "sensor id must be a positive integer" "$ROOT_DIR/sensor.go" ||
  ! grep -Fq "response does not include requested sensor" "$ROOT_DIR/sensor.go"; then
  printf '%s\n' "Sensor requests and responses must preserve requested sensor identity." >&2
  exit 1
fi

for test_name in \
  "TestSensorWithErrorRejectsInvalidRequestedSensorIDs" \
  "TestSensorWithErrorRejectsMismatchedResponseSensorIDs" \
  "TestSensorWithErrorAcceptsMultipleValidResultIDs"; do
  if ! grep -Fq "$test_name" "$ROOT_DIR/sensor_test.go"; then
    printf '%s\n' "Sensor tests must preserve $test_name." >&2
    exit 1
  fi
done

if ! grep -Fq "invalid sensor IDs must fail before HTTP requests" "$ROOT_DIR/sensor_test.go"; then
  printf '%s\n' "Invalid sensor ID tests must prove no HTTP request is made." >&2
  exit 1
fi

if ! grep -Fq '"+1"' "$ROOT_DIR/sensor_test.go" ||
  ! grep -Fq '"１２"' "$ROOT_DIR/sensor_test.go"; then
  printf '%s\n' "Invalid sensor ID tests must reject signed and non-ASCII forms." >&2
  exit 1
fi

DECIMAL_SENSOR_PLAN="$ROOT_DIR/docs/plans/2026-06-14-decimal-sensor-id-validation.md"
if ! grep -Fq "Status: Completed" "$DECIMAL_SENSOR_PLAN" ||
  ! grep -Fq "signed and non-ASCII forms" "$DECIMAL_SENSOR_PLAN" ||
  ! grep -Fq "make check" "$DECIMAL_SENSOR_PLAN"; then
  printf '%s\n' "Decimal sensor ID plan must record completed status and verification." >&2
  exit 1
fi

if ! grep -Fq "Sensor Result ID Validation" "$ROOT_DIR/plans/2026-06-12-001-fix-sensor-result-id-validation-plan.md" ||
  ! grep -Fq "make check" "$ROOT_DIR/plans/2026-06-12-001-fix-sensor-result-id-validation-plan.md"; then
  printf '%s\n' "Sensor result ID validation plan must document repository verification." >&2
  exit 1
fi

for document in "$README" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md"; do
  if ! grep -Fq "non-positive sensor IDs" "$document"; then
    printf '%s\n' "$document must document non-positive sensor ID rejection." >&2
    exit 1
  fi
done

for document in "$README" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md"; do
  if ! grep -Fq "requested sensor identity" "$document"; then
    printf '%s\n' "$document must document requested sensor identity validation." >&2
    exit 1
  fi
done

for document in "$README" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md"; do
  if ! grep -Fq "ASCII decimal" "$document"; then
    printf '%s\n' "$document must document ASCII decimal requested sensor IDs." >&2
    exit 1
  fi
done

if ! grep -Fq "scripts/check-baseline.sh" "$MAKEFILE"; then
  printf '%s\n' "Makefile must run scripts/check-baseline.sh from make check." >&2
  exit 1
fi

for make_contract in \
  'override SHELL := /bin/sh' \
  'override .SHELLFLAGS := -c' \
  '$(error MAKEFILES must be empty; repository verification requires this Makefile to be loaded alone)' \
  'ifneq ($(origin MAKEFILE_LIST),file)' \
  '$(error MAKEFILE_LIST must not be overridden)' \
  'override REPO_ROOT := $(shell path=' \
  'export REPO_ROOT' \
  '/usr/bin/dirname' \
  '/bin/pwd -P' \
  '@cd "$$REPO_ROOT" && for plan in docs/plans/*.md; do \' \
  'cd "$$REPO_ROOT" && test -z "$$(gofmt -l *.go)"' \
  'cd "$$REPO_ROOT" && go vet ./...' \
  'cd "$$REPO_ROOT" && go test ./...' \
  'cd "$$REPO_ROOT" && go test -race ./...' \
  'cd "$$REPO_ROOT" && scripts/test-makefile-root.sh' \
  'cd "$$REPO_ROOT" && scripts/check-baseline.sh'; do
  if ! grep -Fq "$make_contract" "$MAKEFILE"; then
    printf '%s\n' "Makefile must preserve rooted recipe: $make_contract" >&2
    exit 1
  fi
done

for target in "docs:" "fmt:" "lint:" "vet:" "test:" "race:" "build:" "root-test:" "verify:" "check:"; do
  if ! grep -Fq "$target" "$MAKEFILE"; then
    printf '%s\n' "Makefile must expose the $target gate." >&2
    exit 1
  fi
done

if ! grep -Fq "go vet ./..." "$MAKEFILE"; then
  printf '%s\n' "Makefile must run go vet ./... from make vet." >&2
  exit 1
fi

if ! grep -Fq "go test -race ./..." "$MAKEFILE"; then
  printf '%s\n' "Makefile must run go test -race ./... from make race." >&2
  exit 1
fi

if ! grep -Fq "verify: lint vet test race build docs root-test" "$MAKEFILE"; then
  printf '%s\n' "make verify must include the vet and race gates." >&2
  exit 1
fi

for root_contract in \
  'PurpleAir Go' \
  '70 executed target/authority cases' \
  'hostile backticks blocked' \
  'dollar paths failed closed' \
  '1 MAKEFILES preload rejection' \
  '2 MAKEFILE_LIST rejection cases' \
  'MAKEFILE_LIST must not be overridden'; do
  if ! grep -Fq "$root_contract" "$ROOT_DIR/scripts/test-makefile-root.sh"; then
    printf '%s\n' "Makefile root test must preserve: $root_contract" >&2
    exit 1
  fi
done

for root_evidence in \
  'Status: Completed' \
  'nine pre-existing public Make targets plus the root regression gate' \
  '70 executed target and authority cases' \
  'Hostile checkout backticks were blocked and dollar-substitution paths failed closed' \
  '`MAKEFILES`, `SHELL`, and `.SHELLFLAGS` authority were covered' \
  'Command-line and environment `MAKEFILE_LIST` overrides failed closed' \
  'make check'; do
  if ! grep -Fq "$root_evidence" "$ROOT_DIR/docs/plans/2026-06-21-safe-make-root.md"; then
    printf '%s\n' "safe Make root plan must preserve: $root_evidence" >&2
    exit 1
  fi
done

for documented in "go test ./..." "go test -race ./..." "go vet ./..." "make race" "make vet" "make check" "scripts/check-baseline.sh"; do
  if ! grep -Fq "$documented" "$README"; then
    printf '%s\n' "README must document $documented." >&2
    exit 1
  fi
done

expected_workflow=$(mktemp)
trap 'rm -f "$expected_workflow"' EXIT HUP INT TERM
cat >"$expected_workflow" <<'EOF'
name: Check

on:
  push:
    branches: [master]
  pull_request:
  workflow_dispatch:

permissions:
  contents: read

concurrency:
  group: check-${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  go:
    name: Go ${{ matrix.go }} verification
    runs-on: ubuntu-24.04
    timeout-minutes: 10
    strategy:
      fail-fast: false
      matrix:
        go: ["1.25.11", "1.26.4"]
    steps:
      - name: Check out repository
        uses: actions/checkout@df4cb1c069e1874edd31b4311f1884172cec0e10 # v6.0.3
        with:
          persist-credentials: false
      - name: Set up Go
        uses: actions/setup-go@4a3601121dd01d1626a1e23e37211e3254c1c06c # v6.4.0
        with:
          go-version: ${{ matrix.go }}
          cache: true
      - name: Run verification
        run: make check
EOF

if ! cmp -s "$expected_workflow" "$WORKFLOW"; then
  printf '%s\n' "Hosted validation workflow must match the reviewed credential-free contract." >&2
  exit 1
fi

for documented in "GitHub Actions" "no-live-network" "checkout credentials"; do
  if ! grep -Fq "$documented" "$README"; then
    printf '%s\n' "README must document hosted validation: $documented." >&2
    exit 1
  fi
done

for guidance in "make check" "go test -race ./..." "API keys" "live-network"; do
  if ! grep -Fq "$guidance" "$ROOT_DIR/AGENTS.md"; then
    printf '%s\n' "AGENTS.md must preserve contributor guidance: $guidance." >&2
    exit 1
  fi
done

for coordinate_contract in \
  'Lat *float64 `json:"Lat"`' \
  'Lon *float64 `json:"Lon"`' \
  'result %d is missing coordinates'; do
  if ! grep -Fq "$coordinate_contract" "$ROOT_DIR/sensor.go"; then
    printf '%s\n' "Sensor decoder must preserve required coordinate contract: $coordinate_contract" >&2
    exit 1
  fi
done

for coordinate_test in \
  "TestSensorWithErrorRejectsMissingCoordinates" \
  '"missing latitude"' \
  '"missing longitude"' \
  '"later missing result"'; do
  if ! grep -Fq "$coordinate_test" "$ROOT_DIR/sensor_test.go"; then
    printf '%s\n' "Sensor tests must preserve missing-coordinate coverage: $coordinate_test" >&2
    exit 1
  fi
done

for document in "$README" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md"; do
  if ! grep -Fq "explicit" "$document" || ! grep -Fq "coordinate" "$document"; then
    printf '%s\n' "$document must document explicit sensor coordinates." >&2
    exit 1
  fi
done
for module_line in \
  "module github.com/garethpaul/purpleair-go" \
  "go 1.13" \
  "github.com/stretchr/testify v1.5.1"; do
  if ! grep -Fq "$module_line" "$ROOT_DIR/go.mod"; then
    printf '%s\n' "go.mod must keep module baseline: $module_line" >&2
    exit 1
  fi
done

"$ROOT_DIR/scripts/check-module-tidy.sh"
"$ROOT_DIR/scripts/test-module-tidy.sh"

for ignored in "/bin/" "/dist/" "/build/" "*.test" "*.out" ".env" ".env.*" ".idea/" ".vscode/" "*.iml"; do
  if ! grep -Fq "$ignored" "$GITIGNORE"; then
    printf '%s\n' ".gitignore must include $ignored" >&2
    exit 1
  fi
done

if ! tracked_local=$(git -C "$ROOT_DIR" ls-files '.env' '.env.*' '.idea' '.vscode' '*.iml'); then
  printf '%s\n' "Baseline must be able to inspect tracked secret and editor metadata paths." >&2
  exit 1
fi
if [ -n "$tracked_local" ]; then
  printf '%s\n%s\n' "Local secrets or editor metadata must not be tracked:" "$tracked_local" >&2
  exit 1
fi

found_plan=0
for plan in "$DOCS_PLANS"/*.md; do
  [ -e "$plan" ] || continue
  found_plan=1
  status_summary=$(awk '
    /^(## )?Status:/ {
      status_count++
      if ($0 == "Status: Completed" || $0 == "## Status: Completed") {
        completed_count++
      }
    }
    END {
      printf "%d:%d", status_count, completed_count
    }
  ' "$plan")
  if [ "$status_summary" != "1:1" ]; then
    printf '%s\n' "$plan must record exactly one completed status." >&2
    exit 1
  fi
  if ! grep -Fq "make check" "$plan"; then
    printf '%s\n' "$plan must document make check verification." >&2
    exit 1
  fi
done

if [ "$found_plan" -eq 0 ]; then
  printf '%s\n' "docs/plans must contain completed markdown plans." >&2
  exit 1
fi
