# Block 13: Product Gap Closure Roadmap

Status: revised after first Socratic review pass.

Block 13B completion details now live in
`specs/001-sdp-trace-time-series-evidence-substrate/blocks/13b-capture-boundary-state-dx-baseline.md`.

Parent artifacts:

- `docs/cto-adoption-guide.en.md`
- `docs/team-lead-playbook.en.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/12-ci-witness-adoption.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/13b-capture-boundary-state-dx-baseline.md`
- `docs/research/block-12-documentation-review-convergence.md`

## Goal

Close the product gaps that prevent `sdp-trace` from being a credible control
layer for an existing AI SDLC.

The target buyer does not want to replace OpenCode, GSD, Superpowers, custom
agents, prompts, CI, or repository templates. The buyer wants a control layer
that can be added read-only or sidecar-first, then tightened at real control
points.

The roadmap must preserve the product boundary:

- `sdp-trace` records evidence, provenance, trace, gaps, witness state, and
  verifier-derived facts;
- `sdp-gate` or another external policy consumer decides merge, release,
  degradation, readiness, and override outcomes;
- `sdp-report` or another reporting consumer aggregates movement over repos,
  teams, services, harnesses, and time windows.

## Current State

Block 12 provides:

- process-wrapper run artifacts;
- contract-driven local report and gate artifacts;
- missing telemetry artifacts;
- GitHub Actions OIDC-backed CI witness JSON;
- honest `local_observed`, `ci_witnessed`, `cannot_verify`, and
  `audit_grade_gate: cannot_verify` states.

Block 12 does not provide:

- managed harness fail-closed enforcement;
- signed timeline, DSSE, in-toto, or external transparency witness;
- deletion and replay resistance for local run directories;
- native policy override trace events;
- adapter-level internal harness telemetry;
- model or gateway provenance;
- redaction audit trail beyond digest-only safe defaults;
- cross-repository query and degradation analytics;
- non-GitHub CI witness profiles;
- measured overhead and latency budgets.

These are product gaps, not implementation details.

## Product Gap Map

## Operating Modes

The roadmap separates adoption from enforcement. This is a product constraint,
not a packaging detail.

### Observation Mode

Read-only or sidecar-first mode for any harness:

- process wrapper may observe lifecycle and command-level metadata;
- CI/report/gate artifacts may show missing expected evidence;
- unsupported or absent adapters are reported as `not_integrated`,
  `unsupported`, `suppressed`, `missing_telemetry`, or `not_assessed`;
- no automatic blocking is implied by `sdp-trace`;
- this is the day-one path for teams that will not rewrite their harness.

### Advisory Gate Mode

Deterministic local and CI-readable gate artifacts:

- evaluates evidence contracts and required-run expectations;
- emits `pass`, `fail`, `cannot_verify`, and explanatory reason codes as
  verifier-derived facts;
- remains advisory unless an external policy consumer chooses to block;
- cannot claim protected or audit-grade enforcement by itself.

### Protected Gate Mode

CI-controlled enforcement profile:

- requires signed or CI-witnessed checkpoints from a distinct observer boundary;
- requires source, PR, CI, and artifact binding;
- treats missing required evidence as a blocking policy input for `sdp-gate` or
  the customer CI policy;
- appears only after signed checkpoint and source/merge binding exist.

### Managed Harness Mode

Opt-in tightening path for platform-owned harnesses:

- requires wrapper enrollment or adapter registration;
- can fail closed at the wrapper or CI control point;
- must not be required for observation mode;
- must not require teams to rewrite unmanaged harnesses before value appears.

### External Audit Mode

Future profile requiring an external append-only witness, transparency log,
timestamp authority, or customer PKI/private equivalent. CI-only evidence caps
at `ci_witnessed`; it does not become `external_witnessed` without this profile.

## Cross-Cutting Safety Floor

Every block that writes telemetry must apply safe pre-write retention defaults
before any artifact reaches disk.

Hard floor:

