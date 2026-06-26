# Changes

## 2026-06-26 - P1 - Preserve authenticated redirect isolation

- Cloned caller HTTP clients per authenticated request and forced redirect
  rejection without mutating caller transport, timeout, jar, or other policy.
- Prevented custom redirect settings from forwarding organization headers or
  private sensor query credentials to another endpoint.
- Strengthened the redirect regression to use a permissive caller client and a
  private sensor read key.

## 2026-06-26 06:30 UTC - P1 - Implement authenticated Data API client

### Summary

Implemented opt-in authenticated real-time sensor access without changing the
legacy map endpoint, constructors, result model, or compatibility behavior.

### Work completed

- Added `DataAPIClient`, `SensorDataOptions`, `SensorDataResponse`, and the
  fixed typed phase-one sensor field model.
- Sent organization keys only in `X-API-Key` and optional private sensor keys
  only in the request query while redacting both from returned failures and
  their inspectable unwrap chains.
- Added caller context propagation, positive sensor-index validation, the
  package 30-second timeout and redirect rejection, and caller HTTP policy
  preservation.
- Added 1 MiB bounded reads, declared-size preflight, strict single-JSON
  decoding, response-body lifecycle handling, detail-safe statuses, identity,
  timestamp, coordinate, and finite PM2.5 validation.
- Preserved point-aware no-retry behavior and the entire legacy `Client` path.
- Added authenticated adoption, security, baseline, and completed-plan records.

### Threads

- Continued the approved authenticated Data API implementation plan.

### Files changed

- `data_api.go` and `data_api_test.go` — public client and transport coverage.
- README, security, vision, agent, changelog, baseline, and completed plan.

### Validation

- Focused constructor, credential, request, response, lifecycle, redaction,
  and validation tests passed in the official Go 1.25 container.
- Full `make check` passed under Go 1.25.11 and Go 1.26.4, including format,
  vet, unit, race, 70 Make authority cases, baseline, and module-tidy checks.

## 2026-06-25 23:21 PDT - P2 - Design authenticated Data API support

### Summary

Completed an implementation-ready design for opt-in PurpleAir Data API support
without silently changing the legacy endpoint, constructors, or result model.

### Work completed

- Compared silent migration, mode flags, and a separate authenticated client.
- Selected a separate `DataAPIClient` with header-only API-key ownership,
  optional private sensor read keys, typed modern response fields, and no
  automatic retries.
- Defined four TDD implementation tasks, validation and error boundaries,
  deferred scope, and official PurpleAir evidence.
- Updated the roadmap from design to phase-one implementation.

### Threads

- None; repository and official provider evidence were reviewed directly.

### Files changed

- `docs/plans/2026-06-25-authenticated-data-api-design.md` — architecture,
  proposed public API, TDD tasks, sources, and deferred scope.
- `README.md`, `SECURITY.md`, `VISION.md` — adoption, credential, and roadmap
  boundaries.
- `scripts/check-baseline.sh` — durable design and documentation contracts.
- `CHANGES.md` — this cycle record.

### Validation

- Red baseline — failed for the missing design artifact before implementation.
- Official PurpleAir API, key, sensor-read-key, field, points, and response
  guidance — reviewed June 25, 2026.
- `make check` in the official Go 1.25 container — passed, including format,
  vet, unit, race, 70 Make authority cases, baseline, and module-tidy checks.
- `make -f /src/Makefile check` from outside the repository in the official Go
  1.25 container — passed.
- Ten hostile design mutations — rejected missing status, client, header,
  endpoint, sensor-key, points, compatibility, retry, README, and roadmap
  contracts.
- `git diff --check` — passed.

### Bugs / findings

- P2 architecture: the authenticated API has incompatible credentials, endpoint,
  response envelope, field names, and availability semantics; silently reusing
  `Client` or `PurpleAir` would create a breaking and ambiguous boundary.

### Blockers

- No live API key is available or required for this design-only cycle.

### Next action

- Implement Task 1 of the authenticated Data API plan with constructor and
  credential-ownership tests.

