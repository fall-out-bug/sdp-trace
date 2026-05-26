# External Demo Evidence Review

Date: 2026-05-26
Reviewer: Codex GPT-5, sdp-trace closure route
External repository: `fall-out-bug/sdp-trace-demo-jvm-gsd`

## Scope

This review maps current external demo-repository evidence to Spec 007 tasks.
It does not claim buyer-demo completion, merge approval, release approval,
production trust, semantic quality approval, or signed external trust.

## Live Evidence Checked

- Cloned `https://github.com/fall-out-bug/sdp-trace-demo-jvm-gsd.git` at
  `main` head `b101882`.
- Queried PRs #16-#23 with GitHub CLI.
- Ran `bazel test //app:smoke_test --test_output=errors` on `main`; result:
  pass.
- Checked out draft PR #21 and ran `bazel test //app:smoke_test
  --test_output=errors`; result: pass.
- Ran `sdp-trace packet validate --bundle` and `sdp-trace packet check-demo
  --bundle` against feature bundles 1-5; all returned `{"state":"pass"}`.

## Evidence Summary

| Area | Evidence | State |
| --- | --- | --- |
| Setup PR | Demo PR #16 merged; `build-and-test` SUCCESS; `docs/setup-boundary-review.md` exists. | pass |
| Feature packets 1-5 | Demo PRs #16-#20 merged; packets and bundles exist; validate/check-demo pass. | pass with residual `partial` route rows |
| Feature CI | Demo PRs #16-#20 all have `build-and-test` SUCCESS in GitHub. | pass |
| Negative theater PR | Demo PR #21 is open draft with `DEMO-NEGATIVE` title, CI SUCCESS, negative packet and bundle on branch. | pass for draft negative artifact existence |
| First-packet gate | Feature bundles 1-5 pass `packet check-demo`. | pass |
| Buyer rehearsal | No live rehearsal artifact found in the checked external repository or PR metadata. | not_assessed |
| Demo v1 archival tag/note | No tag was visible in the shallow clone; no dedicated archival note was identified. | not_assessed |

## Boundaries

- Packets remain evidence organization only.
- `PC-AGENT-ROUTE` remains `partial` where the packets say route proof is
  local or input-ref based.
- `PC-AUTHORITY`, `PC-ATTESTATION`, and `PC-DECISION` remain `not_assessed`
  where packet rows say so.
- Buyer-demo rehearsal remains open until a retained happy-path plus negative
  rehearsal artifact exists.
