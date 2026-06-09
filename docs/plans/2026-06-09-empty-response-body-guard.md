# Empty Response Body Guard

Status: Completed

## Context

`SensorWithError` already returned errors for blank sensor IDs, request
failures, non-2xx responses, malformed JSON, and empty result sets. A custom
transport could still return a successful HTTP response with no payload, causing
the client to pass an empty response body into JSON decoding.

## Objectives

- Return an explicit error when a successful HTTP response has no body content.
- Add mocked transport coverage that does not call the live PurpleAir endpoint.
- Keep the compatibility `Sensor` wrapper behavior unchanged.
- Document the completed guard in README, SECURITY, VISION, and CHANGES.

## Work Completed

- Added nil and empty-body checks in `SensorWithError`.
- Added a mocked `RoundTripper` regression test for empty response bodies.
- Updated maintenance documentation for the new error path.

## Verification

- `go test ./...`
- `make check`
- `make verify`
- `git diff --check`
