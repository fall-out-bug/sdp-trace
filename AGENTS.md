# Agent Instructions

`sdp-trace` is the portable trust substrate for AI-assisted delivery.

## Purpose
Define traceability, provenance, evidence, gate verdicts, and decision records that work across coding agents and harnesses.

## Boundary
This repo must stay independent from `sdp_lab`, Beads, Operator Mode, and any specific harness runtime.

Allowed:
- JSON schemas
- Markdown docs
- portable examples
- tiny validation/rendering tools when needed

Not allowed:
- dependency on SDP Operator Mode
- dependency on Beads
- dependency on agentloop
- hidden assumptions about Claude, Codex, OpenCode, or GitHub

## Product Language
Use SpecKit-aligned terms first: spec, plan, task, evidence, gate, decision, trace, provenance.

Avoid internal SDP terms unless a doc explicitly maps them.

## Quality Bar
Every claim about a gate or verdict must be evidence-backed or marked `not_assessed`.

No opaque health scores.

## Engineering Stack
Target product code is Go.

No Node.js, npm, JavaScript, TypeScript, or `.mjs` tooling is allowed in
the active product path.

Bash is allowed only as a thin command launcher when Go would add no
product value. Any Bash kept in the active path needs an explicit reason.

New implementation must follow Clean Architecture, Clean Code, TDD, and
keep CRAP below 5 for changed Go code. Do not add TODO or FIXME markers.

## Decomposition Rule
If `AGENTS.md` exceeds 100 lines or any module needs more than 10 skills, the module is too large, under-decomposed, or overengineered.

## Trust Rules
The repository has already failed once by letting prose, task checkboxes, and checked-in JSON overclaim what current verification could not replay. Do not repeat that failure.

- Machine proof wins over prose, task checkboxes, reports, review ledgers, and mirrors.
- Dirty checkout output is local structural evidence only, not external trust.
- Checked-in proof JSON is not authority unless live-verified or externally signed.
- No deferred trust closure: missing external evidence keeps the block open.
- Prose is not authoritative. Use `sdp-trace-claim` tags for claims.
- Source-bound proof requires a clean immutable source commit; if a changed file is a manifest subject, commit it first, then regenerate release proof in a separate commit.
- Do not close task checkboxes, review ledgers, or docs after source-bound proof without another source-bound cycle if those files are manifest subjects.
- Keep mirrored self-trace data synchronized: `assessment-input.json` must mirror self-trace arrays, and hash references must match current files.

## Required Work Loop
Every non-trivial implementation chunk needs:

1. SpecKit delta.
2. Trace coverage when verifier behavior or trust claims change.
3. Test-first behavior when behavior changes.
4. Live verifier state: `pass`, `fail`, `cannot_verify`, or `not_assessed`.
5. Strict review with recorded disposition.
6. Fresh verification before completion claims.

If a chunk cannot be traced or verified yet, mark `not_assessed` or `cannot_verify` with a concrete reason and create a tracked follow-up before closing.

## Block Intake Protocol
When the user says "берем блок в работу", use `sdp-trace-trust-workflow`.
First land current approved work through PR/review, then continue new block work in a fresh worktree.
Prepare SpecKit deltas, run Socratic spec review, fix/record findings, then stop for explicit user approval of the reviewed spec direction before writing implementation code.
After approval, split independent tasks to fast subagents with minimal context, and keep provenance, evidence, trace, pi review, fixes, and final PR review in the loop until no critical or major findings remain.
Do not stop at implementation-only closure for block work unless the user explicitly asks to stop before PR or merge.

## pi Review Rules
For adversarial pi review in this repo, prefer non-OpenAI, non-Anthropic, and non-Google models unless the user explicitly permits otherwise.

- Use MiniMax-M2.7 and ZAI/GLM for multi-file strict review.
- Run separate code, tracing/evidence, and requirements-vs-implementation review planes for trust blocks; repeat them at PR level. Verify review findings against full files before accepting or rejecting them.
- Use Kimi K2P6 as a full bounded reviewer for requirements, fixture semantics, and focused re-review; prefer low/off reasoning for narrow prompts.
- Use MiniMax-M2.5 only for OpenCode demo development, not repo review. Replace hung or empty pi reviews; do not count them as evidence. Record absent GitHub checks as CI `not_assessed`, not green.

## Claim Tags
Use `docs/claim-authoring.md` for authoritative claim syntax.

Current Slice 1 validator intentionally accepts only narrow evidence forms. Do not introduce arbitrary `proof:*`, `state:*`, or `none` evidence unless cross-reference verification has been implemented.

## Commands
Block 10 active development commands must be Go-first:

- Go tests: `go test ./...`
- Schema parse checks: `jq empty schema/*.json`; for schema/contract changes, also check refs, fixture shape, and Go struct/schema alignment.
- Formatting: `gofmt` for changed Go files

Bash verification commands are not product architecture. Keep them only
when they are thin launchers around Go commands or external tools.
