---
title: Sensor Result ID Validation
type: fix
date: 2026-06-12
---

# Sensor Result ID Validation

## Summary

Reject PurpleAir responses whose non-empty result array contains missing,
zero, or negative sensor IDs.

## Problem Frame

`SensorWithContext` currently validates only that `results` is non-empty. A
successful response such as `{"results":[{}]}` therefore returns a zero-value
sensor record as valid data.

## Requirements

- R1. Every decoded result must contain a positive sensor ID.
- R2. Invalid result IDs must return a PurpleAir-specific error without
  returning partial sensor data.
- R3. Responses containing multiple valid sensor records must remain accepted.
- R4. Mocked tests and the repository baseline must preserve result-ID
  validation without live network calls.
- R5. README, SECURITY, VISION, and CHANGES must document the response-integrity
  contract.

## Key Technical Decisions

- **Validate every result:** Reject the complete response when any record has a
  non-positive ID so callers never receive a mixture of trusted and malformed
  records.
- **Do not require equality with the query ID:** PurpleAir responses can contain
  related sensor records, so this change validates identity shape without
  narrowing response membership.
- **Keep validation after decoding:** JSON type errors remain owned by
  `encoding/json`; semantic ID validation runs only on decoded records.

## Implementation Units

### U1. Reject malformed result identities

- **Goal:** Return an indexed error for any result with an ID below one.
- **Files:** `sensor.go`
- **Verification:** Focused mocked response tests.

### U2. Add response-integrity regressions

- **Goal:** Cover missing, zero, and negative IDs plus a multi-result success.
- **Files:** `sensor_test.go`, `scripts/check-baseline.sh`
- **Verification:** `go test ./...`, `go test -race ./...`, and `make check`.

### U3. Document result validation

- **Goal:** Keep public behavior and maintenance guidance aligned.
- **Files:** `README.md`, `SECURITY.md`, `VISION.md`, `CHANGES.md`
- **Verification:** `make check` and `git diff --check`.

## Acceptance Examples

- AE1. Given `{"results":[{}]}`, when a sensor response is decoded, then the
  call returns an invalid sensor ID error and no data. Covers R1 and R2.
- AE2. Given `{"results":[{"ID":-1}]}`, when a sensor response is decoded, then
  the call returns an invalid sensor ID error. Covers R1.
- AE3. Given two results with positive IDs, when the response is decoded, then
  both records are returned. Covers R3.

## Scope Boundaries

- Do not require result IDs to equal the requested query value.
- Do not validate optional measurement fields in this pass.
- Do not change the legacy `Sensor` process-exit behavior or public signatures.

## Risks And Mitigations

- Paired sensor records may use different positive IDs. Validate positivity
  only, rather than exact equality with the request.
- Returning partial data could hide corruption. Reject the entire response on
  the first invalid record and include its index in the error.
