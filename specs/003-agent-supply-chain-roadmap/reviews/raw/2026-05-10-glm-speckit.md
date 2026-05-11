# Socratic Review: Agent Supply Chain Roadmap

---

## 1. Verdict

**REVISE_BEFORE_USER_REVIEW**

The roadmap is structurally complete, honest about gaps, and well-constrained. It can reach approval with targeted fixes. The core product frame is sound. The issues below are fixable without reframing.

---

## 2. Top Socratic Questions

The owner must answer these before implementation scope is approved:

**Q1.** The spec lists 5 open questions (spec:270-278) but no task resolves them and no proposed answer accompanies any of them. Which of these must be answered before Phase 1 begins, and which can survive into implementation?

**Q2.** Three user stories share P0 priority (spec:55, 84, 112), but the plan sequences them across P0-A through P0-D. If resources are constrained to one P0 track at a time, which story delivers the first CTO wow that makes the next story easier to sell?

**Q3.** The evidence theater taxonomy (spec:192-211) identifies 8 theater types. No functional requirement binds the product to detect or report them. Should detection of these theater conditions be a P0 requirement, a P1 requirement, or an emergent property of the evidence-state model?

**Q4.** Research mentions Oh My OpenAgent (research:127-128) as "too broad for P0" but it appears in neither the spec's in-scope nor out-of-scope lists. Is it a future integration target, or should it be explicitly excluded?

**Q5.** The review context mentions Russian-market enterprise adoption. GitFlic appears as a named future adapter (spec:221, plan:35) but no artifact explains why it matters or when. Should Russian-market constraints have a research decision record, or is GitFlic listed only because it is a non-GitHub VCS?

**Q6.** The plan identifies 7 integration risks (plan:86-94) but documents zero mitigations. Is the owner comfortable approving a roadmap where risks are named but not countered, or should each risk have at least a stub mitigation before approval?

**Q7.** The spec says the product must answer one CTO question (spec:17-22). Has this question been tested with a real CTO or CTO proxy, or is it an internally formulated hypothesis? If untested, should validation be a Phase 0 task?

**Q8.** Plan slice P0-B (GitHub adapter model) and plan slice P1-A (Superpowers intent mapping) have no corresponding user stories in the spec. Are these implementation tasks that should be absorbed into existing stories, or do they need their own stories with acceptance criteria?

---

## 3. Findings

| ID | Severity | Cited File:Line | Finding | Why It Matters | Exact Fix |
|---|---|---|---|---|---|
| F-001 | **critical** | spec.md:270-278 | Open questions have no task coverage and no proposed answers. Questions 1, 3, and 5 have zero task that addresses them. | The owner is asked to approve a roadmap with 5 unresolved strategic questions and no analysis to react to. Decision support is missing. | Add proposed answers or analysis to each open question. Add tasks that explicitly resolve the ones that block Phase 1. At minimum, propose a CTO packet format (Q1) and a discovery ordering for pi vs GSD2 (Q3) so the owner can approve or override. |
| F-002 | **critical** | spec.md:192-211 | Evidence theater taxonomy lists 8 theater types but no FR requires the product to detect or report them. | The taxonomy is the spec's strongest differentiator. Without a binding FR, it is advisory prose that implementation can ignore. | Add an FR that requires the product to record theater findings using the taxonomy codes when evidence-state analysis detects the listed conditions. Alternatively, bind each theater type to an existing FR with a cross-reference. |
| F-003 | **major** | plan.md:86-94 | Integration strategy identifies 7 risks but documents 0 mitigations. | Named risks without countermeasures are complaints, not risk management. The GitHub adapter risk ("concepts can leak") is the most dangerous and has no guard. | Add a Mitigation column to the integration strategy table. For GitHub specifically, state the guard: provider-neutral field names in core schema, GitHub-specific names only in adapter layer. |
| F-004 | **major** | tasks.md:25-41 (Phase 1) | FR-010 (evidence states), FR-011 (agent claims separation), FR-012 (signed attestation as top profile), and FR-015 (docs for C-level vs engineers) have no task that explicitly addresses them. | Requirements without tasks are aspirations. Implementation will proceed without them. | Add traceability: either map each FR to existing tasks with explicit references, or add tasks for the uncovered FRs. At minimum, FR-010 and FR-011 need tasks in Phase 1, and FR-015 needs a task in Phase 8. |
| F-005 | **major** | spec.md:55,84,112 vs plan.md:98-143 | Three P0 stories map to four plan slices (P0-A through P0-D), but the mapping is implicit. Slice P0-B has no user story. Slice P1-A (Superpowers) has no user story. | Without explicit mapping, the owner cannot verify that every story has a delivery path and every slice has a validated user need. | Add a traceability table to plan.md mapping stories to slices. Create user stories or acceptance criteria for P0-B and P1-A, or absorb them into existing stories. |
| F-006 | **major** | research.md:127-128 | Oh My OpenAgent is assessed as "too broad for P0" but is not in spec.md scope (in or out). | Ambiguous scope items become scope creep during implementation. | Add Oh My OpenAgent to spec.md out-of-scope list with the reason from research.md, or add it as a named future integration with explicit not-yet-assessed status. |
| F-007 | **minor** | plan.md:9-10 vs tasks.md:25 | The boundary between "roadmap artifact" and "implementation" is stated in plan summary but not enforced in task descriptions. T008 ("Define the minimum CTO packet summary shape") could drift into schema design. | The plan's strongest constraint ("does not authorize schema, Go, CLI, or verifier implementation") has no guardrail at task level. | Add a constraint note to Phase 1 tasks: "Output is Markdown discovery/example only. No JSON schema, Go code, or CLI commands." |
| F-008 | **minor** | spec.md:53-165 | No non-functional requirements exist. A CTO-facing product has no stated expectations for packet generation latency, maximum packet size, accessibility, or localization. | NFRs discovered late create rework. Localization is relevant given Russian-market context. | Add a minimal NFR section: target packet generation latency, maximum packet size, and localization intent. These can be placeholder values for discovery. |
| F-009 | **minor** | research.md:162-173 | Six research gaps are listed but have no owners, no resolution criteria, and no task mapping. | Gaps without owners stay gaps. The Superpowers gap (line 167) and Hermes/OpenClaw gap (line 169) aren't covered until Phase 5-6, but they affect Phase 1 packet shape decisions. | Add ownership and resolution criteria to each gap. At minimum, state which gap blocks which phase. |

