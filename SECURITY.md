# Security Policy

## Supported Versions

The supported security scope for `purpleair-go` is the current default branch, `master`. Older commits, tags, branches, forks, demos, and generated artifacts are not actively supported unless the repository explicitly marks them as maintained.

Project summary: GoLang Parser for PurpleAir

## Reporting a Vulnerability

Please report suspected vulnerabilities through GitHub's private vulnerability reporting or by opening a draft GitHub Security Advisory for `garethpaul/purpleair-go` when that option is available. If GitHub does not show a private reporting option for this repository, contact the repository owner through GitHub and avoid posting exploit details publicly until the issue can be assessed.

Do not open a public issue that includes exploit code, secrets, personal data, or detailed reproduction steps for an unpatched vulnerability.

## What to Include

Helpful reports include:

- the affected file, endpoint, permission, dependency, or workflow
- a concise impact statement explaining what an attacker could do
- reproduction steps using test data and accounts you control
- the branch, commit SHA, platform version, device, runtime, or dependency versions used
- logs, screenshots, or proof-of-concept snippets that demonstrate impact without exposing private data

## Project Security Posture

- This repository appears to be a Go project. The active security scope is the code and documentation on the default branch.
- Review found external API integrations or credential-adjacent configuration; changes in those areas should receive security-focused review before merge.
- Review found network clients, sockets, web APIs, or service endpoints; changes in those areas should receive security-focused review before merge.
- Review found file, document, data, or media parsing flows; changes in those areas should receive security-focused review before merge.
- Dependency manifests detected: go.mod, go.sum. Dependency updates should preserve lockfiles when present and avoid introducing packages without a clear maintenance reason.

## Service and API Notes

For web services, APIs, sockets, or scraping workflows, prioritize reports involving authentication bypass, authorization errors, injection, server-side request forgery, unsafe deserialization, credential leakage, data exposure, or denial-of-service conditions. Use test accounts and minimal proof-of-concept traffic only.

Hosted verification runs formatting, vet, mocked tests, and the race detector
with read-only repository permissions and pinned actions. Tests do not call the
live PurpleAir endpoint.

Custom PurpleAir-compatible endpoints should not embed username/password
credentials in the base URL. `NewClientWithBaseURL` rejects URL userinfo and
falls back to the default endpoint so secrets are not hidden in endpoint
strings.
It also rejects URL fragments so local-only tokens or notes do not hide in
endpoint configuration.
`SensorWithError` should return explicit errors for empty HTTP response bodies
instead of panicking or treating malformed upstream responses as valid data.
It should also reject nil HTTP responses from custom transports before reading
status codes or response bodies.
Transport failures should include PurpleAir-specific request context while
preserving the underlying Go error for callers that inspect error chains.
Caller-provided cancellation and deadlines should propagate to sensor HTTP
requests so applications can stop work before the default client timeout.
The default client uses a 30-second timeout for constructor, nil, and zero-value
clients. Callers may provide a custom `HTTPClient` or a shorter context deadline
without the package replacing their policy.
Sensor responses should stay bounded before JSON parsing so a bad endpoint or
custom transport cannot force unbounded memory reads.
GitHub Actions runs the same no-live-network `make check` gate as local
development with read-only permissions, pinned actions, and checkout credential
persistence disabled. Do not add live PurpleAir calls or credentialed smoke
tests to the workflow without a separate security review.
Decoded sensor results should reject non-positive sensor IDs so malformed
upstream records are not returned as valid zero-value sensor data.

## Dependency and Supply Chain Security

Dependency updates should come from trusted package managers and should keep lockfiles in sync when lockfiles exist. Do not commit credentials, private keys, tokens, generated secrets, or machine-local configuration. If a vulnerability depends on a compromised package, typosquatting risk, insecure transitive dependency, or unsafe build step, include the package name, affected version, and the path through which it is used.

## Safe Research Guidelines

Good-faith research is welcome when it stays within these boundaries:

- use only accounts, devices, data, and infrastructure that you own or have explicit permission to test
- avoid destructive actions, persistence, spam, phishing, social engineering, or denial-of-service testing
- minimize access to personal data and stop testing immediately if private data is exposed
- do not exfiltrate secrets or third-party data; report the minimum evidence needed to verify impact
- keep vulnerability details confidential until the maintainer has assessed the report

## Maintainer Response

The maintainer will review complete reports as availability allows, prioritize issues by exploitability and impact, and coordinate a fix or mitigation when the affected code is still maintained. For sample, archived, or educational repositories, the likely remediation may be documentation, dependency updates, or clearly marking unsupported code rather than a production-style patch release.
