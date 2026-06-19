# Active Stack Nil Context Guard

Status: Completed

## Context

The active sensor-validation stack does not contain the nil-context guard that
was developed on an older parallel branch. `SensorWithContext(nil, sensorID)`
therefore delegates nil-context handling to `net/http` and lacks a stable
PurpleAir-specific error and a regression proving that no request is attempted.

## Requirements

- Reject a nil context with the stable error `purpleair: context is required`.
- Preserve requested sensor-ID validation before context validation.
- Fail before URL construction and HTTP transport execution.
- Preserve cancellation, timeout, response, size, identity, body lifecycle,
  and compatibility-wrapper behavior on the active stack.
- Add mutation-sensitive source, test, static contract, guidance, and completed
  plan evidence.

## Scope Boundaries

- Do not merge or close the older nil-context pull request.
- Do not change public method signatures, request URLs, response parsing,
  dependencies, workflows, or live-network behavior.
- Keep this change stacked on PR #7; do not merge or close either pull request
  without explicit owner authorization.

## Verification Plan

- Run focused nil-context and complete Go tests, formatting, vet, race, build,
  documentation, baseline, and `make check` gates.
- Run the absolute-Makefile gate from an external directory.
- Reject isolated mutations of the guard, stable error, validation order,
  no-request proof, guidance, and completed plan status.
- Audit the exact diff, generated artifacts, credentials, conflicts, modes,
  binaries, large files, dependency and workflow drift, and whitespace before
  commit and push.

## Work Completed

- Restored an explicit nil-context guard on the active sensor-validation stack
  after positive ASCII decimal sensor-ID validation and before request
  construction.
- Added a stable PurpleAir-specific error and a regression proving no transport
  call occurs while malformed sensor IDs retain precedence.
- Added fail-closed source, ordering, test, guidance, and completed-plan
  contracts.

## Verification Results

- Focused nil-context and cancellation tests passed, and `sh -n` plus `dash -n`
  accepted the baseline script.
- A complete isolated `make check` passed with disposable Git metadata before
  the plan was marked completed.
- Repository and external-directory `make check` passed, including formatting,
  vet, unit tests, race tests, documentation, and the fail-closed baseline.
- Seven isolated hostile mutations were rejected across the source guard, stable
  error, regression execution, no-request proof, README guidance, and security
  guidance, plus completed plan status.
