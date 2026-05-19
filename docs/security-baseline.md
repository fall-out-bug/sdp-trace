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
| Source commit | `adf0e0220c7412171d6ca279d52c56ec83f459e4` |
| Integration branch | `codex/016-018-pi-delivery` |
| Go version | `1.22` from `go.mod` |
| Local evidence files | `/tmp/sdp-trace-gosec.json`, `/tmp/sdp-trace-gitleaks.json`, `/tmp/sdp-trace-gitleaks-tracked.json` |

The scanner evidence was produced locally before later integration edits in
this branch. Re-run the commands below before using this ledger in a PR.

## Tool Summary

| Tool | Result | Finding Count | State | Notes |
| --- | --- | ---: | --- | --- |
| `go vet ./...` | pass | 0 | `verified` | Go vet completed locally. |
| `govulncheck ./...` | pass | 0 | `verified` | No known vulnerabilities reported by the local database snapshot. |
| `gosec ./...` | exit 1 | 133 | `needs_triage` | Static findings require code-level review or reviewed suppressions. |
| `gitleaks detect --source . --no-git --redact=100` | exit 1 | 14 | `needs_triage` | Includes local `.codex-subagents/` clutter; not suitable as tracked-source evidence. |
| tracked-source `gitleaks` snapshot | exit 1 | 10 | `needs_triage` | Scanned an archive of tracked `HEAD`; findings are fixture/test candidates but not allowlisted. |

## `gosec` Family Summary

| Rule | Count | Scanner Severity | Triage State | Current Disposition |
| --- | ---: | --- | --- | --- |
| `G304` | 61 | medium | `needs_triage` | Variable file reads/opens. Must verify path provenance and traversal controls per call site. |
| `G301` | 28 | medium | `needs_triage` | Directory permissions `0755`. Must confirm generated evidence readability is intentional for each path. |
| `G306` | 24 | medium | `needs_triage` | File permissions `0644`. Must confirm public-readable evidence output is acceptable for each artifact class. |
| `G204` | 11 | medium | `needs_triage` | Subprocess launches. Must distinguish intended wrapper behavior from untrusted command execution. |
| `G703` | 5 | high | `needs_triage` | Potential path traversal in repo observer / PR review paths. Review before any production-adoption claim. |
| `G101` | 2 | high | `candidate_false_positive` | Semantic strings matched credential heuristics; verify before allowlisting. |
| `G302` | 1 | medium | `needs_triage` | `chmod 0644` on generated evidence file; requires artifact-sensitivity review. |
| `G115` | 1 | high | `candidate_false_positive` | `rune` to `byte` conversion appears guarded by `r <= 127`; verify before suppressing. |

No `gosec` finding is treated as closed by this document. The next code slice
should either fix, narrow, or add a reviewed suppression for each accepted
false positive.

## Tracked `gitleaks` Findings

The tracked-source snapshot reported 10 findings:

| Rule | File | Lines | Triage State | Notes |
| --- | --- | --- | --- | --- |
| `generic-api-key` | `examples/self-trace/evidence-events.json` | 169, 218 | `fixture_candidate` | `dedupe_key` labels; needs allowlist or fixture rewrite. |
| `generic-api-key` | `examples/self-trace/negative-native-policy-field.json` | 178, 227 | `fixture_candidate` | Mirrored self-trace labels; keep mirrors synchronized if edited. |
| `generic-api-key` | `examples/self-trace/assessment-input.json` | 174, 223 | `fixture_candidate` | Mirrored self-trace labels; keep mirrors synchronized if edited. |
| `jwt` | `internal/witness/profiles_test.go` | 959, 1008 | `test_fixture_candidate` | JWT sentinel used by tests; needs reviewed fixture pattern or allowlist. |
| `private-key` | `specs/004-mvp-readiness-hardening/pr-review/ec8db52/packet/inputs/diff.patch` | 45660, 46637 | `historical_fixture_candidate` | Historical review fixture contains private-key marker text; needs policy decision. |

The working-tree scan reported four additional `.codex-subagents/` findings
from local run clutter. Those files are not intended to be tracked and are
already covered by repository hygiene rules.

## Required Follow-Up

- Review all `G703`, `G304`, and `G204` call sites first because they touch
  path, network, external-input, and command-execution boundaries.
- Decide whether fixture-looking `gitleaks` hits should be rewritten or added
  to a reviewed scanner allowlist.
- Add CI jobs for `govulncheck`, `gosec`, and `gitleaks` only after local
  baselines are stable enough to avoid turning known fixture hits into noise.
- Re-run this baseline after any security, trust, path, witness, release-proof,
  repo-observer, or PR-review change.

## Reproduction Commands

```text
go vet ./...
govulncheck ./...
gosec -fmt=json -out=/tmp/sdp-trace-gosec.json ./...
gitleaks detect --source . --no-git --redact=100 --report-format json --report-path /tmp/sdp-trace-gitleaks.json
git archive --format=tar HEAD | tar -xf - -C /tmp/sdp-trace-tracked-scan
gitleaks detect --source /tmp/sdp-trace-tracked-scan --no-git --redact=100 --report-format json --report-path /tmp/sdp-trace-gitleaks-tracked.json
```
