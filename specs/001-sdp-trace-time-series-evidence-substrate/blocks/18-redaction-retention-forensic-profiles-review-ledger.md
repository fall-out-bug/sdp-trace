# Block 18 Review Ledger

Status: spec review pass recorded; implementation review fixes applied; PR-level
review remains open.

## Spec Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| B18-S001 | critical | privacy/security | Digest binding for sealed and external raw references did not require a strong algorithm. | Accepted and fixed. Block 18 now requires SHA-256 or stronger for forensic profiles and treats weak, legacy, unknown, or truncated digests as `cannot_verify`. | `18-redaction-retention-forensic-profiles.md` sections `Event And Artifact Delta`, `Test And Fixture Expectations`; `spec.md` FR-092; `tasks.md` T154/T156 |
| B18-S002 | critical | privacy/security | Redaction authority was a loose field and could be self-claimed. | Accepted and fixed. Authority must reference provenance or accountability identity; self-asserted authority is `cannot_verify` for forensic retention. | `18-redaction-retention-forensic-profiles.md` section `Redaction Policy Contract`; `spec.md` FR-089; `tasks.md` T150/T156 |
| B18-S003 | critical | privacy/security | "Cap to lower trust" was operationally undefined. | Accepted and fixed. Forensic retention now fails insufficient critical evidence while emitting explanatory `capped_to_retention_mode`; cap is not a pass state or upgrade path. | `18-redaction-retention-forensic-profiles.md` section `Verifier Semantics`; `spec.md` FR-091; `tasks.md` T153 |
| B18-S004 | critical | requirements/readiness | Retention modes and the forensic assessment profile were conflated. | Accepted and fixed. The spec now separates embedded retention modes from `forensic_retention` / `--profile forensic-retention` assessment. | `18-redaction-retention-forensic-profiles.md` section `Profile Model`; `spec.md` Key Entities; `tasks.md` T149/T151/T153 |
| B18-S005 | critical | requirements/readiness | The review ledger requirement was missing, making T158 underspecified. | Accepted and fixed. This ledger is now a required artifact and T158 names it explicitly. | `18-redaction-retention-forensic-profiles.md` section `Review Plan`; `tasks.md` T158 |
| B18-S006 | major | tracing/evidence | Raw reference fields were described only in prose. | Accepted and fixed. The spec now includes a schema-shaped raw reference contract with closed access, custody, lifecycle, and unavailable-reason values. | `18-redaction-retention-forensic-profiles.md` section `Event And Artifact Delta`; `tasks.md` T151/T154 |
| B18-S007 | major | tracing/evidence | Critical event family classification was not defined. | Accepted and fixed. The spec now lists default critical event families and requires policy-visible downgrades. | `18-redaction-retention-forensic-profiles.md` section `Event And Artifact Delta`; `spec.md` FR-089/FR-091 |
| B18-S008 | major | privacy/security | Raw reference access state could become stale or mutable after assessment. | Accepted and fixed. Access state is explicitly assessment-time evidence; revocation, compromise, or withdrawal requires superseding events. | `18-redaction-retention-forensic-profiles.md` section `Event And Artifact Delta`; `tasks.md` T154 |
| B18-S009 | major | privacy/security | Built-in safe default policy had no version stability contract. | Accepted and fixed. The spec now requires a versioned built-in policy id and digest with migration behavior. | `18-redaction-retention-forensic-profiles.md` section `Redaction Policy Contract`; `spec.md` FR-088 |
| B18-S010 | major | privacy/security | Preview semantics could diverge from runtime redaction behavior. | Accepted and fixed. Preview must use the same redaction engine and policy resolver and expose rule ids/classes/actions, not matched values. | `18-redaction-retention-forensic-profiles.md` section `CLI And UX` |
| B18-S011 | major | tracing/evidence | Assessment-result schema versioning and read compatibility were not specified. | Accepted and fixed. Block 18 must extend or version `assessment-result` while preserving `assess explain` compatibility. | `18-redaction-retention-forensic-profiles.md` section `Implementation Slices`; `tasks.md` T155 |
| B18-S012 | major | privacy/security | Withholding lacked audit fields and could suppress evidence without accountability. | Accepted and fixed. Withholding must record authority, requestor identity when different, reason code, and justification. | `18-redaction-retention-forensic-profiles.md` section `Redaction Policy Contract`; `spec.md` FR-089 |
| B18-S013 | critical | tracing/evidence | Block 18 introduced retention-mode names that diverged from FR-054. | Accepted and fixed. Block 18 now uses FR-054 retention modes exactly and treats safe/default/sanitized/encrypted/external choices as recording policy profiles or assessment behavior. | `18-redaction-retention-forensic-profiles.md` sections `Profile Model`, `Redaction Policy Contract`; `spec.md` FR-090 Key Entities; `tasks.md` T149-T151 |
| B18-S014 | major | tracing/evidence | Verifier semantics did not enumerate condition outcomes. | Accepted and fixed. The condition table now maps each condition to `pass`, `fail`, `cannot_verify`, or `not_assessed` behavior. | `18-redaction-retention-forensic-profiles.md` section `Verifier Semantics` |
| B18-S015 | major | tracing/evidence | Recorder-emitted fields and verifier-computed fields were not separated. | Accepted and fixed. The spec now names recorder-emitted redaction/retention fields and keeps computed facts in assessment condition rows. | `18-redaction-retention-forensic-profiles.md` section `Event And Artifact Delta` |

