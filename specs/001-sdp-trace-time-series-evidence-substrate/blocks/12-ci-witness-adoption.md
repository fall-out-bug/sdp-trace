# Block 12: CI Witness And Adoption Handoff

Status: implementation-ready after M2.7/GLM spec review.

Documentation review correction:

- Block 12 CI witness is not a signed external trust boundary.
- `gate --witness` reads a JSON witness artifact and must be used only with a
  witness generated inside the protected CI workflow for the same job/artifact
  set.
- A witness JSON file supplied from a developer workspace or committed by an
  agent is not authority.
- External signature, append-only timeline, replay protection, and fail-closed
  managed harness enforcement are explicit follow-up blocks.

## Goal

Turn local contract evidence into a customer-visible CI-witnessed posture without
claiming external trust.

Block 12 adds a generic CI witness record that binds:

- repository identity;
- source commit;
- CI workflow/job/run identity;
- generated report/gate artifact digests;
- discovered run directories;
- local gate result;
- missing external witness state.

This lets a technical executive ask, per repository and commit:

```text
Was this agent-assisted change only locally observed, CI-witnessed, or still
missing the required witness?
```

## Non-Goals

- No public transparency log.
- No Sigstore/Rekor integration.
- No GitHub dependency in recorder/verifier core logic.
- No support claim for any harness, language, or build system.
- No audit-grade release claim from CI alone.
- No Node.js, npm, JavaScript, TypeScript, `.mjs`, or shell product tooling.

## User-Facing Commands

### `sdp-trace witness`

Usage:

```text
sdp-trace witness --kind github-actions --out <file> [--report-dir <dir>] <runs-root-or-run-dir>
```

Behavior:

1. Discover run directories exactly like `report` and `gate`.
2. Hash run manifests and report/gate artifacts when present.
3. Read CI identity from GitHub Actions environment variables.
4. Request and inspect a GitHub Actions OIDC token.
5. Write a witness record to `--out`.
6. Exit `0` only when the GitHub Actions identity is complete.
7. Exit `3` with `status: cannot_verify` when the command is not running inside
   GitHub Actions or required CI identity fields are missing.

Required GitHub Actions environment fields:

- `GITHUB_ACTIONS=true`
- `GITHUB_SHA`
- `GITHUB_RUN_ID`
- `GITHUB_RUN_ATTEMPT`
- `GITHUB_WORKFLOW`
- `GITHUB_JOB`
- `GITHUB_ACTOR`
- `GITHUB_REPOSITORY`
- `GITHUB_REF`
- `GITHUB_SERVER_URL`
- `ACTIONS_ID_TOKEN_REQUEST_URL`
- `ACTIONS_ID_TOKEN_REQUEST_TOKEN`

Partial CI identity behavior:

- if any required field is missing, the witness record is written with
  `status: cannot_verify`;
- the `reason` must name `missing_ci_identity`;
- missing field names must be recorded;
- the command exits `3`.

Replayable unit tests inject an environment map and OIDC token fetcher through
Go code. The CLI reads only the real process environment and uses GitHub's OIDC
request endpoint. Block 12 deliberately does not add `--env-file` or any flag
that can locally manufacture a passing CI witness.

Required witness fields:

- `kind`
- `status`
- `trust_scope`
- `reason`
- `generated_at`
- `source`
- `ci`
- `run_artifacts`
- `report_artifacts`

Allowed values:

- `kind`: `github-actions`
- `status`: `pass`, `cannot_verify`
- `trust_scope`: `ci_witnessed`, `local_observed`

`local_observed` is used only for `status: cannot_verify` records written
outside a complete CI identity. It is not a local witness upgrade.

`status: pass` requires a complete GitHub Actions identity and OIDC token claims
matching repository, ref, and commit SHA. It means only:

```text
GitHub Actions identity and OIDC claims were present and the witness record
binds the observed local artifacts to that CI run and source commit.
```

It does not mean the agent was honest, the tests are sufficient, the release is
audit-grade, or an external witness exists.

GitHub Actions OIDC claim binding:

- `iss` must be `https://token.actions.githubusercontent.com`;
- `repository` must match `GITHUB_REPOSITORY`;
- `ref` must match `GITHUB_REF`;
- `sha` must match `GITHUB_SHA`;
- `aud` must include `sdp-trace`.

### `sdp-trace gate --witness`

Usage:

```text
sdp-trace gate --out <file> --contract <contract.json> --witness <witness.json> <runs-root-or-run-dir>
```

Behavior:

1. Evaluate the existing local contract gate.
2. Load the witness record when `--witness` is provided.
3. Add `ci_witness_gate`.
4. If the witness has `status: pass` and `trust_scope: ci_witnessed`, set
   `ci_witness_gate: pass`.
5. Otherwise set `ci_witness_gate: cannot_verify`.
6. Keep `audit_grade_gate: cannot_verify` until an external witness profile
   exists.

