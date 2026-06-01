# Spec 021 Slice 1 Evidence

Date: 2026-06-01.

## Scope

Implemented Slice 1 for `cmd/sdp-trace` command-surface helpers.

Changed source shape:

- Added `cmd/sdp-trace/command_surface_registry_helpers.go`.
- Added `cmd/sdp-trace/command_surface_metadata_helpers.go`.
- Removed nine numbered helper shards:
  - `cmd/sdp-trace/main_538_commandsurfaceregistrycore.go`
  - `cmd/sdp-trace/main_539_commandsurfaceregistryobserve.go`
  - `cmd/sdp-trace/main_540_commandsurfaceregistryassess.go`
  - `cmd/sdp-trace/main_541_commandsurfaceregistryother.go`
  - `cmd/sdp-trace/main_542_commandsurfaceregistrypacket.go`
  - `cmd/sdp-trace/main_543_commandsurfaceregistry.go`
  - `cmd/sdp-trace/main_580_commandsurfaceflaghelper.go`
  - `cmd/sdp-trace/main_581_commandsurfaceposhelper.go`
  - `cmd/sdp-trace/main_584_commandsurfacereqpos.go`

No command behavior, schema, package boundary, dependency direction, production
trust, release approval, or external attestation change was intended.

## Review

Plan/task review round 1 is recorded in
`specs/021-source-file-locality-cleanup/reviews/plan-task-review-round-1.md`.

Key retained review fixes:

- reverted premature checked confirmation tasks;
- used `ready_for_review` before implementation began;
- split the target into registry and metadata helper files so the slice stays
  cohesive and does not require a mixed code/baseline PR;
- added explicit CRAP/MI coverage and this evidence artifact path.

## Verification

Verified pass:

- `gofmt -w cmd/sdp-trace/command_surface_registry_helpers.go cmd/sdp-trace/command_surface_metadata_helpers.go`
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

Quality notes:

- CRAP strict-less threshold passed after the slice.
- Function MI for every moved function is `100.0`.
- File MI for `command_surface_registry_helpers.go` is `75.4`.
- File MI for `command_surface_metadata_helpers.go` is `100.0`.
- No MI baseline change is needed; CI policy keeps baseline changes separate
  from `cmd`, `internal`, or `tools` Go changes.

Not assessed:

- No production deployment, release approval, external attestation, or merge
  approval was assessed.
