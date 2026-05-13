# Spec Drift Register

This register records known gaps between current implementation, docs, and
historical specs. It is not proof closure. Each row stays open until a fresh
SpecKit delta, review, implementation check, and verifier state close it.

| Area | Current state | Drift | Required next step |
| --- | --- | --- | --- |
| Block 31 harness observation | `assessed_gap` | The spec and ledger now agree that the generic harness observation path is partially implemented, but T226 remains open: the customer-usable first-run OpenCode/GSD observation path is not validated. | Keep demo P0 closure open until T226 emits observed run evidence from the real first-run path with setup metadata, command digest, source commit, time bounds, and artifact digests. |
| Spec 008 invisible flight recorder | `not_assessed` | Spec status and command examples now match the implemented local `observe session` / `packet build-pr` surface. PR #43 `verify` passed for head `a8e581c1352c44ce7298af483b603edf03a9cbc0`, but PR review/sign-off evidence remains pending. | Do not claim PR-backed closure until required review planes and sign-off are observed for the final head. Current docs remain authoritative for command usage. |
| Node-era Block 01/07/08/09 specs | `stale` | Older implementation plans and ledgers mention `npm`, `.mjs`, and `scripts/verify.sh` paths that are no longer active product tooling. | Preserve only as historical artifacts or rewrite closure records around current Go commands. Do not use old commands as live evidence. |
| Roadmap spec 003 | `not_assessed` | Roadmap is marked draft / not implementation, but overlaps behavior now implemented through later blocks. | Treat as product direction only unless a later active spec cites it. Avoid using roadmap prose as implementation authority. |
| Quality gates | `assessed_gap` | CRAP `< 5`, cyclomatic `<= 10`, and cognitive `<= 10` now pass locally for production Go under `cmd`, `internal`, and `tools` and are wired into CI. Function-MI ratchet passes for `cmd`/`internal`; file-MI ratchet passes for `cmd`/`internal`/`tools`; absolute function MI now exits 0. Absolute file MI `> 70` still fails on 15 historical files and is not claimed as pass. | Keep ratcheting historical file-MI baselines down before claiming absolute MI pass. Final-head GitHub CI remains external evidence and must be observed after each push. |
| MI baseline PR packaging | `assessed_gap` | The current checkout introduces MI baselines while also changing production Go files. This is acceptable for the first gate-introduction slice, but after those baselines land the CI policy forbids changing a baseline in the same PR as `cmd/`, `internal/`, or `tools/` Go files. | Split future ratchet-baseline regeneration from product Go changes, or keep the PR blocked by `tools/mibaselinepolicy` until packaging is reviewed. |
| Gate verdict schema behavior | `assessed_gap` | `schema/gate-verdict.schema.json` now accepts `cannot_verify`, allows `external_policy_ref`, requires evidence for assessed verdicts, requires rationale for `not_assessed`/`cannot_verify`, and closes additional properties. This is a trust-contract tightening that is only partially covered by the MVP hardening slice. | Add or cite a reviewed SpecKit delta before claiming the schema behavior is fully closed; until then keep it as an assessed schema-hardening gap with JSON syntax evidence only. |

Policy: a stale spec can explain history, but it cannot close a current trust
claim. If current behavior has no active reviewed spec, record
`not_assessed` or `cannot_verify` and create a new SpecKit delta before claiming
the work is complete.
