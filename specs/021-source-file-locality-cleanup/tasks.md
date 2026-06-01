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

## Active Slice 4 Tasks

- [x] T021-170 Confirm Slice 4 is bounded to `cmd/sdp-trace`
  command-surface schema and runner shards.
- [x] T021-171 Reject broader metadata/registry grouping when pre-change MI
  analysis shows it would require a mixed code/baseline change.
- [x] T021-180 Move schema type definitions into
  `cmd/sdp-trace/command_surface_schema.go`.
- [x] T021-181 Move runner and JSON writer functions into
  `cmd/sdp-trace/command_surface_runner.go`.
- [x] T021-190 Run `gofmt` on changed Go files.
- [x] T021-200 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-210 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-220 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-230 Run three independent reviewer lanes and record Slice 4 evidence
  in `specs/021-source-file-locality-cleanup/reviews/slice-4-evidence.md`.

## Active Slice 5 Tasks

- [x] T021-240 Confirm Slice 5 is bounded to `cmd/sdp-trace`
  command-surface core command metadata shards.
- [x] T021-241 Confirm local pre-change MI analysis keeps the grouped core file
  above the absolute file-MI threshold.
- [x] T021-250 Move selected core command metadata into
  `cmd/sdp-trace/command_surface_core_commands.go`.
- [x] T021-260 Run `gofmt` on changed Go files.
- [x] T021-270 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-280 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-290 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-300 Run three independent reviewer lanes and record Slice 5 evidence
  in `specs/021-source-file-locality-cleanup/reviews/slice-5-evidence.md`.

## Active Slice 6 Tasks

- [x] T021-310 Confirm Slice 6 is bounded to `cmd/sdp-trace`
  command-surface observe command metadata shards.
- [x] T021-311 Confirm local pre-change MI analysis keeps the grouped observe
  file above the absolute file-MI threshold.
- [x] T021-320 Move selected observe command metadata into
  `cmd/sdp-trace/command_surface_observe_commands.go`.
- [x] T021-330 Run `gofmt` on changed Go files.
- [x] T021-340 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-350 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-360 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-370 Run three independent reviewer lanes and record Slice 6 evidence
  in `specs/021-source-file-locality-cleanup/reviews/slice-6-evidence.md`.

## Active Slice 7 Tasks

- [x] T021-380 Confirm Slice 7 is bounded to `cmd/sdp-trace`
  command-surface assess command metadata shards.
- [x] T021-381 Confirm local pre-change MI analysis keeps the grouped assess
  file above the absolute file-MI threshold.
- [x] T021-390 Move selected assess command metadata into
  `cmd/sdp-trace/command_surface_assess_commands.go`.
- [x] T021-400 Run `gofmt` on changed Go files.
- [x] T021-410 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-420 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-430 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-440 Run three independent reviewer lanes and record Slice 7 evidence
  in `specs/021-source-file-locality-cleanup/reviews/slice-7-evidence.md`.

## Active Slice 8 Tasks

- [x] T021-450 Confirm Slice 8 is bounded to `cmd/sdp-trace`
  command-surface packet command metadata shards.
- [x] T021-451 Confirm local pre-change MI analysis keeps the grouped packet
  file above the absolute file-MI threshold.
- [x] T021-460 Move selected packet command metadata into
  `cmd/sdp-trace/command_surface_packet_commands.go`.
- [x] T021-470 Run `gofmt` on changed Go files.
- [x] T021-480 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-490 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-500 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-510 Run three independent reviewer lanes and record Slice 8 evidence
  in `specs/021-source-file-locality-cleanup/reviews/slice-8-evidence.md`.

## Active Slice 9 Tasks

- [x] T021-520 Confirm Slice 9 is bounded to `cmd/sdp-trace`
  command-surface other command metadata shards.
- [x] T021-521 Confirm local pre-change MI analysis keeps the grouped other
  file above the absolute file-MI threshold.
- [x] T021-530 Move selected other command metadata into
  `cmd/sdp-trace/command_surface_other_commands.go`.
- [x] T021-540 Run `gofmt` on changed Go files.
- [x] T021-550 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-560 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-570 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-580 Run three independent reviewer lanes and record Slice 9 evidence
  in `specs/021-source-file-locality-cleanup/reviews/slice-9-evidence.md`.

## Active Slice 10 Tasks

- [x] T021-590 Confirm Slice 10 is bounded to `cmd/sdp-trace`
  command-surface catalog metadata shards.
- [x] T021-591 Reject broader registry/constants/catalog grouping because local
  pre-change MI analysis measured it below the absolute file-MI threshold.
- [x] T021-600 Move selected catalog metadata into
  `cmd/sdp-trace/command_surface_catalog.go`.
- [x] T021-610 Run `gofmt` on changed Go files.
- [x] T021-620 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-630 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-640 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-650 Run three independent reviewer lanes and record Slice 10 evidence
  in `specs/021-source-file-locality-cleanup/reviews/slice-10-evidence.md`.
