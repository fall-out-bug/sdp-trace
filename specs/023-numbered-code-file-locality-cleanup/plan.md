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

## Slice 7

Slice 7 groups `cmd/sdp-trace` wrap, run, preview, and dry-run command shards
into:

- `cmd/sdp-trace/wrap_legacy.go`
- `cmd/sdp-trace/wrap_recorder.go`
- `cmd/sdp-trace/wrap_run.go`
- `cmd/sdp-trace/wrap_run_args.go`
- `cmd/sdp-trace/wrap_preview.go`
- `cmd/sdp-trace/wrap_preview_args.go`
- `cmd/sdp-trace/wrap_preview_payload.go`

The preview grouping was split between command flow, args/contract loading, and
payload rendering after the combined preview file measured below the absolute
file-MI threshold.

## Slice 8

Slice 8 groups `cmd/sdp-trace` query, verify, explain, and query-pack command
shards into:

- `cmd/sdp-trace/query_verify.go`
- `cmd/sdp-trace/query_verify_args.go`
- `cmd/sdp-trace/query_verify_exit.go`
- `cmd/sdp-trace/query_explain.go`
- `cmd/sdp-trace/query_run.go`
- `cmd/sdp-trace/query_dispatch.go`
- `cmd/sdp-trace/query_pack.go`
- `cmd/sdp-trace/query_pack_build.go`
- `cmd/sdp-trace/query_pack_explain.go`
- `cmd/sdp-trace/query_pack_args.go`
- `cmd/sdp-trace/query_pack_validation.go`

Verify, query, and query-pack were split along runner, argument, dispatch, and
validation responsibilities after combined files measured below the absolute
file-MI threshold.

## Slice 9

Slice 9 groups `cmd/sdp-trace` witness command shards into:

- `cmd/sdp-trace/witness_command.go`
- `cmd/sdp-trace/witness_options.go`
- `cmd/sdp-trace/witness_options_parse.go`
- `cmd/sdp-trace/witness_options_build.go`
- `cmd/sdp-trace/witness_flag_set.go`
- `cmd/sdp-trace/witness_required_fields.go`
- `cmd/sdp-trace/witness_kind_validation.go`
- `cmd/sdp-trace/witness_record_builders.go`
- `cmd/sdp-trace/witness_output.go`
- `cmd/sdp-trace/witness_customer_pki.go`

Options and required-field grouping were split after broader groupings measured
below the absolute file-MI threshold.

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
