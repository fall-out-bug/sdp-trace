<!--
Sync Impact Report
Version change: 1.0.0 -> 1.1.0
Modified principles:
- Portable Evidence Substrate -> Portable Evidence Recorder, Not Authority
- Evidence-Backed Claim States -> Evidence State Honesty
- SpecKit Trace Flow -> SpecKit-Compatible Trace Flow
- Go-First Product Path -> Go-First, Small Product Path
- Review, Verification, and Approval Boundaries -> Separate Review, Verification, and Approval
Added sections:
- Authority and Scope Constraints
- Development Workflow
Removed sections:
- Repository Constraints
Templates requiring updates:
- .specify/templates/plan-template.md: updated
- .specify/templates/spec-template.md: updated
- .specify/templates/tasks-template.md: updated
- .specify/templates/commands/*.md: not present
Runtime guidance reviewed:
- README.md: reviewed, no change required
- docs/concepts.md: reviewed, no change required
- docs/agent-entrypoint.md: reviewed, no change required
- docs/reviewer-entrypoint.md: reviewed, no change required
- docs/overclaim-checklist.md: reviewed, no change required
- docs/speckit-compatibility.md: reviewed, no change required
- docs/spec-status-discipline.md: reviewed, no change required
Follow-up TODOs: none
-->

# sdp-trace Constitution

## Core Principles

### I. Portable Evidence Recorder, Not Authority

`sdp-trace` MUST remain a portable evidence substrate for AI-assisted delivery.
It records spec, plan, task, change, evidence, gate facts, decisions, trace, and
provenance. It MUST NOT decide merge approval, release readiness, degradation,
risk acceptance, override approval, or production trust. Those decisions belong
to CI, release governance, customer governance, or another explicit policy
consumer.

Product files MAY include JSON schemas, Markdown docs, portable examples, and
small Go validation or rendering tools. Product files MUST NOT depend on
`sdp_lab`, Beads, Operator Mode, agentloop, hidden agent state, Claude, Codex,
OpenCode, GitHub, or any other specific harness runtime.

Rationale: the product is useful only when downstream agents and governance
layers can consume evidence without inheriting one local workflow or mistaking
the recorder for the approving authority.

### II. Evidence State Honesty

Every claim about a gate, verifier result, external verdict, readiness state,
proof, witness, release claim, or trust outcome MUST be backed by current
evidence or explicitly marked `not_assessed` or `cannot_verify`. Missing
required evidence for a selected profile MUST NOT be converted into `pass`.
Opaque health scores are forbidden.

Machine proof from current Go verifier output wins over prose, task checkboxes,
checked-in reports, review ledgers, mirrors, and checked-in proof JSON unless an
artifact is replayed live or externally signed under the selected profile.
Dirty checkout output is local structural evidence only. Source-bound proof
requires a clean immutable source commit before proof generation.

Rationale: the repository has previously overclaimed replayable verification.
The constitution keeps unknown, missing, stale, and unverifiable evidence
visible instead of smoothing it into green or red status.

### III. SpecKit-Compatible Trace Flow

Use SpecKit-aligned terms first: spec, plan, task, evidence, gate, decision,
trace, and provenance. SpecKit is one supported planning shape, not a runtime
dependency and not the only valid upstream workflow. Tool-specific terms MAY be
mapped in compatibility docs, but core product docs MUST remain portable.

Every non-trivial implementation chunk MUST have a current mapping from user
need to spec, plan, tasks, changed artifacts, evidence, gate facts, decisions,
trace, and provenance. The project-local `sdp-trace-router` skill MUST be
loaded before generic or upstream Spec Kit skills for repository work.

Rationale: Spec Kit adds useful workflow structure, while `sdp-trace` owns the
portable evidence vocabulary and trust boundary.

### IV. Go-First, Small Product Path

Target product code is Go. New product behavior MUST be small, readable,
testable, covered by focused tests, and free of TODO/FIXME markers. Node.js,
npm, JavaScript, TypeScript, and `.mjs` tooling MUST NOT enter the active
product path. Bash MAY be used only as a thin launcher when Go adds no product
value.

Schema or contract changes MUST keep JSON schemas, fixtures, examples, Go
types, docs, and command help aligned. Complexity, CRAP, maintainability, lint,
and drift gates MUST be verified through commands or marked `not_assessed` /
`cannot_verify`; prose promises are not enough.

Rationale: a narrow Go/schema/docs substrate is easier to audit, replay, and
port across harnesses than a mixed toolchain platform.

### V. Separate Review, Verification, and Approval

Verification, review evidence, CI state, PR readiness, merge approval, release
approval, and production trust are separate states. Local commands, CI checks,
and PR review MAY support readiness claims, but they MUST NOT imply approval to
merge or release unless the approving authority is recorded as current evidence.
Missing merge approval MUST remain `merge_approval: not_assessed`.

Trust-sensitive work requires strict review and fresh verification before any
PR-ready claim. Adversarial review is complete only when every retained finding
is resolved or explicitly dispositioned, and missing/off-task/hung review lanes
remain `not_assessed` or `cannot_verify`.

Rationale: approval is an authority decision. Evidence quality improves that
decision, but it does not replace the accountable owner.

## Authority and Scope Constraints

- Current product status is controlled-pilot MVP unless a later live verifier
  profile proves a narrower claim. The repository MUST NOT present itself as a
  production trust authority or universal harness adapter.
- Result state, trust scope, and authority scope MUST remain separate. The
  canonical verifier result states are `observed`, `pass`, `fail`,
  `not_assessed`, and `cannot_verify`; external verdicts and advisory labels
  MUST NOT be substituted for verifier states.
- Local evidence, CI witness evidence, customer authority evidence, and
  external production trust MUST NOT be collapsed into one readiness label.
- Checked-in proof JSON, review reports, ledgers, roadmap rows, and task
  checkboxes are context until replayed by current verifier output or bound by
  accepted external signature.
- Security docs MUST NOT list unverified contact information. Unconfirmed
  reporting channels remain `not_assessed`.
- Scanner, verifier, and reproduction commands in docs MUST be copy-pasteable
  and isolate local state when default external behavior is being claimed.
- Spec Kit git and GitHub issue skills are workflow helpers. They MUST NOT
  bypass project-local branch, review, commit, PR, or merge approval rules.

## Development Workflow

1. Start repository work by reading `AGENTS.md` and loading
   `sdp-trace-router`.
2. For feature work, use a SpecKit spec, plan, and task set under `specs/`
   unless the change is a small obvious local fix with low risk.
3. Before meaningful design or implementation choices, check whether the change
   can be simpler, which edge cases block the simpler path, and whether an
   existing project utility or established open-source pattern already solves
   the problem.
4. For behavior changes, write or update focused tests before implementation
   and verify that failing/passing states are real.
5. For verifier, trust, schema, evidence, profile, witness, packet, gate, or
   release-proof changes, include trace coverage, schema/fixture alignment, and
   strict review before PR-ready claims.
6. After implementation, report changed artifacts, exact verification commands,
   states not assessed, remaining risks, and follow-up work if any.
7. Do not close task checkboxes, review ledgers, or docs after source-bound
   proof if those files are manifest subjects without another source-bound
   proof cycle.

## Governance

This constitution governs Spec Kit workflows in this repository and mirrors the
durable project rules in `AGENTS.md`. If this constitution conflicts with more
specific repository instructions, the more specific project-local rule wins and
this file MUST be amended in the same slice when practical.

Amendments require:

- documented rationale and affected principle list;
- review of `.specify/templates/plan-template.md`,
  `.specify/templates/spec-template.md`, `.specify/templates/tasks-template.md`,
  any `.specify/templates/commands/*.md` files if present, and relevant runtime
  docs;
- semantic version update;
- explicit list of `not_assessed` follow-ups if any dependent artifact was not
  reviewed.

Version policy:

- MAJOR for incompatible governance changes, removed principles, or weaker
  evidence/authority obligations.
- MINOR for new principles, new required workflow sections, or materially
  expanded trust, evidence, or toolchain guidance.
- PATCH for clarifications that preserve current obligations.

**Version**: 1.1.0 | **Ratified**: 2026-05-26 | **Last Amended**: 2026-05-26
