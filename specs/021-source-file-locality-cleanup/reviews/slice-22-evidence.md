# Slice 22 Evidence: Harnessobs Event Scanning And Validation

Status: pass

## Scope

- Package: `internal/harnessobs`
- Target responsibility groups:
  - event source file opening and JSONL scanning
  - scan limits and source hashing
  - line parsing and safe raw-event rejection
  - typed event decoding
  - parsed-event digest validation
  - event identity, reference, and content validation
  - unavailable-field validation

## Local Verification

- implementation: pass
- `gofmt -w internal/harnessobs/event_scan_input.go internal/harnessobs/event_scan_loop.go internal/harnessobs/event_line_parsing.go internal/harnessobs/event_limits.go internal/harnessobs/event_decoding.go internal/harnessobs/event_validation.go internal/harnessobs/event_identity_validation.go internal/harnessobs/event_ref_validation.go internal/harnessobs/event_content_validation.go`: pass
- `go test ./internal/harnessobs`: pass
- `go run ./tools/qualitycheck -mi-under 70.1 -function-mi-under 70.1 internal/harnessobs/event_scan_input.go internal/harnessobs/event_scan_loop.go internal/harnessobs/event_line_parsing.go internal/harnessobs/event_limits.go internal/harnessobs/event_decoding.go internal/harnessobs/event_validation.go internal/harnessobs/event_identity_validation.go internal/harnessobs/event_ref_validation.go internal/harnessobs/event_content_validation.go`: pass
- full repository gates: pass
  - `go test ./...`: pass
  - `go vet ./...`: pass
  - `go run ./tools/doccheck`: pass
  - `go run ./tools/hygienecheck`: pass
  - `jq empty schema/*.json`: pass
  - `git diff --check`: pass
- coverage-backed CRAP and MI baseline gates: pass
  - `go test -count=1 ./... -coverprofile=coverage.out`: pass
  - `go tool cover -func=coverage.out > coverage-func.txt`: pass
  - `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`: pass
  - `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`: pass
  - `go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`: pass
  - `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`: pass
  - `git diff --name-only origin/main...HEAD | go run ./tools/mibaselinepolicy -base-ref origin/main`: pass
- numbered Slice 22 files remaining: pass; selected files `harnessobs_111`
  through `harnessobs_132` removed
- numbered Go files after Slice 22: `815`
  - `internal/harnessobs`: 228
  - `internal/packet`: 200
  - `cmd/sdp-trace`: 195
  - `internal/prreview`: 192

## Reviewer Lanes

- reviewer lane 1: pass; subagent `019e8498-3d75-7862-aaf5-2210090cb4b3`, behavior/diff review, `LGTM`.
- reviewer lane 2: pass after fixes; subagent `019e8498-5bf2-7d03-b4c1-7235584b4e44`, maintainability/navigation, reported one-function microfiles around scanner setup, scan finalization, JSON decode, unsafe rejection, and observed-at validation. Fixed by folding decode/rejection/observed-at helpers into their real responsibility groups and grouping scanner setup with scan finalization in `event_scan_input.go`. Re-review subagent `019e849d-2ef6-73b1-aab1-e8fdf5660227` reported remaining `event_scan_line.go` microfile; fixed by folding it back into `event_scan_loop.go` while preserving MI with `event_scan_input.go`. Re-review subagent `019e849f-a3d8-7ee1-ab6e-bfa38e4763c9` reported remaining `event_source_scanning.go` one-function microfile; fixed by folding `readEvents` into `event_scan_input.go`. Final re-review subagent `019e84a3-1fef-7232-aa9b-32d75d9a26dd` returned `LGTM`.
- reviewer lane 3: pass; subagent `019e8498-7985-79c0-8546-ddb3b00a4e7a`, trust/evidence review, `LGTM`.

## Trust States

- behavior preservation: pass
- CRAP < 5: pass
- changed-file MI > 70: pass
- repo MI baseline/ratchet gates: pass
- zero numbered files in Slice 22 scope: pass
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
