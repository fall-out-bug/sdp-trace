# Feature Specification: MVP Readiness Hardening

**Feature Branch**: `codex/mvp-readiness-spec`
**Created**: 2026-05-10
**Status**: Draft - revised after initial Socratic review, approval pending
**Input**: Repository review found that `sdp-trace` is close to controlled-pilot MVP, but docs and code quality do not yet meet a defensible MVP handoff bar.

## Product Boundary

This block prepares `sdp-trace` for a controlled-pilot MVP handoff. It does not
upgrade the product to external production trust, release approval, merge
approval, or a broad harness compatibility claim.

The goal is narrower and testable:

- public docs must describe the command surface that actually exists;
- English and Russian onboarding docs must not imply different product
  surfaces;
- examples must be clearly real fixtures, pilot evidence packages, placeholders,
  or retired material;
- trust-sensitive Go paths must have measurable lint, complexity, and coverage
  gates;
- any remaining quality gaps must be visible as tracked `not_assessed` or
  deferred scope, not hidden behind "MVP ready" language.

## Definitions

**Controlled-pilot MVP** means a bounded, reviewer-usable evidence substrate for
local/repo-observable traces, fixtures, assessment profiles, and source-bound
engineering checks. It does not mean production release approval, human approval,
merge approval, external production trust, universal harness support, or
guaranteed detection of every unwrapped agent run.

**Trust-sensitive Go path** means any package, function, or CLI path that does
one or more of the following:

Plain-language summary: trust-sensitive paths handle evidence that affects what
`sdp-trace` says a reviewer or downstream policy consumer can rely on. If a path
touches output used to judge traceability, evidence completeness, safety,
authority, or verifier state, it is trust-sensitive.

- creates, validates, normalizes, summarizes, or exports trace, evidence,
  provenance, witness, authority, posture, review, or verifier artifacts;
- maps missing, stale, malformed, unsafe, or unavailable evidence into
  `pass`, `fail`, `observed`, `not_assessed`, or `cannot_verify`;
- renders human-readable trust summaries, next actions, or gate/review facts;
- handles safety-sensitive artifact retention, redaction, command descriptors,
  raw harness output, private paths, provider URLs, prompts, model output, or
  credentials;
- is directly named by MVP docs as part of the current command or evidence
  surface.

**First-class MVP evidence** means an example, fixture, or pilot package that is
linked from README, adoption guides, reviewer/agent entrypoints, or schema docs
as evidence of a current product capability. Historical block records,
supplementary fixtures, and placeholders are not first-class unless a current
entrypoint document uses them as capability evidence.

**Current command surface** means the command families and flags printed by
`go run ./cmd/sdp-trace --help` in the checkout under verification.

**Block ledger** means the durable block record that captures reviewed scope,
verification evidence, review dispositions, remaining gaps, and readiness
state. For this spec-intake phase it is `socratic-review.md` plus the block
record under `specs/001-sdp-trace-time-series-evidence-substrate/blocks/`; for
implementation/PR phases it may be extended by implementation and PR review
ledgers.

**CRAP score** means Change Risk Anti-Patterns score for a function:
`CRAP = CC^2 * (1 - coverage/100)^3 + CC`, where `CC` is cyclomatic complexity
and coverage is function-level test coverage. `gocyclo` alone is not a CRAP gate;
it is only one input.

**Evidence state labels**:

- `real_fixture_local`: committed fixture that is structurally verified by
  current local checks; local only, not external production trust.
- `pilot_evidence`: retained pilot package for a bounded observed slice.
- `placeholder`: intentionally incomplete example that must not support MVP
  proof claims.
- `retired`: historical material kept for context but not current evidence.
- `not_assessed`: outside the selected verification scope.
- `cannot_verify`: selected verification attempted but required evidence,
  consistency, or environment was unavailable or contradictory.
- `assessed_gap`: selected verification scope includes the item and current
  evidence shows the item is missing, not implemented, or not enforced. This is
  a gap, not an out-of-scope `not_assessed` state.
- `deferred_scope`: intentionally deferred out of the selected MVP slice with a
  tracked follow-up. This is a scope decision, not a pass.

## MVP Bar Definition

This block may claim controlled-pilot MVP readiness only when:

- FR-001 through FR-010 are satisfied, or explicitly marked `deferred_scope`
  with a tracked follow-up and excluded from the readiness claim;
