# Block 21: Cross-Repository Degradation Export

Status: spec delta and implementation plan drafted; Socratic spec review
required before implementation approval handoff.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13-product-gap-closure-roadmap.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13b-capture-boundary-state-dx-baseline.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/18-redaction-retention-forensic-profiles.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/19-adapter-event-contract-capture-depth.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/20-forensics-query-pack.md`

## Goal

Provide a portable evidence export that lets a CTO or downstream reporting
consumer compare evidence posture across repositories, teams, services,
harnesses, change types, and time windows without `sdp-trace` issuing a native
degradation verdict.

The product outcome is not a dashboard and not a score. It is a deterministic
aggregate input for `sdp-report`, external BI, or a policy consumer that owns
thresholds and interpretation.

## Problem

Block 20 makes single-run forensic evidence easier to inspect, but it does not
answer the executive question: "Is our delivery evidence posture moving in the
wrong direction across the portfolio?"

The weak framing would be "add degradation analytics." That is dangerous
because it pushes policy thresholds, opaque scoring, and organizational
judgment into `sdp-trace`. The correct framing is "export comparable movement
facts with explicit denominators, evidence digests, stale-input handling, and
`not_assessed` visibility so another layer can decide whether the movement is
acceptable."

## Non-Goals

- No native yes/no degradation, improvement, readiness, merge, release,
  override, incident, audit, legal, or risk-acceptance verdict.
- No opaque health score, trust score, risk score, badge, rank, grade, color
  state, or hidden severity ordering.
- No policy thresholds, alert thresholds, service-level objectives, or
  "acceptable" / "unacceptable" labels.
- No dependency on `sdp-report`, `sdp-gate`, GitHub, GitLab, Jira, Buildkite,
  Jenkins, BigQuery, Snowflake, Datadog, Grafana, Tableau, or any named harness.
- No raw command, test, prompt, model-response, stdout/stderr, source-snippet,
  tool payload, review body, provider URL, private path, or credential export.
- No cross-repository joins on raw personal identifiers. Human accountability
  may be represented only through declared provider-neutral role or team
  dimensions already present in safe upstream artifacts.
- No deterministic replay claim for underlying runs. The export is an aggregate
  evidence index over selected artifacts, not a replay engine.

## Product Boundary

Block 21 may emit aggregate movement facts:

- export id, schema version, selected aggregation profile id/version, producer,
  generated time, and input selection metadata;
- aggregation window start/end and comparison window start/end when selected;
- dimensions: repo, team, service, harness, change type, observer state,
  evidence family, source condition state, and optional provider-neutral labels
  declared safe by upstream artifacts;
- metrics with numerator, denominator, unit, aggregation rule, source artifact
  digest set, and `not_assessed` count;
- state distribution counts for `present`, `issue_observed`, `not_assessed`,
  `cannot_verify`, `missing_telemetry`, `unsupported`, `not_integrated`, and
  `retention_limited`;
- movement values such as current count, previous count, current ratio,
  previous ratio, and signed delta when both windows are comparable;
- input trust state for every selected repository/window input:
  `trusted_input`, `stale_input`, `untrusted_input`, `cannot_verify_input`, or
  `not_assessed_input`;
- refusal rows when required inputs are missing, stale, malformed, untrusted, or
  non-comparable;
- export-level movement summary counts for comparable and non-comparable
  movement rows;
- handoff metadata for `sdp-report` or external BI consumers.

Block 21 must not decide whether a movement is degradation. If a downstream
tool emits that decision, it must be recorded as an external verdict input with
producer and policy reference.

## Input Model

The first export profile is `cross-repo-evidence-posture-v1`.

Required inputs per repository/window:

| Input | Purpose | Missing or malformed state |
| --- | --- | --- |
| Block 20 query-pack result | Source rows for run, witness, retention, adapter capture, task, command, file, test, supersession, and claim evidence | `cannot_verify_input` for that repository/window |
| Repository selection manifest | Declares repo id, safe display label, service/team mapping, selected time window, and expected artifact refs | `cannot_verify_input` |
| Artifact digest manifest | Binds selected input files to SHA-256-or-stronger digests and source baseline when available | `untrusted_input` when digest mismatch is observed; `not_assessed_input` when digest verification is not possible |
| Posture signal manifest | Declares optional closed posture signals not present in the Block 20 row contract, such as witness scope, override presence, late attach, and contract change markers | Missing referenced signals become `not_assessed_input`; malformed signals become `cannot_verify_input` for affected metrics |

Optional inputs:

- external repository/team/service catalog with safe labels;
- external verdict input artifacts from `sdp-gate`, `sdp-report`, or another
  named policy consumer;
- CI or external witness summaries from later witness profiles.

Optional inputs that are absent must create visible `not_assessed_input` or
dimension-level `not_assessed` facts when the selected profile references them.
They must not silently remove rows from denominators.

The repository selection manifest must also carry:

- `freshness_boundary` as an ISO-8601 duration or absolute cutoff timestamp;
- `dimension_exposure_policy`, naming which dimensions may appear in output for
  this export;
- safe label declarations for repo, team, service, harness, and change type;
- selected grouping set id.

Block 21 must validate safe labels itself. Upstream "safe" declarations are
inputs, not authority. Safe labels are limited to provider-neutral slugs matching
`[a-z0-9][a-z0-9._-]{0,63}` unless a future reviewed profile defines a stricter
customer-private equivalent. Labels that fail validation, contain provider URLs,
private paths, email-like strings, token-like strings, or raw personal
identifiers are refused before aggregation.

## Result Contract

The export result should be a stable JSON artifact with:

- `schema_version`;
- `export_profile_id` and `export_profile_version`;
- `export_id`, `producer`, `generated_at`, and selected profile digest when
  available;
- `input_selection` listing every selected repository/window, expected input
  artifact kind, path-redacted artifact id, digest when readable, and input
  trust state;
- `dimensions` using closed names for repo, team, service, harness,
  change_type, time_window, evidence_family, row_evidence_state,
  source_condition_state, observer_state, and input_trust_state;
- `grouping_set_id` and ordered `active_grouping_keys`. The first profile
  supports only `repo_window_v1`, `team_service_window_v1`, and
  `harness_change_window_v1`. Non-active dimensions may be present as row-level
  safe metadata, but they do not affect numerator or denominator grouping;
- `metric_rows`, each with deterministic row id, metric id, metric version,
  numerator, denominator, unit, time window, dimension values, deterministic
  `dimension_key`, source input refs, source artifact digest set hash,
  `source_field_state`, `not_assessed_count`, and input trust state summary;
- `movement_rows`, each with current window value, previous window value, delta,
  metric id/version, `dimension_key`, current metric row ref, previous metric
  row ref, closed-enum comparison basis, comparable boolean, and closed-enum
  non-comparable reason when comparison is refused;
- `movement_summary` with counts by comparable state and non-comparable reason;
- `refusal_rows` for stale, malformed, untrusted, unsafe-label, unsupported, or
  non-comparable inputs, using closed reason codes only;
- `handoff` naming the intended consumer contract without requiring that
  consumer to exist in this repository;
- `output_safety` listing sensitive classes verified absent from serialized
  JSON and explain output.

Closed `refusal_reason` values:

- `stale_input`;
- `malformed_input`;
- `untrusted_input_digest_mismatch`;
- `unsafe_label`;
- `unsupported_input`;
- `missing_required_input`;
- `missing_optional_input`;
- `non_comparable_metric_version`;
- `non_comparable_dimension_key`;
- `non_comparable_denominator_basis`;
- `non_comparable_input_trust_rule`;
- `non_comparable_missing_window`;
- `output_safety_violation`.

Closed `output_safety.sensitive_class` values:

- `raw_command_args`;
- `command_name_or_path`;
- `unsafe_test_identifier`;
- `stdout_stderr_body`;
- `prompt_body`;
- `source_snippet`;
- `tool_payload`;
- `adapter_configuration`;
- `gateway_evidence_ref`;
- `credential_or_token`;
- `authenticated_provider_url`;
- `model_request_response_payload`;
- `raw_review_body`;
- `unsafe_raw_reference_note`;
- `private_filesystem_path`;
- `unsafe_personal_identifier`;
- `unsafe_label`;
- `raw_digest_manifest_path`;
- `free_text_exception_or_refusal_reason`.

Metric ids must use a closed initial catalog. `metric_version` is part of the
profile contract; adding, removing, or changing a metric denominator, numerator,
source field, or aggregation rule requires a new export profile version.

| Metric id | Source field | Numerator | Denominator |
| --- | --- | --- | --- |
| `missing_telemetry_rows` | Block 20 row evidence state | rows with evidence state `missing_telemetry` | all selected query rows in the active grouping set for the same window |
| `not_assessed_rows` | Block 20 row evidence state | rows with evidence state `not_assessed` | all selected query rows in the active grouping set for the same window |
| `cannot_verify_rows` | Block 20 row evidence state | rows with evidence state `cannot_verify` | all selected query rows in the active grouping set for the same window |
| `unsupported_observer_rows` | Block 20 row evidence state or posture signal manifest observer state | rows with evidence state `unsupported` or validated posture signal `observer_state=unsupported` | all selected query rows in the active grouping set for the same window |
| `not_integrated_rows` | Block 20 row evidence state | rows with evidence state `not_integrated` | all selected query rows in the active grouping set for the same window |
| `retention_limited_rows` | Block 20 row evidence state | rows with evidence state `retention_limited` | all selected query rows in the active grouping set for the same window |
| `local_only_evidence_rows` | posture signal manifest witness scope | rows with validated posture signal `witness_scope=local_only` | all selected query rows in the active grouping set for the same window |
| `ci_witnessed_evidence_rows` | posture signal manifest witness scope | rows with validated posture signal `witness_scope=ci_witnessed` | all selected query rows in the active grouping set for the same window |
| `external_witnessed_evidence_rows` | posture signal manifest witness scope | rows with validated posture signal `witness_scope=external_witnessed` | all selected query rows in the active grouping set for the same window |
| `issue_observed_rows` | Block 20 row evidence state | rows with evidence state `issue_observed` | all selected query rows in the active grouping set for the same window |
| `override_rows` | validated external verdict input or posture signal manifest override marker | rows with closed override marker `override_present` | all selected query rows in the active grouping set for the same window |
| `late_attach_rows` | posture signal manifest late attach marker | rows with closed marker `late_attach_observed` | all selected query rows in the active grouping set for the same window |
| `contract_change_rows` | posture signal manifest contract change marker | rows with closed marker `contract_change_observed` | all selected query rows in the active grouping set for the same window |

If a metric's source field is absent, malformed, unsafe, or not integrated for a
selected repository/window, the metric row must still be emitted with
`source_field_state` set to `not_assessed`, `cannot_verify`, or `unsupported`
and the denominator preserved. The implementation must not synthesize witness,
override, late-attach, or contract-change facts from Block 20 fields that do not
carry those facts.

`not_assessed_count` on a metric row counts selected rows in that metric row's
active grouping set and time window whose contribution to that metric is
explicitly `not_assessed`. For metrics derived directly from Block 20 row
evidence state, this is the count of Block 20 rows with evidence state
`not_assessed` in the denominator. For metrics derived from posture signal
manifests or external verdict inputs, this is the count of denominator rows for
which the required source field is absent because the source is
`not_assessed_input`. It is not a count of refused, malformed, stale, or
untrusted inputs; those appear in `refusal_rows` and `input_selection`.

Every metric row must expose raw numerator and denominator. Ratios may be
included only as derived display aids and must never replace the raw counts.

## Aggregation Rules

Aggregation is deterministic and refusal-first:

- Active grouping keys are fixed by the selected `grouping_set_id`; the
  denominator is the count of selected Block 20 query rows whose safe dimension
  values match that ordered grouping key and time window.
- Inputs with digest mismatch are `untrusted_input` and cannot contribute to
  metric numerators or denominators unless the profile explicitly asks for a
  separate untrusted-input count.
- Inputs older than the selected freshness boundary are `stale_input` and must
  be refused or counted only in stale-input metrics.
- Missing required inputs create `cannot_verify_input` refusal rows.
- Missing optional inputs referenced by the profile create `not_assessed_input`
  rows or dimension-level `not_assessed` facts.
- Rows with `not_assessed`, `cannot_verify`, `missing_telemetry`,
  `unsupported`, `not_integrated`, or `retention_limited` remain visible in
  denominators.
- Current and previous windows are comparable only when metric id/version,
  dimension set, denominator basis, and input trust rules match. Otherwise the
  movement row sets `comparable: false`.
- `comparison_basis` uses a closed enum: `same_profile_metric_dimension_window`,
  `non_comparable_metric_version`, `non_comparable_dimension_key`,
  `non_comparable_denominator_basis`, `non_comparable_input_trust_rule`, or
  `non_comparable_missing_window`.
- `non_comparable_reason` uses the same closed enum without free-text payload.
- `refusal_reason` and `output_safety.sensitive_class` use the closed enums in
  the Result Contract; schema and tests must reject any other string.
- No aggregate may use hidden weighting. Counts are either simple row counts or
  explicitly declared distinct counts over closed safe identifiers.
- Source input refs are sorted by input selection id. Each metric row records
  source input refs and a SHA-256 digest over the sorted digest set; full per
  input digests live in `input_selection` to keep metric rows deterministic and
  bounded.

Path-redacted artifact ids must use the format
`artifact:<kind>:<sha256-prefix-16>` and must not contain path separators,
provider hostnames, query strings, usernames, repository slugs, or private
filesystem segments.

## CLI Boundary

The command surface should prefer an explicit export command:

```bash
go run ./cmd/sdp-trace export cross-repo-posture \
  --profile cross-repo-evidence-posture-v1 \
  --selection <selection-manifest.json> \
  --out <cross-repo-posture-export.json>
