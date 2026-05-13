# Implementation Plan: Agent Supply Chain Roadmap

**Branch**: `003-agent-supply-chain-roadmap` | **Date**: 2026-05-10 | **Spec**: [spec.md](spec.md)
**Input**: Product roadmap specification from
`/specs/003-agent-supply-chain-roadmap/spec.md`

## Summary

Create a SpecKit-shaped product roadmap for `sdp-trace` as a neutral agent
supply-chain evidence layer. The roadmap is discovery and product-planning work;
it does not authorize schema, Go, CLI, or verifier implementation.

This roadmap requires `specs/005-product-contract-v0/` before implementation
approval. Integration work cannot be treated as P0 product progress until it
maps to required Change Evidence Packet v0 rows and target transitions.

The first buyer is C-level, usually CTO. The first product path is GitHub-first
but not GitHub-bound. The near-term integration set is:

- GitHub PR evidence packet;
- OpenCode + GSD delivery chain observation;
- `pi` session import discovery;
- GSD2 discovery as a Pi-SDK-based standalone coding agent;
- Superpowers/GSD workflow-intent mapping;
- one general-purpose agent boundary spike with Hermes or OpenClaw;
- signed attestation as the top trust profile after evidence packet semantics
  stabilize.

## Technical Context

**Language/Version**: Markdown SpecKit artifacts only in this roadmap slice.
Future implementation remains Go-first.
**Primary Dependencies**: Existing `sdp-trace` docs, schemas, examples, and
SpecKit conventions.
**Storage**: Product roadmap artifacts under `specs/003-agent-supply-chain-roadmap/`.
**Testing**: Markdown review, `git diff --check`, optional `go test ./...` for
repo baseline.
**Target Platform**: Portable evidence contracts that can later map to GitHub,
GitLab, GitFlic, Gitea/Forgejo, Jenkins, local CLI agents, and signed witness
systems.
**Project Type**: SpecKit product roadmap, not implementation.
**Constraints**: No Node.js, npm, JavaScript, TypeScript, or `.mjs` product-path
changes. No dependency on GitHub, OpenCode, GSD, GSD2, Superpowers, `pi`,
OpenClaw, Hermes, Claude, Codex, or any specific provider.

## Constitution Check

| Rule | Live verifier state | Evidence |
|---|---|---|
| Spec before implementation | `not_assessed` | Roadmap/spec artifacts are present, but this table is not live verifier output. |
| Keep product independent | `not_assessed` | GitHub is described as first adapter, not product ontology; no current verifier claim is made here. |
| Evidence-backed claims only | `not_assessed` | Tool rows remain `not_assessed` until evidence surfaces are inspected by current commands. |
| Preserve missing states | `not_assessed` | Roadmap keeps `not_assessed`, `cannot_verify`, `missing_telemetry`, `unsupported`, and `not_integrated`. |
| No native policy verdicts | `not_assessed` | External consumers decide merge, release, compliance, HR, and risk outcomes. |
| Go-first product path | `not_assessed` | No active product code or non-Go toolchain is added by this roadmap document. |

## Project Structure

```text
specs/003-agent-supply-chain-roadmap/
|-- spec.md
|-- plan.md
|-- research.md
`-- tasks.md
```

Potential future implementation artifacts after separate approval:

```text
schema/
|-- change-host-event.schema.json
|-- agent-supply-chain-record.schema.json
|-- evidence-theater-finding.schema.json
`-- cto-evidence-packet.schema.json

examples/
|-- github-pr-evidence-packet/
|-- opencode-gsd-supply-chain/
|-- pi-session-import/
|-- gsd2-session-import/
`-- general-agent-boundary/

docs/
|-- agent-supply-chain.md
`-- cto-evidence-packet.md
```

## Integration Strategy

