# Block 23 MVP Closure Package

Date: 2026-05-08

Branch: `codex/block-23-mvp-closure`

Current source-bound manifest anchor:

- manifest: `examples/contract-foundation/contract-manifest.example.json`
- `source_commit`: `d57ceeb3eaa3252f6290b36f524d3838f46e5e5d`
- artifact count: 148

The manifest `source_commit` is the source-subject anchor for release-proof
verification. Later closure-package and review-fix commits, including `8bc7f9e`
and the follow-up review-fix commit, are branch evidence commits on top of that
anchor; they are not the source-subject commit named by the manifest.

## Fixed Blockers

| id | status | evidence |
| --- | --- | --- |
| MVP-01 source-bound proof drift | fixed locally | `release-proof` exits 0 on a clean checkout and reports `release_verification_state: "pass"`, `source_commit_status: "matched"`, `source_commit_artifact_status: "matched"` |
| MVP-04 bilingual command/profile docs | fixed locally | Russian-language handoff remains in MVP scope for Block 23; English/Russian entrypoint, CTO, team lead, and customer-question docs updated in branch; command/profile parity scans pass |
| MVP-05 stale CTO/team docs | fixed locally | `docs/cto-adoption-guide.en.md`, `docs/cto-adoption-guide.ru.md`, `docs/team-lead-playbook.en.md`, `docs/team-lead-playbook.ru.md` |
| MVP-06 customer pressure questions | fixed locally | `docs/customer-questions.en.md` and `docs/customer-questions.ru.md`; both files contain all 9 mandatory question rows |
| MVP-07 Block 22 spec drift | fixed locally | Before: Block 22 status said implementation was blocked until explicit approval. After: Block 22 status records implementation and PR review in PR #15. Fixed in review-fix commit. |

## Open Or Deferred

| id | state | reason |
| --- | --- | --- |
| MVP-02 Block 06 open retired-script mirrors | accepted demotion from Block 23 closure; carried to Block 24 | CTO decision on 2026-05-08: the Block 06 method was a toy pilot surface and must not block Block 23 MVP closure. The real gap remains open as Block 24: test `sdp-trace` on a demo repository with CI and trace artifacts. `.beads/issues.jsonl` and `bd ready` remain process evidence only, not source-bound product proof. |
| MVP-03 repository-wide quality | partial | changed releaseproof functions pass Block 23 thresholds; legacy complexity, staged packages, and repo-wide CRAP remain exceptions |
| MVP-08 remote closure | `not_assessed` | branch has not gone through PR, PR-level review, merge, and `origin/main` verification |
| external production trust | `not_assessed` | local source-bound proof is narrower than external production trust |

See `docs/research/block-23-not-assessed-registry.md` for the full registry.

## Verification Table

| command | result |
| --- | --- |
| `go test ./...` | pass |
| `go vet ./...` | pass |
| `staticcheck ./...` | pass |
| `golangci-lint run ./...` | pass |
| `gofmt -l $(rg --files -g '*.go')` | pass, no output |
| `jq empty schema/*.json` | pass |
| `git diff --check HEAD` | pass |
| `go run ./cmd/sdp-trace --help` | pass |
| `go run ./cmd/sdp-trace release-proof --manifest examples/contract-foundation/contract-manifest.example.json --out /tmp/sdp-trace-block23-release-proof.json` | pass |
| selected-state diff between committed `examples/contract-foundation/contract-release-verification.example.json` and fresh `/tmp` release-proof output | pass; compared `release_verification_state`, `trust_scope`, `source_commit_status`, `source_commit_artifact_status`, `source_commit_artifact_counts`, `trusted_contract_release`, `external_trust_profile`, and `source_commit` |
| `git diff --name-only $(jq -r '.source_commit' examples/contract-foundation/contract-manifest.example.json) -- $(jq -r '.artifacts[].path' examples/contract-foundation/contract-manifest.example.json)` | pass, no output |
| command-doc parity scan for all shipped command IDs across English/Russian customer docs | pass |
| profile/state parity scan for `adapter-capture`, `managed-harness`, `forensic-retention`, CI profiles, `not_assessed`, and `cannot_verify` | pass |
| customer-question row count in `docs/customer-questions.en.md` and `docs/customer-questions.ru.md` | pass; 9 rows each |
| retired-command current-closure scan excluding explicitly historical spec context | pass; no hits in current closure docs |
| stale Block 12 scan across customer-facing docs | pass; no hits |
| `bd ready` | pass, no ready work found |

## Command And Profile Coverage Matrix

All rows are covered in English and Russian customer-facing docs with purpose,
minimum invocation, output/trust boundary, and exit-state caveats where relevant.

