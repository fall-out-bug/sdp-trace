# WS-016-C Review: Scanner-Safe Fixtures

**Status: superseded by commit `ce7861a`.**
This per-slice review captures the state before `.gitleaks.toml` was updated
with `[extend] useDefault = true` and before the allowlist format was
consolidated into a single `[allowlist]` block. Do not rely on config structure
claims below.

Date: 2026-05-20
Files: `.gitleaks.toml`, `docs/security-baseline.md`

## Quality

| Finding | Severity | Evidence | Disposition |
|---------|----------|----------|-------------|
| Q1: Allowlist regex `crisis-(glm-critic|judge|kimi-critic)-\d{4}-\d{2}-\d{2}` matches all dedupe keys found by gitleaks. | positive | Verified against findings. | accepted. |
| Q2: Allowlist regex `eyJhbGciOiJFZERTQSJ9\.eyJzdWIiOiJibG9jazIyIn0\.signaturesecret` is an exact match for the test sentinel. | positive | Verified against profiles_test.go. | accepted. |
| Q3: Allowlist path `specs/.*/diff\.patch$` covers historical diff patches. | positive | Matches `specs/004-mvp-readiness-hardening/pr-review/ec8db52/packet/inputs/diff.patch`. | accepted. |
| Q4: gitleaks detect with `--config .gitleaks.toml` returns `no leaks found`. | positive | Verified with fresh command output. | accepted. |

## UX

| Finding | Severity | Evidence | Disposition |
|---------|----------|----------|-------------|
| U1: `.gitleaks.toml` file header explains purpose and warns against broadening. | positive | Self-documenting. | accepted. |

## DX

| Finding | Severity | Evidence | Disposition |
|---------|----------|----------|-------------|
| D1: Allowlist is narrow (3 regexes, 1 path) rather than broad file patterns. | positive | Minimal allowlist reduces drift risk. | accepted. |
| D2: `security-baseline.md` updated to reference the allowlist. | positive | Docs stay synchronized. | accepted. |

## Security

| Finding | Severity | Evidence | Disposition |
|---------|----------|----------|-------------|
| S1: Allowlist does NOT allowlist real secret patterns. Only synthetic test markers. | positive | Regexes are exact or specific to known synthetic strings. | accepted. |
| S2: File path allowlist is scoped to `specs/**/diff.patch` only. | positive | Historical review diffs only. | accepted. |
| S3: `gitleaks` with allowlist still scans all other files for real leaks. | positive | `no leaks found` means no other untriaged findings exist. | accepted. |

## Synthesis

- All findings positive.
- No fixes required.
- No blockers.
