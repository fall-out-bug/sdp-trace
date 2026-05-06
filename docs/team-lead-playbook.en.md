# Team Lead Playbook

Use `sdp-trace` when your team already has an AI coding workflow and needs a shared quality contract.

This playbook covers the current Block 12 product surface: process wrapping,
local report/gate artifacts, and GitHub Actions OIDC witness. It does not yet
promise full harness internals, automatic file mutation tracing, fail-closed
managed harness enforcement, or external signed timelines.

## Daily Use

1. Define the scope.
2. Capture provenance: human, agent, model, tools, and commands.
3. Attach evidence: tests, CI, review comments, files, and diffs.
4. Record accountability: human-held DRI, approver, risk owner, and escalation path.
5. Package an assessment input with evidence, observations, movement data, and `not_assessed` gaps.
6. Record any gate verdict as an external verdict input produced by `sdp-gate` or another policy consumer.

Current capture boundary:

- `wrap` observes the wrapped process lifecycle and command-level events.
- It does not automatically see internal tool calls inside a harness unless the
  harness emits adapter events.
- It does not prove that no one ran an agent outside the wrapper.
- Missing expected evidence must remain `missing_telemetry` or `cannot_verify`,
  not pass.

## Team Defaults

Agree on:

- required evidence per change type
- what external policy blocks merge
- who may approve or override in the policy layer
- which harnesses are supported
- what `not_assessed` means

## Repository Setup

For each repository, add:

- an expected-evidence contract owned by the team;
- a `.sdp-trace-runs/` output location for wrapped local/CI runs;
- a `.sdp-trace-report/` report location for CI artifacts;
- CI steps for `report`, `gate`, `witness`, and `gate --witness`.

The implementation sequence is:

```text
sdp-trace wrap --name <workflow-name> --contract <contract> --output-dir .sdp-trace-runs/<run-id> -- <existing command...>
sdp-trace report --out .sdp-trace-report --contract <contract> .sdp-trace-runs
sdp-trace gate --out .sdp-trace-report/gate-result.json --contract <contract> .sdp-trace-runs
sdp-trace witness --kind github-actions --out .sdp-trace-report/ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs
sdp-trace gate --out .sdp-trace-report/gate-result.json --contract <contract> --witness .sdp-trace-report/ci-witness.json .sdp-trace-runs
```

For GitHub Actions, enable OIDC in the workflow:

```text
permissions:
  id-token: write
  contents: read
```

Without OIDC, `ci_witness_gate` remains `cannot_verify`.

Do not commit `.sdp-trace-report/ci-witness.json` from a developer machine as
trusted evidence. Generate it in CI and retain it as a CI artifact.

## Privacy And Redaction

Default Block 12 report artifacts retain command metadata and stdout/stderr
digests, not raw stdout/stderr bodies. OIDC request tokens are read only to
request GitHub's OIDC token and must not be written to `.sdp-trace-runs/` or
`.sdp-trace-report/`.

Before rollout, agree on:

- whether prompts, source snippets, or tool payloads are allowed at all;
- which outputs must remain digest-only;
- who may approve raw capture for a narrow incident window;
- how redaction decisions are recorded for later investigation.

If a team needs raw prompt, source, or model-response capture, that is a
separate adapter/redaction profile. Do not treat it as enabled by this playbook.
Any future raw-capture profile must redact before persistent write. Block 12
does not provide a raw-capture mode and therefore does not make a post-hoc
redaction safety claim.

## Emergency Changes

Do not hide emergency bypasses. If production urgency requires shipping with
missing telemetry, record the change as a policy override in the external policy
layer and keep the `cannot_verify` or `missing_telemetry` state visible in the
report. A bypass is acceptable only when the organization can later see who
approved it, why, and which evidence was missing.

Block 12 does not yet have a native `policy_override_requested` trace event.
Until that exists, the override record must live in the external policy/change
management system and must reference the report artifacts.

## Offline And Failure Modes

- `wrap`, `report`, and the local contract gate can run without network access.
- `witness --kind github-actions` requires GitHub Actions OIDC and therefore
  cannot pass offline.
- If `witness` exits `3`, inspect `ci-witness.json` for `reason` and
  `missing_identity_fields`.
- If `gate` fails, inspect `gate-result.json`, `missing-telemetry.json`, and the
  per-run `verifier/` outputs.
- Use `sdp-trace dry-run --contract <contract> -- <command...>` to preview the
  command and contract without writing run artifacts.

Gate debugging checklist:

1. Check `gate-result.json` for `local_gate`, `ci_witness_gate`, and
   `audit_grade_gate`.
2. Check `required_evidence` versus `observed_evidence`.
3. Check `missing-telemetry.json` for absent contract evidence.
4. Check `ci-witness.json` for `reason`, `missing_identity_fields`, and OIDC
   state.
5. Check each run's `verifier/` output before changing the contract.

## Reading External Verdicts

External verdicts may use values such as `pass`, `warn`, `fail`, or `not_assessed`, but they are not native `sdp-trace` decisions.

`not_assessed` is not a pass. Missing evidence must stay visible in the assessment input.

## What To Review

Review these artifacts per repo and commit:

- `.sdp-trace-report/summary.json`
- `.sdp-trace-report/evidence-table.json`
- `.sdp-trace-report/missing-telemetry.json`
- `.sdp-trace-report/gate-result.json`
- `.sdp-trace-report/ci-witness.json`

The important questions are:

- Did the expected evidence contract match the work?
- Are all required evidence ids observed?
- Is the trace only `local_observed` or `ci_witnessed`?
- Is any missing telemetry hidden? It should not be.
- Is `audit_grade_gate` still `cannot_verify` because external witness is not
  implemented yet?

## Retention

For investigations, keep `.sdp-trace-report/` and the matching
`.sdp-trace-runs/` for at least the same period as CI logs and review records.
If your organization has incident or audit retention requirements, store the
report directory in an immutable artifact store. Block 12 does not implement
retention by itself.

## Current Follow-Up Product Gaps

- fail-closed managed harness enforcement;
- signed timeline / DSSE / external transparency witness;
- cross-repository dashboard and degradation analytics;
- richer query surface for monthly or incident-wide investigation;
- redaction audit trail beyond digest-only defaults;
- support for non-GitHub CI witness profiles.
- native `policy_override_requested` trace event;
- measured wrapper overhead and latency budgets.
