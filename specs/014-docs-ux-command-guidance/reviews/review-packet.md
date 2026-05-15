# Review Packet: 014-docs-ux-command-guidance Spec Review

## Objective
Adversarial review of the 014 spec and the current docs it intends to improve. Focus on misleading state language, output confusion, profile selection, overclaim prevention, and UX gaps for a cold user.

## Artifact (what is being reviewed)

--- BEGIN SPEC ---
# Feature Specification: Docs UX And Command Guidance

**Feature Branch**: `014-docs-ux-command-guidance`
**Created**: 2026-05-14
**Status**: Draft for PI review
**Input**: UX review found command discoverability, state vocabulary, output-location, profile-selection, and overclaim guidance gaps.

## Problem Statement

The CLI and docs are accurate but high-friction. Users see many commands, many evidence states, several scope vocabularies, and multiple output locations. Critical overclaim warnings are repeated across documents rather than presented as one decision aid.

## Core Claim

This slice may claim:

> The user-facing docs provide a guided path for choosing commands, interpreting evidence states, locating outputs, selecting assessment profiles, and avoiding overclaim.

This slice must not claim:

- new verifier behavior unless implemented separately;
- interactive CLI support unless a command is actually added;
- production trust or authority decisions.

## Required User Stories

### US-001 - Command Choice By Task (P0)

A user can choose the next command from a task-oriented guide rather than reading a long flat reference.

**Independent Test**: Reviewer docs answer "I have a run directory, what now?" and "Which assessment profile applies?" without requiring the full command table.

### US-002 - Evidence State Decision Tree (P0)

A reviewer can distinguish `not_assessed`, `cannot_verify`, `missing_telemetry`, `observed`, `pass`, and `fail` without guessing from exit codes.

**Independent Test**: Docs contain one canonical state table or decision tree, and other docs reference it instead of redefining states inconsistently.

### US-003 - Output Location Map (P1)

A user can tell where each command writes artifacts and what each artifact is for.

**Independent Test**: Docs include an output map for run dirs, reports, query packs, witness outputs, assessment outputs, and release proof outputs.

### US-004 - Overclaim Checklist (P0)

A reviewer has one canonical checklist for forbidden interpretations and trust-scope escalation.

**Independent Test**: Reviewer entrypoint contains the canonical checklist; README and agent entrypoint link to it.

## Functional Requirements

- **FR-001**: Add a task-oriented command guide or restructure existing docs to include one.
- **FR-002**: Add canonical evidence state vocabulary and exit-code mapping.
- **FR-003**: Add "which profile do I use?" decision tree.
- **FR-004**: Add output location reference.
- **FR-005**: Consolidate overclaim rules and link from other docs.
- **FR-006**: Mark any future interactive guide command as a separate implementation follow-up, not implied by docs-only work.

## Acceptance Criteria

- `docs/reviewer-entrypoint.md` has a short task path before long references.
- State vocabulary is consistent across README, concepts, agent entrypoint, reviewer entrypoint, and adoption guide.
- `go run ./tools/doccheck` passes and covers the command claims it owns.
- UX review finds no blocker-level ambiguity in the first-run reviewer path.

## PI Review Prompt

Review whether the proposed docs UX makes command choice and evidence interpretation safer for a cold user. Focus on misleading state language, output confusion, profile selection, and overclaim prevention.

--- END SPEC ---

--- BEGIN PLAN ---
# Implementation Plan: Docs UX And Command Guidance

## Technical Context

**Language**: Markdown; optional small Go doccheck extension
**Dependencies**: Existing docs and command-surface output
**Verification**: doccheck, docs grep for deprecated state terms, cold-reader review, `git diff --check`

## Scope

- Reviewer and operator docs only unless a small doccheck extension is required.
- No CLI behavior changes in the first slice.

## Non-Goals

- Adding an interactive guide command.
- Rewriting all docs.
- Changing exit-code semantics.

## Risks

- Consolidation can hide important caveats if links replace necessary local warnings.
- State vocabulary must reflect actual CLI behavior, not desired behavior.

## Review Plan

Run UX, evidence, and requirements review planes. Verify every state and profile name against current command-surface/docs.

--- END PLAN ---

--- BEGIN TASKS ---
# Tasks: Docs UX And Command Guidance

**Input**: `spec.md`, `plan.md`
**Tests**: doccheck, state-term grep, cold-reader review, `git diff --check`.

## Phase 0 - Review

