# Product Contract v0 Review Prompt

You are an independent reviewer for `sdp-trace` Product Contract v0.

Review target:

- `specs/005-product-contract-v0/spec.md`
- `specs/005-product-contract-v0/plan.md`
- `specs/005-product-contract-v0/example.md`
- `specs/005-product-contract-v0/traceability.md`
- `specs/005-product-contract-v0/tasks.md`
- relevant prior roadmap review findings in
  `specs/003-agent-supply-chain-roadmap/reviews/2026-05-10-pi-socratic-review.md`
- repository rules in `AGENTS.md`

This is a product/spec review only. Do not implement. Do not rewrite the target.

Context:

`sdp-trace` has strong trust substrate work: provenance, evidence, trace,
authority, witness, CI artifact, harness observation, adapter capture, review,
and verifier semantics. The repeated product failure is that substrate features
can be treated as P0 product progress even when no buyer-facing artifact exists.

Product Contract v0 attempts to fix this by making Change Evidence Packet v0 the
first buyer-facing output and requiring every P0 feature to map to stable packet
rows.

Primary review question:

Does Product Contract v0 create a real hard gate that prevents substrate-only
work from being treated as P0 product progress?

Secondary questions:

- Is Change Evidence Packet v0 concrete enough to guide implementation later?
- Is the example packet useful to a CTO, or still too internal?
- Does the traceability matrix preserve current work while honestly exposing
  gaps?
- Does the Russian enterprise baseline work without GitHub, SaaS dashboards,
  public Sigstore/Rekor, raw prompt export, or employee surveillance?
- Are the backlog-gate fields enforceable enough for future roadmap intake?
- Are there hidden overclaims, missing states, or signed-attestation shortcuts?

Return only this structure:

1. Verdict: `APPROVE_FOR_USER_REVIEW`, `REVISE_BEFORE_USER_REVIEW`, or
   `KILL_OR_REFRAME`.
2. What works.
3. Blocking findings table with columns: id, severity (`critical`, `major`,
   `minor`), cited file:line, finding, why it matters, exact fix.
4. Non-blocking concerns.
5. Missing evidence or `not_assessed` areas.
6. Scope-control and overclaim risks.
7. One strongest reason to proceed.
8. One strongest reason not to proceed yet.

Be concrete. Cite the line-numbered packet. Prefer product-contract failures,
buyer artifact gaps, gate enforceability gaps, Russian enterprise realism,
evidence semantics, and implementation-readiness risks over wording polish.
