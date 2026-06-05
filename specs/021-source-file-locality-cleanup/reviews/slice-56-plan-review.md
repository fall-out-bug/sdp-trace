# Slice 56 Plan Review

Status: pass

## Scope

Slice 56 is bounded to gate preview contract helper shards:

- `cmd/sdp-trace/gate_365_previewgatemode.go`
- `cmd/sdp-trace/gate_366_requiredrunids.go`
- `cmd/sdp-trace/gate_367_requiredevidenceidsforcli.go`

Planned cohesive file:

- `cmd/sdp-trace/gate_preview_contract.go`

Explicit exclusions:

- packet and PR review shards (`packet_031` onward)
- command-specific preview rendering outside contract-derived display helpers

## Behavior To Preserve

- `previewGateMode` defaults to `demo.GateModeObservation`.
- Advisory CI is selected when at least one required run has
  `demo.GateModeAdvisoryCI` and no required run has protected-future mode.
- Protected-future mode dominates advisory CI regardless of required-run order.
- Unknown/empty profiles do not change the selected mode.
- `requiredRunIDs` preserves required-run order and omits empty IDs.
- `requiredEvidenceIDsForCLI` preserves required-evidence order and omits empty
  IDs.
- No package boundary, dependency direction, or MI baseline change is planned.

## Planned Regression Evidence

- Add `TestPreviewGateModeSelection`.
- Add `TestRequiredRunIDsOmitEmptyAndKeepOrder`.
- Add `TestRequiredEvidenceIDsForCLIOmitEmptyAndKeepOrder`.

`TestPreviewGateModeSelection` must explicitly cover observation default,
advisory CI fallback, protected-future dominance regardless of order, and
unknown/empty required-run profiles not changing the selected mode.

Focused test existence and execution will use:

```text
go test ./cmd/sdp-trace -list '^(TestPreviewGateModeSelection|TestRequiredRunIDsOmitEmptyAndKeepOrder|TestRequiredEvidenceIDsForCLIOmitEmptyAndKeepOrder)$'
go test ./cmd/sdp-trace -run 'Test(PreviewGateModeSelection|RequiredRunIDsOmitEmptyAndKeepOrder|RequiredEvidenceIDsForCLIOmitEmptyAndKeepOrder)$'
```

## Review Lanes

- scope/correctness reviewer: multi_agent_v1,
  `019e8858-c8b3-7063-a2eb-192f2a8bbc77`; date:
  `2026-06-02T17:13:00+03:00`; prompt class:
  `plan scope/correctness review plus targeted re-review`; timeout:
  `600000ms`; retries: `1`; fallback: `none`; result `LGTM`
- trust/evidence reviewer: multi_agent_v1,
  `019e8858-ccec-7211-9d43-eaf682f92e18`; date:
  `2026-06-02T17:13:00+03:00`; prompt class:
  `plan trust/evidence review plus targeted re-review`; timeout: `600000ms`;
  retries: `1`; fallback: `none`; result `LGTM`
- maintainability/DX reviewer: multi_agent_v1,
  `019e8858-d21f-7ad0-a8be-8d1e3e48dbcc`; date:
  `2026-06-02T17:13:00+03:00`; prompt class:
  `plan maintainability/DX review`; timeout: `600000ms`; retries: `0`;
  fallback: `none`; result `LGTM`
- model/provider: `not_assessed` for all lanes because the harness does not
  expose provider-qualified model IDs
- requested external/provider-qualified lanes: `not_assessed` because no
  callable provider-qualified model surface is exposed in this session

## Plan Review Findings

- scope/correctness lane initial finding: major; unknown/empty profile handling
  was named in the review artifact but omitted from tasks and explicit focused
  evidence. Resolution: added the behavior to T021-3811/T021-3841 and the
  planned `TestPreviewGateModeSelection` assertions.
- trust/evidence lane initial finding: major; planned evidence did not
  explicitly require unknown/empty profile assertions. Resolution: added the
  explicit assertion requirement above.
