---
title: Response Error Context and Body Close Contracts
type: fix
status: planned
date: 2026-06-13
---

# Response Error Context and Body Close Contracts

Status: Planned

## Context

PurpleAir request failures already include operation context and preserve their
underlying cause. Response-body read failures and JSON decode failures are
returned directly, leaving callers to infer which request phase failed. The
client defers response body closure, but tests do not protect that lifecycle
across status, read, size, empty-body, decode, validation, and success exits.

## Priority

Stable phase-specific errors make production failures diagnosable while Go
error wrapping preserves programmatic handling. Explicit body-close contracts
prevent connection leaks if later refactors add early returns or move response
validation.

## Prioritized Engineering Backlog

1. Add response read/decode context and body-close regression coverage now.
2. Replace deprecated `ioutil` APIs only when the module's Go 1.13 compatibility
   floor is intentionally raised.
3. Revisit the legacy fatal `Sensor` wrapper in a separate compatibility plan.

## Objectives

- Wrap response-body read errors with `purpleair: read response body` and `%w`.
- Wrap JSON decode errors with `purpleair: decode response body` and `%w`.
- Preserve `errors.Is` for reader failures and `errors.As` for
  `*json.SyntaxError`.
- Prove every non-nil response body closes on status, read, oversize, blank,
  decode, validation, and successful response paths.
- Preserve request context, timeout, URL, body-size, sensor ID, and requested
  identity behavior.
- Protect the error strings, wrapping operators, regression tests, completed
  plan, and public documentation in the fail-closed baseline.

## Implementation Units

### 1. Response error context

Files:

- `sensor.go`
- `sensor_test.go`

Requirements:

- Add stable phase-specific prefixes only at the body read and JSON decode
  boundaries.
- Wrap the original errors so callers retain standard Go error inspection.
- Do not expose response body contents in errors.

### 2. Response body lifecycle coverage

Files:

- `sensor_test.go`

Requirements:

- Use deterministic custom transports and tracking `io.ReadCloser` fixtures.
- Assert closure after each non-nil response path, including successful decode.
- Cover reader failure with `errors.Is` and malformed JSON with `errors.As`.

### 3. Contracts and documentation

Files:

- `scripts/check-baseline.sh`
- `README.md`
- `SECURITY.md`
- `VISION.md`
- `CHANGES.md`
- `docs/plans/2026-06-13-response-error-context-and-body-close.md`

Requirements:

- Document stable response-phase errors and guaranteed body cleanup.
- Preserve the focused implementation and test contracts against silent
  removal.
- Record completed status and actual verification only after all gates pass.

## Test Scenarios

- A reader returns a sentinel error: the result is nil, the error has the read
  prefix, `errors.Is` finds the sentinel, and the body is closed.
- Malformed JSON: the result is nil, the error has the decode prefix,
  `errors.As` finds `*json.SyntaxError`, and the body is closed.
- Non-2xx, oversized, whitespace-only, empty results, non-positive result IDs,
  missing requested identity, and success each close their bodies in separate
  regression cases.
- Existing cancellation, nil response/body, request validation, timeout, and
  URL tests remain unchanged and green.

## Scope Boundaries

- Do not alter successful response data or sensor matching.
- Do not change the 1 MiB response cap, status handling, or HTTP timeout.
- Do not alter context behavior or duplicate the nil-context guard branch.
- Do not change the Go module floor or dependencies.
- Do not make live PurpleAir requests or require credentials.

## Verification

- `gofmt -w sensor.go sensor_test.go`
- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `go test -cover ./...`
- `make lint`
- `make test`
- `make race`
- `make vet`
- `make build`
- `make verify`
- `make check`
- `go mod verify`
- `git diff --check`
