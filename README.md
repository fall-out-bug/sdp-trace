# sdp-trace

`sdp-trace` is a portable evidence substrate for AI-assisted delivery.

It records what happened, which evidence exists, what is missing, where the
evidence came from, and which human owns the next decision. It does not decide
whether a team may merge, release, accept risk, or claim production trust.

Current product status: controlled-pilot ready for repo-observable evidence,
local trace packages, assessment profiles, witness artifacts, and source-bound
release checks. It is not a broad production trust product and it does not
claim universal harness, model, CI, or air-gapped compatibility.

## Start Here

1. Read [Core Concepts](docs/concepts.md) to understand the contract:
   spec, plan, task, evidence, gate, decision, trace, and provenance.
2. Run the local smoke path:

   ```text
   go test ./...
   go run ./cmd/sdp-trace --help
   go run ./cmd/sdp-trace wrap --name smoke -- /bin/echo ok
   go run ./cmd/sdp-trace verify <run-dir>
   ```

3. Use [Agent Entrypoint](docs/agent-entrypoint.md) for the authoritative
   command and state contract.
4. Use [Reviewer Entrypoint](docs/reviewer-entrypoint.md) for a five-minute
   verification path and overclaim checklist.
5. Use [Documentation Map](docs/README.md) to choose the right next document.

## What It Produces

- trace run artifacts under `.sdp-trace-runs/`;
- report packages under `.sdp-trace-report/`;
- query and forensic query-pack outputs;
- assessment results for supported profiles;
- advisory gate facts for downstream policy consumers;
- CI or customer witness artifacts when required evidence exists;
- source-bound local release proof when manifest subjects match the source
  commit.

Every output is scoped. A local trace does not become CI evidence. A CI witness
does not become production trust. A checked-in JSON file is an audit artifact
until the current verifier replays it or an accepted external signature binds it.

## What It Does Not Do

`sdp-trace` does not:

- replace a harness, agent, ticket tracker, code review process, or CI system;
- detect every unwrapped local agent run;
- decide pass/fail, readiness, degradation, release approval, or risk override;
- convert missing evidence into success;
- depend on `sdp_lab`, Beads, Operator Mode, agentloop, Claude, Codex, OpenCode,
  GitHub, or any specific harness runtime.

Policy decisions belong to `sdp-gate`, CI, release governance, or another
external consumer. `sdp-trace` records evidence and gaps for those consumers.

## Recommended Reading Order

The docs have one primary path. They should not require readers to classify
themselves before understanding the repository.

1. [Documentation Map](docs/README.md)
2. [Core Concepts](docs/concepts.md)
3. [Agent Entrypoint](docs/agent-entrypoint.md)
4. [Reviewer Entrypoint](docs/reviewer-entrypoint.md)
5. [Harness Integration](docs/harness-integration.md)
6. [Schema Reference](schema/README.md)

Governance and rollout documents are supporting references after the core path
is clear.

## Repository Layout

- `cmd/` and `internal/`: Go CLI and verifier implementation.
- `schema/`: portable JSON schema contracts.
- `docs/`: product, command, governance, reviewer, and integration documentation.
- `examples/`: sanitized fixtures and pilot evidence packages.
- `specs/`: SpecKit working artifacts and implementation block records.

## Minimal Flow

```text
spec -> plan -> task -> change -> evidence -> provenance -> accountability -> assessment input
```

External policy consumers can turn assessment input and verifier facts into
decisions. When evidence is missing, the state must remain `not_assessed`,
`cannot_verify`, `missing_telemetry`, or an explicit failure reason.