- no raw prompts by default;
- no raw model responses by default;
- no raw source snippets by default;
- no raw stdout or stderr bodies by default;
- no credentials, tokens, secrets, or OIDC request tokens in persisted
  artifacts;
- argv and command metadata must use safe descriptors, digests, and allowlisted
  basenames rather than raw sensitive values;
- local preview must show the retention mode before a run writes artifacts;
- verifier and query output must expose when retention is `digest_only`,
  `sanitized_excerpt`, `encrypted_raw_ref`, `external_artifact_ref`,
  `not_assessed`, or insufficient for the selected profile.

Block 18 owns configurable retention profiles, redaction audit trails, sealed
raw references, and forensic-grade retention behavior. It does not introduce the
first safety layer; it deepens and verifies a safety layer that must exist from
Block 13B onward.

## Interception Architecture

The product owns several observation boundaries. Each boundary has different
trust and capture limits.

| Boundary | Captures | Cannot Capture Alone | Trust Scope |
|---|---|---|---|
| Process wrapper | process lifecycle, argv descriptor, exit code, stdout/stderr digest descriptors | internal tool calls, prompts, model responses, unwrapped agents | `local_observed` |
| Adapter socket/API | harness lifecycle, task events, tool calls, file mutations, model-call declarations when emitted | events the harness suppresses or cannot observe | `harness_observed` or `agent_reported` depending on authority |
| Tool-level wrapper | selected tool invocations such as test runners, git, or build commands | direct library calls, unwrapped binary paths, shell built-ins | `local_observed` unless separately witnessed |
| VCS/PR observer | commit, branch, PR/MR id, review links, merge event references | local work before commit, internal agent actions | `vcs_observed` |
| CI observer | workflow/job/run identity, source commit/ref, artifact digests, CI-executed tests | local work before CI, agent honesty | `ci_witnessed` |
| External witness | checkpoint existence at time, transparency/audit log reference | whether the underlying work was good or complete | `external_witnessed` |

Gate contracts may require only facts that an active boundary can realistically
observe. If a required observer is absent, the output must say so; it must not
infer the missing event from git, CI, or agent prose.

### G1. Trust Authority Gap

Problem:

CI witness JSON improves the trust scope, but it is not a signed forensic
timeline and does not prevent local deletion, replay, or post-hoc fabrication by
itself.

Customer risk:

A CISO or forensics lead cannot treat the current package as audit-grade proof.

Required closure:

- canonical trace head and per-run nonce;
- monotonic event and checkpoint sequence with gap detection;
- witness checkpoint before gate;
- DSSE or in-toto-style signed checkpoint;
- identity and authority policy for who may sign which observation;
- signer isolation from the agent-controlled workspace;
- replay and source-binding checks;
- optional external append-only witness profile.

### G2. Enforcement Gap

Problem:

If the recorder is optional, a team can bypass it. Block 12 can show missing
evidence, but it does not prevent unmanaged work from reaching a protected gate.

Customer risk:

The CTO receives telemetry only from cooperative teams and cannot distinguish
"team is clean" from "team avoided capture" unless CI contracts are strict.

Required closure:

- managed harness mode;
- fail-closed profile at wrapper, adapter, and CI boundaries;
- expected run manifest or required-run contract anchored by a CI trigger,
  protected workflow, VCS/PR event, or other observer outside the
  agent-controlled workspace;
- native override request event;
- external policy consumer integration points without making `sdp-trace` the
  policy engine.

### G3. Capture Depth Gap

Problem:

The wrapper observes process lifecycle and command-level metadata. It does not
see internal tool calls, model calls, prompt flow, file mutations, or gateway
identity unless a harness or gateway emits events.

Customer risk:

Git and CI already answer part of the story. `sdp-trace` must add visible
provenance and evidence gaps, not just repackage logs.

Required closure:

- stable adapter event contract;
- tool-call and file-mutation event shapes;
- harness lifecycle event shapes;
- model/gateway provenance event shapes;
- test evidence provenance that distinguishes `ci_executed`,
  `harness_observed`, `agent_reported`, and `cannot_verify`;
