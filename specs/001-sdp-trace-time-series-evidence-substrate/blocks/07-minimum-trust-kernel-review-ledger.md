# Block 07 Review Ledger

Status: implementation review in progress; Slice 0 implementation reviewed
Parent: `07-minimum-trust-kernel.md`

This ledger records review findings for Block 07 before implementation. It is repository-visible evidence, not proof authority. The live verifier remains the claim authority after implementation.

## Review Sources

| Source | Scope |
| --- | --- |
| subagent trust-kernel review | external trust gate and MVP scope |
| subagent security review | forgery attack tree and critical controls |
| subagent DX/testing review | verifier UX and package validation split |
| subagent consistency audit | current false closure claims |
| pi / MiniMax-M2.7 | adversarial spec review |
| pi / kimi-k2-thinking | adversarial spec review |
| pi / glm-5.1 | adversarial spec review |
| pi / gpt-5.4-mini | adversarial spec review |
| pi / gpt-5.3-codex-spark | closure check after post-review edits |
| pi / MiniMax-M2.7 | Slice 0 security and forgery review |
| pi / zai/glm-5.1 | Slice 0 SpecKit, bash correctness, and closure review |
| pi / zai/glm-5-turbo | Slice 0 bash/code-quality review |
| pi / openrouter/qwen/qwen3-coder | Slice 0 DX/replayability review |
| pi / kimi-coding/kimi-for-coding | attempted Slice 0 code review; timed out and was replaced |

## Findings

| ID | Severity | Source | Finding | Disposition | Spec or plan change | Closure |
| --- | --- | --- | --- | --- | --- | --- |
| B07-F001 | critical | consistency audit | Current Block 04, T070, and Block 06 closure claims contradict current validation/proof state. | accepted | Plan Slice 2 requires verifier-derived freeze/repair before closure. | closed in spec/plan |
| B07-F002 | critical | trust-kernel review | `production_release_verified.value: true` can become a hidden trust claim outside `trusted_contract_release`. | accepted | Spec requires production and trusted release booleans derive from the same passed state set; plan Slice 3 adds negative fixtures. | closed in spec/plan |
| B07-F003 | critical | security review | Checked-in verification JSON is untrusted unless regenerated or verified. | accepted | Claim authority now belongs to live verifier execution; persisted summaries are audit/cache artifacts. | closed in spec/plan |
| B07-F004 | critical | pi reviews | Verifier bootstrap lacks trust anchor. | accepted | Spec adds Trust Anchor Assumption and states source distribution trust is out-of-band unless external production profile is selected. | closed in spec/plan |
| B07-F005 | critical | pi reviews | Review ledger requirement was circular because ledger was created after implementation. | accepted | Plan moves review ledger to pre-implementation artifact. | closed in spec/plan |
| B07-F006 | critical | pi reviews | External trust profile contradicted by non-strict mode. | accepted | Spec says selected external-trust profile fails when evidence is missing; `not_assessed` only applies under other profiles. | closed in spec/plan |
| B07-F007 | critical | pi reviews | `proof-summary` and `contract-release-verification` risk dual authority. | accepted | Spec subordinates release verification artifacts to proof summary/verifier output. | closed in spec/plan |
| B07-F008 | major | DX/testing review | Baseline validation is coupled to optional incomplete E2E pilot completion. | accepted | Spec separates baseline from slice verification; plan Slice 4 moves slice verification to separate command surface. | closed in spec/plan |
| B07-F009 | major | DX/testing review | E2E package validator conflates package validity with completed proof. | accepted | Plan Slice 4 defines `--mode package` and `--mode complete`. | closed in spec/plan |
| B07-F010 | major | pi reviews | Required vs optional states and `assessed` semantics were undefined. | accepted | Spec adds State Vocabulary and required/optional state tables. | closed in spec/plan |
| B07-F011 | major | pi reviews | Source subject semantics were git-biased and underspecified for non-git environments. | accepted | Spec defines `git_commit_v1`, `directory_snapshot_v1`, and `signed_source_archive_v1`. | closed in spec/plan |
| B07-F012 | major | pi reviews | Prose claim consistency via NLP is infeasible. | accepted | Spec introduces machine-readable claim tags; plan Slice 1 validates tags instead of prose. | closed in spec/plan |
| B07-F013 | major | pi reviews | Negative fixtures are weak unless executed by the verifier. | accepted | Plan Slice 0 requires negative fixture execution in baseline verifier. | closed in spec/plan |
| B07-F014 | major | pi reviews | Cross-reference integrity was late and under-specified. | accepted | Plan adds dedicated Slice 5 with concrete reference checks and negative fixtures. | closed in spec/plan |
| B07-F015 | major | pi reviews | Rollback protection must bind to spec source, not only a version string. | accepted | Spec adds `spec_subject` and binds rollback to source, spec, gate set, proof time, and external audit for production. | closed in spec/plan |
| B07-F016 | minor | pi reviews | Profile names were inconsistent for optional E2E pilot. | accepted | Spec uses `verify:slice`; optional slices are not peer baseline profiles. | closed in spec/plan |
| B07-F017 | minor | pi reviews | Review ledger has no schema yet. | accepted | Spec permits first Markdown ledger and requires `schema/review-ledger.schema.json` in Slice 6. | closed in spec/plan |
| B07-F018 | critical | pi post-review | `production_release_verified` remained ambiguous under `source_bound_local_release`. | accepted | Spec marks production release verification true as forbidden under source-bound local profile and valid only under external production trust. | closed in spec/plan |
| B07-F019 | major | pi post-review | Committed trust-result JSON artifacts could still be mistaken for authoritative verifier output. | accepted | Plan Slice 2 requires regeneration, removal, or explicit untrusted/example labels for trust-result JSON. | closed in spec/plan |
| B07-F020 | informational | pi closure check | Post-edit closure check found no remaining critical or major blockers for B07-F018/B07-F019 and related stale-claim fixes. | recorded | No spec or plan change required. | closed |
| B07-F021 | critical | human review | "Deferred Work" made external trust read like another future promise instead of a current product blocker. | accepted | Spec adds Non-Deferrable Trust Work; plan adds a dedicated external self-release evidence slice and replaces Deferred Work with Non-Goals After Trust Closure. | closed in spec/plan |

