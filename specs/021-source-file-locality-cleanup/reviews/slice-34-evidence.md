# Slice 34 Evidence: Raw Unsafe Rule Semantics

Status: pass

## Scope

- Package: `internal/harnessobs`
- Source shards: `harnessobs_223` through `harnessobs_244`
- Target responsibility groups:
  - map traversal and raw-event skip rules
  - slice traversal
  - string path/token/url safety checks
  - digest-field and raw path-like exemptions
- Excluded:
  - validation enum helpers `harnessobs_245` onward
  - session collect option validation
  - event source resolution
  - raw-event normalization flow

## Local Verification

- plan review lanes: pass
  - lane 1 requirements/scope: `LGTM`
  - lane 2 trust/evidence: `LGTM`
  - lane 3 maintainability/DX: `LGTM`
- implementation: pass
- focused verification: pass
  - `gofmt -w internal/harnessobs/unsafe_value_*.go`: pass
  - `go test ./internal/harnessobs -run 'Test(FindUnsafeRawEventAtReasonCodes|ObserveRejectsUnsafeRawPromptAndDoesNotWriteRun|CollectSessionPropagatesUnsafeRawEventContent|CommandModelSafetyAndSourceDigest)'`: pass
  - `go test ./internal/harnessobs`: pass
  - `go run ./tools/qualitycheck -mi-under 70.1 -function-mi-under 70.1 internal/harnessobs/unsafe_value_*.go`: pass
- full repository gates: pass
  - `go test ./...`: pass
  - `go vet ./...`: pass
  - `golangci-lint run`: pass
  - `go run ./tools/doccheck`: pass
  - `go run ./tools/hygienecheck`: pass
  - `jq empty schema/*.json`: pass
  - `git diff --check`: pass
- coverage-backed CRAP and MI baseline gates: pass
  - `go test -count=1 ./... -coverprofile=coverage.out`: pass
  - `go tool cover -func=coverage.out > coverage-func.txt`: pass
  - `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`: pass
  - `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`: pass
  - `go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools`: pass
  - `go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`: pass
  - `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`: pass
  - `git diff --name-only origin/main...HEAD | go run ./tools/mibaselinepolicy -base-ref origin/main`: pass
- numbered Slice 34 files remaining: pass; selected files `harnessobs_223`
  through `harnessobs_244` removed
- final safety-rule groups:
  - `unsafe_value_traversal.go` for entrypoints, dispatch, and slice traversal
  - `unsafe_value_maps.go` for map traversal and map reason-code dispatch
  - `unsafe_value_raw_skips.go` for raw-event skip exemptions
  - `unsafe_value_strings.go` for string reason, path, and token checks
  - `unsafe_value_urls.go` for authenticated URL checks
  - `unsafe_value_raw_paths.go` for raw path-like field exemptions
  - `unsafe_value_exemptions.go` for digest and raw path-like token exemptions

## Reviewer Lanes

- reviewer lane 1 behavior/correctness: `LGTM`
- reviewer lane 2 trust/evidence: `LGTM` after focused evidence removed
  excluded raw-event normalization-flow coverage.
- reviewer lane 3 maintainability/DX: `LGTM` after one/two-helper microfiles
  were consolidated into traversal and string safety groups and dispatcher
  boilerplate was replaced.

## Trust States

- behavior preservation: pass
- unsafe field path rendering: pass
- raw-event skip semantics: pass
- forbidden/sensitive reason codes: pass
- string path/private-path rejection: pass
- authenticated URL detection: pass
- token/base64 detection: pass
- digest-field exemption: pass
- raw path-like exemption: pass
- CRAP < 5: pass
- changed-file MI > 70: pass
- repo MI baseline/ratchet gates: pass
- zero numbered files in Slice 34 scope: pass
- spec drift: pass
- constitution drift: not_assessed
- product drift: pass
- CleanArch hex: not_assessed
- CleanCode: pass
- SOLID: pass
- DRY: pass
- YAGNI: pass
- production trust: not_assessed
- release approval: not_assessed
- merge approval: not_assessed
