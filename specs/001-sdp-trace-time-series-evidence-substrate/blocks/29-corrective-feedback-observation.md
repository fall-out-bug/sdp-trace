# Block 29: Corrective Feedback Observation

## Problem

Dogfood exposed a trust gap: human corrective feedback to an agent or workflow
can change the delivery path, but current repo observer surfaces only show git,
hook, CI, and prompt-cooperation boundaries. The conversation itself is not
repository evidence, and prompt cooperation is explicitly `agent_reported`.

## Scope

Add a portable local observation surface for corrective feedback and observer
notes.

In scope:

- `sdp-trace observe feedback` records a retained message from a file into a
  JSON event.
- Feedback events carry digest, actor/target safe IDs, summary, retained body,
  `trust_scope=local_structural`, and `proof_state=not_assessed`.
- Repo observer doctor reports a `feedback_events` surface when valid feedback
  event files exist under `.sdp-trace/feedback`; malformed event files are
  `cannot_verify`.
- Documentation states that feedback events are local structural observations,
  not external proof.

Out of scope:

- Automatic capture from Codex, OpenCode, GSD, Slack, email, or chat APIs.
- Treating feedback events as CI proof, approval, or policy decisions.
- Redaction policy beyond size limits and file-based input.

## Acceptance Criteria

- The CLI writes a schema-shaped feedback event from `--message-file` and
  refuses unsafe actor/target tokens.
- The CLI avoids direct message content on argv; message text comes from a
  file.
- `doctor --profile github-actions-git-hooks-v1` includes `feedback_events`
  with `not_assessed` install and proof states for valid observed feedback.
- Existing `agent_prompt` remains `agent_reported_not_proof`.
- Focused Go tests cover event writing, unsafe token rejection, and observer
  surface detection.

## Evidence Boundary

Feedback events can show that a corrective message was deliberately retained in
the repo evidence workspace. They do not prove the original chat transcript is
complete, nor that the agent complied with the feedback. External proof still
requires CI/PR or another witness surface.

## Socratic Review Notes

- What could be overclaimed? Feedback presence could be mistaken for
  compliance. Mitigation: event proof state is always `not_assessed`, and doctor
  next action requires CI/PR binding before external proof.
- What could leak? Feedback may contain sensitive text. Mitigation: input is an
  explicit file, size-limited to 16 KiB; no direct message flag is supported.
- What remains not assessed? Automatic connector capture and full transcript
  completeness are not assessed by this block.

## Implementation Review Notes

- Code/correctness review reported event ID collision risk and missing explicit
  actor/target validation tests. Disposition: fixed with random event ID nonce
  and focused tests.
- Tracing/evidence review reported a critical overclaim: any `.json` file under
  `.sdp-trace/feedback` was treated as observed. Disposition: fixed by parsing
  and validating feedback event shape and message digest before observation;
  malformed files now produce `cannot_verify`.
- Focused re-review found no remaining critical or major feedback-boundary
  issues. Remaining minor note: first invalid JSON file in the feedback
  directory makes the surface `cannot_verify`, which is intentionally stricter
  than ignoring corruption.
