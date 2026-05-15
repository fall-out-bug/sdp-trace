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
For adversarial review in this repo, prefer non-OpenAI, non-Anthropic, and non-Google models unless the user explicitly permits otherwise. Record model, retry, fallback, timeout, and replacement details in the review artifact, not in `AGENTS.md`.
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
2. Launch reviewers in parallel when possible. Keep roles read-only unless the task explicitly assigns a worker.
3. Reject and replace reviewer output that is hung, empty, generic, off-task, or lacks file/line evidence for actionable claims.
4. Verify every finding against full files before accepting or rejecting it.
5. Record disposition as accepted, accepted_fixed, rejected_false_positive, deferred_not_assessed, cannot_verify, or advisory.
6. Re-run relevant verification after fixes.
7. Do not count review as merge/sign-off unless the user explicitly requested that authority and the repo process permits it.
</process>

<evidence_rules>
- A reviewer verdict is advisory work product until inspected.
- A green test run does not prove requirements coverage unless the test actually covers the requirement.
- Absent GitHub checks are CI `not_assessed`, never green.
- Checked-in review JSON is not authority unless live-verified or externally signed.
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
