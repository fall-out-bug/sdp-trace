# Block 26 Review Ledger

Status: initialized for Socratic spec review.

## Spec Review Findings

| ID | Severity | Review plane | Finding | Disposition | Evidence |
| --- | --- | --- | --- | --- | --- |
| S26-LOCAL-01 | major | tracing/evidence self-review | Block 25 had CI artifact requirements but no first-class product observation contract that distinguishes uploaded CI artifacts from checked-in JSON and prose. | accepted_fixed | Block 26 spec introduces producer scope, access state, required artifact families, and checked-in-only failure mode. |
| S26-LOCAL-02 | major | product/customer credibility self-review | A demo can look credible while proving only local or committed artifacts. | accepted_fixed | Block 26 requires demo forward checks with per-feature PR, CI run, artifact bundle, and observation result. |
| S26-LOCAL-03 | major | product-boundary self-review | Turning this into a gate would violate `sdp-trace`'s flight-recorder boundary. | accepted_fixed | Non-goals and goal state that outputs are facts for downstream consumers, not native policy decisions. |
| S26-LOCAL-04 | major | security/privacy self-review | Artifact inspection could leak provider tokens, raw logs, prompts, private URLs, or private paths. | accepted_fixed | Safety non-goals, failure modes, acceptance criteria, and implementation plan require output-safety tests. |
| S26-PC-01 | major | product/customer credibility | The five-feature forward checks described developer workflow more than CTO buyer meaning. | accepted_fixed | Added Customer Implications section with what Block 26 does not claim and state implications. |
| S26-PC-02 | minor | product/customer credibility | Artifact families were not tied to customer-visible trust gaps. | accepted_fixed | Required Artifact Families table now explains the trust gap each family closes. |
| S26-PC-03 | minor | product/customer credibility | Unsafe output was not linked to customer risk. | accepted_fixed | Customer Implications now explains credential, prompt, private path, and data-egress risk without claiming prevention. |
| S26-PC-04 | minor | product/customer credibility | Non-goals did not say how repositories without CI artifacts should be interpreted. | accepted_fixed | Non-goals now distinguish profile `not_assessed` from required-CI `cannot_verify`. |
| S26-PC-05 | minor | product/customer credibility | Customers had to infer meaning of `fail`, `cannot_verify`, and `not_assessed`. | accepted_fixed | Added customer-readable state implication table. |
| S26-PC-06 | minor | product/customer credibility | "valid bundle" fixture naming could imply compliance rather than selected-profile coverage. | accepted_fixed | Implementation plan renames the valid scenario to `ci-uploaded-bundle-complete-coverage`. |
| S26-TE-01 | critical | tracing/evidence | The spec lacked an aggregate state derivation rule. | accepted_fixed | Added State Derivation table for `fail`, `cannot_verify`, `not_assessed`, and `pass`. |
| S26-TE-02 | major | tracing/evidence | `agent_reported_happy_path` allowed either `cannot_verify` or `not_assessed` without a deterministic criterion. | accepted_fixed | Failure mode now uses `cannot_verify` when selected by profile and `not_assessed` only outside profile scope. |
| S26-TE-03 | major | tracing/evidence | Checked-in-only evidence disqualification was not mapped into profile-required `ci_uploaded` behavior. | accepted_fixed | Observation Model now maps lower-authority producer scopes to proof-level mismatch and `cannot_verify`. |
| S26-TE-04 | major | tracing/evidence | Binding and artifact-index states were referenced but not enumerated. | accepted_fixed | Added field vocabulary table with binding and index states. |
| S26-TE-05 | minor | tracing/evidence | `partial` access-state propagation was undefined. | accepted_fixed | State Derivation defines `partial` as per-family and contributing `cannot_verify`. |
| S26-TE-07 | minor | tracing/evidence | `external_artifact_ref` had no failure-mode fixture. | accepted_fixed | Added `external_artifact_ref_unverifiable` failure mode and fixture scenario. |
| S26-TE-08 | minor | tracing/evidence | `pr_ci` risked provider-specific and circular semantics. | accepted_fixed | Renamed canonical family to `change_ci` and clarified that it records artifact-backed change/branch CI evidence, not provider policy. |
| S26-TE-09 | minor | tracing/evidence | Multiple `cannot_verify` reasons could be lost. | accepted_fixed | State Derivation now says reasons accumulate. |
| S26-TE-10 | minor | tracing/evidence | `redaction_scan` and evaluator `output_safety` could be confused. | accepted_fixed | Required Artifact Families now separates artifact-bundle scan evidence from evaluator output-safety. |
| S26-TE-11 | minor | tracing/evidence | Incomplete and dishonest demos both collapsed into `cannot_verify` without distinct reasons. | accepted_fixed | Failure Modes now require distinct reason codes for absent-family and checked-in contradiction. |
| S26-DX-01 | major | DX/replayability | Required Artifact Families looked mandatory for every invocation. | accepted_fixed | Required Artifact Families now states the table is product capability; selected profile controls required families. |
| S26-DX-02 | major | DX/replayability | `harness_observed` producer scope was ambiguous. | accepted_fixed | Added producer-scope glossary. |
| S26-DX-03 | minor | DX/replayability | CLI placement had no selection criterion. | accepted_fixed | Product Surface now states a CLI selection principle. |
| S26-DX-04 | minor | DX/replayability | Fixture filenames encoded expected verdicts and weakened replayability. | accepted_fixed | Implementation plan now requires `input/` and `expected/` split with scenario names. |
| S26-DX-05 | minor | DX/replayability | State vocabularies were not mapped to fields. | accepted_fixed | Added field vocabulary table. |
| S26-DX-06 | minor | DX/replayability | `pr_ci` was provider-centric. | accepted_fixed | Canonical family is now `change_ci`; `pr_ci` may be an alias. |
| S26-DX-07 | info | DX/replayability | Output-safety false positives need diagnosable safe rule ids. | accepted_fixed | Acceptance criteria and implementation plan require ruleset id/digest and safe rule-id failures. |
| S26-SEC-01 | major | security/privacy | Free-text `reasons` and `next_actions` could leak raw artifact content. | accepted_fixed | Observation Model and implementation plan require closed-code safe templating. |
| S26-SEC-02 | major | security/privacy | Redaction/output-safety ruleset was not specified or auditable. | accepted_fixed | Implementation plan requires documented digest-bound default safety ruleset. |
| S26-SEC-03 | minor | security/privacy | Repository/ref/run identifiers may leak private metadata when shared. | accepted_fixed | Non-goals now state observation results are same-boundary unless a future redaction profile is selected. |
| S26-SEC-04 | minor | security/privacy | Malformed fixture handling could echo unsafe input in CI output. | accepted_fixed | Implementation plan requires safe malformed fixture output. |
| S26-SEC-05 | minor | security/privacy | Future network artifact resolution could introduce unsafe URL/token behavior. | accepted_fixed | Non-goals constrain future network resolution and reject implicit network fetching in Block 26. |
| S26-SEC-06 | minor | security/privacy | Unsafe pattern list omitted common CI secret shapes. | accepted_fixed | Implementation plan expands default safety ruleset requirements to JWTs, private keys, cloud credentials, provider tokens, private URLs, and high-entropy values. |