## Slice 0 Implementation Review Findings

| ID | Severity | Source | Finding | Disposition | Spec or plan change | Closure |
| --- | --- | --- | --- | --- | --- | --- |
| B07-S0-F001 | critical | pi / MiniMax-M2.7 security | Dirty checkout with `--allow-dirty` could be misread as trusted verifier output. | accepted | Added `trust_scope`; dirty default is `cannot_verify`, dirty `--allow-dirty` is `local_dirty_structural_only`. | closed by `scripts/verify.sh`, schema, tests |
| B07-S0-F002 | critical | pi / zai/glm-5-turbo bash review | Validation functions could false-pass because multi-command functions were executed under `set +e`. | accepted | Each schema/example and negative-fixture check now returns on first failure. | closed by `scripts/verify.sh` and tests |
| B07-S0-F003 | major | pi / zai/glm-5.1 SpecKit review | Missing gate-set files were silently skipped, producing an incomplete gate digest. | accepted | Added required `gate_set_files_present` verifier state. | closed by `scripts/verify.sh` |
| B07-S0-F004 | critical | pi / zai/glm-5.1 bash review | Missing `jq` could make state recording and aggregation fail toward a false pass. | accepted | Added required tool guards for `git`, `jq`, `node`, and SHA-256 tooling before verification starts. | closed by `scripts/verify.sh` |
| B07-S0-F005 | major | pi / MiniMax-M2.7 security | Checked-in proof-summary example contained real-looking digests and could be mistaken for proof. | accepted | `--example` now emits fixture-only zero commit/digests/time and `trust_scope: "untrusted_shape_only"`. Schema enforces those fields for `untrusted_shape_example`. | closed by schema, verifier, generated example, and closure review |
| B07-S0-F006 | critical | pi / openrouter/qwen/qwen3-coder DX | Byte-identical proof-summary replayability was requested for `generated_at` and git commit fields. | rejected | Live proof summaries intentionally record proof time and selected source subject. Replayability means rerunnable verifier states, not stable cache bytes. | closed as non-requirement |
| B07-S0-F007 | informational | pi / zai/glm-5.1 closure review | Verified `untrusted_shape_example` schema constraints for fixture scope, zero commit, zero digests, fixed timestamp, and `--example` invocation. | recorded | No further change required. | closed |

