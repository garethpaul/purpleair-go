# Make Go Verification Location Independent

Status: Planned

## Context

The Make recipes resolve plan globs, Go packages, source globs, and the baseline
script in the caller's current directory. Invoking the repository Makefile by
absolute path from another directory therefore does not reproduce the
repository's verification gate.

## Objectives

- Resolve the repository root from the loaded Makefile, independent of the
  caller's current directory.
- Run every executable Make recipe from that resolved root.
- Protect the root derivation and rooted recipes with dependency-free,
  mutation-sensitive baseline contracts.
- Preserve formatting, vet, test, race, build-through-test, plan, and scripted
  verification behavior.

## Scope Boundaries

- Do not change Go source, APIs, dependencies, supported Go versions, or hosted
  workflow coverage.
- Do not add live PurpleAir requests, credentials, generated files, or new
  tooling.

## Verification

- every Make alias, including `make check`, from the repository root on Go
  1.25.11 and Go 1.26.4
- `make -f /path/to/Makefile check` from an unrelated directory, including an
  attempted command-line repository-root override
- hostile mutations covering root derivation and every rooted recipe
- `go mod verify`, `git diff --check`, and exact-base dependency, workflow,
  API-surface, secret, captured-prompt, and generated-artifact scans

## Work Planned

- Add an override-protected absolute repository root to the Makefile.
- Prefix the docs, format, vet, test, race, and baseline-check recipes with a
  change to that root.
- Extend the scripted baseline with exact Make and completed-plan contracts.
