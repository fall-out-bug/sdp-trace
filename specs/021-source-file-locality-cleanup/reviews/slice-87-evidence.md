# Slice 87 Evidence

Date: 2026-06-04

Scope: `internal/prreview` portable schema/type shards `prreview_001` through
`prreview_019`.

## Locality

- Removed numbered portable schema/type shards `prreview_001` through
  `prreview_019`.
- Moved constants and package vars into `constants.go`.
- Moved packet/run options into `options.go`.
- Moved packet reference types into `packet_types.go`.
- Moved review profile, run, result, ledger, validation, and plane result types
  into `review_types.go`.
- Behavior logic from `prreview_020` onward is intentionally excluded from this
  slice.

## Source Shape

```text
$ find internal/prreview -maxdepth 1 -type f | sed 's#^#/#' | rg '/prreview_(00[1-9]|01[0-9])_[^/]+\.go$' || true
<no output>

$ git diff --cached --name-only -- internal/prreview/prreview_0[2-9][0-9]_*.go internal/prreview/prreview_1[0-9][0-9]_*.go || true
<no output>
```

## Verification

`pass`:

```text
gofmt -w internal/prreview/*.go
tests='TestPrreviewSchemaConstantsAndPatternsPreserveContracts|TestPrreviewPortableTypesPreserveJSONShape'; listed=$(go test ./internal/prreview -list "$tests" | rg "^($tests)$"); test "$(printf '%s\n' "$listed" | sed '/^$/d' | wc -l | tr -d ' ')" = 2 && go test ./internal/prreview -run "$tests" -count=1 -v
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

Round 1:

- Beauvoir behavior/API: `LGTM`.
- Halley locality/boundary/MI: `LGTM`.
- Peirce tests/evidence: major findings. Focused tests asserted many
  constants only as non-empty and did not cover all moved portable structs.

Fixes:

- `TestPrreviewSchemaConstantsAndPatternsPreserveContracts` now asserts exact
  string values for all moved constants and exact regex pattern strings.
- `TestPrreviewPortableTypesPreserveJSONShape` now covers all moved portable
  structs, including `SafeRef`, `UnavailableField`, `ReviewProfile`,
  `RunPreview`, `PreviewRole`, `RunSet`, `Finding`, `Citation`, `Ledger`, and
  `LedgerFinding`, with present and omitted JSON key checks.
- Focused tests, package tests, MI, full repo gates, and CRAP were rerun after
  the fixes.

Round 2:

- Peirce tests/evidence re-review: `LGTM`.
