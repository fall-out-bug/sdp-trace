---
name: sdp-trace-trust-workflow
description: Run sdp-trace block intake and trust-sensitive implementation workflow from SpecKit delta through review, scoped commits, PR evidence, and final handoff.
---

<objective>
Keep block work source-bound, reviewable, and honest: no implementation-only closure, no deferred trust closure, and no checked-in prose replacing live evidence.
</objective>

<when_to_use>
Use this skill when the user says "берем блок в работу", asks to take a block into work, starts a trust-sensitive feature block, or requests PR-ready closure for a block.
</when_to_use>

<when_not_to_use>
Do not use this skill for quick read-only explanations, purely mechanical formatting, or unrelated external workflow management. If the work touches Beads, deployment, UI, or harness-specific automation, keep that outside the active product path unless the user explicitly asks.
</when_not_to_use>

<intake_protocol>
1. First land or explicitly park current approved work through review/PR discipline.
2. Continue new block work in a fresh worktree when practical.
3. Prepare a SpecKit delta and implementation plan.
4. Run Socratic spec review with `pi-review`; fix or record findings.
5. Stop for explicit user approval of the reviewed spec direction before writing implementation code.
</intake_protocol>

<implementation_loop>
After approval:
1. Split independent tasks into bounded slices with clear file ownership.
2. Use subagents for parallel exploration, implementation, testing, and review when they materially help.
3. Keep each slice traceable to spec, task, evidence, gate, decision, and provenance changes.
4. Use test-first behavior when behavior changes.
5. Run drift checks: spec vs implementation and regression against previous blocks.
6. Record verifier state as `pass`, `fail`, `cannot_verify`, or `not_assessed`.
7. Commit each verified slice with a scoped message.
</implementation_loop>

<pr_protocol>
1. Prepare a PR with code, tracing/evidence, docs, and requirements mapping.
2. Run separate review planes at PR level: code/correctness, tracing/evidence/provenance, requirements-vs-implementation, and security/DX/UX when relevant.
3. Verify reviewer findings against full files before accepting or rejecting them.
4. Re-read the actual diff before finalizing the PR description. Remove any claimed change that was reverted or never made (e.g., baseline updates that were later removed).
5. Ensure the PR is not in Draft state before claiming it is ready for review.
6. Query live GitHub checks for the final head. If checks are absent, record CI as `not_assessed`.
7. Do not claim ready/merge/sign-off until required review and live CI evidence are present.
8. After merge: delete the local worktree (`git worktree remove`) and the feature branch (`git branch -d` or remote deletion). The branch should not persist on origin once merged; keeping it creates drift and implies unfinished work.
</pr_protocol>

<trust_boundaries>
- Dirty checkout output is local structural evidence only.
- Source-bound proof requires a clean immutable source commit.
- If a changed file is a manifest subject, commit it first, then regenerate release proof in a separate commit.
- Do not close task checkboxes, ledgers, or docs after source-bound proof without another source-bound cycle if those files are manifest subjects.
</trust_boundaries>

<anti_rationalization>
Before claiming closure, read `references/anti-rationalizations.md` and reject any shortcut that converts prose, stale JSON, dirty checkout state, or unchecked reviewer output into trust authority.
</anti_rationalization>

<supporting_files>
- Use `templates/block-intake.md` when opening or resuming a block.
- Use `templates/final-evidence-map.md` before PR-ready or complete claims.
- Use `references/anti-rationalizations.md` for stop conditions and common false closures.
</supporting_files>

<success_criteria>
- SpecKit delta and implementation plan were reviewed before implementation approval.
- Each slice maps to spec, task, evidence, gate, decision, and provenance impact.
- Verification state is explicitly `pass`, `fail`, `cannot_verify`, or `not_assessed`.
- Final claims cite fresh command output, live external checks, or explicitly remain open.
</success_criteria>
