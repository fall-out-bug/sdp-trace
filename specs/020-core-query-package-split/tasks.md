# Tasks: Core Query Package Split

Status: draft follow-up prepared by Spec 018 closure; no active implementation
tasks are opened in the current closure route.

## Prepared Backlog

These rows define the future implementation backlog. Convert them to active
task checkboxes only when this follow-up spec is explicitly taken into work.

| ID | Future task | Verification |
| --- | --- | --- |
| T020-001 | Review package split plan before implementation. | Maintainer approval recorded before implementation. |
| T020-002 | Confirm no command behavior change is intended. | Review note names behavior-preserving scope. |
| T020-010 | Inventory core query vs forensic query-pack code. | File/function ownership list. |
| T020-020 | Extract query-pack behavior into a forensic extension package. | `go test ./...`; `go vet ./...`. |
| T020-030 | Update command wiring without changing CLI contracts. | Command tests and `go run ./cmd/sdp-trace command-surface`. |
| T020-040 | Update package ownership and extension boundary docs. | `go run ./tools/doccheck`; `go run ./tools/hygienecheck`. |
| T020-050 | Run repository verification. | `git diff --check` plus required Go checks. |
