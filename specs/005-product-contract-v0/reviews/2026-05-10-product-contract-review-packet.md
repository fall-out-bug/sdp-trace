# Product Contract v0 Review Packet

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

## FILE: specs/005-product-contract-v0/spec.md

```text
     1	# Feature Specification: Product Contract v0
     2
     3	**Feature Branch**: `005-product-contract-v0`
     4	**Created**: 2026-05-10
     5	**Status**: Draft - product contract proposal, focused review pending
     6	**Input**: Repeated roadmap failure: `sdp-trace` has strong trust substrate
     7	work, but no reviewed buyer-facing product contract that blocks substrate-only
     8	backlog from being treated as P0 product progress.
     9
    10	## Product Boundary
    11
    12	Product Contract v0 defines the first buyer-facing output that `sdp-trace`
    13	must be able to produce from its provenance, evidence, trace, verifier, witness,
    14	and observation substrate.
    15
    16	It is not a dashboard, gate, legal contract, compliance certification,
    17	employee-monitoring system, signed-attestation implementation, or new agent
    18	runtime. It is a normative product contract for backlog intake:
    19
    20	> A feature counts as P0 product progress only when it fills, improves, or
    21	> verifies a required row of the Change Evidence Packet v0.
    22
    23	If a feature is useful but does not map to a packet row, it remains substrate,
    24	discovery, or future integration work. That work may still be valuable, but it
    25	must not be presented as closing the first product gap.
    26
    27	## Core Output
    28
    29	The first product output is **Change Evidence Packet v0**.
    30
    31	Canonical artifact:
    32
    33	- portable Markdown report for review, diff, and offline/on-prem use;
    34	- evidence bundle with retained refs, digests, and redaction status;
    35	- optional projections later: static HTML, PR/MR comment, CLI summary, PDF,
    36	  or signed attestation envelope.
    37
    38	The canonical artifact is Markdown plus evidence bundle because the first
    39	enterprise use case must work without GitHub, SaaS dashboards, public Sigstore,
    40	or external log export.
    41
    42	## Input Boundary
    43
    44	Product Contract v0 accepts available evidence, not an idealized complete
    45	delivery record.
    46
    47	Minimum input:
    48
    49	- change identity: repository or source set, commit range, branch, PR/MR id, or
    50	  local change id;
    51	- selected evidence profile: local, change-host-rich, harness-observed, or
    52	  signed;
    53	- available evidence refs: trace events, provenance records, evidence bundle,
    54	  CI artifacts, review records, authority evaluations, harness observations,
    55	  witness records, or external assertions;
    56	- redaction policy and retention status;
    57	- optional human decision-owner hints.
    58
    59	Missing inputs must remain `not_assessed` or `cannot_verify`; they must not be
    60	filled by prose, inference, or checked-in status files alone.
    61
    62	## Required Packet Rows
    63
    64	These rows are the first product contract. Every P0 feature must cite one or
    65	more row ids.
    66
    67	| row id | buyer question | required output |
    68	| --- | --- | --- |
    69	| `PC-CHANGE` | What software change is being assessed? | Change host, repo/source set, branch, PR/MR/local id, commit range, and source refs, or `cannot_verify`. |
    70	| `PC-INITIATOR` | Who or what initiated the change? | Human, upstream general agent, task, issue, prompt boundary, or `not_assessed`. |
    71	| `PC-AGENT-ROUTE` | Through which agent/tool/harness route did the change pass? | Agent runtime, coding tool, harness, adapter, and delegation chain, with source class for each hop. |
    72	| `PC-MUTATION` | What repository or artifact mutation was observed? | Changed paths, commits, generated artifacts, infra config, or package/release refs, with evidence refs. |
    73	| `PC-VERIFICATION` | What verification exists and who witnessed it? | Agent-claimed, harness-observed, CI-witnessed, reviewer-observed, customer-witnessed, signed, or missing evidence rows. |
    74	| `PC-REVIEW` | What review evidence exists? | Reviewer plane, reviewer identity/source class, retained result, independence state, and gaps. |
    75	| `PC-AUTHORITY` | Did observed actions stay within declared authority? | Authority evaluation state or `not_assessed`; no merge, blame, employment, or policy verdict. |
    76	| `PC-THEATER` | Which proof gaps could mislead a buyer? | Evidence theater findings using required reason codes and evidence refs. |
    77	| `PC-ATTESTATION` | Is this packet signed or externally witnessed? | Signed profile, signer/witness refs, freshness, or `not_assessed`/`cannot_verify`. |
    78	| `PC-DECISION` | Who owns the next human decision? | Named role or owner ref for merge, release, risk acceptance, security review, or `cannot_verify`. |
    79	| `PC-RESIDUAL-GAPS` | What remains unknown? | Explicit residual gaps with state, reason, and what evidence would close them. |
    80
    81	## Evidence States
    82
    83	Packet rows MUST preserve these states without collapsing them into a score:
    84
    85	- `pass`
    86	- `partial`
    87	- `fail`
    88	- `not_assessed`
    89	- `cannot_verify`
    90	- `missing_telemetry`
    91	- `unsupported`
    92	- `not_integrated`
    93
    94	Rows may also classify evidence source strength:
    95
    96	- `agent_claimed`
    97	- `harness_observed`
    98	- `change_host_observed`
    99	- `ci_witnessed`
   100	- `reviewer_observed`
   101	- `customer_witnessed`
   102	- `signed`
   103	- `external_assertion`
   104
   105	Source strength is not a trust score. It explains where the claim came from.
   106
   107	## Evidence Theater v0
   108
   109	Packet v0 MUST support these P0 theater findings:
   110
   111	| reason code | condition |
   112	| --- | --- |
   113	| `agent_claimed_verification` | Agent claims tests, checks, or completion, but no independent retained evidence exists. |
   114	| `unbound_intent` | Change exists, but task, issue, prompt boundary, approval, or initiating actor is not bound. |
   115	| `ci_theater` | CI is green or claimed, but the selected evidence profile lacks retained CI witness/artifact coverage for the claim. |
   116	| `scope_theater` | One observed tool/model/harness run is generalized into broad support, readiness, compatibility, or trust. |
   117
   118	These P1 theater findings are deferred unless a selected profile requires them:
   119
   120	- `actor_laundering`
   121	- `review_theater`
   122	- `artifact_theater`
   123	- `human_approval_theater`
   124
   125	Deferral must be visible in `PC-RESIDUAL-GAPS`, not hidden.
   126
   127	## Russian Enterprise Baseline
   128
   129	Product Contract v0 must work in a local or self-hosted enterprise environment.
   130
   131	The baseline profile MUST NOT require GitHub, public SaaS, public transparency
   132	logs, public Sigstore/Rekor, raw prompt export, or broad employee monitoring.
   133
   134	Baseline inputs may be limited to:
   135
   136	- local Git refs;
   137	- GitFlic, GitLab, Gitea/Forgejo, or custom change-host refs when available;
   138	- Jenkins, TeamCity, GitLab CI, GitHub Actions, or custom CI artifacts when
   139	  available;
   140	- local artifact store refs;
   141	- customer private PKI or internal witness records when available;
   142	- digest-only or redacted harness/session evidence.
   143
   144	If those sources are absent, the packet still has value when it makes missing
   145	evidence explicit and names the decision owner. It must not pretend the missing
   146	evidence exists.
   147
   148	## Backlog Gate
   149
   150	Any future P0 backlog item must include:
   151
   152	- `packet_rows`: one or more required row ids from this spec;
   153	- `evidence_surface`: the exact artifact, API, log, schema, fixture, or command
   154	  output to inspect;
   155	- `closure_state`: what moves the item from `not_assessed` to observed,
   156	  partial, unsupported, unsafe, or `cannot_verify`;
   157	- `buyer_effect`: what the CTO packet becomes clearer about after the feature;
   158	- `non_goal`: what the feature still does not prove.
   159
   160	If `packet_rows` is empty, the item cannot be P0 product progress.
   161
   162	## User Scenarios & Testing
   163
   164	### User Story 1 - CTO Reads One Packet (Priority: P0)
   165
   166	A CTO can read a Change Evidence Packet v0 for one software change in under two
   167	minutes and identify the change, agent route, witnessed evidence, agent claims,
   168	theater findings, residual gaps, and next decision owner.
   169
   170	**Why this priority**: The first product proof is not that substrate exists. The
   171	first product proof is that substrate becomes a buyer-readable artifact.
   172
   173	**Independent Test**: Given one synthetic or real change with partial evidence,
   174	the packet separates observed facts, claims, missing evidence, and decision
   175	ownership without requiring raw logs.
   176
   177	### User Story 2 - Product Owner Blocks Substrate-Only Backlog (Priority: P0)
   178
   179	A product owner can reject a proposed P0 feature when it does not cite a packet
   180	row and does not improve the buyer artifact.
   181
   182	**Why this priority**: This prevents the repeated failure where useful substrate
   183	features are treated as product contract progress.
   184
   185	**Independent Test**: Given a proposed integration item, the intake checklist
   186	classifies it as P0 product, substrate, discovery, or future integration based
   187	on packet row mapping.
   188
   189	### User Story 3 - Engineer Maps Existing Capabilities (Priority: P0)
   190
   191	An engineer can inspect `traceability.md` and see which existing `sdp-trace`
   192	capabilities already feed packet rows and which rows remain gaps.
   193
   194	**Why this priority**: The contract must not erase current work. It must show
   195	how current substrate becomes product output.
   196
   197	**Independent Test**: At least one current doc/schema/example path is mapped to
   198	each known row or marked as a gap with required evidence.
   199
   200	### User Story 4 - Russian Enterprise Pilot Works Without GitHub (Priority: P0)
   201
   202	A local/self-hosted enterprise pilot can produce a useful packet with local Git,
   203	self-hosted change-host refs, private CI/artifact refs, and digest-only evidence
   204	even when GitHub APIs and public signing are unavailable.
   205
   206	**Why this priority**: A Russian-market enterprise target cannot depend on
   207	GitHub-native product semantics.
   208
   209	**Independent Test**: The example packet includes a baseline profile path that
   210	does not require GitHub and preserves missing external evidence as gaps.
   211
   212	## Functional Requirements
   213
   214	- **FR-001**: Product Contract v0 MUST define Change Evidence Packet v0 as the
   215	  first buyer-facing output.
   216	- **FR-002**: The canonical packet MUST be portable Markdown plus an evidence
   217	  bundle; static HTML, PR/MR comments, CLI summaries, PDFs, and signed envelopes
   218	  are projections.
   219	- **FR-003**: Product Contract v0 MUST define required packet rows and stable
   220	  row ids.
   221	- **FR-004**: Every P0 feature MUST cite at least one required packet row.
   222	- **FR-005**: A feature that does not cite a packet row MUST NOT be labeled P0
   223	  product progress.
   224	- **FR-006**: Packet rows MUST distinguish observed facts, agent claims,
   225	  independent witness evidence, external assertions, and missing evidence.
   226	- **FR-007**: Packet rows MUST preserve `not_assessed` and `cannot_verify`.
   227	- **FR-008**: Packet v0 MUST include the four P0 evidence theater reason codes.
   228	- **FR-009**: Deferred P1 theater findings MUST be visible as residual gaps when
   229	  relevant.
   230	- **FR-010**: Product Contract v0 MUST include a Russian enterprise baseline
   231	  profile that does not require GitHub, SaaS dashboards, public Sigstore/Rekor,
   232	  raw prompt export, or broad employee monitoring.
   233	- **FR-011**: Tool, harness, model, or change-host support claims MUST NOT be
   234	  made until an evidence surface is inspected and mapped to packet rows.
   235	- **FR-012**: Signed attestation MUST be additive evidence over a packet, not a
   236	  shortcut that upgrades missing evidence into trust.
   237	- **FR-013**: Product Contract v0 MUST include a traceability matrix from current
   238	  substrate capabilities to packet rows.
   239	- **FR-014**: Product Contract v0 MUST include one example packet marked
   240	  example-only, not product proof.
   241	- **FR-015**: Roadmap items in `003-agent-supply-chain-roadmap` MUST be blocked
   242	  from implementation approval until they map to Product Contract v0 rows.
   243
   244	## Success Criteria
   245
   246	- **SC-001**: A CTO proxy can identify change, route, evidence, gaps, theater
   247	  findings, and decision owner from the example packet in under two minutes.
   248	- **SC-002**: A reviewer can reject a substrate-only feature as non-P0 using the
   249	  backlog gate.
   250	- **SC-003**: Existing `sdp-trace` capabilities are preserved as substrate
   251	  inputs and mapped to packet rows rather than dismissed.
   252	- **SC-004**: The baseline packet remains useful when GitHub, public signing,
   253	  and full agent session export are unavailable.
   254	- **SC-005**: A focused Socratic review finds no remaining critical or major
   255	  ambiguity in the contract gate before implementation approval is requested.
   256
   257	## Open Decisions With Proposed Defaults
   258
   259	| decision | proposed default | why |
   260	| --- | --- | --- |
   261	| First artifact | Markdown report plus evidence bundle | Works offline/on-prem, diffable, reviewable, projection-friendly. |
   262	| First projection | Static HTML after packet semantics stabilize | Good for demos, but not canonical. |
   263	| First change-host rich adapter | GitHub | Fastest rich evidence path, but not product ontology. |
   264	| Russian baseline | local Git plus optional GitFlic/GitLab/Jenkins/TeamCity refs | Prevents GitHub-native P0 lock-in. |
   265	| Signed attestation | P2 additive profile | Avoids signing weak evidence as theater. |
   266	| General-purpose agents | in scope only with software-delivery boundary evidence | Prevents employee-monitoring drift. |
```

