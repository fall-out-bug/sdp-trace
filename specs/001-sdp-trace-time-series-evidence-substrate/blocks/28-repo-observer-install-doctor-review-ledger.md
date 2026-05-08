# Block 28 Review Ledger: Repo Observer Install And Doctor UX

Status: Socratic spec review complete; implementation review complete; PR
pending.

## Spec Review Inputs

- Socratic review:
  `specs/001-sdp-trace-time-series-evidence-substrate/blocks/28-repo-observer-install-doctor-socratic.md`
- Spec under review:
  `specs/001-sdp-trace-time-series-evidence-substrate/blocks/28-repo-observer-install-doctor.md`
- UX/DX pi review: `zai/glm-5.1`, raw local scratch
  `.codex-review/block28/ux-dx-zai-glm-5.1.txt`
- Trust/evidence pi review: `minimax/MiniMax-M2.7`, raw local scratch
  `.codex-review/block28/trust-minimax-m2.7.txt`
- Requirements/spec pi review: `openrouter/qwen/qwen3.6-plus`, raw local
  scratch `.codex-review/block28/requirements-qwen3.6-plus.txt`
- Focused re-review after fixes: `openrouter/deepseek/deepseek-v4-pro`, raw
  local scratch `.codex-review/block28/focused-rereview-deepseek-v4-pro.txt`
- Implementation code/correctness pi review: `minimax/MiniMax-M2.7`, raw local
  scratch `.codex-review/block28-implementation/code-correctness-minimax-m2.7.txt`
- Implementation tracing/evidence pi review: `zai/glm-5.1`, raw local scratch
  `.codex-review/block28-implementation/tracing-evidence-zai-glm-5.1.txt`
- Implementation requirements pi review: `openrouter/qwen/qwen3.6-plus`, raw
  local scratch `.codex-review/block28-implementation/requirements-qwen3.6-plus.txt`
- Focused implementation re-review after fixes:
  `openrouter/deepseek/deepseek-v4-pro`, raw local scratch
  `.codex-review/block28-implementation/focused-rereview-deepseek-v4-pro.txt`

