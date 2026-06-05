# Spec Reality Ledger

> Reconciles roadmap rows, spec files, task checkboxes, and live verification
> states for specs 001–022.
>
**Date**: 2026-06-05
**Commit**: `codex/022-post-merge-governance-closure` @
`9438080464f0103144790e0336560deeafc893b6`

> **Note**: This checked-in file records a Spec 022 governance reconciliation
> after PR #60 and PR #63. It is repository evidence, not merge approval
> evidence.
> **Source of truth**: `docs/roadmap.md`, `specs/*/spec.md`, `specs/*/plan.md`, `specs/*/tasks.md`

## Live Verification State (as of 2026-06-01)

| Command | Result | Notes |
|---|---|---|
| `go test -count=1 ./...` | **PASS** | All packages |
| `go vet ./...` | **PASS** | — |
| `go run ./tools/doccheck` | **PASS** | — |
| `go run ./tools/hygienecheck` | **PASS** | — |
| `go run ./tools/schemadoc` | **PASS** | — |
| `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less` | **PASS** | All functions within threshold |
| `go run ./tools/qualitycheck -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools` | **PASS** | No regressions |
| GitHub Actions CI | **PASS** | Latest observed `main` run 26417012453 passed on commit `657a343a5f310538def9afd509e6c610c713cab0` |
| PR #60 live refresh | **verified** | `gh pr view 60 --json number,state,mergeCommit,headRefOid,reviewDecision,statusCheckRollup,url` on 2026-06-01: `state=MERGED`, merge commit `657a343a5f310538def9afd509e6c610c713cab0`, head `977179ba93e577cc51f05e634453deba80c383b6`, `reviewDecision=""`, CI `verify` success at run `26416733118`. |
| PR #63 live refresh | **verified** | `gh pr view 63 --json number,state,mergeCommit,headRefOid,reviewDecision,statusCheckRollup,url` on 2026-06-01: `state=MERGED`, merge commit `1ee2c7af53637c7f43bff4e0e7ef9e34d164908e`, head `9ad4640828a204e66fe056dde0950542ed00f1c4`, `reviewDecision=""`, CI `verify` success at run `26428891169` and `pr-review-evidence-only` success at run `26428891168`. |


## Ledger

| Spec | Title | Roadmap Status | Spec Status | Tasks (Checked / Total) | Key Open Tasks | Live Verification for Claimed Deliverables | Claims Backed by Current Evidence? |
|---|---|---|---|---|---|---|---|
| **001** | sdp-trace Time-Series Evidence Substrate | `draft` (blocked) | Draft | 244 / 244 | — | `go test` PASS, `go vet` PASS, `doccheck` PASS, `schemadoc` PASS. **CRAP PASS**. **MI baseline PASS**. **Absolute MI gap remains assessed**. Block 21 export replay passed; PR #64 `pr-review-evidence-only` workflow evidence passed for Block 32 closure refresh. Block 25 Socratic review ledger exists. Active demo evidence is `fall-out-bug/sdp-trace-demo-jvm-gsd`: T211 is accepted with `kt_jvm_library`, Bazel CI, and semantic `sh_test` endpoint evidence over the compiled app; T212 is accepted as superseded by the v2 packet/bundle CI artifact contract; T213 is accepted with demo PR #25 merged to `main` as `3a9491f734e5214c72014db5d893f125eb254a11`, local Bazel verification, and downloaded-artifact replay over demo CI run `25724386343`; T214 is accepted with demo PR #26 merged to `main` as `a4d1f755552ba1f411af5edcb7d6caf24a9c39bf` and GitHub `build-and-test` run `26447797437`; T215 is accepted with `docs/reviews/block25-jvm-gsd-demo-sanitized-report.md`; T216 is accepted with active-demo role reviews and fixed redaction-scan evidence; T217 is accepted with fresh local verification plus PR #64 `verify` and `pr-review-evidence-only` success at head `537612723f7dad7f6e3a92657336fafa4ce6bbd5` and post-merge `main` CI run `26453881873`. T226 GSD-Redux replay observes route `minimax/MiniMax-M2.5`, setup metadata, command digest, process id, source commit, time bounds, output/normalized digests, and model/interaction/tool/phase evidence; mutation/test remain explicit `not_assessed` for the replayed no-op completed phase. | **Implemented local / trust-bounded** — tasks are reconciled, but spec status is draft and no feature-delivery, test-success, harness-compliance, production trust, or broad GSD/OpenCode support is claimed. |

