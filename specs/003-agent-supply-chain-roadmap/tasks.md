# Tasks: Agent Supply Chain Roadmap

**Input**: Design documents from `/specs/003-agent-supply-chain-roadmap/`
**Prerequisites**: `spec.md`, `plan.md`, `research.md`, Socratic review before
implementation approval
**Tests**: This slice is roadmap-only. Verification is Markdown sanity,
`git diff --check`, and optional repo baseline `go test ./...`.

## Phase 0: Roadmap Draft

- [x] T001 Create isolated worktree for roadmap drafting.
- [x] T002 Draft `spec.md` with CTO value, scope, user stories, functional
  requirements, evidence theater taxonomy, and success criteria.
- [x] T003 Draft `plan.md` with integration strategy, roadmap slices, review
  checkpoints, and non-goals.
- [x] T004 Draft `research.md` with product decisions, integration notes, and
  research gaps.
- [x] T005 Run Socratic/product review of roadmap package before turning any
  item into implementation scope. See
  `reviews/2026-05-10-pi-socratic-review.md`.
- [x] T006 Resolve or explicitly defer critical/major roadmap review findings,
  starting with Product Contract v0 in `specs/005-product-contract-v0/`.
- [x] T007 Stop for explicit user approval of reviewed roadmap direction.
- [x] T008 Do not approve any integration implementation scope until the item
  maps to Product Contract v0 packet rows.
  Retirement status: `retired_superseded`. The reviewed roadmap direction was
  not approved as a standalone implementation source. Later implementation
  authority moved to Specs 005, 006, 007, 017, 018, 019, and the later Block 30
  / Block 32 PR-review evidence path. No future implementation slice is
  approved from Spec 003.

## Phase 1: CTO Evidence Packet Discovery

- [x] T009 Define the minimum CTO packet summary shape for one GitHub PR.
- [x] T010 Define provider-neutral change-host fields that GitHub can populate.
- [x] T011 Define evidence rows for facts, agent claims, CI witness, review
  evidence, missing evidence, and human decision owner.
- [x] T012 Define evidence theater finding rows and reason codes.
- [x] T013 Add one hand-reviewed packet example only after source artifacts are
  identified; keep it marked example/discovery, not product proof.
  Retirement status: `retired_superseded` by Specs 005 and 006. Product
  contract and change-evidence packet artifacts own the concrete packet rows,
  theater reason codes, examples, renderer, validator, and CLI behavior.

## Phase 2: GitHub-First Adapter Specification

- [x] T014 Map GitHub PR, issue, commit, check, review, Actions run, and artifact
  concepts to provider-neutral change-host terms.
- [x] T015 Identify which GitHub API failures become `cannot_verify` vs
  `not_assessed`.
- [x] T016 Define future adapter placeholders for GitLab, GitFlic,
  Gitea/Forgejo, and Jenkins artifact-only flow without claiming support.
  Retirement status: `retired_superseded` by the product packet and PR-review
  evidence paths. Non-GitHub providers remain future adapter placeholders, not
  current support claims.

## Phase 3: OpenCode + GSD Supply-Chain Packet

- [x] T017 Select one existing or new OpenCode/GSD run as the first packet
  candidate.
- [x] T018 Map native OpenCode/GSD normalized events to packet fields.
- [x] T019 Record GSD phase/task metadata as workflow intent only.
- [x] T020 Identify residual missing evidence for mutation, test, PR, review,
  CI, and signed witness states.
- [x] T021 Define what would be required before saying the OpenCode/GSD slice is
  observed.
  Retirement status: `retired_superseded` by Spec 007 and Blocks 24, 25, and
  31. Demo-route proof remains open where those later specs still have explicit
  demo or first-run tasks.

## Phase 4: Pi And GSD2 Discovery

- [x] T022 Inspect `pi` local session storage/export shape.
- [x] T023 Classify `pi` evidence mode as importable, partial, wrapper-only,
  plugin-required, unsafe, unstable, or `not_assessed`.
- [x] T024 Inspect GSD2 session/state surfaces, including runtime state, git
  isolation, verification, crash recovery, and cost/token fields where available.
- [x] T025 Classify GSD2 evidence mode independently from GSD v1.
- [x] T026 Define redaction and retention constraints for both `pi` and GSD2
  before any parser work.
  Retirement status: `retired_superseded`. Later PR-review work uses `pi` as an
  optional runner with explicit `not_assessed` / `cannot_verify` handling.
  GSD2 import remains future discovery and is not current support.

## Phase 5: Superpowers Intent Mapping

- [x] T027 Identify stable Superpowers artifacts or skill-invocation evidence
  across target hosts.
- [x] T028 Map brainstorming, worktree, planning, TDD, review, and verification
  checkpoints to intent facts.
- [x] T029 Define why skill presence does not prove methodology compliance.
  Retirement status: `retired_superseded`. Methodology-intent proof remains
  out of current product support unless a later reviewed spec adds evidence
  contracts.

## Phase 6: General-Purpose Agent Boundary Spike

- [x] T030 Choose Hermes or OpenClaw for the first software-delivery boundary
  spike.
- [x] T031 Define the boundary: channel/session -> general agent -> delegated
  coding agent or direct repo action -> change host -> CI/artifact.
- [x] T032 Identify stable event/session/API evidence exposed by the selected
  agent.
- [x] T033 Define out-of-scope non-software actions explicitly.
- [x] T034 Define privacy and employee-monitoring guardrails for C-level review.
  Retirement status: `retired_superseded`. No general-purpose agent boundary
  spike is approved from this roadmap; employee-monitoring scope remains out of
  product boundary.

## Phase 7: Signed Attestation Capstone

- [x] T035 Define the minimum evidence packet fields that must exist before
  signing can be meaningful.
- [x] T036 Define DSSE/in-toto/Sigstore target profile and customer private
  equivalent profile.
- [x] T037 Define signing failure states, freshness evidence, and identity policy.
- [x] T038 Ensure signed attestation cannot upgrade missing evidence into trust.
  Retirement status: `retired_superseded` by later witness, checkpoint, release
  proof, and OSS compatibility specs. Signed attestation remains additive
  evidence and cannot upgrade missing evidence into trust.

## Phase 8: Review And Approval Checkpoint

- [x] T039 Run separate product-value, evidence-semantics, and integration-order
  review planes.
- [x] T040 Record findings and dispositions in a roadmap review ledger.
- [x] T041 Update roadmap artifacts after accepted findings.
- [x] T042 Ask for explicit approval before creating implementation slices.
  Retirement status: `retired_superseded`. Review findings remain recorded in
  `reviews/2026-05-10-pi-socratic-review.md`; no implementation slice is
  approved from Spec 003.
