# Socratic Review: MVP Readiness Hardening

**Status**: Three-axis review complete; valid critical/major findings fixed in draft, approval pending
**Created**: 2026-05-10

## Required Review Plan

Run independent Socratic review before implementation. Required planes:

1. **Documentation/Product Promise**: Does the spec correctly target MVP handoff
   risks without expanding into a full docs rewrite or overclaiming readiness?
2. **Code Quality/CRAP Gate**: Is the proposed lint, complexity, and coverage
   gate realistic, measurable, and strict enough to satisfy the user's
   `CRAP < 5` intent without hiding current non-compliance?
3. **Trust Boundary**: Does the spec preserve `not_assessed`,
   `cannot_verify`, controlled-pilot scope, and external production trust
   boundaries?

## Intake Findings To Challenge

| ID | Severity | Plane | Finding | Initial disposition |
|---|---|---|---|---|
| DOC-001 | major | documentation | `docs/agent-entrypoint.md` documents `pr-review packet` flags that live `--help` does not advertise. | accepted_fixed_in_spec |
| DOC-002 | major | documentation | Russian command reference omits active command families present in English docs and reviewer entrypoint. | accepted_fixed_in_spec |
| DOC-003 | major | documentation | Placeholder example READMEs appear in first-class example surface without evidence-state labels. | accepted_fixed_in_spec |
| DOC-004 | minor | documentation | README status sentence is accurate but dense for first-screen MVP handoff. | accepted_fixed_in_spec |
| CODE-001 | major | code-quality | `golangci-lint run ./...` fails on telemetry `gosimple` and authority `ineffassign`. | accepted_fixed_in_spec |
| CODE-002 | major | code-quality | `ValidateExportResult` has cyclomatic complexity 66 and cannot satisfy `CRAP < 5` without decomposition. | accepted_fixed_in_spec |
| CODE-003 | major | code-quality | `normalizeOpenCodeRawLine` has cyclomatic complexity 37 and mixes multiple trust-sensitive concerns. | accepted_fixed_in_spec |
| CODE-004 | major | code-quality | MVP-critical packages have zero or very low coverage: `contract`, `policy`, `export`, `trace`, `harnessobs`, `verifier`. | accepted_fixed_in_spec |
| CI-001 | major | verification | Current CI does not enforce lint, complexity, or coverage gates. | accepted_fixed_in_spec |

## Disposition Rules

- Critical or major findings must be fixed in the spec or recorded as blockers.
- A finding may be rejected only with file/command evidence.
- If a gate cannot be implemented in the same MVP-hardening block, record the
  state as `not_assessed` with a concrete follow-up; do not call it passing.
- A clean local run is not GitHub CI evidence.

## Initial Review Output

Reviewer: `minimax/MiniMax-M2.7` via local `pi`
Verdict: `REVISE`

### Critical Findings

| ID | Finding | Disposition |
|---|---|---|
| C-001 | "Trust-sensitive" was undefined and circular, making FR-007 unverifiable. | accepted_fixed: `spec.md` now defines trust-sensitive Go path criteria. |
| C-002 | `internal/posture` was excluded from MVP-critical package coverage despite containing the highest-complexity hotspot. | accepted_fixed: `spec.md` FR-008 and `tasks.md` Phase 5 now include `internal/posture`; `plan.md` expected surfaces already include posture. |

### Major Findings

| ID | Finding | Disposition |
|---|---|---|
| M-001 | "MVP bar" was referenced but never defined. | accepted_fixed: `spec.md` now includes MVP Bar Definition. |
| M-002 | "First-class MVP evidence" was undefined. | accepted_fixed: `spec.md` now defines first-class MVP evidence. |
| M-003 | `internal/verifier` was listed as low coverage but had no coverage task. | accepted_fixed: `tasks.md` now adds focused verifier coverage task. |
| M-004 | Phase 0 exit criteria were not independently verifiable. | accepted_fixed: `plan.md` now defines Phase 0 exit criteria. |
| M-005 | Initial thresholds and minimum gate set were undefined. | accepted_fixed: `spec.md` now defines initial thresholds and minimum gate set. |

### Minor Findings

| ID | Finding | Disposition |
|---|---|---|
| m-001 | `not_assessed` and `cannot_verify` distinction unclear in evidence labels. | accepted_fixed: `spec.md` now defines evidence state labels. |
| m-002 | T009 verification method was vague. | accepted_narrower: kept task wording for implementation review; trust-upgrade checks are now anchored by Product Boundary and MVP Bar. |
| m-003 | Ratchet plan had no milestone. | accepted_fixed: SC-005 now requires first ratchet milestone in implementation ledger. |
| m-004 | Schema validation missing from tasks. | accepted_fixed: T032 explicitly runs `jq empty schema/*.json` plus changed JSON examples. |
| m-005 | "Current" command families undefined. | accepted_fixed: `spec.md` defines current command surface and FR-002 references live help at verification time. |

## Remaining State

Spec direction still requires explicit user approval before documentation, CI, or
Go implementation starts.

## Three-Axis Review

Requested after the initial review. All three reviewer runs returned usable
output and `REVISE`.

### Axis 1: Product/Docs Promise

Reviewer: `minimax/MiniMax-M2.7` via local `pi`
Verdict: `REVISE`

