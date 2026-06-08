# Changes

## 2026-06-08

- Added executable `SensorWithError` examples for mocked success and blank
  sensor ID error paths.
- Added blank sensor ID validation and default timeout coverage for zero-value
  clients.
- Added `SensorWithError` so callers can handle request, response, and JSON parsing failures without a process exit.
- Updated `Sensor` to keep the original API while delegating to the error-returning implementation.
- Replaced the live-network sensor test with mocked HTTP coverage for successful and failed responses.
- Added `make verify` for Go formatting checks and the full test suite.
- Added `make check` as the shared repository verification alias.
- Added mocked coverage for malformed JSON and empty sensor result responses, and made empty result sets return an explicit `SensorWithError` error.
- Added canonical `docs/plans` coverage and made `make verify` require the
  completed baseline plan.
