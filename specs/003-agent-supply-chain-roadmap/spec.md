# Feature Specification: Agent Supply Chain Roadmap

**Feature Branch**: `003-agent-supply-chain-roadmap`
**Created**: 2026-05-10
**Status**: Draft - roadmap artifact, Socratic review completed, revisions pending
**Input**: Product discussion: target C-level buyers, especially CTOs, with
`sdp-trace` as a neutral evidence layer for GitHub-first agentic delivery,
OSS coding tools, OSS harnesses, and general-purpose agents crossing into
software delivery.

## Product Boundary

`sdp-trace` records the agent supply chain of software delivery. It does not
replace a coding agent, harness, Git host, CI system, SIEM, GRC tool, or release
governance process.

The product must answer a CTO-level question quickly:

> Who or what initiated this software change, through which agent or harness
> route, what evidence exists, what was independently witnessed, what was only
> claimed by an agent, what remains `not_assessed` or `cannot_verify`, and which
> human owns the next decision?

For engineers and reviewers, the value is not a separate wow moment. The value
is honest work: fewer manual investigations, fewer unbacked "done" claims, and
less confusion between agent prose, CI facts, and signed evidence.

## Scope

In scope:

- GitHub-first evidence packets for PR/MR-style changes.
- OSS coding tools: `pi`, OpenCode, and GSD2.
- OSS harness or methodology layers: GSD, Superpowers, and GSD2 where it acts as
  both harness and agent.
- General-purpose agents such as OpenClaw and Hermes only when they cross a
  software delivery boundary: repository mutation, GitHub/GitLab/GitFlic
  action, CI run, artifact, infrastructure config, release claim, or
  secret-bearing automation.
- Evidence theater detection as explicit facts and gap states.
- Signed attestation as the top trust profile, not the first adoption step.

Out of scope:

- General monitoring of all personal-agent actions.
- Employee surveillance, disciplinary decisions, or blame assignment.
- Native merge, release, risk, compliance, or employment policy decisions.
- Broad claims that a tool, model, or harness is "supported" from one observed
  run.
- Product dependency on GitHub, OpenCode, GSD, GSD2, Superpowers, `pi`,
  OpenClaw, Hermes, Claude, Codex, or any specific provider.

## User Scenarios & Testing

### User Story 1 - CTO Reads A PR Evidence Packet (Priority: P0)

A CTO can open one evidence packet for a software change and understand the
agent route, evidence strength, missing proof, and human decision owner without
reading raw logs by hand.

**Why this priority**: C-level buyers do not buy JSON schemas. They buy a
short path from agentic delivery uncertainty to decision-grade facts.

**Independent Test**: Given a GitHub PR with an agent-produced change, retained
CI artifacts, and a selected `sdp-trace` profile, a reviewer can produce a
packet that separates observed facts, agent claims, CI-witnessed evidence,
review evidence, missing evidence, and next human decision owner.

**Acceptance Scenarios**:

1. **Given** a PR has commits, CI checks, and `sdp-trace` artifacts, **When**
   the packet is generated, **Then** it records change host identity, PR id,
   source/head refs, commit refs, CI witness state, retained artifact refs, and
   `not_assessed` gaps.
2. **Given** the agent claims tests passed but no CI or retained test artifact
   exists, **When** the packet is generated, **Then** the test claim is recorded
   as agent-claimed verification, not independent evidence.
3. **Given** a human approval is required for the next decision, **When** the
   packet is generated, **Then** the packet records the human decision owner or
   marks ownership `cannot_verify` with a reason.

---

### User Story 2 - OpenCode/GSD Delivery Chain Is Observed (Priority: P0)

A pilot operator can observe a real OpenCode + GSD delivery loop and bind its
raw session signals to a software change packet without hand-authoring proof
events.

**Why this priority**: This is the current closest real dogfood path and tests
whether `sdp-trace` adds value beside an OSS coding tool plus harness loop.

**Independent Test**: A real OpenCode/GSD session can be run or imported through
a reviewed session profile; normalized events preserve model, harness,
interaction, tool, mutation, test, and gap states without raw prompt or response
retention.

**Acceptance Scenarios**:

1. **Given** OpenCode emits native JSONL, **When** `sdp-trace` collects it under
   a reviewed profile, **Then** supported fields normalize into
   `harness-event-v1` facts and unsupported fields remain `not_assessed`.
