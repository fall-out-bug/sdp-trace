# Spec Reality Ledger

> Reconciles roadmap rows, spec files, task checkboxes, and live verification
> states for specs 001–019.
>
**Date**: 2026-05-26
**Commit**: `main` @ `657a343a5f310538def9afd509e6c610c713cab0`

> **Note**: This checked-in file records a post-merge reconciliation after PR
> #60. It is repository evidence, not merge approval evidence.
> **Source of truth**: `docs/roadmap.md`, `specs/*/spec.md`, `specs/*/plan.md`, `specs/*/tasks.md`

## Live Verification State (as of 2026-05-26)

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


## Ledger

| Spec | Title | Roadmap Status | Spec Status | Tasks (Checked / Total) | Key Open Tasks | Live Verification for Claimed Deliverables | Claims Backed by Current Evidence? |
|---|---|---|---|---|---|---|---|
| **001** | sdp-trace Time-Series Evidence Substrate | `draft` (blocked) | Draft | Most checked; **11 open** | T192, T209–T217, T226 | `go test` PASS, `go vet` PASS, `doccheck` PASS, `schemadoc` PASS. **CRAP PASS**. **MI baseline PASS**. **Absolute MI gap remains assessed**. | **Partial** — many blocks implemented and tested locally, but open tasks remain and spec status is draft. |

| **002** | Authority Envelope Boundary Observation | `draft` (old) | Draft — revised, re-review pending | 33 / 35 | T034, T035 | `go test` PASS, schema validation PASS. | **Partial** — implementation and fixtures complete locally; PR closure and merge not done. |
| **003** | Agent Supply Chain Roadmap | `draft` (old) | Draft — roadmap artifact, revisions pending | 5 / 42 | T006–T042 (all phases 1–8) | `go test` PASS, `doccheck` PASS. | **No** — remains a draft roadmap with no implemented deliverables beyond spec files. |
| **004** | MVP Readiness Hardening | `draft` (old) | Draft — revised, approval pending | 40 / 43 | T040, T041, T042 | `go test` PASS, `doccheck` PASS. **CRAP PASS**. **MI baseline PASS**. **Absolute MI gap remains assessed**. | **Partial** — docs and code hygiene work done, but MVP readiness claim remains blocked by open approval/tasks and absolute quality debt. |

| **005** | Product Contract v0 | `draft` (old) | Draft — re-review pending | 13 / 20 | T014–T020 | `go test` PASS, `doccheck` PASS. | **Partial** — contract drafted and linked to roadmap; approval and implementation tasks remain open. |
| **006** | Change Evidence Packet Core | `draft` (old) | Draft — needs Socratic review | 0 / 27 | T001–T027 | `go test` PASS, `doccheck` PASS. | **No** — spec only; no implementation or fixtures committed. |
| **007** | GitHub OSS Demo Packet | `draft` (old) | Draft — needs Socratic review | 5 / 22 | T006–T022 | `go test` PASS, `doccheck` PASS. | **No** — spec and plan drafted; no demo packets or implementation. |
| **008** | Invisible Flight Recorder | `in_progress` (blocked) | Implemented locally; PR/final-head CI evidence pending | 35 / 35 | — | `go test` PASS, `go vet` PASS, `doccheck` PASS. | **Partial** — all tasks checked and local tests pass; PR merge and CI evidence at final head are pending, blocking completion. |
| **009** | Machine-Readable Command Surface | `draft` | Draft for PI review | 14 / 14 | — | `go test` PASS, `doccheck` PASS. | **Partial** — implementation appears complete (tasks checked, tests pass), but spec status remains draft and no PI review approval is recorded. |

| **010** | Command Package Organization | `draft` | Draft for PI review | 15 / 15 | — | `go test` PASS, `doccheck` PASS. | **Partial** — reorganization implemented and verified locally; spec status remains draft, pending PI review/merge. |

| **011** | Schema Documentation Validation | `draft` | Draft for PI review | 14 / 14 | — | `go test` PASS, `doccheck` PASS, `schemadoc` PASS. | **Partial** — checker implemented and CI wired; spec status remains draft, pending PI review/merge. |

| **012** | Repo Hygiene And Artifact Boundary | `draft` | Draft for PI review | 12 / 12 | — | `go test` PASS, `doccheck` PASS, `hygienecheck` PASS. | **Partial** — hygiene check and cleanup implemented; spec status remains draft, pending PI review/merge. |

| **013** | Contributor Onboarding And Verification | `draft` | Draft for PI review | 0 / 12 | T001–T032 | `go test` PASS, `doccheck` PASS. | **No** — spec drafted; no onboarding path or smoke commands implemented. |

