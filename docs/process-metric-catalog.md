# Process Metric Catalog

This catalog defines portable metric names for `sdp-trace` observations over time.

`sdp-trace` records metric movement. It does not decide whether a value is good, bad, sufficient, degrading, ready, blocked, or acceptable. Thresholds and interpretations belong to `sdp-gate` or another external policy consumer.

## Required Dimensions

Every metric sample should include dimensions that are available for the scope:

- `repository`
- `scope`
- `spec_id`
- `block`
- `team`
- `harness`
- `model_family`
- `model_version`
- `stack`
- `build_system`
- `window_start`
- `window_end`

Unavailable dimensions must be omitted from `dimensions` and recorded through `not_assessed` or `unavailable_fields` in the surrounding evidence/provenance package when they matter for interpretation.

## Catalog

| Metric name | Unit | Meaning | Collection method | Evidence source | `not_assessed` rule | Privacy risk |
|---|---|---|---|---|---|---|
| `contract_task_completion_ratio` | ratio | Completed SpecKit tasks divided by in-scope SpecKit tasks for a block or phase. | Count checked tasks in `tasks.md` for the declared scope. | `specs/**/tasks.md`, review notes. | Use `not_assessed` if the task scope is ambiguous or task status cannot be read from committed SpecKit artifacts. | Low. |
| `evidence_coverage_ratio` | ratio | Evidence-backed required artifacts divided by required artifacts for the scope. | Count required evidence items and committed evidence refs. | Examples, validation summaries, command evidence. | Use `not_assessed` if required evidence list is undefined. | Medium when refs point to private systems. |
| `not_assessed_item_count` | count | Count of explicit `not_assessed` entries in the package. | Count `not_assessed[]`, `unavailable_fields[]`, and `not_assessed_reason` occurrences for the scope. | Assessment input, provenance records, metric streams. | Use `not_assessed` if the package cannot be parsed. | Low. |
| `schema_validation_state` | enum | Whether schema validation evidence exists for the scope. Values: `passed`, `failed`, `not_assessed`. | Read validation command evidence. | `go test ./...`, `jq empty schema/*.json`, `go run ./cmd/sdp-trace validate-fixtures <fixture-root>`, validation summary. | Use `not_assessed` if validation command was not run or output is unavailable. | Low. |
| `artifact_safety_scan_state` | enum | Whether committed artifact safety scan evidence exists. Values: `passed`, `failed`, `not_assessed`. | Read safety scan command evidence. | Denylist scan evidence, redaction review summary, committed artifact review notes. | Use `not_assessed` if scan was not run or exclusions changed without review. | Medium because findings may mention sensitive markers. |
| `manifest_digest_match_state` | enum | Whether contract manifest artifact digests matched the checkout. Values: `matched`, `mismatch`, `not_assessed`. | Run or read manifest verification evidence. | `go run ./cmd/sdp-trace release-proof --manifest <manifest> --out <file>`, contract release verification. | Use `not_assessed` if manifest or verification output is missing. | Low. |
| `release_attestation_state` | enum | Current release proof state. Values: `schema_valid`, `digest_verified`, `locally_attested`, `externally_attested`, `production_release_verified`, `not_assessed`. | Read release-proof or external trust verification output. | `go run ./cmd/sdp-trace release-proof --manifest <manifest> --out <file>`, release verification record, external trust evidence. | Use `not_assessed` for any proof state not backed by committed or externally inspectable evidence. | Medium when external signer identity is private. |
| `review_contradiction_count` | count | Number of unresolved contradictions raised by structured reviewers or judges. | Count unresolved findings in review result artifacts. | `specs/**/blocks/*critic*.json`, `*judge*.json`, resolution notes. | Use `not_assessed` if review artifacts are absent or malformed. | Low. |
| `review_independence_state` | enum | Whether review evidence used a different actor/model/provider from authoring. Values: `independent`, `same_actor`, `not_assessed`. | Compare provenance and review metadata. | Review result JSON, provenance records. | Use `not_assessed` if actor/model/provider identity is unavailable. | Low. |
| `source_scope_change_count` | count | Count of changed files or artifacts in the declared scope. | Count changed committed artifacts or trace change refs. | Git diff summary, trace snapshot, evidence events. | Use `not_assessed` if change refs are unavailable. | Low. |
| `command_verification_count` | count | Count of verification commands with recorded evidence. | Count command evidence events with `status: success` or `failure`. | Evidence events, validation summaries. | Use `not_assessed` if command evidence is not recorded. | Low. |
| `harness_identity_coverage_ratio` | ratio | Harness runs with recorded harness identity divided by in-scope harness runs. | Count pilot/self-trace provenance records. | Provenance records, run-cards. | Use `not_assessed` if no harness run occurred or logs are unavailable. | Medium if raw logs contain private prompt content. |
| `model_identity_coverage_ratio` | ratio | Model-involved records with model family/version captured divided by model-involved records. | Count provenance records with model fields. | Provenance records, harness logs. | Use `not_assessed` if provider does not expose stable model identity. | Medium. |
| `stack_detection_coverage_ratio` | ratio | Stack/build-system scopes with evidence-backed detection divided by stack scopes assessed. | Count stack evidence refs per scope. | `BUILD`, `BUILD.bazel`, `MODULE.bazel`, `.bazelrc`, source files, run-cards. | Use `not_assessed` if stack files are unavailable or not inspected. | Low. |
| `build_test_evidence_state` | enum | Whether build/test evidence exists for a scoped run. Values: `passed`, `failed`, `skipped`, `not_assessed`. | Read command evidence for build/test execution. | CI logs, local command summaries, run-card outputs. | Use `not_assessed` if build/test command was not run and no explicit skip reason exists. | Medium when logs include proprietary paths. |
| `unsupported_claim_count` | count | Count of claims made by models/harnesses that are not backed by inspectable evidence. | Count unsupported-claim observations. | Model output evidence, review notes. | Use `not_assessed` if model output is unavailable. | Medium because raw output may include prompts. |

## Rules

1. Metric names in committed examples should come from this catalog unless a spec explicitly introduces a new metric and updates this file in the same change.
2. Do not encode thresholds, target values, traffic-light labels, or pass/fail decisions in this catalog.
3. A metric sample with missing evidence uses `assessment_state: "not_assessed"`, `value: null`, and a concrete `not_assessed_reason`.
4. A metric stream with any partial or `not_assessed` sample or comparison must expose stream-level `assessment_state: "partial"` or `not_assessed`.
5. Cross-team or cross-window comparison is only meaningful when `metric_name`, `unit`, and relevant dimensions match.
