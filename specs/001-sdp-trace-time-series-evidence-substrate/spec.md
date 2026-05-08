# Feature Specification: sdp-trace Time-Series Evidence Substrate

**Feature Branch**: `001-sdp-trace-time-series-evidence-substrate`
**Created**: 2026-04-30
**Status**: Draft
**Input**: User description: "`sdp-trace` must be the repository-observable evidence and trace substrate. Planning may come from SpecKit or another workflow; Beads is only discipline support. The technical executive question is whether the delivery process is moving toward degradation over time; external policy consumers apply policy on top of `sdp-trace` contracts. A CEO must also be able to see which human is accountable for each stage and whether the contract itself was changed under a trusted release process."

## Non-Negotiable Self-Proof Rule

`sdp-trace` MUST be its own first real consumer.

No customer pilot claim, external compatibility verdict, or "trusted contract foundation" claim is acceptable until this repository can validate a committed self-trace package describing the development of this feature itself: spec, plan, tasks, changed artifacts, evidence, provenance, accountability, observations, metric movement, `not_assessed` gaps, Socratic review results, and contract release verification.

If self-trace or self-attestation is missing, the correct product state is not "trusted"; it is `not_assessed` with a concrete reason.

## User Scenarios & Testing

### User Story 1 - technical executive Reviews Process Movement (Priority: P1)

A technical executive reviewing an AI-assisted delivery pilot can inspect accumulated `sdp-trace` artifacts and see what changed over time: evidence quality, scope discipline, correctness signals, review signals, AI behavior, and stack/harness/model slices.

**Why this priority**: This is the product reason for the repository. If the technical executive cannot understand whether the process is improving, stable, degrading, or not assessable, the substrate is not useful.

**Boundary**: `sdp-trace` answers the technical executive question with inspectable movement data: prior/current values, deltas, dimensions, evidence coverage, and `not_assessed` gaps. It does not answer with "yes, we are degrading" or "no, we are not degrading." A yes/no verdict or threshold interpretation belongs to an external policy consumer.

**Independent Test**: A reviewer opens this SpecKit package and the generated examples, then confirms every process signal is backed by evidence or explicitly marked `not_assessed`.

**Acceptance Scenarios**:

1. **Given** a trace snapshot with complete evidence references, **When** a reviewer inspects the metric samples, **Then** each sample links to inspectable evidence or provenance.
2. **Given** missing build/test evidence, **When** an observation is recorded, **Then** the affected sample is marked `not_assessed` instead of inferred.
3. **Given** multiple observations over time, **When** a current window is compared with a previous window, **Then** `sdp-trace` records the prior value, current value, delta, units, dimensions, and evidence coverage without producing a policy verdict.

---

### User Story 2 - External Policy Consumers Use Trace Contracts (Priority: P1)

An external policy consumer implementer can consume `sdp-trace` artifacts as policy inputs without `sdp-trace` deciding pass/fail, readiness, degradation, or override outcomes.

**Why this priority**: if the boundary is vague, recorder and policy layers will duplicate policy logic and confuse users.

**Independent Test**: Read the boundary contract and verify it names `sdp-trace` ownership and external policy consumer ownership separately.

**Acceptance Scenarios**:

1. **Given** an evidence bundle and metric stream, **When** an external policy consumer applies a policy, **Then** the policy decision is external to `sdp-trace`.
2. **Given** an external gate verdict is recorded as evidence, **When** it appears in a trace, **Then** it is represented as an observed verdict input, not as a decision made by `sdp-trace`.
3. **Given** an evidence event carries an externally produced quality, strength, or verdict assertion, **When** it is recorded by `sdp-trace`, **Then** the assertion records its producer, policy reference when available, and external origin rather than becoming a native `sdp-trace` assessment.

---

### User Story 5 - CEO/CIO Verifies Accountability and Contract Integrity (Priority: P1)

A CEO or CIO can inspect a contract release or assessment package and identify the human accountable owner, the approval reference, the escalation path, the risk owner, and whether the checked-out contract matches a signed release manifest.

**Why this priority**: A process that says "the AI did it" is not governable. A schema-valid artifact is also not trustworthy if a person or model can quietly simplify the contract and still call the result valid.

**Boundary**: `sdp-trace` records accountability, release integrity, and verification status. It does not decide whether the organization accepts residual risk; that policy decision belongs to external policy consumer, management, or another external governance process.

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

### User Story 4 - Repository Observer Finds Planning Evidence (Priority: P2)

A repository observer can understand current scope and proof by reading committed planning artifacts without needing a private planning runtime.

**Why this priority**: external trackers and local workflow tools are discipline support. The repository-facing plan and evidence must be inspectable from committed files.

**Independent Test**: A reviewer can start from `/specs/001-sdp-trace-time-series-evidence-substrate/spec.md`, follow `plan.md` and `tasks.md`, and map task status to committed artifacts.

