# Block 23 Review Ledger

Status: Socratic spec review returned `REVISE` across requirements/product,
trace/evidence, and Go-quality/docs planes. Valid critical and major findings
were accepted into the draft spec. Focused re-review is required before
implementation approval.

## Socratic Review Findings

| id | severity | plane | reviewer/source | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- | --- |
| S23-RP-01 | critical | requirements/product | MiniMax-M2.7 | `source_bound_local_release` pass was undefined. | accepted_fixed | Added Definitions for `pass`/`fail` and WS1 JSON/exit-code acceptance criteria. |
| S23-RP-02 | critical | requirements/product | MiniMax-M2.7 | Russian docs were critical without explicit MVP-scope rationale. | accepted_narrower | Added MVP-04 scope caveat: critical only if Russian handoff remains in MVP scope. |
| S23-RP-03 | critical | requirements/product | MiniMax-M2.7 | Block 06 open Beads and dangling proof refs were not enumerated. | accepted_fixed | WS2 now names `sdp-trace-drq.11` and `sdp-trace-drq.12`. |
| S23-RP-04 | major | requirements/product | MiniMax-M2.7 | CRAP gate contradicted non-goals and lacked an executable replacement. | accepted_fixed | WS3 now defines changed-function CRAP formula, scope, threshold, and exception rule. |
| S23-RP-05 | major | requirements/product | MiniMax-M2.7 | Customer questions were not enumerated. | accepted_fixed | WS5 now lists nine mandatory customer questions and deliverable files. |
| S23-RP-06 | major | requirements/product | MiniMax-M2.7 | `not_assessed` and `cannot_verify` were undefined. | accepted_fixed | Added Definitions and WS6 registry requirement. |
| S23-RP-07 | major | requirements/product | MiniMax-M2.7 | WS6 overreached into process governance. | accepted_narrower | WS6 narrowed to closure package and review evidence while keeping repo-required PR/merge acceptance boundaries. |
| S23-RP-08 | major | requirements/product | MiniMax-M2.7 | `gate` command docs could overclaim production trust. | accepted_fixed | WS4 now requires explicit caveats for `gate`, `witness`, and `release-proof`. |
| S23-TE-01 | critical | trace/evidence | ZAI/GLM-5.1 | Source-bound proof anchor and committed-vs-ephemeral artifact boundary were undefined. | accepted_fixed | Added `source-bound anchor` definition and WS1 proof-artifact boundary. |
| S23-TE-02 | critical | trace/evidence | ZAI/GLM-5.1 | Proof cycle could self-invalidate. | accepted_fixed | WS1 now states proof is generated against source commit and proof artifacts are not manifest subjects unless a new cycle is run. |
| S23-TE-03 | critical | trace/evidence | ZAI/GLM-5.1 | Stale ledger false claims could remain distinguished but uncorrected. | accepted_fixed | WS2 now requires annotation or historical separation plus grep check. |
| S23-TE-04 | major | trace/evidence | ZAI/GLM-5.1 | Manifest mismatch types were unclassified. | accepted_fixed | WS1 now requires per-subject drift classification before repair. |
| S23-TE-05 | major | trace/evidence | ZAI/GLM-5.1 | Review artifact output format was undefined. | accepted_fixed | WS6 now defines review artifact fields. |
| S23-TE-06 | major | trace/evidence | ZAI/GLM-5.1 | Verification plan lacked expected command states. | accepted_fixed | Verification Plan now lists expected exit/output states per command. |
| S23-QD-01 | critical | quality/docs | OpenRouter Qwen | CRAP < 5 had no Go tool or algorithm. | accepted_fixed | WS3 now defines CRAP formula and changed-function scope. |
| S23-QD-02 | critical | quality/docs | OpenRouter Qwen | Complexity and coverage thresholds were not numeric. | accepted_fixed | WS3 now names gocyclo and package coverage thresholds. |
| S23-QD-03 | critical | quality/docs | OpenRouter Qwen | Bilingual documentation "cover" was undefined. | accepted_fixed | WS4 now defines command/profile documentation coverage fields. |
| S23-QD-04 | major | quality/docs | OpenRouter Qwen | Workstreams were not truly independent and lacked merge order. | accepted_fixed | WS6 now states expected implementation order and limits parallel worktrees to disjoint write scopes. |
| S23-QD-05 | major | quality/docs | OpenRouter Qwen | Parked/dead code audit used speculative wording. | accepted_fixed | WS3 now requires `rg` import checks plus available Go dead-code tool and disposition. |
| S23-QD-06 | major | quality/docs | OpenRouter Qwen | Verification plan used `rtk` without definition. | accepted_fixed | Verification plan now uses underlying commands directly. |
| S23-QD-07 | major | quality/docs | OpenRouter Qwen | Customer answer map format was undefined. | accepted_fixed | WS5 now names retired customer-question map and `.ru.md`. |
| S23-QD-R1 | critical | quality/docs | OpenRouter MiMo | Closure package and review artifact schema were undefined. | accepted_fixed | WS6 now names closure package files and required table fields. |
| S23-QD-R2 | critical | quality/docs | OpenRouter MiMo | Package-level coverage scope was unbounded and could silently expand the block. | accepted_fixed | WS3 now gates changed production files/functions and reports package coverage for context. |
| S23-QD-R3 | major | quality/docs | OpenRouter MiMo | CRAP coverage operand needed 0-1 normalization and per-function extraction strategy. | accepted_fixed | WS3 now defines normalized coverage and `go tool cover -func` extraction. |
| S23-QD-R4 | major | quality/docs | OpenRouter MiMo | Exception row format and location were undefined. | accepted_fixed | WS3 now defines `block-23-quality-report.md` and exception row schema. |
| S23-QD-R5 | major | quality/docs | OpenRouter MiMo | Slice drift check was undefined. | accepted_fixed | WS6 now defines slice drift check bullets by surface. |
| S23-QD-R6 | major | quality/docs | OpenRouter MiMo | `golangci-lint` configuration ownership was assumed. | accepted_fixed | WS3 now forbids adding/relaxing lint config without approved spec update. |
| S23-QD-R7 | minor | quality/docs | OpenRouter MiMo | Dead-code tool name was vague. | accepted_fixed | WS3 now names `golang.org/x/tools/cmd/deadcode`. |

