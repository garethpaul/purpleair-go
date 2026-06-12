# Default HTTP Timeout Boundary

Status: Completed

## Context

The default PurpleAir client currently waits up to five minutes for a single
sensor request. That duration is excessive for a small JSON lookup and can
amplify upstream stalls into caller resource exhaustion, especially when many
requests are active concurrently.

Callers that need a different service-level objective can already provide a
custom `HTTPClient`, and `SensorWithContext` supports request-specific
deadlines. The package default should therefore be conservative and bounded.

## Priority

Network timeouts are a core reliability boundary. Reducing the default total
request timeout limits stalled connections without removing any explicit
caller override path.

## Prioritized Engineering Backlog

1. Reduce the default total HTTP timeout to 30 seconds now.
2. Consider separate dial, TLS handshake, and response-header limits if the
   client gains a custom transport.
3. Add retry policy only with explicit idempotency and backoff requirements.

## Requirements

- R1. `NewClient`, nil clients, and zero-value clients must use a 30-second
  default total HTTP timeout.
- R2. A caller-provided `HTTPClient` and its timeout must remain unchanged.
- R3. `SensorWithContext` cancellation and deadline behavior must remain
  unchanged.
- R4. Request construction, response size limits, status handling, and decoded
  sensor validation must remain unchanged.
- R5. Tests and static contracts must detect restoration of the five-minute
  timeout or loss of custom-client behavior.
- R6. README, security guidance, vision, and change history must document the
  bounded default and override paths.

## Implementation Units

### U1. Tighten the client default

- **Files:** `client.go`
- Change only `defaultHTTPTimeout` from five minutes to 30 seconds.

### U2. Expand deterministic timeout coverage

- **Files:** `client_test.go`, `scripts/check-baseline.sh`
- Verify constructor, nil/zero-value fallback, and custom client preservation.

### U3. Update maintenance documentation

- **Files:** `README.md`, `SECURITY.md`, `VISION.md`, `CHANGES.md`
- Record the new default and the supported override mechanisms.

## Scope Boundaries

- Do not add retries or a custom transport.
- Do not change public method signatures.
- Do not change response parsing or sensor identity validation.
- Do not remove caller-controlled contexts or custom HTTP clients.

## Verification

- `make check`
- `go test -race ./...`
- `go vet ./...`
- `go mod verify`
- `git diff --check`
- Mutations restoring `5 * time.Minute` or replacing a caller-provided client
  must fail the regression suite.

Completed on 2026-06-12 with `make check`, race-enabled tests, vet, module
verification, and diff hygiene checks passing.