2. **Given** GSD phase or task metadata is present, **When** the evidence packet
   is generated, **Then** it records the phase/task reference as workflow intent,
   not methodology compliance.
3. **Given** private prompt, response, token, or path-like tool metadata exists
   in the native stream, **When** collection runs, **Then** retained output stays
   digest-only or sanitized, and unsafe fields fail before persistence.

---

### User Story 3 - Pi And GSD2 Session Import Are Assessed (Priority: P0)

A pilot operator can determine whether `pi` or GSD2 exposes stable session data
that can be imported without product coupling.

**Why this priority**: `pi` is a minimal agent runtime and GSD2 is a standalone
agent built on the Pi SDK. They may become the cleanest OSS path for portable
agent-session evidence.

**Independent Test**: Discovery artifacts identify stable session/export
surfaces, required fields, missing fields, redaction constraints, and whether
the path is importable, wrapper-only, plugin-required, or `not_assessed`.

**Acceptance Scenarios**:

1. **Given** a `pi` session artifact exists, **When** discovery runs, **Then**
   the artifact is classified as importable, partial, unsafe, unstable, or
   `not_assessed` with evidence refs.
2. **Given** GSD2 controls planning, execution, verification, git isolation, and
   cost/token state, **When** discovery runs, **Then** `sdp-trace` records which
   states can be imported as facts and which remain GSD2-internal claims.
3. **Given** no stable export exists, **When** the roadmap is updated, **Then**
   the row remains `not_assessed` and no support claim is made.

---

### User Story 4 - General-Purpose Agent Boundary Is Audited (Priority: P1)

A CTO or security leader can see when a general-purpose agent crossed into
software delivery and what evidence binds that upstream actor to the resulting
repo, PR, CI, or artifact action.

**Why this priority**: General-purpose agents are increasingly used by
non-technical staff. The risk is not that they chat. The risk is that they
touch code, CI, infrastructure, or release channels without a traceable
software-delivery boundary.

**Independent Test**: A controlled Hermes or OpenClaw style task initiates a
repository or GitHub action, and `sdp-trace` records the upstream channel,
agent/session id where available, delegated coding tool/harness, change host
action, retained evidence, and missing binding states.

**Acceptance Scenarios**:

1. **Given** a general-purpose agent initiates a GitHub/repo action, **When** a
   change packet is generated, **Then** upstream actor/channel/session refs are
   recorded when evidence exists, otherwise attribution remains `not_assessed`.
2. **Given** a general-purpose agent delegates to a coding agent, **When**
   downstream commits or PRs exist, **Then** the packet distinguishes upstream
   intent, delegated execution, Git mutation, and CI witness facts.
3. **Given** the general-purpose agent performs non-software actions only,
   **When** `sdp-trace` evaluates scope, **Then** the action is out of product
   scope and no general monitoring claim is made.

---

### User Story 5 - Signed Attestation Caps The Trust Ladder (Priority: P2)

A governance consumer can require signed evidence packages when local, CI, and
customer witness evidence is not enough.

**Why this priority**: Signed attestation is the top trust profile. It should
cap the ladder after evidence semantics are stable, not block day-one adoption.

**Independent Test**: An evidence packet can be bound to an in-toto/DSSE-style
statement or approved private equivalent without converting local evidence into
production trust by default.

**Acceptance Scenarios**:

1. **Given** an evidence packet has stable refs and digests, **When** signing is
   requested, **Then** the signed statement binds packet digest, source refs,
   witness refs, selected profile, signer identity, and freshness evidence.
2. **Given** required signing evidence is absent, **When** a signed profile is
   selected, **Then** signed-attestation state is `cannot_verify` and the
   package does not claim trusted release.
3. **Given** a customer uses private PKI instead of public Sigstore, **When**
   the profile is configured, **Then** the private-equivalent evidence is
   recorded explicitly and scoped to that customer policy.

## Evidence Theater Taxonomy

The roadmap must make these conditions machine-visible:

- **Agent-claimed verification**: agent says tests passed, but no independent
  retained test evidence exists.
- **Unbound intent**: change exists, but the source task, prompt, issue, or
  approval boundary is missing.
- **Actor laundering**: a general-purpose agent or harness delegates work, but
  Git or PR metadata shows only a bot, shared account, or human committer.
