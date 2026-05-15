# sdp-trace

`sdp-trace` is a portable evidence substrate for AI-assisted delivery.

It records what happened, which evidence exists, what is missing, where the
evidence came from, and which human owns the next decision. It does not decide
whether a team may merge, release, accept risk, or claim production trust.

Current product status: controlled-pilot MVP. You can try local trace packages,
repo-observable evidence, assessment profiles, witness artifacts, and
source-bound release checks. You cannot use this repo as a production trust
authority, release approval system, universal harness adapter, or guarantee that
every unwrapped agent run was detected.

## Start Here

1. Read [Install](docs/install.md) and choose either a release binary or
   source checkout command path.
2. Read [Core Concepts](docs/concepts.md) to understand the contract:
   spec, plan, task, evidence, gate, decision, trace, and provenance.
3. Give [Agent Onboarding](docs/agent-onboarding.md) to any coding agent before
   it works in this repository.
4. Follow [Contributor Quick Start](docs/contributor-quickstart.md) to run the
   canonical local smoke path and verify your environment.
5. Use [Agent Entrypoint](docs/agent-entrypoint.md) for the authoritative
   command and state contract.
6. Use [Reviewer Entrypoint](docs/reviewer-entrypoint.md) for a five-minute
   verification path and overclaim checklist.
7. Use [Documentation Map](docs/README.md) to choose the right next document.

Origin note: `sdp-trace` was extracted from delivery evidence work in
`sdp_lab`. That history is not a runtime dependency and should not be required
context for using this repository.

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

## CTO View

Use `sdp-trace` when the organization needs delivery evidence that survives a
specific agent, harness, or workflow choice. The product gives engineering and
governance teams a portable record of:

- what spec, plan, and task a change came from;
- which evidence exists and which evidence is still missing;
- where the evidence was produced: local, CI, customer authority, or another
  stated scope;
- which gate facts can be handed to policy consumers;
- which human owns the next decision when evidence is incomplete.

The control objective is not "AI says this is safe." The control objective is
"a reviewer or policy system can see the evidence boundary, replay supported
checks, and refuse overclaims." That makes `sdp-trace` useful for CTO-level
questions about agent adoption, release risk, auditability, and vendor
portability without turning this repository into the release authority.

## Developer Experience

The happy path should stay ordinary:

1. Keep your normal spec, plan, task, code review, and CI flow.
2. Wrap or adapt the observable parts of the work.
3. Retain evidence artifacts and explicit gaps.
4. Hand assessment input to the policy consumer that already owns the decision.

Developers should not need to learn an internal SDP runtime, switch agents, or
accept hidden GitHub assumptions. If a workflow cannot yet produce replayable
evidence, record `not_assessed` or `cannot_verify` and keep the policy decision
outside `sdp-trace`.

## What It Does Not Do

`sdp-trace` does not:

- replace a harness, agent, ticket tracker, code review process, or CI system;
- detect every unwrapped local agent run;
- decide pass/fail, readiness, degradation, release approval, or risk override;
- convert missing evidence into success;
- depend on Beads, Operator Mode, agentloop, Claude, Codex, OpenCode, GitHub,
  or any specific harness runtime.

Policy decisions belong to CI, release governance, customer governance, or
another external policy consumer. `sdp-trace` records evidence and gaps for
those consumers.

## Recommended Reading Order

The docs have one primary path. They should not require readers to classify
themselves before understanding the repository.

1. [Install](docs/install.md)
2. [Core Concepts](docs/concepts.md)
3. [Agent Onboarding](docs/agent-onboarding.md)
4. [Contributor Quick Start](docs/contributor-quickstart.md)
5. [Agent Entrypoint](docs/agent-entrypoint.md)
6. [Reviewer Entrypoint](docs/reviewer-entrypoint.md)
7. [Harness Integration](docs/harness-integration.md)
8. [Schema Reference](schema/README.md)

Governance and rollout documents are supporting references after the core path
is clear.

## Repository Layout

- `cmd/` and `internal/`: Go CLI and verifier implementation.
- `schema/`: portable JSON schema contracts.
- `docs/`: product, command, governance, reviewer, and integration documentation.
- `examples/`: sanitized fixtures and pilot evidence packages.
- `specs/`: working specs and implementation block records. Current repository
  records are SpecKit-shaped, but `sdp-trace` can map evidence from other
  planning flows.

## Minimal Flow

```text
spec -> plan -> task -> change -> evidence -> provenance -> accountability -> assessment input
```

External policy consumers can turn assessment input and verifier facts into
decisions. When evidence is missing, the state must remain `not_assessed`,
`cannot_verify`, `missing_telemetry`, or an explicit failure reason.

## Adoption Shape

Start with one repository and one narrow gate. A good first adoption is usually
"retain local and CI evidence for this class of changes" rather than "approve
all AI work." Expand only after reviewers can replay the evidence boundary and
the team has decided who owns merge, release, risk override, and escalation
decisions.
