# Tasks: Numbered Code File Locality Cleanup

Status: complete

## Active Slice 1 Tasks

- [x] T023-001 Inventory remaining numbered active Go files by package.
- [x] T023-002 Confirm Slice 1 is bounded to `cmd/sdp-trace` release-proof
  command shards.
- [x] T023-003 Reject broader release-proof grouping when local pre-change MI
  analysis measures candidate files below the absolute file-MI threshold.
- [x] T023-010 Move selected release-proof declarations into
  `cmd/sdp-trace/release_proof_run.go`,
  `cmd/sdp-trace/release_proof_args.go`, and
  `cmd/sdp-trace/release_proof_policy.go`.
- [x] T023-020 Run `gofmt` on changed Go files.
- [x] T023-030 Run focused Go verification for `cmd/sdp-trace`.
- [x] T023-040 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T023-050 Run CRAP and MI quality gates without changing MI baselines.
- [x] T023-060 Run three independent reviewer lanes and record Slice 1 evidence
  in `specs/023-numbered-code-file-locality-cleanup/reviews/slice-1-evidence.md`.

## Active Slice 2 Tasks

- [x] T023-070 Confirm Slice 2 is bounded to `cmd/sdp-trace` observe command
  adapter and exit-policy shards.
- [x] T023-071 Reject single-file observe grouping when local pre-change MI
  analysis measures the candidate file below the absolute file-MI threshold.
- [x] T023-080 Move selected observe declarations into
  `cmd/sdp-trace/observe_command_adapters.go` and
  `cmd/sdp-trace/observe_exit_policy.go`.
- [x] T023-090 Run `gofmt` on changed Go files.
- [x] T023-100 Run focused Go verification for `cmd/sdp-trace`.
- [x] T023-110 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T023-120 Run CRAP and MI quality gates without changing MI baselines.
- [x] T023-130 Run three independent reviewer lanes and record Slice 2 evidence
  in `specs/023-numbered-code-file-locality-cleanup/reviews/slice-2-evidence.md`.

## Active Slice 3 Tasks

- [x] T023-140 Confirm Slice 3 is bounded to `cmd/sdp-trace` envelope
  summarize command shards.
- [x] T023-141 Reject single-file envelope grouping when local MI analysis
  measures the candidate file below the absolute file-MI threshold.
- [x] T023-150 Move selected envelope declarations into
  `cmd/sdp-trace/envelope_summary_run.go` and
  `cmd/sdp-trace/envelope_summary_args.go`.
- [x] T023-160 Run `gofmt` on changed Go files.
- [x] T023-170 Run focused Go verification for `cmd/sdp-trace`.
- [x] T023-180 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T023-190 Run CRAP and MI quality gates without changing MI baselines.
- [x] T023-200 Run three independent reviewer lanes and record Slice 3 evidence
  in `specs/023-numbered-code-file-locality-cleanup/reviews/slice-3-evidence.md`.

## Active Slice 4 Tasks

- [x] T023-210 Confirm Slice 4 is bounded to `cmd/sdp-trace` export dispatcher
  shards.
- [x] T023-220 Move selected export declarations into
  `cmd/sdp-trace/export_command.go`.
- [x] T023-230 Run `gofmt` on changed Go files.
- [x] T023-240 Run focused Go verification for `cmd/sdp-trace`.
- [x] T023-250 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T023-260 Run CRAP and MI quality gates without changing MI baselines.
- [x] T023-270 Run three independent reviewer lanes and record Slice 4 evidence
  in `specs/023-numbered-code-file-locality-cleanup/reviews/slice-4-evidence.md`.

## Active Slice 5 Tasks

- [x] T023-280 Confirm Slice 5 is bounded to `cmd/sdp-trace` fixture validation
  command shards.
- [x] T023-281 Reject single-file fixture validation grouping when local MI
  analysis measures the candidate file below the absolute file-MI threshold.
- [x] T023-290 Move selected fixture declarations into
  `cmd/sdp-trace/fixture_validation_run.go`,
  `cmd/sdp-trace/fixture_validation_args.go`, and
  `cmd/sdp-trace/fixture_expectation_policy.go`.
