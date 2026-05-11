# Demo Repository Work Plan

Target repository:

- primary: `fall-out-bug/sdp-trace-demo-jvm-gsd`
- optional later polished repo: `fall-out-bug/sdp-trace-demo-github-oss-packet`

## Recommendation

Do not destructively reset current `sdp-trace-demo-jvm-gsd` history in place and
do not create a new repo by default.

Use the existing repo as the primary v2 substrate:

1. Tag current `main` as `demo-v1-observation-baseline`.
2. Add a `demo-v2-packetization` milestone or tracker in the same repo.
3. Select one existing feature PR/history slice as the first packet target.
4. Add packet/bundle/PR-template/CI artifact wiring through normal PRs.
5. Only then add new feature PRs or negative theater PRs.

The current repo's mixed history is not a defect for product proof. It is the
evidence problem `sdp-trace` is supposed to make legible. A new repo is allowed
later for a polished public sales demo, but it should not replace the first
product proof.

The first packet target is selected by evidence richness. If no existing
feature/history slice meets the first-packet minimum bar from `spec.md`, create
a new v2 feature PR in the same repo instead of forcing a weak backfill.

## Demo Roles

The demo repository is the proof target. Feature behavior in CTO-visible demo
PRs must be implemented by OpenCode + GSD +
`minimax-coding-plan/MiniMax-M2.5`, observed by `sdp-trace` as an external
flight recorder.

Codex may:

- install, rebuild, or reinstall `sdp-trace`;
- prepare setup-only recorder and CI infrastructure;
- inspect repository state and GitHub evidence;
- review OpenCode/GSD output and redirect the route with a new task;
- run verification commands when needed;
- generate and validate packets from retained evidence.

Codex must not:

- implement or repair feature behavior directly;
- ask OpenCode/GSD to maintain `sdp-trace` evidence or trace files;
- count setup-only commits as feature delivery;
- backfill missing route/provenance/evidence by editing README, tracker, or
  packet prose.

If a feature needs a code fix, Codex records the finding and sends a new
OpenCode/GSD task. If `sdp-trace` cannot record enough evidence, fix P0 product
blockers in `sdp-trace`, reinstall the recorder, and rerun or otherwise observe
the route with retained evidence before claiming feature proof. `cannot_verify`
may remain only for non-successful packets, non-blocking rows, or P1 and lower
issues that do not prevent the demo claim.

The demo trust scope is local observed unless a row has GitHub/CI, signed,
externally timestamped, or customer-equivalent evidence. Checked-in recorder
artifacts are not authority without live validation or CI-retained resolver and
digest binding. Prompt/setup metadata can support `PC-AGENT-ROUTE: partial`;
`pass` requires stronger machine evidence of recorder integrity and passivity.

## Demo V2 Repository Contract

The repository must show a reviewable GitHub story:

```text
v1 baseline tag -> selected existing feature/history packet or new v2 feature PR
packetization setup PR -> CI artifact -> packet -> review
feature 1 issue/PR or existing history -> packet -> residual gaps
feature 2 issue/PR or existing history -> packet -> residual gaps
feature 3 issue/PR or existing history -> packet -> residual gaps
feature 4 issue/PR or existing history -> packet -> residual gaps
feature 5 issue/PR or existing history -> packet -> residual gaps
DEMO-NEGATIVE draft PR -> packet with theater finding
```

Every feature PR must link:

- issue/task source;
- OpenCode/GSD observation run id;
- CI run id;
- packet artifact;
- evidence bundle artifact;
- review evidence;
- residual gaps.

## Files To Add To Existing Demo Repo

```text
.github/pull_request_template.md
.sdp-trace/demo-profile.json
.sdp-trace/packets/.gitkeep
.sdp-trace/bundles/.gitkeep
docs/demo-tracker.md
```

Modify existing `README.md` and `.github/workflows/ci.yaml` only as needed to
surface packet artifacts. Do not rebuild the app before the packet path exists.

## Setup PR Requirements

Packetization setup PR must establish:

- the selected v1 baseline tag/ref;
- artifact upload for packet and bundle paths;
- PR template with packet checklist;
- demo tracker initialized with current known evidence states.

Setup PR must not implement product features.

Codex-authored setup changes are allowed only when they are explicitly
`setup_only`; they cannot close feature packet rows.

Default setup-only file scope:

- `.github/`;
- `.sdp-trace/`;
- `docs/demo-tracker.md`;
- `.gitignore`, `.ignore`;
- `.opencode/` recorder configuration;
- packet and bundle directory placeholders.

