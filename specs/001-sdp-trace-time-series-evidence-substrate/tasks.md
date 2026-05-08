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
- [x] T070 [US4] Stale historical closure retired: Block 04's old `npm` and script verifier commands are preserved as historical evidence only, and current closure is bounded by Go-first tests, JSON syntax checks, diff checks, and current release-proof output. External production trust remains `not_assessed`; removed source-bound finalization scripts are not current verifier evidence. (Beads mirror: `sdp-trace-cdn.21`)
<!-- sdp-trace-claim: claim=task_closed; subject=T070; state=pass; profile=repo_baseline; evidence=file:docs/research/historical-verifier-retirement-2026-05-07.md -->

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

## Phase 12: Protected Gate Enforcement Profile

**Goal**: Make protected gate evaluation explicit, deterministic, and
fail-closed for missing or invalid protected evidence without making
`sdp-trace` the native merge, release, readiness, degradation, override
approval, or risk-acceptance policy owner.

**Independent Test**: A reviewer can run committed fixtures showing missing
checkpoint evidence, local-development checkpoint evidence, signer-policy
mismatch, missing CI witness binding, stale CI witness evidence, CI witness
mismatch, override presence, and a valid CI-authority protected profile with
deterministic protected gate states and exit codes.

**Closure Evidence**: Block 16 was implemented and merged in PR #6. Evidence
is recorded in `docs/research/block-16-implementation-evidence.md`, including
Go-first verification, schema checks, fixture coverage, implementation review,
and PR-level review disposition with no remaining critical or major findings.

- [x] T126 [US5] Add Block 16 spec and implementation plan for protected gate profile selection, fail-closed evidence requirements, protected trust scope, exit semantics, explain/preview behavior, and no-overclaim boundary.
- [x] T127 [US2] Add a new Block 16 gate-result schema version with selected profile, protected gate state, checkpoint verification summary, protected condition rows, explicit protected trust-scope enum, and Block 14 read compatibility without introducing native policy decision fields.
- [x] T128 [US5] Add protected gate profile evaluation that requires explicit `--profile protected`, signed checkpoint input, trusted-checkpoint policy input, required runs, and required evidence before protected pass is possible.
- [x] T129 [US5] Map checkpoint verification facts into protected profile state, proving local-development signed checkpoints and untrusted checkpoint shapes cannot pass protected gate.
- [x] T130 [US5] Bind CI signer authority to CI witness repository, ref, commit, run id, artifact digest, and freshness checks; missing binding or absent/unbounded freshness data is `cannot_verify`, contradictory binding or stale data is `fail`.
- [x] T131 [US4] Add protected profile CLI flags, deterministic exit-code behavior, usage errors for omitted required gate inputs, and preview rendering for absent protected inputs.
- [x] T132 [US4] Extend `gate explain` and `gate preview` with protected-profile input requirements, checkpoint state, signer authority state, witness binding state, override state, stable reason codes, and next actions.
- [x] T133 [US4] Add safety-sensitive negative tests proving protected gate, explain, and preview output do not print raw command args, stdout/stderr bodies, prompts, source snippets, credentials, OIDC request tokens, model responses, or checkpoint key material.
- [x] T134 [US5] Add committed Block 16 fixtures for missing checkpoint, local-development checkpoint, local-development checkpoint with invalid run binding, missing signer policy, signer mismatch, missing CI witness, absent freshness, stale CI witness, CI source mismatch, CI artifact mismatch, malformed override with a trust-scope failure, valid CI-authority protected profile, and override-present protected profile.
- [x] T135 [US5] Run Go-first verification, schema checks, strict review, pi review, PR-level review, and record Block 16 review disposition before PR closure.

**Checkpoint**: Protected gate users get fail-closed verifier facts and
deterministic CI-friendly exit behavior. `sdp-trace` still does not own the
organization's merge, release, readiness, degradation, override approval, or
risk-acceptance decision.

## Phase 13: Managed Harness Enforcement Profile

**Goal**: Make managed harness evaluation explicit, opt-in, deterministic, and
fail-closed for registered wrapper or adapter boundaries without making
managed enrollment a prerequisite for observation-mode value.

**Independent Test**: A reviewer can run committed fixtures showing unmanaged
run bypass, late enrollment, unauthorized adapter identity, missing capability,
missing harness/tool/file/test telemetry, agent-reported-only test evidence,
unverified suppression, managed witness mismatch, override presence, and a
valid managed profile with deterministic managed gate states and exit codes.

**Activation Gate**: Do not implement managed harness profile behavior until
`blocks/17-managed-harness-enforcement-profile.md` is explicitly approved.
Managed harness pass is verifier-derived evidence for an external CI or policy
owner; it is not a native `sdp-trace` merge, release, readiness, degradation,
override approval, or risk-acceptance decision.

