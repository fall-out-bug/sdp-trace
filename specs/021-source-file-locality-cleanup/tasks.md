# Tasks: Source File Locality Cleanup

Status: draft follow-up prepared by Spec 018 closure; no active implementation
tasks are opened in the current closure route.

## Prepared Backlog

These rows define the future implementation backlog. Convert them to active
task checkboxes only when this follow-up spec is explicitly taken into work.

| ID | Future task | Verification |
| --- | --- | --- |
| T021-001 | Review first cleanup slice and file ownership list. | Maintainer approval recorded before implementation. |
| T021-002 | Confirm the slice is behavior-preserving. | Review note names behavior-preserving scope. |
| T021-010 | Select one package or command family for cleanup. | Slice owner and file list recorded. |
| T021-020 | Move related functions into cohesive behavior-named files. | Focused tests for touched package. |
| T021-030 | Run `gofmt` on changed Go files. | `gofmt -w <changed-go-files>`. |
| T021-040 | Run Go verification. | `go test ./...`; `go vet ./...`. |
| T021-050 | Run docs and hygiene verification. | `go run ./tools/doccheck`; `go run ./tools/hygienecheck`; `git diff --check`. |
