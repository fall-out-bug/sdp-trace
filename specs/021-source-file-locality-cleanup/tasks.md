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

## Active Slice 60 Tasks

- [x] T021-4090 Confirm Slice 60 is bounded to numbered `cmd/sdp-trace` packet
  build-pr input loading shards `packet_054` through `packet_059`.
- [x] T021-4091 Confirm Slice 60 is behavior-preserving: no changes to allowed
  source values, unsupported-source diagnostics, missing-event diagnostics,
  GitHub Actions `GITHUB_EVENT_PATH` fallback rules, fixture explicit event
  path handling, optional checks/artifacts JSON error prefixes, route manifest
  error prefix, route application ordering after hydration, fixture-mode
  hermeticity, package boundary, dependency direction, or baseline change is
  planned.
- [x] T021-4092 Record Slice 60 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-60-plan-review.md`.
- [x] T021-4100 Move packet build-pr input loading helpers into
  `packet_build_pr_input_source.go` and
  `packet_build_pr_input_enrichment.go` after the single-file
  `packet_build_pr_input_loading.go` attempt failed file-level MI at 65.4.
  Record the failed command plus revised responsibility boundary in evidence.
- [x] T021-4110 Run `gofmt` on changed Go files.
- [x] T021-4120 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4121 Run focused packet build-pr input loading regression evidence
  covering exact test existence for
  `TestValidPRInputSourceAcceptsOnlyKnownSources`,
  `TestPREventPathUsesActionsEnvOnlyWhenEventPathMissing`,
  `TestLoadPRInputSourceEventRejectsUnsupportedAndMissingEvent`,
  `TestReadOptionalPREvidenceKeepsErrorPrefixes`, and
  `TestBuildPRInputFromOptionsAppliesOptionalEvidenceAndRoute`; allowed source
  set; unsupported-source diagnostics; missing-event diagnostics; env fallback
  rules; explicit fixture event path handling; checks/artifacts error prefixes;
  route manifest error prefix; successful route application; fixture-mode
  hermeticity; and no baseline changes.
- [x] T021-4130 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-4140 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4150 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  60 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-60-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 61 Tasks

- [x] T021-4160 Confirm Slice 61 is bounded to numbered `cmd/sdp-trace` packet
  build-pr event-to-input mapping shards `packet_060` through `packet_062`.
- [x] T021-4161 Confirm Slice 61 is behavior-preserving: no changes to input
  schema version, prompt-boundary requirement default, GitHub Actions workflow
  run ID source, fixture workflow run ID source, PR field mapping, commit-range
  mapping, changed-files diff URL mapping, package boundary, dependency
  direction, or baseline change is planned.
- [x] T021-4162 Record Slice 61 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-61-plan-review.md`.
- [x] T021-4170 Move packet build-pr event mapping helpers into
  `packet_build_pr_event_mapping.go`. Split further only if MI fails and
  record the failed command plus revised responsibility boundary in evidence.
- [x] T021-4180 Run `gofmt` on changed Go files.
- [x] T021-4190 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4191 Run focused packet build-pr event mapping regression evidence
  covering exact test existence for
  `TestGitHubPRInputFromEventUsesActionsEnvRunID`,
  `TestGitHubPRInputFromEventUsesFixtureRunID`, `TestGitHubPRFromEventMapsPRFields`,
  and `TestGitHubCommitRangeFromEventMapsSHAsAndDiffURL`; schema version;
  prompt-boundary default; actions vs fixture workflow run ID source; PR
  number/URL/title/body ref/author/base ref/head ref/head SHA mapping; commit
  base/head SHA mapping; changed-files diff URL mapping; and no baseline
  changes.
- [x] T021-4200 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-4210 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4220 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  61 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-61-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 62 Tasks

- [x] T021-4230 Confirm Slice 62 is bounded to numbered `cmd/sdp-trace` packet
  build-pr GitHub Actions hydration dispatch shards `packet_063` through
  `packet_064`.
- [x] T021-4231 Confirm Slice 62 is behavior-preserving: no changes to
  fixture-mode no-network behavior, source gating, explicit artifact JSON
  precedence over live discovery, live discovery call conditions, error
  propagation, package boundary, dependency direction, or baseline change is
  planned.
- [x] T021-4232 Record Slice 62 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-62-plan-review.md`.
- [x] T021-4240 Move packet build-pr GitHub Actions hydration dispatch helpers
  into `packet_build_pr_actions_hydration.go`. Split further only if MI fails
  and record the failed command plus revised responsibility boundary in
  evidence.
- [x] T021-4250 Run `gofmt` on changed Go files.
- [x] T021-4260 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4261 Run focused packet build-pr GitHub Actions hydration regression
  evidence covering exact test existence for
  `TestHydrateGitHubActionsEvidenceSkipsFixtureSource`,
  `TestHydrateGitHubActionsArtifactsKeepsExplicitArtifacts`, and
  `TestHydrateGitHubActionsArtifactsBackfillsDiscoveredArtifacts`, and
  `TestHydrateGitHubActionsEvidencePropagatesLiveArtifactErrors`; fixture-mode
  no-network behavior; explicit artifact JSON precedence; successful live
  artifact backfill; live artifact discovery error propagation; and no baseline
  changes.
- [x] T021-4270 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-4280 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4290 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  62 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-62-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 63 Tasks

- [x] T021-4300 Confirm Slice 63 is bounded to numbered `cmd/sdp-trace` packet
  build-pr route manifest helper shards `packet_065` through `packet_066`.
- [x] T021-4301 Confirm Slice 63 is behavior-preserving: no changes to optional
  route JSON read behavior, route field, prompt boundary, integration action,
  and review overwrite semantics, PR identity preservation, CI evidence
  preservation, package boundary, dependency direction, or baseline change is
  planned.
- [x] T021-4302 Record Slice 63 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-63-plan-review.md`.
- [x] T021-4310 Move packet build-pr route manifest helpers into
  `packet_build_pr_route.go`. Split further only if MI fails and record the
  failed command plus revised responsibility boundary in evidence.
- [x] T021-4320 Run `gofmt` on changed Go files.
- [x] T021-4330 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4331 Run focused packet build-pr route regression evidence covering
  exact test existence for `TestReadOptionalPRRouteKeepsOptionalJSONBehavior`
  and `TestApplyPRRouteOnlyOverwritesRouteAndReviewFields`; optional route JSON
  read behavior; route field, prompt boundary, integration action, and review
  overwrite semantics; PR identity preservation; CI evidence preservation; and
  no baseline changes.
- [x] T021-4340 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-4350 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4360 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  63 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-63-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 64 Tasks

- [x] T021-4370 Confirm Slice 64 is bounded to numbered `cmd/sdp-trace` GitHub
  Actions artifact discovery facade and response type shards `packet_067`
  through `packet_068`.
- [x] T021-4371 Confirm Slice 64 is behavior-preserving: no changes to context
  validation before live fetch, fetch error propagation, retained artifact
  filtering, empty-retained-set fail-closed diagnostics, package boundary,
  dependency direction, or baseline change is planned.
- [x] T021-4372 Record Slice 64 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-64-plan-review.md`.
- [x] T021-4380 Move GitHub Actions artifact discovery facade and response
  types into `packet_build_pr_actions_artifacts.go`. Split further only if MI
  fails and record the failed command plus revised responsibility boundary in
  evidence.
- [x] T021-4390 Run `gofmt` on changed Go files.
- [x] T021-4400 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4401 Run focused GitHub Actions artifact discovery regression
  evidence covering exact test existence for
  `TestGitHubActionsArtifactsBackfillsRetainedArtifacts`,
  `TestGitHubActionsArtifactsFailsClosedWithoutRetainedArtifacts`, and
  `TestGitHubActionsArtifactsInvalidContextStopsBeforeFetch`, and
  `TestGitHubActionsArtifactsPropagatesFetchErrors`; context validation before
  live fetch; retained artifact filtering; empty retained artifact fail-closed
  diagnostic; fetch error propagation; and no baseline changes.
- [x] T021-4410 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-4420 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4430 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  64 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-64-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 65 Tasks

- [x] T021-4440 Confirm Slice 65 is bounded to numbered `cmd/sdp-trace` GitHub
  Actions artifact context construction and source selection shards
  `packet_071` through `packet_075`.
- [x] T021-4441 Confirm Slice 65 is behavior-preserving: no changes to API URL
  flag-over-env-over-default precedence, URL validation before returned context
  construction, trailing-slash trimming, `GITHUB_TOKEN` precedence over
  `GH_TOKEN`, missing repo/run and token diagnostics, package boundary,
  dependency direction, or baseline change is planned.
- [x] T021-4442 Record Slice 65 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-65-plan-review.md`.
- [x] T021-4450 Move GitHub Actions artifact context helpers into
  `packet_build_pr_actions_context.go` and
  `packet_build_pr_actions_source.go` after the single-file
  `packet_build_pr_actions_context.go` attempt failed file-level MI at 67.7.
  Record the failed command plus revised responsibility boundary in evidence.
- [x] T021-4460 Run `gofmt` on changed Go files.
- [x] T021-4470 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4471 Run focused GitHub Actions artifact context regression evidence
  covering exact test existence for
  `TestNewGitHubActionsArtifactContextBuildsValidatedContext`,
  `TestGitHubAPIURLSelectionKeepsPrecedenceAndTrimming`,
  `TestGitHubTokenPrefersGitHubTokenThenFallsBack`, and
  `TestValidateGitHubActionsArtifactContextKeepsDiagnostics`; API URL
  precedence and trimming; token precedence/fallback; missing repo/run
  diagnostic; missing token diagnostic; exact test list verification that fails
  when any planned focused test is missing; and no baseline changes.
- [x] T021-4480 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-4490 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4500 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  65 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-65-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 66 Tasks

- [x] T021-4510 Confirm Slice 66 is bounded to numbered `cmd/sdp-trace` GitHub
  Actions API URL validation and trust-target policy shards `packet_076`
  through `packet_085`.
- [x] T021-4511 Confirm Slice 66 is behavior-preserving: no changes to syntax
  diagnostics, embedded-credential rejection before trust-target validation,
  credential rejection winning over HTTPS errors for mixed-invalid URLs,
  loopback-only HTTP allowance, public GitHub host mapping, Enterprise exact
  host binding, malformed server URL fallback behavior, package boundary,
  dependency direction, or baseline change is planned.
- [x] T021-4512 Record Slice 66 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-66-plan-review.md`.
- [x] T021-4520 Move GitHub Actions API URL policy helpers into
  `packet_build_pr_actions_url_policy.go`,
  `packet_build_pr_actions_url_parse.go`, and
  `packet_build_pr_actions_url_host.go` after the single-file
  `packet_build_pr_actions_url_policy.go` attempt failed file-level MI at 61.6.
  Record the failed command plus revised responsibility boundary in evidence.
- [x] T021-4530 Run `gofmt` on changed Go files.
- [x] T021-4540 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4541 Run focused GitHub Actions API URL policy regression evidence
  covering exact test existence for `TestGitHubAPIURLPolicy`,
  `TestGitHubAPIURLPolicyHelpersKeepTrustBoundaries`,
  `TestParseGitHubAPIURLKeepsSyntaxDiagnostics`, and
  `TestGitHubServerHostKeepsFallbackBehavior`; syntax diagnostics; embedded
  credential rejection before trust-target validation; credential rejection
  winning over HTTPS errors for mixed-invalid URLs; HTTPS requirement;
  loopback-only HTTP allowance; public GitHub API host mapping; Enterprise
  exact-host binding; malformed server URL fallback behavior; exact test list
  verification after test implementation that fails when any planned focused
  test is missing; and no baseline changes.
- [x] T021-4550 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T021-4560 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4570 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  66 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-66-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 67 Tasks

- [x] T021-4580 Confirm Slice 67 is bounded to numbered `cmd/sdp-trace` GitHub
  Actions artifact HTTP request/fetch/decode and retained artifact shaping
  shards `packet_086` through `packet_092`.
- [x] T021-4581 Confirm Slice 67 is behavior-preserving: no changes to
  credential fail-closed behavior for malformed or non-HTTPS API URLs,
  loopback-test no-token behavior, GitHub media type headers, non-2xx
  fail-closed diagnostics, JSON decode diagnostics, expired artifact filtering,
  artifact URL precedence over synthesized resolver URLs, missing artifact ID
  empty-resolver behavior, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-4582 Record Slice 67 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-67-plan-review.md`.
- [x] T021-4590 Move GitHub Actions artifact HTTP fetch, request,
  authorization, decode/status, and retained artifact helpers into
  `packet_build_pr_actions_fetch.go`, `packet_build_pr_actions_request.go`,
  `packet_build_pr_actions_authorization.go`,
  `packet_build_pr_actions_decode.go`, and
  `packet_build_pr_actions_retention.go` after the single-file
  `packet_build_pr_actions_http_retention.go` attempt failed file-level MI at
  66.6 and the initial HTTP/retention split failed HTTP file-level MI at 67.5.
  Record the failed commands plus revised responsibility boundary in evidence.
- [x] T021-4600 Run `gofmt` on changed Go files.
- [x] T021-4610 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4611 Run focused GitHub Actions artifact HTTP/retention regression
  evidence covering exact test existence for
  `TestGitHubActionsArtifactsRequestKeepsHTTPContract`,
  `TestGitHubActionsArtifactsAuthorizationFailClosed`,
  `TestSuccessfulHTTPStatusAcceptsOnly2xx`,
  `TestDecodeGitHubActionsArtifactsKeepsDiagnostics`, and
  `TestRetainedGitHubArtifactsKeepsResolverPolicy`; media type header and URL
  path construction; HTTPS-only authorization with malformed/non-HTTPS
  fail-closed behavior; loopback-test no-token behavior; 2xx-only status
  policy; JSON decode diagnostics; expired artifact filtering; artifact URL
  precedence; synthesized resolver URL construction; missing artifact ID empty
  resolver behavior; exact test list verification after test implementation
  that fails when any planned focused test is missing; and no baseline changes.
- [x] T021-4620 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-4630 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4640 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  67 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-67-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 68 Tasks

- [x] T021-4650 Confirm Slice 68 is bounded to remaining numbered
  `cmd/sdp-trace` packet fixture and validation exit helper shards
  `packet_093` through `packet_096`.
- [x] T021-4651 Confirm Slice 68 is behavior-preserving: no changes to
  required PR fixture identity validation, optional empty-path no-op behavior,
  JSON read/unmarshal error propagation, `pass` to zero exit mapping, packet
  validation failure to `cannot_verify`, demo gate failure to `fail`, package
  boundary, dependency direction, or baseline change is planned.
