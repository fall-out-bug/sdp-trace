<think>The user wants me to review this document as an adversarial product/system reviewer. Let me read through the document carefully and identify critical issues.

This is a brief about "sdp-trace" - a telemetry trust layer for AI-assisted delivery. Let me analyze it critically.

Key issues I see:

1. **Adoption model is contradictory** - Section 7 says "Observe-only: no blocking" but Section 4 says "managed mode: verifier fails or blocks". These are incompatible promises.

2. **Trust levels are overstated** - "Strong for observed tool actions" but if the tool wrapper can be bypassed, it's not strong. The harness adapter is labeled "Medium if managed" which is vague.


3. **Anti-forgery claims lack rigor** - Hash linking prevents reordering but doesn't stop deletion or backfill attacks. "Strongest anti-forgery boundary" for remote witness is marketing, not proof.

4. **The "flight recorder" metaphor collapses under scrutiny** - A flight recorder in aviation is tamper-resistant by design; this system has no mechanism to prevent the agent from fabricating telemetry.

5. **Agent cooperation model is wishful thinking** - "Agents cannot simply be asked to write trustworthy telemetry" contradicts allowing agents to enrich telemetry with intent and rationale.

6. **Gateway capture is underspecified** - Section 3 mentions LLM gateway but provides zero details on how it actually prevents tampering.

7. **CTO query surface requires trust that hasn't been demonstrated** - Tasks aren't locked, only recorded. There's no mechanism preventing an agent from rewriting the task record after the fact. Supersession detection is unproven. All the queries depend on telemetry integrity that the system hasn't secured.

8. **Missing concrete enforcement mechanisms** - The anti-forgery section mentions "verifier must not repair these" but there's no explanation of what actually prevents forgery in practice.

9. **Minimum viable changes section is vague** - Section 10 lists "harness adapter contract" as a requirement but the brief never defines what constitutes a harness adapter or its minimum viable form.

10. **Gateway integration is a privacy landmine** - LLM gateway telemetry capturing prompt/response digests could expose sensitive data without any privacy safeguards outlined.

11. **CI witness assumes CI integrity** - Section 5 relies on CI or remote witnesses signing chain heads, but CI itself can be compromised or manipulated.

12. **Open-ended gaps become critical blockers** - The "still missing" section lists fundamental issues like the harness adapter contract and tool capture implementation that would make the whole system unworkable.

13. **"Harness-neutral" claim is unsubstantiated** - The brief never defines what harness neutrality actually means or how it's achieved.

14. **Verdict: CHANGES_REQUIRED** - The adoption model contradiction alone is enough to require rework before the review format can be properly filled out.

Looking at the remaining issues, there are fundamental problems with the threat model and adoption strategy. The agent cooperation model has conflicting requirements that can't both hold, the anti-forgery protections are overstated and incomplete, there's no clear specification for what a "minimum viable" integration actually looks like, and the gap between what the system claims and what's missing creates exploitable holes. The CTO query surface depends on the integrity of a telemetry pipeline that hasn't been secured, making everything downstream unreliable.
</think>

