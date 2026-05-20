# Production Adoption Readiness

`sdp-trace` is a portable evidence substrate for AI-assisted delivery.
This document records what is known, what is pilot-capable, and what remains
`not_assessed` for production adoption beyond a controlled pilot.

## Trust Scope

`sdp-trace` records evidence and gaps. It does not approve merge, release,
risk acceptance, or production trust. Every row below separates what is
known from what remains `not_assessed` or `cannot_verify`.

## Readiness Matrix

| Area | Controlled Pilot | Local Evidence | CI-Witnessed | Source-Bound Release | External Production Trust |
|------|:----------------:|:--------------:|:------------:|:--------------------:|:-------------------------:|
| Repository existence | yes | yes | yes | yes | `not_assessed` |
| Public documentation | yes | yes | yes | yes | `not_assessed` |
| Command surface stability | partial | partial | partial | partial | `not_assessed` |
| Local security baseline | yes | yes | `not_assessed` | `not_assessed` | `not_assessed` |
| CI security baseline | `not_assessed` | `not_assessed` | `not_assessed` | `not_assessed` | `not_assessed` |
| External security audit | `not_assessed` | `not_assessed` | `not_assessed` | `not_assessed` | `not_assessed` |
| Customer adoption evidence | `not_assessed` | `not_assessed` | `not_assessed` | `not_assessed` | `not_assessed` |
| Signed release process | `not_assessed` | `not_assessed` | `not_assessed` | `not_assessed` | `not_assessed` |

Legend:
- **yes**: Evidence exists and is verifiable from source.
- **partial**: Some evidence exists; gaps are recorded.
- **`not_assessed`**: No evidence or external verification available.

## Command Family Readiness

State values follow the command surface schema: `complete`, `partial`, `not_assessed`.

| Command Family | State | Trust Note |
|----------------|-------|------------|
| `assess` | `partial` | Emits verifier facts; missing or stale evidence can produce `cannot_verify`. |
| `checkpoint` | `partial` | Subcommands `create` and `verify` are present; completeness not verified. |
| `envelope` | `complete` | Read-only over refs; reports linked and `not_assessed` areas. |
| `explain` | `complete` | Derived from run artifacts; does not upgrade trust scope. |
| `gate` | `partial` | Emits verifier facts and states; not a native merge/release/risk decision. |
| `harness` | `partial` | Harness event import and validation; limited to supported harness kinds. |
| `interaction` | `partial` | Interaction recording and summarization; relay and import-transcript paths are pilot-only. |
| `observe` | `partial` | First-run harness observation; session profile shape may change. |
| `export` | `partial` | Export cross-repo posture or telemetry profiles for external consumption. |
| `override` | `not_assessed` | Override record creation; authority scope is advisory. |
| `packet` | `partial` | Packet generation and validation; schema versioned but not externally signed. |
| `pr-review` | `partial` | Build, run, synthesize, validate, and summarize automated PR review evidence. |
| `query` | `complete` | Highlights gaps; missing rows are not passes. |
| `query-pack` | `partial` | Produces investigation package; digest-only or redacted data limits reconstruction. |
| `release-proof` | `complete` | `source_bound_local_release` only; dirty/stale source or manifest mismatch fails. |
| `report` | `complete` | Packages observed data and gaps; report presence is not proof of completeness. |
| `verify` | `complete` | Supports local structural assertions only. |
| `witness` | `partial` | CI-bound evidence is not external production trust by itself. |

**Summary**: 6 families are `complete`, 11 are `partial`, 1 is `not_assessed` (18 families).

## Utility Commands

The following commands are stable utilities that support the main families
but are not primary evidence-producing workflows.

| Command | State | Description |
|---------|-------|-------------|
| `command-surface` | `complete` | Emit machine-readable command surface JSON. |
| `doctor` | `complete` | Inspect local environment and contract prerequisites. |
| `dry-run` | `complete` | Show what would run without writing run artifacts. |
| `install` | `complete` | Install portable repo observer files for local git hooks and GitHub Actions artifact upload. |
| `preview` | `complete` | Preview command/contract implications before execution. |
| `run` | `complete` | Run a task-referenced command with an optional contract. |
| `validate-fixtures` | `complete` | Validate checked trace-run fixture directories. |
| `version` | `complete` | Print version. |
| `wrap` | `complete` | Observe one existing command as a trace run. |

**Total command surface**: 18 families + 9 utility commands = 27 commands.

## Spec Completion

Most active specs are `draft` or `in_progress` per [`docs/roadmap.md`](roadmap.md).
No spec has reached `complete`. This is expected for an early-stage project
and does not block controlled pilot use.

## Security Baseline

| Tool | Local Result | Disposition |
|------|-------------|-------------|
| `go vet ./...` | pass | green |
| `govulncheck ./...` | 0 vulnerabilities | green |
| `gosec ./...` | 132 findings | advisory — classified in [`docs/security-baseline.md`](security-baseline.md) |
| `gitleaks detect` (with `.gitleaks.toml`) | 0 findings | triaged in [`docs/security-baseline.md`](security-baseline.md) |

Local ignored clutter (`.worktrees/`, `.codex-subagents/`, `.sdp-trace-runs/`,
`.sdp-trace-report/`) is not repository proof. Scan it with hygiene tooling
locally, but do not treat findings in ignored paths as product evidence.

## `not_assessed` Areas (Explicit)

1. **Live CI status**: GitHub Actions workflows exist but have not been
   verified as green at the current `main` head for this block.
2. **Customer adoption**: No external customer has reported production use.
3. **External security audit**: No third-party security review has been
   commissioned or completed.
4. **Signed release process**: Release artifacts are built locally; no
   externally signed or witnessed release pipeline exists.
5. **Production trust decision**: This repository structures evidence and gaps;
   it does not approve changes, releases, or risk acceptance.

## What This Document Does Not Claim

- No production trust decision has been made.
- No merge or release approval is implied.
- No claim that `gitleaks` or `gosec` exhaustively prove absence of secrets or vulnerabilities.
- No claim that local checks alone support production readiness.
- No external customer adoption evidence.

## Related Documents

- [Adoption Guide](adoption-guide.en.md): implementation model and evidence boundaries.
- [Repository Rollout Playbook](repository-rollout-playbook.en.md): team setup and evidence contracts.
- [Security Baseline](security-baseline.md): scanner triage ledger.
- [Overclaim Checklist](overclaim-checklist.md): forbidden claims and trust-scope rules.

## Verification Commands

Run these locally to reproduce the current baseline:

```bash
# Green gates
go vet ./...
go test -count=1 ./...
govulncheck ./...
go run ./tools/doccheck
go run ./tools/hygienecheck
git diff --check

# Advisory scans
# gosec ./...
# gitleaks detect --source . --config .gitleaks.toml
```

## Maintenance

Update this matrix when:

- a new completed spec is recorded in `docs/roadmap.md`;
- scanner findings are classified or resolved;
- a live CI run completes and is recorded;
- external customer adoption evidence becomes available;
- a GitHub security policy is enabled and configured.

## Change Log

- 2026-05-20: Initial matrix for spec 016. All command families and security
  tools recorded. Explicit `not_assessed` areas named.
