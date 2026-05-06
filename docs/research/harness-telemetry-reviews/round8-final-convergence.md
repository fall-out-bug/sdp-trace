# Round 8: Final Convergence

Status: discussion draft; not committed
Date: 2026-05-05

Inputs:

- `docs/research/agentic-sdlc-evidence-substrate-v4-brief.md`
- `docs/research/harness-telemetry-reviews/round8-cto-buyer-minimax.md`
- `docs/research/harness-telemetry-reviews/round8-platform-harness-owner-glm.md`
- `docs/research/harness-telemetry-reviews/round8-ciso-adversarial-trust-opus.md`
- `docs/research/harness-telemetry-reviews/round8-staff-engineer-dx-kimi.md`
- `docs/research/harness-telemetry-reviews/round8-compliance-forensics-gemini.md`

This file is a human consolidation of Socratic persona outputs. It is
not source-bound proof, not product closure evidence, and not a trusted
release claim.

## Convergence Result

Converged for v0 implementation planning.

All five personas returned:

- verdict: `ACCEPTABLE_WITH_GAPS`;
- can start v0 implementation: yes;
- critical blockers before implementation: no.

The remaining gaps are implementation tasks, explicit V0 limitations, or
future product profiles. They do not require another product-brief round
before drafting the implementation plan.

## Role Results

| Persona | Model | Verdict | Start v0? | Critical Blockers? |
| --- | --- | --- | --- | --- |
| CTO Buyer | MiniMax M2.7 | `ACCEPTABLE_WITH_GAPS` | yes | no |
| Platform Owner | GLM 5.1 | `ACCEPTABLE_WITH_GAPS` | yes | no |
| CISO | Claude Opus Latest | `ACCEPTABLE_WITH_GAPS` | yes | no |
| Staff Engineer | Kimi k2p6 | `ACCEPTABLE_WITH_GAPS` | yes | no |
| Forensics Lead | Gemini Pro Latest | `ACCEPTABLE_WITH_GAPS` | yes | no |

## Accepted V0 Scope

The converged v0 product scope is:

- wrapper-first local observation;
- wrapper composition for existing harness commands;
- Unix socket adapter event ingress;
- optional adapters with capability declarations;
- JSON expected evidence contract;
- pre-write redaction;
- local-only trust explicitly not gate-grade;
- CI witness path for gate-grade trust;
- MissingEvidenceTable;
- `verify`, `query`, `explain`, and `export`;
- audit bundle;
- explicit no-run / missing-run states where VCS/CI can observe them.

## Residual Implementation Tasks

These should feed the implementation plan:

1. Define the concrete V0 JSON schemas for:
   - canonical event;
   - adapter registration;
   - expected evidence contract;
   - MissingEvidenceTable;
   - verifier result;
   - audit bundle.
2. Add adapter event payload definitions for:
   - `harness_identity_observed`;
   - `tool_call_observed`;
   - `model_identity_observed`.
3. Specify CI witness key custody:
   - OIDC job identity or KMS profile;
   - no keys readable by the agent workspace/process;
   - `witness_independence` output.
4. Define CI preflight ownership for `expected_run_absent`.
5. Define default contract fallback for frictionless local `wrap`.
6. Define redaction rule loading order:
   - workspace config;
   - contract profile;
   - built-in defaults.
7. Add emergency override CLI:
   - e.g. `--override-reason <text>`;
   - unsigned declarations remain `human_declared`;
   - signed overrides become `human_signed + partial`.
8. Define `audit-bundle` packaging:
   - event chain;
   - verifier result;
   - MissingEvidenceTable;
   - witness assertions;
   - retention manifest;
   - detached or embedded signature profile.
9. Define test artifact resolution:
   - how digest refs point back to CI artifacts after incident delay.
10. Define platform boundary:
   - V0 targets Unix/macOS;
   - Windows transport is future profile.

## Accepted V0 Limitations

- No retroactive attach to already-running agent processes.
- No raw prompt/response capture by default.
- No org-wide degradation dashboard.
- No local-only audit-grade proof.
- No proof of deleted local runs before VCS/CI/preflight observes an
  expected run.
- No full multi-harness SDK.
- No Windows adapter transport in V0 unless explicitly added later.

## Recommended Next Step

Stop product-brief iteration and draft the v0 implementation plan.

The plan should start with schemas and verifier semantics, not with a
demo harness. The first buildable slice should prove:

1. `sdp-trace wrap` records a local command run with redaction and
   MissingEvidenceTable.
2. Verifier emits four-axis output without local overclaim.
3. Tamper causes `integrity_audit`.
4. CI witness can sign verifier result for a pinned contract and source
   commit.
5. Audit bundle export is reproducible from committed fixtures.
