# Block 10 Implementation Plan: Agentic SDLC Evidence Substrate V0

Status: ready for implementation planning review; implementation not started
Parent: `10-agentic-sdlc-v0-design.md`

This plan prepares parallel implementation. It does not close any trust
claim and does not assert that the described commands already pass.

## Goal

Implement the V0 Go CLI and evidence contracts needed for a separate
demo repository to demonstrate `sdp-trace` on OpenCode + GSD + Bazel +
Kotlin while preserving the portable product boundary from the design
doc.

## Execution Rules

- The Block 10 target architecture is Go-only.
- Do not add Node.js, npm, JavaScript, TypeScript, or `.mjs` tooling.
- Node artifacts are not part of the active product path.
- Bash is allowed only as a thin launcher when Go cannot reasonably
  provide the behavior; any Bash kept must have an explicit reason in
  the implementation notes.
- Do not add TODO or FIXME markers.
- Use Clean Architecture boundaries: domain logic in internal packages,
  CLI as a thin adapter, filesystem/process/socket dependencies behind
  interfaces where tests need seams.
- Use TDD for behavior changes.
- Keep CRAP below 5 for changed Go code.
- Do not claim OpenCode, GSD, Bazel, Kotlin, or CI support as native
  `sdp-trace` product support.
- Keep OpenCode/GSD/Bazel/Kotlin demo artifacts outside this repository;
  this repository may contain only generic contracts, fixtures, and
  verifier behavior.
- Local trace validity must never be described as audit-grade trust.
- Missing observations must be explicit `missing_telemetry`,
  `cannot_verify`, or `not_assessed`.
- Behavior changes need tests before or with implementation.
- Every new trust-affecting output must be verifier-derived.
- Review findings must be recorded before any closure claim.

## Proposed File Responsibilities

Go implementation:

- `go.mod`: Go module root.
- `cmd/sdp-trace/main.go`: CLI entrypoint.
- `internal/trace/event.go`: canonical event model and hashing.
- `internal/trace/gen.go`: schema-bound Go type generation or checked
  hand-written type mapping.
- `internal/trace/store.go`: append-only run directory writer.
- `internal/trace/redaction.go`: pre-write redaction and retention
  descriptors.
- `internal/recorder/wrapper.go`: transparent process wrapper.
- `internal/adapter/socket.go`: Unix socket adapter ingress.
- `internal/contract/contract.go`: expected evidence contract loading.
- `internal/policy/authority.go`: adapter, signer, and witness authority
  policy.
- `internal/verifier/verifier.go`: chain, contract, witness, and missing
  evidence verification.
- `internal/query/query.go`: reviewer-facing query surface.
- `internal/export/audit_bundle.go`: audit bundle export.
- `internal/ci/witness.go`: CI witness record verification and signing
  interface.

Schemas and fixtures:

- `schema/agentic-sdlc-event.schema.json`
- `schema/expected-evidence-contract.schema.json`
- `schema/authority-policy.schema.json`
- `schema/redaction-profile.schema.json`
- `schema/signing-profile.schema.json`
- `schema/missing-evidence-table.schema.json`
- `schema/verifier-result.schema.json`
- `schema/audit-bundle.schema.json`
- `schema/integrity-audit.schema.json`
- `examples/agentic-sdlc/local-wrap-positive/`
- `examples/agentic-sdlc/tamper-negative/`
- `examples/agentic-sdlc/missing-adapter/`
- `examples/agentic-sdlc/ci-witness-positive/`

Demo and docs:

- `docs/agentic-sdlc-v0.md`
- Generic product docs only; demo-run instructions belong in the
  external demo repository.

## Slice A: Schema And Fixture Contract

Owner: schema worker.

Purpose: define the machine contract before the CLI hardens around it.

Tasks:

- Add V0 schemas for event, contract, MissingEvidenceTable, verifier
  result, and audit bundle.
