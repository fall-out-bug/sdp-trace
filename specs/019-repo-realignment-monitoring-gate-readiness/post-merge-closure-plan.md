# Spec 019 Post-Merge Closure Plan

Date: 2026-05-26

## Current State

PR #60 was merged to `main` as commit
`657a343a5f310538def9afd509e6c610c713cab0` on 2026-05-25.

The merge did not close Spec 019. The PR body explicitly deferred:

- WS-019-D: live `wrap` output/schema compatibility.
- WS-019-E: monitoring proof pack.
- WS-019-G: harness/gate shard cleanup.

The PR body also left the review checklist unchecked and stated "Do not merge
without human approval." GitHub PR metadata contains no recorded review approval:
`reviews=[]`, `latestReviews=[]`, and `reviewDecision=""`.

Therefore the post-merge state is:

- `merge_approval`: `not_assessed`.
- `spec_019_completion`: `partial`.
- `phase_0_pre_implementation_approval`: missed; cannot be retroactively
  satisfied as a pre-implementation gate.
- `final_closure_gate`: open.

## Current Handoff Surface

PR #63 is the current integration surface for closure work:

- combined PR #62 post-merge closure with PR #31 PR-review CI integration;
- superseded PR #62 and PR #31 for repo-wide handoff;
- passed final-head CI before merge and merged as
  `1ee2c7af53637c7f43bff4e0e7ef9e34d164908e` on 2026-05-26;
- has no recorded maintainer approval, so `merge_approval` remains
  `not_assessed`.

## 10/10 Closure Work

### S019-P0: Post-Merge Truth Repair

Status: completed.

Deliverables:

- Update `tasks.md` so Phase 0 is recorded as a missed pre-implementation gate,
  not an open task that can still happen before implementation.
- Update `docs/roadmap.md` and `docs/spec-reality-ledger.md` to current `main`
  after PR #60.
- Record PR #61 as stale until closed, rebased, or superseded.

Verification:

- `go run ./tools/doccheck`
- `go run ./tools/hygienecheck`
- `git diff --check`

Exit criteria:

- A maintainer can see that PR #60 was a partial merge and that Spec 019 remains
  open without reading GitHub history.

### S019-P1: Live Wrap Manifest Schema

Status: completed.

Decision:

- Add a dedicated schema for the current live `run.json` manifest emitted by
  `wrap`, with `schema_version` equal to `block10-event-v1`.

Rejected alternatives:

- Force live `wrap` to emit the existing heavy
  `flight-recorder-run.schema.json` shape: larger blast radius and changes the
  current recorder contract.
- Keep the mismatch as an expected failure: insufficient for 10/10 closure.
- Make `check-jsonschema` mandatory in the product path: external validator
  availability should remain optional.

Reversibility: high. A dedicated manifest schema can be deprecated or migrated
without changing live recorder behavior.

Risk if wrong: medium. The wrong schema boundary can make downstream monitoring
integrations validate the wrong artifact.

Confidence: high.

Deliverables:

- Current-run manifest schema under `schema/run-manifest.schema.json`.
- Schema index and README updates.
- OSS compatibility probe changed from "expected fail" to live manifest
  validation.

Verification:

- `jq empty schema/*.json`
- `go test -count=1 ./...`
- `go run ./tools/schemadoc -verify-readme`
- Live `wrap` output validates against the new manifest schema.

### S019-P2: Monitoring And Gate Proof Pack

Status: completed.

Deliverables:

- Reproducible example pack that runs `wrap` or `run`, `verify`, `query`,
  `report`, `gate`, and telemetry export.
- Harness observation example that preserves unsafe input rejection.
- Missing evidence examples that remain `not_assessed` or `cannot_verify`.
- README with copy-pasteable commands and expected states.

Verification:

- Proof pack commands reproduce from a clean temporary directory.
- Unsafe raw prompt/raw response event is rejected before run creation.
- Missing CI witness and audit-grade evidence remain `cannot_verify`.
- `go test -count=1 ./...`
- `git diff --check`

### S019-P3: Harness/Gate Locality Cleanup

Status: completed.

Deliverables:

- Scoped cleanup of selected harness and gate numeric shards.
- Behavior-named files in touched areas.
- Before/after shard count for the owned paths only: 463 to 435.

Verification:

- Monitoring/gate proof pack still passes.
- `go test -count=1 ./...`
- `go vet ./...`
- `go run ./cmd/sdp-trace command-surface`
- `git diff --check`

### S019-P4: Final Trust Review And Closure

Status: split to follow-up Spec 022.

Deliverables:

- Fresh local verification bundle.
- Fresh live CI evidence for final HEAD.
- Adversarial review until exact `LGTM` with zero findings.
- Maintainer approval, rejection, or further split recorded in follow-up
  Spec 022 after the post-merge state is visible.

Verification:

- `go test -count=1 ./...`
- `go vet ./...`
- `go run ./tools/doccheck`
- `go run ./tools/hygienecheck`
- `go run ./tools/schemadoc -verify-readme`
- `jq empty schema/*.json`
- coverage-backed CRAP check
- MI baseline checks
- final live CI query

Current checked-in evidence:

- Local full verification passed before the docs drift repair commit.
- GLM-4.6 full integration review returned LGTM.
- GLM-4.6 PR-review CI hygiene re-review returned LGTM after parser/test fixes.
- PR #63 final-head CI passed before merge.
- Maintainer approval for PR #60 remains `not_assessed`; residual governance
  is split to Spec 022.

Exit criteria:

- Spec 019 residual governance has moved to Spec 022 by maintainer decision in
  the 2026-05-26 closure thread. Spec 019 must not claim retroactive
  pre-implementation approval or PR #60 merge approval.
