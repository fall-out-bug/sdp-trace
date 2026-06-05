# SpecKit Delta: Follow-Up Readiness Hardening

**Parent spec**: `specs/004-mvp-readiness-hardening/spec.md`
**Status**: Draft - Kimi PI review complete, implementation approval pending
**Created**: 2026-05-17
**Input**: Cross-pattern repository review covering formatting drift, releaseproof source-commit validation, stale MI claims, duplication, lint cleanup, command-surface drift, and Go locality/readability.

This is a SpecKit planning delta. Authoritative implementation claims must use
the repository's `sdp-trace-claim` syntax in the relevant implementation,
evidence, or review artifacts. Acceptance criteria in this file define planned
closure contracts; they are not proof that the closure has happened.

## Product Boundary

This delta closes follow-up readiness findings without changing the product boundary. `sdp-trace` remains a portable trust substrate and must not gain dependencies on SDP Operator Mode, Beads, agentloop, `sdp_lab`, GitHub, Pi, Kimi, Codex, OpenCode, or any other harness runtime.

Allowed:

- Markdown specs, task matrices, review ledgers, and claim corrections;
- JSON Schema updates when tied to an existing portable contract;
- small Go validation, rendering, or CLI helpers with focused tests;
- CI/doc checks that replay local evidence.

Not allowed:

- Node.js, npm, JavaScript, TypeScript, or `.mjs` product tooling;
- checked-in proof JSON treated as authority without live replay or external signature;
- docs or task checkboxes claiming closure without matching verifier state;
- broad file reshuffling that only optimizes metrics without improving review navigation.
- TODO/FIXME markers in new Go code.

## Core Claims

This delta may claim only:

> Follow-up readiness findings have a traceable closure plan, scoped implementation slices, and explicit pass/fail/not_assessed criteria.

It must not claim:

- external CI is green until live GitHub checks are queried for the exact head SHA;
- absolute MI `>70` passes unless the live absolute gates pass;
- releaseproof source binding is safe before invalid source refs are rejected by tests and implementation;
- reviewer output is approval before findings are verified against full files.

## Finding Closure Matrix

| Finding | Scope | Expected closure | Verification state before work |
|---|---|---|---|
| Format/import drift | `cmd`, `internal`, `tools` Go files named in review | `gofmt -s` and imports are stable; no formatting-only drift remains | `not_assessed` until replayed |
| Releaseproof accepts unsafe source refs | `internal/releaseproof`, `schema/contract-manifest.schema.json` | source commits accept only immutable commit object refs; branch names, symbolic refs, revspec suffixes, pathspecs, flags, empty/whitespace values are rejected | `not_assessed` until tests and gosec replay pass |
| Stale MI overclaim | `completion-audit.md`, `docs/spec-drift-register.md`, CI docs | docs either remove absolute MI pass language or are backed by live absolute MI pass evidence | `not_assessed` until doccheck and MI commands replay |
| High-signal duplication | observe setup/collect CLI and repeated tests | shared helpers remove duplicated behavior without changing command boundaries | `not_assessed` until `dupl` and focused tests replay |
| Gocritic/unparam/prealloc cleanup | named production cleanup candidates | simple lint findings are fixed or narrowly justified with evidence | `not_assessed` until lint replay |
| Command surface source-of-truth drift | command registry, help, docs, tests | command metadata/help/docs have one checked source or a machine-checkable drift gate | `not_assessed`; may remain a separate approved slice |
| Go locality/readability degradation | numbered one-function files, empty package-only files, package-stutter filenames, boilerplate comments | new layout rules prevent further degradation; high-churn packages are grouped gradually by cohesive behavior | `not_assessed`; broad reshuffle is out of immediate closure unless separately sliced |

## Implementation Slices

### Slice 1 - Delta And Review Matrix

Create this delta, a task plan, PI review packet, and explicit state vocabulary for all findings.

Acceptance:

- each finding maps to files, commands, expected closure, and verifier state;
- external CI is explicitly `not_assessed`;
- no task is marked done because a prose plan exists.

### Slice 2 - Format And Imports

Fix formatting/import drift only.

Initial target files:

- `tools/schemadoc/const.go`
- `tools/doccheck/quickstart_test.go`
- `cmd/sdp-trace/crap_hotspot_test.go`
- `internal/demo/demo_events_payload_int.go`
- `internal/demo/demo_events_payload_string.go`
- `internal/demo/demo_gate_required_evidence.go`

Required checks:

- `gofmt -l $(find cmd internal tools -name '*.go')`
- `golangci-lint run --disable-all --enable=gofmt --enable=goimports ./...`
- `go test -count=1 ./...`

### Slice 3 - Releaseproof Source Commit Hardening

Validate `source_commit` before any `git cat-file`, `git show`, or equivalent releaseproof source inspection.

Target files:

- `internal/releaseproof/source_state.go`
- `internal/releaseproof/artifacts.go`
- `internal/releaseproof/manifest_state.go`
- `internal/releaseproof/releaseproof_test.go`
- `schema/contract-manifest.schema.json`

Acceptance:

