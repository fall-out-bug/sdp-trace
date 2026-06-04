# Slice 83 Evidence

Date: 2026-06-04

Scope: `internal/packet` shared artifact usability helper shards
`packet_145` through `packet_147`.

## Locality

- `packet_145_entryexpired.go`, `packet_146_passrefunverifiable.go`, and
  `packet_147_artifactaccessunverifiable.go` were removed.
- Shared artifact usability helpers now live in
  `internal/packet/artifact_evidence_usability.go`.
- Existing demo-first usability helpers were folded into the same locality to
  avoid replacing numbered shards with another tiny file.
- `cmd/sdp-trace/FAMILY_INDEX.md` is not applicable because this slice is in
  the non-command `internal/packet` package.

## Source Shape

```text
$ find internal/packet -maxdepth 1 -type f | sed 's#^#/#' | rg '/packet_14[5-7]_[^/]+\.go$' || true
<no output>

$ rg -n "func (demoUsableEntry|entryHasResolverAndDigest|syntheticEntryDigest|entryExpired|passRefUnverifiable|artifactAccessUnverifiable)|var unverifiableArtifactAccess" internal/packet/artifact_evidence_usability.go
8:var unverifiableArtifactAccess = map[string]bool{
18:func demoUsableEntry(entry BundleEntry, now time.Time) bool {
23:func entryHasResolverAndDigest(entry BundleEntry) bool {
28:func syntheticEntryDigest(entry BundleEntry) bool {
34:func entryExpired(entry BundleEntry, now time.Time) bool {
47:func passRefUnverifiable(entry BundleEntry) bool {
55:func artifactAccessUnverifiable(access string) bool {

$ git diff --cached --name-only -- internal/packet/packet_148_*.go ... internal/packet/packet_200_*.go internal/prreview/*_[0-9][0-9][0-9]_*.go || true
<no output>
```

## Verification

`pass`:

```text
gofmt -w internal/packet/artifact_evidence_usability.go internal/packet/validation_row_evidence.go internal/packet/packet_test.go
tests='TestEntryExpiredSemantics|TestPassRefUnverifiableSemantics|TestDemoFirstEvidenceUsabilityStillUsesArtifactHelpers'; listed=$(go test ./internal/packet -list "$tests" | rg "^($tests)$"); test "$(printf '%s\n' "$listed" | sed '/^$/d' | wc -l | tr -d ' ')" = 3 && go test ./internal/packet -run "$tests" -count=1 -v
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

#### Boole the 2nd

- Harness: Codex subagent
- Agent id: `019e93e7-4fc8-7b61-a9fb-6a354be49860`
- Model/provider: not_assessed
- Timeout/retries/fallback: 3600000 ms / 0 / none
- Prompt class: implementation behavior review
- Result: `LGTM`

#### Leibniz the 2nd

- Harness: Codex subagent
- Agent id: `019e93e7-cb28-75b1-b603-9ab6d3eb5719`
- Model/provider: not_assessed
- Timeout/retries/fallback: 3600000 ms / 0 / none
- Prompt class: tests and evidence review
- Result: findings

Findings:

- major: `T021-5731` required focused evidence that fails on zero matches, but
  the recorded command used plain `go test -run`, which can pass with zero
  matched tests.

Fixes:

- Reran focused evidence with `go test -list` plus an exact count guard before
  executing the focused tests.
- Updated this evidence file to record the guarded command.

#### Banach the 2nd

- Harness: Codex subagent
- Agent id: `019e93e7-cf33-7350-81e8-b475b567cd36`
- Model/provider: not_assessed
- Timeout/retries/fallback: 3600000 ms / 0 / none
- Prompt class: locality and boundary review
- Result: `LGTM`

### Round 2

#### Leibniz the 2nd

- Harness: Codex subagent
- Agent id: `019e93e7-cb28-75b1-b603-9ab6d3eb5719`
- Model/provider: not_assessed
- Timeout/retries/fallback: 3600000 ms / 0 / none
- Prompt class: tests and evidence re-review
- Result: `LGTM`

## Review Verdict

pass

All three independent implementation-review lanes returned exact `LGTM` after
the focused-test zero-match evidence finding was fixed and re-reviewed.
