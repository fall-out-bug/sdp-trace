# Slice 85 Evidence

Date: 2026-06-04

Scope: `internal/packet` GitHub manifest/default/digest helper shards
`packet_174` through `packet_192`.

## Locality

- Removed numbered manifest/default/digest shards `packet_174` through
  `packet_192`.
- Moved behavior into cohesive GitHub manifest, manifest authority,
  prompt-boundary manifest, packet defaults, resolver, and digest locality
  files under `internal/packet`.
- `cmd/sdp-trace/FAMILY_INDEX.md` is not applicable because this slice is in
  the non-command `internal/packet` package.

## Source Shape

```text
$ find internal/packet -maxdepth 1 -type f | sed 's#^#/#' | rg '/packet_(17[4-9]|18[0-9]|19[0-2])_[^/]+\.go$' || true
<no output>

$ git diff --cached --name-only -- internal/packet/packet_193_*.go internal/packet/packet_194_*.go internal/packet/packet_195_*.go internal/packet/packet_196_*.go internal/packet/packet_197_*.go internal/packet/packet_198_*.go internal/packet/packet_199_*.go internal/packet/packet_200_*.go internal/prreview/*_[0-9][0-9][0-9]_*.go || true
<no output>
```

## Verification

`pass`:

```text
gofmt -w internal/packet/*.go
tests='TestGitHubManifestEntriesPreserveAssemblyOrder|TestGitHubConditionalManifestEntriesPreserveSemantics|TestGitHubBundleEntryAuthorityAndRedactionPreserveSemantics|TestGitHubResidualGapsAndDecisionOwnersPreserveDefaults|TestGitHubResolverAndDigestHelpersPreserveSemantics'; listed=$(go test ./internal/packet -list "$tests" | rg "^($tests)$"); test "$(printf '%s\n' "$listed" | sed '/^$/d' | wc -l | tr -d ' ')" = 5 && go test ./internal/packet -run "$tests" -count=1 -v
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
- Halley locality/boundary/MI: major finding. `digestPlaceholder` moved from
  one-function replacement file into the manifest authority locality file.
- Peirce tests/evidence: major findings. Focused tests now cover retained
  source classes, authority metadata, source refs, digest-only and missing
  prompt-boundary behavior, and empty artifact/integration optional behavior.
  Boundary evidence command now enumerates `packet_193` through `packet_200`
  explicitly instead of using an in-command ellipsis.

Round 2:

- Halley locality/boundary/MI re-review: `LGTM`.
- Peirce tests/evidence re-review: `LGTM`.
