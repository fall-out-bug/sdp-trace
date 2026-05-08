# Block 19 Implementation Evidence

Status: implementation evidence, strict implementation review, and PR-level
review recorded. Source-bound release proof regeneration remains open outside
Block 19 implementation closure.

## Scope

Block 19 adds adapter capture-depth assessment without changing `gate`,
`managed-harness`, or `forensic-retention` semantics.

Implemented surfaces:

- `internal/adaptercapture`: standalone evaluator for `adapter_capture`
  assessment facts.
- `sdp-trace assess --profile adapter-capture`: run input, deterministic exit
  behavior, JSON output, preview, and explain.
- `sdp-trace query --query capture-depth`: read-only capture-depth summary with
  no top-level pass/fail policy verdict.
- `schema/assessment-result.schema.json`: Block 19 assessment-result shape.
- `schema/adapter-event.schema.json` and
  `schema/adapter-capture-run.schema.json`: portable adapter event and run
  input contracts.
- `examples/block19-adapter-capture`: committed assessment-result fixtures for
  same-chain and adapter-bundle binding, missing telemetry, unsupported
  observers, gateway `not_integrated`, late events, test-provenance limits,
  source correlation, task supersession attribution, provider-neutral refs,
  event-level provider-ref safety, conflicting events, and capture-depth
  overclaim rejection.
- `docs/flight-recorder.md` and `docs/harness-integration.md`: adapter capture
  boundary, query surface, and safety guidance.

## Drift Handling

The implementation keeps Block 19 independent from Block 17 by default:
`adapter_capture` can report evidence rows for managed diagnostics, but it does
not authorize adapter identity, prove managed enrollment, or satisfy managed
witness binding by itself.

The implementation also keeps Block 18 redaction/retention boundaries active:
adapter-sensitive metadata is checked for forbidden persisted markers, and
preview/query/explain surfaces must not render raw payloads or credential-like
references.

## Local Verification

Commands run from `/Users/fall_out_bug/projects/vibe_coding/sdp-trace`:

```bash
rtk go test ./internal/adaptercapture ./cmd/sdp-trace ./internal/query
rtk jq empty schema/*.json examples/block19-adapter-capture/*.json
rtk go test ./...
rtk git diff --check
rtk rg -n "TODO|FIXME" <changed-files>
```

Observed states:

- Targeted Go tests: pass, 71 tests across 3 packages.
- JSON parse checks: pass.
- Full Go tests: pass, 147 tests across 14 packages.
- Whitespace diff check: pass.
- TODO/FIXME scan: no code debt markers in changed Block 19 implementation;
  hits are this verification note and the existing `AGENTS.md` policy text.

## Strict Review

Implementation review ran separate code/correctness, tracing/evidence, and
privacy/requirements planes. Initial major findings around FR-103 query facts,
event-level provider-ref safety, and SC-044 fixture breadth were accepted and
fixed. Final focused re-review reported no remaining critical or major findings.

PR-level review ran on PR #9 across code/correctness, tracing/evidence, and
privacy/requirements planes. Final PR-level review reported no critical or
major findings. GitHub CI is `not_assessed`: PR #9 reports no checks.

## Remaining Open Work

- GitHub CI for PR #9 remains historical `not_assessed`, but the repository now
  has tracked CI/check policy follow-up closure under T170.
- Source-bound release proof regeneration is not complete in this evidence
  note. Follow-up remains open as T169; the 2026-05-07 drift audit records the
  current source-bound verifier state as `fail`.
- PR-level minor review notes were accepted as non-blocking for Block 19 merge.
  Follow-up coverage/clarity work is closed under T171.
- The general process rule that recorded trust drift must become tracked work is
  closed under T172.
- Closed-task and block drift audit is recorded under T173. Historical
  Node/npm/script verifier reference cleanup remains open under T174.
