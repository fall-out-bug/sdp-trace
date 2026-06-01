# Plan: Source File Locality Cleanup

Status: in_progress

## Workstreams

### WS-021-A: Slice Selection

Owned files:

- `cmd/sdp-trace` command-surface registry helpers for Slice 1

Deliverable:

- Group the selected command-surface helper shards into
  `cmd/sdp-trace/command_surface_registry_helpers.go` and
  `cmd/sdp-trace/command_surface_metadata_helpers.go`.

### WS-021-B: Behavior-Named Grouping

Owned files:

- selected command-surface registry helper files

Deliverable:

- Move related functions from numbered shards into cohesive behavior-named
  files while preserving tests and public behavior.

### WS-021-C: Verification And Docs

Owned files:

- selected package docs if ownership or dependency direction changes
- `docs/package-ownership-map.md` when needed

Deliverable:

- Verify behavior and update ownership docs only when the cleanup changes
  package-level boundaries.

Slice 1 is expected not to change package boundaries or dependency direction,
so `docs/package-ownership-map.md` requires no content change unless review or
verification finds otherwise.

Slice 2 groups command-surface usage-drift helpers into collection and diff
files. A single combined drift file was rejected because it would fall below the
absolute file-MI threshold and force a mixed code/baseline PR.

Slice 3 groups only command-surface list helpers into
`cmd/sdp-trace/command_surface_list_helpers.go`. A broader argument-helper file
including `isHelp` and `isBoolLiteral` was rejected because local pre-change MI
analysis measured it below the absolute file-MI threshold.

Slice 4 groups command-surface schema type definitions and runner/JSON
functions into `command_surface_schema.go` and `command_surface_runner.go`.
Broader metadata/registry grouping was rejected because local pre-change MI
analysis measured candidate files below the absolute file-MI threshold.

Slice 5 groups the command-surface core command metadata shards into
`cmd/sdp-trace/command_surface_core_commands.go`. The combined file has no
functions and local pre-change MI analysis measured file MI `100.0`.

Slice 6 groups the command-surface observe command metadata shards into
`cmd/sdp-trace/command_surface_observe_commands.go`. The combined file has no
functions and local pre-change MI analysis measured file MI `100.0`.

Slice 7 groups the command-surface assess command metadata shards into
`cmd/sdp-trace/command_surface_assess_commands.go`. The combined file has no
functions and local pre-change MI analysis measured file MI `100.0`.

Slice 8 groups the command-surface packet command metadata shards into
`cmd/sdp-trace/command_surface_packet_commands.go`. The combined file has no
functions and local pre-change MI analysis measured file MI `100.0`.

## Verification

```text
gofmt -w <changed-go-files>
go test ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

CRAP and MI gates are required before any PR claim. If a consolidated file
creates a new MI-baseline entry or stale ratchet behavior, split the slice more
cohesively or move baseline changes into a separate reviewed PR.
