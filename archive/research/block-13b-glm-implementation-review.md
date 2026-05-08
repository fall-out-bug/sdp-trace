## Block 13B SpecKit/DX/Go-first Review

### Finding 1 — Doctor missing required acceptance checks (MAJOR)

**Spec requirement (Doctor Acceptance Criteria):**
> checks wrapper availability, output directory writeability, contract parse, expected-evidence references, and report directory policy

**Implementation:** `buildDoctorReport` only checks contract parse and CI witness identity. Four required checks are absent:

| Required check | Implemented? |
|---|---|
| Wrapper availability | No — `local_wrapper` is a hardcoded `"pass"` string |
| Output directory writeability | No |
| Expected-evidence references | No |
| Report directory policy | No |

Hardcoding `"pass"` for `local_wrapper` without probing the recorder or filesystem is a requirements mismatch. A doctor that reports `pass` without actually checking cannot satisfy "checks wrapper availability."

### Finding 2 — Preview missing required surfaces (MAJOR)

**Spec requirement (Preview Acceptance Criteria):**
> - active boundaries and unsupported or unintegrated boundaries
> - offline implications when the selected profile needs CI or external witness

**Implementation:** `runPreviewCommand` emits `mode`, `command_descriptor`, `contract`, `writes_artifacts`, `safe_retention_modes`, and `warning`. It does **not** emit:

- which observation boundaries are active
- which are `unsupported` or `not_integrated`
- what happens offline when the contract requires `ci_witnessed` or `external_witnessed`

These are explicitly called out as preview acceptance criteria and are absent from the output.

---

### No-critical-or-major items verified

- ObservationState, ObservationBoundary, RetentionMode enums match spec taxonomy exactly.
- Determinism: no wall-clock, random, or hostname fields in doctor or preview output.
- Raw argv/safety: `CommandDescriptor` stores basename + SHA-256 digest only; tests confirm no raw leak in preview/dry-run JSON.
- CI witness prerequisite check correctly enumerates all required GitHub Actions OIDC fields.
- Test coverage is thorough across wrap, verify, explain, query, report, gate, witness, validate-fixtures, and flag parsing.
- No Node/npm path issues (pure Go).

---

## VERDICT: REVISE

Two major requirements mismatches block acceptance:

1. **Doctor** must actually probe wrapper availability, output directory writeability, expected-evidence references, and report directory policy — not hardcode `pass`.
2. **Preview** must include active boundaries, unsupported/unintegrated boundary states, and offline implications for CI/external witness profiles.
