# Socratic Review Prompt: GitHub OSS Demo Packet

You are an independent Socratic reviewer for `sdp-trace`.

Your task: challenge whether the 007 plan is a sound product decision.

Context:

- Buyer: CTO.
- First product artifact: Change Evidence Packet v0.
- User explicitly prefers GitHub + OSS ecosystem for the demo.
- User does not want new demo repositories created merely to avoid the messy
  state of the existing demo repo.
- Existing demo repo is not accepted as finished proof; it is accepted as real
  source material.
- Missing evidence must remain `not_assessed` or `cannot_verify`.
- This is spec/product review only. Do not recommend implementation code.

Review lenses:

1. Product proof: will this produce a CTO-readable proof artifact or just more
   substrate?
2. Demo truth: does using the existing repo create unavoidable confusion, or is
   the mess itself useful evidence?
3. GitHub evidence: does the plan bind issues/PRs/checks/artifacts/reviews
   tightly enough?
4. Theater: does the negative PR/fixture demonstrate misleading evidence
   without poisoning the happy path?
5. Scope: does this avoid enterprise/self-hosted and broad OSS support claims?

Return exactly:

```text
Verdict: <APPROVE_FOR_USER_APPROVAL|REVISE_BEFORE_USER_APPROVAL|KILL_DIRECTION>

Findings:
| id | severity | file/section | finding | exact fix |

Existing repo default:
<direct answer>

Smallest first slice:
<one paragraph>

Claims to refuse:
<bullets>
```
