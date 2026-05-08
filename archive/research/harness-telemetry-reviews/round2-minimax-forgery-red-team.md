<think>The user wants me to review 11 attack scenarios from the perspective of a CISO adversary, classifying each against the current state of the sdp-trace design based on all four reviews. Let me go through each systematically.

The four reviews have identified these gaps:
1. No concrete cross-layer correlation mechanism
2. Local recorder signing has no attestation model
3. Verifier states are overloaded (verdict + trust scope conflated)
4. No mechanism to distinguish live-captured from post-hoc-written telemetry
5. Gateway correlation is temporal only, not causal
6. Key management entirely absent
7. No completeness contract for "expected events"


8. Gateway-only telemetry with no local action verification leaves a critical gap—gateway shows the model was called but not what actually happened, so this degrades to `cannot_verify` without a completeness contract and correlation protocol. Similarly, local-only telemetry without model provenance can't establish that the right model generated the output, also falling to `cannot_verify`. CI-signed but unverified telemetry is worse: a compromised pipeline can sign anything without re-verification, so this cascades to `cannot_verify` or `not_assessed` once the CI threat model is explicit.

Now I'm evaluating the remaining attacks against the same framework. Post-hoc runs and deleted runs both have no mitigation—they either look locally consistent or vanish entirely, leaving us with `not_assessed` for the first and `missing_telemetry` for the second. Event mutation fails the same way since hash linking detects changes but not deletions, so the chain appears broken and we get `cannot_verify`. Reordering is actually prevented though—hash linking plus timestamps catches any attempt to swap event order.

For replay attacks, a slow replay of an identical task can pass if there's no nonce binding to distinguish runs, making this undetectable without a witness. A stolen local key is trickier since the agent could regenerate a valid chain, but since the local recorder signing is local-only and not audit-grade anyway, it never claimed to prevent this. Wrapper bypass is where things get interesting—the harness might not detect it, but if there's a file watcher or git diff layer, it could still catch mutations. The gateway-only versus local-only gap is key here: I need a completeness contract where certain events must exist, and without both layers present, the verifier can't verify either way.
</think>