- [x] T021-4652 Record Slice 68 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-68-plan-review.md`.
- [x] T021-4660 Move PR fixture event shape/loading, shared optional JSON
  loading, and packet exit helpers into `packet_build_pr_fixture_event.go`,
  `packet_build_pr_optional_json.go`, and `packet_validation_exits.go` after the
  single-file `packet_build_pr_fixture_io.go` attempt failed file-level MI at
  66.3. Record the failed command plus revised responsibility boundary in
  evidence.
- [x] T021-4670 Run `gofmt` on changed Go files.
- [x] T021-4680 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4681 Run focused packet fixture and validation-exit regression evidence
  covering exact test existence for `TestLoadPRFixtureEventRequiresIdentity`,
  `TestReadOptionalJSONKeepsOptionalAndErrorBehavior`, and
  `TestPacketExitMappingsKeepTrustSemantics`; required PR number and URL
  validation; optional empty path no-op; JSON read and unmarshal error
  propagation; `pass` to zero exit mapping; packet validation failure to
  `cannot_verify`; demo gate failure to `fail`; exact test list verification
  after test implementation that fails when any planned focused test is
  missing; and no baseline changes.
- [x] T021-4690 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-4700 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4710 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  68 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-68-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 69 Tasks

- [x] T021-4720 Confirm Slice 69 is bounded to numbered `cmd/sdp-trace`
  `pr_review` top-level dispatch and packet subcommand setup shards
  `pr_review_030`, `pr_review_037` through `pr_review_039`, and
  `pr_review_098` through `pr_review_104`, plus packet required-input helper
  shard `pr_review_138`.
- [x] T021-4721 Confirm Slice 69 is behavior-preserving: no changes to
  subcommand routing, usage and missing-subcommand diagnostics, required packet
  flags and defaults, repeated context/verification flag reconstruction,
  optional metadata mapping, positional-argument rejection, missing
  packet-anchor usage errors, packet build failures mapping to
  `cannot_verify`, packet stdout rendering, provenance anchor mapping, package
  boundary, dependency direction, or baseline change is planned.
- [x] T021-4722 Record Slice 69 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-69-plan-review.md`.
- [x] T021-4730 Move top-level `pr-review` dispatch, packet flag metadata,
  packet command parsing/execution, and packet option construction into
  `pr_review_command.go`, `pr_review_packet_flags.go`,
  `pr_review_packet_args.go`, `pr_review_packet_run.go`, and
  `pr_review_packet_options.go` after the combined top-level command/packet
  flag file failed file-level MI at `56.6`, the combined packet
  command/options file failed file-level MI at `68.1`, and the initial packet
  command file with required-input checks failed file-level MI at `68.6`.
  Record failed commands plus revised responsibility boundaries in evidence.
- [x] T021-4740 Run `gofmt` on changed Go files.
- [x] T021-4750 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4751 Run focused `pr-review` command/packet regression evidence
  covering exact test existence for `TestPRReviewHandlersKeepSubcommands`,
  `TestPRReviewPacketFlagsKeepContract`,
  `TestParsePRReviewPacketArgsKeepsUsageBoundaries`, and
  `TestPRReviewPacketOptionsKeepsEvidenceMapping`; subcommand routing; usage
  and missing-subcommand diagnostics; required packet flags and defaults;
  positional-argument rejection; missing packet-anchor usage errors; packet
  build failure `cannot_verify` behavior; repeated context and verification
  flag reconstruction; optional metadata mapping; provenance anchor mapping;
  exact test list verification after test implementation that fails when any
  planned focused test is missing; and no baseline changes.
- [x] T021-4760 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-4770 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4780 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  69 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-69-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 70 Tasks

- [x] T021-4790 Confirm Slice 70 is bounded to numbered `cmd/sdp-trace`
  `pr-review run` execution shards `pr_review_105` through `pr_review_108`.
- [x] T021-4791 Confirm Slice 70 is behavior-preserving: no changes to
  packet/profile loading, work-dir directory validation, repeated
  allowed-runner reconstruction from raw args, preview mode, not-assessed
  reason propagation, parse errors and positional-argument rejection as usage
  errors, packet/profile read failures mapping to `cannot_verify`, preview
  output not implying evidence production, package boundary, dependency
  direction, or baseline change is planned.
- [x] T021-4792 Record Slice 70 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-70-plan-review.md`.
- [x] T021-4800 Move `pr-review run` execution, argument parsing, and output
  rendering into `pr_review_run_command.go`, `pr_review_run_args.go`, and
  `pr_review_run_output.go` after the single-file `pr_review_run_command.go`
  attempt failed file-level MI at `66.1`. Record the failed command plus
  revised responsibility boundary in evidence, and synchronize
  `cmd/sdp-trace/FAMILY_INDEX.md` with the renamed files.
- [x] T021-4810 Run `gofmt` on changed Go files.
- [x] T021-4820 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4821 Run focused `pr-review run` regression evidence covering exact
  test existence for `TestParsePRReviewRunArgsKeepsUsageBoundaries`,
  `TestExecutePRReviewRunKeepsRunnerOptions`, and
  `TestWritePRReviewRunOutputKeepsPreviewBoundary`; parse defaults; parse
  errors and positional argument rejection; packet/profile read failure
  `cannot_verify` behavior; work-dir directory validation; repeated
  allowed-runner reconstruction from raw args; preview mode and not-assessed
  reason propagation; preview output shape not implying produced run evidence;
  exact test list verification after test implementation that fails when any
  planned focused test is missing; and no baseline changes.
- [x] T021-4830 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-4840 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4850 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  70 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-70-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 71 Tasks

- [x] T021-4860 Confirm Slice 71 is bounded to numbered `cmd/sdp-trace`
  `pr-review synthesize` ledger collation shards `pr_review_109` through
  `pr_review_114`.
- [x] T021-4861 Confirm Slice 71 is behavior-preserving: no changes to
  mandatory output path validation, packet/run-set reads, optional
  existing-ledger empty-path no-op behavior, existing ledger as synthesis input
  rather than authority, artifact read failures mapping to `cannot_verify`,
  durable ledger write before stdout mirroring, package boundary, dependency
  direction, or baseline change is planned.
- [x] T021-4862 Record Slice 71 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-71-plan-review.md`.
- [x] T021-4870 Move `pr-review synthesize` execution and argument parsing into
  `pr_review_synthesize_command.go`, and packet/run/existing-ledger input
  loading into `pr_review_synthesis_inputs.go`, after pre-change MI analysis
  measured the planned files at `73.5` and `72.8`. Record the command and
  responsibility boundary in evidence, and synchronize
  `cmd/sdp-trace/FAMILY_INDEX.md` with the renamed files.
- [x] T021-4880 Run `gofmt` on changed Go files.
- [x] T021-4890 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4891 Run focused `pr-review synthesize` regression evidence
  covering exact test existence for
  `TestParsePRReviewSynthesizeArgsKeepsUsageBoundaries`,
  `TestReadPRReviewSynthesisInputsKeepsOptionalLedgerBoundary`, and
  `TestRunPRReviewSynthesizeKeepsLedgerDurability`; mandatory output path
  validation, positional argument rejection, packet read failure
  `cannot_verify` behavior, run-set read failure `cannot_verify` behavior,
  optional existing-ledger empty path, existing ledger read failure propagation,
  durable ledger write before stdout mirroring, exact test list verification
  after test implementation that fails when any planned focused test is missing,
  and no baseline changes.
- [x] T021-4900 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-4910 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4920 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  71 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-71-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 72 Tasks

- [x] T021-4930 Confirm Slice 72 is bounded to numbered `cmd/sdp-trace`
  `pr-review validate` artifact-validation shards `pr_review_115` through
  `pr_review_120`.
- [x] T021-4931 Confirm Slice 72 is behavior-preserving: no changes to
  mandatory output path validation, positional-argument rejection,
  packet/profile/run-set/ledger reads, artifact read failures mapping to
  `cannot_verify`, package-owned validation verdicts, durable validation write
  before stdout mirroring, `reviewValidationExitCode` mapping through
  `exitCannotVerify`, package boundary, dependency direction, or baseline
  change is planned.
- [x] T021-4932 Record Slice 72 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-72-plan-review.md`.
- [x] T021-4940 Move `pr-review validate` execution, argument parsing, output
  path validation, durable validation write, stdout mirroring, and validation
  exit mapping into `pr_review_validate_command.go`, and
  packet/profile/run-set/ledger input loading into
  `pr_review_validation_inputs.go`, after pre-change MI analysis measured the
  planned files at `72.0` and `73.2`. Record the command and responsibility
  boundary in evidence, and synchronize `cmd/sdp-trace/FAMILY_INDEX.md` with
  the renamed files.
- [x] T021-4950 Run `gofmt` on changed Go files.
- [x] T021-4960 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-4961 Run focused `pr-review validate` regression evidence covering
  exact test existence for
  `TestParsePRReviewValidateArgsKeepsUsageBoundaries`,
  `TestReadPRReviewValidationInputsKeepsArtifactBoundaries`, and
  `TestRunPRReviewValidateKeepsDurableVerdictAndExitMapping`; mandatory output
  path validation, positional argument rejection, packet/profile/run-set/ledger
  read failure `cannot_verify` behavior, durable validation write before stdout
  mirroring, non-zero validation verdict mapping through `exitCannotVerify`,
  exact test list verification after test implementation that fails when any
  planned focused test is missing, and no baseline changes.
- [x] T021-4970 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-4980 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-4990 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  72 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-72-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 73 Tasks

- [x] T021-5000 Confirm Slice 73 is bounded to numbered `cmd/sdp-trace`
  `pr-review summarize` human-readable summary shards `pr_review_121` through
  `pr_review_125`.
- [x] T021-5001 Confirm Slice 73 is behavior-preserving: no changes to summary
  text being UX-only output rather than proof or approval, validation/ledger
  reads, artifact read failures mapping to `cannot_verify`,
  positional-argument rejection, optional output path behavior, write-once
  refusal for existing summary files, stdout mirroring when a durable summary
  file is requested, package boundary, dependency direction, or baseline change
  is planned.
- [x] T021-5002 Record Slice 73 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-73-plan-review.md`.
- [x] T021-5010 Move `pr-review summarize` execution and argument parsing into
  `pr_review_summarize_command.go`, and validation/ledger input loading plus
  optional summary-file writing into `pr_review_summary_io.go`, after
  pre-change MI analysis measured the planned files at `73.5` and `76.6`.
  Record the command and responsibility boundary in evidence, and synchronize
  `cmd/sdp-trace/FAMILY_INDEX.md` with the renamed files.
- [x] T021-5020 Run `gofmt` on changed Go files.
- [x] T021-5030 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-5031 Run focused `pr-review summarize` regression evidence covering
  exact test existence for
  `TestParsePRReviewSummarizeArgsKeepsUsageBoundaries`,
  `TestReadPRReviewSummaryInputsKeepsArtifactBoundaries`, and
  `TestRunPRReviewSummarizeKeepsUXOnlyOutputBoundary`; positional argument
  rejection, validation/ledger read failure `cannot_verify` behavior, optional
  output path no-op behavior, existing summary-file refusal as usage failure,
  stdout mirroring when a durable summary file is requested, summary text not
  implying merge approval, exact test list verification after test
  implementation that fails when any planned focused test is missing, and no
  baseline changes.
- [x] T021-5040 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5050 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5060 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  73 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-73-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 74 Tasks

- [x] T021-5070 Confirm Slice 74 is bounded to numbered `cmd/sdp-trace`
  `pr-review check` one-shot workflow shards `pr_review_126` through
  `pr_review_136`.
- [x] T021-5071 Confirm Slice 74 is behavior-preserving: no changes to
  flag-only parsing, required `--out` and packet anchors, packet/profile/readiness
  failures mapping to `cannot_verify`, work-dir directory validation, repeated
  allowed-runner reconstruction from raw args, preview output as non-persisted
  planning data, run-set persistence before ledger and validation publication,
  summary text only after durable artifacts, validation verdict exit mapping
  through `exitCannotVerify`, package boundary, dependency direction, or
  baseline change is planned.
- [x] T021-5072 Record Slice 74 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-74-plan-review.md`.
- [x] T021-5080 Move `pr-review check` command orchestration, flag parsing and
  required inputs, prepare/execute orchestration, preview/summary publication,
  and durable artifact writes into cohesive locality files. Record the command
  and responsibility boundary in evidence, and synchronize
  `cmd/sdp-trace/FAMILY_INDEX.md` with the renamed files.
- [x] T021-5090 Run `gofmt` on changed Go files.
- [x] T021-5100 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-5101 Run focused `pr-review check` regression evidence covering
  exact test existence for `TestPRReviewCheckWritesRunProvenance`,
  `TestPRReviewCheckPreviewDoesNotWriteArtifacts`,
  `TestPRReviewCheckRequiresOutAndPacketInputs`, and
  `TestPRReviewCheckRejectsMissingOrFileWorkDir`; run-set artifact persistence,
  ledger/validation artifact publication, preview output not producing
  artifacts, required output and packet-input usage errors, work-dir failure
  mapping to `cannot_verify`, packet/profile/readiness failures mapping to
  `cannot_verify`, repeated allowed-runner reconstruction from raw args,
  validation verdict exit mapping through `exitCannotVerify`, exact test list
  verification after test implementation that fails when any planned focused
  test is missing, and no baseline changes.
- [x] T021-5110 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5120 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5130 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  74 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-74-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 75 Tasks

- [x] T021-5140 Confirm Slice 75 is bounded to the remaining numbered
  `cmd/sdp-trace` shared helper shards `pr_review_137`, `pr_review_139`,
  `pr_review_142`, and `pr_review_144` through `pr_review_149`, with
  `writeIndentedPayload` treated as a generic CLI JSON output helper because it
  is also used by protected gate output.
- [x] T021-5141 Confirm Slice 75 is behavior-preserving: no changes to stdout
  JSON rendering for `pr-review` and protected gate callers, write-once
  output-file refusal, work-dir diagnostics, validation exit mapping,
  packet/profile load order and partial-input avoidance, repeated raw flag order
  and fallback semantics, runner allow-list normalization, packet-dir
  derivation, package boundary, dependency direction, or baseline change is
  planned.
- [x] T021-5142 Record Slice 75 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-75-plan-review.md`.
- [x] T021-5150 Move shared helper shards into cohesive locality files for
  generic CLI JSON output, filesystem safety, validation exits, packet/profile
  inputs, repeated flag reconstruction, runner allow-lists, and packet-dir
  derivation. Record the shared-helper boundary in evidence, and synchronize
  `cmd/sdp-trace/FAMILY_INDEX.md` with the renamed files.
- [x] T021-5160 Run `gofmt` on changed Go files.
- [x] T021-5170 Run focused Go verification for `cmd/sdp-trace`.
- [x] T021-5171 Run focused shared-helper regression evidence covering exact
  test existence for `TestPRReviewSharedOutputAndFileHelpers`,
  `TestPRReviewSharedPacketProfileAndExitHelpers`,
  `TestPRReviewSharedRepeatedFlagsRunnerSetAndPacketDir`, and
  `TestSharedIndentedPayloadPreservesProtectedGateOutput`; stdout JSON
  rendering with a trailing newline for generic direct use and protected gate
  output, write-once output-file refusal for blank,
  existing file, and directory paths, new-file allowance, work-dir missing and
  file diagnostics, validation exit mapping for `cannot_verify`,
  `coverage_unresolved`, and satisfied states, packet/profile load error order,
  repeated raw flag reconstruction for `--key value` and `--key=value`,
  parsed fallback behavior, runner allow-list comma trimming and empty-entry
  ignoring, packet-dir derivation for directory, file, and missing path, exact
  test list verification after test implementation that fails when any planned
  focused test is missing, and no baseline changes.
- [x] T021-5180 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5190 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5200 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  75 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-75-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 76 Tasks

- [x] T021-5210 Confirm Slice 76 is bounded to numbered `internal/packet`
  packet contract data model and catalog shards `packet_001` through
  `packet_032`.
