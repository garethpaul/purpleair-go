# Custom Base URL Fragment Guard

Status: Completed

## Context

`NewClientWithBaseURL` already trims custom endpoint values and rejects malformed
URLs, non-HTTP schemes, missing hosts, and embedded username/password
credentials. URL fragments are not sent to servers, but allowing them in
endpoint configuration can hide local-only tokens, notes, or state in strings
that look like API URLs.

## Objectives

- Reject custom base URLs that include a fragment.
- Preserve query-parameter support for local proxies and fixture servers.
- Add deterministic unit coverage for fragmented custom endpoints.
- Document the fragment guard in README, VISION, SECURITY, and CHANGES.

## Work Completed

- Updated `isSupportedBaseURL` to require an empty URL fragment.
- Added a fragmented endpoint to the invalid base URL table test.
- Added this completed plan under `docs/plans/`.

## Verification

- `go test ./...`
- `make check`
- `make verify`
- `git diff --check`
