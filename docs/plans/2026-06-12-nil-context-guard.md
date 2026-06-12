# Nil Context Guard

## Status: Completed

## Context

`SensorWithContext` accepts a caller-provided `context.Context` and forwards it
to `http.NewRequestWithContext`. Go requires that context to be non-nil and
panics otherwise. A library caller that accidentally passes nil therefore
crashes its process instead of receiving the method's documented error result.

## Priority

This is a public API panic boundary. It is more urgent than internal cleanup
because one malformed caller input can terminate a service, while the guard is
small, deterministic, and fully testable without credentials or network access.

## Objectives

- Return a stable error when `SensorWithContext` receives a nil context.
- Reject nil before request construction or HTTP client execution.
- Add a regression proving the call does not panic and sends no request.
- Protect implementation, test, documentation, and completed-plan evidence in
  the fail-closed baseline checker.
- Preserve sensor ID validation, cancellation, timeout, response-size,
  response-identity, and successful request behavior.

## Implementation Units

### U1. Public API guard

**Files:** `sensor.go`

**Goal:** Convert the nil-context panic into an ordinary library error.

**Approach:** Check the context before request construction and return a
package-prefixed stable error. Keep existing sensor ID validation order and all
non-nil request behavior unchanged.

### U2. Network-free regression

**Files:** `sensor_test.go`

**Goal:** Prove nil context cannot panic or reach the HTTP transport.

**Approach:** Use a counting transport or server-free client stub, call the
public method with a valid sensor ID and nil context, and assert the exact error,
nil result, and zero requests.

### U3. Durable contracts and documentation

**Files:** `scripts/check-baseline.sh`, `README.md`, `SECURITY.md`, `VISION.md`,
`CHANGES.md`, `AGENTS.md`, `docs/plans/2026-06-12-nil-context-guard.md`

**Goal:** Keep the public panic boundary and its evidence from regressing.

**Approach:** Register the completed plan, require the guard and regression
contracts, and document the nil-context behavior alongside existing caller
cancellation guidance.

## Verification

- focused nil-context regression
- `gofmt` and `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `go build ./...`
- `make lint`, `make test`, `make race`, `make vet`, `make build`, `make verify`,
  and `make check`
- network-isolated exact-file validation on supported Go 1.25 and Go 1.26
- hostile guard, test, plan, and documentation mutations
- workflow YAML parsing, shell syntax, `go mod verify`, and `git diff --check`

## Work Completed

- Added a nil-context guard after existing sensor ID validation and before HTTP
  request construction, returning `purpleair: context is required`.
- Added a public-method regression with a counting transport that proves nil
  context returns a nil result, does not panic, and performs zero requests. A
  combined malformed-ID/nil-context case preserves the existing ID error order.
- Registered this completed plan and protected the guard, exact error,
  regression name, no-request assertion, and README behavior in the baseline.
- Updated contributor, security, vision, README, and changelog documentation.

## Verification Results

- `go test ./... -run TestSensorWithContextRejectsNilContext -count=1` passed.
- The complete local `go test ./...`, `go test -race ./...`, and `go vet ./...`
  paths passed before plan completion.
- The intentionally incomplete plan was rejected by `make check` with the
  exact completed-status error.
- `make lint`, `make test`, `make race`, `make vet`, `make build`, `make verify`,
  and `make check` passed on the completed worktree.
- Network-isolated exact-file `make check` passed on Go 1.25.11 and Go 1.26.4
  after module preloading.
- All 10 hostile guard, error, ordering, regression, plan, and README mutations
  were rejected.
- `go mod verify`, POSIX shell and Dash syntax, workflow YAML parsing, and
  `git diff --check` passed.

## Boundaries

- Do not add retries, live PurpleAir calls, dependencies, or credentialed tests.
- Do not change exported signatures or successful response semantics.
- Do not reorder existing sensor ID errors for calls that also pass nil context.
- Preserve the existing remediation PR and canonical exact-head evidence.
