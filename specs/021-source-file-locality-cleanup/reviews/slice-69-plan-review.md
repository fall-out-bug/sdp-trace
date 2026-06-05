# Slice 69 Plan Review

Status: in_progress
Date: 2026-06-04

## Scope

Slice 69 is bounded to numbered `cmd/sdp-trace` `pr_review` top-level dispatch
and packet subcommand setup shards:

- `cmd/sdp-trace/pr_review_030_handlers.go`
- `cmd/sdp-trace/pr_review_037_packetrequiredflags.go`
- `cmd/sdp-trace/pr_review_038_packetstringflags.go`
- `cmd/sdp-trace/pr_review_039_run.go`
- `cmd/sdp-trace/pr_review_098_runpacket.go`
- `cmd/sdp-trace/pr_review_099_parsepacketargs.go`
- `cmd/sdp-trace/pr_review_100_registerpacketflags.go`
- `cmd/sdp-trace/pr_review_101_buildpacket.go`
- `cmd/sdp-trace/pr_review_102_packetoptions.go`
- `cmd/sdp-trace/pr_review_103_fillpacketidentity.go`
- `cmd/sdp-trace/pr_review_104_fillpacketevidence.go`
- `cmd/sdp-trace/pr_review_138_requirepacketinputs.go`

## Decision Gate

- Simpler/Faster: rename/move only; no behavior change, public API change,
  dependency change, or baseline update.
- Blocking Edge Cases: a combined top-level command/packet flag file failed
  pre-change file MI at `56.6`; a combined packet command/options file failed
  pre-change file MI at `68.1`. Packet build failures, required provenance
  anchors, metadata mapping, and repeated evidence flags are trust-sensitive
  and need focused regression evidence.
- Existing Open Source: no new CLI parser, flag framework, packet builder, or
  review workflow dependency is needed; existing package-local helpers and the
  internal `prreview` package are sufficient.

## Planned File Boundary

- `cmd/sdp-trace/pr_review_command.go`: top-level `pr-review` handler map and
  command dispatch.
- `cmd/sdp-trace/pr_review_packet_flags.go`: packet required flag metadata and
  string flag defaults.
- `cmd/sdp-trace/pr_review_packet_args.go`: packet subcommand parsing,
  positional-argument rejection, and required-input checks.
- `cmd/sdp-trace/pr_review_packet_run.go`: packet build execution and stdout
  rendering.
- `cmd/sdp-trace/pr_review_packet_options.go`: packet option construction,
  provenance anchor mapping, repeated evidence flag reconstruction, and
  `prreview.BuildPacket` handoff.

## Planned Regression Evidence

- Exact focused test existence:
  - `TestPRReviewHandlersKeepSubcommands`
  - `TestPRReviewPacketFlagsKeepContract`
  - `TestParsePRReviewPacketArgsKeepsUsageBoundaries`
  - `TestPRReviewPacketOptionsKeepsEvidenceMapping`
- Focused behavior checks: subcommand routing, usage and missing-subcommand
  diagnostics, required packet flags and defaults, positional-argument
  rejection, missing packet-anchor usage errors, packet build failure
  `cannot_verify` behavior, repeated context and verification flag
  reconstruction, optional metadata mapping, and provenance anchor mapping.
- Standard verification, conditional `golangci-lint run`, and CRAP/MI gates
  remain required before the implementation review.

## Review Rounds

### Round 1

- scope/correctness: minor finding. The regression plan did not explicitly
  cover `--metadata` CLI-to-packet option mapping.
- trust/evidence: `LGTM`
- maintainability/DX: major finding. The planned `pr_review_packet_command.go`
  boundary owned required-input checks but excluded
  `pr_review_138_requirepacketinputs.go`, leaving packet command ownership
  split across a behavior-named file and a stale numbered shard.

### Round 2

- scope/correctness: `LGTM`
- maintainability/DX: `LGTM`

Final plan-review status: `LGTM` across scope/correctness, trust/evidence,
and maintainability/DX after fixes.
