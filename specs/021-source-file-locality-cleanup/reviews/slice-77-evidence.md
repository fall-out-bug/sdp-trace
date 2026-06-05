# Slice 77 Evidence

Status: pass

Date: 2026-06-04

Scope:
- `internal/packet/packet_033_demofirstpacketchecker.go`
- `internal/packet/packet_034_bundlevalidator.go`
- `internal/packet/packet_035_loadbundle.go`
- `internal/packet/packet_036_loadgithubinput.go`
- `internal/packet/packet_037_buildfromgithubinput.go`
- `internal/packet/packet_038_githubpacket.go`
- `internal/packet/packet_039_appendpromptboundaryfinding.go`
- `internal/packet/packet_040_githubbundlemanifest.go`
- `internal/packet/packet_041_classifypromptboundary.go`
- `internal/packet/packet_042_classifypromptmetadata.go`
- `internal/packet/packet_043_classifyprompttext.go`
- `internal/packet/packet_044_promptboundarymetadatamissing.go`
- `internal/packet/packet_045_promptboundarymetadatacomplete.go`
- `internal/packet/packet_046_forbiddenrecorderdutyphrases.go`

Plan/task review:
- Lane 1: `LGTM` (`019e936a-f1fd-7e51-ba62-55150e632407`, Confucius)
- Lane 2: `LGTM` (`019e936a-f541-7c71-b178-b630fbbb648e`, Herschel)
- Lane 3: `LGTM` after fix (`019e936a-fce1-70c3-be81-18499da97328`,
  Plato)
- Findings fixed: source-shape evidence for `packet_033` through `packet_046`
  and explicit `cmd/sdp-trace/FAMILY_INDEX.md` not-applicable handling.

Implementation boundary:
- Moved validation context structs into `internal/packet/validation_contexts.go`.
- Moved JSON loading entrypoints into `internal/packet/input_loading.go`.
- Moved GitHub bundle assembly entrypoint into
  `internal/packet/github_build_entrypoint.go`.
- Moved GitHub packet shell construction into
  `internal/packet/github_packet_shell.go`.
- Moved prompt-boundary theater finding append into
  `internal/packet/github_prompt_finding.go`.
- Moved bundle manifest construction into
  `internal/packet/github_bundle_manifest.go`.
- Moved prompt-boundary classification entrypoint into
  `internal/packet/prompt_boundary_entrypoint.go`.
- Moved prompt-boundary metadata classification and completeness checks into
  `internal/packet/prompt_boundary_metadata.go`.
- Moved prompt-boundary text contamination checks into
  `internal/packet/prompt_boundary_text.go`.
- Moved recorder-duty phrase catalog into
  `internal/packet/prompt_boundary_forbidden_phrases.go`.
- No rendering, packet digesting, validation execution, downstream GitHub row
  helpers, downstream GitHub entry helpers, package boundary, dependency
  direction, or baseline changes were made.
- `cmd/sdp-trace/FAMILY_INDEX.md` was not changed because it indexes command
  families and is not authoritative for non-command `internal/packet` package
  locality.

Focused packet regression evidence:
- Existing test retained: `TestLoadBundleAndGitHubInput`.
- Added `TestLoadBundleAndGitHubInputRejectMalformedJSON`.
- Existing test retained: `TestClassifyPromptBoundary`.
- Added `TestClassifyPromptBoundaryReasons`.
- Added `TestForbiddenRecorderDutyPhrasesPreserveCatalog`.
- Existing test retained: `TestBuildGitHubPromptBoundaryRequiredBlocksAgentRoute`.
- Added `TestBuildFromGitHubInputPreservesPacketShellAndManifest`.
- Added `TestBuildFromGitHubInputContaminatedPromptFindingShape`.
- Existing test retained: `TestBuildFromGitHubInputRedactsSecretLikeResolvers`.

Verification:
- verified: `go test ./internal/packet -list 'Test(LoadBundleAndGitHubInput|LoadBundleAndGitHubInputRejectMalformedJSON|ClassifyPromptBoundary|ClassifyPromptBoundaryReasons|ForbiddenRecorderDutyPhrasesPreserveCatalog|BuildGitHubPromptBoundaryRequiredBlocksAgentRoute|BuildFromGitHubInputPreservesPacketShellAndManifest|BuildFromGitHubInputContaminatedPromptFindingShape|BuildFromGitHubInputRedactsSecretLikeResolvers)$' | awk '/^Test/ {print}'`
- verified: `go test ./internal/packet`
- verified: `rg --files internal/packet | rg 'packet_0(33|34|35|36|37|38|39|40|41|42|43|44|45|46)_'` returned no matches with exit code `1`.
- verified: `test -f internal/packet/validation_contexts.go`
- verified: `test -f internal/packet/input_loading.go`
- verified: `test -f internal/packet/github_build_entrypoint.go`
- verified: `test -f internal/packet/github_packet_shell.go`
- verified: `test -f internal/packet/github_prompt_finding.go`
- verified: `test -f internal/packet/github_bundle_manifest.go`
- verified: `test -f internal/packet/prompt_boundary_entrypoint.go`
- verified: `test -f internal/packet/prompt_boundary_metadata.go`
- verified: `test -f internal/packet/prompt_boundary_text.go`
- verified: `test -f internal/packet/prompt_boundary_forbidden_phrases.go`
- failed then fixed: first MI run reported missing baselines for below-threshold
  `internal/packet/github_bundle_build.go` and
  `internal/packet/prompt_boundary_classification.go`; those broad files were
  split into narrower responsibility files instead of adding MI baselines.
- verified after split: `go test ./...`
- verified after split: `go vet ./...`
- verified after split: `golangci-lint run`
- verified after split: `go run ./tools/doccheck`
- verified after split: `go run ./tools/hygienecheck`
- verified after split: `jq empty schema/*.json`
- verified after split: `git diff --check`
- verified after split: `go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools`
- verified after split: `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`
- verified after split: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`
- cannot_verify extended scope: `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal tools` fails on pre-existing, unrelated `tools/hygienecheck/check_demo_drift.go:checkCurrentDemoRepoDrift`; the repository-local Slice 70, 73, 74, 75, and 76 evidence pattern uses function MI for `cmd internal`, file MI for `cmd internal tools`, and CRAP/gocyclo for `cmd internal tools`.
- verified after split: `go test -count=1 ./... -coverprofile=coverage.out`
- verified after split: `go tool cover -func=coverage.out > coverage-func.txt`
- verified after split: `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`
- verified after split: `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`

Implementation review:
- Lane 1: `LGTM` (`019e9373-4ac1-7d43-a19a-958a4719cb13`, Huygens)
- Lane 2: `LGTM` after evidence fix
  (`019e9373-4eae-75d0-a3af-aef15bd153f0`, Sagan)
- Lane 3: `LGTM` (`019e9373-566c-7903-bb1c-3e0331971276`, Lorentz)

Implementation review findings:
- major: evidence initially implied the function-MI gate covered `cmd internal
  tools`; the extended command was not passing because of unrelated
  pre-existing tools debt.

Fixes after review:
- Ran the extended function-MI command, recorded its unrelated failure as
  `cannot_verify`, and kept the established repository-local MI evidence split
  explicit.
- verified: Lane 2 returned `LGTM` after the evidence correction.