- [ ] T001 Run PI UX review on current docs flow.
- [ ] T002 Verify current state/profile names against command docs.
- [ ] T003 Record accepted/rejected findings.

## Phase 1 - Canonical Guidance

- [ ] T010 Add task-oriented command guide.
- [ ] T011 Add canonical evidence state and exit-code table.
- [ ] T012 Add output location map.
- [ ] T013 Add profile decision tree.
- [ ] T014 Consolidate overclaim checklist.

## Phase 2 - Drift Control

- [ ] T020 Remove or link duplicate overclaim/state prose.
- [ ] T021 Extend doccheck if command/state claims remain duplicated.
- [ ] T022 Verify docs do not introduce new authority claims.

## Phase 3 - Closure

- [ ] T030 Run doccheck and full verification.
- [ ] T031 Run cold-reader UX review.
- [ ] T032 Record remaining advisory UX follow-ups separately.

--- END TASKS ---

--- BEGIN REVIEWER-ENTRYPOINT ---
# Reviewer Entrypoint

Use this path for a first-time reviewer check in under five minutes. For the
full bilingual command/profile surface, see `docs/agent-entrypoint.md` and
`sdp-trace --help`.

For the demo-repository pilot evidence package, read
`examples/pilot-runs/opencode-minimax-kotlin-bazel/README.md` before inspecting
the retained package. Treat that package as an exact observed slice, not broad
OpenCode, MiniMax, Kotlin, or Bazel support.

## Verification Path

From a clean checkout, run:

1. `go test -count=1 ./...`
2. For a source checkout, define `sdp_trace() { go run ./cmd/sdp-trace "$@"; }`.
   After installing a release binary, use `sdp-trace` directly.
3. `sdp_trace --help` for a source checkout, or `sdp-trace --help` for a release binary.
4. `sdp_trace validate-fixtures examples/agentic-sdlc` for a source checkout, or `sdp-trace validate-fixtures examples/agentic-sdlc` for a release binary.
5. Create or inspect a run with
   `sdp_trace wrap --name smoke --output-dir .sdp-trace-runs/smoke -- /bin/echo ok`
   for a source checkout, or the same command with `sdp-trace` after installing
   a release binary.
6. Verify that run with `sdp_trace verify .sdp-trace-runs/smoke` or
   `sdp-trace verify .sdp-trace-runs/smoke`.
7. If documentation changed, compare command examples against `sdp_trace --help`
   or `sdp-trace --help`.

External production trust is not part of this quick path. Use a live
`external_production_trust` profile path before making production trust claims.

## Exit Code Contract

Use `docs/agent-entrypoint.md` as the canonical state, trust-scope, authority
scope, and exit-code contract. The short exit summary is:

- `0`: `observed`, `pass`, or explicitly scoped `not_assessed`
- `1`: `fail`
- `2`: usage error / invalid command invocation
- `3`: `cannot_verify`

If any command returns exit code `3`, inspect the emitted reason and do not
upgrade the state in prose.

## Reviewer Command Surface

This is the reviewer subset for fast orientation. The full command surface is
authoritative in [Agent Entrypoint](agent-entrypoint.md) and `sdp-trace --help`.
When reviewing command docs, compare against both.

- `version`, `wrap`, `run`, `dry-run`, `preview`, `doctor`
- `command-surface`
- `install repo-observer`
- `interaction relay`, `interaction import-transcript`, `interaction summarize`
- `observe setup`, `observe collect`, `observe session`
- `harness observe`, `harness validate`, `harness summarize`
- `envelope summarize`
- `verify`, `explain`, `query`
- `query-pack`, `query-pack explain`
- `export cross-repo-posture`, `export cross-repo-posture explain`, `export telemetry`
- `assess`, `assess preview`, `assess explain`
- `report`, `gate`, `witness`, `release-proof`, `pr-review`
- `packet build-pr`, `packet build-github`, `packet validate`, `packet check-demo`, `packet render`
- `validate-fixtures`

Current assessment profiles:

- `adapter-capture`
- `managed-harness`
- `forensic-retention`
- `ci-artifact-observation`
- `authority-envelope`

Current witness kinds:

- `github-actions`
- `gitlab-ci`
- `buildkite`
- `customer-pki`

Air-gapped evidence is represented through customer policy/private-equivalent
guidance and fixtures, not as a separate `witness --kind` value.

