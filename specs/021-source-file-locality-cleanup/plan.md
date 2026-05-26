# Plan: Source File Locality Cleanup

Status: draft follow-up prepared by Spec 018 closure.

## Workstreams

### WS-021-A: Slice Selection

Owned files:

- one package or command family per slice

Deliverable:

- Pick the highest-value package or command family and define a bounded file
  ownership list before moving code.

### WS-021-B: Behavior-Named Grouping

Owned files:

- selected package or command family

Deliverable:

- Move related functions from numbered shards into cohesive behavior-named
  files while preserving tests and public behavior.

### WS-021-C: Verification And Docs

Owned files:

- selected package docs if ownership or dependency direction changes
- `docs/package-ownership-map.md` when needed

Deliverable:

- Verify behavior and update ownership docs only when the cleanup changes
  package-level boundaries.

## Verification

```text
gofmt -w <changed-go-files>
go test ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

