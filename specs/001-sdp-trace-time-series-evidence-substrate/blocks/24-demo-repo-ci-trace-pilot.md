# Block 24: Demo Repository CI And Trace Pilot

Status: draft SpecKit delta for Socratic review. Implementation is blocked
until this spec is reviewed, valid critical/major findings are resolved, and
the CTO explicitly approves the reviewed direction.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/23-mvp-closure-drift-and-readiness.md`
- `docs/research/block-23-mvp-closure-package.md`
- `docs/customer-questions.en.md`
- `docs/customer-questions.ru.md`
- `docs/agent-entrypoint.md`
- `docs/reviewer-entrypoint.md`

## Goal

Prove `sdp-trace` on a small demo repository through a real CI-backed workflow
with committed trace outputs, without relying on the retired Block 06 toy
runner/validator scripts.

The block must answer the practical customer question Block 23 left open:

> Can a team attach `sdp-trace` to a repository, run CI, and inspect trace,
> evidence, report, gate, and witness artifacts that honestly distinguish
> observed facts from `not_assessed` or `cannot_verify` gaps?

## Scope

The demo repository may be separate from this repository. `sdp-trace` remains
the portable product substrate; demo-specific application code, CI workflow
files, raw build logs, and raw model output belong in the demo repository unless
they are sanitized examples intentionally copied here.

Minimum pilot surface:

- one demo repository with a small app or service;
- CI workflow that runs at least one build/test command;
- `sdp-trace run` or `wrap` around the selected command;
- committed or externally referenced `.sdp-trace-runs/` artifacts;
- `verify`, `explain`, `report`, `gate`, and `witness` outputs;
- at least one trace path showing an observed successful command;
- at least one negative or incomplete path showing `not_assessed`,
  `cannot_verify`, missing telemetry, or witness gaps;
- redaction boundary for logs, command output, paths, tokens, and personal data;
- a customer-readable pilot report mapping artifacts to the Block 23 customer
  questions.

## Non-Goals

- Do not resurrect Block 06 `scripts/*` or `npm` validation as current product
  proof.
- Do not claim external production trust.
- Do not claim production readiness, policy enforcement, customer deployment
  readiness, or enterprise CI support.
- Do not put raw customer data, raw model responses, secrets, CI tokens, private
  paths, or long raw logs into this repository.
- Do not make `sdp-trace` depend on the demo repository, Beads, Operator Mode,
  a specific coding agent, or a specific CI provider.
- Do not turn demo gate output into a native `sdp-trace` policy decision.

## Acceptance Criteria

1. The demo repo run can be replayed or inspected from documented commands and
   artifact refs.
2. CI produces current evidence for at least one `sdp-trace`-wrapped command.
3. `sdp-trace verify` and `explain` are run against the captured run directory
   and their outputs are linked from the pilot report.
4. `sdp-trace report` and `gate` are run against the demo run set and are linked
   without treating gate output as policy enforcement.
5. `sdp-trace witness` is run in CI or records why CI witness evidence remains
   `not_assessed` or `cannot_verify`.
6. The pilot report maps the artifacts to the nine Block 23 customer questions.
7. The report distinguishes:
   - observed command execution;
   - missing telemetry;
   - local-only evidence;
   - CI-witnessed evidence;
   - external production trust that remains `not_assessed`.
8. Redaction checks prove committed artifacts do not contain CI tokens, OIDC
   tokens, provider credentials, private key material, authenticated provider
   URLs, raw model payloads, raw logs, private filesystem paths, or unsafe
   personal identifiers.
9. Any artifact copied back into `sdp-trace` is sanitized, source-attributed,
   and marked as demo evidence rather than source-bound product proof.
10. Review runs across code/correctness, trace/evidence, and
    requirements-vs-implementation planes before implementation closure and
    again at PR level.

## Open Design Questions For Socratic Review

1. Should the demo repository live under the same GitHub owner, or should it be
   intentionally external to prove portability?
2. Which CI provider is the first target: GitHub Actions because it is simplest,
   or a non-GitHub CI because Block 22 already broadened witness semantics?
3. What minimal demo app is credible without becoming a product distraction?
4. Which artifacts should stay only in the demo repo, and which sanitized
   evidence summaries should be copied into `sdp-trace`?
5. What is the smallest negative scenario that proves missing telemetry or
   witness gaps without fabricating a failure?

## Verification Plan

Minimum local checks in `sdp-trace`:

```bash
go test ./...
go vet ./...
staticcheck ./...
golangci-lint run ./...
jq empty schema/*.json
git diff --check HEAD
go run ./cmd/sdp-trace --help
```

Minimum demo checks, to be finalized after Socratic review:

```bash
go run ./cmd/sdp-trace run --task <demo-task> -- <demo-command>
go run ./cmd/sdp-trace verify <demo-run-dir>
go run ./cmd/sdp-trace explain <demo-run-dir>
go run ./cmd/sdp-trace report --out <demo-report-dir> <demo-runs-root>
go run ./cmd/sdp-trace gate --out <demo-gate-result.json> <demo-runs-root>
go run ./cmd/sdp-trace witness --kind <ci-kind> --out <demo-witness.json> <demo-runs-root>
```

Every command must record its actual environment and residual
`not_assessed`/`cannot_verify` states.
