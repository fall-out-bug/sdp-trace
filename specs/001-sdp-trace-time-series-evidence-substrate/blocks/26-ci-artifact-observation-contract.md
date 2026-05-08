# Block 26: CI Artifact Observation Contract

Status: spec draft, pending Socratic review.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/22-additional-ci-enterprise-witness-profiles.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/24-demo-repo-ci-trace-pilot.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/25-compiled-ci-demo-pilot.md`

## Goal

Make CI artifact proof gaps first-class `sdp-trace` observations.

The product must record when provenance, evidence, trace, witness, and report
claims are backed by uploaded CI artifacts, when they are only checked into git,
when they are locally generated, when they are agent-reported, and when the
artifact bundle cannot be inspected. The output is a fact package, not a native
merge, release, readiness, audit, or policy decision.

## Problem

Block 25 exposed a product gap: the repository already has witness profiles,
test provenance checks, and artifact digest binding, but a demo can still look
credible while the important provenance/evidence/trace bundle is absent from CI
artifacts or replaced by checked-in JSON and prose.

That is the wrong failure mode. `sdp-trace` should not fix demo defects to make
the happy path look clean. It should capture facts that make the defect visible:

- a happy-path claim exists but no downloaded CI artifact bundle was inspected;
- `provenance`, `evidence`, or `trace` exists only as committed JSON;
- a CI run exists but artifact upload metadata is missing or incomplete;
- an uploaded bundle exists but lacks required artifact families;
- an uploaded bundle contains stale, mismatched, unsafe, or self-referential
  digest indexes;
- a run claims OpenCode, GSD, MiniMax, compiled JVM tests, PR review, or feature
  delivery without corresponding observed artifact families.

## Non-Goals

- No native policy gate, merge decision, release approval, risk acceptance,
  readiness score, health score, grade, or rank.
- No requirement that every repository use CI artifacts. Repositories without
  uploaded CI artifacts are `not_assessed` when the selected profile does not
  require CI uploads, and `cannot_verify` when CI-uploaded evidence is required
  but absent.
- No dependency on GitHub Actions, OpenCode, GSD, MiniMax, Bazel, Kotlin, or any
  demo repository.
- No validation of private artifact download authentication in this block.
  Access state is recorded as `not_assessed` or `cannot_verify` unless the
  selected profile supplies verifier-readable metadata.
- No automatic protection of private repository names, refs, run ids, workflow
  ids, or commit identifiers. Observation results are intended for the same
  security boundary as the inspected artifacts unless a later profile adds
  path-redacted or hashed identifiers.
- No network artifact fetching in the first implementation. Future network
  resolution must be explicit, HTTPS-only, credential-free in URLs, and bounded
  by configured provider/domain policy before it can be trusted.
- No raw CI logs, raw OIDC tokens, provider tokens, private artifact URLs,
  command bodies, prompts, model responses, or private filesystem paths.
- No automatic repair of demo defects. Negative demo cases are evidence.

## Observation Model

Block 26 adds a provider-neutral CI artifact observation profile. The profile
may be run against a local extracted artifact root, a provider-supplied artifact
manifest, or both.

The profile records:

- selected source repository, ref, commit, run id, run attempt, job id, and
  workflow id when available;
- artifact producer and producer scope: `ci_uploaded`, `checked_in`,
  `local_generated`, `agent_reported`, `harness_observed`,
  `external_artifact_ref`, or `not_assessed`;
- artifact access state: `present`, `absent`, `partial`, `expired`,
  `inaccessible`, `malformed`, `unsafe`, `not_assessed`, or `cannot_verify`;
- required artifact families and their states;
- source/run binding state;
- artifact index state;
- output safety state;
- residual gaps and next evidence needed.

The profile must preserve existing verifier states. Contradictory evidence is
`fail`. Missing or inaccessible evidence is `cannot_verify` when the profile
requires it. Evidence outside the selected profile is `not_assessed`. These
states are facts for consumers; they are not native enforcement.

Producer scope is per artifact family, not only top-level. A single observation
may record a CI-uploaded run directory, a checked-in report, and an
agent-reported review claim at the same time.

Producer scope values mean:

| Producer scope | Meaning |
| --- | --- |
| `ci_uploaded` | The family was observed in an uploaded CI artifact bundle or provider artifact manifest selected by the profile. |
| `checked_in` | The family was observed only in git-committed files. |
| `local_generated` | The family was produced by a local process without a verifier-readable CI or harness boundary. |
| `agent_reported` | The family was claimed by an agent, README, PR body, or prose without observed artifact evidence. |
| `harness_observed` | A recorded harness or runner observed the family, but the selected profile has no CI-upload or external witness binding for it. |
| `external_artifact_ref` | The family is referenced by an external artifact pointer with digest or access metadata. |
| `not_assessed` | The profile did not inspect producer scope for this family. |

When a selected profile requires `ci_uploaded` and the observed family is only
`checked_in`, `local_generated`, `agent_reported`, or `harness_observed`, that
family's producer binding is `mismatch`, its proof-level state contributes
`cannot_verify`, and the reason code must name the lower-authority producer
scope. This is not a native policy failure; it records that the required proof
level was not observed.

## State Derivation

Top-level `artifact_observation_state` is derived deterministically:

| Condition | Top-level state |
| --- | --- |
| Any selected required family, binding, index, or output-safety check has verifier-visible contradiction or unsafe output | `fail` |
| No `fail`, but any selected required family is absent, partial, expired, inaccessible, malformed, unverifiable, or produced below the required proof level | `cannot_verify` |
| No selected required family was assessed because the profile scope does not include it | `not_assessed` |
| Every selected required family is present, bindings match, indexes are valid, and output safety passes | `pass` |

Reasons accumulate. If one family is expired and another is partial, both reason
codes appear; the aggregate state is the highest state from the table.

Field vocabularies:

| Field | Allowed values |
| --- | --- |
| `artifact_observation_state` | `pass`, `fail`, `cannot_verify`, `not_assessed` |
| `producer_scope` | `ci_uploaded`, `checked_in`, `local_generated`, `agent_reported`, `harness_observed`, `external_artifact_ref`, `not_assessed` |
| `artifact_access_state` | `present`, `absent`, `partial`, `expired`, `inaccessible`, `malformed`, `unsafe`, `not_assessed`, `cannot_verify` |
| `family_state` | `pass`, `fail`, `cannot_verify`, `not_assessed` |
| `binding_state` | `matched`, `mismatch`, `absent`, `unverifiable`, `not_assessed` |
| `artifact_index_state` | `valid`, `self_reference`, `digest_mismatch`, `missing`, `unverifiable`, `not_assessed` |
| `output_safety_state` | `pass`, `fail`, `cannot_verify`, `not_assessed` |

`partial` is a per-family access state: at least one artifact in the selected
family is present, but the profile-required set is incomplete. It contributes
`cannot_verify` to the top-level derivation.

Reason strings and next actions must be safe templated messages derived from
closed reason codes and safe categories. They must not embed raw artifact
content, raw paths, raw URLs, log lines, parser input, prompts, model output, or
token-like values.

## Required Artifact Families

For a demo-style product proof package, the profile must be able to classify at
least these artifact families:

| Family | Description | Minimum observation |
| --- | --- | --- |
| `run` | wrapped or recorded run directory | run id, source ref or commit, command/test summary digest; closes the gap where a report exists without an observed run |
| `report` | `sdp-trace` report or assessment output | schema/version and selected profile; closes the gap where prose summarizes output that was not generated |
| `witness` | CI or external witness output | witness kind, trust scope, source/run/artifact binding state; closes the gap where CI identity is inferred from environment text |
| `provenance` | provenance records | actor, harness/model/tool fields or unavailable fields; closes the gap where agent/model participation is asserted but not observed |
| `evidence` | evidence bundle or events | item sources, status origin, digests where available; closes the gap where evidence status is copied from prose |
| `trace` | trace graph or snapshot | spec/plan/task/change/evidence links; closes the gap where feature work has artifacts but no traceable relationship |
| `artifact_index` | digest index or provider manifest | path set, digest set, self-index status; closes the gap where a parallel or stale upload is mistaken for selected artifacts |
| `redaction_scan` | artifact-bundle safety scan result | scanned roots, rule or policy digest, state; closes the gap where unsafe output is uploaded but hidden in artifact contents |
| `review` | review evidence and disposition | review plane, reviewer source, disposition state; closes the gap where "reviewed" is only a PR comment or agent claim |
| `change_ci` | change, merge-request, PR, or branch CI check evidence | change/head SHA, run id, check state or not assessed reason; closes the gap where a branch claims CI without artifact-backed check evidence |

This table is the minimum set the product must be able to classify for demo and
pilot proof packages. It is not the minimum set every invocation must supply.
Required families are profile inputs; unselected families remain
`not_assessed`. The product must not hard-code the demo's exact artifact names
as product concepts. `pr_ci` may be accepted as a backwards-compatible alias,
but `change_ci` is the canonical family name.

`change_ci` records whether the artifact bundle contains evidence for a
change/branch CI claim. It does not re-evaluate the source provider's PR,
merge-request, or branch policy.

`redaction_scan` records whether the candidate artifact bundle contains a
reviewable scan result and whether that scan is bound to the inspected root.
`output_safety` is a separate evaluator check over the observation output before
serialization.

## Failure Modes To Preserve

The first implementation must include fixtures for these cases:

1. `checked_in_only_claim`: required families exist in git but no CI artifact
   bundle is present. Expected state: `cannot_verify` for CI-backed proof, with
   `producer_scope=checked_in`.
2. `ci_bundle_absent`: CI identity exists, but no artifact upload or extracted
   bundle is supplied. Expected state: `cannot_verify`.
3. `ci_bundle_partial`: uploaded bundle exists but omits one or more required
   families. Expected state: `cannot_verify` with per-family gaps.
4. `artifact_index_self_reference`: the index lists itself as a digested entry.
   Expected state: `fail`.
5. `artifact_digest_mismatch`: an indexed artifact digest contradicts selected
   bytes. Expected state: `fail`.
6. `source_run_binding_missing`: source commit, run id, or run attempt binding is
   absent. Expected state: `cannot_verify`.
7. `source_run_binding_mismatch`: source commit, ref, run id, or run attempt
   contradicts the selected input. Expected state: `fail`.
8. `agent_reported_happy_path`: an agent or README claims feature completion,
   tests, PR, CI, evidence, or trace, but the required observed family is absent.
   Expected state: `cannot_verify` when the selected profile requires the
   family; `not_assessed` only when the family is outside profile scope. Reason
   code: `agent_reported_claim_without_observed_family`.
9. `unsafe_artifact_output`: artifact metadata or persisted output contains
   forbidden token-like, prompt, raw log, private URL, or private path material.
   Expected state: `fail` with safe reason metadata.
10. `artifact_expired`: the provider or supplied manifest records that artifacts
    expired before inspection. Expected state: `cannot_verify`, not failure.
11. `external_artifact_ref_unverifiable`: an external artifact reference exists
    but cannot be resolved or digest-checked under the selected profile.
    Expected state: `cannot_verify`, unless the profile explicitly accepts that
    external reference authority.

Incomplete demo features must use reason code
`family_absent_in_ci_bundle`. Dishonest checked-in or prose claims contradicted
by CI artifacts must use reason code
`checked_in_claim_contradicts_ci_artifacts`.

## Customer Implications

Block 26 records whether required artifact families were present in uploaded CI
bundles and whether their digest bindings, source/run bindings, and safety state
are consistent. It does not claim:

- uploaded CI artifacts are secure, tamper-proof, or free from exfiltration;
- artifact upload implies successful build, test, review, or release readiness;
- agent-reported, harness-observed, checked-in, or local evidence satisfies a
  selected CI-uploaded proof requirement;
- a `pass` state is compliance, certification, audit readiness, or release
  approval.

Customer-readable state implications:

| State | Customer implication |
| --- | --- |
| `pass` | The selected observation profile found the required artifact families and bindings. Downstream policy still owns decisions. |
| `fail` | Contradictory or unsafe evidence was observed. Do not treat this bundle as provenance evidence without remediation and a new observation. |
| `cannot_verify` | Required evidence was absent, expired, inaccessible, incomplete, or below the selected proof level. Treat this as no usable proof for that profile, not as proof of absence. |
| `not_assessed` | The profile did not inspect this surface. Ask for profile scope or run a profile that selects the missing family. |

Unsafe output in CI artifacts may indicate credential leakage, prompt exposure,
private path disclosure, or unauthorized data egress. Detection does not imply
prevention; remediation remains with the CI owner and artifact consumer.

## Demo Forward Checks

The next demo iteration should use Block 26 as the truth contract instead of
hand-written proof prose.

Two to five steps ahead, the demo should be able to show:

1. Five feature branches, each with its own PR, CI run, artifact bundle, and
   Block 26 observation result.
2. At least one intentionally incomplete feature where code exists but CI
   artifacts omit provenance or trace. The observation must record the gap.
3. At least one intentionally dishonest feature where a checked-in report claims
   success but uploaded CI artifacts are absent or contradictory. The observation
   must record the contradiction.
4. At least one feature where OpenCode + GSD + MiniMax participation is claimed.
   Unless the artifact family records observable harness/model provenance, the
   model/harness claim remains `agent_reported`, `harness_observed`,
   `not_assessed`, or `cannot_verify`.
5. A customer-readable summary generated from Block 26 observation facts, not
   from manually curated happy-path prose.

The demo may still fail. A failed or incomplete demo is valid evidence if the
failure is captured with deterministic states and safe reasons.

## Product Surface

Block 26 may add:

- a schema for CI artifact observation results;
- a small Go package that evaluates extracted artifact roots and optional
  provider-neutral manifests;
- a CLI profile under `assess`, `query`, or a new observation-oriented command
  if that creates clearer DX;
- fixtures under `examples/block26-ci-artifact-observation/`;
- docs explaining how downstream gates can consume the observation without
  making `sdp-trace` itself an enforcement product.

The implementation must avoid adding Node, npm, JavaScript, or TypeScript to the
active product path.

CLI selection principle: prefer an observation-oriented command or an explicit
`assess --profile ci-artifact-observation` shape. Do not hide this under
forensics query output alone, because the operator must be able to create the
observation before querying it.

## Acceptance Criteria

Block 26 is not complete until all of these are true:

1. The schema records selected source/run identity, producer scope, access
   state, required family states, binding states, index state, output safety,
   reason codes, and next actions.
2. Fixtures cover every failure mode listed above and one valid observed bundle.
3. The evaluator emits `pass`, `fail`, `cannot_verify`, and `not_assessed`
   states without collapsing them into policy decisions.
4. Checked-in-only evidence cannot satisfy a CI-uploaded artifact requirement.
5. Agent-reported or prose-reported happy-path claims cannot satisfy observed CI
   artifact families.
6. A reviewer can run documented commands against committed fixtures and inspect
   exactly which artifact families were present, absent, partial, unsafe,
   stale, mismatched, or not assessed.
7. Output-safety tests prove that raw logs, tokens, OIDC material, provider
   URLs, prompts, model responses, private paths, and unsafe parser input do not
   appear in JSON or explain output. The active safety ruleset id and digest
   must be recorded, and failures must name safe rule ids or classes rather than
   echoing matched content.
8. Block 25 demo specs or future demo specs reference Block 26 observation
   results for artifact truth instead of relying on checked-in JSON or prose.
9. Socratic spec review and implementation review run separate
   product/customer, tracing/evidence, DX/replayability, and security/privacy
   planes with recorded dispositions.

## Review Plan

Socratic spec review must run before implementation approval across:

- product/customer credibility: does this expose the demo truth a CTO would care
  about without pretending to be a gate;
- tracing/evidence: do states, producer scopes, artifact families, and binding
  rules prevent checked-in JSON and prose from overclaiming CI proof;
- DX/replayability: can a repo operator produce and inspect the observation
  without adopting a managed harness or provider-specific runtime;
- security/privacy: does artifact inspection avoid leaking CI/provider/private
  material while still being useful.

Implementation review must repeat these planes and add code/correctness for the
Go evaluator and fixtures. PR-level review must repeat code/correctness,
tracing/evidence, requirements-vs-implementation, and security/privacy.

Any hung, empty, or off-task review is `not_assessed` and must be replaced.
