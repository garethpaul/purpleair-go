# Mocked Sensor Tests

## Status

Completed

## Context

`purpleair-go` had tests, but `sensor_test.go` called the live PurpleAir
endpoint. In this environment that request timed out, and the library path used
`log.Fatal`, which exited the process on request failures.

## Objectives

- Replace live-network tests with deterministic `httptest` coverage.
- Route sensor requests through the client's configured base URL and HTTP client.
- Add an error-returning API while preserving the original `Sensor(sensorId)`
  compatibility wrapper.
- Add a local verification target that checks formatting and tests.

## Verification

- `go test ./...`
- `make verify`
- `git diff --check`

## Follow-Up Candidates

- Add examples for `SensorWithError`.
- Add tests for malformed JSON and empty sensor result sets.
- Consider deprecating `Sensor` after callers can migrate.
