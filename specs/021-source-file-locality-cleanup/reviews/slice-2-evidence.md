# Spec 021 Slice 2 Evidence

Date: 2026-06-01.

## Scope

Implemented Slice 2 for `cmd/sdp-trace` command-surface usage-drift helpers.

Changed source shape:

- Added `cmd/sdp-trace/command_surface_usage_collection.go`.
- Added `cmd/sdp-trace/command_surface_usage_diff.go`.
- Removed eight numbered usage-drift helper shards:
  - `cmd/sdp-trace/main_576_commandsurfacedriftadd.go`
  - `cmd/sdp-trace/main_585_commandsurfacedrift.go`
  - `cmd/sdp-trace/main_586_commandsurfaceregistryusages.go`
  - `cmd/sdp-trace/main_587_commandsurfacehelpusages.go`
  - `cmd/sdp-trace/main_588_commandsurfacediffsets.go`
  - `cmd/sdp-trace/main_589_commandsurfacesorteddiffs.go`
  - `cmd/sdp-trace/main_590_commandsurfacedrifterror.go`
  - `cmd/sdp-trace/main_591_commandsurfacedriftparts.go`

No command behavior, schema, package boundary, dependency direction, production
trust, release approval, external attestation, or MI-baseline change was
intended.

## Decision Notes

Rejected alternative: one combined command-surface drift file. Local metric
analysis measured the combined file below the absolute file-MI threshold, so
Slice 2 was split into collection and diff files.

## Verification

Verified pass:

- `gofmt -w cmd/sdp-trace/command_surface_usage_collection.go cmd/sdp-trace/command_surface_usage_diff.go`
- `go test ./cmd/sdp-trace`
- `go run ./tools/qualitycheck -mi-under 70 cmd/sdp-trace/command_surface_usage_collection.go cmd/sdp-trace/command_surface_usage_diff.go`
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

- `command_surface_usage_collection.go` file MI: `74.3`.
- `command_surface_usage_diff.go` file MI: `71.9`.
- Moved functions have function MI `100.0`.
- CRAP strict-less threshold passed.
- No MI baseline file was changed for this slice.

## Targeted Review

Verified pass:

- Harness: `opencode run`
- Model: `opencode-go/glm-5.1`
- Result: `LGTM`

Review notes:

- Bounded scope: eight numbered usage-drift shards removed, two behavior-named
  files added.
- Behavior preservation: moved function bodies and merged imports reviewed as
  equivalent.
- MI baseline discipline: no `tools/qualitycheck/*baseline*.json` changes in
  this slice.
- Drift: no spec, constitution, product, or repo-wide rename drift found.

Not assessed:

- No production deployment, release approval, external attestation, or merge
  approval was assessed.