- `golangci-lint run ./...` passes, or any exclusion is narrow, documented, and
  reviewed;
- the initial complexity gate passes and a ratchet toward `CRAP < 5` is
  documented;
- MVP-critical packages meet the selected coverage floors or are removed from
  the MVP claim surface;
- CI enforces the selected minimum gate set, or CI enforcement is recorded as an
  `assessed_gap` with a concrete follow-up and the block does not claim CI-gated
  readiness;
- docs and examples preserve controlled-pilot scope and do not claim external
  production trust.

Verification of the readiness claim requires:

- a named maintainer or designated reviewer recorded in the block ledger;
- all gate results recorded in the block ledger with command, timestamp, and
  target commit;
- partial compliance labeled `assessed_gap`, `deferred_scope`,
  `not_assessed`, or `cannot_verify`, not "ready";
- the readiness claim made in the PR description or commit-linked block ledger,
  not only in conversational prose.

Initial thresholds for this block:

- lint: `golangci-lint run ./...` has no active findings;
- complexity/CRAP: establish a per-function CRAP baseline from cyclomatic
  complexity and function coverage before decomposition; after the first
  hardening pass no production function may remain above `gocyclo -over 15 .`
  unless recorded as an `assessed_gap`; after the first ratchet milestone,
  production functions must also pass the reviewed `gocyclo -over 10` equivalent
  gate unless recorded as an `assessed_gap`;
- cognitive complexity: if enforced in this block, the threshold must be
  explicit in the implementation ledger; otherwise cognitive complexity
  enforcement remains `not_assessed` and is not part of the readiness claim;
- Maintainability Index: absolute function/file MI `> 70` may not be claimed
  while historical code remains below threshold. The implementation may enforce
  an MI ratchet in CI only if the baseline is checked in, schema-versioned,
  regenerated only by reviewed ratchet changes, and documented as an exception
  set rather than proof that MI passed. New below-threshold functions and
  regressions against the baseline must fail the ratchet gate.
- coverage: record a pre-change per-package and per-function baseline for all
  MVP-critical packages; `internal/contract`, `internal/policy`,
  `internal/export`, and `internal/posture` each need focused happy/error path
  tests or removal from the MVP claim surface; `internal/trace`,
  `internal/harnessobs`, and `internal/verifier` need focused tests for changed
  or MVP-critical paths; any function with high CRAP must be decomposed, tested,
  or recorded as an `assessed_gap`.

Minimum gate set:

- `go test -count=1 ./...`
- `go test -count=1 ./... -coverprofile=<file>` with package-level coverage review
- per-function CRAP baseline and post-change CRAP review from coverage plus
  cyclomatic complexity
- `golangci-lint run ./...`
- selected `gocyclo` and `gocognit` thresholds
- function-level MI ratchet against a checked-in baseline, or MI remains
  `assessed_gap` with no pass claim
- file-level MI ratchet against a checked-in baseline, or MI remains
  `assessed_gap` with no pass claim
- `jq empty schema/*.json` plus changed JSON examples
- `git diff --check HEAD`

## Current Evidence

The review evidence is local checkout evidence only. The main checkout was
`main...origin/main [ahead 12]` at commit
`5f6706b398d6d68bb9a37be2dee4e6aec1037df3` when reviewed, so these findings
describe that local source state and must be re-run before PR readiness. Any
delta between this local baseline, the feature branch, and the target PR branch
must be recorded in the block ledger before readiness can be claimed.

Passed local checks:

- `go test -count=1 ./...`
- `go test -count=1 ./... -coverprofile=/tmp/sdp-trace-cover.out`
- `jq empty schema/*.json`
- `git diff --check HEAD`

Failed or risk-signaling checks:

- `golangci-lint run ./...` failed with `gosimple` in
  `internal/telemetry/prometheus.go` and `ineffassign` in
  `internal/authority/authority.go`.
- `gocyclo -over 4 .` found multiple functions above the requested CRAP floor,
  including `ValidateExportResult` at cyclomatic complexity 66 and
  `normalizeOpenCodeRawLine` at 37.
- Coverage output showed 0.0% package coverage for `internal/contract`,
  `internal/policy`, and `internal/export`, 2.9% for `internal/trace`, 42.7%
  for `internal/harnessobs`, and 51.1% for `internal/verifier`.

