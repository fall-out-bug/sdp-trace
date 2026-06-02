# Slice 38 Evidence

Status: pass

## Scope

Slice 38 removes numbered `internal/harnessobs` profile loading and session
profile validation shards `harnessobs_291` through `harnessobs_299`.

## Changed Code

- Added `internal/harnessobs/profile_loading.go`.
- Added `internal/harnessobs/session_profile_validation.go`.
- Added `internal/harnessobs/session_profile_raw_events.go`.
- Updated `internal/harnessobs/harnessobs_test.go` with focused session
  profile/raw-event config regression coverage.
- Updated `internal/harnessobs/harnessobs_validate_profile_test.go` with
  focused `LoadProfile` validation handoff coverage.
- Deleted `internal/harnessobs/harnessobs_291_loadprofile.go`.
- Deleted `internal/harnessobs/harnessobs_292_loadsessionprofile.go`.
- Deleted `internal/harnessobs/harnessobs_293_validatesessionprofile.go`.
- Deleted
  `internal/harnessobs/harnessobs_294_validatesessionprofileidentity.go`.
- Deleted `internal/harnessobs/harnessobs_295_validatesessionprofilepaths.go`.
- Deleted `internal/harnessobs/harnessobs_296_validaterequiredsessionpaths.go`.
- Deleted `internal/harnessobs/harnessobs_297_validateraweventconfig.go`.
- Deleted `internal/harnessobs/harnessobs_298_validateraweventpair.go`.
- Deleted `internal/harnessobs/harnessobs_299_unsupportedraweventformat.go`.

## Plan Review

- scope reviewer (`019e8796-cd2c-79e0-8b01-9db6336cb318`): LGTM.
- trust/evidence reviewer (`019e8796-f015-7a11-9a03-0aa0e9251ebf`):
  LGTM.
- maintainability/DX reviewer (`019e8797-b386-7aa3-967a-321ab86711a1`):
  LGTM.

Final plan review verdict: LGTM across all three lanes.

## Focused Verification

verified:

```sh
gofmt -w internal/harnessobs/profile_loading.go \
  internal/harnessobs/session_profile_validation.go \
  internal/harnessobs/session_profile_raw_events.go \
  internal/harnessobs/harnessobs_test.go \
  internal/harnessobs/harnessobs_validate_profile_test.go
go test ./internal/harnessobs -run 'Test(LoadProfileValidatesProfile|ValidateProfile|LoadSessionProfileDefaultsAndRejectsInvalidRawConfig|LoadSessionProfileRejectsUnsafeSetupAction|SetupSessionRequireOptions|CollectSessionNormalizesRawEventsWhenSourceIsMissing|CollectSessionRegeneratesExistingNormalizedSourceFromRawEvents)'
```

Result: pass.

verified:

```sh
go run ./tools/qualitycheck -mi-under 70.1 -function-mi-under 70.1 \
  internal/harnessobs/profile_loading.go \
  internal/harnessobs/session_profile_validation.go \
  internal/harnessobs/session_profile_raw_events.go \
  internal/harnessobs/harnessobs_test.go \
  internal/harnessobs/harnessobs_validate_profile_test.go
```

Result: pass. New code files remain above the absolute file-MI threshold:

- `profile_loading.go`: MI 79.0.
- `session_profile_validation.go`: MI 70.5.
- `session_profile_raw_events.go`: MI 72.6.

Focused regression coverage includes profile load validation handoff, strict
session profile JSON rejection, session profile schema/id rejection, default
stream capture normalization, setup action validation handoff, required session
path whitespace rejection, raw-event source/format pair errors, unsupported
raw-event format rejection, and valid raw-event config acceptance.

## Repository Verification

verified:

```sh
set -e
go test ./...
go vet ./...
if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run; fi
go run ./tools/doccheck
go run ./tools/hygienecheck
jq empty schema/*.json
git diff --check
go test -count=1 ./... -coverprofile=coverage.out
go tool cover -func=coverage.out > coverage-func.txt
go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt
go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less
go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools
go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools
go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal
git diff --name-only origin/main...HEAD | go run ./tools/mibaselinepolicy -base-ref origin/main
rm -f coverage.out coverage-func.txt gocyclo.txt
```

Result: pass.

## Implementation Review

verified:

- behavior/correctness lane (`019e879c-68fb-7e51-850d-de04c1d69597`):
  LGTM.
- trust/evidence/provenance lane (`019e879c-bdbb-7143-b939-eac9f977b3b9`):
  LGTM.
- maintainability/DX lane (`019e879e-0fbf-7521-98ef-cc6a2cd9cd6a`): LGTM.

Final implementation review verdict: LGTM across all three lanes.
