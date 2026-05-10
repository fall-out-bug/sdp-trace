# Full Pi Review: Product Contract v0

Date: 2026-05-10
Scope: `spec.md`, `plan.md`, `example.md`, `traceability.md`, `tasks.md`,
prior `003` review findings, and repository rules.
Status: review completed; critical/major findings unresolved.

## Review Runs

| plane | model | status | raw output |
| --- | --- | --- | --- |
| Product contract and hard gate | `zai/glm-5.1` | usable | `raw/2026-05-10-glm-contract-gate.md` |
| CTO value and Russian enterprise adoption | `openrouter/qwen/qwen3.6-plus` | first attempt `not_assessed`; retry usable | `raw/2026-05-10-qwen-cto-ru-retry.md` |
| Trust and evidence semantics | `minimax/MiniMax-M2.7` | usable, hidden reasoning preamble stripped from stored artifact | `raw/2026-05-10-minimax-trust.md` |
| Implementation readiness and DX | `openrouter/deepseek/deepseek-v4-pro` | usable | `raw/2026-05-10-deepseek-dx.md` |
| Adversarial overclaim and scope control | `openrouter/xiaomi/mimo-v2.5-pro` | usable | `raw/2026-05-10-xiaomi-adversarial.md` |
| Kimi slot from allowlist | `kimi-coding/k2p6` | `not_assessed` | current `pi` model registry did not match this pattern |

## Verdict

`REVISE_BEFORE_USER_REVIEW`

The direction is correct: Change Evidence Packet v0 is the right product
contract layer, and packet-row mapping is the right way to stop substrate-only
P0 work. The current draft is still too easy to bypass. Reviewers converged on
three core issues:

1. A backlog item can cite a packet row without proving forward progress.
2. The packet artifact and evidence bundle are not specified tightly enough for
   implementation or review.
3. Theater, verification, profile, decision-owner, and local-enterprise
   semantics still contain ambiguity that could create evidence theater.

## Consolidated Findings

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| PCV0-001 | critical | hard gate | The backlog gate is currently a citation gate, not a progress gate. | `unresolved_blocker` | `spec.md:150` requires `packet_rows`, but no rule requires a forward transition or measurable buyer improvement. `traceability.md:42` shows the intake template but no validation semantics. |
| PCV0-002 | critical | hard gate | Permanent `not_assessed` can still pass the gate if a row is merely named. | `unresolved_blocker` | `spec.md:59` correctly preserves missing inputs, but `spec.md:150` does not state that repeated `not_assessed` without new evidence is not P0 progress. |
| PCV0-003 | critical | packet artifact | The Markdown packet template is not normative. | `unresolved_blocker` | `spec.md:31` names Markdown plus evidence bundle; `example.md:23` implies metadata and tables, but `spec.md` does not require exact sections, fields, or allowed row-cell values. |
| PCV0-004 | critical | evidence bundle | The evidence bundle format is undefined. | `unresolved_blocker` | `spec.md:34` says retained refs, digests, and redaction status, but no bundle manifest, directory shape, or ref-resolution contract exists. |
| PCV0-005 | critical | evidence semantics | P0 theater reason codes lack derivation and independence rules. | `unresolved_blocker` | `spec.md:111` defines theater conditions; `traceability.md:64` admits reason-code derivation is not implemented. Independent retained evidence is not defined. |
| PCV0-006 | major | evidence states | State vocabulary is incomplete. | `unresolved_blocker` | `spec.md:85` lists `pass`, `partial`, `fail`, `not_assessed`, `cannot_verify`, `missing_telemetry`, `unsupported`, and `not_integrated`, but does not define them. |
| PCV0-007 | major | source strength | Source strength could be turned into an implicit score. | `unresolved_blocker` | `spec.md:94` lists source classes and says they are not a trust score, but does not forbid ranking, aggregation, or projection as confidence. |
| PCV0-008 | major | profile semantics | Evidence profile taxonomy is undefined and inconsistent with the example. | `unresolved_blocker` | `spec.md:49` names `local`, `change-host-rich`, `harness-observed`, `signed`; `example.md:30` uses `local-enterprise-baseline-v0`. Required inputs and maximum achievable row states are not defined per profile. |
| PCV0-009 | major | example semantics | `PC-VERIFICATION` example uses `cannot_verify` despite partial harness evidence. | `unresolved_blocker` | `example.md:42` says harness observed a verification command but marks the row `cannot_verify`; reviewers expect `partial` plus theater/gap finding. |
| PCV0-010 | major | theater semantics | `PC-THEATER` row state treats findings as failure. | `unresolved_blocker` | `example.md:45` marks `PC-THEATER` as `fail` because findings exist. Reviewers argue the row should represent assessment coverage, while findings carry severity. |
| PCV0-011 | major | residual gaps | `PC-RESIDUAL-GAPS` synthesis is undefined and `pass` is ambiguous. | `unresolved_blocker` | `example.md:48` uses `pass` because gaps are recorded, but no rule says the row auto-populates from non-pass rows and active theater findings. |
| PCV0-012 | major | decision ownership | Decision-owner binding is too weak for enterprise use. | `unresolved_blocker` | `spec.md:78` allows named role or owner ref; `example.md:47` asserts service tech lead with `cannot_verify`. Reviewers ask for minimum owner-binding sources under local baseline. |
| PCV0-013 | major | Russian enterprise baseline | The example still depends on OpenCode/GSD, so it does not prove local-only baseline usefulness. | `unresolved_blocker` | `example.md:40` uses OpenCode/GSD under `local-enterprise-baseline-v0`; `spec.md:134` allows local Git, GitFlic/GitLab, Jenkins/TeamCity, private artifacts. No local Git plus internal CI example exists. |
| PCV0-014 | major | authority boundary | `PC-AUTHORITY` needs a limited state vocabulary to avoid policy/blame drift. | `unresolved_blocker` | `spec.md:75` says no policy verdict, but does not restrict authority states. Existing authority substrate should project facts, not compliance or blame. |
| PCV0-015 | major | traceability | Traceability overstates readiness for `PC-AUTHORITY`, `PC-THEATER`, and `PC-RESIDUAL-GAPS`. | `unresolved_blocker` | `traceability.md:17` and `traceability.md:21` say "Good substrate" while the same rows still need packet projection or synthesis rules. |
| PCV0-016 | major | signed attestation | The additive signing rule is right but not enforceable enough. | `unresolved_blocker` | `spec.md:235` says signing is not a shortcut, but does not state that signing a packet cannot upgrade any row's `not_assessed` or `cannot_verify` state. |
| PCV0-017 | major | contract self-status | Draft tasks are checked off without distinguishing draft-complete from reviewed/approved. | `unresolved_blocker` | `tasks.md:11` marks T001-T008 complete. Reviewers worry this repeats "draft artifact equals product progress." Need explicit task/status semantics or a note that checked items are draft-complete only. |
| PCV0-018 | major | review packet completeness | Full-review packet omitted updated `003` plan/tasks, causing one reviewer to question the linkage. | `accepted_narrower` | Actual files do contain the blocker: `003` plan lines 13-15 and 186-190; `003` tasks lines 21-25. The valid narrower fix is to include updated `003` files in any re-review packet. |