Independent review planes used as intake evidence:

- Documentation review: stale command contract, incomplete Russian command
  coverage, placeholder examples, and mixed-language Russian adoption docs.
- Code/CRAP review: complexity and coverage hot spots, missing CI enforcement,
  and lint failures.
- External pi review via `zai/glm-5.1`: verdict `REVISE`; highlighted zero
  coverage in contract/policy/export, extreme cyclomatic complexity, and lint
  failures.

## User Scenarios & Testing

### User Story 1 - First-Time Reviewer Can Trust The Entrypoint Docs (Priority: P1)

A first-time reviewer can compare the documentation entrypoints with the live
CLI help and find no stale command flags, missing current command families, or
contradictory interpretation rules.

**Why this priority**: If the authoritative docs drift from the live CLI, the
repo is asking reviewers to trust stale prose over machine output.

**Independent Test**: A doc freshness check compares the documented command
surface in `docs/agent-entrypoint.md` and `docs/reviewer-entrypoint.md` against
`go run ./cmd/sdp-trace --help`, or the slice records why a machine comparison
is not yet implemented and manually verifies each command line.

**Acceptance Scenarios**:

1. **Given** `go run ./cmd/sdp-trace --help` prints the current command surface,
   **When** a reviewer checks `docs/agent-entrypoint.md`, **Then** documented
   commands and flags match the live surface or are explicitly marked as
   narrative examples outside the authoritative contract.
2. **Given** the Russian command reference exists, **When** a reviewer compares
   it to the English command reference, **Then** current command families such
   as `observe`, `harness`, `authority-envelope`, `ci-artifact-observation`,
   and `pr-review` are not omitted.
3. **Given** a checked-in proof or example is referenced by onboarding docs,
   **When** a reviewer opens it, **Then** its evidence state is explicit:
   `real_fixture_local`, pilot evidence, placeholder, retired, `not_assessed`, or
   `cannot_verify`.

---

### User Story 2 - Pilot Adopter Gets Plain, Accurate Docs In English And Russian (Priority: P1)

A pilot adopter can read the English or Russian adoption path and understand
what `sdp-trace` can try today without decoding internal jargon or getting a
different product promise per language.

**Why this priority**: The MVP handoff target is not only the maintainer. Docs
must let a CTO, reviewer, or pilot engineer distinguish "evidence substrate" from
"release gate" quickly.

**Independent Test**: English and Russian adoption docs are reviewed for
equivalent scope, plain-language first screen, and no hidden production-trust or
harness-compatibility claims.

**Acceptance Scenarios**:

1. **Given** a Russian-speaking reviewer starts with `docs/adoption-guide.ru.md`,
   **When** they read the first screen, **Then** it explains the product in plain
   Russian before using mixed English contract terms.
2. **Given** an English-speaking reviewer starts with `README.md`, **When** they
   read the current status, **Then** they can tell what can be tried now and what
   remains outside MVP without reading historical block docs.
3. **Given** docs mention controlled-pilot readiness, **When** they describe
   missing external evidence, **Then** they keep it as `not_assessed` or
   `cannot_verify`, not a future implied pass.

---

### User Story 3 - Maintainer Can Enforce Code Cleanliness Instead Of Arguing It (Priority: P1)

A maintainer can run a documented local or CI gate and see whether lint,
complexity, and coverage floors meet the MVP bar.

**Why this priority**: A manual review can find current problems, but MVP
readiness requires preventing silent regression.

**Independent Test**: CI or a documented verification command fails on lint
errors and complexity violations above the selected threshold. Trust-critical
coverage floors are explicit.

**Acceptance Scenarios**:

1. **Given** `golangci-lint run ./...` currently fails, **When** this block is
   complete, **Then** the lint gate passes locally and is either enforced in CI
   or explicitly tracked as `assessed_gap` with a concrete follow-up.
2. **Given** `ValidateExportResult` currently has cyclomatic complexity 66,
   **When** this block is complete, **Then** either the function is decomposed
   under an accepted threshold or the threshold/reduction plan is explicit and
   cannot be mistaken for CRAP<5 compliance.
3. **Given** packages such as `internal/contract`, `internal/policy`,
   `internal/export`, and `internal/trace` have very low or zero coverage,
   **When** this block is complete, **Then** focused tests cover their primary
   happy and negative paths or the untested surface is removed/retired from the
   MVP claim.

