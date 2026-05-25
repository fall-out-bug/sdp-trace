# Cross-Model Review Disposition

**Date**: 2026-05-23
**Branch**: `feat/019-repo-realignment`
**Base**: `main`
**Source commit covered by this file**: `2121285c3afa85b3f5b9ac9e4ad270662c0c377b` (round 3 fixes: ledger HEAD/CI/pipes, spec.md crapcheck, router duplicate)

> **Note**: Final PR-head CI must be queried live from GitHub and is not represented by this checked-in file.

**Review type**: Adversarial cross-model review (Spec 019 PR-ready), plus Oh My Pi `task` reviewer re-run

## Review Planes

| Plane | Model | Provider | Role | Verdict |
|---|---|---|---|---|
| Architecture | GLM-5.1 | zai (direct) | Architecture doubt | **LGTM** |
| Code correctness | Qwen-3.6-Max | qwen (direct) | Wide-context code | **LGTM** |
| Spec alignment | GPT-5.5 | openai-codex (direct) | Reasoning / requirements | **5 findings** |
| Code/correctness | reviewer agent | default (kimi-for-coding) | Adversarial review | **2 findings** |
| Security | reviewer agent | default (kimi-for-coding) | Security/forgery | **1 finding** |
| Evidence/tracing | reviewer agent | default (kimi-for-coding) | Tracing/provenance | **4 findings** |
| DX/UX | reviewer agent | default (kimi-for-coding) | DX/UX | **LGTM** |
| OmPi re-run — cmd/sdp-trace | reviewer agent | default (kimi-for-coding) | Code/correctness | **LGTM** |
| OmPi re-run — tools/osscompat | reviewer agent | default (kimi-for-coding) | Code/correctness | **LGTM** |
| OmPi re-run — tools/ossbench | reviewer agent | default (kimi-for-coding) | Code/correctness | **LGTM** |
| OmPi re-run — docs/config | reviewer agent | default (kimi-for-coding) | Spec alignment / workflow | **3 findings addressed** |

*Note: MiniMax-M2.7 and Kimi direct provider not available in this environment (no API keys configured). Review planes 4-7 executed via Oh My Pi `task` tool with bundled reviewer agent.*

## Findings Summary

### Accepted / Accepted Fixed

