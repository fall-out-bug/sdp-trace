# Slice 9 Evidence: Witness Command Locality Cleanup

Status: implemented and pushed; reviewer lanes LGTM; PR checks passed.

## Scope

Slice 9 replaces `cmd/sdp-trace/witness_[0-9]*_*.go` shards with
behavior-named witness command files.

## File Locality

Removed numbered files: 23.

Added behavior-named files:

- `cmd/sdp-trace/witness_command.go`
- `cmd/sdp-trace/witness_options.go`
- `cmd/sdp-trace/witness_options_parse.go`
- `cmd/sdp-trace/witness_options_build.go`
- `cmd/sdp-trace/witness_flag_set.go`
- `cmd/sdp-trace/witness_required_fields.go`
- `cmd/sdp-trace/witness_kind_validation.go`
- `cmd/sdp-trace/witness_record_builders.go`
- `cmd/sdp-trace/witness_output.go`
- `cmd/sdp-trace/witness_customer_pki.go`

Remaining active numbered Go file count after Slice 9: `1075`.

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

- `cmd/sdp-trace/witness_command.go`: `80.0`
- `cmd/sdp-trace/witness_options.go`: `100.0`
- `cmd/sdp-trace/witness_options_parse.go`: `75.3`
- `cmd/sdp-trace/witness_options_build.go`: `86.2`
- `cmd/sdp-trace/witness_flag_set.go`: `100.0`
- `cmd/sdp-trace/witness_required_fields.go`: `80.5`
- `cmd/sdp-trace/witness_kind_validation.go`: `70.3`
- `cmd/sdp-trace/witness_record_builders.go`: `70.9`
- `cmd/sdp-trace/witness_output.go`: `81.9`
- `cmd/sdp-trace/witness_customer_pki.go`: `75.3`

## Review

Pass:

- `opencode-go/glm-5.1`: `LGTM`
- `opencode-go/qwen3.7-max`: `LGTM`
- `opencode-go/deepseek-v4-flash`: `LGTM`

## Trust Notes

- Behavior change: not intended; reviewers reported no findings.
- Package boundary change: not intended.
- MI baseline change: none.
- CRAP threshold: verified locally as strict less than 5.
- PR checks: `verify` and `pr-review-evidence-only` passed on PR #73 after
  push.