- PR/MR and review linkage as first-class references;
- clear `unsupported`, `not_integrated`, `suppressed`, and `missing_telemetry`
  states.

### G4. Privacy And Retention Gap

Problem:

Forensic value competes with IP safety. Raw prompts, outputs, model responses,
source snippets, and stdout can leak secrets. Digest-only evidence can be too
weak for incident reconstruction.

Customer risk:

Developers and security teams will reject capture if redaction happens after
leakage or if retention rules are undefined.

Required closure:

- pre-write redaction profile;
- retention modes: `digest_only`, `sanitized_excerpt`,
  `encrypted_raw_ref`, `external_artifact_ref`, `not_assessed`;
- redaction authority and redaction audit trail;
- profile-based verifier behavior when retention is insufficient.
- retention enforcement checks against retained artifacts, not only a manifest.

### G5. Forensics And Query Gap

Problem:

Raw JSON and per-repo artifacts are not enough for incident review or CTO
rollout monitoring.

Customer risk:

The buyer cannot answer, "what happened across this repo, commit, PR, service,
team, and time window?" without manual reconstruction.

Required closure:

- incident timeline query;
- evidence gap query;
- override and exception query;
- source, PR, CI, and review linkage;
- approval actor and authority references where available;
- retained artifact manifest;
- cross-repo movement export for `sdp-report`.

### G6. Adoption And DX Gap

Problem:

If the product slows teams down, blocks without explainability, or requires
rewriting harnesses, it will be bypassed.

Customer risk:

The control layer becomes ceremony for compliant teams and invisible for the
teams that need oversight.

Required closure:

- dry-run and explain output for every gate-relevant state;
- local preview of what would be captured or sent;
- explicit `offline_dev` state for local work without network witness;
- deterministic verifier output for identical inputs;
- developer-facing error taxonomy separate from verifier state;
- measured wrapper overhead on real demo work;
- bootstrap guidance that does not require harness rewrites;
- non-GitHub CI profiles after the GitHub profile proves the shape;
- explicit unsupported states instead of generic failure.

## Architectural Options

### Option A: Gate-First Control Plane

Build the strongest possible protected CI contract first:

- required evidence contract;
- required run manifest;
- CI witness;
- protected artifact retention;
- fail-closed CI profile;
- native override event;
- queryable gate explanation.

Pros:

- fastest buyer value;
- turns non-cooperating agents into visible missing evidence;
- aligns with the CTO request for a control layer over existing AI SDLC;
- keeps `sdp-trace` independent from specific harnesses.

Cons:

- does not immediately provide rich internal agent telemetry;
- still depends on CI and repository governance;
- forensics remain partial until signed timeline and retention profiles land;
- if implemented before signed checkpoints, it must remain advisory rather than
  protected enforcement.

### Option B: Observer Mesh First

Build adapter, tool, file, harness, gateway, and model observers before
hardening the gate.

Pros:

- richer story than CI and git logs;
- better foundation for provenance and agent behavior analysis;
- attractive to platform owners with managed harnesses.

Cons:

- slower buyer-visible trust closure;
- many harness-specific integration surfaces;
- still bypassable without gate enforcement;
- signing incomplete telemetry can make trust claims less credible.

### Option C: Trust Anchor First

Implement DSSE, in-toto, Sigstore/Rekor or customer-PKI equivalent, and external
witness anchoring before expanding capture or enforcement.

Pros:

- addresses CISO and forensics objections directly;
- creates a strong vocabulary for authority, freshness, and replay protection;
- prevents local trace validity from being mistaken for audit proof.

Cons:

- signatures over weak capture do not create product value;
- higher integration friction;
- may feel like supply-chain ceremony before the CTO sees useful posture data.

### Option D: Reporting First

Build cross-repo dashboards, exports, and degradation analytics over current
artifacts.

Pros:

- directly addresses the CTO question about process movement;
- makes demos easier to understand;
- gives early visibility into repo/team adoption.

Cons:

