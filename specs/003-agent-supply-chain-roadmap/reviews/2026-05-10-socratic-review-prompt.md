# Socratic Review Prompt

You are an independent Socratic reviewer for an `sdp-trace` SpecKit roadmap
package.

Review target:

- `specs/003-agent-supply-chain-roadmap/spec.md`
- `specs/003-agent-supply-chain-roadmap/plan.md`
- `specs/003-agent-supply-chain-roadmap/research.md`
- `specs/003-agent-supply-chain-roadmap/tasks.md`
- repository rules in `AGENTS.md`

This is a roadmap/spec review only. Do not implement. Do not rewrite the
roadmap. Your job is to pressure-test whether the roadmap is good enough to ask
the human owner for approval before implementation.

Product goal:

`sdp-trace` should become an enterprise-grade portable evidence layer for
agentic software delivery. It records provenance, evidence, trace, gaps, and
attestation surfaces across coding agents, harnesses, change hosts, CI,
artifacts, reviewers, and human decision owners. It must remain independent
from specific agent runtimes, GitHub, GSD, OpenCode, Pi, Superpowers, Hermes,
OpenClaw, Claude, Codex, or any vendor.

Known buyer frame:

- Primary buyer: C-level, usually CTO.
- Employee-facing value: honest work, fewer unbacked done claims, less manual
  evidence archaeology.
- OSS is not the project substrate; OSS tools are integration targets.
- Russian-market enterprise adoption matters, but the product must not become a
  local-only fork of GitHub governance.
- Signed attestation is the top trust profile, not the first adoption step.

Critical repo constraints:

- No implementation before reviewed SpecKit direction is explicitly approved.
- Do not collapse `not_assessed`, `cannot_verify`, `missing_telemetry`,
  `unsupported`, or `not_integrated` into pass/fail.
- Do not let `sdp-trace` become a gate, GRC, SIEM, employee surveillance, or
  general agent-monitoring product.
- Tool support claims require inspected evidence surfaces.
- GitHub may be the first adapter, but not the product ontology.
- Harness/methodology artifacts are intent evidence unless separately verified.

Return only this structure:

1. Verdict: `APPROVE_FOR_USER_REVIEW`, `REVISE_BEFORE_USER_REVIEW`, or
   `KILL_OR_REFRAME`.
2. Top Socratic questions: 5-8 questions the owner must answer before
   implementation scope is approved.
3. Findings table with columns: id, severity (`critical`, `major`, `minor`),
   cited file:line, finding, why it matters, exact fix.
4. Missing evidence or `not_assessed` areas.
5. Scope-control risks.
6. One strongest reason to proceed.
7. One strongest reason not to proceed yet.

Be concrete. Cite the line-numbered packet. Prefer requirement gaps, UX risks,
evidence semantics, rollout order, and maintainability risks over wording
polish.
