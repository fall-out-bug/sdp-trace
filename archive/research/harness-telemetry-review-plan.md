# Harness Telemetry Review Plan

Status: discussion draft; not committed
Date: 2026-05-05
Input brief: `archive/research/harness-telemetry-trust-brief.md`
Output directory: `archive/research/harness-telemetry-reviews/`

## Round 1: Concept Attack

Purpose: find product and trust blockers before designing the demo.

Models / roles:

- MiniMax M2.7: CTO buyer + trust-boundary product critic.
- ZAI GLM 5.1: platform architect + schema/verifier consistency critic.
- DeepSeek v4-pro: implementation realism + operational failure-mode critic.
- MiMo v2.5-pro: constrained-model / low-cost-agent critic.

Prompt:

```text
Review the attached Harness-Neutral Telemetry Trust Layer brief.

Role: <role>.

Return only:

VERDICT: CHANGES_REQUIRED | ACCEPTABLE_WITH_GAPS | REJECTED

Critical blockers:
- ...

Major gaps:
- ...

False assumptions:
- ...

Minimum viable changes:
- ...

Questions before demo:
- ...

Attack scenarios not covered:
- ...

Do not produce generic architecture. Tie every finding to telemetry capture,
integration placement, agent cooperation, signing, anti-forgery, adoption, or
CTO usefulness.
```

## Round 2: Forgery Red Team

Purpose: attack signing and anti-forgery.

Roles:

- CISO adversary: malicious developer and compromised agent.
- Platform attacker: bypasses wrapper, replays old runs, tampers with gateway ids.
- CI attacker: tries to make CI sign unverified local telemetry.

Required attack cases:

- fake run generated after failure;
- deleted run;
- event mutation;
- event reordering;
- replay of old valid run;
- stolen local recorder key;
- agent bypassing tool wrapper;
- gateway-only telemetry with no local action evidence;
- local action telemetry with no model provenance;
- CI signs chain head without verifying event chain;
- redaction hides critical evidence;
- team runs unapproved harness.

Output must classify each as:

- prevented;
- detected;
- downgraded to `cannot_verify`;
- downgraded to `not_assessed`;
- not covered.

## Round 3: Adoption / DX

Purpose: pressure gradual adoption and developer acceptance.

Reviewer roles:

- CTO buyer: refuses harness replacement.
- Platform owner: owns wrappers, CI, secrets, gateway, repo templates.
- Staff engineer: rejects slow, leaky, flaky, or ceremony-heavy tooling.
- Forensics lead: must reconstruct incident after the fact.

Questions:

- What is the minimum integration that gives visible value in one week?
- Which integration point should be first: OpenCode/pi/Kilo plugin, LLM gateway, local wrapper, CI?
- Where will teams resist adoption?
- What telemetry is unsafe to capture raw?
- What must be queryable for CTO and incident review?
- What must stay out of scope before the Kotlin+Bazel demo?

## Round 4: Consolidation

Purpose: turn model criticism into decisions.

Human consolidation table:

| Finding | Source | Severity | Disposition | Reason | Demo Impact |
| --- | --- | --- | --- | --- | --- |
| ... | ... | critical/major/minor | accepted/rejected/deferred | ... | blocker/demo requirement/future |

No repository source changes are made in these rounds unless explicitly requested later.