## Re-Review Requirement

Focused re-review was run on the revised Block 23 spec across:

- requirements/product closure;
- trace/evidence and source-bound proof;
- Go quality and bilingual documentation implementability.

Focused re-review returned `APPROVE` on all three planes with no remaining
critical or major findings. Remaining minor cautions:

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| S23-TE-R1 | minor | trace/evidence | `source_commit` absent/empty behavior needed explicit state. | accepted_fixed | Block 23 Definitions/WS1 now require `cannot_verify` with reason. |
| S23-TE-R2 | minor | trace/evidence | Retired-command grep list needed a concrete source. | accepted_fixed | WS2 now starts from `closed-block-task-drift-audit-2026-05-07.md`. |
| S23-TE-R3 | minor | trace/evidence | Manifest subject pre/post drift check should be explicit. | accepted_fixed | Verification Plan now includes manifest subject diff after regeneration. |
| S23-RR-01 | minor | quality/docs | CRAP join may need care for methods and anonymous functions. | deferred_not_assessed | Implementation should record any join ambiguity in `block-23-quality-report.md`; not a spec blocker. |
| S23-RR-02 | minor | quality/docs | `export cross-repo-posture` doc shape must match actual help tree. | accepted_fixed | WS4 requires docs to match current `sdp-trace --help`; implementation must verify help tree. |

Implementation remains blocked until the CTO explicitly approves the reviewed
Block 23 direction.