| **002** | Authority Envelope Boundary Observation | `draft` (old) | merged / post-merge verified | 35 / 35 | — | `go test` PASS, schema validation PASS. Authority fixture replay produced `outside_authority` with downstream-policy next action. PR #64 checks passed at head `537612723f7dad7f6e3a92657336fafa4ce6bbd5`; explicit merge approval was recorded in the 2026-05-26 closure thread; PR #64 merged as `e129515e7c4c7a4a9c4b2b53eb4d3694b41eb2bd`; post-merge `main` CI run `26453881873` passed. | **Implemented and merged** — implementation, fixtures, docs, PR-level closure review, merge approval, merge, and post-merge CI evidence exist for the scoped boundary-observation claim. |
| **003** | Agent Supply Chain Roadmap | `retired_superseded` | Retired roadmap artifact | 42 / 42 | none in task ledger | `doccheck` PASS. Retirement review maps useful direction to later concrete specs and keeps unsupported surfaces out of product claims. | **No product implementation claim** — retired as stale planning, not completed as product. |
| **004** | MVP Readiness Hardening | `draft` (old) | merged / post-merge verified | 43 / 43 | — | `go test` PASS, `doccheck` PASS. **CRAP PASS**. **MI baseline PASS**. **Absolute MI gap remains assessed**. PR-level review and named reviewer sign-off exist for the controlled-pilot scope. Explicit merge approval was recorded in the 2026-05-26 closure thread; PR #64 merged as `e129515e7c4c7a4a9c4b2b53eb4d3694b41eb2bd`; post-merge `main` CI run `26453881873` passed. | **Implemented and merged for controlled pilot** — named sign-off, merge approval, merge, and post-merge CI evidence exist; production readiness and absolute MI pass are not claimed. |

| **005** | Product Contract v0 | `draft` (old) | reviewed contract / accepted gap | 20 / 20 | — | `go test` PASS, `doccheck` PASS. Spec 006 provides packet schema/model, renderer, fixtures, theater tests, and CLI surface. Maintainer decision `accept_gap` was recorded in the 2026-05-26 closure thread for the historical approval gap. | **Reviewed contract / accepted gap** — contract drafted and implementation placeholders are covered by Spec 006 artifacts; the historical approval gap is accepted without claiming production approval, release approval, or external trust. |
| **006** | Change Evidence Packet Core | `draft` (old) | implemented / accepted gap | 27 / 27 | — | `go test` PASS, `doccheck` PASS. Packet schema/model, renderer, CLI validation/render behavior, fixtures, theater tests, closure review, PR #64 evidence, merge evidence, and post-merge CI evidence exist. Maintainer decision `accept_gap` was recorded in the 2026-05-26 closure thread for the historical pre-implementation approval gap. | **Implemented / accepted gap** — implementation and review evidence are mapped; the historical pre-implementation approval gap is accepted without claiming production readiness, release approval, or external trust. |
| **007** | GitHub OSS Demo Packet | `draft` (old) | external demo track approved | 22 / 22 | — | `go test` PASS, `doccheck` PASS. Historical review ledger includes 2026-05-26 focused split re-review. Maintainer decision `approve_existing_demo_track` was recorded in the 2026-05-26 closure thread. External demo repo `fall-out-bug/sdp-trace-demo-jvm-gsd` was cloned at `b101882`; PRs #16-#23 were queried; main and negative PR #21 branch Bazel smoke tests passed; feature bundles 1-5 passed `packet validate` and `packet check-demo`; external tag `demo-v1-observation-baseline` preserves PR #16 base commit; buyer rehearsal artifact covers happy-path PR #20 and negative PR #21. | **Scoped demo track approved** — setup, feature packets, negative draft PR, CI, first-packet gate, v1 archival tag, buyer rehearsal, and existing-demo-track decision have current evidence; release, production trust, compliance, semantic quality, and signed external trust are not claimed. |
| **008** | Invisible Flight Recorder | `in_progress` (blocked) | Implemented locally; PR/final-head CI evidence pending | 35 / 35 | — | `go test` PASS, `go vet` PASS, `doccheck` PASS. | **Partial** — all tasks checked and local tests pass; PR merge and CI evidence at final head are pending, blocking completion. |
| **009** | Machine-Readable Command Surface | `draft` | Draft for PI review | 14 / 14 | — | `go test` PASS, `doccheck` PASS. | **Partial** — implementation appears complete (tasks checked, tests pass), but spec status remains draft and no PI review approval is recorded. |

| **010** | Command Package Organization | `draft` | Draft for PI review | 15 / 15 | — | `go test` PASS, `doccheck` PASS. | **Partial** — reorganization implemented and verified locally; spec status remains draft, pending PI review/merge. |

