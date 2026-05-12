# Implementation Plan: Invisible Flight Recorder

## Phase 0 - Reviewed Spec Direction

1. Add this spec and task list only.
2. Run multi-plane subagent review:
   - requirements;
   - evidence/trust;
   - security/non-interference;
   - DX/UX;
   - implementation risk.
3. Record findings and disposition before product code changes.

## Phase 1 - Minimal Product Slice

Implement the smallest slice that removes manual packet/evidence curation from
the live demo path:

- prompt boundary classifier for recorder-duty contamination;
- recorder/evidence authority manifest extension;
- `packet build-pr` or equivalent CI-owned builder from GitHub context fixture;
- docs that demote manual `github-input.json` packet generation to backfill.

## Phase 2 - Verification And Review

- Add focused Go tests before or with behavior changes.
- Run full repository checks.
- Run post-implementation review planes.
- Fix P0 findings immediately.
- Record P1+ follow-ups if they do not block the invisible flow.

## Phase 3 - PR And Release Surface

- Open PR from this feature branch to `main`.
- Do not commit directly to `main`.
- Confirm CI.
- Update release binary docs/workflow if command surface changes.
- Capture DX/UX review notes for docs and command help.

