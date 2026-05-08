# Block 05 Pi Review Ledger

Status: implementation pi-review findings closed
Parent Spec: `001-sdp-trace-time-series-evidence-substrate`
Beads mirror: `sdp-trace-cdn.22`

## Purpose

This ledger is the committed SpecKit review record for Block 05. Beads mirrors execution tracking, but this file is the repository-observable source for pi-review findings, severity, disposition, and closure evidence.

## Review Inputs

- `blocks/05-customer-pilot-evidence-package.md`
- `blocks/05-customer-pilot-socratic.md`
- `blocks/05-customer-pilot-implementation-plan.md`
- `spec.md`
- `tasks.md`
- Existing pilot target artifacts under retired research artifacts, retired static harness matrix, retired static model matrix, `docs/jvm-bazel-guide.md`, and `examples/jvm-bazel/`

## Review Findings

| ID | Severity | Beads mirror | Finding | Disposition | Evidence |
|---|---|---|---|---|---|
| F001 | critical | `sdp-trace-cdn.22.8` | Native support/compatibility claims leaked into the canonical contract. | Accepted; spec layer revised. | `spec.md` FR-014/FR-036/FR-037 and Block 05 spec now require observed evidence state or external verdict input. |
| F002 | major | `sdp-trace-cdn.22.6` | Native `TBD` and `external_verdict_recorded` row states weakened the evidence boundary. | Accepted; spec layer revised. | Block 05 spec now uses `observed` or `not_assessed` with reason codes; external verdicts are linked evidence inputs. |
| F003 | major | `sdp-trace-cdn.22.10`, `sdp-trace-cdn.22.3` | Beads was the only review-closure evidence. | Accepted; fixed by this committed ledger. | This file records findings, severity, disposition, evidence, and Beads mirrors. |
| F004 | major | `sdp-trace-cdn.22.7`, `sdp-trace-cdn.22.4`, `sdp-trace-cdn.22.22` | Claim-boundary and matrix validation plan was grep-only and too narrow. | Accepted; implementation fixed. | `scripts/validate-pilot-matrices.mjs` checks required matrix columns, required target rows, artifact references for `observed`, reason codes, and banned target-file verdict tokens; `scripts/validate-contracts.sh` runs it. |
| F005 | minor/P3 | `sdp-trace-cdn.22.11` | Customer inputs were not explicitly separated from committed artifacts. | Accepted; spec layer revised. | Block 05 spec and Socratic dialogue state private customer inputs are never committed. |
| F006 | minor/P3 | `sdp-trace-cdn.22.13` | Internal "unsupported claim" wording conflicted with prohibited native verdict wording. | Accepted; spec layer revised. | Block 05 spec and root spec now use `unbacked_claim`. |
| F007 | minor/P3 | `sdp-trace-cdn.22.14` | Tasks used "Superpowers" where "Superpowers-style" was required. | Accepted; task text revised. | T028 now says Superpowers-style harness pattern and states no Superpowers dependency is introduced. |
| F008 | critical | `sdp-trace-cdn.22.12` | Existing JVM+Bazel placeholder records `status: pass` for synthetic behavior. | Accepted; implementation fixed. | `examples/jvm-bazel/evidence-bundle.json` now keeps real Kotlin+Bazel behavior `not_assessed` with `design_fixture_only`/`no_run_artifact`. |
| F009 | major | `sdp-trace-cdn.22.18` | Existing run-card template owns native pass/warn/fail verdicts. | Accepted; implementation fixed. | retired research artifact now uses evidence states, reason codes, artifact references, and optional external-verdict input. |
| F010 | major | `sdp-trace-cdn.22.16` | Harness matrix hides evidence-state boundaries and combines `gsd`/`gsd2`. | Accepted; implementation fixed. | retired static harness matrix now splits `gsd`/`gsd2` and records evidence state, reason code, artifact reference, external verdict reference, gap reason, and next evidence. |
| F011 | major | `sdp-trace-cdn.22.15` | Model matrix omits OpenCode+MiniMax/Kimi/GLM rows. | Accepted; implementation fixed. | retired static model matrix now includes explicit OpenCode+MiniMax, OpenCode+Kimi, and OpenCode+GLM rows as `not_assessed`. |
| F012 | major | `sdp-trace-cdn.22.17` | JVM+Bazel guide overstates build ownership inference in hybrid repos. | Accepted; implementation fixed. | `docs/jvm-bazel-guide.md` now makes Bazel ownership scope-specific and treats `.bazelrc`, Maven/Gradle, and Kotlin dependency metadata as supporting context only. |
| F013 | major | `sdp-trace-cdn.22.21` | Plan hard-coded local RTK/Beads operations as canonical verification. | Accepted; spec plan revised. | Implementation plan now uses portable commands as canonical; `rtk` and Beads remain local operator tools. |
| F014 | minor/P3 | `sdp-trace-cdn.22.19` | `.bazelrc` alone was allowed to prove Bazel ownership. | Accepted; spec layer revised. | Block 05 spec now treats `.bazelrc` as supporting configuration only. |
| F015 | minor/P3 | `sdp-trace-cdn.22.20` | Kotlin dependencies were treated as service-language proof. | Accepted; spec layer revised. | Block 05 spec now separates Kotlin source/rule evidence from dependency metadata. |
| F016 | minor/P3 | `sdp-trace-cdn.22.23` | "First-class stack targets" wording reads like support without evidence. | Accepted; implementation fixed. | `docs/jvm-bazel-guide.md` now says Go and JVM are planned assessment targets and keeps observed behavior row-specific. |
| F017 | major | `sdp-trace-cdn.22.2`, `sdp-trace-cdn.22.1` | Spec gate leaves manifest validation stale after spec/task changes. | Accepted; fixed for spec gate. | Manifest, local DSSE, release verification, and self-attestation artifacts were resynchronized; `npm run validate` passed after the fix. |
| F018 | major | `sdp-trace-cdn.22.5` | `evidence_backed` could over-credit synthetic fixtures. | Accepted; spec layer revised. | Block 05 spec now uses `observed`; synthetic fixtures may support schema/design coverage only. |
| F019 | minor/P3 | `sdp-trace-cdn.22.9` | Package naming blurred Block 05 design and customer handoff outline. | Accepted; spec layer revised. | Block file title now says "Block 05 Design"; retired research artifact remains the T037 handoff outline. |
| F020 | major | `sdp-trace-cdn.22.26` | Ledger overstated implementation closure before implementation pi-review completion. | Accepted; fixed. | Ledger status now records implementation pi-review fixes in progress until this review loop closes. |
| F021 | major | `sdp-trace-cdn.22.25`, `sdp-trace-cdn.22.34` | Harness matrix level taxonomy could read like native capability claims. | Accepted; fixed. | Harness matrix now removes the L0-L4 taxonomy and remains only an evidence-state table. |
| F022 | major | `sdp-trace-cdn.22.27` | Generic run-card template asked for exact prompt without redaction/approval boundary. | Accepted; fixed. | Run-card template now requires redacted prompt, prompt summary, or prompt SHA-256 unless release approval exists. |
| F023 | major | `sdp-trace-cdn.22.24` | External verdict availability was allowed as a native matrix reason code. | Accepted; fixed. | `external_verdict_available` was removed from allowed matrix reason codes. External verdicts remain references. |
| F024 | minor/P3 | `sdp-trace-cdn.22.28` | Customer package checkpoint used readiness wording. | Accepted; fixed. | Checkpoint renamed to `Run preflight`. |
| F025 | major | `sdp-trace-cdn.22.30` | JVM+Bazel fixture omitted Kotlin source/rule evidence placeholder. | Accepted; fixed. | Fixture now includes a `not_assessed` Kotlin source/rule evidence item. |
| F026 | major | `sdp-trace-cdn.22.29` | Matrix validator could over-credit design fixtures for `observed` rows. | Accepted; fixed. | Validator parses JSON artifact references for `observed` rows and rejects placeholder/not_assessed artifacts. |
| F027 | major | `sdp-trace-cdn.22.32` | Run-card validation was parse-only for provenance and trace artifacts. | Accepted; fixed. | OpenCode, harness, and customer package docs now schema-validate trace snapshots and split provenance arrays into schema-validated records. |
| F028 | major | `sdp-trace-cdn.22.31` | Generic run-card template omitted Block 05 required artifacts. | Accepted; fixed. | Template now lists evidence bundle, provenance records, trace snapshot, export limitations, redaction note, and optional handoff artifacts. |
| F029 | minor/P3 | `sdp-trace-cdn.22.33` | Customer package validation used fragile globs for optional handoff files. | Accepted; fixed. | Handoff validation is now conditional per optional file and schema-validates files when present. |
| F030 | major | `sdp-trace-cdn.22.37` | Release proof wording overread contract-foundation proof coverage for Block 05 artifacts. | Accepted; fixed. | Contract manifest now includes Block 05 run-cards, matrices, customer package outline, JVM guide, review ledger, and matrix validator before proof sync. |
| F031 | major | `sdp-trace-cdn.22.36` | Release verification source state differed from self-attestation source state. | Accepted; fixed. | Release verification now records the same source-content mismatch state and counts as self-attestation for the current manifest/source commit. |
| F032 | major | `sdp-trace-cdn.22.35` | Matrix/run-card claim-boundary validator was narrower than ledger claims. | Accepted; fixed. | Validator now scans Block 05 target docs for native verdict/support wording unless the line includes explicit boundary wording. |

## Current Closure State

- Spec-layer findings F001-F007, F013-F015, F017-F019 have fixes applied in SpecKit artifacts and local proof artifacts. Their Beads mirrors are closed.
- Implementation findings F004, F008-F012, F016, and F020-F032 have fixes applied in docs, matrices, fixtures, validation script, and local proof artifacts. Their Beads mirrors are closed.

## Validation Notes

Spec-review validation after fixes:

- `jq empty schema/*.json`: passed.
- `git diff --check`: passed.
- `npm run validate`: passed after contract-foundation manifest, local DSSE, release verification, and self-attestation artifacts were resynchronized.
- External production trust remains `not_assessed`; no `trusted_contract_release: true` claim was introduced.

Implementation-target validation after proof sync:

- `node scripts/validate-pilot-matrices.mjs`: passed.
- `node scripts/validate-json-schema.mjs schema/evidence-bundle.schema.json examples/jvm-bazel/evidence-bundle.json`: passed.
- Targeted banned-token scan for `pass / warn / fail`, JSON `status: pass|warn|fail`, and `first-class stack targets`: no matches.
- `npm run validate`: passed after self-trace hash updates and contract-foundation manifest, local DSSE, release verification, and self-attestation artifact sync.
- `git diff --check`: passed.
