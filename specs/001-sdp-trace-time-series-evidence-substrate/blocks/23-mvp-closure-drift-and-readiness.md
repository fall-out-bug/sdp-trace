# Block 23: MVP Closure Drift And Readiness

Status: draft SpecKit delta for Socratic review. Implementation is blocked until
this spec is reviewed, valid critical/major findings are resolved, and the CTO
explicitly approves the reviewed direction.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/spec.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md`
- retired research artifact
- `docs/agent-entrypoint.md`
- `docs/reviewer-entrypoint.md`

## Goal

Prepare the current MVP backlog for honest closure without hiding proof drift,
documentation drift, code-quality risk, or unanswered customer questions.

The block outcome is not "make everything green by prose." The outcome is a
reviewed, machine-checkable readiness package that separates:

- fixed blockers;
- explicit `not_assessed` or `cannot_verify` states;
- deferred follow-ups that are not part of MVP closure;
- external production trust that remains outside local source-bound proof.

## Problem

The current repository has useful MVP functionality, but closure is not yet
trustworthy enough.

Current review found these closure gaps:

1. The active source-bound release proof fails on current `HEAD` with manifest
   artifact mismatches for `README.md`, `docs/flight-recorder.md`,
   `schema/README.md`, `spec.md`, and `tasks.md`.
2. Local `main` is one commit ahead of `origin/main`; GitHub/remote closure and
   final CI for that commit are `not_assessed`.
3. Block 06 still has open backlog items about a shell-shaped Bazel runner and
   dangling proof-state evidence refs, while historical ledgers mention removed
   Node/npm/script validation paths.
4. Go tests, `go vet`, `staticcheck`, and `golangci-lint` pass, but the repo has
   high-complexity hotspots and low or absent coverage in trust-adjacent
   packages. CRAP < 5 is not currently a measured repository gate.
5. Customer-facing English and Russian docs are stale against the shipped
   Block 22 command/profile surface.
6. Customer-facing docs answer some adoption questions, but not all buyer
   pressure-point questions with explicit artifact-to-answer mappings.
7. Some spec artifacts are stale in the opposite direction: Block 22 spec status
   still says implementation is blocked even though implementation and PR review
   have occurred.

## Non-Goals

- Do not claim `external_production_trust`, `trusted_contract_release`, or
  `production_release_verified`.
- Do not rewrite the product boundary into `sdp-gate`, reporting, dashboarding,
  or workflow orchestration.
- Do not remove historical review ledgers merely because their commands are
  retired. Historical evidence may stay, but it must not be accepted as current
  closure evidence.
- Do not force every legacy function below a complexity threshold in one block.
  MVP closure needs a measured gate, explicit exceptions, and focused fixes for
  trust-adjacent hotspots.
- Do not implement before Socratic review and explicit approval.

## Definitions

- `pass`: command exits `0` and, where JSON output exists, the selected profile
  result field is `pass` or `observed` exactly as defined by that command's
  contract.
- `fail`: command exits `1` or emits a selected profile result of `fail`.
- `cannot_verify`: the selected check was attempted but required evidence,
  environment, or access was unavailable or contradictory enough that the
  selected profile cannot conclude.
- `not_assessed`: the selected check was intentionally outside the current
  assessment scope. Every `not_assessed` entry in the closure package must name
  scope, reason, and follow-up or accepted non-goal.
- `source-bound anchor`: the Git commit named by
  `contract-manifest.example.json.source_commit` plus the manifest-listed
  subject paths and SHA-256 digests. The proof does not claim the containing
  proof-artifact commit is itself the source anchor.
- `trust-adjacent production package`: any production Go package that loads,
  emits, hashes, validates, or explains evidence, provenance, trace events,
  release proof, witness state, assessment state, or gate state. Initial scope:
  `internal/trace`, `internal/verifier`, `internal/releaseproof`,
  `internal/checkpoint`, `internal/recorder`, `internal/adaptercapture`,
  `internal/forensic`, `internal/managed`, `internal/witness`,
  `internal/posture`, `internal/query`, `internal/contract`, `internal/export`,
  and `internal/policy`.
- `changed production function`: a non-test Go function whose implementation,
  control flow, or production dependencies are modified in the Block 23
  implementation diff.
- `changed production file`: a non-test Go file modified in the Block 23
  implementation diff.
- `current closure evidence`: evidence produced by commands present in
  `go run ./cmd/sdp-trace --help` or by standard Go/JQ/git tools named in this
  block. Historical Node/npm/script evidence may remain only when labeled as
  historical and not used as current closure evidence.

## Workstreams

### WS1 - Source-Bound Proof Repair

Repair current source-bound local release proof against an immutable clean
source commit.

Acceptance criteria:

- The five current manifest mismatches are classified before repair as content
  drift, path drift, schema drift, or tool/hash-algorithm drift.
- If `source_commit` is absent or empty in the manifest, release proof must
  return `cannot_verify` with a concrete reason, not infer a local commit.
- `sdp-trace release-proof --manifest examples/contract-foundation/contract-manifest.example.json --out <repo-artifact>` exits `0`, and the JSON output records
  `release_verification_state: "pass"`,
  `trust_scope: "source_bound_local_release"`,
  `source_commit_status: "matched"`, and
  `source_commit_artifact_status: "matched"`.
- Manifest subject drift check prints no changed manifest subjects after proof
  regeneration.
- External production trust remains `not_assessed`; `trusted_contract_release`
  remains `false`.
- Proof regeneration and source-subject changes follow the source-bound cycle:
  source-subject changes land first; proof is generated against that source
  commit; committed proof artifacts must not be treated as manifest subjects for
  the proof they contain unless a new source-bound cycle is run.

### WS2 - Backlog And Block Drift Closure

Resolve or explicitly reclassify closure drift across tasks, block ledgers, and
backlog.

Known Block 06 open items:

- `sdp-trace-drq.11`: Block 06 runner accepts shell-shaped Bazel command.
- `sdp-trace-drq.12`: Block 06 package validator misses proof-state evidence
  refs.

Acceptance criteria:

- Open Block 06 items are either fixed with current Go-first validation or
  explicitly carried as MVP-blocking follow-ups with `not_assessed` closure.
- Task and block ledgers distinguish historical retired commands from current
  closure evidence.
- Historical ledger entries that mention removed active commands are annotated
  with current invalidation status or moved under an explicitly historical
  section. A grep for retired commands in current-closure sections must return
  no matches.
- The retired-command scan starts with the retired verifier list in
  retired research artifact and any additional
  removed commands found by the implementation audit.
- Block 22 status and review ledger agree on implementation and PR-review state.
- No checked task claims a removed active command as current closure evidence.
- Local `main` ahead of `origin/main` is treated as process/release state, not a
  product defect: final remote closure remains `not_assessed` until the change
  is reviewed, pushed through the agreed PR/merge path, and verified on
  `origin/main`.

### WS3 - Go Quality Gate

Turn the current quality review into a measurable gate instead of an informal
claim.

Quality thresholds for this block:

- `gofmt -l $(rg --files -g '*.go')` prints nothing.
- `go test ./...`, `go vet ./...`, `staticcheck ./...`, and
  `golangci-lint run ./...` exit `0`.
- `gocyclo -over 15` prints no production functions in changed files.
- Any trust-adjacent production function with cyclomatic complexity above `15`
  must have focused tests or an explicit exception row.
- Changed production files in trust-adjacent packages must have file coverage at
  or above `70%`, or an exception row with scope, reason, and follow-up.
  Package-level coverage is reported for context, but a one-line change does not
  force full-package remediation unless the changed behavior is untestable at
  file/function scope.
- CRAP is computed for changed production functions only using:
  `CRAP = complexity^2 * (1 - coverage)^3 + complexity`, where coverage is
  the changed function's statement coverage normalized to a `0.0` through `1.0`
  fraction from `go tool cover -func` output. Changed production functions must
  have CRAP `< 5`; legacy functions above that threshold require an exception
  row and may not be called CRAP-clean.

Quality artifact format:

- retired research artifact contains command outputs,
  complexity rows, changed-file/function coverage rows, CRAP rows, and exception
  rows.
- Exception rows use:

  | id | scope | metric | observed | threshold | reason | owner | follow-up |
  | --- | --- | --- | --- | --- | --- | --- |
  | EXAMPLE | `internal/example/foo.go:bar` | `crap` | `6.2` | `<5` | staged legacy branch not changed by behavior | `role:sdp-trace-maintainer` | `sdp-trace-followup-id` |

CRAP extraction strategy:

- Generate `go test -coverprofile=<file> ./...`.
- Use `go tool cover -func=<file>` for per-function coverage percentages.
- Normalize each changed function's percentage by dividing by `100`.
- Join changed production functions to `gocyclo` function complexity by package,
  function name, and file path.

Acceptance criteria:

- Current required checks pass: `go test ./...`, `go vet ./...`,
  `staticcheck ./...`, `golangci-lint run ./...`, `gofmt`, `git diff --check`,
  and JSON syntax checks.
- Complexity and coverage are reported with thresholds and explicit exceptions.
- CRAP < 5 is measured for changed production functions by the formula above.
  Repository-wide CRAP < 5 is not claimed unless every production function has a
  measured score below `5`.
- Trust-adjacent packages with low or absent coverage are either covered by
  focused tests or recorded with a concrete `not_assessed` reason and follow-up.
- Parked/dead-code candidates in `internal/contract`, `internal/export`, and
  `internal/policy` are audited with `rg` import checks plus any available
  `golang.org/x/tools/cmd/deadcode` installation. Each package is either used,
  removed, or explicitly documented as staged non-MVP code.
- `golangci-lint run ./...` uses the repository's current lint configuration or
  the tool defaults if no config exists. Adding or relaxing lint configuration
  is outside this block unless a Socratic spec update approves that change.

### WS4 - Bilingual Tool Documentation

Document the shipped tool and profile surface in English and Russian.

Documentation coverage means each command/profile has:

- English and Russian names or stable command IDs;
- one-sentence purpose;
- minimum viable invocation example;
- output/trust-state boundary;
- exit-code or `not_assessed` / `cannot_verify` behavior when relevant;
- explicit caveat for `gate`, `witness`, and `release-proof` that local or CI
  evidence is not external production trust unless the external profile passes.

Acceptance criteria:

- English and Russian docs cover all current `sdp-trace --help` commands:
  `wrap`, `run`, `dry-run`, `preview`, `doctor`, `verify`, `explain`, `query`,
  `query-pack`, `export cross-repo-posture`, `assess`, `report`, `gate`,
  `witness`, `release-proof`, and `validate-fixtures`.
- English and Russian docs cover current profiles and states:
  `adapter-capture`, `managed-harness`, `forensic-retention`, GitHub Actions,
  GitLab CI, Buildkite, customer PKI, air-gapped fixture guidance,
  `not_assessed`, and `cannot_verify`.
- Stale Block 12-only wording is replaced by MVP-current wording without
  claiming managed harness enforcement, external audit proof, production trust,
  or enterprise support guarantees unless the same paragraph names the exact
  current profile boundary and residual unverified state.

### WS5 - Customer Question Coverage

Create a customer-facing answer map for the known buyer questions.

Mandatory customer questions for MVP handoff:

1. How is this better than CI logs, git diff, and review comments alone?
2. Can it attach sidecar-first without replacing the existing harness?
3. What happens when an agent or developer bypasses the wrapper or adapter?
4. How does it distinguish real process movement from missing telemetry?
5. Which states are `not_assessed` or `cannot_verify`, and what should a CTO do
   with them?
6. What proof exists for a local source-bound release, and what external
   production trust is still absent?
7. Which artifact answers each question in `.sdp-trace-report/`,
   `.sdp-trace-runs/`, examples, or SpecKit docs?
8. What data is intentionally not captured or not committed for privacy?
9. Which CI and enterprise witness profiles are shipped now, and which remain
   unsupported or documentation-only?

Deliverable format:

- retired customer-question map
- retired customer-question map if Russian-language handoff remains in MVP
  scope. If Russian is removed from MVP scope, the English file must explicitly
  record Russian customer questions as `not_assessed` follow-up.

Each question must map to at least one artifact or command, plus a residual
state if the answer depends on unavailable evidence.

Acceptance criteria:

- CTO/team docs answer why `sdp-trace` is better than CI logs, git diff, and
  review comments alone.
- Docs explain how the product separates real process movement from missing
  telemetry and `not_assessed` gaps.
- Docs explain what happens when agents bypass the wrapper or adapter.
- Docs explain what the customer must inspect in `.sdp-trace-report/` and
  which artifact answers each question.
- Answers stay bilingual where customer-facing.

### WS6 - Closure Package And Review Evidence

Produce a closure package that can survive independent review.

Closure package format:

- retired research artifact
- retired research artifact
- retired research artifact
- retired research artifact

`block-23-mvp-closure-package.md` contains:

- source commit and branch state;
- fixed blockers;
- unresolved blockers;
- deferred follow-ups;
- verification command table;
- source-bound proof result;
- documentation parity result;
- customer-question coverage result;
- links to the quality report, registry, and review disposition.

`block-23-not-assessed-registry.md` rows use:

| id | scope | state | reason | evidence | follow-up |
| --- | --- | --- | --- | --- | --- |
| EXAMPLE | `external_production_trust` | `not_assessed` | no external trust profile evidence in this repo | `release-proof` output | accepted non-goal |

`block-23-review-disposition.md` rows use the review-ledger fields from this
spec: reviewer/source, date, plane, finding, disposition, evidence, and
re-review state.

Acceptance criteria:

- The closure package contains: current findings, dispositions,
  `not_assessed`/`cannot_verify` registry, verification command outputs,
  source-bound proof result, documentation parity result, quality report, and
  customer-question coverage map.
- Review artifacts record reviewer/source, date, plane, findings, disposition,
  evidence, and re-review state for requirements/product, trace/evidence, and
  code-quality/docs planes.
- Implementation planning after approval uses independent worktrees only where
  write scopes do not overlap. Expected order:
  WS1 source-bound proof and WS2 ledger drift first; WS3 quality fixes second;
  WS4/WS5 docs after current command/profile surface is stable; WS6 closure
  package last.
- Each slice has focused tests, scoped commit, review disposition, and drift
  check.
- Slice drift check means:
  - `git diff --name-only` matches the approved slice write scope;
  - the relevant verification commands for the touched surface pass;
  - docs slices run command/help parity checks;
  - proof/manifest slices run manifest subject drift checks;
  - no slice closes a task or ledger using stale proof generated before its own
    changes.
- PR-level review repeats code/correctness, tracing/evidence, and
  requirements-vs-implementation planes.
- GitHub checks are recorded as `pass`, `fail`, or `not_assessed`; absent checks
  are never called green.
- After merge, `origin/main` is verified with fresh local checks and the
  source-bound proof state is re-run. CI configuration integrity is outside this
  block and remains `not_assessed` unless separately verified.

## Review Questions

1. Does this spec correctly classify every current MVP closure blocker, or is
   any item too broad for one block?
2. Should open Block 06 issues block MVP closure, or should they become explicit
   non-MVP follow-ups with customer-facing caveats?
3. What is the minimum honest quality bar for MVP: full CRAP measurement, a
   Go-first proxy gate, or `not_assessed` plus tracked follow-up?
4. Which customer questions are mandatory before pilot/customer handoff, and
   which can stay in research-only docs?
5. Is bilingual command documentation required for every command, or only for
   customer-facing workflows?
6. Does source-bound proof repair require a separate block before the docs and
   quality work, or can it be one coordinated MVP closure block with separate
   commits?

## Verification Plan

Minimum local checks:

```bash
go test ./...
go vet ./...
staticcheck ./...
golangci-lint run ./...
gofmt -l $(rg --files -g '*.go')
jq empty schema/*.json
git diff --check HEAD
go run ./cmd/sdp-trace --help
go run ./cmd/sdp-trace release-proof --manifest examples/contract-foundation/contract-manifest.example.json --out examples/contract-foundation/contract-release-verification.example.json
bd ready
```

Expected command states:

- `go test ./...`: exit `0`.
- `go vet ./...`: exit `0`.
- `staticcheck ./...`: exit `0`, no findings.
- `golangci-lint run ./...`: exit `0`, no findings.
- `gofmt -l ...`: exit `0`, no output.
- `jq empty schema/*.json`: exit `0`.
- `git diff --check HEAD`: exit `0`, no output.
- `go run ./cmd/sdp-trace --help`: exit `0`, command list matches docs.
- `release-proof`: exit `0` and JSON states named in WS1.
- `bd ready`: any open issue printed is either fixed before closure or recorded
  in the closure package with MVP scope and residual state.

Additional checks for changed surfaces:

- fixture validation for every touched example directory;
- command/help parity check between docs and `--help`;
- bilingual doc parity scan for shipped tools and profile names;
- coverage and complexity report for changed Go packages;
- manifest-subject drift check after source-bound proof regeneration.
- manifest subject diff before and after proof regeneration, using the manifest
  artifact path list, to prove no subject changed after the proof target commit.

## Initial Finding Disposition

| id | severity | plane | finding | disposition |
| --- | --- | --- | --- | --- |
| MVP-01 | critical | trace/evidence | Current `release-proof` fails on 5 manifest subject mismatches. | unresolved_blocker |
| MVP-02 | major | backlog/drift | Block 06 has two open Beads while historical ledger claims closure through removed scripts. | unresolved_blocker |
| MVP-03 | major | code-quality | CRAP < 5 is not measured; complexity and coverage risks remain. | unresolved_blocker |
| MVP-04 | critical | documentation | Russian docs do not cover shipped CLI/profile surface. | unresolved_blocker |
| MVP-05 | critical | documentation | CTO/team docs are stale against shipped MVP scope. | unresolved_blocker |
| MVP-06 | major | product/questions | Customer pressure-point questions lack explicit answer mapping. | unresolved_blocker |
| MVP-07 | minor | spec-drift | Block 22 spec status still says implementation is blocked. | unresolved_blocker |
| MVP-08 | major | release/process | Local `main` is ahead of `origin/main`; remote closure and final CI are `not_assessed`. | unresolved_blocker |

MVP-04 is critical only if Russian-language handoff remains in MVP scope. If the
CTO removes Russian from MVP closure scope, MVP-04 becomes a major follow-up and
WS4 must be narrowed before implementation approval.
