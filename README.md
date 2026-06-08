# purpleair-go

<!-- README-OVERVIEW-IMAGE -->
![Project overview](docs/readme-overview.svg)

## Overview

`garethpaul/purpleair-go` is a Go project. GoLang Parser for PurpleAir

This README is based on the checked-in source, manifests, scripts, and repository metadata on the `master` branch. The project language mix found during review was: Go (5).

## Repository Contents

- `README.md` - project overview and local usage notes
- `CHANGES.md` - notable maintenance changes
- `Makefile` - local verification entry points
- `go.mod`
- `go.sum`
- `plans` - completed maintenance plans
- `SECURITY.md` - security reporting and disclosure guidance
- `VISION.md` - project direction and maintenance guardrails

Additional scan context:

- Source directories: no top-level source directories detected
- Dependency and build manifests: go.mod, go.sum, Makefile
- Entry points or build surfaces: Makefile
- Test-looking files: client_test.go, sensor_test.go

## Getting Started

### Prerequisites

- Git
- Go

### Setup

```bash
git clone https://github.com/garethpaul/purpleair-go.git
cd purpleair-go
go mod download
```

The setup commands above are derived from repository files. Legacy mobile, Python, or JavaScript samples may require older SDKs or package versions than a modern workstation uses by default.

## Running or Using the Project

- Import `github.com/garethpaul/purpleair-go` from Go code and construct a client with `NewClient()`.
- Use `SensorWithError(sensorID)` for error-returning calls, or the compatibility `Sensor(sensorID)` wrapper for the original behavior.

## Testing and Verification

- `go test ./...`
- `make verify`

`make verify` checks Go formatting and runs the full test suite. Tests use mocked HTTP servers and do not call the live PurpleAir endpoint.

When the required SDK or runtime is unavailable, use static checks and source review first, then verify on a machine that has the matching platform toolchain.

## Configuration and Secrets

- Detected references to PurpleAir. Keep API keys, OAuth credentials, tokens, and account-specific values in local configuration only.

## Security and Privacy Notes

- Review changes touching external API calls or credential-adjacent configuration; examples from the scan include client.go, client_test.go, go.mod, results.go, and 2 more.
- Review changes touching network requests, sockets, or service endpoints; examples from the scan include client.go, client_test.go, sensor.go.
- Review changes touching file, media, JSON, XML, CSV, OCR, or data parsing; examples from the scan include client.go, client_test.go, results.go, sensor.go.

## Maintenance Notes

- See `SECURITY.md` for vulnerability reporting and safe research guidance.
- See `VISION.md` for project direction and contribution guardrails.

## Contributing

Keep changes small and tied to the project that is already present in this repository. For code changes, document the toolchain used, avoid committing generated dependency directories or local configuration, and update this README when setup or verification steps change.
