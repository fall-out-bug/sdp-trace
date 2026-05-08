# Checkpoint Fixtures

These files are schema and UX fixtures for the local Ed25519 checkpoint profile.
They are not source-bound release proof and do not establish protected or
audit-grade trust.

- `signed-checkpoint.json`: positive local signed checkpoint shape.
- `trusted-checkpoint-policy.json`: local development signer allowlist.
- `checkpoint-verification-result.json`: local signed verification result.
- `negative-tampered-checkpoint-verification-result.json`: expected verifier
  states after payload tampering.