## FILE: specs/005-product-contract-v0/plan.md

```text
     1	# Implementation Plan: Product Contract v0
     2
     3	**Branch**: `005-product-contract-v0` | **Date**: 2026-05-10 | **Spec**: [spec.md](spec.md)
     4	**Input**: Product contract reset after Socratic review found that the roadmap
     5	still lacked a concrete buyer-facing artifact and hard intake gate.
     6
     7	## Summary
     8
     9	Create a reviewed Product Contract v0 that makes Change Evidence Packet v0 the
    10	first product output of `sdp-trace` and turns packet row mapping into a hard
    11	backlog gate.
    12
    13	This is not implementation. It is the missing product contract layer above the
    14	existing substrate. The deliverable is a SpecKit package that answers:
    15
    16	- what the first product artifact is;
    17	- what rows it must contain;
    18	- how current substrate capabilities feed those rows;
    19	- how Russian enterprise/local deployment constraints affect P0;
    20	- how future features are accepted or rejected as P0 product work.
    21
    22	## Technical Context
    23
    24	**Language/Version**: Markdown SpecKit artifacts only.
    25	**Primary Dependencies**: Existing `sdp-trace` specs, docs, schemas, examples,
    26	and Socratic review findings from `003-agent-supply-chain-roadmap`.
    27	**Storage**: `specs/005-product-contract-v0/`.
    28	**Testing**: Markdown sanity, `git diff --check`, repo baseline `go test ./...`,
    29	and focused Socratic review before approval.
    30	**Target Platform**: Portable product contract usable across GitHub, GitFlic,
    31	GitLab, local Git, Jenkins/TeamCity, OpenCode/GSD, `pi`, GSD2, Superpowers, and
    32	future signed witness profiles.
    33	**Project Type**: Product contract and roadmap gate.
    34	**Constraints**: No Go, schema, CLI, verifier, dashboard, or adapter
    35	implementation in this slice.
    36
    37	## Constitution Check
    38
    39	| Rule | Status | Evidence |
    40	| --- | --- | --- |
    41	| Spec before implementation | Pass | This package is contract-only and blocks implementation until review. |
    42	| Keep product independent | Pass | Packet rows are provider-neutral; GitHub and GitFlic are sources, not ontology. |
    43	| Evidence-backed claims only | Pass | Support claims require inspected evidence surfaces. |
    44	| Preserve missing states | Pass | Packet rows preserve `not_assessed` and `cannot_verify`. |
    45	| No native policy verdicts | Pass | Packet names decision owners but does not approve merge, release, or compliance. |
    46	| Go-first product path | Pass | No product code or non-Go tooling is added. |
    47
    48	## Project Structure
    49
    50	```text
    51	specs/005-product-contract-v0/
    52	|-- spec.md
    53	|-- plan.md
    54	|-- example.md
    55	|-- traceability.md
    56	`-- tasks.md
    57	```
    58
    59	## What This Is
    60
    61	Product Contract v0 is an acceptance contract for product backlog, not a slogan.
    62
    63	It says:
    64
    65	1. The first buyer-facing artifact is Change Evidence Packet v0.
    66	2. Packet v0 has required row ids.
    67	3. Existing substrate capabilities must map to those rows.
    68	4. Future P0 work must cite the rows it improves.
    69	5. Work that does not cite rows is substrate/discovery/future integration, not
    70	   P0 product progress.
    71
    72	## How To Get There
    73
    74	### Step 1: Fix The Product Output
    75
    76	Choose one canonical output:
    77
    78	- Markdown report plus evidence bundle.
    79
    80	Record static HTML, PR/MR comments, CLI summaries, PDFs, and signed envelopes as
    81	projections. This prevents UI/demo format debates from blocking the contract.
    82
    83	### Step 2: Define Packet Rows
    84
    85	Define required rows in `spec.md`:
    86
    87	- `PC-CHANGE`
    88	- `PC-INITIATOR`
    89	- `PC-AGENT-ROUTE`
    90	- `PC-MUTATION`
    91	- `PC-VERIFICATION`
    92	- `PC-REVIEW`
    93	- `PC-AUTHORITY`
    94	- `PC-THEATER`
    95	- `PC-ATTESTATION`
    96	- `PC-DECISION`
    97	- `PC-RESIDUAL-GAPS`
    98
    99	These row ids become the backlog gate.
   100
   101	### Step 3: Write One Example Packet
   102
   103	Write `example.md` as a concrete packet with partial evidence and explicit gaps.
   104	It must be useful even when many fields are `not_assessed`.
   105
   106	The example is not product proof. It is a product contract example.
   107
   108	### Step 4: Map Current Substrate
   109
   110	Write `traceability.md` mapping current docs, schemas, examples, and internal
   111	packages to packet rows. This prevents throwing away current work while also
   112	showing which rows still lack buyer-facing output.
   113
   114	### Step 5: Block Roadmap Work That Does Not Map
   115
   116	Update the roadmap language so integration work is not P0 by itself. GitHub,
   117	GitFlic, OpenCode/GSD, `pi`, GSD2, Superpowers, and general-purpose agents are
   118	evidence sources for packet rows.
   119
   120	### Step 6: Run Focused Socratic Review
   121
   122	Review only the hard-gate question:
   123
   124	> Does Product Contract v0 actually prevent future substrate-only work from
   125	> being treated as P0 product progress?
   126
   127	Do not proceed to implementation until critical/major findings are fixed or
   128	explicitly blocked and the user approves the reviewed direction.
   129
   130	## Roadmap Reclassification
   131
   132	| item type | P0 product? | condition |
   133	| --- | --- | --- |
   134	| Packet row definition | Yes | It changes the required buyer artifact. |
   135	| Example packet | Yes | It proves the contract can be read and reviewed. |
   136	| Traceability matrix | Yes | It maps substrate to product output. |
   137	| GitHub adapter | Not by itself | P0 only if it fills named packet rows. |
   138	| GitFlic/local Git/Jenkins baseline | Not by itself | P0 only if it proves the packet works without GitHub. |
   139	| OpenCode/GSD import | Not by itself | P0 only if it fills agent route, mutation, or verification rows. |
   140	| `pi`/GSD2 discovery | Discovery | P0 only after evidence surfaces are mapped to packet rows. |
   141	| Signed attestation | P2 profile | It cannot replace missing packet evidence. |
   142	| Dashboard | Later projection | Not needed for Product Contract v0. |
   143
   144	## Review Gates
   145
   146	Before approval:
   147
   148	- `spec.md`, `plan.md`, `example.md`, `traceability.md`, and `tasks.md` exist.
   149	- `003-agent-supply-chain-roadmap` references this contract as the blocker for
   150	  implementation approval.
   151	- Focused Socratic review has usable output.
   152	- Critical/major findings are resolved or explicitly blocked.
   153	- User approves the reviewed Product Contract v0 direction.
   154
   155	Before implementation after approval:
   156
   157	- Implementation plan names exact packet rows.
   158	- Each feature has `packet_rows`, `evidence_surface`, `closure_state`,
   159	  `buyer_effect`, and `non_goal`.
   160	- Schema/Go/CLI work stays Go-first and evidence-preserving.
   161
   162	## Non-Goals
   163
   164	- Do not implement packet generation in this slice.
   165	- Do not add schemas or Go code in this slice.
   166	- Do not build a dashboard.
   167	- Do not select every future adapter.
   168	- Do not claim GitHub, GitFlic, OpenCode/GSD, `pi`, GSD2, or Superpowers support.
   169	- Do not turn signed attestation into day-one product proof.