---

## 4. Missing Evidence and `not_assessed` Areas

| Area | Status | What's Missing |
|---|---|---|
| CTO question validation | `not_assessed` | No evidence the core CTO question (spec:17-22) has been tested with a real CTO buyer. |
| Theater taxonomy completeness | `not_assessed` | No evidence the 8 theater types were reviewed by a security or compliance practitioner. |
| Open question answers | `not_assessed` | All 5 open questions are unresolved. No proposed answers exist for the owner to react to. |
| Research gap: pi session format | `not_assessed` | Explicitly listed as a gap (research:164). Correctly honest. |
| Research gap: GSD2 export format | `not_assessed` | Explicitly listed as a gap (research:165). Correctly honest. |
| Research gap: Superpowers cross-host artifacts | `not_assessed` | Explicitly listed (research:166-167). Blocks Phase 5 but affects packet shape. |
| Research gap: Hermes/OpenClaw event API | `not_assessed` | Explicitly listed (research:168-169). Blocks Phase 6. |
| Research gap: Minimal CTO packet shape | `not_assessed` | Explicitly listed (research:170-171). **This gap blocks Phase 1 and has no task to resolve it.** |
| Research gap: Private PKI profile | `not_assessed` | Explicitly listed (research:172-173). Blocks Phase 7. Acceptable for now. |
| Risk mitigations | `not_assessed` | All 7 integration risks have no documented mitigation. |

---

## 5. Scope-Control Risks

| Risk | Likelihood | Impact | Guard Present? |
|---|---|---|---|
| T008-T012 drift from Markdown discovery into schema design | High | High | No task-level constraint. Plan line 9-10 is summary-level only. |
| One observed OpenCode/GSD run overclaimed as broad support | Medium | Critical | FR-013 and FR-014 exist but have no task. |
| General-agent boundary spike expands into employee monitoring | Medium | Critical | FR-009 exists. T032 and T033 exist. Guard is adequate if tasks execute. |
| GitHub concepts leak into provider-neutral model | High | High | Named as risk (plan:88) but no mitigation documented. |
| Signed attestation used to upgrade weak evidence | Low | Critical | F-012 exists. T037 exists. Guard is adequate. |
| Oh My OpenAgent scope ambiguity | Medium | Medium | Neither in nor out of scope. No guard. |
| 7 integration targets overstretch discovery bandwidth | Medium | Medium | Sequencing exists in plan but 4 are P0. No resource constraint acknowledged. |

---

## 6. One Strongest Reason to Proceed

The roadmap has the honesty infrastructure that most agent-trust products lack: explicit `not_assessed` and `cannot_verify` states, an evidence theater taxonomy, a clear product boundary, and a no-implementation-before-approval constraint. The core product question (spec:17-22) is sharply framed. If the fixes above are applied, this is a defensible foundation for asking a CTO to bet on `sdp-trace` as their evidence layer.

---

## 7. One Strongest Reason Not to Proceed Yet

The spec asks the owner to approve a roadmap with 5 open questions, 6 research gaps, 7 unmitigated risks, and no proposed answers for any of them. The owner is being asked to make strategic decisions without decision support. Specifically, the CTO packet format question (spec:270-271) blocks Phase 1, has no proposed answer, and has no task to produce one. Approving the roadmap in this state means the owner must resolve all strategic questions in their head during the approval meeting, with no analysis to react to. That is not a reviewed roadmap; that is a brainstorm with a approval gate attached.
