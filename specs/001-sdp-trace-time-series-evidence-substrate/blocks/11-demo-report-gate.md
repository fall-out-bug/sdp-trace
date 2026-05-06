# Block 11: Demo-Repo Report And Local Gate

Status: implementation draft; MiniMax-M2.7 review required.

Review status:

- MiniMax-M2.5 review is invalid for repo trust-review; rerun required with
  MiniMax-M2.7 before this block can be treated as reviewed.
- GLM-5.1 SpecKit/DX review: major issues accepted.

Refactor resolution:

- Product code and specs must not hardcode OpenCode, GSD, Bazel, Kotlin, or
  any demo-project evidence kind.
- Demo-specific evidence expectations belong in the external demo repository's
  expected-evidence contract.
- `sdp-trace` report/gate behavior must be driven by generic verifier output
  and contract-declared evidence ids.
- Block 11 no longer defines product-level demo classifiers. Report/gate
  evidence classification is contract-driven.

## Problem

The external demo repository now produces valid local `sdp-trace` run
directories, but the useful story is still trapped in raw run artifacts. A
buyer should not need a separate slide deck, dashboard, or manual narration to
understand what happened.

The demo repository itself must contain the evidence view:

- what ran;
- what passed;
- what is only local evidence;
- what is missing for audit-grade trust;
- whether a local contract gate passes;
- whether an audit-grade release gate can be verified.

## Non-Goals

- No web UI.
- No opaque score.
- No CI/OIDC signing.
- No external witness simulation.
- No dependency on OpenCode, GSD, Bazel, GitHub, or any harness runtime inside
  `sdp-trace` product code.
- No Node.js, npm, JavaScript, TypeScript, `.mjs`, or shell scripts in the
  active product path.

## User-Facing Commands

### `sdp-trace report`

Usage:

```text
sdp-trace report --out <dir> [--contract <contract.json>] <runs-root-or-run-dir>
```

Behavior:

1. Discover run directories exactly like `validate-fixtures`:
   - if the argument itself contains `run.json`, treat it as one run;
   - otherwise inspect direct child directories containing `run.json`.
2. Verify every discovered run and write verifier artifacts back into each run.
3. Write these demo artifacts under `--out`:
   - `summary.json`
   - `timeline.md`
   - `evidence-table.json`
   - `missing-telemetry.json`
4. Do not copy raw stdout/stderr. Keep digest references only.
5. `--out` is required. Omission exits with usage error.
6. Verifier artifact writes are required. A read-only run directory is kept in
   the generated report as `cannot_verify`; the command still writes the
   report artifacts when `--out` is writable, then exits non-zero.

Required `summary.json` fields:

- `generated_at`
- `run_count`
- `observed_count`
- `failed_count`
- `cannot_verify_count`
- `not_assessed_count`
- `trust_scope`
- `audit_grade`
- `audit_grade_reason`
- `runs`

Allowed `summary.json` top-level values:

- `trust_scope`: `local_observed`
- `audit_grade`: `false`
- `audit_grade_reason`: non-empty string

The report-level `trust_scope` is the strongest common scope Block 11 can
prove. It must not be upgraded by command names or agent text.

Required per-run fields:

- `name`
- `run_id`
- `command`
- `exit_code`
- `closure_state`
- `result`
- `trust_scope`
- `completeness`
- `replayability`
- `stdout_digest`
- `stderr_digest`

Per-run allowed values are inherited from the verifier:

- `result`: `observed`, `fail`, `cannot_verify`, `not_assessed`
- `trust_scope`: `local_observed`
- `completeness`: `complete`, `partial`, `missing_telemetry`, `unknown`
- `replayability`: `full`, `partial`, `none`

Run classification fields:

- `kind`: matched contract evidence id, or `unmatched`
- `kind_reason`: short deterministic reason for the classification

Classification rules are intentionally simple and transparent. `sdp-trace`
does not infer harness/build-system semantics from command names. It matches
only evidence requirements declared by `--contract`.

Supported Block 11 evidence requirement shape:

