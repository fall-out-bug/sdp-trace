# Block 13B: Capture Boundary, State Taxonomy, And DX Baseline

Status: completion spec and ledger for follow-on implementation.

Parent artifact:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13-product-gap-closure-roadmap.md`

## Scope

Block 13B defines what `sdp-trace` may observe, how absent observation is
named, and what local diagnostics must explain before later gate blocks rely on
the evidence.

It does not implement protected enforcement, signed checkpoints, managed
harness fail-closed behavior, broad adapter telemetry, query dashboards, or
external audit proof.

## Boundary Spec

Gate contracts may require only evidence that maps to an active observation
boundary. Missing observers must produce an explicit state; they must not be
inferred from agent prose, git history, CI success, or checked-in JSON.

| Boundary | Required identity | Allowed observations | Trust cap without later blocks | Missing-boundary state |
|---|---|---|---|---|
| Process wrapper | local run id, wrapper version, command descriptor, working directory digest | lifecycle, argv descriptor, exit code, stdout/stderr digest descriptors, artifact digests | `local_observed` | `missing_telemetry` when required run is absent |
| Adapter socket/API | adapter id, harness id, registration profile | harness lifecycle, tool calls, model-call declarations, file mutations, task state | `harness_observed` or `agent_reported` depending on authority | `not_integrated`, `unsupported`, or `suppressed` |
| Tool-level wrapper | tool descriptor, parent run id | selected command invocations and result digests | `local_observed` | `unsupported` or `missing_telemetry` |
| VCS/PR observer | provider-neutral repo, commit, branch, PR/MR reference | source and review linkage | `vcs_observed` | `not_integrated` |
| CI observer | provider identity, workflow/job/run id, source ref, artifact digests | CI execution identity and report package binding | `ci_witnessed` | `cannot_verify` when identity or artifact binding is incomplete |
| External witness | witness profile, timestamp/log reference, checkpoint digest | append-only or third-party existence reference | `external_witnessed` | `not_integrated` until Block 22 or customer profile exists |

Boundary rules:

- A weaker boundary cannot upgrade itself by self-reporting a stronger state.
- `ci_witnessed` binds artifacts to CI identity; it does not prove agent
  honesty, test sufficiency, or external audit readiness.
- External policy consumers may block on emitted facts; `sdp-trace` must not
  emit native merge, release, readiness, or degradation verdicts.
- Unmanaged harnesses remain valid observation-mode targets. Their unsupported
  internals are explicit gaps, not implementation failures.

## State Taxonomy

These states are product states, not prose labels. Implementation blocks should
make them machine-enumerable before relying on them in gate output.

| State | Meaning | User action | May pass a required-evidence gate? |
|---|---|---|---|
| `pass` | Required verifier check succeeded for the selected profile. | Retain evidence and continue. | Yes, for that profile only. |
| `fail` | Verifier found evidence contradicting the requirement. | Fix input, contract, or workflow. | No. |
| `cannot_verify` | Verifier lacks required identity, binding, freshness, or artifact data. | Supply the missing verifier input or lower the claim. | No. |
| `not_assessed` | No verifier assessment was selected or implemented for this question. | Select a supported profile or leave the gap visible. | No. |
| `missing_telemetry` | Required observed evidence is absent from active run/report artifacts. | Run through the boundary or change the contract. | No. |
| `not_integrated` | A boundary could provide evidence, but no integration is configured. | Configure the observer or mark it out of scope. | No. |
| `unsupported` | Current product cannot observe the requested evidence type for this boundary. | Remove the requirement or schedule product work. | No. |
| `suppressed` | Evidence was intentionally withheld by policy, redaction, or harness behavior. | Inspect suppression authority and retention policy. | No unless the profile explicitly allows the suppression. |
| `offline_dev` | Local work ran without a network witness requirement being satisfiable. | Treat as local-only and rerun in CI for witness state. | No for CI or external witness requirements. |
| `local_observed` | Evidence came from local wrapper or local verifier output. | Use for reconstruction and developer feedback. | No for protected or audit-grade profiles. |
| `harness_observed` | Harness adapter emitted the observation. | Check adapter authority before using for gates. | Profile-dependent. |
| `agent_reported` | Agent or harness reported a fact without independent observation. | Use as context only unless corroborated. | No for executed-test or protected-gate claims. |
| `vcs_observed` | VCS/PR observer supplied source or review linkage. | Correlate with run and CI evidence. | Profile-dependent. |
| `ci_witnessed` | CI observer bound artifacts to CI identity and source reference. | Retain CI artifact and source binding. | Yes for CI-witness profile only. |
| `external_witnessed` | External append-only or customer-PKI witness bound a checkpoint. | Retain external witness reference. | Yes only for implemented external profile. |

## Doctor Acceptance Criteria

`doctor` or equivalent setup diagnostics should be accepted only if it:

- checks wrapper availability, output directory writeability, contract parse,
  expected-evidence references, and report directory policy;
- detects whether the selected CI witness profile has its required identity
  inputs, including GitHub Actions OIDC for the GitHub profile;
- distinguishes `unsupported`, `not_integrated`, `missing_telemetry`,
  `offline_dev`, and `cannot_verify` in user-facing output;
- reports the selected observation mode and trust cap before any gate advice;
- emits deterministic machine-readable output for identical inputs;
- avoids raw prompts, raw model responses, raw source snippets, raw stdout,
  raw stderr, credentials, secrets, and OIDC request tokens.

## Preview Acceptance Criteria

Local preview must run before artifact write and show:

- command descriptor and allowlisted argv basenames, not sensitive raw args;
- output directories and artifact categories that would be written;
- retention mode for each category: `digest_only`, `sanitized_excerpt`,
  `encrypted_raw_ref`, `external_artifact_ref`, `not_assessed`, or
  insufficient for the selected profile;
- selected evidence contract and required evidence ids;
- active boundaries and unsupported or unintegrated boundaries;
- offline implications when the selected profile needs CI or external witness.

Preview must not claim that a run will pass. It can only state what would be
captured, retained, suppressed, or left unassessed.

## Overhead Measurement Protocol

Purpose: measure wrapper friction before the roadmap claims acceptable DX.

Protocol:

1. Select one demo command that already runs without `sdp-trace`.
2. Run at least five warm baseline executions and record wall time, exit code,
   and artifact size if the command writes artifacts.
3. Run at least five wrapped executions with the same command, same working
   tree state, and same environment class.
4. Record median wall-time delta, p95 wall-time delta, wrapper artifact bytes,
   and number of emitted files.
5. Repeat after enabling preview/doctor only if those commands are part of the
   developer path being measured.

Initial budget:

- wrapper median wall-time overhead target: less than or equal to 5 percent or
  500 ms, whichever is larger;
- preview and doctor target: less than or equal to 2 seconds each on the demo
  repo;
- artifact count and byte growth must be reported, not hidden behind a score.

If the budget is missed, Block 13B may still continue, but the state is
`not_assessed` or `cannot_verify` for acceptable-overhead claims until the
measurement and disposition are recorded.

## Self-Trace, Provenance, And Evidence Checklist

Before Block 13B completion is claimed, the implementation owner must record:

- spec delta referencing this file and the roadmap;
- trace entries for boundary, state taxonomy, doctor, preview, retention floor,
  determinism, and overhead measurement behavior;
- evidence contract fixtures for absent run, unsupported observer,
  unintegrated adapter, suppressed evidence, offline local run, and missing CI
  witness identity;
- provenance records naming actor, command, source commit, run id, report
  directory, CI identity when present, and artifact digests;
- verifier output state as `pass`, `fail`, `cannot_verify`, or `not_assessed`
  for each acceptance criterion;
- review ledger with critical, major, and minor findings plus disposition;
- no source-bound trust claim unless the relevant files are bound to a clean
  immutable source commit and live verifier output.

## No-Overclaim Notes

- Prose in this file is not machine proof.
- Checked-in JSON is not authority unless live-verified or externally signed.
- Dirty worktree verification is local structural evidence only.
- `local_observed` is not protected enforcement.
- `ci_witnessed` is not `external_witnessed`.
- `audit_grade_gate` remains `cannot_verify` until an external witness profile
  exists and verifies.
- Suppressed evidence is not equivalent to safe redaction unless a verifier
  checks the retained artifacts against the retention profile.
- Cross-repository degradation output must expose numerator, denominator,
  dimensions, time window, stale input handling, and `not_assessed` counts; it
  must not collapse them into a health score.

## Ledger

| Area | Required artifact | Current Block 13B doc state |
|---|---|---|
| Boundary spec | This file, Boundary Spec | drafted |
| State taxonomy | This file, State Taxonomy | drafted |
| Doctor acceptance | This file, Doctor Acceptance Criteria | drafted |
| Preview acceptance | This file, Preview Acceptance Criteria | drafted |
| Overhead protocol | This file, Overhead Measurement Protocol | drafted |
| Self-trace checklist | This file, Self-Trace, Provenance, And Evidence Checklist | drafted |
| No-overclaim notes | This file, No-Overclaim Notes | drafted |
| Product implementation | Go code and tests | not_assessed; Worker C scope excludes Go code |
| Live verifier proof | Go-first verification output | not_assessed until implementation exists |
| External trust | External witness profile | not_integrated |