## Slice 1 Implementation Review Findings

| ID | Severity | Source | Finding | Disposition | Spec or plan change | Closure |
| --- | --- | --- | --- | --- | --- | --- |
| B07-S1-F001 | major | subagent claim grammar review | Claim tag grammar must stay strict and avoid NLP over prose; T070 `pass` must fail unless current command evidence passes. | accepted | Slice 1 grammar narrowed to required `claim`, `subject`, `state`, `profile`, and `evidence` fields. | closed by validator and tests |
| B07-S1-F002 | major | subagent trace coverage review | Slice 1 must be covered by self-trace graph, evidence, and provenance artifacts. | accepted | Added Slice 1 nodes, evidence events, provenance records, observation, and accountability node under `examples/self-trace/`. | closed by self-trace validation |
| B07-S1-F003 | major | pi / MiniMax-M2.7 security | Arbitrary `state:*`, `proof:*`, and `none` evidence looked like proof but was not verified by Slice 1. | accepted | Slice 1 now accepts only `command_set:block04-t070` and `state:claim_tags_consistent`; `state=pass` requires executable command evidence. | closed by tests and closure review |
| B07-S1-F004 | critical | pi / MiniMax-M2.7 security | Claim tags are not cryptographically authentic and can be edited by anyone with write access. | accepted as existing trust-anchor boundary | Block 07 Trust Anchor Assumption already states checkout/source distribution trust is external until an external production profile is selected. Slice 1 does not claim cryptographic authenticity. | remains external-trust blocker, not Slice 1 blocker |
| B07-S1-F005 | major | pi / Kimi trace micro-review | Slice 1 trace lacked an accountability node and edges. | accepted | Added `accountability-block-07-slice-1` and accountability edges from claim-consistency changes. | closed by closure review |
| B07-S1-F006 | major | pi / Kimi trace micro-review | Slice 1 trace lacked a dedicated metric stream. | rejected for Slice 1 | No current metric has been defined for Slice 1 movement. Adding a synthetic metric would overclaim observability; metric coverage can be added after a real metric is specified. | closed as non-requirement |
| B07-S1-F007 | informational | pi / ZAI GLM SpecKit review | Claim grammar, prose ignoring, stale T070 handling, and verifier integration matched the Slice 1 spec. | recorded | No change required. | closed |
| B07-S1-F008 | informational | pi / MiniMax and Kimi closure reviews | Final closure checks found no critical or major findings for restricted evidence grammar, T070 pass replay, trace evidence, or accountability. | recorded | No change required. | closed |

## Blocking State

## Slice 2-5 Implementation Review Findings

| ID | Severity | Source | Finding | Disposition | Spec or plan change | Closure |
| --- | --- | --- | --- | --- | --- | --- |
| B07-S2-F001 | critical | subagent Slice 2 audit | T070 and Block 04/06 artifacts overclaimed current closure and proof completion. | accepted | T070 reopened as stale; Block 04 and Block 06 docs now match verifier/package reality. | closed by docs, JSON, and tests |
| B07-S2-F002 | major | subagent Slice 2 audit | Checked-in release verification and self-attestation JSON looked authoritative while manifest verification mismatched. | accepted | Release verification is labeled untrusted/failing local incomplete evidence; self-attestation expects `digest_verified: false`. | closed by schema, JSON, and self-attestation verifier |
| B07-S4-F001 | major | subagent verifier audit | E2E package validator conflated incomplete package shape with completed proof. | accepted | Added `--mode package` and `--mode complete`; committed incomplete package validates only in package mode. | closed by tests |
| B07-S5-F001 | major | subagent trace audit | Cross-reference integrity needed mechanical validation before broader evidence refs can be trusted. | accepted | Added `scripts/validate-cross-references.mjs` and wired it into baseline verifier. | closed by validator and negative fixture |

## Blocking State

The human owner authorized implementation start on 2026-05-03. Slice 0 may be used only as a structural verifier baseline. Slice 1 validates committed authoritative claim tags. Slices 2, 3, 4, and 5 now have verifier/test coverage. Block 07 remains open until external self-release evidence is available or the external production trust profile remains explicitly failed/cannot_verify under live verifier output, and until Slice 7 final ledger/schema closure is completed.
