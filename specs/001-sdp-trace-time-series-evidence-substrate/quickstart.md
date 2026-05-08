# Quickstart: Reviewing the SpecKit Evidence Package

## 1. Start With the Spec

Read:

```text
specs/001-sdp-trace-time-series-evidence-substrate/spec.md
```

Confirm the feature answers the CTO question as evidence-backed process movement: prior/current values, deltas, dimensions, evidence coverage, and `not_assessed` gaps. It must not answer with a built-in yes/no degradation verdict.

## 2. Read the Boundary Contract

Read:

```text
specs/001-sdp-trace-time-series-evidence-substrate/contracts/sdp-trace-sdp-gate-boundary.md
```

Confirm `sdp-trace` owns evidence/provenance/observations/metric streams and `sdp-gate` owns policies, gate decisions, degradation verdicts, readiness, and overrides.

Also confirm external verdicts or evidence-strength assertions are represented as external inputs with producer and origin, not as native `sdp-trace` decisions.

## 3. Inspect the Plan and Tasks

Read:

```text
specs/001-sdp-trace-time-series-evidence-substrate/plan.md
specs/001-sdp-trace-time-series-evidence-substrate/tasks.md
```

Beads issues can mirror execution state, but these SpecKit artifacts are the repo-observable plan.

Check Phase 5 and Phase 5A before trusting any pilot claim. If self-trace or self-attestation tasks are incomplete, the repository has contract scaffolding only.

## 4. Inspect Accountability and Contract Integrity

For any artifact that claims accountability or trusted contract release status, confirm:

- AI actors are recorded only as producers, reviewers, critics, or judges, not as sole accountable owners.
- Accountable identities include machine-readable actor type and use human-held DRI, approver, risk owner, escalation path, approval reference, and line of defense.
- The contract manifest lists schemas, docs, validation scripts, fixtures, source commit, compatibility notes, and previous manifest digest when available.
- The trusted identity policy names expected OIDC issuer, source URI, protected ref, workflow identity, release captain, and required approval refs.
- The release verification record uses `sdp-trace-signature/sigstore-dsse-keyless-v1` or an explicitly documented private equivalent and matches the trusted identity policy.
- Missing signature or digest verification is recorded as `not_assessed` or invalid, not as trusted.

## 5. Run Current Schema Syntax Check

```bash
jq empty schema/*.json
```

This checks current schema JSON syntax.

## 6. Run Pinned Schema Validation When Schema/Example Pairs Exist

New schemas target JSON Schema Draft 2020-12. The Block 10-compatible
validator commands are:

```bash
go test ./...
go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc
```

`examples/github-speckit` remains schema/example material, not a current `validate-fixtures` run package. The generalized command from T036 must exclude `.git/`, `.beads/`, `.sdp-trace-runs/`, `benchmarks/repos/`, temporary directories, editor caches, and generated dependency directories.

## 7. Verify Contract Manifest When It Exists

The planned manifest verification command from T046 must recompute every listed SHA-256 digest and fail when a checked-out schema, contract doc, validation script, or fixture differs from the manifest.

The target signature verification profile is:

```text
sdp-trace-signature/sigstore-dsse-keyless-v1
```

Public environments should verify an in-toto Statement in a DSSE envelope with Sigstore/Cosign keyless identity and transparency-log evidence where available. Private or air-gapped environments may use an approved equivalent, but must record the identity policy and verification result explicitly.

The block is not product proof until self-trace and self-attestation validate. Schema-only examples prove contract shape; local private-equivalent signing proves envelope mechanics; production trust requires the proof state to name its actual trust anchor and immutable source reference.

## 8. Inspect Evidence Safety

For committed examples, confirm artifact references are safe to publish:

- no secrets, credentials, raw customer data, or private prompt contents
- sanitized summaries preserve useful evidence
- SHA-256 hashes are present when committed artifacts are referenced
- redaction notes explain withheld content
- unverified external links use `integrity_status: unverified` or equivalent

## 9. Inspect Pilot Evidence After Runs

Do not inspect pilot evidence as a product proof until `examples/self-trace/assessment-input.json` validates. The repository must first prove it can trace its own development.

Expected sanitized outputs after pilot execution:

```text
retired-research-artifacts/
examples/opencode/
examples/superpowers/
examples/jvm-bazel/
```

Raw local outputs may live under:

```text
.sdp-trace-runs/
```

That path is intentionally ignored by git.

## 10. Check Schema Compatibility

Every new schema and committed example that claims a schema contract must declare or reference a schema version. Breaking changes require updated examples, migration notes, and `sdp-gate` compatibility notes.
