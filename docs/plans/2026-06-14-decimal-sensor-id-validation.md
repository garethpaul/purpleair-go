# Require Decimal Sensor IDs

Status: In Progress

## Context

The client documents sensor IDs as positive decimal integers, but
`strconv.Atoi` also accepts a leading plus sign. That allows a non-decimal
spelling such as `+17937` to reach the network even though the response is
matched against numeric ID `17937`.

## Objectives

- Accept only trimmed ASCII decimal digits before integer conversion.
- Preserve positive-value and overflow rejection before any HTTP request.
- Keep valid sensor requests, response identity matching, and public APIs
  unchanged.
- Add runtime and static contracts that fail if signed or non-ASCII forms are
  accepted again.

## Scope Boundaries

- Do not change endpoint construction, response parsing, dependencies,
  workflows, or supported Go versions.
- Do not add live PurpleAir requests, credentials, or generated artifacts.

## Verification

- focused invalid-requested-sensor-ID tests
- `make check` from the repository root and an unrelated directory
- hostile mutations removing the digit guard, signed-input regression case,
  no-request proof, documentation, or completed-plan evidence
- `go test -race ./...`, `go vet ./...`, `go mod verify`, and `git diff --check`
