# Buyer Demo Rehearsal

Date: 2026-05-26
Reviewer: Codex GPT-5, sdp-trace closure route
External repository: `fall-out-bug/sdp-trace-demo-jvm-gsd`

## Scope

This rehearsal checks one happy-path demo PR and one negative theater PR. It is
not merge approval, release approval, semantic-quality approval, production
trust, or signed external trust.

## Baseline Preservation

The external tag `demo-v1-observation-baseline` was created in
`fall-out-bug/sdp-trace-demo-jvm-gsd`.

```text
refs/tags/demo-v1-observation-baseline^{}:
a8f37aad8500761693feb6ce68517bd65cabc8cc
```

That commit is PR #16's `baseRefOid`, preserving the pre-demo-v2 repository
state before packetization setup.

## Happy Path

Selected happy-path PR: #20, `[demo] info endpoint slice`.

Live GitHub state:

- PR #20 state: `MERGED`
- Head: `5373c8d2b284825e1503c9d6c6df85a283ee0109`
- Merge commit: `1241956aea77a9296dce4bcfbe95846580324be2`
- `build-and-test`: `SUCCESS`

Replay:

```text
sdp-trace packet validate --bundle .sdp-trace/bundles/feature-5/bundle.json
state: pass

sdp-trace packet check-demo --bundle .sdp-trace/bundles/feature-5/bundle.json
state: pass
```

## Negative Theater Path

Selected negative PR: #21, `DEMO-NEGATIVE: agent-claimed verification without
independent artifact`.

Live GitHub state:

- PR #21 state: `OPEN`
- Draft: `true`
- Head: `50b7ed5ca146a32ce289dc9ee29aa34c8919439d`
- `build-and-test`: `SUCCESS`

Replay:

```text
sdp-trace packet validate --bundle .sdp-trace/bundles/negative/bundle.json
state: pass

sdp-trace packet check-demo --bundle .sdp-trace/bundles/negative/bundle.json
state: fail
errors:
- demo first-packet gate requires PC-AGENT-ROUTE evidence from retained structured OpenCode/GSD/MiniMax harness route observation
- demo first-packet gate requires PC-VERIFICATION or PC-REVIEW to be pass, partial, or fail
```

## Rehearsal Result

The happy-path packet is demo-checkable and the negative packet is structurally
valid but intentionally non-green for buyer explanation. This satisfies the
rehearsal contrast without hiding missing evidence or upgrading the negative
case into a pass.

Remaining non-claims:

- Packets are evidence organization only.
- The negative PR remains unmerged and draft.
- No production trust, release readiness, compliance, semantic-quality approval,
  or signed external attestation is claimed.
