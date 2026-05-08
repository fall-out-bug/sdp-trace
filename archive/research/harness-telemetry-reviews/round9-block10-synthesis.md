# Round 9: Block 10 Review Synthesis

Status: review synthesis; corrected planning draft
Date: 2026-05-05

Inputs:

- `round9-block10-glm.md`
- `round9-block10-minimax.md`
- `round9-block10-kimi.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-design.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-implementation-plan.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/10-agentic-sdlc-v0-demo-readiness.md`

This synthesis is not product closure evidence and not a trusted release
claim.

## Result

Block 10 remains suitable for implementation handoff after correction.

Critical review themes were valid and have been incorporated into the
Block 10 design and implementation plan:

- transparent wrapper requires PTY support;
- event hash preimage and canonicalization must be fixed before coding;
- file mutation observation is VCS/workspace snapshot based in V0;
- adapter socket lifecycle needs drain, late-event, and cleanup rules;
- local-only trace wording must not imply external integrity;
- adapter and signer trust upgrades require authority policy;
- CI contract pinning must be outside the assessed agent's writable
  source tree;
- parallel workers need a canonical event type owner;
- demo scripts need structured local `not_assessed` injection.

## Valid Critical Findings Closed Into Spec

| Source | Finding | Resolution |
| --- | --- | --- |
| GLM | Missing PTY architecture | Added transparent wrapper mechanism and PTY tests. |
| GLM | Undefined event hash preimage | Added RFC 8785 canonicalization and hash rules. |
| GLM | No file mutation mechanism | Defined VCS/workspace snapshot observation. |
| GLM | Adapter socket races | Added socket lifecycle, drain, late-event, cleanup behavior. |
| MiniMax | Local keys overclaim integrity | Reworded local traces as honest-recorder sequence continuity. |
| MiniMax | First event lacks external anchor | Added local nonce disclosure and witness requirement. |
| MiniMax | Weak human signing profile | Downgraded unprotected local keys. |
| MiniMax | Ambiguous CI pinning | Required external/non-agent-writable pinning for gate-grade. |
| Kimi | No canonical event type owner | Slice A now owns schema-bound Go event mapping. |
| Kimi | No structured `not_assessed` injection | Added `sdp-trace observe`. |

## Remaining Major Implementation Risks

These are not blockers for the first milestone, but they must be handled
inside implementation slices:

1. Adapter capability handling versus optional contract observers needs
   verifier rules.
2. Legacy Node/Bash active paths must be removed or migrated to Go in
   this batch.
3. OpenCode and GSD live integration remains `not_assessed` until tested
   on the actual demo machine.

Resolved before handoff:

- first-milestone run artifact layout;
- first-milestone event payload fields;
- first-milestone correlation id domains;
- first-milestone validation ordering.
- Go-only target architecture and no new Node/npm tooling.

## Implementation Start Gate

Implementation may start with the first milestone:

```text
local wrap -> event chain -> verify -> missing evidence -> explain -> tamper fail
```

Do not start by building the OpenCode/GSD live demo. The demo should
consume the trace substrate after the local wrapper and verifier are
working.

## Recommended Parallel Launch

Start these slices first:

- Slice A: schema and canonical event type mapping.
- Slice B: CLI skeleton and transparent wrapper, using Slice A types as
  soon as available.
- Slice E: verifier against committed fixtures.

Then launch:

- Slice C: redaction/retention.
- Slice D: adapter socket.
- Slice F: query/export.
- Slice G: CI witness.
- Slice H: OpenCode/GSD/Bazel/Kotlin demo integration.

This order reduces merge conflict risk and prevents the demo from
inventing evidence semantics before the verifier exists.
