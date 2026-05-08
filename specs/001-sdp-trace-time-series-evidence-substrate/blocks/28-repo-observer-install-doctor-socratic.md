# Block 28 Socratic Review: Repo Observer Install And Doctor UX

Review date: 2026-05-08

Spec under review:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/28-repo-observer-install-doctor.md`

Status: Socratic review complete. Implementation remains blocked until the
technical executive explicitly approves the reviewed direction.

## Socratic Questions And Resolutions

### Q1. Does the spec remove hidden setup knowledge, or does it merely document it elsewhere?

**Critic**: The first draft still risked becoming another checklist if the user
had to infer which parts were install status and which parts were proof status.

**Resolution**: Accepted. The spec now requires explicit `install_state` and
`proof_state`. A repository can have `install_state=pass` while
`proof_state=not_assessed` when no CI run or artifact bundle has been observed.
This prevents a clean setup table from masquerading as proof.

### Q2. Is `install repo-observer` too broad for `sdp-trace`'s boundary?

**Critic**: Mutating `.githooks`, `.github/workflows`, `.gitignore`, and local
git config looks like installer behavior rather than observation.

**Resolution**: Accepted as in-scope only because first-mile observation is not
usable without an explicit setup surface. The spec constrains this with dry-run
default, `--write` for mutation, an allowlist of files, no global config, no
commits/pushes/PRs, and explicit local-config output. It remains setup evidence
only, not proof.

### Q3. Are local git hooks framed honestly?

**Critic**: Hooks can be bypassed with `--no-verify`, disabled by local config,
or absent in another clone. The first draft said "installed" but did not force
the output to name bypassability.

**Resolution**: Accepted and fixed. The spec now adds reason code
`local_hooks_bypassable` and requires `doctor` to report bypassability as
`local_structural`, never protected proof.

### Q4. Does `doctor` explain enough for a first-time operator to act without an agent?

**Critic**: The table shape helps, but only if next actions are concrete and
reason codes are closed. Otherwise the operator still needs hidden knowledge.

**Resolution**: Accepted. The spec requires a compact status table, closed
reason codes, and `next_action` per surface. Implementation review must reject
generic messages like "configure CI" if they do not name the missing surface.

### Q5. Are CI artifact states aligned with Block 26?

**Critic**: The first draft's aggregate `observer_state` could collapse "CI
workflow installed" with "CI uploaded artifact observed." That conflicts with
Block 26's proof-level producer semantics.

**Resolution**: Accepted and fixed through the `install_state` / `proof_state`
split. CI workflow presence is setup state; uploaded artifact bundle inspection
is proof state.

### Q6. Does generated output avoid unsafe paths, logs, prompts, and tokens?

**Critic**: The first draft required `repository_root` while also forbidding
private absolute paths in generated status. That was contradictory.

**Resolution**: Critical finding accepted and fixed. Machine output now uses
`repository_id` and `repository_root_ref`, with absolute paths forbidden unless
the user explicitly asks for verbose local human diagnostics. Machine-readable
status must not contain raw absolute repository paths.

### Q7. Are provider assumptions isolated?

**Critic**: The initial profile is GitHub Actions specific. That could turn
GitHub into hidden product architecture.

**Resolution**: Acceptable with constraints. The profile is explicitly named
`github-actions-git-hooks-v1`; future GitLab, Buildkite, customer-managed CI,
and CI-only profiles remain allowed. The product state vocabulary is portable.

### Q8. Are overwrite and force semantics safe enough for existing customer repos?

**Critic**: Installer overwrite behavior can damage existing workflows.

**Resolution**: Accepted. The draft already refuses to overwrite non-matching
files without `--force`, and requires backup or safe diff on force. During
implementation, this must be tested with existing hook/workflow fixtures.

### Q9. Is the machine-readable contract concrete enough for implementation?

**Critic**: The first draft named a schema version but did not require a schema
artifact.

**Resolution**: Accepted and fixed. Acceptance criteria now require a JSON
schema, for example `schema/repo-observer-status.schema.json`, and fixture
validation against it.

## Pi Review Findings

Raw pi-review outputs are local scratch under `.codex-review/block28/` and are
not committed.

Review planes:

- UX/DX first-mile: `zai/glm-5.1`
- Trust/evidence semantics: `minimax/MiniMax-M2.7`
- Requirements-vs-spec consistency: `openrouter/qwen/qwen3.6-plus`
- Focused re-review after fixes: `openrouter/deepseek/deepseek-v4-pro`

Initial pi verdict: `REVISE`.

Focused re-review verdict: `APPROVE`; no critical or major findings remained.
Minor residual notes about safe diff path rendering and existing `core.hooksPath`
overwrite behavior were accepted and folded into the spec.

## Review Findings

| id | severity | plane | finding | disposition | evidence |
| --- | --- | --- | --- | --- | --- |
| S28-UX-01 | major | UX/DX | Single `observer_state` lets a clean install look like proof. | accepted_fixed | Spec now requires separate `install_state` and `proof_state`. |
| S28-SAFE-01 | critical | safety | `repository_root` conflicts with the ban on private absolute paths in machine output. | accepted_fixed | Replaced with `repository_id` and `repository_root_ref`; absolute paths forbidden in machine-readable output. |
| S28-TRUST-01 | major | trust boundary | Local hooks are bypassable but the first draft did not require surfacing that fact. | accepted_fixed | Added `local_hooks_bypassable` and doctor requirement to classify hooks as `local_structural`. |
| S28-CONTRACT-01 | major | machine contract | Schema version was named but no schema artifact was required. | accepted_fixed | Acceptance criteria now require a JSON schema and validated fixtures. |
| S28-UX-02 | critical | UX/DX | `install` and `doctor` used inconsistent profile names. | accepted_fixed | Spec now uses `github-actions-git-hooks-v1` for both commands and defines `install repo-observer` as the setup subcommand. |
| S28-UX-03 | critical | UX/DX | Human output blended install and proof state in one `State` column. | accepted_fixed | Human output table now has separate `Install state` and `Proof state` columns. |
| S28-TRUST-02 | critical | trust/evidence | Trust scopes did not define which scopes can satisfy proof. | accepted_fixed | Added `Trust Scope Semantics` defining proof-satisfying and non-proof scopes. |
| S28-CONTRACT-02 | critical | machine contract | `not_assessed` and `cannot_verify` were listed as trust scopes instead of states. | accepted_fixed | Removed them from trust scope list and added state/scope separation text. |
| S28-CONTRACT-03 | critical | machine contract | `repository_id` derivation was unconstrained and could leak local identity. | accepted_fixed | Spec now defines sanitized remote hashing and forbids path-derived ids. |
| S28-INSTALL-02 | major | repo safety | `.gitignore` mutation strategy was unspecified. | accepted_fixed | Added block-marker insertion and conflict/force behavior. |
| S28-INSTALL-03 | major | repo safety | `--force` backup/diff behavior was nondeterministic. | accepted_fixed | Spec now requires per-file `.bak` plus safe diff summary. |
| S28-CONTRACT-04 | major | machine contract | `next_actions` schema was unspecified. | accepted_fixed | Added object schema and deduplication rule. |
| S28-INSTALL-04 | major | repo safety | Idempotency was unspecified. | accepted_fixed | Added no-op behavior and `already_installed` reason code. |
| S28-DOC-01 | minor | DX | `--help` output contract was unspecified. | accepted_fixed | Acceptance criteria now require help contents and golden fixture. |
| S28-SAFE-02 | minor | safety | Safe diff summaries could still leak paths; existing hook path overwrite behavior was underspecified. | accepted_fixed | Spec now requires repo-relative safe diffs and refuses mismatched `core.hooksPath` without `--force`. |
| S28-PROVIDER-01 | minor | portability | GitHub Actions first profile could be misread as product dependency. | accepted_narrower | Spec already names provider-specific profile and non-goals; future profiles remain explicit. |
| S28-INSTALL-01 | minor | repo safety | `--force` overwrite semantics need concrete tests. | accepted_narrower | Existing spec requires refusal by default and backup/safe diff on force; fixture requirement covers existing hook/workflow cases. |

## Approval Boundary

Implementation is not approved by this review. The reviewed direction can be
presented for explicit approval after the user reads the amended spec and
ledger.
