# Plan: Production Adoption And Security Baseline

Status: draft

## Approach

Keep this slice documentation-first. The goal is not to fix all security
findings immediately; it is to create an auditable triage path and prevent
pilot evidence from being mistaken for production adoption.

## Inputs

- `README.md`
- `docs/adoption-guide.en.md`
- `docs/repository-rollout-playbook.en.md`
- `docs/agent-entrypoint.md`
- `.github/workflows/ci.yml`
- `schema/index.json`
- Security tool outputs from a fresh local run

## Workstreams

### WS-016-A: Adoption Readiness Matrix

Owned files:

- `docs/production-adoption-readiness.md`
- `docs/README.md`
- `README.md` only if the reading path must mention the matrix

Deliverable:

- A matrix that separates controlled-pilot support, local evidence,
  CI-witnessed evidence, source-bound local release, and external production
  trust.

### WS-016-B: Security Scan Triage

Owned files:

- `docs/security-baseline.md`
- optional scanner allowlist/config files if selected

Deliverable:

- Disposition ledger for `govulncheck`, `gosec`, `gitleaks`, and `go vet`.
- Explicit handling of tracked synthetic tokens and local ignored clutter.

### WS-016-C: Scanner-Safe Fixtures

Owned files:

- `examples/self-trace/*`
- `internal/witness/*_test.go`
- historical fixture files only when the spec explicitly allows editing them

Deliverable:

- Fixtures that still test secret redaction behavior but no longer trip
  generic secret scanners, or a narrow reviewed allowlist with rationale.

### WS-016-D: Repository Security Policy

Owned files:

- `.github/SECURITY.md`
- `docs/install.md` only if release/security contact guidance changes

Deliverable:

- Security contact and vulnerability reporting guidance.

## Verification

Minimum local verification after implementation:

```text
go vet ./...
go test -count=1 ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
govulncheck ./...
gitleaks detect on a tracked-file snapshot
```

`gosec ./...` is advisory for this phase. It may become blocking only after a
later spec defines per-call-site dispositions and a stable CI threshold.

## Pi Handoff Notes

Each workstream can be assigned to a separate Pi worker. Workers must not edit
outside their owned files. Security findings are advisory until the integrator
verifies them against full files and fresh commands.