## 2026-06-26 03:45 UTC - P2 - Migrate Sensor callers

### Summary

Marked the pointer-only `Sensor` wrapper deprecated and documented
`SensorWithError` as the preferred default without removing compatibility.
`SensorWithError` is the preferred default for repository and downstream
callers that need actionable lookup failures.

### Work completed

- Added standard Go deprecation documentation that points callers to
  `SensorWithError` for explicit request, response, parsing, and validation
  failures.
- Added before/after migration guidance and retained `SensorWithContext` for
  caller-owned cancellation and deadlines.
- Added baseline and focused test contracts that keep direct `Sensor` calls
  confined to the two compatibility tests.
- Removed the completed caller-migration roadmap item.

### Threads

- None; the focused compatibility/documentation change was completed directly.

### Files changed

- `sensor.go` - standard deprecation GoDoc with unchanged wrapper behavior.
- `sensor_test.go` and `scripts/check-baseline.sh` - migration contracts.
- README, security, vision, agent guidance, and completed plan - caller guidance.

### Validation

- The focused test failed first because `Sensor` lacked the deprecation marker.
- Go 1.26.4 container verification passed `make check`, all Make aliases,
  explicit race tests, external-directory `make check`, 70 Make authority
  cases, the module-tidy matrix, and the focused caller-boundary test.
- Hosted and review evidence is recorded before merge.

### Bugs / findings

- P2: Repository examples had migrated, but the compatibility API did not tell
  downstream callers that it discards actionable errors.

### Blockers

- No API or runtime behavior changes; external callers must choose when to
  migrate from the deprecated wrapper.

### Next action

- Design authenticated PurpleAir API support without silently changing the
  legacy map-response model.

## 2026-06-25 23:53 UTC - P1 - Clear redirect-policy PR for merge

### Summary

Completed final branch triage for pull request #14 and cleared it for merge
after the requested Codex review was unavailable because the nested CLI lacks
authentication.

### Work completed

- Re-read the complete branch diff and confirmed the redirect policy remains
  limited to package-owned HTTP clients.
- Confirmed caller-provided HTTP clients retain their timeout and redirect
  behavior unchanged.
- Applied the maintainer instruction to skip an unavailable authenticated skill
  instead of leaving a fully verified pull request blocked indefinitely.

### Threads

- Reviewed: PurpleAir PR triage — confirmed the branch is current, clean,
  mergeable, and has no competing issue, review, or TODO work.

### Files changed

- `CHANGES.md` — recorded final review evidence and the authenticated-skill
  exception used for pull request #14.

### Validation

- `git diff --check origin/master...HEAD` — passed.
- Manual branch review — confirmed the default client rejects redirects with
  `http.ErrUseLastResponse`, caller clients are preserved, and the regression
  test proves the redirect destination is never contacted.
- Pull request #14 — Go 1.25.11, Go 1.26.4, and all CodeQL checks passed;
  GitHub reports the pull request mergeable with a clean merge state.
- Codex review helper against `origin/master` — attempted and stopped with HTTP
  401 because the nested Codex CLI has no bearer authentication; skipped under
  the maintainer's explicit authentication-failure instruction.

### Bugs / findings

- No additional code defects found during final review.

### Blockers

- None.

### Next action

- Merge pull request #14, synchronize local `master`, and continue with the next
  green maintenance pull request.

## 2026-06-25 23:36 UTC - P1 - Surface legacy endpoint redirects immediately

### Summary

Stopped package-owned HTTP clients from following PurpleAir endpoint redirects
so the legacy default path returns an immediate status error instead of moving
to another host and potentially waiting for the full timeout.

### Work completed

- Added an `http.ErrUseLastResponse` redirect policy to constructor, nil, and
  zero-value fallback clients.
- Preserved caller-provided `HTTPClient` redirect and timeout policy unchanged.
- Added deterministic two-server coverage proving the redirect destination is
  never contacted and the original HTTP 302 reaches existing status handling.
- Documented the legacy default endpoint, current API incompatibility, and the
  separate modernization boundary.

### Threads

