# Spec Closure Route

Status: working audit, current as of 2026-06-01.

This document audits specs 001-022 against the current repository state and
defines the route to close, defer, or repair each spec. It is a planning and
coordination artifact, not proof that any spec is complete.

## Evidence Snapshot

| Evidence | Current State |
| --- | --- |
| Source commit inspected | `codex/022-post-merge-governance-closure` at `9438080464f0103144790e0336560deeafc893b6` |
| PR #64 | merged on 2026-05-26 after explicit maintainer approval; final PR checks passed at head `537612723f7dad7f6e3a92657336fafa4ce6bbd5` |
| Latest `main` CI | pass, run `26453881873` at `e129515e7c4c7a4a9c4b2b53eb4d3694b41eb2bd` |
| Spec directories | 22 |
| SpecKit triplets | 22 / 22 |
| Task checkboxes | 647 / 647 checked |
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
| 001 Time-Series Evidence Substrate | Broad substrate is substantially implemented through historical blocks; `tasks.md` is 244 / 244 checked and `blocks/` has 93 files. Active demo evidence is in `fall-out-bug/sdp-trace-demo-jvm-gsd`: T211 is accepted with `sh_test` semantic behavior evidence over the compiled app, T212 is accepted as superseded by the v2 packet/bundle CI artifact contract, T213 is accepted with demo PR #25 merged to `main` plus downloaded-artifact replay over demo CI run `25724386343`, T214 is accepted with demo PR #26 merged to `main` plus GitHub `build-and-test` run `26447797437`, T215 is accepted with `docs/reviews/block25-jvm-gsd-demo-sanitized-report.md`, T216 is accepted with active-demo role reviews plus fixed redaction-scan evidence, T217 is accepted with fresh local verification plus PR #64 `verify` and `pr-review-evidence-only` success at head `537612723f7dad7f6e3a92657336fafa4ce6bbd5` and post-merge `main` CI run `26453881873`, and T226 is accepted with current OpenCode/GSD-Redux first-run observation evidence. T226 observes `minimax/MiniMax-M2.5`, setup metadata, command digest, process id, source commit, time bounds, output digest, normalized digest, and model/interaction/tool/phase families; mutation and test remain explicit `not_assessed` because the target phase was already complete and execute-phase performed no file-mutation or test action. | Close Spec 001 as historical/demo evidence complete for the scoped pilot path. Do not claim feature delivery, harness compliance, test success, production trust, or broad GSD/OpenCode support from T226. |
| 002 Authority Envelope Boundary Observation | Implementation appears present: authority schemas, docs, fixtures, Go parser/evaluator, and tests are checked in; tasks are 35 / 35. PR-level closure review exists for PR #64, explicit merge approval was recorded in the 2026-05-26 closure thread, PR #64 merged as `e129515e7c4c7a4a9c4b2b53eb4d3694b41eb2bd`, and post-merge `main` CI run `26453881873` passed. | Close Spec 002 as implemented and post-merge verified for its scoped authority-envelope boundary observation. Do not claim external policy approval or production trust. |
| 003 Agent Supply Chain Roadmap | Roadmap/spec discovery artifact only; tasks are 42 / 42 after retirement. It was not implemented as product. | Retired as `retired_superseded` by concrete later specs and blocks. Do not use Spec 003 as implementation authority. |
| 004 MVP Readiness Hardening | Implementation work is checked off: 43 / 43. PR-level review and named reviewer sign-off are recorded for PR #64; explicit merge approval was recorded in the 2026-05-26 closure thread; PR #64 merged as `e129515e7c4c7a4a9c4b2b53eb4d3694b41eb2bd`; post-merge `main` CI run `26453881873` passed. Absolute MI remains an assessed gap, not a pass claim. | Close Spec 004 as controlled-pilot MVP hardening merged and post-merge verified. Do not upgrade controlled-pilot readiness to production readiness or absolute MI pass. |
| 005 Product Contract v0 | Contract design and reviews exist; tasks are 20 / 20 after mapping implementation placeholders to Spec 006 packet artifacts. Maintainer decision `accept_gap` was recorded in the 2026-05-26 closure thread for the historical approval gap. | Close Spec 005 as reviewed contract with accepted historical approval gap. Do not claim production approval, release approval, or external trust. |
| 006 Change Evidence Packet Core | Task ledger now maps 27 / 27 implemented, locally verified, and reviewed tasks to packet schemas, docs, Go package `internal/packet`, CLI surfaces, fixtures, tests, verification commands, closure review, PR/merge evidence, and maintainer decision `accept_gap` for the historical pre-implementation approval gap. | Close Spec 006 as implemented, reviewed, merged, and post-merge verified with accepted historical approval gap. Do not claim production readiness, release approval, or external trust. |
| 007 GitHub OSS Demo Packet | Draft/demo plan plus focused split re-review; tasks are 22 / 22. Maintainer decision `approve_existing_demo_track` approves the existing `fall-out-bug/sdp-trace-demo-jvm-gsd` route, setup PR #16, feature PRs #16-#20 as first slices, negative draft PR #21 for theater explanation, and no separate polished sales repository for this closure. External demo repo evidence covers setup PR, feature packets 1-5, negative draft PR, CI, local Bazel smoke tests, first-packet gate replay, v1 baseline tag, and buyer rehearsal. | Close Spec 007 as scoped existing-demo-track evidence. Rehearsal evidence does not claim release, production trust, compliance, semantic-quality approval, or signed external trust. |
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
| 018 Core/Policy Split And Pi Delivery | Task ledger now maps 11 / 11 implemented, reviewed, and locally verified tasks to command stability, package ownership, extension boundary, source locality, core-first docs, integration verification, maintainer decision `approve_core_extension_direction`, and follow-up Spec 020 / Spec 021. | Close Spec 018 as approved core/extension direction and prepared follow-up implementation specs. Do not claim command removal, separate binaries, production readiness, release approval, or external trust. |
| 019 Repo Realignment, Monitoring, And Gate Readiness | Tasks are 16 / 16. Major workstreams are implemented, PR #63 merged after final-head CI passed, and maintainer decision `split_successor` moves residual governance debt to Spec 022 without claiming retroactive pre-implementation approval or PR #60 merge approval. | Close Spec 019 as split-successor governance state. Do not claim full approval of PR #60, production trust, release approval, or signed external trust. |
| 020 Core Query Package Split | Draft follow-up implementation spec prepared by Spec 018 closure. No active task checkboxes are opened in the current closure route. | Future implementation surface only; do not treat as current closure debt until explicitly taken into work. |
| 021 Source File Locality Cleanup | Draft follow-up implementation spec prepared by Spec 018 closure. No active task checkboxes are opened in the current closure route. | Future implementation surface only; do not treat as current closure debt until explicitly taken into work. |
| 022 Post-Merge Governance Closure | Active governance closure spec on `codex/022-post-merge-governance-closure`. Tasks are 42 / 42 after plan/task review closure, US1 evidence summary, US2 decision-state update, US3 navigation sync, final local verification, and fixed-handoff PR-ready re-review. PR #60 and PR #63 live refresh is recorded; PR #60 merge approval remains `not_assessed`. | Create PR and collect live branch PR/CI evidence. Do not claim retroactive PR #60 approval, production trust, release approval, or external attestation. |

