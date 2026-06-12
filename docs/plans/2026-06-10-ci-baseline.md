# PurpleAir Go CI Baseline

Status: Completed

## Context

The client already has mocked tests, `go vet`, formatting, docs-plan checks, and
a scripted repository baseline through `make check`. The missing guard was
hosted CI that repeats that same no-live-network gate.

## Changes

- Added `.github/workflows/check.yml` for GitHub Actions.
- Used `actions/setup-go` with the stable Go toolchain.
- Ran `make check` in the hosted workflow.
- Extended the baseline script and docs so hosted CI stays visible.

## Verification

- `make check`
- `git diff --check`