**Acceptance Scenarios**:

1. **Given** a fresh clone without private workflow context loaded, **When** a reviewer opens `specs/001-sdp-trace-time-series-evidence-substrate/`, **Then** they can understand the feature, plan, tasks, contract, and evidence expectations.
2. **Given** external or local task trackers exist, **When** they are inspected, **Then** they reference committed planning artifacts as secondary tracking, not the other way around.

## Edge Cases

- A source system cannot expose raw logs: record the missing field as `not_assessed` and keep the run usable.
- A model or harness reports its own identity inconsistently: preserve observed identity and add an `unbacked_claim` item.
- A PR/MR does not exist: evidence events must support local branch, commit, file, command, or manual sources without PR-only assumptions.
- Customer data cannot be committed: examples and summaries must be sanitized while preserving artifact references, hashes, or redaction notes.
- A consuming policy wants thresholds: thresholds belong to external policy consumer or another policy engine, not to `sdp-trace`.
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
- **FR-003**: `sdp-trace` MUST state that external policy consumers can use `sdp-trace` contracts while owning policy evaluation.
- **FR-004**: Every metric sample MUST reference inspectable evidence or be marked `not_assessed`.
- **FR-005**: The metric catalog MUST avoid opaque aggregate health scores.
- **FR-006**: The contract MUST support moving time windows without requiring a fixed baseline.
- **FR-007**: The contract MUST support dimensions for repository, scope, team, harness, model family, model version when available, stack, build system, and time window.
- **FR-008**: The pilot run-card set MUST include OpenCode with MiniMax, Kimi, and GLM model slices.
- **FR-009**: The pilot run-card set MUST include Superpowers-style, `gsd`, `gsd2`, and Oh My OpenAgent harness rows with observed evidence references or explicitly `not_assessed` state and reason codes.
- **FR-010**: The JVM pilot path MUST define the Kotlin+Bazel evidence path; Java+Bazel, Kotlin+Gradle, and synthetic Kotlin+Bazel placeholders are not sufficient proof of real Kotlin+Bazel behavior.
- **FR-011**: Public docs MUST use SpecKit-aligned terms first: spec, plan, task, evidence, gate, decision, trace, provenance.
- **FR-012**: Public docs MUST not imply dependency on Beads, Operator Mode, agentloop, OpenCode, GitHub, Claude, Codex, or any specific harness runtime.
- **FR-013**: Schema and example artifacts MUST be machine-checkable by documented commands.
- **FR-014**: Static compatibility matrices are retired. Evidence for harnesses,
  models, languages, and build tools MUST be recorded at the exact run or
  package level with observed state, artifact reference, gap reason, and next
  required evidence; broad support, readiness, or compatibility claims require
  an external verdict input.
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
- **FR-058**: Gate output MUST distinguish observation, advisory CI, and protected profiles, and MUST keep protected profiles non-pass unless signed checkpoint evidence, trusted signer authority, required evidence, required runs, CI witness binding, and external policy-consumer ownership are explicit.
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
- **FR-071**: Protected gate evaluation MUST be selected explicitly and MUST NOT be the default `sdp-trace gate` behavior.
- **FR-072**: Protected gate evaluation MUST fail closed or return `cannot_verify` when required runs, required evidence, signed checkpoint evidence, trusted-checkpoint policy, signer authority, or CI witness binding is missing or invalid.
- **FR-073**: Protected gate output MUST distinguish verifier-derived `protected_gate` facts from native merge, release, readiness, degradation, override approval, and risk-acceptance decisions.
- **FR-074**: Local observed evidence, local signed checkpoints, and untrusted checkpoint shapes MUST NOT satisfy protected gate trust scope.
- **FR-075**: CI signed protected gate trust MUST require signer authority evidence plus source, ref, commit, run id, and artifact digest binding to the selected run context.
- **FR-076**: Protected profile exit behavior MUST be deterministic: protected pass exits `0`, protected fail exits `1`, protected `cannot_verify` exits `3`, and usage errors exit `2`.
- **FR-077**: Protected gate, explain, and preview output MUST avoid raw command arguments, stdout/stderr bodies, prompts, source snippets, credentials, OIDC request tokens, model responses, and other secret-like values.
- **FR-078**: Protected profile output MUST use a new Block 16 gate-result schema version while preserving `gate explain` read compatibility for Block 14 gate-result artifacts.
- **FR-079**: Managed harness evaluation MUST be selected explicitly and MUST NOT be required for observation-mode users or unmanaged harnesses.
- **FR-080**: Managed harness evaluation MUST fail closed or return `cannot_verify` when required managed policy, adapter registry, registered wrapper enrollment, authorized adapter identity, adapter capability evidence, required harness/tool/file/test events, or managed witness binding is missing or invalid.
- **FR-081**: Managed harness output MUST distinguish verifier-derived `managed_harness_assessment` facts from native merge, release, readiness, degradation, override approval, and risk-acceptance decisions.
- **FR-082**: Agent-reported, local-observed, late-attached, self-claimed adapter, unauthorized adapter, and unregistered wrapper evidence MUST NOT satisfy managed harness profile trust scope.
- **FR-083**: Managed harness trust MUST require a registered wrapper or authorized adapter boundary before child harness execution, policy-authorized adapter capabilities, required event coverage, non-agent-reported executed-test provenance, and witness binding to source, run, policy, adapter or wrapper identity, and artifacts.
- **FR-084**: Managed harness profile exit behavior MUST be deterministic: managed pass exits `0`, managed fail exits `1`, managed `cannot_verify` exits `3`, and usage errors exit `2`.
- **FR-085**: Managed harness assess, explain, and preview output MUST avoid raw command arguments, stdout/stderr bodies, prompts, source snippets, credentials, OIDC request tokens, adapter secrets, gateway tokens, model responses, and checkpoint key material.
- **FR-086**: Managed harness profile output MUST use a new Block 17 assessment-result schema version while preserving `gate explain` read compatibility for Block 14 and Block 16 gate-result artifacts.
- **FR-087**: Redaction and retention assessment profiles MUST be selected explicitly for forensic assessment and MUST NOT enable raw payload persistence by default.
- **FR-088**: The built-in safe default retention mode MUST redact known secret-like payloads before persistence, retain command/output bodies as digests unless a policy permits safer retention, and emit a stable versioned policy id and digest.
- **FR-089**: Redaction policy artifacts MUST define rule ids, allowed FR-054 retention modes, redaction actions, forbidden committed-artifact persistence classes, verifiable authority identity, event-family profile mappings, critical event family classifications, withholding audit fields, and unresolved-redaction profile impact.
- **FR-090**: Flight-recorder events and run manifests MUST expose verifier-readable redaction policy refs, redaction rule refs, redaction authority, input/redacted payload digests, FR-054 retention mode, retention lifecycle, forensic importance, and raw-reference state when applicable.
- **FR-091**: Forensic retention assessment MUST reject digest-only critical event evidence and emit an explanatory cap to a lower retention mode unless the selected redaction policy explicitly classifies that event family as non-critical for the selected profile.
- **FR-092**: Encrypted or external raw references MUST include SHA-256-or-stronger digest binding, reference type, machine-readable access state, access verification time, key custody state where applicable, retention lifecycle, and unavailable reason when access or verification is missing.
- **FR-093**: Missing redaction policy, unresolved redaction, unverifiable raw reference access, missing key custody state, or contradictory raw reference digest MUST produce deterministic `fail` or `cannot_verify` facts for forensic profiles.
- **FR-094**: Redaction, retention, forensic assess, explain, and preview output MUST avoid raw command arguments, stdout/stderr bodies, prompts, source snippets, credentials, OIDC request tokens, adapter secrets, gateway tokens, model responses, key material, and unsafe raw-reference access notes.
- **FR-095**: Adapter event capture MUST be defined by a stable, provider-neutral contract for `run_started`, `task_locked`, `task_superseded`, `tool_call`, `command_started`, `file_mutation`, `model_call_observed`, `test_observed`, and `run_closed`.
- **FR-096**: Adapter events MUST record producer identity, adapter identity, provenance scope, capture state, run id, run nonce, source baseline or commit/tree digest when available, binding mode, same-chain or adapter-bundle digest linkage, correlation refs, redaction/retention metadata, and a deterministic reason when any required field is `not_assessed` or unavailable.
- **FR-097**: Adapter-reported model identity MUST remain `harness_observed` or `agent_reported` unless gateway evidence exists and is bound to the selected run; `sdp-trace` MUST NOT infer gateway authority from model name, harness name, adapter id text, or prose.
- **FR-098**: Test observation provenance MUST distinguish `ci_executed`, `wrapper_executed`, `harness_observed`, `agent_reported`, and `cannot_verify`; agent-reported claims and harness-observed intent MUST NOT appear as executed test evidence without CI or registered wrapper/tool execution evidence.
- **FR-099**: Source, VCS, PR/MR, and review references MUST use provider-neutral shapes and MUST NOT require GitHub, GitLab, Bitbucket, Gerrit, or any specific Git host to be present.
- **FR-100**: Missing adapter events, unsupported observer capabilities, gateway `not_integrated`, late adapter attach, and unavailable prompt/model-response raw capture MUST be explicit machine-readable facts and MUST NOT be collapsed into pass/fail policy decisions.
- **FR-101**: Prompt and raw model-response capture MUST remain unavailable by default; richer retention is allowed only when an explicit Block 18 redaction/retention policy permits it and the retained evidence is verifier-visible.
- **FR-102**: Adapter event examples and product code MUST NOT encode OpenCode, GSD, Bazel, Kotlin, GitHub, Claude, Codex, or any demo-specific harness/model/build-system name as a product concept.
- **FR-103**: Query-facing output for adapter capture MUST expose task drift, task supersession counts, unverified task expansion indicators, unverified claims, missing adapter events, unsupported observers, gateway `not_integrated`, retention limits, and capture-depth caps without claiming scope-creep judgment or forensic completeness beyond available evidence.
- **FR-104**: Adapter event assess, query, preview, and explain output MUST avoid raw command arguments, stdout/stderr bodies, prompts, source snippets, tool-call input/output bodies, adapter configuration, gateway evidence refs, credentials, OIDC request tokens, adapter secrets, gateway tokens, model request/response payloads, PR tokens, authenticated provider URLs, and raw review bodies.
- **FR-105**: Forensics query packs MUST be versioned, explicitly selected, read-only bundles of query rows over existing run, verifier, forensic-retention, and adapter-capture facts; they MUST NOT emit native incident, audit, release, merge, readiness, degradation, override, legal, or risk-acceptance decisions.
- **FR-106**: `forensics-basic-v1` MUST expose safe rows for run summary, chain and witness state, retention and redaction issues, capture depth, command and test timeline using opaque command/test identifiers only, file mutation evidence, task supersession, and unverified claims without requiring reviewers to inspect raw JSONL manually.
- **FR-107**: Forensics query-pack rows MUST preserve `not_assessed`, `cannot_verify`, `missing_telemetry`, `unsupported`, `not_integrated`, and `retention_limited` states with deterministic propagation from upstream Block 09, Block 18, and Block 19 facts, source condition refs where available, closed-enum `evidence_family`, and safe `evidence_gap` descriptions; missing Block 18 or Block 19 facts MUST remain visible instead of disappearing from query output.
- **FR-108**: Digest-only evidence in forensics query-pack output MUST be rendered as existence evidence only and MUST NOT be described as reconstructable evidence unless Block 18 retention facts prove reconstructability for the selected event family; query rows MUST expose `reconstructable: false` when evidence exists but reconstruction is capped.
- **FR-109**: Forensics query-pack result artifacts MUST record schema version, query pack id/version, selected run identity when available, required and optional input artifact SHA-256 digests or path-redacted provider-neutral identifiers, grouped query rows, source fact refs, source condition refs where available, deterministic query-scoped row ids, closed-enum evidence families, row-state derivation, and verified output-safety metadata.
- **FR-110**: Forensics query-pack explain output MUST render the JSON result artifact in stable query-name and row-id order and MUST NOT add claims, summaries, conclusions, hidden severity ordering, ANSI color state, indentation state, or omitted-section state absent from the machine-readable result.
- **FR-111**: Forensics query-pack JSON and explain output MUST avoid raw command arguments, command names, executable paths, script paths, test identifiers unless public-catalog safe, stdout/stderr bodies, prompts, source snippets, tool-call input/output bodies, adapter configuration, gateway evidence refs, credentials, OIDC request tokens, adapter secrets, gateway tokens, PR tokens, authenticated provider URLs, raw model request/response payloads, raw review bodies, unsafe raw-reference access notes, and key material.
- **FR-112**: Cross-repository posture exports MUST aggregate evidence movement facts by repository, team, service, harness, change type, evidence family, row evidence state, input trust state, and time window without emitting native degradation, improvement, readiness, pass/fail, rank, grade, health-score, override approval, incident, audit, legal, release, merge, or risk-acceptance decisions.
- **FR-113**: Every cross-repository metric row MUST include metric id/version, numerator, denominator, unit, dimensions, time window, source artifact digest refs or path-redacted input artifact ids, `not_assessed` count, input trust state summary, and deterministic aggregation rule; ratios MUST NOT replace raw counts.
- **FR-114**: Cross-repository movement rows MUST expose metric id/version, deterministic dimension key, current metric row ref, previous metric row ref, current value, previous value, signed delta, closed-enum comparison basis, comparable boolean, and closed-enum non-comparable reason when metric id/version, dimension key, denominator basis, or input trust rules do not match.
- **FR-115**: Cross-repository aggregation MUST refuse or explicitly mark stale, malformed, digest-mismatched, untrusted, missing required, missing optional, unsupported, unsafe-label, and non-comparable inputs before aggregation; such inputs MUST NOT disappear from denominators or be treated as positive evidence.
- **FR-116**: Cross-repository export JSON and explain output MUST avoid raw command arguments, command names, executable paths, script paths, unsafe test identifiers, stdout/stderr bodies, prompts, source snippets, tool-call input/output bodies, adapter configuration, gateway evidence refs, credentials, OIDC request tokens, adapter secrets, gateway tokens, PR tokens, authenticated provider URLs, raw model request/response payloads, raw review bodies, unsafe raw-reference access notes, private filesystem paths, and raw personal identifiers not declared safe by upstream artifacts.
- **FR-117**: Cross-repository exports MUST validate safe labels, path-redacted artifact ids, refusal reason codes, non-comparable reason codes, and dimension exposure policy inside Block 21; upstream safe declarations and external verdict payloads are inputs, not authority.
- **FR-118**: Witness profiles MUST use a provider-neutral contract that declares profile id/version, identity source, signing boundary, freshness boundary, artifact binding, independence state, unsupported states, and output-safety classes before any provider-specific CI or enterprise witness input can upgrade trust scope.
- **FR-119**: GitLab CI, Buildkite, and customer-private PKI witness profiles MUST share verifier states and reason-code semantics for identity, signer authority, freshness, artifact binding, source binding, run binding, policy binding, independence, unsupported fields, and output safety.
- **FR-120**: No witness profile may upgrade a run to `ci_witnessed` or `external_witnessed` from environment variables, job logs, local signatures, committed JSON, or unchecked artifact paths alone.
- **FR-121**: Customer-private PKI witness verification MUST use public identity, authority-policy, payload-digest, and freshness inputs only; private key material, raw customer identity payloads, and unaudited customer infrastructure state MUST NOT be read, printed, persisted, or required.
- **FR-122**: Air-gapped witness guidance MUST distinguish offline-verifiable public-key, timestamp, and artifact-digest evidence from external checks that remain `not_assessed` or `cannot_verify`; committed witness JSON alone MUST NOT be treated as authority.
- **FR-123**: Witness profile JSON and explain output MUST avoid raw command arguments, unsafe command names or paths, stdout/stderr bodies, prompts, model responses, OIDC request tokens, JWT bodies, CI secrets, private keys, provider tokens, authenticated provider URLs, private filesystem paths, raw job logs, unsafe personal identifiers, free-text parser errors containing input content, and customer directory, LDAP, SAML, cloud, Vault, HSM, KMS, or PKI payloads.
- **FR-124**: CI artifact observation MUST distinguish `ci_uploaded`, `checked_in`, `local_generated`, `agent_reported`, `harness_observed`, `external_artifact_ref`, and `not_assessed` producer scopes before any artifact family can support a CI-backed proof claim.
- **FR-125**: CI artifact observation MUST expose artifact access states for `present`, `absent`, `partial`, `expired`, `inaccessible`, `malformed`, `unsafe`, `not_assessed`, and `cannot_verify`; missing or inaccessible required artifact families MUST NOT be treated as pass.
- **FR-126**: Checked-in JSON, prose, local debug artifacts, job logs, and unchecked artifact paths MUST NOT satisfy a selected `ci_uploaded` artifact-family requirement; they may be recorded only as lower-authority facts with explicit producer scope and reason codes.
- **FR-127**: CI artifact observation MUST evaluate required artifact families for run, report, witness, provenance, evidence, trace, artifact index, redaction scan, review, and change/branch CI evidence when those families are selected by the profile; unselected families remain `not_assessed`.
- **FR-128**: CI artifact observation MUST produce deterministic `fail` facts for contradictory source/run binding, digest mismatch, artifact-index self-reference, unsafe serialized output, or other verifier-visible contradiction, and deterministic `cannot_verify` facts for absent, expired, inaccessible, or incomplete required evidence.
- **FR-129**: Demo or pilot summaries MUST be able to cite CI artifact observation results for artifact truth; a manually curated happy-path report without matching observed artifact families remains `agent_reported`, `checked_in`, `not_assessed`, or `cannot_verify`.
- **FR-130**: CI artifact observation JSON and explain output MUST avoid raw CI logs, command bodies, stdout/stderr bodies, prompts, model responses, OIDC request tokens, JWT bodies, CI secrets, private keys, provider tokens, authenticated provider URLs, private filesystem paths, raw parser input, and unsafe personal or customer payloads.

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
- **Assessment Input**: A package of trace artifacts prepared for a policy engine such as external policy consumer.
- **Pilot Run-Card**: A repeatable harness/model/stack assessment recipe with prompt, expected artifacts, provenance capture, validation, and `not_assessed` rules.
- **Signed Checkpoint**: Detached-signature artifact that binds a flight-recorder run chain head to run, source, task, contract, nonce, and sequence context for replay-resistant verification.
- **Trusted Checkpoint Policy**: Portable policy that names allowed checkpoint signer identities and the authority boundary needed to treat a checkpoint as local signed, CI signed, or externally witnessed evidence.
- **Protected Gate Enforcement Profile**: Explicit gate profile that evaluates protected-use prerequisites as fail-closed verifier facts and deterministic exit behavior for an external CI or policy owner to enforce.
- **Managed Harness Policy**: Portable policy that names approved wrapper identities, adapter identities, capabilities, required managed events, suppression rules, and witness binding requirements for opt-in managed harness enforcement.
- **Adapter Registry**: Portable artifact declaring adapter id, harness id, version, identity state, signing or authority reference, allowed event types, capabilities, and deployment source.
- **Managed Harness Enforcement Profile**: Explicit gate profile that evaluates whether a selected run used an approved managed wrapper or adapter boundary and required managed telemetry without making a native policy decision.
- **Redaction Policy**: Portable artifact that names redaction rules, forbidden committed-artifact persistence classes, allowed FR-054 retention modes, redaction actions, verifiable authority, event-family profile mappings, critical event family classifications, withholding audit fields, and unresolved-redaction impact for verifier-visible privacy and forensic behavior.
- **Retention Mode**: Run-artifact state describing whether evidence is retained as `digest_only`, `sanitized_excerpt`, `encrypted_raw_ref`, `external_artifact_ref`, or `not_assessed` for selected event families.
- **Raw Reference**: SHA-256-or-stronger digest-bound pointer to encrypted or external raw evidence with reference type, access state, access verification time, key custody state when applicable, retention lifecycle, and unavailable reason when verification cannot be completed.
- **Forensic Retention Profile**: Explicit assessment profile that evaluates whether critical flight-recorder evidence is reconstructable under the selected redaction policy without making a legal, incident, readiness, or risk decision.
- **Adapter Event Contract**: Provider-neutral schema and semantics for harness or gateway adapters to report run lifecycle, task, tool, command, file mutation, model identity, test observation, and closure events without making adapter integration mandatory for observation-mode users.
- **Capture Depth**: Verifier-visible state describing which event families were observed, unsupported, missing, not integrated, not assessed, or cannot be verified for a selected run and profile.
- **Gateway Provenance Profile**: Explicit evidence state for model gateway observations that can remain `not_integrated`; it does not upgrade harness-reported or agent-reported model identity without bound gateway evidence.
- **Provider-Neutral Change Reference**: Portable source, VCS, PR/MR, or review reference with source ref, commit or tree digest, change or review id, artifact digest, producer, and `not_assessed` reason fields without binding the contract to a specific Git host.
- **Forensics Query Pack**: Versioned read-only bundle of query rows that assembles safe reviewer-facing forensic views from existing run, verifier, forensic-retention, and adapter-capture facts without making a policy decision.
- **Forensics Query Row**: Machine-readable row with deterministic query-scoped id, query name, evidence state, closed-enum evidence family, source fact refs, source condition refs where available, safe metadata, reconstructability flag, and evidence gap description when the row is `not_assessed` or `cannot_verify`.
- **Cross-Repository Posture Export**: Machine-readable aggregate artifact that groups safe evidence movement facts across selected repositories, services, teams, harnesses, change types, and time windows without deciding whether the movement is degradation or improvement.
- **Cross-Repository Metric Row**: Deterministic aggregate row with closed metric id/version, numerator, denominator, unit, dimensions, time window, source artifact digest refs, `not_assessed` count, input trust state summary, and aggregation rule.
- **Cross-Repository Movement Row**: Comparison row that records current value, previous value, delta, comparison basis, comparable state, and non-comparable reason for a selected metric/dimension/window pair.
- **Posture Signal Manifest**: Optional provider-neutral input artifact that carries closed, safe signals not present in Block 20 query-pack rows, such as witness scope, override presence, late attach, and contract-change markers.
- **Witness Profile**: Provider-neutral contract that defines how a witness source establishes identity, signer authority, freshness, artifact binding, source binding, run binding, policy binding, independence, unsupported states, and output safety before a run can be treated as CI-witnessed or externally witnessed.
- **Normalized Witness Result**: Machine-readable profile result with profile id/version, provider kind, producer, generated time, requested and established trust scope, verifier states, closed reason codes, path-redacted artifact ids, digests, and output-safety verification state.
- **Customer PKI Witness Profile**: Enterprise profile that validates public certificate or public key identity, signer authority policy, payload digest binding, and freshness evidence without reading private key material or depending on a live customer PKI service.
- **CI Artifact Observation**: Provider-neutral fact package that records selected source/run identity, artifact producer scope, artifact access state, required artifact-family states, binding states, index integrity, output safety, reasons, and next actions without making native policy decisions.
- **Artifact Family**: A selected evidence category such as run, report, witness, provenance, evidence, trace, artifact index, redaction scan, review, or PR/CI evidence whose presence, absence, access state, and binding state can be observed independently.
- **E2E Pilot Proof Package**: A sanitized artifact set produced from a real external tool run, containing evidence events, provenance records, observations, metric stream, trace snapshot, assessment input, redaction note, tested-on report, and explicit proof states.
- **External Verdict Input**: A verdict, score, evidence-strength assertion, or decision produced outside `sdp-trace` and recorded as evidence with producer, policy reference, artifact reference, and origin.
- **Flight Recorder Event**: An ordered event in a recorder run, with declared schema version, canonical payload digest, previous event hash, event hash, recorder identity, redaction state, and optional witness reference.
- **Witness Anchor**: A record outside the mutable run artifact that binds run id, source baseline, task hash, recorder version, and chain head so local chain replacement can be detected.
- **Requirement Supersession Event**: An append-only event that changes task or expectation scope by referencing an earlier locked event; it never edits or replaces the earlier event.

