# Slice 29 Evidence: Command Model Safety And Source Digest

Status: passed

## Scope

- Package: `internal/harnessobs`
- Source shards: `harnessobs_198` through `harnessobs_203`
- Target responsibility groups:
  - command model safety: `command_model_safety.go`,
    `command_model_unsafe.go`
  - source file digesting: `source_digest_file.go`
  - source commit discovery: `source_commit.go`, `source_commit_hash.go`

## Local Verification

- plan review lanes: pass
  - lane 1 requirements/scope: `LGTM`
  - lane 2 trust/evidence: `LGTM` after fixing task-checkbox overclaim
  - lane 3 maintainability/DX: `LGTM` after recording local cleanup
    constraints
- implementation: pass
- `gofmt -w internal/harnessobs/command_model_safety.go internal/harnessobs/command_model_unsafe.go internal/harnessobs/source_digest_file.go internal/harnessobs/source_commit.go internal/harnessobs/source_commit_hash.go internal/harnessobs/harnessobs_crap_test.go`: pass
- `go test ./internal/harnessobs -run 'Test(CommandModelSafetyAndSourceDigest|NormalizedWriteAndShellAndSourceCommitBranches|ExtractCommandModel)'`: pass
- `go test ./internal/harnessobs`: pass
- changed-file MI:
  - initial `command_model_safety.go` grouping failed file MI.
  - initial `source_commit.go` with commit-hash validation failed file MI.
  - loop-based `sourceCommitHash` failed function MI.
  - files were split by responsibility, `sourceCommitHash` moved to
    `hex.DecodeString` plus lowercase/length validation, and no baselines were
    changed.
  - `go run ./tools/qualitycheck -mi-under 70.1 -function-mi-under 70.1 internal/harnessobs/command_model_safety.go internal/harnessobs/command_model_unsafe.go internal/harnessobs/source_digest_file.go internal/harnessobs/source_commit.go internal/harnessobs/source_commit_hash.go internal/harnessobs/harnessobs_crap_test.go`: pass
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
- numbered Slice 29 files remaining: pass; selected files `harnessobs_198`
  through `harnessobs_203` removed
- numbered Go files after Slice 29: `744`

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

## Trust States

- behavior preservation: pass
- command model trimming: pass
- unsafe command model character rejection: pass
- unsafe command model path rejection: pass
- overlong command model rejection: pass
- digest file read failure behavior: pass
- digest file SHA-256 behavior: pass
- source commit non-git behavior: pass
- source commit hash shape: pass
- raw-event traversal scope: pass; `harnessobs_204` onward not changed
- CRAP < 5: pass
- changed-file MI > 70: pass
- repo MI baseline/ratchet gates: pass
- zero numbered files in Slice 29 scope: pass
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