Harness observation commands import and validate explicit local harness event
exports. They do not run OpenCode, GSD, MiniMax, GitHub, provider APIs, or any
other harness. Treat missing harness event families as `not_assessed` or
`cannot_verify`, not as feature delivery evidence.

First-run observation commands use a session profile to bound setup and
collection. `observe setup` writes setup metadata before delivery,
`observe collect` normalizes declared harness output after the normal harness
command, and `observe session` is a convenience wrapper for one controlled
command. They do not inject stdin, relay prompts, retain stdout/stderr bodies by
default, or turn missing harness output into a pass.

## Dirty Checkout Behavior

- Clean checkout: verifier trust scope is the stated profile (`repo_baseline_structural`, `source_bound_local_release`, or `external_production_trust`).
- Dirty checkout without a command-supported dirty allowance: required clean-source checks may return `cannot_verify`.
- Dirty structural output may support only the `local_dirty_structural_only`
  authority scope.
- Do not use dirty output to conclude `source_bound_local_release` or `external_production_trust`.

## Not-Assessed Interpretation

`not_assessed` means the selected run did not assess that state.

What it allows:

- Continue using the command output with that state held back.
- Ask for the missing evidence or rerun against a scope that can assess it.

What it does not allow:

- Treating the state as passed.
- Using it as external trust closure.

## Gate, Witness, And Release Caveats

- `pr-review` emits review-record evidence over a frozen PR packet. It can
  report `coverage_satisfied`, `coverage_partial`, `coverage_unresolved`,
  `not_assessed`, or `cannot_verify`, but it does not approve, merge, mark
  ready, release, accept risk, or replace human approval.
- `gate` emits verifier-derived facts and deterministic states. It does not own
  merge, release, readiness, degradation, override approval, or risk acceptance.
- `witness` binds available CI or customer-PKI evidence. A CI witness file is
  not external production trust, a transparency log, or a release approval by
  itself.
- `release-proof` can establish `source_bound_local_release` only when the
  source commit and manifest subjects match. It does not prove
  `external_production_trust`, `trusted_contract_release`, or
  `production_release_verified`.

## What You May State From Output

From verifier results, you may only state:

- Which command/profile was run.
- Which `result` or state values were produced.
- Whether the selected profile concluded with live `pass` or `observed`.
- Which states remain `not_assessed` or `cannot_verify`, with the emitted reason.

You may not state external production trust guarantees until
`external_production_trust` completes with live `pass` and
`production_release_verified` is supported by its dependent evidence chain.

## Quick Reference

| Goal | Command | Typical state boundary |
| --- | --- | --- |
| Local trace verification | `sdp-trace verify <run-dir>` | `observed` supports local structural assertions only |
| Missing evidence review | `sdp-trace query --query missing-evidence <run-dir>` | Missing evidence remains visible, not passed |
| Forensic package review | `sdp-trace query-pack explain --result <file>` | Explanation of retained evidence only |
| Managed harness review | `sdp-trace assess explain --assessment-result <file>` | Assessment facts; external policy owns block/allow |
| First-run harness observation review | `sdp-trace observe collect --profile <session-profile.json> --run <run-dir>` | Session-profile collection; missing declared output is `cannot_verify` |
| Harness event validation | `sdp-trace harness validate --profile <harness-profile.json> --run <run-dir> --out <file>` | Event-family facts; missing required families are not passes |
| Authority envelope review | `sdp-trace assess --profile authority-envelope --authority-package <file> --out <file>` | Authority facts only; external policy owns consequences |
| CI/customer witness review | `sdp-trace witness --kind <kind> --out <file> <runs-root-or-run-dir>` | CI/customer-bound evidence, not production trust by itself |
| Source-bound release review | `sdp-trace release-proof --manifest <file> --out <file>` | Local source-bound proof only |
| Automated PR review evidence | `sdp-trace pr-review check --out review --repo-id <safe-id> --change-ref pr-123 --base <sha> --head <sha> --diff change.diff --profile examples/pr-review/trust-sensitive-default.profile.json` | Review-record completeness only; not merge approval |

## Manual External PR Review Handoff

For `manual_external` PR review planes, a usable `findings_reported` or
`no_findings` status requires retained reviewer output. A bare PR comment or
hand-edited status is not enough.

Reviewer output must be JSON matching `schema/pr-review-result.schema.json` and
must echo the packet digest, plane, and role. The review runner records the raw
output digest as `raw_output_ref`; validation counts the plane only after that
digest-bound output exists.

Minimum handoff steps:

