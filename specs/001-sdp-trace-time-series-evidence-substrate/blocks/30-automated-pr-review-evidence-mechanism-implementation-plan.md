# Block 30 Implementation Plan: Automated PR Review Evidence Mechanism

Status: Implementation plan approved. Implementation has started after Block 30
Socratic review and explicit approval of the reviewed direction.

Spec:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/30-automated-pr-review-evidence-mechanism.md`

## Target Outcome

Add a Go-first PR review evidence mechanism that can assemble a frozen PR packet,
optionally run configured external reviewers, validate reviewer outputs and
coverage, and render a PR-safe summary without becoming a merge gate.

## Work Slices

### Slice 1: Packet And Profile Contracts

Owned surfaces:

- `schema/pr-review-packet.schema.json`
- `schema/pr-review-profile.schema.json`
- Go contract structs under the existing internal package pattern
- fixtures under `examples/pr-review/`

Behavior:

- Validate safe `repo_id`, provider-neutral `change_ref`, base/head SHA,
  unified-diff digest, per-ref context and verification digests, explicit
  `--ci-state`, redaction state, and unavailable fields.
- Validate review profile role declarations, required planes, runner ids,
  timeout shape, output schema refs, and raw-output retention choices.
- Add a `trust-sensitive-default` example profile and profile validation
  fixture.

Verification:

- valid packet/profile fixture;
- missing diff digest negative;
- unsafe repo id / path-derived repo id negative;
- invalid `change_ref` negative;
- missing required plane negative;
- CI absent fixture with `ci_state=not_assessed`.

### Slice 2: `pr-review packet`

Owned surfaces:

- CLI parser and help text for `sdp-trace pr-review packet`
- packet writer
- focused command tests

Behavior:

- Build a packet from explicit file inputs only.
- Compute and record SHA-256 digests.
- Refuse unsafe absolute output refs in machine-readable artifacts.
- Preserve unavailable fields as explicit state, not implicit absence.
- Fail when `--out` exists and is non-empty unless implementation adds
  reviewed `--force` behavior.

Verification:

- packet command writes deterministic JSON;
- changed diff changes packet digest;
- missing optional verification is represented as `not_assessed`;
- packet input guide fixture includes diff, metadata, context, and verification
  examples;
- help text names non-authority boundary.

### Slice 3: External Runner Records

Owned surfaces:

- `schema/pr-review-result.schema.json`
- runner result structs
- command execution wrapper for configured external commands
- fake runner fixtures

Behavior:

- Record `pi` and `opencode` availability, version where available, command
  digest, prompt digest, packet digest, exit state, started/ended timestamps,
  raw output digest, parser status, and model identity when observed.
- Do not require network calls in tests.
- Treat external runners as opaque commands whose stdout must match the role's
  required output schema.
- Map timeout, empty output, off-task structured output, unavailable command,
  parse failure, and mutation detection to explicit states using the spec's
  mapping table.
- Record requested model, observed model, fallback target, and fallback reason.

Verification:

- fake successful `pi` JSON output fixture;
- fake successful OpenCode JSON output fixture;
- fake empty output -> `empty_output` / not counted;
- fake malformed JSON -> `parse_failed` / not counted;
- fake unavailable runner -> `not_assessed`;
- fake OpenCode mutation evidence -> `cannot_verify`;
- no raw output body in committed fixture summary.
- fake runners are Go test helpers or Go-built test binaries, not shell
  scripts.

### Slice 4: `pr-review run`

Owned surfaces:

- CLI parser and help text for `sdp-trace pr-review run`
- role execution loop
- timeout handling
- ignored raw-output directory handling

Behavior:

- Require `--allow-external-runner` for every external runner invoked.
- Bind every role to the current packet digest.
- Support `run --preview` so the operator can inspect planned roles, commands,
  model ids, timeouts, prompt/template digests, and output paths before
  execution.
- Require proactive read-only OpenCode permission enforcement before model
  execution; if unavailable, mark the role `not_assessed` without executing.
- Support clean-tree-required OpenCode guard checks and reviewed dirty-baseline
  mode when a working tree is inspected.
- Preserve failed reviewer attempts as records but do not count them as usable
  plane evidence.

Verification:

- selected runner not allowed -> usage error;
- selected runner unavailable -> review result with `not_assessed`;
- packet digest mismatch -> `cannot_verify`;
- OpenCode write drift -> `cannot_verify`;
- role fallback preserves same plane and packet digest.
- durable command provenance never writes private absolute paths.

### Slice 5: Validation And Summary

Owned surfaces:

- `schema/pr-review-validation.schema.json`
- `sdp-trace pr-review synthesize`
- `sdp-trace pr-review validate`
- `sdp-trace pr-review summarize`
- `sdp-trace pr-review check`
- PR-safe markdown summary renderer

Behavior:

- Verify packet/run/ledger digest consistency.
- Verify required plane coverage.
- Verify critical/major findings have dispositions.
- Generate an initial ledger with `unresolved_review_blocker` dispositions from
  raw reviewer results.
- Preserve conflicts and residual gaps.
- Render `coverage_satisfied`, `coverage_partial`, `coverage_unresolved`,
  `not_assessed`, or `cannot_verify` review coverage state without merge
  approval language.
- Emit `authority_scope=review_record_only` plus explicit
  `not_authorized_by_sdp_trace` merge, release, and risk-acceptance fields.
- Render per-status next actions for every non-usable reviewer result.
- Render "This is not merge approval" as a structural line for
  `coverage_satisfied` and `no_findings` outputs.

Verification:

- all required planes usable -> `coverage_satisfied` when findings are disposed;
- missing required plane -> `coverage_partial`;
- unresolved major finding -> `coverage_unresolved`;
- stale packet digest -> `cannot_verify`;
- no reviewers -> `not_assessed`;
- changed diff requires new packet digest before re-review can close findings;
- summary never says safe-to-merge, approved, trusted, or ready.

### Slice 6: Safety And Regression

Owned surfaces:

- safety-sensitive tests
- docs/entrypoint updates after CLI help changes
- Block 30 review ledger

Behavior:

- Ensure summaries and validation output do not leak raw prompts, model
  responses, command bodies, provider tokens, authenticated URLs, private paths,
  or synthetic secret markers.
- Update entrypoint docs only after command help exists.
- Run implementation and PR-level multi-plane review.

Verification:

- `go test ./...`;
- `jq empty schema/*.json` plus changed fixtures;
- fixture validation for `examples/pr-review/`;
- `git diff --check`;
- focused output-safety tests with synthetic markers that are not echoed on
  failure.

## Review Plan

Before implementation approval:

- Socratic product-boundary review;
- UX/DX review for operator workflow and PR comment shape;
- trace/evidence/provenance review;
- security/privacy/output-safety review;
- implementation feasibility review.

After implementation:

- code/correctness review;
- trace/evidence review;
- requirements-vs-implementation review;
- focused re-review for any schema, safety, or state-mapping fixes;
- PR-level repeat of the same planes before ready.

## Approval Checkpoint

Implementation started only after:

1. Socratic review has usable output;
2. valid critical and major findings are fixed or recorded as blockers;
3. the technical executive explicitly approves the reviewed direction.
