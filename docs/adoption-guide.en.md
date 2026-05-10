# sdp-trace Adoption Guide

`sdp-trace` adds evidence records beside an existing AI-assisted delivery flow:
what happened, which evidence exists, what is missing, and who owns the next
human decision. It does not replace your harness, prompts, agents, CI, review
process, repository templates, or release governance.

The current pilot surface means trace capture, explicit missing telemetry,
assessment profiles, advisory/protected gate facts, CI/customer witness profiles,
forensic query packs, cross-repository posture export, and local source-bound
release proof. It does not mean automatic merge blocking, production release
approval, external audit proof, or guaranteed detection of every unwrapped
agent run.

## What It Provides

For every repository and commit, the organization can inspect:

- which agent or human workflow was observed;
- which task, command, model, harness, and source context were recorded;
- which evidence contract or assessment profile was expected;
- which artifacts were retained, redacted, or digest-only;
- which evidence is missing, `not_assessed`, or `cannot_verify`;
- whether a run is local only, CI-witnessed, customer-PKI witnessed, or still
  documentation/fixture-only;
- whether source-bound local release proof passed without claiming external
  production trust.

There is no opaque score. Missing telemetry stays visible.

## How This Complements CI Logs, Git Diff, And Review Comments

CI logs show command output. Git diff shows file changes. Review comments show
human discussion. `sdp-trace` adds a portable evidence contract across those
surfaces:

- provenance links who or what produced evidence;
- trace runs preserve command and task context;
- assessment profiles explain why evidence is enough, missing, stale, or
  unverifiable;
- witness records bind selected evidence to CI or customer authority when the
  profile has enough data;
- release proof checks manifest subjects against a source commit instead of
  trusting prose.

This is still evidence, not a policy decision. CI, release management, customer
governance, or another external policy consumer decides what to block.

## Implementation Model

The adoption path is sidecar-first:

```text
existing harness / agent / prompt flow
        |
        v
sdp-trace wrap / adapter events
        |
        v
.sdp-trace-runs/
        |
        v
report, query, assess, gate facts
        |
        v
CI or customer witness where available
        |
        v
evidence package per repo and commit
```

Minimum command sequence:

```text
go run ./cmd/sdp-trace wrap --name <workflow-name> --output-dir .sdp-trace-runs/<run-id> -- <existing command...>
go run ./cmd/sdp-trace report --out .sdp-trace-report .sdp-trace-runs
go run ./cmd/sdp-trace gate --out .sdp-trace-report/gate-result.json .sdp-trace-runs
go run ./cmd/sdp-trace witness --kind github-actions --out .sdp-trace-report/ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs
```

If an agent or developer does not run through `sdp-trace wrap` or an adapter,
`sdp-trace` cannot see that local work directly. The detectable signal is at the
expected evidence boundary: required run artifacts, adapter events, witness
bindings, or profile inputs are missing and must remain `missing_telemetry`,
`not_assessed`, or `cannot_verify`.

## Current Profiles And Boundaries

| Surface | What it supports now | Caveat |
| --- | --- | --- |
| `adapter-capture` | Checks adapter event coverage and overclaim risk. | Missing adapter events do not prove no agent was used; they prove the profile lacks evidence. |
| `managed-harness` | Assesses managed policy, adapter registry, run, and witness evidence. | Emits verifier facts; external CI or policy decides block/allow. |
| `forensic-retention` | Checks whether retained/redacted evidence supports reconstruction. | Digest-only or unresolved redaction can block forensic claims. |
| `gate` | Emits advisory/protected gate facts and reasons. | Not a native merge, release, readiness, degradation, override, or risk decision. |
| `witness` | Emits witness artifacts for `github-actions`, `gitlab-ci`, `buildkite`, and `customer-pki`. | CI/customer witness is not external production trust unless the external trust profile passes. |
| `release-proof` | Verifies source-bound local release manifests against a source commit. | `source_bound_local_release` is not `external_production_trust`. |
| Air-gapped guidance | Uses customer policy/private-equivalent evidence patterns. | There is no `witness --kind air-gapped`; unsupported evidence stays `not_assessed` or `cannot_verify`. |

## What To Inspect

- `.sdp-trace-report/summary.json`: run and report summary.
- `.sdp-trace-report/evidence-table.json`: observed evidence rows.
- `.sdp-trace-report/missing-telemetry.json`: required evidence not observed.
- `.sdp-trace-report/gate-result.json`: advisory/protected gate facts and reasons.
- `.sdp-trace-report/ci-witness.json` or another witness artifact: CI/customer binding state.
- `.sdp-trace-runs/<run-id>/`: raw run package, subject to retention and redaction policy.
- `query-pack` output: incident or forensic reconstruction package.
- `release-proof` output: source-bound local release state.
- Workflow docs: spec, plan, tasks, evidence, decisions, and deferred gaps from
  SpecKit, gsd, Superpowers, Oh My OpenAgent, a ticket tracker, or the team's
  custom planning flow.

## Interpreting Missing States

- `not_assessed`: the state was intentionally outside the run scope. Ask whether
  the scope is acceptable or require a follow-up profile.
- `cannot_verify`: the verifier attempted the check but lacked required
  evidence, environment, or consistency. Treat it as fail-closed for trust
  claims.
- `fail`: evidence contradicted the selected profile.
- `observed` or `pass`: the selected local/profile check concluded, but only
  inside its stated trust scope.

## Privacy And Non-Capture

The current pilot surface should not require committed raw customer source, private prompts,
credentials, provider tokens, or raw logs. Prefer digests, sanitized excerpts,
encrypted external references, and explicit redaction notes. If a team needs raw
capture for an incident, it must be a separately approved retention/redaction
profile with a human owner.

## Rollout Inputs

Before rollout, make sure the following are explicit in the repository:

- the expected evidence contract or profile inputs;
- the wrapper or adapter integration command;
- the report and retention policy;
- the witness profile to use, if any;
- the rule that `not_assessed` and `cannot_verify` are not passes;
- the repository's retained report, witness, and assessment artifacts for the
  exact run under review.

For a multi-repository rollout, require teams to publish the report directory
and witness artifacts as retained CI/customer-policy artifacts. Track the ratio
of `observed`, `pass`, `fail`, `not_assessed`, and `cannot_verify` states by
repository. Broader dashboards and policy decisions belong outside `sdp-trace`.
