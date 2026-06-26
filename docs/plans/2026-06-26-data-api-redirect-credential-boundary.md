# Data API Redirect Credential Boundary

Status: Completed

## Problem

`DataAPIClient` rejected redirects only through its default HTTP client. A
caller-provided client with normal redirect behavior could follow a response
while the request carried `X-API-Key` and an optional private `read_key`.

## Decision

Shallow-clone the selected client for each authenticated request and replace
only `CheckRedirect` with the package's fail-closed policy. Preserve the
caller's transport, timeout, jar, and other settings without mutating its client.

## Verification

- The focused test failed first by following the redirect to the destination.
- The regression uses a permissive custom client and private sensor key, then
  requires the original 302 response and zero destination requests.
- Full `make check`, Go race, vet, hosted, and exact-head review evidence is
  recorded before merge.
