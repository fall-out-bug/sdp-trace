# Block 20 Implementation Evidence

Status: implementation verification and strict implementation review recorded.
PR-level review, GitHub checks, merge verification, and post-merge cleanup are
not assessed yet because no PR exists yet.

## Scope

Block 20 adds `forensics-basic-v1`, a read-only query pack over Block 09 run
facts, Block 18 forensic-retention assessment facts, and Block 19 adapter-capture
assessment facts.

Implemented surfaces:

- `internal/query`: query-pack result assembly, deterministic row ids,
  row-state propagation, input artifact digests, source refs, explain rendering,
  output-safety assertions, and malformed-input handling.
- `sdp-trace query-pack --pack forensics-basic-v1 --run <run-dir> --out <file>`.
- `sdp-trace query-pack explain --result <file>`.
- `schema/forensics-query-pack-result.schema.json`: Block 20 query-pack result
  contract with closed row states and evidence families.
- `examples/block20-forensics-query-pack`: fixture matrix and committed
  scenarios for mixed evidence, digest-only caps, missing optional/required
  artifacts, unsupported observers, unresolved redaction, supersession, unverified
  claims, unsafe provider refs, malformed inputs, witness gaps, and file-mutation
  gaps.

## Boundary Checks

Block 20 remains an evidence view:

- no top-level native policy verdict is emitted;
- `pass` / `fail` can appear only as upstream source condition states, not row
  evidence states;
- explain output renders JSON result fields only and does not infer conclusions;
- output safety is asserted against serialized JSON and explain output, not
  against upstream raw material.

Readable malformed required inputs that can be hashed produce deterministic
`cannot_verify` query rows. OS-level unreadable inputs still exit non-zero
because no digest-backed valid result artifact can be produced.

## Local Verification

Commands run from `/Users/fall_out_bug/projects/vibe_coding/sdp-trace-block20`:

```bash
rtk go test -count=1 ./...
rtk jq empty schema/*.json
find examples/block20-forensics-query-pack -name '*.json' \
  ! -path '*/malformed-input/adapter-capture.assessment-result.json' \
  ! -path '*/malformed-run/run.json' \
  ! -path '*/malformed-forensic-retention/forensic-retention.assessment-result.json' \
  -print0 | xargs -0 rtk jq empty
rtk jq -c . examples/block20-forensics-query-pack/fixture-matrix.jsonl >/dev/null
rtk git diff --check
rtk rg -n <deferred-work-marker-pattern> <changed Block 20 files>
```

Observed states:

- Full Go tests: pass, 160 tests across 15 packages.
- Schema syntax checks: pass.
- Block 20 fixture syntax checks: pass, excluding three intentionally malformed
  fixtures listed above.
- Fixture matrix JSONL parse: pass.
- Whitespace diff check: pass.
- Deferred-work marker scan: no matches in changed Block 20 files.

## Strict Review

Implementation review ran separate code/correctness, tracing/evidence, and
requirements-vs-implementation planes.

- Code/correctness: ZAI/GLM-5.1 found a major digest-only reason-code mismatch
  and a minor explain source-ref omission. Both were fixed. Focused re-review
  returned `APPROVE` with no remaining critical or major findings.
- Tracing/evidence: MiniMax-M2.7 found fixture matrix assertion and breadth
  gaps. The matrix now checks source refs, source condition states, evidence
  gaps, every committed scenario directory, malformed forensic retention,
  witness gaps, file-mutation gaps, and supersession source refs. Second focused
  re-review returned `APPROVE`.
- Requirements-vs-implementation: Kimi K2P6 found malformed readable run
  handling, event-family source refs, full safety-class test coverage,
  reconstructability handling, explain source refs, and timeline optional
  artifact rows. Valid findings were fixed. Second focused re-review returned
  `APPROVE` with no remaining critical or major findings.

## Remaining Open Work

- PR-level review is `not_assessed`; no PR exists yet.
- GitHub CI is `not_assessed`; no PR exists yet.
- Post-merge `origin/main` verification and cleanup are `not_assessed`.
