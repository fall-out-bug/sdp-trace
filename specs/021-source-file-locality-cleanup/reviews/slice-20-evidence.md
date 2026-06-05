# Slice 20 Evidence: Harnessobs Raw Signals And Key Lookup

Status: pass

## Scope

- Package: `internal/harnessobs`
- Target responsibility groups:
  - raw signal dispatch and structure traversal
  - raw map/slice/string/scalar signal extraction
  - exact and prefix signal matching
  - recursive key presence traversal
  - generic key lookup across maps and slices
  - string and numeric key extraction
  - timestamp string, numeric, and fallback conversion

## Local Verification

- implementation: pass
- `gofmt -w $(git diff --cached --name-only --diff-filter=ACMR -- '*.go')`: pass
- `go test ./internal/harnessobs`: pass
- `go run ./tools/qualitycheck -mi-under 70.1 -function-mi-under 70.1 internal/harnessobs/raw_signal_dispatch.go internal/harnessobs/raw_signal_values.go internal/harnessobs/raw_signal_collections.go internal/harnessobs/raw_signal_value_keys.go internal/harnessobs/signal_exact_matching.go internal/harnessobs/signal_prefix_matching.go internal/harnessobs/key_presence_entry.go internal/harnessobs/key_presence_collections.go internal/harnessobs/key_lookup_entry.go internal/harnessobs/key_lookup_map.go internal/harnessobs/key_lookup_string.go internal/harnessobs/key_lookup_number.go internal/harnessobs/timestamp_lookup.go internal/harnessobs/timestamp_parse.go`: pass
- full repository gates: pass
  - `go test ./...`: pass
  - `go vet ./...`: pass
  - `go run ./tools/doccheck`: pass
  - `go run ./tools/hygienecheck`: pass
  - `jq empty schema/*.json`: pass
  - `git diff --check`: pass
- coverage-backed CRAP and MI baseline gates: pass
  - `go test -count=1 ./... -coverprofile=coverage.out`: pass
  - `go tool cover -func=coverage.out`: pass
  - `go run ./tools/qualitycheck -gocyclo cmd internal tools`: pass
  - `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`: pass
  - `go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`: pass
  - `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`: pass
  - `git diff --name-only origin/main...HEAD | go run ./tools/mibaselinepolicy -base-ref origin/main`: pass
- numbered Slice 20 files remaining: pass; selected files `harnessobs_065`
  through `harnessobs_091` removed
- numbered Go files after Slice 20: `856`
  - `internal/harnessobs`: 269
  - `internal/packet`: 200
  - `cmd/sdp-trace`: 195
  - `internal/prreview`: 192

## Reviewer Lanes

- reviewer lane 1: pass; subagent `019e847d-01dc-7a71-a513-2aef463938f3`, behavior/spec/product drift, `LGTM`
- reviewer lane 2: pass; subagent `019e847d-2678-7630-9ad1-fa3fb6dc5678`, maintainability/navigation, reported major fragmentation because the initial replacement still split signal, key, and timestamp responsibilities across too many sub-20-line helper files. Fixed by regrouping into responsibility files that keep local helpers beside their callers while preserving the MI gate with concise responsibility comments. Re-review subagent `019e8484-833c-7110-a5ef-4bf12d0bdd95` confirmed the structure fix, then reported stale focused-command evidence. Re-review subagent `019e8486-10ab-7160-8c65-52f109dfe617` found the MI evidence fixed and requested a more focused gofmt command; evidence updated to the staged-file command. Final re-review subagent `019e8487-61c9-7081-a3d5-2ad9c73ff339`: `LGTM`.
- reviewer lane 3: pass; subagent `019e847d-43af-7450-bd14-87209f93d848`, trust/evidence/process, `LGTM`

## Trust States

- behavior preservation: pass
- CRAP < 5: pass
- changed-file MI > 70: pass
- repo MI baseline/ratchet gates: pass
- zero numbered files in Slice 20 scope: pass
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
