# Spec Review: 017 OSS Replacement Compatibility And Benchmarks

Reviewer: Kimi (wide-context review)
Date: 2026-05-20
Artifact: `specs/017-oss-replacement-compatibility-and-benchmarks/` (spec.md, plan.md, tasks.md) plus existing docs `docs/oss-replacement-compatibility.md` and `docs/oss-benchmark-results.md`.

---

## Review

### Correct
- **Explicit OSS-tool boundary discipline.** Decisions clearly state no automatic migration, no production trust from local fixtures, and no benchmark health scores. This aligns with AGENTS.md trust rules.
- **Disallowed tooling prohibition is explicit.** FR-017-06 bans Node.js/npm/JS/TS/.mjs from the product path. No such tooling appears in the spec, plan, or existing product paths.
- **Schema drift is documented, not hidden.** The live `wrap` output vs `flight-recorder-run.schema.json` mismatch is recorded as a known blocker, not papered over.
- **Workstreams are independently assignable.** WS-017-A through WS-017-E have disjoint owned-file sets, and the Pi handoff notes correctly flag schema-touch sensitivity and dependency separation for WS-017-C/WS-017-D.
- **Verification commands are standard and copy-pasteable.** `go test -count=1 ./...`, `go vet ./...`, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `git diff --check` are all reproducible.

---

### Important

#### 1. Reproducibility gap: probe commands are not copy-pasteable
**Severity: Important**
**Location:** `spec.md` Evidence section; `docs/oss-replacement-compatibility.md` Compatibility Probes and Prerequisites sections.

The acceptance criterion states: "Compatibility probes can be rerun from a clean checkout with documented prerequisites." However, neither the spec nor the existing docs provide the exact command lines, flags, or expected stdout/stderr outputs used during the 2026-05-20 probe. The docs list tools and verdicts (e.g., "OPA can express a simplified adapter-capture pass/fail rule") but a reader cannot copy-paste a command to reproduce the result.

**Evidence:**
- `docs/oss-replacement-compatibility.md` states `check-jsonschema` validates flight-recorder fixtures, but gives no `--schema` or `--base-uri` flags.
- `docs/oss-replacement-compatibility.md` states `cosign` can sign/verify a local `run.json` blob with transparency log disabled, but gives no `cosign sign-blob` or `cosign verify-blob` command line.
- `docs/oss-replacement-compatibility.md` states `opa` evaluates a simplified rule, but gives no `opa eval` or `opa test` invocation.

**Resolution:** Add a "Reproduction Commands" section to the compatibility doc (or to the planned `tools/osscompat` output) with exact, subshell-isolated commands and expected exit codes/outputs.

---

#### 2. Verifier state vocabulary mismatch between spec, docs, and schema
**Severity: Important**
**Location:** `spec.md` Evidence section (CUE probe); `docs/oss-replacement-compatibility.md` (prerequisites table and Trust table).

The acceptance criteria mandates four states: `pass`, `fail`, `cannot_verify`, `not_assessed`. These match the canonical `verifierState` enum in `schema/flight-recorder-run.schema.json`:
```json
"state": { "enum": ["pass", "fail", "cannot_verify", "not_assessed"] }
```

However, the spec evidence and docs invent additional states not present in the schema:
- `spec.md` uses `partial` for the CUE JSON Schema import probe.
- `docs/oss-replacement-compatibility.md` uses `blocked` for missing tools in the Prerequisites table.
- `docs/oss-replacement-compatibility.md` uses `local_pass` in the Trust And Verification Status table.

This is a trust boundary violation per AGENTS.md ("Machine proof wins over prose"). Prose is introducing state values that checked-in schema does not recognize.

**Resolution:** Map all probe outcomes to the four schema-defined states:
- `partial` → `cannot_verify` (with reason: "CUE module packaging not implemented").
- `blocked` → `cannot_verify` or `not_assessed` (with reason: "required CLI tool not in PATH").
- `local_pass` → `pass` (with reason: "local fixture only; not production trust").

---

#### 3. FR-017-004 benchmark output requirement not met by current evidence
**Severity: Important**
**Location:** `spec.md` Benchmark Snapshot section; `docs/oss-benchmark-results.md` Benchmark Table section.

FR-017-004 requires: "Add benchmark output with median, min, max, iterations, command, and environment."

The existing evidence reports only median:
- `spec.md` Benchmark Snapshot table has columns: Probe, Median ms, Notes.
- `docs/oss-benchmark-results.md` Benchmark Table has columns: Probe, Median (ms), Notes.

Min, max, and the exact command line invoked are absent. The environment is described in prose above the table but is not bound to each probe row. Because the spec already claims a 20-iteration run was performed, the raw min/max data should exist; omitting it from the documented output means the requirement is not yet satisfied.

**Resolution:** Either (a) update the benchmark table to include min, max, iterations, and exact command per row, or (b) if the data is no longer available, retract the benchmark snapshot and rely on the future `tools/ossbench` deliverable to produce compliant output.

---

### Advisory

#### 4. WS-017-A deliverable wording allows non-executable tooling
**Severity: Advisory**
**Location:** `plan.md` WS-017-A Deliverable.

The deliverable reads: "A command that runs or explains compatibility probes without mutating tracked product artifacts." The phrase "or explains" creates an ambiguity: a tool that merely prints instructions satisfies the letter of the spec while failing the spirit of reproducibility in FR-017-001.

**Resolution:** Change to "A command that runs compatibility probes and reports results without mutating tracked product artifacts." If a probe cannot run due to missing dependencies, it should report `cannot_verify`/`not_assessed`, not fall back to explanation mode.

---

#### 5. Shell/Go prototypes mentioned in objective but absent from requirements
**Severity: Advisory**
**Location:** `spec.md` Objective.

The objective lists "minimal shell/Go prototypes" alongside the OSS replacement candidates. No functional requirement, workstream, or task covers creating or benchmarking these prototypes as deliverables. The benchmark snapshot includes a "shell prototype wrap" baseline, but it is unclear whether this prototype is in-scope work or pre-existing context.

**Resolution:** Either add an FR clarifying the prototype baseline, or remove "minimal shell/Go prototypes" from the objective to prevent scope creep.

---

#### 6. SLSA verifier negative path missing from benchmark output
**Severity: Advisory**
**Location:** `docs/oss-benchmark-results.md` Follow-Ups section.

FR-017-002 includes SLSA verifier negative path as a required probe. FR-017-004 requires benchmark output. The existing benchmark table omits SLSA verifier entirely, listing it only as a follow-up. While this is acknowledged, the spec should not close without either (a) including the measurement or (b) explicitly scoping SLSA verifier out of the benchmark requirement.

**Resolution:** Add a note to FR-017-004 or the Non-Goals clarifying that probes returning `not_assessed` due to missing external services may be omitted from benchmark tables.

---

## Summary

The spec demonstrates strong trust discipline: no replacement approval by accident, no production trust from local fixtures, and explicit schema-drift documentation. However, three Important findings must be addressed before implementation proceeds:

1. The reproduction commands must be copy-pasteable.
2. Verifier states must align with the checked-in schema (`pass`, `fail`, `cannot_verify`, `not_assessed` only).
3. Benchmark output must include min, max, iterations, and exact command lines.

The Advisory items are wording/scope clarifications that improve implementer clarity but do not block work.

**Verdict:** Findings present. Not LGTM.
