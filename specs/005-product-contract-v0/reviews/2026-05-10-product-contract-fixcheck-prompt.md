# Product Contract v0 Fix-Check Prompt

Review only the latest fixes after the reshuffled re-review.

Changed risk areas:

- `packet_state` is now lifecycle metadata with allowed values and an anti-score
  rule.
- `PC-DECISION` now has explicit owner states and must not imply approval.
- Local enterprise baseline now marks expected `PC-AUTHORITY` gaps as profile
  characteristics.
- Baseline evidence bundle display can use a compact Markdown table while the
  generated bundle manifest preserves canonical fields.
- Packet generation ownership is now stated as `sdp-trace` tooling, with manual,
  CI, change-host, or release triggers.
- Theater findings tables now list triggered findings only.
- `ci_theater` trigger text now requires a verification-success claim or
  evidence reference, not any CI mention.

Return:

1. Verdict: `FIXES_ACCEPTED`, `REVISE_FIXES`, or `NOT_ASSESSED`.
2. Findings table with columns: `id`, `severity`, `file/section`, `finding`,
   `exact fix`.
3. Whether the prior major findings in your plane are closed.

Do not reopen unrelated product questions unless the new fix created a critical
or major contradiction.
