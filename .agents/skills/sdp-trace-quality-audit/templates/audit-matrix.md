<audit_matrix_template>
Objective:
- User request:
- Success criteria:
- Out of scope:

Artifact checklist:
| Requirement | Artifact or command | State | Evidence | Notes |
| --- | --- | --- | --- | --- |
| Go tests | `go test -count=1 ./...` | pass/fail/cannot_verify/not_assessed | command output |  |
| Vet | `go vet ./...` | pass/fail/cannot_verify/not_assessed | command output |  |
| Schema JSON | `jq empty schema/*.json` | pass/fail/cannot_verify/not_assessed | command output |  |
| Diff whitespace | `git diff --check` | pass/fail/cannot_verify/not_assessed | command output |  |
| Docs drift | `go run ./tools/doccheck` or scoped equivalent | pass/fail/cannot_verify/not_assessed | command output |  |
| CRAP | coverage-backed CRAP command sequence | pass/fail/cannot_verify/not_assessed | command output |  |
| Security/trust review | OmPi reviewer agent plane | accepted/accepted_fixed/advisory/cannot_verify/not_assessed | review disposition |  |
| DX/docs review | OmPi reviewer agent or manual review | accepted/accepted_fixed/advisory/cannot_verify/not_assessed | review disposition |  |

Findings:
| Severity | Area | Finding | Evidence | Disposition |
| --- | --- | --- | --- | --- |
| Critical/Important/Advisory |  |  |  |  |

Closure:
- Proven:
- Not assessed:
- Cannot verify:
- Deferred follow-up:
</audit_matrix_template>
