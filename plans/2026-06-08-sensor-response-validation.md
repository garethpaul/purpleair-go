# Sensor Response Validation

## Status

Completed

## Context

`SensorWithError` now returns transport, status, and JSON parsing errors, but
successful HTTP responses with malformed JSON or no sensor results need explicit
regression coverage. For a sensor lookup API, an empty `results` array should
not look like a successful sensor response to callers.

## Objectives

- Keep live-network calls out of tests.
- Add mocked coverage for malformed PurpleAir JSON responses.
- Return an explicit error when a successful response has no sensor results.
- Document the strengthened `SensorWithError` error contract.

## Verification

- `go test ./...`
- `make verify`
- `git diff --check`