| surface | covered IDs |
| --- | --- |
| commands | `wrap`, `run`, `dry-run`, `preview`, `doctor`, `verify`, `explain`, `query`, `query-pack`, `export cross-repo-posture`, `assess`, `report`, `gate`, `witness`, `release-proof`, `validate-fixtures` |
| profiles/states | `adapter-capture`, `managed-harness`, `forensic-retention`, `github-actions`, `gitlab-ci`, `buildkite`, `customer-pki`, `air-gapped`, `not_assessed`, `cannot_verify` |
| customer questions | questions 1-9 in `docs/customer-questions.en.md` and `docs/customer-questions.ru.md` |

## Slice Evidence

| slice | expected write scope | actual scope summary | scope match | verification |
| --- | --- | --- | --- | --- |
| WS1 source-bound proof | `internal/releaseproof`, contract foundation examples, source-bound manifest/proof artifacts | releaseproof code/tests plus manifest/proof refresh | yes | release-proof pass, manifest subject drift empty |
| WS2 backlog/block drift | Block 06 ledger, Beads mirror, Block 22 status | Block 06 ledger/Beads reclassification and Block 22 status update | yes | `bd ready`, review disposition, not-assessed registry |
| WS3 quality gate | releaseproof tests and quality report | focused tests, coverage, complexity, CRAP, deadcode exceptions | yes | Go/vet/staticcheck/lint/coverage/gocyclo |
| WS4/WS5 docs/customer questions | entrypoint, reviewer, CTO/team, customer docs | bilingual command/profile docs and customer question maps | yes | command/profile parity and customer-question count scans |
| WS6 closure package | `docs/research/block-23-*` | closure package, quality report, registry, review disposition | yes | review planes and focused re-review |

## Original Manifest Drift Classification

The five pre-repair mismatches were classified before repair as content drift.
They were not path drift, schema drift, or hash-algorithm drift: the paths still
existed and SHA-256 remained the active digest algorithm, but file contents had
changed after the previous manifest refresh.

| path | classification | repair |
| --- | --- | --- |
| `README.md` | content drift | manifest hash refreshed in source-bound cycle |
| `docs/flight-recorder.md` | content drift | manifest hash refreshed in source-bound cycle |
| `schema/README.md` | content drift | manifest hash refreshed in source-bound cycle |
| `specs/001-sdp-trace-time-series-evidence-substrate/spec.md` | content drift | manifest hash refreshed in source-bound cycle |
| `specs/001-sdp-trace-time-series-evidence-substrate/tasks.md` | content drift | manifest hash refreshed in source-bound cycle |

## Source-Bound Proof Result

Clean local run:

| field | value |
| --- | --- |
| `release_verification_state` | `pass` |
| `trust_scope` | `source_bound_local_release` |
| `source_commit_status` | `matched` |
| `source_commit_artifact_status` | `matched` |
| `trusted_contract_release` | `false` |
| `external_trust_profile` | `not_assessed` |

This proves only the local source-bound manifest artifact check. It does not
prove external production trust.

## Documentation And Customer Coverage

The branch updates command/profile documentation and customer-question maps in:

- `README.md`
- `docs/agent-entrypoint.md`
- `docs/reviewer-entrypoint.md`
- `docs/cto-adoption-guide.en.md`
- `docs/cto-adoption-guide.ru.md`
- `docs/team-lead-playbook.en.md`
- `docs/team-lead-playbook.ru.md`
- `docs/customer-questions.en.md`
- `docs/customer-questions.ru.md`

## Supporting Artifacts

- `docs/research/block-23-quality-report.md` (branch evidence, not a manifest subject; not covered by source-bound proof)
- `docs/research/block-23-not-assessed-registry.md` (branch evidence, not a manifest subject; not covered by source-bound proof)
- `docs/research/block-23-review-disposition.md` (branch evidence, not a manifest subject; not covered by source-bound proof)

## Remaining Gate Before Ship

| gate | current state | evidence |
| --- | --- | --- |
| closure package committed | `pass` | `8bc7f9e` plus follow-up review-fix commits |
| implementation review | `pass` for local branch review, PR-level `not_assessed` | MiniMax-M2.7, ZAI GLM-5.1, and Qwen 3.6 Plus reviews were run and valid findings were fixed or dispositioned |
| PR opened | `pass` | PR #16 |
| PR-level review | `pass` | code/correctness, trace/evidence, and requirements-vs-implementation planes ran; valid findings fixed and focused re-review approved |
| GitHub checks | `pass` | GitHub Actions `verify` passed on PR head `30737309d13421c94f3efbbb8650410c1d4da9fa` |
| post-merge `origin/main` verification | `not_assessed` | branch is not merged |
