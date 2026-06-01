# Slice 8 Plan Review: Packet Command Metadata Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Target file: `cmd/sdp-trace/command_surface_packet_commands.go`
- Source shards:
  - `cmd/sdp-trace/main_563_commandsurfacepacketpacket.go`
  - `cmd/sdp-trace/main_564_commandsurfacepacketprreview.go`
  - `cmd/sdp-trace/main_575_commandsurfacepacket.go`

## Decision Gate

- Simpler/Faster: move the three packet metadata shards into one
  behavior-named file without changing values, command behavior, package
  boundaries, or tests.
- Blocking Edge Cases: accidental metadata edits, package init/order changes,
  MI/CRAP regression, or stale evidence could invalidate the cleanup claim.
- Existing Open Source: not applicable; this is a local Go file locality change
  using existing project verification tools and conventions.

## Review

- Requirements fit: pass. The slice is bounded to one command family and avoids
  a repo-wide rename sweep.
- UX impact: not_assessed. No CLI output or behavior change is intended.
- DX impact: pass. The target file removes three numbered shards from the active
  product path.
- Maintainability: pass. The grouped file keeps packet command metadata and the
  group list together.
- Trust boundary: pass. Production trust, release approval, merge approval, and
  external attestation remain not_assessed.

## Planned Verification

- `gofmt -w cmd/sdp-trace/command_surface_packet_commands.go`
- `go test ./cmd/sdp-trace`
- full repository test, vet, doccheck, hygienecheck, schema syntax, diff check
- coverage-backed CRAP and MI gates
- three independent reviewer lanes after implementation
