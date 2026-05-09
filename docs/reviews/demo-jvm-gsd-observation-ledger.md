# Demo JVM GSD Observation Ledger

This ledger records findings from the attempted `sdp-trace-demo-jvm-gsd`
product proof. The demo stop condition is any P0 finding against `sdp-trace`.

## Scope

- Demo repository: `fall-out-bug/sdp-trace-demo-jvm-gsd`
- Local path: `/Users/fall_out_bug/projects/vibe_coding/sdp-trace-demo-jvm-gsd`
- Harness loop: OpenCode + GSD
- Required model route: `minimax-coding-plan/MiniMax-M2.5`
- Product-under-test: `sdp-trace` as a non-interfering observer of harness work
- Required boundary: `sdp-trace` must not alter harness prompts, manually edit
  demo repository trace artifacts, or require post-hoc hand-authored evidence
  to prove the run.

## P0

### P0-001: No non-interfering OpenCode/GSD observation path exists for the required proof

- Status: open
- Severity: P0
- Demo state: stopped before project initialization
- Observed on: 2026-05-09
- Disposition: accepted as Block 31 intake finding; implementation remains
  incomplete until a real or fixture-backed OpenCode/GSD harness export is
  validated without prompt relay or hand-authored proof artifacts.
- Tracking:
  `specs/001-sdp-trace-time-series-evidence-substrate/blocks/31-non-interfering-harness-observation.md`

#### Expected product behavior

`sdp-trace` should be able to observe an OpenCode + GSD + MiniMax-driven demo
cycle without changing the harness prompt stream and without requiring Codex or
an operator to hand-author trace artifacts in the demo repository. For each
feature cycle, the product proof needs enough evidence to reconstruct at least:

- selected model route;
- harness identity and GSD phase/review activity;
- prompts or interaction boundaries without prompt injection into the harness;
- file mutations;
- test and validation evidence;
- PR and merge state;
- explicit `pass`, `fail`, `not_assessed`, or `cannot_verify` states for gaps.

#### Observed behavior

The available `sdp-trace` product surfaces are insufficient for this proof:

- `doctor --profile github-actions-git-hooks-v1` can diagnose that repo observer
  files and CI artifact surfaces are absent, but it does not observe OpenCode or
  GSD internals.
- `install repo-observer --profile github-actions-git-hooks-v1` can propose
  repository hook and CI observer files, but using it would mutate the demo
  repository outside the OpenCode/GSD harness loop.
- `wrap` / `run` can record a command boundary, stdout/stderr digests, and local
  source state, but that is command wrapping, not a non-interfering OpenCode/GSD
  adapter with model call, prompt, review, tool-call, file-mutation, and PR-cycle
  coverage.
- `interaction relay` would route prompt delivery through `sdp-trace`, which
  violates the demo boundary that `sdp-trace` must not interfere with harness
  prompts.
- `assess preview --profile adapter-capture` reports `run: absent` and only
  names the expected evidence shape; no OpenCode/GSD adapter event source is
  installed or available.

#### Reproduction evidence

From the empty demo repository:

```text
$ git status --short --branch
## main...origin/main

$ git log --oneline --decorate --max-count=1
7bfab75 (HEAD -> main, origin/main, origin/HEAD) Reset demo repository to empty state

$ opencode --version
1.14.41

$ opencode models | rg -i "MiniMax-M2.5|minimax"
minimax-coding-plan/MiniMax-M2.5
...

$ /tmp/sdp-trace-demo-observer doctor --profile github-actions-git-hooks-v1
Install state: fail
Proof state: not_assessed
...
github_actions_artifact_bundle | not_assessed | not_assessed | ci_uploaded | filesystem:.sdp-trace/ci | run CI and inspect uploaded artifact bundle
local_wrapped_commands | not_assessed | not_assessed | not_applicable | sdp_trace_runs:not_inspected | outside selected profile; no action required
agent_prompt | not_assessed | not_assessed | agent_reported | agent_prompt:not_inspected | do not rely on prompt instructions as setup proof

$ /tmp/sdp-trace-demo-observer assess preview --profile adapter-capture
{
  "command": "assess preview",
  "selected_profile": "adapter_capture",
  "inputs": {
    "run": "absent"
  },
  "claim": "preview is read-only and does not emit an adapter capture verdict"
}
```

#### Impact

The agreed demo cannot honestly start. Starting it would either:

- rely on manual or post-hoc trace artifact creation;
- use command wrapping or prompt relay as a substitute for harness observation;
- mutate the demo repository with observer setup outside the OpenCode/GSD loop;
- or reduce the proof to CI/repo facts while missing the actual AI harness work.

All four paths would overclaim `sdp-trace` product capability.

#### Required product change

Define and implement a non-interfering OpenCode/GSD observation path. It must
produce or ingest adapter events for model calls, prompt/interaction boundaries,
GSD phase/review activity, tool calls, file mutations, test observations, PR
state, and merge state without requiring prompt changes or hand-authored demo
trace artifacts. Unsupported or unavailable fields must remain explicit
`not_assessed` or `cannot_verify` states.

#### Current fixation

The finding is fixed into the product backlog as Block 31 rather than closed.
The demo remains stopped because no verifier-backed OpenCode/GSD harness export
has been observed yet. Block 31 now provides an initial generic harness
observation command surface:

- `sdp-trace harness observe --profile <harness-profile.json> --source <harness-events.jsonl> --out <run-dir>`
- `sdp-trace harness validate --profile <harness-profile.json> --run <run-dir> --out <validation.json>`
- `sdp-trace harness summarize --validation <validation.json>`

P0-001 can move out of `open` only after Block 31 has:

- reviewed and approved SpecKit direction;
- implemented a portable non-interfering harness observation intake;
- validated a real or fixture-backed OpenCode/GSD export without prompt relay or
  hand-authored proof artifacts;
- preserved every unsupported field as `not_assessed` or `cannot_verify`.

## P1

No P1 findings recorded yet.

## P2

No P2 findings recorded yet.

## P3

No P3 findings recorded yet.
