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

## Active Slice 28 Tasks

- [x] T021-1850 Confirm Slice 28 is bounded to numbered
  `internal/harnessobs` command model extraction and controlled shell field
  parser shards `harnessobs_181` through `harnessobs_197`.
- [x] T021-1851 Confirm Slice 28 is behavior-preserving: no command model
  extraction precedence, shell wrapper detection, flag parsing, quote/escape
  handling, field separator handling, command model safety, package boundary,
  dependency direction, or baseline change is planned.
- [x] T021-1852 Record Slice 28 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-28-plan-review.md`.
- [x] T021-1860 Move selected command model extraction and shell field parser
  declarations into responsibility-named files.
- [x] T021-1870 Run `gofmt` on changed Go files.
- [x] T021-1880 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1881 Run focused command model and shell field parser regression
  evidence covering argv flags, `sh`/`bash -c` extraction, shell-preferred argv,
  quoted prompt ignoring, unsafe model rejection, quoting, escaping, line
  continuation, and trailing escape behavior, or mark any unverified behavior
  `not_assessed` with reason.
- [x] T021-1890 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1900 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1910 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  28 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-28-evidence.md`.

## Active Slice 29 Tasks

- [x] T021-1920 Confirm Slice 29 is bounded to numbered
  `internal/harnessobs` command model safety and source-bound digest helper
  shards `harnessobs_198` through `harnessobs_203`.
- [x] T021-1921 Confirm Slice 29 is behavior-preserving: no command model
  trimming/rejection behavior, digest algorithm, read-error fallback,
  source-commit lookup, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-1922 Record Slice 29 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-29-plan-review.md`.
- [x] T021-1930 Move selected command model safety and source digest
  declarations into responsibility-named files.
- [x] T021-1940 Run `gofmt` on changed Go files.
- [x] T021-1950 Run focused Go verification for `internal/harnessobs`.
- [x] T021-1951 Run focused command model safety and source digest regression
  evidence covering whitespace trimming, unsafe URL/quote/escape/whitespace
  rejection, traversal/absolute path rejection, overlong model rejection,
  digest read failure, digest SHA-256 output, non-git source commit fallback,
  and git source commit hash shape, or mark any unverified behavior
  `not_assessed` with reason.
- [x] T021-1960 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-1970 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-1980 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  29 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-29-evidence.md`.

## Active Slice 30 Tasks

- [x] T021-1990 Confirm Slice 30 is bounded to numbered
  `internal/harnessobs` unsafe value traversal entrypoint and dispatcher shards
  `harnessobs_204` through `harnessobs_208`.
- [x] T021-1991 Confirm Slice 30 is behavior-preserving: no unsafe field path
  rendering, reason-code semantics, raw-event mode behavior, map/slice/string
  delegation, package boundary, dependency direction, or baseline change is
  planned.
- [x] T021-1992 Record Slice 30 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-30-plan-review.md`.
- [x] T021-2000 Move selected unsafe value traversal entrypoint and dispatcher
  declarations into a responsibility-named file.
- [x] T021-2010 Run `gofmt` on changed Go files.
- [x] T021-2020 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2021 Run focused unsafe traversal regression evidence covering
  generic and raw-event unsafe path/reason results, traversal through map,
  slice, and string values, and safe-value empty results, or mark any
  unverified behavior `not_assessed` with reason.
- [x] T021-2030 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2040 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2050 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  30 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-30-evidence.md`.

## Active Slice 31 Tasks

- [x] T021-2060 Confirm Slice 31 is bounded to numbered
  `internal/harnessobs` session setup path, run construction, command metadata,
  time fallback, and session JSON writing shards `harnessobs_209` through
  `harnessobs_216`.
- [x] T021-2061 Confirm Slice 31 is behavior-preserving: no required option
  error text, path-safety behavior, output directory creation, isolation rule
  installation, command digest/model state, session time fallback, JSON output,
  package boundary, dependency direction, or baseline change is planned.
- [x] T021-2062 Record Slice 31 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-31-plan-review.md`.
- [x] T021-2070 Move selected session setup declarations into
  responsibility-named files.
- [x] T021-2080 Run `gofmt` on changed Go files.
- [x] T021-2090 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2091 Run focused session setup regression evidence covering
  required option errors, unsafe path rejection, invalid profile rejection,
  command digest/model state, blank command defaults, session JSON writing, and
  isolation rule installation, or mark any unverified behavior `not_assessed`
  with reason.
- [x] T021-2100 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2110 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2120 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  31 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-31-evidence.md`.

## Active Slice 32 Tasks

- [x] T021-2130 Confirm Slice 32 is bounded to corrective consolidation of
  `internal/harnessobs` artifact JSON, event reference, and digest helper
  microfiles introduced around Slice 27.
- [x] T021-2131 Confirm Slice 32 is behavior-preserving: no JSON formatting,
  event ref rendering/safety, digest algorithm, source digest canonicalization,
  validation digest, command digest, package boundary, dependency direction, or
  baseline change is planned.
- [x] T021-2132 Record Slice 32 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-32-plan-review.md`.
- [x] T021-2140 Consolidate artifact JSON, event reference, and digest helper
  microfiles into four cohesive responsibility-named groups without
  reintroducing numbered files; keep retained event validation split from event
  ref rendering/path safety because the combined file failed file-level MI.
- [x] T021-2150 Run `gofmt` on changed Go files.
- [x] T021-2160 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2161 Run focused serialization/reference/digest regression evidence
  covering JSON output writing, unsafe event ref rejection, event ref rendering,
  raw-line source digest canonicalization, digest mismatch rejection,
  validation output writing, command digest/model state, and source digest
  fallback behavior, or mark any unverified behavior `not_assessed` with
  reason.
- [x] T021-2170 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2180 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2190 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  32 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-32-evidence.md`.