| ID | Severity | Finding | Disposition |
|---|---|---|---|
| AX1-C001 | critical | Current evidence used local `main...origin/main [ahead 12]` state without enough stale-baseline handling. | accepted_fixed: `spec.md` Current Evidence now records commit `5f6706b398d6d68bb9a37be2dee4e6aec1037df3`; tasks add branch/commit delta verification. |
| AX1-C002 | critical | MVP readiness claim had no named claimant/sign-off mechanism. | accepted_fixed: MVP Bar Definition and SC-008 now require named reviewer sign-off in block ledger before ready state. |
| AX1-M001 | major | FR-002 allowed Russian docs to route to English, contradicting bilingual parity. | accepted_fixed: FR-002 now permits routing only as temporary `deferred_scope` with follow-up, not parity. |
| AX1-M002 | major | Trust-sensitive definition was too jargon-heavy for onboarding. | accepted_fixed: plain-language summary added before technical criteria. |
| AX1-M003 | major | Ratchet milestone success criterion was missing/enforcement weak. | accepted_fixed: SC-005 now requires actual CRAP baseline and ratchet milestone. |
| AX1-M004 | major | Schema docs were outside spec scope despite schema validation. | accepted_fixed: SC-009 and plan implementation surfaces now include schema documentation gaps. |
| AX1-M005 | major | T006 lacked exact stale flag file reference. | accepted_fixed: T006 now names `docs/agent-entrypoint.md` and requires rationale. |

### Axis 2: Code Quality, CRAP, CI Gate, Coverage Strategy

Reviewer: `zai/glm-5.1` via local `pi`
Verdict: `REVISE`

| ID | Severity | Finding | Disposition |
|---|---|---|---|
| AX2-C001 | critical | Spec referenced CRAP but never computed CRAP; `gocyclo` is not CRAP. | accepted_fixed: Definitions add CRAP formula; MVP Bar and tasks require per-function CRAP baseline/review. |
| AX2-C002 | critical | Review ledger referenced SC-005 before verifying it existed. | accepted_fixed: SC-005 exists and now explicitly covers CRAP baseline and ratchet milestone. |
| AX2-M001 | major | Coverage floor could be satisfied by token happy/error tests at near-zero coverage. | accepted_fixed: MVP Bar now ties coverage to CRAP and changed/MVP-critical functions; tasks require baseline and deltas. |
| AX2-M002 | major | Plan/tasks phase numbering diverged. | accepted_fixed: plan now splits Phase 4, Phase 5, and Phase 6 to match tasks. |
| AX2-M003 | major | Gate definition was sequenced too late, after decomposition. | accepted_fixed: tasks moved gate/baseline definition into Phase 3 before decomposition. |
| AX2-M004 | major | `gocognit` threshold absent. | accepted_fixed: MVP Bar states gocognit needs explicit implementation-ledger threshold or remains `not_assessed`. |
| AX2-M005 | major | Decomposition target complexity unstated. | accepted_fixed: T020/T021 now require selected maximum complexity/CRAP threshold. |
| AX2-M006 | major | "covered by existing tests" was vacuous for low-coverage packages. | accepted_fixed: plan Phase 3 now requires a focused test if changed package coverage is below selected floor. |

### Axis 3: Trust Boundary And Evidence Semantics

Reviewer: `openrouter/qwen/qwen3.6-plus` via local `pi`
Verdict: `REVISE`

| ID | Severity | Finding | Disposition |
|---|---|---|---|
| AX3-M001 | major | CI enforcement gap was mislabeled `not_assessed` despite being in selected scope. | accepted_fixed: new `assessed_gap` state added; FR-009/T032 use it for missing CI enforcement. |
| AX3-M002 | major | Review verdict was stale after spec revision and intake dispositions remained pending. | accepted_fixed: ledger updated with three-axis review, dispositions, and current approval-pending state. |
| AX3-M003 | major | `block ledger` undefined. | accepted_fixed: `spec.md` now defines block ledger. |
| AX3-m001 | minor | Deferred FRs were conflated with `not_assessed`. | accepted_fixed: `deferred_scope` added and MVP Bar updated. |
| AX3-m002 | minor | `real_fixture` label could imply stronger proof. | accepted_fixed: label changed to `real_fixture_local`. |
| AX3-m003 | minor | Current Evidence lacked commit anchor. | accepted_fixed: commit hash added. |

## Focused Post-Revision Re-Review

All three focused re-review planes returned usable `APPROVE` outputs.

| Plane | Reviewer | Verdict | Notes |
|---|---|---|---|
| Product/docs promise | `minimax/MiniMax-M2.7` | `APPROVE` | Confirmed local-state anchoring, readiness sign-off, bilingual deferred-scope handling, schema-doc scope, and actionable T006 are resolved. |
| Code quality/CRAP | `zai/glm-5.1` | `APPROVE` | Confirmed CRAP formula/baseline, SC-005, phase sequencing, gocognit semantics, decomposition targets, and coverage-vacuity fixes are resolved. Minor note about possible SC-005 truncation was checked against the actual file; SC-005 is complete. |
| Trust boundary/evidence semantics | `openrouter/deepseek/deepseek-v4-pro` | `APPROVE` | Confirmed `assessed_gap`, `deferred_scope`, `real_fixture_local`, block ledger definition, review disposition consistency, and commit anchor are resolved. |

## Post-Revision Verdict

Spec review state: `APPROVED_FOR_USER_DECISION`.

Critical and major findings from the initial Socratic review, three-axis review,
and focused post-revision re-review are exhausted. Implementation remains
blocked until explicit user approval of the reviewed spec direction.
