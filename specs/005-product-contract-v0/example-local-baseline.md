# Change Evidence Packet v0 Example - Local Enterprise Baseline

Status: example-only, not product proof
Contract: Product Contract v0
Profile: `local-enterprise-baseline-v0`

This example demonstrates the minimum useful packet when a Russian enterprise
pilot has local Git and internal CI artifacts, but no GitHub API, no
OpenCode/GSD or other harness evidence, no raw prompt export, and no public
signing.

## Executive Summary

Change `INV-77` modifies invoice export validation. The change is visible from
local Git refs and an internal TeamCity build artifact. No agent route is
available, so agent attribution remains `not_assessed`. CI verification is
partially witnessed by TeamCity artifact metadata, but the packet has no review,
authority envelope, or signed witness. The next merge decision owner is bound to
the internal `billing-maintainers` role from a local policy ref.

This packet is useful because it shows what the customer can verify in a closed
contour and what remains missing. It does not claim agent observability.

## Packet Metadata

| field | value |
| --- | --- |
| packet_id | `cep-2026-05-10-inv-77-local-example` |
| schema | `change-evidence-packet-v0` |
| generated_from | example fixture |
| selected_profile | `local-enterprise-baseline-v0` |
| redaction_policy | no prompt/session data supplied |
| bundle_ref | `bundle:inv-77-local-example` |
| packet_state | `draft` |

## Required Rows

| row id | state | answer | evidence refs | gap / next evidence |
| --- | --- | --- | --- | --- |
| `PC-CHANGE` | `pass` | Local Git commit range `111aaa..222bbb` changes invoice export validation. | `git:111aaa..222bbb`, `artifact:diff-digest:sha256:local-example` | Change-host PR/MR id is `not_assessed`. |
| `PC-INITIATOR` | `not_assessed` | No issue, task, prompt boundary, or initiating actor is bound. | none | Need task-system ref, issue ref, prompt-boundary digest, or external assertion. |
| `PC-AGENT-ROUTE` | `not_assessed` | No agent or harness evidence was supplied. | none | Need harness/session evidence or explicit statement that change was human-only. |
| `PC-MUTATION` | `pass` | File mutation is observed in local Git. | `git:111aaa..222bbb` | Actor/tool/model attribution remains `not_assessed`. |
| `PC-VERIFICATION` | `partial` | TeamCity build artifact exists for the commit range, but test coverage details are not retained in the packet. | `ci:teamcity:build-5482`, `artifact:teamcity-summary:sha256:local-example` | Need retained test report or customer witness evidence for full verification. |
| `PC-REVIEW` | `not_assessed` | No review artifact is retained. | none | Need review record, reviewer plane, identity/source class, and retained result. |
| `PC-AUTHORITY` | `not_assessed` | No selected authority envelope was supplied. | none | Need local policy envelope and selected `policy_id`. |
| `PC-THEATER` | `partial` | Two P0 theater findings are triggered: `unbound_intent`, `ci_theater`. | `theater:unbound_intent`, `theater:ci_theater` | Need source intent binding and retained CI coverage artifact. |
| `PC-ATTESTATION` | `not_assessed` | No signed packet or private customer witness is present. | none | Need private PKI witness or signed packet digest. |
| `PC-DECISION` | `partial` | Merge owner is bound to `billing-maintainers`; release, risk, and security owners are not assessed. | `policy:merge-owner:billing-maintainers` | Need release/risk/security owner policy refs when those decisions are in scope. |
| `PC-RESIDUAL-GAPS` | `pass` | Packet records missing initiator, agent route, review, authority, signing, and CI coverage details. | this packet | `pass` means all non-pass rows and active findings are enumerated per the `PC-RESIDUAL-GAPS` synthesis rule. |

## Theater Findings

| reason code | state | severity | finding | trigger evidence | required closure evidence |
| --- | --- | --- | --- | --- | --- |
| `unbound_intent` | `partial` | major | Change exists, but source task/issue/prompt boundary is not bound. | `git:111aaa..222bbb`, missing `PC-INITIATOR` evidence | Task-system ref, issue ref, or prompt-boundary digest. |
| `ci_theater` | `partial` | major | CI artifact exists, but detailed retained coverage for the selected verification claim is missing. | `ci:teamcity:build-5482`, missing retained test report | Retained TeamCity test report or customer witness. |

## Decision Ownership

| decision | owner state | owner | why |
| --- | --- | --- | --- |
| merge readiness | `owner_bound` | `billing-maintainers` | Bound to local merge-owner policy ref. Alternative: `git:signed-off-by:tech-lead@internal.example` for teams without formal policy refs. |
| release readiness | `not_assessed` | none | No release policy selected. |
| risk acceptance | `not_assessed` | none | No risk policy selected. |
| security review | `not_assessed` | none | No security review policy selected. |

## Evidence Bundle

| ref | retained form | notes |
| --- | --- | --- |
| `git:111aaa..222bbb` | source refs | Local Git evidence. |
| `artifact:diff-digest:sha256:local-example` | digest | Diff body not retained in this example. |
| `ci:teamcity:build-5482` | external_ref | Internal CI ref. |
| `artifact:teamcity-summary:sha256:local-example` | redacted | Build summary retained, detailed test report absent. |
| `policy:merge-owner:billing-maintainers` | external_ref | Local policy ref for merge ownership. |

## What This Packet Does Not Prove

- It does not prove the change was produced by or through an agent.
- It does not prove full test coverage.
- It does not approve merge or release.
- It does not prove authority compliance.
- It does not prove signed attestation or production trust.
- It does not require raw prompt export or public SaaS access.
