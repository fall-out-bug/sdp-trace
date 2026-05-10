# Feature Specification: GitHub OSS Demo Packet

**Feature Branch**: `007-github-oss-demo-packet`
**Created**: 2026-05-10
**Status**: Draft - needs Socratic review before implementation approval
**Input**: 006 Change Evidence Packet Core plus the existing
`sdp-trace-demo-jvm-gsd` repository.

**Dependency**: Do not start demo implementation until
`specs/006-change-evidence-packet-core/` is approved and the product slice has
at least schema, validator, renderer, and CLI validate/render behavior.

## Product Boundary

This slice turns the current demo direction into a GitHub-first OSS ecosystem
demo for the product artifact delivered by 006: Change Evidence Packet v0.

The demo target is a CTO-visible GitHub PR flow:

```text
GitHub issue/task -> OSS agent/harness work -> branch/PR -> CI/checks/artifacts
-> review evidence -> Change Evidence Packet v0 -> CTO reads the packet
```

The demo is not an enterprise closed-contour deployment, signed-attestation
profile, semantic quality evaluator, employee-monitoring system, or proof of
broad OSS tool compatibility.

## Core Demo Claim

The demo may claim:

> For a GitHub PR produced through one demonstrated OSS agent/harness route
> (OpenCode + GSD + MiniMax-M2.5), `sdp-trace` can produce a Change Evidence
> Packet v0 that separates observed facts, agent claims, verification evidence,
> review evidence, theater findings, residual gaps, and next decision ownership.

The demo MUST NOT claim:

- the change is semantically correct unless review/test evidence supports it;
- merge, release, production trust, or compliance approval;
- support for all OpenCode/GSD, `pi`, GSD2, Superpowers, Hermes, OpenClaw, or
  other OSS agents;
- signed trust;
- local/self-hosted enterprise readiness.

## Existing Demo Repository Assessment

Current demo repository:

- GitHub repo: `fall-out-bug/sdp-trace-demo-jvm-gsd`
- Local path: `/Users/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd`
- Stack: Kotlin/JVM, Bazel, GitHub Actions, OpenCode, GSD,
  `minimax-coding-plan/MiniMax-M2.5`
- Current app: Todo REST API with five features.

The current repo is useful as a v1 observation-history artifact. It already has
source code, CI, feature commits, review artifacts, and evidence that
`sdp-trace observe session` can retain model, phase, mutation, tool, and test
signals from OpenCode/GSD runs.

It is not sufficient as the v2 product demo because it does not consistently
bind every feature to:

- a GitHub issue or ticket;
- one PR per feature;
- a retained GitHub CI run and artifact bundle per feature;
- a generated Change Evidence Packet v0 per feature;
- PR-linked review evidence;
- row-level packet states and residual gaps.

## Demo Track Policy

Default recommendation: **continue in the existing
`sdp-trace-demo-jvm-gsd` repository through a v2 packetization track**.

The current repository is not invalid. It is invalid only as a finished CTO
packet demo. As source material, it is exactly the kind of real agent-delivery
history `sdp-trace` should be able to organize into a buyer-readable packet.

Do not destroy or hide the current demo history. Preserve it as
`demo-v1-observation` history because it contains useful evidence about
product-observation fixes, review contamination, and real OpenCode/GSD
friction.

Recommended track:

1. Tag current `main` as `demo-v1-observation-baseline`.
2. Keep current PR/history available as real-world input evidence.
3. Create a v2 packetization milestone in the same repository.
4. Start with one selected existing PR or feature history and generate a
   Change Evidence Packet v0 from it.
5. Add missing GitHub/CI/review/packet wiring through new PRs in the same repo,
   using the 006 renderer and validator.
6. Only create new feature PRs when the 006 packet path can capture them.

Alternative tracks:

| option | when to use | trade-off |
| --- | --- | --- |
| Existing repo v2 packetization | Default | Best product proof because it organizes a real, messy agent-delivery history. |
| Existing repo fresh `demo-v2` root branch | Use if current app/history makes the buyer story too noisy | Cleaner PR story but weaker proof that old agent history can be interpreted. |
| New repo `sdp-trace-demo-github-oss-packet` | Use only for a polished public sales demo | Cleanest narrative, but it avoids the harder product problem already present in v1. |

The first implementation step should not be "rebuild the whole demo". It should
be "produce one honest packet from existing GitHub/demo evidence, with explicit
gaps."

## First Packet Selection

The first CTO-visible packet MUST be selected by evidence richness, not by which
feature looks prettiest.

Initial feature evidence inventory:

| slice | current evidence | first-packet use |
| --- | --- | --- |
| v1 baseline / setup | GitHub repo, CI workflow, Bazel/Kotlin app, observation ledger context | Useful context, not a feature packet. |
| Feature 1 health | README feature row and app/test history; exact PR/check/artifact availability must be inspected | Candidate only if GitHub PR and CI artifacts are still resolvable. |
| Feature 2 create/list todo | Review artifacts exist; exact PR/check/artifact availability must be inspected | Candidate if review and CI refs resolve. |
| Feature 3 complete todo | Review artifact exists; exact PR/check/artifact availability must be inspected | Candidate if PR/check refs resolve. |
| Feature 4 delete todo | Known contaminated review plus re-review exists | Not first packet unless contamination is surfaced explicitly. |
| Feature 5 stats | Recent feature/review artifacts exist; exact PR/check/artifact availability must be inspected | Likely first candidate if PR/check/artifact refs resolve. |

