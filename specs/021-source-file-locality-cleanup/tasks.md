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

## Active Slice 11 Tasks

- [x] T021-660 Confirm Slice 11 is bounded to the remaining numbered
  `cmd/sdp-trace` CLI helper and command-surface metadata/registry shards.
- [x] T021-661 Reject a single combined file because earlier local MI analysis
  showed broader metadata/registry grouping below the absolute file-MI
  threshold.
- [x] T021-670 Move selected declarations into
  `cmd/sdp-trace/cli_arg_helpers.go`,
  `cmd/sdp-trace/command_surface_metadata.go`, and
  `cmd/sdp-trace/command_surface_registry.go`.
- [x] T021-680 Run `gofmt` on changed Go files.
- [x] T021-690 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-700 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-710 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-720 Run three independent reviewer lanes and record Slice 11 evidence
  in `specs/021-source-file-locality-cleanup/reviews/slice-11-evidence.md`.

## Active Slice 12 Tasks

- [x] T021-730 Confirm Slice 12 is bounded to the remaining numbered
  `cmd/sdp-trace` doctor local report/check shards.
- [x] T021-731 Confirm Slice 12 is behavior-preserving: no command behavior,
  JSON output contract, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-732 Record Slice 12 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-12-plan-review.md`.
- [x] T021-740 Move selected doctor declarations into responsibility-named
  files for report assembly, contract checks, writable paths, expected evidence,
  CI prerequisites, preview metadata, and usage text.
- [x] T021-750 Run `gofmt` on changed Go files.
- [x] T021-760 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-770 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-780 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-790 Run three independent reviewer lanes and record Slice 12
  evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-12-evidence.md`.

## Active Slice 13 Tasks

- [x] T021-800 Confirm Slice 13 is bounded to the remaining numbered
  `cmd/sdp-trace` assess command and preview shards.
- [x] T021-801 Confirm Slice 13 is behavior-preserving: no command behavior,
  JSON output contract, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-802 Record Slice 13 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-13-plan-review.md`.
- [x] T021-810 Move selected assess declarations into responsibility-named
  files for command dispatch, profile runs, input loading, artifact writing,
  preview dispatch, preview report shapes, and profile-specific preview data.
- [x] T021-820 Run `gofmt` on changed Go files.
- [x] T021-830 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-840 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-850 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-860 Run three independent reviewer lanes and record Slice 13
  evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-13-evidence.md`.

## Active Slice 14 Tasks

- [x] T021-870 Confirm Slice 14 is bounded to numbered `cmd/sdp-trace` core CLI
  kernel shards only.
- [x] T021-871 Confirm Slice 14 is behavior-preserving: no command behavior,
  output contract, package boundary, dependency direction, or baseline change
  is planned.
- [x] T021-872 Record Slice 14 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-14-plan-review.md`.
- [x] T021-880 Move selected core declarations into responsibility-named files
  for CLI runtime variables, handler registry, top-level dispatch, subcommand dispatch,
  flag validation, JSON payload writing, and exit-code mapping.
- [x] T021-890 Run `gofmt` on changed Go files.
- [x] T021-900 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-910 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-920 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-930 Run three independent reviewer lanes and record Slice 14
  evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-14-evidence.md`.

## Active Slice 15 Tasks

- [x] T021-940 Confirm Slice 15 is bounded to numbered `cmd/sdp-trace` core
  assessment explain, assessment preview setup, and assessment exit-code shards.
- [x] T021-941 Confirm Slice 15 is behavior-preserving: no command behavior,
  output contract, package boundary, dependency direction, or baseline change
  is planned.
- [x] T021-942 Record Slice 15 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-15-plan-review.md`.
- [x] T021-950 Move selected core declarations into responsibility-named files
  for assessment explanation, typed explanation dispatch, preview input status,
  preview remediation actions, and profile exit-code mapping.
- [x] T021-960 Run `gofmt` on changed Go files.
- [x] T021-970 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-980 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-990 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1000 Run three independent reviewer lanes and record Slice 15
  evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-15-evidence.md`.

