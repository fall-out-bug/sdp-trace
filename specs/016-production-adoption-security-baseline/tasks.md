# Tasks: Production Adoption And Security Baseline

Status: draft

## Phase 0 - Review

- [ ] T016-001 Review this spec for scope: production adoption facts must not
  become production trust claims.
- [ ] T016-002 Decide whether `gosec` is blocking or advisory for this phase.
- [ ] T016-003 Decide the scanner policy for synthetic secret fixtures.

## Phase 1 - Pi-Ready Workstreams

- [ ] T016-010 WS-016-A: Create `docs/production-adoption-readiness.md`.
  Pi ownership: adoption docs only.
- [ ] T016-020 WS-016-B: Create `docs/security-baseline.md` with scanner
  dispositions. Pi ownership: security docs and optional scanner config only.
- [ ] T016-030 WS-016-C: Make tracked synthetic fixtures scanner-safe or add a
  narrow reviewed allowlist. Pi ownership: fixture/test files named in the
  security ledger only.
- [ ] T016-040 WS-016-D: Add repository security policy. Pi ownership:
  `SECURITY.md` or `.github/SECURITY.md`.

## Phase 2 - Integration

- [ ] T016-050 Run local verification commands from `plan.md`.
- [ ] T016-060 Update docs index/README links if new docs were added.
- [ ] T016-070 Record remaining `not_assessed` areas: live CI, customer
  adoption, external production trust, and external security audit.