- [x] T136 [US5] Add Block 17 spec and implementation plan for managed harness profile selection, opt-in adoption boundary, registered wrapper enrollment, adapter identity, capability requirements, fail-closed evidence requirements, exit semantics, explain/preview behavior, and no-overclaim boundary.
- [x] T137 [US2] Add a new Block 17 assessment-result schema version with selected profile, managed harness assessment state, managed boundary summary, managed condition rows, explicit managed trust-scope enum, and Block 14/16 read compatibility without introducing native policy decision fields.
- [x] T138 [US5] Add managed policy and adapter registry schemas covering approved wrapper identities, adapter identities, capabilities, allowed event types, deployment source, suppression policy, and witness binding requirements.
- [x] T139 [US5] Add managed harness profile evaluation that requires explicit `--profile managed-harness`, managed policy input, adapter registry input, selected run input, CI or managed witness input, and pre-run registered wrapper or adapter enrollment before managed pass is possible.
- [x] T140 [US5] Prove unmanaged runs, late enrollment, post-hoc adapter registration, post-hoc managed policy creation, self-claimed adapter identity, unauthorized adapter identity, and unregistered wrapper evidence cannot pass managed harness profile.
- [x] T141 [US5] Verify adapter capability and required event coverage for harness lifecycle, tool events, file mutations, and test provenance; missing required coverage emits deterministic condition rows and agent-reported-only test evidence cannot satisfy executed-test requirements.
- [x] T142 [US5] Verify suppression handling so suppressed evidence is visible, non-upgrading, and accepted only when policy-authorized for the selected managed profile.
- [x] T143 [US5] Bind managed witness evidence to source commit, run id, run nonce, managed policy digest, adapter registry digest, enrollment event digest, launch event digest, chain head, adapter or wrapper identity digest, freshness, and output artifact digests; missing binding is `cannot_verify`, contradictory binding is `fail`.
- [x] T144 [US4] Add primary `sdp-trace assess --profile managed-harness` CLI flags, deterministic exit-code behavior, usage errors for omitted required managed inputs, and preview rendering for absent managed inputs.
- [x] T145 [US4] Add `assess explain` and `assess preview` with managed-profile input requirements, wrapper enrollment state, adapter identity state, capability state, event coverage, suppression state, witness binding state, override state, stable reason codes, and next actions while preserving Block 14/16 `gate explain` compatibility.
- [x] T146 [US4] Add safety-sensitive negative tests proving managed assess, explain, and preview output do not print raw command args, stdout/stderr bodies, prompts, source snippets, credentials, OIDC request tokens, adapter secrets, gateway tokens, model responses, or checkpoint key material.
- [x] T147 [US5] Add committed Block 17 fixtures for unmanaged run, late enrollment, unauthorized adapter, missing capability, missing harness event, missing tool event, missing file mutation event, agent-reported test evidence, policy-authorized suppression, suppression without policy, witness missing, witness mismatch, override present, and valid managed profile.
- [x] T148 [US5] Run Go-first verification, schema checks, strict review, pi review, PR-level review, and record Block 17 review disposition before PR closure.

**Checkpoint**: Managed harness users get fail-closed verifier facts for an
approved wrapper or adapter boundary. Unmanaged harness users still get honest
observation-mode gaps instead of being forced into managed enrollment.

## Phase 14: Redaction, Retention, And Forensic Profiles

**Goal**: Make redaction and retention verifier-significant, keep default
recording safe, and add explicit forensic assessment for deliberately retained
sanitized, encrypted, or external evidence.

**Independent Test**: A reviewer can run committed fixtures showing safe
default redaction, digest-only critical evidence rejected by forensic profile,
sanitized excerpt acceptance, encrypted raw reference binding, external artifact
reference binding, unresolved redaction failure or cap, and deterministic
safety-preserving preview/explain output.

**Activation Gate**: Do not implement Block 18 profile behavior until
`blocks/18-redaction-retention-forensic-profiles.md` is explicitly approved.
Forensic retention facts are verifier-derived evidence for external privacy,
incident, compliance, CI, or policy consumers; they are not native legal,
readiness, degradation, merge, release, override, or risk-acceptance decisions.

