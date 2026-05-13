# Agent Instructions

`sdp-trace` is the portable trust substrate for AI-assisted delivery.

## Purpose
Define traceability, provenance, evidence, gate verdicts, and decision records that work across coding agents and harnesses.

## Boundary
This repo must stay independent from `sdp_lab`, Beads, Operator Mode, agentloop, and any specific harness runtime.

Allowed: JSON schemas, Markdown docs, portable examples, and tiny Go validation/rendering tools when needed.

Not allowed: dependencies on SDP Operator Mode, Beads, agentloop, or hidden assumptions about Claude, Codex, OpenCode, GitHub, or any harness.

## Product Language
Use SpecKit-aligned terms first: spec, plan, task, evidence, gate, decision, trace, provenance. Avoid internal SDP terms unless a doc explicitly maps them.

## Quality Bar
Every claim about a gate or verdict must be evidence-backed or marked `not_assessed`. No opaque health scores.

## Engineering Stack
Target product code is Go.

No Node.js, npm, JavaScript, TypeScript, or `.mjs` tooling is allowed in the active product path.

Bash is allowed only as a thin command launcher when Go would add no product value; any active Bash needs an explicit reason.

New Go code must be small, readable, testable, covered by focused tests, and free of TODO/FIXME markers. Put measurable complexity gates in CI or docs, not only in prose.

## Decomposition Rule
Keep this root router under 100 lines. If a module needs more than 10 skills, the module is too large, under-decomposed, or overengineered.

## Trust Rules
The repository has already failed once by letting prose, task checkboxes, and checked-in JSON overclaim what current verification could not replay. Do not repeat that failure.

- Machine proof wins over prose, task checkboxes, reports, review ledgers, and mirrors.
- Dirty checkout output is local structural evidence only, not external trust.
- Checked-in proof JSON is not authority unless live-verified or externally signed.
- No deferred trust closure: missing external evidence keeps the block open.
- Prose is not authoritative. Use `sdp-trace-claim` tags for authoritative claims.
- Source-bound proof requires a clean immutable source commit; if a changed file is a manifest subject, commit it first, then regenerate release proof in a separate commit.
- Do not close task checkboxes, review ledgers, or docs after source-bound proof without another source-bound cycle if those files are manifest subjects.
- Keep mirrored self-trace data synchronized: `assessment-input.json` must mirror self-trace arrays, and hash references must match current files.

## Required Work Loop
Every non-trivial implementation chunk needs a SpecKit delta, Socratic review before approval handoff, trace coverage for verifier/trust changes, test-first behavior for behavior changes, drift checks, live verifier state (`pass`, `fail`, `cannot_verify`, or `not_assessed`), strict review, fresh verification, and a scoped commit.

If a chunk cannot be traced or verified yet, mark `not_assessed` or `cannot_verify` with a concrete reason and create a tracked follow-up before closing.

## Skills Router
Use local project skills for detailed workflows instead of expanding this file:

- `sdp-trace-trust-workflow`: block intake, SpecKit review, implementation slicing, PR/review/merge discipline.
- `pi-review`: adversarial reviewer orchestration, retries, model policy, and disposition rules.
- `sdp-trace-quality-audit`: repository polish, quality gates, completion audits, docs/UX/DX/security review, and final evidence mapping.

When the user says "берем блок в работу", use `sdp-trace-trust-workflow`.

For adversarial pi review in this repo, prefer non-OpenAI, non-Anthropic, and non-Google models unless the user explicitly permits otherwise.

## Claim Tags
Use `docs/claim-authoring.md` for authoritative claim syntax. Current Slice 1 validator intentionally accepts only narrow evidence forms; do not introduce arbitrary `proof:*`, `state:*`, or `none` evidence unless cross-reference verification has been implemented.

## Commands
Use current command contracts in `docs/agent-entrypoint.md` and `docs/reviewer-entrypoint.md`.

- Defaults: `go test ./...`, `jq empty schema/*.json`, `gofmt` for changed Go files, and `git diff --check`.
- For schema/contract changes, also check refs, changed examples, fixture shape, and Go struct/schema alignment.

Bash verification commands are not product architecture. Keep them only when they are thin launchers around Go commands or external tools.
