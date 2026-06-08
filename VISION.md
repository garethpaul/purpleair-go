## PurpleAir Go Vision

PurpleAir Go is a small Go client for fetching and parsing PurpleAir sensor
data from the JSON endpoint.

The repository is useful as a lightweight API wrapper with typed result structs,
a simple client constructor, and tests for basic client configuration.

The goal is to make the client safer and more maintainable while preserving the
straightforward sensor lookup workflow.

The current focus is:

Priority:

- Preserve `NewClient()` and the `Sensor(sensorId)` lookup path
- Keep PurpleAir result fields mapped explicitly
- Avoid hiding HTTP timeouts and user-agent behavior
- Maintain module metadata and basic tests

Next priorities:

- Return errors instead of exiting the process with `log.Fatal`
- Use the client's configured base URL and HTTP client in requests
- Add mocked HTTP tests for successful and failed sensor responses
- Document endpoint assumptions and API availability

Contribution rules:

- One PR = one focused client, model, test, or documentation change.
- Add tests for parsing and error-handling behavior.
- Keep network calls out of deterministic unit tests.
- Document any public API changes.

## Security And Responsible Use

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