- [x] T149 [US5] Add Block 18 spec and implementation plan for redaction policies, FR-054 retention modes, raw reference binding, forensic assessment semantics, CLI preview/explain behavior, assessment-result versioning, review-ledger shape, and no-overclaim boundary.
- [x] T150 [US2] Add a redaction policy schema covering rule ids, detector family, allowed FR-054 retention modes, redaction actions, forbidden committed-artifact persistence classes, verifiable authority, critical event family mappings, unresolved-redaction handling, withholding audit, and policy digest/provenance.
- [x] T151 [US2] Extend flight-recorder event/run contracts with redaction policy refs, redaction rule refs, redaction authority, input/redacted payload digests, FR-054 retention mode, retention lifecycle, forensic importance, and raw-reference state with closed enums.
- [x] T152 [US5] Add built-in safe default redaction and retention behavior with stable policy id/digest, pre-write redaction metadata, digest-first command/output retention, and safety-sensitive negative tests.
- [x] T153 [US5] Add explicit forensic retention assessment that rejects digest-only critical evidence with explanatory `capped_to_retention_mode`, unresolved redaction, missing policy, and unavailable raw references according to selected policy.
- [x] T154 [US5] Add encrypted raw reference and external artifact reference verification for SHA-256-or-stronger digest binding, access state, access verification time, key custody state, retention lifecycle, malformed reference failure, access-unverifiable `cannot_verify`, revocation/compromise supersession, and contradiction failure.
- [x] T155 [US4] Add `assess --profile forensic-retention`, `assess preview --profile forensic-retention`, and `assess explain` support with assessment-result schema versioning, deterministic exit codes, condition rows, next actions, profile-selection trace evidence, and safe output.
- [x] T156 [US4] Add committed Block 18 fixtures for safe default, digest-only critical negative, sanitized excerpt positive, encrypted raw reference positive, external raw reference positive, withhold-to-not-assessed, unresolved redaction, missing policy, missing access state, present-but-unverifiable access state, missing key custody state, weak digest, authority self-claim, and raw-reference digest mismatch.
- [x] T157 [US4] Add safety-sensitive negative tests proving redaction, retention, forensic assess, explain, and preview output do not print raw command args, stdout/stderr bodies, prompts, source snippets, credentials, OIDC request tokens, adapter secrets, gateway tokens, model responses, or key material.
- [x] T158 [US5] Run Go-first verification, schema checks, strict review, pi review, PR-level review, and record Block 18 review disposition in `blocks/18-redaction-retention-forensic-profiles-review-ledger.md` before PR closure.

**Checkpoint**: Users get safe default recording plus explicit forensic
assessment facts. Digest-only or unresolved critical evidence cannot be
overclaimed as reconstructable, and raw evidence remains opt-in, referenced,
and verifier-bound rather than committed by default.

**Closure Evidence**: Block 18 was implemented and merged in PR #8. Evidence is
recorded in `docs/research/block-18-implementation-evidence.md` and
`blocks/18-redaction-retention-forensic-profiles-review-ledger.md`, including
local Go-first verification, schema checks, strict implementation review,
PR-level review disposition with no remaining critical or major findings, and
GitHub CI state recorded as `not_assessed` because PR #8 had no checks.

## Phase 15: Adapter Event Contract And Capture Depth

**Goal**: Make adapter-supplied telemetry portable, schema-bound, safe by
default, and explicit about capture-depth limits so harness integrations expose
provenance that CI and git cannot reconstruct without overclaiming unsupported
observers.

**Independent Test**: A reviewer can validate committed generic adapter-event
fixtures showing run lifecycle, task supersession, tool calls, command
correlation, file mutation correlation, model identity provenance, test
observation provenance, missing adapter events, unsupported observers, gateway
`not_integrated`, and safe output without any demo-specific harness, Git host,
or build-system names.

**Activation Gate**: Do not implement Block 19 adapter event behavior until
`blocks/19-adapter-event-contract-capture-depth.md` is explicitly approved.
Adapter events may improve observation depth and managed-profile evidence, but
they are not native merge, release, readiness, degradation, override approval,
risk-acceptance, audit, or forensic-complete decisions. Prompt and raw
model-response capture remains unavailable by default unless Block 18 retention
and redaction policy explicitly permits it.

- [x] T159 [US5] Add Block 19 spec and implementation plan for adapter event contracts, capture-depth taxonomy, gateway provenance, VCS/PR/review references, test provenance semantics, query output expectations, safety boundaries, review-ledger shape, and no-overclaim boundary.
- [x] T160 [US2] Add portable adapter event schema covering `run_started`, `task_locked`, `task_superseded`, `tool_call`, `command_started`, `file_mutation`, `model_call_observed`, `test_observed`, and `run_closed` with producer identity, provenance scope, capture state, same-chain or adapter-bundle binding, correlation refs, redaction/retention metadata, and closed enums.
- [x] T161 [US2] Extend flight-recorder event/run contracts or add a versioned adapter-event bundle so adapter events can be bound to run id, run nonce, event chain digest, source baseline, source commit or tree digest, and managed policy/registry references without requiring a specific harness runtime.
- [x] T162 [US5] Add test observation provenance values and verifier semantics proving `ci_executed` and wrapper/tool-bound execution can satisfy executed-test evidence, while `harness_observed`, `agent_reported`, and `cannot_verify` cannot be silently upgraded.
- [x] T163 [US5] Add generic VCS, PR/MR, review, and source-event reference shapes with provider-neutral fields for repository/source ref, commit or tree digest, change ref, review ref, artifact digest, producer, and `not_assessed` reason.
- [x] T164 [US5] Add adapter capture-depth evaluation that emits deterministic facts for missing adapter events, unsupported observers, not-integrated gateways, unavailable prompt/model-response raw capture, and capture-depth caps without turning gaps into pass/fail policy decisions.
- [x] T165 [US4] Add query-facing output for task drift, task supersession counts, unverified task expansion indicators, unverified claims, missing adapter events, unsupported observers, gateway `not_integrated`, and capture-depth limits while preserving Block 18 safe default redaction and retention behavior.
- [x] T166 [US4] Add committed Block 19 fixtures for valid generic adapter events, missing required adapter event, unsupported observer, gateway `not_integrated`, agent-reported test claim, harness-observed test correlation without execution proof, file mutation/source correlation, task supersession actor attribution, and provider-neutral PR/review refs.
- [x] T167 [US4] Add safety-sensitive negative tests proving adapter event assess/query/preview/explain output does not print raw command args, stdout/stderr bodies, prompts, source snippets, tool-call input/output bodies, adapter configuration, gateway evidence refs, credentials, OIDC request tokens, adapter secrets, gateway tokens, model request/response payloads, PR tokens, authenticated provider URLs, or raw review bodies.
- [x] T168 [US5] Run Go-first verification, schema checks, strict code/correctness review, tracing/evidence review, requirements-vs-implementation review, PR-level review, and record Block 19 review disposition in `blocks/19-adapter-event-contract-capture-depth-review-ledger.md` before PR closure.