## Functional Requirements

- **FR-001**: The authoritative command documentation MUST match the live Go CLI
  help for current commands and flags. "Match" means every documented
  authoritative command/flag is accepted by the live CLI help surface, and every
  current command family appears in the authoritative entrypoint docs unless it
  is explicitly marked internal or out of MVP scope. Specifically, the
  `pr-review packet` command family in `docs/agent-entrypoint.md` MUST NOT
  document flags absent from live `--help`.
- **FR-002**: The English and Russian command references MUST cover the same
  current command families, where "current" means the live help surface at
  verification time. Routing Russian readers to the English canonical command
  contract is permitted only as a temporary `deferred_scope` state with a
  documented translation follow-up; it must not be presented as bilingual
  parity.
- **FR-003**: Placeholder examples MUST NOT appear as first-class MVP evidence
  without an explicit placeholder, retired, `not_assessed`, or `cannot_verify`
  label.
- **FR-004**: Adoption docs MUST state the controlled-pilot boundary in plain
  language before dense contract terminology.
- **FR-005**: Lint failures from `golangci-lint run ./...` MUST be fixed or
  excluded only with a documented, narrow reason.
- **FR-006**: Complexity gates MUST be measurable in CI or a documented local
  gate. If immediate `CRAP < 5` is not realistic for all existing code, the
  implementation MUST record the current CRAP baseline, an initial numeric
  threshold, and a concrete ratchet milestone in the block ledger rather than
  claiming compliance.
- **FR-007**: Trust-sensitive code paths MUST have focused tests that cover at
  least one happy path and one negative/error path before they are used in MVP
  readiness claims.
- **FR-008**: Coverage floors MUST be explicit for MVP-critical packages:
  `internal/trace`, `internal/contract`, `internal/policy`, `internal/export`,
  `internal/posture`, `internal/harnessobs`, `internal/verifier`, and any
  package changed by this block.
- **FR-009**: CI MUST enforce the selected minimum gate set, or the block ledger
  MUST keep CI enforcement as `assessed_gap` with a concrete reason and
  follow-up. The block must not claim CI-gated readiness while this gap remains.
- **FR-010**: The block MUST NOT claim external production trust, broad harness
  compatibility, release approval, or merge approval.

## Non-Goals

- Do not redesign the product boundary between `sdp-trace`, `sdp-gate`, and
  downstream policy consumers.
- Do not implement new trust semantics beyond quality and documentation
  hardening.
- Do not treat local green tests as GitHub CI success.
- Do not force every historical block artifact to meet current prose standards;
  only first-class MVP entrypoints and examples are in scope.
- If implementation drifts into these non-goals, record it as scope creep in the
  block ledger, revert or split it, and exclude it from MVP readiness claims.

## Success Criteria

- **SC-001**: A reviewer can run the documented command freshness check and see
  no mismatch between live CLI help and authoritative entrypoint docs.
- **SC-002**: A Russian-speaking reviewer can use current docs without missing
  active command families or receiving a different MVP promise.
- **SC-003**: All first-class examples are labeled with their current evidence
  state and are not confused with proof of unsupported surfaces.
- **SC-004**: `golangci-lint run ./...` passes, or any remaining finding is
  explicitly scoped out with a reviewed reason.
- **SC-005**: `gocyclo` or an equivalent measurable complexity gate is part of
  local verification or CI, with an actual CRAP baseline computed from coverage
  plus cyclomatic complexity, an initial threshold, and ratchet toward
  `CRAP < 5`. The first ratchet milestone must be written in the implementation
  ledger; later milestones may be follow-ups but cannot be described as already
  passing.
- **SC-006**: MVP-critical zero-coverage packages have focused tests or are
  removed from the MVP claim surface.
- **SC-007**: Final docs still say controlled-pilot MVP only; external
  production trust remains `not_assessed` unless live external evidence exists.
- **SC-008**: A named reviewer verifies all MVP bar conditions against the target
  commit and records sign-off in the block ledger before the PR is labeled ready.
- **SC-009**: Schema documentation exists for evidence formats referenced in MVP
  entrypoints, or schema documentation gaps are labeled `assessed_gap` with a
  concrete follow-up.