- analytics over weak or bypassable capture can mislead;
- risks creating an opaque product surface before trust semantics are closed;
- can incentivize artifact volume instead of evidence quality.

## Recommended Architecture

Use a staged hybrid:

1. Capture-boundary and DX baseline first, so gate contracts do not require
   unobservable facts.
2. Advisory gate contract for buyer-visible value and deterministic
   explanations.
3. Minimal trust-anchor work before any protected enforcement claim.
4. Protected CI gate profile after signed checkpoint and source/PR/merge binding.
5. Managed harness enforcement as opt-in tightening, not day-one adoption.
6. Redaction and retention before broad capture-depth expansion.
7. Capture-depth expansion only through stable generic event contracts.
8. Query/reporting after the facts are trustworthy enough to aggregate.

This rejects two tempting shortcuts:

- Do not build dashboards over weak evidence and call it degradation detection.
- Do not build signing before the verifier knows whether the signed facts are
  complete enough for the claimed profile.

## Roadmap Blocks

### Block 13A: Product Gap Closure Contract

Purpose:

Turn this roadmap into executable SpecKit deltas and acceptance criteria.

Closes:

- scattered gap ownership;
- unclear sequencing;
- hidden overclaim risk.

Deliverables:

- reviewed roadmap;
- Russian and English customer-facing gap summary;
- block-level acceptance criteria;
- pi review ledger with CTO, Platform, CISO, Staff Engineer, and Forensics
  dispositions.

Acceptance criteria:

- every known product gap maps to a later block or an explicit non-goal;
- no block claims audit-grade, fail-closed, or degradation analytics before the
  required evidence exists;
- persona review has no remaining critical or major findings.

### Block 13B: Capture Boundary, State Taxonomy, And DX Baseline

Purpose:

Commit the physical observation model before gate contracts require evidence.

Closes:

- assumed interception boundary;
- unclear unsupported versus missing telemetry states;
- missing offline and setup diagnostics;
- unknown wrapper overhead.

Deliverables:

- capture-boundary spec covering process wrapper, adapter socket/API,
  tool-level wrapper, VCS/PR observer, CI observer, and external witness;
- machine-enumerable state taxonomy for `unsupported`, `not_integrated`,
  `suppressed`, `missing_telemetry`, `not_assessed`, `cannot_verify`, and
  `offline_dev`;
- `doctor` or equivalent setup diagnostics for wrapper, contract, and CI
  witness prerequisites;
- local preview of captured/retained fields before a run;
- hard-floor pre-write redaction and digest-only defaults for any telemetry
  written by Block 13B or later;
- deterministic verifier-output requirement;
- initial overhead measurement protocol and budget for the demo repo.

Acceptance criteria:

- every later required evidence type maps to a named observation boundary;
- unmanaged harnesses have a documented observation path that does not require
  adapter enrollment;
- identical inputs produce stable verifier and explain output;
- offline local work is represented explicitly, not mislabeled as hidden
  telemetry loss;
- raw prompts, model responses, source snippets, stdout/stderr bodies, tokens,
  secrets, and OIDC request tokens are not persisted by default;
- overhead measurement has a target budget and a repeatable command.

### Block 14: Gate Contract, Explain, And Native Override Event

Purpose:

Make CI/gate behavior operationally useful without overclaiming protected trust.

Closes:

- advisory gate clarity at the CI boundary;
- emergency-change invisibility;
- weak explanation when a gate cannot verify evidence.

Deliverables:

- native gate-result schema;
- optional required-run manifest or required-run section in the expected
  evidence contract, clearly separated between observation and enforcement
  profiles;
- `policy_override_requested` trace event schema and Go support;
- `gate explain` output;
- `gate preview` or equivalent local preview for gate-relevant captured fields;
- developer-facing error taxonomy;
- deterministic output contract for verify/gate/explain;
- CI profile guidance for advisory artifact retention;
- tests for absent run, stale witness, mismatched source, and override record.

Acceptance criteria:

- in observation mode, a missing required run produces `missing_telemetry` or
  `cannot_verify`, not pass;
