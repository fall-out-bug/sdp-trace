# Feature Specification: Change Evidence Packet Core

**Feature Branch**: `006-change-evidence-packet-core`
**Created**: 2026-05-10
**Status**: Draft - needs Socratic review before implementation approval
**Input**: Product Contract v0 needs a product artifact before the GitHub OSS
demo can honestly claim that `sdp-trace` produces buyer-readable evidence.

## Product Boundary

This slice turns Product Contract v0 into a minimal `sdp-trace` product surface:
`Change Evidence Packet v0`.

The packet is a CTO-readable artifact for one GitHub PR/change. It separates
observed facts, agent claims, verification evidence, review evidence, theater
findings, residual gaps, and next decision ownership.

This slice is not:

- a full GitHub crawler;
- enterprise deployment;
- signed attestation;
- employee monitoring;
- semantic code quality scoring;
- support for every OSS agent or harness;
- approval-to-merge.

## Core Product Claim

After this slice, `sdp-trace` may claim:

> Given a structured evidence bundle for one GitHub PR/change, `sdp-trace` can
> validate and render a Change Evidence Packet v0 that preserves missing,
> partial, contradictory, and unverified evidence without upgrading it to trust.

It MUST NOT claim:

- the change is correct;
- the change is approved;
- CI was green unless retained check/artifact evidence supports it;
- review happened unless review evidence is present;
- signed or external trust;
- broad OSS ecosystem support.

## Packet Contract

The packet MUST contain these rows:

- `PC-CHANGE`
- `PC-INITIATOR`
- `PC-AGENT-ROUTE`
- `PC-MUTATION`
- `PC-VERIFICATION`
- `PC-REVIEW`
- `PC-AUTHORITY`
- `PC-THEATER`
- `PC-RESIDUAL-GAPS`
- `PC-DECISION`

Each row MUST have:

- `state`: `pass`, `partial`, `fail`, `cannot_verify`, `not_assessed`, or
  `not_in_scope`;
- `summary`: CTO-readable one-line statement;
- `evidence_refs`: retained evidence references, empty only when the state
  explains why;
- `reason`: required for `partial`, `fail`, `cannot_verify`,
  `not_assessed`, and `not_in_scope`;
- `owner`: next owner for closure or decision.

No row may infer `pass` from absence of negative evidence.

## Evidence Bundle Contract

The evidence bundle manifest MUST bind packet rows to retained sources:

- GitHub PR metadata, or explicit `not_assessed`;
- task source: GitHub issue, PR body task section, or retained task artifact;
- commit range and changed files;
- verification surface: check run, workflow run, artifact, local command output,
  or explicit `cannot_verify`;
- review surface: GitHub review, retained external review, or `not_assessed`;
- agent/harness observation: OpenCode/GSD observation output when available;
- theater finding inputs and reason codes;
- resolver entries for every retained ref;
- digests for retained local artifacts.

Artifact expiry is evidence loss. Expired or missing artifacts MUST become
`cannot_verify` unless another retained ref still supports a weaker `partial`
state.

## Renderer Contract

The first renderer MUST produce canonical Markdown.

Markdown is the authoritative projection for this slice. HTML, PR comments, and
dashboards are later projections.

The rendered packet MUST show:

- packet metadata, including `packet_version`, source change, generation time,
  and `authoring_method`;
- row table with state, summary, evidence refs, reason, and owner;
- residual gaps section;
- theater section, including clean `PC-THEATER: pass` when no finding exists;
- artifact manifest pointer;
- explicit non-approval language.

If a packet is hand-authored before renderer completion, metadata MUST say
`authoring_method: hand_authored_before_tooling`. After the renderer exists,
new product/demo packets MUST be renderer-produced.

## Validator Contract

The validator MUST reject:

- unknown row IDs;
- unknown states;
- required rows missing;
- `pass` rows with no retained evidence refs, except rows explicitly allowed by
  schema with a self-contained rationale;
- `cannot_verify` or `not_assessed` rows without reason;
- evidence refs without resolver entries;
- expired artifact refs presented as `pass`;
- contradictory evidence without `partial` state and residual-gap explanation;
- PR comment/body packet projection marked canonical over the uploaded packet
  artifact.

The validator MAY allow `not_in_scope` for rows outside the declared demo or
product slice, but only with an explicit reason.

## GitHub Evidence Surface

This slice defines a minimal GitHub evidence input shape. It does not need live
GitHub API harvesting in P0.

P0 MAY use captured JSON/fixtures or user-provided refs for:

- PR number, URL, title, body, author, base/head branches, head SHA;
- commit range;
- check run or workflow run URL and conclusion;
- artifact name, URL/path, retention/expiry metadata when available;
- review URL/path and reviewer identity when available.

Live GitHub API import is a follow-up unless implementation proves it is
smaller than fixture/captured input.

## Theater Reason Codes

P0 theater reason codes:

- `agent_claimed_verification`
- `unbound_intent`
- `ci_theater`
- `scope_theater`

Every theater finding MUST cite trigger evidence or missing evidence. A clean
theater row is `PC-THEATER: pass`, not omitted.

## Success Criteria

- **SC-001**: A valid bundle renders a canonical Markdown packet with all
  required rows.
- **SC-002**: Missing verification evidence renders as `cannot_verify` or
  `not_assessed`, never `pass`.
- **SC-003**: Expired artifact evidence cannot support `pass`.
- **SC-004**: Contradictory evidence forces `partial` plus residual-gap text.
- **SC-005**: Theater reason code `agent_claimed_verification` can be rendered
  from fixture input.
- **SC-006**: A clean happy-path fixture renders `PC-THEATER: pass`.
- **SC-007**: The validator distinguishes canonical uploaded packet artifact
  from PR comment/body projection.
- **SC-008**: Demo spec 007 can depend on this slice without hand-authored
  packet proof after renderer completion.

## Out Of Scope

- Signed attestation.
- Enterprise policy/risk scoring.
- GitLab/Bitbucket.
- Broad OSS agent compatibility.
- Automatic review adjudication.
- Merge approval.
- Long-term evidence storage service.