- Define hash preimage using RFC 8785 canonical JSON with `event_hash`
  excluded.
- Define authority policy, signing profile, redaction profile, and
  integrity audit schemas.
- Own canonical Go event type generation or checked hand-written mapping.
  Other slices must consume this package instead of inventing structs.
- Add positive fixture for a local wrapped command.
- Add negative fixture with a mutated event payload.
- Add missing adapter fixture.
- Add CI witness fixture using a demo signing profile clearly marked as
  non-production unless a real CI/OIDC profile is used.
- Add JSON schema validation commands to the fixture validator only after
  fixtures exist.
- Add `scripts/validate-agentic-sdlc-fixtures.sh` as the schema worker's
  responsibility only if a launcher is still necessary. Validation logic
  must live in Go.
- Keep schema and fixture validation in Go commands.

Verification:

```bash
jq empty schema/*.json
go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc
git diff --check
```

Exit criteria:

- All V0 fixtures parse and validate against schemas.
- Negative fixtures remain negative once verifier behavior exists.
- Fixture layout follows the Block 10 run artifact layout exactly.

## Slice B: Go CLI Skeleton And Transparent Wrapper

Owner: CLI/recorder worker.

Purpose: create the minimum Go CLI that can run real commands without
breaking developer workflow.

Tasks:

- Add Go module and CLI command parser.
- Implement:
  - `sdp-trace dry-run`
  - `sdp-trace run`
  - `sdp-trace wrap`
- Preserve child stdin, stdout, stderr, TTY mode, signals, and exit code.
- Use a pseudoterminal when parent stdio is interactive; use direct
  passthrough/pipes for non-interactive execution.
- Export `$SDP_TRACE_SOCKET` after socket creation and before child
  launch.
- Write run artifacts under an explicit output directory or a safe
  ignored default.
- Use the Block 10 run artifact layout:
  - `run.json`
  - `events/*.json`
  - `artifacts/`
  - `verifier/`
  - `export/`
- Emit `recorder_attached`, `run_started`, `command_started`,
  `command_finished`, and `run_closed`.
- Add tests for exit code passthrough and signal/closure recording where
  practical.
- Add binary tests for:
  - non-interactive exit code passthrough;
  - TTY allocation detection using a small fixture program;
  - SIGTERM closure recording where the platform permits it.

Verification:

```bash
go test ./...
go run ./cmd/sdp-trace --help
go run ./cmd/sdp-trace wrap --name echo -- /bin/echo hello
```

Exit criteria:

- Wrapper behavior is transparent for simple commands.
- Local run produces a structurally parseable event chain.
- Local run directory can be consumed by Slice E without adapter or demo
  dependencies.

## Slice C: Redaction, Retention, And Safe Defaults

Owner: privacy/DX worker.

Purpose: prevent the demo from leaking secrets or raw code while making
the output understandable.

Tasks:

- Implement built-in `prewrite_digest_default_v1`.
- Implement `dry-run` output showing what would be retained.
- Store argv as digest plus allowlisted basename by default.
- Store stdout/stderr as digests unless contract opts into sanitized
  excerpts.
- Add redaction manifest digest and retention descriptor to events.
- Emit `redaction_applied` and `retention_lifecycle_observed` when
  relevant.

Verification:

```bash
go test ./internal/trace ./internal/recorder
go run ./cmd/sdp-trace dry-run --contract examples/agentic-sdlc/contracts/basic.json -- /bin/echo --token=secret
```

Exit criteria:

- Raw secret-like values do not reach persisted events in the default
  profile.
- Dry-run says it is simulation only and does not claim future proof.
- Redaction failure emits sanitized `redaction_failed` and downgrades
  completeness.

## Slice D: Adapter Socket Ingress

Owner: adapter worker.

Purpose: support optional OpenCode/GSD telemetry without requiring
harness rewrites.

Tasks:

