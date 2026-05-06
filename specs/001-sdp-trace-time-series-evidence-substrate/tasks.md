# Tasks: sdp-trace Time-Series Evidence Substrate

**Input**: Design documents from `/specs/001-sdp-trace-time-series-evidence-substrate/`
**Prerequisites**: `spec.md`, `plan.md`, `research.md`, `data-model.md`, `contracts/sdp-trace-sdp-gate-boundary.md`
**Tests**: Include schema syntax checks now; schemas target JSON Schema Draft 2020-12 and committed examples validate with pinned local `ajv@8.20.0` through `scripts/validate-json-schema.mjs` once schema/example pairs exist.

**Organization**: Tasks are grouped by user story to preserve independent value and reviewability.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel because it touches different files or produces an independent artifact.
- **[Story]**: User story from `spec.md`.
- **Beads mirror**: Optional secondary tracking issue. Beads does not replace this task list.

## Phase 1: Setup and Canonical SpecKit Package

**Purpose**: Make SpecKit artifacts the repo-observable planning source.

- [x] T001 [US4] Add root README pointer to `specs/001-sdp-trace-time-series-evidence-substrate/spec.md`
- [x] T002 [US4] Link Beads epic `sdp-trace-cdn` and child issues to this spec with `bd update --spec-id`
- [x] T003 [US4] Document in `docs/speckit-compatibility.md` that Beads is secondary execution tracking, not the planning source of truth

**Checkpoint**: A repository observer can start from `specs/001-sdp-trace-time-series-evidence-substrate/` without Beads.

## Phase 2: Foundational Boundary and Extraction

**Purpose**: Prevent policy/runtime coupling before schema work starts.

- [x] T004 [US2] Write source-mapped extraction memo in `specs/001-sdp-trace-time-series-evidence-substrate/research.md` from `sdp_lab` sources (Beads mirror: `sdp-trace-cdn.1`)
- [x] T005 [US2] Finalize `contracts/sdp-trace-sdp-gate-boundary.md` and update `docs/concepts.md` with the same boundary (Beads mirror: `sdp-trace-cdn.2`)
- [x] T006 [US2] Audit `README.md`, `docs/cto-brief.en.md`, `docs/cto-brief.ru.md`, `docs/team-lead-playbook.en.md`, and `docs/team-lead-playbook.ru.md` for language implying `sdp-trace` owns policy decisions (Beads mirror: `sdp-trace-cdn.11`)
- [x] T007 [US2] Replace or narrow gate/decision wording so external verdicts are recorded inputs, not `sdp-trace` decisions (Beads mirror: `sdp-trace-cdn.11`)
- [x] T051 [US4] Rewrite `docs/cto-brief.en.md` and `docs/cto-brief.ru.md` as one-minute CTO decision narratives mapped to SpecKit evidence, with no marketing claims or native `sdp-trace` policy verdicts (Beads mirror: `sdp-trace-cdn.11`)
- [x] T034 [US4] Document JSON Schema Draft 2020-12, pinned local `ajv@8.20.0`, schema IDs, and schema versioning in `schema/README.md` before new schema authoring starts (Beads mirror: `sdp-trace-cdn.8`)
- [x] T040 [US2] Document committed artifact safety rules: sanitization, SHA-256 digests, redaction notes, `integrity_status`, and no secrets or raw customer data in committed examples
- [x] T044 [US5] Design `schema/accountability.schema.json` for human-held DRI, approver, escalation, risk owner, authority scope, approval reference, and line of defense (Beads mirror: `sdp-trace-cdn.2`)
- [x] T045 [US5] Design `schema/risk-classification.schema.json` for observed autonomy, observed impact, classification source, and external oversight assertions without internal pass/fail or oversight derivation (Beads mirror: `sdp-trace-cdn.2`)
- [x] T046 [US5] Design `schema/contract-manifest.schema.json` and example manifest covering schema, doc, validation script, fixture, source commit, approval, compatibility, and previous-manifest digests (Beads mirror: `sdp-trace-cdn.8`)
- [x] T047 [US5] Design `schema/contract-release-verification.schema.json` and signing profile docs for `sdp-trace-signature/sigstore-dsse-keyless-v1` (Beads mirror: `sdp-trace-cdn.8`)
- [x] T048 [US5] Add negative fixtures for AI-as-sole-accountable-owner and modified contract manifest digest mismatch (Beads mirror: `sdp-trace-cdn.8`)
- [x] T049 [US5] Add trusted signer identity policy example and mismatch fixture covering OIDC issuer, source URI, protected ref, workflow identity, release captain, and required approval refs (Beads mirror: `sdp-trace-cdn.8`)
- [x] T050 [US5] Produce one local contract release verification evidence record for the target signing profile shape or approved private equivalent before claiming contract scaffolding complete (Beads mirror: `sdp-trace-cdn.8`)

