# Custom Base URL Validation

## Status: Completed

## Context

`NewClientWithBaseURL` trimmed blank values, but it still accepted malformed,
hostless, or non-HTTP endpoint strings. Those values later failed at request
time instead of falling back to the default PurpleAir JSON endpoint when the
client was constructed.

## Goals

- Reject malformed custom base URLs in `NewClientWithBaseURL`.
- Accept only absolute `http` and `https` URLs with a host.
- Preserve existing query parameter handling for valid fixture or proxy
  endpoints.
- Cover invalid values with deterministic unit tests.

## Work Completed

- Added constructor tests for malformed, non-HTTP, and hostless custom base
  URLs.
- Added an internal supported-base-URL helper in `client.go`.
- Updated README, VISION, and CHANGES with the validation behavior.

## Verification

- `gofmt -w client.go client_test.go`
- `go test ./...`
- `make check`
- `make verify`
- `git diff --check`
