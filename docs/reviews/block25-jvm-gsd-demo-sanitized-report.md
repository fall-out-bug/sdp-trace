# Block 25 JVM/GSD Demo Sanitized Report

Status: evidence report, not approval.

This report records the active `fall-out-bug/sdp-trace-demo-jvm-gsd` Block 25
demo evidence after the v2 packet/bundle track replaced the older artifact-index
shape. It is not a schema fixture, customer questionnaire, sales narrative,
merge approval, release approval, production trust claim, or signed external
attestation.

## Scope

The selected demo target is the JVM/Bazel GSD demo repository
`fall-out-bug/sdp-trace-demo-jvm-gsd`.

Accepted scope:

- Bazel build and test over the demo repository.
- Packet and bundle artifacts uploaded by GitHub Actions.
- Downloaded-artifact digest replay over retained packet/bundle contents.
- Explicit negative states for missing CI OIDC evidence, stale artifact digest,
  and source/run mismatch.

Out of scope and still `not_assessed`:

- production trust;
- owner independence;
- non-GitHub portability;
- release-binary acquisition;
- broad JVM/Bazel compatibility beyond the selected demo repository;
- signed external attestation.

## CI And Artifact Evidence

| Evidence | State |
| --- | --- |
| Demo repository | `fall-out-bug/sdp-trace-demo-jvm-gsd` |
| T213 verifier commit | `3a9491f734e5214c72014db5d893f125eb254a11` |
| T214 negative matrix commit | `a4d1f755552ba1f411af5edcb7d6caf24a9c39bf` |
| T213 artifact replay source run | `25724386343`, `success`, event `push`, head `b1018826a1d7f208c5386a44e33e62a95cbc4d05` |
| T214 matrix CI run | `26447797437`, `success`, event `pull_request`, head `b4dc843b148cccdbb2e98566c0564c96416de5dd` |

Downloaded artifact summary for run `25724386343`:

| Artifact | ID | Size bytes | Expires | State |
| --- | ---: | ---: | --- | --- |
| `change-evidence-packets` | `6939898251` | `14801` | `2026-08-10T09:00:23Z` | available |
| `evidence-bundles` | `6939898549` | `17192` | `2026-08-10T09:00:23Z` | available |

Downloaded artifact summary for run `26447797437`:

| Artifact | ID | Size bytes | Expires | State |
| --- | ---: | ---: | --- | --- |
| `change-evidence-packets` | `7215057226` | `14801` | `2026-08-24T12:22:19Z` | available |
| `evidence-bundles` | `7215057239` | `17192` | `2026-08-24T12:22:19Z` | available |

## Artifact Replay Summary

The v2 artifact verifier added in demo PR #25 covers the trust property behind
the older Block 25 artifact-index wording:

- deterministic JSON index generation;
- no self-index entry;
- retained file digest recomputation;
- retained file size recomputation;
- stale content detection.

Replay evidence:

- downloaded packet/bundle artifacts from run `25724386343` were indexed as 18
  retained entries;
- `scripts/verify-v2-artifact-index.sh` recomputed every retained digest and
  size successfully;
- an intentional mutation of `change-evidence-packets/feature-1.md` failed with
  `digest mismatch for change-evidence-packets/feature-1.md`.

The replay verifies downloaded retained artifact contents. It does not verify
GitHub artifact archive service digests because that metadata is not exposed in
the retained artifact summary; archive service digest is `not_assessed`.

## Negative State Summary

The Block 25 negative matrix added in demo PR #26 preserves the required
negative states independently of clean artifact-index correctness:

| Case | Expected state | Reason code |
| --- | --- | --- |
| `missing-ci-oidc` | `cannot_verify` | `missing_ci_oidc` |
| `stale-artifact-digest` | `fail` | `artifact_digest_mismatch` |
| `source-run-mismatch` | `fail` | `source_run_mismatch` |

GitHub run `26447797437` passed `build-and-test`, including the matrix test and
artifact upload steps.

## Residual Trust States

| Surface | State | Reason |
| --- | --- | --- |
| Production trust | `not_assessed` | Demo artifacts are not production release evidence. |
| Owner independence | `not_assessed` | The demo does not establish independent ownership or approval authority. |
| Non-GitHub portability | `not_assessed` | The selected evidence is GitHub Actions based. |
| Release-binary acquisition | `not_assessed` | No release binary acquisition path was assessed. |
| Broad JVM/Bazel compatibility | `not_assessed` | Evidence is limited to the selected demo repository and target. |
| Signed external attestation | `not_assessed` | No signed external attestation was provided. |

## Redaction Scan

The downloaded artifact roots from run `25724386343` were scanned with the
project-local Block 25 pattern file:

| Field | Value |
| --- | --- |
| Pattern file | `docs/reviews/block25-redaction-patterns.txt` |
| Pattern file SHA-256 | `494d868e528f8a017b0c320aead26ca227d70d2c31d955b1ff0d0b5e77ca52b3` |
| Scanned roots | downloaded `change-evidence-packets` and `evidence-bundles` under the artifact work directory |
| Command | `artifact_root=<downloaded-artifact-root>; rg -n --hidden -i -f docs/reviews/block25-redaction-patterns.txt "$artifact_root"` |
| Exit code | `1` |
| State | `pass`; `rg` returned no matches |

The scan covers high-signal token, password, private key, bearer header,
GitHub token, OpenAI-style key, and AWS access key patterns. It is a redaction
guard over retained artifacts, not proof that every possible secret format is
impossible.

## Verification Commands

The report is backed by these verification classes:

- demo local verification: `bazel test //... --test_output=errors`;
- demo GitHub verification: run `26447797437` `build-and-test` success;
- artifact replay: `scripts/write-v2-artifact-index.sh` followed by
  `scripts/verify-v2-artifact-index.sh` over downloaded packet/bundle artifacts;
- negative replay: intentional retained artifact mutation fails with the
  expected digest mismatch.
- redaction scan: `artifact_root=<downloaded-artifact-root>; rg -n --hidden
  -i -f docs/reviews/block25-redaction-patterns.txt "$artifact_root"`
  returned exit code `1` with no matches.

## Closure Boundary

This report supports Block 25 T215 only. It does not close T216 role review or
T217 final `sdp-trace` PR-level verification. Those tasks remain open until
their own evidence is recorded.
