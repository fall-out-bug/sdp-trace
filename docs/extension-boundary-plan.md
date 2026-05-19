# Extension Boundary Plan

Status: draft

This plan documents how non-core surfaces can move toward extension or
deprecation without deleting behavior or weakening trust-state language.

## Principles

1. Reclassification is documentation and ownership first; it is not command
   removal.
2. `not_assessed` and `cannot_verify` remain explicit. Simplification must not
   turn missing evidence into pass.
3. Any behavior change needs a reviewed implementation spec.
4. Extension commands stay in the binary until a deprecation plan and at least
   one release cycle have completed.

## Lifecycle

| Stage | Meaning | Required Evidence |
| --- | --- | --- |
| classified | Command is listed as core, extension, experimental, or fixture-only. | Stability matrix update. |
| documented separation | First-run docs show the core path first and extension paths as optional. | README/install/adoption docs updated. |
| advisory deprecation | Command is still present but has a documented future-removal notice. | Help text, command-surface note, and reviewed follow-up spec. |
| removal | Command leaves the binary. | Reviewed removal spec, passing tests, updated examples/docs, and one prior release cycle with deprecation notice. |

No command is currently approved for removal.

## Extension Areas

| Area | Commands | Next Decision |
| --- | --- | --- |
| Assessment profiles | `assess` | Stay in binary vs. profile plugin surface. |
| Adapter capture diagnostics | `query --query capture-depth` | Kept in `internal/capturedepth`; decide whether profile owns the command. |
| Protected gates/checkpoints | `gate`, `checkpoint` | Downstream policy consumer vs. in-binary advisory facts. |
| Witness/signing | `witness` | Same binary vs. separate trust/witness tool. |
| Release proof | `release-proof` | Source-bound local trust anchor in binary vs. extension binary. |
| PR evidence | `packet`, `pr-review` | Separate PR tool vs. in-binary workflow. |
| Forensic packaging | `query-pack` | Split from `internal/query` before any package-level core claim. |
| Interaction and harness observation | `interaction`, `observe`, `harness` | Harness-specific adapter package or plugin. |
| Cross-repo posture and telemetry | `export` | Downstream reporting tool vs. in-binary export. |
| Repo observer | `install repo-observer`, `doctor --profile` | Separate installer/doctor utility vs. in-binary helper. |
| Envelope summary | `envelope` | Renderer utility vs. in-binary summary command. |

## Core Commands

Core commands are not deprecation candidates under this plan:

```text
wrap
run
verify
explain
report
query --query missing-evidence
```

Changes to core behavior require updates to `docs/agent-entrypoint.md`, the
schema/contracts affected by the behavior, and the relevant tests.

## Verification

For documentation-only reclassification:

```text
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```
For package or command movement:

```text
go test -count=1 ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```
