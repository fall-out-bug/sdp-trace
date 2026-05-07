**Verdict: APPROVE**

M2 is fixed. No remaining critical or major findings.

**Evidence that M2 is resolved**

- **`jwt_body` / base64-encoded JWT bodies:** `containsSecretLike` now calls `jwtLike`, which tokenizes input on JSON/string delimiters and detects standard three-segment base64 JWT shapes (`eyJ... . ... . ...`). This closes the specific demonstrated attack (raw JWT in `CI.Actor` or any parsed field).
- **`cloud_payload`, `pki_payload`, `free_text_parser_error_with_input`:** These literal safety-class markers are now present in both `containsSecretLike` and `forbiddenOutputPresent`.
- **Pre-parse input sanitization:** `readSafeJSON` rejects envelope, policy, and freshness files before unmarshaling if `containsSecretLike` fires.
- **Per-field envelope sanitization:** `unsafeEnvelopeFields` checks every parsed string field in `Source` and `CI` identity, plus all artifact paths/digests, via `unsafeOutputString` (which invokes `containsSecretLike` and the `@` / `/private/` guards).
- **Pre-write output sanitization:** `finalizeRecordForWrite` marshals the record (minus `OutputSafety`), runs `forbiddenOutputPresent` (which combines `containsSecretLike` + explicit markers + `jwtLike`), and rebuilds a minimal safe record with `StatusFail`/`ReasonUnsafeOutput` if anything is detected.
- **Defense in depth for user-supplied files:** `loadCustomerPublicKey` and `privateKeyInput` both invoke `containsSecretLike` on raw PEM/DER bytes, rejecting private-key material and JWT-shaped content before key loading.

The prior attack path—*malicious envelope containing a raw JWT in `CI.Actor` passes pre-parse validation, passes `unsafeEnvelopeFields`, passes pre-write output scan, and is persisted*—is now blocked at all three layers. No critical or major security findings remain.