- [x] T021-5211 Confirm Slice 76 is behavior-preserving: no changes to
  schema-version values, required row ordering, required decision ordering,
  known state/catalog membership, JSON field names, `omitempty` behavior,
  GitHub PR evidence input shape, package boundary, dependency direction, or
  baseline change is planned.
- [x] T021-5212 Record Slice 76 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-76-plan-review.md`.
- [x] T021-5220 Move packet contract constants/catalogs, core packet types,
  bundle manifest types, GitHub source input types, and validation result types
  into cohesive locality files. Record the data-model boundary in evidence, and
  synchronize `cmd/sdp-trace/FAMILY_INDEX.md` if applicable; otherwise record
  why the command family index is not authoritative for `internal/packet`.
- [x] T021-5230 Run `gofmt` on changed Go files.
- [x] T021-5240 Run focused Go verification for `internal/packet`.
- [x] T021-5241 Run focused packet contract regression evidence covering exact
  test existence for `TestPacketContractCatalogsPreserveTrustSurface`,
  `TestPacketCoreTypesPreserveJSONShape`,
  `TestPacketBundleTypesPreserveJSONShape`, and
  `TestGitHubEvidenceInputTypesPreserveJSONShape`; schema-version constants,
  required row ordering, required decision ordering, state/catalog membership,
  packet/source/projection/row/finding/gap/decision owner JSON fields and
  `omitempty`, bundle manifest/entry/resolver JSON fields and `omitempty`,
  GitHub PR input/commit/check/artifact/review/prompt-boundary/
  prompt-boundary-classification/integration-action/build-result JSON fields
  and `omitempty`, validation JSON fields and `omitempty`, exact test list
  verification after test implementation that fails when any planned focused
  test is missing, and no baseline changes.
- [x] T021-5250 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5260 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5270 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  76 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-76-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 77 Tasks

- [x] T021-5280 Confirm Slice 77 is bounded to numbered `internal/packet`
  packet entrypoint and prompt-boundary behavior shards `packet_033` through
  `packet_046`.
- [x] T021-5281 Confirm Slice 77 is behavior-preserving: no changes to JSON
  loader errors, generated packet/bundle IDs, generated-at UTC formatting,
  default packet profile fields, integration-action extension behavior,
  prompt-boundary verdicts/reasons/route proof effects, contamination theater
  finding shape, bundle manifest digesting, package boundary, dependency
  direction, or baseline changes are planned.
- [x] T021-5282 Record Slice 77 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-77-plan-review.md`.
- [x] T021-5290 Move validator context structs, JSON loading entrypoints,
  GitHub bundle assembly entrypoint helpers, packet shell construction,
  prompt-boundary classification helpers, prompt-boundary theater finding
  append, bundle manifest construction, and recorder-duty phrase catalog into
  cohesive locality files. Record source-shape evidence that numbered
  `packet_033` through `packet_046` files are gone, that the replacement
  locality filenames exist, and that `cmd/sdp-trace/FAMILY_INDEX.md` is not
  applicable because this is non-command `internal/packet` package locality.
- [x] T021-5300 Run `gofmt` on changed Go files.
- [x] T021-5310 Run focused Go verification for `internal/packet`.
- [x] T021-5311 Run focused packet behavior regression evidence covering exact
  test existence for existing prompt-boundary classification and GitHub build
  tests; loader behavior for `LoadBundle` and `LoadGitHubInput`; generated
  packet/bundle IDs and UTC timestamp formatting; integration-action extension
  preservation; contaminated prompt theater finding shape; clean/missing/
  digest-only/malformed prompt-boundary classifications; forbidden phrase
  catalog membership; bundle manifest schema, bundle ID, digest, and entries;
  and no baseline changes.
- [x] T021-5312 Run focused source-shape evidence proving numbered
  `internal/packet/packet_033` through `internal/packet/packet_046` files no
  longer exist and the planned replacement locality files exist.
- [x] T021-5320 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5330 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5340 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  77 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-77-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 78 Tasks

- [x] T021-5350 Confirm Slice 78 is bounded to numbered `internal/packet`
  packet rendering and digest helper shards `packet_047` through `packet_059`.
- [x] T021-5351 Confirm Slice 78 is behavior-preserving: no changes to rendered
  section headers, table column order, markdown escaping for pipes and newlines,
  `none` rendering for blank cells, clean-theater fallback row shape, resolver
  fallback behavior, required row ordering semantics, non-approval fallback
  text, digest prefix/determinism, package boundary, dependency direction, or
  baseline changes are planned.
- [x] T021-5352 Record Slice 78 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-78-plan-review.md`.
- [x] T021-5360 Move packet rendering helpers, row/resolver lookup helpers,
  markdown escaping helper, required-row ordering helper, and packet digest
  helper into cohesive locality files. Record source-shape evidence that
  numbered `packet_047` through `packet_059` files are gone, that replacement
  locality filenames exist, and that `cmd/sdp-trace/FAMILY_INDEX.md` is not
  applicable because this is non-command `internal/packet` package locality.
- [x] T021-5370 Run `gofmt` on changed Go files.
- [x] T021-5380 Run focused Go verification for `internal/packet`.
- [x] T021-5381 Run focused packet rendering/digest regression evidence
  covering exact test existence for clean theater fallback rendering, theater
  finding table rendering, decision table rendering, evidence resolver fallback
  rendering, residual gap and no-gap rendering, non-proof fallback text,
  markdown cell escaping, row lookup missing fallback, required-row ordering
  lookup, packet digest prefix/determinism, top-level `RenderMarkdown` section
  order and section/table headers or source-shape evidence that excluded
  `packet_193` onward orchestration files were untouched, exact source-shape
  verification after implementation, and no baseline changes.
- [x] T021-5390 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5400 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5410 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  78 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-78-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 79 Tasks

- [x] T021-5420 Confirm Slice 79 is bounded to numbered `internal/packet`
  packet validation entrypoint and demo-first gate shards `packet_060` through
  `packet_084`.
- [x] T021-5421 Confirm Slice 79 is behavior-preserving: no changes to general
  `Validate` delegation through `bundleValidator`, demo gate accumulation of
  base validation errors, tool-generated authoring requirement, required
  pass-or-partial row count, required `PC-CHANGE` and `PC-MUTATION` retained
  usable evidence, retained structured OpenCode/GSD/MiniMax route evidence
  requirement, verification-or-review assessed semantics, cannot-verify closure
  cap, route component aliases, synthetic digest rejection, package boundary,
  dependency direction, or baseline changes are planned.
- [x] T021-5422 Record Slice 79 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-79-plan-review.md`.
- [x] T021-5430 Move packet validation entrypoint, demo-first gate entrypoint
  and orchestration, demo gate row/manifest indexing, demo row evidence checks,
  route evidence checks, closure checks, and demo evidence usability helpers
  into cohesive locality files. Record source-shape evidence that numbered
  `packet_060` through `packet_084` files are gone, that replacement locality
  filenames exist, and that `cmd/sdp-trace/FAMILY_INDEX.md` is not applicable
  because this is non-command `internal/packet` package locality.
- [x] T021-5440 Run `gofmt` on changed Go files.
- [x] T021-5450 Run focused Go verification for `internal/packet`.
- [x] T021-5451 Run focused packet validation/demo-gate regression evidence
  covering exact test existence for general `Validate` failure behavior,
  `CheckDemoFirstPacket` accumulation of base `Validate` errors plus
  demo-first errors in the same result,
  minimum four pass-or-partial rows including below-threshold failure,
  self-declared route evidence rejection, digest-bound route evidence
  acceptance, route component normalization, legacy and redux GSD route
  component aliases, all MiniMax route component aliases (`minimax`,
  `minimax-m2.5`, and `minimax-m2`), missing/expired row evidence usability
  failures for `PC-CHANGE` and `PC-MUTATION`,
  verification-or-review assessed semantics, cannot-verify closure cap and
  closure evidence exemption, route evidence source/kind/component/digest
  requirements, exact source-shape verification after implementation, and no
  baseline changes.
- [x] T021-5460 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5470 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5480 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  79 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-79-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 80 Tasks

- [x] T021-5490 Confirm Slice 80 is bounded to numbered `internal/packet`
  general packet validator orchestration, metadata validation, and
  manifest/resolver indexing shards `packet_085` through `packet_101`.
- [x] T021-5491 Confirm Slice 80 is behavior-preserving: no changes to
  `bundleValidator` validation phase order, error accumulation, schema-version
  diagnostics, packet/bundle identity and non-empty checks, packet digest
  required/mismatch behavior, `packet.non_approval` required-field behavior,
  packet state and authoring method catalog checks,
  canonical/non-canonical projection semantics, manifest entry empty-ref
  rejection, manifest retained-form and redaction-status enum diagnostics,
  manifest resolver fallback and resolver-entry override semantics, empty
  resolver-ref ignoring, package boundary, dependency direction, or baseline
  changes are planned.
- [x] T021-5492 Record Slice 80 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-80-plan-review.md`.
- [x] T021-5500 Move validator orchestration, metadata validation helpers,
  projection checks, manifest indexing, resolver indexing, and manifest entry
  enum validation into cohesive locality files. Record source-shape evidence
  that numbered `packet_085` through `packet_101` files are gone, that
  replacement locality filenames exist, and that `cmd/sdp-trace/FAMILY_INDEX.md`
  is not applicable because this is non-command `internal/packet` package
  locality. Record boundary evidence that excluded row validation files
  `packet_102` through `packet_121`, finding/gap/decision-owner and shared row
  lookup helpers `packet_122` through `packet_144`, and shared
  artifact-usability helpers `packet_145` onward are not moved or edited in
  Slice 80.
- [x] T021-5510 Run `gofmt` on changed Go files.
- [x] T021-5520 Run focused Go verification for `internal/packet`.
- [x] T021-5521 Run focused packet metadata/manifest regression evidence
  covering exact test existence for validation phase order and accumulated
  multi-error output across metadata, manifest, row, and finding/gap failures,
  schema-version diagnostics, missing packet and bundle identity fields,
  packet/bundle mismatch, missing and mismatched packet digest, missing
  `packet.non_approval`, packet state and authoring method catalog diagnostics,
  canonical projection kind rejection, non-canonical missing `artifact_ref`
  rejection, manifest entry empty-ref rejection, retained-form and
  redaction-status enum diagnostics, manifest entry resolver fallback,
  resolver-entry override, empty resolver-ref ignoring, exact source-shape
  verification after implementation, focused regression evidence that row
  evidence-ref validation still observes manifest resolver fallback,
  resolver-entry override, and empty resolver-ref ignoring, and no baseline
  changes.
- [x] T021-5530 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5540 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5550 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  80 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-80-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 81 Tasks

- [x] T021-5560 Confirm Slice 81 is bounded to numbered `internal/packet`
  row validation, contradiction attribution, and row/ref/gap lookup shards
  `packet_102` through `packet_121` and `packet_138` through `packet_144`.
- [x] T021-5561 Confirm Slice 81 is behavior-preserving: no changes to row
  validation order, required row presence checks, unknown and duplicate row ID
  diagnostics, row required field diagnostics, missing reason behavior for
  non-pass states, pass-row retained evidence requirements, evidence-ref
  absence and resolver diagnostics, pass-row expired/unverifiable artifact
  diagnostics, contradiction target selection, required-row precedence and
  extension-row fallback ordering for evidence-ref attribution, exact row-ref
  matching, contradiction partial-state and residual-gap requirements,
  gap reason semantics, package boundary, dependency direction, or baseline
  changes are planned.
- [x] T021-5562 Record Slice 81 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-81-plan-review.md`.
- [x] T021-5570 Move row validation orchestration, row indexing and ID checks,
  row field/reason/evidence validation, evidence-ref validation handoff,
  pass-evidence usability handoff, contradiction validation, row/ref lookup,
  and gap lookup helpers into cohesive locality files. Record source-shape
  evidence that numbered `packet_102` through `packet_121` and `packet_138`
  through `packet_144` files are gone, that replacement locality filenames
  exist, and that `cmd/sdp-trace/FAMILY_INDEX.md` is not applicable because
  this is non-command `internal/packet` package locality. Record staged
  boundary evidence that excluded finding/gap/decision-owner validation
  `packet_122` through `packet_136`, shared validator error accumulation
  `packet_137`, shared artifact-usability helpers
  `packet_145` through `packet_147`, GitHub bundle construction `packet_148`
  through `packet_192`, and rendering `packet_193` onward are not moved or edited in
  Slice 81.
- [x] T021-5580 Run `gofmt` on changed Go files.
- [x] T021-5590 Run focused Go verification for `internal/packet`.
- [x] T021-5591 Run focused row validation/contradiction regression evidence
  covering exact test existence for validation order and accumulated row
  diagnostics, missing required rows, unknown row IDs, duplicate row IDs, row
  state/summary/owner/reason diagnostics, pass rows requiring retained evidence
  refs, evidence ref absent-from-manifest and no-resolver diagnostics, pass-row
  expired and unverifiable artifact diagnostics, contradiction target selection
  by explicit row ID including explicit `ContradictsRowID` overriding a
  different evidence-ref fallback target, required-row ordering among shared
  refs where the first `RequiredRows` match wins, required-row precedence over
  extension rows for shared refs, deterministic extension-row fallback
  ordering, exact row-ref matching without prefix/alias inference,
  contradiction partial-state requirement, contradiction residual-gap
  requirement, `gapForRow` requiring a non-empty reason, excluded
  `validateResidualCoverage` behavior still observing `gapForRow` semantics
  for blank and non-empty residual-gap reasons, exact source-shape verification
  after implementation, staged boundary evidence for excluded packet files, and
  no baseline changes.
- [x] T021-5600 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5610 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5620 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  81 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-81-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 82 Tasks

- [x] T021-5630 Confirm Slice 82 is bounded to numbered `internal/packet`
  finding/gap/decision-owner validation and shared validator error accumulation
  shards `packet_122` through `packet_137`.
- [x] T021-5631 Confirm Slice 82 is behavior-preserving: no changes to
  validation phase order, theater finding reason-code diagnostics, trigger
  evidence-ref validation handoff, residual gap unknown-row and missing-reason
  diagnostics, residual coverage exemptions, allowed theater states with
  findings, required decision owner ordering, decision owner trimming,
  duplicate decision-owner last-valid-owner-wins overwrite behavior, decision
  owner state/reason diagnostics, accumulated error formatting, package
  boundary, dependency direction, or baseline changes are planned.
