# Block 32: CI PR Review Integration

Status: Draft spec. Implementation is blocked until Socratic review is complete
and the reviewed direction is explicitly approved.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/30-automated-pr-review-evidence-mechanism.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/30-automated-pr-review-evidence-mechanism-review-ledger.md`
- `.github/workflows/ci.yml`
- `docs/ci-check-policy.md`
- `examples/self-trace/provenance-records.json`

## Goal

Add a repository CI integration for Block 30 PR-review evidence without making
model review a hidden merge approval system.

The product answer should be:

```text
On a pull request, CI can freeze the PR packet, run configured independent
review planes through explicit model runners when secrets are available, persist
safe review evidence artifacts, and report whether review coverage is
satisfied, partial, not_assessed, or cannot_verify.
```

## Product Question

"Can `sdp-trace` run its automated PR-review evidence mechanism in GitHub
Actions using configured GLM, Kimi, and MiniMax reviewers, while preserving the
boundary that review evidence is not approval?"

The answer must include:

- a repeatable GitHub Actions path for packet, model runs, ledger synthesis,
  validation, and summary;
- explicit secret requirements for the selected providers;
- a no-secret or pull-request-from-fork behavior that records `not_assessed`
  rather than pretending review passed;
- retained artifacts for packet, run results, ledger, validation, and summary;
- CI job status rules that do not collapse `coverage_satisfied` into merge,
  release, human approval, or risk acceptance;
- model identity and command provenance for GLM, Kimi, and MiniMax via `pi`;
- safe logs that do not expose API keys, raw prompts, raw model prose, private
  paths, authenticated URLs, or token-shaped content.

## Current State

Block 30 shipped the product mechanism:

- `sdp-trace pr-review packet`
- `sdp-trace pr-review run`
- `sdp-trace pr-review synthesize`
- `sdp-trace pr-review validate`
- `sdp-trace pr-review summarize`
- `sdp-trace pr-review check`

Current GitHub CI does not use it. `.github/workflows/ci.yml` only runs Go
tests, JSON syntax checks, and whitespace checks. Therefore PR review coverage
in CI is currently `not_assessed`, even when ordinary CI is green.

## Runner Choice

First implementation uses `pi`, not OpenCode, for model review execution.

Rationale:

- Existing provenance examples already use `pi` for the desired reviewers:
  - `pi --provider zai --model glm-5.1 --no-tools --no-context-files`
  - `pi --provider kimi-coding --model k2p6 --no-tools --no-context-files`
  - `pi --provider minimax --model MiniMax-M2.7 --no-tools --no-context-files`
- `pi --help` documents CI-suitable environment variables for these providers:
  `ZAI_API_KEY`, `KIMI_API_KEY`, and `MINIMAX_API_KEY`.
- `pi` supports `--no-tools`, `--no-context-files`, `--no-session`, and
  non-interactive `-p`, which is a closer fit for read-only review evidence than
  a coding-agent runner.
- OpenCode remains a future optional runner for teams that want subscription or
  agent-runtime parity, but it requires separate auth/config shape and carries a
  higher non-interference burden.

## Secret Contract

Required GitHub Actions secrets for the default Block 32 model profile:

| Secret | Used for | Required for |
| --- | --- | --- |
| `ZAI_API_KEY` | GLM reviewer through `pi --provider zai --model glm-5.1` | tracing/evidence or DX reviewer slot |
| `KIMI_API_KEY` | Kimi reviewer through `pi --provider kimi-coding --model k2p6` | requirements-vs-implementation reviewer slot |
| `MINIMAX_API_KEY` | MiniMax reviewer through `pi --provider minimax --model MiniMax-M2.7` | code/correctness or trust-boundary reviewer slot |

Optional future secrets:

| Secret | Used for | Status |
| --- | --- | --- |
| `OPENROUTER_API_KEY` | fallback OpenRouter reviewers | out of first implementation unless explicitly enabled |
| `OPENCODE_API_KEY` | OpenCode Go or OpenCode provider auth | out of first implementation |

Secrets must not be committed, echoed, copied into artifacts, or rendered in PR
summaries. CI must pass provider keys only through job environment variables.

## GitHub Actions Shape

The first workflow shape is a separate PR-only job, not a replacement for the
existing `verify` job.

Recommended job name: `pr-review`.

Trigger and trust boundary:

- `pull_request` for same-repository PRs when secrets are available.
- For pull requests where repository secrets are unavailable, run the packet and
  validation path without external runners and record model planes as
  `not_assessed`.
- Do not use `pull_request_target` in the first implementation. It would expose
  secrets to untrusted PR code unless extra checkout and diff controls are
  designed and reviewed.
- The secret-bearing model-review step must not build or execute PR-head
  product code. It must run a trusted `sdp-trace` binary built from the base
  ref or a pinned released artifact, then review the PR diff as data. The normal
  `verify` job may test PR-head code without provider secrets.
- The job/check name must be `pr-review-evidence-only` or another reviewed name
  that carries the non-approval boundary in the GitHub checks UI.
- `docs/ci-check-policy.md` must state that `pr-review-evidence-only` is review
  evidence and must not be configured as a required merge-approval substitute
  until a separate external policy block defines that behavior.

Permissions:

- `contents: read`
- `pull-requests: read`
- `actions: read`
- no pull-request or contents write permission in the first implementation

Steps:

1. Check out the trusted base ref for building `sdp-trace`.
2. Fetch the PR head as data and create a frozen PR diff file without executing
   PR-head code.
3. Install Go and a pinned `pi` version.
4. Build `sdp-trace` from the trusted base ref, or use a pinned released
   `sdp-trace` binary when available.
5. Create a Block 32 PR-review packet with `--ci-state not_assessed` in the
   first implementation. Polling other jobs and recording final ordinary CI
   state is deferred because concurrent job state is race-prone.
6. Run the trusted-base `sdp-trace pr-review run` with the Block 32 CI profile and
   `--allow-external-runner pi` only when all required secrets are present.
7. When any required secret is absent, create/import a run set that records all
   required model planes as `not_assessed` with reason
   `ci_model_review_not_configured`. This avoids leaking which individual
   provider secrets exist.
8. Synthesize ledger, validate, and summarize.
9. Write a GitHub Step Summary with an evidence-only heading and upload safe
   artifacts.
10. Fail the job only when validation is `cannot_verify` or unresolved critical
    or major findings remain under the selected profile policy.

## Review Profile

Add a CI-specific profile fixture separate from
`examples/pr-review/trust-sensitive-default.profile.json`.

Proposed path:

`examples/pr-review/trust-sensitive-ci-pi.profile.json`

Required planes:

- `code_correctness`
- `trace_evidence_provenance`
- `requirements_vs_implementation`

Default role mapping:

| Plane | Provider/model | Role |
| --- | --- | --- |
| `code_correctness` | `minimax/MiniMax-M2.7` | strict code correctness and regression reviewer |
| `trace_evidence_provenance` | `zai/glm-5.1` | trace, evidence, and replayability reviewer |
| `requirements_vs_implementation` | `kimi-coding/k2p6` | requirements-vs-implementation reviewer |

Each role must:

- use runner `pi`;
- set `required_output_schema` to `schema/pr-review-result.schema.json`;
- set `raw_output_retention` to `digest_only`;
- use a deterministic prompt template ref;
- use a default timeout of 600 seconds and a maximum timeout of 900 seconds;
- disable tools, context files, and session persistence;
- pass review input through a prompt file or stdin-like wrapper controlled by
  the trusted-base workflow, not through unescaped shell interpolation.

Prompt templates live under `examples/pr-review/prompt-templates/` for this
first block and are bound into evidence by the existing prompt safe ref and
digest fields.

## Artifacts

The CI job must upload one artifact bundle, for example
`sdp-trace-pr-review`.

Allowed files:

- `packet/packet.json`
- `runs/results.json`
- `ledger.json`
- `validation.json`
- `summary.md`
- prompt template refs and digests, if safe
- command digests
- GitHub Step Summary text matching the same safe summary constraints

Not allowed:

- raw model responses;
- raw prompts containing full private context;
- provider auth files;
- environment dumps;
- unredacted stdout/stderr;
- private local filesystem paths;
- authenticated URLs;
- token-like values.
- raw `pi` stderr.

`digest_only` means retaining the SHA-256 digest and safe ref metadata for raw
runner output while discarding the raw bytes from committed artifacts and
uploaded artifact bundles.

## Failure And State Mapping

| Condition | CI review state | Job behavior |
| --- | --- | --- |
| All required reviewers return usable structured output and no unresolved critical/major findings remain | `coverage_satisfied` | pass |
| Some required reviewer is unavailable because secrets are missing | `not_assessed` or `coverage_partial` | pass only if the workflow is explicitly in no-secret observation mode; never call it green review |
| Runner executable unavailable after setup | `not_assessed` | pass or fail according to profile policy; default is non-blocking for first rollout |
| Model timeout, empty output, parse failure, off-task output, stale packet digest, or mutation evidence | `cannot_verify` | fail |
| Unresolved critical or major finding remains | `coverage_unresolved` | fail |
| GitHub reports no ordinary CI checks | CI state `not_assessed` | do not substitute model review for CI |

First rollout uses no-secret observation mode for missing provider secrets:
all required model planes become `not_assessed` with reason
`ci_model_review_not_configured`, and the Step Summary must say review coverage
was not assessed.

Ordinary CI state in the packet is `not_assessed` for the first implementation.
Final ordinary CI pass/fail remains the separate `verify` job. A later block may
add GitHub Checks polling once stale, pending, and absent check states are
specified as a separate evidence contract.

## Acceptance Criteria

- A PR workflow job can produce the Block 30 review artifact bundle in GitHub
  Actions.
- Same-repository PRs with configured `ZAI_API_KEY`, `KIMI_API_KEY`, and
  `MINIMAX_API_KEY` run GLM, Kimi, and MiniMax through `pi` from a trusted-base
  `sdp-trace` binary, reviewing PR code only as diff/context data.
- Fork or no-secret PRs produce explicit `not_assessed` review evidence instead
  of attempting model calls or reporting review coverage as satisfied.
- Artifacts contain no raw provider keys or unsafe marker classes.
- The workflow does not require Node.js, npm, JavaScript, TypeScript, or `.mjs`
  in the active product path. If installing `pi` requires an external package
  manager, the reason and isolation boundary must be documented in the workflow
  notes; product code remains Go-first.
- The workflow does not use `pull_request_target`.
- PR summaries and logs preserve the non-approval boundary.
- Local verification and GitHub `verify` remain separate evidence from model
  review coverage.
- The GitHub check name and Step Summary visibly say evidence-only / not merge
  approval.
- `docs/ci-check-policy.md` warns maintainers not to treat
  `pr-review-evidence-only` as human approval or production trust.

## Open Questions

- Should a later block add a write-permission PR comment publisher after the
  read-only evidence workflow is proven?
- Should a later policy block make missing model-review secrets fail on
  protected branches after the secrets are installed?
