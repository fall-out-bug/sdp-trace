# Spec Review 017: OSS Replacement Compatibility And Benchmarks

Reviewer: GLM (adversarial spec review)
Date: 2026-05-20

## Findings

### F-017-01 — Advisory: `check-jsonschema` is a Python tool used as a compatibility probe but not in product path

**Location:** `specs/017-.../spec.md` Evidence section; `docs/oss-replacement-compatibility.md` Prerequisites table

**Evidence:** The spec and compatibility doc both list `check-jsonschema` (Python) as a prerequisite probe tool. The AGENTS.md engineering stack says "No Node.js, npm, JavaScript, TypeScript, or `.mjs` tooling is allowed in the active product path" — Python is not in the banned list, and `check-jsonschema` is used as an external probe, not embedded in product code. FR-017-006 explicitly bans only Node.js/npm/JS/TS/.mjs. This is consistent.

**Assessment:** No violation. The spec correctly keeps Python tools outside the Go active path and marks them as external prerequisites that degrade to `blocked` if missing. Noting for clarity only.

---

### F-017-02 — Important: No copy-pasteable verification commands in the compatibility or benchmark docs

**Location:** `docs/oss-replacement-compatibility.md` (full file, 79 lines); `docs/oss-benchmark-results.md` (full file, 85 lines)

**Evidence:** AGENTS.md trust rules require: "Scanner verification commands in docs must be copy-pasteable and reproducible. If a command changes working directory or depends on local state, use subshell isolation (e.g., `(cd /tmp && ...)`)." The compatibility doc lists prerequisites and results but contains **zero** copy-pasteable shell commands or code blocks. A reader cannot reproduce any probe from the doc alone. The benchmark doc similarly has no actual commands — it presents a markdown table of results but no command lines, working directory instructions, or shell-one-liners.

**Impact:** The first acceptance criterion ("Compatibility probes can be rerun from a clean checkout with documented prerequisites") cannot be verified mechanically from the current docs. The documented prerequisites say which tools are needed but not how to invoke them against which files. This is an unverifiable acceptance criterion.

**Recommendation:** Add a section with at least one copy-pasteable command per probe (e.g., `check-jsonschema --schemafile schema/flight-recorder-run.schema.json examples/agentic-sdlc/local-wrap-positive/run.json`). Mark each command as a probe reproduction step.

---

### F-017-03 — Important: WS-017-B deliverable is ambiguous — "implement a decision" is not testable

**Location:** `specs/017-.../plan.md`, WS-017-B Deliverable section

**Evidence:** The deliverable says "Implement the spec decision that live recorder schema compatibility remains blocked until `wrap` output conforms to `flight-recorder-run.schema.json` or a separate current recorder schema is defined." This has two possible interpretations: (a) write a doc recording that it's blocked, or (b) actually fix the wrap output to conform to the schema. The task T017-020 says "Resolve live wrap output vs flight-recorder schema compatibility" which implies active resolution, while the plan deliverable says "remains blocked." The owned files include `schema/*` and `examples/flight-recorder/*`, suggesting schema changes are expected.

**Impact:** A worker assigned WS-017-B cannot determine from the plan whether to fix the drift, document the blocker, or create a new schema. The acceptance criterion ("The live wrap output/schema drift is either fixed or documented as a blocker") permits either outcome, which means the workstream can complete by writing one sentence and changing nothing. This undermines the spec's stated purpose of testing substitution boundaries.

**Recommendation:** Clarify WS-017-B to either: (1) explicitly scope it to documenting the blocker and adding a failing test, or (2) scope it to fixing the wrap output. Do not leave both options open as equally acceptable deliverables.

---

### F-017-04 — Advisory: Workstream independence has a latent coupling through schema files

**Location:** `specs/017-.../plan.md`, all workstreams

**Evidence:** WS-017-A (probe harness) will validate against `schema/*` files. WS-017-B may modify `schema/*`. If WS-017-B changes `flight-recorder-run.schema.json`, then WS-017-A's probe results may change from `fail` to `pass` (or vice versa). The plan correctly notes Pi handoff isolation for WS-017-B but does not state that WS-017-A must not merge before WS-017-B if schema changes land.

**Assessment:** This is a merge-order risk, not a design flaw. The plan's Pi Handoff Notes section partially addresses this by recommending separate workers. Documenting the expected merge order (B before A, or A first with expected probe result flip) would reduce integration risk.

---

