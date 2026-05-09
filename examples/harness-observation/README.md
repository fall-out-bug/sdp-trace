# Harness Observation Fixtures

These fixtures document the Block 31 harness observation contract.

`sdp-trace harness observe` reads explicit local JSONL exports only. It does not
run OpenCode, GSD, provider APIs, GitHub APIs, or any harness runtime.

The first implementation is validated primarily through focused Go tests because
event `source_digest` binds to canonical event JSON. Committed fixture metadata
keeps the expected state matrix visible without making hand-authored proof JSON
authoritative.

## Expected Scenarios

- `generic-complete`: required harness and model families are present; validation
  state is `pass`.
- `zero-event-source`: a supplied source contains no events; required dimensions
  are `not_assessed`.
- `unsafe-raw-prompt`: an event contains a forbidden raw prompt field; observe
  fails before writing a run.
- `source-digest-mismatch`: an event source digest does not match canonical event
  JSON; observe fails before writing a run.

OpenCode/GSD remains a profile exemplar, not a product dependency. Missing
OpenCode/GSD export dimensions remain `not_assessed` or `cannot_verify`.