- [x] T021-5632 Record Slice 82 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-82-plan-review.md`.
- [x] T021-5640 Move theater finding validation, theater row state validation,
  residual gap validation, residual coverage enforcement, decision-owner
  validation, and validator error accumulation into cohesive locality files.
  Record source-shape evidence that numbered `packet_122` through `packet_137`
  files are gone, that replacement locality filenames exist, and that
  `cmd/sdp-trace/FAMILY_INDEX.md` is not applicable because this is non-command
  `internal/packet` package locality. Record staged boundary evidence that
  excluded shared artifact-usability helpers `packet_145` through `packet_147`,
  GitHub bundle construction `packet_148` through `packet_192`, rendering
  `packet_193` onward, and `internal/prreview` numbered files are not moved or
  edited in Slice 82.
- [x] T021-5650 Run `gofmt` on changed Go files.
- [x] T021-5660 Run focused Go verification for `internal/packet`.
- [x] T021-5661 Run focused finding/gap/decision-owner regression evidence
  covering exact named regression tests:
  `TestValidateFindingsGapsDecisionOwnersDiagnosticsInOrder`,
  `TestValidateTheaterFindingDiagnostics`,
  `TestValidateTheaterFindingStatesWithFindings`,
  `TestValidateResidualGapDiagnosticsAndCoverage`,
  `TestValidateDecisionOwnerDiagnosticsAndDuplicateOverwrite`, and
  `TestBundleValidatorAddFormatsErrors`. These tests must cover validation
  phase accumulation with finding/gap/decision-owner errors, including exact
  diagnostic strings and relative `assertErrorsInOrder` ordering for theater
  state diagnostics, required decision owner diagnostics in `requiredDecisions`
  order, theater finding evidence-ref diagnostics, and residual gap
  diagnostics. Also cover theater finding missing and unknown reason-code
  diagnostics, trigger evidence-ref manifest/resolver diagnostics, residual gap
  unknown row and missing reason diagnostics, residual coverage exemption for
  `PC-RESIDUAL-GAPS` and pass rows, missing coverage for non-pass rows,
  theater row pass rejection and allowed partial/fail/cannot_verify states when
  findings are present, missing required decision owner diagnostics, blank
  decision owner decision handling, owner name requirement, duplicate
  decision-owner last-valid-owner-wins overwrite behavior, decision owner
  unknown state and missing non-pass reason diagnostics, accumulated error
  formatting through `bundleValidator.add`, exact source-shape verification
  after implementation, staged boundary evidence for excluded packet/prreview
  files, and no baseline changes.
- [x] T021-5670 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5680 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5690 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  82 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-82-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 83 Tasks

- [x] T021-5700 Confirm Slice 83 is bounded to numbered `internal/packet`
  shared artifact usability helper shards `packet_145` through `packet_147`.
- [x] T021-5701 Confirm Slice 83 is behavior-preserving: no changes to blank
  and whitespace-only `expires_at` non-expiry behavior, malformed timestamp
  expiry behavior, RFC3339 expiry comparison semantics,
  `redaction_status=cannot_verify` and `retained_form=not_retained`
  unverifiable behavior, artifact access classification for `expired`,
  `inaccessible`, `malformed`, `not_assessed`, and `cannot_verify`, default
  artifact access pass-through behavior, demo-first row evidence and retained
  route evidence usability semantics, package boundary, dependency direction,
  or baseline changes are planned.
- [x] T021-5702 Record Slice 83 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-83-plan-review.md`.
- [x] T021-5710 Move artifact usability helpers into a cohesive artifact
  evidence usability locality without creating standalone one-function
  replacement files. Record source-shape evidence that numbered `packet_145`
  through `packet_147` files are gone, that `artifact_evidence_usability.go`
  contains the shared usability helpers, and that `cmd/sdp-trace/FAMILY_INDEX.md`
  is not applicable because this is non-command `internal/packet` package
  locality. Record staged
  boundary evidence that excluded GitHub bundle construction `packet_148`
  through `packet_192`, rendering `packet_193` onward, and `internal/prreview`
  numbered files are not moved or edited in Slice 83.
- [x] T021-5720 Run `gofmt` on changed Go files.
- [x] T021-5730 Run focused Go verification for `internal/packet`.
- [x] T021-5731 Run focused artifact-usability regression evidence covering
  exact named tests `TestEntryExpiredSemantics`,
  `TestPassRefUnverifiableSemantics`, and
  `TestDemoFirstEvidenceUsabilityStillUsesArtifactHelpers`. Run them with a
  focused command that fails on zero matches. These tests must cover blank and
  whitespace-only `ExpiresAt` as not expired, malformed `ExpiresAt` as expired,
  equal-to-now and before-now as expired, after-now as not expired,
  `RedactionStatus=StateCannotVerify`, `RetainedForm=not_retained`, artifact
  access values `expired`, `inaccessible`, `malformed`, `not_assessed`, and
  `StateCannotVerify` as unverifiable, empty/present/unknown artifact access as
  not unverifiable, and demo-first row evidence plus retained route evidence
  rejecting both expired and unverifiable refs through `demoUsableEntry` and
  the shared `entryExpired`/`passRefUnverifiable` helpers.
- [x] T021-5740 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5750 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5760 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  83 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-83-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 84 Tasks

- [x] T021-5770 Confirm Slice 84 is bounded to numbered `internal/packet`
  GitHub source-change and row-construction shards `packet_148` through
  `packet_173`.
- [x] T021-5771 Confirm Slice 84 is behavior-preserving: no changes to
  source-change fields, row IDs, states, summaries, evidence refs, reasons,
  owner assignment, prompt-boundary route/verification classification behavior,
  workflow-run wording, artifact-ref de-duplication, retained artifact
  filtering, review pass/partial/not-assessed behavior, package boundary,
  dependency direction, or baseline changes are planned.
- [x] T021-5772 Record Slice 84 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-84-plan-review.md`.
- [x] T021-5780 Move GitHub source-change and row-construction helpers into
  cohesive responsibility-named files without creating standalone one-function
  replacement files. Record source-shape evidence that numbered `packet_148`
  through `packet_173` files are gone and that `cmd/sdp-trace/FAMILY_INDEX.md`
  is not applicable because this is non-command `internal/packet` package
  locality. Record staged boundary evidence that excluded GitHub manifest entry
  helpers `packet_174` through `packet_192`, rendering `packet_193` onward, and
  `internal/prreview` numbered files are not moved or edited in Slice 84.
- [x] T021-5790 Run `gofmt` on changed Go files.
- [x] T021-5800 Run focused Go verification with `go test ./internal/packet`.
- [x] T021-5801 Run focused GitHub row/source regression evidence covering
  exact named tests `TestGitHubSourceChangeProjection`,
  `TestGitHubChangeAndMutationRowsPreserveCommitRangeSemantics`,
  `TestGitHubInitiatorAndReviewRowsPreserveEvidenceSemantics`,
  `TestGitHubAgentRouteRowsPreservePromptBoundarySemantics`,
  `TestGitHubVerificationRowsPreserveEvidenceSemantics`, and
  `TestGitHubArtifactEvidenceHelpersPreserveSemantics`. Run them with a
  `go test -list` exact-count guard before the focused `go test -run` command
  so zero matches fail. These tests must prove source-change projection,
  change/mutation missing commit-range behavior, initiator row behavior,
  prompt-boundary route fail/cannot-verify/partial behavior, verification
  no-checks/missing-workflow/artifact-binding/non-success/pass behavior,
  artifact-ref de-duplication, retained artifact filtering, and review
  not-assessed/partial/pass behavior.
- [x] T021-5810 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5820 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5830 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  84 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-84-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 85 Tasks

- [x] T021-5840 Confirm Slice 85 is bounded to numbered `internal/packet`
  GitHub manifest/default/digest helper shards `packet_174` through
  `packet_192`.
- [x] T021-5841 Confirm Slice 85 is behavior-preserving: no changes to manifest
  entry refs, source classes, resolvers, retained forms, authority metadata,
  source refs, prompt-boundary resolver/retained-form behavior, agent route
  digest/evidence-kind/component behavior, artifact expiry/digest propagation,
  integration entry actor/authority behavior, residual gap filtering, default
  decision owner states/reasons, redaction markers, digest format, package
  boundary, dependency direction, or baseline changes are planned.
- [x] T021-5842 Record Slice 85 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-85-plan-review.md`.
- [x] T021-5850 Move GitHub manifest/default/digest helpers into cohesive
  responsibility-named files without creating standalone one-function
  replacement files. Record source-shape evidence that numbered `packet_174`
  through `packet_192` files are gone and that `cmd/sdp-trace/FAMILY_INDEX.md`
  is not applicable because this is non-command `internal/packet` package
  locality. Record staged boundary evidence that excluded rendering
  `packet_193` onward and `internal/prreview` numbered files are not moved or
  edited in Slice 85.
- [x] T021-5860 Run `gofmt` on changed Go files.
- [x] T021-5870 Run focused Go verification with `go test ./internal/packet`.
- [x] T021-5871 Run focused GitHub manifest/default/digest regression evidence
  covering exact named tests `TestGitHubManifestEntriesPreserveAssemblyOrder`,
  `TestGitHubConditionalManifestEntriesPreserveSemantics`,
  `TestGitHubBundleEntryAuthorityAndRedactionPreserveSemantics`,
  `TestGitHubResidualGapsAndDecisionOwnersPreserveDefaults`, and
  `TestGitHubResolverAndDigestHelpersPreserveSemantics`. Run them with a
  `go test -list` exact-count guard before the focused `go test -run` command
  so zero matches fail. These tests must prove manifest entry order and refs,
  conditional prompt boundary / PR body / agent route / checks / reviews /
  artifacts / integration entries, authority metadata and source refs,
  bundle-entry defaults and resolver redaction, residual gap filtering,
  decision owner defaults, resolver joining, and digest placeholder format.
- [x] T021-5880 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5890 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5900 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  85 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-85-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 86 Tasks

- [x] T021-5910 Confirm Slice 86 is bounded to numbered `internal/packet`
  Markdown rendering shards `packet_193` through `packet_200`.
- [x] T021-5911 Confirm Slice 86 is behavior-preserving: no changes to
  validation-before-render behavior, error joining, top-level section order,
  packet disclaimer text, metadata field order, Markdown escaping, required-row
  sorting by `RequiredRows`, gap fallback wording, theater clean-row fallback,
  theater finding rendering, package boundary, dependency direction, or baseline
  changes are planned.
- [x] T021-5912 Record Slice 86 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-86-plan-review.md`.
- [x] T021-5920 Move packet Markdown rendering helpers into cohesive
  responsibility-named files without creating standalone one-function
  replacement files. Record source-shape evidence that numbered `packet_193`
  through `packet_200` files are gone and that `cmd/sdp-trace/FAMILY_INDEX.md`
  is not applicable because this is non-command `internal/packet` package
  locality. Record staged boundary evidence that excluded `internal/prreview`
  numbered files are not moved or edited in Slice 86.
- [x] T021-5930 Run `gofmt` on changed Go files.
- [x] T021-5940 Run focused Go verification with `go test ./internal/packet`.
- [x] T021-5941 Run focused packet rendering regression evidence covering exact
  named tests `TestValidateAndRenderHappyPath`,
  `TestRenderCleanTheaterUsesRowState`,
  `TestPacketRenderingHelpersPreserveTables`,
  `TestPacketRenderingSectionHelpersPreserveMetadataRowsAndErrors`,
  `TestRenderMarkdownPreservesTopLevelOrderAndHeaders`, and
  `TestRenderMarkdownRejectsInvalidBundleBeforeProjection`. Run them with a
  `go test -list` exact-count guard before the focused `go test -run` command
  so zero matches fail. These tests must prove validation before projection,
  joined validation errors, top-level order, header/table preservation, packet
  metadata field order, required-row ordering by `RequiredRows`, empty gap
  fallback to `none`, clean-theater fallback state, theater finding rendering,
  evidence table rendering, and non-proof language.
- [x] T021-5950 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-5960 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-5970 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  86 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-86-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 87 Tasks

- [x] T021-5980 Confirm Slice 87 is bounded to numbered `internal/prreview`
  portable schema/type shards `prreview_001` through `prreview_019`.
- [x] T021-5981 Confirm Slice 87 is behavior-preserving: no changes to
  exported names, JSON tags, trust-boundary comments, schema/state/ref/content/
  redaction/plane/runner/status/severity/disposition/coverage/authority
  constants, package regex patterns, error values, package boundary, dependency
  direction, or baseline changes are planned.
- [x] T021-5982 Record Slice 87 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-87-plan-review.md`.
- [x] T021-5990 Move portable prreview constants, vars, option structs, packet
  refs, profile/run/result structs, and ledger/validation structs into cohesive
  responsibility-named files without creating standalone one-function or
  one-type replacement files. Record source-shape evidence that numbered
  `prreview_001` through `prreview_019` files are gone and staged boundary
  evidence that excluded `prreview_020` onward files are not moved or edited in
  Slice 87.
- [x] T021-6000 Run `gofmt` on changed Go files.
- [x] T021-6010 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6011 Run focused portable schema/type regression evidence covering
  exact named tests `TestPrreviewSchemaConstantsAndPatternsPreserveContracts`
  and `TestPrreviewPortableTypesPreserveJSONShape`. Run them with a `go test
  -list` exact-count guard before the focused `go test -run` command so zero
  matches fail. These tests must prove schema/state/ref/content/redaction/plane/
  runner/status/severity/disposition/coverage/authority constants, package regex
  patterns, error values, exported portable type JSON keys, omitted optional
  fields, and trust-boundary structs remain stable.
- [x] T021-6020 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6030 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6040 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  87 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-87-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 88 Tasks

- [x] T021-6050 Confirm Slice 88 is bounded to numbered `internal/prreview`
  packet construction and packet input reference shards `prreview_020` through
  `prreview_037`.
- [x] T021-6051 Confirm Slice 88 is behavior-preserving: no changes to
  validation-before-directory-creation behavior, new-output-directory
  enforcement, packet digest/write behavior, packet ID format, default
  `CreatedBy` and CI state, copied ref kinds/content types, context kind
  inference, verification ref kind rewriting, metadata optionality,
  unavailable-field state/reason strings, package boundary, dependency
  direction, or baseline changes are planned.
- [x] T021-6052 Record Slice 88 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-88-plan-review.md`.
- [x] T021-6060 Move packet construction and packet input reference helpers into
  cohesive responsibility-named files without creating standalone one-function
  replacement files. Record source-shape evidence that numbered `prreview_020`
  through `prreview_037` files are gone and staged boundary evidence that
  excluded `prreview_038` onward numbered files are not moved or edited in
  Slice 88.
- [x] T021-6070 Run `gofmt` on changed Go files.
- [x] T021-6080 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6081 Run focused packet-build regression evidence covering exact
  named tests `TestBuildPacketBindsRefsAndRejectsUnsafeIdentity`,
  `TestBuildPacketRecordsUnavailableInputsAndDigestChangesWithDiff`, and
  `TestPrreviewPacketBuildHelpersPreserveDefaultsRefsAndUnavailableFields`.
  Run them with a `go test -list` exact-count guard before the focused `go test
  -run` command so zero matches fail. These tests must prove validation before
  output directory creation, new-output-directory enforcement without
  overwriting existing contents, packet ID/digest/write behavior, default
  `CreatedBy` and CI state, copied diff/metadata/context/verification ref kinds
  and content types, metadata optionality, context kind inference, verification
  kind rewriting, unavailable-field states/reasons, and digest changes when the
  diff input changes.
- [x] T021-6090 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6100 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6110 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  88 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-88-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 89 Tasks

- [x] T021-6120 Confirm Slice 89 is bounded to numbered `internal/prreview`
  review run orchestration shards `prreview_038` through `prreview_042`.
- [x] T021-6121 Confirm Slice 89 is behavior-preserving: no changes to profile
  validation before option normalization, default `Now`/`WorkDir` behavior,
  preview-only behavior without directory creation or runner invocation,
  new-output-directory enforcement, raw subdirectory creation mode, result
  ordering, packet digest propagation, `results.json` path/shape, error
  propagation, package boundary, dependency direction, or baseline changes are
  planned.
- [x] T021-6122 Record Slice 89 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-89-plan-review.md`.
- [x] T021-6130 Move review run orchestration helpers into cohesive
  responsibility-named files without creating standalone one-function
  replacement files. Record source-shape evidence that numbered `prreview_038`
  through `prreview_042` files are gone and staged boundary evidence that
  excluded ledger/validation/role-execution numbered files are not moved or
  edited in Slice 89.
