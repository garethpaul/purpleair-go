# Default Redirect Policy Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use executing-plans to implement this plan task-by-task.

**Goal:** Return an immediate status error when the package-owned PurpleAir HTTP client receives a redirect instead of following it into an unrelated timeout or response.

**Architecture:** Configure only clients created by `NewClient`, nil-client fallback, and zero-value fallback to return `http.ErrUseLastResponse` from `CheckRedirect`. This leaves the redirect response body under the existing close/error path and preserves caller-provided `HTTPClient` behavior unchanged.

**Tech Stack:** Go 1.13-compatible `net/http`, `httptest`, table-driven unit tests, repository Make gates.

---

## Status: Completed

## Evidence And Decision

- The default endpoint `https://www.purpleair.com/json?show=...` currently responds with HTTP 302 to `purpleair-over-quota-2.appspot.com`; a June 25, 2026 manual probe then timed out without a response body.
- PurpleAir's current public guidance uses the authenticated `api.purpleair.com` API, whose request authentication and response schema are incompatible with this legacy map-response client.
- Go's default `http.Client` policy follows up to ten redirects. The standard library documents that returning `http.ErrUseLastResponse` stops before the next request and returns the redirect response with a nil error, allowing this package's existing non-2xx handling and body closure to run.
- Migrating to the modern keyed API requires a separate public API and model design. Documentation alone would leave callers waiting for the redirect target. Rejecting redirects only in package-owned clients is the smallest compatible repair.

### Task 1: Add failing redirect tests

**Files:**
- Modify: `client_test.go`
- Modify: `sensor_test.go`

**Step 1: Test the package-owned policy**

Require constructor, nil, and zero-value fallback clients to expose a
`CheckRedirect` hook that returns `http.ErrUseLastResponse`.

**Step 2: Test request behavior**

Use two `httptest` servers. The configured endpoint returns HTTP 302 to the
second server. Require `SensorWithError` to return `purpleair: unexpected
status 302` and prove the second server receives no request.

**Step 3: Verify RED**

Run: `go test ./... -run 'TestDefaultHTTPClientRejectsRedirects|TestSensorWithErrorRejectsRedirectsBeforeFollowing'`

Expected: FAIL because package-owned clients currently follow redirects.

### Task 2: Implement the package-owned policy

**Files:**
- Modify: `client.go`

**Step 1: Add the redirect callback**

Set `CheckRedirect` in `defaultHTTPClient()` to a package function that returns
`http.ErrUseLastResponse`.

**Step 2: Verify GREEN**

Run the focused tests from Task 1.

Expected: PASS with the redirect destination untouched and the original 302
body closed through the existing response path.

### Task 3: Document and fully validate

**Files:**
- Modify: `README.md`
- Modify: `SECURITY.md`
- Modify: `VISION.md`
- Modify: `CHANGES.md`
- Modify: `scripts/check-baseline.sh`
- Modify: `docs/plans/2026-06-25-default-redirect-policy.md`

**Step 1: Document endpoint assumptions**

State that the default URL is a legacy PurpleAir map endpoint, not the current
authenticated API, and that package-owned clients surface redirects rather
than following them. Document that callers supplying an `HTTPClient` retain
their chosen redirect policy.

**Step 2: Run repository validation**

Run: `make check`

Expected: formatting, vet, tests, race detector, build-through-test, plans,
root authority, and baseline contracts all pass without live API requests.

**Step 3: Commit**

```bash
git add client.go client_test.go sensor_test.go README.md SECURITY.md VISION.md \
  CHANGES.md scripts/check-baseline.sh \
  docs/plans/2026-06-25-default-redirect-policy.md
git commit -m "fix: reject default client redirects"
```

## Validation Evidence

- The focused tests failed against the original client because all three
  package-owned client paths had nil redirect hooks and the redirected server
  received the sensor request.
- The focused tests passed after `http.ErrUseLastResponse` was configured: the
  caller received `purpleair: unexpected status 302`, the source received one
  request, and the destination received none.
- Go 1.25.11 container `make check` passed formatting, vet, unit tests, race
  detection, build-through-test, completed-plan checks, 70-case Make root
  authority coverage, baseline contracts, and module-tidy mutation tests.
- Pull request #14 passed hosted Go 1.25.11 and Go 1.26.4 verification plus
  CodeQL analysis for Actions and Go.
- The Codex review helper's parallel container gate passed, but the nested
  Codex CLI returned HTTP 401 before review analysis because it is not
  authenticated locally.
- The canonical automated gate uses only mocked HTTP servers; the one live
  availability probe was a separate manual diagnostic and is not part of
  automated validation.
