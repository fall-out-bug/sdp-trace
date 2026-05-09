# Implementation Plan: Authority Envelope Boundary Observation

**Branch**: `002-authority-envelope-boundary-observation` | **Date**: 2026-05-09 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/002-authority-envelope-boundary-observation/spec.md`

## Summary

Add a portable authority envelope concept to `sdp-trace` so trace packages can compare declared actor authority against observed actions. The feature must stay product-generic: it records authority facts, source coverage, and attribution evidence, while external policy consumers decide whether an `outside_authority` fact contaminates a demo, blocks a merge, or requires escalation.

## Technical Context

**Language/Version**: Go verifier CLI, JSON Schema Draft 2020-12, Markdown
**Primary Dependencies**: Go standard library; existing schema/fixture validation paths; `jq` syntax checks
**Storage**: Portable JSON artifacts under `schema/` and examples under `examples/`
**Testing**: `go test ./...`; `jq empty schema/*.json`; fixture validation for new examples; `git diff --check`
**Target Platform**: Repository-portable artifacts, no runtime dependency on a specific harness
**Project Type**: Schema, examples, documentation, and later verifier behavior
**Constraints**: No dependency on Codex, OpenCode, GSD, Claude, GitHub, Beads, Operator Mode, or LLM gateway as a required runtime

## Constitution Check

| Rule | Status | Evidence |
|---|---|---|
| Spec before implementation | Pass | This package is spec-only until reviewed and approved. |
| Keep product independent | Pass | Spec forbids demo/harness-specific policy in product behavior. |
| Evidence-backed claims only | Pass | Actor/model/tool attribution is downgraded when source evidence is missing. |
| Preserve `not_assessed` and `cannot_verify` | Pass | Requirements explicitly require both states. |
| No native policy verdicts | Pass | External policy consumers own contamination/block/merge decisions. |

## Project Structure

```text
specs/002-authority-envelope-boundary-observation/
├── spec.md
├── plan.md
├── data-model.md
├── tasks.md
└── socratic-review.md
```

Expected implementation artifacts after approval:

```text
schema/
├── authority-envelope.schema.json
├── observed-action.schema.json
├── evidence-binding.schema.json
└── authority-evaluation.schema.json

examples/
└── authority-envelope-basic/
    ├── fixture-matrix.json
    ├── valid-within-authority/
    ├── outside-authority-denied-target/
    ├── git-only-attribution-not-assessed/
    ├── gateway-unbound-model-not-assessed/
    └── malformed-policy-cannot-verify/

docs/
└── authority-envelope.md
```

## Contract Decisions To Review

| Decision | Rationale | Review risk |
|---|---|---|
| Use policy-selected evaluation, not global defaults | Avoid encoding one workflow as product behavior. | Users may expect default safety policy; docs must be explicit. |
| Caller supplies `policy_id` | `sdp-trace` records facts for the selected envelope and does not choose among policies. | Missing selected policy must stay `not_assessed`. |
| Reject ambiguous envelopes | Avoid hidden allow-wins or deny-wins behavior. | Policy authors must resolve conflicts before evaluation. |
| Unmatched actions are unassessed, not denied by default | Avoid encoding an implicit deny-all policy. | External policy consumers may still choose to treat `not_assessed` as blocking. |
| Keep git-only mutation separate from actor/model attribution | Git can prove changed files, not tool/model causality. | Product may feel less decisive but stays honest. |
| Treat gateway evidence as insufficient without harness/action binding | Model request alone does not prove file mutation. | Requires an import/binding contract later. |
| Model `outside_authority` as a fact, not a gate verdict | External consumers decide consequences. | Naming may be misread as policy decision; docs must clarify. |
| Start with schema/examples before verifier behavior | Avoid premature policy logic. | Needs strong examples to prevent ambiguous implementation. |
| Evidence refs use safe URI-style strings | Prevent unresolvable free-form references. | Schema must define allowed schemes before examples count as proof. |

## Phases

### Phase 0: Socratic Spec Review

Run independent critic review of this spec package before implementation. Valid critical/major findings must be fixed or recorded as blockers before asking for implementation approval.

### Phase 1: Contract Design

Add JSON schemas for authority envelope, observed action, evidence binding, and authority evaluation. Add schema README entries and examples.

### Phase 2: Fixture Matrix

Add positive and negative fixtures proving git-only, harness-bound, gateway-unbound, malformed-policy, failed-binding, stale-evidence, feedback-without-mutation, and approval-required cases.

### Phase 3: Verifier Behavior

Add Go verifier behavior only after schemas and fixtures are approved. The verifier emits facts, not merge/block decisions.

### Phase 4: Documentation

Document observation sources, source coverage, attribution limits, and the external policy boundary.

## Acceptance Gate Before Implementation

Implementation must not start until:

- `spec.md`, `plan.md`, `data-model.md`, and `tasks.md` exist;
- Socratic pi-review has usable output from independent reviewers;
- critical/major spec findings are resolved or explicitly blocked;
- the user approves the reviewed spec direction.
