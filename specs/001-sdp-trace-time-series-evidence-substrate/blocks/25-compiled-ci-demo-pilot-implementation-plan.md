# Block 25 Implementation Plan: Compiled CI Demo Pilot

Status: draft, pending Socratic review and user approval.

## Execution Rule

Do not implement this plan until `25-compiled-ci-demo-pilot.md` has passed
Socratic review and the reviewed direction is explicitly approved.

## Workspaces

Use separate worktrees or checkouts:

- `sdp-trace`: SpecKit docs, sanitized report, review ledgers, and PR.
- `sdp-trace-demo-ci-pilot`: demo source, Bazel config, CI workflow, scripts,
  README, and demo PR or direct branch depending on repository policy.

No implementation slice may depend on uncommitted generated artifacts as proof.

## Slice 1: Demo Compiled JVM Target

Owner surface:

- demo repo `MODULE.bazel`, `BUILD.bazel`, `app/BUILD.bazel`;
- demo repo `app/src/**`;
- demo repo test source under `app/src/test/**` or equivalent.

Tasks:

1. Replace the primary proof target with a compiled JVM target.
2. Keep Feature Flag / Entitlements behavior small and inspectable.
3. Add semantic tests for plan, entitlement, and seat-overage behavior.
4. Keep source/scope shell checks only if they are clearly secondary.
5. Record `rules_kotlin`, JVM version, and test framework in README and CI
   evidence.

Expected verification:

- `bazel test //...`
- focused review by Head of Engineering lens.

## Slice 2: CI Artifact Generation And Index Integrity

Owner surface:

- demo repo `.github/workflows/sdp-trace-demo.yml`;
- demo repo `scripts/write-artifact-index.sh`;
- demo repo `scripts/verify-artifact-index.sh`;
- demo repo `scripts/redaction-scan.sh`;
- demo repo `.gitignore`.

Tasks:

1. Ensure `.sdp-trace-report/`, `.sdp-trace-runs/`, and `.sdp-trace-tools/`
   remain ignored and are generated only by CI for proof.
2. Generate artifact index without self-indexing.
3. Add an index verification step that recomputes listed digests before upload.
4. Make the index JSON deterministic: relative paths, sorted order, SHA-256
   lowercase hex, and file sizes.
5. Implement the verifier as a self-contained portable shell script that
   detects `sha256sum` on Linux and falls back to `shasum -a 256` on macOS.
6. Preserve `verify`, `explain`, `report`, `gate`, and `witness` outputs.
7. Run redaction scan over clean and negative artifact roots, recording command,
   pattern digest, exit code, and state.
   The scanner output must report only relative filenames, match count, and
   pass/fail/cannot-verify state; it must not emit matched content, line
   fragments, full private paths, the full pattern file, or scanned file
   contents.
8. Keep the no-OIDC job as an intentional `cannot_verify` case.

Expected verification:

- local script unit check using a temporary artifact directory;
- GitHub Actions run success;
- artifact download plus digest recomputation.
- redaction scan pass over downloaded artifact sets.

## Slice 3: Negative Evidence Cases

Owner surface:

- demo repo `scripts/write-dishonest-cases.sh`;
- demo repo generated artifact expectations;
- sanitized `sdp-trace` report after successful CI run.

Tasks:

1. Keep absent or unresolvable source/run binding as `cannot_verify`.
2. Keep contradictory source/run binding as `fail`, or record the contradictory
   variant as `not_assessed` if Block 25 implements only the absent-binding
   variant.
3. Keep stale digest as `fail`, but make it independent of the fixed clean
   artifact index by mutating a non-index artifact after clean index generation
   while leaving the clean index entry unchanged.
4. Keep no-OIDC witness gap as `cannot_verify`.
5. Document that negative cases are intentional evidence-state checks, not
   product failures.

Expected verification:

- downloaded artifacts contain all three cases;
- role review accepts state separation.

## Slice 4: Documentation And Sanitized Report

Owner surface:

- demo repo `README.md`;
- `sdp-trace` retired research artifact;
- `sdp-trace` Block 25 review ledger.

Tasks:

1. Update demo README with compiled-target scope, CI prerequisites, artifact
   retention, and residual `not_assessed` states.
2. Add sanitized Block 25 report and artifact index summary to `sdp-trace` only
   after CI artifacts are downloaded and verified.
3. Preserve evidence-led wording with no questionnaire, sales deck, or coverage
   matrix framing.
4. Record role-review findings and dispositions in
   `blocks/25-compiled-ci-demo-pilot-review-ledger.md`.

Expected verification:

- current-facing wording scan for forbidden framing;
- `go test ./...`;
- `jq empty schema/*.json`;
- `go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc`;
- `git diff --check`.

## Integration Audit

Before PR readiness:

1. Confirm demo repo head, CI run id, and artifact ids match the sanitized
   `sdp-trace` report.
2. Confirm artifact expiration is recorded.
3. Confirm every digest cited in `sdp-trace` was recomputed from downloaded
   artifacts.
4. Confirm no raw logs, raw OIDC material, private URLs, tokens, or local
   filesystem paths were copied into `sdp-trace` or uploaded CI artifacts.
5. Confirm every remaining buyer-relevant gap is `not_assessed` or
   `cannot_verify`, not silently omitted.

## PR And Merge

Use the normal trust workflow:

1. Commit demo repo slices with focused verification.
2. Run demo repo role review and fix valid critical/major findings.
3. Commit `sdp-trace` sanitized report and review ledger.
4. Open or update `sdp-trace` PR.
5. Run PR-level review across code/correctness, tracing/evidence,
   requirements-vs-implementation, and security/privacy.
6. Merge only after fresh CI, role review, and PR-level review have no remaining
   critical or major findings.