**Checkpoint**: Harness and gateway integrations get a stable generic adapter
contract and honest capture-depth facts. Query output may summarize available
evidence, missing telemetry, and unsupported observers, but full reconstruction
is claimed only for evidence actually captured and retained under the selected
profile.

## Phase 16: Post-Block Trust Closure Drift Follow-Ups

**Goal**: Convert recorded trust drift from Blocks 18 and 19 into explicit
tracked work so `not_assessed`, source-bound proof gaps, and PR-review minor
findings do not remain as loose evidence-note prose.

**Independent Test**: A reviewer can inspect this task list and map every
recorded open trust gap from Block 18/19 evidence notes to a concrete task,
owner surface, and non-overclaiming closure condition.

- [x] T169 [US5] Regenerate source-bound release proof after Block 18/19 manifest-subject changes, or record why the current source-bound verifier profile cannot assess it. Closure includes a clean source commit, refreshed active manifest artifact hashes, current `release-proof` output, and explicit `not_assessed` states for external production trust / DSSE steps not assessed by the Go-first local source-bound profile.
- [x] T170 [US5] Add or explicitly document repository CI/check policy for PRs so future PR closure no longer leaves GitHub checks as unowned `not_assessed`. If CI is intentionally absent, add a repo-tracked policy note naming the replacement verification evidence and the remaining trust limitation.
- [x] T171 [US5] Add follow-up fixture or test coverage for PR #9 minor review notes that were accepted as non-blocking: duplicate empty adapter correlation refs, byte-identical pass fixture clarity, and any intentional forward-compatible schema values such as `capture_state: redacted`.
- [x] T172 [US5] Add a drift-to-task rule to the trust workflow so every recorded critical, major, or trust-affecting `not_assessed`/`cannot_verify` drift either blocks closure or lands as an explicit follow-up task with evidence, owner surface, and closure condition.
- [x] T173 [US5] Audit every closed SpecKit task and closed or closure-like block claim for current drift, record machine-checkable pass/fail/not-assessed evidence, and keep source-bound release closure blocked unless the current verifier passes.
- [x] T174 [US5] Replace or explicitly retire historical closed-task verifier references to removed Node/npm/script tooling, including affected self-trace examples, process docs, Block 01/05/06/07/08 ledgers, and manifest subjects. Closure preserves historical evidence as historical while exposing current Go-first verifier commands and remaining `not_assessed` states in `docs/research/historical-verifier-retirement-2026-05-07.md`.

## Phase 17: Forensics Query Pack

**Goal**: Make common forensic reviewer questions executable as safe,
versioned, read-only query packs over existing run, verifier,
forensic-retention, and adapter-capture facts.

**Independent Test**: A reviewer can run `query-pack --pack
forensics-basic-v1` against committed fixtures and see reconstructable evidence,
digest-only caps, redaction issues, capture-depth gaps, task supersession,
unverified claims, unsupported observers, and safe evidence-gap descriptions
without reading raw JSONL manually or receiving a native policy verdict.

**Activation Gate**: Do not implement Block 20 query-pack behavior until
`blocks/20-forensics-query-pack.md` is reviewed through separate Socratic
planes and explicitly approved. Query packs are read-only evidence views. They
must not add raw capture, forensic-completeness badges, legal/audit decisions,
release readiness, merge readiness, risk acceptance, or opaque scores.

