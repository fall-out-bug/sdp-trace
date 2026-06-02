# Slice 46 Evidence

Status: pass

## Scope

Slice 46 is bounded to `internal/harnessobs/harnessobs_348` through
`internal/harnessobs/harnessobs_360`.

Implemented consolidation:

- moved raw normalization orchestration, input validation, and zero-time
  fallback into `internal/harnessobs/raw_normalization.go`
- moved OpenCode raw file opening and scanning into
  `internal/harnessobs/raw_normalization_scan.go`
- moved raw line decoding and unsafe raw-event rejection into
  `internal/harnessobs/raw_normalization_line.go`
- moved normalized source digest calculation into
  `internal/harnessobs/raw_normalization_digest.go`
- moved normalized output file creation and JSONL writing into
  `internal/harnessobs/raw_normalization_writer.go`
- moved shared `blankJSONLLine` into the neutral
  `internal/harnessobs/event_line_parsing.go`
- removed numbered files `harnessobs_348` through `harnessobs_360`

Explicit exclusions:

- generic unsafe raw-value discovery helpers remain in existing non-numbered
  files
- OpenCode event construction helpers remain in existing non-numbered files
- path policy, raw event schema, output JSONL format, package boundary,
  dependency direction, and MI baselines remain unchanged

## Focused Verification

- pass: `go test ./internal/harnessobs -run 'Test(RawNormalizationBranches|NormalizeOpenCodeRawLineBytesRejectsMalformedAndUnsafeRawInput|NormalizeOpenCodeRawLineBytesComputesDigestForEachEvent|WriteNormalizedEventsWritesJSONL|NormalizedWriteAndShellAndSourceCommitBranches|CollectSessionNormalizesRawEventsWhenSourceIsMissing|ValidateZeroEventSourceIsNotAssessed)$'`
- pass: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 internal/harnessobs/raw_normalization.go internal/harnessobs/raw_normalization_scan.go internal/harnessobs/raw_normalization_line.go internal/harnessobs/raw_normalization_digest.go internal/harnessobs/raw_normalization_writer.go internal/harnessobs/event_line_parsing.go`
- pass: `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 internal/harnessobs/raw_normalization.go internal/harnessobs/raw_normalization_scan.go internal/harnessobs/raw_normalization_line.go internal/harnessobs/raw_normalization_digest.go internal/harnessobs/raw_normalization_writer.go internal/harnessobs/event_line_parsing.go`

## Repository Verification

- pass: `go test ./...`
- pass: `go vet ./...`
- pass: `golangci-lint run`
- pass: `go run ./tools/doccheck`
- pass: `go run ./tools/hygienecheck`
- pass: `jq empty schema/*.json`
- pass: `git diff --check`
- pass: `go test -count=1 ./... -coverprofile=coverage.out`
- pass: `go tool cover -func=coverage.out > coverage-func.txt`
- pass: `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`
- pass: `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`
- pass: `go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools`
- pass: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`
- pass: `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`
- pass: `git diff --name-only origin/main...HEAD | go run ./tools/mibaselinepolicy -base-ref origin/main`

No MI baseline files were changed. Focused MI measured the changed/new
production files at or above 70.1.

## Review Lanes

- correctness reviewer: `019e880a-df0a-7441-a766-2e5d011e344d`,
  result `LGTM`.
- trust/evidence/spec-drift reviewer:
  `019e880a-fee9-7eb0-955d-0f7e686ae208`, result `LGTM`.
- maintainability/DX reviewer: `019e880c-f85a-7e32-beca-6289ebde01d9`,
  initial result `minor`; re-reviewer
  `019e8811-04dd-7d12-8d8f-7e6c7a864760`, result `LGTM`.

## Findings

- fixed: split `TestRawNormalizationBranches` from
  `TestRawSignalUtilityBranches` so focused Slice 46 evidence no longer
  includes generic raw signal utility coverage.
- re-reviewed: maintainability/DX lane reviewed the staged diff, full staged
  raw-normalization files, shared line parser, split tests, and Slice 46
  evidence; no remaining findings.
