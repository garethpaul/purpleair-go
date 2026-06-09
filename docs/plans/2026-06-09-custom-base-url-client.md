# Custom Base URL Client

## Status: Completed

## Context

Tests inside this package can set the unexported `baseURL` field directly, but
external callers had no public way to point the client at a local proxy, fixture
server, or PurpleAir-compatible endpoint. That made deterministic integration
setups harder without reaching into package internals.

## Goals

- Add a public `NewClientWithBaseURL(baseURL)` constructor.
- Trim custom base URL input.
- Fall back to the default PurpleAir JSON endpoint when the custom base URL is
  blank.
- Preserve existing query parameters when adding the `show` sensor ID.

## Work Completed

- Added `NewClientWithBaseURL` while preserving `NewClient()`.
- Added unit coverage for custom endpoint query preservation and blank fallback.
- Updated README, VISION, and CHANGES.

## Verification

- `gofmt`
- `go test ./...`
- `make check`
- `make verify`
- `git diff --check`
