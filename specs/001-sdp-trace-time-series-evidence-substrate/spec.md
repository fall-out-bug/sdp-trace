# Feature Specification: sdp-trace Time-Series Evidence Substrate

**Feature Branch**: `001-sdp-trace-time-series-evidence-substrate`
**Created**: 2026-04-30
**Status**: Draft
**Input**: User description: "`sdp-trace` must be the SpecKit-observable evidence and trace substrate. Beads is only discipline support. The CTO question is whether the delivery process is moving toward degradation over time; `sdp-gate` applies policy on top of inherited `sdp-trace` contracts. A CEO must also be able to see which human is accountable for each stage and whether the contract itself was changed under a trusted release process."

## Non-Negotiable Self-Proof Rule

`sdp-trace` MUST be its own first real consumer.

No customer pilot claim, external compatibility verdict, or "trusted contract foundation" claim is acceptable until this repository can validate a committed self-trace package describing the development of this feature itself: spec, plan, tasks, changed artifacts, evidence, provenance, accountability, observations, metric movement, `not_assessed` gaps, Socratic review results, and contract release verification.

If self-trace or self-attestation is missing, the correct product state is not "trusted"; it is `not_assessed` with a concrete reason.

## User Scenarios & Testing

### User Story 1 - CTO Reviews Process Movement (Priority: P1)

A CTO reviewing an AI-assisted delivery pilot can inspect accumulated `sdp-trace` artifacts and see what changed over time: evidence quality, scope discipline, correctness signals, review signals, AI behavior, and stack/harness/model slices.

**Why this priority**: This is the product reason for the repository. If the CTO cannot understand whether the process is improving, stable, degrading, or not assessable, the substrate is not useful.

**Boundary**: `sdp-trace` answers the CTO question with inspectable movement data: prior/current values, deltas, dimensions, evidence coverage, and `not_assessed` gaps. It does not answer with "yes, we are degrading" or "no, we are not degrading." A yes/no verdict or threshold interpretation belongs to `sdp-gate` or another policy consumer.

**Independent Test**: A reviewer opens this SpecKit package and the generated examples, then confirms every process signal is backed by evidence or explicitly marked `not_assessed`.

**Acceptance Scenarios**:

1. **Given** a trace snapshot with complete evidence references, **When** a reviewer inspects the metric samples, **Then** each sample links to inspectable evidence or provenance.
2. **Given** missing build/test evidence, **When** an observation is recorded, **Then** the affected sample is marked `not_assessed` instead of inferred.
3. **Given** multiple observations over time, **When** a current window is compared with a previous window, **Then** `sdp-trace` records the prior value, current value, delta, units, dimensions, and evidence coverage without producing a policy verdict.

---

### User Story 2 - sdp-gate Inherits Trace Contracts (Priority: P1)

A `sdp-gate` implementer can consume `sdp-trace` artifacts as policy inputs without `sdp-trace` deciding pass/fail, readiness, degradation, or override outcomes.

**Why this priority**: `sdp-gate` is built on top of `sdp-trace`; if the boundary is vague, both products will duplicate policy logic and confuse users.

**Independent Test**: Read the boundary contract and verify it names `sdp-trace` ownership and `sdp-gate` ownership separately.

**Acceptance Scenarios**:

1. **Given** an evidence bundle and metric stream, **When** `sdp-gate` applies a policy, **Then** the policy decision is external to `sdp-trace`.
2. **Given** an external gate verdict is recorded as evidence, **When** it appears in a trace, **Then** it is represented as an observed verdict input, not as a decision made by `sdp-trace`.
3. **Given** an evidence event carries an externally produced quality, strength, or verdict assertion, **When** it is recorded by `sdp-trace`, **Then** the assertion records its producer, policy reference when available, and external origin rather than becoming a native `sdp-trace` assessment.

---

### User Story 5 - CEO/CIO Verifies Accountability and Contract Integrity (Priority: P1)

A CEO or CIO can inspect a contract release or assessment package and identify the human accountable owner, the approval reference, the escalation path, the risk owner, and whether the checked-out contract matches a signed release manifest.

**Why this priority**: A process that says "the AI did it" is not governable. A schema-valid artifact is also not trustworthy if a person or model can quietly simplify the contract and still call the result valid.

