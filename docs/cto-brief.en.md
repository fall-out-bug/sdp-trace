# sdp-trace CTO Brief

AI coding is easy to start and hard to govern.

The risk is not that an agent writes code. The risk is that nobody can later explain the scope, provenance, evidence, and quality decision behind the change.

`sdp-trace` is a portable trust layer for AI-assisted delivery. It lets teams keep their existing harness while making changes traceable, evidence-backed, and gateable.

## What It Controls

- Scope: what change was intended.
- Provenance: who or what performed the work.
- Evidence: tests, CI, reviews, commands, diffs, and referenced files.
- Gate verdict: `pass`, `warn`, `fail`, or `not_assessed`.
- Decision record: why the change was accepted, blocked, or overridden.

## What It Does Not Promise

- It does not replace code review.
- It does not guarantee compliance.
- It does not prove code is bug-free.
- It does not require replacing Claude Code, Codex, OpenCode, Cursor, or an internal harness.

## First Adoption Step

Start with one repo and one real AI-assisted change. Produce an evidence bundle and a gate verdict. If the verdict helps reviewers make a better decision, expand to PR gates through `sdp-gate`.