- [x] T175 [US5] Add Block 20 spec and implementation plan for versioned forensics query packs, result contract, query list, safety boundaries, row-state preservation, review-ledger shape, and no-overclaim boundary.
- [x] T176 [US2] Add a portable forensics query-pack result schema covering pack id/version, selected run identity, required and optional input artifact SHA-256 digests or path-redacted provider-neutral identifiers, grouped query rows, source fact refs, source condition refs where available, deterministic query-scoped row ids, row states, closed-enum evidence families, reconstructability flags, safe evidence-gap descriptions, and verified output-safety metadata.
- [x] T177 [US5] Add `forensics-basic-v1` derivation over existing Block 09 run/verifier facts, Block 18 forensic-retention facts, and Block 19 adapter-capture facts without creating a new assessment profile or policy verdict.
- [x] T178 [US5] Add deterministic row semantics and a source-state propagation table for reconstructable evidence, digest-only existence evidence, missing required and optional upstream artifacts, missing telemetry, not-integrated gateways, unsupported observers, unresolved redaction, task supersession, unverified claims, unsafe provider refs, unmapped upstream states, and malformed inputs.
- [x] T179 [US4] Add `query-pack --pack forensics-basic-v1` with required `--out` and explain rendering over the JSON result artifact. Explain output must not add claims absent from the result or encode hidden state through ordering, color, indentation, whitespace, or omitted sections.
- [x] T180 [US4] Add committed Block 20 fixtures and a machine-checkable fixture matrix for mixed positive evidence, digest-only reconstruction cap, missing forensic-retention assessment, missing adapter-capture assessment, unsupported observer, unresolved redaction, task supersession, unverified claim, unsafe provider ref, and malformed input. The matrix must enumerate expected query group, row evidence state, evidence family, source ref shape, and reconstructability state where applicable.
- [x] T181 [US4] Add safety-sensitive negative tests proving forensics query-pack and explain output does not print raw command args, command names, executable paths, script paths, unsafe test identifiers, stdout/stderr bodies, prompts, source snippets, tool-call input/output bodies, adapter configuration, gateway evidence refs, credentials, OIDC request tokens, adapter secrets, gateway tokens, PR tokens, authenticated provider URLs, raw model request/response payloads, raw review bodies, unsafe raw-reference access notes, or key material. Test fixtures must use synthetic values and negative leak assertions must not echo candidate secrets in failure output.
- [x] T182 [US5] Run Socratic spec review across product-boundary, tracing/evidence, and privacy/safety planes; record every valid finding in `blocks/20-forensics-query-pack-review-ledger.md` and fix every critical or major finding before implementation approval handoff.
- [x] T183 [US5] After implementation approval, run Go-first verification, schema checks, strict code/correctness review, tracing/evidence review, requirements-vs-implementation review, PR-level review, and record Block 20 review disposition before PR closure.

**Checkpoint**: Forensic reviewers get a stable query-pack answer to common
investigation questions, but `sdp-trace` still emits evidence views only. Any
legal, incident, audit, release, or risk decision remains downstream.

## Phase 18: Cross-Repository Degradation Export

**Goal**: Provide deterministic cross-repository evidence posture exports for
CTO-level movement analysis without issuing a native degradation verdict.

**Independent Test**: A reviewer can run `export cross-repo-posture --profile
cross-repo-evidence-posture-v1` against committed fixtures and inspect
numerators, denominators, dimensions, time windows, source artifact digests,
`not_assessed` counts, input trust states, movement deltas, and refusal rows
without receiving a health score or yes/no degradation answer.

**Activation Gate**: Do not implement Block 21 export behavior until
`blocks/21-cross-repository-degradation-export.md` is reviewed through separate
Socratic planes and explicitly approved. Cross-repository exports are evidence
movement inputs for `sdp-report`, external BI, or policy consumers. They must
not add policy thresholds, dashboards, opaque scores, degradation verdicts,
readiness decisions, or cross-repository raw personal identifiers.