**Checkpoint**: `sdp-trace` and `sdp-gate` ownership, CTO narrative, validator strategy, and artifact safety rules are clear before new schemas are added.

## Phase 3: User Story 1 - CTO Reviews Process Movement (Priority: P1)

**Goal**: Define time-series observations and metric streams without built-in degradation policy.

**Independent Test**: A sample metric stream can show movement across windows while every sample has evidence or `not_assessed`.

### Contract and Schema Tasks

- [x] T008 [P] [US1] Design `schema/observation.schema.json` for evidence-backed observations (Beads mirror: `sdp-trace-cdn.3`)
- [x] T009 [P] [US1] Design `schema/metric-stream.schema.json` for metric samples, stream comparisons, prior/current values, deltas, evidence coverage, and no interpretation labels (Beads mirror: `sdp-trace-cdn.3`)
- [x] T010 [US1] Add examples under `examples/github-speckit/` showing current-window vs previous-window movement without policy verdicts or degradation labels
- [x] T011 [US1] Define metric catalog in `docs/process-metric-catalog.md` with units, dimensions, evidence source, and `not_assessed` rule (Beads mirror: `sdp-trace-cdn.5`)
- [x] T012 [US1] Update `schema/trace.schema.json` or document a replacement path so trace snapshots can include observations and metric samples

### Verification Tasks

- [x] T013 [US1] Run `jq empty schema/*.json`
- [x] T014 [US1] Record validation output in a sanitized evidence note under `docs/research/`

**Checkpoint**: CTO-facing process movement exists as data, not as a policy verdict.

## Phase 4: User Story 2 - sdp-gate Inherits Trace Contracts (Priority: P1)

**Goal**: Produce the assessment input contract consumed by `sdp-gate`.

**Independent Test**: An assessment input example contains evidence, observations, metric streams, and `not_assessed`, but no pass/fail/degradation decision.

- [x] T015 [P] [US2] Design `schema/evidence-event.schema.json` from portable evidence-event concepts, including redaction, hash, integrity, pending, duplicate, and conflict metadata (Beads mirror: `sdp-trace-cdn.4`)
- [x] T016 [P] [US2] Design `schema/provenance-record.schema.json` for actor/model/harness/tool provenance, SHA-256 payload digests, and optional same-chain `hash_prev` (Beads mirror: `sdp-trace-cdn.4`)
- [x] T017 [US2] Design `schema/assessment-input.schema.json` for policy-engine handoff with no native policy verdicts
- [x] T018 [US2] Add `schema/external-verdict-input.schema.json` and examples that record external verdicts or evidence-strength assertions as externally produced inputs
- [x] T019 [US2] Add `examples/go-service/assessment-input.json` or equivalent portable example
- [x] T041 [US2] Add a negative validation example showing that a native `sdp-trace` assessment input cannot contain pass/fail/readiness/degradation fields
- [x] T042 [US2] Update `schema/README.md` with ownership, external verdict, validation, versioning, and migration rules

**Checkpoint**: `sdp-gate` has a clear inherited input contract.

## Phase 5: Self-Trace v0 - Mandatory Product Proof (Priority: P0)

**Goal**: Use the new `sdp-trace` contracts to trace development of this SpecKit feature before running the customer pilot matrix or making any product trust claim.

**Independent Test**: `examples/self-trace/assessment-input.json` can describe this feature's spec, plan, tasks, commits, evidence, observations, metric stream, and `not_assessed` gaps without a gate decision.

**Beads mirror**: `sdp-trace-cdn.12`

**Hard rule**: If T020-T026 are incomplete, `sdp-trace` has not proven itself. Pilot work remains blocked.

