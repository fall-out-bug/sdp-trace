# Tasks: Source File Locality Cleanup

Status: in_progress

## Active Tasks

- [x] T021-001 Confirm the first cleanup slice is bounded to one command
  family: `cmd/sdp-trace` command-surface registry helpers.
- [x] T021-002 Confirm the slice is behavior-preserving: no command behavior,
  output contract, package boundary, or dependency direction change is planned.
- [x] T021-010 Select exact files for cleanup and record the target grouped file
  in `specs/021-source-file-locality-cleanup/spec.md`.
- [x] T021-020 Move selected helper functions into
  `cmd/sdp-trace/command_surface_registry_helpers.go` and
  `cmd/sdp-trace/command_surface_metadata_helpers.go`.
- [x] T021-030 Run `gofmt` on changed Go files.
- [x] T021-040 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-050 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, and
  `git diff --check`.
- [x] T021-060 Run CRAP and MI quality gates or explicitly record
  `not_assessed` / `cannot_verify` with reason.
- [x] T021-070 Record Slice 1 review and final evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-1-evidence.md`.

## Active Slice 2 Tasks

- [x] T021-080 Confirm Slice 2 is bounded to `cmd/sdp-trace`
  command-surface usage-drift helpers.
- [x] T021-081 Confirm Slice 2 is behavior-preserving: no command behavior,
  output contract, package boundary, dependency direction, or baseline change
  is planned.
- [x] T021-082 Record Slice 2 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-2-plan-review.md`.
- [x] T021-090 Move selected usage collection helpers into
  `cmd/sdp-trace/command_surface_usage_collection.go`.
- [x] T021-091 Move selected usage diff/error helpers into
  `cmd/sdp-trace/command_surface_usage_diff.go`.
- [x] T021-092 Run `gofmt` on changed Go files.
- [x] T021-093 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-094 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-095 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-096 Record Slice 2 review and final evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-2-evidence.md`.

## Active Slice 3 Tasks

- [x] T021-100 Confirm Slice 3 is bounded to `cmd/sdp-trace`
  command-surface list helper shards.
- [x] T021-101 Reject broader argument-helper grouping when pre-change MI
  analysis shows it would require a mixed code/baseline change.
- [x] T021-110 Move selected list helpers into
  `cmd/sdp-trace/command_surface_list_helpers.go`.
- [x] T021-120 Run `gofmt` on changed Go files.
- [x] T021-130 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-140 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-150 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-160 Record Slice 3 review and final evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-3-evidence.md`.
