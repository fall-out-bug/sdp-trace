# Socratic Review: Agent Supply Chain Roadmap

**Reviewer**: Trust & Evidence Semantics | **Date**: 2026-05-10 | **Target**: `specs/003-agent-supply-chain-roadmap/`

---

## 1. Verdict

**REVISE_BEFORE_USER_REVIEW**

The roadmap package demonstrates good structural discipline (SpecKit shape, evidence-first framing, theater taxonomy) but has critical gaps in self-referential trust, detection mechanisms, and scope-control boundaries that could allow the product to drift into evidence theater itself.

---

## 2. Top Socratic Questions

The owner must answer these before implementation scope is approved:

1. **Self-trace gap**: AGENTS.md requires machine proof over prose. The roadmap is entirely prose. How does this roadmap package generate the evidence that it followed its own trust rules, and what is the `assessment-input.json` for this review cycle?

2. **Theater detection vs. theater taxonomy**: Spec.md defines 8 evidence theater types (lines 192-211) but P0 slices P0-A through P0-D (plan.md lines 98-143) do not describe *how* to detect them. Is theater detection in scope for P0 or a deferred concern?

3. **pi/GSD2 discovery method**: Tasks T021-T025 require inspecting `pi` and GSD2 session surfaces, but no discovery method is defined. Is this code inspection, runtime observation, documentation review, or API probing? Who does the inspecting and what stops the inspector from overclaiming?

4. **CTO packet delivery mechanism**: Open question 1 (spec.md line 270) is unanswered but drives P0-A. A PR comment, downloadable archive, static HTML, and CLI summary have fundamentally different UX, schema, and security implications. Which one creates the first product wow, and why?

5. **Signed attestation boundary**: Open question 4 (spec.md line 276) asks about minimum acceptable signed-attestation profiles. But FR-012 says signed attestation "caps the ladder" and P2-A is deferred. Without knowing what "capped" means, how do P0 packets avoid signed-theater risk?

6. **Software-delivery boundary definition**: FR-008 says general-purpose agents are in scope "when they cross a software delivery boundary" but no definition of that boundary exists. What technical or policy mechanism prevents "software delivery boundary" from expanding to mean "any employee action touching a computer"?

7. **One-run != support**: Research.md integration notes (lines 107-131) describe what tools do, not what stable evidence surface each exposes. What exact evidence surface must be inspected before the roadmap claims a row is more than `not_assessed`?

8. **Review closure for `not_assessed` rows**: T003 leaves pi and GSD2 rows `not_assessed` pending discovery. What concrete evidence must surface to close `not_assessed`? Who decides when closure criteria are met?

---

## 3. Findings Table

