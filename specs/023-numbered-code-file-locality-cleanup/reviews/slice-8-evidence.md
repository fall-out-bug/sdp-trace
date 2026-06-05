# Slice 8 Evidence: Query Command Locality Cleanup

Status: implemented and pushed; reviewer lanes LGTM; PR checks passed.

## Scope

Slice 8 replaces `cmd/sdp-trace/query_[0-9]*_*.go` shards with
behavior-named query, verify, explain, and query-pack command files.

## File Locality

Removed numbered files: 21.

Added behavior-named files:

- `cmd/sdp-trace/query_verify.go`
- `cmd/sdp-trace/query_verify_args.go`
- `cmd/sdp-trace/query_verify_exit.go`
- `cmd/sdp-trace/query_explain.go`
- `cmd/sdp-trace/query_run.go`
- `cmd/sdp-trace/query_dispatch.go`
- `cmd/sdp-trace/query_pack.go`
- `cmd/sdp-trace/query_pack_build.go`
- `cmd/sdp-trace/query_pack_explain.go`
- `cmd/sdp-trace/query_pack_args.go`
- `cmd/sdp-trace/query_pack_validation.go`

Remaining active numbered Go file count after Slice 8: `1098`.

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

- `cmd/sdp-trace/query_verify.go`: `78.9`
- `cmd/sdp-trace/query_verify_args.go`: `77.3`
- `cmd/sdp-trace/query_verify_exit.go`: `80.7`
- `cmd/sdp-trace/query_explain.go`: `80.2`
- `cmd/sdp-trace/query_run.go`: `80.2`
- `cmd/sdp-trace/query_dispatch.go`: `73.0`
- `cmd/sdp-trace/query_pack.go`: `100.0`
- `cmd/sdp-trace/query_pack_build.go`: `74.0`
- `cmd/sdp-trace/query_pack_explain.go`: `79.8`
- `cmd/sdp-trace/query_pack_args.go`: `70.2`
- `cmd/sdp-trace/query_pack_validation.go`: `72.6`

## Review

Pass:

- `opencode-go/glm-5.1`: `LGTM`
- `opencode-go/qwen3.7-max`: `LGTM`
- `opencode-go/deepseek-v4-flash`: `LGTM`

Fixes from review:

- `opencode-go/qwen3.7-max` reported a minor `FAMILY_INDEX.md` ordering issue
  in the query section. The query entries were reordered alphabetically and all
  lanes were rerun.

## Trust Notes

- Behavior change: not intended; reviewers reported no findings after fix.
- Package boundary change: not intended.
- MI baseline change: none.
- CRAP threshold: verified locally as strict less than 5.
- PR checks: `verify` and `pr-review-evidence-only` passed on PR #73 after
  push.
