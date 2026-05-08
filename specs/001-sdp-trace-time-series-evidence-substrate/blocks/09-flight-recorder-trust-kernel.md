# Block 09: Flight Recorder Trust Kernel

Status: design accepted for implementation planning; no implementation closure claimed
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Audience: technical executive, CIO, CISO, COO, CPO, Head of Engineering, corporate architecture, implementation agents

## Purpose

Block 09 repositions `sdp-trace` as a flight recorder for agentic development.

The product must not ask teams to fill in an evidence envelope after the fact. It must record the development run while it happens, bind the record to an append-only event chain, and make post-run rewriting detectable.

`sdp-trace` still does not judge code quality, decide readiness, fix failures, or produce policy verdicts. Its job is narrower and harder:

> Record provenance, evidence, and trace from attachment forward so the recorded past cannot be silently rewritten.

The Feature Flag / Entitlements Kotlin+Bazel demo is deliberately not part of this block. The demo is a later stress test for the recorder. Block 09 must first make the recorder product mechanically credible.

## Executive Review Outcome

The executive Socratic review rejected a local-only recorder as insufficient.

Converged findings:

- A local hash chain proves internal consistency, not authenticity.
- A chain controlled by the same actor being reviewed can be deleted, replaced, or recomputed.
- Mid-flight attachment is useful only as an honest scope boundary; it cannot support claims about pre-attachment activity.
- Digest-only evidence is too weak for forensic review unless paired with retained, sanitized, encrypted, or externally referenced artifacts.
- Redaction cannot be hand-waved because append-only retention and secret removal pull in opposite directions.
- technical executive/CISO/COO users need queryable proof surfaces, not raw JSONL archaeology.

Block 09 therefore requires a recorder plus a witness boundary plus a verifier/query surface. A sidecar without a witness is a development aid, not an accountability substrate.

## Product Thesis

`sdp-trace` is valuable only if it reduces the cost of answering:

> What was fixed, by whom or what, from what source state, under what task, with which model and harness, through which commands and file changes, and where is the record missing or unverifiable?

The answer must be based on recorded events, not on a polished final summary.

The correct product claim is:

> From the moment the recorder is attached and witnessed, post-run mutation of recorded provenance, evidence, or trace is detectable.

The incorrect product claims are:

- the recorder proves pre-attachment history
- the recorder proves code quality
- the recorder prevents bad changes
- the recorder produces gate pass/fail decisions
- a local JSONL chain is audit-grade by itself

## Trust Model

### Actors

- `operator`: human or automation that starts the run.
- `agent_harness`: OpenCode, GSD, or another agent runner.
- `model_provider`: the provider or routing surface used by the harness.
- `recorder`: `sdp-trace` process that observes commands, files, artifacts, and imported transcripts.
- `witness`: storage or notary boundary not controlled by the recorded agent process.
- `verifier`: code that checks event-chain structure, witness anchors, artifacts, and declared gaps.
- `consumer`: technical executive, reviewer, CISO, CI policy, or downstream external policy consumer.

### Threats In Scope

- post-run event mutation
- event deletion inside the recorded chain
- event reordering
- final-summary rewriting
- task/prompt expectation rewriting after evidence exists
- command-result rewriting
- file-diff rewriting
- local trace deletion when a witness anchor is expected
- late attachment being misrepresented as full-run proof
- redaction being used to hide unresolved evidence

### Threats Not Solved by Local Mode

- an actor choosing not to start the recorder
- activity before recorder attachment
- kernel-level compromise
- compromised verifier source checkout
- compromised external witness
- model-provider false identity unless provider evidence is captured
- in-process harness actions not exposed to shell, file, transcript, or adapter observation

These unsolved threats must be reported as `not_assessed`, `cannot_verify`, or `fail` depending on the selected profile. They must not be hidden in prose.

## Recorder Modes

### `local_development_recorder`

Records a local hash chain and local artifacts. Useful for development and debugging. It must not support accountability or external-trust claims.

Required output:

- event chain
- run manifest
- local artifact refs
- explicit `trust_scope: local_development_recorder`
- witness states as `not_assessed`

### `witnessed_run_recorder`

Records a local event chain and anchors chain heads to a witness boundary during or immediately after the run.

Minimum acceptable witness for Block 09 implementation:

- append chain head to a file or store outside the run artifact directory
- bind witness entry to run id, source baseline hash, task hash, recorder version, and chain head
- make verifier fail when the event chain and witness entry disagree

This local witness is not external production trust. It is a first implementation step that proves the witness contract. Future profiles may replace it with S3 Object Lock, transparency log, timestamp authority, protected Git ref, customer audit log, or Sigstore/Rekor.

### `externally_witnessed_run`

Future profile. Requires a witness outside the developer-controlled machine or checkout. Block 09 must define the extension point but does not have to implement a production witness.

## Event Chain Requirements

Every event must include:

- `schema_version`
- `run_id`
- `sequence`
- `event_type`
- `event_time`
- `event_payload_digest`
- `prev_event_hash`
- `event_hash`
- `recorder_identity`
- `redaction_state`
- `witness_ref` when anchored

Canonical serialization rules must be deterministic and versioned. Hash algorithms must be declared per event and per chain. If canonicalization changes, the schema version must change.

Required event families:

