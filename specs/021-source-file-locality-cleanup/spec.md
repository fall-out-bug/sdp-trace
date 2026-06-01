# Spec 021: Source File Locality Cleanup

Status: in_progress

## Objective

Replace transitional numbered one-function Go source shards with cohesive,
behavior-named files by package or command family while preserving command
behavior, tests, and quality gates.

## Background

Spec 018 recorded numbered source shards such as `family_123_name.go` as a
transitional organization artifact, not the target repository shape. This spec
turns that policy into bounded cleanup slices.

## Requirements

- FR-021-001: Split cleanup by package or command family; no repo-wide rename
  sweep.
- FR-021-002: Preserve current public command behavior and output contracts.
- FR-021-003: Group functions into cohesive files named after behavior,
  command family, or domain concept.
- FR-021-004: Keep complexity, MI, and CRAP gates as verification signals, not
  as a reason for metric-gaming file moves.
- FR-021-005: Preserve package boundaries and dependency direction from
  `docs/package-ownership-map.md`.

## Non-Goals

- No command behavior change.
- No package split unless another reviewed spec owns it.
- No broad formatting churn outside touched packages.
- No production trust, release approval, or external attestation claim.

## Acceptance Criteria

- Each cleanup slice lists the exact package or command family it touches.
- Tests and quality gates pass after each slice.
- `docs/package-ownership-map.md` remains accurate for touched packages.

## Active Slice 1

Status: implemented locally; post-implementation re-review pending.

Scope: `cmd/sdp-trace` command-surface registry helpers only.

Files selected for grouping:

- `cmd/sdp-trace/main_538_commandsurfaceregistrycore.go`
- `cmd/sdp-trace/main_539_commandsurfaceregistryobserve.go`
- `cmd/sdp-trace/main_540_commandsurfaceregistryassess.go`
- `cmd/sdp-trace/main_541_commandsurfaceregistryother.go`
- `cmd/sdp-trace/main_542_commandsurfaceregistrypacket.go`
- `cmd/sdp-trace/main_543_commandsurfaceregistry.go`
- `cmd/sdp-trace/main_580_commandsurfaceflaghelper.go`
- `cmd/sdp-trace/main_581_commandsurfaceposhelper.go`
- `cmd/sdp-trace/main_584_commandsurfacereqpos.go`

Target files:

- `cmd/sdp-trace/command_surface_registry_helpers.go`
- `cmd/sdp-trace/command_surface_metadata_helpers.go`

Intended behavior boundary: this slice should only consolidate tiny
generated-looking helper shards into cohesive behavior-named files. The
behavior-preserving claim remains `not_assessed` until post-change tests and
review evidence are recorded in
`specs/021-source-file-locality-cleanup/reviews/slice-1-evidence.md`.

## Active Slice 2

Status: implemented locally; targeted review LGTM; PR checks pending.

Scope: `cmd/sdp-trace` command-surface usage drift helpers only.

Files selected for grouping:

- `cmd/sdp-trace/main_576_commandsurfacedriftadd.go`
- `cmd/sdp-trace/main_585_commandsurfacedrift.go`
- `cmd/sdp-trace/main_586_commandsurfaceregistryusages.go`
- `cmd/sdp-trace/main_587_commandsurfacehelpusages.go`
- `cmd/sdp-trace/main_588_commandsurfacediffsets.go`
- `cmd/sdp-trace/main_589_commandsurfacesorteddiffs.go`
- `cmd/sdp-trace/main_590_commandsurfacedrifterror.go`
- `cmd/sdp-trace/main_591_commandsurfacedriftparts.go`

Target files:

- `cmd/sdp-trace/command_surface_usage_collection.go`
- `cmd/sdp-trace/command_surface_usage_diff.go`

Intended behavior boundary: this slice should only consolidate usage-drift
helpers into cohesive behavior-named files. The behavior-preserving claim
remains `not_assessed` until post-change tests and review evidence are recorded
in `specs/021-source-file-locality-cleanup/reviews/slice-2-evidence.md`.

## Active Slice 3

Status: implemented locally; targeted review LGTM; PR checks pending.

Scope: `cmd/sdp-trace` command-surface list helper shards only.

Files selected for grouping:

- `cmd/sdp-trace/main_582_commandsurfaceflaglisthelpers.go`
- `cmd/sdp-trace/main_583_commandsurfaceslicehelpers.go`

Target file:

- `cmd/sdp-trace/command_surface_list_helpers.go`

Intended behavior boundary: this slice should only consolidate simple list
helper functions into a cohesive behavior-named file. `isHelp` and
`isBoolLiteral` were not included because the combined argument-helper file
failed the absolute file-MI threshold during pre-change analysis.

## Active Slice 4

Status: implemented locally; targeted reviews LGTM; PR checks pending.

Scope: `cmd/sdp-trace` command-surface schema and runner shards only.

Files selected for grouping:

- `cmd/sdp-trace/main_536_commandsurfaceschema.go`
- `cmd/sdp-trace/main_544_commandsurfacejson.go`
- `cmd/sdp-trace/main_545_runcommandsurface.go`

Target files:

- `cmd/sdp-trace/command_surface_schema.go`
- `cmd/sdp-trace/command_surface_runner.go`

