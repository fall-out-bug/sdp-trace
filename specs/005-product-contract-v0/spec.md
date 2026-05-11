# Feature Specification: Product Contract v0

**Feature Branch**: `005-product-contract-v0`
**Created**: 2026-05-10
**Status**: Draft - revised after full review, re-review pending
**Input**: Repeated roadmap failure: `sdp-trace` has strong trust substrate
work, but no reviewed buyer-facing product contract that classifies whether
backlog work is product progress or substrate/discovery work.

## Product Boundary

Product Contract v0 defines the first buyer-facing output that `sdp-trace`
must be able to produce from its provenance, evidence, trace, verifier, witness,
and observation substrate.

It is not a dashboard, legal contract, compliance certification,
employee-monitoring system, signed-attestation implementation, approval system,
or new agent runtime. It is a normative product contract for backlog
classification:

> A feature counts as P0 product progress only when it fills, improves, or
> verifies a required row of the Change Evidence Packet v0.

If a feature is useful but does not map to a packet row and does not make a
measurable forward change in that row, it remains substrate, discovery, or
future integration work. That work may still be valuable, but it must not be
presented as closing the first product gap.

## Core Output

The first product output is **Change Evidence Packet v0**.

Canonical artifact:

- portable Markdown report for review, diff, and offline/on-prem use;
- evidence bundle directory with a manifest of retained refs, digests, source
  classes, and redaction status;
- optional projections later: static HTML, PR/MR comment, CLI summary, PDF,
  or signed attestation envelope.

The canonical artifact is Markdown plus evidence bundle because the first
enterprise use case must work without GitHub, SaaS dashboards, public Sigstore,
public transparency logs, raw prompt export, or external log export.

## Packet Markdown Template

Every Change Evidence Packet v0 Markdown report MUST contain these sections in
this order:

1. `Executive Summary`: 5-8 buyer-readable sentences covering change, route,
   evidence, gaps, and next decision owner.
2. `Packet Metadata`: table with `packet_id`, `schema`, `generated_from`,
   `selected_profile`, `redaction_policy`, `bundle_ref`, and `packet_state`.
3. `Required Rows`: table with columns `row id`, `state`, `answer`,
   `evidence refs`, and `gap / next evidence`.
4. `Theater Findings`: table of triggered P0 findings with columns
   `reason code`, `state`, `severity`, `finding`, `trigger evidence`, and
   `required closure evidence`.
5. `Decision Ownership`: table with one row each for merge, release, risk
   acceptance, and security review.
6. `Evidence Bundle`: table linking packet evidence refs to bundle manifest
   entries.
7. `What This Packet Does Not Prove`: explicit non-goals and residual trust
   limits.

Generated packets may add appendices, but projections MUST NOT omit, merge, or
rename the required sections.

When no P0 theater finding is triggered, the `Theater Findings` section remains
present with no finding rows; the `PC-THEATER` required row records that the
assessment ran and found no triggered P0 findings.

`packet_state` is packet lifecycle metadata, not a trust, health, readiness,
approval, or confidence score. Allowed values are:

- `draft`
- `review_ready`
- `reviewed`
- `superseded`

Production trust remains decomposed into packet rows and residual gaps.

## Evidence Bundle Manifest

The evidence bundle is a directory with a required `manifest.json`. Product
Contract v0 does not implement the manifest, but later implementation must use
this minimum shape:

| field | meaning |
| --- | --- |
| `bundle_id` | Stable bundle id referenced by `packet_id`. |
| `packet_digest` | Digest of the canonical Markdown packet. |
| `entries[].ref` | Evidence ref used in packet rows. |
| `entries[].source_class` | Source class such as `git`, `ci`, `harness`, `review`, `witness`, or `external_assertion`. |
| `entries[].digest` | Digest of retained evidence or retained metadata. |
| `entries[].retained_form` | `raw`, `redacted`, `digest_only`, `external_ref`, or `not_retained`. |
| `entries[].redaction_status` | `not_needed`, `redacted`, `digest_only`, `withheld`, or `cannot_verify`. |
| `entries[].resolver` | How a reviewer can resolve or verify the ref, or `not_assessed`. |

