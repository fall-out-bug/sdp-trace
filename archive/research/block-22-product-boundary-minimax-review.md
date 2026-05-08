# Block 22 Product-Boundary Socratic Review

## Findings

### CRITICAL

**1. Customer PKI "without requiring a live customer PKI" is under-specified**
- **Section**: Customer PKI / Product Boundary
- **Problem**: The spec says the customer PKI profile validates "signer authority policy, payload digest binding, and freshness metadata without requiring a live customer PKI" but does not explain how freshness is established without a live service. If the customer PKI profile requires "freshness or timestamp evidence declared by the customer profile," where does that declaration come from, what format is it, and what prevents a self-claimed declaration from being treated as authoritative?
- **Fix Required**: Explicitly name the freshness mechanism for the customer PKI profile. Options include: (a) explicit customer-supplied timestamp file with schema version, (b) declared `valid_until` field in the witness input, or (c) explicit `freshness_policy` field that is NOT a live network call. The spec must state which option applies and reject self-claimed freshness without an external anchor.

---

### MAJOR

**2. "Air-gapped profile" conditional implementation creates ambiguity**
- **Section**: Air-Gapped Profile / Product Boundary
- **Problem**: The spec says air-gapped is "documentation plus fixtures in Block 22 unless a reviewed implementation plan proves a small deterministic validator is useful." This conditional creates two unresolved paths: when is a validator "useful" vs. when is it not? If fixtures exist but no validator exists, what commands can a reviewer actually run to validate the fixtures? If a validator IS added later, does that require a new spec revision?
- **Fix Required**: Explicitly state that the air-gapped profile is documentation-only with no new CLI command, OR define the exact minimal validator that can be added within this block's scope without a separate implementation plan. A reviewer must be able to answer "what can I run?" before approving this block.

**3. Independence state lacks concrete thresholds**
- **Section**: Witness Profile Contract / Independence State
- **Problem**: The contract defines `independence_state` with values including "independent, same pipeline but separate job, same job, local-only, or not_assessed" but provides no concrete criteria for what makes a witness "independent" vs. "same pipeline but separate job." In GitLab CI and Buildkite, a job that runs in the same pipeline as the witness job but in a separate job could still share environment, secrets, or network namespace. The independence boundary is underspecified.
- **Fix Required**: For each profile (GitLab CI, Buildkite, customer PKI), explicitly name what independence state is achievable and what topology would NOT qualify. Example: "GitLab CI cross-job witness in the same pipeline is NOT independent because jobs share CI_REGISTRY and JOB_TOKEN; independence requires separate pipeline refs." This prevents implementers from inventing their own independence thresholds.

**4. Buildkite trust ceiling is unclear**
- **Section**: Buildkite / Non-Goals
- **Problem**: The spec correctly caps Buildkite at "co-located agent and witness topologies that cap trust below `external_witnessed`" but does not name what topology WOULD qualify for `external_witnessed`. Without this, an implementer might claim Buildkite can reach `external_witnessed` under some interpretation.
- **Fix Required**: State explicitly that Buildkite's first implementation caps at `ci_witnessed` maximum and `external_witnessed` is not achievable by the `buildkite-v1` profile. If a future profile could reach `external_witnessed`, name the required topology (e.g., separate Buildkite organization, Buildkite artifact signing integration, or external witness anchor) as a blocked follow-up profile.

**5. CLI boundary does not specify error handling for invalid `--kind`**
- **Section**: CLI Boundary
- **Problem**: The spec shows the command shape `go run ./cmd/sdp-trace witness --kind <profile-kind> --out <file> [--report-dir <dir>] <runs-root-or-run-dir>` but does not specify: (a) whether unknown `--kind` values produce an error or fall back to existing behavior, (b) whether multiple `--kind` values can be supplied, (c) whether the command fails or warns when a profile is selected but required inputs are missing (e.g., customer PKI without required public inputs), or (d) exit code behavior for unsupported profile versions.
- **Fix Required**: Specify deterministic error handling: unknown `--kind` must error with a list of allowed values, missing required inputs for a selected profile must produce `cannot_verify` with a reason, and exit code 2 for usage errors, exit code 3 for `cannot_verify`, exit code 1 for `fail`.

**6. Customer PKI public input handling is underspecified**
- **Section**: Customer PKI / CLI Boundary
- **Problem**: The spec says "If customer PKI needs extra input paths, they must be explicit flags with safe path handling" but does not name the actual flags, required input shapes, or validation behavior. A reviewer cannot evaluate CLI ergonomics without this.
- **Fix Required**: Name the explicit flags for customer PKI inputs. At minimum, this includes: (a) `--customer-pki-authority-policy <path>` for the allowed signer policy, (b) `--customer-pki-certificate` or `--customer-pki-public-key <path>` for declared identity, (c) `--customer-pki-freshness-source <path>` for freshness evidence. State that all paths are rejected if they resolve to private directories, private keys, or runtime-relative paths.

---

### MINOR

**7. `safe_output_classes` in the witness profile contract is unexplained**
- **Section**: Witness Profile Contract / Table
- **Problem**: The contract table includes `safe_output_classes` as a required field but the purpose, allowed values, and relationship to the safety requirements section are not explained in the contract definition. A reviewer reading the table alone cannot derive what safe output classes are or how they differ from `unsupported_states`.
- **Fix Required**: Either (a) add a purpose note in the table clarifying that `safe_output_classes` lists sensitive classes verified absent from JSON and explain output, or (b) move the field to the "Normalized witness result" section where it is actually used, and remove it from the profile contract table.

