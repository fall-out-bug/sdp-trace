# Block 32 Implementation Plan: CI PR Review Integration

Status: Draft. Implementation is blocked until Socratic review is complete and
the reviewed direction is explicitly approved.

## Slice 1: CI Profile And Prompt Contract

Files:

- `examples/pr-review/trust-sensitive-ci-pi.profile.json`
- `examples/pr-review/prompt-templates/*.md`
- `internal/prreview` tests if prompt refs need stricter validation

Work:

- Add a CI profile using `pi` roles for GLM, Kimi, and MiniMax.
- Add deterministic prompt templates that require structured
  `pr-review-result` JSON.
- Ensure role commands do not embed secrets.
- Use command shapes such as `pi --provider zai --model glm-5.1 --no-tools
  --no-context-files --no-session -p <trusted prompt text>` through the existing
  Block 30 external command runner, with command digest recorded by
  `internal/prreview`.
- Add fixture validation for the profile.

Verification:

- `jq empty schema/*.json examples/pr-review/*.json`
- focused Go tests if profile loading behavior changes.

## Slice 2: No-Secret Review Evidence Path

Files:

- `internal/prreview`
- `cmd/sdp-trace/main.go`
- focused CLI tests

Work:

- Add a way for CI to generate explicit `not_assessed` reviewer results for all
  required reviewer planes when any required provider secret is missing, without
  invoking external runners.
- Use reason `ci_model_review_not_configured`.
- Preserve packet digest and required plane identity.

Verification:

- CLI test proving no external command is executed when secrets are absent.
- Validation test proving missing-secret results do not become
  `coverage_satisfied`.

## Slice 3: GitHub Actions Workflow

Files:

- `.github/workflows/pr-review.yml`
- `docs/ci-check-policy.md` or a dedicated CI PR-review doc

Work:

- Add PR-only workflow named `pr-review-evidence-only` for
  packet/run/synthesize/validate/summarize.
- Build `sdp-trace` from the trusted base ref or use a pinned release artifact;
  never build or execute PR-head code in the secret-bearing model-review step.
- Fetch PR-head code only as diff/context data for the review packet.
- Install a pinned `pi` version in the least surprising way available for
  GitHub Actions, and document any Node/npm setup as CI runner setup rather
  than product code.
- Pass `ZAI_API_KEY`, `KIMI_API_KEY`, and `MINIMAX_API_KEY` only through job
  environment variables.
- Do not print raw `pi` stderr into durable artifacts; normalize runner errors
  to safe reason codes.
- Set packet `ci_state` to `not_assessed` in the first implementation.
- Upload safe artifact bundle.
- Avoid `pull_request_target`.
- Keep permissions read-only.
- Write a GitHub Step Summary with an evidence-only title; do not post PR
  comments in the first implementation.

Verification:

- workflow syntax review;
- local dry-run of shell script fragments where feasible;
- PR-level GitHub Actions run after PR creation.

## Slice 4: Safety And Output Redaction

Files:

- `internal/prreview` tests
- workflow docs
- possible command summary tests

Work:

- Add synthetic unsafe markers for provider keys, authenticated URLs, private
  paths, raw prompt markers, and raw model response markers.
- Assert validation/summary/artifact candidates do not echo unsafe markers.
- Ensure raw `pi` stdout and stderr are either parsed into structured result or
  normalized to digest-only/safe reason output.

Verification:

- Go safety tests.
- `git diff --check`.

## Slice 5: Review, PR, And Drift Closure

Work:

- Run Socratic spec review before implementation approval.
- After implementation, run code/correctness, trace/evidence, and
  requirements-vs-implementation review planes.
- Run PR-level review after opening the PR.
- Record missing GitHub checks as `not_assessed`, not green.

Verification:

- `go test ./...`
- `jq empty schema/*.json examples/pr-review/*.json`
- `git diff --check`
- GitHub Actions `verify`
- GitHub Actions `pr-review` state with artifact evidence.

## Risks

- `pi` installation may require Node/npm even though product code remains
  Go-first. The workflow must document this as CI runner setup, not active
  product path, and pin the installed `pi` version.
- Model cost and latency may make always-on review too expensive. An opt-in
  rollout may be better than required checks on day one.
- Fork PRs cannot safely receive repository secrets under `pull_request`.
  The first implementation must preserve `not_assessed` rather than switching
  to `pull_request_target`.
- Raw reviewer output can contain prompt/context fragments. Default retention
  must stay digest-only.
- Same-repository PRs can still contain malicious code. The secret-bearing
  review job must not build or execute PR-head code.