## Active Slice 33 Tasks

- [x] T021-2200 Confirm Slice 33 is bounded to numbered
  `internal/harnessobs` session collection entrypoint and input-loading shards
  `harnessobs_217` through `harnessobs_222`.
- [x] T021-2201 Confirm Slice 33 is behavior-preserving: no collect option
  validation, profile/session loading, profile mismatch rejection, harness
  profile loading, session source unavailable fallback, source collection
  dispatch, time fallback, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-2202 Record Slice 33 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-33-plan-review.md`.
- [x] T021-2210 Move selected session collection entrypoint, context, time, and
  input-loading declarations into cohesive responsibility-named files.
- [x] T021-2220 Run `gofmt` on changed Go files.
- [x] T021-2230 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2231 Run focused session collection regression evidence covering
  collection entrypoint dispatch to source-unavailable fallback, collect option
  validation handoff, profile/session input loading, profile mismatch rejection,
  harness profile loading handoff, and collection time fallback, or mark any
  unverified behavior `not_assessed` with reason.
- [x] T021-2240 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2250 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2260 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  33 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-33-evidence.md`.

## Active Slice 34 Tasks

- [x] T021-2270 Confirm Slice 34 is bounded to numbered
  `internal/harnessobs` raw unsafe traversal rule shards `harnessobs_223`
  through `harnessobs_244`.
- [x] T021-2271 Confirm Slice 34 is behavior-preserving: no unsafe field path
  rendering, raw-event skip semantics, forbidden/sensitive field reason codes,
  string path/token/url classification, digest-field exemptions, raw path-like
  exemptions, package boundary, dependency direction, or baseline change is
  planned.
- [x] T021-2272 Record Slice 34 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-34-plan-review.md`.
- [x] T021-2280 Move selected map/slice/string traversal and safety-rule
  declarations into cohesive responsibility-named files.
- [x] T021-2290 Run `gofmt` on changed Go files.
- [x] T021-2300 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2301 Run focused unsafe-rule regression evidence covering generic
  unsafe traversal, raw-event unsafe reason codes, map child path rendering,
  skippable unretained raw prompt/body fields, path/private-path rejection,
  authenticated URL detection, token/base64 detection and exemptions,
  digest-field detection, and raw path-like field exemptions, or mark any
  unverified behavior `not_assessed` with reason.
- [x] T021-2310 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2320 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2330 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  34 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-34-evidence.md`.

## Active Slice 35 Tasks

- [x] T021-2340 Confirm Slice 35 is bounded to numbered
  `internal/harnessobs` validation enum, safe reference, non-authority, and
  validation decoding helper shards `harnessobs_245` through `harnessobs_255`.
- [x] T021-2341 Confirm Slice 35 is behavior-preserving: no accepted enum
  values, reference safety semantics, operation reference prefix rules,
  digest-mismatch event redaction, non-authority boundary text, validation JSON
  decoding, package boundary, dependency direction, or baseline change is
  planned.
- [x] T021-2342 Record Slice 35 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-35-plan-review.md`.
- [x] T021-2350 Move selected validation lookup, safe reference, non-authority,
  and validation decoding declarations into cohesive responsibility-named files.
- [x] T021-2360 Run `gofmt` on changed Go files.
- [x] T021-2370 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2371 Run focused validation helper regression evidence covering
  profile family/state/rule-key validation, event content/reference validation,
  unsafe operation reference rejection, digest-mismatch event fallback,
  validation JSON decoding, and non-authority boundary rendering, or mark any
  unverified behavior `not_assessed` with reason.
- [x] T021-2380 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2390 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2400 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  35 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-35-evidence.md`.

## Active Slice 36 Tasks

- [x] T021-2410 Confirm Slice 36 is bounded to numbered
  `internal/harnessobs` session collection validation, harness profile/source
  resolution, observed-run materialization, and runtime session command shards
  `harnessobs_256` through `harnessobs_280`.
- [x] T021-2411 Confirm Slice 36 is behavior-preserving: no collect option
  errors, path safety, raw-event normalization dispatch, source-unavailable
  fallback semantics, observed output shape, session collection state/reason,
  command execution IO behavior, process metadata, wait-error propagation,
  package boundary, dependency direction, or baseline change is planned.
- [x] T021-2412 Record Slice 36 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-36-plan-review.md`.
- [x] T021-2420 Move selected collect validation, source resolution, observed
  output, unavailable fallback, and runtime command declarations into cohesive
  responsibility-named files.
- [x] T021-2430 Run `gofmt` on changed Go files.
- [x] T021-2440 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2441 Run focused session collection/runtime regression evidence
  covering collect option validation, unsafe path rejection, raw event source
  normalization, source-unavailable fallback, observed run/event output,
  collected session finalization, command-required errors, command digest/model
  metadata, process metadata, and wait-error propagation, or mark any unverified
  behavior `not_assessed` with reason.
- [x] T021-2450 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2460 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2470 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  36 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-36-evidence.md`.

## Active Slice 37 Tasks

- [x] T021-2480 Confirm Slice 37 is bounded to numbered
  `internal/harnessobs` validation command, input requirement, safe path
  resolution, run-evaluation fallback, and optional output shards
  `harnessobs_281` through `harnessobs_290`.
- [x] T021-2481 Confirm Slice 37 is behavior-preserving: no required option
  messages, whitespace handling, profile/run/out path safety, out-path optional
  behavior, `LoadRun` fallback semantics, validation digest, non-authority
  boundary, JSON output shape, package boundary, dependency direction, or
  baseline change is planned.