Application source, application tests, and functional build behavior are outside
setup-only scope unless a separate independent review records why the change is
not feature behavior.

## Feature PR Requirements

For each feature packet target:

- use existing PR/history when available;
- create a new PR only when a feature has no usable GitHub evidence surface;
- require OpenCode/GSD + MiniMax as the feature implementation route;
- keep `sdp-trace` out of the developer prompt and let it observe from outside;
- route Codex review findings back to OpenCode/GSD instead of direct Codex
  patches;
- audit whether existing/backfilled feature behavior or repairs were
  Codex-authored before using the feature as CTO-visible OpenCode/GSD route
  proof;
- retain prompt text or prompt digest metadata and validate that the prompt did
  not ask OpenCode/GSD to use `sdp-trace`, author evidence, update trace files,
  or close packet rows;
- bind recorder artifacts by resolver and digest through live validation or
  CI-retained artifacts;
- inspect whether historical GitHub Actions artifacts are still available;
- CI artifact names:
  - `feature-<number>-change-evidence-packet`
  - `feature-<number>-evidence-bundle`
- packet path:
  - `.sdp-trace/packets/feature-<number>.md`
- bundle path:
  - `.sdp-trace/bundles/feature-<number>/manifest.json`

The first CTO-visible packet must be generated by the 006 product renderer and
validated by the 006 validator. Hand-authored packets are acceptable only as 006
fixtures before demo work starts; they are not acceptable 007 demo proof.

Packets are `.sdp-trace/packets/feature-<number>.md` as flat files. Bundles are
directories at `.sdp-trace/bundles/feature-<number>/manifest.json` so later
evidence entries can be added without renaming packet files.

## Negative Theater PR Requirements

Create a GitHub draft PR named:

```text
DEMO-NEGATIVE: agent-claimed verification without independent artifact
```

Add the GitHub label `demo-theater`. Retained fixtures may exist for tests, but
they are not the CTO-visible negative demo.

The draft PR should intentionally produce a packet where:

- `PC-VERIFICATION` is not `pass`;
- `PC-THEATER` has a triggered finding;
- triggered reason code is `agent_claimed_verification`;
- `PC-DECISION` does not imply approval;
- README and PR body clearly mark it as a negative demo.

## Demo Tracker Shape

`docs/demo-tracker.md` should contain:

| item | issue | PR | CI | packet | review | theater | decision |
| --- | --- | --- | --- | --- | --- | --- | --- |
| v1 baseline | n/a | n/a | state | n/a | state | n/a | n/a |
| packetization setup | link | link | state | link | state | state | owner state |
| selected feature packet | link/state | link/state | state | link | state | state | owner state |
| feature 1 | link/state | link/state | state | link | state | state | owner state |
| feature 2 | link/state | link/state | state | link | state | state | owner state |
| feature 3 | link/state | link/state | state | link | state | state | owner state |
| feature 4 | link/state | link/state | state | link | state | state | owner state |
| feature 5 | link/state | link/state | state | link | state | state | owner state |
| negative | link/state | draft PR/fixture | state | link | state | triggered | not approved |

Allowed states:

- `not_started`
- `in_progress`
- `pass`
- `partial`
- `fail`
- `not_assessed`
- `cannot_verify`

## Buyer Demo Script

The CTO demo should use one happy-path PR and one negative PR.

Happy-path script:

1. Open the feature PR.
2. Show checks and artifacts.
3. Open `change-evidence-packet.md`.
4. Read `Executive Summary`.
5. Show row table with facts, claims, evidence refs, and gaps.
6. Show Decision Ownership: owner identified, not approval.

Negative script:

1. Open draft negative PR only after the happy-path packet is understood.
2. Show agent/test claim.
3. Show missing independent evidence.
4. Show theater finding.
5. Show why packet prevents overclaim.

## Do Not Do

- Do not present current v1 repo as final CTO packet demo until packets exist.
- Do not hide `not_assessed` rows.
- Do not rely on README prose instead of packet artifacts.
- Do not count local terminal output as GitHub CI evidence.
- Do not merge negative/theater PR as a normal feature.
- Do not claim enterprise/self-hosted support from this demo.
- Do not create a new demo repo merely to avoid explaining the existing
  evidence history.
- Do not show the negative theater PR before a happy-path packet.
- Do not claim risk acceptance or security review is in scope for the Todo demo;
  those decision owners may be `not_in_scope`.