| **011** | Schema Documentation Validation | `draft` | Draft for PI review | 14 / 14 | — | `go test` PASS, `doccheck` PASS, `schemadoc` PASS. | **Partial** — checker implemented and CI wired; spec status remains draft, pending PI review/merge. |

| **012** | Repo Hygiene And Artifact Boundary | `draft` | Draft for PI review | 12 / 12 | — | `go test` PASS, `doccheck` PASS, `hygienecheck` PASS. | **Partial** — hygiene check and cleanup implemented; spec status remains draft, pending PI review/merge. |

| **013** | Contributor Onboarding And Verification | `draft` | implemented-local | 12 / 12 | none in task ledger | `go test` PASS, `doccheck` PASS, `hygienecheck` PASS, smoke path replay PASS. | **Implemented local / review-bounded** — contributor quickstart, README/docs-map links, expected results, failure routing, doccheck coverage, smoke replay, and closure review exist. Formal maintainer closure remains separate. |

| **014** | Docs UX And Command Guidance | `draft` | Draft for PI review | 15 / 15 | — | `go test` PASS, `doccheck` PASS. | **Partial** — docs restructured and canonical contracts added; spec status remains draft, pending PI review/merge. |

| **015** | Spec Governance And Roadmap Navigation | `in_progress` | `in_progress` — implementation active on branch | 16 / 16 | — | `go test` PASS, `doccheck` PASS. | **Partial** — roadmap and taxonomy implemented; spec status `in_progress` because final multi-LLM review and merge are pending. |

| **016** | Production Adoption And Security Baseline | `in_progress` | `in_progress` | 10 / 10 | — | `go test` PASS, `go vet` PASS, security docs exist. | **Partial** — readiness matrix and security baseline created; `gosec` findings classified as advisory. PR #59 merged, but repeat review and external audit remain `not_assessed`. |

| **017** | OSS Replacement Compatibility And Benchmarks | `in_progress` | `in_progress` | 11 / 11 | none in task ledger | `go test` PASS (`tools/osscompat`, `tools/ossbench`). **CRAP PASS**. **MI PASS**. Supply-chain probes (in-toto, Cosign, SLSA) explicitly preserve optional-tool absence as `not_assessed` and conformance as `cannot_verify`. Live `wrap` manifest has `schema/run-manifest.schema.json`; optional external validation is `not_assessed` when `check-jsonschema` is absent. | **Implemented local / trust-bounded** — compatibility harness, benchmark tooling, controlled supply-chain probes, and docs index updates exist. External supply-chain trust remains out of scope. |

