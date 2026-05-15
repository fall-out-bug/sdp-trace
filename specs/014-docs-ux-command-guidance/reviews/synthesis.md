# Socratic Review Synthesis: 014-docs-ux-command-guidance

## Reviewers
- DeepSeek (opencode/deepseek-v4-flash-free) — reasoning review
- MiniMax (minimax/MiniMax-M2.7) — trust/adversarial review
- Qwen (opencode/qwen3.6-plus-free) — wide-context review

## Cross-Model Consensus Findings

### 🔴 Critical — State vocabulary fragmentation (all three reviewers)
**Finding:** The spec's US-002 lists `missing_telemetry` as a state reviewers must distinguish, but the canonical "State And Exit Code Contract" in `docs/agent-entrypoint.md:273-284` and `docs/reviewer-entrypoint.md:38-41` defines only 5 states: `observed`, `pass`, `fail`, `not_assessed`, `cannot_verify`.

Additional orphan tokens found across docs:
- `missing_telemetry` — 70+ matches across 5 docs
- `not_integrated` — matches in `agent-onboarding.md`, `harness-integration.md`, `flight-recorder.md`
- `warn` — defined in `docs/concepts.md:88` (External Verdict), absent from entrypoints
- `coverage_satisfied/partial/unresolved` — used in `reviewer-entrypoint.md:122-123` for pr-review, not in state contract
- `unsupported`, `retention_limited`, `stale`, `outside_authority`, `complete/partial` — scattered

**Impact:** US-002 acceptance criterion ("State vocabulary is consistent across README, concepts, agent entrypoint, reviewer entrypoint, and adoption guide") is unachievable without first defining the canonical vocabulary boundary.

**Disposition:** Accepted. The spec delta must either (a) expand the canonical state contract to classify all tokens (result states vs telemetry labels vs completeness markers vs command-specific sub-states), or (b) scrub non-canonical tokens from non-spec docs and replace with canonical equivalents.

### 🔴 Critical — Profile selection has disjoint taxonomies (DeepSeek + Qwen)
**Finding:** Three profile-like vocabularies coexist with zero cross-reference:
1. Trust profile IDs: `repo_baseline_structural`, `source_bound_local_release`, `external_production_trust`
2. Assessment profiles: `adapter-capture`, `managed-harness`, etc.
3. Witness kinds: `github-actions`, `gitlab-ci`, `buildkite`, `customer-pki`

Plus `local_dirty_structural_only` (authority scope) and `repo_baseline` (naming mismatch in `claim-authoring.md:15`).

**Impact:** FR-003 ("which profile do I use?" decision tree) is essential but the spec does not explicitly task the mapping between these taxonomies.

**Disposition:** Accepted. Add explicit task to create a single decision aid mapping trust profile IDs ↔ assessment profiles ↔ witness kinds.

### 🔴 Critical — Overclaim rules duplicated, not canonical (MiniMax + Qwen)
**Finding:** Forbidden/overclaim claims appear in at least 4 documents with overlapping but non-identical wording:
- `agent-entrypoint.md:297-305` — "Forbidden Claims" (5 bullets)
- `reviewer-entrypoint.md:119-146` — "Gate, Witness, And Release Caveats" + "What You May State From Output"
- `agent-onboarding.md:75-76` — inline caveats
- `adoption-guide.en.md:11-13` — opening caveats

**Impact:** US-004 acceptance criterion ("Reviewer entrypoint contains the canonical checklist; README and agent entrypoint link to it") fails because no single canonical checklist exists.

**Disposition:** Accepted. FR-005 must produce one canonical overclaim checklist file; all other docs link to it. Critical warnings may retain one-line inline summaries with links.

### 🟡 Major — No output location map (DeepSeek + Qwen)
**Finding:** Output destinations are scattered across README, adoption guide, agent entrypoint command table. No single table maps: command → default output path → format → purpose → trust boundary.

**Disposition:** Accepted. T012 must produce an explicit output location reference (table or diagram).

### 🟡 Major — Reviewer entrypoint is not task-first (DeepSeek + Qwen)
**Finding:** `reviewer-entrypoint.md` is ~190 lines. The flat command list (lines 46-65) precedes the Quick Reference task table (lines 148-161). A cold user sees raw commands before task guidance.

**Disposition:** Accepted. T010 must restructure reviewer entrypoint so task path comes first, or spin out the flat reference to a separate file.

### 🟡 Major — Acceptance criteria under-specified (Qwen)
**Finding:** "Short task path" and "consistent vocabulary" lack measurable thresholds. Doccheck scope for docs-only changes is unclear.

**Disposition:** Accepted. Add measurable AC: e.g., grep audit shows zero orphan state tokens outside canonical doc.

## Disposition Summary

| Finding | Severity | Disposition | Owner |
|---|---|---|---|
| State vocabulary fragmentation | Critical | Accepted, requires spec delta | Spec author |
| Profile taxonomy disjoint | Critical | Accepted, requires new task | Spec author |
| Overclaim rules duplicated | Critical | Accepted, requires canonical file | Spec author |
| No output location map | Major | Accepted | Implementation |
| Reviewer entrypoint not task-first | Major | Accepted | Implementation |
| Acceptance criteria under-specified | Major | Accepted | Spec author |
| Link-based consolidation risk | Medium (plan) | Accepted, add mitigation | Plan |
| State-term drift verification missing | Medium (tasks) | Accepted, add grep task | Tasks |

## Required Spec Delta Before Implementation

1. **Canonical State Contract:** Define the canonical vocabulary boundary. Which tokens are result states with exit codes? Which are telemetry labels? Which are command-specific sub-states? Add this to the spec as a new requirement or clarify FR-002.
2. **Profile Taxonomy Map:** Add explicit requirement/task for a single decision aid mapping all profile/scope/witness taxonomies.
3. **Measurable Acceptance Criteria:** Add grep-based or doccheck-based measurable thresholds.
4. **Plan Risk Mitigation:** Add mitigation for link-based consolidation (inline one-line summaries + link).
5. **Tasks Update:** Add T023 (grep audit for orphan state tokens).

## Verdict
**Spec direction is sound** — it correctly identifies the UX problems. **However, implementation must not start until the canonical vocabulary boundary is defined**, or the acceptance criteria for US-002 will remain unverifiable.

Reviewer state: `review_complete_pending_spec_delta`
