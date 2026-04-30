# Implementation Plan: sdp-trace Time-Series Evidence Substrate

**Branch**: `001-sdp-trace-time-series-evidence-substrate` | **Date**: 2026-04-30 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-sdp-trace-time-series-evidence-substrate/spec.md`

## Summary

Reframe `sdp-trace` as the portable evidence, provenance, observation, and time-series trace substrate for AI-assisted delivery. The implementation creates SpecKit-visible artifacts first, then evolves schemas, examples, docs, and pilot run-cards so `sdp-gate` can later apply policy without `sdp-trace` owning gate decisions.

## Technical Context

**Language/Version**: JSON Schema, Markdown, shell validation commands
**Primary Dependencies**: `jq`; JSON Schema validator to be selected in this plan
**Storage**: Files committed to the repository; local ignored pilot outputs under `.sdp-trace-runs/`
**Testing**: `jq empty schema/*.json`; future JSON Schema validation over committed examples
**Target Platform**: Portable repository artifacts
**Project Type**: Schema and documentation substrate
**Performance Goals**: Validation commands should run locally in seconds over committed artifacts
**Constraints**: No dependency on `sdp_lab`, Beads, Operator Mode, agentloop, OpenCode, GitHub, Claude, Codex, or any specific harness runtime
**Scale/Scope**: Pilot covers OpenCode model slices, selected harness families, and JVM/Kotlin/Bazel evidence paths

## Constitution Check

The repository rules act as the constitution for this feature.

| Rule | Status | Evidence |
|---|---|---|
| Use SpecKit terms first | Pass | Spec, plan, task, evidence, gate, decision, trace, provenance are used throughout this package. |
| Keep `sdp-trace` independent | Pass | Plan excludes runtime dependencies on `sdp_lab`, Beads, Operator Mode, agentloop, and harness runtimes. |
| Evidence-backed claims only | Pass | Requirements force evidence references or `not_assessed`. |
| No opaque health scores | Pass | Metric catalog explicitly excludes aggregate health scores. |
| Do not implement without plan/spec | Pass | This SpecKit package is the implementation gate. |

## Project Structure

### Documentation (this feature)

```text
specs/001-sdp-trace-time-series-evidence-substrate/
├── spec.md
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── sdp-trace-sdp-gate-boundary.md
└── tasks.md
```

### Repository Artifacts Expected Later

```text
schema/
├── evidence-event.schema.json
├── observation.schema.json
├── metric-stream.schema.json
└── assessment-input.schema.json

examples/
├── opencode/
├── superpowers/
├── github-speckit/
├── self-trace/
└── jvm-bazel/

docs/
├── concepts.md
├── evidence-policy.md
├── harness-compatibility-matrix.md
├── model-compatibility.md
└── jvm-bazel-guide.md
```

**Structure Decision**: Keep feature planning under `specs/`; keep product docs under `docs/`; keep portable schemas under `schema/`; keep concrete examples under `examples/`; keep ignored raw run outputs under `.sdp-trace-runs/`.

## Phase 0: Research

Research is captured in [research.md](research.md). It maps useful `sdp_lab` ideas into portable `sdp-trace` terms and rejects runtime or policy coupling.

## Phase 1: Design

Design outputs:

- [data-model.md](data-model.md) defines the entities and relationships.
- [contracts/sdp-trace-sdp-gate-boundary.md](contracts/sdp-trace-sdp-gate-boundary.md) defines product ownership boundaries.
- [quickstart.md](quickstart.md) defines how a repository observer validates the artifacts.

## Phase 2: Task Breakdown

Task breakdown lives in [tasks.md](tasks.md). Beads issues mirror these tasks for execution discipline only:

| SpecKit task area | Beads mirror |
|---|---|
| Extract portable ideas from `sdp_lab` | `sdp-trace-cdn.1` |
| Define `sdp-trace` / `sdp-gate` boundary | `sdp-trace-cdn.2` |
| Design time-series observation contract | `sdp-trace-cdn.3` |
| Design evidence event and provenance contract | `sdp-trace-cdn.4` |
| Define process metric catalog | `sdp-trace-cdn.5` |
| Create harness and model run-cards | `sdp-trace-cdn.6` |
| Add JVM/Kotlin/Bazel fixture plan | `sdp-trace-cdn.7` |
| Add schema validation and fixture strategy | `sdp-trace-cdn.8` |
| Self-trace this feature with the new contracts | `sdp-trace-cdn.12` |
| Build customer pilot evidence package outline | `sdp-trace-cdn.9` |
| Update compatibility matrices from evidence | `sdp-trace-cdn.10` |
| Reframe CTO and team docs | `sdp-trace-cdn.11` |

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|---|---|---|
| Multiple schema artifacts | Separate evidence events, observations, metric streams, and assessment inputs have different consumers and validation needs. | One large schema would make `sdp-gate` inheritance less clear and would force unrelated fields into every artifact. |
| Pilot matrix across harness/model/stack | Customer pilot explicitly asks for OpenCode, MiniMax/Kimi/GLM, Superpowers/GSD-style harnesses, and JVM/Kotlin/Bazel. | A single Codex or JVM baseline would not prove portability. |

## Self-Trace Milestone

`sdp-trace` starts tracing itself after the minimal contract set exists and before broad pilot execution.

Minimum prerequisite tasks:

- boundary and public language aligned: T005-T007
- observation and metric stream schemas drafted: T008-T009
- evidence, provenance, and assessment input schemas drafted: T015-T017
- at least one valid example and one `not_assessed` example exist

Self-trace output:

```text
examples/self-trace/
├── evidence-events.json
├── provenance-records.json
├── observations.json
├── metric-stream.json
├── trace-snapshot.json
└── assessment-input.json
```

The first self-trace is not a gate decision. It is a consumer test proving the contracts can describe the development of this SpecKit feature itself.

## Verification Commands

Initial verification:

```bash
jq empty schema/*.json
```

Full validation is a planned task because the repository currently documents that JSON Schema validation is not pinned yet.
