# WS-016-B Review: Security Scan Triage

Date: 2026-05-20
Files: `docs/security-baseline.md`, `docs/README.md`

## Quality

| Finding | Severity | Evidence | Disposition |
|---------|----------|----------|-------------|
| Q1: G304 count 61 matches gosec output. All other counts verified against gosec output. | positive | Counts match. | accepted. |
| Q2: G204 deferred advisory lists "Audit every caller" but does not name specific callers to audit. | minor | `internal/trace/source_snapshot.go` callers are hardcoded; `internal/repoobserver/git.go` and `tools/mibaselinepolicy/git.go` need caller enumeration. | **accepted_fixed** — add file/line references for callers that pass external input. |
| Q3: gitleaks finding #4 (`internal/witness/profiles.go:586`) is marked false positive, but the same sentinel string appears in test files as well. Verified: only one production occurrence in profiles.go. | info | Single production occurrence confirmed. | accepted. |

## UX

| Finding | Severity | Evidence | Disposition |
|---------|----------|----------|-------------|
| U1: Classification table uses 8 columns and may wrap on narrow screens. | minor | Markdown table; standard rendering. | advisory. |
| U2: "Local Ignored Clutter Policy" section uses clear bullet list and explicit paths. | positive | Easy to scan. | accepted. |
| U3: gitleaks findings table includes line numbers and disposition rationale. | positive | Traceable to source. | accepted. |

## DX

| Finding | Severity | Evidence | Disposition |
|---------|----------|----------|-------------|
| D1: Verification commands section provides exact copy-pasteable commands. | positive | Reproducible. | accepted. |
| D2: Change log uses ISO date and references spec number. | positive | Consistent. | accepted. |
| D3: docs/README.md link added in correct alphabetical position (Security before Spec Roadmap). | positive | Natural reading order. | accepted. |

## Security

| Finding | Severity | Evidence | Disposition |
|---------|----------|----------|-------------|
| S1: Document explicitly states it is "advisory" and "does not prove absence of vulnerabilities." | positive | No overclaim. | accepted. |
| S2: G101 classified as false positive with rationale. Verified: `ReasonLocalHooksBypassable` is a string constant, `authorityPreviewSafety` is a policy map. | positive | Checked against source. | accepted. |
| S3: G115 classified as false positive with guard condition `r <= 127`. Verified in source. | positive | Checked against source. | accepted. |
| S4: G204 and G703 deferred as advisory rather than falsely claimed fixed. | positive | Honest classification. | accepted. |
| S5: Local ignored clutter policy correctly excludes `.worktrees/`, `.codex-subagents/`, etc. | positive | Matches `.gitignore`. | accepted. |

## Synthesis

- All findings minor, positive, or advisory.
- One fix: add caller references for G204 deferred advisory.
- No blockers.