```text
CISO ADVERSARY ATTACK CLASSIFICATION

==============================================================
ATTACK: fake posthoc run
Status: NOT_COVERED (not_assessed by default)

The local hash chain is internally consistent whether written live
or manufactured post-session. No timestamp anchor, no monotonic
capture proof, no recorder process attestation. The brief names this
threat (Section 6) but provides zero detection mechanism. Under
AGENTS.md trust rule "source-bound proof requires clean immutable
commit", a posthoc chain fails only if the chain artifacts are also
manifest subjects. If they are not, the chain passes `local_only`
without degradation.

Specific failure: no wire protocol distinguishes "recorder wrote
event N at time T" from "agent wrote event N into chain at time T+Δ".
No replay-of-recorder-proof exists.

Mitigation absent: requires monotonic external time source injection
at capture time or witness-backed chain head before session closes.
Neither exists in current design.
==============================================================

ATTACK: deleted run
Status: NOT_COVERED (missing_telemetry emitted, not prevented)

A deleted run produces zero telemetry. Verifier sees an absence.
The correct output is `missing_telemetry` with no ability to prove
deletion versus never-started. This is honest behavior—the brief
explicitly marks missing telemetry as a verifier state—but it means
a deleted run cannot be distinguished from a never-attached run, and
the product provides zero deterrence against deletion.

CTO impact: "missing_telemetry" is the only signal. It cannot answer
"was this run deleted or was it never recorded" without a witness
that observed the run's existence.

Mitigation absent: remote witness provides anti-backfill only for
runs it has seen. It does not provide existence-guarantee for runs
that never attached.
==============================================================

ATTACK: event mutation
Status: DOWNGRADED_TO cannot_verify

Hash linking detects that the chain has changed. The verifier sees
a broken hash chain and emits `cannot_verify`. It does not repair
the events. Correct behavior, but the attack succeeds at the goal:
the verifier cannot assess the mutated run. The result is a
`cannot_verify` state that tells the CTO "something happened but we
cannot confirm what."

Critical gap: the brief does not specify whether the verifier
produces a partial assessment (which events are intact vs broken) or
a binary `cannot_verify` with no event-level granularity. No partial
assessment model is defined.

Mitigation status: detected (hash break), not repaired (correct),
but outcome is an uninformative block state, not a security alert.
==============================================================

ATTACK: reordering
Status: PREVENTED

Hash linking with timestamp ordering prevents in-chain reordering.
If an event's position in the chain changes, its hash predecessor
changes, breaking the chain. The brief correctly handles this
(Section 6 threat named correctly, hash linking is the right
countermeasure).

Residual gap: hash linking works within the local chain. If events
are deleted and re-inserted with new timestamps, the chain looks clean
but the event set is corrupted. This is event deletion + insertion,
not pure reordering, and falls to the "deleted run" case above.
==============================================================

ATTACK: replay old valid run
Status: DOWNGRADED_TO cannot_verify (slow replay) / NOT_COVERED (fast replay with nonce spoof)

Without a unique-per-run binding, a task with the same content hash
can be re-executed and replayed as a new run. The brief names this
threat (Section 6) but provides no nonce generation protocol, no
task-unique-identifier derivation, and no binding integrity mechanism.

Slow replay (same task, later time): detected only if remote witness
uses a nonce. No nonce protocol exists → replay passes as a new run
until witnessed. Result: `local_only` for unattached replay runs.

Fast replay (same session, before chain closes): no mechanism to
prevent the same events from being recorded twice or replayed within
the same chain.

Mitigation status: partially covered only for externally-witnessed
runs where the witness injects a per-run nonce. All other runs are
vulnerable.
==============================================================

ATTACK: stolen local key
Status: DOWNGRADED_TO not_assessed (by design admission)

The brief correctly states local-only signing "cannot provide
audit-grade trust" and that "local recorder signing is not audit-grade
if user controls host" (signing model, Section 5). This attack is
therefore within the acknowledged scope: local-only telemetry is never
claimed to resist a local key compromise.

The attack succeeds at full chain forgery because the agent can:
1. Obtain the local signing key
2. Re-sign the entire modified chain
3. Present a locally-verifiable chain that checks out

Result: `local_only` with `pass` from local verifier, with no
detectable signal that the key was compromised. This is "not_assessed"
by explicit design admission, not a gap.

CTO warning required: the brief must explicitly state that local-only
runs are not trust-grade and must never be described as "verified"
in gate decisions. Current prose does not enforce this distinction.
==============================================================

ATTACK: wrapper bypass
Status: DOWNGRADED_TO cannot_verify (detection gap) / DOWNGRADED_TO not_assessed (blind spot)

Two sub-cases with different outcomes:

Sub-case A: Agent uses a subprocess or inline execution that bypasses
the tool/shell wrapper but the file watcher still observes mutations.
The run produces: gateway telemetry (if attached), file watcher events
but no shell/command events, and tool-level events. The verifier sees
an incomplete event set. Without a completeness contract, this emits
`cannot_verify`. The gap is detectable but the attack is not prevented.

Sub-case B: Agent bypasses the wrapper AND mutates files through a
channel the file watcher does not observe (e.g., memory-mapped writes,
atomic swaps as noted by DeepSeek review). The run produces gateway
telemetry only. Result: `local_only` (no local events) + `cannot_verify`
(because local events expected per contract are absent). The verifier
has no signal to distinguish "no agent action" from "action without
local trace."

The wrapper bypass problem is the primary escape vector. The
integration placement table rates tool wrapper capture as "Strong
for observed tool actions" but does not rate it for "unobserved tool
actions." Strength is only stated for the happy path.

Minimum required counter: a completeness contract stating which events
are required vs optional per run type. Without this, bypass is
indistinguishable from a short run.
==============================================================

ATTACK: gateway-only telemetry
Status: DOWNGRADED_TO cannot_verify

Gateway telemetry proves a model call occurred. It does not prove
the agent acted on the response, or that the response caused file
changes, or that any command ran. The brief names "gateway telemetry
without local action telemetry" as a threat (Section 6) but provides
no correlation protocol or causal link.

The verifier must emit `cannot_verify` for local actions that lack
local telemetry, or must define a gateway-only trust path. Currently
it does neither. The only legal output without local action telemetry
is `cannot_verify` per the trust rules ("machine proof wins").

Residual risk: a run with only gateway telemetry will appear
provenance-complete to a human reading a merged report because the
gateway events look authoritative. The brief has no display-layer
warning that "gateway-only runs cannot support gate decisions."
This display-layer overclaim is the actual attack surface.
==============================================================

ATTACK: local-only telemetry without model provenance
Status: DOWNGRADED_TO cannot_verify

The inverse of the above: file and shell events prove actions occurred
but cannot prove which model was used, whether a model was called at
all, or whether the agent's reasoning influenced the outcome. An
agent that runs commands without model involvement (or uses an
unapproved model not routed through the gateway) produces local-only
telemetry that looks provenance-complete.

Result: `local_only` with no model identity. This is not `pass`.
The display layer must show `local_only` clearly, not treat it as
equivalent to a full-provenance run.

The "unapproved harness" attack (the 11th attack) is a specific case:
if the harness is unapproved and produces no gateway telemetry, the
CTO query "is the team drifting into unapproved harnesses" produces
a silent negative answer—the system shows no harness data, not
"unapproved harness detected." This is an absence-of-evidence-as-good
failure mode that directly defeats the CTO degradation question.
==============================================================

ATTACK: CI signs unverified telemetry
Status: DOWNGRADED_TO cannot_verify (correct) / NOT_COVERED (signing model gap)

CI signing is currently "not_assessed" by the brief's own admission
(Section 9: "CI signing recipe" is listed as missing). The signing
model names this threat correctly but provides no mitigation.

The correct behavior is: CI must re-run the verifier and include the
verifier output in what it signs. CI signs: (chain head, verifier
output, build/test evidence). Without verifier re-execution in CI,
the signing model allows CI to sign anything, which is the attack.

The brief does not specify this. Until it does, any CI-signed run
using the current model is `cannot_verify` because the CI signature
covers an unverifiable payload.

Additional gap: CI pipeline can be misconfigured or compromised.
The brief does not define a CI integrity attestation model (what
prevents a malicious CI job from running the agent and signing its
own chain?). The "CI identity" in the signing table is assumed to
be trustworthy, which it is not.
==============================================================

ATTACK: redaction hides evidence
Status: NOT_COVERED (detectable if witnessed, not detectable if local-only)

Two sub-cases:

Sub-case A: Redaction occurs before remote witness attaches. Witness
anchors the pre-redaction chain head. Post-redaction chain head
diverges. Witness does not detect divergence (it anchored earlier
head). The run appears as two valid chains with different heads,
one of which is witnessed. Without a full chain comparison protocol,
the verifier cannot identify which chain is the redacted one.

Sub-case B: Redaction occurs in local-only run. The local chain is
mutated (see event mutation). Verifier detects broken hash chain,
emits `cannot_verify`. This is the correct detection signal but
the result is an opaque block.

Sub-case C: Legitimate redaction (agent removes sensitive content).
The brief does not define a redaction policy—what makes redaction
admissible vs adversarial. An agent can claim any redaction is
legitimate privacy action. The verifier has no policy to evaluate
this claim.

The redaction state exists in Block 09 but the detection logic and
policy are not defined. Current state: not covered for local-only,
partially covered for witnessed (requires chain comparison protocol
that is not specified).
==============================================================

ATTACK: unapproved harness
Status: NOT_COVERED (absence-of-evidence failure mode)

The 11th attack case: a team member runs an agent through an
unapproved harness. No telemetry is produced. The verifier emits
`missing_telemetry` or `not_assessed`. This is indistinguishable from
a clean run through an approved harness with perfect telemetry.
The system cannot distinguish "no harness" from "harness with no
telemetry."

CTO impact: The query "is the team drifting into unapproved harnesses,
agents, prompts, or skills?" has no answer. The system shows absence
of unapproved harness evidence as absence of evidence, not as a
detection event. A CTO looking at a dashboard sees all approved runs
with clean telemetry. The unapproved run never appears.

Mitigation would require: a harness registry, a telemetry mandation
check, and a passive or active scan for unapproved harness activity.
None of these exist. The brief lists "harness registry / approved
agent governance" as missing and explicitly acknowledges this gap.

This is the attack that most directly defeats the stated CTO value
proposition. Every other gap matters less than this one.
==============================================================

SUMMARY TABLE
==============================================================
Attack                    | Classification          | Ruling
--------------------------|-------------------------|-------
fake posthoc run          | NOT_COVERED             | not_assessed
deleted run               | NOT_COVERED             | missing_telemetry
event mutation            | DETECTED → degrades     | cannot_verify
reordering                | PREVENTED               | (no action needed)
replay old valid run      | PARTIALLY_COVERED      | cannot_verify / not_assessed
stolen local key          | NOT_ASSESSED by design  | not_assessed (admitted)
wrapper bypass            | PARTIALLY_COVERED      | cannot_verify / not_assessed
gateway-only telemetry    | DOWNGRADED             | cannot_verify
local-only no model prov  | DOWNGRADED             | cannot_verify / local_only
CI signs unverified       | NOT_COVERED            | cannot_verify (correct output,
                          |                         | but signing model is vapor)
redaction hides evidence  | NOT_COVERED            | cannot_verify (witnessed) /
                          |                         | not_assessed (local-only)
unapproved harness        | NOT_COVERED            | absence-of-evidence
                          |                         | failure mode
==============================================================

CRITICAL THRESHOLD: 7 of 11 attacks are NOT_COVERED or
PARTIALLY_COVERED. The product cannot make any trust-grade claims
until at minimum: completeness contract, correlation protocol,
display-layer enforcement of trust-scope labels, and CI signing
recipe are defined. The two attacks that most directly defeat the
CTO use case—gateway-only without local, and unapproved harness
with no telemetry—are covered only by "not_assessed." The CTO
dashboard will show clean runs for both attack cases unless the
display layer is hardened separately.
```