A packet row that cites an evidence ref absent from the bundle manifest is
`cannot_verify`.

For `local-enterprise-baseline-v0`, the Markdown packet MAY show a compact
Evidence Bundle table with only `ref`, `source_class`, `retained_form`, and
`resolver`. The generated bundle manifest still preserves the canonical fields
above when tooling is available. Closed-contour resolver values MAY include
`internal_ref`, `not_resolvable_outside_customer`, or `not_assessed`.

## Input Boundary

Product Contract v0 accepts available evidence, not an idealized complete
delivery record.

Minimum viable packet input:

- change identity: repository/source set plus commit range, PR/MR id, or local
  change id;
- selected evidence profile;
- packet metadata and evidence bundle manifest, even if most entries are
  `not_retained`;
- at least one evidence ref for `PC-CHANGE`, or `PC-CHANGE` is
  `cannot_verify`;
- `PC-DECISION` row with owner refs or explicit `cannot_verify`/`not_assessed`.

Prompt/session evidence is optional in the Russian enterprise baseline. A local
Git plus internal CI packet is valid if it preserves missing agent-route and
initiator evidence as explicit gaps.

Missing inputs must remain `not_assessed` or `cannot_verify`; they must not be
filled by prose, inference, or checked-in status files alone.

## Packet Generation

The canonical packet is generated by `sdp-trace` tooling from retained evidence
inputs. Product Contract v0 does not implement that tooling, but future
implementation must treat generated packets as the product artifact and
hand-authored examples as contract examples only.

Generation MAY be triggered manually, by CI, by a change-host webhook, or by a
release process. The Russian enterprise baseline does not require continuous
generation; per-change or per-release packet generation is sufficient.

## Required Packet Rows

These rows are the first product contract. Every P0 feature must cite one or
more row ids and show forward progress on at least one cited row.

| row id | buyer question | required output |
| --- | --- | --- |
| `PC-CHANGE` | What software change is being assessed? | Change host, repo/source set, branch, PR/MR/local id, commit range, and source refs, or `cannot_verify`. |
| `PC-INITIATOR` | Who or what initiated the change? | Human, upstream general agent, task, issue, prompt boundary, or `not_assessed`. |
| `PC-AGENT-ROUTE` | Through which agent/tool/harness route did the change pass? | Agent runtime, coding tool, harness, adapter, and delegation chain, with source class for each hop. |
| `PC-MUTATION` | What repository or artifact mutation was observed? | Changed paths, commits, generated artifacts, infra config, or package/release refs, with evidence refs. |
| `PC-VERIFICATION` | What verification exists and who witnessed it? | Agent-claimed, harness-observed, CI-witnessed, reviewer-observed, customer-witnessed, signed, or missing evidence rows. |
| `PC-REVIEW` | What review evidence exists? | Reviewer plane, reviewer identity/source class, retained result, independence state, and gaps. |
| `PC-AUTHORITY` | Did observed actions stay within declared authority? | Authority fact state or `not_assessed`; no merge, blame, employment, compliance, or policy verdict. |
| `PC-THEATER` | Which proof gaps could mislead a buyer? | Evidence theater findings using required reason codes and trigger evidence refs. |
| `PC-ATTESTATION` | Is this packet signed or externally witnessed? | Signed profile, signer/witness refs, freshness, or `not_assessed`/`cannot_verify`. |
| `PC-DECISION` | Who owns the next human decision? | Owner refs for merge, release, risk acceptance, and security review, or explicit missing states. |
| `PC-RESIDUAL-GAPS` | What remains unknown? | Explicit residual gaps with state, reason, source row, and what evidence would close them. |

