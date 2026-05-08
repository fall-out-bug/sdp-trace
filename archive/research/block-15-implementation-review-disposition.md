# Block 15 Implementation Review Disposition

Date: 2026-05-06

## Review Runs

- MiniMax-M2.7 strict multi-file review: stopped as hung/empty output after
  repeated waits; not counted as review evidence.
- ZAI/GLM strict multi-file review: unavailable because no API key was
  configured; not counted as review evidence.
- Kimi K2P6 one-file micro-review over `internal/checkpoint/checkpoint.go`:
  found critical/major issues in `VerifySet` signature/digest coverage, run-set
  binding, and signer policy binding.
- DeepSeek V4 Flash multi-file security/platform review: found critical issues
  in sequence previous-digest enforcement and missing nonce handling.
- DeepSeek V4 Flash re-review after fixes: no remaining critical or major
  findings for checkpoint core.
- DeepSeek V4 Flash PR-level review: found one major envelope-validation gap;
  fixed and re-verified.

## Findings And Disposition

| Severity | Finding | Disposition |
|---|---|---|
| critical | `VerifySet` checked sequence linkage without verifying checkpoint signature, payload digest, or run binding. | Fixed by making `VerifySet` call `Verify` for each checkpoint before accepting sequence linkage. |
| critical | `VerifySet` did not enforce all checkpoints belong to the same run. | Fixed with run id consistency checks. |
| critical | `Create` allowed `sequence > 0` without `previous_checkpoint_digest`. | Fixed with `validateSequenceLink` in create and verify paths. |
| critical | Missing run nonce could become an empty-string binding pass. | Fixed by making payload derivation fail when `recorder_attached.run_nonce` is missing. |
| major | Policy signer id could be accepted without public-key binding. | Fixed by requiring policy public key in schema and returning `cannot_verify` when code receives a policy signer without one. |
| major | Provided key pair was not checked for public/private consistency before signing. | Fixed with explicit public/private key consistency validation. |
| major | `Verify` did not reject unsupported checkpoint `schema_version` before trusting overlapping struct fields. | Fixed with envelope validation for schema version, profile, hash algorithm, and canonicalization. |

## Residual Minor Items

- CLI verifies a single checkpoint artifact. Set-level sequence verification is
  implemented in Go and covered by tests, but no multi-checkpoint CLI mode is
  exposed in this block.
- Replay freshness remains `not_assessed` because no external timestamp,
  transparency log, or customer witness profile exists in Block 15.

No remaining critical or major findings are known after the re-review,
PR-level review fix, and fresh local verification.
