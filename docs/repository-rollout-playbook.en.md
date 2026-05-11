# Repository Rollout Playbook

Use `sdp-trace` when your team already has an AI coding workflow and needs a
shared evidence contract without replacing the current harness.

Use [Agent Onboarding](agent-onboarding.md) as the first link for coding agents.
This playbook is for wiring one repository into the evidence path.

This playbook covers the current pilot product surface: wrapping, task-linked
runs, previews, local verification, reports, queries, query packs, assessment
profiles, advisory/protected gate facts, CI/customer witness profiles,
source-bound release proof, and fixture validation. It does not promise full
harness internals, automatic detection of every bypass, external production
trust, or policy decisions.

## Daily Use

1. Define the spec, plan, task, and expected evidence.
2. Capture provenance: human, agent, model, tools, commands, and source context.
3. Run through `wrap`, `run`, or an adapter when the workflow can be observed.
4. Attach evidence: tests, CI, review comments, files, diffs, and retained artifacts.
5. Record accountability: human-held DRI, approver, risk owner, and escalation path.
6. Package report, query, assessment, gate facts, and witness artifacts where available.
7. Keep `not_assessed` and `cannot_verify` states visible.

Current capture boundary:

- `wrap` observes the wrapped process lifecycle and command-level events.
- Adapter profiles can assess richer harness events when the harness emits them.
- `sdp-trace` does not prove that no one ran an agent outside the wrapper.
- Missing expected evidence must remain `missing_telemetry`, `not_assessed`, or
  `cannot_verify`, not pass.

## Team Defaults

Agree on:

- required evidence per change type;
- which assessment profiles apply to which workflow;
- which external policy blocks merge or release;
- who may approve or override in the policy layer;
- which harnesses and CI systems are supported;
- what `not_assessed` means in customer handoff.

## Repository Setup

For each repository, add:

- an expected evidence contract owned by the team;
- `.sdp-trace-runs/` for wrapped local/CI runs;
- `.sdp-trace-report/` for report artifacts;
- optional adapter registry, managed policy, redaction policy, and witness policy files;
- CI steps for `report`, `gate`, and the selected `witness` kind.

Minimum implementation sequence:

```text
sdp-trace wrap --name <workflow-name> --output-dir .sdp-trace-runs/<run-id> -- <existing command...>
sdp-trace report --out .sdp-trace-report .sdp-trace-runs
sdp-trace gate --out .sdp-trace-report/gate-result.json .sdp-trace-runs
sdp-trace witness --kind github-actions --out .sdp-trace-report/ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs
```

Useful local checks:

```text
sdp-trace doctor
sdp-trace verify .sdp-trace-runs/<run-id>
sdp-trace explain .sdp-trace-runs/<run-id>
sdp-trace query --query missing-evidence .sdp-trace-runs/<run-id>
sdp-trace query-pack --pack forensics-basic-v1 --run .sdp-trace-runs/<run-id> --out query-pack.json
```

## Profiles To Use

| Need | Command |
| --- | --- |
| Adapter coverage and overclaim review | `sdp-trace assess --profile adapter-capture --out assessment.json --run .sdp-trace-runs/<run-id>` |
| Managed harness profile | `sdp-trace assess --profile managed-harness --out assessment.json --contract contract.json --run .sdp-trace-runs/<run-id> --adapter-registry registry.json --managed-policy policy.json --managed-witness witness.json` |
| Forensic retention profile | `sdp-trace assess --profile forensic-retention --out assessment.json --run .sdp-trace-runs/<run-id> --redaction-policy redaction.json` |
| Assessment explanation | `sdp-trace assess explain --assessment-result assessment.json` |
| Source-bound release proof | `sdp-trace release-proof --manifest contract-manifest.json --out release-proof.json` |

`managed-harness` emits verifier facts and exit behavior. It does not decide
merge/release readiness. Missing managed witness evidence generally remains
`cannot_verify`.

## Witness Profiles

Supported `witness --kind` values are:

- `github-actions`
- `gitlab-ci`
- `buildkite`
- `customer-pki`

For GitHub Actions, enable OIDC in the workflow:

```text
permissions:
  id-token: write
  contents: read
```

Without required identity or binding evidence, witness output must remain
`cannot_verify`. Do not commit a witness file from a developer machine as trusted
evidence. Generate it in CI or under the customer-approved PKI process and retain
it as a protected artifact.

Air-gapped evidence is not a separate command kind. Treat it as customer
policy/private-equivalent guidance: explicit authority policy, payload digest,
freshness or timestamp evidence, and retained audit references. If any required
piece is unavailable, record `not_assessed` or `cannot_verify`.

## Gate Debugging

`gate` output is verifier-derived evidence, not a native policy decision.

Debugging checklist:

1. Check `gate-result.json` for selected mode, required runs, required evidence, and reason rows.
2. Check `.sdp-trace-report/missing-telemetry.json` for absent contract evidence.
3. Check witness output for source, run, freshness, and identity binding state.
4. Check `assess explain` output for profile-specific conditions.
5. Check each run's verifier output before changing the contract.

## Privacy And Redaction

Default to safe-to-commit artifacts: metadata, digests, sanitized excerpts, and
external references. Do not commit raw customer source, private prompts,
credentials, provider tokens, authenticated URLs, OIDC request tokens, or raw
logs. If raw capture is required for an incident, define a retention/redaction
profile first and name a human owner.

## Emergency Changes

Do not hide emergency bypasses. If production urgency requires shipping with
missing telemetry, record the override in the external policy/change management
system and keep `cannot_verify` or `missing_telemetry` visible in the report.
`sdp-trace` records evidence; it does not approve the risk.

## What To Review

Review these artifacts per repo and commit:

- `.sdp-trace-report/summary.json`
- `.sdp-trace-report/evidence-table.json`
- `.sdp-trace-report/missing-telemetry.json`
- `.sdp-trace-report/gate-result.json`
- the selected witness artifact, such as `.sdp-trace-report/ci-witness.json`
- assessment result files for `adapter-capture`, `managed-harness`, or `forensic-retention`
- query-pack output for incident or forensic review
- release-proof output for source-bound local release claims

The important questions are:

- Did the expected evidence contract match the work?
- Are all required evidence ids observed?
- Is the trace only local, or is it bound to CI/customer evidence?
- Is any missing telemetry hidden?
- Which states are `not_assessed` or `cannot_verify`, and who owns follow-up?
- Is any paragraph accidentally treating local gate, witness, or release proof as external production trust?

## Retention

For investigations, keep `.sdp-trace-report/`, selected witness artifacts, and
the matching `.sdp-trace-runs/` for at least the same period as CI logs and
review records. If your organization has incident or audit retention
requirements, store the report directory in an immutable artifact store.

## Current Follow-Up Product Gaps

- native dashboards and policy decisions;
- guaranteed detection of every unwrapped agent run;
- external production trust without a selected passing external profile;
- universal air-gapped witness command;
- raw prompt/source/model-response capture without a separate redaction profile;
- measured wrapper overhead and latency budgets.
