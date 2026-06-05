# Slice 52 Evidence

Status: pass

## Scope

Slice 52 consolidated demo witness expectation and retained artifact digest
construction from `cmd/sdp-trace/gate_345` through `cmd/sdp-trace/gate_348`.

Removed numbered files:

- `cmd/sdp-trace/gate_345_demowitnessexpectation.go`
- `cmd/sdp-trace/gate_346_demowitnessartifacts.go`
- `cmd/sdp-trace/gate_347_demowitnessartifact.go`
- `cmd/sdp-trace/gate_348_sha256file.go`

Added cohesive files:

- `cmd/sdp-trace/protected_witness_expectation.go`
- `cmd/sdp-trace/protected_witness_digest.go`

Explicit exclusions kept out of this slice:

- protected preview status/action helpers (`gate_349` through `gate_351`)
- override request handling (`gate_352` onward)
- shared JSON/text helpers (`gate_360` onward)
- preview mode and required-ID helper shards (`gate_365` onward)

## Plan Review

- plan review artifact:
  `specs/021-source-file-locality-cleanup/reviews/slice-52-plan-review.md`
- scope reviewer: multi_agent_v1,
  `019e8858-c8b3-7063-a2eb-192f2a8bbc77`, result `LGTM`
- trust/evidence reviewer: multi_agent_v1,
  `019e8858-ccec-7211-9d43-eaf682f92e18`, result `LGTM`
- maintainability/DX reviewer: multi_agent_v1,
  `019e8858-d21f-7ad0-a8be-8d1e3e48dbcc`, result `LGTM`
- model/provider: `not_assessed` for all lanes because the harness does not
  expose provider-qualified model IDs

## Behavior Evidence

New focused tests:

- `TestDemoWitnessExpectationUsesFirstRunAndRetainedDigests`
- `TestDemoWitnessExpectationPropagatesRunArtifactReadErrors`

Failed-first note: deleting `run.json` changed discovery into
`no run directories found`; the retained artifact error-path test therefore
uses malformed retained `run.json` bytes to exercise `trace.OpenRunArtifact`
decode error propagation.

Focused test existence:

```text
go test ./cmd/sdp-trace -list '^(TestLoadProtectedWitnessExpectationRejectsMissingRuns|TestDemoWitnessExpectationUsesFirstRunAndRetainedDigests|TestDemoWitnessExpectationPropagatesRunArtifactReadErrors|TestGateCommandFailsForWitnessArtifactMismatch|TestProtectedGatePassesWithCISignedCheckpointAndBoundWitnessCLI)$'
```

Result: pass.

Focused execution:

```text
go test ./cmd/sdp-trace -run 'Test(LoadProtectedWitnessExpectationRejectsMissingRuns|DemoWitnessExpectationUsesFirstRunAndRetainedDigests|DemoWitnessExpectationPropagatesRunArtifactReadErrors|GateCommandFailsForWitnessArtifactMismatch|ProtectedGatePassesWithCISignedCheckpointAndBoundWitnessCLI)$'
```

Result: pass.

Package regression:

```text
go test ./cmd/sdp-trace
```

Result: pass.

## Repository Verification

```text
go test ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
jq empty schema/*.json
git diff --check
```

Result: pass.

## Quality Gates

Focused MI for added files:

```text
go run ./tools/qualitycheck -mi-under 70 cmd/sdp-trace/protected_witness_expectation.go cmd/sdp-trace/protected_witness_digest.go
```

Result: pass. Added-file MI values: `protected_witness_expectation.go` 73.9,
`protected_witness_digest.go` 73.7.

Full quality gates:

```text
go test -count=1 ./... -coverprofile=coverage.out
go tool cover -func=coverage.out > coverage-func.txt
go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt
go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less
go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools
go run ./tools/qualitycheck -fail-only -function-mi-under 70 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal
go run ./tools/qualitycheck -fail-only -mi-under 70 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools
```

Result: pass.

## Drift

- numbered Go file count before Slice 52: 525
- numbered Go file count after Slice 52: 521
- spec drift: pass; plan/tasks match the implemented file split
- constitution drift: pass; no Node/JS/tooling dependencies added
- product drift: pass; witness expectation construction rules preserved
- baseline drift: pass; no MI baseline changes

## Review Lanes

- implementation scope/correctness reviewer: multi_agent_v1, agent
  `019e8858-c8b3-7063-a2eb-192f2a8bbc77`, model/provider `not_assessed`
  (not exposed by harness), date 2026-06-02, prompt class
  `implementation-scope-correctness`, timeout 120000 ms waits, retries 0,
  fallback `not_used`, result `LGTM`
- implementation trust/evidence reviewer: multi_agent_v1, agent
  `019e8858-ccec-7211-9d43-eaf682f92e18`, model/provider `not_assessed`
  (not exposed by harness), date 2026-06-02, prompt class
  `implementation-trust-evidence`, timeout 120000 ms waits, retries 0,
  fallback `not_used`, result `LGTM`
- implementation maintainability/DX reviewer: multi_agent_v1, agent
  `019e8858-d21f-7ad0-a8be-8d1e3e48dbcc`, model/provider `not_assessed`
  (not exposed by harness), date 2026-06-02, prompt class
  `implementation-maintainability-dx`, timeout 120000 ms waits, retries 0,
  fallback `not_used`, result `LGTM`
- requested external/provider-qualified lanes remain `not_assessed` because no
  callable provider-qualified model surface is exposed in this session
