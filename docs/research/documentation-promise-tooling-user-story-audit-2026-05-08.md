# Documentation Promise, Tooling, and User Story Audit

Date: 2026-05-08

Scope:

- Current checkout: `c2be538` / `origin/main`, after Block 24.
- Local hygiene commit under review: `f4ad9c5 Fix Go hygiene audit findings`.
- Documentation surfaces sampled: `README.md`, `docs/agent-entrypoint.md`, `docs/reviewer-entrypoint.md`, `docs/customer-questions.en.md`, `schema/README.md`, and the active SpecKit package under `specs/001-sdp-trace-time-series-evidence-substrate/`.
- Tooling checked live with current Go CLI and filesystem state.

## Hygiene Commit Handling

`f4ad9c5 Fix Go hygiene audit findings` is not an ancestor of current `origin/main`, but its touched-file patch is already present on current `main`.

Evidence:

- `git merge-base --is-ancestor f4ad9c5 HEAD` returned `ancestor=1`, so the commit object itself is not in history.
- `git diff --stat f4ad9c5..HEAD -- <f4ad9c5 touched files>` returned no file diff.
- `git cherry-pick f4ad9c5` on a branch from current `origin/main` produced an empty cherry-pick.

Delivery message:

> The hygiene fix was accounted for by content, not by preserving the old local commit object. Current `main` already contains the same touched-file changes, so there is no honest separate PR to open and no empty commit should be created.

## Live Tooling Baseline

Current command surface is Go-first:

- `go run ./cmd/sdp-trace --help`
- `go test ./...`
- `jq empty schema/*.json`
- `go run ./cmd/sdp-trace validate-fixtures <fixture-root>`
- `git diff --check`

Observed validation behavior:

```text
go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc
examples/agentic-sdlc/local-wrap-positive => observed
examples/agentic-sdlc/tamper-negative => fail
```

```text
go run ./cmd/sdp-trace validate-fixtures examples/github-speckit
no run directories found
exit status 2
```

Filesystem state:

- `scripts/` does not exist.
- No `scripts/validate-json-schema.mjs`, `scripts/validate-contracts.sh`, `scripts/check-artifact-safety.sh`, `scripts/verify-contract-manifest.sh`, `scripts/verify-self-attestation.sh`, or `scripts/finalize-source-bound-release.sh` exists in the current checkout.

## Findings

### Major: Active SpecKit docs promised retired Node/AJV validation

Evidence:

- Before this audit patch, `specs/001-sdp-trace-time-series-evidence-substrate/plan.md`, `tasks.md`, `research.md`, and `socratic-resolution-notes.md` still named pinned AJV/script validation as current.
- After this audit patch, `specs/001-sdp-trace-time-series-evidence-substrate/plan.md:14-17`, `plan.md:94`, and `tasks.md:5` name the current Go-first validation baseline.

Why it matters:

Before this patch, this directly violated the current repo contract that active product-path tooling is Go-first and that Node/npm/JavaScript tooling is not allowed in the active product path. It also broke User Story 4 because a fresh repository observer following SpecKit docs landed on tools that no longer exist.

Disposition:

Accepted and patched in this branch. Earlier AJV/script records are now described as historical evidence, not the current command contract.

### Major: Quickstart pointed at a fixture path that the current validator rejects

Evidence:

- Before this audit patch, `specs/001-sdp-trace-time-series-evidence-substrate/quickstart.md` told users to run `go run ./cmd/sdp-trace validate-fixtures examples/github-speckit`.
- The command fails with `no run directories found`.
- `go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc` succeeds.
- After this audit patch, `specs/001-sdp-trace-time-series-evidence-substrate/quickstart.md:57-67` points at `examples/agentic-sdlc` and states that `examples/github-speckit` is not a current `validate-fixtures` run package.

Why it matters:

This is a direct DX break in the reader path. It also weakens the User Story 4 promise that a reviewer can start from the SpecKit package and map task status to committed artifacts.

Disposition:

Accepted and patched in this branch.

### Major: Product docs cited removed `scripts/*` commands as live evidence sources

Evidence:

- Before this audit patch, `docs/process-metric-catalog.md` named `scripts/validate-contracts.sh`, `scripts/check-artifact-safety.sh`, `scripts/verify-contract-manifest.sh`, and `scripts/verify-self-attestation.sh` as evidence sources.
- Before this audit patch, `docs/contract-release-signing.md` told users to run `scripts/finalize-source-bound-release.sh`.
- Current checkout has no `scripts/` directory.
- After this audit patch, `docs/process-metric-catalog.md:33-36` cites current Go verifier/release-proof evidence surfaces and `docs/contract-release-signing.md:58-64` points at `go run ./cmd/sdp-trace release-proof`.

Why it matters:

These are not only historical ledgers; they are user-facing product docs. A CTO or tool builder following them will fail before reaching the current Go verifier. This undermines the "machine proof wins over prose" rule because prose points to non-existent proof machines.

Disposition:

Accepted and patched in this branch. Historical script references can stay only inside clearly historical audit notes.

### Minor: README and customer handoff docs are mostly aligned with current boundaries

Evidence:

- `README.md:7-9` says current proof authority comes from live Go verifier output and checked-in summaries are audit artifacts unless replayed or externally signed.
- `README.md:25-31` explicitly says `sdp-trace` does not decide pass/fail, readiness, or degradation and must prove itself first.
- `docs/customer-questions.en.md:9-18` maps customer questions to current Go commands and caveats, including missing telemetry, `not_assessed`, `cannot_verify`, witness kinds, and local-vs-production trust.

Why it matters:

This is the right public message shape. The stale areas are concentrated in SpecKit/current-facing product references, not in the latest customer question surface.

Disposition:

No fix needed in this pass, except keeping these docs as the source of tone and boundary language when repairing stale SpecKit text.

## User Story Alignment

| Story | Current doc alignment | Risk |
| --- | --- | --- |
| US1: CTO reviews process movement | Mostly aligned in README and customer questions. Metric catalog now points at current verifier evidence. | Low |
| US2: `sdp-gate` inherits trace contracts | Public docs preserve boundary: `sdp-trace` records inputs, external consumers decide policy. | Low |
| US5: CEO/CIO accountability and contract integrity | Conceptual promise remains present, and release-signing now points at current `release-proof`. External production trust remains intentionally not claimed. | Low |
| US3 / US3A: pilot and first real slice | Current run-card and customer pilot outline docs now avoid retired Node/script validation and require `not_assessed` when no Go verifier profile exists. Historical block ledgers still record old command evidence as history. | Low |
| US4: repository observer finds SpecKit evidence | Improved in this branch: active SpecKit docs now route to Go-first validation and the passing fixture root. Current run-card docs no longer route users to retired scripts. | Low |

## Recommended Fix Plan

1. Completed in this branch: patch active SpecKit files to state the current Go-first validation baseline:
   - `go test ./...`
   - `jq empty schema/*.json`
   - `go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc`
   - `git diff --check`
2. Completed in this branch: replace current-facing `scripts/*` references in `docs/process-metric-catalog.md` and `docs/contract-release-signing.md` with current Go verifier/release-proof commands.
3. Completed in this branch: patch live run-card and customer pilot instruction docs to avoid retired Node/script validators and mark arbitrary package schema validation `not_assessed` when no current Go verifier exists.
4. Leave block ledgers and historical research notes intact unless a current-facing page links to them as live instructions.
5. Completed after patching:
   - `go test ./...`
   - `jq empty schema/*.json`
   - `go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc`
   - `git diff --check`
