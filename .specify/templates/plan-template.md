# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]

**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.22 unless the feature is docs/schema-only.

**Primary Dependencies**: Go standard library by default; any new dependency
requires explicit justification and a simpler/open-source alternatives check.

**Storage**: Repository files only unless the feature explicitly requires
another portable artifact format.

**Testing**: Focused Go tests for behavior changes, plus required repo checks
from `AGENTS.md`.

**Target Platform**: Portable CLI, schema, docs, and examples for local agent
and harness consumption.

**Project Type**: Go CLI/tooling plus JSON schemas, Markdown docs, and portable
examples.

**Performance Goals**: [domain-specific, e.g., 1000 req/s, 10k lines/sec, 60 fps or NEEDS CLARIFICATION]

**Constraints**: Harness-independent evidence recorder, not an approval
authority; evidence-backed or explicitly `not_assessed` / `cannot_verify`; no
opaque health scores; no Node/npm/JS/TS in active product path.

**Scale/Scope**: [domain-specific, e.g., 10k users, 1M LOC, 50 screens or NEEDS CLARIFICATION]

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Evaluate these gates from `.specify/memory/constitution.md`:

- Portable evidence recorder, not authority: no dependency on a specific
  harness/runtime and no native merge, release, risk, degradation, or production
  trust decision.
- Evidence state honesty: gate/verifier/external verdict/readiness claims are
  backed by current evidence or marked `not_assessed` / `cannot_verify`.
- SpecKit-compatible trace flow: requirements map to plan, tasks, evidence,
  gates, decisions, trace, and provenance without making SpecKit a runtime
  dependency.
- Go-first, small product path: product code stays Go-first; Bash remains a
  thin launcher; Node/npm/JS/TS are out of active product scope; schema,
  fixtures, examples, Go types, docs, and command help stay aligned.
- Separate review, verification, and approval: local verification, CI, review
  evidence, PR readiness, merge approval, release approval, and production
  trust remain distinct.
- Authority and scope constraints: result state, trust scope, and authority
  scope use the canonical vocabulary from `docs/agent-entrypoint.md`.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
cmd/        # Product CLI entrypoints
internal/   # Product packages
tools/      # Small repository validation/rendering tools
schema/     # JSON schemas and schema docs
docs/       # User, agent, and reviewer documentation
examples/   # Portable examples and fixtures
specs/      # SpecKit specs, plans, tasks, and evidence
```

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
