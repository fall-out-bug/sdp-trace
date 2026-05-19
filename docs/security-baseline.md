# Security Baseline

Status: draft

This document records local security probes for the adoption-readiness work.
It is a triage ledger, not a production trust claim. Scanner output can show
findings that need review; it cannot prove absence of vulnerabilities or
secrets.

## Probe Metadata

| Field | Value |
| --- | --- |
| Date | 2026-05-20 |
| Scanner target | Local branch snapshot after final review fixes |
| Go version | `1.22` from `go.mod` |
| `gosec` version | `dev` build from local OSS tool cache |

The scanner evidence was produced locally on this branch. Re-run the commands
below before using this ledger in a later PR; raw scanner JSON is intentionally
not checked in.

## Tool Summary

| Tool | Result | Finding Count | State | Notes |
| --- | --- | ---: | --- | --- |
| `go vet ./...` | pass | 0 | `verified` | Go vet completed locally. |
| `govulncheck ./...` | pass | 0 | `verified` | No known vulnerabilities reported by the local database snapshot. |
| `gosec ./...` | exit 1 | 132 | `needs_triage` | Static findings require code-level review or reviewed suppressions. |
| `gitleaks detect --source . --no-git --redact=100`, default config | exit 1 | 14 | `triaged` | Includes local `.codex-subagents/` clutter and tracked fixture candidates. |
| tracked-source `gitleaks` snapshot, default config | exit 1 | 10 | `triaged` | Scanned an archive of tracked `HEAD`; findings are reviewed fixture/test candidates. |
| working-tree `gitleaks`, `.gitleaks.toml` | pass | 0 | `verified` | Reviewed allowlist suppresses only known fixture/local-clutter patterns. |
| tracked-source `gitleaks`, `.gitleaks.toml` | pass | 0 | `verified` | Reviewed allowlist suppresses only known fixture patterns. |

## `gosec` Family Summary

| Rule | Count | Scanner Severity | Triage State | Current Disposition |
| --- | ---: | --- | --- | --- |
| `G304` | 60 | medium | `needs_triage` | Variable file reads/opens. Must verify path provenance and traversal controls per call site. |
| `G301` | 28 | medium | `needs_triage` | Directory permissions `0755`. Must confirm generated evidence readability is intentional for each path. |
| `G306` | 24 | medium | `needs_triage` | File permissions `0644`. Must confirm public-readable evidence output is acceptable for each artifact class. |
| `G204` | 11 | medium | `needs_triage` | Subprocess launches. Must distinguish intended wrapper behavior from untrusted command execution. |
| `G703` | 5 | high | `needs_triage` | Path traversal via taint analysis. This rule ID is emitted by the local `gosec` `dev` build used for the baseline. |
| `G101` | 2 | high | `candidate_false_positive` | Semantic strings matched credential heuristics; verify before allowlisting. |
| `G302` | 1 | medium | `needs_triage` | `chmod 0644` on generated evidence file; requires artifact-sensitivity review. |
| `G115` | 1 | high | `candidate_false_positive` | `rune` to `byte` conversion appears guarded by `r <= 127`; verify before suppressing. |

One reviewed `G304` false positive in `internal/capturedepth` has a scoped
`#nosec` suppression because `runDir` is a caller-selected local evidence root
and query output does not echo provider refs. The remaining `gosec` findings
are not treated as closed by this document. The next code slice should either
fix, narrow, or add a reviewed suppression for each accepted false positive.

## Tracked `gitleaks` Findings

The tracked-source snapshot with the default config reported 10 findings:

| Rule | File | Lines | Triage State | Notes |
| --- | --- | --- | --- | --- |
| `generic-api-key` | `examples/self-trace/evidence-events.json` | 169, 218 | `allowlisted_fixture` | `dedupe_key` labels; allowlist is path+regex scoped. |
| `generic-api-key` | `examples/self-trace/negative-native-policy-field.json` | 178, 227 | `allowlisted_fixture` | Mirrored self-trace labels; allowlist is path+regex scoped. |
| `generic-api-key` | `examples/self-trace/assessment-input.json` | 174, 223 | `allowlisted_fixture` | Mirrored self-trace labels; allowlist is path+regex scoped. |
| `jwt` | `internal/witness/profiles_test.go` | 959, 1008 | `allowlisted_test_fixture` | JWT sentinel used by tests; allowlist is path+regex scoped. |
| `private-key` | `specs/004-mvp-readiness-hardening/pr-review/ec8db52/packet/inputs/diff.patch` | 45660, 46637 | `allowlisted_historical_fixture` | Historical review fixture contains private-key marker text; allowlist is path+regex scoped. |

The working-tree scan reported four additional `.codex-subagents/` findings
from local run clutter. Those files are not intended to be tracked and are
already covered by repository hygiene rules.

## Required Follow-Up

- Review all `G703`, `G304`, and `G204` call sites first because they touch
  path, network, external-input, and command-execution boundaries.
- Re-run `gitleaks` with `.gitleaks.toml` after each fixture or scanner-rule
  change; do not broaden allowlist patterns without review.
- Add CI jobs for `govulncheck`, `gosec`, and `gitleaks` only after local
  baselines are stable enough to avoid turning known fixture hits into noise.
- Re-run this baseline after any security, trust, path, witness, release-proof,
  repo-observer, or PR-review change.

## Reproduction Commands

```text
go vet ./...
govulncheck ./...
gosec -fmt=json -out <scan-output-dir>/gosec.json ./...
gitleaks detect --source . --no-git --redact=100 --config .gitleaks.toml
git archive --format=tar HEAD | tar -xf - -C <tracked-scan-dir>
cp .gitleaks.toml <tracked-scan-dir>/.gitleaks.toml
gitleaks detect --source <tracked-scan-dir> --no-git --redact=100 --config <tracked-scan-dir>/.gitleaks.toml
```
