# Spec 023: Numbered Code File Locality Cleanup

Status: in_progress

## Objective

Replace remaining numbered Go source shards such as `family_123_name.go` with
cohesive behavior-named files across active code packages while preserving
behavior, tests, quality gates, and package boundaries.

## Background

Spec 021 removed the `main_[0-9]*.go` command-surface shards but repository
inventory still shows numbered Go files in active product code:

- `cmd/sdp-trace`: 430 files
- `internal/harnessobs`: 350 files
- `internal/packet`: 200 files
- `internal/prreview`: 192 files

The total remaining active numbered Go file count at intake is 1172.

## Requirements

- FR-023-001: Split cleanup by package, command family, or cohesive behavior;
  no repo-wide rename sweep.
- FR-023-002: Preserve current public command behavior and output contracts.
- FR-023-003: Group functions into cohesive files named after behavior,
  command family, or domain concept.
- FR-023-004: Keep complexity, MI, and CRAP gates as verification signals, not
  as a reason for metric-gaming file moves.
- FR-023-005: Preserve package boundaries and dependency direction from
  `docs/package-ownership-map.md`.

## Non-Goals

- No command behavior change.
- No package split unless another reviewed spec owns it.
- No broad formatting churn outside touched packages.
- No production trust, release approval, or external attestation claim.

## Acceptance Criteria

- Each cleanup slice lists the exact package or command family it touches.
- Tests and quality gates pass after each slice.
- Each slice has three independent staged-diff reviewer lanes before commit.
- Remaining active numbered Go file count decreases monotonically.

## Active Slice 1

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` release-proof command shards only.

Files selected for grouping:

- `cmd/sdp-trace/release_151_run.go`
- `cmd/sdp-trace/release_154_parseargs.go`
- `cmd/sdp-trace/release_155_flagsandexits.go`

Target files:

- `cmd/sdp-trace/release_proof_run.go`
- `cmd/sdp-trace/release_proof_args.go`
- `cmd/sdp-trace/release_proof_policy.go`

Rejected grouping:

- A single `release_proof_command.go` file was rejected because local
  pre-change MI analysis measured file MI `65.3`, below the absolute threshold.
- Combining args and exits was rejected because local pre-change MI analysis
  measured file MI `69.2`, below the absolute threshold.

Intended behavior boundary: this slice should only move release-proof command
runner, argument parsing, required flags, and exit-code declarations into
behavior-named files. No CLI behavior, JSON field, schema contract, or command
metadata value should change.

## Active Slice 2

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` observe command adapter and exit policy shards only.

Files selected for grouping:

- `cmd/sdp-trace/observe_012_runcommand.go`
- `cmd/sdp-trace/observe_013_runharness.go`
- `cmd/sdp-trace/observe_028_harnessstateexits.go`

Target files:

- `cmd/sdp-trace/observe_command_adapters.go`
- `cmd/sdp-trace/observe_exit_policy.go`

Rejected grouping:

- A single observe command file was rejected because local pre-change MI
  analysis measured file MI `64.1`, below the absolute threshold.

Intended behavior boundary: this slice should only move observe command adapter
functions and harness state exit-code policy into behavior-named files. No CLI
behavior, JSON field, schema contract, or command metadata value should change.

## Active Slice 3

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` envelope summarize command shards only.

Files selected for grouping:

- `cmd/sdp-trace/envelope_173_run.go`
- `cmd/sdp-trace/envelope_174_requiredflags.go`
- `cmd/sdp-trace/envelope_174_writeoptionaljsonfile.go`
- `cmd/sdp-trace/envelope_175_parsesummarizeargs.go`

Target files:

- `cmd/sdp-trace/envelope_summary_run.go`
- `cmd/sdp-trace/envelope_summary_args.go`

Rejected grouping:

- A single `envelope_command.go` file was rejected because local MI analysis
  measured file MI `69.0`, below the absolute threshold.

Intended behavior boundary: this slice should only move the envelope summarize
runner, argument parsing, required flags, and optional JSON output helper into
behavior-named command files. No CLI behavior, JSON field, schema contract, or
command metadata value should change.

## Active Slice 4

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` export command dispatcher shards only.

Files selected for grouping:

- `cmd/sdp-trace/export_452_run.go`
- `cmd/sdp-trace/export_453_telemetryrequested.go`
- `cmd/sdp-trace/export_454_crossrepopostureexplainrequested.go`
- `cmd/sdp-trace/export_455_crossrepoposturerequested.go`

Target file:

- `cmd/sdp-trace/export_command.go`

Intended behavior boundary: this slice should only move the export dispatcher
and export subcommand predicates into a behavior-named command file. No CLI
behavior, JSON field, schema contract, or command metadata value should change.

## Active Slice 5

Status: implemented locally; targeted reviews LGTM; PR checks pending.

Scope: `cmd/sdp-trace` fixture validation command shards only.

Files selected for grouping:

- `cmd/sdp-trace/fixture_472_run.go`
- `cmd/sdp-trace/fixture_473_rootarg.go`
- `cmd/sdp-trace/fixture_474_validatefixtureruns.go`
- `cmd/sdp-trace/fixture_475_validatefixturerun.go`
- `cmd/sdp-trace/fixture_476_expectationfailed.go`
- `cmd/sdp-trace/fixture_477_expectedresultfailed.go`
- `cmd/sdp-trace/fixture_478_unexpectedresultfailed.go`

Target files:

- `cmd/sdp-trace/fixture_validation_run.go`
- `cmd/sdp-trace/fixture_validation_args.go`
- `cmd/sdp-trace/fixture_expectation_policy.go`

Rejected grouping:

- A single `fixture_validation_command.go` file was rejected because local MI
  analysis measured file MI `65.4`, below the absolute threshold.
- Keeping fixture root argument parsing in `fixture_validation_run.go` was
  rejected because local MI analysis measured file MI `68.6`, below the
  absolute threshold.

Intended behavior boundary: this slice should only move fixture validation
runner, root selection, per-run validation, and fixture expectation policy into
behavior-named command files. No CLI behavior, verifier artifact behavior, JSON
field, schema contract, or command metadata value should change.
