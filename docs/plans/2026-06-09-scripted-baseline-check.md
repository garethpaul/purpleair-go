# Scripted Baseline Check

## Status: Completed

## Context

The repository had a Makefile gate for Go formatting, mocked tests, and
completed plan metadata, but it did not have a scriptable repository baseline
guard for required files, module metadata, verification docs, and local
metadata hygiene.

## Objectives

- Keep `make check` as the root verification command.
- Add a script-level baseline guard for required repository files.
- Check completed docs-plan metadata without needing to inspect the Makefile
  loop.
- Keep local secrets and editor metadata out of the Go module.

## Work Completed

- Added `scripts/check-baseline.sh`.
- Wired the script into `make check` after the existing verification gate.
- Added `*.iml` to local metadata ignore coverage.
- Updated README, VISION, and CHANGES.

## Verification

- `scripts/check-baseline.sh`
- `go test ./...`
- `make check`
- `git diff --check`

## Follow-Up Candidates

- Add a `go vet` gate if the maintained Go version baseline is modernized.
- Add release tagging only if downstream callers depend on module versions.
