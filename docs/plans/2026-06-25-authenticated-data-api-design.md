# Authenticated PurpleAir Data API Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use executing-plans to implement this plan task-by-task.

**Goal:** Add opt-in support for PurpleAir's authenticated real-time sensor API without changing the legacy map endpoint, constructors, result model, or compatibility behavior.

**Architecture:** Introduce a separate `DataAPIClient` with its own credential, request, response, and error types. Keep the existing `Client`, `PurpleAir`, `Result`, `NewClient`, `NewClientWithBaseURL`, `Sensor`, `SensorWithError`, and `SensorWithContext` paths unchanged so adoption is explicit and source-compatible.

**Tech Stack:** Go 1.13-compatible standard library, `net/http`, `context`, `encoding/json`, `httptest`, table-driven tests, repository Make gates.

---

Status: Completed

## Evidence And Decision

The current `Client` is coupled to the historical
`https://www.purpleair.com/json?show=` map response and the large legacy
`PurpleAir`/`Result` model. Existing constructors carry no credential and
preserve custom compatible endpoints and HTTP clients.

Official PurpleAir guidance reviewed on June 25, 2026 establishes that:

- The current base is `https://api.purpleair.com/v1/`; a single real-time
  request uses `api.purpleair.com/v1/sensors/{sensor_index}`.
- A read API key is sent in the `X-API-Key` header and is tied to an organization.
- The API uses a points-based system.
- Private sensor data additionally requires a sensor read key, distinct from
  the organization's API read key.
- The response contains API/data timestamps and a `sensor` object rather than
  the legacy map `results` array.
- The API returns raw measurements and supports selected response fields; AQI
  conversion remains caller-owned.

Primary references:

- https://community.purpleair.com/t/making-api-calls-with-the-purpleair-api/180
- https://community.purpleair.com/t/about-the-purpleair-api/7145
- https://community.purpleair.com/t/sensor-indexes-and-read-keys/4000
- https://community.purpleair.com/t/api-fields-descriptions/4652
- https://community.purpleair.com/t/creating-api-keys/3951

Three approaches were compared:

1. Silently move `NewClient()` to the modern API. Rejected because it changes
   credentials, endpoint availability, schema, and returned data.
2. Add authentication and mode flags to `Client`. Rejected because one type
   would represent incompatible wire models and ambiguous zero values.
3. Add a separate authenticated client. Selected because the legacy Client
   remains unchanged and adoption is explicit.

The legacy Client remains unchanged throughout this plan.

Phase one covers one real-time sensor and a fixed typed field set. It does not
include groups, history, write keys, AQI conversion, or arbitrary fields. It
uses no automatic retries because point consumption and retry policy belong to
callers.

## Proposed Public API

Create `data_api.go` with these Go 1.13-compatible shapes:

```go
const defaultDataAPIBaseURL = "https://api.purpleair.com/v1"

type DataAPIClient struct {
    HTTPClient *http.Client
    baseURL    string
    readAPIKey string
}

type SensorDataOptions struct {
    SensorReadKey string
}

type SensorDataResponse struct {
    APIVersion    string     `json:"api_version"`
    Timestamp     int64      `json:"time_stamp"`
    DataTimestamp int64      `json:"data_time_stamp"`
    Sensor        SensorData `json:"sensor"`
}

type SensorData struct {
    SensorIndex int      `json:"sensor_index"`
    Name        *string  `json:"name"`
    LastSeen    *int64   `json:"last_seen"`
    Latitude    *float64 `json:"latitude"`
    Longitude   *float64 `json:"longitude"`
    PM25ATM     *float64 `json:"pm2.5_atm"`
}

func NewDataAPIClient(readAPIKey string) (*DataAPIClient, error)
func (c *DataAPIClient) SensorData(
    ctx context.Context,
    sensorIndex int,
    options SensorDataOptions,
) (*SensorDataResponse, error)
```

The constructor trims and rejects a blank API read key. The key stays in an
unexported field, is sent only in `X-API-Key`, and never appears in URLs, logs,
or returned errors. The optional sensor read key is added as `read_key` only for
one private-sensor request and is excluded from errors.

The request asks for `name,last_seen,latitude,longitude,pm2.5_atm`;
`sensor_index` is returned by default and is not added to `fields`. Pointer
fields distinguish missing provider values from valid zero coordinates or
measurements.

## Transport And Error Boundary

- Reject nil context, blank API key, and non-positive sensor index before I/O.
- Build the path from the validated integer, never caller path text.
- Keep the package-owned 30-second timeout and redirect rejection; preserve a
  caller-provided `HTTPClient` exactly.