```text
VERDICT: CHANGES_REQUIRED

Critical blockers:

1. ADOPTION MODEL CONTRADICTION (adoption, anti-forgery):
   Section 7 says "Observe-only: attach recorder... no blocking." Section 4 says
   "managed mode: verifier fails or blocks." These are mutually exclusive promises.
   A CTO who deploys observe-only and sees no blocking cannot later claim anti-forgery
   enforcement. The product cannot be both optional wrapper and trust boundary
   simultaneously without explicit state transitions. You need a formal mode
   contract, not two paragraphs that contradict.

2. AGENT COOPERATION MODEL FAILS ITS OWN PREMISE (agent cooperation, telemetry capture):
   Thesis states: "Agents cannot simply be asked to write trustworthy telemetry."
   Section 4 then allows: "Agents may enrich telemetry with intent, rationale, plan,
   claims, self-report." You are asking agents to write telemetry. Your threat model
   assumes cooperative agents for enrichment but adversarial agents for signing. This
   asymmetry is unresolved. An agent that enriches intent can also redact intent.
   What stops self-reported enrichment from being the vector for redaction attacks
   listed in anti-forgery?

3. ANTI-FORGERY MODEL IS INCOMPLETE (anti-forgery):
   Lists threats correctly but provides zero mitigations for:
   - Event deletion (hash linking stops reordering, not deletion)
   - Backfilled telemetry (no timestamp authority before local recorder starts)
   - Partial telemetry with successful-looking summary (no completeness check)
   - Task rewrite after command evidence (no task immutability contract)
   "Verifier must not repair these" is the right instinct, but you must show what
   the verifier can detect, not just what it cannot repair. Currently the verifier
   emits states but cannot prove absence.

Major gaps:

4. NO MINIMUM VIABLE INTEGRATION DEFINITION (integration placement, adoption):
   Section 9 lists "tool/shell/file capture implementation" and
   "OpenCode/pi/Kilo integration plan" as missing. Section 3A lists 6 integration
   options with "review questions" but no decision tree. Section 3A ends with
   "Working hypothesis: layered, not either/or" with 4 bullet points that are
   architectural hand-waving, not an integration sequence. You cannot ask a CTO
   to buy "layered" when the first layer's API contract does not exist.

5. HARNESS ADAPTER CONTRACT UNDEFINED (integration placement, telemetry capture):
   Section 9 lists "harness adapter contract" as missing. This is not a minor gap.
   If you cannot define what a harness adapter must emit, you cannot claim any
   telemetry captured at the harness layer is verifiable. The trust strength column
   says "Medium if managed" but "managed" is undefined. The adapter is the primary
   integration point and it does not exist as a spec.

6. GATEWAY TELEMETRY PRIVACY MODEL MISSING (telemetry capture, integration placement):
   Section 3A raises "may capture sensitive prompts" as a risk. The architecture
   treats prompt/response digests as strong model provenance evidence. If a CTO
   runs an agent against proprietary codebase prompts, those digests are a data
   exfiltration vector. You have no privacy boundary defined. Signing independent
   of harness telemetry does not solve this; it separates the privacy violation
   from the evidence.

7. EXTERNAL WITNESS PROTOCOL UNDEFINED (signing, anti-forgery):
   "Remote append-only witness, timestamp service, or transparency log" are three
   different things with radically different trust properties. A transparency log
   (like CT Log) is auditable by third parties. A private append-only store is
   just another mutable database controlled by the vendor. The "strongest
   anti-forgery boundary" claim depends entirely on which one you mean. You cannot
   make this claim without a protocol.

False assumptions:

8. TASK LOCKING (CTO usefulness, anti-forgery):
   Section 8 asks "What task was locked?" as a required CTO query. The brief
   provides no task locking mechanism. A recorded task is not a locked task. The
   "task lifecycle" captured by a harness adapter is a log of what the agent
   said it was doing, not a cryptographically enforceable task boundary. An agent
   that rewrites its own task definition mid-run breaks every downstream query
   without detection.

9. VERIFIER STATE COMPLETENESS (anti-forgery):
   The anti-forgery section lists threats the verifier "must not repair" but does
   not enumerate what the verifier can detect. You have states: pass, fail,
   cannot_verify, not_assessed, missing_telemetry, local_only, witnessed,
   externally_witnessed. What verifier state fires when:
   - Events exist but are hash-broken?
   - Task was rewritten post-evidence?
   - Gateway telemetry exists with no local action telemetry?
   - Local action telemetry exists with no gateway provenance?
   "not_assessed" cannot be the sink for all unspecified conditions.

10. CI SIGNING TRUST (signing):
    "CI or remote witness signs chain head" is presented as a trust upgrade. CI
    can sign unverified telemetry (the brief names this threat). The mitigant is
    not stated. If CI runs an agent that produces telemetry and then signs that
    telemetry, you have signing-as-trust without verification. The brief
    acknowledges this in threats but does not address it in the signing model.

Minimum viable changes:

11. Define one integration point concretely (adoption, integration placement):
    Pick the minimum viable layer. Recommended: shell wrapper for commands
    and exit states, because bash is cross-harness. Write the schema for what
    a shell wrapper emits: cwd, argv digest, start time, finish time, exit code.
    Show that this layer can detect its own bypass (shell command not wrapped).
    This becomes the baseline for "strong for build/test evidence."

12. Write the mode contract (adoption, anti-forgery):
    Explicitly define observe-only vs. managed vs. witnessed modes with:
    - What triggers mode transitions
    - What verifier states are legal per mode
    - What claims can be made per mode
    - What a CTO can trust in each mode

13. Define harness adapter output schema (integration placement, telemetry capture):
    Even if imperfect, specify the minimum events the adapter must emit and
    the minimum metadata per event. This becomes the contract that "medium if
    managed" can be evaluated against.

14. Write the privacy boundary for gateway telemetry (telemetry capture):
    Either exclude prompt/response content from the trust model entirely
    (use only request-id correlation), or define a PII scrubbing boundary
    before digest. The brief currently treats gateway telemetry as strong
    provenance without resolving the privacy conflict.

15. Add a completeness invariant (anti-forgery):
    Define what "complete run" means in terms of required events. The verifier
    must be able to state whether a run is complete or partial. Currently it
    can only state pass/fail/per-event conditions, not run-level integrity.

Questions before demo:

16. If pi or OpenCode bypasses the adapter (logs a different tool call than
    what the agent executed), how does the verifier detect this? You say harness
    adapter captures "tool requested" but an adversarial harness can lie about
    what was requested.

17. How does the CTO query surface handle a run where all evidence is present
    but the agent's plan was wrong? "Task completed, all tests passed, scope
    unchanged, witnessed" still produced the wrong feature. Is that a trust
    failure or an engineering failure? The brief conflates delivery trust with
    correctness trust.

18. What stops a team from running the agent in unmanaged mode, then attaching
    sdp-trace as a wrapper after the fact to produce a clean report? You mention
    "late attach" but do not define the detection mechanism.

Attack scenarios not covered:

19. Silent bypass via subagent delegation:
    Harness adapter captures top-level agent lifecycle. Agent spawns a
    subprocess or subagent that performs the actual mutations. Main harness
    only sees "subagent completed" with no file/command evidence. This is
    the primary escape vector for any harness-level integration.

20. Test artifact poisoning:
    Local recorder captures test commands and exit codes. Agent or harness
    modifies test suite after evidence capture but before CI runs. CI sees
    passing tests but not the state of tests at recording time. No replay
    integrity protection exists.

21. Gateway correlation attack:
    LLM gateway logs a request-id. Agent performs a different action than
    what the prompt intended. Local recorder shows the action. Gateway proves
    the prompt was issued. No link exists between prompt intent and action
    outcome, but the merged report looks provenance-complete.

22. Witness chain fork:
    Remote witness anchors chain head at time T. Agent rewrites local events
    to change the chain content while keeping the same head hash (hash collision
    or deliberate mutation before anchor). Witness confirms "head existed"
    but does not confirm "head matched this event sequence."

23. Scope inflation post-evidence:
    Agent records allowed scope S1. Executes within scope. Captures evidence.
    Then rewrites task to redefine scope boundaries, making S1 a subset of S2.
    Verifier sees "changes inside allowed scope" because scope was redefined
    after evidence. No mechanism to detect task supersession post-evidence
    capture versus legitimate task evolution.
```
