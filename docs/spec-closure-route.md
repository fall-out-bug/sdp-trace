# Spec Closure Route

Status: working audit, current as of 2026-05-26.

This document audits specs 001-019 against the current repository state and
defines the route to close, defer, or repair each spec. It is a planning and
coordination artifact, not proof that any spec is complete.

## Evidence Snapshot

| Evidence | Current State |
| --- | --- |
| Source commit inspected | PR #64 head `76fc3ac5ee61bc3b6c0c2554c53c67f14a9520b7` plus local closure-route reconciliation updates |
| Open GitHub PRs | PR #64 open, not draft, merge state `CLEAN`; final observed head checks passed before these local updates |
| Latest `main` CI | pass, run `26436136055` |
| Spec directories | 19 |
| SpecKit triplets | 19 / 19 |
| Task checkboxes | 571 / 605 checked |
| Formal roadmap `complete` rows | 0 |

Commands used for this audit:

- `find specs -maxdepth 2 -type f \( -name spec.md -o -name plan.md -o -name tasks.md \)`
- `rg -n '^- \[ \]' specs/*/tasks.md`
- `go run ./cmd/sdp-trace command-surface`
- `gh pr list --state open --limit 50`
- `gh run list --branch main --limit 5`

## Closure Semantics

Use the axes in `docs/spec-status-discipline.md`.

Do not close a spec because task boxes are checked. A closeable spec needs:

- task ledger reconciled with current files;
- implementation evidence for the claimed scope;
- review disposition or explicit `not_assessed`;
- merge state represented against current `main`;
- trust claims backed by current evidence or marked out of scope.

## Audit Table

| Spec | Current Reality | Closure Route |
| --- | --- | --- |
| 001 Time-Series Evidence Substrate | Broad substrate is substantially implemented through historical blocks; `tasks.md` is 236 / 244 checked and `blocks/` has 93 files. Not closed: T211-T217 and T226 remain open. | Split into closed historical blocks vs live blockers. Block 21 and Block 32 closure refreshes are recorded; Block 25 demo-repository evidence and T226 remain real demo / first-run observation work. Do not try to close 001 wholesale. |
| 002 Authority Envelope Boundary Observation | Implementation appears present: authority schemas, docs, fixtures, Go parser/evaluator, and tests are checked in; tasks are 34 / 35. PR-level closure review exists for PR #64. Missing merge/post-merge closure remains T035. | Keep T035 open until merge, fresh CI, local verification, PR review, and post-merge verification are all represented. |
| 003 Agent Supply Chain Roadmap | Roadmap/spec discovery artifact only; tasks are 42 / 42 after retirement. It was not implemented as product. | Retired as `retired_superseded` by concrete later specs and blocks. Do not use Spec 003 as implementation authority. |
| 004 MVP Readiness Hardening | Most implementation work is checked off: 42 / 43. PR-level review and named reviewer sign-off are recorded for PR #64. Explicit merge approval remains open as T042; absolute MI remains an assessed gap, not a pass claim. | Keep T042 open until explicit merge approval exists. Do not upgrade controlled-pilot readiness to production readiness or absolute MI pass. |
| 005 Product Contract v0 | Contract design and reviews exist; tasks are 19 / 20 after mapping implementation placeholders to Spec 006 packet artifacts. Historical explicit approval remains `not_assessed`. | Do not claim Product Contract approval until maintainers accept, waive, or explicitly preserve the approval gap. |
| 006 Change Evidence Packet Core | Task ledger now maps 26 / 27 implemented, locally verified, and reviewed tasks to packet schemas, docs, Go package `internal/packet`, CLI surfaces, fixtures, tests, this branch's verification commands, closure review, and PR #64 evidence. Historical pre-implementation approval remains open / `not_assessed`. | Do not close as complete until maintainers decide whether the historical approval gap is accepted, waived, or kept as a permanent trace boundary. |
| 007 GitHub OSS Demo Packet | Draft/demo plan plus focused split re-review; tasks are 7 / 22. Phase 1 review is closed, but explicit demo-track approval and Phase 2 demo-repository packetization work remain open. | Keep open. Do not close as implemented from this repository. The local demo repo is not present in this environment, so T009-T022 remain `not_assessed` until real demo-repository evidence exists. |
| 008 Invisible Flight Recorder | Tasks are 26 / 26 and implementation/review artifacts exist. Roadmap still says PR/final-head CI evidence pending. | Verify whether the PR/merge happened historically. If already merged, update closure state from current `main`; if not, run final-head CI/review equivalent and close as implemented-local or historical. |
| 009 Machine-Readable Command Surface | Tasks are 14 / 14 and command-surface output works. Review directory absent, but tasks claim PI review and implementation review happened. | Reconstruct missing review evidence or mark review evidence `not_assessed`. Run command-surface drift checks against docs and close as implemented if review gap is accepted. |
| 010 Command Package Organization | Tasks are 14 / 14; review artifacts exist; current main contains the command-family split. | Run current behavior snapshot checks and close. This is likely closeable with no product code if review artifacts are sufficient. |
| 011 Schema Docs Generation | Tasks are 14 / 14; `tools/schemadoc`, schema index, README verification, and CI hooks exist. | Run `go run ./tools/schemadoc` and `go run ./tools/schemadoc -verify-readme`; if pass, close with local/CI evidence. |
| 012 Repo Hygiene And Artifact Boundary | Tasks are 12 / 12; `tools/hygienecheck` exists and CI uses hygiene checks. | Run hygienecheck and close if current main passes. Record any ignored local clutter as local-only, not repo proof. |
| 013 Contributor Onboarding And Verification | Tasks are 12 / 12 after reconciliation. Contributor quickstart, README/docs-map links, expected results, failure routing, doccheck coverage, smoke replay, and closure review exist. | Treat as implemented-local with local onboarding evidence. Formal maintainer closure remains separate from task-box completion. |
| 014 Docs UX And Command Guidance | Tasks are 15 / 15; docs such as reviewer entrypoint, output map, profile guide, overclaim checklist exist. | Run doccheck plus a focused docs UX/review refresh. Likely closeable after current evidence is recorded. |
| 015 Spec Governance And Roadmap | Tasks are 17 / 17 after the status-discipline update; roadmap and status discipline exist. | Close after this audit route is linked and verified. The remaining risk is stale roadmap facts after future merges; keep freshness rule active. |
| 016 Production Adoption Security Baseline | Tasks are 10 / 10; docs exist; PR #59 is now merged per roadmap. External audit, customer adoption, signed release, and production trust remain `not_assessed`. | Close only as controlled-pilot/security-baseline docs. Do not close production adoption. Update wording to separate local security baseline from external production trust. |
| 017 OSS Replacement Compatibility And Benchmarks | Tasks are 11 / 11. Tooling, benchmark harnesses, controlled supply-chain probes, docs index, and roadmap updates exist. Optional external tools remain `not_assessed` when absent; supply-chain conformance remains `cannot_verify`. | Treat as implemented-local with explicit optional-tool boundaries. Do not upgrade to external supply-chain trust without new evidence. |
| 018 Core/Policy Split And Pi Delivery | Task ledger now maps 9 / 11 implemented and locally verified tasks to command stability, package ownership, extension boundary, source locality, core-first docs, and integration verification. Maintainer review and follow-up implementation specs remain open / `not_assessed`. | Keep maintainer review open as `not_assessed` until explicit approval exists; prepare follow-up specs only after that direction is accepted. |
| 019 Repo Realignment, Monitoring, And Gate Readiness | Tasks are 11 / 16. Major workstreams are implemented, PR #63 merged after final-head CI passed, and Phase 0 approval gates plus T019-120 remain open. | Treat as post-merge governance debt. Maintainers must decide whether to accept the already-merged work, reject it, or split remaining governance closure into a successor spec. |