- a witness for a different repository, ref, commit, run id, or artifact digest
  fails or becomes `cannot_verify`;
- override state is visible, linked to source/report artifacts, and never
  upgrades `audit_grade_gate`;
- emergency override can be recorded with one explicit CLI action or external
  reference and never hides the missing evidence;
- Block 14 output is advisory unless Block 15 signed checkpoint evidence is
  present and an external policy consumer decides to enforce;
- all behavior is implemented in Go with tests.

### Block 15: Signed Checkpoint And Replay Resistance

Purpose:

Separate structural local validity from signed witnessed existence.

Closes:

- local deletion/replay/post-hoc fabrication gap;
- plain JSON witness limitation;
- lack of signed timeline checkpoint.

Deliverables:

- canonical checkpoint statement schema;
- DSSE envelope profile using in-toto-style statement;
- signer authority policy schema;
- CI OIDC signing profile;
- verifier for identity, authority, source binding, artifact digest binding,
  checkpoint freshness, and chain head;
- signer isolation requirement: the signer cannot run inside the
  agent-controlled process, workspace, or environment that produced the local
  events;
- monotonic event and checkpoint sequence verification;
- explicit cap: Block 15 can support `ci_witnessed`, not
  `external_witnessed`;
- negative fixtures for replay, wrong signer, wrong authority, stale source,
  missing prior checkpoint, and tampered chain.

Acceptance criteria:

- local chain validity remains `local_observed`;
- signed CI checkpoint can support `ci_witnessed`;
- external audit-grade remains `cannot_verify` unless an external witness
  profile is present;
- old valid telemetry cannot be replayed against a new commit without a
  verifier-visible failure or `cannot_verify`.

### Block 16: Protected Gate Enforcement Profile

Purpose:

Allow protected CI gates to consume signed checkpoints and source/PR/merge
binding without making `sdp-trace` the policy engine.

Closes:

- enforcement gap at the protected CI boundary;
- witness-before-merge ambiguity;
- expectation manifest controlled by the same actor as the agent.

Deliverables:

- protected gate profile consuming Block 15 checkpoints;
- required-run expectation anchored by CI trigger, VCS/PR metadata, protected
  workflow configuration, or external policy input;
- PR/MR identity reference in gate and witness records;
- merge-event temporal binding: checkpoint must be shown to predate merge
  through CI/VCS evidence or remain `cannot_verify`;
- approval and override external-reference fields;
- fail-closed policy input for `sdp-gate` or customer CI, not a native
  `sdp-trace` merge verdict.

Acceptance criteria:

- a post-hoc checkpoint cannot satisfy witness-before-merge;
- PR, commit, CI run, artifact digests, and checkpoint hash are correlated;
- protected enforcement is unavailable or `cannot_verify` without Block 15
  signed checkpoint evidence;
- `sdp-trace` emits policy inputs, not native merge/release decisions.

### Block 17: Managed Harness Enforcement Profile

Purpose:

Give platform owners a real control point for teams using managed harnesses.

Closes:

- optional wrapper bypass in managed environments;
- unclear unsupported versus missing telemetry states.

Deliverables:

- managed harness profile;
- wrapper enrollment marker;
- adapter registration and authorization policy;
- fail-closed mode for missing required adapter lifecycle events;
- degraded mode for unmanaged harnesses;
- `doctor` or equivalent environment check for wrapper/adaptor setup;
- tests for outside-wrapper, late-attach, adapter disconnect, unauthorized
  adapter, and suppressed telemetry.

Acceptance criteria:

- managed profile fails closed when required control-point evidence is missing;
- unmanaged profile reports `missing_telemetry` or `not_integrated` without
  pretending capture is complete;
- observation mode still produces first evidence for unmanaged harnesses without
  enrollment;
- `sdp-trace` still does not require teams to rewrite their harness.

### Block 18: Redaction, Retention, And Forensic Profiles

Purpose:

Make capture safe enough for real repositories and useful enough for incident
review.

Closes:

- raw data leakage risk;
- weak digest-only reconstruction;
- missing redaction audit trail;
- unclear retention responsibility.

