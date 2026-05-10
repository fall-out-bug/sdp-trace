# Pi Socratic Review: Agent Supply Chain Roadmap

Date: 2026-05-10
Scope: `spec.md`, `plan.md`, `research.md`, `tasks.md`, and repository rules.
Status: review completed; critical/major findings unresolved.

## Review Runs

| plane | model | status | raw output |
| --- | --- | --- | --- |
| SpecKit consistency and approval gate | `zai/glm-5.1` | usable | `raw/2026-05-10-glm-speckit.md` |
| Trust and evidence semantics | `minimax/MiniMax-M2.7` | usable, hidden reasoning preamble stripped from stored artifact | `raw/2026-05-10-minimax-trust.md` |
| CTO value and Russian enterprise adoption | `openrouter/qwen/qwen3.6-plus` | usable | `raw/2026-05-10-qwen-cto-market.md` |
| Integration sequencing and DX feasibility | `openrouter/deepseek/deepseek-v4-pro` | usable | `raw/2026-05-10-deepseek-integration.md` |
| Kimi slot from allowlist | `kimi-coding/k2p6` | `not_assessed` | current `pi` model registry did not match this pattern |

## Verdict

`REVISE_BEFORE_USER_REVIEW`

The direction is viable, but it is not ready for explicit roadmap approval.
All reviewers converged on the same core issue: the roadmap has good evidence
semantics, but the first CTO-facing packet remains too abstract. The approval
gate should not ask the owner to approve implementation scope until the first
packet shape, theater detection contract, local-market adapter posture, and
discovery closure criteria are sharper.

## Consolidated Findings

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| SR-001 | critical | product value | CTO packet format and artifact are undefined. | `unresolved_blocker` | `spec.md` asks which packet format creates the first wow, but leaves PR comment, archive, static HTML, Markdown report, and CLI summary open (`spec.md:270`). `plan.md` says P0-A must produce one sample packet, but not its concrete surface (`plan.md:98`). |
| SR-002 | critical | evidence semantics | Evidence theater is a taxonomy, not yet a binding detection/reporting contract. | `unresolved_blocker` | `spec.md` names eight theater conditions (`spec.md:192`) and scope says detection is in scope (`spec.md:40`), but FRs and tasks do not define minimum detection rows or closure criteria. |
| SR-003 | major | Russian enterprise adoption | GitHub-first is too dominant for a Russian-market enterprise target. | `unresolved_blocker` | Scope starts with GitHub-first packets (`spec.md:32`), FR-003 only names GitFlic as future capability (`spec.md:220`), and Phase 2 keeps non-GitHub providers as placeholders (`tasks.md:39`). |
| SR-004 | major | risk management | Integration risks are named without mitigations. | `unresolved_blocker` | The integration table lists risk for every target but has no mitigation column (`plan.md:86`). |
| SR-005 | major | traceability | Open questions and research gaps are not mapped to tasks, owners, or closure criteria. | `unresolved_blocker` | Open questions remain in `spec.md:268`; research gaps remain in `research.md:162`; tasks list discovery items but no gap-to-task-to-evidence mapping (`tasks.md:53`). |
| SR-006 | major | discovery method | `pi`, GSD2, Superpowers, Hermes/OpenClaw discovery methods are underspecified. | `unresolved_blocker` | Tasks say to inspect surfaces (`tasks.md:55`, `tasks.md:58`, `tasks.md:66`, `tasks.md:78`) but do not state whether inspection is runtime observation, source inspection, docs review, API probing, or fixture capture. |
| SR-007 | major | scope control | General-purpose agent boundary needs an enforceable software-delivery boundary. | `unresolved_blocker` | Scope excludes broad monitoring (`spec.md:45`) and FR-009 excludes general monitoring (`spec.md:234`), but no minimum evidence condition distinguishes software-delivery boundary from general employee-agent activity. |
| SR-008 | major | trust semantics | Signed attestation "top trust profile" needs operational meaning before enterprise discussion. | `unresolved_blocker` | Signed attestation is deferred (`spec.md:41`, `plan.md:165`), but open question 4 and Phase 7 leave minimum private-equivalent profile unresolved (`spec.md:276`, `tasks.md:85`). |
| SR-009 | major | product/review evidence | Review evidence itself should be recorded with a lightweight manifest, but full `assessment-input.json` is not required for this roadmap draft. | `accepted_narrower` | MiniMax overreached by implying full self-trace mirror mechanics are required for roadmap prose. The valid narrower point is that this review package needs a review manifest with files, models, commands, and `not_assessed` reviewer states. This file and the raw outputs provide that starting point. |
| SR-010 | minor | source quality | External source entries lack `last_checked` metadata. | `deferred_not_assessed` | `research.md:133` lists source URLs but no check dates. Useful before external-facing materials, not a blocker for internal roadmap review. |

## Top Socratic Questions For The Owner

1. What is the first CTO packet surface: PR comment, static HTML, Markdown,
   archive, CLI summary, or another artifact?
2. What is the minimum packet content that creates buyer value even when many
   rows are `not_assessed`?
3. Should Russian-market P0 include GitFlic/local Git/Jenkins-style artifact
   flow beside GitHub, or is GitHub-only acceptable for the pilot?
4. Which evidence theater findings are P0 rows in the first packet, and which
   are explicitly deferred?
5. What exact evidence closes a `pi` or GSD2 discovery row from `not_assessed`
   to importable, partial, wrapper-only, unsafe, or unstable?
6. What minimum technical boundary keeps general-purpose agent tracking inside
   software delivery and out of employee surveillance?
7. Is signed attestation additive evidence over an already meaningful packet,
   or a separate enterprise profile that may be required by some buyers from
   day one?

## Required Before Approval

Before asking for explicit roadmap approval:

1. Resolve the CTO packet format and add a first example shape.
2. Bind evidence theater taxonomy to minimum P0 packet rows or explicitly defer
   specific theater categories.
3. Add GitFlic/local Git/Jenkins-style Russian-market discovery posture, or
   explicitly state why it is not P0.
4. Add mitigations to the integration risk table.
5. Add gap-to-task-to-evidence closure mapping for open questions and research
   gaps.
6. Define discovery methods for `pi`, GSD2, Superpowers, and the selected
   general-purpose agent spike.
7. Define software-delivery boundary minimum evidence conditions.

## Evidence Commands

- `pi --no-tools --no-context-files --no-session --model zai/glm-5.1 --thinking high ...`
- `pi --no-tools --no-context-files --no-session --model minimax/MiniMax-M2.7 --thinking high ...`
- `pi --no-tools --no-context-files --no-session --model openrouter/qwen/qwen3.6-plus --thinking high ...`
- `pi --no-tools --no-context-files --no-session --model openrouter/deepseek/deepseek-v4-pro --thinking high ...`

All four completed with usable outputs. Each run also emitted a startup warning
that `kimi-coding/k2p6` from local settings did not match the current model
registry; this affects only the unused Kimi slot.
