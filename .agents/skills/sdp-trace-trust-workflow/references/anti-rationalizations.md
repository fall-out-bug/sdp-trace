<anti_rationalizations>
Use this list as a stop-check during block intake, implementation, PR prep, and final handoff.

| Rationalization | Required response |
| --- | --- |
| "The checked-in proof JSON says pass." | Treat it as artifact context only until live-verified or externally signed. |
| "The task checkbox is closed, so the block is done." | Prose/task state is not authority; replay evidence or mark open. |
| "Dirty checkout shows what changed, so source-bound proof is enough." | Dirty checkout is local structural evidence only; source-bound proof needs clean immutable source. |
| "Docs can be updated after proof without re-running." | If docs/ledgers are manifest subjects, changing them requires another source-bound cycle. |
| "CI was green before the last commit." | Query final-head checks; otherwise CI is `not_assessed`. |
| "Reviewer approved it." | Reviewer output is advisory until every finding is checked against full files and commands. |
| "No evidence means no problem." | Missing evidence keeps the state `not_assessed` or `cannot_verify`; it is not pass. |
| "Tests passed, so requirements are covered." | Tests are evidence only for behavior they actually exercise. |
| "This is just a small trust change." | Any gate/verdict/provenance change needs trace coverage and explicit verifier state. |
| "We can clean it up later." | Deferred trust closure is not closure; create a tracked follow-up and keep the block open if required evidence is missing. |
| "The review artifact is marked superseded, so it's fine to keep it." | Delete stale artifacts. A header or marker does not prevent stale claims from being cited as evidence. |
| "The command works if you run it in the right directory." | Verification commands must be copy-pasteable and reproducible. Use subshell isolation if a command must change directory. |
| "The scanner count is close enough." | Scanner counts must match the exact output of the documented command. Update docs or commands until they match. |
| "We can leave the email placeholder until someone verifies it." | Unverified contact information in security docs is a finding. Mark `not_assessed` or remove until verified. |
| "The configured scan passes, so the default-config scan doesn't matter." | Both configured and default-config scans must be documented and consistent. Default-config counts are evidence of what external scanners will see. |
</anti_rationalizations>

<red_flags>
- A final answer says ready/complete but includes no fresh verification commands.
- A review ledger claims fixes for files absent from the current diff.
- External release, GitHub, OIDC, or customer authority is implied from local artifacts.
- `not_assessed` or `cannot_verify` is hidden in prose instead of named.
- A generic skill tries to close a trust-sensitive block outside this workflow.
- Stale review artifacts with incorrect counts or outdated findings are left in the repository with only a "superseded" header or marker.
- Scanner verification commands are not copy-pasteable, change the caller's working directory without subshell isolation, or produce different counts depending on where they are run.
- Unverified email addresses or contact channels are listed in security docs without an explicit `not_assessed` marker.
- A review is treated as complete while the reviewer still reports findings of any severity (non-zero finding count).
</red_flags>