## Active Slice 16 Tasks

- [x] T021-1010 Confirm Slice 16 is bounded to the remaining numbered
  `cmd/sdp-trace` core shards: release proof, witness target, export/posture,
  and fixture expectation helpers.
- [x] T021-1011 Confirm Slice 16 is behavior-preserving: no command behavior,
  output contract, package boundary, dependency direction, or baseline change
  is planned.
- [x] T021-1012 Record Slice 16 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-16-plan-review.md`.
- [x] T021-1020 Move selected core declarations into responsibility-named files
  for release proof writing, witness target parsing, export dispatch helpers,
  telemetry export, cross-repo posture export/explain, and fixture expectation
  metadata.
- [x] T021-1030 Run `gofmt` on changed Go files.
- [x] T021-1040 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-1050 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1060 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1070 Run three independent reviewer lanes and record Slice 16
  evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-16-evidence.md`.

## Active Slice 17 Tasks

- [x] T021-1080 Confirm Slice 17 is bounded to numbered
  `internal/harnessobs` type, option, context, scanner, and lookup-map shards
  `harnessobs_011` through `harnessobs_033`.
- [x] T021-1081 Confirm Slice 17 is behavior-preserving: no exported type
  shape, JSON tag, package boundary, dependency direction, or baseline change
  is planned.
- [x] T021-1082 Record Slice 17 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-17-plan-review.md`.
- [x] T021-1090 Move selected harnessobs declarations into
  responsibility-named files for options/context, session model, validation
  lookup sets, event reference checks, existing path specs, shell field
  scanning, and isolation rule installers.
- [x] T021-1100 Run `gofmt` on changed Go files.
- [x] T021-1110 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1120 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1130 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1140 Run three independent reviewer lanes and record Slice 17
  evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-17-evidence.md`.

## Active Slice 18 Tasks

- [x] T021-1150 Confirm Slice 18 is bounded to numbered
  `internal/harnessobs` observe and session setup entrypoint shards
  `harnessobs_034` through `harnessobs_045`.
- [x] T021-1151 Confirm Slice 18 is behavior-preserving: no exported function
  signature, output JSON shape, path-safety behavior, package boundary,
  dependency direction, or baseline change is planned.
- [x] T021-1152 Record Slice 18 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-18-plan-review.md`.
- [x] T021-1160 Move selected observe/session setup declarations into
  responsibility-named files for observe entrypoint, option validation, path
  resolution, source loading, event writing, run construction, observation
  context/time, and session setup entrypoint.
- [x] T021-1170 Run `gofmt` on changed Go files.
- [x] T021-1180 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1190 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1200 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1210 Run three independent reviewer lanes and record Slice 18
  evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-18-evidence.md`.

## Active Slice 19 Tasks

- [x] T021-1220 Confirm Slice 19 is bounded to numbered
  `internal/harnessobs` OpenCode normalization and session command model event
  shards `harnessobs_046` through `harnessobs_064`.
- [x] T021-1221 Confirm Slice 19 is behavior-preserving: no event family
  detection semantics, event IDs, observed-at behavior, actor fallback,
  command-model digesting, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-1222 Record Slice 19 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-19-plan-review.md`.
- [x] T021-1230 Move selected OpenCode/session command declarations into
  responsibility-named files for raw-line normalization, OpenCode event
  construction, family maps/rules/order, observed-at/actor extraction, session
  command facts/events/time, and normalized event construction.
- [x] T021-1240 Run `gofmt` on changed Go files.
- [x] T021-1250 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1260 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1270 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1280 Run three independent reviewer lanes and record Slice 19
  evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-19-evidence.md`.

## Active Slice 20 Tasks

- [x] T021-1290 Confirm Slice 20 is bounded to numbered
  `internal/harnessobs` raw signal traversal, recursive key lookup, and
  timestamp extraction shards `harnessobs_065` through `harnessobs_091`.
