# Block 30 Socratic Review: Automated PR Review Evidence Mechanism

Review date: 2026-05-09

Spec under review:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/30-automated-pr-review-evidence-mechanism.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/30-automated-pr-review-evidence-mechanism-implementation-plan.md`

Status: Socratic review complete. Implementation remains blocked until the
technical executive explicitly approves the reviewed direction.

Raw pi-review outputs are local scratch under `.codex-review/block30/` and are
not committed.

## Review Planes

Initial review used separate planes:

- UX/DX operator workflow: `zai/glm-5.1`
- product boundary and overclaim risk: `minimax/MiniMax-M2.7`
- trace/evidence/provenance: `openrouter/qwen/qwen3.6-plus`
- security/privacy/output safety: `openrouter/deepseek/deepseek-v4-pro`
- implementation feasibility: `openrouter/xiaomi/mimo-v2.5-pro`

The first MiniMax product-boundary attempt produced an empty artifact and was
not counted. A shorter retry produced usable findings and is the counted
product-boundary review.

Focused re-review after fixes used:

- UX/DX: `zai/glm-5.1` - APPROVE
- trace/evidence/provenance: `openrouter/qwen/qwen3.6-plus` - APPROVE
- security/privacy/output safety: `openrouter/deepseek/deepseek-v4-pro` - APPROVE
- product boundary and overclaim risk: `minimax/MiniMax-M2.7` - APPROVE

The local `pi` startup warning about the configured Kimi default model was not
review evidence and did not affect the explicitly selected models above.

## Socratic Questions And Resolutions

### Q1. Does the mechanism produce a review record, or does it accidentally become a merge gate?

**Critic**: The first draft used `complete`, `blocked`, `accepted_blocking`,
and `unresolved_blocker`, which were easy for external systems to collapse into
merge policy even though the prose disclaimed approval authority.

**Resolution**: Accepted and fixed. Coverage states are now
`coverage_satisfied`, `coverage_partial`, and `coverage_unresolved`.
Dispositions use review-record language such as `accepted_review_blocking` and
`unresolved_review_blocker`. Every validation output must include
`authority_scope=review_record_only` plus `merge_decision`,
`release_decision`, and `risk_acceptance` set to
`not_authorized_by_sdp_trace`.

### Q2. Can an operator actually run the workflow without hidden rituals?

**Critic**: The first draft required `validate --ledger <file>` but no command
created the ledger. It also accepted `--diff`, `--metadata`, `--context`, and
`--verification` without defining their input shape.

**Resolution**: Accepted and fixed. The command surface now includes
`pr-review synthesize` for initial ledger generation and `pr-review check` for
the common packet-run-synthesize-validate-summarize path. The spec now includes
a Packet Input Guide with required file formats, an explicit `--ci-state` flag,
and safe `repo_id` / `change_ref` patterns.

### Q3. Are review packet and citation bindings strong enough to prevent stale evidence?

**Critic**: `context_refs` and `verification_refs` lacked per-ref digests, and
citations could point at moving files rather than the frozen packet.

**Resolution**: Accepted and fixed. Safe refs now require `id`, `kind`, `ref`,
`digest_sha256`, `content_type`, and `redaction_state`. Every finding citation
must reference a packet context ref or diff ref and include a diff hunk id or a
source digest plus line range. Re-review after a diff change requires a new
packet digest.

### Q4. Are `not_assessed` and `cannot_verify` deterministic?

**Critic**: The first draft said unusable runner output became
`not_assessed` or `cannot_verify`, which left state mapping open to local
interpretation.

**Resolution**: Accepted and fixed. The spec now contains a failure mapping
table. Runner unavailable is `not_assessed`; timeouts, empty output, parse
failure, off-task structured output, stale packet, or unverifiable mutation
produce non-usable states and required-plane `cannot_verify` effects.

### Q5. Is OpenCode read-only review actually enforced?

**Critic**: The first draft relied too much on before/after git status, which
detects mutation after giving the agent too much authority.

**Resolution**: Accepted and fixed. OpenCode must use a permission profile that
denies write, edit, delete, and external mutation operations before the model
starts. If read-only enforcement is unavailable, the role remains
`not_assessed` and the agent is not executed. Working-tree checks remain only a
verification step.

### Q6. Is `pi` a hidden product dependency?

**Critic**: The first draft called `pi` the primary backend while also claiming
no hosted-model dependency.

**Resolution**: Accepted and fixed. `pi` is now the first supported optional
external runner. The mechanism remains usable for packet validation, ledger
validation, and summary rendering with reviewer execution `not_assessed` when
no runner is available.

### Q7. Are raw outputs, prompts, paths, and markers safe?

**Critic**: Synthetic-marker safety tests and raw-output refs were too vague.
Command provenance could leak absolute local paths.

**Resolution**: Accepted and fixed. Safe refs have closed redaction states.
Command provenance must use safe refs and digests instead of absolute local
paths. The spec defines marker classes for prompt, token, private path,
authenticated URL, and model-response fixtures, and requires `validate`,
`summarize`, and failure paths not to echo marker values.

## Approval Boundary

This Socratic review approves the amended spec direction for approval handoff.
It does not approve implementation. Implementation may start only after the
technical executive explicitly approves the reviewed direction.
