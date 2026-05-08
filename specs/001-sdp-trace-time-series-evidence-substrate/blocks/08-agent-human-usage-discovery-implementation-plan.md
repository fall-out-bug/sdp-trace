# Block 08 Implementation Plan: Agent and Human Usage Discovery

> **For agentic workers:** REQUIRED EXECUTION MODE: implement each task through fresh `gpt-5.3-codex-spark` subagents with review between tasks. Use `superpowers:subagent-driven-development`. Steps remain repo-tracked in SpecKit artifacts and the Block 08 review ledger.

**Goal:** implement repository-visible first-use entrypoints for agents and human reviewers so they can discover the correct verifier profile, command surface, current trust boundary, and forbidden claims without adding runtime behavior.

**Architecture:** docs-first implementation with one agent entrypoint doc, one human reviewer doc, minimal README surfacing, and a deterministic validation layer that keeps onboarding text aligned with the live verifier surface. Trust semantics remain owned by `scripts/verify.sh` and Block 07; Block 08 only exposes them clearly and prevents overclaim drift.

**Tech Stack:** Markdown docs, existing shell verifier surface, `rg`, small shell validation script(s), existing `scripts/validate-contracts.sh`, SpecKit artifacts, existing release-proof workflow.

---

Status: implemented; implementation-review closure and source-bound proof sync completed, external trust remains intentionally open
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`

## Goal

Implement Block 08 as repository-readable discovery surfaces and validation guardrails without introducing new verifier features, new CLI UX, new agent workflow, or any broader trust claim than Block 07 currently proves.

## Consensus

Consensus is recorded in:

- `blocks/08-agent-human-usage-discovery.md`
- `blocks/08-agent-human-usage-discovery-review-ledger.md`

Implementation scope is intentionally narrow:

- add user-facing discovery docs for agent and human first-use paths
- add deterministic validation that keeps those docs aligned with the existing verifier surface
- add only minimal navigation so the new docs are discoverable
- keep external production trust explicitly blocked

Out of scope:

- new verifier profiles
- new `scripts/verify.sh` flags or command aliases
- any runtime dependency on Codex, Pi, Beads, OpenCode, or another harness
- broad README rewrite or any wording that suggests trust is easier than Block 07 proves

## File Responsibilities

- `docs/agent-entrypoint.md`: portable agent first-use contract for profile selection, command contract, JSON/text output expectations, and forbidden claims.
- `docs/reviewer-entrypoint.md`: human first-use path for five-minute verification, clean-vs-dirty interpretation, external trust blocker, and `not_assessed` handling.
- `README.md`: minimal neutral links to the new discovery docs from the existing start surface.
- `scripts/validate-discovery-entrypoints.sh`: deterministic validation for required sections, exact command surface references, and forbidden overclaim wording in discovery docs.
- `scripts/test-discovery-entrypoints.sh`: focused positive/negative tests for the Block 08 discovery validator.
- `scripts/validate-contracts.sh`: call the Block 08 discovery validator as part of the repo validation surface.
- `specs/001-sdp-trace-time-series-evidence-substrate/spec.md`: update only if the root FR/SC set needs explicit Block 08 implementation outputs.
- `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md`: close T087-T089 after docs and validator coverage are real; close T090 after the fresh-agent pre-implementation review; keep implementation-review closure in the Block 08 ledger separately.
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/08-agent-human-usage-discovery.md`: update status only if implementation changes the spec-facing lifecycle marker; do not change trust semantics.
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/08-agent-human-usage-discovery-review-ledger.md`: append implementation-review findings and closure evidence.
- Release proof artifacts under `examples/contract-foundation/` and `examples/self-trace/`: update only if touched files are manifest subjects and the source-bound cycle requires regenerated proof.

## Task 1: Scope Freeze and Implementation Gate

- Confirm implementation starts from the accepted Block 08 design artifacts and the committed review ledger.
- Keep implementation limited to docs, deterministic validation, and minimal navigation.
- Preserve the existing command surface exactly:
  - `npm run verify:baseline`
  - `npm run verify:source-bound`
  - `npm run verify:external-trust`
  - `scripts/verify.sh --profile baseline|source-bound|external-trust [--json] [--allow-dirty] [--version]`
- Preserve current trust facts:
  - clean-checkout activation basis: baseline `pass`, source-bound `pass`, external-trust `fail`
  - dirty-checkout baseline without `--allow-dirty`: `cannot_verify`
  - dirty-checkout baseline with `--allow-dirty`: `pass` only in `local_dirty_structural_only`
  - dirty-checkout source-bound and external-trust: `cannot_verify`

Verification:

```bash
rg -n "Goal|Consensus|File Responsibilities|Task 1" specs/001-sdp-trace-time-series-evidence-substrate/blocks/08-agent-human-usage-discovery-implementation-plan.md
npm run verify:baseline -- --allow-dirty
scripts/verify.sh --profile source-bound --allow-dirty
scripts/verify.sh --profile external-trust --allow-dirty
```

Expected result: the plan stays within Block 08 scope and the implementation gate uses live verifier behavior rather than prose assumptions.

## Task 2: Agent Entrypoint Doc

- Create `docs/agent-entrypoint.md`.
- Make it a portable agent-facing discovery contract, not a harness-specific recipe.
- Include:
  - which profile to choose from which claim
  - the fixed command surface
  - when to use `--json`
  - what `pass`, `fail`, `not_assessed`, and `cannot_verify` mean for the selected profile
  - explicit forbidden claims
  - explicit statement that checked-in `proof-summary` JSON is not authority until live-verified
- Reuse Block 08 vocabulary exactly: `repo_baseline_structural`, `source_bound_local_release`, `external_production_trust`, `not_assessed`, `local_dirty_structural_only`.
- Do not mention Codex-specific behavior as part of the contract.

Verification:

```bash
rg -n "repo_baseline_structural|source_bound_local_release|external_production_trust|--json|not_assessed|cannot_verify|forbidden" docs/agent-entrypoint.md
rg -n "Codex|OpenCode|Beads|Claude|Pi" docs/agent-entrypoint.md
```

Expected result: the doc contains the required profile/command/claim guidance and does not turn portability into hidden harness assumptions.

## Task 3: Human Reviewer Entrypoint Doc

- Create `docs/reviewer-entrypoint.md`.
- Make it a five-minute first-use path for technical executive/CISO/reviewer audiences.
- Include:
  - clean-checkout first-pass commands
  - explicit warning that `verify:external-trust` currently exits `1` from a clean checkout
  - explicit warning not to use `--allow-dirty` for source-bound or external-trust conclusions
  - clean-vs-dirty interpretation
  - what `not_assessed` means and what action it allows
  - the current external trust blocker and why its states are downstream of `external_trust_profile_selected: fail`
  - a short “what you may say from this result” boundary
- Keep the doc human-readable, but bind every allowed statement back to verifier states rather than soft trust language.

Verification:

```bash
rg -n "verify:baseline|verify:source-bound|verify:external-trust|clean checkout|--allow-dirty|not_assessed|external_trust_profile_selected|production_release_verified" docs/reviewer-entrypoint.md
rg -n "production-ready|supported|compatible|ready for production|trusted release" docs/reviewer-entrypoint.md
```

Expected result: a reviewer can follow the path without confusing a known external-trust blocker, dirty-checkout limitation, or `not_assessed` state for hidden success.

## Task 4: Minimal Navigation

- Add only minimal neutral surfacing in `README.md`.
- Place links in the existing “Start Here” area or equivalent navigational section.
- Use labels that describe scope honestly, such as:
  - `Agent entrypoint (current verifier contract)`
  - `Reviewer entrypoint (current proof scope)`
- Do not rewrite the README narrative or suggest broader trust/readiness than current verifier output supports.

Verification:

```bash
rg -n "Agent entrypoint|Reviewer entrypoint|current verifier contract|current proof scope" README.md
rg -n "production-ready|supported|compatible|ready for customer" README.md
```

Expected result: the new docs are discoverable from the repo front door without turning the README into a marketing or closure surface.

## Task 5: Deterministic Discovery-Doc Validation

- Create `scripts/validate-discovery-entrypoints.sh`.
- Validate:
  - required docs exist
  - required section markers exist in both new docs
  - all three canonical commands appear exactly
  - required profile names appear exactly
  - required dirty-checkout and `not_assessed` guidance is present
  - forbidden overclaim phrases do not appear in the Block 08 discovery docs and the new README links
- Keep the validator deterministic and grep-based. Do not add NLP or model judgment.
- Wire it into `scripts/validate-contracts.sh`.

Verification:

```bash
scripts/validate-discovery-entrypoints.sh
rg -n "validate-discovery-entrypoints" scripts/validate-contracts.sh
```

Expected result: the validator passes on the good docs and becomes part of the normal repo validation path.

## Task 6: Validator Tests

- Create `scripts/test-discovery-entrypoints.sh`.
- Cover:
  - happy path on the committed docs
  - failure when one required command line is removed
  - failure when a forbidden trust/readiness phrase is injected
  - failure when a required profile name is replaced with a non-canonical alias
- Use temporary files/directories; do not mutate committed docs during the test.
- Keep failures named and reproducible.

Verification:

```bash
scripts/test-discovery-entrypoints.sh
```

Expected result: the validator has at least one positive case and explicit negative cases, so Block 08 guardrails are not just best-effort prose.

## Task 7: SpecKit Sync and Source-Bound Proof Discipline

- Update `tasks.md` only after the docs and validator are implemented and validated:
  - close T087 only after the design review is both implemented and repository-visible
  - close T088 only after the agent entrypoint doc and validator are real
  - close T089 only after the human entrypoint doc and validator are real
  - close T090 after the fresh-agent pre-implementation review and before implementation planning
- Record implementation-review findings (including remediation evidence and rerun evidence) in the Block 08 review ledger as a separate implementation closure gate.
- Update `spec.md` only if the root spec needs explicit Block 08 implementation outputs or success criteria.
- If any touched file is a manifest subject:
  - commit the source-subject changes first
  - regenerate manifest / local DSSE / self-attestation / release verification in a separate source-bound cycle
  - rerun `npm run verify:source-bound`
- Keep `external_production_trust` open unless real external evidence exists.

Verification:

```bash
rg -n "T087|T088|T089|T090|Block 08" specs/001-sdp-trace-time-series-evidence-substrate/tasks.md specs/001-sdp-trace-time-series-evidence-substrate/spec.md
npm run verify:baseline
npm run verify:source-bound
```

Expected result: SpecKit status matches implemented discovery surfaces, and source-bound proof is either regenerated correctly or explicitly left pending until the required commit split is done.

## Task 8: Implementation Review and Closure

- Run implementation pi review against:
  - `docs/agent-entrypoint.md`
  - `docs/reviewer-entrypoint.md`
  - `README.md`
  - `scripts/validate-discovery-entrypoints.sh`
  - `scripts/test-discovery-entrypoints.sh`
  - `scripts/validate-contracts.sh`
- Use:
  - GLM for multi-file spec/DX consistency
  - MiniMax for trust-boundary and overclaim review
  - Kimi only for one-file human-doc micro-reviews
- Record every valid finding in `blocks/08-agent-human-usage-discovery-review-ledger.md`.
- Fix every valid finding, including minor ones.
- Rerun validation and proof commands from a clean checkout before claiming closure.

Verification:

```bash
scripts/test-discovery-entrypoints.sh
scripts/validate-discovery-entrypoints.sh
npm run validate
git diff --check
npm run verify:baseline
```

Expected result: Block 08 implementation closes with repository-visible docs, deterministic validation, implementation-review evidence in the ledger, and no new trust overclaims.
