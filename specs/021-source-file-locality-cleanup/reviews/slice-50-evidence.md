# Slice 50 Evidence

Status: pass

## Scope

Slice 50 consolidated gate preview and protected target selection from
`cmd/sdp-trace/gate_325` through `cmd/sdp-trace/gate_333`.

Removed numbered files:

- `cmd/sdp-trace/gate_325_gatepreviewreport.go`
- `cmd/sdp-trace/gate_326_protectedgatepreviewreport.go`
- `cmd/sdp-trace/gate_327_rungatepreview.go`
- `cmd/sdp-trace/gate_328_parsegatepreviewargs.go`
- `cmd/sdp-trace/gate_329_gatepreviewstringflags.go`
- `cmd/sdp-trace/gate_330_buildgatepreviewreport.go`
- `cmd/sdp-trace/gate_331_runprotectedgatepreview.go`
- `cmd/sdp-trace/gate_332_newprotectedgatepreviewreport.go`
- `cmd/sdp-trace/gate_333_protectedrundir.go`

Added cohesive files:

- `cmd/sdp-trace/gate_preview_cli.go`
- `cmd/sdp-trace/gate_preview_args.go`
- `cmd/sdp-trace/gate_preview_standard.go`
- `cmd/sdp-trace/gate_preview_reports.go`
- `cmd/sdp-trace/gate_preview_protected.go`
- `cmd/sdp-trace/protected_gate_run_dir.go`

Explicit exclusions kept out of this slice:

- protected checkpoint trust matching (`gate_334` through `gate_344`)
- demo witness construction (`gate_345` onward)
- protected preview status/action helpers (`gate_349` through `gate_351`)
- override request handling (`gate_352` onward)
- shared JSON/text helpers (`gate_360` onward)
- preview mode and required-ID helper shards (`gate_365` onward)

## Plan Review

- plan review artifact:
  `specs/021-source-file-locality-cleanup/reviews/slice-50-plan-review.md`
- scope reviewer: multi_agent_v1,
  `019e8858-c8b3-7063-a2eb-192f2a8bbc77`, result `LGTM`
- trust/evidence reviewer: multi_agent_v1,
  `019e8858-ccec-7211-9d43-eaf682f92e18`, result `LGTM`
- maintainability/DX reviewer: multi_agent_v1,
  `019e8858-d21f-7ad0-a8be-8d1e3e48dbcc`, result `LGTM`
- model/provider: `not_assessed` for all lanes because the harness does not
  expose provider-qualified model IDs

## Behavior Evidence

New or tightened focused tests:

- `TestGatePreviewStandardReportShape`
- `TestGatePreviewParseAndContractFailurePaths`
- `TestProtectedGatePreviewInputFailurePaths`
- `TestGatePreviewIsReadOnlyAndDoesNotPrintSecretLikeValues` now rejects any
  standard gate verdict field
- `TestGatePreviewReportsWitnessArtifactMismatch` now rejects any standard gate
  verdict field

Focused test existence:

```text
go test ./cmd/sdp-trace -list '^(TestProtectedGatePreviewRendersAbsentInputsWithoutWriting|TestProtectedGateRequiresSingleRunDir|TestGatePreviewStandardReportShape|TestGatePreviewIsReadOnlyAndDoesNotPrintSecretLikeValues|TestGatePreviewReportsWitnessArtifactMismatch|TestGatePreviewParseAndContractFailurePaths|TestProtectedGatePreviewInputFailurePaths)$'
```

Result: pass.

Focused execution:

```text
go test ./cmd/sdp-trace -run 'Test(ProtectedGatePreviewRendersAbsentInputsWithoutWriting|ProtectedGateRequiresSingleRunDir|GatePreviewStandardReportShape|GatePreviewIsReadOnlyAndDoesNotPrintSecretLikeValues|GatePreviewReportsWitnessArtifactMismatch|GatePreviewParseAndContractFailurePaths|ProtectedGatePreviewInputFailurePaths)$'
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

Initial MI check for the coarser `gate_preview_cli.go` split failed:

```text
go run ./tools/qualitycheck -mi-under 70 cmd/sdp-trace/gate_preview_cli.go cmd/sdp-trace/gate_preview_reports.go cmd/sdp-trace/gate_preview_protected.go cmd/sdp-trace/protected_gate_run_dir.go
```

Result: fail. `cmd/sdp-trace/gate_preview_cli.go` had file MI 68.3.

Resolution: split argument parsing and standard report building into
`gate_preview_args.go` and `gate_preview_standard.go`.

Focused MI after split:

```text
go run ./tools/qualitycheck -mi-under 70 cmd/sdp-trace/gate_preview_cli.go cmd/sdp-trace/gate_preview_args.go cmd/sdp-trace/gate_preview_standard.go cmd/sdp-trace/gate_preview_reports.go cmd/sdp-trace/gate_preview_protected.go cmd/sdp-trace/protected_gate_run_dir.go
```

Result: pass. Lowest file MI among added files was
`gate_preview_protected.go` at 73.2.

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

- numbered Go file count before Slice 50: 545
- numbered Go file count after Slice 50: 536
- spec drift: pass; plan/tasks were updated to record the MI-driven split
- constitution drift: pass; no Node/JS/tooling dependencies added
- product drift: pass; preview remains read-only and verdict-free
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