```

## FILE: specs/005-product-contract-v0/example.md

```text
     1	# Change Evidence Packet v0 Example
     2
     3	Status: example-only, not product proof
     4	Contract: Product Contract v0
     5	Profile: Russian enterprise baseline plus optional rich change-host refs
     6
     7	This example shows the intended buyer artifact shape. It is not generated by
     8	current product code and must not be used as evidence that packet generation is
     9	implemented.
    10
    11	## Executive Summary
    12
    13	Change `PAY-1842` modified payment retry logic. The repository mutation is
    14	observable from Git refs. The agent route is partially observed from an
    15	OpenCode/GSD session run. Test verification is agent-claimed and partially
    16	harness-observed, but no retained CI artifact is available. Review ownership is
    17	not assessed. The next decision belongs to the service tech lead before merge
    18	or release.
    19
    20	This packet does not approve merge, release, compliance, production trust, or
    21	employee action. It records evidence, gaps, and decision ownership.
    22
    23	## Packet Metadata
    24
    25	| field | value |
    26	| --- | --- |
    27	| packet_id | `cep-2026-05-10-pay-1842-example` |
    28	| schema | `change-evidence-packet-v0` |
    29	| generated_from | example fixture |
    30	| selected_profile | `local-enterprise-baseline-v0` |
    31	| redaction_policy | digest-only prompts and model outputs |
    32	| packet_state | `not_assessed` for production trust |
    33
    34	## Required Rows
    35
    36	| row id | state | answer | evidence refs | gap / next evidence |
    37	| --- | --- | --- | --- | --- |
    38	| `PC-CHANGE` | `pass` | Local Git commit range `abc123..def456` changes `services/payments/retry.go` and tests. | `git:abc123..def456`, `artifact:diff-digest:sha256:example` | Change-host PR/MR id is `not_assessed`. |
    39	| `PC-INITIATOR` | `not_assessed` | Task id `PAY-1842` is named, but original issue or prompt boundary is not bound. | `external_assertion:task-id:PAY-1842` | Need issue ref, approved task packet, or prompt-boundary digest. |
    40	| `PC-AGENT-ROUTE` | `partial` | OpenCode/GSD session is observed. Model and harness events are digest-only. Upstream general agent is not assessed. | `harness-run:opencode-gsd-example`, `trace:session-command-digest` | Need upstream initiator binding if a general-purpose agent delegated the work. |
    41	| `PC-MUTATION` | `pass` | File mutation is observed in Git and correlated to harness tool/mutation events. | `git:diff`, `harness-event:mutation:sha256:example` | None for mutation existence. |
    42	| `PC-VERIFICATION` | `cannot_verify` | Agent claimed tests passed. Harness observed a verification command. No retained CI artifact is available. | `harness-event:test:sha256:example`, `agent_claim:test-pass` | Need retained CI artifact or customer witness evidence. |
    43	| `PC-REVIEW` | `not_assessed` | No independent review artifact is retained. | none | Need review record, reviewer plane, identity/source class, and retained result. |
    44	| `PC-AUTHORITY` | `not_assessed` | No selected authority envelope was supplied for this packet. | none | Need authority envelope and selected `policy_id`. |
    45	| `PC-THEATER` | `fail` | Two theater findings are present: `agent_claimed_verification`, `unbound_intent`. | `theater:agent_claimed_verification`, `theater:unbound_intent` | Need independent test witness and source intent binding. |
    46	| `PC-ATTESTATION` | `not_assessed` | No signed packet or customer private witness is present. | none | Need signed packet digest and signer/witness policy. |
    47	| `PC-DECISION` | `cannot_verify` | Expected owner is service tech lead, but no authority reference is bound. | `external_assertion:owner:service-tech-lead` | Need owner ref from change-host, task system, or policy. |
    48	| `PC-RESIDUAL-GAPS` | `pass` | Packet records missing intent, review, CI witness, authority, and signed evidence. | this packet | Gaps must remain visible in projections. |
    49
    50	## Theater Findings
    51
    52	| reason code | severity | finding | evidence | required closure evidence |
    53	| --- | --- | --- | --- | --- |
    54	| `agent_claimed_verification` | major | Agent claims tests passed, but no independent retained CI/customer witness exists. | `agent_claim:test-pass`, `harness-event:test:sha256:example` | Retained CI artifact, customer witness, or signed verification record. |
    55	| `unbound_intent` | major | Git change and task id exist, but source task/issue/prompt boundary is not bound. | `external_assertion:task-id:PAY-1842` | Issue/task ref or prompt-boundary digest. |
    56	| `ci_theater` | not_assessed | No CI state was available, so CI theater is not assessed rather than failed. | none | CI run ref and artifact coverage. |
    57	| `scope_theater` | pass | The packet does not claim broad OpenCode/GSD support from this run. | packet text | Keep support claims out of summary and projections. |
    58
    59	## Decision Owner
    60
    61	| decision | owner state | owner | why |
    62	| --- | --- | --- | --- |
    63	| merge readiness | `cannot_verify` | service tech lead | Owner is asserted but not bound to policy or change-host metadata. |
    64	| release readiness | `not_assessed` | none | Release decision is outside this packet. |
    65	| risk acceptance | `not_assessed` | none | No security/risk policy was selected. |
    66
    67	## Attachments
    68
    69	| ref | retained form | notes |
    70	| --- | --- | --- |
    71	| `git:abc123..def456` | source refs | Local Git evidence. |
    72	| `artifact:diff-digest:sha256:example` | digest | Diff body not retained in this example. |
    73	| `harness-run:opencode-gsd-example` | digest-only event bundle | Raw prompts and model outputs are not retained. |
    74	| `theater:agent_claimed_verification` | packet row | Derived from missing independent test witness. |
    75	| `theater:unbound_intent` | packet row | Derived from missing source intent binding. |
    76
    77	## What This Packet Does Not Prove
    78
    79	- It does not prove the change is safe.
    80	- It does not approve merge or release.
    81	- It does not prove broad OpenCode/GSD support.
    82	- It does not prove the named human owner has approved the change.
    83	- It does not prove signed attestation or production trust.
    84	- It does not monitor employee activity outside this software delivery change.
