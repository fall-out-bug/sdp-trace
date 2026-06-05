# Spec 021 Slice 6 Evidence

Date: 2026-06-01.

## Scope

Implemented Slice 6 for `cmd/sdp-trace` command-surface observe command
metadata shards.

Changed source shape:

- Added `cmd/sdp-trace/command_surface_observe_commands.go`.
- Removed five numbered shards:
  - `cmd/sdp-trace/main_555_commandsurfaceobserveinteraction.go`
  - `cmd/sdp-trace/main_556_commandsurfaceobserveobserve.go`
  - `cmd/sdp-trace/main_557_commandsurfaceobserveharness.go`
  - `cmd/sdp-trace/main_558_commandsurfaceobserveenvelope.go`
  - `cmd/sdp-trace/main_572_commandsurfaceobserve.go`

No command behavior, command metadata value, schema contract, package boundary,
dependency direction, production trust, release approval, external attestation,
or MI-baseline change was intended.

## Decision Notes

Local pre-change metric analysis measured the combined observe metadata file at
file MI `100.0`, so no MI baseline change was needed.

## Verification

Verified pass:

- `gofmt -w cmd/sdp-trace/command_surface_observe_commands.go`
- `go test ./cmd/sdp-trace`
- `go run ./tools/qualitycheck -mi-under 70 cmd/sdp-trace/command_surface_observe_commands.go`
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

Quality notes:

- `command_surface_observe_commands.go` file MI: `100.0`.
- The moved metadata file has no functions, so function MI and CRAP are not
  applicable to the new file.
- No MI baseline file was changed for this slice.

## Targeted Reviews

Verified pass:

- Harness: `opencode run`
- Model: `opencode-go/glm-5.1`
- Result: `LGTM`
- Harness: `opencode run`
- Model: `kimi-for-coding/k2p6`
- Result: `LGTM`
- Harness: `opencode run`
- Model: `opencode-go/deepseek-v4-pro`
- Result: `LGTM`

Review notes:

- Behavior preservation: observe command metadata values and
  `commandSurfaceObserveGroup` list order remain equivalent after the move.
- Scope: no unrelated code, schema, package boundary, dependency, or baseline
  change found.
- Evidence: production, release, external attestation, and merge approval remain
  `not_assessed`.
- Generated artifact discipline: `coverage.out`, `coverage-func.txt`, and
  `gocyclo.txt` were removed before staging.
- `opencode-go/minimax-m3` was not counted because the harness wrote temporary
  files during review despite read-only instructions.

Not assessed:

- No production deployment, release approval, external attestation, or merge
  approval was assessed.