- [x] T020 [US1] Add `examples/self-trace/evidence-events.json` covering commits, touched files, commands run, schema checks, SpecKit task status, redaction notes, and integrity status
- [x] T021 [US1] Add `examples/self-trace/provenance-records.json` covering human actor, Codex session, command summaries, payload digests, and missing fields as `not_assessed`
- [x] T022 [US1] Add `examples/self-trace/observations.json` recording that SpecKit artifacts exist, Beads mirrors SpecKit, schema syntax passed, and boundary docs remain partial until T005-T007 complete
- [x] T023 [US1] Add `examples/self-trace/metric-stream.json` with at least evidence completeness, `not_assessed` count, schema validity, and SpecKit task completion samples
- [x] T024 [US2] Add `examples/self-trace/trace-snapshot.json` linking spec, plan, tasks, changes, evidence, observations, metric samples, and any external verdict inputs as external
- [x] T025 [US2] Add `examples/self-trace/assessment-input.json` as the first policy-engine handoff package, without pass/fail/degradation decision
- [x] T026 [US4] Record self-trace validation notes in `docs/research/self-trace-v0-summary.md`
- [x] T052 [US4] Add a self-trace validation command that validates `examples/self-trace/assessment-input.json` and every referenced self-trace artifact from a fresh checkout (Beads mirror: `sdp-trace-cdn.12`)
- [x] T053 [US4] Record the crisis Socratic review artifacts as external review evidence in the self-trace package (Beads mirror: `sdp-trace-cdn.12`)
- [x] T054 [US5] Ensure self-trace accountability names human-held DRI, approver, risk owner, and escalation for the repository proof itself (Beads mirror: `sdp-trace-cdn.12`)
- [x] T055 [US1] Add self-trace metric samples for contract task completion, evidence coverage, `not_assessed` count, schema validity, and review contradiction count (Beads mirror: `sdp-trace-cdn.12`)
- [x] T056 [US2] Add a negative self-trace fixture proving a self-trace package with native pass/fail/readiness/degradation fields fails validation (Beads mirror: `sdp-trace-cdn.12`)

**Checkpoint**: `sdp-trace` can describe its own development at v0 before claiming external pilot readiness. Without this checkpoint, the product remains conceptually unproven.

## Phase 5A: Self-Attested Contract Release (Priority: P0)

**Goal**: Prove the contract release against an immutable source reference and explicit attestation level before treating the contract as trusted.

**Independent Test**: A fresh checkout can run one command that reports `schema_valid`, `digest_verified`, `locally_attested`, `externally_attested`, and `production_release_verified` as separate states, with missing states recorded as `not_assessed`.

**Beads mirror**: `sdp-trace-cdn.13`

- [x] T057 [US5] Replace local manifest `source_commit` placeholders with an immutable git commit or signed source reference before trusted-release claims (Beads mirror: `sdp-trace-cdn.13`)
- [x] T058 [US5] Add `scripts/verify-self-attestation.sh` or equivalent release verifier that reports separate proof states instead of one overloaded trusted flag (Beads mirror: `sdp-trace-cdn.13`)
- [x] T059 [US5] Add a self-attestation evidence record under `examples/self-trace/` covering manifest digest, DSSE envelope binding, signer identity, freshness, and rollback status (Beads mirror: `sdp-trace-cdn.13`)
- [x] T060 [US5] Add negative self-attestation fixtures for wrong source commit, wrong signer, wrong trusted identity policy, stale manifest, missing external attestation, and modified verification artifact (Beads mirror: `sdp-trace-cdn.13`)
- [x] T061 [US4] Record a self-attestation summary under `docs/research/self-attestation-summary.md` with exact commands and residual `not_assessed` items (Beads mirror: `sdp-trace-cdn.13`)

**Checkpoint**: Local development proof may remain useful, but product trust claims remain blocked until self-attestation records the actual proof state without hiding missing external trust anchors.

## Phase 5B: Source-Bound Release Finalization and External Trust Design (Priority: P0)

**Goal**: Define the next release-finalization step that can close local source-content proof while keeping external production trust explicit and evidence-backed.

**Independent Test**: A repository observer can read Block 04 spec and Socratic dialogue, then distinguish a source-bound local release from an externally trusted production release without relying on Beads context.

**Beads mirror**: `sdp-trace-cdn.21`

