# Slice 32 Evidence: Correct Serialization Microfile Drift

Status: planned

## Scope

- Package: `internal/harnessobs`
- Corrective target:
  - artifact JSON writing microfiles
  - event reference rendering/safety/validation microfiles
  - digest helper microfiles
- Excluded:
  - command parser
  - session setup
  - session collection
  - raw-event unsafe rule semantics
  - numbered shards `harnessobs_217` onward

## Local Verification

- plan review lanes: pass
  - lane 1 requirements/scope: `LGTM`
  - lane 2 trust/evidence: `LGTM`
  - lane 3 maintainability/DX: `LGTM`
- implementation: pending
- implementation: pass
- focused verification: pass
  - `gofmt -w internal/harnessobs/artifact_json.go internal/harnessobs/event_refs.go internal/harnessobs/event_ref_validation.go internal/harnessobs/digest_helpers.go`: pass
  - `go test ./internal/harnessobs -run 'Test(ObserveValidateCompleteHarnessExport|NormalizeOpenCodeRawLineBytesComputesDigestForEachEvent|SetupSessionCommandRejectsModelAndWritesDigest|CommandModelSafetyAndSourceDigest|CollectSessionWritesObservedRun|ObserveRejectsDigestMismatch|LoadRunRejectsUnsafeEventRefs|ValidateWritesOutPathWhenPasses)'`: pass
  - `go test ./internal/harnessobs`: pass
  - `go run ./tools/qualitycheck -mi-under 70.1 -function-mi-under 70.1 internal/harnessobs/artifact_json.go internal/harnessobs/event_refs.go internal/harnessobs/event_ref_validation.go internal/harnessobs/digest_helpers.go`: pass
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
- microfile drift correction: pass; artifact/event/digest helper files
  consolidated from 13 files to 4 cohesive files without reintroducing numbered
  files:
  - `artifact_json.go` for JSON artifact data rendering and file writing
  - `event_refs.go` for event artifact ref rendering and path safety
  - `event_ref_validation.go` for retained event validation checks
  - `digest_helpers.go` for SHA-256, source, validation, and command digests
  `event_ref_validation.go` remains split from `event_refs.go` because the
  combined event ref rendering/path-safety/validation file failed file-level
  MI, while this split keeps a cohesive validation responsibility instead of
  returning to one-helper microfiles.

## Reviewer Lanes

- reviewer lane 1 behavior/correctness: `LGTM`
- reviewer lane 2 trust/evidence: `LGTM` after final four files and
  MI-driven event validation split were recorded in durable evidence.
- reviewer lane 3 maintainability/DX: `LGTM`

## Trust States

- behavior preservation: pass
- JSON writing format: pass
- event reference rendering: pass
- unsafe event reference rejection: pass
- retained event validation: pass
- raw-line digest canonicalization: pass
- source file digest fallback: pass
- validation digest generation: pass
- command digest generation: pass
- CRAP < 5: pass
- changed-file MI > 70: pass
- repo MI baseline/ratchet gates: pass
- no new numbered files: pass
- microfile drift: pass
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
