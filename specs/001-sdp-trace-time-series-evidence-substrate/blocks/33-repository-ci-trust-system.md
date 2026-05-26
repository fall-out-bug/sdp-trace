# Block 33: Repository CI Trust System

Status: Draft spec. Implementation is blocked until Socratic review is complete
and the reviewed direction is explicitly approved.

Parent artifacts:

- `.github/workflows/ci.yml`
- `.github/workflows/pr-review.yml`
- `docs/ci-check-policy.md`
- `docs/agent-entrypoint.md`
- `docs/repository-rollout-playbook.en.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/32-ci-pr-review-integration.md`
- `examples/contract-foundation/contract-manifest.example.json`

## Goal

Make repository CI for `sdp-trace` prove the repository's code, schemas,
fixtures, docs, and self-trace artifacts are internally coherent without
claiming merge approval, release readiness, production trust, or customer-repo
policy coverage.

The product answer should be:

```text
On every PR, required CI can prove fast code and contract consistency. Separate
evidence-only workflows can publish safe trace/review/release facts. Each check
name and artifact states exactly what was verified and what remains
not_assessed or cannot_verify.
```

## Product Question

"Can `sdp-trace` dogfood its own evidence discipline in CI strongly enough that
a maintainer can tell which repo claims are actually verified, without reading
the whole repository by hand?"

The answer must include:

- a required fast verification layer for code, broad JSON/YAML syntax, format,
  and whitespace;
- a contract validation layer for schema/example alignment, self-trace mirror
  state, claim tags, fixture matrices, and source-bound manifest-subject drift,
  with unavailable validators explicitly kept outside required CI;
- an evidence-artifact boundary for safe trace/review/report artifacts without
  making those artifacts part of branch-protection approval;
- a source-bound release boundary that keeps release verification out of
  ordinary PR CI until a later reviewed block selects it;
- explicit state mapping for `pass`, `fail`, `not_assessed`, and
  `cannot_verify`;
- branch-protection guidance that does not mistake evidence-only jobs for human
  approval;
- a roadmap boundary for customer-repository CI templates, kept outside this
  first block.

## Current State

The current required workflow is too narrow for a trust substrate:

- `.github/workflows/ci.yml` runs `go test ./...`, `jq empty
  schema/*.json examples/block19-adapter-capture/*.json`, and `git diff
  --check`.
- JSON syntax only covers one example folder, not the repo's schema and fixture
  surface.
- YAML workflow syntax is not checked locally or in CI.
- Go formatting is enforced indirectly only if tests or compiler errors expose
  it.
- There is no dedicated CI contract layer for self-trace mirror/hash sync,
  claim-tag validation, schema/example drift, fixture matrix drift, or
  source-bound manifest-subject drift.
- Block 32 adds `pr-review-evidence-only`, but that is model-review evidence
  for a frozen PR packet, not the repository's general quality gate.
- Release proof and external trust remain separate from ordinary PR checks.

Therefore a green `verify` check currently means "basic code tests and a narrow
JSON syntax check passed", not "the repository's trust artifacts are coherent".

## Current Validator Inventory

Block 33 may wire only validators that are concrete and locally runnable at the
start of implementation. Unavailable surfaces must stay out of required jobs;
they may be documented as future `not_assessed` evidence-only work with tracked
follow-up.

