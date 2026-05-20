# Command Stability Matrix

Status: draft

This matrix classifies the current `sdp-trace command-surface` into the
smallest stable adoption path and optional extension surfaces. The
classification does not remove commands or change runtime behavior.

## Tiers

| Tier | Meaning |
| --- | --- |
| core | Stable evidence substrate for first adoption. |
| extension | Functional surface for a specific policy, witness, harness, release, forensic, or repo-observer integration. |
| experimental | Discovery, preview, or diagnostic surface whose contract is not the first-adoption path. |
| fixture-only | Repository fixture or CI helper surface, not an adopter workflow. |
| not_assessed | Present but not yet classified. |

## Core Path

| Command | Core Scope | Reason |
| --- | --- | --- |
| `wrap` | Observe one existing command as a trace run. | First step for local evidence capture. |
| `run` | Run a task-referenced command with optional contract evidence. | Task-bound capture path. |
| `verify` | Verify a recorded run directory. | Replays structural evidence locally. |
| `explain` | Render one run for human review. | Makes verifier facts readable without changing state. |
| `report` | Build `.sdp-trace-report/` from one run or run root. | Packages observed and missing evidence. |
| `query --query missing-evidence` | Read or replay the missing-evidence table. | Narrow query path for core evidence gaps. |

The core path records evidence and gaps only. It does not approve merge,
release, readiness, risk, or production trust.

## Extension Path

| Command | Extension Area | Reason |
| --- | --- | --- |
| `query --query capture-depth` | Adapter capture diagnostics | Depends on adapter-capture semantics and now lives in `internal/capturedepth`. |
| `query-pack` | Forensic query packages | Incident/forensic packaging, not first adoption. |
| `assess` | Assessment profiles | Profile-specific policy facts. |
| `gate` | Advisory/protected gate facts | Downstream policy consumer decides block/allow. |
| `checkpoint` | Protected checkpoint proof | Policy/checkpoint integration. |
| `witness` | CI/customer witness evidence | Requires CI or customer authority context. |
| `release-proof` | Source-bound local release proof | Useful trust anchor, but not external production trust. |
| `packet` | Change Evidence Packet bundles | PR packet integration. |
| `pr-review` | Automated PR review evidence pipeline | Review automation, not core capture. |
| `interaction` | Interaction trace import/summary | Harness/interaction integration. |
| `observe` | Harness observation setup/collection | Harness-specific adapter discovery. |
| `harness` | Harness event validation/summary | Harness-specific event processing. |
| `export` | Cross-repo posture / telemetry | Portfolio reporting surface. |
| `envelope` | Delivery trace envelope summary | Rendering/summary extension. |
| `install repo-observer` | Repo observer installation | Writes repo-local observer files. |

## Experimental And Diagnostic

| Command | Reason |
| --- | --- |
| `command-surface` | Agent discovery surface; schema prefix is stable but full shape is not a user adoption contract. |
| `version` | Useful diagnostic output; not yet contract-locked beyond human display. |
| `dry-run` | Preview-only command execution shape. |
| `preview` | Read-only contract preview; unavailable checks stay `not_assessed`. |
| `doctor` | Environment diagnostic. Local `offline_dev` is expected outside CI. |
| `override` | Override request event protocol is not yet a stable command contract. |

## Fixture-Only

| Command | Reason |
| --- | --- |
| `validate-fixtures` | Validates checked fixture directories for repository maintenance and CI. |

## Current Not-Assessed Commands

No command family from the current `command-surface` is unclassified. Some
individual profiles or flags still emit `not_assessed` states when evidence is
outside scope; that is verifier state, not command-tier classification.

## Maintenance Rule

When `sdp-trace command-surface` adds, removes, or changes a command family,
update this matrix and run:

```text
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```
