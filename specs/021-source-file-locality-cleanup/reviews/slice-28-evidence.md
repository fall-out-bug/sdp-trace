# Slice 28 Evidence: Command Model And Shell Fields

Status: pass

## Scope

- Package: `internal/harnessobs`
- Source shards: `harnessobs_181` through `harnessobs_197`
- Target responsibility groups:
  - command model extraction from argv and `sh`/`bash -c`:
    `command_model_extraction.go`, `command_model_shell.go`
  - command model flag scanning: `command_model_args.go`,
    `command_model_arg.go`, `command_model_next_arg.go`,
    `command_model_prefixed_arg.go`
  - controlled shell field scanning: `shell_fields.go`,
    `shell_field_scan.go`, `shell_field_escape.go`,
    `shell_field_quote.go`, `shell_field_opening_quote.go`,
    `shell_field_unquoted.go`, `shell_field_finish.go`

## Local Verification

- plan review lanes: pass
  - lane 1 requirements/scope: `LGTM` equivalent
  - lane 2 trust/evidence: `LGTM`
  - lane 3 maintainability: `LGTM` after plan fixes
  - non-evidence plan attempts: two opencode-go/glm-5.1 scope lanes are
    `cannot_verify` and not counted
- implementation: pass
- `gofmt -w internal/harnessobs/command_model_extraction.go internal/harnessobs/command_model_shell.go internal/harnessobs/command_model_args.go internal/harnessobs/command_model_arg.go internal/harnessobs/command_model_next_arg.go internal/harnessobs/command_model_prefixed_arg.go internal/harnessobs/shell_fields.go internal/harnessobs/shell_field_scan.go internal/harnessobs/shell_field_escape.go internal/harnessobs/shell_field_quote.go internal/harnessobs/shell_field_opening_quote.go internal/harnessobs/shell_field_unquoted.go internal/harnessobs/shell_field_finish.go`: pass
- `go test ./internal/harnessobs -run 'Test(ExtractCommandModel|ShellFieldsControlledSyntax|SetupSessionWritesSessionRunWithCommand|SetupSessionCommandRejectsModelAndWritesDigest|NormalizedWriteAndShellAndSourceCommitBranches)'`: pass
- `go test ./internal/harnessobs`: pass
- changed-file MI:
  - initial broad scanner/argv grouping failed for `command_model_args.go`,
    `command_model_arg_values.go`, and `shell_field_scanner_stages.go`.
  - `shell_field_unquoted.go` and then `shell_field_quote.go` also failed
    during split refinement; files were split further and baselines were not
    changed.
  - `go run ./tools/qualitycheck -mi-under 70.1 -function-mi-under 70.1 internal/harnessobs/command_model_extraction.go internal/harnessobs/command_model_shell.go internal/harnessobs/command_model_args.go internal/harnessobs/command_model_arg.go internal/harnessobs/command_model_next_arg.go internal/harnessobs/command_model_prefixed_arg.go internal/harnessobs/shell_fields.go internal/harnessobs/shell_field_scan.go internal/harnessobs/shell_field_escape.go internal/harnessobs/shell_field_quote.go internal/harnessobs/shell_field_opening_quote.go internal/harnessobs/shell_field_unquoted.go internal/harnessobs/shell_field_finish.go`: pass
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
- staged pre-review checks: pass
  - `git diff --cached --check`: pass
  - `git diff --cached --name-only | go run ./tools/mibaselinepolicy -base-ref origin/main`: pass
- numbered Slice 28 files remaining: pass; selected files `harnessobs_181`
  through `harnessobs_197` removed
- numbered Go files after Slice 28: `750`

## Reviewer Lanes

- reviewer lane 1 behavior/correctness: `LGTM`; opencode-go/deepseek-v4-pro
  via OpenCode, 2026-06-02, prompt class
  `implementation-review/behavior-correctness`.
- reviewer lane 2 trust/evidence: `LGTM`; opencode-go/mimo-v2.5-pro via
  OpenCode, 2026-06-02, prompt class
  `implementation-review/trust-evidence`.
- reviewer lane 3 maintainability/DX: `LGTM`; opencode-go/qwen3.7-max via
  OpenCode, 2026-06-02, prompt class
  `implementation-review/maintainability-dx`.

## Non-Evidence Implementation Attempts

- An initial opencode-go/deepseek-v4-pro behavior lane returned zero findings
  and `LGTM` with explanatory text. It is not counted for the exact-output
  closure; the re-verdict above returned exactly `LGTM`.
- A kimi-for-coding/k2p6 trust/evidence lane on 2026-06-02 was stopped after a
  prolonged review without a final verdict. It is `cannot_verify` and is not
  counted as reviewer evidence.
- An opencode-go/minimax-m3 maintainability lane on 2026-06-02 exited with an
  unknown certificate verification error before a verdict. It is
  `cannot_verify` and is not counted as reviewer evidence.

## Trust States

- behavior preservation: pass
- argv model extraction: pass
- shell wrapper extraction: pass
- shell-preferred argv precedence: pass
- quoted prompt ignored for model extraction: pass
- unsafe model rejection: pass
- shell quote handling: pass
- shell escape handling: pass
- line continuation handling: pass
- trailing escape preservation: pass
- command model safety scope: pass; `harnessobs_198` onward not changed
- CRAP < 5: pass
- changed-file MI > 70: pass
- repo MI baseline/ratchet gates: pass
- zero numbered files in Slice 28 scope: pass
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
