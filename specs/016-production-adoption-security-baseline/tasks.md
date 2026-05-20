# Tasks: Production Adoption And Security Baseline

Status: in_progress

## Phase 0 - Review

- [x] T016-001 Review this spec for scope: production adoption facts must not
  become production trust claims.
- [x] T016-002 Verify the spec decision that `gosec` is advisory for this
  phase is reflected in docs and CI language.
- [x] T016-003 Verify tracked synthetic secret fixtures use only narrow
  reviewed path-and-regex allowlists or scanner-safe rewrites.

## Phase 1 - Pi-Ready Workstreams

- [x] T016-010 WS-016-A: Create `docs/production-adoption-readiness.md`.
  Pi ownership: adoption docs only.
- [x] T016-020 WS-016-B: Create `docs/security-baseline.md` with scanner
  dispositions. Pi ownership: security docs and optional scanner config only.
- [x] T016-030 WS-016-C: Make tracked synthetic fixtures scanner-safe or add a
  narrow reviewed allowlist. Pi ownership: fixture/test files named in the
  security ledger only.
- [x] T016-040 WS-016-D: Add repository security policy. Pi ownership:
  `.github/SECURITY.md`.

## Phase 2 - Integration

- [x] T016-050 Run local verification commands from `plan.md`.
- [x] T016-060 Update docs index/README links if new docs were added.
- [x] T016-070 Record remaining `not_assessed` areas: live CI, customer
  adoption, external production trust, and external security audit.
