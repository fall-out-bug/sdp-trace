# Spec 021 Slice 5 Evidence

Date: 2026-06-01.

## Scope

Implemented Slice 5 for `cmd/sdp-trace` command-surface core command metadata
shards.

Changed source shape:

- Added `cmd/sdp-trace/command_surface_core_commands.go`.
- Removed eight numbered shards:
  - `cmd/sdp-trace/main_548_commandsurfacecorebasic.go`
  - `cmd/sdp-trace/main_549_commandsurfacecorewrap.go`
  - `cmd/sdp-trace/main_550_commandsurfacecorerun.go`
  - `cmd/sdp-trace/main_551_commandsurfacecorepreview.go`
  - `cmd/sdp-trace/main_552_commandsurfacecoredoctor.go`
  - `cmd/sdp-trace/main_553_commandsurfacecoreinstall.go`
  - `cmd/sdp-trace/main_554_commandsurfacecorefixtures.go`
  - `cmd/sdp-trace/main_571_commandsurfacecore.go`

No command behavior, command metadata value, schema contract, package boundary,
dependency direction, production trust, release approval, external attestation,
or MI-baseline change was intended.

## Decision Notes

Local pre-change metric analysis measured the combined core metadata file at
file MI `100.0`, so no MI baseline change was needed.

## Verification

Verified pass:

- `gofmt -w cmd/sdp-trace/command_surface_core_commands.go`
- `go test ./cmd/sdp-trace`
- `go run ./tools/qualitycheck -mi-under 70 cmd/sdp-trace/command_surface_core_commands.go`
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

- `command_surface_core_commands.go` file MI: `100.0`.
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
- Model: `opencode-go/minimax-m3`
- Result: `LGTM`

Review notes:

- Behavior preservation: core command metadata values and `commandSurfaceCore`
  list order remain equivalent after the move.
- Scope: no unrelated code, schema, package boundary, dependency, or baseline
  change found.
- Evidence: production, release, external attestation, and merge approval remain
  `not_assessed`.
- Generated artifact discipline: `coverage.out`, `coverage-func.txt`, and
  `gocyclo.txt` were removed before staging.

Not assessed:

- No production deployment, release approval, external attestation, or merge
  approval was assessed.
