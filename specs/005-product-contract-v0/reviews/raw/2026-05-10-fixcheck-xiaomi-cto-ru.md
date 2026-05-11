# Fix-Check Review: Product Contract v0

**Reviewer role**: CTO / Russian enterprise fix-check
**Date**: 2026-05-10
**Scope**: RU-001 through RU-004 major findings + baseline credibility

---

## Verdict

**`FIXES_ACCEPTED`**

All four major findings are closed. The fixes are concrete, visible in the spec, and do not introduce new contradictions. The Russian enterprise baseline remains credible without GitHub or public SaaS.

---

## Findings Table

| id | severity | file/section | finding | exact fix |
| --- | --- | --- | --- | --- |
| FC-001 | minor | `spec.md` Evidence Bundle Manifest | `retained_form` and `redaction_status` still have overlapping value sets (`redacted`, `digest_only` appear in both). The compact display table sidesteps this for the reader, but the schema hasn't collapsed the semantic overlap. | Not blocking: the compact display addresses the practical overhead. Schema collapse can happen at implementation time. |
| FC-002 | minor | `spec.md` Packet Generation | The section states `sdp-trace` tooling generates packets but does not say what happens when the tooling is not yet implemented. A CTO asking "what do I do next week?" gets no answer. | Not blocking: the contract explicitly says implementation is future work. The examples show the intended shape. A "getting started" guide belongs in implementation docs, not the spec. |
| FC-003 | none | - | No finding. | - |
| FC-004 | none | - | No finding. | - |

No new major or critical findings.

---

## RU-001 through RU-004 Closure Status

| prior id | prior severity | status | evidence in current spec |
| --- | --- | --- | --- |
| RU-001 | major | **closed** | Compact Evidence Bundle table (`ref`, `source_class`, `retained_form`, `resolver`) for baseline display. Full manifest preserved when tooling is available. CTO reader overhead is reduced. |
| RU-002 | major | **closed** | `Packet Generation` section added: `sdp-trace` tooling, triggers are manual/CI/change-host/release, baseline requires only per-change or per-release. |
| RU-003 | major | **closed** | Profile taxonomy now explicitly says: "Rows commonly `not_assessed` or `partial` for a profile are expected profile characteristics, not product defects." `PC-AUTHORITY` is listed under expected gaps for `local-enterprise-baseline-v0`. |
| RU-004 | major | **closed** | `PC-DECISION` section now allows baseline binding via local role assignment, `git:signed-off-by`, internal wiki, or team convention. Example shows `Signed-off-by` variant. Formal policy refs required only when authority envelope is in profile scope. |

---

## Baseline Credibility Check

The `local-enterprise-baseline-v0` profile remains credible for a Russian enterprise without GitHub or public SaaS:

1. **Inputs**: local Git, GitFlic/GitLab/Gitea/Forgejo, Jenkins/TeamCity, local artifact store, customer private PKI - all self-hosted or optional.
2. **No hard dependencies**: no GitHub API, no public Sigstore/Rekor, no raw prompt export, no SaaS dashboards, no broad employee monitoring.
3. **Redaction**: defaults to `digest_only` for prompts/model/session data. Customer-specific redaction MAY follow internal classification schemes (e.g., 152-FZ). Packet records policy but does not enforce it.
4. **Expected gaps are explicit**: `PC-INITIATOR`, `PC-AGENT-ROUTE`, `PC-REVIEW`, `PC-AUTHORITY`, `PC-ATTESTATION` are listed as commonly `not_assessed` - this is a profile characteristic, not a product deficiency.
5. **Decision binding is realistic**: merge owner can be bound via `Signed-off-by`, internal wiki, or team convention, not only formal policy docs.
6. **Packet generation is lightweight**: per-change or per-release, no continuous pipeline required.

The example packet (`example-local-baseline.md`) demonstrates a valid local Git plus TeamCity packet with two triggered theater findings (`unbound_intent`, `ci_theater`) and a bound merge owner. The packet is useful because it makes missing evidence explicit rather than hiding it.

---

## Recommendation

All major findings are closed. No new blockers. Ready for explicit user approval of the reviewed Product Contract v0 direction.