```json
{
  "id": "agent_session_observed",
  "event_type": "command_started",
  "payload_field": "wrapper_name",
  "payload_equals": "agent-session"
}
```

Demo-specific ids and payload values belong in the external demo repository's
contract, not in product code.

Required `evidence-table.json` fields:

- `runs`: array of the per-run objects above.

Required `missing-telemetry.json` fields:

- `missing_audit_evidence`
- `missing_harness_evidence`
- `notes`

Required `timeline.md` format:

```text
# SDP Trace Timeline

| Run | Kind | Result | Trust Scope | Command | Exit |
|-----|------|--------|-------------|---------|------|
```

Rows are sorted lexicographically by run directory name.

### `sdp-trace gate`

Usage:

```text
sdp-trace gate --out <file> <runs-root-or-run-dir>
sdp-trace gate --out <file> --contract <contract.json> <runs-root-or-run-dir>
```

Behavior:

1. Verify every discovered run.
2. Evaluate a local contract gate:
   - every `required_evidence[].id` from the selected contract must be locally
     observed;
   - all included runs must be verifier `observed`;
   - all included runs must have closure state `completed`;
   - no run may be `fail` or `cannot_verify`.
3. `--out` is required. Omission exits with usage error.
4. Always evaluate audit-grade posture separately:
   - if there is no CI/OIDC witness and no external witness, audit-grade result
     is `cannot_verify`.

Required `gate-result.json` fields:

- `local_gate`
- `audit_grade_gate`
- `reasons`
- `required_evidence`
- `observed_evidence`
- `gate_conditions`
- `missing_audit_evidence`
- `runs`

Allowed gate values:

- `pass`
- `fail`
- `cannot_verify`

Gate `runs` uses the same per-run object schema as `summary.json`.

`required_evidence` contains only contract-declared evidence ids.
Gate-wide enforcement conditions such as `all_runs_observed` and
`all_runs_completed` are reported in `gate_conditions` and must not be mixed
with contract evidence ids.

Invalid or malformed run behavior:

- if any discovered run cannot be loaded or verified, `local_gate` is
  `fail`;
- `audit_grade_gate` remains `cannot_verify`;
- the malformed run must appear in `runs` with `result: cannot_verify` when a
  run name can be determined.

Audit evidence is fixed in Block 11 because this block does not implement
CI/OIDC or external witnesses:

```json
{
  "missing_audit_evidence": [
    "ci_oidc_witness",
    "external_witness_checkpoint"
  ]
}
```

This is not circular: it shows that the local demo can pass while the same
evidence package cannot be represented as audit-grade release evidence.

## Demo Repository Layout

Recommended output path:

```text
.sdp-trace-demo/
  summary.json
  evidence-table.json
  missing-telemetry.json
  timeline.md
  gate-result.json
```

These files are demonstration artifacts. They are not external trust evidence.

## Trust Semantics

Block 11 must keep the distinction explicit:

- local observed evidence can pass the local contract gate;
- local observed evidence cannot pass an audit-grade release gate;
- raw stdout/stderr may contain useful tool data, but until parsed and redacted
  it remains digest-only in this block;
- CI/OIDC witness and external checkpoint remain missing telemetry.

The report must not say or imply that local evidence is production/audit-grade.

## Acceptance Criteria

1. `go test ./...` passes.
2. `sdp-trace report --out <dir> --contract <contract> <runs-root>` creates all
   four report artifacts.
3. `sdp-trace gate --out <file> --contract <contract> <runs-root>` creates
   `gate-result.json`.
4. The external demo repo's contract-declared run set produces:
   - local gate: `pass`;
   - audit-grade gate: `cannot_verify`;
   - missing audit evidence includes CI/OIDC witness and external witness.
5. Generated artifacts contain no raw stdout/stderr payloads.
6. No Node/npm/JS/TS/shell product tooling is introduced.
7. Missing or malformed `--out` handling is tested.
8. Report/gate generated JSON validates with `jq empty`.
9. A tampered run in the input set makes the local contract gate fail.
