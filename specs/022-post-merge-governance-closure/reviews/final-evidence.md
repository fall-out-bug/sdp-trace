# Spec 022 Final Evidence Map

Date: 2026-06-01

Branch: `codex/022-post-merge-governance-closure`

## Local Verification

| Check | Command | State | Notes |
| --- | --- | --- | --- |
| Tests | `go test -count=1 ./...` | pass | All packages passed. |
| Vet | `go vet ./...` | pass | No vet failures. |
| Docs command surface | `go run ./tools/doccheck` | pass | No doccheck failures. |
| Repository hygiene | `go run ./tools/hygienecheck` | pass | No hygiene failures. |
| Schema docs | `go run ./tools/schemadoc` | pass | No schemadoc failures. |
| JSON syntax | `for f in schema/*.json; do jq empty "$f"; done` | pass | All schema JSON parsed. |
| Diff whitespace | `git diff --check` | pass | No whitespace errors. |
| Coverage profile | `go test -count=1 ./... -coverprofile=coverage.out` | pass | Coverage profile generated for CRAP check. |
| CRAP strict-less | `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less` | pass | All reported functions are below strict `< 5` threshold. |
| Cyclomatic/cognitive | `go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools` | pass | No threshold failures. |
| File MI baseline | `go run ./tools/qualitycheck -fail-only -mi-under 70 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools` | pass | No baseline regression. |
| Function MI baseline | `go run ./tools/qualitycheck -fail-only -function-mi-under 70 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal` | pass | No baseline regression. |
| Absolute file MI `> 70` | `go run ./tools/qualitycheck -fail-only -mi-under 70.1 cmd internal tools` | pass | No file-level absolute MI failure observed in this run. |
| Absolute function MI `> 70` | `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 cmd internal tools` | fail | Historical function `tools/hygienecheck/check_demo_drift.go:10 checkCurrentDemoRepoDrift` reports MI `61.5`; this is an assessed gap, not a pass claim. |

## Live External Evidence

| Surface | State | Evidence |
| --- | --- | --- |
| PR #60 | verified | `gh pr view 60 --json number,state,mergeCommit,headRefOid,reviewDecision,statusCheckRollup,url` on 2026-06-01 returned `state=MERGED`, merge commit `657a343a5f310538def9afd509e6c610c713cab0`, head `977179ba93e577cc51f05e634453deba80c383b6`, empty `reviewDecision`, and CI `verify` success. |
| PR #63 | verified | `gh pr view 63 --json number,state,mergeCommit,headRefOid,reviewDecision,statusCheckRollup,url` on 2026-06-01 returned `state=MERGED`, merge commit `1ee2c7af53637c7f43bff4e0e7ef9e34d164908e`, head `9ad4640828a204e66fe056dde0950542ed00f1c4`, empty `reviewDecision`, CI `verify` success, and `pr-review-evidence-only` success. |

## Trust Boundaries

- PR #60 merge approval remains `not_assessed`.
- Spec 022 does not claim retroactive pre-implementation approval.
- Spec 022 does not claim production trust, release approval, or external
  attestation.
- Absolute function MI `> 70` remains an assessed quality gap because one
  historical function is below `70.1`.
