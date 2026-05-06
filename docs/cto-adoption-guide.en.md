# sdp-trace CTO Adoption Guide

`sdp-trace` is a control layer for existing AI-assisted delivery. It does not
replace your harness, prompts, agents, CI, review process, or repository
templates.

In the current Block 12 implementation, "control layer" means evidence capture,
explicit missing telemetry, and CI-witnessed reporting. It does not yet mean
automatic merge blocking, org-wide degradation dashboards, external notarized
proof, or guaranteed detection of every unwrapped agent run.

## What The CTO Gets

For every repository and commit, the organization can see:

- which agent or human workflow was observed;
- which evidence contract was expected;
- which commands and test/build steps were observed;
- which artifacts were retained as digests;
- which evidence is missing;
- whether the trace is only local, CI-witnessed, or externally witnessed;
- whether the local contract gate passes, whether CI witness is present, and
  whether audit-grade evidence is still `cannot_verify`.

There is no opaque score. Missing telemetry stays visible.

## Implementation Model

The adoption path is sidecar-first:

```text
existing harness / agent / prompt flow
        |
        v
sdp-trace wrap
        |
        v
local run artifacts
        |
        v
report + local gate
        |
        v
CI witness
        |
        v
CTO/team evidence package per repo and commit
```

The first integration contract is small:

- put an expected-evidence contract in each repo;
- wrap the existing harness command;
- generate report/gate artifacts in CI;
- add a CI witness record;
- keep local-only and missing telemetry states explicit.

If an agent or developer does not run through `sdp-trace wrap` or an adapter,
`sdp-trace` cannot see that local work directly. The first detectable signal is
at the expected evidence boundary: CI/report/gate will show that required run
artifacts or contract evidence are missing. That is useful control posture
evidence, but it is not the same as a complete agent activity log.

Minimal CI command sequence:

```text
sdp-trace report --out .sdp-trace-report --contract <contract> .sdp-trace-runs
sdp-trace gate --out .sdp-trace-report/gate-result.json --contract <contract> .sdp-trace-runs
sdp-trace witness --kind github-actions --out .sdp-trace-report/ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs
sdp-trace gate --out .sdp-trace-report/gate-result.json --contract <contract> --witness .sdp-trace-report/ci-witness.json .sdp-trace-runs
```

## Trust Levels

- `local_observed`: useful for reconstruction, not gate-grade.
- `ci_witnessed`: CI has bound the evidence package to a repository, commit,
  workflow, job, and run id.
- `external_witnessed`: future profile for independent timestamp/log witness.

CI witness improves buyer value because it prevents a local-only story from
being mistaken for a gate-grade trace. It still does not prove agent honesty or
release quality by itself.

For GitHub Actions, the workflow must grant OIDC access (`id-token: write`).
Without OIDC, `sdp-trace witness` records `cannot_verify` rather than pretending
that environment variables are enough.

Block 12 CI witness is not external trust. It is a CI-generated JSON artifact
that binds report/run digests to GitHub Actions OIDC claims when generated in a
protected workflow. It is not a public transparency log, DSSE envelope, or
court-ready signed timeline. Do not accept a witness file committed by an agent
or developer as authority; generate it inside CI and store it as a protected CI
artifact.

Policy interpretation should follow this table:

| State | What it can support | What it cannot support |
|---|---|---|
| `local_observed` | Local reconstruction and developer feedback | Merge/release trust by itself |
| `ci_witnessed` | CI-bound evidence package for a repo/commit | Agent honesty, test sufficiency, external audit proof |
| `external_witnessed` | Future external timestamp/log profile | Not implemented in Block 12 |

## What Is Not Visible Yet

Block 12 does not yet provide:

- automatic detection that an agent was used outside the wrapper;
- internal tool-call telemetry for harnesses without an adapter;
- raw prompt/model response capture;
- file mutation or VCS event capture beyond current recorder events;
- signed timeline or append-only transparency log;
- automatic degradation analytics across repositories;
- a dashboard or query surface beyond generated artifacts and existing local
  query commands.

These are explicit gaps, not hidden pass states.

## Handoff To Engineering

Give the team lead:

- the expected-evidence contract template;
- the wrapper command;
- the CI witness command;
- the report directory policy;
- the rule that `cannot_verify` is not a pass.

The CTO should review the generated `.sdp-trace-report/` per repo and commit,
not raw JSON scattered across developer machines.

For a multi-repository rollout, require teams to publish the CI-generated report
directory as a retained CI artifact. Start by tracking the ratio of
`local_observed`, `ci_witnessed`, `cannot_verify`, and missing evidence per
repository. That is the first truthful degradation signal; broader
`sdp-report` analytics are a follow-up product layer.
