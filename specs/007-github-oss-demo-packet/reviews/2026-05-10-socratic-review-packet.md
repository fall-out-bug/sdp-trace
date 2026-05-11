# Socratic Review Packet: GitHub OSS Demo Packet

Date: 2026-05-10
Scope: `specs/007-github-oss-demo-packet/`

## Target Decision

Should `sdp-trace` continue the CTO demo inside the existing
`fall-out-bug/sdp-trace-demo-jvm-gsd` repository through a v2 packetization
track, instead of creating a new clean demo repository by default?

## Review Target Files

- `specs/007-github-oss-demo-packet/spec.md`
- `specs/007-github-oss-demo-packet/plan.md`
- `specs/007-github-oss-demo-packet/demo-repo-plan.md`
- `specs/007-github-oss-demo-packet/tasks.md`
- `specs/005-product-contract-v0/spec.md`
- `specs/005-product-contract-v0/reviews/2026-05-10-rereview.md`

## Current Demo Repo Facts

Repository: `fall-out-bug/sdp-trace-demo-jvm-gsd`

Local path: `/Users/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd`

Observed current state:

```text
git status --short --branch
## main...origin/main
```

The current README says the repo is a minimal Todo REST API with five features
and explicitly states that `sdp-trace` records harness/process evidence but does
not infer semantic quality.

Recent history includes feature PRs and review artifacts, including a known
contaminated Feature 4 review that was later marked and re-reviewed. This is
important: the existing repo is not clean, but the product should be able to
make this kind of history legible rather than hide it.

The `sdp-trace` observation ledger records that live OpenCode/GSD/MiniMax runs
eventually produced observed model, phase, mutation, tool, and test evidence.
It also records that PR/review/merge/gap families were explicitly
`not_assessed` for earlier observation claims until those stages were captured.

## Product Contract Context

Product Contract v0 says the first buyer-facing artifact is Change Evidence
Packet v0: Markdown plus evidence bundle.

Each packet must include:

- Executive Summary
- Packet Metadata
- Required Rows
- Theater Findings
- Decision Ownership
- Evidence Bundle
- What This Packet Does Not Prove

Important row ids:

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

P0 product progress requires `packet_rows`, `evidence_surface`,
`start_state`, `target_transition`, `buyer_effect`, and `non_goal`.

## Known Risk

Creating a new clean repo would make a prettier sales story, but could dodge the
harder product question: can `sdp-trace` organize a real messy agent-delivery
history into a buyer-readable evidence packet?

Keeping the existing repo could make the CTO demo noisy if the first packet
tries to explain too much history at once.

## Requested Review Output

Review the 007 package as a product/spec decision. Do not propose code.

Return:

1. Verdict: `APPROVE_FOR_USER_APPROVAL`, `REVISE_BEFORE_USER_APPROVAL`, or
   `KILL_DIRECTION`.
2. Findings table with columns: `id`, `severity`, `file/section`, `finding`,
   `exact fix`.
3. A direct answer to whether existing-repo packetization is the right default.
4. The smallest acceptable first implementation slice.
5. Any proof claims the demo must still refuse.

Severity:

- `critical`: must fix before user approval.
- `major`: should fix before user approval unless explicitly deferred.
- `minor`: polish.