```

## FILE: specs/005-product-contract-v0/traceability.md

```text
     1	# Product Contract v0 Traceability Matrix
     2
     3	This matrix maps current `sdp-trace` substrate to Change Evidence Packet v0
     4	rows. It is intentionally conservative: a current capability may feed a packet
     5	row without proving the full row end to end.
     6
     7	## Packet Row Coverage
     8
     9	| packet row | current substrate inputs | current coverage | remaining gap |
    10	| --- | --- | --- | --- |
    11	| `PC-CHANGE` | `schema/delivery-trace-envelope.schema.json`, `internal/repoobserver`, `examples/block28-repo-observer/` | Partial | Need packet-level change identity contract for local Git, GitHub, GitFlic, GitLab, and artifact-only flows. |
    12	| `PC-INITIATOR` | `schema/decision-record.schema.json`, trace/provenance examples, harness session metadata | Partial | Need source task/issue/prompt-boundary binding row and safe redaction rule. |
    13	| `PC-AGENT-ROUTE` | `internal/harnessobs`, `examples/harness-observation/`, `schema/adapter-event.schema.json`, `schema/adapter-capture-run.schema.json` | Partial | Need route projection from raw substrate into buyer-readable chain. |
    14	| `PC-MUTATION` | `internal/harnessobs`, `internal/repoobserver`, authority observed actions, Git evidence examples | Partial | Need packet row that separates mutation existence from actor/tool/model attribution. |
    15	| `PC-VERIFICATION` | `schema/evidence-event.schema.json`, `schema/evidence-bundle.schema.json`, `internal/harnessobs`, `internal/ciartifact`, block17 and block26 examples | Partial | Need packet row that separates agent-claimed, harness-observed, CI-witnessed, and missing verification. |
    16	| `PC-REVIEW` | `internal/prreview`, `examples/pr-review/`, reviewer entrypoint docs | Partial | Need packet row for reviewer plane, independence state, retained result, and absent review state. |
    17	| `PC-AUTHORITY` | `schema/authority-envelope.schema.json`, `schema/authority-evaluation.schema.json`, `internal/authority`, `docs/authority-envelope.md` | Good substrate | Need packet projection that says authority state without making policy decisions. |
    18	| `PC-THEATER` | Existing negative fixtures, managed harness failures, adapter capture overclaim failures, review ledger discipline | Partial | Need first-class theater reason-code rows and derivation rules. |
    19	| `PC-ATTESTATION` | `schema/witness-profile-result.schema.json`, `internal/witness`, `internal/releaseproof`, `docs/contract-release-signing.md`, block15/16 examples | Good substrate | Need additive packet profile language and private/customer witness baseline. |
    20	| `PC-DECISION` | `schema/accountability.schema.json`, `docs/accountability-model.md`, decision records | Partial | Need buyer-facing owner row that does not become approval or blame. |
    21	| `PC-RESIDUAL-GAPS` | `not_assessed`/`cannot_verify` patterns across schemas and examples | Good substrate | Need packet-level gap summary grouped by decision relevance. |
    22
    23	## Existing Work Reclassification
    24
    25	| current area | product role under Contract v0 | not product progress until |
    26	| --- | --- | --- |
    27	| Flight recorder / trace substrate | Evidence source for multiple rows | It renders into packet rows. |
    28	| Evidence bundle schemas | Attachment and evidence-row source | Packet references the bundle and exposes missing states. |
    29	| Harness observation | Agent route, mutation, verification evidence source | It fills `PC-AGENT-ROUTE`, `PC-MUTATION`, or `PC-VERIFICATION`. |
    30	| Authority envelope | Authority row source | It is projected into `PC-AUTHORITY` without policy verdicts. |
    31	| CI artifact observation | Verification/witness source | It fills `PC-VERIFICATION` with retained CI witness refs. |
    32	| PR review profiles | Review row source | It fills `PC-REVIEW` with plane, result, and independence state. |
    33	| Witness and release proof | Attestation row source | It fills `PC-ATTESTATION` as additive evidence only. |
    34	| Adapter capture | Route/source capture | It fills route or mutation rows without broad support claims. |
    35	| Query packs | Investigation aid | A packet uses or links query results. |
    36	| Dashboard/report UI | Projection | Packet semantics are already stable. |
    37
    38	## Backlog Intake Template
    39
    40	Every P0 candidate must include this block:
    41
    42	```text
    43	packet_rows:
    44	  - PC-...
    45	evidence_surface:
    46	  - <artifact, API, schema, fixture, or command output>
    47	closure_state:
    48	  not_assessed -> <observed|partial|unsupported|unsafe|cannot_verify>
    49	buyer_effect:
    50	  <what becomes clearer in the packet>
    51	non_goal:
    52	  <what this still does not prove>
    53	```
    54
    55	If the candidate cannot fill `packet_rows`, it is not P0 product progress.
    56
    57	## Known Gaps
    58
    59	| gap | blocks | closure evidence |
    60	| --- | --- | --- |
    61	| No generated packet command | Implementation after approval | Go implementation that renders packet from retained inputs. |
    62	| No packet schema | Implementation after contract approval | `change-evidence-packet.schema.json` or equivalent Go model. |
    63	| No GitFlic/local enterprise fixture | Russian baseline confidence | Example fixture using local Git plus self-hosted/change-host refs. |
    64	| Theater reason-code derivation is not implemented | `PC-THEATER` automation | Documented rules and tests for P0 theater codes. |
    65	| Decision owner row is not bound | `PC-DECISION` confidence | Policy, task, or change-host owner ref with missing-state handling. |
    66	| Static HTML projection absent | Demo polish | Projection generated from canonical packet without changing semantics. |
