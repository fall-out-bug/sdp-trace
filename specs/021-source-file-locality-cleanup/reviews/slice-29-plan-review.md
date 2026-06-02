# Slice 29 Plan Review: Command Model Safety And Source Digest

Status: passed

## Proposed Slice

- Scope: `internal/harnessobs/harnessobs_198` through
  `internal/harnessobs/harnessobs_203`.
- Intended grouping:
  - `command_model_safety.go` for `safeCommandModel`
  - `command_model_unsafe.go` for `unsafeCommandModelIdentity`,
    `unsafeCommandModelChars`, and `unsafeCommandModelPath`
  - `source_digest_file.go` for `digestFile`, reusing existing
    package-local `sha256Hex` if focused regression evidence confirms
    unchanged SHA-256 behavior
  - `source_commit.go` for `sourceCommit`
  - `source_commit_hash.go` for package-local commit-hash validation if
    focused regression evidence confirms unchanged hash-shape behavior
- Explicitly excluded: unsafe raw-event traversal (`harnessobs_204` onward),
  session setup, collection, validation, and isolation helpers.

## Review Lanes

- lane 1 requirements/scope: `LGTM`; opencode-go/deepseek-v4-pro via
  OpenCode, 2026-06-02, prompt class `plan-review/scope`.
- lane 2 trust/evidence: `LGTM`; opencode-go/deepseek-v4-flash via OpenCode,
  2026-06-02, prompt class `plan-review/trust-evidence`, after fixing the
  task-checkbox overclaim found by opencode-go/mimo-v2.5-pro.
- lane 3 maintainability/DX: `LGTM`; opencode-go/qwen3.7-max via OpenCode,
  2026-06-02, prompt class `plan-review/maintainability-dx`, after recording
  the local cleanup constraints.

## Findings

- lane 3 maintainability found two minor local cleanup opportunities inside the
  selected files:
  - `digestFile` duplicates existing package-local `sha256Hex`.
  - `sourceCommit` compiles the source commit hash regexp on every invocation.
- Resolution: fold both into the implementation as behavior-preserving local
  cleanup inside the selected files, then rerun focused regression and
  maintainability re-review. Focused plan re-review returned `LGTM`.
- lane 2 trust/evidence found that T021-1920 and T021-1921 were checked before
  plan review closure. Resolution: return those task checkboxes to pending
  until all plan-review lanes close. Focused trust re-review returned `LGTM`.

## Non-Evidence Attempts

- An opencode-go/mimo-v2.5-pro trust/evidence re-verdict found the
  task-checkbox overclaim above. It is recorded as a finding and not counted as
  reviewer closure.
