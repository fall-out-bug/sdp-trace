# Slice 30 Evidence: Unsafe Value Traversal Entrypoints

Status: passed

## Scope

- Package: `internal/harnessobs`
- Source shards: `harnessobs_204` through `harnessobs_208`
- Target responsibility group:
  - unsafe value traversal entrypoints and dispatcher:
    `unsafe_value_traversal.go`
- Excluded:
  - session setup and collection: `harnessobs_209` onward
  - map/slice/string-specific unsafe rule shards: `harnessobs_223` onward

## Local Verification

- plan review lanes: pass
  - lane 1 requirements/scope: `LGTM` after exact re-verdict
  - lane 2 trust/evidence: `LGTM`
  - lane 3 maintainability/DX: `LGTM`
- implementation: pass
- `gofmt -w internal/harnessobs/unsafe_value_traversal.go`: pass
- `go test ./internal/harnessobs -run 'TestFindUnsafe(At|RawEventAt)ReasonCodes'`: pass
- `go test ./internal/harnessobs`: pass
- changed-file MI:
  - `go run ./tools/qualitycheck -mi-under 70.1 -function-mi-under 70.1 internal/harnessobs/unsafe_value_traversal.go`: pass
  - file maintainability index: `73.9`
  - `findUnsafeValueAt` maintainability index: `74.6`
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
  - `git diff --name-only origin/main...HEAD | go run ./tools/mibaselinepolicy -base-ref origin/main`: pass
- numbered Slice 30 files remaining: pass; selected files `harnessobs_204`
  through `harnessobs_208` removed
- numbered Go files after Slice 30: `739`

## Reviewer Lanes

- reviewer lane 1 behavior/correctness: `LGTM`;
  opencode-go/deepseek-v4-flash via OpenCode, 2026-06-02, prompt class
  `implementation-review/behavior-correctness`, verdict-only after
  opencode-go/deepseek-v4-pro inspected the staged diff, reran focused tests,
  and returned a no-finding non-exact response.
- reviewer lane 2 trust/evidence: `LGTM`; opencode-go/mimo-v2.5-pro via
  OpenCode, 2026-06-02, prompt class
  `implementation-review/trust-evidence`, exact re-verdict after an initial
  no-finding non-exact response.
- reviewer lane 3 maintainability/DX: `LGTM`; opencode-go/qwen3.7-max via
  OpenCode, 2026-06-02, prompt class
  `implementation-review/maintainability-dx`.

## Non-Evidence Attempts

- The initial behavior/correctness lane found no issue and reran focused tests,
  but returned explanatory prose before `LGTM`; it is not counted as closure.
- A follow-up behavior/correctness re-verdict inspected the diff and reran the
  same focused tests, but also emitted tool-log output before `LGTM`; it is not
  counted as closure.
- The initial trust/evidence lane found no issue, but returned explanatory
  prose before `LGTM`; it is not counted as closure.

## Trust States

- behavior preservation: pass
- generic unsafe path/reason traversal: pass
- raw-event unsafe path/reason traversal: pass
- map delegation: pass
- slice delegation: pass
- string delegation: pass
- safe-value empty result behavior: pass
- session setup scope: pass; `harnessobs_209` onward not changed
- rule-specific unsafe semantics scope: pass; `harnessobs_223` onward not
  changed
- CRAP < 5: pass
- changed-file MI > 70: pass
- repo MI baseline/ratchet gates: pass
- zero numbered files in Slice 30 scope: pass
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
