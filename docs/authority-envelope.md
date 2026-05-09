# Authority Envelope

Authority envelope assessment records whether an observed action stayed within the authority selected for a task. It is a fact surface, not a native gate.

Use it when a trace package needs to show:

- what task and `policy_id` were selected;
- what actor role and action types were declared;
- what action was observed from git, harness, CI, gateway, manual import, or external assertion evidence;
- whether action existence, actor attribution, tool attribution, and model attribution are independently verified, `not_assessed`, or `cannot_verify`;
- which evidence refs support the result.

## Boundary

`sdp-trace` does not decide whether `outside_authority` contaminates a demo, blocks a merge, triggers discipline, or changes release readiness. Those decisions belong to external policy consumers.

Missing selected policy or missing applicable rules produce `not_assessed`. Malformed envelopes, conflicting allow/deny rules, stale evidence, inaccessible evidence, unsafe evidence refs, and failed bindings produce `cannot_verify`.

Git-only evidence can prove a path mutation and commit binding. It cannot prove tool, harness, model, or human attribution by itself. Gateway evidence can prove that a model request existed, but it cannot prove mutation causality without a verified harness/action binding.

## Command

```bash
go run ./cmd/sdp-trace assess \
  --profile authority-envelope \
  --authority-package examples/authority-envelope-basic/outside-authority-denied-target/authority-package.json \
  --out authority-evaluation.json
```

Explain a result:

```bash
go run ./cmd/sdp-trace assess explain --assessment-result authority-evaluation.json
```

Preview required inputs without emitting an evaluation:

```bash
go run ./cmd/sdp-trace assess preview \
  --profile authority-envelope \
  --authority-package examples/authority-envelope-basic/outside-authority-denied-target/authority-package.json
```

## Package Shape

An authority package contains:

- `selected_policy_id`
- `actors`
- `authority_envelopes`
- `observed_actions`
- optional `evidence_bindings`
- optional `evidence_resolution`

The selected policy is caller supplied. If multiple envelopes exist, `sdp-trace` does not choose the newest, strictest, or broadest envelope.

Evidence references use safe URI-style strings:

- `file:<relative-path>`
- `artifact:<artifact-id>#<path-or-json-pointer>`
- `git:<commit-sha>#<path>`
- `external:<opaque-id>`

Required unresolved `external:` refs remain `cannot_verify`.

## Fixture Coverage

Fixtures live under `examples/authority-envelope-basic/` and cover:

- valid `within_authority`
- `outside_authority` for a denied target
- git-only attribution staying `not_assessed`
- gateway request without mutation model attribution
- malformed/conflicting policy as `cannot_verify`
- failed binding as `cannot_verify`
- stale evidence as `cannot_verify`
- feedback without inferred mutation
- one observed action evaluated under two selected `policy_id` values
