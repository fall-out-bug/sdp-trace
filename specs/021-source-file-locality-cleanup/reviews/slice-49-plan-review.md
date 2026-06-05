# Slice 49 Plan Review

Status: pass

## Scope

Slice 49 is bounded to `cmd/sdp-trace/gate_312` through
`cmd/sdp-trace/gate_324`, covering gate explain CLI and read-only rendering
only.

Planned consolidation:

- explain command parsing, gate-result loading, schema validation, and command
  entrypoint into `cmd/sdp-trace/gate_explain_cli.go`
- summary and protected checkpoint/condition rendering into
  `cmd/sdp-trace/gate_explain_renderer.go`
- required-run, witness-binding, missing-evidence, and override collection
  rendering into
  `cmd/sdp-trace/gate_explain_collections.go`
- shared reason and next-action rendering into the neutral
  `cmd/sdp-trace/explain_common_collections.go` because assessment explain
  renderers also call those helpers
- split beyond those files only if the MI gate fails, and record the failed
  command plus the narrower responsibility boundary in Slice 49 evidence

Explicit exclusions:

- gate preview (`gate_325` onward)
- protected run-dir/trust matching (`gate_333` onward)
- override request handling (`gate_352` onward)
- shared JSON/text helpers (`gate_360` onward)
- internal demo/trace/checkpoint package behavior, schemas, fixtures, and MI
  baselines

## Decision Gate

- Simpler/Faster: move the existing read-only explain functions into cohesive
  non-numbered renderer files without changing signatures or output contracts.
- Blocking Edge Cases: explanation must not re-run gate evaluation, must keep
  legacy and protected fields distinct, must return `cannot_verify` for
  missing, malformed, and unsupported persisted schemas, and must not expose raw
  run commands that may contain secrets.
- Existing Open Source: no new library is justified; this is package-local Go
  CLI formatting over existing persisted `demo.GateResult` data.

## Planned Verification

- focused test existence: run exact list checks for each planned test:
  `go test ./cmd/sdp-trace -list '^TestGateExplainRendersProtectedFields$'`,
  `go test ./cmd/sdp-trace -list '^TestGateExplainParseUsage$'`,
  `go test ./cmd/sdp-trace -list '^TestGateExplainUnsupportedArtifactCannotVerify$'`,
  `go test ./cmd/sdp-trace -list '^TestGateExplainMalformedArtifactCannotVerify$'`,
  `go test ./cmd/sdp-trace -list '^TestGateExplainDoesNotPrintRawSecretLikeCommand$'`,
  `go test ./cmd/sdp-trace -list '^TestGateExplainRendersLegacyAndCollectionFields$'`,
  and `go test ./cmd/sdp-trace -list '^TestGateExplainRestatesPersistedVerdictsWithoutReevaluation$'`
- focused execution: `go test ./cmd/sdp-trace -run 'Test(GateExplainRendersProtectedFields|GateExplainParseUsage|GateExplainUnsupportedArtifactCannotVerify|GateExplainMalformedArtifactCannotVerify|GateExplainDoesNotPrintRawSecretLikeCommand|GateExplainRendersLegacyAndCollectionFields|GateExplainRestatesPersistedVerdictsWithoutReevaluation)$'`
- planned-new focused tests: `TestGateExplainRendersLegacyAndCollectionFields`
  and `TestGateExplainRestatesPersistedVerdictsWithoutReevaluation`
- required focused assertions include missing and malformed gate-result artifact
  load failures returning `exitCannotVerify`
- required focused assertions include read-only persisted verdict preservation:
  explain output must restate the persisted gate fields without recomputing,
  re-running, or upgrading them
- repository: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `jq empty schema/*.json`, and `git diff --check`
- quality: CRAP strict-less and file/function MI gates without baseline changes

## Review Lanes

- scope reviewer: harness `multi_agent_v1`, agent
  `019e884a-5189-7fa0-a1ed-4da96e2a0eb0`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T15:28:42+03:00`, prompt class
  `slice-49-plan-scope`, timeout `600000ms`, retries `1`, fallback `none`,
  result `LGTM`.
- trust/evidence reviewer: harness `multi_agent_v1`, agent
  `019e884a-5648-7691-ad6f-27f8a8307297`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T15:28:42+03:00`, prompt class
  `slice-49-plan-trust-evidence`, timeout `600000ms`, retries `1`, fallback
  `none`, result `LGTM`.
- maintainability/DX reviewer: harness `multi_agent_v1`, agent
  `019e884a-5c72-7fb3-a60e-4f02143da0c0`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T15:28:42+03:00`, prompt class
  `slice-49-plan-maintainability-dx`, timeout `600000ms`, retries `1`,
  fallback `none`, result `LGTM`.
- unavailable requested external/provider-qualified lanes must be recorded as
  `not_assessed` with the reason instead of implied through generic reviewer
  wording

## Findings

- fixed: missing/malformed gate-result artifact load failures are explicitly
  preserved and verified as `cannot_verify`.
- fixed: focused test existence evidence requires exact per-test `go test
  -list '^TestName$'` checks.
- fixed: read-only persisted verdict preservation is a named focused test
  requirement.
- fixed: shared reason/next-action renderers move to neutral
  `explain_common_collections.go`, not a gate-local file.
