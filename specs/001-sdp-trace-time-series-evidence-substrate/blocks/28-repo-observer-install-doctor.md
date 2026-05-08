# Block 28: Repo Observer Install And Doctor UX

Status: Implemented locally; post-implementation pi re-review approved; PR
pending.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/24-demo-repo-ci-trace-pilot.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/25-compiled-ci-demo-pilot.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/26-ci-artifact-observation-contract.md`
- `docs/agent-entrypoint.md`
- `docs/ci-artifact-observation-downstream.md`

## Goal

Make repository observation setup an explicit `sdp-trace` product surface.

A first-time repository operator must be able to install, inspect, and explain
repo-level observation without knowing hidden rituals such as `core.hooksPath`,
CI artifact upload conventions, hook trust limits, or which claims remain
`not_assessed`.

The product must answer:

```text
What is observing this repository now?
What is not observing it?
Which evidence sources can support proof under the selected profile?
Which gaps remain not_assessed or cannot_verify, and what exact next action
would change that state?
```

## Dogfood Finding

The Kotlin/Bazel demo dogfood exposed a product UX failure. We first moved
trace/evidence instructions from prompts into repository scripts and CI hooks,
but the corrected setup still depended on unstated operator knowledge:

- someone had to know that local git hooks require `git config core.hooksPath`;
- someone had to know which files are tracked setup and which run outputs are
  ignored evidence;
- someone had to know that local hook output is structural evidence, not
  uploaded CI proof;
- someone had to know that GitHub Actions artifact upload is unobserved until a
  real run uploads artifacts;
- someone had to know that agent summaries and README prose never close proof
  gaps.

If those facts live only in an agent's memory, fieldbook, or demo README, the
product has failed. `sdp-trace` must surface them through CLI output and
machine-readable status.

## Problem

The repository already has `wrap`, `doctor`, CI artifact observation, witness
profiles, and provenance/evidence vocabularies. Those are useful after the user
knows how to wire a repo, but they do not provide a clear first-mile workflow.

The missing product surface is not another demo prompt. It is a repo-observer
setup contract:

- install the local observation pieces or print the exact manual changes;
- inspect whether those pieces are active in the current checkout;
- inspect whether CI observation is configured and whether it has produced
  uploaded artifacts;
- classify every setup dimension as `pass`, `fail`, `not_assessed`, or
  `cannot_verify`;
- explain trust scope without inventing policy verdicts.

Without this, a demo can look correct only because the operator already knows
the intended ritual. That is not acceptable for customer or CTO review.

## Non-Goals

- No native merge gate, release gate, readiness verdict, health score, risk
  score, or buyer-facing claim.
- No requirement that every repository use git hooks or GitHub Actions.
  Unsupported or unselected observer surfaces remain `not_assessed`.
- No global Git configuration changes. Local checkout configuration is allowed
  only with explicit user-visible output.
- No automatic commit, push, PR creation, branch protection change, secret
  creation, or GitHub repository mutation.
- No hard dependency on any single CI provider in the product core.
  Provider-specific logic is isolated to named profiles; the initial profile
  targets GitHub Actions.
- No dependency on OpenCode, GSD, MiniMax, Kotlin, Bazel, or any demo
  repository.
- No Node.js, npm, JavaScript, TypeScript, or `.mjs` product path.
- No raw prompt bodies, model responses, raw CI logs, private filesystem paths,
  authenticated URLs, token-like values, or secrets in generated status output.
- No claim that local git hooks are tamper-proof, mandatory, or equivalent to
  protected CI evidence.

## Product Boundary

Block 28 installs and diagnoses observation surfaces. It does not decide policy.

`sdp-trace` may say:

- "local git hook observation is installed in this checkout";
- "CI workflow file is present";
- "uploaded CI artifacts have not been observed yet";
- "agent-reported summary cannot satisfy proof";
- "this profile remains `not_assessed` until the first CI run uploads the
  selected artifact family."

`sdp-trace` must not say:

- "the repository is trusted";
- "the feature is delivered";
- "the PR is safe to merge";
- "the CTO should accept this";
- "local hook evidence is production proof."

## CLI Contract

Setup command:

```text
sdp-trace install repo-observer \
  --profile github-actions-git-hooks-v1 \
  [--repository-id <safe-id>] \
  [--write] \
  [--out <file>]
