# Slice 48 Evidence

Status: pass

## Scope

Slice 48 is bounded to `cmd/sdp-trace/gate_302` through
`cmd/sdp-trace/gate_311`.

Implemented consolidation:

- moved protected gate run/resolve orchestration and result writing into
  `cmd/sdp-trace/protected_gate_core.go`
- moved protected checkpoint replay and evaluation input construction into
  `cmd/sdp-trace/protected_gate_evaluation.go`
- moved required checkpoint/policy/witness input structure and JSON loading
  into `cmd/sdp-trace/protected_gate_inputs.go`
- moved protected row loading and witness expectation loading into
  `cmd/sdp-trace/protected_gate_loaders.go`
- removed numbered files `gate_302` through `gate_311`

Explicit exclusions:

- protected checkpoint trust matching (`gate_334` through `gate_344`)
- demo witness construction (`gate_345` onward)
- gate explain (`gate_312` through `gate_324`)
- preview (`gate_325` onward)
- override request handling (`gate_352` onward)
- shared JSON/text helpers (`gate_360` onward)
- internal checkpoint/demo/trace package behavior, schemas, fixtures, and MI
  baselines

## MI-Triggered Split

The initial three-file target failed focused MI:

- failed: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 cmd/sdp-trace/protected_gate_core.go cmd/sdp-trace/protected_gate_inputs.go cmd/sdp-trace/protected_gate_loaders.go`
- failure: `protected_gate_core.go` MI `64.5`

The first implementation split moved result writing into a one-helper
`protected_gate_writer.go` file. Maintainability review flagged that as
recreating the mini-file pathology, so result writing was folded back into
`protected_gate_core.go`; focused tests and MI still passed.

Final split keeps responsibility-level files rather than one-helper shards:
run/resolve/result writing, checkpoint/evaluation input construction, required
trust inputs, and row/witness expectation loading.

## Focused Verification

- pass: `go test ./cmd/sdp-trace -list 'Test(ProtectedGateRejectsLocalSignedCheckpointCLI|ProtectedGatePassesWithCISignedCheckpointAndBoundWitnessCLI|ProtectedGateCoreFailurePaths|ProtectedGateCoreWriteAndEvaluationBranches)$'`
- pass: `go test ./cmd/sdp-trace -run 'Test(ProtectedGateRejectsLocalSignedCheckpointCLI|ProtectedGatePassesWithCISignedCheckpointAndBoundWitnessCLI|ProtectedGateCoreFailurePaths|ProtectedGateCoreWriteAndEvaluationBranches)$'`
- pass: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 cmd/sdp-trace/protected_gate_core.go cmd/sdp-trace/protected_gate_evaluation.go cmd/sdp-trace/protected_gate_inputs.go cmd/sdp-trace/protected_gate_loaders.go`
- pass: `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 cmd/sdp-trace/protected_gate_core.go cmd/sdp-trace/protected_gate_evaluation.go cmd/sdp-trace/protected_gate_inputs.go cmd/sdp-trace/protected_gate_loaders.go`

Focused tests include explicit assertions for `PolicyProvided: true` and UTC
evaluation time on the protected gate evaluation input path.

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
  `019e8842-aab2-7090-b7a1-d605a51ff131`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T15:21:21+03:00`, prompt class
  `slice-48-implementation-correctness`, timeout `600000ms`, retries `0`,
  fallback `none`, result `LGTM`.
- trust/evidence/spec-drift reviewer: harness `multi_agent_v1`, agent
  `019e8842-f08b-79c0-aa8f-7170c6dd9830`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T15:21:21+03:00`, prompt class
  `slice-48-implementation-trust-evidence`, timeout `600000ms`, retries `0`,
  fallback `none`, result `LGTM`.
- maintainability/DX reviewer: harness `multi_agent_v1`, agent
  `019e8842-f3cf-7ab0-ab7d-12ed34853923`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T15:21:21+03:00`, prompt class
  `slice-48-implementation-maintainability-dx`, timeout `600000ms`,
  retries `1`, fallback `none`, result `LGTM`.
- requested external/provider-qualified lanes: `not_assessed`; unavailable in
  current callable tool surface for this slice. Local `multi_agent_v1` lanes
  record harness, agent id, date, prompt class, timeout, retries, fallback, and
  result when completed.

## Findings

- fixed: folded one-helper `protected_gate_writer.go` back into
  `protected_gate_core.go`; focused tests and MI stayed passing.
- none remaining after three implementation review lanes returned `LGTM`.
