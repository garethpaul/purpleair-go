# Hosted Go Validation

Status: Completed

## Context

The client had formatting, vet, mocked unit, build-through-test, and scripted
baseline gates, but pushes and pull requests did not run them. The tests were
also not exercising Go's race detector as part of the canonical local command.

## Work Completed

- Added `make race` with `go test -race ./...` and wired it into `make verify`
  and `make check`.
- Added a fixed-runner GitHub Actions matrix for the current supported Go patch
  releases, Go 1.25.11 and Go 1.26.4.
- Limited the workflow token to read-only contents access and pinned checkout
  and Go setup actions to reviewed commits, disabled persisted checkout
  credentials, and retained manual dispatch for maintainers.
- Extended the baseline script to preserve race coverage, exact Go versions,
  action pins, permissions, the canonical command, and this completed plan.

## Verification

- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `make check`
- `git diff --check`