1. Build or reuse a frozen packet directory with `packet/packet.json`.
2. Give the reviewer the packet digest, plane, role id, diff ref, context refs,
   and validation criteria.
3. Store the reviewer JSON output in a file outside the packet directory.
4. Use a profile role whose `command` prints that JSON file, then run
   `sdp-trace pr-review run --packet <packet-dir> --profile <profile> --out <runs-dir>`.
5. Run `sdp-trace pr-review synthesize`, `validate`, and `summarize` against the
   resulting run set and ledger.

If the reviewer output is absent, empty, off-task, malformed, lacks retained raw
output, or targets a different packet digest, record the plane as
`not_assessed` or `cannot_verify`. Do not treat it as sign-off.

This entrypoint is intentionally minimal and is intended to prevent over-claiming
from reproducible verifier output.

--- END REVIEWER-ENTRYPOINT ---

--- BEGIN AGENT-ENTRYPOINT (relevant sections) ---
## State And Exit Code Contract

- `0`: `observed`, `pass`, or explicitly scoped `not_assessed`
- `1`: `fail`
- `2`: usage error / invalid command invocation
- `3`: `cannot_verify`

- `observed`: verifier evidence met required checks for the selected local profile.
- `pass`: selected profile concluded successfully where the command contract uses pass/fail states.
- `fail`: verifier evidence conflicted or was insufficient for required checks.
- `not_assessed`: state was outside the run scope; it does not imply success or evidence.
- `cannot_verify`: verifier could not execute a required check or lacked required evidence.

A checked-in proof JSON is an audit artifact, not authority. Authority is replayed
only from live Go verifier output and the canonical command/state contract above.

## Air-Gapped Fixture Guidance

Air-gapped evidence is a fixture and customer-policy pattern, not a native
`witness --kind air-gapped` command. Use `customer-pki` or an accepted private
equivalent with explicit authority policy, payload digest, freshness evidence,
and retained audit references. If those are absent, record `not_assessed` or
`cannot_verify`; do not claim external production trust.

## Forbidden Claims

Do not emit these in this repo surface:

- `external_production_trust=true` without a live `external_production_trust` profile pass.
- `trusted_contract_release=true` without live external trust closure.
- `production_release_verified=true` outside a concluded `external_production_trust` run.
- Claims that treat `repo_baseline_structural` or `source_bound_local_release` outputs as production trust.
- Dirty-checkout structural output as source-bound or external-trust evidence.
## Profile Selection

Each assertion is anchored to one of these profile IDs:

- `repo_baseline_structural`
- `source_bound_local_release`
- `external_production_trust`

Do not infer profile from role. Choose the profile directly from the claim:

- `repo_baseline_structural`: structural command, fixture, and local trace integrity.
- `source_bound_local_release`: local manifest, source commit, artifact digest, and DSSE/source-bound checks.
- `external_production_trust`: external identity, protected source, transparency or customer audit evidence, approval, and production release verification.

Dirty-checkout baseline output is only valid under the
`local_dirty_structural_only` authority scope. It is not a profile ID and must
not be used to close `source_bound_local_release` or
`external_production_trust`.

## Result, Trust Scope, And Authority Scope

Keep these vocabularies separate:

- Result state: the verifier outcome for a selected command or profile, such as
  `observed`, `pass`, `fail`, `not_assessed`, or `cannot_verify`.
- Trust scope: the evidence boundary recorded by a run, witness, or assessment,
  such as `local_observed` or `ci_witnessed`.
- Authority scope: the reporting boundary for a committed package, such as
  `demo_pilot_only`.

Known trust scopes used by the current pilot docs:

- `local_observed`: local run/report evidence was captured and checked, but it
  is not CI-witnessed or external production trust.
- `ci_witnessed`: available CI identity and artifact binding evidence supported
  the witness profile for the exact CI topology under assessment.
- `external_witnessed`: external witness evidence was supplied and accepted by
  the selected profile.

Known authority scopes used by current docs:

--- END AGENT-ENTRYPOINT ---

--- BEGIN README (relevant sections) ---
## Start Here

1. Read [Install](docs/install.md) and choose either a release binary or
   source checkout command path.
2. Read [Core Concepts](docs/concepts.md) to understand the contract:
   spec, plan, task, evidence, gate, decision, trace, and provenance.
3. Give [Agent Onboarding](docs/agent-onboarding.md) to any coding agent before
   it works in this repository.
