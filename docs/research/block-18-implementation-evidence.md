# Block 18 Implementation Evidence

Status: implementation evidence recorded; strict review and PR-level review are
still open under T158.

## Scope

Block 18 adds explicit forensic retention assessment without changing `gate`
semantics or managed-harness behavior.

Implemented surfaces:

- `internal/forensic`: standalone evaluator for
  `forensic_retention` assessment facts.
- `sdp-trace assess --profile forensic-retention`: run + redaction-policy input,
  deterministic exit behavior, JSON output, preview, and explain.
- `schema/assessment-result.schema.json`: `oneOf` for Block 17 managed and
  Block 18 forensic result shapes.
- `schema/redaction-policy.schema.json`: portable redaction policy shape with
  FR-054 retention modes, redaction actions, authority, critical families, and
  withholding policy.
- `schema/flight-recorder-event.schema.json` and
  `schema/flight-recorder-run.schema.json`: Block 18 redaction policy refs,
  structured authority, pre/redacted digests, retention lifecycle, raw
  reference binding, and forensic condition rows.
- `examples/block18-forensic-retention`: committed assessment-result fixtures.
- `docs/flight-recorder.md`: Go-first current commands and Block 18 forensic
  retention boundary.

## Drift Handling

Two drift points were found and fixed during implementation:

- Block 18 initially introduced retention names that diverged from FR-054. The
  implementation keeps `digest_only`, `sanitized_excerpt`,
  `encrypted_raw_ref`, `external_artifact_ref`, and `not_assessed` as the only
  retention modes.
- `docs/flight-recorder.md` referenced historical Node `.mjs` query scripts.
  The doc now uses current Go-first commands and treats query names as product
  surfaces rather than active Node tooling.

## Local Verification

Commands run from `/Users/fall_out_bug/projects/vibe_coding/sdp-trace`:

```bash
rtk go test ./...
rtk jq empty schema/*.json examples/block18-forensic-retention/*.json examples/flight-recorder/forensic-digest-only-negative/run.json examples/flight-recorder/forensic-digest-only-negative/events/003-test-output-observed.json examples/flight-recorder/redaction-unresolved/run.json examples/flight-recorder/redaction-unresolved/events/003-test-output-observed.json
rtk git diff --check
go run ./cmd/sdp-trace assess explain --assessment-result examples/block18-forensic-retention/digest-only-critical-fail.assessment-result.json
```

Observed states:

- Go tests: pass, 129 tests across 13 packages.
- JSON parse checks: pass.
- Diff whitespace check: pass.
- Explain output showed `critical_evidence_digest_only` and
  `capped_to_retention_mode: digest_only` without raw secret-like values.

Follow-up verification after strict review fixes:

```bash
rtk go test ./...
rtk jq empty schema/*.json examples/block18-forensic-retention/*.json examples/flight-recorder/forensic-digest-only-negative/run.json examples/flight-recorder/forensic-digest-only-negative/events/003-test-output-observed.json examples/flight-recorder/redaction-unresolved/run.json examples/flight-recorder/redaction-unresolved/events/003-test-output-observed.json
rtk git diff --check HEAD
rtk rg -n "TO[D]O|FIX[M]E" $(git ls-files -m -o --exclude-standard)
```

Observed states:

- Go tests: pass, 132 tests across 13 packages.
- JSON parse checks: pass.
- Diff whitespace check: pass.
- Marker scan: no matches; command exited `1` because `rg` found no
  matches.

## Remaining Open Work

- T158 strict implementation review and PR-level review are not complete.
- GitHub CI is not assessed in this local workspace.
- Source-bound release proof regeneration is not complete in this evidence
  note.
