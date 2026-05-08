## Review: Proposed Rule/Skill Updates

### Finding 1 — AGENTS.md additions violate own decomposition rule (MAJOR)

The retro proposes adding workflow reminders to AGENTS.md (parent integration audit, empty/hung pi output, PR merge/cleanup from main worktree). These are **skill-level operational workflow** details, not agent boundary/purpose/constraint rules. AGENTS.md already delegates block work to the skill (`"use sdp-trace-trust-workflow"`). Duplicating skill workflow there violates AGENTS.md's own line 84:

> *If `AGENTS.md` exceeds 100 lines or any module needs more than 10 skills, the module is too large, under-decomposed, or overengineered.*

Current AGENTS.md is ~95 lines. Even "compact" additions land at the edge. The correct home is SKILL.md exclusively.

**Fix**: Drop the AGENTS.md proposal entirely. Put all new workflow clauses in the skill only. AGENTS.md already has the trigger and the delegation — that is sufficient.

### Finding 2 — Proposed skill additions are sound and fill real gaps

The five proposed SKILL.md additions (integration audit after subagents, negative leak tests, `--no-tools --no-context-files` default, empty artifact hygiene, merge-from-main-worktree) all address concrete failures from the retro and **do not conflict** with the Go-only stack or existing rules. The 12-step expected-behavior sequence is a natural expansion of the existing block intake protocol.

No issues with these.

### Finding 3 — One missing actionable: negative-leak-test scope trigger

The retro says "for safety-sensitive outputs, add negative leak tests using secret-like markers" but does not define **who decides** what counts as safety-sensitive at implementation time. The existing SKILL.md has no classification step for claim surface sensitivity.

**Fix**: Add a one-liner trigger in the skill, e.g.: *"If a chunk changes CLI output rendering, command preview, dry-run, or any path where user-supplied secrets appear in argv/env, classify it safety-sensitive and require at least one negative assertion (marker must not appear in output)."*

---

## VERDICT: REVISE

- **Drop AGENTS.md additions** — all new content goes to SKILL.md only.
- **Add scope trigger** for safety-sensitive negative tests.
- Everything else in the retro proposal is clean and can land as-is.
