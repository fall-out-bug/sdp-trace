# Round 9: Block 10 GLM Review

Status: pi review output summary
Date: 2026-05-05
Model: `zai/glm-5.1`
Role: Platform / Harness Owner

Reviewed:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-design.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-implementation-plan.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-demo-readiness.md`

This is a review artifact, not source-bound proof or closure evidence.

## Critical Findings

### C1: No pseudoterminal architecture for transparent wrapper

The design required TTY behavior, colors, signals, and interactive tool
compatibility, but the implementation plan did not require PTY support.
Testing only `/bin/echo` would hide the failure.

Disposition: accepted. Block 10 now requires PTY mode for interactive
stdio and explicit tests.

### C2: Event hash construction undefined

The documents listed `previous event hash` and `event hash` but did not
define the exact hash preimage or canonicalization method.

Disposition: accepted. Block 10 now specifies RFC 8785 canonical JSON
and excludes `event_hash` from the hash preimage.

### C3: No mechanism for `file_mutation_observed`

The event was required, but the recorder is only a process wrapper and
the design did not specify filesystem watcher, VCS snapshot, or syscall
interception.

Disposition: accepted. V0 now defines file mutation observation as
pre/post VCS/workspace snapshot comparison, not syscall-level capture.

### C4: Socket lifecycle races with child process exit

The adapter socket lifecycle did not define drain behavior, late events,
partial frames, or shutdown ordering.

Disposition: accepted. Block 10 now defines per-run socket directory,
bounded drain, late message behavior, parse failure states, and cleanup
scope.

## Major Findings

- `wrap` vs `run` contract boundary unclear.
- Contract lock mechanics missing.
- Schema and CLI slices risked rework through provisional structs.
- `command_finished` event payloads were underspecified.
- Socket cleanup on crash/signal was unspecified.
- Fully unavailable external tools could make the demo vacuous.
- Observer id generation and role semantics were unclear.

Disposition: accepted in part. The Block 10 docs now define default
local contracts for `wrap`, contract lock mechanics, canonical type
ownership in Slice A, fixture replay mode, and stricter socket lifecycle.
Event-specific payload detail remains an implementation task for Slice A.

## Minor Findings

- Go module coexistence with current Node/bash validators needs care.
- `dry-run` output format should be made testable.
- Correlation id semantics are undefined.
- Absolute retention expiry examples can rot.
- Adapter capabilities vs optional contract observers need verifier
  rules.
- `--tamper` must not destructively mutate canonical demo output.

Disposition: accepted in part. Tamper behavior and retention warning are
now explicit. Remaining items are retained as implementation details for
schema/verifier slices.