4. Follow [Contributor Quick Start](docs/contributor-quickstart.md) to run the
   canonical local smoke path and verify your environment.
5. Use [Agent Entrypoint](docs/agent-entrypoint.md) for the authoritative
   command and state contract.
6. Use [Reviewer Entrypoint](docs/reviewer-entrypoint.md) for a five-minute
   verification path and overclaim checklist.
7. Use [Documentation Map](docs/README.md) to choose the right next document.

Origin note: `sdp-trace` was extracted from delivery evidence work in
`sdp_lab`. That history is not a runtime dependency and should not be required
context for using this repository.

## What It Produces

- trace run artifacts under `.sdp-trace-runs/`;
- report packages under `.sdp-trace-report/`;
- query and forensic query-pack outputs;
- assessment results for supported profiles;
- advisory gate facts for downstream policy consumers;
- CI or customer witness artifacts when required evidence exists;
- source-bound local release proof when manifest subjects match the source
  commit.


--- END README ---


## Contract (rules the artifact must satisfy)

## Trust Rules
The repository has already failed once by letting prose, task checkboxes, and checked-in JSON overclaim what current verification could not replay. Do not repeat that failure.

- Machine proof wins over prose, task checkboxes, reports, review ledgers, and mirrors.
- Dirty checkout output is local structural evidence only, not external trust.
- Checked-in proof JSON is not authority unless live-verified or externally signed.
- No deferred trust closure: missing external evidence keeps the block open.
- Prose is not authoritative. Use `sdp-trace-claim` tags for authoritative claims.
- Source-bound proof requires a clean immutable source commit; if a changed file is a manifest subject, commit it first, then regenerate release proof in a separate commit.
- Do not close task checkboxes, review ledgers, or docs after source-bound proof without another source-bound cycle if those files are manifest subjects.
- Keep mirrored self-trace data synchronized: `assessment-input.json` must mirror self-trace arrays, and hash references must match current files.

## Required Work Loop
Every non-trivial implementation chunk needs a SpecKit delta, Socratic review before approval handoff, trace coverage for verifier/trust changes, test-first behavior for behavior changes, drift checks, live verifier state (`pass`, `fail`, `cannot_verify`, or `not_assessed`), strict review, fresh verification, and a scoped commit.

If a chunk cannot be traced or verified yet, mark `not_assessed` or `cannot_verify` with a concrete reason and create a tracked follow-up before closing.

## Pi Harness Setup
This repository includes project-level Pi configuration in `.pi/settings.json` and prompt templates in `.pi/prompts/`.

At the start of every session in this repo, run `/sdp-trace-boot` to load the mandatory process reminder.

## Quality Bar
## Quality Bar
Every claim about a gate or verdict must be evidence-backed or marked `not_assessed`. No opaque health scores.

## Engineering Stack
Target product code is Go.

No Node.js, npm, JavaScript, TypeScript, or `.mjs` tooling is allowed in the active product path.

Bash is allowed only as a thin command launcher when Go would add no product value; any active Bash needs an explicit reason.

New Go code must be small, readable, testable, covered by focused tests, and free of TODO/FIXME markers. Put measurable complexity gates in CI or docs, not only in prose.

## Claim Tags
## Claim Tags
Use `docs/claim-authoring.md` for authoritative claim syntax. Current Slice 1 validator intentionally accepts only narrow evidence forms; do not introduce arbitrary `proof:*`, `state:*`, or `none` evidence unless cross-reference verification has been implemented.

## Commands
Use current command contracts in `docs/agent-entrypoint.md` and `docs/reviewer-entrypoint.md`.


## Required Work Loop
## Required Work Loop
Every non-trivial implementation chunk needs a SpecKit delta, Socratic review before approval handoff, trace coverage for verifier/trust changes, test-first behavior for behavior changes, drift checks, live verifier state (`pass`, `fail`, `cannot_verify`, or `not_assessed`), strict review, fresh verification, and a scoped commit.

If a chunk cannot be traced or verified yet, mark `not_assessed` or `cannot_verify` with a concrete reason and create a tracked follow-up before closing.

## Pi Harness Setup

## Review Prompt (apply this lens)
Review whether the proposed docs UX makes command choice and evidence interpretation safer for a cold user. Focus on misleading state language, output confusion, profile selection, and overclaim prevention. Assume the author is overconfident. Do not validate. Do not summarize. Return only actionable issues with file/line or artifact references, or state that you cannot find any after checking the contract.