## Findings And Disposition

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| S28-UX-01 | major | UX/DX | Single `observer_state` lets a clean install look like proof. | accepted_fixed | Spec now requires `install_state` for setup and `proof_state` for live evidence. |
| S28-SAFE-01 | critical | safety | `repository_root` conflicts with the ban on private absolute paths in machine output. | accepted_fixed | Spec now uses safe `repository_id` and `repository_root_ref`; raw absolute paths are forbidden in machine status. |
| S28-TRUST-01 | major | trust boundary | Local hooks are bypassable but the first draft did not require surfacing that fact. | accepted_fixed | Spec adds `local_hooks_bypassable` and requires `doctor` to report hook evidence as `local_structural`. |
| S28-CONTRACT-01 | major | machine contract | Schema version was named but no schema artifact was required. | accepted_fixed | Acceptance criteria require a JSON schema and fixture validation. |
| S28-UX-02 | critical | UX/DX | `install` and `doctor` used inconsistent profile names. | accepted_fixed | Spec now uses `github-actions-git-hooks-v1` for both commands and defines `install repo-observer` as setup subcommand. |
| S28-UX-03 | critical | UX/DX | Human output blended install and proof state in one `State` column. | accepted_fixed | Human output table now has separate `Install state` and `Proof state` columns. |
| S28-TRUST-02 | critical | trust/evidence | Trust scopes did not define which scopes can satisfy proof. | accepted_fixed | Added proof-satisfying and non-proof trust scope semantics. |
| S28-CONTRACT-02 | critical | machine contract | `not_assessed` and `cannot_verify` were listed as trust scopes instead of states. | accepted_fixed | Removed them from trust scope list and added state/scope separation text. |
| S28-CONTRACT-03 | critical | machine contract | `repository_id` derivation was unconstrained and could leak local identity. | accepted_fixed | Spec now derives safe ids from sanitized remote hashes and forbids path-derived ids. |
| S28-INSTALL-02 | major | repo safety | `.gitignore` mutation strategy was unspecified. | accepted_fixed | Added marked block insertion and conflict/force behavior. |
| S28-INSTALL-03 | major | repo safety | `--force` backup/diff behavior was nondeterministic. | accepted_fixed | Spec now requires per-file `.bak` and safe diff summary. |
| S28-CONTRACT-04 | major | machine contract | `next_actions` schema was unspecified. | accepted_fixed | Added object schema and deduplicated aggregation rule. |
| S28-INSTALL-04 | major | repo safety | Idempotency was unspecified. | accepted_fixed | Added no-op behavior and `already_installed` reason code. |
| S28-DOC-01 | minor | DX | `--help` output contract was unspecified. | accepted_fixed | Acceptance criteria now require help contents and golden fixture. |
| S28-SAFE-02 | minor | safety | Focused re-review noted safe diff path and hook path overwrite ambiguity. | accepted_fixed | Spec now requires repo-relative safe diffs and refuses mismatched `core.hooksPath` without `--force`. |
| S28-PROVIDER-01 | minor | portability | GitHub Actions first profile could be misread as product dependency. | accepted_narrower | Spec keeps provider-specific profile naming and future-profile boundary; no further change needed. |
| S28-INSTALL-01 | minor | repo safety | `--force` overwrite behavior needs tests. | accepted_narrower | Spec already requires no silent overwrite and backup/safe diff on force; implementation must fixture-test existing files. |
| I28-CONFIG-01 | critical | requirements | Generated `.sdp-trace/config.json` omitted repository id, installed file manifest, and install metadata. | accepted_fixed | `sdpTraceConfig` now writes `repository_id`, `installed_files`, and stable `install_metadata`; package tests inspect the config. |
| I28-FORCE-01 | critical | repo safety | `--force` did not expose backup/diff-summary behavior. | accepted_fixed | Forced overwrites create `.bak` files and return repository-relative `force_diff_summary`; tests assert no absolute path leak. |
| I28-HOOKS-01 | critical | requirements | `core.hooksPath` mismatch was classified as double `cannot_verify` instead of misconfigured `fail`. | accepted_fixed | `hooksPathSurface` now reports `install_state=fail`, `proof_state=cannot_verify`, reason `hooks_path_mismatch`. |
| I28-REASON-01 | major | machine contract | Implementation introduced reason codes outside the closed spec vocabulary. | accepted_fixed | Custom generated-file and repository-id reason codes were removed from output and schema; example tests enforce closed reason codes. |
| I28-FIXTURE-01 | major | tracing/evidence | Example status fixtures had too few surfaces and could overclaim `install_state=pass`. | accepted_fixed | Fixtures were regenerated from live CLI output and now contain 13 surfaces matching implementation shape. |
| I28-SURFACE-01 | major | requirements | Minimum surfaces for PR/check binding and local wrapped commands were missing. | accepted_fixed | Added `pr_check_binding` and `local_wrapped_commands` surfaces as explicit `not_assessed` entries. |
| I28-WORKFLOW-01 | major | tracing/evidence | CI workflow surface reused `ci_artifact_upload_present` reason code. | accepted_fixed | Workflow presence now uses `ci_workflow_present`; artifact upload remains a separate surface. |
| I28-SAFE-01 | major | output safety | Absolute hook-path display could leak misleading path fragments. | accepted_fixed | `safeRef` now redacts absolute paths with `unsafe_absolute_path_redacted`; path containment uses `filepath.Rel`. |
| I28-SSH-01 | minor | identity safety | SSH origin URL sanitization kept the `git@` user component in the hash input. | accepted_fixed | SSH-style origin sanitization strips the user prefix before hashing. |
| I28-GITIGNORE-01 | critical | code review | Reviewer claimed `.gitignore` end-marker substring was off by one. | rejected_false_positive | In Go `strings.Index` returns the start of the marker; adding marker length is required to include the end marker. Removing it would corrupt replacement. |
| PR28-FORCE-01 | major | PR requirements | PR-level review found `force_diff_summary` had metadata but no safe before/after summary. | accepted_fixed | `DiffSummary` now includes safe sha/byte/line before/after summaries without raw file content. |
| PR28-FORCE-02 | major | PR requirements | Forced `core.hooksPath` replacement was not represented in `force_diff_summary`. | accepted_fixed | Forced hook-path replacement now appends a safe `git_config:core.hooksPath` summary with redacted prior value and `.githooks` target. |
| PR28-SCOPE-01 | major | PR requirements | Outside-profile PR/check and local wrapped-command surfaces used concrete trust scopes. | accepted_fixed | They now use `not_applicable`, `outside_profile_scope`, and exact no-action next action; spec was amended to use `not_applicable` for trust scope. |
| PR28-ORIGIN-01 | minor | PR requirements | Repository id derivation did not strip URL fragments. | accepted_fixed | `sanitizeOrigin` strips fragments before hashing. |
| PR28-HUMAN-01 | minor | PR requirements | Human output did not explicitly warn that `core.hooksPath` is local checkout configuration. | accepted_fixed | Human output now prints a local-checkout configuration note. |

## Remaining Approval State

No unresolved critical or major Socratic, spec pi-review, or implementation
pi-review findings remain. Focused implementation re-review returned `APPROVE`.

PR creation, PR-level review, and merge/post-merge verification are still
pending.