- Create Unix socket and export `$SDP_TRACE_SOCKET` to child process.
- Create socket in a per-run private temp directory.
- Accept adapter registration and adapter lifecycle events.
- Validate allowed event types after registration.
- Implement bounded drain after child exit before `run_closed`.
- Reject or mark late adapter messages after `run_closed`.
- Remove stale socket files only inside the recorder-owned per-run temp
  directory.
- Emit missing rows for unsupported, absent, suppressed, disconnected, or
  unauthorized adapters.
- Add small fixture/client test that sends adapter events during a
  wrapped command.

Verification:

```bash
go test ./internal/adapter ./internal/recorder
go run ./cmd/sdp-trace wrap --name adapter-fixture -- examples/agentic-sdlc/adapter-client.sh
```

Exit criteria:

- Adapter events appear in the chain only with explicit identity state.
- Self-claimed adapter identity does not upgrade to gate-grade trust.

## Slice E: Verifier, MissingEvidenceTable, Explain

Owner: verifier worker.

Purpose: turn raw traces into deterministic trust-state output.

Tasks:

- Implement `sdp-trace verify <run-dir>`.
- Verify sequence, event hash, previous hash, run closure, contract lock,
  retention states, and witness tuple when present.
- Emit four-axis verifier output.
- Emit MissingEvidenceTable from the contract.
- Implement `sdp-trace explain <run-dir>` for:
  - missing run directory;
  - no events;
  - corrupt chain;
  - local-only trace;
  - missing contract;
  - contract locked late;
  - unauthorized signer;
  - storage overflow.
- Write `integrity_audit` when the chain is corrupted instead of
  appending to the corrupted chain.
- Implement `sdp-trace observe --run <run-dir> --state <state> --event
  <event> --reason <reason>` for local demo/preflight observations. This
  command cannot upgrade beyond local/not_assessed states without
  authority policy.

Verification:

```bash
go test ./internal/verifier ./internal/query
go run ./cmd/sdp-trace verify examples/agentic-sdlc/local-wrap-positive
go run ./cmd/sdp-trace explain examples/agentic-sdlc/tamper-negative
```

Exit criteria:

- Local positive fixture verifies as local-only, not gate-grade.
- Tamper fixture fails with a stable machine-readable reason.
- Missing adapter fixture produces explicit missing rows.
- Verifier reads only the shared run artifact layout and does not depend
  on recorder internals.

## Slice F: Query And Audit Bundle Export

Owner: forensics worker.

Purpose: make the trace useful without forcing reviewers to read JSON.

Tasks:

- Implement `sdp-trace query` for:
  - `run-summary`;
  - `timeline`;
  - `missing-evidence`;
  - `commands`;
  - `files`;
  - `tests`;
  - `redactions`;
  - `witness`;
  - `overrides`;
  - `retention`.
- Implement `sdp-trace export <run-dir> --format audit-bundle`.
- Audit bundle includes event chain, chain head, verifier result,
  MissingEvidenceTable, witness assertions, retention manifest, and
  integrity audit record if present.

Verification:

```bash
go test ./internal/query ./internal/export
go run ./cmd/sdp-trace query examples/agentic-sdlc/local-wrap-positive --query missing-evidence
go run ./cmd/sdp-trace export examples/agentic-sdlc/local-wrap-positive --format audit-bundle
```

Exit criteria:

- Forensics can reconstruct the run timeline and gaps from CLI output.
- Audit bundle is reproducible from fixture inputs.

## Slice G: CI Witness And Gate-Grade Boundary

Owner: CI/trust worker.

Purpose: demonstrate where local evidence becomes gate-usable.

Tasks:

- Implement `sdp-trace ci-verify`.
- Bind source digest, contract digest, run id, chain head, verifier
  version, verifier result, CI identity, timestamp, and witness
  independence.
- Support a demo signing profile and keep it clearly lower trust than
  real CI/OIDC unless actual OIDC evidence is present.
- Treat repo-committed CI YAML pinning as insufficient when the assessed
  change can modify that YAML.
