→ Read docs/reviewer-entrypoint.md
→ Read docs/agent-entrypoint.md
→ Read README.md
✱ Grep "missing_telemetry" in . · 70 matches
→ Read docs/concepts.md
✱ Grep "coverage_satisfied|coverage_partial|coverage_unresolved" in . · 19 matches
Let me check for state vocabulary drift and a few more documents referenced in the spec.
→ Read docs/agent-onboarding.md [offset=65, limit=20]
→ Read docs/adoption-guide.en.md [offset=75, limit=20]
✱ Grep "not_integrated" in docs · 5 matches
→ Read docs/concepts.md [offset=85, limit=15]
✱ Grep "Profile Selection|profile decision|Which.*profile" in docs · 2 matches
Now let me verify one more specific thing — how the spec's own acceptance criteria stack against current reality:
I have all the evidence. Here is the adversarial review.
---
## PI Review: 014-docs-ux-command-guidance
### C1 — `missing_telemetry` is an orphan state across 5 docs, absent from the canonical contract
**Severity: critical**
The spec lists `missing_telemetry` in US-002 as a state to distinguish. However, it is **not defined** in either `docs/agent-entrypoint.md` (lines 273-284) or `docs/reviewer-entrypoint.md` (lines 38-41) under "State And Exit Code Contract." It has no exit-code mapping. Despite that, it appears as a valid state in:
- `README.md:133`
- `docs/agent-onboarding.md:73`
- `docs/adoption-guide.en.md:84`
- `docs/harness-integration.md:37`
- `specs/014-docs-ux-command-guidance/spec.md:34` (the spec itself!)
A cold user reading the reviewer entrypoint first learns 5 states with 4 exit codes. Then they hit `missing_telemetry` in other docs with zero guidance on what exit code it produces, whether it's a result state, a trust scope, or a third category. **The acceptance criterion "State vocabulary is consistent across README, concepts, agent entrypoint, reviewer entrypoint, and adoption guide" fails immediately in the current state.** The implementation plan (T011) must either add `missing_telemetry` to the canonical state contract with exit code mapping *or* remove it from the non-entrypoint docs and replace with `not_assessed`/`cannot_verify`.
---
### C2 — `not_integrated` same orphan pattern
**Severity: critical**
`not_integrated` appears as a valid state in `docs/agent-onboarding.md:73`, `docs/harness-integration.md:37`, and `docs/flight-recorder.md:99` but has no definition, exit code, or mention in either entrypoint's state contract. Previous PI reviews (`specs/005-product-contract-v0/reviews/raw/2026-05-10-minimax-trust.md:42`) already flagged this as "never defined." The current spec does not address it. If the slice intends to fix state vocabulary, `not_integrated` must be either added to the canonical contract or scrubbed from all non-spec docs and replaced.
---
### C3 — `coverage_satisfied/coverage_partial/coverage_unresolved` are ad-hoc states without contract entry
**Severity: critical**
`docs/reviewer-entrypoint.md:122` and `docs/agent-entrypoint.md:221` use these as `pr-review` states, but they are **absent** from "State And Exit Code Contract" in both entrypoints. A cold reader cannot tell:
- What exit code maps to `coverage_satisfied`?
- Is `coverage_partial` a pass or a fail?
- How does `coverage_unresolved` relate to `cannot_verify` or `fail`?
If these are valid sub-states of the `pr-review` command, the state contract must list them with exit code mapping, or the command docs must explicitly state they are **not** verifier states and belong to a separate taxonomy.
---
### C4 — Three separate profile vocabularies with zero cross-reference mapping
**Severity: critical**
A cold user encounters **three** different "profile" taxonomies:
1. **Trust profile IDs** (`docs/agent-entrypoint.md:15-19`): `repo_baseline_structural`, `source_bound_local_release`, `external_production_trust`
2. **Assessment profiles** (`docs/reviewer-entrypoint.md:67-73`): `adapter-capture`, `managed-harness`, `forensic-retention`, `ci-artifact-observation`, `authority-envelope`
3. **Witness kinds** (`docs/reviewer-entrypoint.md:75-80`): `github-actions`, `gitlab-ci`, `buildkite`, `customer-pki`
There is **no** table, decision tree, or prose mapping explaining which assessment profiles contribute to which trust profile ID. Does `assess --profile authority-envelope` produce evidence for `external_production_trust`? Can `adapter-capture` ever close `source_bound_local_release`? A user asked "Which assessment profile applies?" (the exact US-001 test) cannot answer this from current docs. FR-003 ("Add 'which profile do I use?' decision tree") is essential but the plan tasks do not explicitly call out that the mapping between these three taxonomies must be created.
---
### C5 — Concepts.md defines a *fourth* state vocabulary with `warn` and without `observed`
**Severity: major**
`docs/concepts.md:86-93` (External Verdict) lists `pass`, `warn`, `fail`, `not_assessed`, `cannot_verify`. This differs from both entrypoints:
- Adds `warn` (absent from entrypoints)
- Omits `observed` (present in both entrypoints)
- Omission of `missing_telemetry` from concepts is correct (it's absent from the contract), but the presence of `warn` introduces yet another undocumented state term
A cold user trying to learn "evidence states" from concepts finds a different set than the reviewer entrypoint. FR-002 must either reconcile this or explicitly demarcate External Verdict states as a separate taxonomy.
---
### C6 — Reviewer entrypoint flat command list precedes task guidance
**Severity: major**
The spec requires (FR-001) "a task-oriented command guide or restructure existing docs to include one" with acceptance criterion "reviewer-entrypoint has a short task path before long references."
Current `docs/reviewer-entrypoint.md`:
- Lines 46-65: long flat "Reviewer Command Surface" list (20+ bullet entries)
- Lines 148-161: "Quick Reference" goal-to-command table (the task-oriented guide)
The long reference comes **first**. A cold user scanning the page encounters the raw command list before any task-oriented path. Reorder so the Quick Reference table moves above the flat command list, or replace the flat list with a reference link.
---
### C7 — No output location map despite US-003 requirement
**Severity: major**
`README.md:36-48` lists output types (`.sdp-trace-runs/`, `.sdp-trace-report/`, query packs, etc.) but does not map which *commands* write to which default paths versus requiring `--out`. The agent entrypoint command table (lines 182-225) shows `--out <file>` on most commands but a cold user must scan the entire table and infer the layout.
Example of the problem: `report` writes to `--out <dir>`, `gate` writes to `--out <file>`, witness uses `--report-dir` for input and `--out` for output. There is no single table showing: "`verify` reads `<run-dir>` and writes to it; `report` writes to `--out <dir>`; `gate` writes to `--out <file>`; witness reads `--report-dir` and writes `--out`; assessment profiles write `--out <file>`."
---
### C8 — Duplicated overclaim forbiddens with no canonical checklist
**Severity: major**
US-004 requires "Reviewer entrypoint contains the canonical checklist; README and agent entrypoint link to it."
Current state:
- `docs/agent-entrypoint.md:297-305`: "Forbidden Claims" (5 items)
- `docs/reviewer-entrypoint.md:135-146`: "What You May State From Output" (overlapping but differently worded)
- Neither document links to the other as the canonical source
The reviewer entrypoint is listed as the canonical location in the spec, but it is not labeled as such and the agent entrypoint lists more forbidden items (lines 301-304 list `trusted_contract_release` and `production_release_verified` separately, which the reviewer entrypoint condenses into one).
---
### C9 — Verification path requires Go toolchain as step 1, blocking cold reviewers
**Severity: medium**
`docs/reviewer-entrypoint.md:16` starts with `go test -count=1 ./...`. A cold reviewer who does not have Go installed (e.g., a compliance reviewer, a PM, a governance stakeholder) is blocked at step 1. The entrypoint claims to be "for a first-time reviewer check in under five minutes" but the first step assumes full development environment. Consider a fallback: "If you don't have Go installed, start with `sdp-trace --help` and the output map below."