## Evidence States

Packet rows MUST preserve these states without collapsing them into a score:

| state | meaning |
| --- | --- |
| `pass` | The row was assessed and the required evidence for the selected profile is present. For summary rows, `pass` means the summary was computed according to its rules, not that no risk exists. |
| `partial` | Some relevant evidence exists, but the selected profile still has material gaps. |
| `fail` | Evidence contradicts the required claim, a required check fails, or a packet/projection would mislead if treated as successful. |
| `not_assessed` | The row was not assessed because the selected profile, inputs, or task scope did not require or provide evidence. |
| `cannot_verify` | Evidence is claimed or referenced but cannot be read, resolved, integrity-checked, or bound to the row. |
| `missing_telemetry` | Required telemetry or event data was expected for the selected profile but was not produced or retained. |
| `unsupported` | The selected tool, host, profile, or evidence surface is known not to provide the required data. |
| `not_integrated` | `sdp-trace` has relevant substrate capability, but that capability was not connected or invoked for this packet. |

Rows may also classify evidence source strength:

- `agent_claimed`
- `harness_observed`
- `change_host_observed`
- `ci_witnessed`
- `reviewer_observed`
- `customer_witnessed`
- `signed`
- `external_assertion`

Source strength classes are categorical, not ordinal. Projections MUST NOT
rank, aggregate, color-score, or present source classes as trust scores,
confidence levels, maturity levels, or readiness levels.

## Profile Taxonomy

| profile id | required inputs | rows that can normally reach `pass` | rows commonly `not_assessed` or `partial` | notes |
| --- | --- | --- | --- | --- |
| `local-enterprise-baseline-v0` | local Git/source refs, packet metadata, evidence bundle manifest | `PC-CHANGE`, `PC-MUTATION`, `PC-RESIDUAL-GAPS` | `PC-INITIATOR`, `PC-AGENT-ROUTE`, `PC-REVIEW`, `PC-AUTHORITY`, `PC-ATTESTATION` | Does not require GitHub, SaaS, raw prompts, public signing, or harness sessions. Internal CI/private artifacts may improve `PC-VERIFICATION`. |
| `change-host-rich-v0` | change-host PR/MR refs, review/check metadata, source refs | `PC-CHANGE`, `PC-REVIEW`, parts of `PC-DECISION` | `PC-AGENT-ROUTE`, `PC-INITIATOR` | GitHub, GitFlic, GitLab, Gitea/Forgejo, or custom hosts map to provider-neutral rows. |
| `harness-observed-v0` | harness/session profile, retained event refs, redaction policy | `PC-AGENT-ROUTE`, `PC-MUTATION`, parts of `PC-VERIFICATION` | `PC-REVIEW`, `PC-ATTESTATION` | Harness evidence proves only observed session facts, not methodology compliance. |
| `signed-v0` | packet digest, bundle digest, signer/witness refs, signer policy, freshness evidence | `PC-ATTESTATION` | underlying rows keep their original states | Signing is additive evidence. It never upgrades `not_assessed` or `cannot_verify` row states. |

Rows commonly `not_assessed` or `partial` for a profile are expected profile
characteristics, not product defects.

Profiles may be combined. Combined profiles evaluate each row independently.
A combined profile does not upgrade a row from `not_assessed`,
`cannot_verify`, `unsupported`, or `not_integrated` unless that same row has
qualifying evidence under another selected profile. No global trust verdict is
computed.

## Evidence Theater v0

Packet v0 MUST support these P0 theater findings:

