# Spec Review Synthesis: 017 OSS Replacement Compatibility And Benchmarks

Date: 2026-05-20
Reviewer planes: GLM-5.1 (zai), Kimi-k2-thinking (kimi-coding), MiniMax-M2.7 (minimax — failed 404, retried not required per model-policy fallback)
Context: `reviews/017-spec-review-pack.md`

## Findings Summary

| ID | Severity | Finding | Disposition | Resolution |
|---|---|---|---|---|
| F-017-02 | Important | No copy-pasteable verification commands in compatibility/benchmark docs | **accepted** | Add Reproduction Commands section to `docs/oss-replacement-compatibility.md` |
| F-017-03 | Important | WS-017-B deliverable ambiguous: "implement a decision" not testable | **accepted** | Clarify in `plan.md`: scope to documenting blocker + adding failing test, OR fixing wrap output; not both open |
| F-017-05 | Important | Benchmark table not reproducible; missing min/max/commands per FR-017-004 | **accepted** | Add note retracting one-shot table until `tools/ossbench` lands; or add min/max if data available |
| F-017-07 | Important | Status vocabulary mismatch: doc uses 6 states, AC specifies 4 | **accepted** | Map doc states to canonical 4-state schema; update AC or doc for consistency |
| F-017-04 | Advisory | Latent WS coupling through schema files; merge order undocumented | **advisory** | Add merge-order note to `plan.md` Pi Handoff Notes |
| F-017-08 | Advisory | WS-017-A/E require new Go tools from scratch | **advisory** | Noted; no spec change needed |
| Kimi-A4 | Advisory | WS-017-A deliverable "or explains" allows non-executable tooling | **accepted** | Tighten deliverable wording in `plan.md` |
| Kimi-A5 | Advisory | Shell/Go prototypes mentioned in objective but absent from requirements | **advisory** | Noted; objective is descriptive, benchmark table already covers shell prototype |
| Kimi-A6 | Advisory | SLSA verifier missing from benchmark table | **advisory** | Add note to Follow-Ups; SLSA is `not_assessed` for external dependency |

## Disposition Rules Applied
- All Important findings accepted for fix before implementation proceeds.
- Advisory findings accepted if zero-cost or clarity-improving; otherwise recorded.
- No findings rejected as false positive.
- MiniMax reviewer failed with 404; per model-policy, GLM and Kimi provide sufficient independent coverage for spec review. No retry required.

## Verification After Fix
After edits, re-read modified files to ensure:
1. All probe commands are copy-pasteable and subshell-isolated where needed.
2. Status values match canonical schema enum (`pass`, `fail`, `cannot_verify`, `not_assessed`).
3. WS-017-B deliverable has a single clear scope.
4. Benchmark table no longer claims authority it cannot back with min/max/commands.

## Verdict Before Fix
**Not LGTM** — Important findings must be resolved.

## Fixes Applied
The following fixes were applied in separate commits prior to this synthesis.
They are not included in the diff of this file alone; inspect the branch history
for the actual edits.

1. `plan.md`: WS-017-A deliverable tightened; WS-017-B deliverable clarified to doc+test only; merge-order note added.
2. `docs/oss-replacement-compatibility.md`: Status vocabulary aligned to canonical 4 states; `Reproduction Commands` section added with subshell-isolated copy-pasteable commands.
3. `docs/oss-benchmark-results.md`: Commands added; min/max marked `—` with explicit note that one-shot data was not preserved; FR-017-004 gap documented.
4. `spec.md`: Benchmark snapshot annotated as provisional.

## Re-Verification
- `go test -count=1 ./...`: pass
- `go vet ./...`: pass
- `go run ./tools/doccheck`: pass
- `go run ./tools/hygienecheck`: pass
- `git diff --check`: pass

## Verdict After Fix
**LGTM** — All Important findings resolved. Advisory findings accepted or recorded. Spec review complete.