**Boundary**: `sdp-trace` records accountability, release integrity, and verification status. It does not decide whether the organization accepts residual risk; that policy decision belongs to `sdp-gate`, management, or another external governance process.

**Independent Test**: A reviewer validates examples proving that AI actors can appear as producers/reviewers but cannot be the sole accountable owner or approver, and that a modified contract artifact fails manifest digest verification.

**Acceptance Scenarios**:

1. **Given** a contract release manifest, **When** a reviewer inspects the accountability section, **Then** it names a human-held DRI, approver, risk owner, escalation path, approval reference, and line of defense.
2. **Given** a model-generated artifact, **When** it is recorded as evidence, **Then** the AI actor appears in provenance, not as the sole accountable owner or approver.
3. **Given** a schema file is changed after manifest generation, **When** manifest verification runs, **Then** the release is not accepted as a trusted contract release.
4. **Given** signature verification is unavailable in a private environment, **When** the verification result is recorded, **Then** it is explicitly `not_assessed` or invalid and cannot be silently treated as trusted.
5. **Given** an OIDC-authenticated signer does not match the trusted identity policy, **When** release verification runs, **Then** the signature does not establish a trusted contract release.

---

### User Story 3 - Pilot Evaluates Harness, Model, and JVM Stack Slices (Priority: P1)

A pilot operator can run evidence-focused assessments across OpenCode, Superpowers-style harnesses, `gsd`, `gsd2`, Oh My OpenAgent, MiniMax, Kimi, GLM, and JVM/Kotlin/Bazel targets.

**Why this priority**: The customer pilot explicitly needs these slices. Unsupported claims here would destroy trust.

**Independent Test**: Run-card artifacts define exact expected outputs, provenance fields, unbacked-claim capture, and `not_assessed` behavior before any external compatibility verdict is recorded.

**Acceptance Scenarios**:

1. **Given** an OpenCode run with Kimi, GLM, or MiniMax, **When** the run completes, **Then** model identity, harness identity, evidence artifacts, and unbacked claims are recorded.
2. **Given** a Kotlin+Bazel target, **When** stack detection runs, **Then** Bazel ownership is based on scope-specific target evidence such as `BUILD`, `BUILD.bazel`, `MODULE.bazel`, `WORKSPACE`, or `WORKSPACE.bazel`, with `.bazelrc` treated as supporting configuration only.
3. **Given** a harness cannot export tool logs, **When** a pilot row is updated, **Then** the missing capability remains `not_assessed` with a reason code such as `missing_export` or `discovery_required`.

---

### User Story 3A - Pilot Operator Proves First Real Slice (Priority: P0)

A pilot operator can run or wrap an actual OpenCode + MiniMax execution against a scoped Kotlin+Bazel target and receive a validated `sdp-trace` package that records the tested-on environment, evidence, provenance, metric sample or stream data, trace snapshot, assessment input, and residual `not_assessed` gaps.

**Why this priority**: Run-cards and schemas do not prove product value. The first defensible value is an evidence-backed E2E path for the exact stack the customer will test: OpenCode, MiniMax, Kotlin, and Bazel.

**Boundary**: `sdp-trace` may ship a reference runner that shells out to external tools and records their observed behavior. It must not make OpenCode, MiniMax, Kotlin, Bazel, Bazelisk, or provider SDKs product dependencies. It must not produce native pass/fail, readiness, compatibility, support, or degradation verdicts.

**Independent Test**: A reviewer can run the Block 06 runner in an environment with OpenCode, MiniMax access, and a Kotlin+Bazel target, then validate the generated evidence package. A committed sanitized report states exactly what was tested and what remains `not_assessed`.

**Acceptance Scenarios**:

