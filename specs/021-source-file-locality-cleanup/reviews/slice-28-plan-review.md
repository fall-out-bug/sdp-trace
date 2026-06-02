# Slice 28 Plan Review: Command Model And Shell Fields

Status: passed

## Proposed Slice

- Scope: `internal/harnessobs/harnessobs_181` through
  `internal/harnessobs/harnessobs_197`.
- Intended grouping:
  - `command_model_extraction.go` for `extractCommandModel`
  - `command_model_shell.go` for `shellCommandString` and
    `shellCommandShape`
  - command model argv helpers for `extractCommandModelArgs`,
    `commandModelArg`, `nextCommandModelArg`, and `prefixedCommandModelArg`,
    split narrower if MI requires
  - `shell_fields.go` for `shellFields`, including the existing controlled
    shell parser comment currently parked under `prefixedCommandModelArg`
  - pre-MI scanner-stage candidate `shell_field_scanner_stages.go` for
    `scan`, `consumeEscaped`, `startsEscape`, `consumeQuoted`, `shellQuote`,
    `consumeUnquoted`, `shellFieldSeparator`, `finish`, and `flush`; split
    narrower only if MI requires
  - if split is required, keep `shellQuote` with quote/escape stages and
    `shellFieldSeparator` with unquoted field consumption
- Explicitly excluded: command model safety, source commit, session setup, raw
  event safety, and validation helpers (`harnessobs_198` onward).

## Review Lanes

- lane 1 requirements/scope: `LGTM` equivalent; opencode-go/deepseek-v4-pro
  via OpenCode, 2026-06-02, prompt class `plan-review/scope`.
- lane 2 trust/evidence: `LGTM`; kimi-for-coding/k2p6 via OpenCode,
  2026-06-02, prompt class `plan-review/trust-evidence`.
- lane 3 maintainability: `LGTM`; opencode-go/minimax-m3 via OpenCode,
  2026-06-02, prompt class `plan-review/maintainability-dx`.

## Findings

- lane 3 maintainability found the initial scanner-stage plan too loose:
  move the `shellFields` parser comment next to `shellFields`, make
  `shell_field_scanner_stages.go` the explicit pre-MI candidate, and bind
  `shellQuote`/`shellFieldSeparator` to their primary consumers if a split is
  required. The plan above has been updated and focused re-review returned
  `LGTM`.

## Non-Evidence Attempts

- An opencode-go/glm-5.1 scope lane on 2026-06-02 exited without a visible
  verdict. It is `cannot_verify` and is not counted as reviewer evidence.
- A replacement opencode-go/glm-5.1 scope lane on 2026-06-02 was stopped after
  prolonged review without a verdict. It is `cannot_verify` and is not counted
  as reviewer evidence.