- [x] T021-2482 Record Slice 37 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-37-plan-review.md`.
- [x] T021-2490 Move selected validation command, input resolution, fallback
  evaluation, non-blank, and optional output declarations into cohesive
  responsibility-named files.
- [x] T021-2500 Run `gofmt` on changed Go files.
- [x] T021-2510 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2511 Run focused validation command regression evidence covering
  required option errors, whitespace-only option rejection, unsafe profile/run
  path rejection, unsafe output path rejection, optional output omission,
  validation output writing, source-unavailable cannot-verify fallback,
  validation digest, non-authority boundary, and normal validation from loaded
  run, or mark any
  unverified behavior `not_assessed` with reason.
- [x] T021-2520 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2530 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2540 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  37 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-37-evidence.md`.

## Active Slice 38 Tasks

- [x] T021-2550 Confirm Slice 38 is bounded to numbered
  `internal/harnessobs` profile loading, session profile identity/path
  validation, required session path, and raw-event config shards
  `harnessobs_291` through `harnessobs_299`.
- [x] T021-2551 Confirm Slice 38 is behavior-preserving: no strict JSON load
  behavior, profile validation handoff, session profile schema/id errors,
  stream capture normalization handoff, setup action/isolation validation
  handoff, required session path whitespace handling, raw-event format/source
  pairing errors, raw-event format support, package boundary, dependency
  direction, or baseline change is planned.
- [x] T021-2552 Record Slice 38 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-38-plan-review.md`.
- [x] T021-2560 Move selected profile loading and session/raw-event config
  validation declarations into cohesive responsibility-named files.
- [x] T021-2570 Run `gofmt` on changed Go files.
- [x] T021-2580 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2581 Run focused profile/session validation regression evidence
  covering profile load validation handoff, strict session profile JSON
  rejection, session profile schema/id rejection, default stream capture
  normalization, setup action validation handoff, required session path
  whitespace rejection, raw-event source/format pair errors, unsupported
  raw-event format rejection, and valid raw-event config acceptance, or mark any
  unverified behavior `not_assessed` with reason.
- [x] T021-2590 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2600 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2610 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  38 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-38-evidence.md`.

## Active Slice 39 Tasks

- [x] T021-2620 Confirm Slice 39 is bounded to numbered
  `internal/harnessobs` stream capture, setup action validation, and isolation
  rule validation shards `harnessobs_300` through `harnessobs_309`.
- [x] T021-2621 Confirm Slice 39 is behavior-preserving: no stream capture
  defaulting or unsupported-mode errors, setup action max/id/kind validation,
  isolation rule list iteration, isolation id/pattern/target/kind validation,
  unsafe pattern semantics, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-2622 Record Slice 39 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-39-plan-review.md`.
- [x] T021-2630 Move selected stream capture, setup action, and isolation rule
  validation declarations into cohesive responsibility-named files.
- [x] T021-2640 Run `gofmt` on changed Go files.
- [x] T021-2650 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2651 Run focused session rule validation regression evidence
  covering default stream capture, unsupported stream capture modes,
  too many setup actions, unsafe setup action IDs, unsupported setup action
  kinds, valid setup action kinds, isolation rule list iteration, unsafe
  isolation IDs, unsafe/blank/newline isolation patterns, unsafe target paths,
  unsupported isolation kinds, and valid isolation rules, or mark any
  unverified behavior `not_assessed` with reason.
- [x] T021-2660 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2670 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2680 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  39 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-39-evidence.md`.

## Active Slice 40 Tasks

- [x] T021-2690 Confirm Slice 40 is bounded to numbered
  `internal/harnessobs` isolation target resolution, line/JSON rule
  materialization, readback verification, and isolation result digest shards
  `harnessobs_310` through `harnessobs_334`.
- [x] T021-2691 Confirm Slice 40 is behavior-preserving: no changes to
  profile-relative isolation path safety, parent/filename validation, installer
  dispatch, line rule idempotence, optional line/object loading, blank JSON
  handling, JSON read-deny mutation, readback cannot-verify semantics, digest
  assignment, package boundary, dependency direction, or baseline change is
  planned.
- [x] T021-2692 Record Slice 40 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-40-plan-review.md`.
- [x] T021-2700 Move selected isolation installation and readback declarations
  into cohesive responsibility-named files.
- [x] T021-2710 Run `gofmt` on changed Go files.
- [x] T021-2720 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2721 Run focused isolation installation/readback regression evidence
  covering profile-relative isolation path traversal, clean profile-relative
  path joining, parent containment, unsafe filenames, unsupported installer
  kind, line rule idempotence, optional missing/blank line reads, line writing,
  line-read error propagation, JSON read-deny creation, missing/blank/invalid
  JSON object loading, nested object preservation, readback pass/cannot-verify,
  unsupported readback kind, digest assignment, and full setup-session isolation
  installation, or mark any unverified behavior `not_assessed` with reason.
- [x] T021-2730 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2740 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2750 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  40 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-40-evidence.md`.

## Active Slice 41 Tasks

- [x] T021-2760 Confirm Slice 41 is bounded to numbered
  `internal/harnessobs` loaded session run loading/validation and shared
  existing JSON reader shards `harnessobs_335` through `harnessobs_339`.
- [x] T021-2761 Confirm Slice 41 is behavior-preserving: no changes to safe
  existing file validation, permissive JSON unmarshal semantics, strict JSON
  unknown-field rejection, strict trailing-data rejection, session run schema
  validation, profile ID safety, package boundary, dependency direction, or
  baseline change is planned.
