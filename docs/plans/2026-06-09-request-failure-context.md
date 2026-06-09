# Request Failure Context

Status: Completed

## Context

`SensorWithError` returned transport failures directly from `http.Client.Do`.
Those errors are useful, but without PurpleAir-specific context callers must
infer that the failure happened while fetching sensor data. A custom transport
that returns a nil response also exercises this request-failure path through
Go's HTTP client.

## Objectives

- Wrap HTTP request failures with `purpleair: request failed` while preserving
  the original error for callers that use `errors.Is` or `errors.As`.
- Add a defensive nil-response guard before reading response fields.
- Add `make lint` and `make build` aliases for the repository-standard Go gate
  sequence.
- Document the completed guard in README, SECURITY, VISION, and CHANGES.

## Work Completed

- Wrapped `http.Client.Do` errors in `SensorWithError` with request context.
- Added a nil-response regression test using a custom RoundTripper.
- Added a defensive nil-response guard before checking the response body.
- Added Makefile `lint` and `build` aliases and routed `verify` through lint,
  test, build, and docs.

## Verification

- Red `go test ./...` with the request-failure context regression.
- `gofmt -w sensor.go sensor_test.go`
- `make lint`
- `make test`
- `make build`
- `make docs`
- `make check`
- `git diff --check`