- [x] T021-6140 Run `gofmt` on changed Go files.
- [x] T021-6150 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6151 Run focused review-run regression evidence covering exact named
  tests `TestRunReviewArtifactPipelineRedactsUnsafeReviewerText`,
  `TestRunReviewRecordsRunnerFailureStatesAndPromptDigest`,
  `TestRunReviewPreviewReturnsPreviewOnly`,
  `TestRunReviewPreservesValidationDefaultsAndOutputContracts`,
  `TestRunReviewNotAssessedReasonDoesNotInvokeRunner`,
  `TestRunReviewCannotVerifyUnreadablePromptTemplate`,
  `TestRunReviewPromptIncludesPacketEvidence`, and
  `TestRunReviewMapsTimeoutToTimedOut`. Run them with a `go test -list`
  exact-count guard before the focused `go test -run` command so zero matches
  fail. These tests must cover run orchestration, preview short-circuit,
  role-result ordering, output/raw artifact writing, runner-failure states,
  prompt digest/evidence propagation, not-assessed bypass behavior, and
  timeout mapping. They must also prove profile validation happens before
  output directory creation, default `Now`/`WorkDir` behavior is applied,
  existing output directories are rejected without overwriting contents,
  `results.json` is written at the expected path with the expected packet
  digest/results shape, and role execution errors propagate without silently
  writing a successful results artifact.
- [x] T021-6160 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6170 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6180 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  89 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-89-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 90 Tasks

- [x] T021-6190 Confirm Slice 90 is bounded to numbered `internal/prreview`
  ledger synthesis shards `prreview_043` through `prreview_048`.
- [x] T021-6191 Confirm Slice 90 is behavior-preserving: no changes to
  reviewer-run/result separation, packet digest propagation, ledger schema
  version, finding sorting by ID, existing-ledger disposition carry-forward,
  duplicate-ID last-existing-wins lookup behavior, generated fallback finding
  IDs, severity normalization, summary sanitization, evidence-ref preservation,
  empty-ledger behavior, package boundary, dependency direction, or baseline
  changes are planned.
- [x] T021-6192 Record Slice 90 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-90-plan-review.md`.
- [x] T021-6200 Move ledger synthesis helpers into cohesive
  responsibility-named files without creating standalone one-function
  replacement files. Record source-shape evidence that numbered `prreview_043`
  through `prreview_048` files are gone and staged boundary evidence that
  excluded validation, summary/rendering, IO, sanitizer/default-disposition,
  and role-execution numbered files are not moved or edited in Slice 90.
- [x] T021-6210 Run `gofmt` on changed Go files.
- [x] T021-6220 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6221 Run focused ledger-synthesis regression evidence covering exact
  named tests `TestValidationAndLedgerLifecyclePreserveTrustSemantics`,
  `TestLedgerDispositionCarryForward`, and
  `TestPrreviewLedgerSynthesisPreservesOrderingCarryForwardAndSanitization`.
  Run them with a `go test -list` exact-count guard before the focused `go test
  -run` command so zero matches fail. These tests must cover packet digest and
  schema propagation, finding flattening from reviewer results, stable sorting
  by ID, fallback finding IDs, prior disposition carry-forward, default
  disposition selection, duplicate-ID existing ledger lookup behavior, severity
  normalization, summary sanitization, evidence-ref preservation, and empty
  ledger behavior.
- [x] T021-6230 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6240 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6250 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  90 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-90-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 91 Tasks

- [x] T021-6260 Confirm Slice 91 is bounded to numbered `internal/prreview`
  validation orchestration, digest validation, and required-plane selection
  shards `prreview_049` through `prreview_060`.
- [x] T021-6261 Confirm Slice 91 is behavior-preserving: no changes to
  authority-boundary defaults, CI state propagation, unique reason/next-action
  normalization, packet/run/ledger digest mismatch reasons, required-plane
  not-assessed fallback, best usable result selection across retries,
  required-plane sorting by plane name, plane validation note aggregation,
  usable-count accumulation, package boundary, dependency direction, or
  baseline changes are planned.
- [x] T021-6262 Record Slice 91 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-91-plan-review.md`.
- [x] T021-6270 Move validation orchestration and required-plane selection
  helpers into cohesive responsibility-named files without creating standalone
  one-function replacement files. Record source-shape evidence that numbered
  `prreview_049` through `prreview_060` files are gone and staged boundary
  evidence that excluded plane-ranking/model-identity, ledger-finding
  validation, coverage-state, summary/rendering, IO, prompt, and reviewer
  execution numbered files are not moved or edited in Slice 91.
- [x] T021-6280 Run `gofmt` on changed Go files.
- [x] T021-6290 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6291 Run focused validation-orchestration regression evidence
  covering exact named tests `TestValidateReviewStatesAndAuthorityBoundary`,
  `TestValidationAndLedgerLifecyclePreserveTrustSemantics`,
  `TestValidateUsesBestPlaneResultAcrossRetries`,
  `TestValidateCannotVerifyPerResultPacketDigestMismatch`,
  `TestValidateCoverageStatesForNoReviewersUnresolvedAndStaleDigest`, and
  `TestPrreviewValidationOrchestrationPreservesDigestRequiredPlaneAndAuthorityContracts`.
  Run them with a `go test -list` exact-count guard before the focused `go test
  -run` command so zero matches fail. These tests must cover validation result
  authority defaults, CI state propagation, unique reason/next-action
  normalization, packet/run/ledger digest mismatch reasons, required-plane
  fallback and sorting, best usable result selection, plane note aggregation,
  usable-count impact on coverage, and stale per-result digest handling.
- [x] T021-6300 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6310 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6320 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  91 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-91-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 92 Tasks

- [x] T021-6330 Confirm Slice 92 is bounded to numbered `internal/prreview`
  validation ranking, model-identity enforcement, and ledger-finding
  validation shards `prreview_061` through `prreview_068`.
- [x] T021-6331 Confirm Slice 92 is behavior-preserving: no changes to
  findings-over-clean usable result ranking, cannot-verify status ordering,
  not-assessed floor semantics, model identity mismatch status/reason/
  next-action behavior, critical/major unresolved blocker detection, summary
  sanitization, finding citation resolvability, package boundary, dependency
  direction, or baseline changes are planned.
- [x] T021-6332 Record Slice 92 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-92-plan-review.md`.
- [x] T021-6340 Move validation ranking, model-identity, and ledger-finding
  validation helpers into cohesive responsibility-named files without creating
  standalone one-function replacement files. Record source-shape evidence that
  numbered `prreview_061` through `prreview_068` files are gone and staged
  boundary evidence that excluded coverage-state, summary/rendering, IO,
  prompt, reviewer execution, and packet option validation numbered files are
  not moved or edited in Slice 92.
- [x] T021-6350 Run `gofmt` on changed Go files.
- [x] T021-6360 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6361 Run focused validation-ranking and ledger-finding regression
  evidence covering exact named tests
  `TestValidateReviewStatesAndAuthorityBoundary`,
  `TestValidationAndLedgerLifecyclePreserveTrustSemantics`,
  `TestValidateUsesBestPlaneResultAcrossRetries`,
  `TestValidateCoverageStatesForNoReviewersUnresolvedAndStaleDigest`, and
  `TestPrreviewValidationRankingModelAndLedgerFindingContracts`. Run them with
  a `go test -list` exact-count guard before the focused `go test -run`
  command so zero matches fail. These tests must cover usable/non-usable
  ranking, cannot-verify status classification, model identity mismatch
  projection, unresolved critical/major finding detection, summary
  sanitization, citation resolvability, validation reason accumulation, and
  coverage impact.
- [x] T021-6370 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6380 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6390 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  92 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-92-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 93 Tasks

- [x] T021-6400 Confirm Slice 93 is bounded to numbered `internal/prreview`
  coverage-state, model-identity support, and summary rendering shards
  `prreview_069` through `prreview_078`.
- [x] T021-6401 Confirm Slice 93 is behavior-preserving: no changes to
  cannot-verify precedence, no-review not-assessed behavior, partial coverage
  behavior, unresolved critical/major coverage behavior, satisfied coverage
  behavior, model mismatch fallback semantics, authority-boundary summary
  language, safe text handling in rendered next actions and findings, package
  boundary, dependency direction, or baseline changes are planned.
- [x] T021-6402 Record Slice 93 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-93-plan-review.md`.
- [x] T021-6410 Move coverage-state, model-identity support, and summary
  rendering helpers into cohesive responsibility-named files without creating
  standalone one-function replacement files. Record source-shape evidence that
  numbered `prreview_069` through `prreview_078` files are gone and staged
  boundary evidence that excluded IO, option validation, reviewer execution,
  prompt, sanitizer, and citation numbered files are not moved or edited in
  Slice 93.
- [x] T021-6420 Run `gofmt` on changed Go files.
- [x] T021-6430 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6431 Run focused coverage/model/summary regression evidence covering
  exact named tests `TestSynthesizeAndValidateCoverageSatisfied`,
  `TestValidateCannotVerifyUnexplainedModelMismatch`,
  `TestValidateCoverageStatesForNoReviewersUnresolvedAndStaleDigest`,
  `TestValidationAndSummaryRedactUnsafeMarkerClasses`, and
  `TestPrreviewCoverageModelAndSummaryRenderingContracts`. Run them with a
  `go test -list` exact-count guard before the focused `go test -run` command
  so zero matches fail. These tests must cover coverage-state precedence,
  not-assessed/partial/unresolved/satisfied states, model identity mismatch and
  fallback metadata semantics, authority-boundary summary language, safe next
  action rendering, safe finding rendering, and no approval overclaiming in
  summaries.
- [x] T021-6440 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6450 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6460 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  93 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-93-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 94 Tasks

- [x] T021-6470 Confirm Slice 94 is bounded to numbered `internal/prreview`
  JSON artifact IO and typed artifact reader shards `prreview_079` through
  `prreview_086`.
- [x] T021-6471 Confirm Slice 94 is behavior-preserving: no changes to JSON
  indentation/trailing newline behavior, parent-directory creation mode, file
  write mode, blank `WriteJSON` path semantics, packet/run-set directory
  default paths, profile validation, run-set validation, decode/read error
  propagation, package boundary, dependency direction, or baseline changes are
  planned.
- [x] T021-6472 Record Slice 94 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-94-plan-review.md`.
- [x] T021-6480 Move JSON artifact IO and typed artifact readers into cohesive
  responsibility-named files without creating standalone one-function
  replacement files. Record source-shape evidence that numbered `prreview_079`
  through `prreview_086` files are gone and staged boundary evidence that
  excluded option validation, reviewer execution, prompt, sanitizer, citation,
  validation, and summary numbered files are not moved or edited in Slice 94.
- [x] T021-6490 Run `gofmt` on changed Go files.
- [x] T021-6500 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6501 Run focused IO regression evidence covering exact named tests
  `TestWriteJSONAndReadRunSetUseDirectoryContracts`,
  `TestPacketProfileAndSmallHelpers`, `TestReadRunSetRejectsDuplicateRunIDs`,
  and `TestPrreviewArtifactIOReadWriteContracts`. Run them with a `go test
  -list` exact-count guard before the focused `go test -run` command so zero
  matches fail. These tests must cover blank write path no-op, parent directory
  creation with `0o755` mode where mode bits are exposed, JSON file writes
  with `0o644` mode where mode bits are exposed, indented JSON with trailing
  newline, packet and run-set directory default paths, typed ledger and
  validation readers, invalid JSON propagation, missing file propagation,
  profile validation, and duplicate run-set rejection.
- [x] T021-6510 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6520 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6530 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  94 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-94-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 95 Tasks

- [x] T021-6540 Confirm Slice 95 is bounded to numbered `internal/prreview`
  packet option, run-set, and review profile validation shards `prreview_087`
  through `prreview_099`.
- [x] T021-6541 Confirm Slice 95 is behavior-preserving: no changes to
  validation before packet directory creation, exact packet/profile/run-set
  validation error messages, duplicate review-run ID rejection, schema-version
  optionality and mismatch errors, required profile field errors, role field and
  runner errors, required-plane-without-role errors, package boundary,
  dependency direction, or baseline changes are planned.
- [x] T021-6542 Record Slice 95 plan/task review in
  `specs/021-source-file-locality-cleanup/reviews/slice-95-plan-review.md`.
- [x] T021-6550 Move packet option, run-set, and profile validation helpers
  into cohesive responsibility-named files without creating standalone
  one-function replacement files. Record source-shape evidence that numbered
  `prreview_087` through `prreview_099` files are gone and staged boundary
  evidence that excluded role execution, runner command handling, prompt, JSON
  artifact IO, validation/summary, sanitizer, citation, and packet construction
  numbered files are not moved or edited in Slice 95.
- [x] T021-6560 Run `gofmt` on changed Go files.
- [x] T021-6570 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6571 Run focused validation regression evidence covering exact named
  tests `TestBuildPacketBindsRefsAndRejectsUnsafeIdentity`,
  `TestValidateProfileRejectsMalformedProfiles`,
  `TestReadRunSetRejectsDuplicateRunIDs`,
  `TestRunReviewPreservesValidationDefaultsAndOutputContracts`, and
  `TestPrreviewOptionSchemaValidationContracts`. Run them with a `go test
  -list` exact-count guard before the focused `go test -run` command so zero
  matches fail. These tests must cover packet out/diff requirements, repo ID
  and change-ref patterns, commit SHA validation, CI state validation, run-set
  missing and duplicate review IDs, profile schema optionality and mismatch,
  profile required fields, role required fields, invalid runner rejection,
  required-plane-to-role mapping, and validation-before-output-directory
  creation. The focused Slice 95 regression must assert exact error strings for
  packet out/diff, repo ID, change-ref, commit SHA, CI state, run-set missing
  and duplicate review IDs, profile schema mismatch and optionality, required
  profile fields, role field and runner errors, and required-plane-without-role
  errors.
- [x] T021-6580 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6590 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6600 Run three independent reviewer lanes, fix every finding, repeat
  affected lanes until each reviewer returns exactly `LGTM`, and record Slice
  95 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-95-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 96 Tasks

- [x] T021-6610 Confirm Slice 96 is bounded to numbered `internal/prreview`
  role execution and OpenCode baseline/mutation guard shards `prreview_100`
  through `prreview_124`.
- [x] T021-6611 Confirm Slice 96 is behavior-preserving: no changes to role
  order, timeout handling, raw result writing, exact status and reason strings
  for runner timeout/empty output/parse failure/unavailable/failed, prompt
  evidence `cannot_verify` states, disallowed runner errors, not-assessed
  override behavior, OpenCode read-only enforcement, dirty baseline handling,
  mutation detection, command digest behavior, package boundary, dependency
  direction, or baseline changes are planned.
- [x] T021-6612 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-96-plan-review.md`.
- [x] T021-6620 Move role execution, result completion, runner preparation,
  command/error execution, OpenCode baseline guard, and working-tree baseline
  helpers into cohesive responsibility-named files without creating standalone
  one-function replacement files. Record source-shape evidence that numbered
  `prreview_100` through `prreview_124` files are gone and staged boundary
  evidence that excluded parser/raw-result writing `prreview_125` through
  `prreview_133`, preview, packet copying, prompt rendering/sanitization/
  citation, validation/summary, and previously consolidated validation files are
  not moved or edited in Slice 96.
