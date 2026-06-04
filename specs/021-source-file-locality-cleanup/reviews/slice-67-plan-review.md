# Slice 67 Plan Review

Status: in_progress
Date: 2026-06-04

## Scope

Slice 67 is bounded to numbered `cmd/sdp-trace` GitHub Actions artifact HTTP
request/fetch/decode and retained artifact shaping shards:

- `cmd/sdp-trace/packet_086_fetchgithubactionsartifacts.go`
- `cmd/sdp-trace/packet_087_githubactionsartifactsrequest.go`
- `cmd/sdp-trace/packet_088_githubactionsartifactsauthorization.go`
- `cmd/sdp-trace/packet_089_successfulhttpstatus.go`
- `cmd/sdp-trace/packet_090_decodegithubactionsartifacts.go`
- `cmd/sdp-trace/packet_091_retainedgithubartifacts.go`
- `cmd/sdp-trace/packet_092_githubartifactresolver.go`

## Decision Gate

- Simpler/Faster: rename/move only; no behavior change, public API change,
  dependency change, or baseline update.
- Blocking Edge Cases: a single combined HTTP/retention file failed
  pre-change file MI at `66.6`; credential fail-closed behavior and retained
  artifact resolver policy are trust-sensitive and need focused regression
  evidence.
- Existing Open Source: no new parsing, HTTP client abstraction, or artifact
  resolver library is needed; existing Go standard library and project-local
  helpers are sufficient.

## Planned File Boundary

- `cmd/sdp-trace/packet_build_pr_actions_fetch.go`: live artifact fetch.
- `cmd/sdp-trace/packet_build_pr_actions_request.go`: request construction.
- `cmd/sdp-trace/packet_build_pr_actions_authorization.go`: HTTPS-only token
  attachment.
- `cmd/sdp-trace/packet_build_pr_actions_decode.go`: status policy and payload
  decode.
- `cmd/sdp-trace/packet_build_pr_actions_retention.go`: retained artifact
  filtering and resolver construction.

## Planned Regression Evidence

- Exact focused test existence:
  - `TestGitHubActionsArtifactsRequestKeepsHTTPContract`
  - `TestGitHubActionsArtifactsAuthorizationFailClosed`
  - `TestSuccessfulHTTPStatusAcceptsOnly2xx`
  - `TestDecodeGitHubActionsArtifactsKeepsDiagnostics`
  - `TestRetainedGitHubArtifactsKeepsResolverPolicy`
- Focused behavior checks: media type header, URL path construction,
  HTTPS-only authorization, malformed/non-HTTPS authorization fail-closed
  behavior, loopback-test no-token behavior, 2xx-only status policy, JSON
  decode diagnostics, expired artifact filtering, artifact URL precedence,
  synthesized resolver URL construction, missing artifact ID empty resolver
  behavior.
- Standard verification and CRAP/MI gates remain required before the
  implementation review.

## Review Rounds

### Round 1

- scope/correctness: `LGTM`
- trust/evidence: major finding. Focused evidence listed HTTPS-only and
  malformed/non-HTTPS authorization behavior, but did not explicitly require
  loopback-test no-token regression evidence even though it is a credential
  leakage trust boundary.
- maintainability/DX: major finding. Repository verification task omitted the
  AGENTS-required `golangci-lint run` availability check.

### Round 2

- trust/evidence: `LGTM`
- maintainability/DX: `LGTM`

Final plan-review status: `LGTM` across scope/correctness, trust/evidence,
and maintainability/DX after fixes.

## Implementation Adjustment

During implementation, `packet_build_pr_actions_http.go` measured file MI
`67.5`, below the absolute threshold. The final implementation split HTTP
processing into fetch, request, authorization, and decode/status files while
keeping the reviewed behavior scope unchanged.