| ID | Severity | Location | Finding | Why It Matters | Fix |
|----|----------|----------|---------|----------------|-----|
| F01 | **Critical** | AGENTS.md:47-56, tasks.md:18-21 | Self-trace gap: The roadmap contains no evidence of its own compliance with AGENTS.md trust rules. `assessment-input.json` is not defined for this review cycle. | Prose roadmap + unchecked prose review = evidence theater for the roadmap itself. Machine proof wins over prose, but no machine proof exists for this package. | Add a lightweight self-trace manifest: git commit refs for each reviewed file, reviewer model/version, verification command outputs (markdown lint, git diff check, link validation), and explicit `not_assessed` declarations for any review plane that was not live-verified. |
| F02 | **Critical** | spec.md:192-211, plan.md:98-143 | Theater detection is taxonomied but not specified. P0 slices describe packet shape and discovery, not theater detection logic. | The product could generate packets that look like evidence but contain no theater-finding logic. CTOs would see a green-ish packet without understanding that theater detection was not implemented. | Add theater detection to P0-A or explicitly defer it with tracked open issue. Define minimum viable theater detection for P0 (at minimum: check for CI ref without retained artifact, check for review without reviewer independence evidence). |
| F03 | **Major** | spec.md:270, plan.md:99-109 | CTO packet format (Open Question 1) is unanswered but drives P0-A exit criteria. | A PR comment format vs a downloadable archive vs CLI summary vs static HTML have different trust semantics, schema requirements, and retention implications. "CTO wow" without format definition risks a demo that doesn't survive scrutiny. | Resolve Open Question 1 explicitly in spec.md before P0-A begins. Default to the format with the clearest retention and signature path unless evidence requires otherwise. |
| F04 | **Major** | spec.md:274, plan.md:154-163 | General-purpose agent boundary spike (Open Question 2) names Hermes or OpenClaw but not why. No selection rationale exists. | The boundary spike tests whether general-purpose agents can be traced when they cross into software delivery. Picking the wrong first spike wastes effort or creates misleading precedent. | Add a one-paragraph rationale for first-spike selection: which agent is most likely to produce a real boundary crossing in a testable timeframe, and what evidence surface will be inspected? |
| F05 | **Major** | spec.md:54, plan.md:43-44 | "pi session import discovery" and "GSD2 discovery" are framed as P0 but tasks T021-T025 show no defined discovery method. | Discovery without a defined method is likely to produce either shallow documentation review or premature runtime probing. Either outcome is an evidence theater risk. | Define the discovery method: (1) code inspection of session storage paths, (2) runtime observation under review profile, or (3) documentation review. Assign to a named reviewer with explicit evidence scope. |
| F06 | **Major** | spec.md:243-248, research.md:162-173 | FR-014 says rows stay `not_assessed` until evidence surfaces are inspected, but research gaps (pi session format, GSD2 state DB) are not mapped to tasks or reviewers. | `not_assessed` without closure criteria becomes a permanent holding state, not an honest evidence gap. | Add explicit mapping: research gap -> task ID -> who verifies -> what evidence must surface -> what closes the gap. |
| F07 | **Major** | spec.md:168-191, plan.md:165-177 | "Signed attestation caps the trust ladder" is a metaphor without an operational definition. FR-012 says it caps the ladder; Open Question 4 asks about minimum acceptable profiles. | "Capped" could mean: (a) signing is the highest trust level, (b) signing closes the ladder with no higher level, or (c) signing is additive but not prerequisite. Without clarity, signed profiles risk becoming theater that makes weak evidence look official. | Define "capped" explicitly: is signing additive evidence or exclusive? Define minimum packet completeness before signing is meaningful. |
| F08 | **Major** | spec.md:16-26, AGENTS.md:10-21 | The product boundary says `sdp-trace` "does not replace" a GRC tool, but FR-009 says general-purpose agent monitoring outside software-delivery boundaries is out of scope. No technical mechanism enforces this boundary. | A CTO buying this for "agent supply chain" evidence could deploy it as a broad employee monitoring tool by defining "software delivery boundary" to include any action touching a repo. The spec's prose guardrails are unenforceable. | Add a technical or policy mechanism: either (a) define minimum evidence density for "software delivery boundary" (e.g., must include at least one change-host action, CI action, or artifact mutation), or (b) add explicit prose that the product does not audit personal-agent behavior outside a tracked PR/CI/artifact context. |
| F09 | **Minor** | plan.md:86-95 | Integration strategy table lists risks but no mitigations. GitHub-specific concept leakage is identified as a risk but not mitigated. | Risks without mitigations become accepted risks, which become product debt. | Add one-line mitigation per risk: e.g., for GitHub leakage: "Use provider-neutral field names in schema; GitHub concepts map only in adapter layer." |
| F10 | **Minor** | tasks.md:23-32 | T012 says "add one hand-reviewed packet example only after source artifacts are identified; keep it marked example/discovery, not product proof." This is good but not enforced. | Without enforcement, a discovery example could be copy-pasted into sales materials as proof. | Add a file header or frontmatter tag to any discovery examples: `status: discovery-only; not-verified: true`. Add a CI check that rejects `status: verified` on any file without a traceable source commit. |
| F11 | **Minor** | research.md:133-161 | External sources section has real URLs but no last-checked dates. URLs for GitHub Agentic Workflows (2026-02-13) may be stale or moved. | Stale links in research become stale claims in product direction. | Add `last_checked` metadata to each source entry. Add a quarterly link-check task. |

---

## 4. Missing Evidence or `not_assessed` Areas

| Area | Current State | Missing Evidence | What Would Close It |
|------|---------------|------------------|----------------------|
| pi session storage/export shape | `not_assessed` (tasks.md T021) | Any pi local session artifact | One real `pi` session export file inspected by named reviewer, classified as importable/partial/unsafe |
| GSD2 state DB/session export format | `not_assessed` (tasks.md T023) | GSD2 runtime state artifact | One real GSD2 session or state DB inspected, with redaction safety verified |
| Superpowers artifact stability | `not_assessed` (research.md gap 3) | Stable Superpowers artifact across Codex/OpenCode/Copilot CLI/Claude Code hosts | One reviewed run producing a preserved skill invocation or checkpoint artifact |
| Hermes/OpenClaw event/session API | `not_assessed` (research.md gap 4) | Stable boundary event or session export | One confirmed boundary crossing event recorded by the selected agent |
| Minimum CTO packet format | `not_assessed` (Open Question 1) | Decision on PR comment vs archive vs HTML vs Markdown | Explicit decision rationale in spec.md |
| Self-trace for this review cycle | `not_assessed` | Reviewer model, verification commands, git refs, disposition records | Self-trace manifest added to `reviews/2026-05-10-socratic-review-packet.md` |
| Theater detection mechanism | `not_assessed` | Defined detection logic for each theater type | Minimum viable detection spec added to P0-A or explicitly deferred |

