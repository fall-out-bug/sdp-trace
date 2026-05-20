# Review: Slice 5 — Integration and Final Polish

**Reviewer:** Kimi (non-OpenRouter wide-context review)  
**Date:** 2026-05-20  
**Scope:** `specs/017-oss-replacement-compatibility-and-benchmarks/tasks.md`, `spec.md`, `plan.md`; `docs/roadmap.md`, `docs/oss-replacement-compatibility.md`, `docs/oss-benchmark-results.md`; `tools/osscompat/*`, `tools/ossbench/*`; examples and schema drift fixtures.

---

## Correct

- **Verification commands pass.** `go test -count=1 ./tools/osscompat/... ./tools/ossbench/...`, `go vet`, `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, and `git diff --check` all return clean.
- **No disallowed tooling.** No Node.js, npm, JavaScript, TypeScript, or `.mjs` appears in the product path. Python-based `check-jsonschema` is used only as an external probe, consistent with FR-017-006.
- **No TODO/FIXME markers** in new Go files (`tools/osscompat/*`, `tools/ossbench/*`).
- **Honest status transitions.** T017-080 ("Update roadmap and docs index after accepted implementation") is correctly left unchecked because the slice is still in review.
- **Trust claims are scoped with disclaimers.** Docs repeatedly mark results as "local experiments only," "non-authoritative," and `not_assessed` where external evidence is missing. No production Sigstore/Rekor/SLSA trust is claimed from local fixtures.
- **Wrap output/schema drift is preserved as structural evidence.** `examples/flight-recorder/wrap-output-drift/wrap-output.txt` captures the verbatim non-JSON output, and `tools/osscompat/probe_test.go:145-154` (`TestWrapOutputIsNotJSONObject`) asserts the mismatch.
- **`tools/ossbench` satisfies FR-017-004 structurally.** The harness emits `min_ms`, `max_ms`, `median_ms`, `iterations`, and command/environment metadata. `TestRunBenchmark_CustomCommand` and `TestStats` provide coverage.

---

## Fixed

*(No fixes applied by this reviewer; findings are reported for correction by the author.)*

---

## Blocker

*(None identified.)*

---

## Important

### 1. `gofmt` non-compliance in `tools/ossbench`
- **Axis:** DX  
- **Files:** `tools/ossbench/bench.go`, `tools/ossbench/bench_test.go`  
- **Evidence:** `gofmt -l tools/ossbench/*.go` returns both files. Diff shows struct-field misalignment (`benchmarkResult` JSON tag column alignment off by one space; `TestStats` struct field alignment off).  
- **Required by:** AGENTS.md ("`gofmt` for changed Go files" is a default command).  
- **Resolution:** Run `gofmt -w tools/ossbench/bench.go tools/ossbench/bench_test.go`.

### 2. Reproduction commands reference a non-existent directory
- **Axis:** Quality + DX  
- **Files:** `docs/oss-replacement-compatibility.md:59`, `docs/oss-benchmark-results.md:97`  
- **Evidence:** Both docs reference `examples/flight-recorder/local-wrap-positive/run.json`. The actual directory is `examples/flight-recorder/local-positive/` (confirmed via `ls examples/flight-recorder/` and `find examples/flight-recorder -name run.json`).  
- **Impact:** Copy-pasteable commands in the docs will fail with a path-not-found error, undermining the reproducibility acceptance criterion.  
- **Resolution:** Replace `local-wrap-positive` with `local-positive` in both files.

### 3. Stale FR-017-004 unsatisfaction claim in benchmark docs
- **Axis:** Quality + UX  
- **File:** `docs/oss-benchmark-results.md:31-33`  
- **Evidence:** "This table is a **local markdown ledger only** and does not satisfy FR-017-004 until `tools/ossbench` produces reproducible structured output with full statistics."  
- **Problem:** `tools/ossbench` was implemented and task T017-050 is checked. The tool *does* produce reproducible structured output with full statistics (`min_ms`, `max_ms`, `median_ms`, `iterations`, `all_ms`, command, and environment). The sentence is now factually false and misleads readers into believing FR-017-004 is still unmet.  
- **Resolution:** Update the note to state that `tools/ossbench` now exists and satisfies the structural requirement, while the historical table remains provisional because it predates the harness.

### 4. Roadmap Capability Index cross-reference drift
- **Axis:** Quality + DX  
- **File:** `docs/roadmap.md:96`  
- **Evidence:** Capability Index lists Spec 017 as `Draft`. The Active Specs table (`docs/roadmap.md:35`) lists it as `in_progress`. The spec source-of-truth (`specs/017-.../spec.md:3`) also says `in_progress`.  
- **Problem:** Readers navigating by the Capability Index receive a stale status that contradicts the primary tables. T017-080 (roadmap update) is intentionally unchecked, but the existing checked-in file still contains the stale value.  
- **Resolution:** Update Capability Index entry to `in_progress` (or add an inline note referencing T017-080 if the update is intentionally deferred).

---

## Advisory

### A. Minor status inconsistency in prototype docs
- **Axis:** DX  
- **Files:** `docs/oss-policy-prototype.md:3` (`Status: draft`), `docs/oss-supply-chain-prototype.md:3` (`Status: draft`)  
- **Evidence:** Parent spec and sibling docs (`oss-replacement-compatibility.md`, `oss-benchmark-results.md`) are `in_progress`.  
- **Note:** These prototypes are explicitly scoped as local experiments, so `draft` may be intentional. Consider aligning to `in_progress` for consistency, or keep with a parenthetical rationale.

### B. Automated vs manual probe result mapping is slightly divergent
- **Axis:** Quality  
- **File:** `docs/oss-replacement-compatibility.md` (Compatibility Probes table) vs `tools/osscompat/probe.go`  
- **Evidence:** The doc lists `check-jsonschema` live-wrap drift as `fail`, but `tools/osscompat/probe.go:74-78` (`runJSONSchemaWrapDrift`) returns `stateCannotVerify` with reason "live wrap output/schema drift is documented as a blocker." The doc reflects the *manual* reproduction result; the tool reflects the *automated* probe limitation.  
- **Note:** This is not a bug per se, but it may confuse consumers who expect the tool output to mirror the doc table exactly. Consider adding a footnote to the doc table or probe description clarifying that automated probes return `cannot_verify` when they cannot safely invoke the external validator without mutating state.

### C. New docs lack `sdp-trace-claim` tags on status tables
- **Axis:** Security  
- **Files:** `docs/oss-replacement-compatibility.md` (Trust And Verification Status table), `docs/oss-benchmark-results.md` (benchmark claims), `docs/oss-policy-prototype.md` (Probe Result table), `docs/oss-supply-chain-prototype.md` (Probe Results table)  
- **Evidence:** Roadmap Claim-Tag Enforcement Scope says "Required: New authoritative claims in any file created or materially modified after spec 015." None of the new docs contain `sdp-trace-claim` tags. Prose disclaimers ("local experiments only," "non-authoritative") are present and mitigate human misinterpretation, but the status tables read as structured claims.  
- **Note:** If these tables are intended as commentary rather than authoritative claims, consider an explicit non-authoritative preamble (e.g., "The following table is observational commentary, not a machine-authoritative claim set") or add claim tags for the rows that assert `pass`/`fail`/`not_assessed`.

---

## Summary

- **Blocker:** 0
- **Important:** 4
- **Advisory:** 3

All findings are correctable without architecture changes. The slice is structurally sound: harnesses exist, tests pass, trust claims are conservatively scoped, and the wrap/schema drift is preserved as evidence rather than hidden. Once the stale docs, broken path reference, and `gofmt` issues are resolved, the slice will be in good shape for final PR review.