## Required Before User Approval

1. Add a forward-progress rule: a P0 item must improve at least one cited row by
   adding new evidence, narrowing uncertainty, changing a row state, or
   documenting unsupported/not-integrated with closure evidence. Mere citation
   does not count.
2. Define all evidence states and forbid source-strength ranking or scoring.
3. Define the normative Markdown packet template and evidence bundle manifest.
4. Define profile taxonomy: `local-enterprise-baseline-v0`,
   `change-host-rich-v0`, `harness-observed-v0`, `signed-v0`; include required
   inputs and maximum achievable states per row.
5. Define theater derivation rules for the four P0 reason codes, including
   "independent" and "retained" evidence.
6. Fix example row semantics: `PC-VERIFICATION`, `PC-THEATER`,
   `PC-RESIDUAL-GAPS`, and `PC-DECISION`.
7. Add a local-only Russian enterprise example: local Git plus Jenkins/TeamCity
   or equivalent internal CI/artifact refs, no OpenCode/GSD dependency.
8. Define `PC-AUTHORITY` states as fact states only, not compliance/blame
   verdicts.
9. Clarify that signed attestation does not change underlying row states.
10. Update task/status language so draft-complete artifacts are not confused
    with reviewed/approved product progress.

## What Survived Review

- Product Contract v0 is the right structural move.
- Change Evidence Packet v0 is a valid first buyer-facing artifact.
- Markdown plus evidence bundle is a defensible canonical artifact for
  offline/on-prem enterprise use.
- Packet-row mapping is the right backlog gate, once hardened with forward
  progress semantics.
- Current `sdp-trace` substrate is valuable; it needs projection into packet
  rows rather than another substrate-only feature pass.

## Review Evidence Notes

- The first Qwen CTO/RU-market attempt produced a tool-call request instead of
  review findings and is `not_assessed`.
- The Qwen retry used the same review plane with a shorter bounded prompt and
  produced usable findings.
- MiniMax output included a hidden reasoning preamble; that preamble was
  stripped from the stored raw artifact before this ledger was written.
- All usable reviewer outputs were treated as evidence input, not approval.
