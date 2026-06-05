# Slice 76 Evidence

Status: pass

Date: 2026-06-04

Scope:
- `internal/packet/packet_001_const.go` through
  `internal/packet/packet_032_validation.go`
- `internal/packet/packet_test.go`
- `specs/021-source-file-locality-cleanup/plan.md`
- `specs/021-source-file-locality-cleanup/tasks.md`
- `specs/021-source-file-locality-cleanup/reviews/slice-76-plan-review.md`

Plan/task review:
- Lane 1: `LGTM` after fix (`019e9359-811f-7023-8a9c-6283c481770f`,
  Copernicus)
- Lane 2: `LGTM` after fix (`019e9359-84ca-7070-b507-545ab955ebe5`,
  Zeno)
- Lane 3: `LGTM` after fix (`019e9359-87fe-7e21-bc5b-901057a12ff9`,
  Poincare)
- Finding fixed: the plan/tasks now require focused JSON-shape evidence for
  `PromptBoundaryClassification` and `BuildPRResult`.

Implementation boundary:
- Moved packet schema constants, required rows, trust state catalogs, retained
  form catalogs, redaction catalogs, theater reason catalogs, and required
  decision ordering into `internal/packet/contract_catalogs.go`.
- Moved packet/source/projection/row/theater finding/residual gap/decision
  owner data model types into `internal/packet/core_types.go`.
- Moved bundle manifest, bundle entry, resolver entry, and bundle data model
  types into `internal/packet/bundle_types.go`.
- Moved GitHub PR evidence input, prompt boundary, prompt-boundary
  classification, integration action, build result, PR, commit range, check,
  artifact, and review data model types into
  `internal/packet/github_input_types.go`.
- Moved validation result type into `internal/packet/validation_types.go`.
- No validation, GitHub bundle building, prompt-boundary classification
  behavior, rendering, digesting, loading, package boundary, dependency
  direction, or baseline changes were made.
- `cmd/sdp-trace/FAMILY_INDEX.md` was not changed because it indexes command
  families and is not authoritative for `internal/packet` package-local data
  model locality.

Focused packet regression evidence:
- Added `TestPacketContractCatalogsPreserveTrustSurface`.
- Added `TestPacketCoreTypesPreserveJSONShape`.
- Added `TestPacketBundleTypesPreserveJSONShape`.
- Added `TestGitHubEvidenceInputTypesPreserveJSONShape`.
- Strengthened catalog evidence to assert exact catalog membership, not only
  required-value presence.
- Strengthened JSON-shape evidence to assert exact struct `json` tags and
  required-field vs `omitempty` status with reflection.
- Preserved existing serialized shape where an otherwise empty
  `GitHubPREvidenceInput.PromptBoundary` field is emitted as
  `"prompt_boundary":{}` because the parent field is a non-pointer struct;
  its `omitempty` tag does not omit the zero-value nested struct.

Verification:
- verified: `go test ./internal/packet -list 'Test(PacketContractCatalogsPreserveTrustSurface|PacketCoreTypesPreserveJSONShape|PacketBundleTypesPreserveJSONShape|GitHubEvidenceInputTypesPreserveJSONShape)$' | awk '/^Test/ {print}' > /tmp/slice-76-tests.txt && diff -u <(printf 'TestPacketContractCatalogsPreserveTrustSurface\nTestPacketCoreTypesPreserveJSONShape\nTestPacketBundleTypesPreserveJSONShape\nTestGitHubEvidenceInputTypesPreserveJSONShape\n') /tmp/slice-76-tests.txt`
- verified: `go test ./internal/packet -run 'Test(PacketContractCatalogsPreserveTrustSurface|PacketCoreTypesPreserveJSONShape|PacketBundleTypesPreserveJSONShape|GitHubEvidenceInputTypesPreserveJSONShape)$'`
- verified: `go test ./...`
- verified: `go vet ./...`
- verified: `golangci-lint run`
- verified: `go run ./tools/doccheck`
- verified: `go run ./tools/hygienecheck`
- verified: `jq empty schema/*.json`
- verified: `git diff --check`
- verified: `go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools`
- verified: `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`
- verified: `go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`
- verified: `go test -count=1 ./... -coverprofile=coverage.out`
- verified: `go tool cover -func=coverage.out > coverage-func.txt`
- verified: `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`
- verified: `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`
- verified: `find cmd internal tools -type f -name '*_[0-9]*_*.go' | wc -l | tr -d ' '` returned `360`.

Implementation review:
- Lane 1: `LGTM` (`019e9362-68c9-7ed1-8160-f9682047362c`, Einstein)
- Lane 2: `LGTM` after fixes (`019e9362-6c69-7193-9e53-3c6e881cf144`,
  Laplace)
- Lane 3: `LGTM` (`019e9362-7048-7f02-95fd-3811a21795a3`, Epicurus)

Implementation review findings:
- major: catalog regression evidence checked only required-value presence, so
  unexpected catalog values would not fail the test.
- major: JSON-shape evidence checked serialized populated samples but not exact
  struct tags, so adding `omitempty` to required fields could pass.

Fixes after review:
- Added exact map membership assertions for packet catalogs.
- Added reflection-based exact JSON tag assertions for required and optional
  fields across the Slice 76 data model types.
- verified after fix: exact focused test-list command.
- verified after fix: focused packet regression tests.
- verified after fix: `go test ./...`
- verified after fix: `go vet ./...`
- verified after fix: `golangci-lint run`
- verified after fix: `go run ./tools/doccheck`
- verified after fix: `go run ./tools/hygienecheck`
- verified after fix: `jq empty schema/*.json`
- verified after fix: `git diff --check`
- verified after fix: `git diff --cached --check`
- verified after fix: CRAP/MI quality gate bundle.
- verified after fix: Lane 2 returned `LGTM` after residual evidence wording
  fix.
