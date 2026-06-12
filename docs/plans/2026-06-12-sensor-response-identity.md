# Sensor Response Identity Validation

Status: Completed

## Context

The client validates that returned result IDs are positive, but it accepts a
payload whose results all belong to different sensors than the requested ID.
That can silently associate stale, proxied, or malformed data with the caller's
requested sensor.

## Priority

Sensor identity is part of the API contract. Returning data for another sensor
is more dangerous than returning an explicit error because callers may publish
or act on measurements under the wrong location.

## Objectives

- Require requested sensor IDs to be positive decimal integers.
- Require at least one result ID to match the requested sensor.
- Continue accepting additional positive result IDs for paired sensor payloads.
- Add mocked regressions for malformed requests, mismatched responses, and
  matching multi-result responses.
- Protect the implementation, tests, documentation, and completed plan through
  the repository baseline.

## Verification

- `gofmt -w sensor.go sensor_test.go`
- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `make check`
- `git diff --check`

## Work Completed

- Requested IDs are trimmed, parsed as decimal integers, and rejected before
  network access when they are malformed or non-positive.
- Every decoded result must retain a positive ID, and at least one result must
  match the requested sensor identity.
- Additional positive IDs remain accepted for paired sensor responses.
- The baseline protects the implementation, three regression test contracts,
  canonical documentation, and this completed plan.

## Verification Results

- `gofmt -w sensor.go sensor_test.go` completed with no remaining format diff.
- `go test ./...` passed.
- Invalid requested IDs were verified to fail before the configured HTTP
  transport is called.
- `go vet ./...` passed.
- `go test -race ./...` passed.
- `make fmt`, `make lint`, `make vet`, `make test`, `make race`, `make build`,
  `make docs`, `make verify`, and `make check` passed on Go 1.25.3.
- `make check` passed in the official Go 1.25.11 and Go 1.26.4 containers
  after marking the read-only bind mount as a Git safe directory.
- All 10 hostile plan, request parsing, response identity, regression, and
  documentation mutations were rejected.
- `git diff --check` passed.