- None. Repository tracing, official Go redirect semantics, current PurpleAir
  API guidance, and a manual endpoint probe provided sufficient direct evidence.

### Files changed

- `client.go` and `client_test.go` — configured and tested package-owned
  redirect rejection while preserving caller clients.
- `sensor_test.go` — proved redirects stop before a second host request.
- `README.md`, `SECURITY.md`, `VISION.md`, and
  `docs/plans/2026-06-25-default-redirect-policy.md` — documented endpoint and
  policy assumptions.
- `scripts/check-baseline.sh` — protected the source, tests, plan, and guidance.

### Validation

- Focused redirect tests — failed before implementation because all default
  clients lacked a redirect hook and the destination server was contacted;
  passed after the `ErrUseLastResponse` policy was added.
- Go 1.25.11 container `make check` — passed formatting, vet, unit tests,
  race detection, build-through-test, completed-plan checks, 70-case Make root
  authority coverage, baseline contracts, and module-tidy mutation tests.
- Pull request #14 — hosted Go 1.25.11 and Go 1.26.4 verification plus CodeQL
  analysis for Actions and Go all passed.
- Codex review helper against `origin/master` — parallel container `make check`
  passed, but the nested Codex CLI stopped before analysis with HTTP 401
  because no local Codex identity is authenticated.
- Manual non-credentialed availability probe — the legacy default URL returned
  HTTP 302 to `purpleair-over-quota-2.appspot.com`; following that redirect
  timed out after 15 seconds with no response body.

### Bugs / findings

- Go's default client follows redirects. That converted the legacy endpoint's
  immediate redirect into unrelated follow-up behavior and hid the actionable
  original status from callers.

### Blockers

- The required Codex review cannot complete until the nested Codex CLI is
  authenticated; do not merge before a clean review.

### Next action

- Authenticate the nested Codex CLI, rerun branch review against `master`, and
  merge pull request #14 only if that review is clean.

## 2026-06-25

- **Timestamp:** 2026-06-25 10:15 UTC
- **Priority:** P2 response integrity.
- **Summary:** Rejected PurpleAir results that omit latitude or longitude
  instead of silently interpreting missing JSON fields as coordinate zero.
- **Work:** Added indexed missing-coordinate errors, preserved invalid-ID error
  precedence, updated successful fixtures to state coordinates explicitly,
  extended baseline contracts, and documented the decoder boundary.
- **Threads:** None.
- **Files changed:** `sensor.go`, `sensor_test.go`, `example_test.go`,
  `scripts/check-baseline.sh`, `README.md`, `SECURITY.md`, `VISION.md`,
  `CHANGES.md`, and
  `docs/plans/2026-06-25-required-sensor-coordinate-fields.md`.
- **Validation:** The focused test failed before implementation and passed
  afterward; full package tests passed in a Go 1.25 container. The first
  `make check` exposed a documentation-contract wording mismatch, the second
  exposed container Git ownership protection, and the final ephemeral
  safe-directory run passed the complete canonical gate.
- **Finding:** Go JSON decoding converted absent numeric coordinate fields to
  valid-looking zero values, allowing incomplete upstream records through the
  existing finite and range checks.
- **Blockers:** None. The host has no Go toolchain, so Docker supplied the
  isolated validation environment.
- **Next action:** Complete `make check`, open a focused pull request, run
  Codex review and hosted checks, and merge only if clean.

## 2026-06-25

- **Timestamp:** 2026-06-25 06:14 UTC
- **Priority:** P2 correctness hardening.
- **Summary:** Rejected finite PurpleAir sensor coordinates outside valid
  geographic latitude and longitude ranges.
- **Work:** Added per-result range validation, regression coverage for each
  invalid direction and a later invalid result, exact endpoint acceptance
  coverage, public usage/security documentation, and an implementation plan.
- **Threads:** None.
- **Files changed:** `sensor.go`, `sensor_test.go`, `README.md`,
  `SECURITY.md`, `CHANGES.md`, and
  `docs/plans/2026-06-25-sensor-coordinate-range-validation.md`.
- **Validation:** The focused test failed before implementation and passed
  afterward; `make lint vet test race build root-test` and `make check`
  passed in a Go 1.25 container.
