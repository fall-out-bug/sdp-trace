# Spec 020: Core Query Package Split

Status: draft follow-up prepared by Spec 018 closure.

## Objective

Split forensic query-pack code out of `internal/query` so the core
`query --query missing-evidence` path can be reasoned about as a small core
package without mixed forensic ownership.

## Background

Spec 018 approved the core command set as `wrap`, `run`, `verify`, `explain`,
`report`, and `query --query missing-evidence`. It also classified `query-pack`
as a forensic extension and recorded that remaining query-pack code in
`internal/query` should be split before any package-level minimal-core claim.

## Requirements

- FR-020-001: Preserve current command behavior for `query --query
  missing-evidence` and `query-pack`.
- FR-020-002: Move forensic query-pack implementation details out of
  `internal/query` into a clearly named extension package.
- FR-020-003: Keep `internal/query` focused on the core missing-evidence query
  and shared query primitives needed by core commands.
- FR-020-004: Update package ownership docs and command documentation after the
  split.
- FR-020-005: Preserve `not_assessed` and `cannot_verify` states; do not turn
  forensic gaps into pass claims.

## Non-Goals

- No command removal.
- No separate binary or plugin split.
- No behavior change to query output shape without a reviewed follow-up.
- No production trust, release approval, or external attestation claim.

## Acceptance Criteria

- `internal/query` no longer owns forensic query-pack implementation details.
- The forensic extension package is listed in `docs/package-ownership-map.md`.
- `go test ./...`, `go vet ./...`, `go run ./tools/doccheck`,
  `go run ./tools/hygienecheck`, and `git diff --check` pass.

