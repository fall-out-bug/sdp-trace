# Slice 4 Adversarial Review (GLM): WS-017-C/D Policy and Supply-Chain Prototypes

Reviewer: GLM adversarial pass
Date: 2026-05-20
Scope: 7 files across `examples/oss-policy/`, `examples/oss-supply-chain/`, and `docs/`

## Verdict: 3 Important, 3 Advisory — not LGTM

---

## Findings

### F1. Important | UX | Cosign block changes CWD without subshell isolation

**File:** `examples/oss-supply-chain/README.md:35`

```bash
cd /tmp
cosign generate-key-pair
```

`cd /tmp` is a bare CWD change in a documentation command block. AGENTS.md line 44 requires:

> Scanner verification commands in docs must be copy-pasteable and reproducible. If a command changes working directory or depends on local state, use subshell isolation (e.g., `(cd /tmp && ...)`).

**Fix:** Wrap the entire Cosign block in a subshell:

```bash
(cd /tmp && \
  cosign generate-key-pair && \
  echo '{"run":"test"}' > run.json && \
  cosign sign-blob --key cosign.key --yes run.json > run.json.sig && \
  cosign verify-blob --key cosign.pub --signature run.json.sig run.json)
```

---

### F2. Important | UX | SLSA verifier command uses repo-root-relative path after `cd /tmp`

**File:** `examples/oss-supply-chain/README.md:44-48`

```bash
slsa-verifier verify-artifact \
  --provenance-path examples/oss-supply-chain/local-dsse.json \
  --source-uri local/test \
  /dev/null
```

A reader who follows the Cosign block above (or the in-toto block referencing `/tmp`) will have CWD `/tmp`. The relative path `examples/oss-supply-chain/local-dsse.json` resolves nowhere from `/tmp`.

**Fix:** Either use an absolute placeholder, a subshell with explicit CWD, or add a comment noting the assumed repo-root CWD:

```bash
# From repo root:
slsa-verifier verify-artifact \
  --provenance-path examples/oss-supply-chain/local-dsse.json \
  --source-uri local/test \
  /dev/null
```

Or wrap with subshell: `(cd "$REPO_ROOT" && slsa-verifier verify-artifact ...)`

---

### F3. Important | UX | in-toto command references key that does not exist

**File:** `examples/oss-supply-chain/README.md:24-30`

```bash
in-toto-run \
  --step-name test-wrap \
  --products /dev/null \
  --key /tmp/test-key.pem \
  -- /bin/true
```

No prior step generates `/tmp/test-key.pem`. The copy-paste reader gets an error: key file not found. The Cosign block generates its own key (`cosign generate-key-pair`), but the in-toto block does not.

**Fix:** Add a key generation step before the in-toto block:

```bash
# Generate a throwaway key for local testing
openssl genpkey -algorithm RSA -out /tmp/test-key.pem 2>/dev/null
```

Or use `in-toto-keygen` if available.

---

### F4. Advisory | Quality | Non-canonical state `local_pass` in example READMEs

**Files:** `examples/oss-policy/README.md:3`, `examples/oss-supply-chain/README.md:3`

Both use `Status: \`local_pass\``. The canonical verifier states per AGENTS.md:49 are `pass`, `fail`, `cannot_verify`, `not_assessed`. `local_pass` is not canonical.

Context: these are document-level metadata, not verifier output, and the `local_` prefix distinguishes them from external `pass`. However, the backtick formatting presents them as machine-readable states, which could confuse a reader scanning for canonical states.

**Recommendation:** Use `not_assessed` (since there is no external verifier run), or change formatting to prose: `Status: locally tested, not externally verified`.

---

### F5. Advisory | Quality | Incomplete negative-path failure reason for SLSA verifier

**File:** `examples/oss-supply-chain/README.md:51`

> Expected: failure because no Rekor entry exists.

The DSSE fixture at `local-dsse.json:13` has a truncated signature (`"sig": "MEUCIQDtVdxRBBT..."`). The SLSA verifier will fail for multiple reasons: invalid signature, missing Rekor entry, untrusted key. Stating only "no Rekor entry" is incomplete.

**Recommendation:** Expand to: `Expected: failure (truncated signature, no Rekor entry, untrusted key).`

---

### F6. Advisory | DX | Minor article error in doc prose

**Files:** `docs/oss-policy-prototype.md:8`, `examples/oss-policy/README.md:6`

Both say "an simplified adapter-capture profile" — should be "a simplified adapter-capture profile."

---

## Correct (with evidence)

- **No overclaim.** Both `docs/oss-policy-prototype.md` and `docs/oss-supply-chain-prototype.md` explicitly state "does not approve replacing sdp-trace" semantics. Non-goals sections are thorough and honest.
- **Canonical verifier states in probe tables.** All probe-result tables use `pass`, `fail`, `cannot_verify`, `not_assessed` exclusively.
- **Honest `not_assessed`.** `docs/oss-supply-chain-prototype.md` marks "SLSA/Rekor production trust" as `not_assessed` with concrete reason: "No live external provenance or Rekor inclusion."
- **Valid DSSE payload.** Base64-decoded `local-dsse.json` payload is a structurally valid in-toto SLSA v1 statement with intentionally fake digest (`abcdef...`) appropriate for a negative-path fixture.
- **Correct Rego policy.** `adapter.rego` uses valid Rego v1 syntax. Default `pass := false`, rule body checks `trace_id != ""` and `count(input.provenance) <= 3`. Test fixture has non-empty trace_id and 2 provenance entries, so `pass = true` is correct.
- **No secrets in fixtures.** `local-dsse.json` contains a truncated fake signature and a local key ID. No private key material.
- **Substitution boundaries are clearly drawn.** Both example READMEs and both docs have explicit "What replaces / What needs glue / What remains sdp-trace-specific" sections.
- **Spec links are valid.** Relative links to `../../specs/017-oss-replacement-compatibility-and-benchmarks/` and `../specs/...` resolve correctly.

---

## Security

No findings. No secrets, no private key material, no unsafe commands, no misleading trust claims in the documentation files. The `local-dsse.json` truncated signature is appropriate for a negative-path test fixture.

---

## Summary

| # | Severity | Axis | File:Line | Summary |
|---|----------|------|-----------|---------|
| F1 | Important | UX | `examples/oss-supply-chain/README.md:35` | `cd /tmp` without subshell isolation violates AGENTS.md:44 |
| F2 | Important | UX | `examples/oss-supply-chain/README.md:44-48` | SLSA verifier path breaks if CWD changed by earlier block |
| F3 | Important | UX | `examples/oss-supply-chain/README.md:24-30` | in-toto command references non-existent key file |
| F4 | Advisory | Quality | `examples/oss-policy/README.md:3`, `examples/oss-supply-chain/README.md:3` | Non-canonical `local_pass` state |
| F5 | Advisory | Quality | `examples/oss-supply-chain/README.md:51` | Incomplete negative-path failure reason |
| F6 | Advisory | DX | `docs/oss-policy-prototype.md:8`, `examples/oss-policy/README.md:6` | "an simplified" → "a simplified" |

**3 Important, 3 Advisory. Not LGTM. Fix F1–F3 for copy-pasteability compliance, then re-review.**
