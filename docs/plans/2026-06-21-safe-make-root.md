# Safe Makefile Root Resolution

Status: Completed

## Context

Caller-controlled `MAKEFILE_LIST` redirected formatting, vetting, tests, race
checks, documentation checks, and baseline validation outside the reviewed Go
checkout despite the protected `REPO_ROOT` assignment.

## Scope Boundaries

- Do not change the PurpleAir client API, network behavior, timeouts, response
  boundaries, or dependency versions.
- Preserve credential-free, no-live-network validation.
- Preserve the hosted Go 1.25 and Go 1.26 matrix.

## Work Completed

- Reject command-line and environment replacement of `MAKEFILE_LIST`.
- Canonicalize the checked-in Makefile directory through quoted POSIX tools.
- Add coverage for all nine pre-existing public Make targets plus the root regression gate.
- Include the root policy in `make verify` and `make check`.

## Verification Completed

- `make fmt`, `make vet`, `make test`, `make race`, `make build`, `make docs`,
  `make root-test`, `make verify`, and `make check` passed on Go 1.25.3.
- All 30 target and `REPO_ROOT` override cases passed from a checkout path with
  spaces and an apostrophe.
- Command-line and environment `MAKEFILE_LIST` overrides failed closed.
- Go source, API behavior, module metadata, and dependency sums were unchanged.
