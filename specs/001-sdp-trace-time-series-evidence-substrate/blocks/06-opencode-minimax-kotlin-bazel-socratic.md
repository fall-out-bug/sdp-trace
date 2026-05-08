# Block 06 Socratic Dialogue: OpenCode + MiniMax + Kotlin+Bazel E2E Proof

Date: 2026-05-01
Block: `06-opencode-minimax-kotlin-bazel-e2e-proof`
Beads mirror: `sdp-trace-drq`

## Consensus Candidate

Block 06 must prove one real product value path: OpenCode + MiniMax over a scoped Kotlin+Bazel target, producing a validated `sdp-trace` evidence package and tested-on report. A packaging-only fixture is insufficient.

## Consensus Result

Consensus is accepted for spec review.

The selected scope is a reference E2E proof path:

1. Use OpenCode and MiniMax as external tested-on tools.
2. Use a scoped Kotlin+Bazel target as the source under assessment.
3. Provide a runner that captures evidence without adding OpenCode, MiniMax, Kotlin, or Bazel as dependencies.
4. Commit only sanitized reports and valid `sdp-trace` artifacts.
5. Keep missing MiniMax credentials, missing OpenCode export, missing Bazel execution, and missing source access as explicit `not_assessed` proof states.
6. Require assessed Bazel command execution before calling Block 06 complete.

Implementation is blocked until pi-review findings on the spec artifacts are recorded, mirrored in Beads, and closed.

## Q1: Is a packager enough to prove product value?

**Critic**: A packager proves formatting, not product behavior. It can turn fake inputs into valid JSON.

**Answer**: No. The first product proof must include an actual OpenCode + MiniMax run. Packaging is necessary only after the real run produces observed evidence.

**Resolution**: Block 06 cannot close with fixture-only artifacts. It needs a committed sanitized run report from a real E2E execution.

## Q2: Does adding a runner violate the repository boundary?

**Critic**: A runner can turn `sdp-trace` into a runtime harness and create hidden dependency on OpenCode.

**Answer**: A narrow reference runner is allowed if it shells out to external tools, records tested-on environment, and keeps all external tools optional. The product dependency remains the portable `sdp-trace` contract.

**Resolution**: Do not add OpenCode, MiniMax, Kotlin, Bazel, Bazelisk, or provider SDKs to `package.json`. The runner fails or records `not_assessed` when the external tool is absent.

## Q3: What proves MiniMax specifically?

**Critic**: Saying "MiniMax" in the command is weak if the provider aliases or routing are ambiguous.

**Answer**: The proof must separate model listing from model access. `opencode models` can prove that a MiniMax id is listed, but access is verified only by a successful run with that model or authenticated provider evidence with an artifact reference.

**Resolution**: `minimax_model_listed`, `minimax_access_verified`, and `opencode_minimax_run_completed` are separate proof states. If OpenCode hides observed provider details, the run report must record that limitation.

## Q4: What proves Kotlin+Bazel specifically?

**Critic**: A repo-level `.bazelrc` or a Kotlin file somewhere in a monorepo does not prove the assessed service is Kotlin+Bazel.

**Answer**: Target-specific evidence is required: a Bazel package/target tied to the selected scope, Kotlin source or Kotlin Bazel rule evidence tied to that target, and workspace context from `MODULE.bazel`, `WORKSPACE`, or `WORKSPACE.bazel`.

**Resolution**: The runner must record target evidence before the model run. When Bazel is available, it should use `bazel query` or record the exact BUILD rule and Kotlin source labels inspected. Missing target evidence keeps `kotlin_bazel_target_identified` `not_observed` or `not_assessed`.

## Q5: Must Bazel build/test execute for the first proof?

**Critic**: If Bazel is not installed, the product can still run OpenCode and package evidence. Is that enough?

**Answer**: It is enough to prove the OpenCode+MiniMax evidence-capture path, but not enough to assess build/test execution. Since the target customer will test Kotlin+Bazel, the first complete product proof must include Bazel or Bazelisk command evidence in an environment that has it.

**Resolution**: `bazel_commands_executed` is its own proof state. It can be `not_assessed` in partial runs, but Block 06 cannot be called complete until at least one committed sanitized proof records assessed Bazel command execution evidence.

## Q5A: Can the model choose the Bazel command to execute?

**Critic**: Executing a model-suggested command in a customer repository is unsafe and can mutate state or run a broad target.

**Answer**: No. The operator must provide or approve the Bazel command before execution. The model may suggest candidate commands, but those candidates are evidence only until approved.

**Resolution**: Full proof requires `--bazel-command` tied to `--bazel-target`. Without it, the package remains partial.

## Q6: Should one observed slice update all compatibility matrices?

**Critic**: A successful OpenCode + MiniMax + Kotlin+Bazel run could be misread as broad OpenCode, MiniMax, Kotlin, or Bazel support.

**Answer**: No. One proof only updates the exact observed row and the exact tested-on environment.

**Resolution**: Matrices may reference the committed report only for the specific OpenCode + MiniMax + Kotlin+Bazel slice. Kimi, GLM, other harnesses, and other stacks remain `not_assessed`.

## Q7: What if the run cannot complete because credentials are missing?

**Critic**: A failed provider setup can still produce evidence. But calling that product proof would be misleading.

**Answer**: Missing credentials or model access can produce a valid incomplete package, but not a completed product proof.

**Resolution**: The runner exits non-zero for incomplete proof and records the missing proof state. The implementation task remains open until a real run package exists.

## Q8: What is the smallest acceptable committed evidence package?

**Critic**: Too many artifact requirements will slow the first proof. Too few will make the proof unauditable.

**Answer**: Minimum committed package:

- run report
- redaction note
- proof states
- evidence events
- provenance records
- observations
- metric stream
- trace snapshot
- assessment input
- validation output or reproducible validation command

**Resolution**: Keep the package narrow, but include enough artifacts for a future external policy consumer consumer to inspect.

## Q8A: How do we prevent raw output from being committed accidentally?

**Critic**: A caller-controlled output path can put raw provider output under a tracked directory before validation runs.

**Answer**: Reject tracked output paths before execution. Raw output must live under `.sdp-trace-runs/` or another path proven ignored by `git check-ignore`.

**Resolution**: Output path safety is a preflight requirement, not a post-run cleanup step.

## Q9: What is the main UX risk?

**Critic**: The user may run the command, get a wall of JSON, and still not know if the product worked.

**Answer**: The runner must print a concise proof-state table and the committed report must start with the same table.

**Resolution**: Human-readable report first, JSON evidence second.

## Q10: What consensus is required before implementation?

**Critic**: The block touches runtime wrappers, product claims, matrices, and validation. Scope can sprawl.

**Answer**: Implementation is limited to one tested slice and one evidence package shape. No broad harness runtime and no policy verdicts.

**Resolution**: Proceed only after spec pi review is complete, the review ledger is updated, every valid spec finding is closed, and the implementation plan preserves the one-slice boundary.
