# Sensor Process Exit Boundary

Status: Completed

## Problem

The compatibility `Client.Sensor` method calls `log.Fatal` when
`SensorWithError` fails. A library lookup error can therefore terminate the
embedding process, bypassing caller cleanup, recovery, and error policy.

## Scope

- Preserve the public `Sensor(string) *PurpleAir` signature.
- Return `nil` from the compatibility wrapper when the error-returning lookup
  fails instead of exiting the process.
- Keep `SensorWithError` and `SensorWithContext` as the diagnostic APIs.
- Add mutation-sensitive tests, static contracts, and synchronized guidance.
- Do not change HTTP requests, validation, response parsing, result models,
  dependency versions, or live-network behavior.

## Implementation Units

### U1: Remove the process-exit side effect

Update `sensor.go` so the compatibility wrapper delegates to
`SensorWithError`, returns its successful result, and returns `nil` on failure
without logging or terminating the caller.

Test scenarios:

- An invalid sensor ID returns `nil` and the test process continues.
- A valid mocked response still returns the requested sensor data.
- A mutation restoring `log.Fatal`, panic, or a non-nil error fallback is
  rejected.

### U2: Preserve the compatibility contract

Extend `sensor_test.go`, `scripts/check-baseline.sh`, and project guidance with
the named no-process-exit boundary and the preferred error-returning API.

Test scenarios:

- The focused wrapper tests and complete Go/Make gates pass.
- Mutations removing either success or failure coverage, guidance, or completed
  plan evidence are rejected.
- The absolute Makefile path continues to pass outside the checkout.

## Validation

- Run focused tests, formatting, vet, race tests, build, and the full Make gate
  from both repository and external working directories.
- Run isolated hostile mutations for source, tests, guidance, and plan status.
- Audit the exact diff, generated artifacts, credentials, conflict markers,
  binaries, large files, file modes, and whitespace before committing.

## Risks

- Callers that relied on `Sensor` terminating the process will now receive
  `nil`; callers needing failure detail must use `SensorWithError` or
  `SensorWithContext`.
- Mocked tests do not exercise live PurpleAir availability or credentials.
- This change is stacked on PR #6, which must remain open and merge first.

## Work Completed

- Removed fatal logging from the pointer-only compatibility wrapper and return
  `nil` when the delegated lookup fails.
- Added direct failure-continuation and successful mocked-response tests while
  preserving the existing method signature.
- Added a method-scoped static contract and synchronized README, security,
  vision, and changelog guidance.

## Verification Completed

- Repository and external-directory `make check` passed, including gofmt,
  `go vet`, unit tests, `go test -race`, build-through-test, documentation, and
  the fail-closed baseline script.
- All repository and external-directory Make gates passed.
- Seven isolated hostile mutations were rejected across fatal logging,
  non-nil failure fallback, success delegation, failure coverage, success
  coverage, guidance, and completed plan status.
- Exact diff, generated-artifact, credential, conflict-marker, binary,
  large-file, mode, and whitespace audits passed.
- The deterministic tests used mocked HTTP servers and made no live PurpleAir
  request or credentialed call.
