<!--
Sync Impact Report
Version change: template -> 1.0.0
Modified principles:
- Placeholder principles -> Portable Evidence Substrate
- Placeholder principles -> Evidence-Backed Claim States
- Placeholder principles -> SpecKit Trace Flow
- Placeholder principles -> Go-First Product Path
- Placeholder principles -> Review, Verification, and Approval Boundaries
Added sections:
- Repository Constraints
- Development Workflow
Removed sections: none
Templates requiring updates:
- .specify/templates/plan-template.md: updated for sdp-trace gates and paths
- .specify/templates/tasks-template.md: updated for Go-first tasks and tests
- .specify/templates/spec-template.md: reviewed, no required change
Follow-up items: none
-->

# sdp-trace Constitution

## Core Principles

### I. Portable Evidence Substrate

`sdp-trace` MUST remain a portable trust substrate for traceability,
provenance, evidence, gate verdicts, and decision records. Product files MAY
include JSON schemas, Markdown docs, portable examples, and small Go
validation or rendering tools. Product files MUST NOT depend on `sdp_lab`,
Beads, Operator Mode, agentloop, hidden agent state, or any specific harness.

Rationale: the product is useful only if downstream agents and governance
layers can consume the artifacts without inheriting one local runtime.

### II. Evidence-Backed Claim States

Every claim about a gate, verdict, readiness state, proof, or trust outcome
MUST be evidence-backed or explicitly marked `not_assessed` or
`cannot_verify`. Opaque health scores are forbidden. Machine proof beats prose,
task checkboxes, checked-in reports, and checked-in proof JSON unless the proof
is live-verified or externally signed.

Rationale: the repository has previously overclaimed replayable verification;
the constitution must keep missing evidence visible instead of smoothing it
into green or red status.

### III. SpecKit Trace Flow

Use SpecKit-aligned terms first: spec, plan, task, evidence, gate, decision,
trace, and provenance. Every non-trivial implementation chunk MUST have an
explicit SpecKit delta and a current mapping from requirements to changed
artifacts and verification evidence. The project-local `sdp-trace-router`
skill MUST be loaded before generic or upstream Spec Kit skills.

Rationale: Spec Kit adds useful workflow structure, but `sdp-trace` trust
rules define the product boundary and evidence semantics.

### IV. Go-First Product Path

Target product code is Go. New product behavior MUST be small, readable,
testable, covered by focused tests, and free of TODO/FIXME markers. Node.js,
npm, JavaScript, TypeScript, and `.mjs` tooling MUST NOT enter the active
product path. Bash MAY be used only as a thin command launcher when Go adds no
product value.

Rationale: the active product surface is a portable Go/schema/docs substrate,
not a mixed toolchain platform.

### V. Review, Verification, and Approval Boundaries

Fresh local verification is required before claiming pass status. Default
checks are `go test ./...`, `go vet ./...`, `jq empty schema/*.json`, and
`git diff --check`; schema or contract changes require reference, fixture, and
Go struct/schema alignment checks. PR review, green CI, and local verification
are review evidence, not merge approval. Missing approval MUST remain
`merge_approval: not_assessed`.

Rationale: readiness, review evidence, and approval-to-merge are distinct
states with different authority.

## Repository Constraints

- Product docs and examples MUST preserve evidence-state honesty:
  `verified`, `not_assessed`, `assumed`, `blocked`, `failed`,
  `cannot_verify`, and narrower source states when used by a spec.
- Security docs MUST NOT list unverified contact information; unconfirmed
  reporting channels remain `not_assessed`.
- Source-bound proof requires a clean immutable source commit before proof
  generation. If a manifest subject changes after proof generation, proof MUST
  be regenerated from a new source-bound cycle.
- Spec Kit git and GitHub issue skills are installed as workflow helpers.
  They MUST NOT bypass project-local branch, review, commit, PR, or merge
  approval rules.

## Development Workflow

1. Start sdp-trace work by reading `AGENTS.md` and loading
   `sdp-trace-router`.
2. For feature work, use a SpecKit spec, plan, and task set under `specs/`
   unless the change is a small obvious local fix with low risk.
3. For behavior changes, write or update focused tests before implementation
   and verify that failing/passing states are real.
4. For verifier, trust, schema, evidence, or release-proof changes, include
   trace coverage and strict review before PR-ready claims.
5. After implementation, report changed artifacts, exact verification commands,
   states not assessed, remaining risks, and follow-up work if any.

## Governance

This constitution mirrors the durable project rules in `AGENTS.md` for Spec Kit
workflows. If this constitution conflicts with more specific repository
instructions, the more specific project-local rule wins and this file MUST be
amended. Amendments require a documented rationale, affected-template review,
and a semantic version update:

- MAJOR for incompatible governance or principle changes.
- MINOR for new principles or materially expanded required workflow.
- PATCH for clarifications that preserve current obligations.

**Version**: 1.0.0 | **Ratified**: 2026-05-26 | **Last Amended**: 2026-05-26