- [x] T023-300 Run `gofmt` on changed Go files.
- [x] T023-310 Run focused Go verification for `cmd/sdp-trace`.
- [x] T023-320 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T023-330 Run CRAP and MI quality gates without changing MI baselines.
- [x] T023-340 Run three independent reviewer lanes and record Slice 5 evidence
  in `specs/023-numbered-code-file-locality-cleanup/reviews/slice-5-evidence.md`.

## Active Slice 6 Tasks

- [x] T023-350 Confirm Slice 6 is bounded to `cmd/sdp-trace` interaction
  command shards and the directly shared CLI flag requirement helper.
- [x] T023-351 Reject transcript import single-file grouping when local MI
  analysis measures the candidate file below the absolute file-MI threshold.
- [x] T023-360 Move selected interaction declarations into behavior-named
  command files.
- [x] T023-370 Run `gofmt` on changed Go files.
- [x] T023-380 Run focused Go verification for `cmd/sdp-trace`.
- [x] T023-390 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T023-400 Run CRAP and MI quality gates without changing MI baselines.
- [x] T023-410 Run three independent reviewer lanes and record Slice 6 evidence
  in `specs/023-numbered-code-file-locality-cleanup/reviews/slice-6-evidence.md`.

## Active Slice 7 Tasks

- [x] T023-420 Confirm Slice 7 is bounded to `cmd/sdp-trace` wrap, run,
  preview, and dry-run command shards.
- [x] T023-421 Reject combined preview grouping when local MI analysis measures
  the candidate file below the absolute file-MI threshold.
- [x] T023-430 Move selected wrap declarations into behavior-named command
  files.
- [x] T023-440 Run `gofmt` on changed Go files.
- [x] T023-450 Run focused Go verification for `cmd/sdp-trace`.
- [x] T023-460 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T023-470 Run CRAP and MI quality gates without changing MI baselines.
- [x] T023-480 Run three independent reviewer lanes and record Slice 7 evidence
  in `specs/023-numbered-code-file-locality-cleanup/reviews/slice-7-evidence.md`.

## Active Slice 8 Tasks

- [x] T023-490 Confirm Slice 8 is bounded to `cmd/sdp-trace` query, verify,
  explain, and query-pack command shards.
- [x] T023-491 Reject combined verify, query, and query-pack groupings when
  local MI analysis measures candidate files below the absolute file-MI
  threshold.
- [x] T023-500 Move selected query declarations into behavior-named command
  files.
- [x] T023-510 Run `gofmt` on changed Go files.
- [x] T023-520 Run focused Go verification for `cmd/sdp-trace`.
- [x] T023-530 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T023-540 Run CRAP and MI quality gates without changing MI baselines.
- [x] T023-550 Run three independent reviewer lanes and record Slice 8 evidence
  in `specs/023-numbered-code-file-locality-cleanup/reviews/slice-8-evidence.md`.

## Active Slice 9 Tasks

- [x] T023-560 Confirm Slice 9 is bounded to `cmd/sdp-trace` witness command
  shards.
- [x] T023-561 Reject broader witness option and validation groupings when
  local MI analysis measures candidate files below the absolute file-MI
  threshold.
- [x] T023-570 Move selected witness declarations into behavior-named command
  files.
- [x] T023-580 Run `gofmt` on changed Go files.
- [x] T023-590 Run focused Go verification for `cmd/sdp-trace`.
- [x] T023-600 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T023-610 Run CRAP and MI quality gates without changing MI baselines.
- [x] T023-620 Run three independent reviewer lanes and record Slice 9 evidence
  in `specs/023-numbered-code-file-locality-cleanup/reviews/slice-9-evidence.md`.

## Active Slice 10 Tasks

- [x] T023-630 Confirm Slice 10 is bounded to `cmd/sdp-trace` doctor
  repo-observer and install command shards.
- [x] T023-631 Reject broader install-args grouping when local MI analysis
  measures the candidate file below the absolute file-MI threshold.
- [x] T023-640 Move selected doctor declarations into behavior-named command
  files.
- [x] T023-650 Run `gofmt` on changed Go files.
- [x] T023-660 Run focused Go verification for `cmd/sdp-trace`.
- [x] T023-670 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T023-680 Run CRAP and MI quality gates without changing MI baselines.
- [x] T023-690 Run three independent reviewer lanes and record Slice 10
  evidence in
  `specs/023-numbered-code-file-locality-cleanup/reviews/slice-10-evidence.md`.
