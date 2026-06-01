# Spec 021 Slice 2 Plan Review

Date: 2026-06-01.

Scope: `cmd/sdp-trace` command-surface usage-drift helpers.

## Decision Gate

- Simpler/Faster: A single `command_surface_usage_drift.go` file would remove
  more numbered shards with fewer files.
- Blocking Edge Cases: The single-file shape measured below the absolute file
  MI threshold and would force a mixed code/baseline PR, which CI rejects.
- Existing Open Source: Not applicable; this is source locality cleanup inside
  existing Go package code, with no new parser, workflow engine, or dependency.

## Selected Slice

Group usage-drift helpers into two cohesive files:

- `cmd/sdp-trace/command_surface_usage_collection.go`
- `cmd/sdp-trace/command_surface_usage_diff.go`

Rejected alternative:

- One combined drift file: rejected because analysis measured file MI below
  threshold.

Review result: ready to implement. The slice is bounded, behavior-preserving in
intent, and does not require a baseline change.
