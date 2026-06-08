# PurpleAir Go Baseline

## Status: Completed

## Context

`purpleair-go` is a small Go client for PurpleAir sensor lookups. Recent
maintenance replaced live-network tests with mocked coverage, added
`SensorWithError`, and strengthened response validation for malformed JSON and
empty sensor result sets.

## Objectives

- Preserve the original `Sensor(sensorID)` compatibility wrapper.
- Provide an error-returning `SensorWithError(sensorID)` path for callers that
  need request, status, JSON, or empty-result errors.
- Keep deterministic tests on mocked HTTP servers instead of the live PurpleAir
  endpoint.
- Keep `make check` as the local formatting, test, and plan-verification gate.
- Record the completed baseline under `docs/plans`.

## Work Completed

- Added mocked client and sensor response tests.
- Added explicit errors for malformed JSON and empty sensor responses.
- Added `make check` and `make verify` wrappers.
- Extended `make verify` to require this canonical completed plan.

## Verification

- `make check`
- `make verify`
- `go test ./...`
- `git diff --check`

## Follow-Up Candidates

- Add examples for `SensorWithError`.
- Document client timeout and endpoint assumptions.
- Consider deprecating `Sensor` after callers can migrate.