| **014** | Docs UX And Command Guidance | `draft` | Draft for PI review | 15 / 15 | — | `go test` PASS, `doccheck` PASS. | **Partial** — docs restructured and canonical contracts added; spec status remains draft, pending PI review/merge. |

| **015** | Spec Governance And Roadmap Navigation | `in_progress` | `in_progress` — implementation active on branch | 16 / 16 | — | `go test` PASS, `doccheck` PASS. | **Partial** — roadmap and taxonomy implemented; spec status `in_progress` because final multi-LLM review and merge are pending. |

| **016** | Production Adoption And Security Baseline | `in_progress` | `in_progress` | 10 / 10 | — | `go test` PASS, `go vet` PASS, security docs exist. | **Partial** — readiness matrix and security baseline created; `gosec` findings classified as advisory. PR #59 merged, but repeat review and external audit remain `not_assessed`. |

| **017** | OSS Replacement Compatibility And Benchmarks | `in_progress` | `in_progress` | 9 / 11 | T017-040, T017-080 | `go test` PASS (`tools/osscompat`, `tools/ossbench`). **CRAP PASS**. **MI PASS**. Supply-chain probes (in-toto, Cosign, SLSA) explicitly preserved as `cannot_verify` with reproducible reasons and test coverage. Live `wrap` manifest has `schema/run-manifest.schema.json`; optional external validation is `not_assessed` when `check-jsonschema` is absent. | **Partial** — compatibility harness and benchmark tooling exist and pass unit tests and quality gates. Supply-chain external evidence and roadmap/index closure remain open. |

| **018** | Core/Policy Split And Pi Delivery Plan | `draft` | `in_review` | 0 / 11 | T018-001–T018-070 | `go test` PASS, `doccheck` PASS. | **No** — planning spec in review; no implementation tasks completed. |
| **019** | Repo Realignment, Monitoring, And Gate Readiness | `in_progress` | post-merge partial | 14 / 16 | T019-004 (post-merge approval), T019-120 | `go test` PASS, `go vet` PASS, `doccheck` PASS, `hygienecheck` PASS, `schemadoc` PASS, **CRAP PASS**, **MI baseline PASS**, CI PASS on `main` at `657a343a5f310538def9afd509e6c610c713cab0`. Live `wrap` manifest schema added as `schema/run-manifest.schema.json`. Monitoring/gate proof pack is replayed by `cmd/sdp-trace/spec019_proof_pack_cli_test.go`. Scoped harness/gate numeric shard count moved from 463 to 435. Current-branch Qwen3.6 Plus alternative-LLM review is LGTM after findings were addressed. | **Partial** — PR #60 merged completed workstreams A/B/C/F/H and post-merge closure completed D/E/G. Post-merge approval and current-branch live CI remain open. GitHub PR metadata has no recorded review approval; `merge_approval` remains `not_assessed`. |

## Qualification Notes

### Quality Gate Status
- **CRAP**: All functions within threshold of 5 under strict-less mode. Quality gate passes.
- **MI Baseline**: No regressions. All changed production files maintain or improve MI. Quality gate passes.
- **Absolute MI**: Current absolute `-mi-under 70.1` and
  `-function-mi-under 70.1` checks still report historical files/functions
  below the absolute bar. This is an assessed quality gap, not a baseline
  regression.
- **PR #61**: Open and failing as of 2026-05-26. The branch predates PR #60 and
  is stale; it must be closed, rebased, or superseded before any merge.

### Remaining Blockers
- **Spec 017**: Supply-chain external evidence remains `cannot_verify` where optional tools or external trust anchors are unavailable.
- **Spec 019 T019-004**: Post-merge approval is `not_assessed`.
- **Spec 019 T019-120**: Final post-merge branch review is complete; live CI evidence remains open until branch/PR publication.
### Specs with Completed Tasks but Draft Status
Specs **009, 010, 011, 012, 014** have all tasks checked and local verification passes, yet their `spec.md` status remains `draft`. The discrepancy means implementation exists but has not been formally approved/merged. They must not be treated as `complete` until the spec status transitions and any pending PI review or merge is recorded.

### Blocked Specs
- **001**: Blocked on self-attestation proof (signed release process) and open tasks. Historical baseline quality failures were remediated; absolute MI remains an assessed gap.
- **008**: Blocked on PR/final-head CI evidence.
- **017**: Blocked on remaining roadmap/docs closure and external supply-chain evidence.

### No Completed Specs
No spec is in the `complete` status. The roadmap “Completed Specs” section correctly lists *(none yet)*.

---

*Generated as part of WS-019-A. For questions or updates, edit this file and rerun the verification commands listed in the live-verification table.*
