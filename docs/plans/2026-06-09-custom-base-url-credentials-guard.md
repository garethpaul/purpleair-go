# Custom Base URL Credentials Guard

## Status: Completed

## Context

`NewClientWithBaseURL` accepts alternate endpoints for local proxies, fixture
servers, and PurpleAir-compatible services. It already rejects malformed,
hostless, and non-HTTP URLs, but URLs with embedded userinfo credentials still
passed validation and could hide secrets in endpoint strings.

## Goals

- Reject custom base URLs that include username/password userinfo.
- Preserve valid absolute `http` and `https` endpoints with hosts.
- Keep query parameter support for local proxy or fixture URLs.
- Cover credential-bearing URLs with deterministic unit tests.

## Work Completed

- Updated `isSupportedBaseURL` to reject URLs with embedded userinfo.
- Added invalid-value test coverage for a credential-bearing custom endpoint.
- Updated README, VISION, SECURITY, and CHANGES notes for the credential guard.

## Verification

- `gofmt -w client.go client_test.go`
- `go test ./...`
- `make check`
- `make verify`
- `git diff --check`
