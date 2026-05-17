# Tasks: Follow-Up Readiness Hardening

**Input**: `followup-hardening-spec.md`
**Prerequisites**: PI spec review, disposition of critical/major findings, explicit approval before implementation code changes.
**Verification vocabulary**: `pass`, `fail`, `cannot_verify`, `not_assessed`.

## Phase 0 - Spec And PI Review

- [x] FH001 Review `followup-hardening-spec.md` with PI architecture/trust doubt.
- [x] FH002 Review `followup-hardening-spec.md` with Kimi wide-context review when Kimi credentials are available.
- [x] FH003 Record all reviewer model/provider/harness, timeout, retry, fallback, and disposition details.
- [x] FH004 Fix or block all critical/major spec findings before implementation.
- [x] FH005 Stop for explicit approval before Go, schema, CI, or product-doc implementation changes.

## Phase 1 - Format And Imports

- [x] FH010 Run the format/import commands from Slice 2 and capture the exact failing files.
- [x] FH011 Apply only formatting/import fixes to the Slice 2 target files.
- [x] FH012 Re-run Slice 2 verification and record `pass` or `fail`.
- [x] FH013 Commit Slice 2 separately.

## Phase 2 - Releaseproof Security Hardening

- [x] FH020 Add test-first negative cases for unsafe `source_commit` refs.
- [x] FH021 Implement immutable commit object validation before source inspection.
- [x] FH022 Align schema constraints with the implementation contract.
- [x] FH023 Fix or narrowly justify remaining gosec rows.
- [x] FH024 Run Slice 3 verification and record `pass`, `fail`, or `cannot_verify`.
- [x] FH025 Commit Slice 3 separately.

## Phase 3 - Duplication Cleanup

- [x] FH030 Confirm current `dupl` output before editing.
- [x] FH031 Extract only high-signal observe CLI shared flow.
- [x] FH032 Extract repeated observe-collect test helpers and hygiene table runner only where behavior is identical.
- [x] FH033 Run Slice 4 verification and record `pass` or `fail`.
- [x] FH034 Commit Slice 4 separately.

## Phase 4 - Lint Cleanup

- [x] FH040 Confirm current gocritic/unparam/prealloc output before editing.
- [x] FH041 Fix production-safe findings in the Slice 5 target files.
- [x] FH042 Reject prealloc suggestions that reduce readability and record why.
- [x] FH043 Run Slice 5 verification and record `pass` or `fail`.
- [x] FH044 Commit Slice 5 separately.

## Phase 5 - Maintainability And Docs Claims

- [x] FH050 Replay absolute file/function MI commands before changing docs.
- [x] FH051 Decide whether this slice closes absolute MI or corrects stale pass language.
- [x] FH052 Update stale docs and CI consistency references.
- [x] FH053 If implementing MI refactors, apply TDD one file at a time and avoid metric-only shuffling.
- [x] FH054 Run Slice 6 verification and record `pass`, `fail`, or `cannot_verify`.
- [x] FH055 Commit Slice 6 separately.

## Phase 6 - Final Evidence And Review

- [x] FH060 Run the full final gate from `followup-hardening-spec.md`.
- [x] FH061 Run PI review planes for security/trust, code/correctness, and docs/evidence drift.
- [x] FH062 Verify every actionable reviewer finding against full files before accepting or rejecting it.
- [x] FH063 Record external GitHub CI as `not_assessed` until live checks are queried for the exact head SHA.
- [x] FH064 Prepare final evidence map with commands, timestamps, commit SHAs, reviewer dispositions, and remaining `not_assessed` areas.