- [x] T021-2762 Record Slice 41 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-41-plan-review.md`.
- [x] T021-2770 Move selected loaded session run and existing JSON reader
  declarations into cohesive responsibility-named files.
- [x] T021-2780 Run `gofmt` on changed Go files.
- [x] T021-2790 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2791 Run focused session loading and JSON reader regression
  evidence covering successful loaded session run parsing, unsupported session
  schema versions, unsafe loaded session profile IDs, unsafe/missing existing
  files, permissive JSON unmarshal behavior, strict JSON unknown-field
  rejection, strict JSON trailing-data rejection, and strict malformed JSON
  rejection, or mark any unverified behavior `not_assessed` with reason.
- [x] T021-2800 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2810 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2820 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  41 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-41-evidence.md`.

## Active Slice 42 Tasks

- [x] T021-2830 Confirm Slice 42 is bounded to numbered
  `internal/harnessobs` session run construction shards `harnessobs_340`
  through `harnessobs_342`.
- [x] T021-2831 Confirm Slice 42 is behavior-preserving: no changes to session
  run schema/profile/path/raw-event field copying, sorted setup action IDs,
  command/process/collection default states, source commit state pass-through,
  collection reason, timestamp formatting, package boundary, dependency
  direction, or baseline change is planned. Source commit discovery/proof is
  `not_assessed` in this slice because `harnessobs_343` remains excluded.
- [x] T021-2832 Record Slice 42 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-42-plan-review.md`.
- [x] T021-2840 Move selected session run construction declarations into a
  cohesive responsibility-named file.
- [x] T021-2850 Run `gofmt` on changed Go files.
- [x] T021-2860 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2861 Run focused session run construction regression evidence
  covering setup action ID sorting, schema/profile/path/raw-event field
  copying, cannot-verify defaults for command/process/collection, source commit
  state pass-through from the excluded source helper, collection reason,
  timestamp formatting, and full setup-session saved run behavior. Source
  commit discovery/proof remains `not_assessed` for this slice, or mark any
  other unverified behavior `not_assessed` with reason.
- [x] T021-2870 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2880 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2890 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  42 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-42-evidence.md`.

## Active Slice 43 Tasks

- [x] T021-2900 Confirm Slice 43 is bounded to numbered
  `internal/harnessobs` source commit state mapping shard `harnessobs_343`.
- [x] T021-2901 Confirm Slice 43 is behavior-preserving: no changes to git
  command execution, source commit hash validation, empty commit handling,
  pass/cannot_verify state mapping, package boundary, dependency direction, or
  baseline change is planned.
- [x] T021-2902 Record Slice 43 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-43-plan-review.md`.
- [x] T021-2910 Move source commit state mapping into the cohesive existing
  source commit file.
- [x] T021-2920 Run `gofmt` on changed Go files.
- [x] T021-2930 Run focused Go verification for `internal/harnessobs`.
- [x] T021-2931 Run focused source commit state regression evidence covering
  non-git empty commit mapping to `cannot_verify`, valid commit pass mapping
  when available, and source commit hash validation, or mark any unverified
  behavior `not_assessed` with reason.
- [x] T021-2940 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-2950 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-2960 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  43 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-43-evidence.md`.

## Active Slice 44 Tasks

- [x] T021-2970 Confirm Slice 44 is bounded to numbered
  `internal/harnessobs` event source reading handoff shard `harnessobs_344`.
- [x] T021-2971 Confirm Slice 44 is behavior-preserving: no changes to profile
  loading, event scanning, source digest return, error propagation, package
  boundary, dependency direction, path-safety behavior, raw normalization, or
  baseline change is planned.
- [x] T021-2972 Record Slice 44 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-44-plan-review.md`.
- [x] T021-2980 Move event source reading handoff into the cohesive event scan
  input file.
- [x] T021-2990 Run `gofmt` on changed Go files.
- [x] T021-3000 Run focused Go verification for `internal/harnessobs`.
- [x] T021-3001 Run focused event source reading regression evidence covering
  successful profile load plus event scan, profile load failure propagation,
  source read failure propagation, source digest return, and collect-session
  missing-source cannot_verify behavior, or mark any unverified behavior
  `not_assessed` with reason.
- [x] T021-3010 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3020 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3030 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  44 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-44-evidence.md`.

## Active Slice 45 Tasks

- [x] T021-3040 Confirm Slice 45 is bounded to numbered
  `internal/harnessobs` profile-relative source/output path safety shards
  `harnessobs_345` through `harnessobs_347`.
- [x] T021-3041 Confirm Slice 45 is behavior-preserving: no changes to
  absolute path rejection, URL-like path rejection, traversal rejection, profile
  base directory joining, existing-file validation, output-file validation,
  package boundary, dependency direction, raw normalization execution, or
  baseline change is planned.
- [x] T021-3042 Record Slice 45 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-45-plan-review.md`.
- [x] T021-3050 Move profile-relative source/output path safety into the
  cohesive session profile paths file.
- [x] T021-3060 Run `gofmt` on changed Go files.
- [x] T021-3070 Run focused Go verification for `internal/harnessobs`.
- [x] T021-3071 Run focused profile-relative path regression evidence covering
  profile-directory joining, traversal rejection, absolute path rejection,
  URL-like path rejection, existing-file requirement, output-file parent
  handling, session collect handoff, and raw-normalization path
  preflight/handoff only, setup isolation handoff, and isolation rule
  validation handoff, or mark any unverified behavior `not_assessed` with
  reason. Raw-event normalization execution (`normalizeRawEvents` /
  `harnessobs_348` onward) remains `not_assessed` for this slice.
- [x] T021-3080 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3090 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3100 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  45 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-45-evidence.md`.