```

Default behavior is dry-run. Without `--write`, the command prints and optionally
writes a plan, but it does not mutate the repository.

With `--write`, the command may create or update only explicitly listed files:

- `.sdp-trace/README.md`
- `.sdp-trace/config.json`
- `.githooks/pre-commit`
- `.githooks/post-commit`
- `.githooks/pre-push`
- `.github/workflows/sdp-trace-observe.yml`
- `.gitignore` entries needed to keep volatile run outputs out of git

The command may set local repository config:

```text
git config core.hooksPath .githooks
```

If it does so, the output must say this is local checkout configuration and is
not committed into the repository.

Diagnosis command:

```text
sdp-trace doctor --profile github-actions-git-hooks-v1 --out <file>
```

`doctor --profile github-actions-git-hooks-v1` inspects the current checkout and
emits both human-readable output and a machine-readable status document.
`install repo-observer` names the setup subcommand. The `--profile` value names
the observation profile for both install and doctor. `doctor` defaults to
`github-actions-git-hooks-v1` only when `.sdp-trace/config.json` records that
profile; otherwise the user must provide `--profile` explicitly.

`--out` writes the machine-readable JSON status document. Human-readable output
always goes to stdout.

`doctor` must not require `install` to have been run. If no observer setup
exists, it reports absent surfaces with exact next actions.

## Human Output Requirements

`install repo-observer` and
`doctor --profile github-actions-git-hooks-v1` must print a compact status table
with these columns:

| Surface | Install state | Proof state | Trust scope | Evidence source | Next action |
| --- | --- | --- | --- | --- | --- |

Minimum surfaces:

- local git hooks path;
- tracked hook scripts;
- pre-commit observation;
- post-commit observation;
- pre-push observation;
- CI workflow file;
- CI artifact upload step;
- last observed CI artifact bundle;
- PR/check binding;
- local wrapped commands;
- agent/prose summaries.

Minimum trust scopes:

- `local_structural`
- `ci_uploaded`
- `external_witnessed`
- `agent_reported`
- `not_applicable`

State values such as `pass`, `fail`, `not_assessed`, and `cannot_verify` are not
trust scopes. They belong only in install/proof state columns.

The output must explicitly list "not observed" surfaces. Human output must
render `install_state` and `proof_state` as distinguishable axes, not a single
blended column. A clean install with no CI run yet is not all green; it should
show local surfaces as installed and CI artifact proof as `not_assessed` or
`cannot_verify` depending on the selected profile.

Example shape:

```text
Surface                    Install state  Proof state     Trust scope       Next action
local hooks path            pass           not_assessed    local_structural  none
tracked hook scripts        pass           not_assessed    local_structural  none
CI workflow                 pass           not_assessed    ci_uploaded       push branch and inspect run
CI artifact upload          pass           not_assessed    ci_uploaded       run CI and download artifact manifest
last CI artifact bundle     not_assessed   cannot_verify   ci_uploaded       supply extracted artifact root
agent summaries             not_assessed   not_assessed    agent_reported    supply corroborating artifact family
```

## Trust Scope Semantics

Trust scope describes evidence provenance. State describes assessment outcome.
The two dimensions must not be collapsed.

Proof-satisfying scopes for a selected profile:

- `ci_uploaded`: proof-satisfying only when the selected profile requires CI
  uploaded artifacts and the artifact family is observed in an uploaded CI
  bundle or provider artifact manifest.
- `external_witnessed`: proof-satisfying only when the selected profile requires
  an external witness and the witness binding is observed and valid.

Non-proof scopes:

- `local_structural`: local setup or local hook evidence. Useful for setup and
  diagnosis, but not sufficient for CI-uploaded or externally witnessed proof.
- `agent_reported`: agent, README, PR prose, or model summary without
  corroborating observed artifact evidence. It never satisfies proof by itself.

Surfaces not required by the selected profile are listed with
`install_state=not_assessed`, `proof_state=not_assessed`, trust scope
`not_applicable`, reason code `outside_profile_scope`, and next action
`outside selected profile; no action required`.

## Machine Output Requirements

The machine-readable output should use schema version
`block28-repo-observer-status-v1`.

Required top-level fields:

- `schema_version`
- `profile`
- `repository_id`
- `repository_root_ref`
- `git_head`
- `git_branch`
- `install_state`
- `proof_state`
- `surfaces`
- `gaps`
- `next_actions`
- `generated_at`

`repository_id` is a safe caller-supplied or derived identifier. It must not be
an absolute filesystem path, authenticated URL, or private provider URL.
`repository_root_ref` is either `"current_repository"` or another safe generated
id. Raw absolute repository paths are allowed in human terminal output only when
the user explicitly asks for verbose local diagnostics; they must not appear in
machine-readable status.

If `--repository-id` is supplied, it must match the safe identifier grammar used
by repository examples. If it is omitted, the command derives a stable safe id
by stripping credentials and fragments from `origin` remote URL, hashing the
sanitized value with SHA-256, and rendering a fixed-length prefix such as
`repo_0123456789abcdef`. If no remote is configured, it derives from the current
repository marker, not from an absolute filesystem path.

`install_state` is derived from configured observation surfaces:

| Condition | State |
| --- | --- |
| Required selected setup is present and no selected required surface is missing | `pass` |
| A selected surface is present but contradictory, unsafe, or misconfigured | `fail` |
| A required selected surface cannot be inspected | `cannot_verify` |
| A surface is outside selected profile scope | `not_assessed` |

`proof_state` is derived from live evidence under the selected proof profile and
uses the Block 26 distinction: outside selected proof scope is `not_assessed`;
selected but missing, inaccessible, expired, incomplete, or below required proof
level is `cannot_verify`.

| Condition | State |
| --- | --- |
| All selected required live evidence was observed at the selected trust scope | `pass` |
| Observed live evidence contradicts the selected setup or contains unsafe output | `fail` |
| Required live evidence is missing, inaccessible, expired, or below required proof level | `cannot_verify` |
| No live run/artifact has been supplied or selected for this profile | `not_assessed` |

This split is required so a clean install does not look like proof. A freshly
installed repo with no CI run can have `install_state=pass` and
`proof_state=not_assessed`.

Each `surface` item must include:

- `surface_id`
- `state`
- `trust_scope`
- `evidence_source`
- `observed_path` or `observed_ref`, when safe
- `reason_code`
- `next_action`

Reason codes must be closed vocabulary. Free-text errors, raw paths outside the
repository root, raw URLs, raw logs, model outputs, prompts, and token-like
values must not be embedded.

Closed-vocabulary reason codes:

- `hooks_path_absent`
- `hooks_path_set`
- `hooks_path_mismatch`
- `hook_script_absent`
- `hook_script_present`
- `hook_output_not_observed`
- `local_hooks_bypassable`
- `already_installed`
- `ci_workflow_absent`
- `ci_workflow_present`
- `ci_artifact_upload_absent`
- `ci_artifact_upload_present`
- `ci_artifact_bundle_not_observed`
- `ci_artifact_bundle_observed`
- `agent_reported_not_proof`
- `outside_profile_scope`
- `unsafe_output_refused`
- `manual_step_required`

Additional reason codes require a spec amendment.

`next_actions` is an array of objects:

```json
{
  "surface_id": "ci_artifact_bundle",
  "action_text": "run CI and supply extracted artifact root",
  "blocking": true
}
```

It is a deduplicated aggregation of per-surface `next_action` values. It must
not contain raw provider errors, paths, URLs, logs, prompts, or token-like
values.

## Install Behavior

The installer must be deterministic and reviewable:

1. Detect repository root and fail with `cannot_verify` outside a git
   repository.
2. Detect existing `.githooks` and CI workflow files.
3. In dry-run mode, emit a plan and write nothing unless `--out` is supplied.
4. In `--write` mode, create or update only allowed files.
5. Refuse to overwrite non-matching existing hook or workflow files unless
   `--force` is provided.
6. If `--force` is provided, preserve a per-file `.bak` copy and include a safe
   diff summary in the output plan. If either backup or diff generation fails,
   abort with `cannot_verify` and reason code `unsafe_output_refused`.
7. Set local `core.hooksPath` only when `--write` is provided.
8. Print exactly which surfaces are installed, which are pending, and which
   remain outside profile scope.
9. Be idempotent. If the target state already exists, file content matches the
   generated template, and `core.hooksPath` is already correct, `--write` is a
   no-op. It reports `pass` with reason code `already_installed` and creates no
   backup.

Install output is setup evidence only. It does not prove feature delivery,
build/test success, PR review, CI pass, artifact upload, or release trust.

`.gitignore` updates must use a marked block:

```text
# sdp-trace begin
...
# sdp-trace end
```

If the block already exists unchanged, installation is a no-op. If the block
exists but diverges, installation refuses without `--force`. Generated hook
scripts are placed in `.githooks/` and should be tracked in git so that cloning
the repository preserves the scripts; each checkout still needs `install` or an
equivalent local `core.hooksPath` setting.

If local `core.hooksPath` already points somewhere other than `.githooks`,
`--write` must report the mismatch and refuse to change it unless `--force` is
provided. Forced changes must show the prior and new hook path as safe relative
or generated references, not raw absolute paths.

Generated `.sdp-trace/config.json` contains the selected profile, repository id,
installed file manifest, and install metadata. Generated `.sdp-trace/README.md`
explains installed observation surfaces and points the operator to
`sdp-trace doctor --profile github-actions-git-hooks-v1`.

Generated hook scripts must have a written template contract. For the initial
profile, they must capture only safe git event metadata, staged-file summaries,
exit codes from safe structural checks such as `git diff --check`, and
repository-relative output paths. They must not persist command bodies, full
environment dumps, prompts, model output, secrets, raw logs, or private absolute
paths. Hook observation failure must not silently upgrade proof state; the hook
surface reports `cannot_verify` or `fail` according to the closed reason code.
Safe diff summaries must use repository-relative paths and must not include raw
absolute filesystem paths.

## Doctor Behavior

`doctor --profile github-actions-git-hooks-v1` must inspect:

- local git repository state;
- local `core.hooksPath`;
- tracked hook script presence and executable bit;
- ignored volatile output directories;
- configured CI workflow file presence;
- presence of an artifact-upload step, such as `actions/upload-artifact` in
  GitHub Actions or an equivalent CI-provider artifact-publishing step matching
  a documented pattern table;
- optional extracted CI artifact root if supplied;
- whether any observed artifact bundle can satisfy selected proof level.

Doctor must distinguish:

- installed but never exercised;
- exercised locally but not CI-uploaded;
- CI workflow present but artifact upload absent;
- CI artifact upload configured but no run inspected;
- agent/prose claim present without observed artifact family.
- local hook installation that can be bypassed with ordinary git mechanisms
  such as `--no-verify`; this is expected and must be reported as
  `local_structural`, not as protected proof.

## Profile Semantics

Initial profile:

```text
github-actions-git-hooks-v1
```

Required selected setup surfaces:

- local hooks install state;
- GitHub Actions workflow file;
- GitHub Actions artifact upload declaration.

Live CI artifact inspection is not required for `install` to pass. It is
`not_assessed` for the initial profile until a run or extracted artifact root is
supplied. A stricter proof mode that turns absent uploaded artifacts into
`cannot_verify` is deferred to a later named profile and must not affect Block 28
initial fixtures.

Future profiles may support GitLab CI, Buildkite, customer-managed CI, or
no-hook CI-only setups. They must preserve the same state vocabulary and avoid
provider-specific hidden rituals.

## Safety Requirements

- Generated hook scripts must not persist secrets, env dumps, full command
  bodies, prompts, model responses, raw CI logs, or private absolute paths.
- Generated CI workflow must upload only selected observation files, not the
  entire repository or arbitrary logs.
- Any output path rendered in machine output must be repository-relative or a
  safe generated id.
- If an unsafe value would be printed or serialized, the command must fail
  closed with `cannot_verify` and safe reason code `unsafe_output_refused`.
- Hook output directories must be ignored by default; committed examples must be
  small fixtures only.

## Acceptance Criteria

1. A first-time operator can run one install command and see exactly what will
   be installed before mutation.
2. After `--write`, `doctor --profile github-actions-git-hooks-v1` reports local
   hook setup without requiring agent-authored instructions.
3. `doctor` shows CI artifact proof as `not_assessed` until a CI run or artifact
   root is supplied; it does not render an all-green table.
4. Agent summaries and README claims are classified as `agent_reported` and
   cannot satisfy proof-level surfaces.
5. Local hook output is classified as `local_structural`, not `ci_uploaded` or
   external proof.
6. Existing hook or CI files are not overwritten silently.
7. The command surface is documented in `docs/agent-entrypoint.md` and `--help`
   in the same implementation change. `--help` must list subcommands, accepted
   flags, default dry-run behavior, profile semantics, and the documentation
   path; this must be covered by a golden CLI fixture.
8. The machine-readable status has a JSON schema, for example
   `schema/repo-observer-status.schema.json`, and fixtures are validated against
   it.
9. Fixtures cover absent hooks, hooks path mismatch, bypassable local hooks,
   missing CI workflow,
   missing upload-artifact step, present workflow with no run, agent-reported
   happy path, and a valid installed local setup.
10. No implementation path depends on Node.js/npm/JavaScript tooling.
11. No machine-readable output includes raw prompt, model response, token-like
    value, private URL, raw CI log body, or private absolute path.

## Required Review Plan

Before implementation approval, run Socratic spec review with at least these
questions:

1. Does the spec remove hidden setup knowledge, or does it merely document it in
   another place?
2. Is `install repo-observer` too broad for `sdp-trace`'s boundary, or is it the
   minimal UX surface required for first-mile observation?
3. Are local git hooks framed honestly as local structural evidence rather than
   proof?
4. Does `doctor` explain enough for a first-time repo operator to act without an
   agent?
5. Are CI artifact states aligned with Block 26, especially
   `not_assessed` vs `cannot_verify`?
6. Does generated output avoid leaking unsafe paths, logs, prompts, and tokens?
7. Are provider assumptions isolated to profile semantics rather than hard-coded
   product behavior?
8. Are overwrite and force semantics safe enough for existing customer repos?

Critical or major findings must be fixed or explicitly blocked before
implementation approval.