```

`--profile`, `--selection`, and `--out` are required. `--validate-only` may
validate the selection manifest, artifact reachability, digest shape, freshness
boundary, and safe labels without writing an export artifact.

A future explain mode may render the JSON artifact, but it must not add
conclusions absent from the JSON. Explain output uses this stable traversal
order: header fields, input selection, movement summary, refusal rows, metric
rows by row id, movement rows by row id, handoff, and output safety. The first
profile does not support truncation, pagination, hidden filters, or omitted
sections; large artifacts should be consumed as JSON by downstream tools.
Explain must re-run the same output-safety check on the rendered bytes before
printing.

Stdout JSON mode is not part of the first profile because aggregate exports are
likely to be captured in logs and BI pipelines by default. The command must not
print export payloads, raw parser exceptions, stack traces, private paths, or
unsafe upstream values to stdout or stderr.

Partial exports that contain `not_assessed_input`, `stale_input`, or refused
repository/window rows may exit zero when a valid export artifact is written.
Usage errors, unreadable selection manifests, serialization failure, and cases
where no valid export artifact can be written exit non-zero.

Developer-facing error taxonomy:

| Error code | Exit | Stderr boundary |
| --- | --- | --- |
| `usage_error` | non-zero | flag name and safe reason only |
| `selection_unreadable` | non-zero | path-redacted selection id only |
| `out_unwritable` | non-zero | path-redacted output id only |
| `no_export_artifact` | non-zero | closed reason code only |
| `partial_export_with_refusals` | zero | count of refusal rows only |
| `validate_only_failed` | non-zero | closed reason code and row id only |

## Safety Requirements

Aggregate export is safety-sensitive because it joins evidence across
repositories.

Export JSON and explain output must not print:

- raw command arguments, command names, executable paths, script paths, or test
  identifiers unless public-catalog safe;
- stdout/stderr bodies;
- prompts;
- source snippets;
- tool-call input/output bodies;
- adapter configuration;
- gateway evidence refs;
- credentials, tokens, signed URLs, OIDC request tokens, adapter secrets,
  gateway tokens, PR tokens, or authenticated provider URLs;
- raw model request or response payloads;
- raw review bodies;
- raw-reference access notes that could expose storage location secrets;
- private filesystem paths;
- raw personal identifiers that are not declared safe role/team/service
  dimensions by selected upstream artifacts.
- raw digest-manifest paths or unredacted artifact ids;
- free-text refusal, exception, or non-comparable reason strings.

Negative leak tests must use synthetic marker values and must not echo candidate
secret values in failure output. Tests must include reserved synthetic prefixes
for each sensitive class, regex-class checks for token-like strings, URL/path
class checks, and hashed negative assertions so failure messages identify only
the class and digest of the leaked marker.

External verdict inputs are never exported as payloads. Block 21 may count only
closed, validated fields needed for `override_rows`; malformed, unsafe, or
payload-bearing external verdict inputs become `cannot_verify_input`.

The selected `dimension_exposure_policy` is a safety boundary. Team, service,
repo, harness, or change-type dimensions that are organizationally sensitive
for a customer or pilot must be omitted by policy or replaced by safe slugs
before aggregation. Block 21 must refuse dimensions outside the exposure policy
instead of relying on display-layer hiding.

## Acceptance Criteria

- A committed valid fixture exports at least two repositories and two time
  windows with metric rows grouped by repo, team, service, harness, change type,
  evidence family, and row evidence state.
- Every metric row includes numerator, denominator, unit, time window,
  dimensions, source artifact digest refs, `not_assessed_count`, and input trust
  state summary.
- Movement rows expose current value, previous value, delta, comparison basis,
  metric id/version, `dimension_key`, current/previous metric row refs,
  comparable boolean, and closed refusal reason when non-comparable.
- Export-level `movement_summary` exposes comparable and non-comparable counts
  by closed reason code.
- Stale, digest-mismatched, malformed, missing required, missing optional, and
  non-comparable inputs are represented as explicit rows or refusal facts.
- Block 21 validates safe labels itself and rejects provider URLs, private
  paths, email-like strings, token-like strings, raw personal identifiers, and
  labels outside the selected dimension exposure policy.
- No output field is named or described as degradation, improvement, health,
  readiness, pass/fail, rank, grade, red/yellow/green, blocked, acceptable, or
  unacceptable unless it is inside a clearly named external verdict input
  record.
- Schema and tests reject hidden weighting, missing denominator, missing source
  artifact digest refs for readable inputs, free-text metric ids, free-text
  dimensions, free-text refusal reasons, free-text non-comparable reasons, and
  native policy verdict fields.
- Explain output renders only JSON artifact fields in stable row order and
  adds no conclusions, summaries, hidden severity order, ANSI color state, or
  omitted-section state.
- Safety tests prove aggregate export and explain output do not persist the
  sensitive classes listed above.

## Implementation Plan

Slice A: Contract and fixtures.

- Add `schema/cross-repo-posture-export.schema.json`.
- Add `examples/block21-cross-repo-posture/` with valid mixed input,
  stale-input, digest-mismatch, missing-required, missing-optional,
  non-comparable-window, and unsafe-label fixtures.
- Add a fixture matrix that names expected metric rows, movement rows, refusal
  rows, dimensions, input trust states, numerator/denominator values, and
  source digest refs.
- Include posture signal fixtures for witness scope, overrides, late attach,
  and contract changes so those metrics are not invented from Block 20 rows.

Slice B: Aggregator.

- Add a small Go package that reads selection manifests and Block 20 query-pack
  result artifacts, validates posture signal manifests, verifies readable input
  digests, maps source rows to closed metric ids, and writes deterministic
  aggregate rows.
- Refuse or mark stale/untrusted inputs before aggregation.
- Preserve all non-pass evidence states in denominators.

Slice C: CLI and explain.

- Add `export cross-repo-posture` with required `--profile`, `--selection`, and
  `--out`.
- Add explain rendering over the JSON artifact only if the JSON contract and
  safety tests are already passing.
- Keep command output path-oriented and avoid stdout JSON mode for the first
  profile.

Slice D: Safety and review evidence.

- Add negative leak tests for every sensitive class.
- Add schema/fixture alignment tests and changed-example validation.
- Run code/correctness, tracing/evidence, and requirements-vs-implementation
  reviews after implementation.
- Repeat the three planes at PR level before ready/merge.

## Activation Gate

Do not implement Block 21 behavior until this spec delta has passed Socratic
review across product-boundary, tracing/evidence, and privacy/safety planes and
the reviewed direction is explicitly approved by the user.

Block 21 implementation depends on Block 20 result semantics. The current repo
contains Block 20 code and fixtures, but the Block 20 spec file still says it is
awaiting explicit implementation approval. That contradiction must be resolved
or explicitly accepted as historical drift before Block 21 implementation can
claim full dependency confidence.

## Review Ledger Shape

Review findings are recorded in
`specs/001-sdp-trace-time-series-evidence-substrate/blocks/21-cross-repository-degradation-export-review-ledger.md`
with severity, plane, finding, disposition, and evidence. Critical and major
findings must be fixed or explicitly blocked before approval handoff.

## No-Overclaim Notes

- Cross-repo export facts are movement inputs, not degradation decisions.
- Denominators and `not_assessed` counts are product requirements, not display
  details.
- Stale or untrusted inputs are not weak positives; they are explicit refused
  or capped inputs.
- `sdp-report` or BI can interpret the export, but that interpretation is
  outside this repository unless recorded as an external verdict input.
- Checked-in export JSON is not authority unless live-verified or externally
  signed.
