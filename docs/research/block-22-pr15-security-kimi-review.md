**Verdict: BLOCK**

Three major findings remain that violate the spec or create safety regressions. No critical exploits are present, but the combination of a digest-length bug, a safety-sanitization bypass, and a schema regression prevents approval.

---

## Major

**M1 — `strongDigest` rejects SHA-384/512, contradicting “SHA-256-or-stronger”**
- **File/func:** `internal/witness/profiles.go:699-704` (`strongDigest`)
- **Spec evidence:** CLI Boundary requires `--customer-pki-payload-digest <sha256-or-stronger-digest>`; Safety/Artifact Binding requires “SHA-256-or-stronger digests.”
- **Issue:** `len(value) != 64` hard-codes SHA-256. A valid SHA-384 (96 hex chars) or SHA-512 (128 hex chars) payload digest is incorrectly treated as weak and causes `fail` / `witness_artifact_digest_mismatch`.
- **Required fix:** Change length check to `len(value) < 64` (or `>= 64` with valid hex) to accept any stronger digest.

**M2 — `finalizeRecordForWrite` copies unvalidated envelope fields into the “safe” record after detecting unsafe output**
- **File/func:** `internal/witness/profiles.go:540-551` (`finalizeRecordForWrite`)
- **Spec evidence:** Safety Requirements: “No raw token, OIDC JWT, CI secret, private key … may be persisted or printed;” `witness_unsafe_output_candidate` must produce `fail`.
- **Issue:** When `forbiddenOutputPresent(raw)` is true, the function rebuilds a minimal record but then copies `record.ProfileID`, `record.ProfileVersion`, and `record.ProviderKind` from the original tainted record. If the unsafe bytes were in `ProfileID` (or another copied field), the sanitized record still contains the secret and writes it to disk.
- **Required fix:** Remove the three assignments (`safe.ProfileID = record.ProfileID`, etc.) so the safe record uses only the closed defaults from `baseRecord`, or validate the copied strings against the closed enum before copying.

**M3 — New schema requires `profile_states` and `output_safety` for all kinds, but the GitHub Actions path is not updated in this diff**
- **File/line:** `schema/witness-profile-result.schema.json` required list includes `profile_states`, `output_safety`, `profile_id`, etc.; `cmd/sdp-trace/main.go:1373` still calls `witness.WriteGitHubActions`.
- **Spec evidence:** Schema defines required fields for all `kind` values including `github-actions`; normalized witness result contract requires `profile_states` and `output_safety` for every profile.
- **Issue:** The diff expands `Record` and the JSON schema with new required fields, but does not modify the GitHub Actions implementation (not present in the patch). If the existing `WriteGitHubActions` does not populate these fields, its output fails schema validation.
- **Required fix:** Either update the GitHub Actions witness builder in this PR to populate the new required fields, or scope the new schema/required list to Block-22-only profile kinds so the existing path remains valid.

---

## Minor

**m1 — `ReasonKeyCustodyNA` is declared but never emitted**
- **File:** `internal/witness/witness.go:59`
- **Issue:** The constant exists, but `BuildCustomerPKI` never returns `StatusNotAssessed` / `ReasonKeyCustodyNA` when key custody is missing; it silently defaults to `"not_assessed"`.

**m2 — Customer PKI public-key loader is hard-coded to Ed25519**
- **File/func:** `internal/witness/profiles.go:643-670` (`loadCustomerPublicKey`)
- **Spec evidence:** Customer PKI profile requires “declared public certificate or public key identity” with no algorithm restriction.
- **Issue:** Only Ed25519 keys/certificates are accepted. RSA/ECDSA customer certificates are rejected.

**m3 — `unsafeInputPath` does not block access to private directories**
- **File/func:** `internal/witness/profiles.go:588-596` (`unsafeInputPath`)
- **Spec evidence:** CLI Boundary: “refuse … private directories, private key paths, provider tokens, or customer directories.”
- **Issue:** The function rejects symlinks, traversal (`..`), URLs, and filenames containing `private.key`, but allows absolute paths such as `/private/…` or `/etc/…`.

**m4 — Schema root permits `additionalProperties`**
- **File:** `schema/witness-profile-result.schema.json:5`
- **Issue:** The witness result schema sets `"additionalProperties": true` at the root, which is looser than the closed profile contract described in the spec.