## Open Task Classification

Current task ledger state is 571 / 605 checked, leaving 34 open task boxes.

| Category | Specs / tasks | Closure meaning |
| --- | --- | --- |
| Approval or maintainer gates | 002 T035; 004 T042; 005 T014; 006 T003; 018 T018-001/T018-070; 019 T019-001/T019-002/T019-003/T019-004/T019-120 | Cannot be converted to `pass` locally. Needs explicit maintainer acceptance, waiver, rejection, merge/post-merge evidence, or successor-spec split. |
| External demo / first-run work | 001 T211-T217, T226; 007 T009-T022 | Real remaining work outside this repo's current local artifacts. Needs demo-repository execution, retained evidence, and review. |
| Review / PR evidence closure | none currently isolated as review-only after this refresh | Future review gaps must still be closed with fresh review, live PR evidence, and explicit trust-boundary wording, not task-box cleanup alone. |
| Retire or re-scope stale planning | 007 if demo direction is superseded | Needs a maintainer decision: preserve as draft, retire as superseded, or create a smaller successor spec. |

## Recommended Route

1. **Finish verification for repaired ledgers**: specs 006 and 018 now map
   existing implementation artifacts to their task files, but fresh local
   verification and review evidence still need to be recorded before closure.
2. **Close low-risk implemented specs**: 010, 011, 012, 014, 015 can likely be
   closed with docs/tool verification and review-evidence reconstruction.
3. **Close trust-sensitive implemented specs**: 002, 008, 009, 016 need focused
   review and explicit trust-boundary wording before closure.
4. **Resolve partial/blocked specs**: 004, 005, 017, 019 need targeted closure
   work or explicit deferral.
5. **Retire or re-scope stale planning specs**: 003 and 007 should not remain
   open forever if their product direction has been superseded.
6. **Do not close 001 as one giant spec**: split it into historical completed
   block evidence plus a small successor list for the remaining open tasks.

## Immediate Next Slice

Continue with specs 006 and 018:

- run the relevant local verification commands;
- record review gaps as `not_assessed`, not as implicit approval;
- update closure state only after verification and review evidence are current.

The stale-ledger contradiction has been reduced. The remaining risk is more
specific: implemented-local work may still be missing approval, review, or
fresh PR-level evidence.

<!-- sdp-trace-claim: claim=profile_passed; subject=spec-closure-route-audit; state=pass; profile=repo_baseline_structural; evidence=state:claim_tags_consistent -->
