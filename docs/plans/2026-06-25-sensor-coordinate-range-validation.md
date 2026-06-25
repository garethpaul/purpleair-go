# Sensor Coordinate Range Validation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use executing-plans to implement this plan task-by-task.

**Goal:** Reject PurpleAir sensor results whose latitude or longitude falls outside valid geographic ranges.

**Architecture:** Extend the existing per-result numeric validation in `decodeSensorResults`. Keep finite-number checks first, then reject latitude outside `[-90, 90]` or longitude outside `[-180, 180]` before appending the result.

**Tech Stack:** Go 1.13-compatible standard library, `testify/assert`, mocked HTTP transports

Status: Completed

---

### Task 1: Specify Geographic Boundaries

**Files:**
- Modify: `sensor_test.go`

**Step 1: Write the failing test**

Add a table-driven `TestSensorWithErrorRejectsOutOfRangeCoordinates` covering
latitude below -90, latitude above 90, longitude below -180, longitude above
180, and a later invalid result after a valid requested sensor.

**Step 2: Run test to verify it fails**

Run: `go test ./... -run TestSensorWithErrorRejectsOutOfRangeCoordinates`

Expected: FAIL because the current decoder accepts finite coordinates outside
geographic ranges.

### Task 2: Enforce Coordinate Ranges

**Files:**
- Modify: `sensor.go`
- Test: `sensor_test.go`

**Step 1: Add minimal validation**

Reject any result with latitude outside `[-90, 90]` or longitude outside
`[-180, 180]`, reporting the result index and out-of-range coordinates.

**Step 2: Run focused test**

Run: `go test ./... -run TestSensorWithErrorRejectsOutOfRangeCoordinates`

Expected: PASS.

**Step 3: Preserve boundary values**

Add or extend coverage proving exactly `-90`, `90`, `-180`, and `180` remain
accepted.

### Task 3: Document and Validate

**Files:**
- Modify: `README.md`
- Modify: `SECURITY.md`
- Modify: `CHANGES.md`
- Modify: `docs/plans/2026-06-25-sensor-coordinate-range-validation.md`

**Step 1: Document the boundary**

State that decoded results must contain finite latitude/longitude values within
valid geographic ranges.

**Step 2: Run repository verification**

Run: `gofmt -w sensor.go sensor_test.go`

Run: `make check`

Expected: formatting, vet, tests, race tests, build, baseline, and module-tidy
checks pass in a Go-capable environment.

**Step 3: Mark the plan complete**

Set `Status: Completed` only after the focused test and `make check` pass.

## Verification Completed

- RED: `go test ./... -run TestSensorWithErrorRejectsOutOfRangeCoordinates -count=1` failed before implementation because finite out-of-range coordinates were accepted.
- GREEN: focused rejection and exact-boundary acceptance tests passed after the decoder guard was added.
- Component gates: `make lint vet test race build root-test` passed in the Go 1.25 container before final plan completion.
- Canonical verification: `make check` passed in the Go 1.25 container after the plan and maintenance record were finalized.
