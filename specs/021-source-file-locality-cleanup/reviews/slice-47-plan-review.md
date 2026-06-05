# Slice 47 Plan Review

Status: pass

## Scope

Slice 47 is bounded to `cmd/sdp-trace/gate_271` through
`cmd/sdp-trace/gate_289`, covering checkpoint CLI command handling only.

Planned consolidation:

- checkpoint subcommand routing and handler map into
  `cmd/sdp-trace/checkpoint_command.go`
- checkpoint create flag definitions, required flags, parsing, key loading,
  creation handoff, JSON write, and stdout rendering into
  `cmd/sdp-trace/checkpoint_create_cli.go`
- checkpoint verify flag definitions, parsing, checkpoint/policy loading,
  verification rendering, and exit-code mapping into
  `cmd/sdp-trace/checkpoint_verify_cli.go`
- split beyond those three files only if the MI gate fails, and record the
  failed command plus the narrower responsibility boundary in Slice 47 evidence

Explicit exclusions:

- protected-gate checkpoint policy/witness logic (`gate_302` onward)
- shared JSON/text file helpers (`gate_360` onward)
- internal checkpoint package behavior, schemas, fixtures, and MI baselines

## Decision Gate

- Simpler/Faster: move the existing functions into cohesive non-numbered
  checkpoint CLI files without changing signatures or command contracts.
- Blocking Edge Cases: checkpoint create and verify have different trust
  semantics, stderr/stdout behavior, required flag contracts, optional policy
  behavior, and exit-code mapping; combining them into one large file would risk
  MI regression and weaker review locality.
- Existing Open Source: no new library is justified; this is package-local Go
  CLI organization around existing flag helpers and checkpoint package APIs.

## Planned Verification

- focused: `go test ./cmd/sdp-trace -run 'Test(CheckpointVerifyHelpersCoverPolicyAndExitBranches|CheckpointCreateAndVerifyCLI|CheckpointCreateFlagValidation|CheckpointCreateFailurePaths|CheckpointVerifyRejectsPositionalArgs|CheckpointVerifyInputLoadFailures|RunCheckpointRejectsMissingOrUnknownSubcommand)$'`
- planned-new focused tests: `TestCheckpointCreateFlagValidation`,
  `TestCheckpointCreateFailurePaths`,
  `TestCheckpointVerifyRejectsPositionalArgs`,
  `TestCheckpointVerifyInputLoadFailures`, and
  `TestRunCheckpointRejectsMissingOrUnknownSubcommand`
- repository: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `jq empty schema/*.json`, and `git diff --check`
- quality: CRAP strict-less and file/function MI gates without baseline changes

## Review Lanes

- scope reviewer: harness `multi_agent_v1`, agent
  `019e8814-ff94-7572-a168-ded23b1731b7`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T14:29:32+03:00`, prompt class
  `slice-47-plan-scope`, timeout `600000ms`, retries `0`, fallback `none`,
  result `LGTM`.
- trust/evidence reviewer: harness `multi_agent_v1`, agent
  `019e8815-065d-75e2-906a-f1e18cf515b3`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T14:29:32+03:00`, prompt class
  `slice-47-plan-trust-evidence`, timeout `600000ms`, retries `1`, fallback
  `none`, result `LGTM`.
- maintainability/DX reviewer: harness `multi_agent_v1`, agent
  `019e8815-0ad7-7771-81cd-b7e65b292ce8`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T14:29:32+03:00`, prompt class
  `slice-47-plan-maintainability-dx`, timeout `600000ms`, retries `1`,
  fallback `none`, result `LGTM`.
- unavailable requested external/provider-qualified lanes must be recorded as
  `not_assessed` with the reason instead of implied through generic reviewer
  wording

## Findings

- fixed: planned-new focused tests are explicitly labelled as planned-new
  rather than implied existing tests.
- fixed: implementation review evidence must record reviewer metadata and
  unavailable external/provider-qualified lanes as `not_assessed`.
- fixed: planned target files are named explicitly, with further splits allowed
  only if MI fails and evidence records the reason.
- fixed: focused regression plan now includes checkpoint create/write failures
  and checkpoint/policy input load failures.