- `run_started`
- `source_baseline_recorded`
- `task_locked`
- `expectation_locked`
- `model_identity_observed`
- `harness_identity_observed`
- `command_started`
- `command_finished`
- `file_state_observed`
- `file_mutation_observed`
- `test_output_observed`
- `artifact_captured`
- `requirement_superseded`
- `redaction_applied`
- `witness_anchor_recorded`
- `run_closed`
- `recorder_interrupted`

Requirement changes are never edits. They are `requirement_superseded` events that reference the earlier task or expectation event.

## Evidence Capture Requirements

Block 09 must not rely on digest-only evidence for every proof surface.

Each captured artifact must choose one evidence-retention mode:

- `digest_only`: only hash and metadata retained; useful but low forensic value.
- `sanitized_excerpt`: selected safe excerpts retained plus full raw digest.
- `encrypted_raw_ref`: raw artifact encrypted or sealed; key held outside the run artifact.
- `external_artifact_ref`: raw artifact stored in an external log or CI artifact with digest and access note.
- `not_assessed`: evidence unavailable; reason and next required evidence mandatory.

For shell command events, the recorder must capture at minimum:

- command argv as structured array, with redaction metadata
- argv digest over canonical structured argv, not a shell-joined string
- working directory
- start and end time
- exit code or signal
- stdout and stderr retention mode
- pre and post git status summary
- pre and post source tree or diff digest for the declared scope

For file mutation events, the recorder must bind mutations to an observed source scope. If attribution to a specific command is unavailable, that field must be `not_assessed` rather than inferred.

## Redaction Model

Redaction is part of the proof model, not a formatting pass.

Block 09 must specify and test:

- redaction before persistence for known secret patterns
- event-level redaction metadata
- retained digest of the redacted payload
- optional sealed raw payload reference when allowed
- redaction authority field
- reason code
- verifier behavior for unresolved redaction

Forbidden:

- storing raw secrets in committed examples
- claiming a redacted event proves the original raw value unless a sealed raw reference exists
- silently replacing event payloads after chain closure

If redaction cannot be completed safely, the state is `cannot_verify` or `not_assessed` with a reason. It is not a pass.

## Verifier Semantics

The verifier must distinguish structure, witness, completeness, and forensic usefulness.

Required states:

| State | Meaning |
| --- | --- |
| `event_chain_structurally_valid` | Event hashes, sequence, and previous-hash links verify under declared canonicalization. |
| `event_chain_witnessed` | Declared witness entry matches run id and chain head. |
| `source_baseline_recorded` | Source baseline event exists and is bound to a git commit, tree, or directory snapshot. |
| `task_locked` | Task/prompt/expectation event exists before dependent command/model events. |
| `late_attach_boundary_explicit` | Pre-attachment history is marked as `not_assessed` when applicable. |
| `model_identity_recorded` | Requested and observed model identity are recorded, or the gap is explicit. |
| `command_events_bound` | Command start/end events are paired and have exit state and artifact refs. |
| `file_mutations_bound` | Final diff or mutation set is bound to source baseline and run closure. |
| `redaction_resolved` | Redaction state is safe, explicit, and verifier-readable. |
| `run_closed` | Run closure event exists and references final chain head and artifact set. |

Exit behavior:

- `0`: selected profile passes.
- `1`: verifier ran and required evidence failed or mismatched.
- `2`: invalid invocation or unsupported profile.
- `3`: verifier could not verify because required tool, source, witness, or artifact access was unavailable.

Profiles:

- `flight_recorder_local`: chain structure and local artifact shape only; witness states may be `not_assessed`.
- `flight_recorder_witnessed`: chain structure plus witness agreement; missing witness fails.
- `flight_recorder_forensic`: witnessed chain plus evidence-retention requirements for reviewer reconstruction; digest-only critical events fail unless explicitly accepted by profile.

## Query Surface

Raw JSONL is not an acceptable product surface for technical executive or incident review.

Block 09 must define query outputs for:

- run summary
- source/task/model/harness provenance
- late attach gaps
- requirement supersession timeline
- command timeline
- failed or missing command evidence
- file mutation summary
- test evidence summary
- redaction unresolved items
- witness verification result

Queries may be CLI-first. They must use verifier-derived states and must not produce policy verdicts.

Example questions:

- What task was locked before the first command?
- Did any requirement change after a failed command?
- Which commands were actually observed?
- Which files changed after attachment?
- Which claims are impossible because evidence is `not_assessed`?
- Is the chain witnessed or only locally consistent?

## Out of Scope

- external production witness implementation
- full OpenCode/GSD adapter implementation
- Feature Flag / Entitlements demo repo
- code-quality scoring
- readiness or compatibility verdicts
- broad harness/model support claims
- kernel-level anti-tamper guarantees
- mandatory enterprise deployment or endpoint control

## Acceptance Conditions

Block 09 spec is ready for implementation when:

- the threat model names what local mode cannot prove
- local chain consistency is not described as audit-grade trust
- witness anchoring is required for accountability profiles
- event schema and canonicalization requirements are explicit
- redaction behavior has verifier states
- late attach is treated as a visible gap, not a convenience that preserves full-run proof
- query surfaces are defined before demo work starts
- implementation tasks are small enough for agent handoff

Block 09 implementation is complete only when:

- positive and negative fixtures exist for local and witnessed chains
- tampering with an event is detected
- deleting or changing a witness entry is detected in witnessed mode
- late attachment produces an explicit `not_assessed` boundary
- requirement supersession cannot rewrite the original task event
- redaction unresolved state fails the relevant profile
- verifier and query commands are documented
- review findings are recorded and closed or left blocking
