# PurpleAir Go CI Baseline

Status: Completed

## Context

The client already has mocked tests, `go vet`, formatting, docs-plan checks, and
a scripted repository baseline through `make check`. The missing guard was
hosted CI that repeats that same no-live-network gate.

## Changes

- Added `.github/workflows/check.yml` for GitHub Actions.
- Used pinned checkout and Go setup actions with a Go 1.25.11 and Go 1.26.4
  matrix, read-only repository permissions, and no persisted checkout
  credentials.
- Ran `make check` in the hosted workflow.
- Extended the baseline script and docs so the exact hosted CI contract stays
  visible and mutation-resistant.

## Verification

- `make check`
- `git diff --check`
