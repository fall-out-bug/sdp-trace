# Tasks: Product Contract v0

**Input**: Design documents from `/specs/005-product-contract-v0/`
**Prerequisites**: `003-agent-supply-chain-roadmap` Socratic review findings
**Tests**: Contract-only slice. Verification is Markdown sanity,
`git diff --check`, repo baseline `go test ./...`, and focused Socratic review
before approval.

## Phase 0: Contract Draft

- [x] T001 Create `specs/005-product-contract-v0/`.
- [x] T002 Draft `spec.md` defining Product Contract v0, Change Evidence Packet
  v0, required packet rows, evidence states, theater v0, Russian enterprise
  baseline, and P0 classification rule.
- [x] T003 Draft `plan.md` explaining what the contract is and how to get from
  current substrate to reviewed product classification.
- [x] T004 Draft `example.md` with one concrete example packet marked
  example-only, not product proof.
- [x] T005 Draft `traceability.md` mapping current substrate capabilities to
  packet rows and known gaps.

## Phase 1: Roadmap Linkage

- [x] T006 Update `003-agent-supply-chain-roadmap` so P0 integration work must
  pass Product Contract v0 classification before implementation approval.
- [x] T007 Reclassify GitHub, GitFlic/local Git, OpenCode/GSD, `pi`, GSD2,
  Superpowers, and general-purpose agent work as evidence sources for packet
  rows, not standalone P0 product outcomes.
- [x] T008 Add Product Contract v0 to the approval checkpoint before any
  implementation scope approval.

## Phase 2: Focused Review

- [x] T009 Build focused review packet for `005-product-contract-v0`.
- [x] T010 Run Socratic pi-review focused on whether the contract creates a real
  P0 classification rule against substrate-only P0 work.
- [x] T011 Record review findings and dispositions in
  `reviews/2026-05-10-full-pi-review.md`.
- [x] T012 Apply first full-review fixes: packet template, bundle manifest,
  profile taxonomy, source-strength anti-ranking rule, theater row semantics,
  row-specific rules, and local enterprise baseline example.
- [x] T013 Run reshuffled role/model re-review and record dispositions in
  `reviews/2026-05-10-rereview.md`.
- [ ] T014 Stop for explicit user approval of reviewed Product Contract v0.

## Phase 3: Future Implementation Planning After Approval

- [ ] T015 Define `change-evidence-packet-v0` schema or Go model.
- [ ] T016 Define packet renderer inputs and safe redaction behavior.
- [ ] T017 Add baseline local/self-hosted fixture.
- [ ] T018 Add rich GitHub fixture only as one evidence source.
- [ ] T019 Add focused theater reason-code derivation tests.
- [ ] T020 Add CLI or command surface only after schema/model and fixtures are
  reviewed.

## Completion Rule

This slice is not complete until Product Contract v0 is re-reviewed and
approved. Checked draft tasks are draft-complete, not product approval.
Implementation tasks T015-T020 are placeholders for later planning, not
authorization to write code.
