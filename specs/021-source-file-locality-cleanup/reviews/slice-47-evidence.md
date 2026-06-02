# Slice 47 Evidence

Status: pass

## Scope

Slice 47 is bounded to `cmd/sdp-trace/gate_271` through
`cmd/sdp-trace/gate_289`.

Implemented consolidation:

- moved checkpoint subcommand routing and handler map into
  `cmd/sdp-trace/checkpoint_command.go`
- moved checkpoint create execution and checkpoint write handoff into
  `cmd/sdp-trace/checkpoint_create_cli.go`
- moved checkpoint create flag definitions into
  `cmd/sdp-trace/checkpoint_create_flag_defs.go`
- moved checkpoint create parsing and flagset setup into
  `cmd/sdp-trace/checkpoint_create_flags.go`
- moved checkpoint verify execution and exit-code mapping into
  `cmd/sdp-trace/checkpoint_verify_cli.go`
- moved checkpoint verify parsing and required-input checks into
  `cmd/sdp-trace/checkpoint_verify_flags.go`
- moved checkpoint and policy input loading into
  `cmd/sdp-trace/checkpoint_verify_inputs.go`
- removed numbered files `gate_271` through `gate_289`

Explicit exclusions:

- protected-gate checkpoint policy/witness logic (`gate_302` onward)
- shared JSON/text file helpers (`gate_360` onward)
- internal checkpoint package behavior, schemas, fixtures, and MI baselines

## MI-Triggered Split

The initial three-file target failed focused MI:

- failed: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 cmd/sdp-trace/checkpoint_command.go cmd/sdp-trace/checkpoint_create_cli.go cmd/sdp-trace/checkpoint_verify_cli.go`
- failure: `checkpoint_create_cli.go` MI `61.8`
- failure: `checkpoint_verify_cli.go` MI `59.7`

The second split failed focused MI for one file:

- failed: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 ...`
- failure: `checkpoint_create_flags.go` MI `68.7`

Final split keeps responsibility-level files rather than one-helper shards:
command routing, create execution, create flag definitions, create parsing,
verify execution, verify parsing, and verify input loading.

## Focused Verification

- pass: `go test ./cmd/sdp-trace -run 'Test(CheckpointVerifyHelpersCoverPolicyAndExitBranches|CheckpointCreateAndVerifyCLI|CheckpointCreateFlagValidation|CheckpointCreateFailurePaths|CheckpointVerifyRejectsPositionalArgs|CheckpointVerifyInputLoadFailures|RunCheckpointRejectsMissingOrUnknownSubcommand)$'`
- pass: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 cmd/sdp-trace/checkpoint_command.go cmd/sdp-trace/checkpoint_create_cli.go cmd/sdp-trace/checkpoint_create_flag_defs.go cmd/sdp-trace/checkpoint_create_flags.go cmd/sdp-trace/checkpoint_verify_cli.go cmd/sdp-trace/checkpoint_verify_flags.go cmd/sdp-trace/checkpoint_verify_inputs.go`
- pass: `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 cmd/sdp-trace/checkpoint_command.go cmd/sdp-trace/checkpoint_create_cli.go cmd/sdp-trace/checkpoint_create_flag_defs.go cmd/sdp-trace/checkpoint_create_flags.go cmd/sdp-trace/checkpoint_verify_cli.go cmd/sdp-trace/checkpoint_verify_flags.go cmd/sdp-trace/checkpoint_verify_inputs.go`

## Repository Verification

- pass: `go test ./...`
- pass: `go vet ./...`
- pass: `golangci-lint run`
- pass: `go run ./tools/doccheck`
- pass: `go run ./tools/hygienecheck`
- pass: `jq empty schema/*.json`
- pass: `git diff --check`
- pass: `go test -count=1 ./... -coverprofile=coverage.out`
- pass: `go tool cover -func=coverage.out > coverage-func.txt`
- pass: `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`
- pass: `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`
- pass: `go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools`
- pass: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`
- pass: `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`
- pass: `git diff --name-only origin/main...HEAD | go run ./tools/mibaselinepolicy -base-ref origin/main`

No MI baseline files were changed. Focused MI measured the changed/new
production files at or above 70.1 after the recorded split.

## Review Lanes

- correctness reviewer: harness `multi_agent_v1`, agent
  `019e8836-17b3-7db2-bb24-068e699c18f9`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T15:04:43+03:00`, prompt class
  `slice-47-implementation-correctness`, timeout `600000ms`, retries `1`
  after prior usage-limit errors, fallback `none`, result `LGTM`.
- trust/evidence/spec-drift reviewer: harness `multi_agent_v1`, agent
  `019e8836-1cd1-7ae0-8e6c-80a0beb63ec3`, model/provider
  `not_assessed` (not exposed by harness), date
  `2026-06-02T15:04:43+03:00`, prompt class
  `slice-47-implementation-trust-evidence`, timeout `600000ms`, retries `1`
  after prior usage-limit errors, fallback `none`, result `LGTM`.
- maintainability/DX reviewer: harness `multi_agent_v1`, agent
  `019e8836-21b5-72f2-8a06-4eb8ea2530b2`, model/provider
  `not_assessed` (not exposed by harness), date
  `2026-06-02T15:04:43+03:00`, prompt class
  `slice-47-implementation-maintainability-dx`, timeout `600000ms`,
  retries `1` after prior usage-limit errors, fallback `none`, result
  `LGTM`.
- requested external/provider-qualified lanes: `not_assessed`; unavailable in
  current callable tool surface for this slice. Local `multi_agent_v1` lanes
  record harness, agent id, date, prompt class, timeout, retries, fallback, and
  result when completed.

## Findings

- none after three implementation review lanes returned `LGTM`.
