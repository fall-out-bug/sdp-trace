# Slice 49 Evidence

Status: pass

## Scope

Slice 49 is bounded to `cmd/sdp-trace/gate_312` through
`cmd/sdp-trace/gate_324`.

Implemented consolidation:

- moved gate explain entrypoint into `cmd/sdp-trace/gate_explain_cli.go`
- moved gate explain argument parsing and gate-result artifact loading/schema
  validation into `cmd/sdp-trace/gate_explain_inputs.go`
- moved summary and protected checkpoint/condition rendering into
  `cmd/sdp-trace/gate_explain_renderer.go`
- moved required-run, witness-binding, missing-audit-evidence, and override
  collection rendering into `cmd/sdp-trace/gate_explain_collections.go`
- moved shared reason and next-action rendering into neutral
  `cmd/sdp-trace/explain_common_collections.go`
- removed numbered files `gate_312` through `gate_324`

Explicit exclusions:

- gate preview (`gate_325` onward)
- protected run-dir/trust matching (`gate_333` onward)
- override request handling (`gate_352` onward)
- shared JSON/text helpers (`gate_360` onward)
- internal demo/trace/checkpoint package behavior, schemas, fixtures, and MI
  baselines

## MI-Triggered Split

The initial four-file target failed focused MI:

- failed: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 cmd/sdp-trace/gate_explain_cli.go cmd/sdp-trace/gate_explain_renderer.go cmd/sdp-trace/gate_explain_collections.go cmd/sdp-trace/explain_common_collections.go`
- failure: `gate_explain_cli.go` MI `69.4`

Final split keeps responsibility-level files rather than one-helper shards:
command entrypoint, artifact input parsing/loading, summary/protected
rendering, gate-specific collections, and shared common explain collections.

## Focused Verification

- pass: `go test ./cmd/sdp-trace -list '^TestGateExplainRendersProtectedFields$' | grep -qx 'TestGateExplainRendersProtectedFields'`
- pass: `go test ./cmd/sdp-trace -list '^TestGateExplainParseUsage$' | grep -qx 'TestGateExplainParseUsage'`
- pass: `go test ./cmd/sdp-trace -list '^TestGateExplainUnsupportedArtifactCannotVerify$' | grep -qx 'TestGateExplainUnsupportedArtifactCannotVerify'`
- pass: `go test ./cmd/sdp-trace -list '^TestGateExplainMalformedArtifactCannotVerify$' | grep -qx 'TestGateExplainMalformedArtifactCannotVerify'`
- pass: `go test ./cmd/sdp-trace -list '^TestGateExplainDoesNotPrintRawSecretLikeCommand$' | grep -qx 'TestGateExplainDoesNotPrintRawSecretLikeCommand'`
- pass: `go test ./cmd/sdp-trace -list '^TestGateExplainRendersLegacyAndCollectionFields$' | grep -qx 'TestGateExplainRendersLegacyAndCollectionFields'`
- pass: `go test ./cmd/sdp-trace -list '^TestGateExplainRestatesPersistedVerdictsWithoutReevaluation$' | grep -qx 'TestGateExplainRestatesPersistedVerdictsWithoutReevaluation'`
- pass: `go test ./cmd/sdp-trace -run 'Test(GateExplainRendersProtectedFields|GateExplainParseUsage|GateExplainUnsupportedArtifactCannotVerify|GateExplainMalformedArtifactCannotVerify|GateExplainDoesNotPrintRawSecretLikeCommand|GateExplainRendersLegacyAndCollectionFields|GateExplainRestatesPersistedVerdictsWithoutReevaluation)$'`
- pass: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 cmd/sdp-trace/gate_explain_cli.go cmd/sdp-trace/gate_explain_inputs.go cmd/sdp-trace/gate_explain_renderer.go cmd/sdp-trace/gate_explain_collections.go cmd/sdp-trace/explain_common_collections.go`
- pass: `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 cmd/sdp-trace/gate_explain_cli.go cmd/sdp-trace/gate_explain_inputs.go cmd/sdp-trace/gate_explain_renderer.go cmd/sdp-trace/gate_explain_collections.go cmd/sdp-trace/explain_common_collections.go`

Focused tests cover parse usage, missing/malformed artifact load
`cannot_verify`, unsupported schema `cannot_verify`, legacy protected-field
absence, protected checkpoint/condition details, collections, reasons, next
actions, read-only persisted verdict preservation, and secret non-disclosure.

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
  `019e8852-7f51-7152-bafb-2953ec3b0c5f`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T15:37:16+03:00`, prompt class
  `slice-49-implementation-correctness`, timeout `600000ms`, retries `1`,
  fallback `none`, result `LGTM`.
- trust/evidence/spec-drift reviewer: harness `multi_agent_v1`, agent
  `019e8852-c87c-73b1-bb84-3073265e95c4`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T15:37:16+03:00`, prompt class
  `slice-49-implementation-trust-evidence`, timeout `600000ms`, retries `0`,
  fallback `none`, result `LGTM`.
- maintainability/DX reviewer: harness `multi_agent_v1`, agent
  `019e8852-cc0e-77a2-b3ce-8b4424264990`, model/provider `not_assessed`
  (not exposed by harness), date `2026-06-02T15:37:16+03:00`, prompt class
  `slice-49-implementation-maintainability-dx`, timeout `600000ms`,
  retries `0`, fallback `none`, result `LGTM`.
- requested external/provider-qualified lanes: `not_assessed`; unavailable in
  current callable tool surface for this slice. Local `multi_agent_v1` lanes
  record harness, agent id, date, prompt class, timeout, retries, fallback, and
  result when completed.

## Findings

- fixed: added explicit malformed gate-result artifact `cannot_verify`
  regression and exact test-existence evidence after correctness review found
  the first evidence overclaimed malformed artifact coverage.
- none remaining after three implementation review lanes returned `LGTM`.
