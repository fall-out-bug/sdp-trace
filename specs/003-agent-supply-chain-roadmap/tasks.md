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
- [ ] T006 Resolve or explicitly defer critical/major roadmap review findings,
  starting with Product Contract v0 in `specs/005-product-contract-v0/`.
- [ ] T007 Stop for explicit user approval of reviewed roadmap direction.
- [ ] T008 Do not approve any integration implementation scope until the item
  maps to Product Contract v0 packet rows.

## Phase 1: CTO Evidence Packet Discovery

- [ ] T009 Define the minimum CTO packet summary shape for one GitHub PR.
- [ ] T010 Define provider-neutral change-host fields that GitHub can populate.
- [ ] T011 Define evidence rows for facts, agent claims, CI witness, review
  evidence, missing evidence, and human decision owner.
- [ ] T012 Define evidence theater finding rows and reason codes.
- [ ] T013 Add one hand-reviewed packet example only after source artifacts are
  identified; keep it marked example/discovery, not product proof.

## Phase 2: GitHub-First Adapter Specification

- [ ] T014 Map GitHub PR, issue, commit, check, review, Actions run, and artifact
  concepts to provider-neutral change-host terms.
- [ ] T015 Identify which GitHub API failures become `cannot_verify` vs
  `not_assessed`.
- [ ] T016 Define future adapter placeholders for GitLab, GitFlic,
  Gitea/Forgejo, and Jenkins artifact-only flow without claiming support.

## Phase 3: OpenCode + GSD Supply-Chain Packet

- [ ] T017 Select one existing or new OpenCode/GSD run as the first packet
  candidate.
- [ ] T018 Map native OpenCode/GSD normalized events to packet fields.
- [ ] T019 Record GSD phase/task metadata as workflow intent only.
- [ ] T020 Identify residual missing evidence for mutation, test, PR, review,
  CI, and signed witness states.
- [ ] T021 Define what would be required before saying the OpenCode/GSD slice is
  observed.

## Phase 4: Pi And GSD2 Discovery

- [ ] T022 Inspect `pi` local session storage/export shape.
- [ ] T023 Classify `pi` evidence mode as importable, partial, wrapper-only,
  plugin-required, unsafe, unstable, or `not_assessed`.
- [ ] T024 Inspect GSD2 session/state surfaces, including runtime state, git
  isolation, verification, crash recovery, and cost/token fields where available.
- [ ] T025 Classify GSD2 evidence mode independently from GSD v1.
- [ ] T026 Define redaction and retention constraints for both `pi` and GSD2
  before any parser work.

## Phase 5: Superpowers Intent Mapping

- [ ] T027 Identify stable Superpowers artifacts or skill-invocation evidence
  across target hosts.
- [ ] T028 Map brainstorming, worktree, planning, TDD, review, and verification
  checkpoints to intent facts.
- [ ] T029 Define why skill presence does not prove methodology compliance.

## Phase 6: General-Purpose Agent Boundary Spike

- [ ] T030 Choose Hermes or OpenClaw for the first software-delivery boundary
  spike.
- [ ] T031 Define the boundary: channel/session -> general agent -> delegated
  coding agent or direct repo action -> change host -> CI/artifact.
- [ ] T032 Identify stable event/session/API evidence exposed by the selected
  agent.
- [ ] T033 Define out-of-scope non-software actions explicitly.
- [ ] T034 Define privacy and employee-monitoring guardrails for C-level review.

## Phase 7: Signed Attestation Capstone

- [ ] T035 Define the minimum evidence packet fields that must exist before
  signing can be meaningful.
- [ ] T036 Define DSSE/in-toto/Sigstore target profile and customer private
  equivalent profile.
- [ ] T037 Define signing failure states, freshness evidence, and identity policy.
- [ ] T038 Ensure signed attestation cannot upgrade missing evidence into trust.

## Phase 8: Review And Approval Checkpoint

- [ ] T039 Run separate product-value, evidence-semantics, and integration-order
  review planes.
- [ ] T040 Record findings and dispositions in a roadmap review ledger.
- [ ] T041 Update roadmap artifacts after accepted findings.
- [ ] T042 Ask for explicit approval before creating implementation slices.
