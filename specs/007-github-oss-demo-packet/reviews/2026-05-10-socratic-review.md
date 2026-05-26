# Socratic Review Ledger: GitHub OSS Demo Packet

Date: 2026-05-10
Scope: `specs/007-github-oss-demo-packet/`
Status: Prior Socratic review completed for the combined product/demo draft;
split into 006 product core and 007 demo. Focused re-review of the split is
pending.

## Review Runs

| plane | model | status | raw output |
| --- | --- | --- | --- |
| Product proof | `openrouter/qwen/qwen3.6-plus` | usable | `raw/2026-05-10-qwen-product-proof.md` |
| Demo truth and contamination | `openrouter/deepseek/deepseek-v4-pro` | usable | `raw/2026-05-10-deepseek-demo-truth.md` |
| GitHub evidence | `openrouter/xiaomi/mimo-v2.5-pro` | usable | `raw/2026-05-10-xiaomi-github-evidence.md` |
| Theater and scope | `zai/glm-5.1` | usable | `raw/2026-05-10-glm-theater-scope.md` |
| Fix-check: product proof | `openrouter/qwen/qwen3.6-plus` | usable for combined draft | terminal capture, 2026-05-10 |
| Fix-check: GitHub evidence | `openrouter/xiaomi/mimo-v2.5-pro` | usable for combined draft | terminal capture, 2026-05-10 |
| Fix-check: theater and scope | `openrouter/deepseek/deepseek-v4-pro` | usable for combined draft | terminal capture, 2026-05-10 |

`pi` emitted a startup warning that `kimi-coding/k2p6` did not match the local
model registry. Kimi was not assigned a reviewer slot in this cycle.

## Verdict

`REVIEW_REFRESH_REQUIRED_AFTER_SPLIT`

All reviewers agreed on the direction: continue in the existing
`sdp-trace-demo-jvm-gsd` repository as the default. The existing repo's messy
history is the product problem. A clean new repository would be useful later
for a polished public demo, but it should not replace the first proof that
`sdp-trace` can make real agent-delivery evidence legible.

The original combined draft was not approval-ready because it lacked a minimum bar
for the first CTO-visible packet, did not pin the negative theater target, and
left retroactive GitHub artifact/provenance rules too loose.

Focused re-review accepted the combined-draft fixes. After the split, the main
design change is stricter: 007 no longer allows hand-authored CTO-visible demo
packets because 006 now owns the pre-renderer fixture path. This split needs a
focused re-review before asking for approval.

Xiaomi raised three minor clarification items; all three were applied, with the
second item superseded by the 006/007 split:

- excluded contaminated features must not appear on CTO-visible packet/tracker
  surfaces unless the contamination is surfaced;
- if no existing feature meets the first-packet bar, 007 creates a new v2 PR
  only after the 006 renderer/validator exists;
- "exits `not_assessed`" is defined inline as assessed `pass`, `partial`, or
  `fail`.

