# Block 30 Review Ledger: Automated PR Review Evidence Mechanism

Status: Spec review, implementation review, and PR-level review completed; no
remaining critical or major findings are recorded for Block 30.

Raw review outputs are local scratch under `.codex-review/block30/` and are not
committed.

## Review Attempt Status

| plane | model | status | evidence status |
| --- | --- | --- | --- |
| UX/DX operator workflow | `zai/glm-5.1` | initial REVISE, focused re-review APPROVE | counted |
| product boundary / overclaim | `minimax/MiniMax-M2.7` | first attempt empty, retry CONDITIONAL APPROVAL, focused re-review APPROVE | retry counted; first attempt not_assessed |
| trace/evidence/provenance | `openrouter/qwen/qwen3.6-plus` | initial REVISE, focused re-review APPROVE | counted |
| security/privacy/output safety | `openrouter/deepseek/deepseek-v4-pro` | initial REVISE, focused re-review APPROVE | counted |
| implementation feasibility | `openrouter/xiaomi/mimo-v2.5-pro` | initial REVISE | counted |

## Critical And Major Findings

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| UX-01 | critical | UX/DX | No command produced or updated the disposition ledger required by `validate` and `summarize`. | accepted_fixed | Added `pr-review synthesize`; implementation plan now generates initial `unresolved_review_blocker` dispositions. Focused UX re-review approved. |
| UX-02 | critical | UX/DX | The spec did not explain how an operator produces packet inputs from a real PR. | accepted_fixed | Added Packet Input Guide with diff, metadata, context, verification, and explicit `--ci-state`. Focused UX re-review approved. |
| UX-03 | major | UX/DX | `repo_id` and `change_ref` formats were undefined. | accepted_fixed | Added safe regexes for `repo_id` and `change_ref`. |
| UX-04 | major | UX/DX | There was no preview before external model execution. | accepted_fixed | Added `pr-review run --preview` with planned roles, models, prompt digests, command digests, timeouts, and output paths. |
| UX-05 | major | UX/DX | Review profiles required hand-authoring without default, template, or validation surface. | accepted_fixed | Spec now requires profile validation without running reviewers and an example `trust-sensitive-default` profile fixture. |
| UX-06 | major | UX/DX | `--out` directory and file behavior was undefined for re-runs. | accepted_fixed | Spec now refuses non-empty/existing outputs unless reviewed `--force` behavior is added; cross-packet directory reuse is `cannot_verify`. |
| UX-07 | major | UX/DX | `off_task` detection was undefined. | accepted_fixed | Failure mapping now defines `off_task` as parsed output declaring wrong packet, plane, or role. |
| UX-08 | major | UX/DX | The common workflow required multiple manual commands with no single command path. | accepted_fixed | Added `pr-review check` common-path wrapper. |
| TE-01 | critical | trace/evidence | Packet context and verification refs lacked per-ref digests. | accepted_fixed | Safe-ref shape now requires `digest_sha256` for diff, context, verification, prompt, and raw-output refs. Focused trace re-review approved. |
| TE-02 | critical | trace/evidence | OpenCode mutation checks were ambiguous on a dirty working tree. | accepted_fixed | Added `clean_required` default and explicit `dirty_baseline` mode with `cannot_verify` intersections. |
| TE-03 | critical | trace/evidence | Runner failure mapping conflated `not_assessed` and `cannot_verify`. | accepted_fixed | Added deterministic failure mapping table. |
| TE-04 | major | trace/evidence | Fallback model provenance did not distinguish intended and actual models. | accepted_fixed | Added `requested_model`, `observed_model`, `fallback_for_model`, and `fallback_reason`. |
| TE-05 | major | trace/evidence | Re-review after fixes was not bound to a new packet digest. | accepted_fixed | Spec now requires a new packet and digest after any diff change; stale re-review closure is `cannot_verify`. |
| TE-06 | major | trace/evidence | Synthetic marker safety contract was underspecified. | accepted_fixed | Added marker classes, injection points, and non-echo requirements for validation, summary, and failure paths. |
| SEC-30-01 | major | security/privacy | OpenCode read-only enforcement relied on post-run mutation detection. | accepted_fixed | Spec now requires proactive read-only permission enforcement before execution; unavailable enforcement means `not_assessed`. Focused security re-review approved. |
| FE-01 | major | implementation feasibility | `context_refs` and `verification_refs` had no exact item shape. | accepted_fixed | Safe-ref shape now defines id, kind, ref, digest, content type, and redaction state. |
| FE-02 | major | implementation feasibility | The `pi` prompt contract was not deterministic enough to test. | accepted_fixed | Spec now requires prompt template refs or deterministic template placeholders and records final prompt digest. |
| FE-03 | major | implementation feasibility | Reviewer status enum overlapped semantically. | accepted_fixed | Removed generic `usable`; only `findings_reported` and `no_findings` count as usable. |
| FE-04 | major | implementation feasibility | `raw_output_ref` did not have a concrete type. | accepted_fixed | Raw output refs now use the safe-ref shape with closed redaction state. |
| FE-05 | major | implementation feasibility | External runner stdout contract was undefined for fake-runner tests. | accepted_fixed | Spec says external runners are opaque commands whose stdout must match the role's required output schema; plan requires sample fake outputs. |
| FE-06 | major | implementation feasibility | The user's enabled `pi` model list was treated like a machine contract. | accepted_fixed | Profile is now declarative; local `pi` settings are not inferred by `sdp-trace`. |
| FE-07 | major | implementation feasibility | Disallowed runner behavior conflicted between usage error and skip. | accepted_fixed | Spec now fails fast with usage error before any external runner is invoked. |
| FE-08 | major | implementation feasibility | Model identity fields were optional/ambiguous. | accepted_fixed | Model fields are required strings with `not_assessed` sentinel and requested-vs-observed mismatch handling. |
| FE-09 | major | implementation feasibility | The spec did not define who produces the ledger. | accepted_fixed | Added `pr-review synthesize`. |
| FE-10 | major | implementation feasibility | Synthetic marker tests lacked injection strategy. | accepted_fixed | Added marker class and injection table. |
| PB-01 | critical | product boundary | `complete` coverage state looked like an approval signal. | accepted_fixed | Renamed to `coverage_satisfied`; added authority fields and non-approval rendering. Focused product-boundary re-review approved. |
| PB-02 | critical | product boundary | Disposition vocabulary looked like merge-blocking policy. | accepted_fixed | Renamed blocking vocabulary to review-record terms and added structural authority boundary. |
| PB-03 | critical | product boundary | `pi` was described as a primary backend while non-dependency was claimed. | accepted_fixed | Reframed `pi` as first supported optional external runner. |
| PB-04 | major | product boundary | Summary could be collapsed into a gate by external consumers. | accepted_fixed | Summary must render authority structurally and avoid merge/readiness language. |
| PB-05 | major | product boundary | "AI does not approve" framing was prose, not a contract. | accepted_fixed | Spec now states the mechanism contractually prevents AI merge approval and every validation output carries authority fields. |