- [x] T184 [US5] Add Block 21 spec and implementation plan for cross-repository posture export, input model, result contract, metric catalog, stale/untrusted input handling, safety boundaries, review-ledger shape, and no-overclaim boundary.
- [x] T185 [US2] Add a portable cross-repository posture export schema covering profile id/version, selected repository/window inputs, input trust states, grouping set id, active grouping keys, metric rows, movement rows, movement summary, refusal rows, source input refs, source artifact digest set hash, closed metric ids, closed dimension names, closed refusal reasons, closed non-comparable reasons, closed output-safety sensitive classes, denominator requirements, metric-row `not_assessed` counts, handoff metadata, and verified output-safety metadata.
- [x] T186 [US5] Add `cross-repo-evidence-posture-v1` aggregation over Block 20 query-pack result artifacts, repository selection manifests, artifact digest manifests, and posture signal manifests without creating a policy threshold, degradation verdict, dashboard score, or native readiness decision.
- [x] T187 [US5] Add deterministic aggregation and movement semantics for missing telemetry, local-only evidence, CI-witnessed evidence, external-witnessed evidence, failed or issue-observed verifier states, overrides, late attach, unsupported observers, not-integrated observers, retention limits, contract changes, stale inputs, untrusted inputs, unsafe labels, missing required inputs, missing optional inputs, absent posture signal fields, and non-comparable windows.
- [x] T188 [US4] Add `export cross-repo-posture --profile cross-repo-evidence-posture-v1` with required `--selection` and `--out`. Explain rendering, if added, must render only the JSON artifact and must not add conclusions absent from the result or encode hidden state through ordering, color, indentation, whitespace, or omitted sections.
- [x] T189 [US4] Add committed Block 21 fixtures and a machine-checkable fixture matrix for valid multi-repo movement, stale input, digest mismatch, missing required input, missing optional input, non-comparable windows, unsupported observer rows, local-only versus CI-witnessed rows from posture signal manifests, external-witnessed rows from posture signal manifests, override rows, late-attach rows, contract-change rows, unsafe labels, unsafe digest-manifest paths, and unsafe external verdict payloads. The matrix must enumerate expected metric rows, numerator, denominator, active grouping keys, dimension key, input trust state, source input refs, source digest set hash, movement comparability, non-comparable reason, movement summary, and refusal rows.
- [x] T190 [US4] Add safety-sensitive negative tests proving cross-repository export and explain output does not print raw command args, command names, executable paths, script paths, unsafe test identifiers, stdout/stderr bodies, prompts, source snippets, tool-call input/output bodies, adapter configuration, gateway evidence refs, credentials, OIDC request tokens, adapter secrets, gateway tokens, PR tokens, authenticated provider URLs, raw model request/response payloads, raw review bodies, unsafe raw-reference access notes, private filesystem paths, or unsafe personal identifiers. Test fixtures must use synthetic values and negative leak assertions must not echo candidate secrets in failure output.
- [x] T191 [US5] Run Socratic spec review across product-boundary, tracing/evidence, and privacy/safety planes; record every valid finding in `blocks/21-cross-repository-degradation-export-review-ledger.md` and fix every critical or major finding before implementation approval handoff.
- [ ] T192 [US5] After implementation approval, run Go-first verification, schema checks, strict code/correctness review, tracing/evidence review, requirements-vs-implementation review, PR-level review, and record Block 21 review disposition before PR closure.

**Checkpoint**: CTO-level consumers get comparable movement facts across
repositories, but `sdp-trace` still emits evidence substrate exports only. Any
degradation, readiness, alerting, or portfolio-risk decision remains downstream.

## Phase 19: Additional CI And Enterprise Witness Profiles

**Goal**: Add provider-neutral CI and enterprise witness profiles without
treating GitHub Actions as the hidden witness model or letting environment
variables alone upgrade trust.

**Independent Test**: A reviewer can run documented witness profile commands or
fixture validation against committed GitLab CI, Buildkite, customer PKI, and
air-gapped examples, then inspect identity source, signing boundary, freshness
boundary, artifact binding, independence state, requested trust scope,
established trust scope, unsupported states, and `not_assessed` or
`cannot_verify` gaps without reading raw provider logs or private customer
material.

**Activation Gate**: Do not implement Block 22 witness profile behavior until
`blocks/22-additional-ci-enterprise-witness-profiles.md` is reviewed through
separate Socratic planes and explicitly approved. Block 22 may normalize
witness evidence and verifier states, but it must not create policy verdicts,
CI enforcement ownership, enterprise support claims, or environment-only trust
upgrades.

- [x] T193 [US5] Add Block 22 spec and review ledger for provider-neutral witness profile semantics, GitLab CI, Buildkite, customer PKI, air-gapped guidance, fixture matrix, safety boundaries, and explicit no-overclaim rules.
- [x] T194 [US2] Add a provider-neutral witness profile/result contract covering profile id/version, identity source, signing boundary, freshness boundary, artifact binding, source/run/policy binding, independence state, unsupported states, requested trust scope, established trust scope, closed reason codes, safe artifact refs, digests, and output-safety state.
- [x] T195 [US2] Add GitLab CI witness profile behavior and fixtures for valid witness evidence, environment-only non-upgrade, source mismatch, missing identity, stale freshness, artifact digest mismatch, unsupported profile version, and malformed inputs.
- [x] T196 [US2] Add Buildkite witness profile behavior and fixtures for valid witness evidence, organization/pipeline/build/job identity, same-job or agent-reported-only topology caps, missing independent signer, artifact digest mismatch, stale freshness, and malformed inputs.
- [x] T197 [US2] Add customer PKI witness profile behavior and fixtures using public certificate or public key identity, authority policy, payload digest binding, and freshness evidence; reject signer mismatch, expired or unverifiable identity, weak digest, missing freshness, and any private-key input.
- [x] T198 [US4] Add air-gapped witness profile documentation and fixtures that distinguish offline-verifiable public-key, timestamp, and artifact-digest evidence from external checks that remain `not_assessed` or `cannot_verify`.
- [x] T199 [US4] Update witness CLI docs and command contracts for closed Block 22 `--kind` values, explicit customer-PKI input flags if needed, safe path handling, deterministic exit states, and explain output that renders only verifier facts.
- [x] T200 [US4] Add safety-sensitive negative tests proving witness JSON and explain output do not print or persist CI tokens, OIDC tokens, JWT bodies, private key material, provider tokens, authenticated provider URLs, raw job logs, private filesystem paths, unsafe personal identifiers, free-text parser errors containing input content, or customer directory, LDAP, SAML, cloud, Vault, HSM, KMS, or PKI payloads.
- [x] T201 [US5] Run Socratic spec review across product-boundary, tracing/evidence, and enterprise/security planes; record every valid finding in `blocks/22-additional-ci-enterprise-witness-profiles-review-ledger.md` and fix every critical or major finding before implementation approval handoff.
- [x] T202 [US5] After implementation approval, run Go-first verification, schema checks, strict code/correctness review, tracing/evidence review, requirements-vs-implementation review, PR-level review, and record Block 22 review disposition before PR closure.