- [x] T021-1291 Confirm Slice 20 is behavior-preserving: no raw signal
  semantics, key matching semantics, timestamp fallback behavior, package
  boundary, dependency direction, or baseline change is planned.
- [x] T021-1292 Record Slice 20 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-20-plan-review.md`.
- [x] T021-1300 Move selected raw signal/key/timestamp declarations into
  responsibility-named files for raw signal dispatch, raw map/slice/string/
  scalar extraction, signal matching, key presence traversal, generic key
  lookup, string/numeric lookup, and timestamp parsing.
- [x] T021-1310 Run `gofmt` on changed Go files.
- [x] T021-1320 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1330 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1340 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1350 Run three independent reviewer lanes and record Slice 20
  evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-20-evidence.md`.

## Active Slice 21 Tasks

- [x] T021-1360 Confirm Slice 21 is bounded to numbered
  `internal/harnessobs` token safety, run loading, validation loading/summary,
  and profile validation shards `harnessobs_092` through `harnessobs_110`.
- [x] T021-1361 Confirm Slice 21 is behavior-preserving: no mutation-tool
  classification, safe token rendering, run/event loading, validation summary,
  profile validation, degradation rule behavior, package boundary, dependency
  direction, or baseline change is planned.
- [x] T021-1362 Record Slice 21 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-21-plan-review.md`.
- [x] T021-1370 Move selected token/run/validation declarations into
  responsibility-named files for mutation detection, safe token rendering, run
  loading, validation summary/loading, profile metadata/family validation, and
  degradation rule validation.
- [x] T021-1380 Run `gofmt` on changed Go files.
- [x] T021-1390 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1400 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1410 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1420 Run three independent reviewer lanes and record Slice 21
  evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-21-evidence.md`.

## Active Slice 22 Tasks

- [x] T021-1430 Confirm Slice 22 is bounded to numbered
  `internal/harnessobs` event source scanning, safe decode, and parsed-event
  validation shards `harnessobs_111` through `harnessobs_132`.
- [x] T021-1431 Confirm Slice 22 is behavior-preserving: no JSONL scanning,
  source hashing, limit handling, safe event rejection, typed event decoding,
  digest validation, event identity/ref/content/unavailable-field behavior,
  package boundary, dependency direction, or baseline change is planned.
- [x] T021-1432 Record Slice 22 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-22-plan-review.md`.
- [x] T021-1440 Move selected event scanning and validation declarations into
  responsibility-named files for source scanning, line parsing, safe decoding,
  parsed-event digest validation, identity/ref/content validation, and
  unavailable-field validation.
- [x] T021-1450 Run `gofmt` on changed Go files.
- [x] T021-1460 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1470 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1480 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1490 Run three independent reviewer lanes and record Slice 22
  evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-22-evidence.md`.

## Active Slice 23 Tasks

- [x] T021-1500 Confirm Slice 23 is bounded to numbered
  `internal/harnessobs` evaluation and dimension-composition shards
  `harnessobs_133` through `harnessobs_142`.
- [x] T021-1501 Confirm Slice 23 is behavior-preserving: no validation state,
  schema-version mismatch, event-family count, dimension ordering, state
  ranking, validation digest, package boundary, dependency direction, or
  baseline change is planned.
- [x] T021-1502 Record Slice 23 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-23-plan-review.md`.
- [x] T021-1510 Move selected evaluation declarations into
  responsibility-named files for evaluation dispatch, validation assembly,
  dimension collection, dimension ordering, and dimension state composition.
- [x] T021-1520 Run `gofmt` on changed Go files.
- [x] T021-1530 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1540 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1550 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1560 Run three independent reviewer lanes and record Slice 23
  evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-23-evidence.md`.

## Active Slice 24 Tasks

- [x] T021-1570 Confirm Slice 24 is bounded to numbered
  `internal/harnessobs` existing path-safety shards `harnessobs_143` through
  `harnessobs_149`.
- [x] T021-1571 Confirm Slice 24 is behavior-preserving: no traversal
  rejection, symlink resolution, working-directory containment, expected
  file/directory type behavior, package boundary, dependency direction, or
  baseline change is planned.