| Surface | Current validator | Block 33 status |
| --- | --- | --- |
| Go code behavior | `go test ./...` | ready for required `verify` |
| Go formatting | `gofmt -l $(git ls-files '*.go')` with non-empty output failing | ready for required `verify` |
| Whitespace | `git diff --check`, before any generated-output step | ready for required `verify` |
| JSON syntax | `find schema examples -name '*.json' -print0 \| xargs -0 jq empty` | ready for required `verify`; syntax only |
| GitHub Actions workflow syntax | `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml` | ready for required `verify` if review accepts this pinned Go-based external linter |
| Committed fixture replay | package tests under `go test ./...`; `sdp-trace validate-fixtures` exists but is not a whole-repo fixture-matrix validator today | partial; no separate required contract job until scoped command is reviewed |
| Self-trace mirror/hash sync | no single current Go command identified in this intake | not_assessed; future validator |
| Claim-tag validation | docs define syntax; no current Go validator identified in this intake | not_assessed; future validator |
| PR-review profile/prompt ref validation without model execution | profile loading exists inside `internal/prreview`; no CLI validation contract exists | not_assessed; future validator unless added after separate approval |
| Source-bound manifest-subject drift | documented source-bound cycle uses JQ/git with a selected source commit | evidence-only/manual; not required in Block 33 |

## Non-Goals

Block 33 must not:

- generate customer-repository CI templates;
- define downstream merge, release, risk, or approval policy;
- make `pr-review-evidence-only` a required human-review substitute;
- use source-bound local release proof as external production trust;
- add Node.js, npm, JavaScript, TypeScript, or `.mjs` product tooling;
- auto-commit generated proof artifacts from CI;
- hide missing checks, missing secrets, or skipped profiles behind a pass.

## CI Layers

### Layer 1: Required Fast Verify

Purpose: keep every PR mechanically buildable and catch broad syntax drift
quickly.

Required checks:

- `go test ./...`
- `gofmt -l $(git ls-files '*.go')` with non-empty output failing
- `git diff --check`
- `find schema examples -name '*.json' -print0 | xargs -0 jq empty`
- `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml`

Trust boundary:

- Proves only fast local code/syntax hygiene.
- Does not prove source-bound release, external trust, model-review coverage,
  or customer-repo readiness.

### Layer 2: Contract Validate

Purpose: prove repo contracts are internally coherent, not just syntactically
valid.

Validation surfaces:

- schema refs and committed examples are loadable and aligned with Go structs
  where a Go contract exists;
- committed fixture matrices are replayed by the Go verifier instead of only
  parsed by `jq`;
- self-trace mirrors are synchronized with source arrays;
- committed hash references match current files;
- claim tags use documented syntax and allowed evidence forms;
- contract manifest subjects have no unreported drift relative to the selected
  source-bound cycle;
- PR-review profiles and prompt-template refs are valid without executing model
  runners.

Trust boundary:

- Proves internal repository contract consistency for selected profiles.
- Required jobs must not include surfaces that would report `not_assessed`.
  Unimplemented validators stay outside required scope and are tracked as future
  work.
- Evidence-only contract reports may include `not_assessed` rows only when they
  name the follow-up owner.
- Does not update release proof artifacts by itself.
- First Block 33 implementation should harden Layer 1 and may document or add
  evidence-only `contract-validate`; it must not make `contract-validate`
  required until all included surfaces are concrete pass/fail/cannot-verify
  checks.

### Layer 3: Evidence Artifact Workflows

Purpose: produce safe retained facts that maintainers can inspect without
reading raw logs.

Evidence workflows:

- `pr-review-evidence-only` from Block 32 for automated model-review evidence;
- future `trace-evidence-artifacts` for report, gate, missing-telemetry,
  selected assessments, and query-pack outputs over a controlled self-trace or
  CI run package. This is explicitly deferred out of Block 33 implementation.

Trust boundary:

- These workflows publish evidence artifacts and Step Summaries.
- They are not merge approval, release readiness, production trust, or risk
  acceptance.
- Missing provider secrets, missing run artifacts, or skipped profiles remain
  `not_assessed`.
- Any `cannot_verify` state in a trust-sensitive evidence workflow must make the
  GitHub check fail unless that workflow is explicitly documented as a
  non-required observation mode and the Step Summary visibly says
  `cannot_verify`.

### Layer 4: Source-Bound Release Boundary

Purpose: keep source-bound local release claims out of ordinary PR CI until the
source-bound cycle is explicitly selected.