**Checkpoint**: CI and enterprise users get consistent witness profile
semantics across GitHub Actions, GitLab CI, Buildkite, customer PKI, and
air-gapped guidance, while `sdp-trace` still reports verifier facts only. Policy
decisions, enterprise support declarations, and external audit conclusions
remain downstream.

## Phase 20: Demo Repository CI And Trace Pilot

**Goal**: Replace the retired Block 06 toy pilot method with a real
demo-repository pilot that runs `sdp-trace` through CI and produces inspectable
trace, evidence, report, gate, and witness artifacts.

**Independent Test**: A repository observer can follow the Block 24 spec and
pilot report, inspect the demo repository CI run, and verify which artifacts are
observed, local-only, CI-witnessed, `not_assessed`, or `cannot_verify` without
relying on the retired Block 06 `scripts/*` or `npm` validation path.

**Activation Gate**: Do not implement demo-repository work until
`blocks/24-demo-repo-ci-trace-pilot.md` is reviewed through separate Socratic
planes and explicitly approved. Block 24 may prove demo execution and trace
inspectability, but it must not claim external production trust, production
readiness, policy enforcement, or customer deployment readiness.

- [ ] T203 [US5] Run Socratic spec review for Block 24 across product/demo
  credibility, trace/evidence, CI/witness, and privacy/safety planes; record and
  fix every critical or major finding before implementation approval.
- [ ] T204 [US4] Create or select a demo repository and document ownership,
  CI provider, app scope, command surface, artifact retention, and privacy
  boundary without adding a runtime dependency from `sdp-trace` to the demo.
- [ ] T205 [US4] Add a CI-backed `sdp-trace` demo run that captures at least
  one wrapped command, `verify`, `explain`, `report`, `gate`, and `witness`
  output, with explicit `not_assessed` or `cannot_verify` states for missing
  external trust.
- [ ] T206 [US4] Add a demo pilot report mapping trace/report/gate/witness
  artifacts to the nine Block 23 customer questions and separating observed
  movement from missing telemetry.
- [ ] T207 [US5] Add redaction and safety checks proving committed or linked
  demo artifacts do not expose tokens, private paths, raw logs, raw model
  payloads, customer data, or unsafe personal identifiers.
- [ ] T208 [US5] Run implementation and PR-level reviews across
  code/correctness, trace/evidence, and requirements-vs-implementation planes
  before claiming Block 24 closure.

**Checkpoint**: Block 24 is the first real demo-repository proof path. Until it
closes, `sdp-trace` may claim Block 23 MVP command/readiness documentation and
source-bound local release proof, but not demo-repository pilot closure.

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
- **Phase 12**: Depends on Block 14 gate/explain/preview behavior and Block 15 signed-checkpoint replay resistance. It may make protected profile facts fail-closed and CI-friendly, but it must not make native policy decisions.
- **Phase 13**: Depends on Block 13B capture-boundary taxonomy and Block 16 protected gate schema/CLI compatibility. It may make managed wrapper or adapter enforcement fail-closed only for explicitly selected managed profile runs, and it must preserve observation-mode value for unmanaged harnesses.
- **Phase 14**: Depends on Block 09 flight-recorder retention/redaction semantics and Block 17 `assess` command surface. It may make forensic retention facts fail-closed only for explicitly selected forensic profiles, while preserving safe default recording for ordinary users.
- **Phase 15**: Depends on Block 18 redaction/retention safety and Block 17 managed harness adapter policy semantics. It may expand capture depth through generic adapter events, but it must keep prompt/model-response raw capture unavailable by default and expose capture-depth gaps as `missing_telemetry`, `unsupported`, `not_integrated`, `not_assessed`, or `cannot_verify`.
- **Phase 16**: Depends on Blocks 18 and 19 merge evidence. It must not claim
  source-bound or CI-backed closure until the corresponding proof exists; it
  exists to turn recorded drift into tracked work.
- **Phase 17**: Depends on Blocks 18 and 19 because it derives forensic query
  rows from redaction/retention and adapter-capture facts. It must not expand
  raw capture, add a new policy decision, expose unsafe artifact refs, or hide
  missing/capped evidence rows.
