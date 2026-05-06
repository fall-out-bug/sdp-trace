```text
VERDICT: CHANGES_REQUIRED

Critical blockers:
- The brief introduces at least 7 new verifier states (`missing_telemetry`, `local_only`,
  `witnessed`, `externally_witnessed`, plus reuse of existing 4) without defining a
  state machine: permitted transitions, mutual exclusions, composition rules when
  multiple layers report conflicting states for the same event. The current Slice 1
  validator has a narrow, fixed evidence vocabulary. Adding states without a schema
  that the baseline verifier can actually check will reproduce the overclaim problem
  called out in AGENTS.md trust rule 1. (anti-forgery, telemetry capture)
- Section 3/3A proposes "hash linking" for local event chains and "chain head"
  signing but specifies no hash algorithm, no chain structure, no fork/merge
  semantics, and no tolerance for out-of-order events from async tool wrappers.
  A hand-waved hash chain is worse than no chain because it creates a false trust
  signal. (signing, telemetry capture)
- The signing model (Section 5) says "if the agent can access the signing key,
  signature does not prove truth" but then proposes a "local recorder-signed" level
  where the recorder runs on the same host as the agent. No threat boundary is
  drawn between agent process and recorder process. If both are user-space on the
  same machine, the local recorder signature is agent-accessible by definition.
  (signing, anti-forgery)

Major gaps:
- Section 3A lists "LLM gateway / proxy" capture including "prompt hash, response
  hash" but Section 9 acknowledges no privacy model exists. Prompt hashes are
  invertible for short prompts; full prompt capture creates legal/privacy liability
  in enterprise settings. The brief must commit to a position: digest-only at the
  gateway, or no capture of prompt bodies, before this is reviewable.
  (integration placement, telemetry capture)
- No schema artifact is specified. The brief proposes event types, evidence roles,
  trust scopes, and source declarations per event, but gives no JSON schema, no
  required fields, no versioning strategy. Without a concrete schema, the
  "harness-neutral" contract is untestable. (telemetry capture, adoption)
- Section 4 says agent-written telemetry is "allowed for" plan/rationale/claims but
  gives no mechanism to distinguish agent-written from recorder-written events at
  the schema level. If both are in the same event stream, a consumer cannot apply
  different trust weights without a mandatory `source` discriminator that the
  verifier enforces. (agent cooperation, anti-forgery)
- Section 8 lists 12+ CTO queries but the brief has no query surface schema, no
  index model, and no defined relationship between query output and verifier state.
  This is a product wishlist, not an architecture. (adoption, CTO usefulness)
- The adoption ladder (Section 7) implies the verifier output changes meaning per
  rung ("missing telemetry becomes local fail" at rung 2) but does not define how
  the verifier knows which rung it is operating in, who sets the rung, or how rung
  transitions are governed. (adoption, anti-forgery)
- No subagent/session boundary model. OpenCode, pi, and Kilo all support
  subagents, parallel sessions, and retry loops. The brief mentions "subagent
  boundaries" in a table cell but provides no event structure for session nesting,
  parent-child correlation, or causality across retry boundaries.
  (integration placement, telemetry capture)

False assumptions:
- "Agents cannot simply be asked to write trustworthy telemetry" (Section 2) is
  correct, but the brief then assumes recorder components *can* be trusted to
  capture events independently. If the agent controls the tool layer, the shell
  layer, and the file system, a local recorder is co-located with the adversary.
  The threat model must be explicit about which threats are in scope for local-only
  vs. CI-witnessed vs. externally-witnessed runs. (anti-forgery, signing)
- The brief assumes "the architecture should be layered, not either/or" but does
  not address the combinatorial explosion of verifier states when 4+ layers each
  independently report pass/fail/missing. A single event could plausibly produce
  dozens of compound states with no defined precedence. (anti-forgery)
- Assumes harness plugins can capture "tool calls" reliably. In practice, pi and
  OpenCode tool calls go through different code paths (shell, file edit, search),
  and some tool implementations invoke subtools internally. The "tool call" boundary
  is not a stable abstraction across harnesses. (integration placement)

Minimum viable changes:
- Define a concrete event schema (JSON Schema) with mandatory fields: `event_id`,
  `timestamp`, `source` (agent | recorder | gateway | ci | witness), `evidence_role`
  (observed | self_reported | derived | missing), and `verifier_state` from a
  closed enum. Validate it with the existing baseline verifier.
- Define the verifier state machine: closed enum of states, permitted transitions,
  composition rules for multi-layer conflicts. Commit to no more than 8 states
  total; justify each.
- Remove or defer gateway prompt/response capture to a follow-up brief. Commit
  to digest-only with stated algorithm for Slice 2.
- Add a concrete threat model table: threat actor (agent, user, CI compromise,
  external), capability, in-scope/out-of-scope, which signing level is required
  to resist.
- Define the adoption rung as a governed field in the run manifest, not prose.
  The verifier must read it and adjust pass/fail thresholds accordingly.
- Add a subagent/session correlation schema: `session_id`, `parent_session_id`,
  `retry_of` (optional), `harness` (required).

Questions before demo:
- Which specific harness (pi, OpenCode, Kilo) will the first adapter target, and
  does that harness expose the tool-call and subagent APIs the brief assumes exist?
- What is the minimum event count for a "task completed" chain that the verifier
  would pass? Without a lower bound, a chain with one event can trivially pass.
- Who is the intended deployer of the local recorder: the developer, the team
  lead, or the CI pipeline? The trust assumptions differ for each.
- Has any harness maintainer confirmed that their plugin API can emit the events
  this brief requires? If not, the "harness-neutral" claim is premature.
- What happens when two layers disagree: gateway says model X was called, tool
  recorder says no file edits happened. Which state does the verifier emit?

Attack scenarios not covered:
- Agent spawns a child process that bypasses the tool/shell wrapper entirely,
  then the agent reports the wrapper-captured events as the full story. No layer
  detects the gap because the file watcher only sees the final state, not the
  process tree. (anti-forgery, telemetry capture)
- Recorder process is killed mid-run; agent continues and commits. Partial chain
  is indistinguishable from a short legitimate run. (signing, anti-forgery)
- Agent produces a valid chain, then amends the git commit and regenerates the
  file watcher digest to match. If the remote witness only sees the chain head
  hash, it cannot detect that the underlying events were rewritten before the
  chain head was recomputed. This is the "source-bound proof" problem from
  AGENTS.md applied to telemetry chains. (signing, anti-forgery)
- Two harnesses run concurrently on the same repo. Events interleave in the
  file watcher and shell wrapper. No harness_id session boundary separates them.
  (telemetry capture, integration placement)
- Agent redacts "sensitive" events from the chain after capture, re-links hashes,
  and presents the clean chain to the verifier. No append-only guarantee is
  enforced locally. (anti-forgery, signing)
- CI signs the chain head without re-running the verifier itself. A compromised
  or misconfigured CI pipeline becomes a trust stamping machine for garbage.
  (signing, adoption)
```
