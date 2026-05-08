# Block 24 Socratic Review: Demo Repository CI And Trace Pilot

Status: Socratic review and focused re-review complete. Implementation is
blocked until the CTO explicitly approves the reviewed direction.

## Initial Socratic Questions And Resolutions

### Q1: Is a same-owner separate repository enough?

**Critic**: Same-owner GitHub proof can be dismissed as a controlled fixture,
not true customer portability.

**Resolution**: Same-owner separate repo is enough for Block 24 because the
current gap is repository attachment and CI trace inspectability, not customer
tenant independence. The report must say portability across another owner,
different access controls, and non-GitHub providers remains `not_assessed`.
The report must also name what a different owner must provide: CI workflow edit
access, artifact retention policy, a `sdp-trace` source or release ref, and OIDC
or equivalent witness permission if CI witness is in scope.

### Q2: Does GitHub Actions first overfit the demo?

**Critic**: Block 22 broadened witness semantics. Starting with GitHub Actions
could make GitHub look like the hidden product model.

**Resolution**: GitHub Actions is the right first live CI target because the
current repository already uses GitHub and the shipped CLI has a
`github-actions` witness profile. Block 24 must keep GitLab, Buildkite,
customer PKI, and air-gapped execution as `not_assessed`, not implied support.

### Q3: Is a Feature Flag / Entitlements Kotlin+Bazel app credible without
expanding Block 24 too far?

**Critic**: A trivial app can repeat the Block 06 mistake, but a full
Kotlin+Bazel/agent demo can sprawl beyond a CI trace pilot.

**Resolution**: Use the repo-native external demo target: a small Feature Flag /
Entitlements Kotlin+Bazel service. Keep the domain path narrow, but require a
real Bazel test over Kotlin source so the demo does not look like another toy
proof. The credible value is still the end-to-end CI trace path: wrap, verify,
explain, report, gate, witness, redaction scan, artifact index, and
customer-question mapping from a separate repo. The pilot report must not claim
broader JVM/Bazel, OpenCode, GSD, or model-agent support.
To avoid a low-signal demo, the report must include a "CI Alone vs sdp-trace"
contrast showing at least one structured fact preserved by `sdp-trace` that raw
CI logs alone do not preserve as durable evidence.

### Q4: Should Block 24 use `wrap` or `run`?

**Critic**: `run --task` gives better task traceability, while `wrap` is closer
to "attach to any command".

**Resolution**: Use `wrap` for the first pass because it minimizes contract
setup and answers the customer's attachment question. If implementation uses
`run`, it must record the task ref and contract choice explicitly. Either way,
`verify`, `explain`, `report`, `gate`, and `witness` must operate on the same
captured run root.

### Q5: How do we prove a negative path without fabricating failure?

**Critic**: A deliberately failed test can look like product failure, while a
fake witness file can create artificial evidence.

**Resolution**: Use missing authority rather than failed app behavior. A
no-OIDC GitHub Actions job or local-only witness attempt can produce
`cannot_verify` with reason `missing_ci_oidc` or missing identity/binding fields
without falsifying command execution.
The negative job must run as a separate job with no OIDC permission and must
capture the `witness` exit code and JSON artifact. A `cannot_verify` result means
the verifier attempted an in-scope witness check but lacked required identity or
binding evidence; it is not a pass and not automatically a product failure.

### Q6: Where do raw artifacts belong?

**Critic**: Copying `.sdp-trace-runs/` and `.sdp-trace-report/` wholesale into
this repository risks committing raw logs, paths, tokens, or user data.

**Resolution**: Raw artifacts stay in the demo repo or CI artifact store unless
they pass an explicit redaction review. This repository should record sanitized
links, digest indexes, state summaries, and short excerpts only where needed for
review.

### Q7: Does `gate` output create a policy claim?

**Critic**: A customer may read a gate result as a merge/release decision.

**Resolution**: The report must state that `gate` output is verifier-derived
fact output and not a native `sdp-trace` policy decision. Any downstream allow,
block, risk, readiness, or production-trust decision remains outside this repo.

### Q8: Is building `sdp-trace` from source in CI enough?

**Critic**: Building from source proves developer workflow, not installable
product UX.

**Resolution**: It is enough for Block 24 if the source ref is explicit and the
tested command surface is current. Installable binary/release distribution stays
`not_assessed` unless an approved release artifact exists before implementation.

### Q9: How do we keep demo evidence from becoming source-bound proof?

**Critic**: A `ci_witnessed` demo artifact can be misread as product release
proof.

**Resolution**: Every Block 24 artifact copied into this repository must carry
`authority_scope=demo_pilot_only`. The source-bound release state remains
separate and external production trust remains `not_assessed`.

### Q10: What makes redaction claims verifiable?

**Critic**: "Redaction scan" as prose is not evidence.

**Resolution**: The implementation must run a concrete denylist scan using
`archive/research/block-24-redaction-denylist.patterns`, record the pattern-file
digest, scanned roots, exit state, and match count, and treat a missing scan as
`cannot_verify`.

## Review Plan

Run separate review planes:

1. product/demo credibility;
2. trace/evidence and claim-boundary review;
3. CI/witness semantics review;
4. privacy/safety review.

Record every finding in
`24-demo-repo-ci-trace-pilot-review-ledger.md`. Valid critical or major
findings block approval until fixed or explicitly marked `unresolved_blocker`.