- Distinguish `independent`, `same_job`, `same_container`,
  `same_process`, and `not_assessed` witness independence.
- Reject replay against another source digest or contract digest.
- Downgrade same-process or same-job witness independence where
  applicable.

Verification:

```bash
go test ./internal/ci ./internal/verifier
go run ./cmd/sdp-trace ci-verify --run examples/agentic-sdlc/ci-witness-positive --contract examples/agentic-sdlc/contracts/basic.json
```

Exit criteria:

- CI witness positive fixture can reach `ci_witnessed` only when the
  tuple matches.
- Local-only traces cannot be promoted by local signatures.

## Slice H: External Demo Repo Integration Contract

Owner: demo integration worker.

Purpose: define the minimum contract that an external demo repository
uses to exercise the generic product CLI without adding demo artifacts
to this repository.

Tasks:

- Define the CLI commands and artifact layout an external demo repo must
  use: `wrap`, `report`, `gate`, and `witness`.
- Define the generic expected-evidence contract fields needed by a
  harness/build-system demo without naming OpenCode, GSD, Bazel, or
  Kotlin as product dependencies.
- Require external demo repos to record unavailable harness/build-system
  dependencies as `not_assessed` or `missing_telemetry`, never as pass.
- Keep tamper, local-only, and CI-witness scenarios as product verifier
  behavior over generic fixtures.
- Store demo-specific source, build files, harness config, OpenCode/GSD
  runs, Bazel output, and Kotlin code only in the external demo repo.

Verification:

```bash
go test ./...
go run ./cmd/sdp-trace --help
go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc
```

Exit criteria:

- The product repo exposes a stable integration contract for the external
  demo repo.
- No OpenCode, GSD, Bazel, Kotlin, or demo-run artifacts are required in
  the product repository.

## Slice I: Review, Ledger, And Closure Preparation

Owner: review coordinator.

Purpose: keep implementation honest before closure.

Tasks:

- Run pi reviews on design and implementation artifacts with GLM,
  MiniMax, and Kimi using the review plan.
- Record valid findings in a Block 10 review ledger.
- Fix critical and major findings before implementation closure.
- Keep minor findings visible if intentionally deferred.
- Run fresh verification before any completion claim.

Verification:

```bash
jq empty schema/*.json
go test ./...
git diff --check
```

Exit criteria:

- No critical or major valid review finding remains unresolved.
- Completion language is limited to commands that were actually run.

## Parallelization Guidance

Safe parallel work sets:

- Slice A starts first and owns schemas, canonical Go event mapping, and
  the shared run artifact layout.
- Slice B can scaffold CLI in parallel, but event writing must import the
  Slice A canonical event package. Provisional event structs are not
  allowed beyond throwaway tests.
- Slice E can start from Slice A fixtures, but must not invent a
  different run directory layout.
- Slice C can run with Slice B after event persistence is stubbed.
- Slice D is independent after CLI env propagation exists.
- Slice E should consume Slice A and B outputs; it can start with
  fixtures before full CLI capture exists.
- Slice F can start from verifier output fixtures.
- Slice G should wait for verifier tuple definitions.
- Slice H can prepare scripts and demo contracts while core CLI lands,
  but must not claim a real demo until run evidence exists.
  Slice H must not influence first-milestone recorder/verifier scope.

Conflict boundaries:

- Only one worker should own `cmd/sdp-trace/main.go` command wiring at a
  time.
- Schema changes must be coordinated through Slice A.
- Demo docs must not broaden claims beyond verifier behavior.

## First Milestone

The first buildable milestone is:

```text
local wrap -> event chain -> verify -> missing evidence -> explain -> tamper fail
```

This is enough to show the product shape before CI witness and adapter
integration are complete.

First-milestone validation order:

```bash
jq empty schema/*.json
go test ./...
go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc
git diff --check
```

The first milestone must not depend on Node or npm.