Selection rule:

1. Inventory every existing feature/history slice for available GitHub issue or
   PR body task source, PR metadata, commit range, GitHub Actions run, artifact
   availability, review evidence, and OpenCode/GSD observation evidence.
2. Choose the slice with the most retained, resolvable evidence surfaces.
3. Record the selection rationale in packet metadata `generated_from`.
4. If no existing slice can meet the first-packet minimum bar, the first packet
   MUST come from a new v2 feature PR under the packetization track.

Minimum bar for the first CTO-visible packet:

- at least four required packet rows are `pass` or `partial`;
- `PC-CHANGE` and `PC-MUTATION` have retained, resolvable evidence refs;
- at least one of `PC-VERIFICATION`, `PC-REVIEW`, or `PC-AGENT-ROUTE` exits
  `not_assessed`, meaning the row is assessed as `pass`, `partial`, or `fail`;
- no more than one row is `cannot_verify` without a concrete closure path in
  `PC-RESIDUAL-GAPS`;
- every remaining missing row has a one-line reason.

007 MUST NOT use hand-authored packets as demo proof. Hand-authored packets are
allowed only inside 006 pre-renderer fixture validation. The demo starts after
the 006 renderer exists, so CTO-visible demo packets must be tool-generated.

## Contamination Handling

Known contamination is evidence, not dirt to hide.

If a selected packet target has known contamination history, the packet MUST do
one of:

- include the contamination in `PC-RESIDUAL-GAPS` and, when it could mislead a
  buyer, trigger a theater finding; or
- exclude the feature from CTO-visible packet selection and record the reason
  in the demo tracker.

Exclusion without a theater finding is acceptable only when the excluded
feature never appears in any CTO-visible packet or tracker surface shown to the
buyer.

Feature 4 in the existing demo repo has known review contamination and
re-review history. It MUST NOT be the first CTO-visible packet unless the packet
surfaces that history explicitly.

## Provenance And Retention Rules

Task source evidence must be time-bound and resolvable.

- A GitHub issue is the preferred task source.
- A PR body task section is acceptable only with the PR body revision or binding
  commit/ref recorded when available.
- A retained task artifact must have a digest, source ref, and resolver entry in
  the evidence bundle manifest.
- Retroactive task binding MUST record binding commit SHA or packet-generation
  timestamp. Retroactive bindings cannot upgrade `PC-INITIATOR` beyond
  `partial`.

GitHub Actions artifact retention is part of the proof:

- v2 demo workflows MUST retain packet and bundle artifacts for at least 180
  days;
- retroactive packets MUST inspect artifact availability at packet-generation
  time;
- expired artifacts produce `cannot_verify` with reason `artifact_expired`, or
  `partial` only when another retained evidence ref still supports the row;
- GitHub CI green is not `PC-VERIFICATION: pass` unless the packet has retained
  check/artifact evidence in the bundle.

When GitHub CI, harness observation, review artifact, or packet data
contradict, the affected row is `partial`; both refs are listed, and
`PC-RESIDUAL-GAPS` names the contradiction and closure evidence.

## Required Demo Evidence

Every feature PR in the v2 demo MUST have:

| evidence | minimum requirement |
| --- | --- |
| Task source | GitHub issue, PR body task section, or retained task artifact. |
| Agent route | OpenCode/GSD observation run with selected model route. |
| Mutation | GitHub commit range plus observed mutation evidence when available. |
| Verification | GitHub Actions check and retained artifact, or explicit `cannot_verify`. |
| Review | GitHub PR review, retained external review artifact, or `not_assessed`. |
| Packet | Generated `change-evidence-packet.md` and evidence bundle manifest. |
| Theater | P0 theater assessment result, including no-finding state when clean. |
| Decision owner | next owner state, not approval. |

## Required Packet Rows

Every demo packet MUST include all Product Contract v0 rows:

- `PC-CHANGE`
- `PC-INITIATOR`
- `PC-AGENT-ROUTE`
- `PC-MUTATION`
- `PC-VERIFICATION`
- `PC-REVIEW`
- `PC-AUTHORITY`
- `PC-THEATER`
- `PC-ATTESTATION`
- `PC-DECISION`
- `PC-RESIDUAL-GAPS`

The demo should expect some rows to be `not_assessed`. That is acceptable when
the packet says so clearly, but an all-gap packet is not a successful CTO demo.

## Demo Feature Set

Use a deliberately small app. The app exists to produce inspectable PRs, not to
impress through domain complexity.

Recommended v2 feature sequence:

