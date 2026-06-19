# Reject Declared Oversized Responses Before Reading

Status: Completed

## Context

Sensor responses are bounded to 1 MiB with `io.LimitReader`, which protects
unknown-length and chunked bodies. When a server declares a larger positive
`Content-Length`, the client can reject it immediately instead of reading up to
the same limit first.

## Objectives

- Reject a declared response length above `maxSensorResponseBytes` before the
  first body read.
- Preserve the existing stable oversized-response error and deferred body close.
- Keep the bounded read path for zero, unknown, misleading, or absent lengths.
- Add a reader-sensitive regression proving preflight rejection does not read.
- Protect implementation, tests, documentation, completed status, and exact
  verification evidence in the scripted baseline.

## Scope Boundaries

- Do not raise the response limit or trust a declared length as proof a body is
  small enough.
- Do not change request URLs, credentials, timeouts, decoding, or sensor
  identity validation.
- Do not add dependencies or change the supported Go matrix.

## Verification

- focused preflight and existing oversized-body tests
- all Make gates including `make check` on Go 1.25.11 and Go 1.26.4
- hostile mutations covering comparison boundaries, read avoidance, closure,
  fallback bounded reads, docs, status, and evidence
- `go test -race ./...`, `go vet ./...`, `go mod verify`, and
  `git diff --check`
- exact-base dependency, workflow, API-surface, secret, captured-prompt, and
  generated-artifact scans

## Work Completed

- Added a success-response preflight that rejects a declared `Content-Length`
  above `maxSensorResponseBytes` before reading the response body.
- Preserved status handling, deferred body closure, and the bounded read path
  for unknown, absent, accepted, or misleading declarations.
- Added reader-sensitive tests for zero-read oversized rejection and the exact
  accepted boundary.
- Added scripted contracts and synchronized README, security, vision, and
  change-log guidance.

## Verification Results

- Focused preflight, exact-boundary, fallback oversized-body, and body-close
  tests passed.
- Canonical `make check` passed locally and on Go 1.25.11 and Go 1.26.4 in
  network-disabled, read-only source containers.
- Eight hostile mutations were rejected across the comparison boundary,
  preflight removal, zero-read assertion, body closure, unknown-length
  fallback, documentation contract, completed status, and exact evidence.
- `go mod verify` and `git diff --check` passed; exact-base protected-file,
  secret, captured-prompt, and generated-artifact scans found no changes or
  findings outside the intended paths.
