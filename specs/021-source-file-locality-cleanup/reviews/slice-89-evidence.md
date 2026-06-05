# Slice 89 Evidence

Date: 2026-06-04

Scope: `internal/prreview` review run orchestration shards `prreview_038`
through `prreview_042`.

## Locality

- Removed numbered review run orchestration shards `prreview_038` through
  `prreview_042`.
- Moved public run-review entrypoint and option defaults into `run_review.go`.
- Moved run execution and role iteration into `run_review_execution.go`.
- Ledger synthesis, validation logic, low-level role execution, prompt
  generation, packet validation, and IO helper implementation are intentionally
  excluded from this slice.

## Source Shape

```text
$ find internal/prreview -maxdepth 1 -type f | sed 's#^#/#' | rg '/prreview_0(3[8-9]|4[0-2])_[^/]+\.go$' || true
<no output>

$ git diff --name-only | rg '^internal/prreview/prreview_(0(4[3-9]|[5-9][0-9])|1[0-9][0-9])_' || true
<no output>
```

## Verification

`pass`:

```text
gofmt -w internal/prreview/*.go
tests='TestRunReviewArtifactPipelineRedactsUnsafeReviewerText|TestRunReviewRecordsRunnerFailureStatesAndPromptDigest|TestRunReviewPreviewReturnsPreviewOnly|TestRunReviewPreservesValidationDefaultsAndOutputContracts|TestRunReviewNotAssessedReasonDoesNotInvokeRunner|TestRunReviewCannotVerifyUnreadablePromptTemplate|TestRunReviewPromptIncludesPacketEvidence|TestRunReviewMapsTimeoutToTimedOut'; listed=$(go test ./internal/prreview -list "$tests" | rg "^($tests)$"); test "$(printf '%s\n' "$listed" | sed '/^$/d' | wc -l | tr -d ' ')" = 8 && go test ./internal/prreview -run "$tests" -count=1 -v
go test ./internal/prreview
go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal
go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools
```

`pass`:

```text
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

- Lane A behavior/correctness: Beauvoir `LGTM`.
- Lane B locality/boundary/MI: Halley `LGTM`.
- Lane C tests/evidence: Peirce `LGTM`.