**8. Fixture matrix does not cover cross-profile integration scenarios**
- **Section**: Fixture Matrix
- **Problem**: The fixture matrix covers individual profile scenarios (valid GitLab, valid Buildkite, valid customer PKI) but does not cover cross-scenario interactions that are likely in real enterprise environments, such as: (a) a GitLab CI run that also has Buildkite artifacts, (b) a customer PKI witness that references GitLab CI source binding, or (c) an air-gapped environment that attempts to consume a Buildkite witness artifact.
- **Fix Required**: Add at least one cross-scenario fixture or explicitly state that cross-profile combinations are out of scope for Block 22. If out of scope, the acceptance criteria should not imply comprehensive cross-profile testing.

**9. Jenkins is mentioned as a follow-up but has no explicit non-goal statement**
- **Section**: Goal / Non-Goals
- **Problem**: The Goal section mentions "Jenkins remains a documented follow-up profile candidate" but the Non-Goals section does not explicitly exclude Jenkins from Block 22 implementation scope. This creates ambiguity: could an implementer claim Jenkins support as part of "enterprise CI" without violating the spec?
- **Fix Required**: Add to Non-Goals: "Jenkins CI support is explicitly out of scope for Block 22. Jenkins follows as a separate reviewed profile candidate because its plugin, controller, agent, and credential topologies create more overclaim risks."

**10. Block 21 integration path is not referenced**
- **Section**: Product Boundary / Dependencies
- **Problem**: Block 22 depends on Block 21 for cross-repository witness posture consumption, and Block 21 depends on witness profiles from Block 22. However, neither the Block 22 spec nor the acceptance criteria reference how cross-repository posture exports would handle mixed witness profile sources (GitLab CI repo + GitHub Actions repo + Buildkite service).
- **Fix Required**: Add a note stating whether Block 21 posture exports can aggregate across different witness profile kinds and what the independence/trust scope rules are for mixed sources. At minimum, state that posture exports use the normalized witness result from Block 22 and that the trust scope is the minimum across source profiles.

**11. "No raw token, OIDC JWT, CI secret" exclusion scope is unclear for JSON output**
- **Section**: Safety Requirements
- **Problem**: The safety requirements correctly exclude raw tokens, OIDC JWTs, CI secrets, and private material from output. However, in the context of the witness profile contract, witness inputs MAY contain OIDC JWTs or CI tokens as inputs to the witness normalization process (even if the output must exclude them). The spec does not clarify whether the input parsing layer must handle OIDC JWT parsing without persisting the parsed JWT body.
- **Fix Required**: Clarify that input parsing for OIDC JWTs or CI tokens is allowed for identity extraction (subject, issuer, audience) but the raw JWT body, token body, or secret value must not be persisted in the normalized result or any intermediate artifact.

---

## Summary Table

| ID | Severity | Section | Fix Required |
|----|----------|---------|--------------|
| 1 | CRITICAL | Customer PKI / Product Boundary | Name explicit freshness mechanism without live PKI service |
| 2 | MAJOR | Air-Gapped Profile | Clarify whether CLI validator exists or is documentation-only |
| 3 | MAJOR | Independence State | Define concrete independence thresholds per profile |
| 4 | MAJOR | Buildkite trust ceiling | State explicitly that Buildkite caps at ci_witnessed maximum |
| 5 | MAJOR | CLI error handling | Specify exit codes, unknown kind behavior, and missing input behavior |
| 6 | MAJOR | Customer PKI CLI flags | Name explicit flag shapes for required public PKI inputs |
| 7 | MINOR | safe_output_classes | Explain purpose in contract table or move to result section |
| 8 | MINOR | Fixture matrix | Add cross-profile fixture or explicitly state out of scope |
| 9 | MINOR | Jenkins reference | Add explicit Jenkins non-goal statement |
| 10 | MINOR | Block 21 integration | Reference cross-repository posture aggregation for mixed sources |
| 11 | MINOR | OIDC JWT input handling | Clarify input parsing vs. output persistence scope |

---

## Verdict

**REVISE**

Block 22 addresses a legitimate product gap: GitHub Actions is not the universal CI environment and the witness boundary should not assume it. The provider-neutral contract approach is sound, and the safety requirements are appropriately strict.

However, the spec is not ready for implementation approval because:

1. **CRITICAL finding #1** (customer PKI without live PKI freshness mechanism) would allow an implementer to accept self-claimed freshness declarations as authoritative, which directly contradicts the "no environment variables alone" rule.

2. **MAJOR finding #2** (air-gapped conditional implementation) creates an undefined path where reviewers cannot answer "what can I run?" before approving.

3. **MAJOR findings #3, #4, #5, #6** collectively mean an implementer would lack sufficient guidance on independence thresholds, Buildkite trust ceilings, CLI error handling, and customer PKI input shapes—all of which are required before a PR can be reviewed against the acceptance criteria.

The acceptance criteria themselves are well-structured (especially SC-052 through SC-054) and the fixture matrix coverage is comprehensive for individual profiles. The revision needed is to close the underspecified freshness, air-gapped validator, CLI ergonomics, and independence threshold gaps before implementation starts.
