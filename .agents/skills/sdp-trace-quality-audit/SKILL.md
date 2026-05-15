---
name: sdp-trace-quality-audit
description: Audit and polish sdp-trace quality gates, docs, DX/UX, security, spec drift, and completion evidence when repository-wide quality or readiness is requested.
---

<objective>
Turn broad repository-polish requests into evidence-backed closure without overclaiming trust, CI, review, or metric status.
</objective>

<when_to_use>
Use this skill for requests mentioning repository polish, CRAP, cognitive complexity, Maintainability Index, Clean Code, Clean Architecture, security review, DX/UX review, docs completeness, spec drift, or "work without spec".
</when_to_use>

<when_not_to_use>
Do not use this skill to bypass SpecKit deltas for feature work, to claim live CI from checked-in artifacts, or to perform product deployment/release authority. Use `sdp-trace-trust-workflow` for block implementation and `pi-review` for adversarial review.
</when_not_to_use>

<principles>
- Treat the user's checklist as deliverables, not vibes.
- Machine proof beats prose, but only for the requirement it actually covers.
- Checked-in CI/proof prose is not live authority; final-head CI must be queried after every push.
- For user-facing MI `> 70`, replay with `-mi-under 70.1` and `-function-mi-under 70.1` so rounded `70.0` rows do not pass by accident.
- CI baseline ratchets prevent regressions; they are not the same as an absolute MI closure claim.
- File splitting can satisfy metrics while hurting navigation; record that as a Clean Code/DX review point instead of hiding it.
- Preserve advisory findings separately from blockers.
</principles>

<process>
1. Start with `git status --short`, current branch/head, and PR/check state if a PR exists.
2. Restate the objective as concrete success criteria.
3. Build a prompt-to-artifact checklist mapping every requested gate/review/doc item to commands, files, reviewers, or live PR evidence.
4. Run live gates:
   - `go test -count=1 ./...`
   - `go vet ./...`
   - `go run ./tools/doccheck`
   - `git diff --check`
   - `rg -n "TODO|FIXME" cmd internal tools || true`
   - `go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools`
   - `go run ./tools/qualitycheck -fail-only -mi-under 70.1 cmd internal tools`
   - `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 cmd internal tools`
   - coverage-backed CRAP: `go test -count=1 ./... -coverprofile=coverage.out`, `go tool cover -func=coverage.out`, `go run ./tools/qualitycheck -gocyclo cmd internal tools`, then `go run ./tools/crapcheck -threshold 5 -strict-less`.
5. Run or collect separate security/trust, spec/no-spec, and Clean Code/DX/UX review planes. Treat subagent output as advisory until checked against files and commands.
6. Update durable docs only for stable policy, open gaps, and audit mappings. Do not check in exact "CI passed on this head" proof that a later doc commit will invalidate.
7. Commit in scoped slices and query final-head PR CI after each push.
</process>

<completion_audit>
Before declaring done, verify each checklist row with current artifacts. Use `pass`, `fail`, `cannot_verify`, `not_assessed`, or scoped variants such as `pass_with_advisory_findings`. Do not mark the objective complete while any explicit requirement is unverified or a blocker remains.
</completion_audit>

<outputs>
Keep outputs compact:
- changed files and commits
- exact commands run
- live PR check URL if available
- blockers vs advisory follow-ups
- trust boundaries that remain `not_assessed`
</outputs>

<supporting_files>
Use `templates/audit-matrix.md` to map broad quality requests to concrete artifacts, commands, reviewer planes, and unresolved trust states.
</supporting_files>

<red_flags>
- The audit reports a percentage or health score without named evidence.
- A command is listed as passing without fresh session output.
- A docs-only change claims verifier behavior changed.
- Security/trust review is folded into generic code review.
- Advisory findings are hidden because gates passed.
</red_flags>