## Active Slice 46 Tasks

- [x] T021-3110 Confirm Slice 46 is bounded to numbered
  `internal/harnessobs` raw-event normalization execution shards
  `harnessobs_348` through `harnessobs_360`.
- [x] T021-3111 Confirm Slice 46 is behavior-preserving: no changes to
  supported raw format gating, raw/output same-file rejection, zero-time
  fallback, scanner limits, malformed JSONL errors, unsafe-input rejection,
  normalized source digest calculation, output parent creation, JSONL write
  format, package boundary, dependency direction, or baseline change is
  planned.
- [x] T021-3112 Record Slice 46 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-46-plan-review.md`.
- [x] T021-3120 Move raw-event normalization execution into cohesive
  raw-normalization files.
- [x] T021-3130 Run `gofmt` on changed Go files.
- [x] T021-3140 Run focused Go verification for `internal/harnessobs`.
- [x] T021-3141 Run focused raw-normalization regression evidence covering
  supported/unsupported format gating, raw/output same-file rejection,
  zero-time fallback, scanner limits, blank JSONL lines, malformed JSONL,
  unsafe raw event rejection, normalized source digest calculation, output
  parent creation, JSONL write format, and session collection normalization.
- [x] T021-3150 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3160 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3170 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  46 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-46-evidence.md`.

## Active Slice 47 Tasks

- [x] T021-3180 Confirm Slice 47 is bounded to numbered `cmd/sdp-trace`
  checkpoint CLI shards `gate_271` through `gate_289`.
- [x] T021-3181 Confirm Slice 47 is behavior-preserving: no changes to
  checkpoint create/verify subcommand routing, flag names/defaults, required
  flag errors, positional argument rejection, private-key and checkpoint JSON
  loading, optional policy semantics, output JSON rendering, checkpoint create
  stdout format, verify exit-code mapping, package boundary, dependency
  direction, or baseline change is planned.
- [x] T021-3182 Record Slice 47 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-47-plan-review.md`.
- [x] T021-3190 Move checkpoint CLI command handling into cohesive checkpoint
  CLI files: `checkpoint_command.go`, `checkpoint_create_cli.go`, and
  `checkpoint_verify_cli.go`. Split further only if MI fails and record the
  failed command plus the revised responsibility boundary in evidence.
- [x] T021-3200 Run `gofmt` on changed Go files.
- [x] T021-3210 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-3211 Run focused checkpoint CLI regression evidence covering
  missing/unknown checkpoint subcommands, create flag validation, checkpoint
  creation/write success and failures, verify positional rejection, verify
  required-input validation, checkpoint/policy load handling, optional policy
  absence, result rendering, and verify exit-code mapping.
- [x] T021-3220 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3230 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3240 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  47 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-47-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 48 Tasks

- [x] T021-3250 Confirm Slice 48 is bounded to numbered `cmd/sdp-trace`
  protected-gate core shards `gate_302` through `gate_311`.
- [x] T021-3251 Confirm Slice 48 is behavior-preserving: no changes to
  protected gate setup failure handling, required checkpoint/policy/witness
  input semantics, JSON decode error handling, contract/row/run-dir/witness
  expectation error codes, protected checkpoint replay handoff,
  `PolicyProvided: true`, UTC evaluation time, result JSON writing/rendering,
  gate exit-code behavior, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-3252 Record Slice 48 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-48-plan-review.md`.
- [x] T021-3260 Move protected-gate core execution and input loading into
  `protected_gate_core.go`, `protected_gate_inputs.go`, and
  `protected_gate_loaders.go`. Split further only if MI fails and record the
  failed command plus the revised responsibility boundary in evidence.
- [x] T021-3270 Run `gofmt` on changed Go files.
- [x] T021-3280 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-3281 Run focused protected-gate core regression evidence covering
  missing/malformed required trust inputs, contract/row/run-dir/witness
  expectation failures, protected checkpoint replay handoff, protected gate
  evaluation inputs, explicit `PolicyProvided: true`, UTC evaluation time,
  result writing failure, result rendering, and exit-code behavior. Include
  test-existence evidence proving all planned focused test names are present.
- [x] T021-3290 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3300 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3310 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  48 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-48-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 49 Tasks

- [x] T021-3320 Confirm Slice 49 is bounded to numbered `cmd/sdp-trace`
  gate-explain shards `gate_312` through `gate_324`.
- [x] T021-3321 Confirm Slice 49 is behavior-preserving: no changes to explain
  usage errors, gate-result artifact loading, supported schema validation,
  missing/malformed gate-result artifact `cannot_verify`, unsupported schema
  `cannot_verify`, read-only behavior, legacy protected-field absence output,
  protected checkpoint/condition detail lines,
  required-run/witness/override/missing-evidence/reason/next-action rendering,
  secret non-disclosure by not printing raw run commands, package boundary,
  dependency direction, or baseline change is planned.
