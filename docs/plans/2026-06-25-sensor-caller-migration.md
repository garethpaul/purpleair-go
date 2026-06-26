# Sensor Caller Migration

## Status: Completed

## Problem

The repository's executable examples already use `SensorWithError`, but the
pointer-only `Sensor` compatibility wrapper is not formally deprecated and the
roadmap still asks maintainers to migrate callers. Without a checked migration
boundary, future examples or production helpers can silently discard request,
response, and parsing errors again.

## Requirements

- Add standard Go deprecation documentation to `Sensor` without removing or
  changing its pointer-only compatibility behavior.
- Make `SensorWithError` the documented default for callers that do not need a
  custom context, and retain `SensorWithContext` for cancellation/deadlines.
- Keep direct `Sensor` calls confined to the two compatibility tests.
- Add static and mutation-sensitive contracts for the deprecation and caller
  boundary.
- Synchronize README, vision, security, agent, and change-log guidance.

## Implementation Plan

1. Add the failing baseline contract and hostile caller/deprecation mutations.
2. Add `// Deprecated: Use SensorWithError...` to the compatibility wrapper.
3. Publish a short migration section with before/after error handling.
4. Remove the completed roadmap item while preserving the public method.
5. Run focused Go tests, the full Make gate, race tests, review, and hosted CI.

## Verification Completed

- The focused Go test failed first because `sensor.go` did not contain the
  standard deprecation marker.
- `Sensor` now delegates exactly as before while directing callers to
  `SensorWithError` for explicit failure handling.
- Repository examples and production helpers contain no direct `Sensor` calls;
  the two compatibility tests remain intentionally exempt.
- README, security, vision, agent, and change-log guidance preserve
  `` `SensorWithError` is the preferred default `` and the roadmap item is
  complete.
- Go 1.26.4 container verification passed `make check`, `make lint`, `make
  test`, `make race`, `make build`, and external-directory `make check`.
- The same run passed 70 Make authority cases, the module-tidy matrix, Go vet,
  all unit and race tests, and the focused caller-boundary test.
- Diff, artifact, review, and hosted evidence is added before merge.