- [x] T062 [US5] Add `blocks/04-release-finalization-external-trust.md` defining source-bound local finalization plus external Sigstore/Rekor and customer PKI trust design (Beads mirror: `sdp-trace-cdn.21`)
- [x] T063 [US5] Add `blocks/04-release-finalization-socratic.md` recording challenge/resolution dialogue for local DSSE, dirty worktree proof, external trust anchors, customer PKI, pilot claims, and circular digest risk (Beads mirror: `sdp-trace-cdn.21`)
- [x] T064 [US5] Add implementation plan for a source-bound local finalization guard that refuses dirty-tree source-content proof and verifies manifest artifact digests against a committed source reference; release-proof regeneration remains a separate future step (Beads mirror: `sdp-trace-cdn.21`)
- [x] T065 [US5] Extend self-attestation result or adjacent schema to record external trust profile, transparency evidence, protected ref status, workflow identity status, approval status, and production release verification without generic trust shortcuts (Beads mirror: `sdp-trace-cdn.21`)
- [x] T066 [US5] Add external trust negative fixtures for missing Rekor/customer audit evidence, OIDC issuer mismatch, source URI mismatch, protected ref mismatch, workflow identity mismatch, approval mismatch, and expired freshness (Beads mirror: `sdp-trace-cdn.21`)
- [x] T067 [US4] Update release docs and self-attestation summary to use distinct terms: source-bound local release and externally trusted production release (Beads mirror: `sdp-trace-cdn.21`)
- [x] T068 [US4] Run Socratic review and obtain explicit consensus before implementing Block 04 code or changing release artifacts (Beads mirror: `sdp-trace-cdn.21`)
- [x] T069 [US4] Run pi review for Block 04 specs and implementation, register every finding in Beads, and close every valid finding including P3/minor items (Beads mirror: `sdp-trace-cdn.21`)
- [ ] T070 [US4] Stale historical closure: re-verify `rtk npm run validate`, `rtk git diff --check`, `rtk scripts/verify-self-attestation.sh --all`, and source-bound finalization guard fixtures before Block 04 can be treated as currently closed (Beads mirror: `sdp-trace-cdn.21`)
<!-- sdp-trace-claim: claim=task_closed; subject=T070; state=stale; profile=repo_baseline; evidence=state:claim_tags_consistent -->

**Checkpoint**: Block 04 implementation artifacts exist, but current closure is stale until T070 is re-verified by live commands. External trust remains `not_assessed` unless real Sigstore/Rekor or accepted customer PKI evidence is committed.

## Phase 6: User Story 3 - Pilot Evaluates Harness, Model, and JVM Stack Slices (Priority: P1)

**Goal**: Create repeatable pilot run-cards and evidence paths for the customer-requested matrix.

**Independent Test**: Each run-card lists prompt, expected artifacts, provenance fields, `unbacked_claim` capture, validation, and `not_assessed` behavior.

- [x] T027 [P] [US3] Add OpenCode run-card covering MiniMax, Kimi, and GLM in `docs/research/opencode-model-run-card.md` (Beads mirror: `sdp-trace-cdn.6`)
- [x] T028 [P] [US3] Add harness run-card for the Superpowers-style harness pattern, `gsd`, `gsd2`, and Oh My OpenAgent in `docs/research/harness-run-card.md` without introducing a Superpowers dependency (Beads mirror: `sdp-trace-cdn.6`)
- [x] T029 [US3] Add Kotlin+Bazel pilot fixture plan in `docs/research/kotlin-bazel-fixture-plan.md` (Beads mirror: `sdp-trace-cdn.7`)
- [x] T030 [US3] Update `docs/jvm-bazel-guide.md` with Kotlin+Bazel-specific evidence requirements
- [x] T031 [US3] Add or update `examples/jvm-bazel/` with a Kotlin+Bazel evidence bundle or fixture placeholder that is explicitly `not_assessed` until run; do not call the Kotlin+Bazel gap closed until a committed run artifact exists
- [x] T032 [US3] Update `docs/harness-compatibility-matrix.md` as a legacy-named evidence matrix with evidence state, reason code, artifact reference, gap reason, and next required evidence only (Beads mirror: `sdp-trace-cdn.10`)
- [x] T033 [US3] Update `docs/model-compatibility.md` as a legacy-named evidence matrix with evidence state, reason code, artifact reference, gap reason, and next required evidence only (Beads mirror: `sdp-trace-cdn.10`)
- [x] T071 [US3] Add Block 05 spec, Socratic dialogue, and implementation plan for customer pilot run-cards and evidence-package scope (Beads mirror: `sdp-trace-cdn.22`)
- [x] T072 [US3] Run pi review for Block 05 spec artifacts and register every valid finding in Beads, including minor/P3 findings (Beads mirror: `sdp-trace-cdn.22`)
- [x] T073 [US3] Close every valid Block 05 spec-gate review finding before implementation starts; implementation-target findings remain open under T027-T033/T037 (Beads mirror: `sdp-trace-cdn.22`)
- [x] T074 [US3] Run pi review for Block 05 implementation and register every valid finding in Beads, including minor/P3 findings (Beads mirror: `sdp-trace-cdn.22`)
- [x] T075 [US4] Verify Block 05 validation commands and close `sdp-trace-cdn.22` only after implementation review findings are closed (Beads mirror: `sdp-trace-cdn.22`)
- [x] T076 [US4] Add committed Block 05 pi-review ledger with findings, severity, disposition, evidence, and optional Beads mirror IDs (Beads mirror: `sdp-trace-cdn.22`)
- [x] T077 [US4] Add deterministic pilot matrix validation and wire it into `scripts/validate-contracts.sh` (Beads mirror: `sdp-trace-cdn.22`)

