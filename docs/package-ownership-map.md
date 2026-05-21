# Package Ownership Map

Status: machine review complete for Spec 018; maintainer human review not_assessed

This map separates core substrate packages from extension packages. It is a
design boundary for simplification work, not a package-removal plan.

## Core Packages

| Package | Used By | Boundary |
| --- | --- | --- |
| `internal/trace` | `wrap`, `run`, `verify`, `explain`, `report`, query/report readers | Shared trace artifact model and source metadata helpers. |
| `internal/recorder` | `wrap`, `run` | Command observation and run artifact writing. |
| `internal/verifier` | `verify`, `explain`, `report`, core missing-evidence query | Local verifier facts and missing-evidence table. |
| `internal/contract` | `run`, `dry-run`, `preview` | Contract loading and hashing; used by core and preview surfaces. |
| `internal/query` | `query --query missing-evidence`, query-pack internals | Core missing-evidence query plus forensic package support. |

`internal/query` no longer imports `internal/adaptercapture`; the adapter
capture-depth query was split to `internal/capturedepth` in this slice. The
remaining query-pack code is still a mixed package and should be split before
claiming a minimal core library package.

## Extension Packages

| Package | Area |
| --- | --- |
| `internal/capturedepth` | `query --query capture-depth` adapter-capture diagnostics. |
| `internal/adaptercapture` | Adapter-capture assessment profile. |
| `internal/managed` | Managed-harness assessment profile. |
| `internal/forensic` | Forensic-retention assessment profile. |
| `internal/ciartifact` | CI-artifact observation assessment profile. |
| `internal/authority` | Authority-envelope assessment profile. |
| `internal/checkpoint` | Checkpoint creation and verification. |
| `internal/policy` | Protected gate/checkpoint policy helpers. |
| `internal/witness` | CI and customer-PKI witness artifacts. |
| `internal/releaseproof` | Source-bound local release proof. |
| `internal/packet` | Change Evidence Packet bundles. |
| `internal/prreview` | PR review evidence pipeline. |
| `internal/interaction` | Interaction event import and summary. |
| `internal/harnessobs` | Harness observation and event normalization. |
| `internal/posture` | Cross-repo posture aggregation. |
| `internal/telemetry` | Prometheus-style telemetry export. |
| `internal/export` | Export command dispatch. |
| `internal/repoobserver` | Repo observer installation and doctor profile checks. |

## Demo And Fixture Support

| Package | Area |
| --- | --- |
| `internal/demo` | Fixture/demo support for examples and tests. Not a first-adoption runtime dependency. |

## Tools

All `tools/*` packages are repository verification helpers, not product runtime
packages:

| Tool | Purpose |
| --- | --- |
| `tools/doccheck` | Documentation drift checks against command-surface. |
| `tools/hygienecheck` | Repository hygiene checks. |
| `tools/schemadoc` | JSON schema documentation generation. |
| `tools/qualitycheck` | Cyclomatic/cognitive/MI quality checks. |
| `tools/crapcheck` | CRAP score check from coverage and cyclomatic data. |
| `tools/mibaselinepolicy` | Maintainability-index baseline policy checks. |

## Dependency Direction

Target direction:

```text
cmd/sdp-trace
  -> core packages
  -> extension packages

extension packages -> core packages
core packages      -> no extension package imports
```

Current checked state:

- `internal/query` imports `internal/verifier`, but no longer imports
  `internal/adaptercapture`.
- `internal/capturedepth` imports `internal/adaptercapture` and is classified
  as extension.
- Cross-extension imports should stay rare and explicit.

## Source File Locality

Numbered one-function Go source shards are transitional. They helped previous
mechanical organization work, but they are not the target CleanCode,
CleanArchitecture, SOLID, or modern Go shape for this repository.

Target cleanup direction:

- group functions into cohesive files named after behavior, command family, or
  domain concept;
- keep command wiring in `cmd/sdp-trace` thin and move reusable behavior into
  internal packages before grouping;
- avoid package stutter and generated-looking ordinal prefixes in human-owned
  source files;
- split cleanup by package or command family so review remains meaningful;
- preserve tests, command contracts, MI/complexity baselines, and dependency
  direction during each slice.

This map does not approve a repo-wide rename. A later implementation spec must
own each package-level rename/grouping slice and define the verification
commands before files are moved.

## Open Split Work

- Split query-pack code out of `internal/query` if a future core-library claim
  requires package-level purity rather than command-level purity.
- Keep witness, release-proof, PR review, and repo-observer in the same binary
  for this phase. Separate binaries or plugins require a later reviewed spec.