---

## 5. Scope-Control Risks

1. **CTO packet format creep**: If P0-A starts without a resolved Open Question 1, the format will be chosen by whoever writes the first example, potentially creating a PR-comment format that doesn't survive CTO scrutiny.

2. **One-run-overclaim cycle**: OpenCode + GSD (Slice P0-C) is the "closest real dogfood path." If one OpenCode/GSD run produces a packet, the temptation to call it "OpenCode/GSD support" is high. Research.md explicitly warns against this (line 89: "One observed profile can be overclaimed as broad support") but no enforcement mechanism exists.

3. **Signed attestation theater**: P2-A is deferred, but FR-012 says signed attestation "caps the ladder." If signed packets are presented to CTOs before evidence semantics are proven stable, signing makes incomplete evidence look official. The product must resist pressure to ship signing as a shortcut to weak evidence.

4. **Employee monitoring drift**: The spec says `sdp-trace` is "not employee surveillance" (out of scope, spec.md line 47) but does not define a technical boundary. A CTO could deploy `sdp-trace` to observe all agentic actions in an organization and call it "software delivery boundary monitoring." The prose guardrail is insufficient.

5. **GitHub-as-ontology leakage**: GitHub is the first adapter, but FR-002 says "product concepts MUST NOT be GitHub-specific." Without explicit field naming discipline (e.g., `change_host` instead of `github_repo`), every future adapter will fight GitHub assumptions.

6. **Discovery theater**: P0-D tasks T021-T025 require discovery work but do not define who does it, how they document findings, or what stops them from calling shallow documentation review "evidence." Discovery without trace is theater.

---

## 6. Strongest Reason to Proceed

**The evidence theater taxonomy is honest and useful.** Spec.md lines 192-211 name the exact failure modes that make agentic delivery opaque: agent-claimed verification without retained evidence, unbound intent, actor laundering, review theater, CI theater, artifact theater, human approval theater, and scope theater. This taxonomy gives CTOs and engineers a shared vocabulary for honest evidence assessment. If implemented with genuine `not_assessed` discipline, it directly addresses "fewer manual investigations, fewer unbacked 'done' claims, and less confusion between agent prose, CI facts, and signed evidence" (spec.md line 26).

---

## 7. Strongest Reason Not to Proceed Yet

**Self-trace gap (F01)**: AGENTS.md trust rules require machine proof over prose. The roadmap package is entirely prose. A Socratic review of prose by prose is exactly the failure mode AGENTS.md warns against: "Do not repeat that failure" (AGENTS.md line 47). Before asking the human owner to approve implementation scope, the roadmap package must demonstrate it can follow its own trust rules. Without a self-trace manifest, the review cannot be trusted to assess trust rules.

---

## 8. Disposition Summary

| Finding | Disposition |
|---------|-------------|
| F01 Critical - Self-trace gap | **Must fix**: Add self-trace manifest to review packet before user approval |
| F02 Critical - Theater detection not specified | **Must fix or explicitly defer**: Add theater detection to P0 scope or create tracked deferral |
| F03 Major - CTO packet format unresolved | **Must fix**: Resolve Open Question 1 before P0-A begins |
| F04 Major - General-agent spike rationale missing | **Should fix**: Add one-paragraph selection rationale |
| F05 Major - Discovery method undefined | **Should fix**: Define pi/GSD2 discovery method |
| F06 Major - Research gaps not mapped to tasks | **Should fix**: Add explicit gap-to-task-to-closure mapping |
| F07 Major - "Capped ladder" metaphor unoperational | **Should fix**: Define signed attestation semantics explicitly |
| F08 Major - Software-delivery boundary unenforceable | **Should fix**: Add technical or policy boundary definition |
| F09 Minor - Risk mitigations missing | **Nice to have**: Add one-line mitigations to integration table |
| F10 Minor - Discovery example not enforced | **Nice to have**: Add status header and CI check for discovery files |
| F11 Minor - External sources have no last-checked dates | **Nice to have**: Add `last_checked` metadata |

**Blocking issues for user approval**: F01, F02, F03 (3 critical+major findings must be resolved or explicitly deferred before the human owner is asked to approve implementation scope).
