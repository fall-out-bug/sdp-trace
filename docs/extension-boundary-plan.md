# Extension Boundary Plan

Status: in_review for Spec 018; maintainer human review not_assessed

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

All extension areas stay in the current binary for this phase. The next action
is classification, package-boundary cleanup, or a follow-up implementation spec;
it is not an execution-time decision to remove or split commands.

| Area | Commands | Spec Decision |
| --- | --- | --- |
| Assessment profiles | `assess` | Extension surface in current binary; no plugin split in this phase. |
| Adapter capture diagnostics | `query --query capture-depth` | Extension surface owned by adapter-capture diagnostics; not a core query. |
| Protected gates/checkpoints | `gate`, `checkpoint` | Extension surface that emits advisory facts for downstream policy consumers. |
| Witness/signing | `witness` | Extension surface in current binary; separate trust tool requires a later spec. |
| Release proof | `release-proof` | Extension surface; source-bound local trust anchor, not core adoption. |
| PR evidence | `packet`, `pr-review` | Extension workflow in current binary; separate PR tool requires a later spec. |
| Forensic packaging | `query-pack` | Extension surface; split from `internal/query` before any package-level core claim. |
| Interaction and harness observation | `interaction`, `observe`, `harness` | Extension surface for harness-specific adapters. |
| Cross-repo posture and telemetry | `export` | Extension surface for downstream reporting. |
| Repo observer | `install repo-observer`, `doctor --profile` | Extension helper in current binary; separate installer requires a later spec. |
| Envelope summary | `envelope` | Extension renderer utility in current binary. |

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
