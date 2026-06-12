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
  "docs/plans/2026-06-12-nil-context-guard.md" \
  "scripts/check-baseline.sh"; do
  require_file "$path"
done

if ! grep -Fq "defaultHTTPTimeout = 30 * time.Second" "$ROOT_DIR/client.go" ||
  ! grep -Fq "TestZeroValueClientUsesDefaultTimeout" "$ROOT_DIR/client_test.go" ||
  ! grep -Fq "TestNilClientUsesDefaultTimeout" "$ROOT_DIR/client_test.go" ||
  ! grep -Fq "TestClientPreservesCallerProvidedHTTPTimeout" "$ROOT_DIR/client_test.go"; then
  printf '%s\n' "Client timeout tests must preserve the 30-second default and caller overrides." >&2
  exit 1
fi

for document in "$README" "$ROOT_DIR/SECURITY.md" "$ROOT_DIR/VISION.md" "$ROOT_DIR/CHANGES.md"; do
  if ! grep -Fq "30-second" "$document"; then
    printf '%s\n' "$document must document the 30-second default HTTP timeout." >&2
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

if ! grep -Fq "if ctx == nil" "$ROOT_DIR/sensor.go" ||
  ! grep -Fq "purpleair: context is required" "$ROOT_DIR/sensor.go"; then
  printf '%s\n' "SensorWithContext must reject nil context before request construction." >&2
  exit 1
fi

nil_context_line=$(grep -nF "if ctx == nil" "$ROOT_DIR/sensor.go" | cut -d: -f1)
request_line=$(grep -nF "http.NewRequestWithContext" "$ROOT_DIR/sensor.go" | cut -d: -f1)
if [ -z "$nil_context_line" ] || [ -z "$request_line" ] || [ "$nil_context_line" -ge "$request_line" ]; then
  printf '%s\n' "Nil context must be rejected before request construction." >&2
  exit 1
fi

if ! grep -Fq "TestSensorWithContextRejectsNilContext" "$ROOT_DIR/sensor_test.go" ||
  ! grep -Fq "nil context must fail before HTTP requests" "$ROOT_DIR/sensor_test.go" ||
  ! grep -Fq "sensor id validation must remain before nil context validation" "$ROOT_DIR/sensor_test.go"; then
  printf '%s\n' "Sensor tests must prove nil context returns an error without HTTP requests." >&2
  exit 1
fi

if ! grep -Fq "strconv.Atoi(sensorId)" "$ROOT_DIR/sensor.go" ||
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

if ! grep -Fq "scripts/check-baseline.sh" "$MAKEFILE"; then
  printf '%s\n' "Makefile must run scripts/check-baseline.sh from make check." >&2
  exit 1
fi

for target in "docs:" "fmt:" "lint:" "vet:" "test:" "race:" "build:" "verify:" "check:"; do
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

if ! grep -Fq "verify: lint vet test race build docs" "$MAKEFILE"; then
  printf '%s\n' "make verify must include the vet and race gates." >&2
  exit 1
fi

for documented in "go test ./..." "go test -race ./..." "go vet ./..." "make race" "make vet" "make check" "scripts/check-baseline.sh"; do
  if ! grep -Fq "$documented" "$README"; then
    printf '%s\n' "README must document $documented." >&2
    exit 1
  fi
done

if ! grep -Fq "nil context" "$README"; then
  printf '%s\n' "README must document the SensorWithContext nil-context error." >&2
  exit 1
fi

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
for module_line in \
  "module github.com/garethpaul/purpleair-go" \
  "go 1.13" \
  "github.com/stretchr/testify v1.5.1"; do
  if ! grep -Fq "$module_line" "$ROOT_DIR/go.mod"; then
    printf '%s\n' "go.mod must keep module baseline: $module_line" >&2
    exit 1
  fi
done

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
