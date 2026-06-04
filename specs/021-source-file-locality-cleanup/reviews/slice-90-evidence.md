# Slice 90 Evidence

Date: 2026-06-04

Scope: `internal/prreview` ledger synthesis shards `prreview_043` through
`prreview_048`.

## Locality

- Removed numbered ledger synthesis shards `prreview_043` through
  `prreview_048`.
- Moved ledger synthesis, existing-ledger lookup, and run finding flattening
  into `ledger_synthesis.go`.
- Moved review finding projection, fallback finding IDs, and disposition
  carry-forward into `ledger_finding_projection.go`.
- Validation logic, summary rendering, file IO, safe text/default disposition
  helper implementation, and lower-level role execution are intentionally
  excluded from this slice.

## Source Shape

```text
$ find internal/prreview -maxdepth 1 -type f | sed 's#^#/#' | rg '/prreview_0(4[3-8])_[^/]+\.go$' || true
<no output>

$ git diff --name-only | rg '^internal/prreview/prreview_(0(4[9]|[5-9][0-9])|1[0-9][0-9])_' || true
<no output>
```

## Verification

`pass`:

```text
gofmt -w internal/prreview/*.go
tests='TestValidationAndLedgerLifecyclePreserveTrustSemantics|TestLedgerDispositionCarryForward|TestPrreviewLedgerSynthesisPreservesOrderingCarryForwardAndSanitization'; listed=$(go test ./internal/prreview -list "$tests" | rg "^($tests)$"); test "$(printf '%s\n' "$listed" | sed '/^$/d' | wc -l | tr -d ' ')" = 3 && go test ./internal/prreview -run "$tests" -count=1 -v
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
