---
name: pi-review
description: Orchestrate adversarial sdp-trace review planes, reviewer model policy, retries, finding verification, and disposition records.
---

<objective>
Produce usable independent review evidence for sdp-trace without counting hung, empty, off-task, or unchecked reviewer output as proof.
</objective>

<when_to_use>
Use this skill for Socratic review, implementation review, PR review, security/trust review, requirements-vs-implementation review, trace/evidence review, or requests naming pi-review.
</when_to_use>

<model_policy>
For adversarial review in this repo, prefer non-OpenAI, non-Anthropic, and non-Google models unless the user explicitly permits otherwise. Use current, provider-qualified model IDs verified by `pi --list-models <family>` in the same session when model freshness matters or when a reviewer previously failed. Do not use stale reviewer defaults such as `qwen3-coder`, `deepseek-chat-v3.1`, or `glm-4.6` unless every newer candidate is unavailable and the fallback is recorded as degraded. Record model, provider, retry, fallback, timeout, and replacement details in the review artifact, not in `AGENTS.md`.
</model_policy>

<review_planes>
For trust-sensitive work, run separate planes:
- code/correctness
- tracing/evidence/provenance
- requirements-vs-implementation
- security/forgery/overclaim when trust, credentials, authority, external inputs, or verdicts changed
- DX/UX when command surfaces, docs, packets, reports, or human workflows changed
</review_planes>

<process>
1. Build a bounded context pack with the objective, changed files, relevant specs, rules, and verification commands.
2. For non-trivial trust claims, use `workflows/claim-doubt-cycle.md` before accepting the claim.
3. Launch reviewers in parallel when possible, preferably through `codex-subagent panel run pi` for durable run IDs, logs, and structured results. Keep roles read-only unless the task explicitly assigns a worker. Pin each reviewer to a current provider-qualified model when using Pi.
4. Reject and replace reviewer output that is hung, empty, generic, off-task, or lacks file/line evidence for actionable claims.
5. Verify every finding against full files before accepting or rejecting it.
6. Record disposition as accepted, accepted_fixed, rejected_false_positive, deferred_not_assessed, cannot_verify, or advisory.
7. Re-run relevant verification after fixes.
8. For PR-ready trust work, run review rounds iteratively against the full diff. Fix every finding of any severity before the next round. Repeat until the reviewer outputs exactly `LGTM` (zero findings). Do not treat a partial review or non-zero finding count as sufficient.
9. Delete stale review artifacts rather than marking them superseded; a header or marker is insufficient to prevent stale claims from being cited as evidence.
10. Do not count review as merge/sign-off unless the user explicitly requested that authority and the repo process permits it.
</process>

<evidence_rules>
- A reviewer verdict is advisory work product until inspected.
- A green test run does not prove requirements coverage unless the test actually covers the requirement.
- Absent GitHub checks are CI `not_assessed`, never green.
- Checked-in review JSON is not authority unless live-verified or externally signed.
- Scanner counts claimed in docs must match the actual output of copy-pasteable commands. Verify command reproducibility (e.g., `gitleaks` auto-loads `.gitleaks.toml` from CWD; default-config scans require directory isolation or config exclusion).
- Stale review artifacts with incorrect counts or outdated findings must be deleted, not annotated.
</evidence_rules>

<output_format>
Summarize findings first, ordered by severity. Include file/line evidence, disposition, verification command, and remaining `not_assessed` areas.
</output_format>

<synthesis_hygiene>
When producing or updating `reviews/synthesis.md` (or any review ledger):
- Cross-check every claimed fix against the actual diff. If a file was reverted or removed from the diff, remove the claim about it.
- Do not list a verification command as passing unless it was actually run and produced fresh output in this session.
- Record CI status as `not_assessed` until live GitHub checks are queried for the final head.
- A green local test run is not a green CI run. Distinguish them explicitly.
</synthesis_hygiene>

<supporting_files>
- Use `workflows/claim-doubt-cycle.md` for any gate, verdict, provenance, release-proof, or non-obvious safety claim.
- Use `references/model-policy.md` when selecting pi/Codex/OpenCode/Kimi/GLM/MiniMax roles.
- Use `templates/review-disposition.md` for durable synthesis or handoff records.
- For implementation delegation rather than review, use `sdp-trace-pi-handoff`.
</supporting_files>

<red_flags>
- The reviewer received the author's conclusion instead of artifact plus contract.
- A reviewer was allowed to edit files during an allegedly independent review.
- Cross-model output is accepted without checking file/line evidence.
- A model/harness failure is hidden instead of recorded as fallback, replacement, `cannot_verify`, or `not_assessed`.
- The synthesis says "approved" while Critical/Important findings remain unresolved.
- Stale review artifacts with incorrect counts are left in the repository with only a "superseded" header.
- Scanner verification commands are not copy-pasteable, change the caller's working directory, or produce different counts depending on where they are run.
</red_flags>
