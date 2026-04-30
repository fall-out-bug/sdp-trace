# Feature Specification: sdp-trace Time-Series Evidence Substrate

**Feature Branch**: `001-sdp-trace-time-series-evidence-substrate`
**Created**: 2026-04-30
**Status**: Draft
**Input**: User description: "`sdp-trace` must be the SpecKit-observable evidence and trace substrate. Beads is only discipline support. The CTO question is whether the delivery process is moving toward degradation over time; `sdp-gate` applies policy on top of inherited `sdp-trace` contracts."

## User Scenarios & Testing

### User Story 1 - CTO Reviews Process Movement (Priority: P1)

A CTO reviewing an AI-assisted delivery pilot can inspect accumulated `sdp-trace` artifacts and see what changed over time: evidence quality, scope discipline, correctness signals, review signals, AI behavior, and stack/harness/model slices.

**Why this priority**: This is the product reason for the repository. If the CTO cannot understand whether the process is improving, stable, degrading, or not assessable, the substrate is not useful.

**Independent Test**: A reviewer opens this SpecKit package and the generated examples, then confirms every process signal is backed by evidence or explicitly marked `not_assessed`.

**Acceptance Scenarios**:

1. **Given** a trace snapshot with complete evidence references, **When** a reviewer inspects the metric samples, **Then** each sample links to inspectable evidence or provenance.
2. **Given** missing build/test evidence, **When** an observation is recorded, **Then** the affected sample is marked `not_assessed` instead of inferred.
3. **Given** multiple observations over time, **When** a current window is compared with a previous window, **Then** `sdp-trace` records the movement as data without producing a policy verdict.

---

### User Story 2 - sdp-gate Inherits Trace Contracts (Priority: P1)

A `sdp-gate` implementer can consume `sdp-trace` artifacts as policy inputs without `sdp-trace` deciding pass/fail, readiness, degradation, or override outcomes.

**Why this priority**: `sdp-gate` is built on top of `sdp-trace`; if the boundary is vague, both products will duplicate policy logic and confuse users.

**Independent Test**: Read the boundary contract and verify it names `sdp-trace` ownership and `sdp-gate` ownership separately.

**Acceptance Scenarios**:

1. **Given** an evidence bundle and metric stream, **When** `sdp-gate` applies a policy, **Then** the policy decision is external to `sdp-trace`.
2. **Given** an external gate verdict is recorded as evidence, **When** it appears in a trace, **Then** it is represented as an observed verdict input, not as a decision made by `sdp-trace`.

---

### User Story 3 - Pilot Evaluates Harness, Model, and JVM Stack Slices (Priority: P1)

A pilot operator can run evidence-focused assessments across OpenCode, Superpowers-style harnesses, `gsd`, `gsd2`, Oh My OpenAgent, MiniMax, Kimi, GLM, and JVM/Kotlin/Bazel targets.

**Why this priority**: The customer pilot explicitly needs these slices. Unsupported claims here would destroy trust.

**Independent Test**: Run-card artifacts define exact expected outputs, provenance fields, unsupported-claim capture, and `not_assessed` behavior before any compatibility claim is made.

**Acceptance Scenarios**:

1. **Given** an OpenCode run with Kimi, GLM, or MiniMax, **When** the run completes, **Then** model identity, harness identity, evidence artifacts, and unsupported claims are recorded.
2. **Given** a Kotlin+Bazel target, **When** stack detection runs, **Then** Bazel ownership is based on files such as `BUILD`, `BUILD.bazel`, `MODULE.bazel`, or `.bazelrc`, not Maven or Gradle assumptions.
3. **Given** a harness cannot export tool logs, **When** a pilot row is updated, **Then** the missing capability remains `not_assessed` or `TBD`.

---

### User Story 4 - Repository Observer Finds SpecKit Evidence (Priority: P2)

A repository observer can understand current scope and proof by reading SpecKit artifacts without needing Beads.

**Why this priority**: Beads is a discipline tool. The repository-facing plan and evidence must live in committed SpecKit files.

**Independent Test**: A reviewer can start from `/specs/001-sdp-trace-time-series-evidence-substrate/spec.md`, follow `plan.md` and `tasks.md`, and map task status to committed artifacts.

**Acceptance Scenarios**:

1. **Given** a fresh clone without Beads context loaded, **When** a reviewer opens `specs/001-sdp-trace-time-series-evidence-substrate/`, **Then** they can understand the feature, plan, tasks, contract, and evidence expectations.
2. **Given** Beads issues exist, **When** they are inspected, **Then** they reference this SpecKit spec as secondary tracking, not the other way around.