## Residual Minor Notes

Minor reviewer notes about command naming style, optional future generic review
namespace, and exact raw-output retention duration are recorded as
non-blocking implementation considerations. They do not block the reviewed spec
direction because the command namespace is explicitly PR-specific for source
binding, raw-output retention is already constrained by safe refs/redaction, and
generic review is a future extension.

## Final Spec Review State

- Critical findings remaining: 0
- Major findings remaining: 0
- Implementation approved: yes
- Approval required from: technical executive

## Implementation Review: T220-T225 Slice

Date: 2026-05-09

Raw review outputs are local scratch under `.codex-review/block30-implementation/`
and are not committed.

| plane | model | status | evidence status |
| --- | --- | --- | --- |
| code/correctness | `minimax/MiniMax-M2.7` | REVISE with mixed valid findings and false positives | counted after verification |
| trace/evidence/provenance | `openrouter/qwen/qwen3.6-plus` | REVISE | counted |
| requirements-vs-implementation | `kimi-coding/k2p6` | REVISE | counted |
| focused re-review | `openrouter/deepseek/deepseek-v4-pro` | APPROVE | counted |

### Implementation Findings

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| IMPL-30-01 | critical | trace/evidence | `coverage_unresolved` did not produce non-zero validation/check exit. | accepted_fixed | Added `reviewValidationExitCode`; focused re-review approved. |
| IMPL-30-02 | major | trace/evidence | `pr-review check` did not persist `runs/results.json`. | accepted_fixed | `check` now writes run provenance under `<out>/runs/results.json`; focused re-review approved. |
| IMPL-30-03 | major | trace/evidence | Packet digest replay semantics were implicit. | accepted_fixed | Added explicit canonical `packetDigest` helper clearing `packet_digest` before hashing and replay test; focused re-review approved. |
| IMPL-30-04 | major | trace/evidence | Duplicate `review_run_id` values were not rejected on imported run sets. | accepted_fixed | Added `validateRunSet` and duplicate-run test; focused re-review approved. |
| IMPL-30-05 | major | trace/evidence | Validation schema allowed arbitrary plane result status strings. | accepted_fixed | Constrained `plane_results[].status` enum; focused re-review approved. |
| IMPL-30-06 | major | requirements | Existing `--out <file>` targets could be overwritten by synthesize/validate/summarize. | accepted_fixed | Added `requireOutputFile`/`refuseExistingFile` checks; focused re-review approved. |
| IMPL-30-07 | major | requirements | Runner output parsing did not reject unknown JSON fields. | accepted_fixed | Switched reviewer output parser to `json.Decoder.DisallowUnknownFields`; focused re-review approved. |
| IMPL-30-08 | major | requirements | `check` lacked `--work-dir` parity with `run`. | accepted_fixed | Added `--work-dir` to `check` and preflight directory validation; focused re-review approved. |
| IMPL-30-09 | major | requirements | Packet actor provenance was only implicit. | accepted_fixed | Added `--created-by` to `packet` and `check`, defaulting to `sdp-trace-cli`; focused re-review approved. |
| IMPL-30-10 | major | code/correctness | `pr-review run` did not preflight invalid `--work-dir`. | accepted_fixed | Added `requireDirectory` before invoking runners; covered by focused re-review. |
| IMPL-30-11 | critical | code/correctness | Reviewer claimed `cmd/sdp-trace/pr_review_cli_test.go` referenced undefined `writeText`. | rejected_false_positive | Current file uses `writeFileStringForPRReviewTest`; `go test ./cmd/sdp-trace` passes. |
| IMPL-30-12 | minor | code/correctness | Repeated flag helper can retain duplicate values. | deferred_not_assessed | Non-blocking; repeated refs are preserved and covered. Deduplication is not required by the Block 30 contract. |