## Dispositions

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| 006-P1 | major | product proof | Product tasks T009-T015 lacked P0 classification fields. | `accepted_fixed` | `tasks.md` now maps T009-T015 to packet rows, evidence surface, start state, transition, buyer effect, and non-goal. |
| 006-P2 | major | product proof | No row-state threshold for demo success. | `accepted_fixed` | `spec.md` now defines first-packet minimum bar: four rows pass/partial, `PC-CHANGE` and `PC-MUTATION` resolvable, one of verification/review/agent-route exits `not_assessed`. |
| 006-P3 | major | product proof | Negative demo reason code was ambiguous. | `accepted_fixed` | `spec.md` and `demo-repo-plan.md` now make `agent_claimed_verification` the primary first negative example. |
| 006-P4 | minor | product proof | Locally-generated packets could become permanent proof. | `accepted_fixed_superseded_by_split` | 006 owns pre-renderer fixture validation; 007 requires CTO-visible demo packets to be renderer-generated and validator-checked. |
| 006-P5 | minor | product proof | Evidence contradiction rule missing. | `accepted_fixed` | `spec.md` now marks contradictory surfaces as `partial` and requires both refs in residual gaps. |
| 006-P6 | minor | product proof | PR comment lifecycle/canonicality missing. | `accepted_fixed` | `plan.md` now says uploaded packet artifact is canonical and PR comment/body link is informational. |
| CT-001 | critical | demo truth | First packet target unspecified. | `accepted_fixed` | `spec.md` now has feature evidence inventory and first-packet selection/minimum-bar rules. |
| CT-002 | critical | demo truth | Feature 4 contamination not handled in spec. | `accepted_fixed` | `spec.md` now requires contamination surfaced in residual gaps/theater or explicit exclusion; Feature 4 cannot be first packet unless surfaced. |
| CT-003 | critical | demo truth | First packet could be honest but useless. | `accepted_fixed` | First-packet minimum bar added. If no existing feature meets it, first packet comes from a new v2 PR. |
| CT-004 | major | demo truth | Reset decision rule lacked "too noisy" threshold. | `accepted_fixed` | `plan.md` reset table switches to fresh root branch only if all existing features fail first-packet rule and minimum bar. |
| CT-005 | major | demo truth | Negative PR could poison first buyer impression. | `accepted_fixed` | `spec.md` and `demo-repo-plan.md` require `DEMO-NEGATIVE:` prefix, `demo-theater` label, separate handling, and showing it after happy-path packet. |
| CT-006 | minor | demo truth | Setup PR could green without proving artifact path. | `accepted_fixed` | `plan.md` requires setup-only artifact or packet validation command; setup is not feature proof. |
| CT-007 | minor | demo truth | Feature evidence inventory missing. | `accepted_fixed` | `spec.md` now includes initial inventory table. |
| G1 | critical | GitHub evidence | Retroactive task binding could backdate provenance. | `accepted_fixed` | `spec.md` adds provenance rules for PR body/task artifact binding and caps retroactive initiator confidence at `partial`. |
| G2 | critical | GitHub evidence | GitHub artifact retention/expiry ignored. | `accepted_fixed` | `spec.md` requires 180-day retention for v2 packet artifacts and `artifact_expired` handling for retroactive packets. |
| G3 | major | GitHub evidence | First packet selection criteria missing. | `accepted_fixed` | First-packet evidence-richness selection rule added. |
| G4 | major | GitHub evidence | Draft PR and retained fixture treated as equivalent negative demo. | `accepted_fixed` | Negative CTO demo must be a GitHub draft PR; fixtures are tests only. |
| G5 | major | GitHub evidence | PR-visible packet pointer optional. | `accepted_fixed` | `plan.md` requires PR comment or PR body link for demo PRs. |
| G6 | major | GitHub evidence | Historical CI citation vs fresh CI was ambiguous. | `accepted_fixed` | `spec.md` now distinguishes retroactive artifact availability and v2 retention. |
| G7 | major | GitHub evidence | Setup/first packet ordering ambiguous. | `accepted_fixed` | `plan.md` clarifies setup PR is infrastructure only; first packet follows selection rule. |
| G8 | minor | GitHub evidence | Packet/bundle path convention unclear. | `accepted_fixed` | `demo-repo-plan.md` defines flat packet files and directory-per-feature bundles. |
| G9 | minor | GitHub evidence | Hand-authored/tool-generated transition unclear. | `accepted_fixed_superseded_by_split` | 006 owns pre-renderer fixtures; 007 depends on 006 renderer/validator. |
| G10 | minor | GitHub evidence | "generate" overclaimed before renderer exists. | `accepted_fixed_superseded_by_split` | 007 starts after renderer completion, so CTO-visible demo packets are tool-generated. |
| SR-001 | major | theater/scope | Core claim implied broad OSS support. | `accepted_fixed` | `spec.md` now says one demonstrated route: OpenCode + GSD + MiniMax-M2.5. |
| SR-002 | major | theater/scope | Negative reason code ambiguous. | `accepted_fixed` | Primary negative reason code fixed to `agent_claimed_verification`. |
| SR-003 | major | theater/scope | No happy-path clean theater criterion. | `accepted_fixed` | Added SC-006 requiring at least one happy-path `PC-THEATER: pass`. |
| SR-004 | major | theater/scope | Product prerequisites not linked to generated proof. | `accepted_fixed_superseded_by_split` | 006 is now a separate prerequisite spec; 007 starts only after 006 renderer/validator behavior exists. |
| SR-005 | minor | theater/scope | Todo app implied enterprise ceremony. | `accepted_fixed` | `demo-repo-plan.md` says risk/security owners may be `not_in_scope`. |
| SR-006 | minor | theater/scope | Sales-demo wording drift. | `accepted_narrower` | Existing wording keeps "polished public sales demo" only as later non-default option; current demo remains CTO product proof. |
| SR-007 | minor | theater/scope | Task dependency edges missing. | `accepted_fixed_superseded_by_split` | 007 tasks now begin by confirming 006 availability and no longer contain product implementation tasks. |