- [x] T021-3322 Record Slice 49 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-49-plan-review.md`.
- [x] T021-3330 Move gate explain CLI and rendering into
  `gate_explain_cli.go`, `gate_explain_renderer.go`,
  `gate_explain_collections.go`, and neutral shared
  `explain_common_collections.go` for reason/next-action helpers used by both
  gate and assessment explain. Split further only if MI fails and record the
  failed command plus the revised responsibility boundary in evidence.
- [x] T021-3340 Run `gofmt` on changed Go files.
- [x] T021-3350 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-3351 Run focused gate-explain regression evidence covering parse
  usage, missing/malformed artifact load `cannot_verify`, unsupported schema
  `cannot_verify`, legacy protected-field absence, protected
  checkpoint/condition details, collection rendering, reasons, next actions,
  read-only persisted verdict preservation without recomputation or upgrade,
  and secret non-disclosure. Include exact per-test `go test -list
  '^TestName$'` evidence proving all planned focused test names are present.
- [x] T021-3360 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3370 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3380 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  49 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-49-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 50 Tasks

- [x] T021-3390 Confirm Slice 50 is bounded to numbered `cmd/sdp-trace`
  gate-preview and protected-target shards `gate_325` through `gate_333`.
- [x] T021-3391 Confirm Slice 50 is behavior-preserving: no changes to preview
  read-only behavior, target arity usage errors, contract load failure
  behavior, standard preview report fields, witness mismatch reporting without
  issuing any gate verdict fields, protected preview input status mapping,
  protected preview `cannot_verify` setup exit for unreadable/malformed inputs,
  secret non-disclosure by not printing raw run commands, protected single-run
  selection semantics, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-3392 Record Slice 50 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-50-plan-review.md`.
- [x] T021-3400 Move gate preview and protected target selection into
  `gate_preview_cli.go`, `gate_preview_args.go`,
  `gate_preview_standard.go`, `gate_preview_reports.go`,
  `gate_preview_protected.go`, and neutral `protected_gate_run_dir.go`.
  Keep protected preview status/action helpers `gate_349` through `gate_351`
  as out-of-scope dependencies for this slice.
  Record the initial MI failure for the coarser CLI file plus the revised
  responsibility boundary in evidence.
- [x] T021-3410 Run `gofmt` on changed Go files.
- [x] T021-3420 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-3421 Run focused gate-preview regression evidence covering exact
  test existence, parse usage, contract load failures, standard preview
  report shape fields (`required_runs`, `required_evidence`,
  `witness_inspectable`, `witness_mismatches`, `claim`), standard preview
  read-only behavior, explicit absence of gate verdict fields in standard and
  witness-mismatch preview output, witness mismatch reporting, protected
  preview absent/unreadable/malformed input statuses and exit codes through
  existing `gate_349` through `gate_351` helpers, protected single-run
  selection, and secret non-disclosure.
- [x] T021-3430 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3440 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3450 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  50 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-50-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 51 Tasks

- [x] T021-3460 Confirm Slice 51 is bounded to numbered `cmd/sdp-trace`
  protected checkpoint trust matching shards `gate_334` through `gate_344`.
- [x] T021-3461 Confirm Slice 51 is behavior-preserving: no changes to
  checkpoint upgrade eligibility, CI-isolated signer authority requirement,
  signer id/authority/public-key policy binding, witness protected-trust
  status, witness source wildcard behavior, artifact exact-count/exact-digest
  matching, protected gate output contracts, package boundary, dependency
  direction, or baseline change is planned.
- [x] T021-3462 Record Slice 51 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-51-plan-review.md`.
- [x] T021-3470 Move protected checkpoint trust matching into
  `protected_checkpoint_trust.go`, `protected_checkpoint_signer.go`,
  `protected_witness_match.go`, and `protected_witness_artifacts.go`. Split
  further only if MI fails and record the failed command plus revised
  responsibility boundary in evidence.
- [x] T021-3480 Run `gofmt` on changed Go files.
- [x] T021-3490 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-3491 Run focused protected-trust regression evidence covering exact
  test existence, protected gate pass with CI-signed checkpoint/bound witness,
  local-signed checkpoint rejection, explicit checkpoint-fail non-upgrade,
  signer id/authority/public-key mismatch rejection, committed fixture
  protected rows, direct witness protected input matching, optional source
  wildcard behavior, exact artifact count/path/digest matching, and protected
  trust helper hotspot coverage.
- [x] T021-3500 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3510 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3520 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  51 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-51-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 52 Tasks

- [x] T021-3530 Confirm Slice 52 is bounded to numbered `cmd/sdp-trace`
  demo witness expectation and artifact digest shards `gate_345` through
  `gate_348`.
- [x] T021-3531 Confirm Slice 52 is behavior-preserving: no changes to
  deriving witness expectations from observed run directories, first-run ID
  anchoring, `<run-dir-base>/run.json` artifact path shape, retained-byte
  SHA-256 digest calculation, discovery/open error propagation through
  protected witness expectation loading, package boundary, dependency
  direction, or baseline change is planned.
- [x] T021-3532 Record Slice 52 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-52-plan-review.md`.
- [x] T021-3540 Move demo witness expectation construction into
  `protected_witness_expectation.go` and `protected_witness_digest.go`. Split
  further only if MI fails and record the failed command plus revised
  responsibility boundary in evidence.
- [x] T021-3550 Run `gofmt` on changed Go files.
- [x] T021-3560 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-3561 Run focused witness-expectation regression evidence covering
  exact test existence, missing-run rejection, first-run ID anchoring,
  per-run `<run-dir-base>/run.json` paths, retained-byte SHA-256 digests,
  retained `run.json` open/read error propagation, standard witness mismatch
  failure, protected gate CI-signed/bound witness pass, and no baseline changes.
- [x] T021-3570 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3580 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3590 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  52 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-52-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 53 Tasks

- [x] T021-3600 Confirm Slice 53 is bounded to numbered `cmd/sdp-trace`
  protected preview input status/action shards `gate_349` through `gate_351`.
- [x] T021-3601 Confirm Slice 53 is behavior-preserving: no changes to blank
  input `absent`, missing/permission-denied `present_unreadable`, malformed
  JSON `present_malformed`, readable JSON `present_readable`, stable
  checkpoint/checkpoint_policy/witness next-action ordering, protected preview
  `cannot_verify` setup exit for unreadable/malformed inputs, package boundary,
  dependency direction, or baseline change is planned.
