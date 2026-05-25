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
4. Run Socratic spec review with adversarial review; fix or record findings.
5. Commit or explicitly park the reviewed spec/review handoff before delegating or writing implementation code.
6. Stop for explicit user approval of the reviewed spec direction unless the user already requested unattended worker delegation through OmPi `task` tool subagents.
</intake_protocol>

<implementation_loop>
After approval:
1. Split independent tasks into bounded slices with clear file ownership.
2. For large approved blocks, use OmPi `task` tool subagents in isolated branches instead of loading the whole implementation loop into the active context.
3. Use subagents for parallel exploration, implementation, testing, and review when they materially help.
4. Keep each slice traceable to spec, task, evidence, gate, decision, and provenance changes.
5. Use test-first behavior when behavior changes.
6. Run drift checks: spec vs implementation and regression against previous blocks.
7. Record verifier state as `pass`, `fail`, `cannot_verify`, or `not_assessed`.
8. Commit each verified slice with a scoped message.
9. Ensure verification commands in docs are copy-pasteable and reproducible. Use subshell isolation for commands that must change directory or depend on ambient state. Document both configured and default-config scanner behavior when they differ.
</implementation_loop>

<pr_protocol>
1. Prepare a PR with code, tracing/evidence, docs, and requirements mapping.
2. Run separate review planes at PR level through the OmPi `reviewer` agent via the `task` tool: code/correctness, tracing/evidence/provenance, requirements-vs-implementation, and security/DX/UX when relevant.
3. Additional review planes may use external model APIs when explicitly configured in `.omp/model-policy.yml`, recorded in the review artifact with model, provider, harness, date, and prompt class.
4. For trust-sensitive PRs, run iterative adversarial review rounds against the full diff. Fix every finding of any severity before the next round. Repeat until the reviewer outputs exactly `LGTM` (zero findings).
4. Verify reviewer findings against full files before accepting or rejecting them.
5. Re-read the actual diff before each review round and before finalizing the PR description. Remove any claimed change that was reverted or never made (e.g., baseline updates that were later removed).
6. Delete stale review artifacts; do not rely on headers or markers to neutralize stale claims.
7. Ensure the PR is not in Draft state before claiming it is ready for review.
8. Query live GitHub checks for the final head. If checks are absent, record CI as `not_assessed`.
9. Do not claim ready/merge/sign-off until required review and live CI evidence are present.
10. After merge: delete the local worktree (`git worktree remove`) and the feature branch (`git branch -d` or remote deletion). The branch should not persist on origin once merged; keeping it creates drift and implies unfinished work.
</pr_protocol>

<trust_boundaries>
- Dirty checkout output is local structural evidence only.
- Source-bound proof requires a clean immutable source commit.
- If a changed file is a manifest subject, commit it first, then regenerate release proof in a separate commit.
- Do not close task checkboxes, ledgers, or docs after source-bound proof without another source-bound cycle if those files are manifest subjects.
- Security docs must not list unverified contact information. Email addresses and reporting channels must be confirmed deliverable or explicitly marked `not_assessed`.
- Scanner counts in docs must match the exact output of the documented command. Update docs or commands until they match.
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
