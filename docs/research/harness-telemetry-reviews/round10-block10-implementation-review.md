# Round 10: Block 10 Implementation Review

Status: implementation review ledger
Date: 2026-05-05

This is a review/disposition artifact. It is not source-bound proof,
not product closure evidence, and not a trusted release claim.

## Scope

Reviewed first milestone implementation:

```text
local wrap -> event chain -> verify -> missing evidence -> explain -> tamper fail
```

Reviewed implementation paths:

- `cmd/sdp-trace/`
- `internal/trace/`
- `internal/recorder/`
- `internal/verifier/`
- `internal/query/`
- `internal/export/`
- `examples/agentic-sdlc/`
- Block 10 design and implementation docs

## Initial Review Findings

### GLM / Platform Review

Initial unresolved major/minor findings:

- dead second event model;
- non-deterministic timestamp-based source snapshot digest;
- inconsistent verifier error contract;
- dry-run silently falling back to default contract on contract load
  failure;
- unused verifier artifact reader;
- dual hash preimage policy risk;
- flag parser gaps;
- socket env exposed without socket;
- non-signal strings in signal field;
- unbounded copy helper;
- missing tests for `run`, non-zero exit, and `cannot_verify`.

Disposition: accepted where valid. Invalid or future-scope items were
not carried into first milestone.

### MiniMax / CISO Review

Initial valid findings:

- source snapshot digest must not be a timestamp nonce;
- local-only replayability must not overclaim;
- contract digest must be checked during verification;
- fixture expectations must not live inside mutable run directories.

Disposition: accepted.

An intermediate MiniMax pass raised broad findings that `wrap` executes
user commands. Those were rejected as invalid against the product
boundary: `sdp-trace wrap` intentionally executes the exact
user-provided command and records local evidence.

### Kimi / Micro-Review

Valid findings:

- dry-run with `--use-default-contract` was broken;
- boolean flag parsing was incomplete;
- unknown flags were not rejected;
- dry-run usage errors used stdout;
- explicit output directories could be clobbered;
- `verify` could create missing run directories;
- default contract digest was not verified;
- `runQuery` used a magic exit code;
- dead `flagSet.out`;
- subcommand help behavior was unclear.

Disposition: accepted and fixed.

## Corrections Applied

- Removed Node/npm/JS/`.mjs` artifacts and shell scripts from the active
  product path.
- Added Go module and first milestone Go CLI.
- Added local event chain model, deterministic canonicalization, payload
  digest, and event hash verification.
- Removed the unused second event model.
- Replaced timestamp source digest with explicit local source snapshot
  state.
- Added output directory clobber protection.
- Removed fake `SDP_TRACE_SOCKET` export until adapter socket exists.
- Added strict flag parsing and boolean literal handling.
- Fixed dry-run contract handling and stderr behavior.
- Added contract digest verification for explicit and default contracts.
- Moved fixture expectations outside run directories.
- Added tests for:
  - `wrap -> verify -> query -> explain`;
  - `run` command;
  - dry-run;
  - fixture validation;
  - expected negative fixture;
  - unexpected negative fixture;
  - flag parser errors;
  - non-zero command exit propagation;
  - tamper detection;
  - missing evidence;
  - `cannot_verify` manifest/event load paths.

## Final Review Result

Final pi review:

- GLM: `NO CRITICAL OR MAJOR FINDINGS`
- MiniMax: `NO CRITICAL OR MAJOR FINDINGS`
- Kimi: `NO CRITICAL OR MAJOR FINDINGS`

## Final Verification Commands

Commands run after corrections:

```text
go test ./...
jq empty schema/*.json
go run ./cmd/sdp-trace --help
go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc
go run ./cmd/sdp-trace verify examples/agentic-sdlc/local-wrap-positive
go run ./cmd/sdp-trace explain examples/agentic-sdlc/local-wrap-positive
git diff --check
```

The final answer must report actual command results from the current
workspace, not this ledger as proof authority.
