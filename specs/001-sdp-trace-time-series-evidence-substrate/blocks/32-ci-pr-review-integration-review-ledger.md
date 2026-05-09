# Block 32 Review Ledger: CI PR Review Integration

Status: draft spec created; Socratic review pending.

## Recipe Search Notes

Local recipes found:

- `~/.pi/agent/settings.json` enables `zai/glm-5.1`,
  `minimax/MiniMax-M2.7`, and `kimi-coding/k2p6`.
- `~/.pi/agent/models.json` defines provider/model metadata for `zai`,
  `minimax`, and `kimi-coding`.
- `examples/self-trace/provenance-records.json` records prior successful
  provenance command shapes:
  - `pi --provider zai --model glm-5.1 --no-tools --no-context-files`
  - `pi --provider kimi-coding --model k2p6 --no-tools --no-context-files`
  - `pi --provider minimax --model MiniMax-M2.7 --no-tools --no-context-files`
- `sdp_lab/prompts/skills/spec-interrogate/SKILL.md` uses provider rotation
  between `zai/glm-5.1`, `kimi-coding/k2p6`, and `minimax/MiniMax-M2.7`.
- `pi --help` lists `ZAI_API_KEY`, `KIMI_API_KEY`, and `MINIMAX_API_KEY` as
  environment-variable credentials.

External recipes checked:

- Pi provider docs confirm API-key providers can be configured through
  environment variables; relevant variables are `ZAI_API_KEY`,
  `KIMI_API_KEY`, and `MINIMAX_API_KEY`.
- Kimi Code docs confirm `KIMI_API_KEY` is intended for CI/CD style injection.
- MiniMax OpenCode docs confirm OpenCode can be configured for MiniMax, but
  Block 32 keeps OpenCode out of first implementation because `pi` is already
  the existing review runner pattern.

## Socratic Review Findings

Raw review outputs are local scratch under `.codex-review/block32-socratic/`
and are not committed.

| ID | Severity | Plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| S32-PB-01 | critical | product boundary / merge confusion | Spec said review evidence is not approval, but did not prevent a green `pr-review` check from being mistaken for required merge approval. | accepted_fixed | Spec now requires the job/check name `pr-review-evidence-only`, Step Summary evidence-only wording, and `docs/ci-check-policy.md` warning that the check is not human approval or production trust. |
| S32-PB-02 | major | secret boundary | The original shape checked out, built, and ran PR-head code in the same job that would hold provider secrets, allowing malicious PR code to exfiltrate keys. | accepted_fixed | Spec now requires a trusted-base `sdp-trace` binary and treats PR-head code only as diff/context data in the secret-bearing model-review step. |
| S32-PB-03 | major | UX / evidence visibility | Artifact-only output made review evidence hard for PR authors and reviewers to see. | accepted_narrower_fixed | First implementation remains read-only but must write a GitHub Step Summary with an evidence-only heading. PR comments are deferred to a later write-permission block. |
| S32-IMPL-01 | critical | implementation feasibility | Ordinary CI state was underspecified and race-prone while the review job runs concurrently with `verify`. | accepted_fixed | First implementation records packet `ci_state` as `not_assessed`; ordinary CI remains the separate `verify` job. Checks polling is deferred to a later evidence contract. |
| S32-IMPL-02 | major | secret handling | Missing-secret behavior conflicted between per-plane `not_assessed` and `coverage_partial`, and could reveal which provider secrets exist. | accepted_fixed | Spec now records all required model planes as `not_assessed` with reason `ci_model_review_not_configured` when any required secret is absent. |
| S32-IMPL-03 | major | runner setup | `pi` installation, version pinning, input passing, timeout, and stderr safety were underspecified. | accepted_fixed | Spec and plan now require pinned `pi`, trusted prompt input handling, 600s default/900s max timeout, and safe stderr normalization/digest-only retention. |
| S32-IMPL-04 | major | buildability | Reviewer questioned whether Block 30 has a `pi` runner adapter. | rejected_false_positive | `internal/prreview` already has `RunnerPI`, allowed-runner gating, external command execution, command digesting, and tests for `RunnerPI`. Block 32 needs CI profile/workflow wiring, not a new runner type. |
| S32-TE-01 | major | trace/evidence | `ci_secret_unavailable` and `digest_only` were not anchored enough as evidence states. | accepted_fixed | Reason changed to `ci_model_review_not_configured`; `digest_only` is now defined as retaining SHA-256 digest and safe ref metadata while discarding raw bytes. |

## Socratic Review State

- Product-boundary review: `zai/glm-5.1` returned `REWORK`; valid findings
  accepted and fixed in the draft.
- Implementation/secret review: `kimi-coding/k2p6` returned `REWORK`; valid
  findings accepted and fixed or rejected with code evidence.
- Trace/evidence review: first MiniMax attempt was off-task and is
  `not_assessed`; retry returned `REWORK` with medium findings that overlap the
  accepted fixes above.
