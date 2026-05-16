<final_evidence_map_template>
Changed artifacts:
- Source:
- Schemas:
- Tests/fixtures:
- Docs:
- Reviews/ledgers:

Verification:
- `go test ./...`:
- `go vet ./...`:
- `jq empty schema/*.json`:
- `git diff --check`:
- Additional gate commands:
- Live CI/checks:

Trust states:
- pass:
- fail:
- cannot_verify:
- not_assessed:

Review disposition:
- accepted_fixed:
- rejected_false_positive:
- deferred_not_assessed:
- cannot_verify:
- advisory:

Cleanup:
- Worktree removed:
- Feature branch deleted (local and remote):
- Post-merge main verified (`git pull`, `go test ./...`):

Closure statement:
- What is proven:
- What is not assessed:
- What cannot be verified:
- Required follow-up before closure, if any:
</final_evidence_map_template>
