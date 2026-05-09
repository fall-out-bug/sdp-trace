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
- Acceptable setup boundary: before the delivery work starts, the customer may
  run one initialization command, select one observation profile/configuration
  file, and install one reviewed wrapper or hook if the profile requires it.
  After that setup, the path must be "set up and forget": no extra prompt steps,
  hand-export steps, adapter authoring, or harness-loop friction during
  delivery.

## P0

### P0-001: No usable first-run OpenCode/GSD observation path exists for the customer case

- Status: open
- Severity: P0
- Demo state: stopped before project initialization
- Observed on: 2026-05-09
- Disposition: unresolved product blocker. Block 31 added a generic harness
  import surface, but it does not yet give a customer-usable first-run path for
  the OpenCode/GSD case that the demo is meant to prove.
- Tracking:
  `specs/001-sdp-trace-time-series-evidence-substrate/blocks/31-non-interfering-harness-observation.md`

#### Expected product behavior

`sdp-trace` should give a customer-usable first-run path for an OpenCode + GSD +
MiniMax-driven delivery cycle. A customer may do a bounded setup step before the
work starts, such as repository initialization, wrapper installation, or profile
configuration. After that, the observation path must be "set up and forget": the
customer should be able to run the real harness workflow without building an
adapter, rewriting the harness prompt stream, hand-exporting trace material, or
hand-authoring proof artifacts in the demo repository during delivery.

For each feature cycle, the product proof needs enough evidence to reconstruct
at least:

- selected model route;
- harness identity and GSD phase/review activity;
- prompts or interaction boundaries without prompt injection into the harness;
- file mutations;
- test and validation evidence;
- PR and merge state;
- explicit `pass`, `fail`, `not_assessed`, or `cannot_verify` states for gaps.

The demo must prove that this path works for the customer case, not merely that
`sdp-trace` can validate a pre-shaped JSONL file after someone else has solved
the observation problem.

#### Observed behavior

The available `sdp-trace` product surfaces are insufficient for this customer
path:

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
- `harness observe` can ingest an explicit `harness-event-v1` JSONL export, but
  it does not give the customer a supported way to run next to OpenCode/GSD and
  produce that export from the real workflow.

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
- require the customer or demo agent to build the missing observation adapter
  before product value is available;
- or reduce the proof to CI/repo facts while missing the actual AI harness work.

All paths would overclaim `sdp-trace` product capability. A demo that hides this
behind prepared artifacts would prove demo choreography, not product
workability.

#### Required product change

Define and implement a customer-usable first-run observation path for the
OpenCode/GSD case. It may require bounded pre-work setup, but it must let a
customer run the real OpenCode/GSD workflow after setup and produce a
verifier-backed observation package for model calls, prompt/interaction
boundaries, GSD phase/review activity, tool calls, file mutations, test
observations, PR state, and merge state without requiring prompt changes, a
customer-built adapter, hand-export chores, or hand-authored demo trace
artifacts during delivery. Unsupported or unavailable fields must remain
explicit `not_assessed` or `cannot_verify` states, but buyer-critical
delivery-loop evidence must be observed rather than replaced by disclaimers.

Allowed observation mechanisms are process-boundary capture, stdout/stderr
digests or retained-safe excerpts, declared log tailing, declared output
directory watching, and filesystem artifact reads. `sdp-trace` must not inject
stdin, rewrite harness arguments after setup, mutate undeclared environment,
hide PATH rewrites, or interpose on provider network calls.

#### Current status

The finding remains an unresolved product blocker. Block 31 merged only an
initial generic harness observation command surface:

- `sdp-trace harness observe --profile <harness-profile.json> --source <harness-events.jsonl> --out <run-dir>`
- `sdp-trace harness validate --profile <harness-profile.json> --run <run-dir> --out <validation.json>`
- `sdp-trace harness summarize --validation <validation.json>`

That surface is not enough for the customer case because it assumes a correctly
shaped harness export already exists. P0-001 can move out of `open` only after
`sdp-trace` has a supported first-run path for the actual OpenCode/GSD workflow,
for example:

```text
sdp-trace observe session --profile opencode-gsd --out <run-dir> -- <harness-command>
```

The exact CLI may change, but the closure contract must not. The path must:

- allow bounded pre-work setup, then run beside the real harness command without
  prompt relay or additional in-loop operator steps;
- produce the observation evidence needed by `harness observe` / `validate` /
  `summarize` without a customer-built adapter;
- bind setup metadata, command digest, process id or unavailable reason,
  start/end time bounds, source commit, and output artifact digests to the run;
- validate the real customer-case workflow, not only a pre-shaped fixture;
- preserve every unsupported field as `not_assessed` or `cannot_verify` without
  weakening the core claim that the delivery loop was observed.

#### Scope clarification on 2026-05-09 after Block 31 merge

Status remains `open`.

The separate demo agent is blocked because this first-run product path does not
exist. The Block 31 merge in `sdp-trace` does not by itself fix the reported
problem; it only creates a generic intake surface. Treating that surface as
sufficient would shift product work onto the customer.

Until `sdp-trace` can run against the actual customer-case workflow and produce
verifier-backed observation evidence, any claim that P0-001 is fixed would be an
overclaim. Current closure state: `unresolved_blocker`.

## P1

No P1 findings recorded yet.

## P2

No P2 findings recorded yet.

## P3

No P3 findings recorded yet.
