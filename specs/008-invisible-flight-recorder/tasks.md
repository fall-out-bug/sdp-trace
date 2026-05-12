# Tasks: Invisible Flight Recorder

## Spec Review

- [x] T001 Run requirements review.
- [x] T002 Run evidence/trust review.
- [x] T003 Run security/non-interference review.
- [x] T004 Run DX/UX review.
- [x] T005 Run implementation-risk review.
- [x] T006 Record review findings and disposition.

## Product Implementation

- [x] T010 Add prompt boundary classifier and tests.
- [x] T011 Extend evidence/packet authority metadata so recorder, CI,
  operator, and integration artifacts are distinguishable.
- [x] T012 Add CI-owned PR packet builder from GitHub context/API fixture data.
- [x] T013 Add fixture tests for successful PR packet generation.
- [x] T014 Add negative tests for missing check/artifact/route/prompt boundary
  evidence.
- [x] T015 Add CLI entrypoint/help for the invisible path.
- [x] T016 Add prompt-boundary tests for clean, contaminated, digest-only,
  missing, and malformed metadata.
- [x] T017 Add fixture test proving `packet build-pr` writes bundle JSON,
  rendered Markdown, and result summary.
- [x] T018 Add tests proving checked-in stale packet inputs are not used as
  live CI authority.
- [x] T019 Add secret-redaction tests for GitHub token/header-like inputs.

## Documentation And Release

- [x] T020 Update `docs/agent-entrypoint.md`.
- [x] T021 Update `docs/change-evidence-packet.md`.
- [x] T022 Update install/release docs if new commands affect binary users.
- [x] T023 Run DX/UX docs review.

## Verification And PR

- [x] T030 Run `go test ./... -count=1`.
- [x] T031 Run `jq empty schema/*.json` plus changed examples.
- [x] T032 Run `git diff --check`.
- [x] T033 Run release binary build and checksum verification.
- [x] T034 Run post-implementation multi-plane review.
- [x] T035 Open PR to `main` with review and verification evidence.
