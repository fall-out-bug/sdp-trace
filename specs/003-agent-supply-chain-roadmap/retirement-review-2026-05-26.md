# Retirement Review: Agent Supply Chain Roadmap

Date: 2026-05-26
Reviewer: Codex GPT-5, sdp-trace closure route
Scope: Spec 003 task ledger closure as stale planning, not implementation
closure.

## Decision

Retire Spec 003 as `retired_superseded`.

Spec 003 was a roadmap/discovery artifact. It was reviewed and found not ready
for roadmap approval. Later specs and blocks implemented or narrowed the useful
parts through concrete product contracts, packet behavior, review evidence, CI
workflow evidence, OSS compatibility probes, and governance ledgers.

This retirement does not claim Spec 003 was implemented as product. It prevents
the stale roadmap from remaining an active implementation source.

## Supersession Map

| Spec 003 area | Current owner |
| --- | --- |
| CTO packet shape and evidence rows | Specs 005 and 006 |
| Evidence theater rows and reason codes | Spec 006 |
| GitHub-first demo packet workflow | Spec 007, still approval/demo-repo gated |
| OpenCode/GSD demo proof | Blocks 24, 25, and 31, still demo/first-run gated |
| `pi` runner evidence | Block 30 and Block 32 PR-review evidence |
| Signed / witness profile concepts | Witness, checkpoint, release proof, and OSS compatibility blocks |
| Repo governance / closure status | Specs 015, 018, and 019 |

## Review Disposition

The 2026-05-10 Socratic review returned `REVISE_BEFORE_USER_REVIEW`. The
valid direction was preserved, but approval did not happen on Spec 003 itself.
The closure-route decision is to retire the roadmap rather than resurrect a
large abstract plan.

Open review findings are dispositioned as:

- `retired_superseded` when later concrete specs own the product surface;
- `not_assessed` where the related later spec is still open, especially demo
  repository and first-run observation work;
- `out_of_scope` for general employee monitoring and broad unsupported tool
  claims.

## Boundaries

- No implementation slice is approved from Spec 003.
- No broad GitHub, GitFlic, GSD2, Superpowers, Hermes, OpenClaw, Sigstore, or
  customer-PKI support claim is made by this retirement.
- Demo work remains open in Spec 007 and Block 25 / Block 31 tasks.
- Approval and merge gates in other specs remain independent and are not
  satisfied by this retirement.
