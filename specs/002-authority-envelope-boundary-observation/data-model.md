# Data Model: Authority Envelope Boundary Observation

## Entity: AuthorityEnvelope

Declares what an actor is allowed or denied to do for a task.

Fields:

- `schema_version` [required]: contract version.
- `task_id` [required]: stable task or phase identifier.
- `policy_id` [required]: selected authority policy.
- `authority_scope` [required]: scope such as `task`, `phase`, `repository`, or `custom`.
- `actor_ref` [required]: reference to an `ActorDeclaration.actor_id`.
- `allowed_events` [required]: event types allowed by default.
- `denied_events` [required]: event types denied by default.
- `target_rules` [required]: target-specific event rules.
- `approval_requirements` [optional]: actions requiring explicit approval evidence.
- `effective_from_event_id` [optional]: event boundary for superseding policies.
- `supersedes_policy_id` [optional]: previous policy reference.

Rules:

- Absence of an envelope means authority is `not_assessed`.
- The caller supplies the selected `policy_id`; `sdp-trace` does not choose among competing envelopes.
- The envelope records declared authority; it is not proof that the actor obeyed it.
- Target patterns use slash-separated repository-relative glob syntax. `*` matches within one path segment and `**` matches recursively.
- Top-level event rules and target rules participate in the same conflict check. If an event appears in both an allow and deny position for the same resolved target, if a target rule disagrees with top-level event rules, or if overlapping target rules disagree, the envelope is invalid and evaluation is `cannot_verify`.
- If an observed action matches no top-level event rule and no target rule in the selected envelope, authority evaluation for that action is `not_assessed`.
- Target patterns are provider-neutral and must not require GitHub-specific paths.

### TargetRule

Fields:

- `rule_id` [required]: stable rule id.
- `target_pattern` [required]: repository-relative glob pattern or provider-neutral artifact reference.
- `allowed_events` [required]: event types allowed for matching targets.
- `denied_events` [required]: event types denied for matching targets.

Rules:

- Multiple target rules are matched independently.
- Ambiguous or conflicting matches make the envelope `cannot_verify`.
- No allow-wins or deny-wins default is permitted.

### ApprovalRequirement

Fields:

- `requirement_id` [required]: stable requirement id.
- `event_type` [optional]: event type requiring approval.
- `target_rule_ref` [optional]: target rule id requiring approval.
- `approval_evidence_ref` [required]: evidence reference for the approval.

Rules:

- Missing required approval evidence makes the evaluated action `outside_authority` when the action is otherwise observed.
- Malformed or inaccessible approval evidence makes evaluation `cannot_verify`.

## Entity: ActorDeclaration

Declares an actor for authority evaluation.

Fields:

- `actor_id` [required]: stable identity in the trace package.
- `actor_type` [required]: `human`, `ai_agent`, `ci_system`, `service_account`, `bot`, `unknown`, or `custom`.
- `declared_role` [required]: role label such as `observer`, `implementer`, `reviewer`, `release_manager`, `ci`, or `custom`.
- `harness` [optional]: harness identity.
- `model` [optional]: model identity.
- `model_attribution_source` [optional]: `gateway_verified`, `harness_observed`, `agent_reported`, `not_assessed`, or `cannot_verify`.
- `operation_id` [optional]: harness or gateway operation id.

Rules:

- `model` without a non-agent-reported source cannot establish model authority for a mutation.
- A bot or shared account can be an actor, but human attribution remains separate.
- `ai_agent` means an actor whose action selection is materially driven by LLM inference. `bot` means deterministic or pre-programmed automation unless model evidence says otherwise.
- Role labels have no built-in authority. An `observer` is not read-only unless the selected authority envelope says so.

## Entity: ObservedAction

Records an action observed from one or more sources.

Fields:

- `event_id` [required]: stable event id.
- `task_id` [optional]: related task id when known.
- `event_type` [required]: `observe`, `review`, `feedback`, `direct_mutation`, `commit`, `merge`, `ci_run`, `harness_tool_call`, `gateway_request`, or an extension event type prefixed with `custom:`.
- `target` [optional]: path, artifact id, resource id, or provider-neutral reference.
- `source_type` [required]: `git`, `pr_api`, `ci_artifact`, `harness_log`, `llm_gateway`, `manual_import`, or `external_assertion`.
- `evidence_refs` [required]: inspectable evidence references.
- `actor_id` [optional]: actor id when observed.
- `operation_id` [optional]: operation id when observed.
- `observed_at` [required]: timestamp or `not_assessed`.
- `observation_confidence` [required]: `single_source`, `corroborated`, `external_assertion_only`, `not_assessed`, or `cannot_verify`.

