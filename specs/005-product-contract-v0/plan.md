# Implementation Plan: Product Contract v0

**Branch**: `005-product-contract-v0` | **Date**: 2026-05-10 | **Spec**: [spec.md](spec.md)
**Input**: Product contract reset after Socratic review found that the roadmap
still lacked a concrete buyer-facing artifact and a P0 classification rule.

## Summary

Create a reviewed Product Contract v0 that makes Change Evidence Packet v0 the
first product output of `sdp-trace` and turns packet row mapping into the rule
for classifying P0 product progress.

This is not implementation. It is the missing product contract layer above the
existing substrate. The deliverable is a SpecKit package that answers:

- what the first product artifact is;
- what rows it must contain;
- how current substrate capabilities feed those rows;
- how Russian enterprise/local deployment constraints affect P0;
- how future features qualify, or do not qualify, as P0 product work.

## Technical Context

**Language/Version**: Markdown SpecKit artifacts only.
**Primary Dependencies**: Existing `sdp-trace` specs, docs, schemas, examples,
and Socratic review findings from `003-agent-supply-chain-roadmap`.
**Storage**: `specs/005-product-contract-v0/`.
**Testing**: Markdown sanity, `git diff --check`, repo baseline `go test ./...`,
and focused Socratic review before approval.
**Target Platform**: Portable product contract usable across GitHub, GitFlic,
GitLab, local Git, Jenkins/TeamCity, OpenCode/GSD, `pi`, GSD2, Superpowers, and
future signed witness profiles.
**Project Type**: Product contract and roadmap classification.
**Constraints**: No Go, schema, CLI, verifier, dashboard, or adapter
implementation in this slice.

## Constitution Check

| Rule | Status | Evidence |
| --- | --- | --- |
| Spec before implementation | Pass | This package is contract-only and requires review before implementation planning. |
| Keep product independent | Pass | Packet rows are provider-neutral; GitHub and GitFlic are sources, not ontology. |
| Evidence-backed claims only | Pass | Support claims require inspected evidence surfaces. |
| Preserve missing states | Pass | Packet rows preserve `not_assessed` and `cannot_verify`. |
| No native policy verdicts | Pass | Packet names decision owners but does not approve merge, release, or compliance. |
| Go-first product path | Pass | No product code or non-Go tooling is added. |

## Project Structure

```text
specs/005-product-contract-v0/
|-- spec.md
|-- plan.md
|-- example.md
|-- example-local-baseline.md
|-- traceability.md
`-- tasks.md
```

## What This Is

Product Contract v0 is an acceptance contract for product backlog, not a slogan.

It says:

1. The first buyer-facing artifact is Change Evidence Packet v0.
2. Packet v0 has required row ids.
3. Existing substrate capabilities must map to those rows.
4. Future P0 work must cite the rows it improves and name the target
   transition.
5. Work that does not cite rows or improve a row is substrate, discovery, or
   future integration, not P0 product progress.

## How To Get There

### Step 1: Fix The Product Output

Choose one canonical output:

- Markdown report plus evidence bundle.

Record static HTML, PR/MR comments, CLI summaries, PDFs, and signed envelopes as
projections. This keeps UI/demo format debates from delaying the contract.

### Step 2: Define Packet Rows

Define required rows in `spec.md`:

- `PC-CHANGE`
- `PC-INITIATOR`
- `PC-AGENT-ROUTE`
- `PC-MUTATION`
- `PC-VERIFICATION`
- `PC-REVIEW`
- `PC-AUTHORITY`
- `PC-THEATER`
- `PC-ATTESTATION`
- `PC-DECISION`
- `PC-RESIDUAL-GAPS`

These row ids become the P0 classification surface.

### Step 3: Write Example Packets

Write `example.md` as a concrete harness-observed packet with partial evidence
and explicit gaps. Write `example-local-baseline.md` as the Russian enterprise
baseline where local Git and internal CI are present but agent-route evidence is
absent.

The examples are not product proof. They are product contract examples.

### Step 4: Map Current Substrate

Write `traceability.md` mapping current docs, schemas, examples, and internal
packages to packet rows. This prevents throwing away current work while also
showing which rows still lack buyer-facing output.

### Step 5: Reclassify Roadmap Work That Does Not Map

Update the roadmap language so integration work is not P0 by itself. GitHub,
GitFlic, OpenCode/GSD, `pi`, GSD2, Superpowers, and general-purpose agents are
evidence sources for packet rows.

### Step 6: Run Focused Socratic Review

Review the classification question:

> Does Product Contract v0 make it impossible to classify substrate-only work
> as P0 product progress without naming packet rows, evidence surfaces, and
> target transitions?

Do not proceed to implementation until critical/major findings are fixed or
explicitly deferred with rationale and the user approves the reviewed
direction.

## Roadmap Reclassification

| item type | P0 product? | condition |
| --- | --- | --- |
| Packet row definition | Yes | It changes the required buyer artifact. |
| Example packet | Yes | It proves the contract can be read and reviewed. |
| Traceability matrix | Yes | It maps substrate to product output. |
| GitHub adapter | Not by itself | P0 only if it fills named packet rows. |
| GitFlic/local Git/Jenkins baseline | Not by itself | P0 only if it proves the packet works without GitHub. |
| OpenCode/GSD import | Not by itself | P0 only if it fills agent route, mutation, or verification rows. |
| `pi`/GSD2 discovery | Discovery | P0 only after evidence surfaces are mapped to packet rows. |
| Signed attestation | P2 profile | It cannot replace missing packet evidence. |
| Dashboard | Later projection | Not needed for Product Contract v0. |

## Review And Approval Checkpoints

Before approval:

- `spec.md`, `plan.md`, `example.md`, `traceability.md`, and `tasks.md` exist.
- `example-local-baseline.md` shows a closed-contour enterprise profile.
- `003-agent-supply-chain-roadmap` references this contract as the classification
  rule before implementation approval.
- Focused Socratic review has usable output.
- Critical/major findings are resolved or explicitly deferred with rationale.
- User approves the reviewed Product Contract v0 direction.

Before implementation after approval:

- Implementation plan names exact packet rows.
- Each feature has `packet_rows`, `evidence_surface`, `start_state`,
  `target_transition`, `buyer_effect`, and `non_goal`.
- Schema/Go/CLI work stays Go-first and evidence-preserving.

## Non-Goals

- Do not implement packet generation in this slice.
- Do not add schemas or Go code in this slice.
- Do not build a dashboard.
- Do not select every future adapter.
- Do not claim GitHub, GitFlic, OpenCode/GSD, `pi`, GSD2, or Superpowers support.
- Do not turn signed attestation into day-one product proof.