- accept only immutable commit object identifiers supported by the contract;
- reject empty, whitespace, `HEAD`, branch names such as `main`, `abc^{tree}`, `abc:path`, `--flag`, path-like values, and any ref that can resolve to a non-commit object;
- each remaining gosec row is either fixed or narrowly justified with `// #nosec G...` and local evidence.
- add red tests for the rejection table before changing `source_state.go` or other implementation files.
- validate the source commit in pure Go before `git cat-file`, `git show`, or any other `exec.Command` invocation.
- keep schema and Go validators aligned. If schema changes are needed, either add a reusable `$defs` entry such as `gitCommitSHA` or constrain `source_commit` inline and document why.

Negative tests must cover at least:

- empty string, whitespace-only, and tab/newline-only values;
- `HEAD`, `main`, `origin/main`, and `refs/heads/main`;
- `abc^{tree}`, `abc:path`, `abc~1`, and `abc..def`;
- `--flag`, `-e`, and `--help`;
- path-like values such as `../etc/passwd` and `foo/bar`;
- 39-character hex, 41-character hex, and 64-character hex values;
- uppercase hex when the validator is intentionally lowercase-only;
- non-hex characters in a 40-character string.

Required checks:

- `go test -count=1 ./internal/releaseproof ./cmd/sdp-trace`
- `golangci-lint run --disable-all --enable=gosec ./...`
- `govulncheck ./...`

### Slice 4 - High-Signal Duplication Cleanup

Close only duplication that improves review navigation or reduces command/test drift.

Production targets:

- shared observe CLI flow between `cmd/sdp-trace/harness_observe_cli.go` and `cmd/sdp-trace/observe_setup_cli.go`;
- keep command boundary behavior unchanged.

Test targets:

- shared helpers for repeated observe-collect rejection tests;
- one hygiene table runner where repeated table logic is genuinely identical.

Required checks:

- `golangci-lint run --disable-all --enable=dupl ./...`
- focused package tests for changed packages.

### Slice 5 - Gocritic, Unparam, And Prealloc Cleanup

Make production-safe lint cleanups first; skip prealloc changes that reduce readability.

Initial target files:

- `internal/forensic/forensic_tables.go`
- `cmd/sdp-trace/main_546_commandsurfacedrift.go`
- `cmd/sdp-trace/wrap_preview_args.go`
- `internal/posture/posture_time.go`

Required checks:

- `golangci-lint run --disable-all --enable=gocritic --enable=unparam --enable=prealloc ./...`
- `go test -count=1 ./...`

### Slice 6 - Maintainability And Docs Overclaim Closure

Make the MI claim honest. Either absolute MI `>70` actually passes, or docs explicitly describe the live ratchet/baseline state without pass language.

Default closure path: correct stale docs to match live verifier state. Actual
absolute MI improvement is secondary and may be attempted only when it does not
delay doc honesty. If live absolute MI replay fails, `completion-audit.md` and
related docs must be corrected in the same scoped commit as the verification
evidence. Do not split doc closure away from the verifier state it describes.

Initial files to audit:

- `specs/004-mvp-readiness-hardening/completion-audit.md`
- `docs/spec-drift-register.md`
- `.github/workflows/ci.yml`

If pursuing absolute MI improvement, first targets are:

- `tools/doccheck/quickstart.go`
- `tools/doccheck/registry.go`
- `cmd/sdp-trace/main_546_commandsurfacedrift.go`
- command-surface registry files.

Required checks:

- `go run ./tools/qualitycheck -fail-only -mi-under 70.1 cmd internal tools`
- `go run ./tools/qualitycheck -fail-only -function-mi-under 70.1 cmd internal tools`
- `go run ./tools/doccheck`

## Deferred Or Separate Specs

The following findings are real but should not be mixed into the six closure slices unless separately approved:

- command surface single source of truth across registry, help, docs, and tests;
- broad Go locality cleanup of numbered one-function files;
- package layout decomposition for `cmd/sdp-trace`, `internal/harnessobs`, `internal/packet`, and `internal/prreview`;
- comment-noise cleanup across trust-boundary boilerplate.

These should become separate SpecKit deltas with their own PI review because they can touch hundreds of files and can easily blur evidence closure with readability preference.

## Final Gate

After each implementation slice:

- `go test -count=1 ./...`
- `go vet ./...`
- `golangci-lint run ./...`
- `git diff --check`
- scoped commit with the slice name and verifier state.

Before closing the whole package:

- `go test -count=1 ./... -coverprofile=coverage.out`
- `go tool cover -func=coverage.out > coverage-func.txt`
- `go run ./tools/qualitycheck -gocyclo cmd internal tools > gocyclo.txt`
- `go run ./tools/crapcheck -cover-func coverage-func.txt -gocyclo gocyclo.txt -threshold 5 -strict-less`
- `jq empty schema/*.json`
- `go run ./tools/hygienecheck`
- `go run ./tools/schemadoc`
- `go run ./tools/schemadoc -verify-readme`
- `govulncheck ./...`
- adversarial PI review planes: security/trust, code/correctness, docs/evidence drift.

External GitHub CI remains `not_assessed` until queried live for the exact final head SHA.

## PI Review Prompt

Review this SpecKit delta for whether it gives future agents enough information to close the reported readiness findings without overclaiming verification or mixing unrelated changes. Focus on:

- whether each slice is small enough for scoped commits;
- whether releaseproof source-commit hardening has adequate negative tests;
- whether MI and CI claims are evidence-safe;
- whether broad Go locality cleanup is correctly separated from immediate closure;
- whether any step depends on a specific harness such as Pi, Kimi, Codex, OpenCode, GitHub, Beads, or SDP Operator Mode.
