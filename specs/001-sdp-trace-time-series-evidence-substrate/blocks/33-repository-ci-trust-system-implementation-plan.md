# Block 33 Implementation Plan: Repository CI Trust System

Status: Draft. Implementation is blocked until Socratic review is complete and
the reviewed direction is explicitly approved.

## Slice 1: Fast Required Verify

Files:

- `.github/workflows/ci.yml`
- `docs/ci-check-policy.md`

Work:

- Keep one fast required `verify` job.
- Expand JSON syntax checks from one example folder to the relevant committed
  schema and example surface.
- Add Go formatting verification without changing files in CI.
- Add workflow YAML syntax verification for `.github/workflows/*.yml` with the
  pinned Go-based `actionlint` command from the spec.
- Keep the job name and docs clear that this is code/syntax hygiene, not release
  proof or external trust.

Verification:

- local execution of the exact commands added to CI;
- `go test ./...`;
- `gofmt -l $(git ls-files '*.go')`;
- `find schema examples -name '*.json' -print0 | xargs -0 jq empty`;
- `go run github.com/rhysd/actionlint/cmd/actionlint@v1.7.7 .github/workflows/*.yml`;
- `git diff --check`.

## Slice 2: Contract Validate Job Design

Files:

- `.github/workflows/ci.yml` or a separate `.github/workflows/contract-validate.yml`
- `docs/ci-check-policy.md`

Work:

- Keep `contract-validate` evidence-only or documented-only in Block 33 unless
  Socratic review approves at least two concrete validators from the inventory.
- Wire only validators that already exist as committed Go code or reviewed
  standard tooling at the start of implementation.
- Do not include unavailable validators in required job scope. Track them as
  future work instead of emitting green `not_assessed` branch-protection checks.
- Use static path targets or full-directory commands that run the same locally
  and in CI. Diff-scoped validation may be an optimization only after the full
  command is defined.
- Candidate first-pass evidence-only checks:
  - package-level fixture replay already covered by `go test ./...`;
  - broad schema/example JSON syntax from Slice 1;
  - future self-trace mirror/hash sync, claim-tag validation, PR-review profile
    ref validation, and source-bound manifest-subject drift remain out of scope
    until concrete validators are added.

Verification:

- focused local commands for each validator;
- no broad Bash framework unless it is a thin launcher around Go or standard
  validation tools;
- no `not_assessed` surface in any required job.

## Slice 3: Evidence-Only Workflow Boundaries

Files:

- `docs/ci-check-policy.md`
- `.github/workflows/pr-review.yml`

Work:

- Preserve `pr-review-evidence-only` as a separate evidence job.
- Document that missing model-review secrets remain `not_assessed`.
- Defer `trace-evidence-artifacts` to a later reviewed block. Block 33 may
  document the boundary but must not add the workflow.
- If a later block adds it, it must upload only safe artifacts and must not
  depend on raw local paths, raw prompts, raw model outputs, provider tokens, or
  unchecked PR-head execution.

Verification:

- workflow syntax check;
- artifact allowlist review;
- safety marker tests if new artifact rendering code is introduced.

## Slice 4: Source-Bound Release Boundary

Files:

- `docs/ci-check-policy.md`

Work:

- Separate source-bound release verification from ordinary PR CI in docs.
- Defer a `source-bound-release-verify` workflow to a later reviewed block.
- If a later block adds it, the first trigger should be `workflow_dispatch`;
  generic PR CI must not run source-bound release proof.
- Do not auto-regenerate or commit release proof from CI.
- Keep `source_bound_local_release` narrower than external production trust.

Verification:

- docs state `cannot_verify` when source-bound context is missing;
- no new source-bound workflow in Block 33 unless separately re-approved.

## Slice 5: Review, PR, And Promotion

Work:

- Run Socratic spec review across:
  - product boundary / UX of check names and branch protection;
  - trace/evidence correctness and state mapping;
  - DX/maintainability and CI cost/latency.
- Fix or explicitly block every valid critical and major finding.
- After implementation approval, split implementation so baseline `verify`
  changes can land independently from contract/evidence documentation.
- Run implementation review and PR-level review across code/correctness,
  trace/evidence, and requirements-vs-implementation.
- Keep merge approval separate from green CI and review evidence.

Verification:

- `go test ./...`
- broad JSON syntax check selected by reviewed spec
- workflow YAML syntax check
- `git diff --check`
- GitHub Actions `verify`
- GitHub Actions `contract-validate` only if it remains evidence-only and has
  concrete validators
- evidence-only workflow state recorded without overclaiming

## Initial Recommendation

Implement Block 33 in two PRs after spec approval:

1. **PR 1 / T247 + part of T250**: baseline required CI hardening. Broaden JSON,
   add gofmt check, add workflow YAML syntax, and update CI policy docs.
2. **PR 2 / T248-T250**: contract/evidence/release boundary documentation and,
   only if the inventory supports it, evidence-only `contract-validate`.
   `trace-evidence-artifacts` and `source-bound-release-verify` remain future
   workflow blocks.

T251 spans both PRs: each PR needs local verification, implementation review,
PR-level review, and GitHub workflow evidence.

This avoids mixing ordinary CI hygiene with proof workflow semantics and makes
review easier.

## Risks

- A larger required CI can slow every PR and push maintainers to bypass checks.
- Naming a workflow "trust" can mislead reviewers into treating it as release
  approval.
- Source-bound proof can become self-invalidating if docs or task ledgers record
  exact head SHA claims that change in the next commit.
- Adding shell-heavy validators would violate the Go-first product direction and
  make the repo harder to replay.
- Running evidence workflows on PR-head code with secrets would repeat the Block
  32 secret-boundary failure mode.