`missing_audit_evidence` behavior:

- Without a passing CI witness:
  - `ci_oidc_witness`
  - `external_witness_checkpoint`
- With a passing CI witness:
  - `external_witness_checkpoint`

## Adoption Handoff

The customer implementation path is:

1. Team lead places an expected-evidence contract in the repository.
2. Platform owner wraps the existing harness command:

   ```text
   sdp-trace wrap --name <existing-harness> --contract <contract> --output-dir .sdp-trace-runs/<run-id> -- <existing command...>
   ```

3. CI runs the normal test/build workflow.
4. CI generates report and local gate:

   ```text
   sdp-trace report --out .sdp-trace-report --contract <contract> .sdp-trace-runs
   sdp-trace gate --out .sdp-trace-report/gate-result.json --contract <contract> .sdp-trace-runs
   ```

   `report --out` creates or updates the report directory.

5. CI writes witness:

   ```text
   sdp-trace witness --kind github-actions --out .sdp-trace-report/ci-witness.json --report-dir .sdp-trace-report .sdp-trace-runs
   ```

6. CI reruns gate with witness:

   ```text
   sdp-trace gate --out .sdp-trace-report/gate-result.json --contract <contract> --witness .sdp-trace-report/ci-witness.json .sdp-trace-runs
   ```

7. The technical executive-facing artifact is `.sdp-trace-report/`.

## Capture Boundary

Block 12 uses the existing process wrapper boundary:

- the wrapper observes wrapped process lifecycle and command-level metadata;
- it propagates the wrapped process exit code and still writes run artifacts
  when possible;
- it does not automatically observe internal tool calls inside an arbitrary
  harness;
- it does not prove that an agent was never run outside the wrapper;
- missing expected run artifacts are visible only when the expected-evidence
  contract or CI workflow requires them.

Harnesses with plugin APIs should add adapter events in a later block. Harnesses
without plugin APIs remain `local_observed`/partial unless the process or tool
layer can be wrapped.

## Security And Forensics Boundaries

Block 12 artifacts are JSON files. They are useful operational evidence, but
they are not a signed forensic timeline.

Current known limits:

- `.sdp-trace-runs/` is mutable unless the customer protects it with CI artifact
  storage or another immutable store;
- deleted local runs are not detectable unless CI expected them and records
  missing evidence;
- replay of old local artifacts is not fully prevented without an external
  timestamp/log profile;
- `ci-witness.json` is not a DSSE envelope and is not externally signed;
- `gate --witness` must not be used to trust a witness file from an untrusted
  workspace;
- redaction audit trails and retention policy are deployment responsibilities in
  Block 12.
- native `policy_override_requested` trace events are not implemented in Block
  12; emergency overrides must remain external policy records that reference the
  generated report artifacts.
- Block 12 has no raw-capture mode. Any later raw prompt/source/model-response
  capture profile must redact before persistent write.

## Operational States

- No wrapper/adapter output where evidence is expected:
  `missing_telemetry`.
- Local wrapper output with no CI witness:
  `local_observed`.
- CI witness without OIDC:
  `cannot_verify`.
- CI witness with matching OIDC claims:
  `ci_witnessed`.
- External audit-grade proof:
  not implemented in Block 12.

## Schema Delta

Block 12 adds:

- `schema/ci-witness.schema.json`
- gate-result v2 fields in the Go output:
  - `ci_witness_gate`
  - `witness`
  - conditional `missing_audit_evidence`

`ci-witness.schema.json` must define:

- required top-level fields;
- `kind`, `status`, and `trust_scope` enums;
- source commit/repository/ref object;
- CI provider identity object;
- GitHub OIDC claim projection for passing witnesses;
- artifact digest arrays;
- `additionalProperties: false` for the witness envelope.

The existing compatibility `gate-verdict.schema.json` is not the native Block
12 gate-result schema. Block 12 implementation tests validate behavior through
Go structs and `jq empty`; a dedicated native gate-result schema remains a
follow-up unless the current block introduces one.

## Acceptance Criteria

1. `go test ./...` passes.
2. `sdp-trace witness --kind github-actions --out <file> <runs-root>` writes
   `status: cannot_verify` and exits `3` outside GitHub Actions.
3. With a complete GitHub Actions env fixture, witness generation returns
   `status: pass`, `trust_scope: ci_witnessed`, and includes source commit,
   repository, workflow, job, run id, and artifact digests.
4. `gate --witness <passing-witness>` sets `ci_witness_gate: pass`.
5. `gate --witness <cannot-verify-witness>` sets
   `ci_witness_gate: cannot_verify`.
6. `audit_grade_gate` remains `cannot_verify` in both cases.
7. Demo report/gate/witness JSON validates with `jq empty`.
8. No product code names OpenCode, GSD, Bazel, Kotlin, or any demo-specific
   evidence kind.
