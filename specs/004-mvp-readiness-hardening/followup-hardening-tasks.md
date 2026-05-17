# Tasks: Follow-Up Readiness Hardening

**Input**: `followup-hardening-spec.md`
**Prerequisites**: PI spec review, disposition of critical/major findings, explicit approval before implementation code changes.
**Verification vocabulary**: `pass`, `fail`, `cannot_verify`, `not_assessed`.

## Phase 0 - Spec And PI Review

- [ ] FH001 Review `followup-hardening-spec.md` with PI architecture/trust doubt.
- [x] FH002 Review `followup-hardening-spec.md` with Kimi wide-context review when Kimi credentials are available.
- [x] FH003 Record all reviewer model/provider/harness, timeout, retry, fallback, and disposition details.
- [x] FH004 Fix or block all critical/major spec findings before implementation.
- [ ] FH005 Stop for explicit approval before Go, schema, CI, or product-doc implementation changes.

## Phase 1 - Format And Imports

- [ ] FH010 Run the format/import commands from Slice 2 and capture the exact failing files.
- [ ] FH011 Apply only formatting/import fixes to the Slice 2 target files.
- [ ] FH012 Re-run Slice 2 verification and record `pass` or `fail`.
- [ ] FH013 Commit Slice 2 separately.

## Phase 2 - Releaseproof Security Hardening

- [ ] FH020 Add test-first negative cases for unsafe `source_commit` refs.
- [ ] FH021 Implement immutable commit object validation before source inspection.
- [ ] FH022 Align schema constraints with the implementation contract.
- [ ] FH023 Fix or narrowly justify remaining gosec rows.
- [ ] FH024 Run Slice 3 verification and record `pass`, `fail`, or `cannot_verify`.
- [ ] FH025 Commit Slice 3 separately.

## Phase 3 - Duplication Cleanup

- [ ] FH030 Confirm current `dupl` output before editing.
- [ ] FH031 Extract only high-signal observe CLI shared flow.
- [ ] FH032 Extract repeated observe-collect test helpers and hygiene table runner only where behavior is identical.
- [ ] FH033 Run Slice 4 verification and record `pass` or `fail`.
- [ ] FH034 Commit Slice 4 separately.

## Phase 4 - Lint Cleanup

- [ ] FH040 Confirm current gocritic/unparam/prealloc output before editing.
- [ ] FH041 Fix production-safe findings in the Slice 5 target files.
- [ ] FH042 Reject prealloc suggestions that reduce readability and record why.
- [ ] FH043 Run Slice 5 verification and record `pass` or `fail`.
- [ ] FH044 Commit Slice 5 separately.

## Phase 5 - Maintainability And Docs Claims

- [ ] FH050 Replay absolute file/function MI commands before changing docs.
- [ ] FH051 Decide whether this slice closes absolute MI or corrects stale pass language.
- [ ] FH052 Update stale docs and CI consistency references.
- [ ] FH053 If implementing MI refactors, apply TDD one file at a time and avoid metric-only shuffling.
- [ ] FH054 Run Slice 6 verification and record `pass`, `fail`, or `cannot_verify`.
- [ ] FH055 Commit Slice 6 separately.

## Phase 6 - Final Evidence And Review

- [ ] FH060 Run the full final gate from `followup-hardening-spec.md`.
- [ ] FH061 Run PI review planes for security/trust, code/correctness, and docs/evidence drift.
- [ ] FH062 Verify every actionable reviewer finding against full files before accepting or rejecting it.
- [ ] FH063 Record external GitHub CI as `not_assessed` until live checks are queried for the exact head SHA.
- [ ] FH064 Prepare final evidence map with commands, timestamps, commit SHAs, reviewer dispositions, and remaining `not_assessed` areas.
