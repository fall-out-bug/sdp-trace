# Contributor Quick Start

Use this page to verify that your local environment can build, run, and check
`sdp-trace` before deeper work.

## Prerequisites

- [Install](install.md): choose the release binary or the source-checkout path.
- [Core Concepts](concepts.md): read the vocabulary boundary before interpreting
  any verifier output.

## Canonical Smoke Path

Run these commands in order from the repository root. They are the minimum
local verification steps; each command must complete before the next one is
meaningful.

```text
go test -count=1 ./...
go run ./cmd/sdp-trace --help
go run ./cmd/sdp-trace doctor
go run ./cmd/sdp-trace wrap --name smoke --output-dir .sdp-trace-runs/smoke -- /bin/echo ok
go run ./cmd/sdp-trace verify .sdp-trace-runs/smoke
go run ./cmd/sdp-trace explain .sdp-trace-runs/smoke
go run ./cmd/sdp-trace report --out .sdp-trace-report .sdp-trace-runs/smoke
go run ./cmd/sdp-trace query --query missing-evidence .sdp-trace-runs/smoke
```

On Windows, replace `/bin/echo ok` with a command available in your shell, for
example `cmd /c echo ok` in Command Prompt.

## Expected Results

| Step | Expected state | What it proves |
|------|----------------|----------------|
| `go test` | pass | Go toolchain and local build are healthy. |
| `--help` | exit `0` (help displayed) | CLI compiles and prints current command surface. |
| `doctor` | `offline_dev` (local development); `pass` only in CI or with `--profile` | Local wrapper and output directories are ready. CI prerequisites may show `cannot_verify`; that is expected for local development. |
| `wrap` | exit `0` (run recorded) | The tool can record a trace run. |
| `verify` | `observed` | The run directory is structurally valid. |
| `explain` | human-readable summary | The run can be rendered for review. |
| `report` | report files written | The run can be packaged for review. |
| `query --query missing-evidence` | missing-evidence table | Core evidence gaps can be inspected without an assessment profile. |

## Failure Routing

When a step fails or returns an unexpected state, use this table to choose the
next diagnostic command.

| Symptom | Likely cause | Next step |
|---------|--------------|-----------|
| `go test` fails | Go/toolchain setup or checkout corruption | Run `go version`, check `GOPATH`, inspect test output for package-level errors. |
| `--help` fails or prints no commands | Build failure or dirty checkout | Run `go build ./cmd/sdp-trace` and inspect compilation errors. |
| `doctor` fails with local errors | Permissions or missing directories | Inspect the JSON output for the failing `control_point`; check that `.sdp-trace-runs` is writable. |
| `wrap` fails | CLI usage error | Verify the double dash `--` is present and the command after it exists; check that `--output-dir` is writable. |
| `verify` fails | Corrupted or missing run directory | Check that the run directory exists and contains `run.json`. |
| `explain` fails | Run directory unreadable | Re-run `verify` on the same directory to get structured errors. |
| `query --query missing-evidence` fails | Missing verifier artifact or replay error | Re-run `verify` and inspect the run directory path. |

## Before You Write Tasks Or Claims

Read [Claim Authoring](claim-authoring.md) before writing task checkboxes,
proof prose, or verdict language in this repository. Untagged prose is not
machine-authoritative.

## Full Command Surface

For the authoritative command, state, trust-scope, and exit-code contract, see
[Agent Entrypoint](agent-entrypoint.md). This quick start shows only the
smallest core verification path; it does not duplicate the full command table
or optional extension surfaces.

## Trust Scope Note

Smoke path results are `local_observed` evidence only. They do not become
CI evidence, production trust, or release approval without a subsequent
witness or assessment step scoped to the correct profile.