**Checkpoint**: Pilot scope is executable without unsupported compatibility claims.

## Phase 6A: User Story 3A - First Real Product Proof (Priority: P0)

**Goal**: Prove one real E2E value path with OpenCode + MiniMax + Kotlin+Bazel without turning external tools into product dependencies.

**Independent Test**: In an environment with OpenCode, MiniMax model access, and a Kotlin+Bazel target, a pilot operator can run the Block 06 command and validate a generated `sdp-trace` evidence package. A committed sanitized report states the tested-on environment and residual `not_assessed` gaps.

**Beads mirror**: `sdp-trace-drq`

- [x] T078 [US3] Add Block 06 spec, Socratic dialogue, and implementation plan for OpenCode + MiniMax + Kotlin+Bazel E2E proof (Beads mirror: `sdp-trace-drq`)
- [x] T079 [US3] Run pi review for Block 06 spec artifacts, record every valid finding in `blocks/06-opencode-minimax-kotlin-bazel-review-ledger.md`, and mirror every valid finding in Beads including minor/P3 findings (Beads mirror: `sdp-trace-drq`)
- [x] T080 [US3] Close every valid Block 06 spec-gate review finding in the committed review ledger before implementation starts; Beads closure mirrors the ledger (Beads mirror: `sdp-trace-drq`)
- [x] T081 [US3] Add reference runner `scripts/run-opencode-minimax-kotlin-bazel-proof.sh` that shells out to external OpenCode/Bazel tools without adding repository dependencies (Beads mirror: `sdp-trace-drq`)
- [x] T082 [US3] Add `scripts/validate-e2e-pilot-package.sh` and proof-state checks for generated E2E pilot packages (Beads mirror: `sdp-trace-drq`)
- [x] T083 [US3] Produce a committed sanitized OpenCode + MiniMax + Kotlin+Bazel proof package from a real run, or keep Block 06 incomplete with explicit `not_assessed` reason (Beads mirror: `sdp-trace-drq`)
- [x] T084 [US3] Update OpenCode/Kotlin+Bazel docs and exact matrix rows only from the committed Block 06 proof package (Beads mirror: `sdp-trace-drq`)
- [x] T085 [US4] Run pi review for Block 06 implementation, record every valid finding in `blocks/06-opencode-minimax-kotlin-bazel-review-ledger.md`, and mirror every valid finding in Beads including minor/P3 findings (Beads mirror: `sdp-trace-drq`)
- [x] T086 [US4] Verify `npm run validate`, `jq empty schema/*.json`, `git diff --check`, and Block 06 package validation before closing the review ledger and mirroring closure to `sdp-trace-drq` (Beads mirror: `sdp-trace-drq`)
- [x] T091 [US4] Regenerate source-bound local release artifacts for the current committed source subject, verify `source_bound_local_release`, and keep `external_production_trust` open unless real Sigstore/Rekor or accepted customer PKI evidence is committed.
- [x] T092 [US4] Add Block 07 Slice 7 review-ledger schema, machine-checkable closure ledger, validator wiring, and final verifier-derived closure statement that keeps external production trust explicitly blocked without blocking Block 08 discovery.

**Checkpoint**: The repository has at least one real, validated, sanitized product proof package for the exact first tested slice, or the product proof remains explicitly incomplete.

## Phase 7: User Story 4 - Repository Observer Finds SpecKit Evidence (Priority: P2)

