# Research: Post-Merge Governance Closure

## Decision: Bind planning directly to Spec 022 despite branch mismatch

Rationale: The user explicitly requested `speckit-plan spec 22`. The Spec Kit
setup script requires a numbered feature branch and failed on
`codex/install-github-speckit`, but the target spec path is unambiguous.

Alternatives considered:

- Switch branches before planning: rejected because it risks moving away from
  the current committed Spec Kit governance changes.
- Abort until a numbered branch exists: rejected because it would add process
  without improving the plan content.

## Decision: Docs-only governance closure

Rationale: Spec 022 is about preserving and closing governance evidence after a
missed pre-merge gate. It does not change runtime behavior, command contracts,
schemas, or product code.

Alternatives considered:

- Add Go validation for closure ledgers now: rejected as over-scope for a
  closure planning spec.
- Create a new JSON decision artifact: rejected because existing Markdown
  ledgers are the current governance surfaces.

## Decision: Live PR/CI refresh before any `complete` claim

Rationale: Checked-in ledgers are context, not authority. Before claiming Spec
022 complete, the worker must refresh PR #60 / PR #63 live state when available
and record unavailable live state as `not_assessed` or `cannot_verify`.

Alternatives considered:

- Trust checked-in ledgers only: rejected because it can repeat stale closure
  drift.
- Require live GitHub state as a hard blocker in all environments: rejected
  because local/offline work should be able to record `not_assessed` or
  `cannot_verify` explicitly.

## Decision: Synchronize three governance surfaces

Rationale: Spec 022 closure affects the decision surface, reality ledger, and
navigation surface. Updating only one would reintroduce roadmap/ledger drift.

Alternatives considered:

- Update only `docs/closure-decision-ledger.md`: rejected because readers use
  roadmap and spec reality ledger for current state.
- Add a separate closure artifact and link to it: rejected as unnecessary until
  closure evidence is too large for the three existing surfaces.

## Decision: No contracts directory

Rationale: This plan does not introduce or change an external interface. The
contract is repository documentation consistency, captured in existing docs and
tasks.

Alternatives considered:

- Add a Markdown contract under `contracts/`: rejected because it would create
  an extra authority surface for no API/schema change.