Rules:

- Git source can prove path mutation and commit binding.
- Harness source can prove tool-call, cwd, operation id, and affected paths when present.
- Gateway source can prove provider/model/request id only when linked to an observed action.
- If `task_id` is absent and no repository-level envelope is selected, authority evaluation for the action is `not_assessed`.
- `observation_confidence` describes one observed action's source support. It does not decide authority state and has no blocking semantics.

## Entity: EvidenceBinding

Links observations across sources.

Fields:

- `binding_id` [required]: stable binding id.
- `left_event_id` [required]: first event.
- `right_event_id` [required]: second event.
- `binding_type` [required]: `same_path`, `same_commit`, `same_operation`, `same_gateway_request`, `same_artifact_digest`, or `custom`.
- `binding_state` [required]: `verified`, `not_assessed`, or `cannot_verify`.
- `matched_fields` [required]: field names compared for this binding.
- `evidence_ref` [required]: evidence supporting the binding.

Rules:

- Bindings are required before model/gateway evidence can support mutation attribution.
- Failed or missing bindings must not be silently ignored.
- A binding is `verified` only when every declared `matched_fields` value agrees across both source events.
- A binding is `not_assessed` when one or more required source events is absent.
- A binding is `cannot_verify` when source events exist but cannot be parsed, integrity-checked, compared, or contain disagreeing values.
- Time-only binding is insufficient for mutation attribution unless the schema declares an explicit tolerance and another stable field also matches.

## Entity: AuthorityEvaluation

Compares an observed action against a selected authority envelope.

Fields:

- `evaluation_id` [required]: stable id.
- `event_id` [required]: observed action being evaluated.
- `policy_id` [required]: authority envelope used.
- `state` [required]: `within_authority`, `outside_authority`, `not_assessed`, or `cannot_verify`.
- `reason_code` [required]: machine-readable reason.
- `matched_rule_ref` [optional]: target or event rule that matched.
- `actor_attribution` [required]: `verified`, `not_assessed`, or `cannot_verify`.
- `tool_attribution` [required]: `verified`, `not_assessed`, or `cannot_verify`.
- `model_attribution` [required]: `verified`, `not_assessed`, or `cannot_verify`.
- `source_coverage` [required]: list of source types used, such as `git`, `harness_log`, or `llm_gateway`.
- `evidence_refs` [required]: evidence used for the evaluation.

Rules:

- `outside_authority` is a product fact, not a native gate verdict.
- If no envelope is available, `state` is `not_assessed`.
- If a required envelope or evidence source is malformed or inaccessible, `state` is `cannot_verify`.
- Policy consumers own consequences such as block, quarantine, contamination, or escalation.
- `source_coverage` is descriptive only. It is not an evidence-strength score and has no native policy meaning.

## Evidence Reference Format

Evidence references are safe URI-style strings. Initial allowed schemes:

- `file:<relative-path>` for committed package-local files.
- `artifact:<artifact-id>#<path-or-json-pointer>` for CI or retained artifact bundle members.
- `git:<commit-sha>#<path>` for immutable source references.
- `external:<opaque-id>` for external systems where the trace package cannot expose raw material.

Rules:

- Relative paths are resolved from the trace package root.
- Evidence refs must not contain credentials, authenticated URLs, raw prompts, raw model output, or private source snippets.
- A reference that cannot be resolved by the selected evidence-resolution context is `cannot_verify` when required for evaluation. If no resolver exists for an `external:` reference, required evidence that depends on that reference is `cannot_verify`.

## Event Ordering

Supersession and stale-evidence checks use explicit event ids and source refs, not wall-clock inference alone.

Rules:

- If an envelope declares `effective_from_event_id`, actions before that event are not evaluated under that envelope unless the caller explicitly selects a repository-level policy with documented scope.
- If ordering between an envelope and action cannot be established from event ids, source commits, or artifact metadata, evaluation is `cannot_verify`.
- Wall-clock timestamps may support diagnostics, but must not be the only ordering evidence for authority evaluation.

## State Transition Notes

Authority state is not monotonic. A new evidence source can revise an earlier `not_assessed` actor attribution to `verified`, or a malformed policy can move a prior candidate evaluation to `cannot_verify`.

Published trace packages should append superseding evaluations rather than rewriting prior evaluations silently.