- [x] T021-6630 Run `gofmt` on changed Go files.
- [x] T021-6640 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6641 Run focused role execution regression evidence covering exact
  named tests `TestRunReviewRecordsRunnerFailureStatesAndPromptDigest`,
  `TestRunReviewPreservesValidationDefaultsAndOutputContracts`,
  `TestApplyRunnerErrorClassifiesUnavailableAndFailure`,
  `TestRunReviewNotAssessedReasonDoesNotInvokeRunner`,
  `TestRunReviewCannotVerifyUnreadablePromptTemplate`,
  `TestRunReviewPromptIncludesPacketEvidence`, and
  `TestRunReviewMapsTimeoutToTimedOut`. Run them with a `go test -list`
  exact-count guard before the focused `go test -run` command so zero matches
  fail. These tests must cover runner allow-list behavior, prompt digest/ref
  preservation, empty/malformed/offtask outputs, OpenCode read-only missing,
  OpenCode mutation detection, raw output retention, not-assessed override,
  unreadable prompt template handling, prompt evidence injection, timeout
  mapping, runner error classification, profile role order preservation,
  dirty-baseline handling as `not_assessed` with `working_tree_dirty`, and a
  role with no command preserving `StatusNotAssessed` and
  `runner_command_not_configured` without invoking a runner.
- [x] T021-6650 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6660 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6670 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 96 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-96-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 97 Tasks

- [x] T021-6680 Confirm Slice 97 is bounded to numbered `internal/prreview`
  parsed reviewer output and raw result retention shards `prreview_125` through
  `prreview_133`.
- [x] T021-6681 Confirm Slice 97 is behavior-preserving: no changes to strict
  JSON decoding, packet/plane/role mismatch handling, off-task status/reason,
  parsed run ID/runner/model/status defaults, sanitizer invocation, execution
  metadata attachment, command digest and prompt ref propagation, raw-output
  digest shape, digest-only retention, retained raw output path/permissions,
  package boundary, dependency direction, or baseline changes are planned.
- [x] T021-6682 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-97-plan-review.md`.
- [x] T021-6690 Move parsed reviewer output normalization and raw result
  retention helpers into cohesive responsibility-named files without creating
  standalone one-function replacement files. Record source-shape evidence that
  numbered `prreview_125` through `prreview_133` files are gone and staged
  boundary evidence that excluded preview/copying `prreview_134` onward, prompt
  rendering/sanitization/citation, validation/summary, role execution, and
  previously consolidated validation files are not moved or edited in Slice 97.
- [x] T021-6700 Run `gofmt` on changed Go files.
- [x] T021-6710 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6711 Run focused parsed-output/raw-retention regression evidence
  covering exact named tests
  `TestRunReviewArtifactPipelineRedactsUnsafeReviewerText`,
  `TestRunReviewRecordsRunnerFailureStatesAndPromptDigest`,
  `TestRunReviewPreservesValidationDefaultsAndOutputContracts`, and
  `TestPacketProfileAndSmallHelpers`. Run them with a `go test -list`
  exact-count guard before the focused `go test -run` command so zero matches
  fail. These tests must cover strict reviewer JSON decoding, including
  rejection of reviewer output with an unknown JSON field, wrong
  packet/plane/role off-task handling, exact parsed defaults for missing
  `review_run_id`, `runner`, requested/observed model fields, model
  family/version, and status defaulting from findings/no-findings, sanitizer
  invocation, execution metadata attachment, command digest and prompt ref
  propagation, raw-output digest-only retention without persisted bytes,
  retained raw output path/digest/file mode `0o600` where mode bits are exposed,
  and default reviewer status with/without findings.
- [x] T021-6720 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6730 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6740 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 97 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-97-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 98 Tasks

- [x] T021-6750 Confirm Slice 98 is bounded to numbered `internal/prreview`
  preview, prompt digest, copied input, packet digest, and output directory
  helper shards `prreview_134` through `prreview_148`.
- [x] T021-6751 Confirm Slice 98 is behavior-preserving: no changes to
  preview-only behavior, preview role order, command digest, prompt digest/ref
  shape, empty prompt refs, prompt read error handling, copied input naming,
  content type, digest/ref shape, copied input file permissions, canonical
  packet digest replay, missing output path and existing output directory error
  strings, package boundary, dependency direction, or baseline changes are
  planned.
- [x] T021-6752 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-98-plan-review.md`.
- [x] T021-6760 Move preview, prompt digest, copied input, packet digest, and
  output directory helpers into cohesive responsibility-named files without
  creating standalone one-function replacement files. Record source-shape
  evidence that numbered `prreview_134` through `prreview_148` files are gone
  and staged boundary evidence that excluded CI/runner/status/disposition
  helpers `prreview_149` through `prreview_156`, citation resolution
  `prreview_157` onward, unsafe text/safe ID helpers, prompt rendering,
  sanitizer, validation/summary, role execution, parsed output handling, and
  previously consolidated validation files are not moved or edited in Slice 98.
- [x] T021-6770 Run `gofmt` on changed Go files.
- [x] T021-6771 Add focused regression coverage named
  `TestBuildPacketCopiesInputsAndComputesStableDigests` before running the
  focused exact-count guard. The test must stat copied context and verification
  input files for mode `0o644` where mode bits are exposed, read copied bytes,
  assert exact copied file names, refs, IDs, kinds, content types, digest
  values, and replay packet digest with `packet_digest` cleared.
- [x] T021-6780 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6781 Run focused preview/input/digest regression evidence covering
  exact named tests `TestRunReviewPreviewReturnsPreviewOnly`,
  `TestRunReviewPreservesValidationDefaultsAndOutputContracts`,
  `TestRunReviewCannotVerifyUnreadablePromptTemplate`,
  `TestRunReviewPromptIncludesPacketEvidence`,
  `TestBuildPacketCopiesInputsAndComputesStableDigests`, and
  `TestPacketProfileAndSmallHelpers`. Run them with a `go test -list`
  exact-count guard before the focused `go test -run` command so zero matches
  fail. These tests must cover preview-only no artifact writes, preview role
  order and command/prompt digests, empty prompt refs, prompt read error
  handling, copied context/verification input naming/content-type/digest/ref
  shape and file mode `0o644` where mode bits are exposed, stable packet digest
  replay with `packet_digest` cleared, missing output path, existing output
  directory rejection, and read-dir failure propagation.
- [x] T021-6790 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6800 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6810 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 98 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-98-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 99 Tasks

- [x] T021-6820 Confirm Slice 99 is bounded to numbered `internal/prreview`
  CI/runner validation, plane result, reviewer status/action, disposition, and
  severity helper shards `prreview_149` through `prreview_156`.
- [x] T021-6821 Confirm Slice 99 is behavior-preserving: no changes to accepted
  CI states, accepted runner IDs, usable reviewer status semantics, retained raw
  output requirements, plane-result degradation to `cannot_verify`, exact
  reviewer reason/next-action strings, default disposition mapping,
  unknown-severity fallback, package boundary, dependency direction, or baseline
  changes are planned.
- [x] T021-6822 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-99-plan-review.md`.
- [x] T021-6830 Move CI/runner validation, plane result, reviewer
  status/action, disposition, and severity helpers into cohesive
  responsibility-named files without creating standalone one-function
  replacement files. Record source-shape evidence that numbered `prreview_149`
  through `prreview_156` files are gone and staged boundary evidence that
  excluded citation resolution `prreview_157` onward, unsafe text/safe ID and
  content-type helpers `prreview_169` onward, prompt rendering, sanitizer,
  validation/summary orchestration, role execution, parsed output handling, and
  Slice 98 packet/preview/input helpers are not moved or edited in Slice 99.
- [x] T021-6840 Run `gofmt` on changed Go files.
- [x] T021-6841 Add focused helper regression coverage named
  `TestPrreviewStatusDispositionHelpersPreserveContracts` before running the
  focused exact-count guard. The test must assert accepted/rejected CI states,
  accepted/rejected runners, usable statuses, retained raw output requirements,
  plane-result degradation to `cannot_verify`, exact reviewer reason/next-action
  strings, default disposition mapping, and unknown severity fallback.
- [x] T021-6850 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6851 Run focused validation/status/disposition regression evidence
  with an exact-count guard for named tests covering:
  `TestValidateReviewStatesAndAuthorityBoundary`,
  `TestValidateCannotVerifyUsableStatusWithoutRetainedOutput`,
  `TestPrreviewValidationOrchestrationPreservesDigestRequiredPlaneAndAuthorityContracts`,
  `TestPrreviewValidationRankingModelAndLedgerFindingContracts`,
  `TestLedgerDispositionCarryForward`,
  `TestPrreviewLedgerSynthesisPreservesOrderingCarryForwardAndSanitization`,
  and `TestPrreviewStatusDispositionHelpersPreserveContracts`. These tests must
  cover accepted/rejected CI states, accepted/rejected runners, usable statuses
  with and without retained raw output, exact reviewer reason/next-action
  strings, default disposition mapping, and unknown severity fallback.
- [x] T021-6860 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6870 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6880 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 99 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-99-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 100 Tasks

- [x] T021-6890 Confirm Slice 100 is bounded to numbered `internal/prreview`
  citation resolver shards `prreview_157` through `prreview_168`.
- [x] T021-6891 Confirm Slice 100 is behavior-preserving: no changes to
  citation anchor requirements, diff alias and diff-ref ID matching,
  context/verification ref ID lookup, unknown-ref source-digest fallback,
  digest-only citation acceptance, diff/context/verification location rules,
  package boundary, dependency direction, or baseline changes are planned.
- [x] T021-6892 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-100-plan-review.md`.
- [x] T021-6900 Move citation anchor, resolver dispatch,
  diff/context/verification matching, safe ref ID lookup, and citation location
  helpers into cohesive responsibility-named files without creating standalone
  one-function replacement files. Record source-shape evidence that numbered
  `prreview_157` through `prreview_168` files are gone and staged boundary
  evidence that excluded unsafe text helpers `prreview_169` through
  `prreview_172`, command/context/content/safe ID helpers `prreview_173`
  onward, prompt rendering, sanitizer, validation/summary orchestration, role
  execution, parsed output handling, and previous packet/preview/status helpers
  are not moved or edited in Slice 100.
- [x] T021-6910 Run `gofmt` on changed Go files.
- [x] T021-6920 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6921 Run focused citation regression evidence with an exact-count
  guard for named tests `TestCitationResolvableCharacterization`,
  `TestPrreviewValidationRankingModelAndLedgerFindingContracts`, and
  `TestPrreviewLedgerSynthesisPreservesOrderingCarryForwardAndSanitization`.
  These tests must cover empty citation rejection, diff alias and diff-ref ID
  matching, context and verification ref lookup, unknown ref source-digest
  fallback, unknown ref without source digest rejection, digest-only citation
  acceptance, context citation by diff hunk, context citation by source digest,
  verification citation by source digest, and exact diff/context/verification
  location rules.
- [x] T021-6930 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-6940 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-6950 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 100 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-100-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 101 Tasks

- [x] T021-6960 Confirm Slice 101 is bounded to numbered `internal/prreview`
  unsafe text, unique string, digest, context/content type, safe ID, and default
  string shards `prreview_169` through `prreview_182`.
- [x] T021-6961 Confirm Slice 101 is behavior-preserving: no changes to unsafe
  marker/pattern redaction, empty text behavior, unique string safe
  dedupe/sorting, command digest empty/non-empty and NUL-separated hashing,
  context kind mapping, content type mapping, normalized extension fallback,
  safe ID mapping/fallback, default string fallback, package boundary,
  dependency direction, or baseline changes are planned.
- [x] T021-6962 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-101-plan-review.md`.
- [x] T021-6970 Move unsafe text, unique string, command digest, context/content
  type, normalized extension, safe ID, and default string helpers into cohesive
  responsibility-named files without creating standalone one-function
  replacement files. Record source-shape evidence that numbered `prreview_169`
  through `prreview_182` files are gone and staged boundary evidence that
  excluded generic `copy` and prompt rendering/sanitizer helpers
  `prreview_183` onward, citation resolution, validation/summary
  orchestration, role execution, parsed output handling, and previous
  packet/preview/status helpers are not moved or edited in Slice 101.
- [x] T021-6980 Run `gofmt` on changed Go files.
- [x] T021-6981 Add focused helper regression coverage named
  `TestPrreviewSmallHelpersPreserveContracts` before running the focused
  exact-count guard. The test must assert empty command digest, NUL-separated
  non-empty command digest hashing, whitespace-only `defaultString` fallback,
  normalized extension allow-list and `.txt` fallback, content types for
  JSON/Markdown/text, and context kind mapping for task Markdown, ordinary
  Markdown, JSON, and source excerpt inputs.
- [x] T021-6990 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-6991 Run focused unsafe/helper regression evidence with an
  exact-count guard for named tests `TestValidationAndSummaryRedactUnsafeMarkerClasses`,
  `TestPrreviewLedgerSynthesisPreservesOrderingCarryForwardAndSanitization`,
  `TestRunReviewPreviewReturnsPreviewOnly`,
  `TestBuildPacketCopiesInputsAndComputesStableDigests`,
  `TestPacketProfileAndSmallHelpers`, `TestSafeID`, and
  `TestPrreviewSmallHelpersPreserveContracts`. These tests must cover unsafe
  marker/pattern redaction, unique safe dedupe/sorting, command digest, context
  kind, content type, normalized extension, safe ID mapping/fallback, and
  default string fallback.
- [x] T021-7000 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7010 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7020 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 101 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-101-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 102 Tasks

- [x] T021-7030 Confirm Slice 102 is bounded to remaining numbered
  `internal/prreview` helpers `prreview_183` through `prreview_192`.
- [x] T021-7031 Confirm Slice 102 is behavior-preserving: no changes to public
  `Copy`, prompt template rendering, prompt evidence rendering, packet ref
  path/digest checks, prompt token replacement, reviewer result sanitization,
  redaction string spelling, package boundary, dependency direction, or MI
  baselines are planned.
- [x] T021-7032 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-102-plan-review.md`.
- [x] T021-7040 Move generic copy, prompt rendering/evidence/ref-read,
  prompt replacement, reviewer result sanitization, and redaction constant
  helpers into cohesive responsibility-named files without creating standalone
  one-function replacement files. Record source-shape evidence that numbered
  `prreview_183` through `prreview_192` files are gone and that previously
  consolidated `prreview_001` through `prreview_182` source files are not
  moved or edited in Slice 102. Record staged boundary evidence that schema
  files, CLI files, dependency manifests, role execution, parsed output,
  validation/summary, citation, and packet construction files are not touched.
- [x] T021-7050 Run `gofmt` on changed Go files.
- [x] T021-7051 Add focused prompt/sanitizer regression coverage named
  `TestPrreviewPromptSanitizerAndCopyHelpersPreserveContracts`. The test must
  assert public `Copy` behavior, template absence and read-failure behavior,
  prompt packet JSON rendering, prompt evidence ordering and fences, unsafe
  packet ref path rejection, digest mismatch `cannot_verify`, narrow prompt
  token replacement, reviewer result text/evidence-ref sanitization, and
  `redactedUnsafeReviewerText` spelling.
