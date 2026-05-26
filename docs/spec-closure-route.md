# Spec Closure Route

Status: working audit, current as of 2026-05-26.

This document audits specs 001-019 against the current repository state and
defines the route to close, defer, or repair each spec. It is a planning and
coordination artifact, not proof that any spec is complete.

## Evidence Snapshot

| Evidence | Current State |
| --- | --- |
| Source commit inspected | `1ee2c7af53637c7f43bff4e0e7ef9e34d164908e` plus uncommitted docs updates |
| Open GitHub PRs | none found by `gh pr list --state open` |
| Latest `main` CI | pass, run `26436136055` |
| Spec directories | 19 |
| SpecKit triplets | 19 / 19 |
| Task checkboxes | 470 / 605 checked |
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
| 001 Time-Series Evidence Substrate | Broad substrate is substantially implemented through historical blocks; `tasks.md` is 232 / 244 checked and `blocks/` has 93 files. Not closed: T192, T209-T217, T226, T244 remain open. | Split into closed historical blocks vs live blockers. Do not try to close 001 wholesale. Close Block 21/25/32 leftovers with fresh review evidence, then decide whether 001 becomes `historical` plus successor specs. |
| 002 Authority Envelope Boundary Observation | Implementation appears present: authority schemas, docs, fixtures, Go parser/evaluator, and tests are checked in; tasks are 33 / 35. Missing PR-level review and merge/post-merge closure in T034-T035. | Run current local verification and focused PR-level review against authority artifacts. If no findings remain, update spec/roadmap to `complete` or `historical implemented`, with trust scope limited to local verifier behavior. |
| 003 Agent Supply Chain Roadmap | Roadmap/spec discovery artifact only; tasks are 5 / 42. It was not implemented as product. | Do not close as implemented. Either retire as superseded by specs 005-007/017/018/019, or reopen only if a current supply-chain roadmap is still needed. |
| 004 MVP Readiness Hardening | Most implementation work is checked off: 40 / 43. Local quality gates have since improved, but final PR-level review, named reviewer sign-off, and merge approval remain open. | Re-run MVP readiness against current main. If current gates pass, replace stale CRAP/MI blocker text with current evidence, then close T040-T042 through review/sign-off or explicitly defer the MVP readiness claim. |
| 005 Product Contract v0 | Contract design and reviews exist; tasks are 13 / 20. Implementation tasks T015-T020 remain open, although some packet schemas/code now exist elsewhere. | Reconcile against current packet implementation before writing code. Either mark product-contract design approved and hand implementation ownership to 006, or update 005 to reflect the implemented packet surface. |
| 006 Change Evidence Packet Core | Task ledger says 0 / 27, but repository contains packet schemas, docs, Go package `internal/packet`, and packet CLI surfaces. The spec ledger is stale. | Run a spec-vs-implementation audit. Map existing packet artifacts to T004-T027, identify real gaps, then update tasks and run focused review. This is high priority because downstream specs depend on 006. |
| 007 GitHub OSS Demo Packet | Draft/demo plan only; tasks are 5 / 22 and depend on 006. | Keep deferred until 006 is reconciled. Re-scope before implementation because the product direction has shifted toward core evidence substrate and PR-review integration. |
| 008 Invisible Flight Recorder | Tasks are 26 / 26 and implementation/review artifacts exist. Roadmap still says PR/final-head CI evidence pending. | Verify whether the PR/merge happened historically. If already merged, update closure state from current `main`; if not, run final-head CI/review equivalent and close as implemented-local or historical. |
| 009 Machine-Readable Command Surface | Tasks are 14 / 14 and command-surface output works. Review directory absent, but tasks claim PI review and implementation review happened. | Reconstruct missing review evidence or mark review evidence `not_assessed`. Run command-surface drift checks against docs and close as implemented if review gap is accepted. |
| 010 Command Package Organization | Tasks are 14 / 14; review artifacts exist; current main contains the command-family split. | Run current behavior snapshot checks and close. This is likely closeable with no product code if review artifacts are sufficient. |
| 011 Schema Docs Generation | Tasks are 14 / 14; `tools/schemadoc`, schema index, README verification, and CI hooks exist. | Run `go run ./tools/schemadoc` and `go run ./tools/schemadoc -verify-readme`; if pass, close with local/CI evidence. |
| 012 Repo Hygiene And Artifact Boundary | Tasks are 12 / 12; `tools/hygienecheck` exists and CI uses hygiene checks. | Run hygienecheck and close if current main passes. Record any ignored local clutter as local-only, not repo proof. |
| 013 Contributor Onboarding And Verification | Tasks are 0 / 12, but onboarding docs already exist from adjacent work. Not trace-bound to this spec. | Audit existing contributor quickstart/install/onboarding docs against 013 tasks. Either check off what is truly covered or retire 013 as superseded by docs UX/core-first specs. |
| 014 Docs UX And Command Guidance | Tasks are 15 / 15; docs such as reviewer entrypoint, output map, profile guide, overclaim checklist exist. | Run doccheck plus a focused docs UX/review refresh. Likely closeable after current evidence is recorded. |
| 015 Spec Governance And Roadmap | Tasks are 17 / 17 after the status-discipline update; roadmap and status discipline exist. | Close after this audit route is linked and verified. The remaining risk is stale roadmap facts after future merges; keep freshness rule active. |
| 016 Production Adoption Security Baseline | Tasks are 10 / 10; docs exist; PR #59 is now merged per roadmap. External audit, customer adoption, signed release, and production trust remain `not_assessed`. | Close only as controlled-pilot/security-baseline docs. Do not close production adoption. Update wording to separate local security baseline from external production trust. |
| 017 OSS Replacement Compatibility And Benchmarks | Tasks are 9 / 11. Tooling and benchmark harnesses exist. Supply-chain prototype remains open; roadmap/docs index update remains open. | Finish or explicitly defer T017-040, update docs index/roadmap for T017-080, then close with `osscompat` states preserving optional-tool `not_assessed`. |
| 018 Core/Policy Split And Pi Delivery | Tasks are 0 / 11, but the repository already contains command stability, package ownership, extension boundary, and core-first docs. The task ledger is stale. | Run a spec-vs-implementation audit like 006. Map existing docs to T018 workstreams, check command-surface output, then close or rewrite remaining tasks. |
| 019 Repo Realignment, Monitoring, And Gate Readiness | Tasks are 11 / 16 after PR #63 integration. Major workstreams are implemented; Phase 0 HITL approval and T019-120 remain open. | Treat as post-merge governance debt. Review the `post-merge-closure-plan.md`, decide whether maintainers accept the already-merged work, then close HITL gates or mark them `not_assessed` with a successor spec. |

