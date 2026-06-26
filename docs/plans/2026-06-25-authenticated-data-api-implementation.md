# Authenticated PurpleAir Data API Implementation

Status: Completed

## Context

The approved design selected a separate authenticated client because the
current PurpleAir Data API has different credentials, endpoints, response
fields, availability, and point semantics from the historical map endpoint.

## Requirements

- Keep the legacy `Client`, constructors, response structs, and sensor methods
  unchanged.
- Require an organization read API key and send it only in `X-API-Key`.
- Keep optional private sensor read keys scoped to one request and out of
  returned errors.
- Request one positive sensor index and the fixed typed phase-one field set.
- Preserve caller context and HTTP clients while retaining package timeout and
  redirect defaults.
- Bound responses to 1 MiB and strictly decode them, close every body, and
  preserve primary failures.
- Validate status, sensor identity, timestamps, coordinate pairs and ranges,
  and finite PM2.5 values.
- Use no automatic retries, AQI conversion, caching, or live credential tests.

## Work Completed

- Added the separate public authenticated client and typed response model.
- Added header and request-scoped credential ownership with detail-safe errors.
- Removed URL-bearing `url.Error` wrappers before exposing request causes so
  private sensor read keys remain absent from inspectable error chains.
- Added bounded transport, strict decoding, lifecycle, and response validation.
- Added table-driven local transport tests without provider credentials or
  live network access.
- Updated adoption, security, roadmap, agent, baseline, and changelog guidance.

## Verification Completed

- Reproduced the missing client as focused compile failures before production
  code was added.
- Passed focused constructor, request, response, status, lifecycle, context,
  redaction, and validation tests in the official Go 1.25 container.
- Ran `make check` in the official Go 1.25.11 container.
- Ran `make check` in the official Go 1.26.4 container.
- Confirmed `git diff --check` passes.

## Residual Risk

PurpleAir owns API field availability, point costs, quotas, and status
semantics. Re-check current provider documentation before expanding fields or
adding caller retry policy. No credentialed live request belongs in the
deterministic repository gate.
