# Quickstart: Post-Merge Governance Closure

Use this path when implementing Spec 022.

## 1. Refresh live PR/CI state

Run from the repository root when GitHub access is available:

```bash
gh pr view 60 --json number,state,mergeCommit,headRefOid,reviewDecision,statusCheckRollup
gh pr checks 60
gh pr view 63 --json number,state,mergeCommit,headRefOid,reviewDecision,statusCheckRollup
gh pr checks 63
```

If GitHub access is unavailable, record the live refresh as `not_assessed` or
`cannot_verify` with the concrete reason. Do not infer live state from checked-in
ledgers alone.

For the `codex/022-post-merge-governance-closure` worktree, GitHub CLI access
is available. PR #60 and PR #63 live refresh is therefore required before any
Spec 022 `complete` claim.

## 2. Cite existing decision evidence

Confirm `split_successor` is present and consistent in:

```text
docs/closure-decision-ledger.md
docs/spec-reality-ledger.md
docs/roadmap.md
specs/019-repo-realignment-monitoring-gate-readiness/plan.md
specs/019-repo-realignment-monitoring-gate-readiness/tasks.md
specs/019-repo-realignment-monitoring-gate-readiness/post-merge-closure-plan.md
```

## 3. Decide residual remediation state

Record one of:

- no residual remediation remains, with cited evidence; or
- residual work requires reviewed successor specs before implementation.

Reviewed successor specs require a retained review artifact for their
`spec.md`, `plan.md`, and `tasks.md`; a triplet existing on disk is not enough.

Do not reopen the already-recorded accept/reject/split decision unless a new
maintainer decision explicitly supersedes `split_successor`.

## 4. Update synchronized closure surfaces

Update these together:

```text
docs/closure-decision-ledger.md
docs/spec-reality-ledger.md
docs/roadmap.md
```

The three surfaces must report the same Spec 022 closure state.

## 5. Verify

Run:

```bash
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check
```

If the change expands beyond docs, also run:

```bash
go test ./...
go vet ./...
go run ./tools/schemadoc
```
