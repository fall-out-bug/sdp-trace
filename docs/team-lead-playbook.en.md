# Team Lead Playbook

Use `sdp-trace` when your team already has an AI coding workflow and needs a shared quality contract.

## Daily Use

1. Define the scope.
2. Capture provenance: human, agent, model, tools, and commands.
3. Attach evidence: tests, CI, review comments, files, and diffs.
4. Record accountability: human-held DRI, approver, risk owner, and escalation path.
5. Package an assessment input with evidence, observations, movement data, and `not_assessed` gaps.
6. Record any gate verdict as an external verdict input produced by `sdp-gate` or another policy consumer.

## Team Defaults

Agree on:

- required evidence per change type
- what external policy blocks merge
- who may approve or override in the policy layer
- which harnesses are supported
- what `not_assessed` means

## Reading External Verdicts

External verdicts may use values such as `pass`, `warn`, `fail`, or `not_assessed`, but they are not native `sdp-trace` decisions.

`not_assessed` is not a pass. Missing evidence must stay visible in the assessment input.
