# Client Input and Timeout Guards

## Status: Completed

## Context

`purpleair-go` has deterministic mocked tests for sensor lookup success,
transport failures, status errors, malformed JSON, and empty result sets. Two
client safety gaps remained: blank sensor IDs still reached request
construction, and a zero-value `Client` used `http.DefaultClient`, which has no
timeout.

## Objectives

- Reject blank sensor IDs before issuing HTTP requests.
- Keep zero-value and nil clients on the same default timeout baseline as
  `NewClient()`.
- Preserve the existing `NewClient()`, `Sensor`, and `SensorWithError` API
  shapes.
- Make the docs plan verification target cover every completed plan under
  `docs/plans`.

## Work Completed

- Added shared default base URL and HTTP timeout helpers.
- Updated `SensorWithError` to trim and reject blank sensor IDs.
- Made nil or zero-value clients use a default HTTP client with the standard
  five-minute timeout.
- Added tests for blank sensor IDs, zero-value clients, and nil client helpers.
- Updated README, VISION, CHANGES, and the docs verification target.

## Verification

- `gofmt -w *.go`
- `go test ./...`
- `make check`
- `make verify`
- `git diff --check`

## Follow-Up Candidates

- Add examples for `SensorWithError`.
- Document endpoint availability and current PurpleAir API assumptions.
