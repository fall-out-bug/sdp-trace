# sdp-trace for CTOs in One Minute

AI-assisted delivery increases speed, but it weakens control: a week later it is hard to know what was promised, who or what did the work, what evidence exists, what is missing, and who owns the next step.

`sdp-trace` is not here to say "everything is fine" or "we are degrading." It exists to record a verifiable chain:

```text
idea -> spec -> task -> change -> evidence -> provenance -> accountability -> metric movement -> verified contract
```

For the CTO, this means the process can be inspected over time. Not through opinions or opaque scores, but through prior/current/delta values, evidence coverage, and explicit `not_assessed` gaps.

For the CEO, this means the accountable owner is not AI. Every significant artifact has a human-held DRI, approver, risk owner, and escalation path.

For the CIO, this means the contract cannot be silently simplified. Schemas, docs, validation scripts, and fixtures are covered by a contract manifest with digest verification and a release verification profile.

`sdp-trace` deliberately does not make policy decisions: pass/fail, readiness, degradation, thresholds, and overrides belong to `sdp-gate` or another external policy consumer.

Block 01 builds the contract scaffold: evidence contracts, accountability, manifest verification, signing profile, negative fixtures, and proof that missing data remains `not_assessed`.

The product is not allowed to ask a customer for trust until it traces itself. The next proof must show this repository's own spec, plan, tasks, changes, evidence, provenance, accountability, reviews, metrics, and missing data under the same contracts.

Repository evidence starts at `specs/001-sdp-trace-time-series-evidence-substrate/`, `schema/README.md`, `examples/contract-foundation/`, and `docs/research/block-01-validation-summary.md`.