## Recommended Route

1. **Repair stale ledgers first**: specs 006 and 018 have implementation
   artifacts despite unchecked tasks. Closing anything downstream before this
   reconciliation will keep producing contradictory roadmap state.
2. **Close low-risk implemented specs**: 010, 011, 012, 014, 015 can likely be
   closed with docs/tool verification and review-evidence reconstruction.
3. **Close trust-sensitive implemented specs**: 002, 008, 009, 016 need focused
   review and explicit trust-boundary wording before closure.
4. **Resolve partial/blocked specs**: 004, 005, 017, 019 need targeted closure
   work or explicit deferral.
5. **Retire or re-scope stale planning specs**: 003, 007, 013 should not remain
   open forever if their product direction has been superseded.
6. **Do not close 001 as one giant spec**: split it into historical completed
   block evidence plus a small successor list for the remaining open tasks.

## Immediate Next Slice

Start with specs 006 and 018:

- build a mapping table from each unchecked task to current files;
- mark each task as `implemented`, `partial`, `not_started`, or `superseded`;
- run the relevant local verification commands;
- record review gaps as `not_assessed`, not as implicit approval;
- update `docs/roadmap.md` only after that reconciliation.

This slice removes the worst source of confusion: task ledgers that claim no
work while the repository clearly contains product artifacts.

<!-- sdp-trace-claim: claim=profile_passed; subject=spec-closure-route-audit; state=pass; profile=repo_baseline_structural; evidence=state:claim_tags_consistent -->
