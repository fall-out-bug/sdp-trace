# Spec Review Ledger: Invisible Flight Recorder

Date: 2026-05-12

Scope: pre-implementation SpecKit review for `specs/008-invisible-flight-recorder`.

## Review Inputs

- Requirements plane: `reviews/raw/requirements-deepseek.md`
- Evidence/trust plane: `reviews/raw/evidence-deepseek.md`
- Security/non-interference plane: `reviews/raw/security-deepseek.md`
- DX/UX plane: `reviews/raw/dx-deepseek.md`
- Implementation-risk plane: `reviews/raw/implementation-risk-deepseek.md`

## P0/P1 Disposition

| Finding | Severity | Disposition |
| --- | --- | --- |
| Forbidden recorder-duty phrase tests were implicit. | P0 | Added prompt-boundary classification states and FR-013. |
| Contamination classification lacked explicit acceptance criteria. | P0 | Added `clean`, `contaminated`, `digest_only`, `missing`, and malformed fail-closed behavior. |
| CI output set was underspecified. | P0 | Added command contract and FR-014 requiring bundle JSON, Markdown, and result summary. |
| `PC-VERIFICATION` binding to workflow/artifact ids was underspecified. | P0 | Added FR-015 and stale packet override prevention. |
| Missing or contradictory evidence failure tests were underspecified. | P0 | Kept FR-009 and added FR-017 diagnostics requirement. |
| Non-interference and actor authority boundaries were ambiguous. | P0 | Added observer contract boundaries and evidence authority metadata. |
| Manual evidence contamination prevention lacked authority rules. | P0 | Added actor/write-authority/source-state requirements and stale override prevention. |
| GitHub token/secret exposure risk was not explicit. | P0 | Added FR-016. |
| Integration actions needed separate representation tests. | P1 | Kept FR-010 and added authority categories for `integration`. |
| DX command examples and default profile behavior were unclear. | P1 | Added command contract with `--` separator and default live-demo profile example. |
| Manual `build-github` lower authority needed clearer docs. | P1 | Added command contract and documentation requirements. |
| Release binary/docs workflow test was missing. | P1 | Kept T022/T033 and FR-012. |

## Approval State

Spec direction is approved for implementation with the constraints above. Any
implementation that requires prompt mentions, manual packet edits, or manually
authored `.sdp-trace`/`.evidence` files for route proof must be treated as a
blocking defect.
