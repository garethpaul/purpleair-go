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
- Canonicalize the checked-in Makefile directory with pinned POSIX tools and export it only to recipes.
- Add executed coverage for all nine pre-existing public Make targets plus the root regression gate.
- Include the root policy in `make verify` and `make check`.

## Verification Completed

- 70 executed target and authority cases kept quality commands inside the checkout.
- Hostile checkout backticks were blocked and dollar-substitution paths failed closed.
- `MAKEFILES`, `SHELL`, and `.SHELLFLAGS` authority were covered. Go commands remain repository-owned and `GO` environment values are non-authoritative.
- Command-line and environment `MAKEFILE_LIST` overrides failed closed.
- `make check` remains the complete repository gate and no runtime source changed.
