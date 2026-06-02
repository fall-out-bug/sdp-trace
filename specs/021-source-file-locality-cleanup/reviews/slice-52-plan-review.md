# Slice 52 Plan Review

Status: pass

## Scope

Slice 52 is bounded to `cmd/sdp-trace/gate_345` through
`cmd/sdp-trace/gate_348`, covering demo witness expectation and artifact digest
construction for protected gate witness inputs.

Planned consolidation:

- witness expectation construction and run artifact collection into
  `cmd/sdp-trace/protected_witness_expectation.go`
- retained `run.json` opening and SHA-256 digest calculation into
  `cmd/sdp-trace/protected_witness_digest.go`
- split beyond those files only if the MI gate fails, and record the failed
  command plus the narrower responsibility boundary in Slice 52 evidence

Explicit exclusions:

- protected preview status/action helpers (`gate_349` through `gate_351`)
- override request handling (`gate_352` onward)
- shared JSON/text helpers (`gate_360` onward)
- preview mode and required-ID helper shards (`gate_365` onward)
- internal demo/trace package behavior, schemas, fixtures, and MI baselines

## Decision Gate

- Simpler/Faster: move existing expectation/digest helpers into cohesive
  non-numbered files without changing signatures or call sites.
- Blocking Edge Cases: expectations must be derived from discovered run
  directories, not supplied witness summaries; the first discovered run ID must
  anchor the expected run ID; each discovered run must contribute
  `<run-dir-base>/run.json`; digests must be SHA-256 over retained file bytes;
  discovery/open errors must continue to surface as `cannot_verify` through the
  protected witness expectation loader.
- Existing Open Source: no new library is justified; standard library SHA-256
  and existing `trace.OpenRunArtifact` cover the behavior.

## Planned Verification

- focused test existence: exact per-test `go test ./cmd/sdp-trace -list
  '^TestName$'` checks for
  `TestLoadProtectedWitnessExpectationRejectsMissingRuns`,
  `TestDemoWitnessExpectationUsesFirstRunAndRetainedDigests`,
  `TestDemoWitnessExpectationPropagatesRunArtifactReadErrors`,
  `TestGateCommandFailsForWitnessArtifactMismatch`, and
  `TestProtectedGatePassesWithCISignedCheckpointAndBoundWitnessCLI`
- focused execution: `go test ./cmd/sdp-trace -run 'Test(LoadProtectedWitnessExpectationRejectsMissingRuns|DemoWitnessExpectationUsesFirstRunAndRetainedDigests|DemoWitnessExpectationPropagatesRunArtifactReadErrors|GateCommandFailsForWitnessArtifactMismatch|ProtectedGatePassesWithCISignedCheckpointAndBoundWitnessCLI)$'`
- planned-new focused test:
  `TestDemoWitnessExpectationUsesFirstRunAndRetainedDigests` and
  `TestDemoWitnessExpectationPropagatesRunArtifactReadErrors`
- focused evidence must cover missing-run rejection, first-run ID anchoring,
  per-run `<run-dir-base>/run.json` path construction, retained-byte SHA-256
  digest calculation, retained `run.json` open/read error propagation, standard
  witness mismatch failure, and protected gate pass with CI-signed checkpoint
  plus bound witness
- repository: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`,
  `jq empty schema/*.json`, and `git diff --check`
- quality: CRAP strict-less and file/function MI gates without baseline changes

## Review Lanes

- scope reviewer: multi_agent_v1, agent
  `019e8858-c8b3-7063-a2eb-192f2a8bbc77`, model/provider `not_assessed`
  (not exposed by harness), date 2026-06-02, prompt class `plan-scope`,
  timeout 120000 ms waits, retries 1 after retained artifact read-error
  evidence clarification, fallback `not_used`, result `LGTM`
- trust/evidence reviewer: multi_agent_v1, agent
  `019e8858-ccec-7211-9d43-eaf682f92e18`, model/provider `not_assessed`
  (not exposed by harness), date 2026-06-02, prompt class
  `plan-trust-evidence`, timeout 120000 ms waits, retries 1 after retained
  artifact read-error evidence clarification, fallback `not_used`, result
  `LGTM`
- maintainability/DX reviewer: multi_agent_v1, agent
  `019e8858-d21f-7ad0-a8be-8d1e3e48dbcc`, model/provider `not_assessed`
  (not exposed by harness), date 2026-06-02, prompt class
  `plan-maintainability-dx`, timeout 120000 ms waits, retries 1 after retained
  artifact read-error evidence clarification, fallback `not_used`, result
  `LGTM`
- unavailable requested external/provider-qualified lanes must be recorded as
  `not_assessed` with the reason instead of implied through generic reviewer
  wording

## Findings

- scope lane initial major: planned evidence covered missing-run discovery
  failure but not retained `run.json` open/read error propagation; fixed by
  adding planned-new
  `TestDemoWitnessExpectationPropagatesRunArtifactReadErrors` and re-reviewed
  to `LGTM`