## Success Criteria

### Measurable Outcomes

- **SC-001**: A repository observer can find the canonical feature spec, plan, and tasks under `specs/001-sdp-trace-time-series-evidence-substrate/`.
- **SC-002**: At least one contract document explicitly separates `sdp-trace` data ownership from external policy ownership.
- **SC-003**: The implementation plan identifies local or external task trackers as secondary tracking for committed planning artifacts.
- **SC-004**: No new public artifact claims `sdp-trace` decides degradation, readiness, gate pass/fail, or override; technical executive-facing docs phrase the answer as evidence-backed movement data unless an external verdict is explicitly named.
- **SC-005**: The pilot plan contains explicit run-card coverage for OpenCode+MiniMax, OpenCode+Kimi, OpenCode+GLM, and Kotlin+Bazel.
- **SC-006**: The schema validation plan documents Draft 2020-12, a pinned validator command, exclusions for ignored/local outputs, and validation of committed `sdp-trace` JSON artifacts.
- **SC-007**: Self-trace examples include sanitized artifact references, SHA-256 digests where artifacts are committed, and explicit `integrity_status` for unverified external references.
- **SC-008**: The boundary contract and data model define how external consumers declare schema versions and how breaking changes are signaled.
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
- **SC-037**: Block 16 fixtures prove that missing checkpoint evidence, local-development checkpoint evidence, missing signer policy, signer mismatch, missing CI witness binding, absent or stale CI witness freshness evidence, and CI witness mismatch cannot pass the protected profile.
- **SC-038**: A reviewer can run protected gate, explain, and preview commands against committed Block 16 fixtures and see protected-profile input gaps, checkpoint state, signer authority state, witness binding state, override state, and deterministic next actions without reading raw JSON manually.
- **SC-039**: Safety-sensitive protected gate tests prove that secret-like command arguments, stdout/stderr bodies, prompts, source snippets, credentials, OIDC request tokens, model responses, and checkpoint key material are not printed or persisted.
- **SC-040**: Block 17 fixtures prove that unmanaged runs, late enrollment, post-hoc policy or registry artifacts, unauthorized adapter identity, adapter disconnect, missing adapter capabilities, missing required harness/tool/file/test telemetry, agent-reported-only test evidence, unverified suppression, missing managed witness binding, stale managed witness evidence, and managed witness mismatch cannot pass the managed harness profile.
- **SC-041**: A reviewer can run managed harness assess, explain, and preview commands against committed Block 17 fixtures and see managed-profile input gaps, wrapper enrollment state, adapter identity state, capability state, event coverage, suppression state, witness binding state, override state, and deterministic next actions without reading raw JSON manually.
- **SC-042**: Safety-sensitive managed harness tests prove that raw command arguments, stdout/stderr bodies, prompts, source snippets, credentials, OIDC request tokens, adapter secrets, gateway tokens, model responses, and checkpoint key material are not printed or persisted.
- **SC-043**: Block 18 fixtures prove that safe default redaction, digest-only critical evidence, sanitized excerpts, encrypted raw references, external artifact references, unresolved redaction, withholding, missing policy, weak digest, unavailable access, missing key custody, authority self-claim, and raw-reference digest mismatch produce deterministic forensic retention facts.
- **SC-044**: Block 19 fixtures prove that generic adapter events, same-chain and adapter-bundle binding, missing adapter events, unsupported observer capabilities, gateway `not_integrated`, late adapter attach, agent-reported test claims, harness-observed test correlation without execution proof, file mutation/source correlation, task supersession actor attribution, and provider-neutral PR/review refs produce deterministic capture-depth facts.
- **SC-045**: Safety-sensitive adapter event assess, query, preview, and explain tests prove that raw command arguments, stdout/stderr bodies, prompts, source snippets, tool-call input/output bodies, adapter configuration, gateway evidence refs, credentials, OIDC request tokens, adapter secrets, gateway tokens, model request/response payloads, PR tokens, authenticated provider URLs, and raw review bodies are not printed or persisted.
- **SC-046**: Block 20 fixtures prove that `forensics-basic-v1` query-pack output shows mixed positive evidence, digest-only reconstruction caps, missing forensic-retention facts, missing adapter-capture facts, unsupported observers, unresolved redaction, task supersession, unverified claims, unsafe provider refs, and malformed inputs as deterministic rows with preserved evidence states, closed-enum evidence families, source refs, reconstructability flags, and safe evidence-gap descriptions.
- **SC-047**: A reviewer can run the documented forensics query-pack command against committed fixtures and identify what is reconstructable, what is only digest-bound, what is missing, and what evidence gap was observed without reading raw JSONL manually.
- **SC-048**: Safety-sensitive forensics query-pack and explain tests prove that raw command arguments, command names, executable paths, script paths, unsafe test identifiers, stdout/stderr bodies, prompts, source snippets, tool-call input/output bodies, adapter configuration, gateway evidence refs, credentials, OIDC request tokens, adapter secrets, gateway tokens, PR tokens, authenticated provider URLs, raw model request/response payloads, raw review bodies, unsafe raw-reference access notes, and key material are not printed or persisted, that output-safety classes are verified against serialized output, and that negative leak assertions do not echo candidate secrets in failure output.
- **SC-049**: Block 21 fixtures prove that `cross-repo-evidence-posture-v1` exports at least two repositories and two time windows with metric rows grouped by repo, team, service, harness, change type, evidence family, row evidence state, and input trust state.
- **SC-050**: A reviewer can run the documented cross-repository posture export command against committed fixtures and inspect numerator, denominator, time window, dimensions, source artifact digest refs, `not_assessed` counts, input trust states, current/previous values, deltas, and non-comparable refusal reasons without receiving a native degradation verdict.
- **SC-051**: Safety-sensitive cross-repository export and explain tests prove that raw command arguments, command names, executable paths, script paths, unsafe test identifiers, stdout/stderr bodies, prompts, source snippets, tool-call input/output bodies, adapter configuration, gateway evidence refs, credentials, OIDC request tokens, adapter secrets, gateway tokens, PR tokens, authenticated provider URLs, raw model request/response payloads, raw review bodies, unsafe raw-reference access notes, private filesystem paths, unsafe personal identifiers, unsafe labels, raw digest-manifest paths, and free-text parser/refusal reasons are not printed or persisted.
- **SC-052**: Block 22 fixtures prove that GitLab CI, Buildkite, customer PKI, and air-gapped witness profiles share provider-neutral witness semantics for valid, missing, stale, mismatched, environment-only, unsupported, malformed, and unsafe inputs.
- **SC-053**: A reviewer can run documented witness profile commands or fixture validation against committed Block 22 fixtures and identify identity source, signing boundary, freshness boundary, artifact binding, independence state, requested trust scope, established trust scope, unsupported states, and `not_assessed` or `cannot_verify` gaps without reading raw provider logs or private customer material.
- **SC-054**: Safety-sensitive witness profile tests prove that witness JSON and explain output do not print or persist CI tokens, OIDC tokens, JWT bodies, private key material, provider tokens, authenticated provider URLs, raw job logs, private filesystem paths, unsafe personal identifiers, free-text parser errors containing input content, or customer directory, LDAP, SAML, cloud, Vault, HSM, KMS, or PKI payloads.
- **SC-055**: Block 26 fixtures prove that CI artifact observation distinguishes valid uploaded bundles, checked-in-only claims, absent bundles, partial bundles, self-referential indexes, digest mismatches, missing source/run binding, contradictory source/run binding, agent-reported happy paths, unsafe artifact output, and expired artifacts with deterministic `pass`, `fail`, `cannot_verify`, or `not_assessed` states.
- **SC-056**: A reviewer can run documented CI artifact observation commands against committed Block 26 fixtures and identify which required artifact families were present, absent, partial, expired, inaccessible, malformed, unsafe, checked-in-only, agent-reported, or not assessed without reading raw provider logs.
- **SC-057**: Safety-sensitive CI artifact observation tests prove that JSON and explain output do not print or persist raw CI logs, command bodies, stdout/stderr bodies, prompts, model responses, OIDC tokens, JWT bodies, CI secrets, private keys, provider tokens, authenticated provider URLs, private filesystem paths, raw parser input, or unsafe personal/customer payloads.

## Assumptions

- external policy consumers may consume `sdp-trace` artifacts but live outside this product boundary.
- Beads remains useful for local work tracking, but Beads is not a product dependency and is not the repo observer's source of truth.
- The initial implementation may be schema and documentation heavy before adding tiny validation tools.
- Customer pilot artifacts may need sanitization before committing summaries to the repository.
- Current schemas already declare JSON Schema Draft 2020-12; this feature standardizes that draft for new schemas unless a future major version changes it.
- Until `sdp-trace` reaches v1.0, schema changes may still be breaking, but every breaking change must update examples, compatibility notes, and downstream consumer handoff documentation.
- Public examples may use synthetic human-held roles for DRI and approver fields; customer pilots must map those fields to the customer's accepted identity or approval system.
- Public Sigstore/Rekor verification is the target profile, but private or air-gapped pilots may use an approved equivalent if the manifest, DSSE envelope, identity policy, and verification result remain explicit.