- [x] T021-3602 Record Slice 53 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-53-plan-review.md`.
- [x] T021-3610 Move protected preview input status/action helpers into
  `protected_preview_inputs.go`. Split further only if MI fails and record the
  failed command plus revised responsibility boundary in evidence.
- [x] T021-3620 Run `gofmt` on changed Go files.
- [x] T021-3630 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-3631 Run focused protected-preview input regression evidence
  covering exact test existence, status mapping branches, stable action order,
  explicit permission-denied to `present_unreadable` mapping, protected preview
  absent input rendering without writes, unreadable/malformed input
  `cannot_verify` exits, readable input pass-through, and no protected verdict
  field emission.
- [x] T021-3640 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3650 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3660 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  53 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-53-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 54 Tasks

- [x] T021-3670 Confirm Slice 54 is bounded to numbered `cmd/sdp-trace`
  override request CLI shards `gate_352` through `gate_359`.
- [x] T021-3671 Confirm Slice 54 is behavior-preserving: no changes to
  accepting only `override request`, parser error handling, missing required
  flag diagnostics, positional text rejection, append-before-print behavior,
  append failure exit code/stderr, stable override payload keys, optional
  `external_reference`, UTC RFC3339Nano `created_at`, advisory/non-upgrading
  override semantics, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-3672 Record Slice 54 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-54-plan-review.md`.
- [x] T021-3680 Move override request CLI handling into
  `override_request.go`. Split further only if MI fails and record the failed
  command plus revised responsibility boundary in evidence.
- [x] T021-3690 Run `gofmt` on changed Go files.
- [x] T021-3700 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-3701 Run focused override request regression evidence covering
  exact test existence for
  `TestOverrideRequestAppendsFlightRecorderEvent`,
  `TestOverrideRequestPersistsExternalReferencePayload`,
  `TestOverrideRequestRejectsInvalidRequestArgsBeforeAppend`,
  `TestOverrideRequestAppendFailureDoesNotPrintEvent`, and
  `TestGateOutputIncludesOverrideWithoutPassingMissingEvidence`; successful
  event append; stable persisted payload keys; optional
  `external_reference`; rejection of non-`request` subcommands; unknown-flag
  parser errors; positional text rejection; missing required flag diagnostics;
  append failure behavior; advisory/non-upgrading gate output; and no baseline
  changes.
- [x] T021-3710 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3720 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3730 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  54 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-54-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 55 Tasks

- [x] T021-3740 Confirm Slice 55 is bounded to numbered `cmd/sdp-trace` shared
  artifact IO shards `gate_360` through `gate_364`.
- [x] T021-3741 Confirm Slice 55 is behavior-preserving: no changes to JSON
  read/unmarshal error propagation, JSON parent-directory creation,
  two-space pretty JSON plus trailing newline, JSON `os.WriteFile` requested
  `0o644` mode subject to process umask, text
  parent-directory creation, sibling temp-file naming, temp cleanup,
  close-on-write-error behavior, chmod-before-rename, atomic rename
  publication, package boundary, dependency direction, or baseline change is
  planned.
- [x] T021-3742 Record Slice 55 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-55-plan-review.md`.
- [x] T021-3750 Move shared artifact IO helpers into `artifact_json_io.go` and
  `artifact_text_io.go` after the initial `artifact_io.go` consolidation fails
  file-level MI; record the failed command plus revised responsibility boundary
  in evidence.
- [x] T021-3760 Run `gofmt` on changed Go files.
- [x] T021-3770 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-3771 Run focused artifact IO regression evidence covering exact
  test existence for `TestReadJSONFilePropagatesReadAndDecodeErrors`,
  `TestWriteJSONFileCreatesPrettyJSONWithNewline`, and
  `TestWriteTextFileAtomicPublishesCompleteTextAndCleansTemp`,
  `TestFinishAtomicTextWriteNormalizesModeBeforeRename`, and
  `TestWriteAndCloseTempTextReturnsWriteErrorOnClosedFile`, and
  `TestWriteTextFileAtomicRemovesTempOnRenameFailure`; read error propagation;
  decode error propagation; pretty JSON newline; parent directory creation;
  JSON file is readable through the just-created path and has no executable
  bits; sibling temp-file cleanup after successful rename; temp cleanup on
  rename failure after a temp file exists; atomic text overwrite publication;
  chmod-before-rename behavior; write-failure error propagation through a
  closed temp file; and no baseline changes.
- [x] T021-3780 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3790 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3800 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  55 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-55-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 56 Tasks

- [x] T021-3810 Confirm Slice 56 is bounded to numbered `cmd/sdp-trace` gate
  preview contract helper shards `gate_365` through `gate_367`.
- [x] T021-3811 Confirm Slice 56 is behavior-preserving: no changes to
  observation default mode, advisory CI mode selection, protected-future
  dominance over advisory CI regardless of required-run order, unknown/empty
  required-run profile ignoring, required run ID ordering with empty IDs
  omitted, required evidence ID ordering with empty IDs omitted, package
  boundary, dependency direction, or baseline change is planned.
- [x] T021-3812 Record Slice 56 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-56-plan-review.md`.
- [x] T021-3820 Move gate preview contract helpers into
  `gate_preview_contract.go`. Split further only if MI fails and record the
  failed command plus revised responsibility boundary in evidence.
- [x] T021-3830 Run `gofmt` on changed Go files.
- [x] T021-3840 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-3841 Run focused gate preview contract regression evidence covering
  exact test existence for `TestPreviewGateModeSelection`,
  `TestRequiredRunIDsOmitEmptyAndKeepOrder`, and
  `TestRequiredEvidenceIDsForCLIOmitEmptyAndKeepOrder`; observation default;
  advisory CI fallback; protected-future dominance regardless of order;
  unknown/empty profile ignoring; ID order preservation; empty ID omission; and
  no baseline changes.
