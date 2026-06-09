# Response Body Size Guard

Status: Completed

## Context

`SensorWithError` reads PurpleAir-compatible HTTP responses before JSON
decoding. The expected single-sensor payload is small, but custom endpoints or
transports could return a very large body and force unnecessary memory use
before the client discovers the response is unusable.

## Objectives

- Bound sensor response body reads before JSON parsing.
- Return an explicit PurpleAir error when the body exceeds the cap.
- Cover the oversized-body path with a deterministic custom transport test.
- Document the guard in README, VISION, SECURITY, and CHANGES.

## Work Completed

- Added a 1 MiB `maxSensorResponseBytes` cap in `SensorWithError`.
- Switched the body read to `io.LimitReader` with one sentinel byte so
  over-limit responses can be detected.
- Added unit coverage for oversized responses using a mocked HTTP transport.
- Updated top-level documentation and security notes for the bounded read.

## Verification

- `gofmt -w sensor.go sensor_test.go`
- `go test ./...`
- `make lint`
- `make test`
- `make build`
- `make verify`
- `make check`
- `scripts/check-baseline.sh`
- `git diff --check`