**Goal**: Make the evidence package self-explanatory from committed files.

**Independent Test**: A reviewer can follow `quickstart.md`, validate schemas, and identify what remains `not_assessed`.

- [x] T035 [US4] Add pass and `not_assessed` fixtures for new schemas under `examples/`
- [x] T036 [US4] Add validation command that uses pinned local `ajv@8.20.0` and excludes `.git`, `.beads`, `.sdp-trace-runs`, `benchmarks/repos/`, temporary directories, editor caches, and generated dependency directories (Beads mirror: `sdp-trace-cdn.8`)
- [x] T037 [US4] Build customer pilot evidence package outline in `docs/research/customer-pilot-evidence-package.md` (Beads mirror: `sdp-trace-cdn.9`)
- [x] T038 [US4] Verify `jq empty schema/*.json`
- [x] T039 [US4] Verify all committed examples are parseable JSON where applicable
- [x] T043 [US4] Verify schema-version fields and compatibility notes for all committed examples that claim a schema contract

**Checkpoint**: The repository itself explains the plan, proof, gaps, and execution path.

## Phase 8: Agent and Human Usage Discovery (Queued After Block 07)

**Goal**: Make the first-use path clear for both agents and human reviewers without creating claims that outrun the verifier.

**Activation Gate**: Do not elaborate or implement this phase until Block 07 has either closed under live verifier evidence or remains open with a verifier-derived blocking state that this phase can explain honestly.

- [x] T087 [US4] After Block 07, run design review for `blocks/08-agent-human-usage-discovery.md` with separate agent-path and human-path critiques
- [x] T088 [US4] Define the agent entrypoint: profile selection, verifier command surface, evidence emission rules, and forbidden claims
- [x] T089 [US4] Define the human entrypoint: five-minute verification path, proof-scope explanation, dirty checkout warning, and `not_assessed` interpretation
- [x] T090 [US4] Review Block 08 with a fresh agent and record every valid finding before implementation planning

**Checkpoint**: Agents and humans can independently discover how to use `sdp-trace`, and both paths converge on the same verifier-derived proof states.

## Phase 9: Flight Recorder Trust Kernel (Before External Demo)

**Goal**: Make `sdp-trace` a credible flight recorder for agentic development before using the Feature Flag / Entitlements Kotlin+Bazel demo to stress the mechanics.

**Independent Test**: A repository observer can validate committed flight-recorder fixtures showing local chain validity, witnessed chain validity, tamper detection, late-attach gaps, requirement supersession, redaction states, and reviewer query output without relying on a demo repo.

**Activation Gate**: Do not start a new external Feature Flag / Entitlements demo until Block 09 has a recorder kernel that distinguishes local consistency from witnessed accountability evidence.

- [x] T093 [US4] Add Block 09 design artifact for the flight-recorder trust kernel, including threat model, local/witnessed/external profile split, event-chain requirements, redaction model, verifier semantics, and query surface
- [x] T094 [US4] Add Block 09 Socratic review synthesis with executive-role findings and explicit dispositions for local-only chain, mid-flight attach, bypass, redaction, digest-only evidence, and query-surface risks
- [x] T095 [US4] Add Block 09 implementation plan sliced for agent handoff before any demo-repo work starts
- [x] T096 [US4] Add Block 09 review ledger and close or explicitly block every critical/major spec-gate finding before implementation begins
- [x] T097 [US4] Add flight-recorder event, run, and witness schemas plus positive/negative fixtures
- [x] T098 [US4] Implement local chain verifier with mutation, deletion, and reordering negative tests
- [x] T099 [US4] Implement witnessed-run verifier with missing/mismatched witness negative tests
- [x] T100 [US4] Implement late-attach and requirement-supersession fixtures and verifier behavior
- [x] T101 [US4] Implement evidence-retention and redaction-state fixtures plus verifier behavior
- [x] T102 [US4] Implement reviewer query surface for run summary, provenance, gaps, requirements, commands, file mutations, tests, redactions, and witness state
- [x] T103 [US4] Wire committed Block 09 fixture validation into repository validation after verifier behavior exists
- [x] T104 [US4] Run strict implementation review for Block 09 and record every valid finding before activating the external Feature Flag / Entitlements demo

**Checkpoint**: `sdp-trace` can prove recorder-chain integrity and witness agreement for committed fixtures, while clearly marking local-only, late-attach, redaction, and unavailable evidence limits.

