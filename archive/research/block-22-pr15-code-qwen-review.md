## Block 22 PR-Level Code/Correctness Review

### Verdict: **APPROVE**

No critical or major findings remain. The implementation faithfully implements the spec. Below are the minor findings.

---

### MINOR Findings

**m1 — `ProfileStates`/`OutputSafety` marked `omitempty` but schema `required`**

| File | Line | Evidence |
|---|---|---|
| `internal/witness/witness.go` | ~83 | `ProfileStates *ProfileStates \`json:"profile_states,omitempty"\`` and `OutputSafety *OutputSafety \`json:"output_safety,omitempty"\`` |
| `schema/witness-profile-result.schema.json` | line 16 | Both declared in `"required"` array |

All builder paths populate these before returning, so they are never actually omitted in practice. However, the Go struct tag (`omitempty`) contradicts the JSON Schema (`required`). If future code accidentally leaves one nil, the serialized output would silently violate the schema.

**Fix:** Remove `omitempty` from these two struct fields, or remove them from the schema `required` array.

---

**m2 — `"ci_same_job"` hardcoded instead of using a named constant**

| File | Line | Evidence |
|---|---|---|
| `internal/witness/profiles.go` | ~350 | `record.ProfileStates = defaultProfileStates(stateCannotVerify, "ci_same_job")` |

The code defines `independenceCIJob = "ci_isolated_job"` and `independenceExternal = "external_independent"` as named constants, but uses the bare string `"ci_same_job"` (also a valid `independence_state` enum value defined in the schema). Inconsistent with the constant convention.

**Fix:** Add `independenceSameJob = "ci_same_job"` constant and use it.

---

**m3 — No end-to-end CLI test exercising `--witness-envelope` flag**

| File | Evidence |
|---|---|
| `cmd/sdp-trace/main_test.go` | `TestWitnessCommandBuildkiteRequiresExplicitEnvelope` clears env and verifies `cannot_verify` without envelope, but no test passes `--witness-envelope <file>` and validates a successful CI profile result through the CLI path. |

The `--witness-envelope` flag is declared in `main.go` and plumbing exists through `ProfileOptions.EnvelopePath → BuildCIEnvelopeProfile`, but the integration layer has no test exercising a valid envelope file with the `--kind buildkite` or `--kind gitlab-ci` CLI path.

**Fix:** Add a CLI integration test that provides a valid `--witness-envelope` and asserts `exit 0` with `pass` status.

---

**m4 — Redundant double lowercasing in `containsSecretLike` / `forbiddenOutputPresent`**

| File | Evidence |
|---|---|
| `internal/witness/profiles.go` | `forbiddenOutputPresent` does `text := strings.ToLower(string(raw))` then calls `containsSecretLike(raw)` which does `lower := strings.ToLower(string(raw))` internally. |

Two allocations of the same `ToLower` string per output safety scan. Minor performance waste; not a correctness bug.

**Fix:** Pass `text` into `containsSecretLike` or have `containsSecretLike` return a version that avoids re-lowercasing.

---

### Summary

All spec-required behaviors are implemented and tested:
- Provider-neutral profile contract with closed reason codes ✓
- GitLab CI / Buildkite envelope normalization with topology capping ✓
- Customer PKI with Ed25519 freshness verification, revocation, key custody ✓
- Air-gapped excluded from `--kind` (docs+fixture only) ✓
- Three-layer input/output safety (pre-parse, per-field, pre-write) ✓
- JWT-shaped payload detection at all layers ✓
- Symlink/path-safety rejection ✓
- Environment-only cannot upgrade trust ✓
- Run-id replay resistance ✓
- Artifact digest binding ✓
- CLI exit codes (`0`/`1`/`2`/`3`) ✓
- Schema/Go struct alignment ✓