Intended behavior boundary: this slice should only move schema type definitions
and command-surface runner/JSON functions into behavior-named files. Metadata,
registry, command family definitions, and argument helper shards remain outside
this slice because broader grouping failed pre-change MI analysis.

## Active Slice 5

Status: implemented locally; targeted reviews LGTM; PR checks pending.

Scope: `cmd/sdp-trace` command-surface core command metadata shards only.

Files selected for grouping:

- `cmd/sdp-trace/main_548_commandsurfacecorebasic.go`
- `cmd/sdp-trace/main_549_commandsurfacecorewrap.go`
- `cmd/sdp-trace/main_550_commandsurfacecorerun.go`
- `cmd/sdp-trace/main_551_commandsurfacecorepreview.go`
- `cmd/sdp-trace/main_552_commandsurfacecoredoctor.go`
- `cmd/sdp-trace/main_553_commandsurfacecoreinstall.go`
- `cmd/sdp-trace/main_554_commandsurfacecorefixtures.go`
- `cmd/sdp-trace/main_571_commandsurfacecore.go`

Target file:

- `cmd/sdp-trace/command_surface_core_commands.go`

Intended behavior boundary: this slice should only move core command-surface
metadata values and the core command list into one behavior-named file. No CLI
behavior, JSON field, schema contract, or command metadata value should change.

## Active Slice 6

Status: implemented locally; targeted reviews LGTM; PR checks pending.

Scope: `cmd/sdp-trace` command-surface observe command metadata shards only.

Files selected for grouping:

- `cmd/sdp-trace/main_555_commandsurfaceobserveinteraction.go`
- `cmd/sdp-trace/main_556_commandsurfaceobserveobserve.go`
- `cmd/sdp-trace/main_557_commandsurfaceobserveharness.go`
- `cmd/sdp-trace/main_558_commandsurfaceobserveenvelope.go`
- `cmd/sdp-trace/main_572_commandsurfaceobserve.go`

Target file:

- `cmd/sdp-trace/command_surface_observe_commands.go`

Intended behavior boundary: this slice should only move observe command-surface
metadata values and the observe command list into one behavior-named file. No
CLI behavior, JSON field, schema contract, or command metadata value should
change.

## Active Slice 7

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` command-surface assess command metadata shards only.

Files selected for grouping:

- `cmd/sdp-trace/main_559_commandsurfaceassessassess.go`
- `cmd/sdp-trace/main_560_commandsurfaceassessreport.go`
- `cmd/sdp-trace/main_561_commandsurfaceassessgate.go`
- `cmd/sdp-trace/main_562_commandsurfaceassesscheckpoint.go`
- `cmd/sdp-trace/main_573_commandsurfaceassess.go`

Target file:

- `cmd/sdp-trace/command_surface_assess_commands.go`

Intended behavior boundary: this slice should only move assess command-surface
metadata values and the assess command list into one behavior-named file. No
CLI behavior, JSON field, schema contract, or command metadata value should
change.

## Active Slice 8

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` command-surface packet command metadata shards only.

Files selected for grouping:

- `cmd/sdp-trace/main_563_commandsurfacepacketpacket.go`
- `cmd/sdp-trace/main_564_commandsurfacepacketprreview.go`
- `cmd/sdp-trace/main_575_commandsurfacepacket.go`

Target file:

- `cmd/sdp-trace/command_surface_packet_commands.go`

Intended behavior boundary: this slice should only move packet command-surface
metadata values and the packet command list into one behavior-named file. No
CLI behavior, JSON field, schema contract, or command metadata value should
change.

## Active Slice 9

Status: implemented and pushed; targeted reviews LGTM; PR checks passed.

Scope: `cmd/sdp-trace` command-surface other command metadata shards only.

Files selected for grouping:

- `cmd/sdp-trace/main_565_commandsurfaceotherverify.go`
- `cmd/sdp-trace/main_566_commandsurfaceotherquery.go`
- `cmd/sdp-trace/main_567_commandsurfaceotherwitness.go`
- `cmd/sdp-trace/main_568_commandsurfaceotherrelease.go`
- `cmd/sdp-trace/main_569_commandsurfaceotheroverride.go`
- `cmd/sdp-trace/main_570_commandsurfaceotherexport.go`
- `cmd/sdp-trace/main_574_commandsurfaceother.go`

Target file:

- `cmd/sdp-trace/command_surface_other_commands.go`

Intended behavior boundary: this slice should only move other command-surface
metadata values and the other command list into one behavior-named file. No
CLI behavior, JSON field, schema contract, or command metadata value should
change.

## Active Slice 10

Status: implemented locally; targeted reviews LGTM; PR checks pending.

Scope: `cmd/sdp-trace` command-surface catalog metadata shards only.

Files selected for grouping:

- `cmd/sdp-trace/main_577_commandsurfaceprofiles.go`
- `cmd/sdp-trace/main_578_commandsurfacestates.go`

Target file:

- `cmd/sdp-trace/command_surface_catalog.go`

Intended behavior boundary: this slice should only move profile, witness-kind,
and state metadata values into one behavior-named catalog file. No CLI
behavior, JSON field, schema contract, or command metadata value should change.
