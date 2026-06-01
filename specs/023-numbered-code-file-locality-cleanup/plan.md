# Plan: Numbered Code File Locality Cleanup

Status: in_progress

## Workstreams

### WS-023-A: Inventory And Slice Selection

Owned files:

- active numbered Go source files under `cmd`, `internal`, and `tools`

Deliverable:

- Select bounded slices by cohesive command family or package behavior.

### WS-023-B: Behavior-Named Grouping

Owned files:

- selected numbered source shards for the active slice

Deliverable:

- Move related functions from numbered shards into behavior-named files while
  preserving package boundaries and public behavior.

### WS-023-C: Verification And Evidence

Owned files:

- active slice source files
- active slice review and evidence artifacts

Deliverable:

- Verify behavior, CRAP, MI, drift, and reviewer evidence before commit.

## Slice 1

Slice 1 groups `cmd/sdp-trace` release-proof command shards into:

- `cmd/sdp-trace/release_proof_run.go`
- `cmd/sdp-trace/release_proof_args.go`
- `cmd/sdp-trace/release_proof_policy.go`

Single-file release-proof grouping was rejected because it failed the absolute
file-MI threshold.

## Slice 2

Slice 2 groups `cmd/sdp-trace` observe command adapter and exit-policy shards
into:

- `cmd/sdp-trace/observe_command_adapters.go`
- `cmd/sdp-trace/observe_exit_policy.go`

Single-file observe grouping was rejected because it failed the absolute
file-MI threshold.

## Slice 3

Slice 3 groups `cmd/sdp-trace` envelope summarize command shards into:

- `cmd/sdp-trace/envelope_summary_run.go`
- `cmd/sdp-trace/envelope_summary_args.go`

Single-file envelope grouping was rejected because it failed the absolute
file-MI threshold.

## Slice 4

Slice 4 groups `cmd/sdp-trace` export dispatcher shards into:

- `cmd/sdp-trace/export_command.go`

The slice is bounded to command dispatch and predicate helpers; telemetry and
cross-repo posture command implementations remain outside this slice.

## Slice 5

Slice 5 groups `cmd/sdp-trace` fixture validation shards into:

- `cmd/sdp-trace/fixture_validation_run.go`
- `cmd/sdp-trace/fixture_validation_args.go`
- `cmd/sdp-trace/fixture_expectation_policy.go`

Single-file fixture validation grouping and a runner+root-arg grouping were
rejected because they failed the absolute file-MI threshold.

## Slice 6

Slice 6 groups `cmd/sdp-trace` interaction command shards into:

- `cmd/sdp-trace/interaction_command.go`
- `cmd/sdp-trace/interaction_relay.go`
- `cmd/sdp-trace/interaction_relay_args.go`
- `cmd/sdp-trace/interaction_transcript_import.go`
- `cmd/sdp-trace/interaction_transcript_import_args.go`
- `cmd/sdp-trace/interaction_summary.go`
- `cmd/sdp-trace/cli_flag_requirements.go`

The transcript import grouping was split between import execution and argument
parsing after the combined file measured below the absolute file-MI threshold.
The `requireOnlyFlagsCode` helper moves to a behavior-named CLI requirements
file because it is already shared outside interaction commands.

## Verification

```text
gofmt -w <changed-go-files>
go test ./cmd/sdp-trace
go test ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
jq empty schema/*.json
git diff --check
```

CRAP and MI gates are required before any PR claim. If a consolidated file
creates a new MI-baseline entry or stale ratchet behavior, split the slice more
cohesively or move baseline changes into a separate reviewed PR.
