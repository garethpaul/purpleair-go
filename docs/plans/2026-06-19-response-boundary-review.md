# PurpleAir response boundary review

Status: Completed

## Scope

Review the complete open pull-request stack through the request construction,
HTTP transport, response lifecycle, bounded body read, JSON decoding, sensor
identity validation, compatibility wrapper, reusable client, and hosted Go
verification boundaries.

## Findings

- `http.Client.Do` wraps transport failures in `url.Error`, whose string includes
  the full request URL. A custom PurpleAir-compatible endpoint with an API key
  in its query string therefore copied that key into ordinary logged errors.
- Response body close failures were ignored even after an otherwise successful
  read and decode, so connection finalization failures were indistinguishable
  from success.
- The 1 MiB body limit bounded input bytes but `json.Unmarshal` allocated the
  complete `[]Result` before validation. A compact array of empty objects could
  therefore amplify a small response into hundreds of thousands of large Go
  structs.

## Design

Return a stable request-failure message that never renders the request URL while
retaining the original error as the unwrap cause. Close every non-nil body and
report close failures only when no earlier error exists. Decode the top-level
response into a raw results field, then stream at most 1,024 result objects into
the public model before sensor identity validation.

## Verification

- Red-first fake-transport tests proved query API key disclosure, ignored close
  errors, and unbounded result acceptance on the reviewed aggregate head.
- Focused tests cover randomized query secrets, `errors.Is` preservation,
  close-error precedence, exact result-count boundaries, malformed/non-finite
  numeric payloads, and concurrent client reuse.
- `go test ./...`, `go test -race ./...`, `go vet ./...`, and `make check` pass.
- Hosted Go 1.25.11, Go 1.26.4, and CodeQL gates must pass before merge.

## Residual risk

All HTTP behavior is verified with fake transports and local test servers. No
live PurpleAir API key, endpoint, sensor, provider response, availability, or
rate-limit behavior is exercised by this review.
