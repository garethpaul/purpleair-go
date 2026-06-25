# Required Sensor Coordinate Fields Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use executing-plans to implement this plan task-by-task.

**Goal:** Reject PurpleAir results that omit either `Lat` or `Lon` instead of silently treating the missing JSON field as coordinate zero.

**Architecture:** Preserve the public `Result` model and bounded streaming decoder. Decode each bounded result element through `json.RawMessage`, verify coordinate field presence with pointer fields, then decode the existing public struct and apply the established finite/range validation.

**Tech Stack:** Go standard library, `testify/assert`, Dockerized current Go verification, mocked HTTP servers.

---

Status: Completed

### Task 1: Add failing missing-field coverage

**Files:**
- Modify: `sensor_test.go`

Add a table-driven test for missing latitude, missing longitude, and a later incomplete result after a valid requested sensor. Assert a nil response and an indexed `missing coordinates` error.

Run in Docker:

```bash
docker run --rm -v "$PWD:/src" -w /src golang:1.25 go test ./... -run TestSensorWithErrorRejectsMissingCoordinates -count=1
```

Expected: FAIL because absent numeric fields currently become zero values.

### Task 2: Enforce field presence

**Files:**
- Modify: `sensor.go`

Decode each result as a bounded `json.RawMessage`, unmarshal a small presence probe with `*float64` latitude/longitude fields, reject either nil pointer, then unmarshal the existing `Result` value and retain all current numeric validation.

Run the focused test and expect PASS.

### Task 3: Document and fully validate

**Files:**
- Modify: `README.md`
- Modify: `SECURITY.md`
- Modify: `VISION.md`
- Modify: `CHANGES.md`
- Modify: `docs/plans/2026-06-25-required-sensor-coordinate-fields.md`

Document that decoded results require explicit coordinate fields. Run `gofmt`, the focused test, and Dockerized `make check`. Mark this plan completed only after all gates pass.

## Verification Completed

- RED: the focused missing-coordinate test failed because omitted numeric JSON
  fields decoded to zero and were returned as valid coordinates.
- GREEN: missing latitude, missing longitude, and a later incomplete result are
  rejected with their result index.
- Regression: invalid result IDs retain precedence over coordinate presence,
  while successful and downstream-error fixtures now state coordinates
  explicitly.
- Component gates: formatting, `go test ./...`, `go vet ./...`, and
  `go test -race ./...` passed in the Go 1.25 container.
- The first canonical gate attempt failed because the new baseline wording
  required `VISION.md` to say `explicit coordinates`; the document was aligned
  before rerunning the unchanged source and tests.
- The second canonical gate attempt reached the final metadata scan but
  container Git rejected the host-mounted checkout as dubious ownership; the
  ephemeral container was rerun with `/src` explicitly marked safe.
- Canonical verification: Dockerized `make check` passed.
