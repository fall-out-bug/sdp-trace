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
</anti_rationalizations>

<red_flags>
- A final answer says ready/complete but includes no fresh verification commands.
- A review ledger claims fixes for files absent from the current diff.
- External release, GitHub, OIDC, or customer authority is implied from local artifacts.
- `not_assessed` or `cannot_verify` is hidden in prose instead of named.
- A generic skill tries to close a trust-sensitive block outside this workflow.
</red_flags>
