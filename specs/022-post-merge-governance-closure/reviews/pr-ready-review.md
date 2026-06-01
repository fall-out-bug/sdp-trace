# Spec 022 PR-Ready Review

Date: 2026-06-01

Branch: `codex/022-post-merge-governance-closure`

Prompt class: full-diff PR-ready review for spec drift, constitution drift,
product drift, CRAP `< 5`, MI `> 70`, Clean Architecture hex boundaries,
Clean Code, SOLID, DRY, YAGNI, trust overclaiming, and review/evidence
completeness.

## Round 1

| Lane | Harness | Model | State | Disposition |
| --- | --- | --- | --- | --- |
| DeepSeek | `opencode run` | `openrouter/deepseek/deepseek-v4-pro` | findings | Found stale task-count and timestamp drift. Findings fixed in this branch. |
| Z.AI | `opencode run` | `opencode-go/glm-5.1` fallback for requested Z.AI lane | findings | Found stale task-count, task-state, and merge-approval ambiguity. Findings fixed in this branch. |
| Kimi | `opencode run` | `kimi-for-coding/k2p6` | `cannot_verify` | Lane hung after reading stale committed diff and never produced a final verdict. Process was stopped; output is not used as review authority. |

## Round 1 Retained Findings

| ID | Severity | Finding | Disposition |
| --- | --- | --- | --- |
| PRR1-M1 | major | Cross-surface task counts still reported `27 / 42`, `632 / 647`, `15 open`, and `T028-T042` after T028-T033 were checked. | Fixed by updating `docs/closure-decision-ledger.md`, `docs/spec-reality-ledger.md`, `docs/roadmap.md`, `docs/spec-closure-route.md`, and `docs/open-task-breakdown.md`; after final local verification the current state is `39 / 42`, `644 / 647`, 3 open tasks, and `T039/T041/T042`. |
| PRR1-M2 | major | `docs/open-task-breakdown.md` described completed US3 navigation tasks as still in progress. | Fixed by marking T028-T033 complete and pointing remaining work to final closure checks. |
| PRR1-m1 | minor | `docs/closure-decision-ledger.md` merge-approval row could be misread as covering PR #60. | Fixed by qualifying the approval as PR #64 only and keeping PR #60 merge approval `not_assessed`. |
| PRR1-m2 | minor | `docs/roadmap.md` and `docs/spec-reality-ledger.md` had stale final-update or live-verification dates. | Fixed to 2026-06-01. |
| PRR1-m3 | minor | T037 asked for `not_assessed` records for docs-only work while final evidence had actual test/vet/schemadoc passes. | Fixed by making T037 require the actual checks used for final closure. |

## Verified Review Planes

| Plane | State | Notes |
| --- | --- | --- |
| Spec drift | fixed | Round 1 stale task-count drift fixed. |
| Constitution drift | verified | No retained findings. |
| Product drift | verified | No product code changed in this slice. |
| CRAP `< 5` | verified | Local strict-less CRAP check passed and is recorded in `final-evidence.md`. |
| MI baseline | verified | Baseline checks passed and are recorded in `final-evidence.md`. |
| Absolute file MI `> 70` | verified | Passed in local final-evidence run. |
| Absolute function MI `> 70` | failed | Historical `tools/hygienecheck/check_demo_drift.go:10 checkCurrentDemoRepoDrift` remains MI `61.5`; recorded as assessed gap, not pass. |
| Clean Architecture hex | not_applicable | No product architecture code changed. |
| Clean Code / SOLID / DRY / YAGNI | verified | No retained findings after docs drift fix. |
| Trust overclaiming | verified | PR #60 merge approval remains `not_assessed`; no retroactive approval, production trust, release approval, or external attestation is claimed. |

## Process Improvement

The Kimi lane and an earlier external review incident showed that external
review commands can mutate or reset local work. `.agents/skills/sdp-trace-trust-workflow/SKILL.md`
now requires external review lanes to run only against a committed handoff or
isolated copy/worktree, and to record mutation/reset lanes as compromised
evidence rather than review authority.
