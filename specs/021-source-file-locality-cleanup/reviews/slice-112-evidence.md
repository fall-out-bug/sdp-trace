# Slice 112 Evidence

Date: 2026-06-05T00:57:31Z

Scope:
- Consolidated `internal/packet/demo_first_state_helpers.go` into `internal/packet/demo_first_closure_requirements.go`.
- Extended `TestCheckDemoRequiresVerificationOrReviewAssessed` to cover the assessed-state matrix.
- Updated SpecKit plan/tasks and plan review artifact for Slice 112.

Rejected path:
- Movement value/comparison consolidation into `posture_validate_movement_row.go` was attempted as feasibility only.
- It failed file-level MI with `maintainability index 65.0` for `internal/posture/posture_validate_movement_row.go`.
- The movement path was reverted before this implementation scope was finalized.

Source shape:
- `internal/packet/demo_first_state_helpers.go`: removed.
- `internal/packet/demo_first_closure_requirements.go`: now owns `rowAssessed` next to its only consumer.
- `internal/packet/packet_test.go`: focused test covers pass, partial, fail, cannot_verify, not_assessed, and not_in_scope.

Focused verification:

```text
gofmt -w internal/packet/demo_first_closure_requirements.go internal/packet/packet_test.go
test "$(go test ./internal/packet -list 'TestCheckDemoRequiresVerificationOrReviewAssessed$' | grep -c '^TestCheckDemoRequiresVerificationOrReviewAssessed$')" -eq 1
go test ./internal/packet -run 'TestCheckDemoRequiresVerificationOrReviewAssessed$' -count=1 -v
```

Result: pass.

Repository verification:

```text
go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal
go run ./tools/qualitycheck -fail-only -mi-under 70.1 -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools
go test ./...
go vet ./...
golangci-lint run
go run ./tools/doccheck
go run ./tools/hygienecheck
jq empty schema/*.json
git diff --check
go test -count=1 ./... -coverprofile=coverage.out
go tool cover -func=coverage.out > coverage-func.txt
go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt
go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less
rm -f coverage.out coverage-func.txt gocyclo.txt
```

Result: pass.

Drift checks:
- Spec drift: pass; implementation matches Slice 112 packet closure scope.
- Constitution/product drift: pass; no schema, public packet contract, or harness-specific dependency added.
- CRAP: pass; strict `< 5` gate passed.
- MI: pass; no MI baseline changes.
- CleanArch hex: pass; package boundary unchanged.
- CleanCode/SOLID/DRY/YAGNI: pass; single-use helper moved to its only responsibility owner with no new abstraction.

Implementation review:

| Reviewer | Agent ID | Harness | Model/provider | Prompt class | Result |
|---|---|---|---|---|---|
| Beauvoir | `019e9406-f078-7fd2-b8d0-e22ac17a1e3a` | Codex subagent | not_assessed | implementation review | LGTM |
| Peirce | `019e9406-f40c-79f1-904e-54d0f0b73866` | Codex subagent | not_assessed | implementation review | LGTM |
| Halley | `019e9406-f7c2-7f80-80d9-86f7cf7e0c22` | Codex subagent | not_assessed | implementation review | finding |

Initial review findings:
- Halley major: `TestCheckDemoRequiresVerificationOrReviewAssessed` used substring-based `hasError`, so it did not preserve the exact first-packet gate diagnostic required by T021-7790.
- Halley minor: focused evidence recorded `go test -list` output but not the exact-count assertion required by T021-7790.

Fix:
- Updated `TestCheckDemoRequiresVerificationOrReviewAssessed` to use exact diagnostic matching.
- Updated focused evidence to record the exact-count guard command.
- Re-ran focused verification and the full repository verification bundle.

Re-review:

| Reviewer | Agent ID | Harness | Model/provider | Prompt class | Result |
|---|---|---|---|---|---|
| Beauvoir | `019e9406-f078-7fd2-b8d0-e22ac17a1e3a` | Codex subagent | not_assessed | implementation re-review | LGTM |
| Peirce | `019e9406-f40c-79f1-904e-54d0f0b73866` | Codex subagent | not_assessed | implementation re-review | LGTM |
| Halley | `019e9406-f7c2-7f80-80d9-86f7cf7e0c22` | Codex subagent | not_assessed | implementation re-review | LGTM |

Review state: pass.
