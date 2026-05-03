# sdp-trace Accountability Model

`sdp-trace` records accountability so AI-assisted delivery remains governable.

## Rule

AI actors may produce, review, critique, or judge artifacts. They cannot be the sole accountable owner, approver, risk owner, or escalation owner.

Accountable identities are machine-readable:

```json
{
  "identity_ref": "role:sdp-trace-release-captain",
  "actor_type": "human_role"
}
```

Allowed accountable actor types:

- `human_user`
- `human_role`
- `human_group`

## Three Lines Mapping

| Line | Responsibility in `sdp-trace` |
|---|---|
| First | Delivery team owns execution evidence and recording quality. |
| Second | Contract, risk, or gate owner challenges controls and policy fit. |
| Third | Independent assurance samples whether the process works as claimed. |

`sdp-trace` records line-of-defense and identity facts. It does not decide whether separation is sufficient; that policy belongs to `sdp-gate` or another governance process.

## Effective Accountability

Evidence can inherit accountability from a containing package, but assessment inputs cannot claim completeness or trusted-release readiness if referenced evidence has no direct or inherited human accountability.

## Accountability Claim

`accountability_claim` narrows what the owner is accountable for:

- `recording_only`: accountable for accurate recording, not source truth
- `content_approval`: accountable for reviewed content
- `risk_acceptance`: accountable for accepted residual risk
- `release_approval`: accountable for releasing the contract artifact