## Edge Cases

- A source system cannot expose raw logs: record the missing field as `not_assessed` and keep the run usable.
- A model or harness reports its own identity inconsistently: preserve observed identity and add an unsupported-claim item.
- A PR/MR does not exist: evidence events must support local branch, commit, file, command, or manual sources without PR-only assumptions.
- Customer data cannot be committed: examples and summaries must be sanitized while preserving artifact references, hashes, or redaction notes.
- A consuming policy wants thresholds: thresholds belong to `sdp-gate` or another policy engine, not to `sdp-trace`.

## Requirements

### Functional Requirements

- **FR-001**: `sdp-trace` MUST define portable contracts for evidence, provenance, observations, metric samples, metric streams, trace snapshots, and assessment inputs.
- **FR-002**: `sdp-trace` MUST NOT decide process pass/fail, merge readiness, degradation, override, or policy outcomes.
- **FR-003**: `sdp-trace` MUST state that `sdp-gate` is built on top of `sdp-trace` and inherits its contracts while owning policy evaluation.
- **FR-004**: Every metric sample MUST reference inspectable evidence or be marked `not_assessed`.
- **FR-005**: The metric catalog MUST avoid opaque aggregate health scores.
- **FR-006**: The contract MUST support moving time windows without requiring a fixed baseline.
- **FR-007**: The contract MUST support dimensions for repository, scope, team, harness, model family, model version when available, stack, build system, and time window.
- **FR-008**: The pilot run-card set MUST include OpenCode with MiniMax, Kimi, and GLM model slices.
- **FR-009**: The pilot run-card set MUST include Superpowers-style, `gsd`, `gsd2`, and Oh My OpenAgent harness rows as evidence-backed or explicitly `not_assessed`.
- **FR-010**: The JVM pilot path MUST close the Kotlin+Bazel evidence gap; Java+Bazel and Kotlin+Gradle are not sufficient proof.
- **FR-011**: Public docs MUST use SpecKit-aligned terms first: spec, plan, task, evidence, gate, decision, trace, provenance.
- **FR-012**: Public docs MUST not imply dependency on `sdp_lab`, Beads, Operator Mode, agentloop, OpenCode, GitHub, Claude, Codex, or any specific harness runtime.
- **FR-013**: Schema and example artifacts MUST be machine-checkable by documented commands.
- **FR-014**: Compatibility matrices MUST only claim support when backed by committed examples, run artifacts, or documented `not_assessed` status.

### Key Entities

- **Evidence Event**: One observed proof item from a source such as CI, command output, file inspection, review, scanner, harness log, model output, or manual sign-off.
- **Provenance Record**: Actor, model, harness, tool, command, artifact, timestamp, and source chain metadata when available.
- **Observation**: A dated statement about process state or behavior, backed by one or more evidence events.
- **Metric Sample**: A numeric, boolean, categorical, or count value measured for a dimension set and time window.
- **Metric Stream**: Ordered metric samples over time for the same metric name and comparable dimensions.
- **Trace Snapshot**: A point-in-time graph linking specs, plans, tasks, changes, evidence, observations, external verdict inputs, and decisions.
- **Assessment Input**: A package of trace artifacts prepared for a policy engine such as `sdp-gate`.
- **Pilot Run-Card**: A repeatable harness/model/stack assessment recipe with prompt, expected artifacts, provenance capture, validation, and `not_assessed` rules.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A repository observer can find the canonical feature spec, plan, and tasks under `specs/001-sdp-trace-time-series-evidence-substrate/`.
- **SC-002**: At least one contract document explicitly separates `sdp-trace` data ownership from `sdp-gate` policy ownership.
- **SC-003**: The implementation plan identifies every current Beads task as secondary tracking for a SpecKit task or artifact.
- **SC-004**: No new public artifact claims `sdp-trace` decides degradation, readiness, gate pass/fail, or override.
- **SC-005**: The pilot plan contains explicit run-card coverage for OpenCode+MiniMax, OpenCode+Kimi, OpenCode+GLM, and Kotlin+Bazel.
- **SC-006**: The schema validation plan documents a command that excludes benchmark checkouts and validates committed `sdp-trace` JSON artifacts.

## Assumptions

- `sdp-gate` will consume `sdp-trace` artifacts but will live in a separate product/repository boundary.
- Beads remains useful for local work tracking, but Beads is not a product dependency and is not the repo observer's source of truth.
- The initial implementation may be schema and documentation heavy before adding tiny validation tools.
- Customer pilot artifacts may need sanitization before committing summaries to the repository.
