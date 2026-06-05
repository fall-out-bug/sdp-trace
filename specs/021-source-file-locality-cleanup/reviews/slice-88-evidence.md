# Slice 88 Evidence

Date: 2026-06-04

Scope: `internal/prreview` packet construction and packet input reference
shards `prreview_020` through `prreview_037`.

## Locality

- Removed numbered packet construction/ref shards `prreview_020` through
  `prreview_037`.
- Moved build orchestration and packet write behavior into `packet_build.go`.
- Moved packet identity/provenance attachment into `packet_identity.go`.
- Moved default `CreatedBy` / CI state handling into `packet_defaults.go`.
- Moved packet ref aggregation into `packet_refs.go`.
- Moved copied diff/metadata/context/verification ref conversion into
  `packet_ref_inputs.go`.
- Moved explicit unavailable-field construction into
  `packet_unavailable_fields.go`.
- Run execution, ledger/validation logic, option validation, generic copy
  utilities, prompt generation, and lower-level IO helpers are intentionally
  excluded from this slice.

## Source Shape

```text
$ find internal/prreview -maxdepth 1 -type f | sed 's#^#/#' | rg '/prreview_0(2[0-9]|3[0-7])_[^/]+\.go$' || true
<no output>

$ git diff --cached --name-only | rg '^internal/prreview/prreview_(0(3[8-9]|[4-9][0-9])|1[0-9][0-9])_' || true
<no output>
```

## Verification

`pass`:

```text
gofmt -w internal/prreview/*.go
tests='TestBuildPacketBindsRefsAndRejectsUnsafeIdentity|TestBuildPacketRecordsUnavailableInputsAndDigestChangesWithDiff|TestPrreviewPacketBuildHelpersPreserveDefaultsRefsAndUnavailableFields'; listed=$(go test ./internal/prreview -list "$tests" | rg "^($tests)$"); test "$(printf '%s\n' "$listed" | sed '/^$/d' | wc -l | tr -d ' ')" = 3 && go test ./internal/prreview -run "$tests" -count=1 -v
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

- Round 1 Lane A behavior/correctness: Beauvoir `LGTM`.
- Round 1 Lane B locality/boundary/MI: Halley `LGTM`.
- Round 1 Lane C tests/evidence: Peirce reported a major test-evidence gap:
  unsafe identity rejection did not prove the output directory was not created
  before identity validation.
- Fix: `TestBuildPacketBindsRefsAndRejectsUnsafeIdentity` now records the
  unsafe output path and asserts `os.IsNotExist` after `unsafe_repo_id`.
- Round 2 Lane C tests/evidence: Peirce `LGTM`.