- [x] T021-3850 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3860 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3870 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  56 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-56-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 57 Tasks

- [x] T021-3880 Confirm Slice 57 is bounded to numbered `cmd/sdp-trace` packet
  command surface shards `packet_031` through `packet_032`.
- [x] T021-3881 Confirm Slice 57 is behavior-preserving: no changes to packet
  subcommand names, handler bindings, required flag names, required flag
  diagnostic messages, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-3882 Record Slice 57 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-57-plan-review.md`.
- [x] T021-3890 Move packet command surface helpers into
  `packet_command_surface.go`. Split further only if MI fails and record the
  failed command plus revised responsibility boundary in evidence.
- [x] T021-3900 Run `gofmt` on changed Go files.
- [x] T021-3910 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-3911 Run focused packet command surface regression evidence covering
  exact test existence for `TestPacketHandlersExposeExpectedSubcommands` and
  `TestPacketRequiredFlagsKeepNamesAndDiagnostics`; exact subcommand set; all
  five handler bindings by function identity; help/known packet command smoke
  coverage; required flag order and diagnostics; and no baseline changes.
- [x] T021-3920 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-3930 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-3940 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  57 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-57-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 58 Tasks

- [x] T021-3950 Confirm Slice 58 is bounded to numbered `cmd/sdp-trace` packet
  build-pr flow and artifact publication shards `packet_040` through
  `packet_050`.
- [x] T021-3951 Confirm Slice 58 is behavior-preserving: no changes to packet
  dispatch missing-subcommand diagnostic, build-pr flag defaults,
  required flag validation, JSON `cannot_verify` output semantics, output file
  names, validation and live-gate error aggregation, artifact write labels,
  artifact write order, first-write-failure short-circuiting, output-directory
  creation with mode subject to process umask, package boundary, dependency
  direction, or baseline change is planned.
- [x] T021-3952 Record Slice 58 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-58-plan-review.md`.
- [x] T021-3960 Move packet build-pr command helpers into
  `packet_command_dispatch.go`, `packet_build_pr_run.go`,
  `packet_build_pr_options.go`, `packet_build_pr_result.go`,
  `packet_build_pr_artifact_render.go`, `packet_build_pr_artifact_write.go`,
  and `packet_build_pr_artifact_files.go` after the single-file
  `packet_build_pr_command.go` attempt failed file-level MI at 56.6 and the
  three-file split still failed at 67.7/64.1. Record the failed commands plus
  revised responsibility boundary in evidence.
- [x] T021-3970 Run `gofmt` on changed Go files.
- [x] T021-3980 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-3981 Run focused packet build-pr flow regression evidence covering
  exact test existence for
  `TestPacketDispatchRequiresKnownPacketSubcommand`,
  `TestParsePacketBuildPROptionsKeepsFlagDefaultsAndRequiredOut`,
  `TestBuildPacketPRResultKeepsPathsAndAggregatesGateErrors`,
  `TestRunPacketBuildPRWritesCannotVerifyJSONForInputFailure`,
  `TestRenderPacketPRMarkdownDowngradesResultOnFailure`,
  `TestWritePacketPRArtifactsWritesCannotVerifyJSONOnRenderFailure`,
  `TestWritePacketPRFilesCreatesOutputDirectoryAndArtifacts`,
  `TestWritePacketPRArtifactFilesStopsAfterFirstFailure`, and existing
  `TestPacketBuildPRFixtureCLIWritesArtifacts`; packet dispatch
  missing-subcommand behavior; build-pr flag default parsing; required
  flag validation; output path names; validation/live-gate error aggregation;
  JSON `cannot_verify` output paths for input reconstruction and rendering
  failures; cannot_verify result state for validation/live-gate failures;
  output-directory creation with mode subject to process umask; artifact file
  labels and ordering; artifact writer first-failure short-circuiting; CLI
  build-pr smoke coverage; and no baseline changes.
- [x] T021-3990 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-4000 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4010 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  58 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-58-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 59 Tasks

- [x] T021-4020 Confirm Slice 59 is bounded to numbered `cmd/sdp-trace` packet
  build-pr live gate error shards `packet_051` through `packet_053`.
- [x] T021-4021 Confirm Slice 59 is behavior-preserving: no changes to row ID
  lookup, `PC-AGENT-ROUTE` pass/partial acceptance, `PC-VERIFICATION` pass-only
  acceptance, route-before-verification error ordering, diagnostic strings and
  row reason inclusion, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-4022 Record Slice 59 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-59-plan-review.md`.
- [x] T021-4030 Move packet build-pr live gate helpers into
  `packet_build_pr_gate_errors.go`. Split further only if MI fails and record
  the failed command plus revised responsibility boundary in evidence.
- [x] T021-4040 Run `gofmt` on changed Go files.
- [x] T021-4050 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4051 Run focused packet build-pr live gate regression evidence
  covering exact test existence for
  `TestPacketBuildPRGateErrorsPreserveRouteAndVerificationOrder`,
  `TestPacketBuildPRRouteErrorsAcceptPassAndPartial`, and
  `TestPacketBuildPRVerificationErrorsRequirePass`; row ID lookup; route
  pass/partial acceptance; route failure diagnostic with row reason;
  verification pass-only acceptance; verification failure diagnostic with row
  reason; route-before-verification ordering; and no baseline changes.
- [x] T021-4060 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-4070 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4080 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  59 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-59-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.