- Require a successful JSON response, cap reads at 1 MiB, reject trailing JSON,
  and close every non-nil body while preserving the primary failure.
- Require the returned sensor index to match, timestamps to be non-negative,
  coordinates to be both absent or both finite/in range, and PM2.5 to be finite.
- Return detail-safe status errors for 401, 403, 404, 429, and 5xx. Use no
  automatic retries and never include either credential or provider body.

## Task 1: Add Constructor And Credential Tests

**Files:**
- Create: `data_api.go`
- Create: `data_api_test.go`

**Step 1: Write failing tests**

Cover blank keys, whitespace normalization, fixed base URL, package timeout and
redirect behavior, and caller-owned HTTP client retention.

**Step 2: Verify RED**

Run: `go test ./... -run 'TestNewDataAPIClient|TestDataAPIClientHTTPClient'`

Expected: FAIL because `DataAPIClient` does not exist.

**Step 3: Implement and verify GREEN**

Return `purpleair: API read key is required` for blank input without echoing it.

```bash
go test ./... -run 'TestNewDataAPIClient|TestDataAPIClientHTTPClient'
git add data_api.go data_api_test.go
git commit -m "feat: add authenticated data API client"
```

## Task 2: Add Authenticated Request Tests

**Files:**
- Modify: `data_api.go`
- Modify: `data_api_test.go`

**Step 1: Write failing request tests**

Use `httptest.Server` to require `GET`, exact sensor path, fixed fields,
`X-API-Key`, optional `read_key`, no credential in the URL, and context
propagation.

**Step 2: Implement, verify, and commit**

```bash
go test ./... -run 'TestDataAPISensorRequest|TestDataAPIPrivateSensorRequest'
git add data_api.go data_api_test.go
git commit -m "feat: request authenticated sensor data"
```

## Task 3: Add Response And Failure Tests

**Files:**
- Modify: `data_api.go`
- Modify: `data_api_test.go`

**Step 1: Write failing response tests**

Cover the response envelope, mismatched IDs, partial/out-of-range coordinates,
non-finite numbers, negative timestamps, oversized/empty bodies,
malformed/trailing JSON, nil responses, body closure, redirects, status errors,
cancellation, and detail-safe messages.

**Step 2: Implement bounded decoding and commit**

Do not unmarshal modern data into `PurpleAir`, `Results`, or `Result`.

```bash
go test ./... -run 'TestDataAPISensorResponse|TestDataAPISensorFailure'
git add data_api.go data_api_test.go
git commit -m "feat: validate authenticated sensor responses"
```

## Task 4: Document Adoption And Validate

**Files:**
- Modify: `README.md`
- Modify: `SECURITY.md`
- Modify: `VISION.md`
- Modify: `CHANGES.md`
- Modify: `AGENTS.md`
- Modify: `scripts/check-baseline.sh`
- Modify: `docs/plans/2026-06-25-authenticated-data-api-design.md`

Document key ownership, points, fixed fields, raw PM data, no AQI conversion,
no retries, and the unchanged legacy client.

Run: `make check`

Expected: formatting, vet, tests, race tests, build-through-test, module tidy,
Make authority tests, baseline contracts, and no live network access pass.

```bash
git add README.md SECURITY.md VISION.md CHANGES.md AGENTS.md \
  scripts/check-baseline.sh docs/plans/2026-06-25-authenticated-data-api-design.md
git commit -m "docs: define authenticated API adoption"
```

## Explicitly Deferred

- Any change to the legacy public API or result fields.
- Arbitrary fields or dynamic result maps.
- Multi-sensor, group, historical, CSV, write-key, or mutation endpoints.
- AQI/conversion calculations, retries, backoff, point-budget management,
  caching, persistence, live API tests, or checked-in credentials.

## Verification Completed

- Reviewed the legacy client, model, transport, tests, baseline, and Go 1.13
  floor.
- Compared three approaches and selected a separate authenticated client.
- Cross-checked authentication, endpoint, private sensor read keys, response
  envelope, points, and field behavior against the official sources above.
- Added a red-first baseline contract and ran `make check` after completing the
  design and documentation.
- Ran the complete check target both from the repository root and through an
  absolute Makefile path outside the repository in the official Go 1.25 image.
- Confirmed ten hostile mutations fail closed across the design, README, and
  roadmap contracts, then passed `git diff --check`.
- Phase one was subsequently implemented as the separate authenticated client
  recorded in `2026-06-25-authenticated-data-api-implementation.md`.

## Residual Risk

The official API and field catalog can evolve, point costs and availability are
provider-owned, and no credentialed live request belongs in deterministic
tests. Re-check current official documentation during implementation.
