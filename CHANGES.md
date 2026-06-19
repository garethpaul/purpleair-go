# Changes

## 2026-06-19

- Redacted request URLs from transport error strings while preserving the
  original error chain for cancellation, deadlines, and diagnostics.
- Returned response-body close failures after otherwise successful lookups and
  preserved primary error precedence on failed lookups.
- Bounded incremental result decoding to 1,024 entries to prevent compact JSON
  arrays from amplifying into excessive `Result` allocations.
- Added property-style secret-redaction coverage and concurrent client reuse
  coverage under the race detector without live PurpleAir requests.

## 2026-06-17

- Restored the active-stack nil context guard with a stable error, preserved
  sensor-ID validation order, and no-request regression coverage.

## 2026-06-16

- Added the Sensor process exit boundary so the pointer-only compatibility
  lookup returns `nil` on errors instead of terminating the embedding process.

## 2026-06-13

- Wrapped response read and JSON decode failures with PurpleAir-specific phase
  context while preserving underlying Go errors for `errors.Is` and
  `errors.As`.
- Added deterministic coverage proving all non-nil response bodies are closed
  across status, read, size, blank, decode, validation, and success paths.
- Rejected an oversized declared Content-Length before the first response body
  read while preserving the bounded fallback for absent or inaccurate lengths.

## 2026-06-12

- Reduced the default total sensor HTTP timeout from five minutes to a
  30-second boundary for constructor, nil, and zero-value clients.
- Added exact fallback and caller-provided client preservation coverage while
  retaining `SensorWithContext` deadline overrides.
- Rejected missing or non-positive sensor IDs so malformed upstream records
  cannot be returned as valid data.
- Added mocked coverage for invalid identities and multiple positive sensor ID
  results.
- Required positive decimal request IDs and rejected responses that do not
  preserve the requested sensor identity in at least one result.
- Required requested sensor IDs to use ASCII decimal digits, rejecting signed
  and non-ASCII spellings before network access.

## 2026-06-10

- Added credential-free, pinned hosted verification on Go 1.25.11 and Go
  1.26.4 that runs the local no-live-network `make check` baseline.
- Added `make race` and wired `go test -race ./...` into the canonical check.
- Added a `make vet` static analysis gate and wired it into `make verify` and
  `make check`.
- Added `SensorWithContext` so callers can cancel sensor requests or apply
  deadlines while preserving wrapped context errors.

## 2026-06-09

- Added a sensor response body size guard before JSON parsing.
- Wrapped HTTP request failures with PurpleAir-specific context while
  preserving the original transport error.
- Added `scripts/check-baseline.sh` and local metadata coverage for required
  files, Go module metadata, completed plan metadata, and verification docs.
- Added a nil HTTP response guard so custom transports return an error instead
  of panicking in `SensorWithError`.
- Added `make lint` and `make build` aliases to match the shared repository
  verification workflow.

## 2026-06-08

- Added executable `SensorWithError` examples for mocked success and blank
  sensor ID error paths.
- Added `NewClientWithBaseURL` for local proxies, fixture servers, and alternate
  PurpleAir-compatible endpoints.
- Added custom base URL validation so malformed or non-HTTP endpoints fall back
  to the default PurpleAir JSON endpoint.
- Added custom base URL credential validation so URLs with embedded userinfo
  fall back to the default PurpleAir JSON endpoint.
- Added custom base URL fragment validation so local-only tokens or notes do
  not hide in endpoint strings.
- Added an empty response-body guard so `SensorWithError` returns an error
  instead of decoding an empty body.
- Added blank sensor ID validation and default timeout coverage for zero-value
  clients.
- Added `SensorWithError` so callers can handle request, response, and JSON parsing failures without a process exit.
- Updated `Sensor` to keep the original API while delegating to the error-returning implementation.
- Replaced the live-network sensor test with mocked HTTP coverage for successful and failed responses.
- Added `make verify` for Go formatting checks and the full test suite.
- Added `make check` as the shared repository verification alias.
- Added mocked coverage for malformed JSON and empty sensor result responses, and made empty result sets return an explicit `SensorWithError` error.
- Added canonical `docs/plans` coverage and made `make verify` require the
  completed baseline plan.