```

## FILE: specs/005-product-contract-v0/tasks.md

```text
     1	# Tasks: Product Contract v0
     2
     3	**Input**: Design documents from `/specs/005-product-contract-v0/`
     4	**Prerequisites**: `003-agent-supply-chain-roadmap` Socratic review findings
     5	**Tests**: Contract-only slice. Verification is Markdown sanity,
     6	`git diff --check`, repo baseline `go test ./...`, and focused Socratic review
     7	before approval.
     8
     9	## Phase 0: Contract Draft
    10
    11	- [x] T001 Create `specs/005-product-contract-v0/`.
    12	- [x] T002 Draft `spec.md` defining Product Contract v0, Change Evidence Packet
    13	  v0, required packet rows, evidence states, theater v0, Russian enterprise
    14	  baseline, and backlog gate.
    15	- [x] T003 Draft `plan.md` explaining what the contract is and how to get from
    16	  current substrate to a reviewed product gate.
    17	- [x] T004 Draft `example.md` with one concrete example packet marked
    18	  example-only, not product proof.
    19	- [x] T005 Draft `traceability.md` mapping current substrate capabilities to
    20	  packet rows and known gaps.
    21
    22	## Phase 1: Roadmap Linkage
    23
    24	- [x] T006 Update `003-agent-supply-chain-roadmap` so P0 integration work is
    25	  blocked by Product Contract v0 approval.
    26	- [x] T007 Reclassify GitHub, GitFlic/local Git, OpenCode/GSD, `pi`, GSD2,
    27	  Superpowers, and general-purpose agent work as evidence sources for packet
    28	  rows, not standalone P0 product outcomes.
    29	- [x] T008 Add Product Contract v0 to the review gate before any implementation
    30	  scope approval.
    31
    32	## Phase 2: Focused Review
    33
    34	- [ ] T009 Build focused review packet for `005-product-contract-v0`.
    35	- [ ] T010 Run Socratic pi-review focused on whether the contract creates a real
    36	  hard gate against substrate-only P0 work.
    37	- [ ] T011 Record review findings and dispositions.
    38	- [ ] T012 Fix or explicitly block all critical/major findings.
    39	- [ ] T013 Stop for explicit user approval of reviewed Product Contract v0.
    40
    41	## Phase 3: Future Implementation Planning After Approval
    42
    43	- [ ] T014 Define `change-evidence-packet-v0` schema or Go model.
    44	- [ ] T015 Define packet renderer inputs and safe redaction behavior.
    45	- [ ] T016 Add baseline local/self-hosted fixture.
    46	- [ ] T017 Add rich GitHub fixture only as one evidence source.
    47	- [ ] T018 Add focused theater reason-code derivation tests.
    48	- [ ] T019 Add CLI or command surface only after schema/model and fixtures are
    49	  reviewed.
    50
    51	## Completion Rule
    52
    53	This slice is not complete until Product Contract v0 is reviewed and approved.
    54	Implementation tasks T014-T019 are placeholders for later planning, not
    55	authorization to write code.
