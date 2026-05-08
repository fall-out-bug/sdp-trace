# Block 23 Not-Assessed Registry

Date: 2026-05-08

| id | scope | state | reason | evidence | follow-up |
| --- | --- | --- | --- | --- | --- |
| B23-NA-01 | `external_production_trust` | `not_assessed` | no external Sigstore/Rekor, customer PKI audit, protected source, workflow identity, transparency, approval, freshness, or production release evidence is verified by the local profile | `release-proof` output keeps `external_trust_profile: "not_assessed"` and `trusted_contract_release: false` | accepted non-goal for Block 23 |
| B23-NA-02 | GitHub PR checks for Block 23 branch | `not_assessed` | no PR exists yet for the current branch head | local command results in `block-23-quality-report.md` | assess during PR gate |
| B23-NA-03 | repository-wide CRAP `<5` | `not_assessed` | only changed releaseproof functions were measured for CRAP | `block-23-quality-report.md` CRAP rows | future quality-hardening block |
| B23-NA-04 | staged non-MVP packages `internal/contract`, `internal/export`, `internal/policy` | `not_assessed` | packages are currently unreachable by `deadcode` and have no tests; they are not used as current closure evidence | `deadcode` and `rg` results in `block-23-quality-report.md` | decide remove, wire, or document later |
| B23-NA-05 | remote `origin/main` verification | `not_assessed` | Block 23 branch is not merged | `git status --branch` on `codex/block-23-mvp-closure` | verify after PR merge |
| B23-NA-06 | external model PR-level review | `not_assessed` | PR-level code, trace/evidence, and requirements review cannot run before PR exists | this registry | run at PR gate |
| B23-NA-07 | Block 06 Beads mirror source-bound proof | `not_assessed` | `.beads/issues.jsonl` and `bd ready` are process/backlog evidence, not manifest subjects and not shipped `sdp-trace` commands. CTO accepted demotion from Block 23 closure scope on 2026-05-08 because Block 06 was a toy pilot method. | `.beads/issues.jsonl`, `bd ready`, Block 24 spec | carried to Block 24 demo-repository CI/trace pilot |
| B23-NA-08 | local release-proof signature and identity-policy verification | `not_assessed` | `signature_status` and `identity_policy_status` are deferred because local signing key material and policy verification are not provisioned in the local source-bound profile | `examples/contract-foundation/contract-release-verification.example.json` | future local signature/policy verification profile |
