# Block 22: Tracing/Evidence Socratic Review Findings

## Critical Findings

### C01 — No shared closed reason-code registry
**Section**: Witness Profile Contract → normalized witness result  
**FR trace**: FR-119  
**Problem**: The spec states profiles "share verifier states and reason-code semantics" but never enumerates the shared reason codes. Without a closed enum or registry, every profile can invent ad-hoc reason strings, breaking FR-119 compliance and making Cross-Repo Export (Block 21) aggregation impossible.  
**Fix required**: Add an explicit closed reason-code registry (table or enum) covering at least: missing identity, identity mismatch, missing freshness, stale freshness, missing signer authority, signer mismatch, missing artifact digest, artifact digest mismatch, missing source binding, source mismatch, missing run binding, run mismatch, missing policy binding, policy mismatch, environment-only insufficient, unsupported version, malformed input, unsafe output candidate, private-key input rejected. Map each to the verifier state (`fail`, `not_assessed`, `cannot_verify`) it produces.

### C02 — Witness CLI exit codes undefined
**Section**: CLI Boundary  
**FR trace**: FR-120, FR-076, FR-084 (parallel precedent)  
**Problem**: Block 16 and Block 17 define deterministic exit codes (0 pass, 1 fail, 2 usage error, 3 cannot_verify). The witness command has no exit-code contract. CI consumers and Block 21 export cannot rely on deterministic exit behavior.  
**Fix required**: Define exit codes for `sdp-trace witness`: `0` for established trust scope satisfying the profile, `1` for verifier `fail`, `2` for usage/argument errors, `3` for `cannot_verify` or `not_assessed` with no profile failure. Align with Block 16/17 precedent.

### C03 — Trust scope determination rules missing
**Section**: Witness Profile Contract → normalized witness result  
**FR trace**: FR-118, FR-120  
**Problem**: The spec lists trust scope values (`local_observed`, `ci_witnessed`, `external_witnessed`, etc.) but never defines the rules that determine which scope is established. A verifier needs a decision table: given identity state + signer authority state + freshness state + artifact binding state + independence state → established trust scope. Without this, profiles can silently overclaim or underreport.  
**Fix required**: Add a trust-scope determination matrix table showing the minimum verifier-state combination for each trust scope level. Explicitly state that no combination of environment-variable-only inputs may produce `ci_witnessed` or `external_witnessed`.

## Major Findings

### M01 — Customer PKI explicit input contract missing
**Section**: Profile Requirements → Customer PKI + CLI Boundary  
**FR trace**: FR-121  
**Problem**: The spec says customer PKI uses "declared public certificate or public key identity" and "authority policy" but does not define the input shape: file paths, inline PEM, certificate subject strings, or JSON policy structure. The CLI boundary says "extra input paths, explicit flags" but names none. Fixtures cannot be designed without this.  
**Fix required**: Define explicit customer PKI input flags (e.g. `--public-cert`, `--public-key`, `--authority-policy`, `--payload-digest`, `--freshness-evidence`) with safe path handling rules, allowed formats (PEM/JSON/DER), and explicit refusal of implicit directory scanning or private-key discovery.

### M02 — Air-gapped profile identity missing
**Section**: Profile Requirements → Air-Gapped Profile + CLI Boundary  
**FR trace**: FR-122, SC-053  
**Problem**: The spec says "documentation plus fixtures" but does not state whether `air-gapped-v1` is a valid `--kind` identifier. Without a profile_id, Block 21 cross-repo export cannot reference an air-gapped witness row, and SC-053 "reviewer can run documented witness profile commands" is unverifiable.  
**Fix required**: Either assign `air-gapped-v1` as a valid `--kind` with deterministic `cannot_verify` behavior for all external checks, or explicitly document it as documentation-only with no CLI command—and state which SC/FR the documentation alone satisfies versus which requires a verifier.

### M03 — Source/run/policy binding states referenced in result but absent from profile contract table
**Section**: Witness Profile Contract (table omits source binding, run binding, policy binding) vs. normalized result (requires them)  
**FR trace**: FR-118, FR-119  
**Problem**: The contract table declares `artifact_binding` and `independence_state` but the normalized result exposes six additional binding states: identity, signer authority, freshness, artifact binding, source binding, run binding, policy binding. These are profile-specific evaluation dimensions that need to appear in the profile contract or be explicitly derived from it.  
**Fix required**: Add `source_binding`, `run_binding`, and `policy_binding` to the profile contract table with the same structure as `artifact_binding`: what they cover, how they are validated, and what states they produce.

### M04 — Independence state closed values inconsistent with trust-scope vocabulary
**Section**: Witness Profile Contract table (`independence_state`) vs. Product Boundary vocabulary  
**FR trace**: FR-118  
**Problem**: The contract table lists independence as "independent, same pipeline but separate job, same job, local-only, or not_assessed" but these overlap with trust-scope values (`local_observed`, `ci_witnessed`, `external_witnessed`) without clarifying the relationship. Independence is a witness topology dimension; trust scope is an evidence-strength classification. The spec needs to distinguish them.  
**Fix required**: Define independence_state as a closed enum of witness topology (e.g., `external_independent`, `ci_isolated_job`, `ci_shared_pipeline`, `local_only`, `not_assessed`, `cannot_verify`) and add a separate paragraph explaining how independence_state constrains but does not equal trust scope.