```

## FILE: specs/003-agent-supply-chain-roadmap/reviews/2026-05-10-pi-socratic-review.md

```text
     1	# Pi Socratic Review: Agent Supply Chain Roadmap
     2
     3	Date: 2026-05-10
     4	Scope: `spec.md`, `plan.md`, `research.md`, `tasks.md`, and repository rules.
     5	Status: review completed; critical/major findings unresolved.
     6
     7	## Review Runs
     8
     9	| plane | model | status | raw output |
    10	| --- | --- | --- | --- |
    11	| SpecKit consistency and approval gate | `zai/glm-5.1` | usable | `raw/2026-05-10-glm-speckit.md` |
    12	| Trust and evidence semantics | `minimax/MiniMax-M2.7` | usable, hidden reasoning preamble stripped from stored artifact | `raw/2026-05-10-minimax-trust.md` |
    13	| CTO value and Russian enterprise adoption | `openrouter/qwen/qwen3.6-plus` | usable | `raw/2026-05-10-qwen-cto-market.md` |
    14	| Integration sequencing and DX feasibility | `openrouter/deepseek/deepseek-v4-pro` | usable | `raw/2026-05-10-deepseek-integration.md` |
    15	| Kimi slot from allowlist | `kimi-coding/k2p6` | `not_assessed` | current `pi` model registry did not match this pattern |
    16
    17	## Verdict
    18
    19	`REVISE_BEFORE_USER_REVIEW`
    20
    21	The direction is viable, but it is not ready for explicit roadmap approval.
    22	All reviewers converged on the same core issue: the roadmap has good evidence
    23	semantics, but the first CTO-facing packet remains too abstract. The approval
    24	gate should not ask the owner to approve implementation scope until the first
    25	packet shape, theater detection contract, local-market adapter posture, and
    26	discovery closure criteria are sharper.
    27
    28	## Consolidated Findings
    29
    30	| id | severity | plane | finding | disposition | evidence |
    31	| --- | --- | --- | --- | --- | --- |
    32	| SR-001 | critical | product value | CTO packet format and artifact are undefined. | `unresolved_blocker` | `spec.md` asks which packet format creates the first wow, but leaves PR comment, archive, static HTML, Markdown report, and CLI summary open (`spec.md:270`). `plan.md` says P0-A must produce one sample packet, but not its concrete surface (`plan.md:98`). |
    33	| SR-002 | critical | evidence semantics | Evidence theater is a taxonomy, not yet a binding detection/reporting contract. | `unresolved_blocker` | `spec.md` names eight theater conditions (`spec.md:192`) and scope says detection is in scope (`spec.md:40`), but FRs and tasks do not define minimum detection rows or closure criteria. |
    34	| SR-003 | major | Russian enterprise adoption | GitHub-first is too dominant for a Russian-market enterprise target. | `unresolved_blocker` | Scope starts with GitHub-first packets (`spec.md:32`), FR-003 only names GitFlic as future capability (`spec.md:220`), and Phase 2 keeps non-GitHub providers as placeholders (`tasks.md:39`). |
    35	| SR-004 | major | risk management | Integration risks are named without mitigations. | `unresolved_blocker` | The integration table lists risk for every target but has no mitigation column (`plan.md:86`). |
    36	| SR-005 | major | traceability | Open questions and research gaps are not mapped to tasks, owners, or closure criteria. | `unresolved_blocker` | Open questions remain in `spec.md:268`; research gaps remain in `research.md:162`; tasks list discovery items but no gap-to-task-to-evidence mapping (`tasks.md:53`). |
    37	| SR-006 | major | discovery method | `pi`, GSD2, Superpowers, Hermes/OpenClaw discovery methods are underspecified. | `unresolved_blocker` | Tasks say to inspect surfaces (`tasks.md:55`, `tasks.md:58`, `tasks.md:66`, `tasks.md:78`) but do not state whether inspection is runtime observation, source inspection, docs review, API probing, or fixture capture. |
    38	| SR-007 | major | scope control | General-purpose agent boundary needs an enforceable software-delivery boundary. | `unresolved_blocker` | Scope excludes broad monitoring (`spec.md:45`) and FR-009 excludes general monitoring (`spec.md:234`), but no minimum evidence condition distinguishes software-delivery boundary from general employee-agent activity. |
    39	| SR-008 | major | trust semantics | Signed attestation "top trust profile" needs operational meaning before enterprise discussion. | `unresolved_blocker` | Signed attestation is deferred (`spec.md:41`, `plan.md:165`), but open question 4 and Phase 7 leave minimum private-equivalent profile unresolved (`spec.md:276`, `tasks.md:85`). |
    40	| SR-009 | major | product/review evidence | Review evidence itself should be recorded with a lightweight manifest, but full `assessment-input.json` is not required for this roadmap draft. | `accepted_narrower` | MiniMax overreached by implying full self-trace mirror mechanics are required for roadmap prose. The valid narrower point is that this review package needs a review manifest with files, models, commands, and `not_assessed` reviewer states. This file and the raw outputs provide that starting point. |
    41	| SR-010 | minor | source quality | External source entries lack `last_checked` metadata. | `deferred_not_assessed` | `research.md:133` lists source URLs but no check dates. Useful before external-facing materials, not a blocker for internal roadmap review. |
    42
    43	## Top Socratic Questions For The Owner
    44
    45	1. What is the first CTO packet surface: PR comment, static HTML, Markdown,
    46	   archive, CLI summary, or another artifact?
    47	2. What is the minimum packet content that creates buyer value even when many
    48	   rows are `not_assessed`?
    49	3. Should Russian-market P0 include GitFlic/local Git/Jenkins-style artifact
    50	   flow beside GitHub, or is GitHub-only acceptable for the pilot?
    51	4. Which evidence theater findings are P0 rows in the first packet, and which
    52	   are explicitly deferred?
    53	5. What exact evidence closes a `pi` or GSD2 discovery row from `not_assessed`
    54	   to importable, partial, wrapper-only, unsafe, or unstable?
    55	6. What minimum technical boundary keeps general-purpose agent tracking inside
    56	   software delivery and out of employee surveillance?
    57	7. Is signed attestation additive evidence over an already meaningful packet,
    58	   or a separate enterprise profile that may be required by some buyers from
    59	   day one?
    60
    61	## Required Before Approval
    62
    63	Before asking for explicit roadmap approval:
    64
    65	1. Resolve the CTO packet format and add a first example shape.
    66	2. Bind evidence theater taxonomy to minimum P0 packet rows or explicitly defer
    67	   specific theater categories.
    68	3. Add GitFlic/local Git/Jenkins-style Russian-market discovery posture, or
    69	   explicitly state why it is not P0.
    70	4. Add mitigations to the integration risk table.
    71	5. Add gap-to-task-to-evidence closure mapping for open questions and research
    72	   gaps.
    73	6. Define discovery methods for `pi`, GSD2, Superpowers, and the selected
    74	   general-purpose agent spike.
    75	7. Define software-delivery boundary minimum evidence conditions.
    76
    77	## Evidence Commands
    78
    79	- `pi --no-tools --no-context-files --no-session --model zai/glm-5.1 --thinking high ...`
    80	- `pi --no-tools --no-context-files --no-session --model minimax/MiniMax-M2.7 --thinking high ...`
    81	- `pi --no-tools --no-context-files --no-session --model openrouter/qwen/qwen3.6-plus --thinking high ...`
    82	- `pi --no-tools --no-context-files --no-session --model openrouter/deepseek/deepseek-v4-pro --thinking high ...`
    83
    84	All four completed with usable outputs. Each run also emitted a startup warning
    85	that `kimi-coding/k2p6` from local settings did not match the current model
    86	registry; this affects only the unused Kimi slot.
```
