# Socratic Review Packet

Generated: 2026-05-10
Branch: codex/product-roadmap-speckit


## FILE: AGENTS.md

```text
     1	# Agent Instructions
     2
     3	`sdp-trace` is the portable trust substrate for AI-assisted delivery.
     4
     5	## Purpose
     6	Define traceability, provenance, evidence, gate verdicts, and decision records that work across coding agents and harnesses.
     7
     8	## Boundary
     9	This repo must stay independent from `sdp_lab`, Beads, Operator Mode, and any specific harness runtime.
    10
    11	Allowed:
    12	- JSON schemas
    13	- Markdown docs
    14	- portable examples
    15	- tiny validation/rendering tools when needed
    16
    17	Not allowed:
    18	- dependency on SDP Operator Mode
    19	- dependency on Beads
    20	- dependency on agentloop
    21	- hidden assumptions about Claude, Codex, OpenCode, or GitHub
    22
    23	## Product Language
    24	Use SpecKit-aligned terms first: spec, plan, task, evidence, gate, decision, trace, provenance.
    25
    26	Avoid internal SDP terms unless a doc explicitly maps them.
    27
    28	## Quality Bar
    29	Every claim about a gate or verdict must be evidence-backed or marked `not_assessed`.
    30
    31	No opaque health scores.
    32
    33	## Engineering Stack
    34	Target product code is Go.
    35
    36	No Node.js, npm, JavaScript, TypeScript, or `.mjs` tooling is allowed in the active product path.
    37
    38	Bash is allowed only as a thin command launcher when Go would add no product value. Any Bash kept in the active path needs an explicit reason.
    39
    40	New Go code must be small, readable, testable, covered by focused tests, and free of TODO/FIXME markers.
    41	Put measurable complexity gates in CI or docs, not only in prose.
    42
    43	## Decomposition Rule
    44	If `AGENTS.md` exceeds 100 lines or any module needs more than 10 skills, the module is too large, under-decomposed, or overengineered.
    45
    46	## Trust Rules
    47	The repository has already failed once by letting prose, task checkboxes, and checked-in JSON overclaim what current verification could not replay. Do not repeat that failure.
    48
    49	- Machine proof wins over prose, task checkboxes, reports, review ledgers, and mirrors.
    50	- Dirty checkout output is local structural evidence only, not external trust.
    51	- Checked-in proof JSON is not authority unless live-verified or externally signed.
    52	- No deferred trust closure: missing external evidence keeps the block open.
    53	- Prose is not authoritative. Use `sdp-trace-claim` tags for claims.
    54	- Source-bound proof requires a clean immutable source commit; if a changed file is a manifest subject, commit it first, then regenerate release proof in a separate commit.
    55	- Do not close task checkboxes, review ledgers, or docs after source-bound proof without another source-bound cycle if those files are manifest subjects.
    56	- Keep mirrored self-trace data synchronized: `assessment-input.json` must mirror self-trace arrays, and hash references must match current files.
    57
    58	## Required Work Loop
    59	Every non-trivial implementation chunk needs:
    60
    61	1. SpecKit delta.
    62	2. Socratic spec review before approval handoff.
    63	3. Trace coverage when verifier behavior or trust claims change.
    64	4. Test-first behavior when behavior changes.
    65	5. Drift checks: spec vs implementation and regression against previous blocks.
    66	6. Live verifier state: `pass`, `fail`, `cannot_verify`, or `not_assessed`.
    67	7. Strict review with recorded disposition, fresh verification, and scoped commit.
    68
    69	If a chunk cannot be traced or verified yet, mark `not_assessed` or `cannot_verify` with a concrete reason and create a tracked follow-up before closing.
    70
    71	## Block Intake Protocol
    72	When the user says "take block into work", use `sdp-trace-trust-workflow`.
    73	First land current approved work through PR/review, then continue new block work in a fresh worktree.
    74	Prepare SpecKit deltas, run Socratic spec review, fix/record findings, then stop for explicit user approval of the reviewed spec direction before writing implementation code.
    75	After approval, split independent tasks to fast subagents with minimal context; each slice needs scoped commit, review, drift/regression check, and integration audit.
    76	Do not stop at implementation-only closure for block work unless the user explicitly asks to stop before PR or merge.
    77	Prepare and review the PR with code, tracing/evidence, and requirements-vs-implementation planes before ready/merge.
    78
    79	## Review Rules
    80	For adversarial pi review in this repo, prefer non-OpenAI, non-Anthropic, and non-Google models unless the user explicitly permits otherwise.
    81
    82	- Run separate code, tracing/evidence, and requirements-vs-implementation review planes for trust blocks; repeat them at PR level.
    83	- Verify review findings against full files before accepting or rejecting them.
    84	- Replace hung, empty, or off-task reviews; do not count them as evidence.
    85	- Record absent GitHub checks as CI `not_assessed`, not green.
    86	- Keep model selection, retry, fallback, and timeout details in `sdp-trace-trust-workflow` and `pi-review`, not in this root router.
    87
    88	## Claim Tags
    89	Use `docs/claim-authoring.md` for authoritative claim syntax.
    90
    91	Current Slice 1 validator intentionally accepts only narrow evidence forms. Do not introduce arbitrary `proof:*`, `state:*`, or `none` evidence unless cross-reference verification has been implemented.
    92
    93	## Commands
    94	Use current command contracts in `docs/agent-entrypoint.md` and `docs/reviewer-entrypoint.md`.
    95
    96	- Defaults: `go test ./...`, `jq empty schema/*.json`, `gofmt` for changed Go files, and `git diff --check`.
    97	- For schema/contract changes, also check refs, changed examples, fixture shape, and Go struct/schema alignment.
    98
    99	Bash verification commands are not product architecture. Keep them only when they are thin launchers around Go commands or external tools.
```

## FILE: specs/003-agent-supply-chain-roadmap/spec.md

```text
     1	# Feature Specification: Agent Supply Chain Roadmap
     2
     3	**Feature Branch**: `003-agent-supply-chain-roadmap`
     4	**Created**: 2026-05-10
     5	**Status**: Draft - roadmap artifact, Socratic review pending
     6	**Input**: Product discussion: target C-level buyers, especially CTOs, with
     7	`sdp-trace` as a neutral evidence layer for GitHub-first agentic delivery,
     8	OSS coding tools, OSS harnesses, and general-purpose agents crossing into
     9	software delivery.
    10
    11	## Product Boundary
    12
    13	`sdp-trace` records the agent supply chain of software delivery. It does not
    14	replace a coding agent, harness, Git host, CI system, SIEM, GRC tool, or release
    15	governance process.
    16
    17	The product must answer a CTO-level question quickly:
    18
    19	> Who or what initiated this software change, through which agent or harness
    20	> route, what evidence exists, what was independently witnessed, what was only
    21	> claimed by an agent, what remains `not_assessed` or `cannot_verify`, and which
    22	> human owns the next decision?
    23
    24	For engineers and reviewers, the value is not a separate wow moment. The value
    25	is honest work: fewer manual investigations, fewer unbacked "done" claims, and
    26	less confusion between agent prose, CI facts, and signed evidence.
    27
    28	## Scope
    29
    30	In scope:
    31
    32	- GitHub-first evidence packets for PR/MR-style changes.
    33	- OSS coding tools: `pi`, OpenCode, and GSD2.
    34	- OSS harness or methodology layers: GSD, Superpowers, and GSD2 where it acts as
    35	  both harness and agent.
    36	- General-purpose agents such as OpenClaw and Hermes only when they cross a
    37	  software delivery boundary: repository mutation, GitHub/GitLab/GitFlic
    38	  action, CI run, artifact, infrastructure config, release claim, or
    39	  secret-bearing automation.
    40	- Evidence theater detection as explicit facts and gap states.
    41	- Signed attestation as the top trust profile, not the first adoption step.
    42
    43	Out of scope:
    44
    45	- General monitoring of all personal-agent actions.
    46	- Employee surveillance, disciplinary decisions, or blame assignment.
    47	- Native merge, release, risk, compliance, or employment policy decisions.
    48	- Broad claims that a tool, model, or harness is "supported" from one observed
    49	  run.
    50	- Product dependency on GitHub, OpenCode, GSD, GSD2, Superpowers, `pi`,
    51	  OpenClaw, Hermes, Claude, Codex, or any specific provider.
    52
    53	## User Scenarios & Testing
    54
    55	### User Story 1 - CTO Reads A PR Evidence Packet (Priority: P0)
    56
    57	A CTO can open one evidence packet for a software change and understand the
    58	agent route, evidence strength, missing proof, and human decision owner without
    59	reading raw logs by hand.
    60
    61	**Why this priority**: C-level buyers do not buy JSON schemas. They buy a
    62	short path from agentic delivery uncertainty to decision-grade facts.
    63
    64	**Independent Test**: Given a GitHub PR with an agent-produced change, retained
    65	CI artifacts, and a selected `sdp-trace` profile, a reviewer can produce a
    66	packet that separates observed facts, agent claims, CI-witnessed evidence,
    67	review evidence, missing evidence, and next human decision owner.
    68
    69	**Acceptance Scenarios**:
    70
    71	1. **Given** a PR has commits, CI checks, and `sdp-trace` artifacts, **When**
    72	   the packet is generated, **Then** it records change host identity, PR id,
    73	   source/head refs, commit refs, CI witness state, retained artifact refs, and
    74	   `not_assessed` gaps.
    75	2. **Given** the agent claims tests passed but no CI or retained test artifact
    76	   exists, **When** the packet is generated, **Then** the test claim is recorded
    77	   as agent-claimed verification, not independent evidence.
    78	3. **Given** a human approval is required for the next decision, **When** the
    79	   packet is generated, **Then** the packet records the human decision owner or
    80	   marks ownership `cannot_verify` with a reason.
    81
    82	---
    83
    84	### User Story 2 - OpenCode/GSD Delivery Chain Is Observed (Priority: P0)
    85
    86	A pilot operator can observe a real OpenCode + GSD delivery loop and bind its
    87	raw session signals to a software change packet without hand-authoring proof
    88	events.
    89
    90	**Why this priority**: This is the current closest real dogfood path and tests
    91	whether `sdp-trace` adds value beside an OSS coding tool plus harness loop.
    92
    93	**Independent Test**: A real OpenCode/GSD session can be run or imported through
    94	a reviewed session profile; normalized events preserve model, harness,
    95	interaction, tool, mutation, test, and gap states without raw prompt or response
    96	retention.
    97
    98	**Acceptance Scenarios**:
    99
   100	1. **Given** OpenCode emits native JSONL, **When** `sdp-trace` collects it under
   101	   a reviewed profile, **Then** supported fields normalize into
   102	   `harness-event-v1` facts and unsupported fields remain `not_assessed`.
   103	2. **Given** GSD phase or task metadata is present, **When** the evidence packet
   104	   is generated, **Then** it records the phase/task reference as workflow intent,
   105	   not methodology compliance.
   106	3. **Given** private prompt, response, token, or path-like tool metadata exists
   107	   in the native stream, **When** collection runs, **Then** retained output stays
   108	   digest-only or sanitized, and unsafe fields fail before persistence.
   109
   110	---
   111
   112	### User Story 3 - Pi And GSD2 Session Import Are Assessed (Priority: P0)
   113
   114	A pilot operator can determine whether `pi` or GSD2 exposes stable session data
   115	that can be imported without product coupling.
   116
   117	**Why this priority**: `pi` is a minimal agent runtime and GSD2 is a standalone
   118	agent built on the Pi SDK. They may become the cleanest OSS path for portable
   119	agent-session evidence.
   120
   121	**Independent Test**: Discovery artifacts identify stable session/export
   122	surfaces, required fields, missing fields, redaction constraints, and whether
   123	the path is importable, wrapper-only, plugin-required, or `not_assessed`.
   124
   125	**Acceptance Scenarios**:
   126
   127	1. **Given** a `pi` session artifact exists, **When** discovery runs, **Then**
   128	   the artifact is classified as importable, partial, unsafe, unstable, or
   129	   `not_assessed` with evidence refs.
   130	2. **Given** GSD2 controls planning, execution, verification, git isolation, and
   131	   cost/token state, **When** discovery runs, **Then** `sdp-trace` records which
   132	   states can be imported as facts and which remain GSD2-internal claims.
   133	3. **Given** no stable export exists, **When** the roadmap is updated, **Then**
   134	   the row remains `not_assessed` and no support claim is made.
   135
   136	---
   137
   138	### User Story 4 - General-Purpose Agent Boundary Is Audited (Priority: P1)
   139
   140	A CTO or security leader can see when a general-purpose agent crossed into
   141	software delivery and what evidence binds that upstream actor to the resulting
   142	repo, PR, CI, or artifact action.
   143
   144	**Why this priority**: General-purpose agents are increasingly used by
   145	non-technical staff. The risk is not that they chat. The risk is that they
   146	touch code, CI, infrastructure, or release channels without a traceable
   147	software-delivery boundary.
   148
   149	**Independent Test**: A controlled Hermes or OpenClaw style task initiates a
   150	repository or GitHub action, and `sdp-trace` records the upstream channel,
   151	agent/session id where available, delegated coding tool/harness, change host
   152	action, retained evidence, and missing binding states.
   153
   154	**Acceptance Scenarios**:
   155
   156	1. **Given** a general-purpose agent initiates a GitHub/repo action, **When** a
   157	   change packet is generated, **Then** upstream actor/channel/session refs are
   158	   recorded when evidence exists, otherwise attribution remains `not_assessed`.
   159	2. **Given** a general-purpose agent delegates to a coding agent, **When**
   160	   downstream commits or PRs exist, **Then** the packet distinguishes upstream
   161	   intent, delegated execution, Git mutation, and CI witness facts.
   162	3. **Given** the general-purpose agent performs non-software actions only,
   163	   **When** `sdp-trace` evaluates scope, **Then** the action is out of product
   164	   scope and no general monitoring claim is made.
   165
   166	---
   167
   168	### User Story 5 - Signed Attestation Caps The Trust Ladder (Priority: P2)
   169
   170	A governance consumer can require signed evidence packages when local, CI, and
   171	customer witness evidence is not enough.
   172
   173	**Why this priority**: Signed attestation is the top trust profile. It should
   174	cap the ladder after evidence semantics are stable, not block day-one adoption.
   175
   176	**Independent Test**: An evidence packet can be bound to an in-toto/DSSE-style
   177	statement or approved private equivalent without converting local evidence into
   178	production trust by default.
   179
   180	**Acceptance Scenarios**:
   181
   182	1. **Given** an evidence packet has stable refs and digests, **When** signing is
   183	   requested, **Then** the signed statement binds packet digest, source refs,
   184	   witness refs, selected profile, signer identity, and freshness evidence.
   185	2. **Given** required signing evidence is absent, **When** a signed profile is
   186	   selected, **Then** signed-attestation state is `cannot_verify` and the
   187	   package does not claim trusted release.
   188	3. **Given** a customer uses private PKI instead of public Sigstore, **When**
   189	   the profile is configured, **Then** the private-equivalent evidence is
   190	   recorded explicitly and scoped to that customer policy.
   191
   192	## Evidence Theater Taxonomy
   193
   194	The roadmap must make these conditions machine-visible:
   195
   196	- **Agent-claimed verification**: agent says tests passed, but no independent
   197	  retained test evidence exists.
   198	- **Unbound intent**: change exists, but the source task, prompt, issue, or
   199	  approval boundary is missing.
   200	- **Actor laundering**: a general-purpose agent or harness delegates work, but
   201	  Git or PR metadata shows only a bot, shared account, or human committer.
   202	- **Review theater**: review exists, but reviewer independence, runner,
   203	  read-only state, model identity, or evidence retention is missing.
   204	- **CI theater**: checks are green but do not cover the changed risk or were not
   205	  retained as evidence for the selected claim.
   206	- **Artifact theater**: proof JSON or Markdown exists, but it is stale,
   207	  unreplayed, unsigned for the selected profile, or not bound to source.
   208	- **Human approval theater**: approval prose exists, but authority, role,
   209	  approval reference, or decision owner is missing.
   210	- **Scope theater**: one observed model/tool/path is generalized into broad
   211	  compatibility, support, readiness, or trust.
   212
   213	## Functional Requirements
   214
   215	- **FR-001**: `sdp-trace` MUST model the software-delivery agent supply chain as
   216	  facts across upstream initiator, agent runtime, harness/methodology, coding
   217	  tool, change host, CI, review, artifact, witness, and human decision owner.
   218	- **FR-002**: The first change-host adapter MUST be GitHub, but product concepts
   219	  MUST NOT be GitHub-specific.
   220	- **FR-003**: Change-host records MUST support future GitLab, GitFlic, Gitea,
   221	  Forgejo, and custom VCS/MR providers without changing evidence semantics.
   222	- **FR-004**: `sdp-trace` MUST distinguish post-hoc import, wrapper/sidecar
   223	  observation, and native plugin/hook evidence.
   224	- **FR-005**: Post-hoc import and wrapper/sidecar observation MUST be P0
   225	  adoption paths; native plugins/hooks MAY be P1+ only after discovery proves
   226	  value.
   227	- **FR-006**: Workflow layers such as GSD, Superpowers, and Oh My OpenAgent MUST
   228	  be recorded as intent, phase, role, task, or checkpoint facts unless separate
   229	  evidence proves compliance.
   230	- **FR-007**: GSD2 MUST be evaluated separately from GSD because it is a
   231	  standalone agent/runtime built on the Pi SDK, not only a harness layer.
   232	- **FR-008**: General-purpose agents MUST be in scope only when they cross a
   233	  software-delivery boundary.
   234	- **FR-009**: General-purpose agent monitoring outside software-delivery
   235	  boundaries MUST be explicitly out of scope.
   236	- **FR-010**: Evidence packet output MUST preserve `pass`, `fail`,
   237	  `not_assessed`, `cannot_verify`, `missing_telemetry`, `unsupported`, and
   238	  `not_integrated` states without collapsing them into a score.
   239	- **FR-011**: Agent claims MUST be recorded separately from CI-witnessed,
   240	  change-host-observed, harness-observed, gateway-observed, customer-witnessed,
   241	  or signed evidence.
   242	- **FR-012**: Signed attestation MUST be represented as a top trust profile, not
   243	  as a prerequisite for local or pilot adoption.
   244	- **FR-013**: Every roadmap item that names a tool MUST state the exact evidence
   245	  surface to inspect before claiming support.
   246	- **FR-014**: Missing stable export, missing credentials, absent CI artifacts,
   247	  unavailable session logs, or blocked API access MUST keep the relevant row
   248	  `not_assessed` or `cannot_verify`.
   249	- **FR-015**: Product docs MUST state that C-level buyers consume summaries and
   250	  risk facts, while engineers and reviewers consume traceable evidence rows.
   251
   252	## Success Criteria
   253
   254	- **SC-001**: A CTO can read one packet summary and identify the agent route,
   255	  independent evidence, missing proof, and human decision owner.
   256	- **SC-002**: A reviewer can explain why one OpenCode/GSD run is observed
   257	  without claiming broad OpenCode/GSD support.
   258	- **SC-003**: A `pi`/GSD2 discovery row cannot move out of `not_assessed` until
   259	  a stable session/export evidence path is inspected.
   260	- **SC-004**: A general-purpose agent can be recorded as upstream software
   261	  delivery initiator without turning `sdp-trace` into a general employee
   262	  monitoring product.
   263	- **SC-005**: Evidence theater cases are visible as facts, not hidden behind a
   264	  green score or prose summary.
   265	- **SC-006**: Signed attestation can cap the trust ladder while local and
   266	  CI-witnessed packets remain honestly scoped.
   267
   268	## Open Questions
   269
   270	1. Which exact CTO packet format creates the first product wow: PR comment,
   271	   downloadable archive, static HTML, Markdown report, or CLI summary?
   272	2. Which general-purpose agent should be the first boundary spike: Hermes or
   273	   OpenClaw?
   274	3. Should GSD2 discovery run before or after `pi` session import discovery, given
   275	   that GSD2 is built on the Pi SDK but adds stronger workflow semantics?
   276	4. What is the minimum acceptable signed-attestation profile for a customer that
   277	   cannot use public Sigstore?
   278	5. Which future non-GitHub change host matters first for product direction:
   279	   GitLab, GitFlic, Gitea/Forgejo, or Jenkins-only artifact flow?
```

## FILE: specs/003-agent-supply-chain-roadmap/plan.md

```text
     1	# Implementation Plan: Agent Supply Chain Roadmap
     2
     3	**Branch**: `003-agent-supply-chain-roadmap` | **Date**: 2026-05-10 | **Spec**: [spec.md](spec.md)
     4	**Input**: Product roadmap specification from
     5	`/specs/003-agent-supply-chain-roadmap/spec.md`
     6
     7	## Summary
     8
     9	Create a SpecKit-shaped product roadmap for `sdp-trace` as a neutral agent
    10	supply-chain evidence layer. The roadmap is discovery and product-planning work;
    11	it does not authorize schema, Go, CLI, or verifier implementation.
    12
    13	The first buyer is C-level, usually CTO. The first product path is GitHub-first
    14	but not GitHub-bound. The near-term integration set is:
    15
    16	- GitHub PR evidence packet;
    17	- OpenCode + GSD delivery chain observation;
    18	- `pi` session import discovery;
    19	- GSD2 discovery as a Pi-SDK-based standalone coding agent;
    20	- Superpowers/GSD workflow-intent mapping;
    21	- one general-purpose agent boundary spike with Hermes or OpenClaw;
    22	- signed attestation as the top trust profile after evidence packet semantics
    23	  stabilize.
    24
    25	## Technical Context
    26
    27	**Language/Version**: Markdown SpecKit artifacts only in this roadmap slice.
    28	Future implementation remains Go-first.
    29	**Primary Dependencies**: Existing `sdp-trace` docs, schemas, examples, and
    30	SpecKit conventions.
    31	**Storage**: Product roadmap artifacts under `specs/003-agent-supply-chain-roadmap/`.
    32	**Testing**: Markdown review, `git diff --check`, optional `go test ./...` for
    33	repo baseline.
    34	**Target Platform**: Portable evidence contracts that can later map to GitHub,
    35	GitLab, GitFlic, Gitea/Forgejo, Jenkins, local CLI agents, and signed witness
    36	systems.
    37	**Project Type**: SpecKit product roadmap, not implementation.
    38	**Constraints**: No Node.js, npm, JavaScript, TypeScript, or `.mjs` product-path
    39	changes. No dependency on GitHub, OpenCode, GSD, GSD2, Superpowers, `pi`,
    40	OpenClaw, Hermes, Claude, Codex, or any specific provider.
    41
    42	## Constitution Check
    43
    44	| Rule | Status | Evidence |
    45	|---|---|---|
    46	| Spec before implementation | Pass | This branch contains roadmap/spec artifacts only. |
    47	| Keep product independent | Pass | GitHub is first adapter, not product ontology. |
    48	| Evidence-backed claims only | Pass | Tool rows remain `not_assessed` until evidence surfaces are inspected. |
    49	| Preserve missing states | Pass | Roadmap keeps `not_assessed`, `cannot_verify`, `missing_telemetry`, `unsupported`, and `not_integrated`. |
    50	| No native policy verdicts | Pass | External consumers decide merge, release, compliance, HR, and risk outcomes. |
    51	| Go-first product path | Pass | No active product code or non-Go toolchain is added. |
    52
    53	## Project Structure
    54
    55	```text
    56	specs/003-agent-supply-chain-roadmap/
    57	|-- spec.md
    58	|-- plan.md
    59	|-- research.md
    60	`-- tasks.md
    61	```
    62
    63	Potential future implementation artifacts after separate approval:
    64
    65	```text
    66	schema/
    67	|-- change-host-event.schema.json
    68	|-- agent-supply-chain-record.schema.json
    69	|-- evidence-theater-finding.schema.json
    70	`-- cto-evidence-packet.schema.json
    71
    72	examples/
    73	|-- github-pr-evidence-packet/
    74	|-- opencode-gsd-supply-chain/
    75	|-- pi-session-import/
    76	|-- gsd2-session-import/
    77	`-- general-agent-boundary/
    78
    79	docs/
    80	|-- agent-supply-chain.md
    81	`-- cto-evidence-packet.md
    82	```
    83
    84	## Integration Strategy
    85
    86	| Integration | First evidence mode | Why | Risk |
    87	|---|---|---|---|
    88	| GitHub PR | Post-hoc import plus CI artifact refs | Fastest CTO-visible packet. | GitHub-specific concepts can leak into product model. |
    89	| OpenCode + GSD | Wrapper/sidecar plus native JSONL import | Already closest real dogfood chain. | One observed profile can be overclaimed as broad support. |
    90	| `pi` | Session import discovery | Minimal runtime may expose cleaner session evidence. | Stable export shape not yet assessed. |
    91	| GSD2 | Session import discovery plus wrapper feasibility | Combines tool and harness on Pi SDK. | Treating GSD2 like GSD would miss runtime-owned state. |
    92	| Superpowers | Workflow-intent mapping | Strong methodology/checkpoint surface. | Compliance should not be inferred from skill presence. |
    93	| Hermes/OpenClaw | Boundary spike only | Tests non-technical staff/general-agent risk. | Scope can sprawl into general agent monitoring. |
    94	| Signed attestation | Top trust profile | Governance capstone. | Premature signing can make weak evidence look stronger. |
    95
    96	## Roadmap Slices
    97
    98	### Slice P0-A: CTO Evidence Packet Shape
    99
   100	Define the packet contract, summary language, evidence rows, theater findings,
   101	and decision-owner fields. This slice can start as docs/examples before schema.
   102
   103	Exit criteria:
   104
   105	- One sample packet maps a PR to facts, claims, missing evidence, and next
   106	  decision owner.
   107	- Packet text does not claim merge, release, compliance, or production trust.
   108	- Every claim row has evidence refs or an explicit missing state.
   109
   110	### Slice P0-B: GitHub Change-Host Adapter Model
   111
   112	Define GitHub as the first change-host adapter without making GitHub the product
   113	ontology.
   114
   115	Exit criteria:
   116
   117	- GitHub concepts map to provider-neutral change-host fields.
   118	- Missing GitHub API access is `cannot_verify` or `not_assessed`.
   119	- Future GitLab/GitFlic/Gitea rows are named as planned adapters, not current
   120	  support.
   121
   122	### Slice P0-C: OpenCode + GSD Supply-Chain Packet
   123
   124	Use the existing OpenCode/GSD observation path as the first real software
   125	delivery chain.
   126
   127	Exit criteria:
   128
   129	- Native OpenCode/GSD events bind to a change packet without hand-authored proof.
   130	- GSD phase/task facts are workflow intent unless separately verified.
   131	- Missing mutation/test/PR/CI facts remain visible.
   132
   133	### Slice P0-D: Pi And GSD2 Discovery
   134
   135	Inspect `pi` and GSD2 session/export surfaces and classify import feasibility.
   136
   137	Exit criteria:
   138
   139	- Discovery reports identify stable artifacts, missing fields, and redaction
   140	  risks.
   141	- `pi` and GSD2 rows remain `not_assessed` until a real artifact is inspected.
   142	- GSD2 is treated as a standalone agent/runtime, not just GSD v1 methodology.
   143
   144	### Slice P1-A: Superpowers Workflow Intent
   145
   146	Map Superpowers skills/checkpoints to intent evidence only.
   147
   148	Exit criteria:
   149
   150	- Brainstorming, worktree, plan, TDD, review, and verification checkpoints can be
   151	  referenced when artifacts exist.
   152	- Presence of a skill does not prove compliance.
   153
   154	### Slice P1-B: General-Purpose Agent Boundary
   155
   156	Pick Hermes or OpenClaw for a single controlled software-delivery boundary
   157	spike.
   158
   159	Exit criteria:
   160
   161	- The spike records upstream channel/agent/session refs when available.
   162	- Delegation from general agent to coding agent is represented as a chain.
   163	- Non-software actions are explicitly out of scope.
   164
   165	### Slice P2-A: Signed Attestation Profile
   166
   167	Bind stable evidence packets to signed statements after packet semantics are
   168	reviewed.
   169
   170	Exit criteria:
   171
   172	- Signing binds packet digest, source refs, witness refs, selected profile,
   173	  identity, and freshness evidence.
   174	- Missing signed evidence blocks signed-trust claims with `cannot_verify`.
   175	- Customer private equivalent is recorded as scoped policy evidence, not as a
   176	  universal support claim.
   177
   178	## Review Gates
   179
   180	Before any implementation:
   181
   182	- Complete this roadmap package: `spec.md`, `plan.md`, `research.md`,
   183	  `tasks.md`.
   184	- Run Socratic/product review focused on C-level value, scope control, evidence
   185	  semantics, and integration order.
   186	- Resolve or block critical/major findings.
   187	- Ask for explicit approval of the reviewed roadmap direction.
   188
   189	Before any adapter implementation:
   190
   191	- Identify exact evidence surface for the selected tool.
   192	- Add fixture shape before parser behavior.
   193	- Define redaction/retention safety constraints.
   194	- Define what remains `not_assessed`.
   195
   196	Before any support claim:
   197
   198	- Run real tool/session evidence.
   199	- Retain safe artifacts.
   200	- Validate packet generation.
   201	- Record residual gaps.
   202
   203	## Non-Goals
   204
   205	- Do not build a dashboard in this roadmap slice.
   206	- Do not add schemas, Go code, or CLI commands in this roadmap slice.
   207	- Do not start native plugins before import/wrapper discovery proves value.
   208	- Do not monitor general-purpose agents outside software-delivery boundaries.
   209	- Do not turn signed attestation into a shortcut around weak evidence.
```

## FILE: specs/003-agent-supply-chain-roadmap/research.md

```text
     1	# Research: Agent Supply Chain Roadmap
     2
     3	## Decision: Reframe Around Agent Supply Chain
     4
     5	**Decision**: Position the roadmap around the agent supply chain of software
     6	delivery, not generic "AI observability" or "agent governance."
     7
     8	**Rationale**: The CTO-level question is not whether an agent existed. The
     9	question is whether a software change has a traceable route from upstream intent
    10	to repo mutation, review, CI, artifact, and human decision owner.
    11
    12	**Alternatives considered**:
    13
    14	- LLM observability framing: too broad and already served by tracing platforms.
    15	- AI governance/GRC framing: too high-level and risks becoming policy instead
    16	  of evidence.
    17	- Coding-agent replacement framing: wrong boundary; `sdp-trace` should observe
    18	  tools, not replace them.
    19
    20	## Decision: GitHub-First, Not GitHub-Bound
    21
    22	**Decision**: Use GitHub as the first change-host adapter because it is the
    23	fastest route to a C-level evidence packet, but keep product fields
    24	provider-neutral.
    25
    26	**Rationale**: GitHub now has first-party agent surfaces, third-party Claude and
    27	Codex agents, custom Copilot agents, and agentic workflows. That makes GitHub a
    28	good first proof surface and a competitive risk.
    29
    30	**Risk**: A GitHub-only model would miss upstream general-agent intent and would
    31	be painful to port to GitLab, GitFlic, Gitea/Forgejo, or Jenkins-centered flows.
    32
    33	**Implication**: Use `change_host`, `change_ref`, `review_ref`, `ci_ref`, and
    34	`artifact_ref` concepts instead of hard-coding GitHub terms in core semantics.
    35
    36	## Decision: Start With Import And Wrapper Modes
    37
    38	**Decision**: P0 evidence modes are post-hoc import and wrapper/sidecar
    39	observation. Native plugins/hooks are P1+.
    40
    41	**Rationale**: Native plugins produce richer evidence but are costly, fragile,
    42	and tool-specific. Post-hoc import and wrapper modes let us validate product
    43	value without coupling `sdp-trace` to every tool.
    44
    45	**Implication**:
    46
    47	- OpenCode/GSD: wrapper plus raw JSONL import.
    48	- `pi`: session import discovery first.
    49	- GSD2: session import discovery plus wrapper feasibility.
    50	- Superpowers: artifact/intent mapping.
    51	- Hermes/OpenClaw: boundary spike, not full plugin.
    52
    53	## Decision: Treat Harnesses As Intent Sources First
    54
    55	**Decision**: GSD, Superpowers, and similar methodology layers should be
    56	recorded as intent, phase, role, task, and checkpoint evidence unless separate
    57	observations prove compliance.
    58
    59	**Rationale**: A plan or skill invocation does not prove the agent followed the
    60	plan. Treating methodology presence as compliance would recreate evidence
    61	theater.
    62
    63	**Implication**: A packet may say "GSD phase declared" or "Superpowers
    64	verification checkpoint requested." It must not say "methodology complied"
    65	without additional evidence.
    66
    67	## Decision: Evaluate GSD2 Separately
    68
    69	**Decision**: GSD2 gets a separate discovery row from GSD v1.
    70
    71	**Rationale**: Public docs describe GSD2 as a standalone CLI coding agent built
    72	on the Pi SDK with direct runtime control, git isolation, context management,
    73	cost/token state, and crash recovery. That is a different evidence surface from
    74	GSD v1 as a harness/methodology layer.
    75
    76	**Implication**: GSD2 discovery should inspect runtime-owned state, not only
    77	generated plans or phase docs.
    78
    79	## Decision: Limit General-Purpose Agent Scope
    80
    81	**Decision**: General-purpose agents are in scope only at software-delivery
    82	boundaries.
    83
    84	**Rationale**: OpenClaw/Hermes-style agents are increasingly used through chat,
    85	gateway, cron, skills, and tools. They can be used by non-technical staff and
    86	can cross into code, repos, CI, and infrastructure. That creates a real C-level
    87	risk. But observing every personal-agent action would turn `sdp-trace` into a
    88	general monitoring product.
    89
    90	**Implication**: `sdp-trace` records only crossings into repositories, change
    91	hosts, CI, artifacts, infra config, release claims, or secret-bearing software
    92	automation.
    93
    94	## Decision: Signed Attestation Is The Top Profile
    95
    96	**Decision**: Signed attestation caps the trust ladder after packet semantics
    97	stabilize.
    98
    99	**Rationale**: Signing weak or incomplete evidence can make theater look
   100	official. The product should first make gaps explicit, then bind mature packets
   101	to DSSE/in-toto/Sigstore or customer private equivalents.
   102
   103	**Implication**: P0 packets may be local or CI-witnessed. P2 signed packets must
   104	record identity, policy, source refs, witness refs, packet digest, and freshness
   105	evidence.
   106
   107	## Integration Notes
   108
   109	These notes are discovery pointers, not support claims. Re-verify before using
   110	them in external-facing materials.
   111
   112	- GitHub third-party agents currently document Claude and Codex as supported
   113	  agents, with agent sessions creating PRs and requesting review.
   114	- GitHub custom agents are Copilot agent profiles, not arbitrary OSS-agent
   115	  runtime registration.
   116	- GitHub Agentic Workflows run agent-driven repository automation through
   117	  GitHub Actions and emphasize read-only-by-default workflow execution.
   118	- OpenCode exposes terminal-native AI coding behavior plus MCP, LSP, GitHub
   119	  Copilot experimental support, and self-hosted provider configuration.
   120	- `pi` is an agent harness mono repo with a coding-agent CLI, agent runtime,
   121	  unified LLM API, and explicit public session-sharing motivation.
   122	- GSD v1 is a meta-prompting/context/spec-driven system for multiple coding
   123	  harnesses.
   124	- GSD2 is a standalone coding agent built on the Pi SDK.
   125	- Superpowers is a multi-host skills/workflow methodology; useful as checkpoint
   126	  and intent evidence.
   127	- Oh My OpenAgent is a high-autonomy OpenCode harness layer with multi-agent
   128	  orchestration; useful as a future harness-intent source but too broad for P0.
   129	- OpenClaw and Hermes are general-purpose agents with gateways, channels,
   130	  tools, memory, skills, or scheduled automation; useful only for boundary
   131	  spikes into software delivery.
   132
   133	## External Sources
   134
   135	Re-check these before turning discovery notes into product or sales claims.
   136
   137	- GitHub custom agents:
   138	  https://docs.github.com/en/copilot/how-tos/copilot-sdk/use-copilot-sdk/custom-agents
   139	- GitHub Copilot third-party agents in VS Code:
   140	  https://code.visualstudio.com/docs/copilot/agents/third-party-agents
   141	- GitHub Agentic Workflows technical preview:
   142	  https://github.blog/changelog/2026-02-13-github-agentic-workflows-are-now-in-technical-preview
   143	- GitHub Agentic Workflows safe outputs:
   144	  https://github.github.com/gh-aw/reference/safe-outputs/
   145	- OpenCode:
   146	  https://github.com/opencode-ai/opencode
   147	- Pi:
   148	  https://github.com/badlogic/pi-mono
   149	- GSD:
   150	  https://github.com/gsd-build/get-shit-done
   151	- GSD2:
   152	  https://github.com/gsd-build/gsd-2
   153	- Superpowers:
   154	  https://github.com/obra/superpowers
   155	- Oh My OpenAgent:
   156	  https://ohmyopenagent.com/docs
   157	- OpenClaw:
   158	  https://github.com/openclaw/openclaw
   159	- Hermes Agent:
   160	  https://github.com/NousResearch/hermes-agent
   161
   162	## Research Gaps
   163
   164	- Exact `pi` local session storage/export format.
   165	- Exact GSD2 state database/session export format and redaction safety.
   166	- Whether Superpowers emits stable artifacts across Codex, OpenCode, Copilot
   167	  CLI, and Claude Code hosts.
   168	- Whether Hermes or OpenClaw has a stable event/session API suitable for a
   169	  safe boundary spike.
   170	- Minimal GitHub evidence packet shape that produces a CTO wow without needing
   171	  a dashboard.
   172	- Minimal signed private-equivalent profile for customers that cannot use public
   173	  Sigstore/Rekor.
```

## FILE: specs/003-agent-supply-chain-roadmap/tasks.md

```text
     1	# Tasks: Agent Supply Chain Roadmap
     2
     3	**Input**: Design documents from `/specs/003-agent-supply-chain-roadmap/`
     4	**Prerequisites**: `spec.md`, `plan.md`, `research.md`, Socratic review before
     5	implementation approval
     6	**Tests**: This slice is roadmap-only. Verification is Markdown sanity,
     7	`git diff --check`, and optional repo baseline `go test ./...`.
     8
     9	## Phase 0: Roadmap Draft
    10
    11	- [x] T001 Create isolated worktree for roadmap drafting.
    12	- [x] T002 Draft `spec.md` with CTO value, scope, user stories, functional
    13	  requirements, evidence theater taxonomy, and success criteria.
    14	- [x] T003 Draft `plan.md` with integration strategy, roadmap slices, review
    15	  gates, and non-goals.
    16	- [x] T004 Draft `research.md` with product decisions, integration notes, and
    17	  research gaps.
    18	- [ ] T005 Run Socratic/product review of roadmap package before turning any
    19	  item into implementation scope.
    20	- [ ] T006 Resolve or explicitly block critical/major roadmap review findings.
    21	- [ ] T007 Stop for explicit user approval of reviewed roadmap direction.
    22
    23	## Phase 1: CTO Evidence Packet Discovery
    24
    25	- [ ] T008 Define the minimum CTO packet summary shape for one GitHub PR.
    26	- [ ] T009 Define provider-neutral change-host fields that GitHub can populate.
    27	- [ ] T010 Define evidence rows for facts, agent claims, CI witness, review
    28	  evidence, missing evidence, and human decision owner.
    29	- [ ] T011 Define evidence theater finding rows and reason codes.
    30	- [ ] T012 Add one hand-reviewed packet example only after source artifacts are
    31	  identified; keep it marked example/discovery, not product proof.
    32
    33	## Phase 2: GitHub-First Adapter Specification
    34
    35	- [ ] T013 Map GitHub PR, issue, commit, check, review, Actions run, and artifact
    36	  concepts to provider-neutral change-host terms.
    37	- [ ] T014 Identify which GitHub API failures become `cannot_verify` vs
    38	  `not_assessed`.
    39	- [ ] T015 Define future adapter placeholders for GitLab, GitFlic,
    40	  Gitea/Forgejo, and Jenkins artifact-only flow without claiming support.
    41
    42	## Phase 3: OpenCode + GSD Supply-Chain Packet
    43
    44	- [ ] T016 Select one existing or new OpenCode/GSD run as the first packet
    45	  candidate.
    46	- [ ] T017 Map native OpenCode/GSD normalized events to packet fields.
    47	- [ ] T018 Record GSD phase/task metadata as workflow intent only.
    48	- [ ] T019 Identify residual missing evidence for mutation, test, PR, review,
    49	  CI, and signed witness states.
    50	- [ ] T020 Define what would be required before saying the OpenCode/GSD slice is
    51	  observed.
    52
    53	## Phase 4: Pi And GSD2 Discovery
    54
    55	- [ ] T021 Inspect `pi` local session storage/export shape.
    56	- [ ] T022 Classify `pi` evidence mode as importable, partial, wrapper-only,
    57	  plugin-required, unsafe, unstable, or `not_assessed`.
    58	- [ ] T023 Inspect GSD2 session/state surfaces, including runtime state, git
    59	  isolation, verification, crash recovery, and cost/token fields where available.
    60	- [ ] T024 Classify GSD2 evidence mode independently from GSD v1.
    61	- [ ] T025 Define redaction and retention constraints for both `pi` and GSD2
    62	  before any parser work.
    63
    64	## Phase 5: Superpowers Intent Mapping
    65
    66	- [ ] T026 Identify stable Superpowers artifacts or skill-invocation evidence
    67	  across target hosts.
    68	- [ ] T027 Map brainstorming, worktree, planning, TDD, review, and verification
    69	  checkpoints to intent facts.
    70	- [ ] T028 Define why skill presence does not prove methodology compliance.
    71
    72	## Phase 6: General-Purpose Agent Boundary Spike
    73
    74	- [ ] T029 Choose Hermes or OpenClaw for the first software-delivery boundary
    75	  spike.
    76	- [ ] T030 Define the boundary: channel/session -> general agent -> delegated
    77	  coding agent or direct repo action -> change host -> CI/artifact.
    78	- [ ] T031 Identify stable event/session/API evidence exposed by the selected
    79	  agent.
    80	- [ ] T032 Define out-of-scope non-software actions explicitly.
    81	- [ ] T033 Define privacy and employee-monitoring guardrails for C-level review.
    82
    83	## Phase 7: Signed Attestation Capstone
    84
    85	- [ ] T034 Define the minimum evidence packet fields that must exist before
    86	  signing can be meaningful.
    87	- [ ] T035 Define DSSE/in-toto/Sigstore target profile and customer private
    88	  equivalent profile.
    89	- [ ] T036 Define signing failure states, freshness evidence, and identity policy.
    90	- [ ] T037 Ensure signed attestation cannot upgrade missing evidence into trust.
    91
    92	## Phase 8: Review And Approval Gate
    93
    94	- [ ] T038 Run separate product-value, evidence-semantics, and integration-order
    95	  review planes.
    96	- [ ] T039 Record findings and dispositions in a roadmap review ledger.
    97	- [ ] T040 Update roadmap artifacts after accepted findings.
    98	- [ ] T041 Ask for explicit approval before creating implementation blocks.
```
