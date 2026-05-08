# Block 09 Socratic Review: Flight Recorder Trust Kernel

Status: review synthesis; implementation not started
Parent: `09-flight-recorder-trust-kernel.md`

## Review Setup

The concept was reviewed adversarially through `pi` with executive roles:

- CIO
- technical executive
- CFO
- COO
- CPO
- CISO
- Head of Engineering
- Corporate Architect

Model routes used included:

- `zai/glm-5.1`
- `minimax/MiniMax-M2.7`
- `kimi-coding/k2p6`
- `openrouter/deepseek/deepseek-v4-pro`
- `openrouter/deepseek/deepseek-v4-flash`
- `openrouter/qwen/qwen3.6-plus`
- `openrouter/xiaomi/mimo-v2.5-pro`

Some outputs attempted to inspect or invent local implementation files despite read-only/no-tool prompting. Code-specific claims from those outputs are not treated as evidence. Only converged product and architecture criticisms are recorded here.

## Converged Verdict

All usable executive reviews returned `revise`.

No role accepted a local-only sidecar hash chain as sufficient for accountability or technical executive/gate use.

## Critical Findings

### F09-C01: Local Hash Chain Is Not a Trust Boundary

Finding:

A local JSONL hash chain proves internal consistency only. If the recorded actor can delete, replace, or recompute the chain, the chain is not an accountability artifact.

Required resolution:

Block 09 must define a witness boundary and a verifier state that distinguishes local consistency from witnessed recording.

### F09-C02: Mid-Flight Attachment Is a Gap, Not Proof

Finding:

Mid-flight attachment is useful only when it honestly marks pre-attachment history as unassessed. It cannot support claims about what happened before attachment.

Required resolution:

Late attach must create a visible `not_assessed` boundary. Profiles that require full-run proof must fail if attachment occurs after the declared start.

### F09-C03: Voluntary Recording Is Bypassable

Finding:

If using the recorder is optional, the actors who most need oversight can avoid it. The recorder should not pretend to solve this.

Required resolution:

Block 09 must separate recorder proof from external enforcement. `sdp-trace` may emit missing/tampered trace evidence; a downstream gate decides whether that blocks work.

### F09-C04: Redaction Can Destroy Either Safety or Proof

Finding:

Storing raw secrets in an append-only log is unsafe. Redacting before hashing can make the proof weaker. Redacting after hashing may preserve evidence of secret exposure or make deletion impossible.

Required resolution:

Define redaction states, redaction authority, payload retention modes, sealed raw references, and verifier behavior for unresolved redaction.

### F09-C05: Digest-Only Evidence Is Often Not Forensic

Finding:

A digest of stdout, stderr, or diff proves a payload existed only if the payload can be retrieved or independently compared. Digest-only logs are too weak for reviewer reconstruction.

Required resolution:

Define evidence retention modes and require richer retention for profiles that claim forensic usefulness.

### F09-C06: Raw JSONL Is Not an Executive Product Surface

Finding:

technical executive, CISO, COO, and incident reviewers need queryable answers, not manual JSONL inspection.

Required resolution:

Block 09 must define query surfaces for run summary, gaps, task changes, command evidence, file mutations, tests, redactions, and witness state.

## Socratic Questions That Must Be Answered by Implementation

1. Can a developer produce a valid-looking run log for a run that never happened?
2. If the local log is deleted, where is the witness record that proves the run existed?
3. If the recorder attaches late, how does a reviewer see the pre-attachment gap without reading prose?
4. If a requirement changes after a failed command, can the original task be rewritten?
5. If stdout contained a secret, what exactly is hashed, retained, redacted, and verifiable?
6. If the model identity is a string provided by the harness, what is the verifier state: `pass`, `not_assessed`, or `cannot_verify`?
7. What does a technical executive learn from the trace that cannot be reconstructed from git history and CI logs alone?
8. How does the recorder report its blind spots without turning them into soft success?
9. What is the overhead of command/file capture on a real Kotlin+Bazel run?
10. What downstream gate can consume missing/tampered trace evidence without `sdp-trace` becoming the policy engine?

## Demo Implications

The Feature Flag / Entitlements Kotlin+Bazel demo must not start until Block 09 has a recorder kernel.

The demo must later prove:

- OpenCode + GSD + MiniMax 2.5 can be observed without modifying the product repo to depend on `sdp-trace`.
- Tampering with an event breaks verification.
- Changing the witness entry breaks witnessed verification.
- Late attach produces visible `not_assessed` scope.
- Requirement changes are superseding events.
- Redaction preserves safety and verifier honesty.
- Reviewer queries answer practical questions in minutes.
- Kotlin+Bazel overhead is measured, not assumed.

## Disposition

Proceed to implementation planning only after the Block 09 design records:

- local vs witnessed profile split
- event schema and canonicalization
- witness contract
- evidence retention modes
- redaction verifier semantics
- late attach semantics
- query surface
- negative fixtures
