# Slice 86 Evidence

Date: 2026-06-04

Scope: `internal/packet` Markdown rendering shards `packet_193` through
`packet_200`.

## Locality

- Removed numbered packet rendering shards `packet_193` through `packet_200`.
- Moved behavior into cohesive packet rendering, packet section rendering, and
  existing theater rendering locality files under `internal/packet`.
- `cmd/sdp-trace/FAMILY_INDEX.md` is not applicable because this slice is in
  the non-command `internal/packet` package.

## Source Shape

```text
$ find internal/packet -maxdepth 1 -type f | sed 's#^#/#' | rg '/packet_(19[3-9]|200)_[^/]+\.go$' || true
<no output>

$ git diff --cached --name-only -- internal/prreview/*_[0-9][0-9][0-9]_*.go || true
<no output>
```

## Verification

`pass`:

```text
gofmt -w internal/packet/*.go
tests='TestValidateAndRenderHappyPath|TestRenderCleanTheaterUsesRowState|TestPacketRenderingHelpersPreserveTables|TestPacketRenderingSectionHelpersPreserveMetadataRowsAndErrors|TestRenderMarkdownPreservesTopLevelOrderAndHeaders|TestRenderMarkdownRejectsInvalidBundleBeforeProjection'; listed=$(go test ./internal/packet -list "$tests" | rg "^($tests)$"); test "$(printf '%s\n' "$listed" | sed '/^$/d' | wc -l | tr -d ' ')" = 6 && go test ./internal/packet -run "$tests" -count=1 -v
```

`pass`:

```text
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

Round 1:

- Beauvoir behavior/correctness: `LGTM`.
- Halley locality/boundary/MI: `LGTM`.
- Peirce tests/evidence: `LGTM`.
