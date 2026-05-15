# Review Disposition: 015-spec-governance-and-roadmap

Review objective: Adversarial Socratic spec review before implementation approval.
- Artifact: `specs/015-spec-governance-and-roadmap/spec.md`
- Contract: AGENTS.md trust rules; docs/claim-authoring.md; project convention that checked-in prose is not authority.
- Review plane: spec/tracing/evidence/provenance
- Model/harness: Internal adversarial review (pi harness) — external GLM/Qwen/DeepSeek planes deferred to PR stage.
- Prompt class: claim-doubt-cycle
- Timeout/retries/fallback: N/A for internal review; external planes to be launched at PR.

## Findings

| Severity | Finding | Evidence | Disposition | Verification |
| --- | --- | --- | --- | --- |
| Important | Overclaim risk: spec claims "repository has a lightweight roadmap" but defines no machine-verifiable evidence for this claim. | spec.md Core Claim | accepted | Add explicit acceptance criteria and verification command (e.g., doccheck + file existence). |
| Important | Lifecycle labels vs trust verdicts: US-002 defines `complete` but does not explain how to prevent readers from conflating status with trust verdict. | spec.md US-002 | accepted | Add caveat paragraph or link to claim-authoring.md distinguishing status from evidence state. |
| Important | Historical evidence boundary: US-003 requires distinction but proposes no concrete mechanism (directory naming, claim tags, exclusion list). | spec.md US-003 | accepted | Add concrete rule: historical block records remain in `blocks/` but roadmap marks them `historical` with no live status. |
| Advisory | Claim-tag enforcement scope vague: US-004/FR-004 require a plan but do not define Markdown scopes or exemption rules. | spec.md US-004, FR-004 | accepted | Scope to new/touched files unless separately approved; list exempt historical paths. |
| Advisory | Missing capability mapping: US-001 requires one-page roadmap but spec does not list capabilities or discovery method. | spec.md US-001 | accepted | Derive capabilities from existing spec titles and command/docs surface; document in roadmap. |
| Advisory | No roadmap freshness mechanism: Risk section notes stale roadmap but defines no owner or refresh cycle. | spec.md Risks | accepted | Add lightweight rule: roadmap updated when new spec opened or status changed; owner = spec author. |
| Advisory | Status transitions undefined: Seven statuses listed but no transition rules (e.g., blocked -> complete?). | spec.md US-002 | accepted | Add minimal transition guidance or keep it informal for Slice 1. |
| Advisory | Roadmap artifact format unspecified: FR-001 says "short roadmap/navigation artifact" but no filename/location/format. | spec.md FR-001 | accepted | Propose `docs/roadmap.md` or `ROADMAP.md`; await user approval. |

## Unresolved states
- not_assessed: External multi-LLM review (GLM/Qwen/DeepSeek) — deferred to PR stage due to harness availability.
- cannot_verify: (none)

## Synthesis
- Required fixes before implementation:
  1. Add machine-verifiable acceptance criteria to Core Claim.
  2. Add lifecycle label caveat linking to claim-authoring.md.
  3. Define concrete historical evidence boundary rule.
  4. Scope claim-tag enforcement to new/touched files with exemption list.
  5. Propose roadmap filename/location.
- Advisory follow-ups:
  - Define status transition rules in a later governance slice if needed.
  - Add roadmap freshness owner rule.
- What this review does not prove:
  - That the roadmap will be accurate (depends on manual curation).
  - That lifecycle labels will be adopted by all specs (depends on contributor discipline).
  - That external LLM reviewers will find no additional issues.