## Phase 10: Gate Contract, Explain, And Native Override Event

**Goal**: Make advisory gate output operationally useful without claiming
protected enforcement or audit-grade trust.

**Independent Test**: A reviewer can run committed fixtures showing missing
required runs, mismatched CI witness binding, visible override request state,
deterministic `gate explain`, deterministic read-only `gate preview`, and
unchanged `audit_grade_gate: cannot_verify`.

**Activation Gate**: Do not implement protected fail-closed behavior in this
phase. Block 14 may emit facts for an external policy consumer; it must not
make merge, release, readiness, degradation, or risk-acceptance decisions.

- [x] T105 [US4] Add Block 14 spec and implementation plan for advisory gate contract, required runs, explain, preview, and override event scope
- [x] T106 [US2] Extend the contract model with optional required runs and tests for absent, unmatched, and protected-future run requirements
- [x] T107 [US2] Add or extend Draft 2020-12 gate-result schema covering schema version, gate mode, trust cap, required-run states, witness bindings, override requests, reasons, next actions, gate conditions, and run summaries
- [x] T108 [US2] Extend gate-result output with schema version, gate mode, trust cap, required-run states, witness bindings, override requests, reasons, and next actions
- [x] T109 [US5] Add CI witness binding checks for repository, ref, commit, run id, and artifact digest mismatches without upgrading audit-grade trust
- [x] T110 [US5] Add `policy_override_requested` flight-recorder event support with producer/origin fields and CLI/external-reference path that records override requests without converting missing evidence to pass
- [x] T111 [US4] Add deterministic `gate explain` output with next actions for missing telemetry, cannot-verify witness, stale witness, source mismatch, and override-present states
- [x] T112 [US4] Add deterministic read-only `gate preview` output for gate-relevant fields, selected mode, trust cap, required runs, evidence ids, witness inspectability, and locally detectable witness mismatches
- [x] T113 [US4] Add committed Block 14 fixtures for missing required run, unmatched run, stale witness, source mismatch, artifact mismatch, valid override, malformed override, and protected-future requirement
- [x] T114 [US4] Add safety-sensitive negative tests proving gate, explain, and preview output do not print raw command args, stdout/stderr bodies, prompts, source snippets, credentials, OIDC request tokens, or model responses
- [x] T115 [US4] Run Go-first verification, schema checks, strict review, pi review, and record Block 14 review disposition before PR closure

**Checkpoint**: Advisory gate users can see what is missing, what is witnessed,
what was overridden, what remains unverified, and what to do next. Protected
enforcement and audit-grade trust remain explicitly blocked until later signed
checkpoint and external policy-consumer work.

## Phase 11: Signed Checkpoint And Replay Resistance

**Goal**: Make a run checkpoint independently replayable enough for a verifier
to detect post-hoc mutation, source/run replay, stale checkpoint reuse, and
signer-policy mismatch without claiming protected enforcement.

**Independent Test**: A reviewer can run committed fixtures showing a valid
local signed checkpoint, payload tampering, replay against another run,
sequence gap/duplicate failure, signer-policy mismatch, and unchanged
`audit_grade_gate: cannot_verify`.

**Activation Gate**: Do not implement protected fail-closed behavior in this
phase. Block 15 may emit signed-checkpoint verification facts; it must not make
merge, release, readiness, degradation, override approval, or risk-acceptance
decisions.

- [x] T116 [US5] Add Block 15 spec and implementation plan for signed checkpoint artifacts, replay binding, monotonic sequence checks, signer authority policy, and no-overclaim gate posture
- [x] T117 [US5] Add signed checkpoint and trusted checkpoint policy schemas with local development, CI isolated job, and external witness authority states
- [x] T118 [US5] Add Go checkpoint payload derivation from run artifacts, including run id, run nonce, chain head, event count, source snapshot, task hash, contract digest, and previous checkpoint digest
- [x] T119 [US5] Add detached Ed25519 checkpoint signing and verification for the local development profile with deterministic payload digest checks
- [x] T120 [US5] Add replay-resistance verification for run id, nonce, source snapshot, task hash, contract digest, event count, and chain head mismatches
- [x] T121 [US5] Add monotonic checkpoint-set verification for duplicate, missing, and descending sequence evidence
- [x] T122 [US5] Add trusted-checkpoint policy verification that fails unauthorized signers and leaves missing policy `not_assessed`
- [x] T123 [US4] Add `checkpoint create` and `checkpoint verify` CLI commands with safety-sensitive negative output tests
- [x] T124 [US4] Add gate-level tests proving local signed checkpoint evidence does not convert `protected_future` or `audit_grade_gate` to pass
- [x] T125 [US5] Add committed Block 15 fixtures and record Go-first verification, schema checks, strict review, pi review, and review disposition before PR closure