| reason code | trigger condition | required trigger evidence |
| --- | --- | --- |
| `agent_claimed_verification` | Agent or harness text claims tests/checks/completion, but no independent retained verification artifact exists. | Claim ref plus absence of qualifying independent retained artifact. |
| `unbound_intent` | Change exists, but task, issue, prompt boundary, approval, or initiating actor is not bound. | `PC-CHANGE` evidence plus missing or unverified `PC-INITIATOR` evidence. |
| `ci_theater` | A CI status, check result, or build artifact is referenced as evidence of verification success, but the selected evidence profile lacks retained coverage for the specific verification claim. | CI/check claim ref plus missing retained artifact or coverage ref. |
| `scope_theater` | One observed tool/model/harness run is generalized into broad support, readiness, compatibility, or trust. | Product/backlog/support claim plus only one narrow observed evidence surface. |

`Independent` means the verification artifact was produced or witnessed by a
system or person that is not the same agent runtime making the claim.
`Retained` means the artifact or its verifying metadata is present in the
evidence bundle manifest or resolvable external ref; memory, chat prose, and
unretained terminal output do not count.

Theater findings are claims. Each finding MUST cite the trigger evidence or the
specific missing evidence that caused it. A theater finding without a cited
trigger is itself an overclaim.

`PC-THEATER` row state uses these rules:

| condition | row state |
| --- | --- |
| P0 theater assessment ran and no findings were triggered | `pass` |
| P0 theater assessment ran and one or more findings were triggered | `partial` |
| Selected profile did not assess theater findings | `not_assessed` |
| Required trigger refs are malformed, stale, or contradictory | `cannot_verify` |
| A projection or packet claims "no theater" while trigger evidence exists | `fail` |

These P1 theater findings are deferred unless a selected profile requires them:

- `actor_laundering`
- `review_theater`
- `artifact_theater`
- `human_approval_theater`

Deferral must be visible in `PC-RESIDUAL-GAPS`, not hidden.

## Row-Specific Rules

### PC-VERIFICATION

If agent text claims tests passed and no independent retained verification
artifact exists, `PC-VERIFICATION` cannot be `pass`.

If harness evidence shows a verification command but no independent retained CI
or customer witness exists, `PC-VERIFICATION` is `partial`, and theater findings
may include `agent_claimed_verification` or `ci_theater`.

### PC-AUTHORITY

`PC-AUTHORITY` projection states are limited to:

- `within_declared_authority`
- `exceeded_declared_authority`
- `not_assessed`
- `cannot_verify`

These states describe the relationship between observed actions and a selected
declared authority envelope. They do not describe compliance, blame,
discipline, employment fitness, merge approval, release approval, or policy
adherence.

### PC-DECISION

`PC-DECISION` must cover merge, release, risk acceptance, and security review.
Each decision has its own owner ref and owner state.

Owner refs may come from change-host metadata, task-system ownership, internal
role registry, approval policy, customer private PKI subject, or explicit
external assertion. External assertions remain `cannot_verify` until bound to a
policy or authority source.

Decision Ownership table entries use these owner states:

- `owner_bound`: owner ref is bound to a retained source for that decision.
- `owner_asserted`: owner is named, but only by an unbound external assertion.
- `not_assessed`: no owner evidence was assessed for that decision.
- `cannot_verify`: an owner ref is claimed but cannot be resolved or checked.
- `not_in_scope`: the selected profile explicitly excludes that decision.

`owner_bound` means the next decision owner is identifiable. It does not mean
the decision has been made, the change is ready, or approval has been granted.

For `local-enterprise-baseline-v0`, owner binding MAY use a retained local role
assignment, Git `Signed-off-by` trailer, internal wiki reference, or team
convention ref. Formal policy refs are required only when the selected profile
includes an authority envelope.

### PC-RESIDUAL-GAPS

`PC-RESIDUAL-GAPS` is computed from:

- every required packet row whose state is not `pass`;
- every active theater finding;
- every deferred P1 theater category relevant to the selected profile;
- every evidence ref that is absent from the bundle manifest.

The row can be `pass` only when all known gaps are enumerated with source row,
state, reason, and required closure evidence. If gap synthesis cannot inspect
required rows or bundle refs, the row is `cannot_verify`.