## Remaining Risk

No known critical or major finding remains intentionally open from the combined
draft review. The split itself still needs focused re-review because it changes
sequencing and removes the pre-renderer hand-authored demo path.

## 2026-05-11 Operating Model Correction Review

Scope:

- Codex role boundary for demo feature work.
- OpenCode + GSD + `minimax-coding-plan/MiniMax-M2.5` as the CTO-visible
  feature implementation route.
- `sdp-trace` as passive flight recorder outside the developer prompt.
- Setup-only Codex exception.
- P0/P1+ product blocker handling.

Review runs:

| plane | reviewer/runtime | status | raw output |
| --- | --- | --- | --- |
| requirements/evidence/security | Codex subagent `019e17af-e8b4-7833-9311-94b042d937b5` | usable | notification captured in thread |
| focused Socratic | `pi` with `openrouter/deepseek/deepseek-v4-pro` | usable | `raw/2026-05-11-operating-model-deepseek.md` |
| focused fix-check | `pi` with `openrouter/deepseek/deepseek-v4-pro` | usable for first accepted fixes | `raw/2026-05-11-operating-model-fixcheck-deepseek.md` |
| replacement fix-check | Codex subagent `019e17b7-a8f7-7360-9336-e90bc1d2eb9c` | usable | notification captured in thread |
| implementation review | `pi` with `deepseek/deepseek-chat` | usable only as secondary signal; line references were imprecise | local scratch `/tmp/sdp-trace-review/packet-check-demo-fixcheck.md` |
| implementation review | Codex subagent `019e1821-a351-71d0-a0ee-91f089482f46` | usable | notification captured in thread |
| implementation review | Codex subagent `019e1825-4cd5-7ec0-b4e2-2e6284ef5588` | usable | notification captured in thread |
| implementation fix-check | Codex subagent `019e1828-1eca-7f31-943b-b7bd1ab0e583` | usable | notification captured in thread |

Findings and dispositions:

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| OM-001 | critical | requirements / route proof | Demo success could still pass without proving the OpenCode/GSD/MiniMax route because the first-packet bar allowed `PC-AGENT-ROUTE` to remain `not_assessed`/`cannot_verify`. | `accepted_fixed` | `spec.md` now requires first-packet `PC-AGENT-ROUTE` to be `pass` or `partial` with recorder-observed OpenCode/GSD/MiniMax evidence; SC-009 repeats the success gate. |
| OM-002 | major | provenance / contamination | Existing/backfilled feature packets lacked a Codex-authored feature contamination audit. | `accepted_fixed` | `spec.md` FR-015, `tasks.md` T014, and `demo-repo-plan.md` require the audit before using history as CTO-visible route proof. |
| OM-003 | major | blocker handling | `demo-repo-plan.md` allowed continuing with `cannot_verify` after a P0 route/provenance blocker. | `accepted_fixed` | `spec.md`, `plan.md`, and `demo-repo-plan.md` now block feature proof until P0 route/provenance blockers are fixed and rerun or otherwise observed. |
| OM-004 | major | prompt boundary | Prompt separation was only required for the first selected feature, not every feature. | `accepted_fixed` | `spec.md` FR-016/SC-008 and `tasks.md` T015-T019 require retained prompt text or digest validation for every feature. |
| OM-005 | critical | observation integrity | Checked-in recorder JSON could be mistaken for authority without live validation, CI retention, signing, or timestamping. | `accepted_fixed_narrowed` | `spec.md` and `demo-repo-plan.md` now define `local_observed` trust scope, require live-validated or CI-retained recorder artifacts with resolver refs and digests, and prohibit audit-grade integrity claims without stronger evidence. |
| OM-006 | critical | recorder passivity | Recorder passivity was asserted without machine-verifiable proof. | `accepted_fixed_narrowed` | `spec.md` and `tasks.md` now say prompt/setup metadata can support only `PC-AGENT-ROUTE: partial`; `pass` requires stronger machine evidence, and residual gaps must name unverified passivity. |
| OM-007 | critical | success gate | First-packet minimum bar was not explicitly gated before buyer rehearsal. | `accepted_fixed` | `tasks.md` T021 adds a first-packet gate before rehearsal; `spec.md` FR-019/SC-010 require the gate. |
| OM-008 | major | theater evidence | Negative theater depended on an undefined automatic assessor. | `accepted_fixed_narrowed` | `plan.md` and `tasks.md` require 006-generated or 006-validated theater findings; if 006 validates a supplied finding rather than detecting it automatically, the packet must state that limitation. |
| OM-009 | major | setup-only boundary | Setup-only Codex changes were prose labels rather than checkable boundaries. | `accepted_fixed` | `spec.md` FR-018, `tasks.md` T012, and `demo-repo-plan.md` define machine-checkable or independently reviewed setup-only scope. |
| OM-010 | major | implementation / route evidence | `packet check-demo` could accept generic harness evidence if the manifest resolver merely mentioned OpenCode, GSD, and MiniMax. | `accepted_fixed` | `internal/packet/packet.go` now requires `evidence_kind: harness_route_observation` and `observed_components` covering OpenCode, GSD, and MiniMax; `schema/evidence-bundle-manifest.v0.schema.json` and `cmd/sdp-trace/packet_cli_test.go` cover the structured fields. |
| OM-011 | major | implementation / evidence freshness | `packet check-demo` could accept expired `PC-CHANGE` or `PC-MUTATION` refs for `partial` rows when expiry was represented only by `expires_at`. | `accepted_fixed` | `internal/packet/packet.go` passes `now` into `demoUsableEntry` and rejects `entryExpired`; `cmd/sdp-trace/packet_cli_test.go` covers an expired partial `PC-CHANGE` ref. |

