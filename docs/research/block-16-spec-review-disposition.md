# Block 16 Spec Review Disposition

Date: 2026-05-06

Scope: `specs/001-sdp-trace-time-series-evidence-substrate/blocks/16-protected-gate-enforcement-profile.md`

## Review Runs

| Role | Model | Status | Disposition |
|---|---|---|---|
| Adversarial trust reviewer | MiniMax-M2.7 | not_counted | Attempted to request tools despite no-tools review mode. Replaced. |
| SpecKit consistency reviewer | ZAI/GLM-5.1 | counted | One major schema/CLI enum finding plus minor consistency findings. |
| Full reviewer | Kimi K2P6, low reasoning | counted | Critical determinism and coverage findings for state mapping, ordering, CLI behavior, and freshness. |
| Platform CI owner | DeepSeek V3.2 | counted | Critical operational findings for input handling, condition visibility, freshness, and override independence. |
| Developer DX reviewer | Qwen 3.5 Plus endpoint | not_counted | Endpoint returned 404. Replaced. |
| Compliance/audit reviewer | MiniMax-M2.7 Highspeed | counted | Major findings for explain/preview shape, schema migration, and override condition clarity. |
| Replacement DX reviewer | Qwen 3 Coder Plus | counted | Migration, input handling, explain/preview, trust-scope clarity findings. |
| Replacement adversarial reviewer | DeepSeek R1 0528 | counted | Freshness, external witness, override severity, local replay, and reason-code findings. |

Counted review planes: 6.

Default reviewer families avoided OpenAI, Anthropic/Claude, and Google models.

## Finding Disposition

| Finding area | Severity accepted | Disposition |
|---|---|---|
| Top-level `protected_gate` enum vs lower-level `not_integrated` | major | Updated spec: top-level `protected_gate` never emits `not_integrated`; lower-level `not_integrated` maps to `cannot_verify` unless a failure is present. |
| Missing protected-profile inputs | major | Updated spec: omitted required gate flags are usage errors; readable artifacts that cannot satisfy verifier produce `fail` or `cannot_verify`. |
| Default gate accidentally evaluating protected conditions | major | Added acceptance criterion that default gate without `--profile protected` does not evaluate protected conditions. |
| Condition row visibility under dominant failure | major | Updated spec: all protected condition rows remain visible even when top-level state is dominated by fail/cannot_verify. |
| Condition and reason ordering ambiguity | major | Updated spec: condition rows use fixed condition-id order; reasons and next actions use severity order with condition-id tie breaker. |
| Stable machine-readable reason codes | minor accepted as major for DX | Updated spec: condition rows include `reason_code`, and reason codes must be stable identifiers. |
| CI witness field-level mismatch mapping | major | Updated spec: absent required witness fields are `cannot_verify`; present contradictory fields are `fail`; empty digest list is `cannot_verify`; non-empty mismatch is `fail`. |
| Freshness absent, unbounded, stale, or contradictory | major | Updated spec: absent/unbounded/null/unknown freshness is `cannot_verify`; expired or contradictory freshness is `fail`; neither can satisfy protected profile. |
| Local-development signed checkpoint ambiguity | major | Updated spec: valid local-development checkpoint authority is protected `fail` because local signed trust is outside protected pass scope. |
| Override condition ambiguity and severity dominance | major | Added override state table and rule that override `cannot_verify` cannot reduce severity of prior failures or upgrade the profile. |
| Preview output shape | major | Updated spec: preview renders input inspectability statuses and next actions; it does not verify signatures, signer authority, or witness binding. |
| Explain output shape | major | Updated spec: explain renders selected profile, checkpoint summary, protected conditions, Block 14 conditions, and stable reason-coded next actions. |
| Block 14 migration/read compatibility | major | Updated spec: Block 16 requires new gate-result schema version and `gate explain` read compatibility for Block 14 artifacts. |
| Path-derived sensitive leakage | minor accepted | Updated implementation tasks to include path-derived sensitive output assertions. |
| External witness future risk | major partially accepted | Block 16 keeps external witness outside protected pass until a later approved verifier exists. External witness replay fixtures are deferred because external witness verifier is explicitly not integrated in this block. |

## Remaining Non-Blocking Notes

- Reviewers disagreed on whether absent freshness should be `fail` or
  `cannot_verify`. The spec keeps absent/unbounded freshness as
  `cannot_verify` because no contradictory evidence exists, but this remains
  fail-closed for protected CI use through exit `3`.
- External witness support remains `not_integrated`; Block 16 must not create
  aspirational external witness pass behavior.
- Review process requirements are intentionally not encoded in the Block 16
  product spec. They belong to operating workflow guidance.