Deliverables:

- redaction profile schema;
- pre-write redaction implementation;
- retention manifest;
- retention enforcement check comparing manifest to actually retained
  artifacts;
- redaction audit event;
- encrypted external raw reference shape;
- profile-based verifier downgrade when retention is insufficient;
- tests for secret-like argv/stdout, unresolved redaction, sealed raw refs,
  and retention manifest mismatch.

Acceptance criteria:

- default profile stores no raw prompts, model responses, source snippets,
  stdout, stderr, tokens, or secrets;
- opt-in richer retention is explicit and verifier-visible;
- unresolved redaction cannot support forensic-complete profiles;
- retention gaps appear in query output.

### Block 19: Adapter Event Contract And Capture Depth

Purpose:

Expose agent SDLC provenance that CI and git cannot reconstruct.

Closes:

- missing internal tool-call telemetry;
- model/harness provenance gap;
- file mutation and task drift visibility gap.

Deliverables:

- stable adapter event schema for `run_started`, `task_locked`,
  `task_superseded`, `tool_call`, `command_started`, `file_mutation`,
  `model_call_observed`, `test_observed`, `run_closed`;
- test observation provenance values including `ci_executed`,
  `harness_observed`, `agent_reported`, and `cannot_verify`;
- source/VCS event references without binding to a specific Git host;
- PR/MR and review references without binding to a specific Git host;
- gateway provenance profile that can remain `not_integrated`;
- prompt and model-response raw capture remains unavailable by default unless
  Block 18 retention/redaction profile explicitly allows it;
- examples that are generic and do not encode OpenCode, GSD, Bazel, or Kotlin
  as product concepts;
- query output for task drift, scope creep, unverified claims, and unsupported
  observers.

Acceptance criteria:

- adapter-reported model identity remains `harness_observed` or
  `agent_reported` unless gateway evidence exists;
- test claims from an agent never appear as executed test evidence;
- missing adapter events are explicit;
- file mutation evidence can be correlated with source commit and run id;
- no product code contains demo-specific harness or build-system names.

### Block 20: Forensics Query Pack

Purpose:

Make the product useful a month after an incident.

Closes:

- raw JSON investigation burden;
- missing signed timeline view;
- weak repo/commit/PR/CI linkage.

Deliverables:

- `query timeline`;
- `query gaps`;
- `query evidence`;
- `query overrides`;
- `query provenance`;
- `query pr`;
- `query approvals`;
- artifact retention summary;
- export format for external incident and compliance systems.

Acceptance criteria:

- reviewer can reconstruct the available evidence for task, actors,
  harness/model identity state, commands, file mutations, tests, missing
  evidence, redactions, witnesses, PR links, approvals, and overrides without
  reading raw event files;
- rows explicitly show `not_assessed` and provenance state where capture was
  unavailable or unintegrated;
- query output never emits policy verdicts unless an external producer is
  named;
- every row links back to event hashes and artifact digests.
- deterministic re-execution of the agent session is not claimed; the product
  provides a replayable evidence timeline, not a guarantee that side effects can
  be reproduced.

### Block 21: Cross-Repository Degradation Export

Purpose:

Provide the evidence substrate for CTO-level process movement analysis.

Closes:

- no multi-repo degradation input;
- no trendable output for `sdp-report`.

Deliverables:

- machine-readable aggregate export by repo, team, service, harness,
  change type, and time window;
- metrics for missing telemetry, local-only evidence, CI-witnessed evidence,
  external-witnessed evidence, failed verifier states, overrides, late attach,
  unsupported observers, and contract changes;
- no opaque health score;
- documented handoff to `sdp-report` or external BI tooling.

Acceptance criteria:

- export can answer whether evidence posture is improving or worsening without
  `sdp-trace` issuing a yes/no degradation verdict;
- metrics include numerator, denominator, time window, dimensions, source
  artifact digests, and `not_assessed` counts;
- aggregation refuses or marks stale/untrusted inputs.

