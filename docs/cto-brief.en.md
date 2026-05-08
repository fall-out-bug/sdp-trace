# sdp-trace For CTOs In One Minute

AI-assisted delivery increases speed, but it weakens control: a week later it is hard to know what was promised, who or what did the work, what evidence exists, what is missing, and who owns the next step.

`sdp-trace` is not here to say "everything is fine" or "we are degrading." It exists to record a verifiable chain:

```text
idea -> spec -> task -> change -> evidence -> provenance -> accountability -> metric movement -> verified contract
```

For the CTO, this means the process can be inspected over time. Not through opinions or opaque scores, but through prior/current/delta values, evidence coverage, and explicit `not_assessed` gaps.

For the CEO, this means the accountable owner is not AI. Every significant artifact has a human-held DRI, approver, risk owner, and escalation path.

For the CIO, this means the contract cannot be silently simplified. Schemas,
docs, validation commands, fixtures, and release-proof artifacts can be checked
against explicit contract manifests and authority scopes.

`sdp-trace` deliberately does not make policy decisions: pass/fail, readiness, degradation, thresholds, and overrides belong to `sdp-gate` or another external policy consumer.

The current product surface supports controlled pilots: local trace capture,
report packages, missing-telemetry visibility, assessment profiles, CI/customer
witness artifacts when evidence exists, and source-bound local release proof.

It does not claim broad production trust, universal harness compatibility, or
automatic detection of every unwrapped local agent run.

Start with `docs/README.md`, then use `docs/cto-adoption-guide.en.md` as a
supporting reference. Development specs and research notes are useful for audit
history, but they are not the onboarding path.
