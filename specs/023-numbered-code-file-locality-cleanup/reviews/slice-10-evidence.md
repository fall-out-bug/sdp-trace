# Slice 10 Evidence: Doctor Repo-Observer Locality Cleanup

Status: implemented and pushed; reviewer lanes LGTM; PR checks passed.

## Scope

Slice 10 replaces doctor repo-observer and install command shards with
behavior-named doctor command files. Local doctor report/check shards remain
outside this slice.

## File Locality

Removed numbered files: 13.

Added behavior-named files:

- `cmd/sdp-trace/doctor_command.go`
- `cmd/sdp-trace/doctor_repo_observer.go`
- `cmd/sdp-trace/doctor_install.go`
- `cmd/sdp-trace/doctor_install_args.go`
- `cmd/sdp-trace/doctor_install_options.go`

Remaining active numbered Go file count after Slice 10: `1062`.

Count command:

```text
find cmd internal tools -type f -name '*.go' | rg '(^|/)[A-Za-z]+_[0-9]+_' | sort | wc -l
```

Staged-index cross-check:

```text
git ls-files 'cmd/**/*.go' 'internal/**/*.go' 'tools/**/*.go' | rg '(^|/)[A-Za-z]+_[0-9]+_' | sort | wc -l
```

Both commands returned `1062` in this worktree after staging Slice 10.

## Local Verification

Pass:

- `go test ./cmd/sdp-trace`
- `go test ./...`
- `go vet ./...`
- `go run ./tools/doccheck`
- `go run ./tools/hygienecheck`
- `jq empty schema/*.json`
- `git diff --check`
- `go test -count=1 ./... -coverprofile=coverage.out`
- `go tool cover -func=coverage.out > coverage-func.txt`
- `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`
- `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`
- `go run ./tools/qualitycheck -fail-only -mi-under 70 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools`
- `go run ./tools/qualitycheck -fail-only -function-mi-under 70 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal`
- `git diff --name-only origin/main...HEAD | go run ./tools/mibaselinepolicy -base-ref origin/main`

## File MI

- `cmd/sdp-trace/doctor_command.go`: `73.4`
- `cmd/sdp-trace/doctor_repo_observer.go`: `70.7`
- `cmd/sdp-trace/doctor_install.go`: `71.8`
- `cmd/sdp-trace/doctor_install_args.go`: `74.7`
- `cmd/sdp-trace/doctor_install_options.go`: `78.4`

## Review

Requested OpenCode lanes:

- `opencode-go/glm-5.1`: `cannot_verify`; review command timed out before
  verdict.
- `opencode-go/qwen3.7-max`: `cannot_verify`; review command timed out before
  verdict.
- `opencode-go/deepseek-v4-flash`: `cannot_verify`; first review command
  failed with a certificate error and retry timed out before verdict.

Fallback subagent lanes:

- `Zeno`: `LGTM`
- `Euler`: `LGTM`
- `Nietzsche`: initial major finding on count-command reproducibility; fixed by
  adding the exact count command and staged-index cross-check; re-review
  returned `LGTM`.

## Trust Notes

- Behavior change: not intended; fallback reviewers reported no behavior drift.
- Package boundary change: not intended.
- MI baseline change: none.
- CRAP threshold: verified locally as strict less than 5.
- PR checks: `verify` and `pr-review-evidence-only` passed on PR #73 after
  push.