### Block 22: Additional CI And Enterprise Witness Profiles

Purpose:

Avoid overfitting witness behavior to GitHub Actions.

Closes:

- GitHub-only CI witness limitation.

Deliverables:

- provider-neutral witness interface;
- GitLab CI profile;
- Buildkite or Jenkins profile, selected by customer/demo need;
- customer PKI/private-equivalent profile;
- air-gapped profile documentation.

Acceptance criteria:

- provider profiles share verifier semantics;
- every profile states identity source, signing boundary, freshness boundary,
  artifact binding, and unsupported states;
- no profile can upgrade trust by environment variables alone.

## Recommended Order

1. Block 13A: finish reviewed roadmap and customer gap summary.
2. Block 13B: capture boundary, state taxonomy, doctor, preview, determinism,
   offline state, and overhead budget.
3. Block 14: advisory gate contract, gate-result schema, explain, native
   override.
4. Block 15: signed checkpoint, authority policy, replay resistance.
5. Block 16: protected gate enforcement profile with source/PR/merge binding.
6. Block 17: managed harness enforcement profile.
7. Block 18: redaction and retention profile.
8. Block 19: adapter event contract and capture depth.
9. Block 20: forensics query pack.
10. Block 21: cross-repo degradation export.
11. Block 22: additional CI and enterprise witness profiles.

The ordering intentionally places Block 18 before broad capture-depth rollout:
deeper capture without redaction and retention rules is not deployable in real
customer repositories.

## Demo Readiness Path

The OpenCode + GSD + Bazel + Kotlin demo should advance only through portable
contracts:

1. current local observed demo evidence;
2. Block 13B preview, doctor, deterministic explain, and overhead measurement;
3. Block 14 advisory gate and visible override/missing-run behavior;
4. Block 15 signed CI checkpoint;
5. Block 16 protected gate profile with source/PR/merge binding;
6. Block 17 managed-harness profile for a demo-managed wrapper path;
7. Block 18 safe redaction and retention profile;
8. Block 19 optional adapter events for OpenCode/GSD, kept in the demo repo or
   as generic adapter examples.

The demo must show real repo artifacts, not a slide show:

- report package per commit;
- gate result;
- witness/checkpoint;
- missing telemetry case;
- override case;
- tamper/replay negative case;
- incident timeline query;
- aggregate export input.

## Product Non-Goals For This Roadmap

- No built-in opaque health score.
- No native `sdp-trace` merge or release verdict.
- No dashboard before the query/export contract is stable.
- No raw prompt/source/model-response capture before redaction and retention
  profiles.
- No harness-specific product dependency.
- No external audit-grade claim from local or CI-only evidence.
- No Node.js, npm, JavaScript, TypeScript, `.mjs`, or Bash product path.
- No deterministic re-execution claim for arbitrary agent side effects; the
  forensic product is an evidence timeline unless a future profile says
  otherwise.

## Socratic Review Questions

CTO buyer:

- Can this be adopted without changing the existing harness?
- What does the CTO see when the agent avoids telemetry?
- Does the plan answer process movement without an opaque score?
- Which block creates buyer-visible value first?

Platform / Harness Owner:

- Where are the actual control points?
- What works for harnesses with no plugin API?
- What fails closed, and what only reports degraded posture?
- Can teams debug setup failures without reading raw JSON?

CISO / Adversarial Trust:

- Where are signing keys and authority policies outside the agent-controlled
  workspace?
- What prevents replay, deletion, and post-hoc trace generation?
- Which states are local-only, CI-witnessed, and externally witnessed?
- Can any self-reported field upgrade trust?

Staff Engineer / DX Skeptic:

- What is the day-one friction?
- What is captured by default, and can it leak secrets?
- How are emergency changes recorded?
- Can false positives be explained deterministically?

Compliance / Forensics Lead:

- Can an investigator reconstruct timeline, evidence, gaps, redactions, and
  approvals a month later?
- Are report rows linked to event hashes and artifact digests?
- Is retention explicit?
- Can the run be proven to have existed before merge?
