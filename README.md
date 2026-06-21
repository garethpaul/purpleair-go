# purpleair-go

<!-- README-OVERVIEW-IMAGE -->
![Project overview](docs/readme-overview.svg)

## Overview

`garethpaul/purpleair-go` is a Go project. GoLang Parser for PurpleAir

This README is based on the checked-in source, manifests, scripts, and repository metadata on the `master` branch. The project language mix found during review was: Go (5).

## Repository Contents

- `README.md` - project overview and local usage notes
- `CHANGES.md` - notable maintenance changes
- `Makefile` - local verification entry points
- `.github/workflows/check.yml` - hosted current-Go verification matrix
- `go.mod`
- `go.sum`
- `docs/plans` - canonical completed maintenance plans
- `plans` - completed maintenance plans
- `scripts/check-baseline.sh` - repository maintenance baseline guard
- `SECURITY.md` - security reporting and disclosure guidance
- `VISION.md` - project direction and maintenance guardrails

Additional scan context:

- Source directories: no top-level source directories detected
- Dependency and build manifests: go.mod, go.sum, Makefile
- Entry points or build surfaces: Makefile
- Test-looking files: client_test.go, sensor_test.go

## Getting Started

### Prerequisites

- Git
- Go

### Setup

```bash
git clone https://github.com/garethpaul/purpleair-go.git
cd purpleair-go
go mod download
```

The setup commands above are derived from repository files. Legacy mobile, Python, or JavaScript samples may require older SDKs or package versions than a modern workstation uses by default.

## Running or Using the Project

- Import `github.com/garethpaul/purpleair-go` from Go code and construct a client with `NewClient()`.
- Use `NewClientWithBaseURL(baseURL)` when a local proxy, fixture server, or
  alternate PurpleAir-compatible endpoint is needed.
- Use `SensorWithError(sensorID)` for error-returning calls. The compatibility
  `Sensor(sensorID)` wrapper preserves its pointer-only signature but returns
  `nil` on failure under the Sensor process exit boundary.
- Use `SensorWithContext(ctx, sensorID)` when callers need cancellation or a
  deadline shorter than the client's HTTP timeout.
- `SensorWithError(sensorID)` returns errors for blank sensor IDs, request
  failures, nil HTTP responses, empty response bodies, non-2xx responses,
  oversized response bodies, malformed JSON, and successful responses that
  contain no sensor results or results with non-positive sensor IDs. Sensor IDs
  must contain only ASCII decimal digits and represent positive integers;
  signed and non-ASCII forms are rejected before any request. Each response
  must preserve the requested sensor identity in at least one result.
- `NewClient()`, nil clients, and zero-value clients use a 30-second total HTTP
  timeout by default. Assign a custom `HTTPClient` or use `SensorWithContext`
  when a caller needs a different deadline.
- Blank custom base URLs fall back to the default PurpleAir JSON endpoint, and
  existing query parameters are preserved when the `show` sensor ID is added.
- Custom base URLs must be absolute `http` or `https` URLs with a host; invalid
  values fall back to the default PurpleAir JSON endpoint.
- Custom base URLs must not embed username/password credentials; use local
  configuration or request headers outside the checked-in URL instead.
- Custom base URLs must not include URL fragments; keep local-only tokens or
  notes out of endpoint strings.
- `SensorWithError` reports transport failures without rendering the request
  URL, so API keys in custom endpoint query strings are not copied into logs;
  the original Go error remains available through `errors.Is` and `errors.As`.
- Response read and JSON decode failures include PurpleAir-specific phase
  context while preserving `errors.Is` and `errors.As`; all non-nil HTTP
  response bodies are closed on successful and failed lookups. A close failure
  is returned after an otherwise successful lookup, while an earlier request,
  status, read, size, decode, or validation error keeps precedence.
- Responses with a declared Content-Length above 1 MiB are rejected before the
  first body read; unknown, absent, and misleading lengths remain protected by
  the bounded read path. Result arrays are decoded incrementally and limited to
  1,024 entries so compact JSON cannot expand into an unbounded slice of large
  result structs.
- `SensorWithContext` propagates the caller context to the HTTP request and
  preserves cancellation and deadline errors through that wrapper.
- The active-stack nil context guard returns `purpleair: context is required`
  before request construction while preserving sensor-ID validation order.

## Testing and Verification

- `go test ./...`
- `go test -race ./...`
- `make lint`
- `make race`
- `make vet`
- `make test`
- `make build`
- `make check`
- `make verify`
- `scripts/check-baseline.sh`

