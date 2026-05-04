# Block 06 Pi Review Ledger

Block: `06-opencode-minimax-kotlin-bazel-e2e-proof`
Beads mirror: `sdp-trace-drq`
Status: spec review findings closed; implementation review findings closed

## Review Rules

- Every valid finding is registered in this ledger.
- Every valid finding is mirrored in Beads as a child issue under `sdp-trace-drq`.
- Minor/P3 findings are treated as required fixes.
- Implementation cannot start until all spec-gate findings are closed.
- Implementation cannot close until all implementation-review findings are closed.

## Spec Review Findings

| ID | Beads mirror | Severity | Finding | Disposition | Evidence |
|---|---|---|---|---|---|
| F001 | `sdp-trace-drq.1` | critical | Block 06 completion can under-prove Kotlin+Bazel if Bazel command execution remains `not_assessed`. | Closed. | Block 06 now requires `bazel_commands_executed` observed for full completion; `rtk npm run validate` passed. |
| F002 | `sdp-trace-drq.2` | major | Bazel command UX was underspecified and could execute model-suggested commands. | Closed. | Runner command now requires `--bazel-target` and `--bazel-command`; model-suggested commands are evidence only; `rtk npm run validate` passed. |
| F003 | `sdp-trace-drq.3` | major | Proof states lacked a machine-readable home and boolean wording risked native verdict semantics. | Closed. | Package now requires `evidence/proof-states.json` with `observed`, `not_observed`, or `not_assessed`; `rtk npm run validate` passed. |
| F004 | `sdp-trace-drq.4` | major | MiniMax listing could be treated as access verification without evidence. | Closed. | `minimax_model_listed`, `minimax_access_verified`, and run completion are separate proof states; `rtk npm run validate` passed. |
| F005 | `sdp-trace-drq.5` | major | OpenCode permission isolation, timeout, and dirty-tree checks were missing. | Closed. | Runner requirements now include no dangerous permission bypass, timeout, no-edit prompt, and git status before/after; `rtk npm run validate` passed. |
| F006 | `sdp-trace-drq.6` | major | Kotlin+Bazel detection was marker-based rather than target-based. | Closed. | Spec now requires Bazel target evidence tied to scope and `bazel query` or exact BUILD rule/source labels; `rtk npm run validate` passed. |
| F007 | `sdp-trace-drq.7` | major | Raw output path was caller-controlled and could be tracked. | Closed. | `--out` must be `.sdp-trace-runs/` or `git check-ignore` verified before execution; `rtk npm run validate` passed. |
| F008 | `sdp-trace-drq.8` | major | Spec gate omitted full validation and stale manifest handling. | Closed. | Implementation plan now requires `npm run validate` or explicit stale proof state before spec acceptance; manifest and DSSE proof refreshed; `rtk npm run validate` passed. |
| F009 | `sdp-trace-drq.9` | major | Tasks privileged Beads over committed review ledger. | Closed. | T079/T080/T085/T086 now name the committed review ledger as primary and Beads as mirror; `rtk npm run validate` passed. |
| F010 | `sdp-trace-drq.10` | minor | Minor wording and portability issues: metric movement overpromise, stale ledger line, and command portability. | Closed. | Spec now says metric sample/stream data, ledger records findings, and tasks use portable commands; `rtk npm run validate` passed. |

## Implementation Review Findings

| ID | Beads mirror | Severity | Finding | Disposition | Evidence |
|---|---|---|---|---|---|
| F011 | `sdp-trace-drq.11` | major | Runner accepted shell-shaped `--bazel-command` values after only checking that the target string appeared. | Closed. | Runner now rejects non-`bazel`/`bazelisk` commands and shell metacharacters before execution; `scripts/test-e2e-runner.sh` covers rejection; `npm run test:e2e-pilot` passed. |
| F012 | `sdp-trace-drq.12` | major | Package validator required proof-state `evidence_refs` but did not verify that they resolve to committed evidence event ids. | Closed. | `scripts/validate-e2e-pilot-package.sh` now resolves every proof-state evidence ref against `evidence/evidence-events.json`; `scripts/test-e2e-pilot-package.sh` covers dangling refs; `npm run test:e2e-pilot` passed. |

## Implementation Evidence Before Review

- Reference runner: `scripts/run-opencode-minimax-kotlin-bazel-proof.sh`
- Package validator: `scripts/validate-e2e-pilot-package.sh`
- Runner tests: `scripts/test-e2e-runner.sh`
- Package validator tests: `scripts/test-e2e-pilot-package.sh`
- Committed proof package: `examples/pilot-runs/opencode-minimax-kotlin-bazel/`
- Tested-on report: `docs/research/opencode-minimax-kotlin-bazel-proof-report.md`
- Latest package validation: `scripts/validate-e2e-pilot-package.sh examples/pilot-runs/opencode-minimax-kotlin-bazel`
