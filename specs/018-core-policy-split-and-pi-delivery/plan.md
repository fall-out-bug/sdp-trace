# Plan: Core/Policy Split And Pi Delivery

Status: draft

## Workstreams

### WS-018-A: Command Stability Matrix

Owned files:

- `docs/command-stability-matrix.md`
- `docs/agent-entrypoint.md`
- `docs/README.md`

Deliverable:

- Apply the spec-approved command tiers: the named core commands stay core,
  every other current command family is extension, experimental, or
  fixture-only, and no command remains unclassified.

### WS-018-B: Core First Docs

Owned files:

- `README.md`
- `docs/install.md`
- `docs/contributor-quickstart.md`
- `docs/adoption-guide.en.md`
- `docs/adoption-guide.ru.md`

Deliverable:

- First-run docs present the smaller core path before optional pilot surfaces.

### WS-018-C: Package Ownership Map

Owned files:

- `docs/package-ownership-map.md`
- optional generated inventory under `tools/` if implemented in Go

Deliverable:

- Map `cmd`, `internal`, and `tools` packages to core or extension status.

### WS-018-D: Deprecation/Extension Plan

Owned files:

- `docs/extension-boundary-plan.md`
- affected specs only if they need status changes

Deliverable:

- A safe sequence for moving non-core surfaces after a later implementation
  spec; this phase keeps extension surfaces in the current binary.

## Verification

Docs-only changes:

```text
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

Code/tooling changes:

```text
go test -count=1 ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

## Pi Handoff Notes

Do not assign the same docs file to multiple workers. WS-018-A and WS-018-B
both touch docs navigation, so run them sequentially or give one worker an
explicit read-only dependency on the other's output.