- **Phase 18**: Depends on Block 20 because it aggregates query-pack result
  rows across repository/window selections. It must not collapse movement facts
  into degradation, readiness, health-score, alert, rank, grade, or risk
  decisions, and it must refuse or explicitly mark stale, untrusted, missing,
  and non-comparable inputs.
- **Phase 19**: Depends on Blocks 15 and 16 for signed-checkpoint and protected
  witness semantics, Block 17 for managed witness binding, and Block 21 for
  cross-repository witness posture consumption. It must not treat GitHub
  Actions as the hidden witness model, claim broad enterprise CI support, or
  upgrade trust from environment variables alone.
- **Phase 20**: Depends on Block 23 closure and reopens the real pilot gap left
  by the retired Block 06 method. It must use current `sdp-trace` command
  surfaces, CI evidence, and trace artifacts instead of retired `scripts/*` or
  `npm` validation.

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
- T138 can run in parallel with T137 after T136 because managed policy and adapter registry schemas are separate from gate-result versioning.
- T140, T141, and T142 can run in parallel after T139 if their write scopes are split between boundary enrollment, event coverage, and suppression handling.
- T144 and T145 can run in parallel after T137 and T139 because CLI exit behavior and explanation rendering have separate test surfaces.
- T160 and T163 can run in parallel after T159 because adapter event shape and provider-neutral VCS/PR/review refs are separate schema surfaces.
- T162 and T164 can run in parallel after T160 and T161 because test provenance semantics and capture-depth evaluation have separate verifier surfaces.
- T165 and T167 can run in parallel after T164 because query rendering and safety-output assertions have separate command/test surfaces.
- T176 and T180 can run in parallel after T175 because schema/result contract
  and fixture matrix are separate surfaces.
- T177 and T179 can run in parallel after T176 because derivation and CLI
  rendering can keep disjoint write scopes if they share only the schema
  contract.
- T178 and T181 can run in parallel after T177 because row semantics and safety
  leak assertions exercise separate behavior paths.
- T185 and T189 can run in parallel after T184 because schema/result contract
  and fixture-matrix design are separate surfaces.
- T186 and T188 can run in parallel after T185 because aggregation and CLI
  wiring can keep disjoint write scopes if they share only the schema contract.
- T187 and T190 can run in parallel after T186 because movement semantics and
  safety leak assertions exercise separate behavior paths.
- T195 and T196 can run in parallel after T194 because GitLab CI and Buildkite
  profile behavior can keep disjoint fixture and normalization scopes.
- T197 and T198 can run in parallel after T194 because customer PKI validation
  and air-gapped documentation exercise separate witness surfaces.
- T199 and T200 can run in parallel after T195-T198 because command-contract
  documentation and safety leak assertions have separate verification surfaces.
- T204 and T206 can run in parallel after T203 because demo-repository setup and
  report template work have separate write scopes.
- T205 and T207 can run in parallel after T204 because CI trace capture and
  redaction/safety checks have separate verification surfaces.

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
13. Complete Block 16 protected gate profile only after explicit spec approval, keeping enforcement ownership with external CI or policy consumers.
14. Complete Block 17 managed harness profile only after explicit spec approval, keeping managed enrollment opt-in and preserving unmanaged observation-mode value.
15. Complete Block 18 redaction, retention, and forensic profiles only after explicit spec approval, keeping raw evidence opt-in and keeping forensic/legal/risk decisions outside `sdp-trace`.
16. Complete Block 19 adapter event contract and capture-depth expansion only after explicit spec approval, keeping generic adapter events provider-neutral and making missing or unsupported telemetry visible instead of forensic-complete by implication.
17. Complete Phase 16 trust-closure follow-ups before making any broader source-bound release or CI-backed trust claim that depends on Blocks 18/19.
18. Complete Block 20 forensics query pack only after explicit spec approval,
    preserving query-pack output as read-only evidence views over existing
    facts and keeping downstream forensic, legal, audit, release, and risk
    decisions outside `sdp-trace`.
19. Complete Block 21 cross-repository posture export only after explicit spec
    approval, preserving export output as movement facts with raw
    numerator/denominator evidence and keeping degradation interpretation
    outside `sdp-trace`.
20. Complete Block 22 additional CI and enterprise witness profiles only after
    explicit spec approval, preserving provider-neutral witness semantics and
    preventing environment-variable-only trust upgrades.
21. Complete Block 23 MVP closure before treating the retired Block 06 pilot
    method as demoted from MVP closure.
22. Complete Block 24 demo-repository CI and trace pilot before claiming a real
    demo-repo pilot is closed.

### Evidence Discipline

- Do not record harness/model/stack support, readiness, or compatibility as native `sdp-trace` outcomes. Record observed evidence state or an explicitly external verdict input.
- Do not add policy thresholds to `sdp-trace`.
- Keep raw pilot outputs ignored until sanitized.
- Every public claim must link to a file, command, example, or `not_assessed` entry.