## Review Evidence State

- MiniMax-M2.7 privacy/security review: assessed.
- Kimi K2P6 requirements/readiness review: assessed.
- ZAI/GLM-5.1 tracing/evidence review: first run returned no usable review; replacement run completed and was assessed.
- Local tracing/evidence synthesis: used only to apply and cross-check the concrete schema/task fixes above; it does not replace a later implementation or PR-level pi review.

## Implementation Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| B18-I001 | major | tracing/evidence | `docs/flight-recorder.md` still used old recorder profile names, making Block 18 docs drift from the active schemas. | Accepted and fixed. Docs now use `local_development_recorder`, `witnessed_run_recorder`, and `externally_witnessed_run`, and state that `forensic_retention` is an assessment profile rather than a recorder profile or retention mode. | `docs/flight-recorder.md` section `Recorder Profiles` |
| B18-I002 | major | schema/correctness | `assessment-result` reused managed-only condition states for forensic condition rows. | Accepted and fixed. Managed and forensic condition schemas now have separate state enums; forensic rows allow only `pass`, `fail`, `cannot_verify`, or `not_assessed`. | `schema/assessment-result.schema.json` `$defs.managedConditionState`, `$defs.forensicConditionState` |
| B18-I003 | major | UX/correctness | Forensic preview claimed runtime redaction-engine equivalence that the implementation did not execute. | Accepted and fixed. Preview now declares `not_executed_in_preview`, remains read-only, and the spec labels runtime equivalence as a future requirement only when actually executed. | `cmd/sdp-trace/main.go` forensic preview output; `18-redaction-retention-forensic-profiles.md` section `CLI And UX` |
| B18-I004 | major | tracing/evidence | The verifier missed rule references and withholding audit evidence, so policy application could be under-specified. | Accepted and fixed. Events now carry rule refs into evaluation; unknown or action-mismatched rules fail, missing apply-rule refs are `cannot_verify`, and withholding requires verified authority plus reason and justification. | `internal/forensic/forensic.go`; `internal/forensic/forensic_test.go`; `schema/redaction-policy.schema.json` |
| B18-I005 | major | requirements-vs-implementation | Self-claimed redaction authority produced `fail`, but the spec says forensic authority that cannot be verified is `cannot_verify`. | Accepted and fixed. Self-claimed redaction authority now yields `cannot_verify`; fixture naming and README were updated to match. | `internal/forensic/forensic.go`; `examples/block18-forensic-retention/authority-self-claim-cannot-verify.assessment-result.json` |
| B18-I006 | major | tracing/evidence/schema | Forensic condition arrays required at least eight rows but did not force exactly one row for each condition id. | Accepted and fixed. Assessment-result and run schemas now require all eight forensic condition ids and set `minItems: 8` plus `maxItems: 8`, preventing duplicate-id ambiguity. | `schema/assessment-result.schema.json`; `schema/flight-recorder-run.schema.json` |

## Current Review Evidence State

- MiniMax-M2.7 code/correctness re-review: one claimed major about write-error exit-code inconsistency was rejected as a false positive after checking that managed and forensic assess both return `1` on write failure; no accepted critical or major code finding remains from that pass.
- ZAI/GLM-5.1 tracing/evidence/schema re-review: two claimed majors were rejected as false positives after checking the full files (`evidenceRetention` exists and run-level `forensic_conditions` is wired); one accepted major about duplicate forensic condition ids was fixed and re-reviewed.
- Kimi K2P6 requirements-vs-implementation re-review: no critical or major findings in the provided diff.
- Strict implementation review fixes above are applied and locally verified.

## PR-Level Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| B18-PR001 | major | tracing/evidence/schema | PR-level review found that `redaction-policy.schema.json` required policy fields that the Go `Policy` struct neither represented nor consumed, and that an empty policy could cascade into `redaction_rule_unknown` fail instead of remaining `cannot_verify`. | Accepted and fixed. The Go policy contract now models redaction actions, forbidden persistence classes, authority, profile mappings, and unresolved-redaction impact; incomplete policy contracts are `cannot_verify`; rule coverage is `cannot_verify` rather than `fail` when the selected policy is missing. | `internal/forensic/forensic.go`; `internal/forensic/forensic_test.go`; `examples/block18-forensic-retention/missing-policy-cannot-verify.assessment-result.json` |

## PR Review Evidence State

- MiniMax-M2.7 PR-level code/correctness review: no critical or major correctness regressions; minor observation about `argv` becoming optional for retained command metadata.
- ZAI/GLM-5.1 PR-level tracing/evidence/schema review: one accepted major, B18-PR001, fixed locally after review.
- Kimi K2P6 PR-level requirements-vs-implementation review: no critical or major findings.
- GitHub CI is `not_assessed`: GitHub reported no checks on PR #8.