Current result: focused fix-check reported no critical or major findings for
OM-001 through OM-004. A `pi` final fix-check for OM-005 through OM-009 hung and
returned empty output; it is unusable and not counted as evidence. Replacement
subagent fix-check reported `NO CRITICAL OR MAJOR FINDINGS` for OM-005 through
OM-009.

Implementation fix-check for the 007 first-packet gate reported
`NO CRITICAL OR MAJOR FINDINGS` after OM-010 and OM-011 were fixed. Local live
verification passed with `go test ./...`, `jq empty schema/*.json`,
`git diff --check`, `sdp-trace packet validate --bundle
examples/change-evidence-packet/happy-path.bundle.json`, and `sdp-trace packet
check-demo --bundle examples/change-evidence-packet/happy-path.bundle.json`.
CI, external witness, signed attestation, and production trust remain
`not_assessed` for this local dirty worktree.

## 2026-05-11 Demo Repository Audit

Target: `/home/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd`.
Branch: `demo-v2-packetization`, not `main`.

Observed dirty files:

- isolation/evidence setup: `.gitignore`, `.ignore`, `.opencode/opencode.json`,
  `.bazelversion`, `.evidence/local-build-test/...`;
- planning/docs: `README.md`, `.planning/phases/01-project-skeleton/*.md`;
- build files: `MODULE.bazel`, `MODULE.bazel.lock`, `WORKSPACE`,
  `app/BUILD.bazel`.

Local verification:

- `bazel build //...`: pass.
- `bazel test //... --test_output=errors`: pass.
- `git diff --check`: pass.
- Agent-visible scan across `README.md`, `.planning`, `.github`, Bazel files,
  `app`, `.opencode`, `.ignore`, and `.gitignore`: no `sdp-trace` or
  `sdp_trace` matches.

Disposition:

- `.ignore` plus `.opencode/opencode.json` is consistent with the current
  isolation direction, but the residual gap remains open until `sdp-trace`
  installs/verifies the isolation contract itself.
- `.evidence/local-build-test/manifest.json` correctly records local dirty
  scope and `cannot_verify` for CI witness, audit-grade gate, and
  `RG-OPENCODE-CONTEXT-IGNORE-001`.
- Build-file changes are outside the default setup-only allowlist. They appear
  to be build-system repair/setup rather than application feature behavior, but
  they require an independent setup-only review record before being treated as
  setup evidence. They must not close CTO-visible feature packet rows.
- No Codex-authored demo application feature behavior is accepted from this
  audit. Any feature repair must go through OpenCode/GSD + MiniMax under passive
  `sdp-trace` observation.