### M05 — Fixture matrix omits expected-outputs-per-fixture enumeration
**Section**: Fixture Matrix  
**FR trace**: SC-052, SC-053  
**Problem**: The spec lists fixture cases but says "Each fixture must state expected verifier result, trust scope, reason codes..." without actually enumerating these expectations. Reviewers cannot verify fixture completeness without a table mapping each fixture to its expected outputs.  
**Fix required**: Add a fixture expectation table (or reference an appendix table in the spec) with columns: fixture name, profile_id, verifier state, trust scope, reason codes, identity state, signing boundary state, freshness state, artifact binding state, output-safety state. At minimum enumerate the 16+ fixtures listed.

### M06 — No freshness evaluation semantics defined
**Section**: Profile Requirements (all profiles reference freshness_boundary)  
**FR trace**: FR-119  
**Problem**: "Freshness" is invoked by every profile but the spec never defines what makes freshness valid, stale, or missing. Is it an expired timestamp? A missing nonce? A sequence gap? Without evaluation semantics, each profile can interpret freshness differently, violating FR-119 shared semantics.  
**Fix required**: Add a Freshness Evaluation section defining: valid (timestamp within policy window AND associated with correct run AND not superseded), stale (timestamp outside policy window or superseded by later evidence), missing (no timestamp/nonce/sequence provided per profile rules). Map each to the appropriate verifier state.

### M07 — Cross-surface consumption boundary underspecified
**Section**: Product Boundary + cross-reference to Blocks 16/17/21  
**FR trace**: FR-119, SC-052  
**Problem**: The spec says gate, protected gate, managed harness, and cross-repo surfaces "can consume the normalized witness result without becoming policy owners" but does not define the consumption mapping. Which fields flow into Block 16 protected gate? Which into Block 21 export grouping?  
**Fix required**: Add a Cross-Surface Consumption section with a mapping table: profile field → Block 16 gate-result field, → Block 17 assessment-result field, → Block 21 export grouping key. State explicitly which fields are consumed as verifier facts vs. ignored by each surface.

## Minor Findings

### N01 — Profile versioning format and compatibility rules undefined
**Section**: Witness Profile Contract table (`profile_version`)  
**Problem**: The table says "Semantic changes require a new version" but does not define version format (`semver`? `integer`?), additive vs. breaking change rules, or backward-compatibility behavior.  
**Fix required**: Define profile version as `major.minor`, specify that minor changes are additive field additions that verifiers may ignore, and major changes require profile_id version bump (e.g., `gitlab-ci-v2`).

### N02 — safe_output_classes not populated with concrete values
**Section**: Witness Profile Contract table (`safe_output_classes`)  
**Problem**: The Safety Requirements section lists 15+ forbidden output classes, but the contract table leaves `safe_output_classes` as an abstract declaration without concrete values. Verification cannot assert absence without defined classes.  
**Fix required**: Populate `safe_output_classes` with the concrete forbidden classes from the Safety Requirements section (CI tokens, OIDC tokens, JWT bodies, private keys, provider tokens, authenticated URLs, raw job logs, private paths, unsafe personal identifiers, free-text parser errors with input content, customer infrastructure payloads). Define the verifier that checks output serialization against these classes.

### N03 — Buildkite signed-boundary vs. agent-reported distinction needs tighter framing
**Section**: Profile Requirements → Buildkite  
**Problem**: The spec says "agent-reported metadata versus signed or authority-bound witness facts" and "co-located agent and witness topologies cap trust below external_witnessed" but does not state what constitutes a "signed" Buildkite fact in the absence of a Buildkite SDK dependency.  
**Fix required**: Clarify that Buildkite v1 profiles treat all Buildkite-provided metadata as agent-reported/ci_witnessed maximum, and external_witnessed requires an independent signer (e.g., external checkpoint, customer PKI attestation) bound to the Buildkite run identity.

### N04 — No explicit fixture for Buildkite stale freshness or missing independent signer
**Section**: Fixture Matrix vs. Tasks T196  
**Problem**: T196 lists "stale freshness" and "missing independent signer" as Buildkite fixture requirements, but these do not appear in the spec's fixture matrix table.  
**Fix required**: Add "Buildkite stale freshness" and "Buildkite missing independent signer" rows to the fixture matrix with explicit expected verifier states.

### N05 — No explicit fixture for GitLab missing identity or artifact digest mismatch
**Section**: Fixture Matrix vs. Tasks T195  
**Problem**: T195 lists "missing identity" and T196 references Buildkite artifact digest mismatch, but the fixture matrix does not explicitly enumerate "GitLab missing identity" or a separate "GitLab artifact digest mismatch" fixture (only "GitLab source commit mismatch" is listed).  
**Fix required**: Add "GitLab missing identity" and "GitLab artifact digest mismatch" rows to the fixture matrix with explicit expected verifier states.

---

## Verdict: **REVISE**

Three critical gaps—missing reason-code registry, undefined CLI exit codes, and absent trust-scope determination rules—block deterministic verifier semantics across all four profiles. Seven major gaps around input contracts, profile identification, binding-state symmetry, fixture enumeration, freshness evaluation, cross-surface consumption, and independence/trust-scope disambiguation must be resolved before implementation can begin with FR-119/FR-120 compliance.
