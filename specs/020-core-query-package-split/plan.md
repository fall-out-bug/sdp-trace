# Plan: Core Query Package Split

Status: draft follow-up prepared by Spec 018 closure.

## Workstreams

### WS-020-A: Package Boundary Inventory

Owned files:

- `internal/query/`
- `docs/package-ownership-map.md`

Deliverable:

- Identify which files and functions serve core missing-evidence queries and
  which serve forensic query-pack behavior.

### WS-020-B: Forensic Query Package Extraction

Owned files:

- `internal/query/`
- new forensic query extension package under `internal/`
- affected `cmd/sdp-trace` query-pack wiring

Deliverable:

- Move query-pack-specific behavior into an extension package while preserving
  command behavior and tests.

### WS-020-C: Documentation And Verification

Owned files:

- `docs/package-ownership-map.md`
- `docs/command-stability-matrix.md`
- `docs/extension-boundary-plan.md`

Deliverable:

- Update ownership and extension documentation to match the package split.

## Verification

```text
go test ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