- **Review theater**: review exists, but reviewer independence, runner,
  read-only state, model identity, or evidence retention is missing.
- **CI theater**: checks are green but do not cover the changed risk or were not
  retained as evidence for the selected claim.
- **Artifact theater**: proof JSON or Markdown exists, but it is stale,
  unreplayed, unsigned for the selected profile, or not bound to source.
- **Human approval theater**: approval prose exists, but authority, role,
  approval reference, or decision owner is missing.
- **Scope theater**: one observed model/tool/path is generalized into broad
  compatibility, support, readiness, or trust.

## Functional Requirements

- **FR-001**: `sdp-trace` MUST model the software-delivery agent supply chain as
  facts across upstream initiator, agent runtime, harness/methodology, coding
  tool, change host, CI, review, artifact, witness, and human decision owner.
- **FR-002**: The first change-host adapter MUST be GitHub, but product concepts
  MUST NOT be GitHub-specific.
- **FR-003**: Change-host records MUST support future GitLab, GitFlic, Gitea,
  Forgejo, and custom VCS/MR providers without changing evidence semantics.
- **FR-004**: `sdp-trace` MUST distinguish post-hoc import, wrapper/sidecar
  observation, and native plugin/hook evidence.
- **FR-005**: Post-hoc import and wrapper/sidecar observation MUST be P0
  adoption paths; native plugins/hooks MAY be P1+ only after discovery proves
  value.
- **FR-006**: Workflow layers such as GSD, Superpowers, and Oh My OpenAgent MUST
  be recorded as intent, phase, role, task, or checkpoint facts unless separate
  evidence proves compliance.
- **FR-007**: GSD2 MUST be evaluated separately from GSD because it is a
  standalone agent/runtime built on the Pi SDK, not only a harness layer.
- **FR-008**: General-purpose agents MUST be in scope only when they cross a
  software-delivery boundary.
- **FR-009**: General-purpose agent monitoring outside software-delivery
  boundaries MUST be explicitly out of scope.
- **FR-010**: Evidence packet output MUST preserve `pass`, `fail`,
  `not_assessed`, `cannot_verify`, `missing_telemetry`, `unsupported`, and
  `not_integrated` states without collapsing them into a score.
- **FR-011**: Agent claims MUST be recorded separately from CI-witnessed,
  change-host-observed, harness-observed, gateway-observed, customer-witnessed,
  or signed evidence.
- **FR-012**: Signed attestation MUST be represented as a top trust profile, not
  as a prerequisite for local or pilot adoption.
- **FR-013**: Every roadmap item that names a tool MUST state the exact evidence
  surface to inspect before claiming support.
- **FR-014**: Missing stable export, missing credentials, absent CI artifacts,
  unavailable session logs, or blocked API access MUST keep the relevant row
  `not_assessed` or `cannot_verify`.
- **FR-015**: Product docs MUST state that C-level buyers consume summaries and
  risk facts, while engineers and reviewers consume traceable evidence rows.

## Success Criteria

- **SC-001**: A CTO can read one packet summary and identify the agent route,
  independent evidence, missing proof, and human decision owner.
- **SC-002**: A reviewer can explain why one OpenCode/GSD run is observed
  without claiming broad OpenCode/GSD support.
- **SC-003**: A `pi`/GSD2 discovery row cannot move out of `not_assessed` until
  a stable session/export evidence path is inspected.
- **SC-004**: A general-purpose agent can be recorded as upstream software
  delivery initiator without turning `sdp-trace` into a general employee
  monitoring product.
- **SC-005**: Evidence theater cases are visible as facts, not hidden behind a
  green score or prose summary.
- **SC-006**: Signed attestation can cap the trust ladder while local and
  CI-witnessed packets remain honestly scoped.

## Open Questions

1. Which exact CTO packet format creates the first product wow: PR comment,
   downloadable archive, static HTML, Markdown report, or CLI summary?
2. Which general-purpose agent should be the first boundary spike: Hermes or
   OpenClaw?
3. Should GSD2 discovery run before or after `pi` session import discovery, given
   that GSD2 is built on the Pi SDK but adds stronger workflow semantics?
4. What is the minimum acceptable signed-attestation profile for a customer that
   cannot use public Sigstore?
5. Which future non-GitHub change host matters first for product direction:
   GitLab, GitFlic, Gitea/Forgejo, or Jenkins-only artifact flow?