### F-017-05 — Important: Benchmark table is not reproducible and FR-017-004 requires reproducibility

**Location:** `docs/oss-benchmark-results.md`; `specs/017-.../spec.md` FR-017-004

**Evidence:** FR-017-004 requires "Add benchmark output with median, min, max, iterations, command, and environment." The benchmark doc provides median, iterations, and environment but is missing **min**, **max**, and the **actual command lines** used. The doc also admits the benchmark harness does not yet exist: "A Go tool under `tools/ossbench/` could automate the 20-iteration protocol." Task T017-050 tracks this, but the spec already contains a benchmark table in the Evidence section with numbers that cannot be reproduced.

**Impact:** The benchmark numbers in the spec's Evidence section are one-shot local measurements with no reproduction path. This violates AGENTS.md's principle that "Machine proof wins over prose" — the benchmark table is currently prose-level evidence.

**Recommendation:** Either: (1) add min/max and exact commands to the benchmark table, or (2) mark the current benchmark table as `local_one_shot` evidence and add a clear note that T017-050 must land before the numbers are used for any decision.

---

### F-017-06 — Advisory: `.pi/extensions/sdp-trace-boot.ts` is TypeScript in the repo

**Location:** `.pi/extensions/sdp-trace-boot.ts`

**Evidence:** The file is a TypeScript extension for the Pi harness, using `import type { ExtensionAPI } from "@earendil-works/pi-coding-agent"`. AGENTS.md bans "Node.js, npm, JavaScript, TypeScript, or `.mjs` tooling ... in the active product path." The `.pi/` directory is harness configuration, not the active product path (which is `cmd/`, `internal/`, `tools/`). This is consistent with the AGENTS.md Boundary rule which allows "tiny Go validation/rendering tools when needed" in the product path but does not restrict harness extensions.

**Assessment:** No violation. The `.pi/` directory is infrastructure, not product code. Noting for clarity.

---

### F-017-07 — Advisory: Acceptance criterion "probe results distinguish pass/fail/cannot_verify/not_assessed" has no automated check

**Location:** `specs/017-.../spec.md`, Acceptance Criteria

**Evidence:** The second acceptance criterion says probe results must distinguish four states. The compatibility doc table uses `pass`, `fail`, `cannot_verify`, `partial`, `local_pass`, and `blocked` — six distinct states rather than the four specified. The states `partial`, `local_pass`, and `blocked` are not in the accepted set. This is either an over-specification in the AC (too narrow) or an inconsistency in the doc (too broad).

**Recommendation:** Normalize the compatibility doc's status values to the canonical four-state model (`pass`, `fail`, `cannot_verify`, `not_assessed`). Do not expand the AC to include `partial`, `local_pass`, or `blocked`; the Slice 1 validator intentionally accepts only narrow evidence forms.

---

### F-017-08 — Advisory: No existing `tools/osscompat/` or `tools/ossbench/` directories — workstreams start from zero

**Location:** `tools/` directory listing; `specs/017-.../plan.md`

**Evidence:** The current `tools/` directory contains six existing tools (crapcheck, doccheck, hygienecheck, mibaselinepolicy, qualitycheck, schemadoc). Neither `tools/osscompat/` nor `tools/ossbench/` exist yet. WS-017-A and WS-017-E both require new Go tools from scratch. This is not a defect but a scoping note: both workstreams represent non-trivial Go implementation effort. The plan should confirm the workers are aware of the expected tool structure (main.go, focused tests, no TODO/FIXME markers per AGENTS.md).

---

## Summary

| ID | Severity | Summary |
|---|---|---|
| F-017-01 | Advisory | `check-jsonschema` (Python) used as external probe only — no violation |
| F-017-02 | Important | No copy-pasteable verification commands in compatibility/benchmark docs |
| F-017-03 | Important | WS-017-B deliverable ambiguous: "implement a decision" is not testable |
| F-017-04 | Advisory | Latent WS coupling through schema files; merge order undocumented |
| F-017-05 | Important | Benchmark table not reproducible; missing min/max/commands per FR-017-004 |
| F-017-06 | Advisory | `.ts` file in `.pi/` is harness config, not product path — no violation |
| F-017-07 | Advisory | Status vocabulary mismatch: doc uses 6 states, AC specifies 4 |
| F-017-08 | Advisory | WS-017-A and WS-017-E require new Go tools from scratch |

**Verdict: Not LGTM.** Three Important findings require resolution before this spec proceeds to implementation.