| Integration | First evidence mode | Why | Risk |
|---|---|---|---|
| GitHub PR | Post-hoc import plus CI artifact refs | Fastest CTO-visible packet. | GitHub-specific concepts can leak into product model. |
| OpenCode + GSD | Wrapper/sidecar plus native JSONL import | Already closest real dogfood chain. | One observed profile can be overclaimed as broad support. |
| `pi` | Session import discovery | Minimal runtime may expose cleaner session evidence. | Stable export shape not yet assessed. |
| GSD2 | Session import discovery plus wrapper feasibility | Combines tool and harness on Pi SDK. | Treating GSD2 like GSD would miss runtime-owned state. |
| Superpowers | Workflow-intent mapping | Strong methodology/checkpoint surface. | Compliance should not be inferred from skill presence. |
| Hermes/OpenClaw | Boundary spike only | Tests non-technical staff/general-agent risk. | Scope can sprawl into general agent monitoring. |
| Signed attestation | Top trust profile | Governance capstone. | Premature signing can make weak evidence look stronger. |

## Roadmap Slices

### Slice P0-A: CTO Evidence Packet Shape

Define the packet contract, summary language, evidence rows, theater findings,
and decision-owner fields. This slice can start as docs/examples before schema.

Exit criteria:

- One sample packet maps a PR to facts, claims, missing evidence, and next
  decision owner.
- Packet text does not claim merge, release, compliance, or production trust.
- Every claim row has evidence refs or an explicit missing state.

### Slice P0-B: GitHub Change-Host Adapter Model

Define GitHub as the first change-host adapter without making GitHub the product
ontology.

Exit criteria:

- GitHub concepts map to provider-neutral change-host fields.
- Missing GitHub API access is `cannot_verify` or `not_assessed`.
- Future GitLab/GitFlic/Gitea rows are named as planned adapters, not current
  support.

### Slice P0-C: OpenCode + GSD Supply-Chain Packet

Use the existing OpenCode/GSD observation path as the first real software
delivery chain.

Exit criteria:

- Native OpenCode/GSD events bind to a change packet without hand-authored proof.
- GSD phase/task facts are workflow intent unless separately verified.
- Missing mutation/test/PR/CI facts remain visible.

### Slice P0-D: Pi And GSD2 Discovery

Inspect `pi` and GSD2 session/export surfaces and classify import feasibility.

Exit criteria:

- Discovery reports identify stable artifacts, missing fields, and redaction
  risks.
- `pi` and GSD2 rows remain `not_assessed` until a real artifact is inspected.
- GSD2 is treated as a standalone agent/runtime, not just GSD v1 methodology.

### Slice P1-A: Superpowers Workflow Intent

Map Superpowers skills/checkpoints to intent evidence only.

Exit criteria:

- Brainstorming, worktree, plan, TDD, review, and verification checkpoints can be
  referenced when artifacts exist.
- Presence of a skill does not prove compliance.

### Slice P1-B: General-Purpose Agent Boundary

Pick Hermes or OpenClaw for a single controlled software-delivery boundary
spike.

Exit criteria:

- The spike records upstream channel/agent/session refs when available.
- Delegation from general agent to coding agent is represented as a chain.
- Non-software actions are explicitly out of scope.

### Slice P2-A: Signed Attestation Profile

Bind stable evidence packets to signed statements after packet semantics are
reviewed.

Exit criteria:

- Signing binds packet digest, source refs, witness refs, selected profile,
  identity, and freshness evidence.
- Missing signed evidence keeps signed-trust claims at `cannot_verify`.
- Customer private equivalent is recorded as scoped policy evidence, not as a
  universal support claim.

## Review And Approval Checkpoints

Before any implementation:

- Complete this roadmap package: `spec.md`, `plan.md`, `research.md`,
  `tasks.md`.
- Complete and review `specs/005-product-contract-v0/`.
- Map every proposed P0 integration to one or more Product Contract v0 packet
  rows.
- Run Socratic/product review focused on C-level value, scope control, evidence
  semantics, and integration order.
- Resolve or block critical/major findings.
- Ask for explicit approval of the reviewed roadmap direction.

Before any adapter implementation:

- Identify exact evidence surface for the selected tool.
- Add fixture shape before parser behavior.
- Define redaction/retention safety constraints.
- Define what remains `not_assessed`.

Before any support claim:

- Run real tool/session evidence.
- Retain safe artifacts.
- Validate packet generation.
- Record residual gaps.

## Non-Goals

- Do not build a dashboard in this roadmap slice.
- Do not add schemas, Go code, or CLI commands in this roadmap slice.
- Do not start native plugins before import/wrapper discovery proves value.
- Do not monitor general-purpose agents outside software-delivery boundaries.
- Do not turn signed attestation into a shortcut around weak evidence.