`make vet` runs `go vet ./...`, and `make race` runs `go test -race ./...`.
`make check` delegates to `make verify`, which checks Go formatting, vet, unit
and race tests, the Go build-through-test gate, and completed plans under
`docs/plans`.
Tests and executable examples use mocked HTTP servers and do not call the live
PurpleAir endpoint, including response validation edge cases.
GitHub Actions runs the same gate on Go 1.25.11 and Go 1.26.4 with read-only
permissions and pinned actions.

The baseline script checks required files, module metadata, completed docs-plan
metadata, verification documentation, and local secret/editor metadata hygiene.
GitHub Actions runs the same no-live-network `make check` gate on pushes, pull
requests, and manual dispatches without persisting checkout credentials.

When the required SDK or runtime is unavailable, use static checks and source review first, then verify on a machine that has the matching platform toolchain.

## Configuration and Secrets

- Detected references to PurpleAir. Keep API keys, OAuth credentials, tokens, and account-specific values in local configuration only.

## Security and Privacy Notes

- Review changes touching external API calls or credential-adjacent configuration; examples from the scan include client.go, client_test.go, go.mod, results.go, and 2 more.
- Review changes touching network requests, sockets, or service endpoints; examples from the scan include client.go, client_test.go, sensor.go.
- Review changes touching file, media, JSON, XML, CSV, OCR, or data parsing; examples from the scan include client.go, client_test.go, results.go, sensor.go.
- `NewClientWithBaseURL` rejects URLs with embedded userinfo credentials so
  endpoint configuration does not hide secrets in the base URL.
- `NewClientWithBaseURL` rejects URL fragments so local-only tokens or notes
  do not hide in endpoint configuration.

## Maintenance Notes

- Make gates reject caller-controlled `MAKEFILE_LIST` and `REPO_ROOT` values
  before running Go validation or documentation checks.

- See `SECURITY.md` for vulnerability reporting and safe research guidance.
- See `VISION.md` for project direction and contribution guardrails.
- See `docs/plans/2026-06-08-purpleair-go-baseline.md` for the canonical
  deterministic client-test baseline.
- See `docs/plans/2026-06-08-client-input-and-timeout-guards.md` for the
  sensor ID and timeout guard baseline.
- See `docs/plans/2026-06-12-default-http-timeout-boundary.md` for the bounded
  30-second default and caller override contract.
- See `docs/plans/2026-06-12-sensor-response-identity.md` for positive request
  IDs and requested sensor identity validation.
- See `docs/plans/2026-06-08-sensor-with-error-examples.md` for the executable
  `SensorWithError` examples baseline.
- See `docs/plans/2026-06-09-custom-base-url-client.md` for the custom endpoint
  constructor guard.
- See `docs/plans/2026-06-09-custom-base-url-validation.md` for the custom
  endpoint URL validation guard.
- See `docs/plans/2026-06-09-custom-base-url-credentials-guard.md` for the
  custom endpoint credential guard.
- See `docs/plans/2026-06-09-custom-base-url-fragment-guard.md` for the custom
  endpoint fragment guard.
- See `docs/plans/2026-06-09-empty-response-body-guard.md` for the
  `SensorWithError` empty-body error guard.
- See `docs/plans/2026-06-09-request-failure-context.md` for request-failure
  context, nil-response handling, and repository gate aliases.
- See `docs/plans/2026-06-09-response-body-size-guard.md` for the
  `SensorWithError` response body size guard.
- See `docs/plans/2026-06-13-response-error-context-and-body-close.md` for
  wrapped response-phase errors and response body lifecycle contracts.
- See `docs/plans/2026-06-09-scripted-baseline-check.md` for the scripted
  repository baseline guard and local metadata checks.
- See `docs/plans/2026-06-10-go-vet-verification-gate.md` for the static
  analysis verification gate.
- See `docs/plans/2026-06-10-ci-baseline.md` for the GitHub Actions baseline.
- See `docs/plans/2026-06-10-hosted-go-validation.md` for the current-Go matrix
  and canonical race detector gate.
- See `docs/plans/2026-06-10-sensor-context-cancellation.md` for caller-driven
  request cancellation and deadline support.
- See `docs/plans/2026-06-19-response-boundary-review.md` for request URL
  secrecy, close-error precedence, bounded result decoding, and concurrent
  client reuse evidence.

## Contributing

Keep changes small and tied to the project that is already present in this repository. For code changes, document the toolchain used, avoid committing generated dependency directories or local configuration, and update this README when setup or verification steps change.