1. `F1`: health endpoint and Bazel/CI baseline.
2. `F2`: create/list todo.
3. `F3`: complete todo.
4. `F4`: delete todo.
5. `F5`: stats endpoint.

Each feature must be its own issue, branch, PR, CI run, packet, and review
record. Setup work must be separate from feature work.

## Negative / Theater Demo

The demo needs one controlled theater example, separate from the happy-path
feature PRs.

Recommended example:

- create a GitHub draft PR named with prefix `DEMO-NEGATIVE:`;
- label it `demo-theater`;
- make the agent claim verification;
- omit or break the independent verification artifact binding;
- generate a packet that marks `PC-VERIFICATION` as `partial` or
  `cannot_verify`;
- trigger `agent_claimed_verification` as the primary first negative example;
- keep the PR unmerged or mark it as demo-negative.

This is important because the product promise is not "green checks look nice";
the promise is "the packet makes misleading evidence visible."

`ci_theater` may be a later negative example, but it is not the primary first
negative demo.

## User Scenarios & Testing

### User Story 1 - CTO Reads A GitHub PR Packet (Priority: P0)

A CTO opens one GitHub PR and can find a packet artifact or PR comment that
answers: what changed, which agent route produced it, what evidence exists,
what is missing, and who owns the next decision.

**Independent Test**: Given a feature PR URL, a reviewer can locate the packet
and answer those questions in under five minutes without reading raw logs.

### User Story 2 - Engineer Inspects Evidence Bundle (Priority: P0)

An engineer can inspect the evidence bundle manifest and resolve packet refs to
GitHub PR data, CI artifact data, OpenCode/GSD observation output, and review
evidence.

**Independent Test**: Every `evidence refs` value in the packet exists in the
bundle manifest or the row is `cannot_verify`.

### User Story 3 - Product Owner Sees Theater (Priority: P0)

A product owner can see that a demo-negative PR is not hidden. The packet names
the specific theater reason code and required closure evidence.

**Independent Test**: The negative PR packet contains a triggered theater
finding and does not claim merge/release readiness.

### User Story 4 - Demo Does Not Overclaim OSS Support (Priority: P0)

The demo shows OpenCode/GSD support for the observed profile and names all other
OSS tools as future evidence surfaces unless they have their own retained runs.

**Independent Test**: README and packets do not claim broad `pi`, GSD2,
Superpowers, Hermes, OpenClaw, or enterprise support from a single OpenCode/GSD
run.

## Functional Requirements

- **FR-001**: The v2 demo MUST use GitHub as the first change-host evidence
  surface.
- **FR-002**: The v2 demo MUST generate Change Evidence Packet v0 artifacts, not
  only raw observation ledgers.
- **FR-003**: Each feature PR MUST map to all required packet rows.
- **FR-004**: Each packet MUST include an evidence bundle manifest.
- **FR-005**: Each packet MUST preserve `not_assessed` and `cannot_verify`
  states.
- **FR-006**: Each feature PR MUST have a GitHub Actions check or an explicit
  `cannot_verify` verification state.
- **FR-007**: Each feature PR MUST have review evidence or explicit
  `not_assessed`.
- **FR-008**: The demo MUST include one controlled theater example.
- **FR-009**: The demo MUST NOT treat packet generation as merge, release,
  compliance, production trust, or quality approval.
- **FR-010**: Existing `sdp-trace-demo-jvm-gsd` history MUST be preserved before
  any v2 packetization track or optional sales-demo branch starts.

## Success Criteria

- **SC-001**: At least one feature PR has a generated Change Evidence Packet v0
  that a CTO proxy can understand in under five minutes.
- **SC-002**: All five feature PRs have packet artifacts or are explicitly
  marked `not_assessed` in the demo tracker.
- **SC-003**: The negative/theater PR demonstrates a non-green packet state
  without being confused with product failure.
- **SC-004**: No demo artifact claims broad OSS support, enterprise readiness,
  signed trust, semantic quality, merge approval, or release readiness.
- **SC-005**: A reviewer can reproduce packet evidence refs from GitHub PR
  metadata, CI artifacts, observation output, and review artifacts.
- **SC-006**: At least one happy-path feature PR has `PC-THEATER: pass`, showing
  that theater assessment can produce a clean result when evidence is properly
  bound.

## Open Decisions With Proposed Defaults

| decision | proposed default | why |
| --- | --- | --- |
| Demo repo path | Existing repo `sdp-trace-demo-jvm-gsd`, v2 packetization milestone | Proves `sdp-trace` can organize real agent-delivery history instead of only a sterile new demo. |
| First OSS route | OpenCode + GSD + MiniMax-M2.5 | Already observed and closest to current proof. |
| First packet projection | GitHub PR comment plus uploaded artifact | CTO can inspect from PR; artifact preserves offline evidence. |
| Number of happy-path feature PRs | Five | Matches existing demo framing and prior demo contract. |
| Negative demo | One draft PR with missing independent verification | Shows evidence theater without poisoning happy path. |
| Enterprise/local profile | Out of scope for this demo | Future paid profile; not demo P0. |
