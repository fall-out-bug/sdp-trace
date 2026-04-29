# Team Lead Playbook

Use `sdp-trace` when your team already has an AI coding workflow and needs a shared quality contract.

## Daily Use

1. Define the scope.
2. Capture provenance: human, agent, model, tools, and commands.
3. Attach evidence: tests, CI, review comments, files, and diffs.
4. Run or record gate checks.
5. Publish a verdict: `pass`, `warn`, `fail`, or `not_assessed`.
6. Record the decision and any override reason.

## Team Defaults

Agree on:

- required evidence per change type
- what blocks merge
- who may override
- which harnesses are supported
- what `not_assessed` means

## Reading Verdicts

- `pass`: evidence satisfies the gate.
- `warn`: evidence is incomplete or risk exists, but the team may proceed.
- `fail`: gate criteria are not met.
- `not_assessed`: the gate could not make a defensible decision.

`not_assessed` is not a pass.
