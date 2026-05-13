# CI Check Policy

Pull requests are expected to report GitHub Actions checks from
`.github/workflows/ci.yml`.

Required CI evidence:

- `go test -count=1 ./...`
- `go run ./tools/mibaselinepolicy -base-ref <base-sha>` on pull requests
- `golangci-lint run ./...`
- `jq empty schema/*.json examples/block19-adapter-capture/*.json examples/self-trace/proof-summary.example.json tools/qualitycheck/function-mi-baseline.json tools/qualitycheck/file-mi-baseline.json`
- `git diff --check`

Quality gate policy:

| Gate | Threshold | CI state | Evidence rule |
| --- | ---: | --- | --- |
| CRAP | `< 5` per Go function in `cmd`, `internal`, and `tools` | CI-enforced by `Go quality gates` | Compute from `go test -count=1 ./... -coverprofile`, `go tool cover -func`, and `go run ./tools/qualitycheck -gocyclo` output through `go run ./tools/crapcheck -strict-less`. Do not claim pass while local CRAP rows exceed the threshold. |
| Cyclomatic complexity | `<= 10` per Go function in `cmd`, `internal`, and `tools` | CI-enforced by `Go quality gates` | `go run ./tools/qualitycheck -fail-only -cyclo-over 10 cmd internal tools` must pass. `-fail-only` keeps passing CI logs quiet; failures still print to stderr. |
| Cognitive complexity | `<= 10` per Go function in `cmd`, `internal`, and `tools` | CI-enforced by `Go quality gates` | `go run ./tools/qualitycheck -fail-only -cognitive-over 10 cmd internal tools` must pass. `-fail-only` keeps passing CI logs quiet; failures still print to stderr. |
| Function Maintainability Index | ratchet toward `>= 70` per production function in `cmd` and `internal` | CI-enforced by `Go quality gates` | `go run ./tools/qualitycheck -fail-only -function-mi-under 70 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal` must pass. Historical below-threshold functions may not regress; new below-threshold functions fail. This is not an absolute MI `> 70` pass claim. |
| File Maintainability Index | ratchet toward `> 70` per production file in `cmd`, `internal`, and `tools` | CI-enforced by `Go quality gates` | `go run ./tools/qualitycheck -fail-only -mi-under 70 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools` must pass. Historical below-threshold files may not regress; new below-threshold files fail. This is not an absolute MI `> 70` pass claim. |
| Spec drift / work without spec | no implementation-only trust change | review-required, not CI-enforced | Review must compare changed behavior against the active spec, plan, and task. If no applicable spec delta exists, record `cannot_verify` or `not_assessed` instead of pass. |
| CleanCode / CleanArchitecture review | behavior remains small, cohesive, and dependency-direction neutral | review-required, not CI-enforced | Review changed Go packages for avoidable complexity, hidden product coupling, and boundary violations. Record findings with evidence refs; absence of review is `not_assessed`. |
| Security review | trust, path, network, secret, and authority changes reviewed | review-required, not CI-enforced | Any release proof, witness, harness, credential, external input, or path-handling change needs a security review row. Missing evidence is `cannot_verify` if review was required for the slice. |
| DX review | documented commands match live CLI and failure modes are actionable | review-required, not CI-enforced | Compare affected docs and examples against `go run ./cmd/sdp-trace --help` or subcommand help. Missing command evidence is `cannot_verify` for command-surface changes. |
| UX review | human-facing packet, report, and explanation output remains readable and non-misleading | review-required, not CI-enforced | For changed rendered output or reviewer/operator workflows, inspect the generated artifact or mark the plane `not_assessed` with rationale. |
| Docs completeness | command and trust-scope docs updated with behavior changes | review-required, not CI-enforced | For command-surface changes, compare `docs/agent-entrypoint.md` with `go run ./cmd/sdp-trace --help` and update affected docs in the same slice. |

Selected coverage floors for this slice:

| Surface | Floor | Current gate state |
| --- | ---: | --- |
| MVP-critical packages under `internal/trace`, `internal/contract`, `internal/policy`, `internal/export`, `internal/posture`, `internal/harnessobs`, and `internal/verifier` | `>= 75%` package statement coverage | Locally satisfied in the 2026-05-12 run recorded in `specs/004-mvp-readiness-hardening/implementation-ledger.md`; CI records coverage input for CRAP but does not enforce package floors yet. |
| Changed production packages under `cmd` and `internal` | `>= 75%` package statement coverage | Locally satisfied in the same run. |
| Go tooling under `tools` | `>= 75%` package statement coverage | Locally satisfied in the 2026-05-12 run after `tools/crapcheck`, `tools/mibaselinepolicy`, and `tools/qualitycheck` refactors; absolute tool MI remains an `assessed_gap`. |

Regenerate `tools/qualitycheck/function-mi-baseline.json` or
`tools/qualitycheck/file-mi-baseline.json` only as a reviewed ratchet change.
The baselines are not proof that MI passed; they are the current exception set
for historical functions/files below the absolute threshold. After a baseline
exists on the target branch, pull requests may not change it in the same diff as
production Go files under `cmd/`, `internal/`, or `tools/`; CI fails that
combination with `go run ./tools/mibaselinepolicy -base-ref <base-sha>`, so a
product change cannot add its own MI exception. First checked-in baselines are
reviewed as part of the gate-introduction slice, not as later ratchet
regeneration.

Maintainability Index is a local ratchet/refactor heuristic, not an external
quality proof. The formula is implemented in `tools/qualitycheck/halstead.go`;
file-level MI aggregates whole Go files, so same-package file splits can improve
the metric without changing behavior. Use MI alongside CRAP, complexity,
coverage, review evidence, and spec-drift checks before making quality claims.

If GitHub does not report checks for a PR, record CI as `not_assessed`; do not
treat local verification as a substitute for remote CI evidence. Local
verification may support implementation review, but CI-backed closure requires
the workflow result or an explicit repo-tracked replacement policy.