## Open Task Classification

Current task ledger state is 647 / 647 checked, leaving 0 open Spec 022 task
boxes.

| Category | Specs / tasks | Closure meaning |
| --- | --- | --- |
| Approval or maintainer gates | none in current task ledger | Spec 022 is PR-ready governance closure evidence. Approval, trust, and merge states remain separate until PR review and live CI evidence exist. |
| Review / PR evidence closure | none currently isolated as review-only after this refresh | Future review gaps must still be closed with fresh review, live PR evidence, and explicit trust-boundary wording, not task-box cleanup alone. |
| Retire or re-scope stale planning | none currently isolated as stale after this refresh | Future stale specs should be retired only with an explicit supersession map. |

## Recommended Route

The implemented-local and stale-ledger backlog has no open Spec 022 task
checkboxes. The next work is:

1. **PR evidence**: create the Spec 022 PR and query live branch PR/CI state.
2. **Successor surfaces**: Spec 020 and Spec 021 remain prepared as draft
   follow-up specs. Do not treat them as active implementation debt until
   maintainers explicitly take them into work.
3. **External demo surface**: T226 is closed for the current first-run
   OpenCode/GSD-Redux observation path. Keep mutation/test as explicit
   `not_assessed` for the replayed no-op phase and do not upgrade that closure
   into feature-delivery or test-success evidence.

## Decision Surface

The durable row-by-row decision ledger lives in
[`docs/closure-decision-ledger.md`](closure-decision-ledger.md).

| Decision needed | Open tasks | Current evidence | Valid closure outcomes |
| --- | --- | --- | --- |
| Spec 019 post-merge governance | T019-001/T019-002/T019-003/T019-004; Spec 022 closed for PR-ready evidence | Maintainer decision `split_successor` moves remaining governance debt to active Spec 022; PR #60 merge approval remains `not_assessed` after 2026-06-01 live refresh. | Split-successor state remains current; Spec 022 is PR-ready but must still collect live branch PR/CI evidence before any stronger closure claim. |
| T226 first-run OpenCode/GSD observation | 001 T226 | GSD-Redux replacement observes setup/model/source/digest/tool/phase facts. Mutation/test remain `not_assessed` because the replayed phase was already complete and execute-phase performed no mutation or test action. | Closed as observed first-run path with unavailable dimensions retained, not green. |

The remaining risk is active Spec 022 closure work, future successor-spec
execution, repository drift, and any later maintainer decision to take Spec 020
or Spec 021 into work.

<!-- sdp-trace-claim: claim=profile_passed; subject=spec-closure-route-audit; state=pass; profile=repo_baseline_structural; evidence=state:claim_tags_consistent -->