- **Finding:** The decoder already rejected non-finite numbers but treated
  finite impossible locations such as latitude 91 as valid sensor data.
- **Blockers:** Required Codex review remains unavailable because the managed
  CLI requires an active ChatGPT login and rejects the configured API-key-only
  authentication path.
- **Next action:** Open a focused pull request, run hosted checks, and retry the
  required Codex review before considering merge.

## 2026-06-21

- Hardened all nine pre-existing Make gates against `MAKEFILE_LIST` and
  `REPO_ROOT` redirection without changing the PurpleAir client API.

## 2026-06-19

- Redacted request URLs from transport error strings while preserving the
  original error chain for cancellation, deadlines, and diagnostics.
- Returned response-body close failures after otherwise successful lookups and
  preserved primary error precedence on failed lookups.
- Bounded incremental result decoding to 1,024 entries to prevent compact JSON
  arrays from amplifying into excessive `Result` allocations.
- Added property-style secret-redaction coverage and concurrent client reuse
  coverage under the race detector without live PurpleAir requests.

## 2026-06-17

- Restored the active-stack nil context guard with a stable error, preserved
  sensor-ID validation order, and no-request regression coverage.

## 2026-06-16

- Added the Sensor process exit boundary so the pointer-only compatibility
  lookup returns `nil` on errors instead of terminating the embedding process.

## 2026-06-13

- Wrapped response read and JSON decode failures with PurpleAir-specific phase
  context while preserving underlying Go errors for `errors.Is` and
  `errors.As`.
- Added deterministic coverage proving all non-nil response bodies are closed
  across status, read, size, blank, decode, validation, and success paths.
- Rejected an oversized declared Content-Length before the first response body
  read while preserving the bounded fallback for absent or inaccurate lengths.

## 2026-06-12

- Reduced the default total sensor HTTP timeout from five minutes to a
  30-second boundary for constructor, nil, and zero-value clients.
- Added exact fallback and caller-provided client preservation coverage while
  retaining `SensorWithContext` deadline overrides.
- Rejected missing or non-positive sensor IDs so malformed upstream records
  cannot be returned as valid data.
- Added mocked coverage for invalid identities and multiple positive sensor ID
  results.
- Required positive decimal request IDs and rejected responses that do not
  preserve the requested sensor identity in at least one result.
- Required requested sensor IDs to use ASCII decimal digits, rejecting signed
  and non-ASCII spellings before network access.

## 2026-06-10

- Added credential-free, pinned hosted verification on Go 1.25.11 and Go
  1.26.4 that runs the local no-live-network `make check` baseline.
- Added `make race` and wired `go test -race ./...` into the canonical check.
- Added a `make vet` static analysis gate and wired it into `make verify` and
  `make check`.
- Added `SensorWithContext` so callers can cancel sensor requests or apply
  deadlines while preserving wrapped context errors.

## 2026-06-09

- Added a sensor response body size guard before JSON parsing.
- Wrapped HTTP request failures with PurpleAir-specific context while
  preserving the original transport error.
- Added `scripts/check-baseline.sh` and local metadata coverage for required
  files, Go module metadata, completed plan metadata, and verification docs.
- Added a nil HTTP response guard so custom transports return an error instead
  of panicking in `SensorWithError`.
- Added `make lint` and `make build` aliases to match the shared repository
  verification workflow.

## 2026-06-08

- Added executable `SensorWithError` examples for mocked success and blank
  sensor ID error paths.
- Added `NewClientWithBaseURL` for local proxies, fixture servers, and alternate
  PurpleAir-compatible endpoints.
- Added custom base URL validation so malformed or non-HTTP endpoints fall back
  to the default PurpleAir JSON endpoint.
- Added custom base URL credential validation so URLs with embedded userinfo
  fall back to the default PurpleAir JSON endpoint.
- Added custom base URL fragment validation so local-only tokens or notes do
  not hide in endpoint strings.
- Added an empty response-body guard so `SensorWithError` returns an error
  instead of decoding an empty body.
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