| **018** | Core/Policy Split And Pi Delivery Plan | `draft` | reviewed / follow-ups prepared | 11 / 11 | — | `go test` PASS, `doccheck` PASS, `hygienecheck` PASS. Maintainer decision `approve_core_extension_direction` was recorded in the 2026-05-26 closure thread. Follow-up Spec 020 and Spec 021 are filed as draft implementation specs. | **Reviewed direction / follow-ups prepared** — workstreams A-E, integration verification, maintainer review, and follow-up spec filing are complete; command removal, separate binaries, production readiness, release approval, and external trust are not claimed. |
| **020** | Core Query Package Split | `draft` | follow-up prepared | 0 / 0 active checkboxes | — | SpecKit triplet exists. Future task backlog is present as non-active table rows. | **Prepared future implementation spec** — not current closure debt until explicitly taken into work. |
| **021** | Source File Locality Cleanup | `draft` | follow-up prepared | 0 / 0 active checkboxes | — | SpecKit triplet exists. Future task backlog is present as non-active table rows. | **Prepared future implementation spec** — not current closure debt until explicitly taken into work. |
| **019** | Repo Realignment, Monitoring, And Gate Readiness | `in_progress` | split successor | 16 / 16 | — | `go test` PASS, `go vet` PASS, `doccheck` PASS, `hygienecheck` PASS, `schemadoc` PASS, **CRAP PASS**, **MI baseline PASS**, CI PASS on `main` at `657a343a5f310538def9afd509e6c610c713cab0`. PR #60 live refresh on 2026-06-01 confirmed `state=MERGED`, merge commit `657a343a5f310538def9afd509e6c610c713cab0`, CI `verify` success, and empty `reviewDecision`. Integration PR #63 live refresh confirmed `state=MERGED`, merge commit `1ee2c7af53637c7f43bff4e0e7ef9e34d164908e`, CI `verify` success, and `pr-review-evidence-only` success. Live `wrap` manifest schema added as `schema/run-manifest.schema.json`. Monitoring/gate proof pack is replayed by `cmd/sdp-trace/spec019_proof_pack_cli_test.go`. Scoped harness/gate numeric shard count moved from 463 to 435. Maintainer decision `split_successor` moves residual governance debt to Spec 022. | **Split successor** — PR #60 merged completed workstreams A/B/C/F/H and post-merge closure completed D/E/G. PR #63 added PR-review CI enforcement and pinned workflow action checks. Missed Phase 0 gates remain historical process evidence; PR #60 merge approval remains `not_assessed`; Spec 022 owns the residual governance closure. |
| **022** | Post-Merge Governance Closure | `in_progress` | PR-ready | 42 / 42 | — | SpecKit triplet exists and was activated on branch `codex/022-post-merge-governance-closure`. Plan/task review Round 1 found critical/major issues; fixes were committed in `fc8a1e1` and `be2b19e`; Round 2 recorded no retained critical or major findings in `specs/022-post-merge-governance-closure/reviews/plan-task-review-round-2.md`. US1 evidence summary is committed in `docs/spec-reality-ledger.md`; US2 decision state is committed in `docs/closure-decision-ledger.md` and `docs/open-task-breakdown.md`; US3 navigation sync is committed in `docs/roadmap.md`, `docs/spec-closure-route.md`, and closure surfaces; final local verification is recorded in `specs/022-post-merge-governance-closure/reviews/final-evidence.md` and was replayed after rebasing onto current `main`; fixed-handoff PR-ready re-review returned LGTM across DeepSeek, Kimi, and Z.AI fallback lanes; live PR #60 and PR #63 refresh is verified as recorded above. | **PR-ready governance closure** — evidence summary, decision state, navigation sync, local verification, and fixed-handoff adversarial review are recorded on PR #72. The closure must not claim retroactive PR #60 approval, production trust, release approval, or external attestation. PR #60 merge approval remains `not_assessed`; final-head PR #72 CI/merge evidence is required before closure. |

## Qualification Notes

### Quality Gate Status
- **CRAP**: All functions within threshold of 5 under strict-less mode. Quality gate passes.
- **MI Baseline**: No regressions. All changed production files maintain or improve MI. Quality gate passes.
- **Absolute MI**: Current absolute `-mi-under 70.1` and
  `-function-mi-under 70.1` checks still report historical files/functions
  below the absolute bar. This is an assessed quality gap, not a baseline
  regression.
- **PR #61**: Closed as superseded by PR #62. The branch predated PR #60 and
  was not the correct path for post-merge baseline or quality closure.
- **PR #62**: Superseded by merged integration PR #63. Its CI pass remains
  historical evidence for that head, not the current handoff surface.
- **PR #31**: Superseded by merged integration PR #63. Its live checks are
  historical evidence, but the branch is no longer the repo-wide completion
  surface because PR #63 includes the CI PR-review integration changes.

### Remaining Blockers
- **Spec 017**: Supply-chain external evidence remains `cannot_verify` where optional tools or external trust anchors are unavailable.
- **Spec 019 governance debt**: Pre-implementation review and approval gates
  were missed before PR #60 merged; they remain process evidence and cannot be
  retroactively marked complete. Residual governance is actively handled by
  Spec 022. PR #60 merge approval remains `not_assessed` after the 2026-06-01
  live refresh because GitHub `reviewDecision` is empty.
- **Demo / first-run specs**: 001 Block 25/T226 and 007 buyer rehearsal are
  represented by real external demo-repository or first-run observation
  evidence, not local ledger cleanup. Future non-no-op GSD runs still need
  separate mutation/test evidence before those dimensions can move beyond
  `not_assessed`.
### Specs with Completed Tasks but Draft Status
Specs **009, 010, 011, 012, 014** have all tasks checked and local verification passes, yet their `spec.md` status remains `draft`. The discrepancy means implementation exists but has not been formally approved/merged. They must not be treated as `complete` until the spec status transitions and any pending PI review or merge is recorded.

### Blocked Specs
- **001**: Blocked on open demo/first-run tasks. Historical baseline quality failures were remediated; absolute MI remains an assessed gap.
- **008**: Blocked on PR/final-head CI evidence.
- **017**: Task ledger closed as implemented-local; external supply-chain trust remains out of scope and requires new evidence before any stronger claim.

### No Completed Specs
No spec is in the `complete` status. The roadmap “Completed Specs” section correctly lists *(none yet)*.

---

*Generated as part of WS-019-A. For questions or updates, edit this file and rerun the verification commands listed in the live-verification table.*