1. **Given** OpenCode and a configured MiniMax model id, **When** the runner executes against a scoped Kotlin+Bazel target, **Then** the output records OpenCode version, requested model id, command, source reference, scoped stack markers, produced artifacts, and export limitations.
2. **Given** Bazel or Bazelisk is unavailable, **When** the runner prepares the package, **Then** Bazel command execution is recorded as `not_assessed` with a reason and Block 06 remains incomplete rather than recording success.
3. **Given** the runner emits a package, **When** repository validation runs, **Then** evidence, provenance, trace, metric, and assessment artifacts validate or the package is rejected as incomplete.
4. **Given** a committed sanitized proof report exists, **When** matrix rows are updated, **Then** only the exact OpenCode + MiniMax + Kotlin+Bazel slice may cite that evidence.

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
- A model or harness reports its own identity inconsistently: preserve observed identity and add an `unbacked_claim` item.
- A PR/MR does not exist: evidence events must support local branch, commit, file, command, or manual sources without PR-only assumptions.
- Customer data cannot be committed: examples and summaries must be sanitized while preserving artifact references, hashes, or redaction notes.
- A consuming policy wants thresholds: thresholds belong to `sdp-gate` or another policy engine, not to `sdp-trace`.
- Duplicate, out-of-order, or conflicting evidence arrives: preserve every event with source identity, dedupe key, observed timestamp, and conflict relationship; do not collapse conflicts into a single success/failure claim.
- Pending evidence later resolves: append a superseding event or sample revision with provenance; do not mutate already-published committed examples silently.
- Evidence references point to private systems or raw prompts: commit only sanitized summaries, content hashes, redaction notes, and access-neutral references; credentials and secrets must not appear in committed artifacts.
- An AI actor produced or reviewed an artifact: preserve it in provenance, but require a human-held DRI, approver, risk owner, and escalation path for accountable artifacts.
- A person or model changes a schema, validation script, fixture, or contract doc after release: the contract manifest digest check must fail until a new authorized release manifest is produced.
- Public Rekor or Sigstore verification is unavailable in a private customer environment: record the selected equivalent or `not_assessed`; do not call the release trusted unless the selected identity policy is satisfied.
- A contract release is old but still schema-valid: record freshness or rollback status separately from JSON validation.
- A metric stream mixes assessed and unassessed samples: set stream-level `assessment_state` to `partial` or `not_assessed` so movement comparisons are not over-trusted.
- A signer is technically valid but not authorized for contract releases: reject trusted release status unless signer identity matches the trusted identity policy.
- A safety scan passes but proprietary data risk remains: record scan status as lower-bound automated evidence, and require the accountable evidence owner to approve committed customer-facing examples.
- OpenCode is installed but the requested MiniMax model id is not available: record model availability as `not_observed` with observed provider output, or `not_assessed` when the check cannot run; do not claim the E2E proof is complete.
- OpenCode returns a successful answer but no exportable structured session: preserve stdout/stderr digests and sanitized summary, and record export limitation as evidence.
- Bazel is not installed in the proof environment: Kotlin+Bazel scope detection may be assessed, but Bazel command execution remains `not_assessed`.
- One OpenCode + MiniMax + Kotlin+Bazel run succeeds: only that exact tested slice can move to observed evidence; other models, harnesses, and stacks remain `not_assessed`.
- A recorder attaches after development work has already started: pre-attachment history is explicitly `not_assessed` or `cannot_verify`; the trace MUST NOT imply full-run provenance.
- A local event chain verifies internally but has no witness anchor: the chain may support local development reconstruction only; it MUST NOT support accountability, external trust, or audit-grade claims.
- A requirement changes after commands or test evidence already exist: the original task remains immutable and the change is recorded as a superseding event rather than rewriting expectations.
- Command output or agent transcript contains a secret: the recorder must choose a safe retention/redaction state and must not silently preserve raw secrets in committed artifacts or claim unavailable raw evidence as verified.

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
- **FR-009**: The pilot run-card set MUST include Superpowers-style, `gsd`, `gsd2`, and Oh My OpenAgent harness rows with observed evidence references or explicitly `not_assessed` state and reason codes.
- **FR-010**: The JVM pilot path MUST define the Kotlin+Bazel evidence path; Java+Bazel, Kotlin+Gradle, and synthetic Kotlin+Bazel placeholders are not sufficient proof of real Kotlin+Bazel behavior.
- **FR-011**: Public docs MUST use SpecKit-aligned terms first: spec, plan, task, evidence, gate, decision, trace, provenance.
- **FR-012**: Public docs MUST not imply dependency on `sdp_lab`, Beads, Operator Mode, agentloop, OpenCode, GitHub, Claude, Codex, or any specific harness runtime.
- **FR-013**: Schema and example artifacts MUST be machine-checkable by documented commands.
- **FR-014**: Compatibility matrices are legacy-named evidence matrices. They MUST record observed evidence state, artifact references, gap reasons, and next required evidence; they MUST NOT claim support, readiness, or compatibility as native `sdp-trace` outcomes.
- **FR-015**: Every schema artifact MUST declare JSON Schema Draft 2020-12 and a schema version; every committed example MUST declare or reference the schema version it follows once full validation is enabled.
- **FR-016**: Evidence and provenance artifact references MUST be safe to commit: no secrets, credentials, raw customer data, or private prompt contents; sanitized artifacts MUST preserve hash, redaction, and access notes.
- **FR-017**: `sdp-trace` MUST distinguish native observations from external verdict or evidence-strength assertions through explicit producer and origin fields.
- **FR-018**: Native movement data MUST include comparable window values and deltas; interpretation labels such as degrading, blocked, pass, fail, ready, or not ready MUST remain external verdicts.
- **FR-019**: Accountable artifacts MUST include human-held identity objects for DRI, approver, escalation path, risk owner, authority scope, accountability claim, approval reference, and line of defense.
- **FR-020**: AI actors MAY be recorded as producers, reviewers, critics, or judges, but MUST NOT be the sole accountable owner, approver, risk owner, or escalation owner. Accountable identity objects MUST include machine-readable actor type.
- **FR-021**: `sdp-trace` MUST define a risk classification that records observed autonomy, observed impact, classification source, and externally declared oversight assertions without deriving required oversight internally.
- **FR-022**: `sdp-trace` MUST define a contract manifest that records schema, documentation, validation script, fixture, source commit, approval, compatibility, and previous-manifest digests.
- **FR-023**: The target contract release signing profile MUST be `sdp-trace-signature/sigstore-dsse-keyless-v1`, using an in-toto Statement, DSSE envelope, and Sigstore/Cosign keyless verification where available.
- **FR-024**: A checkout MUST NOT be treated as a trusted contract release unless manifest digest verification, artifact digest verification, source-content verification, signature verification, identity-policy verification, freshness verification, and required external trust evidence are all assessed and successful.
- **FR-024a**: A non-trusted local or partial release proof MAY record missing verification states as `not_assessed`, but those `not_assessed` states MUST keep `trusted_contract_release` false.
- **FR-025**: `sdp-trace` MUST define a trusted identity policy for contract release signing, including OIDC issuer, source URI, protected ref, workflow identity, release captain, and required approval evidence.
- **FR-026**: Contract manifests MUST declare exactly one freshness mechanism: `valid_until` or `freshness_policy`; omission fails validation.
- **FR-027**: Metric streams MUST expose stream-level assessment state when any sample or comparison is partial or `not_assessed`.
- **FR-028**: Contract scaffolding completion evidence MUST include one contract release verification result for the selected signing profile shape or approved private equivalent; this does not by itself establish product trust.
- **FR-029**: Before any external pilot claim, `sdp-trace` MUST validate a committed self-trace assessment input for this feature under the same schemas downstream users are expected to consume.
- **FR-030**: `sdp-trace` MUST distinguish contract scaffolding from product proof. Schema validity, digest verification, local signing, external attestation, and production release verification MUST be recorded as separate evidence states.
- **FR-031**: A contract release proof MUST use an immutable source reference. A placeholder such as `working-tree-*` MAY appear only in local development fixtures and MUST NOT support a trusted product claim.
- **FR-032**: `sdp-trace` MUST distinguish a source-bound local release from an externally trusted production release.
- **FR-033**: Source-bound local release finalization MUST verify that the selected source reference contains every manifest artifact path with matching digest before `source_commit_artifacts_verified` can be assessed true.
- **FR-034**: External production trust MUST name the selected trust profile, such as Sigstore/Rekor or customer PKI/private equivalent, and record transparency, audit, timestamp, approval, protected ref, and workflow identity evidence as assessed or `not_assessed`.
- **FR-035**: `trusted_contract_release: true` MUST be derivable from explicit assessed proof states; it MUST NOT be set by schema validity, local DSSE, or private key possession alone.
- **FR-036**: Pilot run-cards MUST be treated as evidence recipes, not completed observed behavior evidence; a run-card alone MUST NOT move a model, harness, stack, or customer pilot row out of `not_assessed`.
- **FR-037**: Pilot matrices MUST expose evidence state, reason code, artifact reference, gap reason, and next required evidence for each row; supported, ready, compatible, pass, fail, or warn wording is allowed only as an externally produced verdict input with producer, origin, and policy reference when available.
- **FR-038**: OpenCode model rows for MiniMax, Kimi, and GLM MUST remain `not_assessed` until committed sanitized run artifacts record observed model identity, harness identity, prompt, source reference, produced artifacts, and export limitations.
- **FR-039**: Kotlin+Bazel pilot artifacts MUST distinguish evidence-design placeholders from real stack proof; synthetic or placeholder fixtures MUST keep real Kotlin+Bazel run behavior `not_assessed` with a reason code until committed run evidence exists.
- **FR-040**: The first real product proof MUST target OpenCode + MiniMax + Kotlin+Bazel before broad pilot compatibility claims are made.
- **FR-041**: A Block 06 reference runner MAY shell out to OpenCode and Bazel/Bazelisk but MUST NOT add OpenCode, MiniMax, Kotlin, Bazel, Bazelisk, or provider SDKs as repository dependencies.
- **FR-042**: The Block 06 runner MUST record separate proof states for OpenCode availability, MiniMax model listing, MiniMax access verification, target-based Kotlin+Bazel identification, OpenCode+MiniMax run completion, Bazel command execution, `sdp-trace` package validation, and sanitized report commitment.
- **FR-043**: A packaging-only fixture MUST NOT satisfy Block 06 product proof. Completion requires a committed sanitized report from a real OpenCode + MiniMax run with assessed Bazel command execution or the block remains incomplete.
- **FR-044**: Missing MiniMax credentials, missing model access, missing OpenCode export, missing Bazel/Bazelisk, missing operator-approved Bazel command, or unavailable source access MUST be recorded as explicit incomplete or `not_assessed` proof states.
- **FR-045**: The Block 06 proof package MUST include machine-readable proof states, evidence events, provenance records, observations, metric stream, trace snapshot, assessment input, redaction note, and tested-on report.
- **FR-046**: Raw provider output, raw customer source, credentials, private prompts, and private logs MUST NOT be committed as Block 06 evidence.
- **FR-047**: Matrix updates based on Block 06 MUST cite the committed proof package and MUST NOT generalize one observed slice into broad support, readiness, compatibility, pass/fail, or degradation claims.
- **FR-048**: `sdp-trace` MUST support a flight-recorder event model that records source baseline, task/expectation lock, model/harness identity, command events, file mutation evidence, test evidence, redaction state, witness anchor, and run closure as ordered event-chain data.
- **FR-049**: Flight-recorder events MUST use deterministic canonicalization, declared schema version, declared hash algorithm, `prev_event_hash`, and `event_hash`; changing, deleting, or reordering committed events MUST be verifier-detectable.
- **FR-050**: A local flight-recorder chain MUST NOT be treated as accountability or audit-grade evidence unless it is bound to a witness anchor outside the mutable run artifact.
- **FR-051**: Flight-recorder verifier profiles MUST distinguish local structural chain validity, witnessed run validity, and forensic usefulness; each profile MUST emit `pass`, `fail`, `not_assessed`, or `cannot_verify` states with machine-readable reasons.
- **FR-052**: Mid-flight recorder attachment MUST create a visible `not_assessed` boundary for pre-attachment history; no profile may infer provenance for activity before attachment.
- **FR-053**: Requirement, task, prompt, or expectation changes after run start MUST be represented as superseding events that preserve the original locked event and link to the replacement event.
- **FR-054**: Flight-recorder evidence retention MUST distinguish `digest_only`, `sanitized_excerpt`, `encrypted_raw_ref`, `external_artifact_ref`, and `not_assessed`; profiles that require forensic reconstruction MUST reject insufficient retention for critical events.
- **FR-055**: Flight-recorder redaction MUST be verifier-visible and MUST distinguish safe redaction, sealed raw evidence, unresolved redaction, and unverifiable redaction; unresolved redaction MUST fail profiles that require forensic or accountability evidence.
- **FR-056**: Flight-recorder query surfaces MUST expose run summary, provenance, late-attach gaps, requirement supersession timeline, command timeline, file mutations, test evidence, redaction issues, and witness state without producing policy verdicts.
- **FR-057**: Advisory gate contracts MAY declare required runs separately from required evidence; missing required runs MUST produce `missing_telemetry` or `cannot_verify`, not `pass`.
- **FR-058**: Gate output MUST distinguish observation, advisory CI, and future protected profiles, and MUST keep protected profiles `cannot_verify` until signed checkpoint evidence and an external policy consumer exist.
- **FR-059**: CI witness binding MUST compare available repository, ref, commit, run id, and artifact digest data against the current gate input; mismatches MUST produce deterministic `fail` or `cannot_verify` reasons.
- **FR-060**: `policy_override_requested` trace events MUST be visible as override records and MUST NOT convert missing evidence to pass or upgrade audit-grade trust.
- **FR-061**: Gate explain output MUST provide deterministic human-readable reasons and next actions for missing telemetry, cannot-verify witness, stale witness, source mismatch, and override-present states.
- **FR-062**: Gate preview output MUST be read-only, deterministic, and explicit about gate-relevant fields, selected mode, trust cap, required runs, evidence ids, and local witness inspectability.
- **FR-063**: Gate, explain, and preview output MUST avoid raw command arguments, stdout/stderr bodies, prompts, source snippets, credentials, OIDC request tokens, model responses, and other secret-like values.
- **FR-064**: `sdp-trace` MUST NOT turn advisory gate facts into native merge, release, readiness, degradation, override approval, or risk-acceptance decisions.
- **FR-065**: `sdp-trace` MUST support a signed checkpoint artifact that binds a run id, run nonce, event chain head, event count, source snapshot digest/state, task hash, contract digest, checkpoint sequence, signing profile, payload digest, detached signature, and signer identity.
- **FR-066**: Signed checkpoint verification MUST recompute the canonical payload digest and detached signature before using any checkpoint binding as evidence.
- **FR-067**: Signed checkpoint verification MUST fail when a checkpoint is replayed against a different run id, nonce, source snapshot, task hash, contract digest, event count, or event chain head.
- **FR-068**: Signed checkpoint verification MUST expose monotonic checkpoint sequence checks and fail duplicate, missing, or descending sequence evidence.
- **FR-069**: Signed checkpoint signer authority MUST be checked against an explicit trusted-checkpoint policy when supplied; missing policy MUST remain `not_assessed`, and policy mismatch MUST fail.
- **FR-070**: Local development checkpoint signatures MUST NOT upgrade protected gate, audit-grade, release, readiness, degradation, override approval, or risk-acceptance state.