### Implementation Review State

- Critical findings remaining: 0
- Major findings remaining: 0
- Local verification after fixes: `go test ./cmd/sdp-trace ./internal/prreview`, `go test ./...`, `jq empty schema/*.json examples/pr-review/trust-sensitive-default.profile.json`, `git diff --check`
- PR-level review: assessed on PR #28; accepted trace/evidence blockers were
  fixed and focused trace/evidence re-review returned `APPROVE`.

## PR-Level Review: PR #28

Date: 2026-05-09

Raw review outputs are local scratch under `.codex-review/block30-pr/` and are
not committed.

| plane | model | status | evidence status |
| --- | --- | --- | --- |
| code/correctness | `minimax/MiniMax-M2.7` | no confirmed critical or major blocker in reviewed output | counted |
| trace/evidence/provenance | `openrouter/qwen/qwen3.6-plus` | REVISE | counted |
| requirements-vs-implementation | `openrouter/xiaomi/mimo-v2.5-pro` | REVISE with minor findings | counted |
| focused trace/evidence re-review | `openrouter/qwen/qwen3.6-plus` | APPROVE | counted |

### PR-Level Findings

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| PR30-TE-01 | major | trace/evidence | `ci_state` was captured in packets but not consumed or rendered by validation and summary output. | accepted_fixed | Added `Validation.CIState`, validation schema requirement, summary rendering, and regression coverage. |
| PR30-TE-02 | major | trace/evidence | OpenCode runner handling only recorded `read_only_enforced`; it did not establish a working-tree baseline or detect post-run mutations. | accepted_fixed | OpenCode roles now default to `clean_required`, refuse dirty starts, record safe baseline digest/count, and return `cannot_verify` with `mutation_detected` on mutation. |
| PR30-TE-03 | major | trace/evidence | Imported run sets validated only top-level packet digest, allowing individual stale reviewer results to count. | accepted_fixed | Validation now checks every `ReviewerResult.PacketDigest` and returns `cannot_verify` with `result_packet_digest_mismatch:<runID>` for stale results. |
| PR30-TE-04 | minor | trace/evidence | `coverage_partial` exits zero and may be mistaken for successful closure. | accepted_documented | This is intentional: partial review is not merge authorization and the validation JSON/summary carry explicit `coverage_partial`; non-zero is reserved for unresolved blockers and `cannot_verify`. |
| PR30-TE-05 | minor | trace/evidence | `required_output_schema` is declared but not used for generic runtime JSON Schema validation. | accepted_documented | The parser enforces the concrete Block 30 Go contract with unknown-field rejection; generic schema dispatch is outside this slice. |
| PR30-RI-01 | minor | requirements-vs-implementation | T236 was still unchecked while PR-level review was in progress. | accepted_fixed | T236 is checked only after PR-level fixes, focused re-review, fresh verification, and this ledger update. |
| PR30-RI-02 | minor | requirements-vs-implementation | `mutation_detected` was not listed in the failure mapping table. | accepted_fixed | Added the OpenCode mutation row to the failure mapping. |

### PR-Level Review State

- Critical findings remaining: 0
- Major findings remaining: 0
- Focused re-review: `openrouter/qwen/qwen3.6-plus` returned `APPROVE`
  after the PR-level fixes.
- Local verification after PR-level fixes: `go test ./internal/prreview
  ./cmd/sdp-trace`, `go test ./...`, `jq empty schema/*.json
  examples/pr-review/trust-sensitive-default.profile.json`, `git diff --check`
- GitHub CI: `verify` passed for PR #28 after the follow-up commit.
