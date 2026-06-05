# Slice 27 Evidence: Harnessobs Serialization And Digests

Status: pass

## Scope

- Package: `internal/harnessobs`
- Source shards: `harnessobs_174` through `harnessobs_180`
- Target responsibility groups:
  - artifact JSON writing: `artifact_json.go`, `artifact_json_data.go`,
    `artifact_json_file.go`
  - event reference rendering and validation: `event_refs.go`,
    `event_ref_path.go`, `event_ref_safety.go`
  - digest helpers: `digest_sum.go`, `source_digest.go`,
    `source_digest_canonical.go`, `validation_digest.go`,
    `command_digest.go`

## Local Verification

- plan review lanes: pass
  - lane 1 requirements/scope: `LGTM` equivalent
  - lane 2 trust/evidence: `LGTM`
  - lane 3 maintainability: `LGTM`
- implementation: pass
- `gofmt -w internal/harnessobs/artifact_json.go internal/harnessobs/artifact_json_data.go internal/harnessobs/artifact_json_file.go internal/harnessobs/event_refs.go internal/harnessobs/event_ref_path.go internal/harnessobs/event_ref_safety.go internal/harnessobs/digest_sum.go internal/harnessobs/source_digest.go internal/harnessobs/source_digest_canonical.go internal/harnessobs/validation_digest.go internal/harnessobs/command_digest.go`: pass
- `go test ./internal/harnessobs -run 'Test(ObserveValidateCompleteHarnessExport|NormalizeOpenCodeRawLineBytesComputesDigestForEachEvent|SetupSessionCommandRejectsModelAndWritesDigest|CollectSessionWritesObservedRun|ObserveRejectsDigestMismatch|LoadRunRejectsUnsafeEventRefs|ValidateWritesOutPathWhenPasses)'`: pass
- `go test ./internal/harnessobs`: pass
- changed-file MI:
  - initial broad grouping command failed for `artifact_json.go`,
    `event_refs.go`, and `digest_helpers.go`; files were split further and
    baselines were not changed.
  - `go run ./tools/qualitycheck -mi-under 70.1 -function-mi-under 70.1 internal/harnessobs/artifact_json.go internal/harnessobs/artifact_json_data.go internal/harnessobs/artifact_json_file.go internal/harnessobs/event_refs.go internal/harnessobs/event_ref_path.go internal/harnessobs/event_ref_safety.go internal/harnessobs/digest_sum.go internal/harnessobs/source_digest.go internal/harnessobs/source_digest_canonical.go internal/harnessobs/validation_digest.go internal/harnessobs/command_digest.go`: pass
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
  - `go tool cover -func=coverage.out > coverage-func.txt`: pass; total
    statement coverage `88.0%`
  - `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`: pass
  - `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`: pass
  - `go run ./tools/qualitycheck -fail-only -cyclo-over 10 -cognitive-over 10 cmd internal tools`: pass
  - `go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`: pass
  - `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`: pass
  - `git diff --cached --name-only | go run ./tools/mibaselinepolicy -base-ref origin/main`: pass
- numbered Slice 27 files remaining: pass; selected files `harnessobs_174`
  through `harnessobs_180` removed
- numbered Go files after Slice 27: `767`

## Reviewer Lanes

- reviewer lane 1 behavior/correctness: `LGTM`; opencode-go/glm-5.1 via
  OpenCode, 2026-06-01, prompt class `implementation-review/behavior`.
- reviewer lane 2 trust/evidence: `LGTM`; kimi-for-coding/k2p6 via
  OpenCode, 2026-06-01, prompt class `implementation-review/trust-evidence`.
- reviewer lane 3 maintainability/DX: `LGTM`; opencode-go/minimax-m3 via
  OpenCode, 2026-06-01, prompt class
  `implementation-review/maintainability-dx`.

## Trust States

- behavior preservation: pass
- JSON writing format: pass
- event reference rendering: pass
- unsafe event reference rejection: pass
- raw-line digest canonicalization: pass
- validation digest generation: pass
- command digest generation: pass
- CRAP < 5: pass
- changed-file MI > 70: pass
- repo MI baseline/ratchet gates: pass
- zero numbered files in Slice 27 scope: pass
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
