# Reject Declared Oversized Responses Before Reading

Status: Pending

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

Pending implementation.

## Verification Results

Pending implementation and validation.