Block 33 boundary:

- no new required source-bound release workflow;
- no automatic proof regeneration from CI;
- no production-trust or release-readiness claim;
- document the future workflow as `workflow_dispatch` first, not generic PR CI;
- if a manifest-subject drift check is added later, it must name the selected
  source commit and emit `cannot_verify` when the source-bound context is
  missing.

Trust boundary:

- `source_bound_local_release: pass` is narrower than external production
  trust.
- Missing signature, stale manifest subjects, dirty source state, or missing
  verifier profile keeps release proof `fail` or `cannot_verify`.
- CI may verify existing proof artifacts, but generation/regeneration remains a
  reviewed source-bound cycle, not an automatic CI side effect.

## State And Job Mapping

| State | Meaning in repo CI | Required job behavior |
| --- | --- | --- |
| `pass` | Selected checks completed inside their stated trust scope. | Pass. |
| `fail` | Evidence contradicted the selected contract or tests failed. | Fail. |
| `not_assessed` | Surface was intentionally outside this workflow scope or unavailable by rollout policy. | Never part of required job scope. Allowed only in evidence-only outputs with a tracked follow-up owner. |
| `cannot_verify` | Workflow attempted a required check but lacked evidence, safe execution, or consistency. | Fail for required CI and trust-sensitive evidence workflows. Observation-mode evidence jobs may pass only if the Step Summary visibly reports `cannot_verify`. |

No job may translate absent GitHub checks, missing model-review secrets, missing
source-bound proof, or skipped profiles into green proof language.

## Branch Protection Guidance

Initial required checks:

- `verify`

Evidence-only after Block 33, not initially required:

- `contract-validate`

Evidence-only, not required for generic merge:

- `pr-review-evidence-only`
- future `trace-evidence-artifacts`
- future `source-bound-release-verify`

Promotion criteria for `contract-validate`:

- every included validator is concrete and locally runnable;
- no included validator can produce `not_assessed`;
- five consecutive PR or `main` runs complete without false-positive failure;
- no maintainer override was needed during those runs;
- a separate review approves the branch-protection change.

Backout criteria:

- if a promoted required `contract-validate` check produces more than one
  confirmed false-positive failure in seven days, demote it back to
  evidence-only until a reviewed fix lands.

Rationale: required branch checks should protect buildability and contract
consistency. Evidence-only workflows should inform reviewers and release
owners without replacing human approval or downstream policy.

## Acceptance Criteria

- CI docs name required checks, evidence-only checks, and release verification
  checks separately.
- `.github/workflows/ci.yml` no longer validates only one historical example
  folder for JSON syntax.
- Go formatting and workflow YAML syntax are checked in CI.
- The reviewed validator inventory names which current validators are ready now
  and which remain deferred outside required CI.
- Source-bound release verification is separated from ordinary PR checks and is
  deferred unless a later reviewed block selects it.
- `pr-review-evidence-only` remains evidence-only and is not required as a
  substitute for human approval.
- CI output and docs avoid stale head-SHA proof claims.
- Required jobs never include unavailable `not_assessed` surfaces; evidence-only
  jobs may record `not_assessed`, and required attempted checks that cannot
  verify fail.
- Customer-repository CI templates are explicitly deferred to a later block.

## Reviewed Decisions To Confirm

- `contract-validate` starts evidence-only or documented-only in Block 33. It
  becomes required only after the promotion criteria above are met.
- Block 33 required implementation scope is Layer 1 hardening unless Socratic
  review approves at least two concrete contract validators from the inventory.
- `source-bound-release-verify` is deferred to a later block or manual-only
  design; it is not part of ordinary PR CI in Block 33.
- `trace-evidence-artifacts` is deferred to a later block after baseline repo CI
  and contract-validation semantics are stable.
- Layer 1 latency budget: target under 5 minutes on `ubuntu-latest`.
- Evidence-only contract validation latency budget: target under 10 minutes.