## 2026-05-11 Feature 1 Readiness Slice

Feature: `GET /ready` in
`/home/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd`.

Route and evidence:

- Initial `codex-subagent run opencode` route hung and was cancelled. It is
  recorded as friction and not counted as delivery evidence.
- Direct OpenCode route with `opencode/minimax-m2.5-free` completed under
  passive recorder at
  `.evidence/feature-readiness-opencode-direct-free/run`.
- OpenCode review-fix route tightened `/ready` smoke coverage at
  `.evidence/feature-readiness-review-fix/run`.
- Final local verification ran under recorder at
  `.evidence/feature-readiness-final-verification/run` and produced
  `local_gate: pass`, `ci_witness_gate: cannot_verify`, and
  `audit_grade_gate: cannot_verify`.
- The product builder generated
  `.sdp-trace/bundles/feature-1/bundle.json` and
  `.sdp-trace/packets/feature-1.md`.
- `sdp-trace packet validate --bundle .sdp-trace/bundles/feature-1/bundle.json`:
  pass.
- `sdp-trace packet check-demo --bundle .sdp-trace/bundles/feature-1/bundle.json`:
  pass.

Disposition:

- `PC-AGENT-ROUTE` is `partial`, not `pass`, because route proof is local-only
  and prompt body retention is digest-only.
- CI witness, signed attestation, authority, and audit-grade trust remain
  `not_assessed` or `cannot_verify`; they are not claimed.
- The feature slice demonstrates `sdp-trace` as local flight recorder plus packet
  generator/checker for an OpenCode/MiniMax demo route, not merge approval or
  production trust.

## Current Recommended First Slice

1. Complete 006 Change Evidence Packet Core first.
2. Keep the existing `sdp-trace-demo-jvm-gsd` repository.
3. Tag current `main` as `demo-v1-observation-baseline`.
4. Add packetization setup PR with artifact path and tracker only.
5. Select the first packet target by evidence inventory and minimum bar.
6. If no existing feature meets the bar, create one new v2 feature PR under the
   packetization track.
7. Produce one CTO-readable Change Evidence Packet v0 before attempting the
   full five-feature demo.

## 2026-05-26 Focused Split Re-Review

Reviewer: Codex GPT-5, sdp-trace closure route
Scope: focused review of the split 007 spec package after 006 implementation
and closure-route reconciliation.

Files reviewed:

- `spec.md`
- `plan.md`
- `demo-repo-plan.md`
- `tasks.md`
- this review ledger

Verdict: `LGTM_FOR_SPEC_REVIEW_SCOPE`

The split-specific concern recorded on 2026-05-10 is addressed for the spec
package: 007 no longer allows hand-authored CTO-visible packets, requires 006
renderer/validator output for demo proof, preserves the OpenCode/GSD/MiniMax
route boundary, and keeps Codex-authored feature behavior out of CTO-visible
route proof.

No critical or major spec-review finding remains open for Phase 1.

Dispositions:

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| SPLIT-001 | major | product proof | 007 needed focused re-review after 006/007 split removed hand-authored CTO-visible packets. | `accepted_fixed` | `spec.md` and `plan.md` require 006-generated and 006-validated packets for demo proof. |
| SPLIT-002 | major | demo truth | Demo implementation must not start merely because 006 artifacts exist locally; explicit approval and demo-repo evidence are still required. | `accepted_boundary` | `tasks.md` keeps T008 and Phase 2 tasks open. |
| SPLIT-003 | major | evidence boundary | Checked-in recorder JSON or historical review prose could be mistaken for demo authority. | `accepted_fixed` | `spec.md` and `demo-repo-plan.md` require live validation or CI-retained resolver/digest binding, with `partial` / `cannot_verify` for weaker evidence. |

Current blocked states:

- T008 remains open: explicit approval of demo track and first implementation
  slice is not represented.
- T009-T022 remain open: the demo repository is not present in this worktree
  environment, so demo-environment availability, tagging, PR setup, feature
  packets, negative theater PR, and buyer rehearsal are `not_assessed`.
- External CI, signed attestation, and production trust remain `not_assessed`.

Verification:

- inspected the current 007 spec package and existing review ledger;
- inspected current task state;
- checked local demo repository availability at
  `/home/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd`: not present
  in this environment.