- Focused re-review: `zai/glm-5.1` returned `APPROVE`; no remaining critical
  or major blockers were reported for the five accepted blocker classes.

## Implementation Review Findings

Raw review outputs are local scratch under `.codex-review/block32-implementation/`
and are not committed.

| ID | Severity | Plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| I32-CODE-01 | critical | code/correctness | Reviewer claimed `--not-assessed-reason` bypassed runner allowance checks. | rejected_false_positive | `--not-assessed-reason` is a no-execution path for recording missing CI model-review evidence. It never invokes the runner, so `runner_not_allowed` is not the relevant execution gate. `TestRunReviewNotAssessedReasonDoesNotInvokeRunner` proves no command runs. |
| I32-CODE-02 | critical | code/correctness | Prompt template read failures were silently swallowed by `renderPrompt`, risking empty prompts. | accepted_fixed | `renderPrompt` now returns an error and `runRole` records `cannot_verify` with `prompt_ref_cannot_verify`; `TestRunReviewCannotVerifyUnreadablePromptTemplate` covers the path. |
| I32-CODE-03 | major | code/correctness / CI permissions | Workflow/spec used `actions: write` while describing it as artifact-upload-only. | accepted_fixed | Workflow and spec now use `actions: read`, preserving no write permission for PR/content surfaces. |
| I32-TE-01 | minor | trace/evidence | Non-digest raw-output retention path still writes bytes while marking refs `digest_only`. | deferred_not_assessed | Pre-existing Block 30 behavior outside the CI digest-only path. Block 32 CI profile uses `raw_output_retention: digest_only`, and tests assert raw bytes are not persisted for that mode. |
| I32-TE-02 | minor | trace/evidence | Prompt template refs are trusted profile paths without a broader path-containment policy. | deferred_not_assessed | Profiles remain trusted inputs in this block. CI uses trusted-base profile files, not PR-head profile code. |
| I32-REQ-01 | minor | requirements-vs-implementation | Missing prompt-template failure path lacked direct test coverage. | accepted_fixed | Added `TestRunReviewCannotVerifyUnreadablePromptTemplate`. |
| I32-REQ-02 | minor | requirements-vs-implementation | Simple `{{key}}` replacement is intentionally limited. | deferred_not_assessed | Accepted as sufficient for Block 32 deterministic prompt fields; richer template semantics are future work. |

## Implementation Review State

- Code/correctness: `minimax/MiniMax-M2.7` returned `REWORK`; valid findings
  fixed or rejected with evidence.
- Trace/evidence: `zai/glm-5.1` returned `APPROVE` with minor notes.
- Requirements-vs-implementation: `kimi-coding/k2p6` attempt returned repeated
  tool-read requests instead of reviewing the supplied packet and is
  `not_assessed`; replacement
  `openrouter/qwen/qwen3.6-plus` returned `APPROVE` with minor notes.
- Critical findings remaining: 0
- Major findings remaining: 0

## PR Review Findings

Raw review outputs are local scratch under `.codex-review/block32-pr/` and are
not committed.

| ID | Severity | Plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| PR32-REQ-01 | major | requirements-vs-implementation | Reviewer claimed `actions/checkout@v6`, `actions/setup-go@v6`, `actions/setup-node@v6`, and `actions/upload-artifact@v6` do not exist and would fail the workflow. | rejected_false_positive | Live upstream checks on 2026-05-10 showed official v6 releases exist for the referenced actions; downgrading would be stale advice. |
| PR32-TE-01 | major | tracing/evidence | T242 was checked off without proving the full uploaded artifact pipeline redacts unsafe structured reviewer text from `results.json`, `ledger.json`, `validation.json`, and `summary.md`. | accepted_fixed | `parseReviewerOutput` now sanitizes structured reviewer findings before `results.json` is written; `TestRunReviewArtifactPipelineRedactsUnsafeReviewerText` injects synthetic token, prompt, authenticated URL, and private-path markers and asserts the workflow upload artifact candidates do not contain them. |
| PR32-TE-02 | major | tracing/evidence | The `pr-review-pi.sh` wrapper passes rendered prompts as a `pi -p` command argument, which can hit argument length limits for large prompts. | rejected_narrower | Current Block 32 prompt templates intentionally pass metadata and packet identifiers only, not PR diff bodies. The reviewer identified a valid future robustness risk, but it is not a blocker for this CI profile. |

## PR Review State

- Requirements-vs-implementation: `openrouter/qwen/qwen3.6-plus` returned
  `REWORK`; one major was rejected as stale after live upstream verification.
- Tracing/evidence: `zai/glm-5.1` returned `REWORK`; valid T242 artifact
  pipeline finding accepted and fixed.
- Code/correctness: `minimax/MiniMax-M2.7` output was not clean findings-only
  JSON and is recorded as lower-quality review evidence; no additional
  critical or major blocker was accepted from that output.
- Critical findings remaining: 0
- Major findings remaining: 0
