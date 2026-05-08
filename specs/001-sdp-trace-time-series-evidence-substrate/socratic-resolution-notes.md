# Socratic Resolution Notes

Date: 2026-04-30
Source review: prior internal spec critique (archived outside this repository)

## Scope

Three clean-context critics (`zai/glm-5.1`, `minimax/MiniMax-M2.7`, `kimi-coding/k2p6`) reviewed this SpecKit package. This note records the author-side resolutions for the converged blockers before a judge pass.

## Resolutions

### R1: technical executive Degradation Question Without Native Verdict

Status: resolved

Resolution: `sdp-trace` answers the technical executive question with movement data, not a yes/no degradation verdict. The spec now states that native movement records include prior/current values, deltas, units, dimensions, evidence coverage, and `not_assessed` gaps. Interpretation labels such as degrading, improving, pass, fail, ready, blocked, or not ready remain external verdicts.

Changed artifacts:

- `spec.md`
- `plan.md`
- `data-model.md`
- `quickstart.md`
- `contracts/external-policy-consumer-boundary.md`

### R2: Policy-Adjacent Fields and External Verdicts

Status: resolved

Resolution: `sdp-trace` does not assign evidence strength. Source-provided strength, quality, verdict, score, readiness, degradation, pass/fail, or override values are recorded as external assertions or external verdict inputs with explicit producer, origin, policy reference when available, artifact reference, and provenance.

Changed artifacts:

- `spec.md`
- `data-model.md`
- `research.md`
- `tasks.md`
- `contracts/external-policy-consumer-boundary.md`

### R3: JSON Schema Draft and Validator Timing

Status: resolved

Resolution: New schemas target JSON Schema Draft 2020-12. The original `ajv-cli@5.0.0` strategy was superseded on 2026-05-01 after final review found a dependency audit risk, and later retired from the active product path. Current active validation is Go-first: `go test ./...`, `jq empty schema/*.json`, `go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc`, and `git diff --check`. Earlier AJV/script validation records are historical evidence, not the current command contract.

Changed artifacts:

- `plan.md`
- `tasks.md`
- `research.md`
- `quickstart.md`
- `schema/README.md`

### R4: Evidence Reference Security and Integrity

Status: resolved

Resolution: Committed evidence references must be safe to publish. They must not contain secrets, credentials, raw customer data, or private prompt contents. Sanitized examples preserve summaries, hashes, redaction notes, and `integrity_status`. SHA-256 digests provide content continuity but are not authentication signatures; signing and write authorization remain external unless a future schema version adds them.

Changed artifacts:

- `spec.md`
- `plan.md`
- `data-model.md`
- `research.md`
- `quickstart.md`
- `tasks.md`

### R5: Schema Versioning and Backward Compatibility

Status: resolved

Resolution: New schemas use stable `$id` values and semver schema versions. Additive optional fields are minor changes; required field removals, enum semantic changes, or ownership-boundary changes are major changes. `schema/trace.schema.json` remains a compatibility path until a replacement path and migration note are committed. external policy consumers must declare supported schema versions.

Changed artifacts:

- `spec.md`
- `plan.md`
- `data-model.md`
- `research.md`
- `quickstart.md`
- `contracts/external-policy-consumer-boundary.md`
- `schema/README.md`

## Deferred

No policy gate is added here. This review improves the SpecKit contract only. Enforcement remains a later schema/tooling task, and pass/fail policy remains owned by external policy consumer.