- [x] T021-7060 Run focused Go verification with `go test ./internal/prreview`.
- [x] T021-7061 Run focused prompt/sanitizer regression evidence with an
  exact-count guard for named tests
  `TestRunReviewRecordsRunnerFailureStatesAndPromptDigest`,
  `TestRunReviewCannotVerifyUnreadablePromptTemplate`,
  `TestRunReviewPromptIncludesPacketEvidence`,
  `TestValidationAndSummaryRedactUnsafeMarkerClasses`, and
  `TestPrreviewPromptSanitizerAndCopyHelpersPreserveContracts`. These tests
  must cover run-review prompt digest/evidence behavior, unreadable prompt
  template failure, prompt helper contracts, reviewer unsafe text sanitization,
  redaction spelling, and public `Copy` behavior.
- [x] T021-7070 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7080 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7090 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 102 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-102-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 103 Tasks

- [x] T021-7100 Confirm Slice 103 is bounded to the remaining repo-wide
  numbered Go filename `cmd/sdp-trace/spec019_proof_pack_cli_test.go`.
- [x] T021-7101 Confirm Slice 103 is behavior-preserving: no changes to proof
  pack CLI replay behavior, fixture paths/content, digest assertion semantics,
  test names, package boundary, dependency direction, or MI baselines are
  planned.
- [x] T021-7102 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-103-plan-review.md`.
- [x] T021-7110 Rename the test file and same-file helper names from
  spec-number language to domain responsibility language without touching
  product code, schema files, examples, fixtures, CLI behavior, or spec 019
  docs. Record repo-wide source-shape evidence that no numbered Go filenames
  remain. Record staged boundary evidence that product code, schema files,
  examples, fixtures, CLI behavior files, spec 019 docs, public contract files,
  and dependency manifests were not moved or edited beyond the renamed test
  file.
- [x] T021-7120 Run `gofmt` on changed Go files.
- [x] T021-7130 Run focused Go verification with a `go test -list`
  exact-count guard for `TestSpec019MonitoringGateProofPack`,
  `TestSpec019HarnessProofPack`, `TestSpec019TelemetryProofPack`, and
  `TestSpec019HarnessFixtureDigest`, then run
  `go test ./cmd/sdp-trace -run 'TestSpec019MonitoringGateProofPack|TestSpec019HarnessProofPack|TestSpec019TelemetryProofPack|TestSpec019HarnessFixtureDigest' -count=1 -v`.
- [x] T021-7140 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7150 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7160 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 103 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-103-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 104 Tasks

- [x] T021-7170 Confirm Slice 104 is bounded to `cmd/sdp-trace` `flagset_*`
  parser helper files and the existing `flagset_parse.go` responsibility
  boundary.
- [x] T021-7171 Confirm Slice 104 is behavior-preserving: no changes to `--`
  rest handling, positional argument preservation, unknown flag errors,
  `--flag=value`, bare boolean defaults, boolean literal consumption, missing
  string value errors, command-specific flag registration, package boundary,
  dependency direction, or MI baselines are planned.
- [x] T021-7172 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-104-plan-review.md`.
- [x] T021-7180 Consolidate `flagset_bool_value_at.go`,
  `flagset_consume.go`, `flagset_consume_arg.go`,
  `flagset_consume_bool.go`, `flagset_consume_flag.go`,
  `flagset_consume_no_equals.go`, `flagset_consume_string.go`, and
  `flagset_split.go` into cohesive parser responsibility files without
  changing command-specific flag registration, CLI subcommand behavior,
  `flagset.go` storage/accessor methods, public docs outside SpecKit
  artifacts, schema files, examples, fixtures, or dependencies. A small split
  across `flagset_parse.go`, dispatch, value, and boolean parser files is
  allowed only to preserve MI > 70 without recreating one-function shards.
  Record source-shape evidence that those one-function parser helper files are
  gone and staged boundary evidence that excluded areas are untouched.
- [x] T021-7190 Run `gofmt` on changed Go files.
- [x] T021-7191 Extend focused flagset regression coverage so
  `TestFlagSetParsesQuotedStringValue` asserts successful `--string=value`
  parsing and `TestFlagSetParsesBooleanValues` asserts boolean literal
  consumption from the next argv element.
- [x] T021-7200 Run focused flag parser verification with a `go test -list`
  exact-count guard for existing tests `TestFlagSetParsesQuotedStringValue`,
  `TestFlagSetRejectsMissingStringValueAtEnd`,
  `TestFlagSetRejectsUnknownFlagWithEquals`,
  `TestFlagSetParsesBooleanValues`,
  `TestFlagSetRepeatsFlagsOverwriteAndKeepOrder`, and
  `TestFlagSetCapturesEverythingAfterSeparator`, then run those tests with
  `go test ./cmd/sdp-trace -run <guard> -count=1 -v`. These tests must cover
  string value parsing including `--string=value`, missing string value
  diagnostics, unknown flag errors, boolean parsing forms including next-argv
  literal consumption, repeated flag overwrite/order behavior, and `--`
  separator rest capture.
- [x] T021-7210 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7220 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7230 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 104 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-104-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 105 Tasks

- [x] T021-7240 Confirm Slice 105 is bounded to `cmd/sdp-trace` standard gate
  argument parser helper files: `gate_parse_args.go`, `gate_target_arg.go`,
  `gate_output_path.go`, and `gate_string_flags.go`.
- [x] T021-7241 Confirm Slice 105 is behavior-preserving: no changes to
  `flagSet` parsing, standard gate execution, protected gate execution,
  preview, explain, gate exit-code logic, public command surface docs, schema
  files, examples, fixtures, dependencies, package boundary, dependency
  direction, or MI baselines are planned.
- [x] T021-7242 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-105-plan-review.md`.
- [x] T021-7250 Replace the four scoped helper shards
  `gate_parse_args.go`, `gate_target_arg.go`, `gate_output_path.go`, and
  `gate_string_flags.go` with one cohesive gate argument parser responsibility
  file, using `gate_parse_args.go` as the destination name, without creating
  replacement one-function shards. Record source-shape evidence that all four
  scoped helper shards are gone or replaced by that one cohesive responsibility
  file and staged boundary evidence that excluded gate files and public
  surfaces are untouched.
- [x] T021-7260 Run `gofmt` on changed Go files.
- [x] T021-7270 Add or extend focused regression coverage for `parseGateArgs`
  so it directly asserts successful parsing of the standard/protected gate
  string flags, exact missing target diagnostic, exact multiple target
  diagnostic, exact missing `--out` diagnostic, and unknown flag diagnostic.
- [x] T021-7280 Run focused gate parser verification with a `go test -list`
  exact-count guard expecting exactly 7 tests:
  `TestParseGateArgsPreservesContract`,
  `TestReportAndGateCommands`, `TestGateCommandAcceptsWitness`,
  `TestProtectedGateRequiresCheckpointPolicyAndWitnessFlags`,
  `TestProtectedGateMalformedNamedInputIsUsageError`,
  `TestDefaultGateDoesNotEmitProtectedFields`, and
  `TestProtectedGatePassesWithCISignedCheckpointAndBoundWitnessCLI`, then run
  those tests with `go test ./cmd/sdp-trace -run <guard> -count=1 -v`.
- [x] T021-7290 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7300 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7310 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 105 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-105-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 106 Tasks

- [x] T021-7320 Confirm Slice 106 is bounded to `cmd/sdp-trace` gate exit-code
  helper files: `gate_exit_code.go`, `gate_exit_states.go`,
  `gate_has_state.go`, `gate_protected_exit_code.go`,
  `gate_protected_exit_codes.go`, and `gate_state_exit_code.go`.
- [x] T021-7321 Confirm Slice 106 is behavior-preserving: no changes to gate
  evaluation, protected gate input loading, parser behavior, preview, explain,
  report rendering, public command surface docs, schema files, examples,
  fixtures, dependencies, package boundary, dependency direction, or MI
  baselines are planned.
- [x] T021-7322 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-106-plan-review.md`.
- [x] T021-7330 Replace the six scoped helper shards with cohesive gate
  exit-code responsibility files, using `gate_exit_code.go` as the primary
  destination and allowing a protected exit mapping split only if needed to
  preserve MI > 70, without creating replacement one-function shards. Record
  source-shape evidence that the six scoped helper shards are gone or replaced
  by the cohesive responsibility files and staged boundary evidence that
  excluded gate files and public surfaces are untouched.
- [x] T021-7340 Run `gofmt` on changed Go files.
- [x] T021-7350 Run focused gate exit-code verification with a `go test -list`
  exact-count guard expecting exactly 4 tests:
  `TestGateExitCodeChecksRequiredRunStatesDirectly`,
  `TestGateExitCodeAggregatesNonProtectedStates`,
  `TestGateExitCodeUsesProtectedGateWhenSelected`, and
  `TestGateAndFixtureHelpers`, then run those tests with
  `go test ./cmd/sdp-trace -run <guard> -count=1 -v`. The focused set must
  explicitly cover protected `pass`, `fail`, `cannot_verify`, `not_assessed`,
  and unknown-state fallback exit-code behavior.
- [x] T021-7360 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7370 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7380 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 106 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-106-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 107 Tasks

- [x] T021-7390 Confirm Slice 107 is bounded to `cmd/sdp-trace` gate
  subcommand dispatch helper files: `gate_subcommand_handlers.go`,
  `gate_subcommand_run.go`, and the existing routing destination
  `gate_run.go`.
- [x] T021-7391 Confirm Slice 107 is behavior-preserving: no changes to
  standard gate execution, protected gate execution, preview behavior, explain
  behavior, argument parsing, gate exit-code logic, public command surface
  docs, schema files, examples, fixtures, dependencies, package boundary,
  dependency direction, or MI baselines are planned.
- [x] T021-7392 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-107-plan-review.md`.
- [x] T021-7400 Consolidate `gate_subcommand_handlers.go` and
  `gate_subcommand_run.go` into `gate_run.go` without creating replacement
  one-function shards. Record source-shape evidence that the two helper shards
  are gone and staged boundary evidence that excluded gate files and public
  surfaces are untouched.
- [x] T021-7410 Run `gofmt` on changed Go files.
- [x] T021-7420 Add or extend focused regression coverage for
  `runGateSubcommand` so it directly asserts preview dispatch, explain dispatch,
  non-subcommand fallback, and unknown-subcommand fallback to parent gate flag
  parsing (`handled=false`, code `0`) from the shared optional subcommand
  dispatcher.
- [x] T021-7430 Run focused gate subcommand dispatch verification with a
  `go test -list` exact-count guard expecting exactly 6 tests:
  `TestRunGateSubcommandPreservesDispatch`,
  `TestReportAndGateCommands`,
  `TestGateExplainParseUsage`,
  `TestGatePreviewIsReadOnlyAndDoesNotPrintSecretLikeValues`,
  `TestGatePreviewParseAndContractFailurePaths`, and
  `TestProtectedGatePreviewRendersAbsentInputsWithoutWriting`, then run those
  tests with `go test ./cmd/sdp-trace -run <guard> -count=1 -v`.
- [x] T021-7440 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7450 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7460 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 107 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-107-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 108 Tasks

- [x] T021-7470 Confirm Slice 108 is bounded to `cmd/sdp-trace`
  `gate_standard_run.go`, `gate_preview_standard.go`, and the cohesive
  destination `gate_standard.go`.
- [x] T021-7471 Confirm Slice 108 is behavior-preserving: no changes to
  standard gate evaluation, standard preview behavior, protected gate
  execution, protected preview behavior, explain, argument parsing, gate
  exit-code logic, public command surface docs, schema files, examples,
  fixtures, dependencies, package boundary, dependency direction, or MI
  baselines are planned.
- [x] T021-7472 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-108-plan-review.md`.
- [x] T021-7480 Consolidate `gate_standard_run.go` and
  `gate_preview_standard.go` into `gate_standard.go` without creating
  replacement one-function shards. Record source-shape evidence that both
  source helper shards are gone and staged boundary evidence that excluded gate
  files and public surfaces are untouched.
- [x] T021-7490 Run `gofmt` on changed Go files.
- [x] T021-7495 Add or extend focused regression coverage for
  `runStandardGate` so it directly asserts indented JSON output with trailing
  newline and stderr/nonzero-exit propagation for a `demo.WriteGate` error.
- [x] T021-7500 Run focused standard gate verification with a `go test -list`
  exact-count guard expecting exactly 8 tests:
  `TestRunStandardGatePreservesOutputAndErrors`,
  `TestReportAndGateCommands`,
  `TestGateCommandAcceptsWitness`,
  `TestGateCommandFailsForWitnessArtifactMismatch`,
  `TestGateCommandCannotVerifyWhenWitnessOmitsRunArtifact`,
  `TestDefaultGateDoesNotEmitProtectedFields`, and
  `TestGatePreviewStandardReportShape`, and
  `TestGatePreviewReportsWitnessArtifactMismatch`, then run those tests with
  `go test ./cmd/sdp-trace -run <guard> -count=1 -v`.
- [x] T021-7510 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7520 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7530 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 108 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-108-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 109 Tasks

- [x] T021-7540 Confirm Slice 109 is bounded to `internal/interaction`
  event type catalog files: `interaction_event_types.go`,
  `interaction_event_types_data.go`, `interaction_event_validation_type.go`,
  and the cohesive destination `interaction_event_catalog.go`.
- [x] T021-7541 Confirm Slice 109 is behavior-preserving: no changes to event
  type values, ordering, validation behavior, import/export behavior, schemas,
  examples, fixtures, dependencies, package boundary, dependency direction, or
  MI baselines are planned.
- [x] T021-7542 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-109-plan-review.md`.
- [x] T021-7550 Consolidate `interaction_event_types.go`,
  `interaction_event_types_data.go`, and
  `interaction_event_validation_type.go` into
  `interaction_event_catalog.go` without creating replacement micro-shards.
  Record source-shape evidence that the source helper shards are gone and
  staged boundary evidence that excluded `internal/interaction` behavior files
  and public surfaces are untouched.
- [x] T021-7560 Run `gofmt` on changed Go files.
- [x] T021-7570 Add focused regression coverage for `EventTypes()` proving it
  returns event types in the existing order and returns a defensive copy that
  callers cannot use to mutate the backing catalog, and extend moved validation
  coverage so unsupported event type and friction-class mismatch diagnostics are
  preserved.
- [x] T021-7580 Run focused interaction event catalog verification with a
  `go test -list` exact-count guard expecting exactly 2 tests:
  `TestEventTypesReturnsStableCopy` and
  `TestValidateEventRejectsInvalidFields`, then run them with
  `go test ./internal/interaction -run <guard> -count=1 -v`.
- [x] T021-7590 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7600 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7610 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 109 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-109-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 110 Tasks

- [x] T021-7620 Confirm Slice 110 is bounded to `internal/posture` metric row
  validation files: `posture_validate_metric_counts.go`,
  `posture_validate_metric_identity.go`, and
  `posture_validate_metric_row.go`.
- [x] T021-7621 Confirm Slice 110 is behavior-preserving: no changes to
  malformed metric row detection, export validation behavior, schemas,
  examples, fixtures, dependencies, package boundary, dependency direction, or
  MI baselines are planned.
- [x] T021-7622 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-110-plan-review.md`.
- [x] T021-7630 Consolidate `posture_validate_metric_counts.go` and
  `posture_validate_metric_identity.go` into
  `posture_validate_metric_row.go` without creating replacement micro-shards.
  Record source-shape evidence that the metric helper shards are gone and
  staged boundary evidence that excluded posture behavior files and public
  surfaces are untouched.
