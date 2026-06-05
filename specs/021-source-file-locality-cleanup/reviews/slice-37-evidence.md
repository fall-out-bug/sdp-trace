# Slice 37 Evidence

Status: pass

## Scope

Slice 37 removes numbered `internal/harnessobs` validation command shards
`harnessobs_281` through `harnessobs_290`.

## Changed Code

- Added `internal/harnessobs/validation_command.go`.
- Added `internal/harnessobs/validation_requirements.go`.
- Added `internal/harnessobs/validation_paths.go`.
- Added `internal/harnessobs/validation_evaluation.go`.
- Updated `internal/harnessobs/harnessobs_test.go` with focused validation
  command regression coverage.
- Deleted `internal/harnessobs/harnessobs_281_validate.go`.
- Deleted `internal/harnessobs/harnessobs_282_validatevalidateinputs.go`.
- Deleted `internal/harnessobs/harnessobs_283_requirevalidateoptions.go`.
- Deleted `internal/harnessobs/harnessobs_284_requirenonblank.go`.
- Deleted `internal/harnessobs/harnessobs_285_resolvevalidateinputs.go`.
- Deleted `internal/harnessobs/harnessobs_286_resolvevalidatesourcepaths.go`.
- Deleted `internal/harnessobs/harnessobs_287_resolvevalidateoutpath.go`.
- Deleted `internal/harnessobs/harnessobs_288_evaluationfromrun.go`.
- Deleted `internal/harnessobs/harnessobs_289_fallbacksourceunavailable.go`.
- Deleted `internal/harnessobs/harnessobs_290_writevalidationifrequested.go`.

## Plan Review

- initial scope lane (`019e8788-8f70-7cb2-bc40-bdcb718f7d0c`): finding
  fixed; final scope re-review (`019e878a-4226-76d3-82c0-22ecea63317e`):
  LGTM.
- initial trust/evidence lane (`019e8788-fbf2-7bd3-9b10-7f410835ad77`):
  finding fixed; final trust/evidence re-review
  (`019e878a-629b-7bb0-ade8-5867a8641499`): LGTM.
- initial maintainability/DX lane
  (`019e8789-17d9-7720-a091-23ccd5febec3`): finding fixed; final
  maintainability/DX re-review (`019e878a-9b18-7b93-bc79-0a83552cf9b3`):
  LGTM.
- MI-driven grouping update re-review
  (`019e878f-0f9e-7a83-9690-a95e918f2402`): LGTM.

## Focused Verification

verified:

```sh
gofmt -w internal/harnessobs/harnessobs_test.go \
  internal/harnessobs/validation_command.go \
  internal/harnessobs/validation_requirements.go \
  internal/harnessobs/validation_paths.go \
  internal/harnessobs/validation_evaluation.go
go test ./internal/harnessobs -run 'Test(ObserveValidateCompleteHarnessExport|ValidateRequiredOptionErrors|ValidateWritesCannotVerifyOutForMissingRun|ValidateWritesOutPathWhenPasses|ValidateCannotVerifyWhenRunFileInvalid|ValidateRejectsUnsafePaths|ValidateRejectsUnsafeOutBasename)'
```

Result: pass.

verified:

```sh
go run ./tools/qualitycheck -mi-under 70.1 -function-mi-under 70.1 \
  internal/harnessobs/validation_command.go \
  internal/harnessobs/validation_requirements.go \
  internal/harnessobs/validation_paths.go \
  internal/harnessobs/validation_evaluation.go \
  internal/harnessobs/harnessobs_test.go
```

Result: pass. New code files remain above the absolute file-MI threshold.

Focused regression coverage includes required option errors, whitespace-only
option rejection, unsafe profile/run path rejection, unsafe output path
rejection, optional output omission, validation output writing,
source-unavailable cannot-verify fallback, validation digest, non-authority
boundary, and normal validation from a loaded run.

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

- behavior/correctness lane (`019e8791-7bfc-7472-bcef-6395849418a1`):
  LGTM.
- trust/evidence/provenance lane (`019e8791-cf44-75c3-944b-51ed40d9ca62`):
  LGTM.
- maintainability/DX lane (`019e8793-195b-73b0-aef8-643f7f94105c`): LGTM.

Final implementation review verdict: LGTM across all three lanes.
