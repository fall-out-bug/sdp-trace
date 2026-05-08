# Block 06 Implementation Plan: OpenCode + MiniMax + Kotlin+Bazel E2E Proof

Status: accepted for implementation; spec pi-review findings closed
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Beads mirror: `sdp-trace-drq`

## Goal

Implement the first real product proof path: a reference runner that executes or wraps OpenCode + MiniMax against a scoped Kotlin+Bazel target, emits a sanitized `sdp-trace` evidence package, validates it, and records the exact tested-on environment.

## Consensus

Consensus is recorded for a real E2E proof. A packaging-only fixture is not enough. OpenCode, MiniMax, Kotlin, Bazel, and Bazelisk remain external tested-on tools, not repository dependencies.

## File Responsibilities

- `scripts/run-opencode-minimax-kotlin-bazel-proof.sh`: reference runner and evidence capture wrapper.
- `scripts/validate-e2e-pilot-package.sh`: validates a generated E2E proof package.
- retired research artifact: human-readable tested-on report after a real run.
- `examples/pilot-runs/opencode-minimax-kotlin-bazel/`: committed sanitized proof package after a real run.
- retired research artifact: update with the new runner command and exact MiniMax proof boundary.
- retired research artifact: update with E2E proof requirements and Bazel execution proof state.
- retired static harness matrix: update only the exact OpenCode row when committed evidence exists.
- retired static model matrix: update only the exact MiniMax row when committed evidence exists.
- `specs/001-sdp-trace-time-series-evidence-substrate/spec.md`: sync Block 06 functional and success criteria.
- `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md`: add T078-T086 and keep them open until evidence exists.
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/06-opencode-minimax-kotlin-bazel-review-ledger.md`: committed pi-review ledger.

## Task 1: Spec Review Gate

- Run pi review against:
  - `blocks/06-opencode-minimax-kotlin-bazel-e2e-proof.md`
  - `blocks/06-opencode-minimax-kotlin-bazel-socratic.md`
  - `blocks/06-opencode-minimax-kotlin-bazel-implementation-plan.md`
- Record every valid finding in the Block 06 review ledger.
- Mirror every valid finding as a child Beads issue under `sdp-trace-drq`.
- Fix every valid spec finding before implementation.

Verification:

```bash
npm run validate
jq empty schema/*.json
git diff --check
rg -n "packaging-only|not_assessed|OpenCode|MiniMax|Kotlin\\+Bazel|dependency|proof state" specs/001-sdp-trace-time-series-evidence-substrate/blocks/06-opencode-minimax-kotlin-bazel-*.md
```

Expected result: spec artifacts preserve the E2E proof boundary and do not claim product proof without a real run. If manifest/self-attestation proof is stale during spec review, the review ledger must record that state as `not_assessed` with the exact failing command before spec acceptance.

## Task 2: Runner Preflight

- Add `scripts/run-opencode-minimax-kotlin-bazel-proof.sh`.
- Parse `--repo`, `--scope`, `--bazel-target`, `--bazel-command`, `--model`, `--out`, and `--timeout-seconds`.
- Reject `--out` unless it is under `.sdp-trace-runs/` or proven ignored by `git check-ignore`.
- Require `opencode --version`.
- Check `opencode models` for the requested MiniMax model id when possible and record this as `minimax_model_listed`.
- Check target-specific Kotlin+Bazel evidence:
  - `bazel query "$bazel_target"` when Bazel is available
  - exact `BUILD` or `BUILD.bazel` rule tied to the target
  - `MODULE.bazel`, `WORKSPACE`, or `WORKSPACE.bazel` in repository root or owning ancestor
  - `.kt`, `.kts`, `kt_jvm_*`, or Kotlin rule evidence tied to the target
- Print a proof-state table.
- Write raw capture files only under the requested local output directory.

Verification:

```bash
scripts/run-opencode-minimax-kotlin-bazel-proof.sh --help
scripts/run-opencode-minimax-kotlin-bazel-proof.sh --repo /tmp/missing --scope service --bazel-target //service:test --bazel-command "bazel test //service:test" --model minimax-coding-plan/MiniMax-M2.5 --out .sdp-trace-runs/test-missing
```

Expected result: help prints usage; missing repo exits non-zero and does not create a false completed proof.

## Task 3: OpenCode + MiniMax Execution

- Generate the deterministic no-verdict prompt from the Block 06 spec.
- Record `git status --porcelain` before the run when the assessed repository is a git repository.
- Run OpenCode with no dangerous permission bypass and a bounded timeout.
- Run:

```bash
opencode run --model "$model" --format json --dir "$repo" "$prompt"
```

- Capture stdout/stderr separately.
- Record the requested model id, OpenCode version, command, start time, end time, exit code, and output digests.
- Record `minimax_access_verified` as `observed` only after a successful MiniMax run or authenticated provider evidence with artifact reference.
- Record `git status --porcelain` after the run and fail the proof if OpenCode mutated the assessed repository.
- If the run fails after the check runs, emit an incomplete package with `opencode_minimax_run_completed` as `not_observed`; if the check cannot run, emit `not_assessed`; then exit non-zero.

Verification:

```bash
opencode run --help
```

Expected result: implementation uses supported OpenCode `run` flags from local help output.

## Task 4: Bazel Command Evidence

- Detect `bazel` or `bazelisk` in PATH.
- Require operator-provided `--bazel-command` for full proof.
- Refuse to execute a model-suggested command automatically.
- Run the operator-provided command only when it references `--bazel-target` or a narrower target in the same scope.
- If unavailable, record `bazel_commands_executed` as `not_assessed` with reason `bazel_unavailable`.
- Do not treat missing Bazel as model failure, but do keep Block 06 incomplete.

Verification:

```bash
command -v bazel || true
command -v bazelisk || true
```

Expected result: the command records whether Bazel or Bazelisk is available in the current proof environment.

## Task 5: Evidence Package Generation

- Emit:
  - `README.md`
  - `run-report.md`
  - `redaction-note.md`
  - `evidence/proof-states.json`
  - `evidence/evidence-events.json`
  - `evidence/provenance-records.json`
  - `evidence/observations.json`
  - `evidence/metric-stream.json`
  - `evidence/trace-snapshot.json`
  - `handoff/assessment-input.json`
- Reuse existing schemas where possible.
- Mark every unavailable external proof as `not_assessed`.
- Include SHA-256 digests for sanitized artifacts.
- Set package completion to complete only when every required proof state is `observed`, including `bazel_commands_executed`.

Verification:

```bash
find .sdp-trace-runs -maxdepth 4 -type f | sort
```

Expected result: generated package includes the required files and raw output stays in ignored local paths.

## Task 6: Package Validator

- Add `scripts/validate-e2e-pilot-package.sh`.
- Validate JSON parseability.
- Validate evidence/provenance/trace/metric/assessment artifacts against existing schemas where applicable.
- Check required proof states exist in `evidence/proof-states.json`.
- Check full completion requires `bazel_commands_executed` and all other required states as `observed`.
- Check raw output files are not present in committed example directories.
- Wire the validator into `scripts/validate-contracts.sh` conditionally when `examples/pilot-runs/opencode-minimax-kotlin-bazel/` exists.

Verification:

```bash
scripts/validate-e2e-pilot-package.sh examples/pilot-runs/opencode-minimax-kotlin-bazel
```

Expected result: committed proof package validates after a real run exists; missing package exits non-zero during development.

## Task 7: Committed Sanitized Report

- After a real OpenCode + MiniMax + Kotlin+Bazel run with assessed Bazel command execution exists, copy only sanitized artifacts into `examples/pilot-runs/opencode-minimax-kotlin-bazel/`.
- Add retired research artifact with:
  - tested-on environment
  - command summary
  - proof-state table
  - artifact references
  - residual `not_assessed` items
  - statement that no policy verdict is produced by `sdp-trace`

Verification:

```bash
rg -n "OpenCode|MiniMax|Kotlin|Bazel|not_assessed|no policy verdict" retired-research-artifact examples/pilot-runs/opencode-minimax-kotlin-bazel
```

Expected result: report is readable without raw logs or secrets.

## Task 8: Matrix and Docs Sync

- Update the OpenCode/MiniMax rows only when committed evidence exists.
- Keep Kimi, GLM, other harnesses, and non-tested stacks `not_assessed`.
- Update run-card docs with the reference runner command.
- Update SpecKit tasks T078-T086 as each task completes.

Verification:

```bash
node scripts/validate-pilot-matrices.mjs
rg -n "opencode-minimax-kotlin-bazel" retired static harness matrix retired static model matrix retired-research-artifact
```

Expected result: matrices cite the exact proof package and do not broaden the claim.

## Task 9: Implementation Review and Closure

- Run:

```bash
npm run validate
jq empty schema/*.json
git diff --check
scripts/validate-e2e-pilot-package.sh examples/pilot-runs/opencode-minimax-kotlin-bazel
```

- Run pi review on the implementation.
- Register every valid finding in Beads, including minor/P3 findings.
- Fix every valid finding.
- Rerun validation.
- Close `sdp-trace-drq` only after all Block 06 tasks and findings are closed.

Expected result: Block 06 proves one real product value path and the repository states the remaining unproven values explicitly.