- [x] T021-1572 Record Slice 24 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-24-plan-review.md`.
- [x] T021-1580 Move selected existing path-safety declarations into
  responsibility-named files for public existing path checks and existing path
  resolution/type validation.
- [x] T021-1590 Run `gofmt` on changed Go files.
- [x] T021-1600 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1601 Run focused path-safety regression evidence covering traversal
  rejection, URL-like rejection, symlink resolution, working-directory
  containment, and file/directory type checks, or mark any unverified behavior
  `not_assessed` with reason.
- [x] T021-1610 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1620 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1630 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  24 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-24-evidence.md`.

## Active Slice 25 Tasks

- [x] T021-1640 Confirm Slice 25 is bounded to numbered
  `internal/harnessobs` output file and parent path-safety shards
  `harnessobs_150` through `harnessobs_161`.
- [x] T021-1641 Confirm Slice 25 is behavior-preserving: no output traversal
  rejection, output basename validation, missing-parent resolution, symlink
  containment, working-directory relative conversion, package boundary,
  dependency direction, or baseline change is planned.
- [x] T021-1642 Record Slice 25 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-25-plan-review.md`.
- [x] T021-1650 Move selected output file and parent path declarations into
  responsibility-named files for output file validation and parent path
  resolution/containment.
- [x] T021-1660 Run `gofmt` on changed Go files.
- [x] T021-1670 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1671 Run focused output path regression evidence covering output
  traversal rejection, basename validation, missing-parent resolution, symlink
  resolution, working-directory containment, and successful validation output
  writing, or mark any unverified behavior `not_assessed` with reason.
- [x] T021-1680 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1690 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1700 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  25 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-25-evidence.md`.

## Active Slice 26 Tasks

- [x] T021-1710 Confirm Slice 26 is bounded to numbered
  `internal/harnessobs` output directory safety and emptiness shards
  `harnessobs_162` through `harnessobs_173`.
- [x] T021-1711 Confirm Slice 26 is behavior-preserving: no output directory
  traversal rejection, existing output symlink containment, missing parent
  containment, working-directory escape detection, empty/non-empty outdir
  behavior, package boundary, dependency direction, or baseline change is
  planned.
- [x] T021-1712 Record Slice 26 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-26-plan-review.md`.
- [x] T021-1720 Move selected output directory declarations into
  responsibility-named files for output directory validation, symlink/parent
  containment, and empty-or-missing directory checks.
- [x] T021-1730 Run `gofmt` on changed Go files.
- [x] T021-1740 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1741 Run focused output directory regression evidence covering
  traversal rejection, existing symlink escape, parent symlink escape,
  path-escape classification, missing/empty directory acceptance, and non-empty
  directory rejection, or mark any unverified behavior `not_assessed` with
  reason.
- [x] T021-1750 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1760 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1770 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  26 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-26-evidence.md`.

## Active Slice 27 Tasks

- [x] T021-1780 Confirm Slice 27 is bounded to numbered
  `internal/harnessobs` artifact serialization, event reference, and digest
  helper shards `harnessobs_174` through `harnessobs_180`.
- [x] T021-1781 Confirm Slice 27 is behavior-preserving: no JSON formatting,
  event reference validation, digest canonicalization, command digest,
  package boundary, dependency direction, or baseline change is planned.
- [x] T021-1782 Record Slice 27 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-27-plan-review.md`.
- [x] T021-1790 Move selected artifact serialization, event reference, and
  digest declarations into responsibility-named files.
- [x] T021-1800 Run `gofmt` on changed Go files.
- [x] T021-1810 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1811 Run focused serialization/reference/digest regression evidence
  covering JSON writing, safe/unsafe event refs, raw-line source digest
  canonicalization, validation digest generation, and command digest generation,
  or mark any unverified behavior `not_assessed` with reason.
- [x] T021-1820 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1830 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1840 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  27 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-27-evidence.md`.
