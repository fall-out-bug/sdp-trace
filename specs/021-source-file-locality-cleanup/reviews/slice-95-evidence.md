# Slice 95 Evidence

Date: 2026-06-05

Scope:
- Consolidated `internal/prreview/prreview_087` through
  `internal/prreview/prreview_099` packet option, run-set, and review profile
  validation helpers into responsibility-named files.
- Added focused regression coverage for exact option/schema validation error
  contracts.
- Excluded role execution, runner command handling, prompt, JSON artifact IO,
  validation/summary, sanitizer, citation, and packet construction internals.

Changed source shape:
- Added `internal/prreview/packet_option_validation.go`.
- Added `internal/prreview/packet_identity_validation.go`.
- Added `internal/prreview/runset_validation.go`.
- Added `internal/prreview/profile_validation.go`.
- Added `internal/prreview/profile_role_validation.go`.
- Deleted numbered files `internal/prreview/prreview_087_validatepacketoptions.go`
  through `internal/prreview/prreview_099_validaterequiredplaneroles.go`.

Focused verification:

```sh
gofmt -w internal/prreview/*.go &&
tests='TestBuildPacketBindsRefsAndRejectsUnsafeIdentity|TestValidateProfileRejectsMalformedProfiles|TestReadRunSetRejectsDuplicateRunIDs|TestRunReviewPreservesValidationDefaultsAndOutputContracts|TestPrreviewOptionSchemaValidationContracts';
listed=$(go test ./internal/prreview -list "$tests" | rg "^($tests)$");
test "$(printf '%s\n' "$listed" | sed '/^$/d' | wc -l | tr -d ' ')" = 5 &&
go test ./internal/prreview -run "$tests" -count=1 -v &&
go test ./internal/prreview &&
go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal &&
go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools
```

Result: verified pass.

Repository verification:

```sh
go test ./... &&
go vet ./... &&
golangci-lint run &&
go run ./tools/doccheck &&
go run ./tools/hygienecheck &&
jq empty schema/*.json &&
git diff --check
```

Result: verified pass.

CRAP and source-boundary verification:

```sh
go test -count=1 ./... -coverprofile=coverage.out &&
go tool cover -func=coverage.out > coverage-func.txt &&
go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt &&
go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less &&
rm -f coverage.out coverage-func.txt gocyclo.txt
```

Result: verified pass.

Source-boundary verification:

```sh
! find internal/prreview -maxdepth 1 -type f | sed 's#^#/#' | rg '/prreview_0(8[7-9]|9[0-9])_[^/]+\.go$' &&
! git diff --name-only | rg '^internal/prreview/prreview_(1[0-9][0-9])_'
```

Result: verified pass.

Implementation review:
- Lane A: harness `multi_agent_v1`, agent id
  `019e9406-f078-7fd2-b8d0-e22ac17a1e3a`, model/provider `not_assessed`
  (tool does not expose exact inherited model/provider), date 2026-06-05,
  prompt class implementation review, timeout 300000 ms, retries 0, fallback
  none, result `LGTM`.
- Lane B: harness `multi_agent_v1`, agent id
  `019e9406-f40c-79f1-904e-54d0f0b73866`, model/provider `not_assessed`
  (tool does not expose exact inherited model/provider), date 2026-06-05,
  prompt class implementation review, timeout 300000 ms, retries 0, fallback
  none, result `LGTM`.
- Lane C: harness `multi_agent_v1`, agent id
  `019e9406-f7c2-7f80-80d9-86f7cf7e0c22`, model/provider `not_assessed`
  (tool does not expose exact inherited model/provider), date 2026-06-05,
  prompt class implementation review, timeout 300000 ms, retries 0, fallback
  none, result `LGTM`.
