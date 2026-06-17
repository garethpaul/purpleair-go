# Active Stack Nil Context Guard

Status: Planned

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
