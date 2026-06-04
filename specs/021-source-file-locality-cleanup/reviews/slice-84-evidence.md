# Slice 84 Evidence

Date: 2026-06-04

Scope: `internal/packet` GitHub source-change and row-construction shards
`packet_148` through `packet_173`.

## Locality

- Removed numbered row/source shards `packet_148` through `packet_173`.
- Moved the behavior into cohesive GitHub source, row, route, verification,
  artifact evidence, artifact retention, default-row, and review-row locality
  files under `internal/packet`.
- `cmd/sdp-trace/FAMILY_INDEX.md` is not applicable because this slice is in
  the non-command `internal/packet` package.

## Source Shape

```text
$ find internal/packet -maxdepth 1 -type f | sed 's#^#/#' | rg '/packet_(14[8-9]|15[0-9]|16[0-9]|17[0-3])_[^/]+\.go$' || true
<no output>

$ git diff --cached --name-only -- internal/packet/packet_174_*.go ... internal/packet/packet_200_*.go internal/prreview/*_[0-9][0-9][0-9]_*.go || true
<no output>
```

## Verification

`pass`:

```text
gofmt -w internal/packet/github_*.go internal/packet/packet_test.go
tests='TestGitHubSourceChangeProjection|TestGitHubChangeAndMutationRowsPreserveCommitRangeSemantics|TestGitHubInitiatorAndReviewRowsPreserveEvidenceSemantics|TestGitHubAgentRouteRowsPreservePromptBoundarySemantics|TestGitHubVerificationRowsPreserveEvidenceSemantics|TestGitHubArtifactEvidenceHelpersPreserveSemantics'; listed=$(go test ./internal/packet -list "$tests" | rg "^($tests)$"); test "$(printf '%s\n' "$listed" | sed '/^$/d' | wc -l | tr -d ' ')" = 6 && go test ./internal/packet -run "$tests" -count=1 -v
go test ./internal/packet
go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal
go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools
go test ./...
go vet ./...
golangci-lint run
go run ./tools/doccheck && go run ./tools/hygienecheck && jq empty schema/*.json && git diff --check
go test -count=1 ./... -coverprofile=coverage.out
go tool cover -func=coverage.out > coverage-func.txt
go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt
go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less
```

Temporary files `coverage.out`, `coverage-func.txt`, and `gocyclo.txt` were
removed after CRAP verification.

## Reviewer Rounds

### Round 1

#### Mill the 2nd

- Harness: Codex subagent
- Agent id: `019e93f7-d548-7d03-ac2f-1b1e2633f330`
- Model/provider: not_assessed
- Timeout/retries/fallback: 3600000 ms / 0 / none
- Prompt class: implementation behavior review
- Result: `LGTM`

#### Wegener the 2nd

- Harness: Codex subagent
- Agent id: `019e93f7-d8fd-7640-8846-e2d52e8b6ce8`
- Model/provider: not_assessed
- Timeout/retries/fallback: 3600000 ms / 0 / none
- Prompt class: tests and evidence review
- Result: `LGTM`

#### Plato the 2nd

- Harness: Codex subagent
- Agent id: `019e93f7-dcb6-7173-8f9d-d2ccb6c9ee0b`
- Model/provider: not_assessed
- Timeout/retries/fallback: 3600000 ms / 0 / none
- Prompt class: locality and boundary review
- Result: findings

Findings:

- major: `github_source_change.go` was a standalone one-function replacement
  for removed `packet_148_githubsourcechange.go`.

Fixes:

- Removed `github_source_change.go`.
- Moved `githubSourceChange` into the existing GitHub packet projection
  locality in `github_packet_shell.go`, which already contains `githubPacket`.
- Re-ran focused tests, package tests, MI gates, repository gates, and CRAP.

### Round 2

#### Plato the 2nd

- Harness: Codex subagent
- Agent id: `019e93f7-dcb6-7173-8f9d-d2ccb6c9ee0b`
- Model/provider: not_assessed
- Timeout/retries/fallback: 3600000 ms / 0 / none
- Prompt class: locality and boundary re-review
- Result: `LGTM`

## Review Verdict

pass

All three independent implementation-review lanes returned exact `LGTM` after
the standalone one-function replacement finding was fixed and re-reviewed.