## Russian Enterprise Baseline

Product Contract v0 must work in a local or self-hosted enterprise environment.

The baseline profile MUST NOT require GitHub, public SaaS, public transparency
logs, public Sigstore/Rekor, raw prompt export, or broad employee monitoring.

Baseline inputs may be limited to:

- local Git refs;
- GitFlic, GitLab, Gitea/Forgejo, or custom change-host refs when available;
- Jenkins, TeamCity, GitLab CI, GitHub Actions, or custom CI artifacts when
  available;
- local artifact store refs;
- customer private PKI or internal witness records when available;
- digest-only or redacted harness/session evidence when available.

Redaction policy for the baseline profile defaults to `digest_only` for prompts,
model outputs, and session data. Customer-specific redaction MAY follow internal
classification schemes; the packet records the selected policy but does not
enforce it.

If those sources are absent, the packet still has value when it makes missing
evidence explicit and names the decision owner state. It must not pretend the
missing evidence exists.

## P0 Classification Rule

Any future P0 backlog item must include:

- `packet_rows`: one or more required row ids from this spec;
- `evidence_surface`: the exact artifact, API, log, schema, fixture, or command
  output to inspect;
- `start_state`: the current row state or explicit `not_assessed`;
- `target_transition`: what changes in the row after the work;
- `buyer_effect`: what the CTO packet becomes clearer about after the feature;
- `non_goal`: what the feature still does not prove.

The target transition must show forward progress on at least one cited row.
Forward progress means one of:

- a row state moves from unknown/missing to `partial`, `pass`, `fail`,
  `unsupported`, or `not_integrated` with evidence;
- a row state moves from `partial` to `pass` or `fail` by closing or confirming
  the material gap;
- an evidence surface becomes inspectable and produces retained refs;
- a theater reason code gains a derivation rule and trigger evidence;
- a row becomes more precise by separating observed facts from claims;
- a projection adds a required non-goal or residual gap that a previous packet
  omitted, and that gap names a specific row, evidence surface, and closure
  condition.

Repeating `not_assessed -> not_assessed` or `cannot_verify -> cannot_verify`
without a new evidence surface, narrower claim, or explicit unsupported state
does not qualify as P0 product progress.

`buyer_effect` MUST cite the packet row, packet section, or theater reason code
that becomes clearer. Vague claims such as "more visibility" do not qualify.

## User Scenarios & Testing

### User Story 1 - CTO Reads One Packet (Priority: P0)

A CTO can read a Change Evidence Packet v0 for one software change in under two
minutes and identify the change, agent route, witnessed evidence, agent claims,
theater findings, residual gaps, and next decision owner.

**Why this priority**: The first product proof is not that substrate exists. The
first product proof is that substrate becomes a buyer-readable artifact.

**Independent Test**: Given one synthetic or real change with partial evidence,
the packet separates observed facts, claims, missing evidence, and decision
ownership without requiring raw logs.

### User Story 2 - Product Owner Classifies Substrate-Only Work (Priority: P0)

A product owner can classify a proposed P0 feature as substrate or discovery
when it names a packet row but does not make forward progress on that row.

**Why this priority**: This prevents the repeated failure where useful substrate
features are treated as product contract progress.

**Independent Test**: Given a proposed integration item, the intake checklist
classifies it as P0 product, substrate, discovery, or future integration based
on packet row mapping and target transition.

### User Story 3 - Engineer Maps Existing Capabilities (Priority: P0)

An engineer can inspect `traceability.md` and see which existing `sdp-trace`
capabilities already feed packet rows and which rows remain gaps.

**Why this priority**: The contract must not erase current work. It must show
how current substrate becomes product output.

**Independent Test**: At least one current doc/schema/example path is mapped to
each known row or marked as a gap with required evidence.

### User Story 4 - Russian Enterprise Pilot Works Without GitHub (Priority: P0)

