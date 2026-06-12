# Sensor Context Cancellation

Status: Completed

## Context

`SensorWithError` always created a request without a caller context, leaving
applications dependent on the client's five-minute timeout. Callers could not
cancel work during shutdown or apply a shorter per-request deadline while
retaining the library's response validation and error context.

## Objectives

- Add a context-aware sensor lookup method without breaking existing callers.
- Propagate cancellation and deadlines to the HTTP request.
- Preserve `errors.Is` checks through PurpleAir-specific request errors.
- Keep deterministic tests free of live network calls.

## Work Completed

- Added `SensorWithContext(ctx, sensorID)` as the context-aware request path.
- Kept `SensorWithError(sensorID)` as a compatibility wrapper using
  `context.Background()`.
- Added a custom-transport regression test that verifies context values reach
  the request and `context.Canceled` remains discoverable through the error.
- Documented cancellation and deadline behavior in README, SECURITY, VISION,
  and CHANGES.

## Verification

- `gofmt -w sensor.go sensor_test.go`
- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `make check`
- Replaced the caller context with `context.Background()` in a mutation check
  and confirmed the context propagation test failed.
- `git diff --check`