- [x] T021-7640 Run `gofmt` on changed Go files.
- [x] T021-7650 Run focused posture validation verification with a
  `go test -list` exact-count guard expecting exactly 1 test:
  `TestValidateMetricRowShapeRejectsMalformedRows`, then run it with
  `go test ./internal/posture -run <guard> -count=1 -v`. The metric row test
  must cover the moved count predicates for negative numerator, negative
  denominator, invalid unit, and negative not-assessed count.
- [x] T021-7660 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7670 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7680 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 110 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-110-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 111 Tasks

- [x] T021-7690 Confirm Slice 111 is bounded to `internal/posture` movement row
  identity validation files: `posture_validate_movement_identity.go` and
  `posture_validate_movement_row.go`.
- [x] T021-7691 Confirm Slice 111 is behavior-preserving: no changes to
  malformed movement identity detection, malformed movement row detection,
  comparison-basis behavior, non-comparable reason handling, movement summary
  validation, export validation behavior, schemas, examples, fixtures,
  dependencies, package boundary, dependency direction, or MI baselines are
  planned.
- [x] T021-7692 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-111-plan-review.md`.
- [x] T021-7700 Consolidate `posture_validate_movement_identity.go` into
  `posture_validate_movement_row.go` without creating replacement
  micro-shards. Record source-shape evidence that the movement identity helper
  shard is gone and staged boundary evidence that excluded posture behavior
  files and public surfaces are untouched.
- [x] T021-7710 Run `gofmt` on changed Go files.
- [x] T021-7720 Run focused posture movement validation verification with a
  `go test -list` exact-count guard expecting exactly 2 tests:
  `TestPostureHelperCoverageForMovementAndSafety` and
  `TestValidateMovementRowRejectsMalformedRows`, then run them with
  `go test ./internal/posture -run <guard> -count=1 -v`.
- [x] T021-7730 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7740 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7750 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 111 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-111-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 112 Tasks

- [x] T021-7760 Confirm Slice 112 is bounded to `internal/packet`
  demo-first closure assessment files: `demo_first_state_helpers.go` and
  `demo_first_closure_requirements.go`.
- [x] T021-7761 Confirm Slice 112 is behavior-preserving: no changes to
  demo-first route evidence, row evidence, pass-count logic, closure cap logic,
  packet schemas, examples, fixtures, dependencies, package boundary,
  dependency direction, CRAP < 5, MI > 70, or MI baselines are planned.
- [x] T021-7762 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-112-plan-review.md`.
- [x] T021-7770 Consolidate `demo_first_state_helpers.go` into
  `demo_first_closure_requirements.go` without creating replacement
  micro-shards. Record source-shape evidence that the state helper shard is
  gone and staged boundary evidence that excluded packet behavior files and
  public surfaces are untouched.
- [x] T021-7780 Run `gofmt` on changed Go files.
- [x] T021-7790 Run focused packet demo-first closure verification with a
  `go test -list` exact-count guard expecting exactly 1 test:
  `TestCheckDemoRequiresVerificationOrReviewAssessed`, then run it with
  `go test ./internal/packet -run <guard> -count=1 -v`. The test must preserve
  the exact first-packet gate diagnostic and cover the `rowAssessed` state
  matrix: pass, partial, and fail count as assessed for the verification/review
  closure requirement while cannot_verify, not_assessed, and not_in_scope do
  not.
- [x] T021-7800 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7810 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7820 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 112 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-112-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 113 Tasks

- [x] T021-7830 Confirm Slice 113 is bounded to `internal/packet` bundle
  manifest indexing files: `validation_manifest_index.go` and
  `validation_manifest_entries.go`.
- [x] T021-7831 Confirm Slice 113 is behavior-preserving: no changes to
  resolver entry indexing, manifest entry validation, row validation,
  contradiction validation, residual gap validation, packet schemas, examples,
  fixtures, dependencies, package boundary, dependency direction, CRAP < 5, MI
  > 70, or MI baselines are planned.
- [x] T021-7832 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-113-plan-review.md`.
- [x] T021-7840 Consolidate `validation_manifest_index.go` into
  `validation_manifest_entries.go` without creating replacement micro-shards.
  Record source-shape evidence that the manifest index shard is gone and staged
  boundary evidence that excluded packet behavior files and public surfaces are
  untouched.
- [x] T021-7850 Run `gofmt` on changed Go files.
- [x] T021-7860 Run focused packet validation verification with a `go test
  -list` exact-count guard expecting exactly 2 tests:
  `TestValidateAcceptsResolverEntryForManifestEvidence` and
  `TestValidateManifestResolverIndexingFeedsRowEvidenceRefs`, then run them with
  `go test ./internal/packet -run <guard> -count=1 -v`. The focused test set
  must preserve that `indexManifest` still indexes both manifest entries and
  resolver entries.
- [x] T021-7870 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7880 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7890 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 113 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-113-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 114 Tasks

- [x] T021-7900 Confirm Slice 114 is bounded to `internal/packet` validation
  result entrypoint files: `validation_types.go` and
  `validation_entrypoint.go`.
- [x] T021-7901 Confirm Slice 114 is behavior-preserving: no changes to
  validation orchestration, demo-first validation, manifest/row/resolver
  validation, packet schemas, examples, fixtures, dependencies, package
  boundary, dependency direction, CRAP < 5, MI > 70, or MI baselines are
  planned.
- [x] T021-7902 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-114-plan-review.md`.
- [x] T021-7910 Consolidate `validation_types.go` into
  `validation_entrypoint.go` without creating replacement micro-shards. Record
  source-shape evidence that the validation type shard is gone and staged
  boundary evidence that excluded packet behavior files and public surfaces are
  untouched.
- [x] T021-7920 Run `gofmt` on changed Go files.
- [x] T021-7930 Run focused packet validation type verification with a
  `go test -list` exact-count guard expecting exactly 2 tests:
  `TestGitHubEvidenceInputTypesPreserveJSONShape` and
  `TestValidateRejectsBundleMismatch`, then run them with
  `go test ./internal/packet -run <guard> -count=1 -v`. The focused test set
  must preserve the exported `Validation` JSON shape and `Validate` return
  behavior. `TestGitHubEvidenceInputTypesPreserveJSONShape` must assert that
  `json.Marshal(Validation{})` preserves the `state` key and omits `errors`.
- [x] T021-7940 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-7950 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-7960 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 114 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-114-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 115 Tasks

- [x] T021-7970 Confirm Slice 115 is bounded to `internal/query` query-pack
  output safety files: `querypack_safety_type.go` and `querypack_safety.go`.
- [x] T021-7971 Confirm Slice 115 is behavior-preserving: no changes to
  query-pack result shape, builder behavior, safety provider/sensitive class
  catalogs, schemas, examples, fixtures, dependencies, package boundary,
  dependency direction, CRAP < 5, MI > 70, or MI baselines are planned.
- [x] T021-7972 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-115-plan-review.md`.
- [x] T021-7980 Consolidate `querypack_safety_type.go` into
  `querypack_safety.go` without creating replacement micro-shards. Record
  source-shape evidence that the safety type shard is gone and staged boundary
  evidence that excluded query behavior files and public surfaces are
  untouched.
- [x] T021-7990 Run `gofmt` on changed Go files.
- [x] T021-8000 Run focused query-pack output safety verification with a
  `go test -list` exact-count guard expecting exactly 2 tests:
  `TestQueryPackSchemaMirrorsImplementationEnums` and
  `TestForensicsBasicPackSafetyClassesAreVerifiedAgainstOutput`, then run them with
  `go test ./internal/query -run <guard> -count=1 -v`. The focused test set
  must preserve generated output safety values and must explicitly assert the
  `QueryPackOutputSafety` JSON tag contract: `json.Marshal(QueryPackOutputSafety{})`
  includes `verified_absent_sensitive_classes` and omits
  `redaction_policy_digest`, while a non-empty digest is emitted as
  `redaction_policy_digest`.
- [x] T021-8010 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-8020 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-8030 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 115 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-115-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 116 Tasks

- [x] T021-8040 Confirm Slice 116 is bounded to `internal/query` query-pack
  input artifact files: `querypack_artifact_type.go` and
  `querypack_artifact_schema.go`.
- [x] T021-8041 Confirm Slice 116 is behavior-preserving: no changes to
  artifact digesting, artifact reader control flow, input loading, builder
  result assembly, query-pack result shape, schemas, examples, dependencies,
  package boundary, dependency direction, CRAP < 5, MI > 70, or MI baselines
  are planned.
- [x] T021-8042 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-116-plan-review.md`.
- [x] T021-8050 Consolidate `querypack_artifact_type.go` into
  `querypack_artifact_schema.go` without creating replacement micro-shards.
  Record source-shape evidence that the artifact type shard is gone and staged
  boundary evidence that excluded query behavior files and public surfaces are
  untouched.
- [x] T021-8060 Run `gofmt` on changed Go files.
- [x] T021-8070 Run focused query-pack input artifact verification with a
  `go test -list` exact-count guard expecting exactly 2 tests:
  `TestForensicsBasicPackDerivesRowsWithoutPolicyVerdict` and
  `TestForensicsBasicPackPreservesUnreadableOptionalArtifactAsCannotVerifyRows`,
  then run them with `go test ./internal/query -run <guard> -count=1 -v`. The
  focused test set must preserve generated input artifact values and explicitly
  assert the `QueryPackInputArtifact` JSON tag contract:
  `json.Marshal(QueryPackInputArtifact{})` includes `role`, `path_redacted_id`,
  and `artifact_required`, omits empty `sha256` and `schema_version`, and emits
  non-empty `sha256` and `schema_version`.
- [x] T021-8080 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-8090 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-8100 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 116 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-116-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 117 Tasks

- [x] T021-8110 Confirm Slice 117 is bounded to `internal/query` query-pack row
  files: `querypack_row_type.go` and `querypack_row_factory.go`.
- [x] T021-8111 Confirm Slice 117 is behavior-preserving: no changes to row
  ordering, row source mapping, condition conversion, summary rows, explanation
  rendering, query-pack result shape, schemas, examples, dependencies, package
  boundary, dependency direction, CRAP < 5, MI > 70, or MI baselines are
  planned.
- [x] T021-8112 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-117-plan-review.md`.
- [x] T021-8120 Consolidate `querypack_row_type.go` into
  `querypack_row_factory.go` without creating replacement micro-shards. Record
  source-shape evidence that the row type shard is gone and staged boundary
  evidence that excluded query behavior files and public surfaces are
  untouched.
- [x] T021-8130 Run `gofmt` on changed Go files.
- [x] T021-8140 Run focused query-pack row verification with a `go test -list`
  exact-count guard expecting exactly 2 tests:
  `TestForensicsBasicPackDerivesRowsWithoutPolicyVerdict` and
  `TestExplainForensicsPackRendersStableSafeRows`, then run them with
  `go test ./internal/query -run <guard> -count=1 -v`. The focused test set
  must preserve generated query row values and explicitly assert the
  `QueryPackRow` JSON tag contract: required row fields emit under their
  current JSON names, optional empty fields are omitted, non-empty optional
  fields emit, `reconstructable` preserves pointer boolean semantics, and
  `related_rows` omits empty slices while emitting non-empty slices.
- [x] T021-8150 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-8160 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-8170 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 117 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-117-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 118 Tasks

- [x] T021-8180 Confirm Slice 118 is bounded to `internal/adaptercapture` valid
  event fixture spec files: `adaptercapture_valid_event_spec_type.go` and
  `adaptercapture_valid_event_specs.go`.
- [x] T021-8181 Confirm Slice 118 is behavior-preserving: no changes to valid
  event generation, event ordering, required event type lists, adapter-capture
  assessment behavior, schemas, examples, dependencies, package boundary,
  dependency direction, CRAP < 5, MI > 70, or MI baselines are planned.
- [x] T021-8182 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-118-plan-review.md`.
- [x] T021-8190 Consolidate `adaptercapture_valid_event_spec_type.go` into
  `adaptercapture_valid_event_specs.go` without creating replacement
  micro-shards. Record source-shape evidence that the spec type shard is gone
  and staged boundary evidence that excluded adapter-capture behavior files and
  public surfaces are untouched.
- [x] T021-8200 Run `gofmt` on changed Go files.
- [x] T021-8210 Run focused adapter-capture fixture verification with a
  `go test -list` exact-count guard expecting exactly 2 tests:
  `TestEvaluatePassesForBoundGenericAdapterEvents` and
  `TestEvaluateCannotVerifyMissingRequiredAdapterEvent`, then run them with
  `go test ./internal/adaptercapture -run <guard> -count=1 -v`. The focused
  test set must preserve generated valid adapter events from `validEventSpecs`
  and preserve required-event behavior without changing the separate required
  event type list.
- [x] T021-8220 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-8230 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-8240 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 118 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-118-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.

## Active Slice 119 Tasks

- [x] T021-8250 Confirm Slice 119 is bounded to `internal/authority` event type
  validation files: `authority_event_type.go`,
  `authority_event_set_validation.go`, and `authority_vars.go`.
- [x] T021-8251 Confirm Slice 119 is behavior-preserving: no changes to
  event-set validation, pre-decision blocker ordering, target-rule validation,
  custom event prefix semantics, authority evaluation behavior, schemas,
  examples, dependencies, package boundary, dependency direction, CRAP < 5, MI >
  70, or MI baselines are planned.
- [x] T021-8252 Run three independent plan/task reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record the plan review in
  `specs/021-source-file-locality-cleanup/reviews/slice-119-plan-review.md`.
- [x] T021-8260 Consolidate `authority_event_type.go` into
  `authority_event_set_validation.go` without creating replacement micro-shards.
  Keep `authority_vars.go` as the source of `standardEventTypes`. Record
  source-shape evidence that the event type shard is gone and staged boundary
  evidence that excluded authority behavior files and public surfaces are
  untouched.
- [x] T021-8270 Run `gofmt` on changed Go files.
- [x] T021-8280 Run focused authority event type verification with a
  `go test -list` exact-count guard expecting exactly 2 tests:
  `TestValidateEnvelopeReasons` and `TestEvaluateActionReasonOrdering`, then
  run them with `go test ./internal/authority -run <guard> -count=1 -v`. The
  focused test set must preserve unsupported event diagnostics, standard event
  acceptance, and `custom:` event acceptance through existing envelope/action
  evaluation paths. If current coverage lacks a `custom:` acceptance case,
  extend one of the two guarded tests before implementation review; do not add
  a third focused test name unless this task is re-reviewed.
- [x] T021-8290 Run repository verification: `go test ./...`, `go vet ./...`,
  `golangci-lint run` when available or record `not_assessed`/`cannot_verify`
  with reason, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq
  empty schema/*.json`, and `git diff --check`.
- [x] T021-8300 Run CRAP and MI quality gates without changing MI baselines.
- [x] T021-8310 Run three independent implementation reviewer lanes, fix every
  finding, repeat affected lanes until each reviewer returns exactly `LGTM`,
  and record Slice 119 evidence in
  `specs/021-source-file-locality-cleanup/reviews/slice-119-evidence.md`,
  including harness, agent id, model/provider if available, date, prompt class,
  timeout, retries, fallback, result, and any unavailable external/provider
  lanes as `not_assessed`.
