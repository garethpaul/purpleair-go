## PurpleAir Go Vision

PurpleAir Go is a small Go client for fetching and parsing PurpleAir sensor
data from the JSON endpoint.

The repository is useful as a lightweight API wrapper with typed result structs,
a simple client constructor, and tests for basic client configuration.

The goal is to make the client safer and more maintainable while preserving the
straightforward sensor lookup workflow.

Current baseline: `make check` verifies Go formatting, `go vet`, mocked unit
tests, the Go build-through-test gate, and completed `docs/plans` coverage
without calling the live PurpleAir endpoint.

The current focus is:

Priority:

- Preserve `NewClient()` and the `Sensor(sensorId)` lookup path
- Support caller-controlled sensor request cancellation and deadlines
- Keep PurpleAir result fields mapped explicitly
- Avoid hiding HTTP timeouts and user-agent behavior
- Maintain module metadata, mocked HTTP tests, `make lint`, `make vet`,
  `make race`, `make build`, and `make check`
- Run the canonical gate on current supported Go patch releases in hosted CI
  with read-only permissions and pinned actions
- Keep a scriptable baseline guard for required files and local metadata
- Keep executable examples for the preferred error-returning API
- Keep completed maintenance plans under `docs/plans`
- Reject blank sensor IDs before making HTTP requests
- Keep alternate endpoint configuration explicit through constructors
- Keep custom endpoint URLs constrained to absolute HTTP(S) URLs
- Reject custom endpoint URLs that embed username/password credentials
- Reject custom endpoint URLs that include fragments
- Return explicit errors for nil HTTP responses from custom transports
- Return explicit errors for empty HTTP response bodies
- Bound sensor response body reads before JSON parsing
- Wrap transport failures with PurpleAir-specific request context

Next priorities:

- Migrate callers from `Sensor` to `SensorWithError`
- Document endpoint assumptions and API availability

Contribution rules:

- One PR = one focused client, model, test, or documentation change.
- Add tests for parsing and error-handling behavior.
- Keep network calls out of deterministic unit tests.
- Document any public API changes.

## Security And Responsible Use

Canonical security policy and reporting:

- [`SECURITY.md`](SECURITY.md)

Sensor data can imply location and environmental conditions. The client should
avoid surprising data collection, should make requested sensor IDs explicit, and
should not log full responses by default.

## What We Will Not Merge (For Now)

- Process exits from library error paths
- Hidden global HTTP configuration
- Live-network-only tests
- Silent changes to public result fields

This list is a roadmap guardrail, not a permanent rule.
Strong user demand and strong technical rationale can change it.