A local/self-hosted enterprise pilot can produce a useful packet with local Git,
self-hosted change-host refs, private CI/artifact refs, and digest-only evidence
even when GitHub APIs, public signing, and harness/session evidence are
unavailable.

**Why this priority**: A Russian-market enterprise target cannot depend on
GitHub-native product semantics.

**Independent Test**: The local-baseline example includes local Git plus
internal CI/artifact refs and preserves missing agent-route evidence as gaps.

## Functional Requirements

- **FR-001**: Product Contract v0 MUST define Change Evidence Packet v0 as the
  first buyer-facing output.
- **FR-002**: The canonical packet MUST be portable Markdown plus an evidence
  bundle; static HTML, PR/MR comments, CLI summaries, PDFs, and signed envelopes
  are projections.
- **FR-003**: Product Contract v0 MUST define required packet rows and stable
  row ids.
- **FR-004**: Every P0 feature MUST cite at least one required packet row and a
  target transition for that row.
- **FR-005**: A feature that does not cite a packet row or does not make forward
  row progress MUST NOT be labeled P0 product progress.
- **FR-006**: Packet rows MUST distinguish observed facts, agent claims,
  independent witness evidence, external assertions, and missing evidence.
- **FR-007**: Packet rows MUST preserve `not_assessed` and `cannot_verify`.
- **FR-008**: Packet v0 MUST include the four P0 evidence theater reason codes
  and derivation triggers.
- **FR-009**: Deferred P1 theater findings MUST be visible as residual gaps when
  relevant.
- **FR-010**: Product Contract v0 MUST include a Russian enterprise baseline
  profile that does not require GitHub, SaaS dashboards, public Sigstore/Rekor,
  raw prompt export, or broad employee monitoring.
- **FR-011**: Tool, harness, model, or change-host support claims MUST NOT be
  made until an evidence surface is inspected and mapped to packet rows.
- **FR-012**: Signed attestation MUST be additive evidence over a packet and
  MUST NOT upgrade underlying row states.
- **FR-013**: Product Contract v0 MUST include a traceability matrix from current
  substrate capabilities to packet rows.
- **FR-014**: Product Contract v0 MUST include one harness-observed example and
  one local-enterprise-baseline example, both marked example-only and not product
  proof.
- **FR-015**: Roadmap items in `003-agent-supply-chain-roadmap` MUST NOT be
  classified as P0 product implementation work until they map to Product
  Contract v0 rows and target transitions.

## Success Criteria

- **SC-001**: A CTO proxy can identify change, route, evidence, gaps, theater
  findings, and decision owner from the example packet in under two minutes.
- **SC-002**: A reviewer can classify a substrate-only feature as non-P0 using
  row mapping and target transition.
- **SC-003**: Existing `sdp-trace` capabilities are preserved as substrate
  inputs and mapped to packet rows rather than dismissed.
- **SC-004**: The baseline packet remains useful when GitHub, public signing,
  full agent session export, and OpenCode/GSD evidence are unavailable.
- **SC-005**: A focused Socratic review finds no remaining critical or major
  ambiguity in the product classification rule before implementation approval
  is requested.

## Open Decisions With Proposed Defaults

| decision | proposed default | why |
| --- | --- | --- |
| First artifact | Markdown report plus evidence bundle | Works offline/on-prem, diffable, reviewable, projection-friendly. |
| First projection | Static HTML after packet semantics stabilize | Good for demos, but not canonical. |
| First change-host rich adapter | GitHub | Default for fastest rich-evidence path; not a product-ontology commitment. |
| Russian baseline | local Git plus optional GitFlic/GitLab/Jenkins/TeamCity refs | Prevents GitHub-native P0 lock-in. |
| Signed attestation | P2 additive profile | Avoids signing weak evidence as theater. |
| General-purpose agents | in scope only with software-delivery boundary evidence | Prevents employee-monitoring drift. |