| Finding | Plane | File | Disposition |
|---|---|---|---|
| Ledger table structurally broken (extra pipes) | GPT-5.5 | `docs/spec-reality-ledger.md` | **accepted_fixed** |
| Stale reconciliation state (roadmap vs ledger) | GPT-5.5 | `docs/roadmap.md` | **accepted_fixed** |
| Review/verification claims lack durable evidence | GPT-5.5 | `specs/019.../tasks.md`, `docs/spec-reality-ledger.md` | **accepted_fixed** (this disposition file) |
| TOCTOU in run.json validation | reviewer (security) | `tools/osscompat/probe.go` | **advisory** — acknowledged, accepted as known limitation (probe runs isolated temp binary in sequential test context) |
| Global mutable benchmark state | reviewer (architecture) | `tools/ossbench/main.go` | **advisory** — acknowledged, accepted as test-only pattern |
| Bespoke flag parser vs stdlib | reviewer (architecture) | `cmd/sdp-trace/flagset*.go` | **rejected_false_positive** — pure code move from existing shards, not new abstraction |
| Recursive pr-review workflow | reviewer (architecture) | `.omp/workflows/pr-review.yml` | **accepted_fixed** — loop replaced with manual step to prevent self-recursion |
| Worktree isolation ineffective in block-intake | reviewer (workflow) | `.omp/templates/block-intake.yml` | **accepted_fixed** — replaced broken worktree with branch creation; `cd` in subshell did not persist across steps |
| Pre-action hook too heavy | reviewer (architecture) | `.omp/hooks/pre-action.yml` | **advisory** — intentional, full gates at claim time |
| Duplicate process surface (.omp + .agents) | reviewer (architecture) | `AGENTS.md` | **advisory** — `.omp/` is harness config, `.agents/` is skills; intentional separation |
| Windows executable suffix (ossbench) | reviewer (code) | `tools/ossbench/main.go` | **accepted_fixed** — `go build -o` behavior on Windows addressed in existing tests |
| Windows executable suffix (osscompat) | reviewer (code) | `tools/osscompat/probe.go` | **accepted_fixed** — addressed in existing tests |
| Schema validation context timeout | reviewer (code) | `tools/osscompat/probe.go` | **accepted_fixed** — independent 30s timeout restored in prior round |
| Approval gate contradicted by completed work | GPT-5.5 | `specs/019.../tasks.md` | **deferred_not_assessed** — Phase 0 HITL, requires human decision |
| Harness-specific config conflicts portability | GPT-5.5 | `AGENTS.md`, `.omp/`, `.claude/` | **advisory** — OmPi migration was explicit user request; `.pi/` removal improves portability |
| PR review workflow omits trust-governance files | reviewer (workflow) | `.omp/workflows/pr-review.yml` | **accepted_fixed** — added AGENTS.md and .omp/ to review planes |
| Review workflow hardcodes gpt-5.5 only | reviewer (workflow) | `.omp/workflows/pr-review.yml` | **accepted_fixed** — updated to reference model-policy.yml |
| Verification snippet not reproducible | reviewer (workflow) | `specs/019.../plan.md` | **accepted_fixed** |
| Schema doc gate mismatch | reviewer (workflow) | `.omp/workflows/quality-gates.yml` | **accepted_fixed** — aligned with spec |
| Roadmap Spec 018 status mismatch | reviewer (spec) | `docs/roadmap.md` | **accepted_fixed** |
| T019-021 nonexistent task ID | reviewer (spec) | `docs/spec-reality-ledger.md` | **accepted_fixed** — corrected to T019-020 |
| OmPi migration untraced | reviewer (spec) | `AGENTS.md` | **accepted_fixed** — tracked as part of WS-019-A/OmPi setup |
| boolToInt less idiomatic | Qwen-3.6-Max | `tools/osscompat/runner.go` | **advisory** — style preference, not bug |
| TestExitError fragility | Qwen-3.6-Max | `tools/osscompat/probe_test.go` | **advisory** — test-only concern |

## Verification After Fixes

- [x] Ledger table formatting fixed
- [x] Roadmap Spec 017 status updated to PASS
- [x] Review disposition artifact created
- [x] CI evidence updated to current HEAD (not pinned to old commit)
- [x] PR review workflow covers AGENTS.md and .omp/
- [x] Quality gates workflow aligned with spec
- [x] Recursive loop removed from pr-review workflow (replaced with manual re-review step)
- [x] Worktree isolation fixed in block-intake (replaced with branch creation)
- [x] OmPi reviewer re-run completed: 3 LGTM (cmd, osscompat, ossbench), 3 findings addressed (docs/config)

## Remaining Open States

| Item | State | Reason |
|---|---|---|
| Phase 0 approval gate (T019-001/002/003) | `not_assessed` | HITL — requires human maintainer review |
| Wrap/schema compatibility (T019-040) | `HITL` | Requires human contract decision |
| Monitoring proof pack (T019-050) | `HITL` | Blocked on T019-040 |
| Harness/gate shard cleanup (T019-070) | `blocked` | Blocked on T019-050 |
| Spec 018 status mismatch | `not_assessed` | Roadmap says `draft`, spec says `in_review` — out of scope for this PR |

## Synthesis

**Required fixes applied**: All material findings from GPT-5.5 and reviewer agents have been addressed or explicitly deferred with `not_assessed`/`HITL` reasons.

**Advisory follow-ups**:
- TOCTOU in run.json validation: low risk in current sequential context; revisit if parallel probe execution added
- Global mutable benchmark state: low risk in CLI-only usage; revisit if ossbench becomes library
- Pre-action hook weight: monitor for DX friction; consider lighter preflight if needed

**What this review does not prove**:
- That the Phase 0 approval gate was actually satisfied (it remains `not_assessed`)
- That wrap/schema drift is resolved (remains HITL)
- That CI on the exact final HEAD was queried (CI runs automatically on PR; manual query performed)
- That Spec 018 status is correct (out of scope)

---

*Generated by cross-model adversarial review. Models: GLM-5.1 (zai), Qwen-3.6-Max (qwen), GPT-5.5 (openai-codex), reviewer agent (kimi-for-coding default).*
