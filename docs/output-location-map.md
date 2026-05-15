# Output Location Map

This table maps each command family to its default output location, format,
and trust boundary. For the full command contract, see
`docs/agent-entrypoint.md`.

## Run Artifacts

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `wrap` | `--output-dir .sdp-trace-runs/<name>/` | JSON + metadata | Record one command as a trace run | Local observation only |
| `run` | `--output-dir .sdp-trace-runs/<task-ref>/` | JSON + metadata | Task-linked trace run | Local observation; missing contract evidence visible |
| `observe setup` | `--out <run-dir>/` | JSON | Setup metadata before harness run | Session-profile bounded |
| `observe collect` | normalizes into `<run-dir>/` | JSON | Harness output after run | `cannot_verify` if declared output missing |
| `observe session` | `--out <run-dir>/` | JSON | Convenience wrapper for setup + collect | Same as setup + collect |
| `harness observe` | `--out <run-dir>/` | JSON | Import local harness export | Reads explicit files only; unsafe content fails before write |
| `harness validate` | `--out <validation.json>` | JSON | Validate observed harness events | Emits evidence facts; missing required event families are not passes |

## Environment And Setup

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `doctor` | stdout | Markdown + JSON | Inspect local environment and prerequisites | Structural readiness; offline or missing prerequisites can produce `cannot_verify` |
| `doctor --profile <profile>` | `--out <file>` | JSON | Inspect repo observer installation and proof state | Local hooks/config are `local_structural`; CI artifact proof remains `not_assessed` until observed |
| `install repo-observer` | `--out <file>` | JSON | Install portable repo observer files | Dry-run by default; with `--write`, writes documented allowlist only |

## Reports And Summaries

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `report` | `--out .sdp-trace-report/` | JSON + markdown | Package observed data and gaps | Report presence is not proof of completeness |
| `gate` | `--out .sdp-trace-report/gate-result.json` | JSON | Advisory gate facts | Not a native merge/release/risk decision |
| `explain` | stdout | Markdown | Human-readable run explanation | Does not upgrade trust scope |
| `harness summarize` | stdout | Markdown | Human summary of harness validation | Non-authoritative |
| `assess explain` | stdout | Markdown | Explain assessment result | Unsupported schema may give `cannot_verify` |
| `query-pack explain` | stdout | Markdown | Explain forensic query-pack | No new evidence created |
| `envelope summarize` | `--out summary.json` | JSON | Summarize delivery trace envelope | Read-only over refs |
| `interaction summarize` | `--out summary.json` | JSON | Summarize interaction events | Friction counts are facts, not scores |

## Query And Assessment

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `query` | stdout | JSON | Missing evidence or capture depth | Missing rows are not passes |
| `query-pack` | `--out <file>` | JSON | Forensic query package | Limited by retained/redacted evidence |
| `assess --profile <profile>` | `--out <file>` | JSON | Profile-specific assessment | Facts only; policy owns block/allow |
| `assess preview` | stdout | JSON | Preview required inputs | Read-only; does not evaluate authority |

## Witness And Release

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `witness --kind <kind>` | `--out <file>` | JSON | CI or customer witness artifact | CI witness is not production trust by itself |
| `release-proof` | `--out <file>` | JSON | Source-bound local release proof | Narrower than external production trust |

## Cross-Repo And Telemetry

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `export cross-repo-posture` | `--out <file>` | JSON | Cross-repo evidence posture | Degradation decisions remain outside |
| `export telemetry` | `--out <file\|->` | Prometheus text | Telemetry export | Dashboards/alerts remain downstream |

## PR Review

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `pr-review packet` | `--out <dir>/` | JSON + files | Build frozen PR packet | Packet digest binds to inputs |
| `pr-review run` | `--out <dir>/` | JSON | Run review planes | Raw output digest recorded as `raw_output_ref` |
| `pr-review synthesize` | `--out <file>` | JSON | Synthesize runs | Aggregation only |
| `pr-review validate` | `--out <file>` | JSON | Validate against ledger | Completeness check |
| `pr-review summarize` | `--out <file>` | JSON | Summarize validation | Not merge approval |
| `pr-review check` | `--out <dir>/` | JSON + files | End-to-end PR review | Review-record completeness only |

## Packet Commands

| Command | Default output | Format | Purpose | Trust boundary |
| --- | --- | --- | --- | --- |
| `packet build-pr` | `--out <dir>/` | JSON + files | Build live PR packet | `PC-VERIFICATION` must bind to workflow evidence |
| `packet build-github` | `--out <file>` | JSON | Build from fixture | Backfill/fixture authority only |
| `packet validate` | stdout | JSON | Validate bundle | Structural validation |
| `packet check-demo` | stdout | JSON | Demo-check bundle | Limited to first-packet minimum bar |
| `packet render` | `--out <file>` | Markdown | Render bundle | Row states and residual gaps are not approval |