## External Socratic Review

Initial external Socratic review was run through `pi`:

- product/customer credibility: MiniMax-M2.7, `CHANGES_REQUESTED`, fixed above;
- tracing/evidence: OpenRouter Qwen, `CHANGES_REQUESTED`, fixed above;
- DX/replayability: ZAI/GLM-5.1, `CHANGES_REQUESTED`, fixed above;
- security/privacy: OpenRouter DeepSeek, `CHANGES_REQUESTED`, fixed above.

Focused re-review was run after fixes:

- product/customer credibility: MiniMax-M2.7, `APPROVE`; no remaining critical
  or major findings;
- tracing/evidence: OpenRouter Qwen, `APPROVE`; no remaining critical or major
  findings;
- DX/replayability: ZAI/GLM-5.1, `APPROVE`; no remaining critical or major
  findings;
- security/privacy: OpenRouter DeepSeek, `APPROVE`; no remaining critical or
  major findings.

Block 26 remains a reviewed SpecKit direction. Implementation code still needs
explicit approval before WS1-WS5 begin.

## Implementation Review Findings

Implementation is in progress after explicit approval.

Current implementation evidence before external implementation review:

- Product surface: `assess --profile ci-artifact-observation
  --artifact-manifest <file> --out <file>` plus `assess preview` and
  `assess explain`.
- Go evaluator: `internal/ciartifact`.
- Schema: `schema/ci-artifact-observation.schema.json`, wired into
  `schema/assessment-result.schema.json`.
- Fixtures: `examples/block26-ci-artifact-observation/fixture-matrix.json` and
  valid uploaded-bundle manifest input.
- Docs: `docs/agent-entrypoint.md`.
- Local verification before implementation review:
  - `go test ./...`: pass.
  - `jq empty schema/*.json
    examples/block26-ci-artifact-observation/fixture-matrix.json
    examples/block26-ci-artifact-observation/input/ci-uploaded-bundle-complete-coverage/artifact-manifest.json`:
    pass.
  - `git diff --check HEAD`: pass.

External implementation review pending across code/correctness,
tracing/evidence, requirements-vs-implementation, and security/privacy planes.
