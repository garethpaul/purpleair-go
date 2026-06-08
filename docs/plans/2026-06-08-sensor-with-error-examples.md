# SensorWithError Examples

## Status: Completed

## Context

`purpleair-go` now exposes `SensorWithError` so callers can handle request,
status, JSON, blank-input, and empty-result failures without the process exit
behavior of the original `Sensor` compatibility wrapper. The follow-up was to
add examples that document the preferred API while staying deterministic.

## Objectives

- Add executable examples for `SensorWithError`.
- Keep examples off the live PurpleAir endpoint by using `httptest`.
- Cover both a successful mocked sensor response and a blank sensor ID error.
- Keep `make check` as the formatting, test, and completed-plan gate.

## Work Completed

- Added `example_test.go` with `ExampleClient_SensorWithError` and
  `ExampleClient_SensorWithError_error`.
- Updated README, VISION, and CHANGES notes for executable examples.
- Relied on Go's example test runner through the existing `go test ./...`
  verification path.

## Verification

- `gofmt -w *.go`
- `go test ./...`
- `make check`
- `make verify`
- `git diff --check`
