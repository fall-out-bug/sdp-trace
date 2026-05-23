# Project Context: sdp-trace

Portable trust substrate for AI-assisted delivery.

## Start Here

1. **Load `AGENTS.md`** before any work. The Skills Router section is authoritative.
2. **Route through `sdp-trace-router`** skill before any generic global skill.
3. **Do NOT improvise** process steps outside loaded skills.

## Quality Gates (run before every claim)

```bash
go test -count=1 ./...
go vet ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
go run ./tools/schemadoc
# Coverage-backed CRAP
go test -count=1 ./... -coverprofile=coverage.out
go tool cover -func=coverage.out > coverage-func.txt
go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt
go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less
# MI baseline
go run ./tools/qualitycheck -mi-baseline tools/qualitycheck/file-mi-baseline.json cmd internal tools
go run ./tools/qualitycheck -function-mi-baseline tools/qualitycheck/function-mi-baseline.json cmd internal tools
# Whitespace
git diff --check
# JSON
for f in schema/*.json; do jq empty "$f" || exit 1; done
```

## Trust Rules (non-negotiable)

- Machine proof wins over prose, checkboxes, and ledgers.
- Checked-in proof JSON is not authority unless live-verified.
- No deferred trust closure: missing external evidence keeps the block open.
- Source-bound proof requires a clean immutable source commit.
- Every claim is `pass`, `fail`, `cannot_verify`, or `not_assessed`.

## Block Intake

When the user says "берем блок в работу":
1. Invoke `sdp-trace-trust-workflow` skill.
2. Follow `intake_protocol` step by step.
3. Do NOT ask exploratory questions before intake completes.

## Review

On review requests, invoke `pi-review` skill.
Prefer non-OpenAI, non-Anthropic, non-Google models.
Iterative adversarial review: fix every finding, repeat until **LGTM** (zero findings).

## Engineering Stack

- Target: Go.
- No Node.js, npm, JavaScript, TypeScript, or `.mjs` in active product path (`cmd/`, `internal/`, `tools/`).
- Bash: thin launcher only, with explicit reason.
- New code: small, readable, testable, TODO/FIXME-free.

## Project Skills (`.agents/skills/`)

| Skill | When to use |
|---|---|
| `sdp-trace-router` | Entrypoint for all work |
| `sdp-trace-trust-workflow` | Block intake, implementation, PR discipline |
| `pi-review` | Adversarial review orchestration |
| `sdp-trace-pi-handoff` | External worker delegation |
| `sdp-trace-quality-audit` | Quality gates, docs, DX/UX audit |