**Checkpoint**: `sdp-trace` can verify that a signed checkpoint belongs to the
selected run context and was not replayed from another run. Protected
enforcement and audit-grade trust remain blocked until external policy-consumer
and external witness work exist.

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1**: Start immediately.
- **Phase 2**: Depends on Phase 1.
- **Phase 3**: Depends on boundary clarity, T034 validator strategy, T040 artifact safety, and T044-T050 accountability/signing contract foundations from Phase 2.
- **Phase 4**: Depends on Phase 2 and can run partly in parallel with Phase 3.
- **Phase 5**: Depends on minimal schemas from Phases 3-4 and runs before external pilot claims.
- **Phase 6**: Depends on self-trace and self-attestation learning from Phase 5 and Phase 5A and must not record support, readiness, or compatibility as native `sdp-trace` outcomes.
- **Phase 6A**: Depends on Phase 6 run-card design and must prove one real OpenCode + MiniMax + Kotlin+Bazel slice before any product packaging or pilot-readiness claim.
- **Phase 7**: Depends on schema and pilot artifacts from Phases 3-6.
- **Phase 8**: Depends on Block 07 live verifier outcome. It must not mask missing trust evidence with onboarding or documentation polish.
- **Phase 9**: Depends on Block 08 discovery and executive Socratic review. It must land the recorder trust kernel before any new external demo-repo execution; the demo is a stress test, not the place to invent trust mechanics.
- **Phase 10**: Depends on Block 13B capture-boundary and state taxonomy, and reuses Block 11/12 report, gate, and CI witness behavior. It must remain advisory until later signed checkpoint and external policy-consumer enforcement work exists.
- **Phase 11**: Depends on Block 14 advisory gate output and Block 09/13B run-chain semantics. It supplies signed-checkpoint facts but must not implement protected enforcement or external witness trust.

### Parallel Opportunities

- T008 and T009 can run in parallel.
- T015 and T016 can run in parallel.
- T044 and T045 can run in parallel.
- T046 and T047 can run in parallel after T034 and T040.
- T049 can run after T047.
- T050 can run after T046, T047, and T049.
- T020 through T023 can run in parallel once minimal schemas exist.
- T027 and T028 can run in parallel.
- T032 and T033 can run in parallel after pilot evidence exists.
- T081 and T082 can run in parallel after T080 closes, because the runner and validator have separate write scopes.
- T106 and T109 can run in parallel after T105 because required-run semantics and witness binding have separate implementation surfaces.
- T111 and T112 can run in parallel after T108 because explain and preview write separate command paths over the same output contract.

## Implementation Strategy

### MVP First

1. Complete Phase 1.
2. Complete Phase 2.
3. Complete T034, T040, and T044-T050 before new evidence/assessment schema authoring.
4. Complete T008, T009, T015, T016, T017, and T018.
5. Add one valid example, one `not_assessed` example, one negative policy-verdict example, one negative AI-accountability example, and one negative modified-manifest example.
6. Complete Self-Trace v0 tasks T020-T026 and T052-T056.
7. Complete Self-Attested Contract Release tasks T057-T061 or explicitly record external attestation as `not_assessed`.
8. Complete Block 05 run-cards without claiming observed behavior.
9. Complete Block 06 for OpenCode + MiniMax + Kotlin+Bazel before product packaging or pilot-readiness claims.
10. Complete Block 09 flight-recorder trust kernel before starting a new Feature Flag / Entitlements external demo.
11. Complete Block 14 advisory gate contract and explanation work before any protected enforcement or signed-checkpoint gate claim.
12. Complete Block 15 signed-checkpoint replay verification before any protected gate profile claims that checkpoint evidence exists.

### Evidence Discipline

- Do not record harness/model/stack support, readiness, or compatibility as native `sdp-trace` outcomes. Record observed evidence state or an explicitly external verdict input.
- Do not add policy thresholds to `sdp-trace`.
- Keep raw pilot outputs ignored until sanitized.
- Every public claim must link to a file, command, example, or `not_assessed` entry.