### Key Entities

- **Accountability Record**: Human-held DRI, approver, escalation path, risk owner, authority scope, accountability claim, approval reference, and line-of-defense metadata attached to accountable artifacts.
- **Risk Classification**: Observed autonomy, observed impact, classification source, and externally declared oversight metadata used by downstream policy consumers.
- **Contract Manifest**: Versioned list of contract artifacts and SHA-256 digests covering schemas, docs, validation scripts, fixtures, source commit, approval refs, compatibility notes, and previous manifest digest when available.
- **Contract Release Verification**: Evidence record that manifest digests, signature profile, signer identity policy, freshness or rollback status, and signature verification status were checked.
- **Trusted Identity Policy**: Contract for authorized release signer identity, protected source ref, workflow identity, release captain, required approval evidence, and private-equivalent verification profile where applicable.
- **Evidence Event**: One observed proof item from a source such as CI, command output, file inspection, review, scanner, harness log, model output, or manual sign-off.
- **Provenance Record**: Actor, model, harness, tool, command, artifact, timestamp, and source chain metadata when available.
- **Observation**: A dated statement about process state or behavior, backed by one or more evidence events.
- **Metric Sample**: A numeric, boolean, categorical, or count value measured for a dimension set and time window.
- **Metric Stream**: Ordered metric samples over time for the same metric name and comparable dimensions.
- **Trace Snapshot**: A point-in-time graph linking specs, plans, tasks, changes, evidence, observations, external verdict inputs, and decisions.
- **Assessment Input**: A package of trace artifacts prepared for a policy engine such as `sdp-gate`.
- **Pilot Run-Card**: A repeatable harness/model/stack assessment recipe with prompt, expected artifacts, provenance capture, validation, and `not_assessed` rules.
- **Signed Checkpoint**: Detached-signature artifact that binds a flight-recorder run chain head to run, source, task, contract, nonce, and sequence context for replay-resistant verification.
- **Trusted Checkpoint Policy**: Portable policy that names allowed checkpoint signer identities and the authority boundary needed to treat a checkpoint as local signed, CI signed, or externally witnessed evidence.
- **E2E Pilot Proof Package**: A sanitized artifact set produced from a real external tool run, containing evidence events, provenance records, observations, metric stream, trace snapshot, assessment input, redaction note, tested-on report, and explicit proof states.
- **External Verdict Input**: A verdict, score, evidence-strength assertion, or decision produced outside `sdp-trace` and recorded as evidence with producer, policy reference, artifact reference, and origin.
- **Flight Recorder Event**: An ordered event in a recorder run, with declared schema version, canonical payload digest, previous event hash, event hash, recorder identity, redaction state, and optional witness reference.
- **Witness Anchor**: A record outside the mutable run artifact that binds run id, source baseline, task hash, recorder version, and chain head so local chain replacement can be detected.
- **Requirement Supersession Event**: An append-only event that changes task or expectation scope by referencing an earlier locked event; it never edits or replaces the earlier event.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A repository observer can find the canonical feature spec, plan, and tasks under `specs/001-sdp-trace-time-series-evidence-substrate/`.
- **SC-002**: At least one contract document explicitly separates `sdp-trace` data ownership from `sdp-gate` policy ownership.
- **SC-003**: The implementation plan identifies every current Beads task as secondary tracking for a SpecKit task or artifact.
- **SC-004**: No new public artifact claims `sdp-trace` decides degradation, readiness, gate pass/fail, or override; CTO-facing docs phrase the answer as evidence-backed movement data unless an external verdict is explicitly named.
- **SC-005**: The pilot plan contains explicit run-card coverage for OpenCode+MiniMax, OpenCode+Kimi, OpenCode+GLM, and Kotlin+Bazel.
- **SC-006**: The schema validation plan documents Draft 2020-12, a pinned validator command, exclusions for ignored/local outputs, and validation of committed `sdp-trace` JSON artifacts.
- **SC-007**: Self-trace examples include sanitized artifact references, SHA-256 digests where artifacts are committed, and explicit `integrity_status` for unverified external references.
- **SC-008**: The boundary contract and data model define how `sdp-gate` inherits schema versions and how breaking changes are signaled.
- **SC-009**: Accountability examples identify human-held DRI, approver, risk owner, escalation path, approval reference, and line of defense; a negative example with AI as sole accountable owner fails validation.
- **SC-010**: A contract manifest example validates and includes SHA-256 digests for schemas, docs, validation scripts, fixtures, source commit, compatibility notes, and previous manifest digest when available.
- **SC-011**: A contract release verification example records `sdp-trace-signature/sigstore-dsse-keyless-v1`, manifest digest status, signer identity policy, signature status, and freshness or rollback status.
- **SC-012**: A negative modified-contract fixture fails manifest verification for the intended digest mismatch reason.
- **SC-013**: A trusted identity policy example validates and a signer identity mismatch fails trusted-release verification.
- **SC-014**: A metric stream example with mixed assessed and `not_assessed` samples exposes stream-level `assessment_state: partial` or `not_assessed`.
- **SC-015**: Contract scaffolding evidence includes one release verification result for the target signing profile shape or approved private equivalent; otherwise signing remains `not_assessed` and the block cannot claim contract scaffolding complete.
- **SC-016**: A committed `examples/self-trace/assessment-input.json` validates and links this feature's spec, plan, tasks, changed files, commands, Socratic review artifacts, accountability, observations, and metric streams without native pass/fail/degradation decisions.
- **SC-017**: A repository observer can run a documented self-trace validation command from a fresh checkout and reproduce the proof state for `schema_valid`, `digest_verified`, `locally_attested`, `externally_attested`, and `production_release_verified`.
- **SC-018**: Customer pilot run-cards remain blocked until SC-016 passes; any earlier pilot readiness claim is invalid.
- **SC-019**: A Block 04 release-finalization spec and Socratic dialogue let a reviewer distinguish source-bound local finalization from external production trust before implementation starts.
- **SC-020**: A future source-bound local release can assess `source_commit_artifacts_verified: true` only after the manifest subject artifacts are committed and the source-bound finalization guard verifies that the selected 40-character source commit contains those artifact paths with matching digests. DSSE envelope and self-attestation result regeneration are separate release-proof freshness steps and do not change the source artifact set.
- **SC-021**: Block 05 run-cards cover OpenCode+MiniMax, OpenCode+Kimi, OpenCode+GLM, Superpowers-style harnesses, `gsd`, `gsd2`, Oh My OpenAgent, and Kotlin+Bazel with required artifacts, provenance fields, `unbacked_claim` capture, validation, and `not_assessed` rules.
- **SC-022**: Harness and model compatibility matrices show evidence state, reason code, artifact reference, gap reason, and next required evidence; no row records observed behavior unless a committed sanitized run artifact or evidence summary supports it.
- **SC-023**: The customer pilot evidence package outline defines safe customer inputs, sanitized outputs, package shape, validation commands, review checkpoints, and residual `not_assessed` reporting without raw customer data or pilot readiness claims.
- **SC-024**: Block 06 spec and Socratic artifacts define OpenCode + MiniMax + Kotlin+Bazel as the first real product proof and reject packaging-only proof.
- **SC-025**: A repository observer can run a documented Block 06 command shape without installing repository dependencies for OpenCode, MiniMax, Kotlin, Bazel, or Bazelisk.
- **SC-026**: The Block 06 runner emits separate proof states for OpenCode availability, MiniMax model listing, MiniMax access verification, target-based Kotlin+Bazel identification, OpenCode+MiniMax run completion, Bazel command execution, package validation, and sanitized report commitment.
- **SC-027**: A committed sanitized Block 06 report names the exact tested-on OpenCode version, MiniMax model id, source reference, Bazel target, Bazel command, command summary, validation result, and residual `not_assessed` fields.
- **SC-028**: `scripts/validate-e2e-pilot-package.sh` rejects a committed proof package that omits required proof states in `evidence/proof-states.json` or includes raw output.
- **SC-029**: Model and harness matrices cite the Block 06 proof package only for the exact observed OpenCode + MiniMax + Kotlin+Bazel slice.
- **SC-030**: If a real OpenCode + MiniMax run cannot be completed, Block 06 remains open and the incomplete package cannot be used as product proof.
- **SC-031**: A Block 09 design, Socratic review synthesis, and implementation plan define the flight-recorder trust kernel before any new external demo repo work starts.
- **SC-032**: Flight-recorder fixtures prove that event mutation, event deletion, event reordering, missing witness, witness mismatch, late attachment, task rewrite, and unresolved redaction are detected or explicitly reported by verifier states.
- **SC-033**: A reviewer can run documented flight-recorder query commands against committed fixtures and identify source baseline, task, model/harness identity, command timeline, file mutations, test evidence, redaction state, witness state, and `not_assessed` gaps without reading raw JSONL manually.
- **SC-034**: Block 14 fixtures prove that absent required runs, unmatched runs, stale or mismatched CI witnesses, and protected-future requirements produce deterministic non-pass states.
- **SC-035**: A reviewer can run `gate explain` and `gate preview` against committed Block 14 fixtures and see required-run gaps, witness binding state, override records, trust cap, and next actions without reading raw JSON manually.
- **SC-036**: Safety-sensitive gate, explain, and preview tests prove that secret-like command arguments, stdout/stderr bodies, prompts, source snippets, credentials, OIDC request tokens, and model responses are not printed or persisted.

## Assumptions

- `sdp-gate` will consume `sdp-trace` artifacts but will live in a separate product/repository boundary.
- Beads remains useful for local work tracking, but Beads is not a product dependency and is not the repo observer's source of truth.
- The initial implementation may be schema and documentation heavy before adding tiny validation tools.
- Customer pilot artifacts may need sanitization before committing summaries to the repository.
- Current schemas already declare JSON Schema Draft 2020-12; this feature standardizes that draft for new schemas unless a future major version changes it.
- Until `sdp-trace` reaches v1.0, schema changes may still be breaking, but every breaking change must update examples, compatibility notes, and downstream `sdp-gate` handoff documentation.
- Public examples may use synthetic human-held roles for DRI and approver fields; customer pilots must map those fields to the customer's accepted identity or approval system.
- Public Sigstore/Rekor verification is the target profile, but private or air-gapped pilots may use an approved equivalent if the manifest, DSSE envelope, identity policy, and verification result remain explicit.
