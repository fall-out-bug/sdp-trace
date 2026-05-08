# Block 26 Implementation Plan: CI Artifact Observation Contract

Status: draft, blocked on Socratic spec review approval.

## Work Slices

| Slice | Owner surface | Output | Verification |
| --- | --- | --- | --- |
| WS0 spec review | Block 26 spec, implementation plan, review ledger | reviewed direction with dispositions | no unresolved critical/major Socratic findings |
| WS1 schema and examples | `schema/ci-artifact-observation.schema.json`, `examples/block26-ci-artifact-observation/` | valid and negative fixture matrix | `jq empty` plus focused Go/schema alignment review |
| WS2 evaluator | `internal/ciartifact/` or equivalent Go package | deterministic family/access/binding/index/safety states | focused Go tests for all fixture cases |
| WS3 CLI and explain | `cmd/sdp-trace/` | observation command or explicit assess profile | CLI tests for JSON and human-readable output |
| WS4 docs and demo handoff | `docs/`, Block 25 or future demo specs | demo truth contract and downstream gate guidance | drift review against Blocks 22, 24, 25 |
| WS5 review and PR | review ledger, PR body | code/tracing/requirements/security reviews | fresh local verification and PR-level review |

## Data Model Draft

The result artifact should include:

- `schema_version`
- `selected_profile`
- `artifact_observation_state`
- `authority_scope`
- `selected_source`
- `selected_run`
- `producer_scope`
- `artifact_access_state`
- `required_families`
- `artifact_families`
- `bindings`
- `artifact_index`
- `output_safety`
- `reasons`
- `next_actions`
- `safety_ruleset`

Closed states should reuse the existing vocabulary where possible:
`pass`, `fail`, `cannot_verify`, `not_assessed`, `present`, `absent`,
`partial`, `expired`, `inaccessible`, `malformed`, `unsafe`, `matched`,
`mismatch`.

Reason payloads must be built from closed reason codes, safe family ids, safe
state values, and safe rule ids. The evaluator must not copy raw parser input,
artifact bytes, raw paths, raw URLs, log lines, command bodies, prompts, or
model responses into `reasons` or `next_actions`.

The default output-safety ruleset must be documented and digest-bound. It must
cover token-like strings, JWTs, SSH/private key markers, cloud-provider
credentials, provider-token markers, private artifact URLs, private filesystem
paths, prompt/model-response markers, raw job-log markers, and high-entropy
secret-like values where practical. Rule failures must identify only rule ids
or safe classes.

## Fixture Matrix

Fixtures must be split into replayable inputs and expected outputs:

- `examples/block26-ci-artifact-observation/input/<scenario>/`
- `examples/block26-ci-artifact-observation/expected/<scenario>.observation-result.json`

Required scenario names:

- `ci-uploaded-bundle-complete-coverage`
- `checked-in-only-claim`
- `ci-bundle-absent`
- `ci-bundle-partial`
- `artifact-index-self-reference`
- `artifact-digest-mismatch`
- `source-run-binding-missing`
- `source-run-binding-mismatch`
- `agent-reported-happy-path`
- `unsafe-artifact-output`
- `artifact-expired`
- `external-artifact-ref-unverifiable`

Each fixture must cite the condition id and reason code that caused the top-level
state. A fixture that is intentionally malformed for parser tests must live in a
separate `malformed/` directory and must not be treated as schema-valid example
evidence.

Malformed fixture tests must not echo malformed file contents, raw input, or
private local paths in CI output. They should report only scenario id, safe
reason code, and expected state.

## Review Checkpoints

Do not implement WS1-WS5 until WS0 completes and the reviewed direction is
approved.

After WS2, run a tracing/evidence focused review before wiring CLI output. This
prevents a pretty CLI from hiding a weak state model.

After WS3, run a security/privacy focused review over serialized JSON and
explain output before adding docs that rely on the output.

Before PR ready/merge:

- `go test ./...`
- `jq empty schema/*.json`
- fixture syntax validation, excluding intentionally malformed fixtures;
- `git diff --check HEAD`;
- drift scan against Blocks 22, 24, and 25;
- PR-level review across code/correctness, tracing/evidence,
  requirements-vs-implementation, and security/privacy.

## Demo Handoff

The next demo repo work should not claim five-feature happy path until every
feature branch has:

- PR ref and head SHA;
- CI run id and attempt;
- uploaded artifact bundle or explicit unavailable state;
- Block 26 observation result;
- trace/provenance/evidence families or explicit per-family gaps;
- review disposition for any intentionally broken feature.

The demo can show broken features. The product requirement is that the broken
state is observable and cannot be mistaken for CI-backed proof.
