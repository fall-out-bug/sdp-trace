# JVM Benchmark Bootstrap Runbook

This runbook prepares local checkouts for Phase A read-only intake.

## Scope

Only `sdp-trace` owns this benchmark bootstrap. `sdp-gate` consumes sanitized trace outputs later.

## Prerequisites

- `git`
- `jq`
- network access to GitHub

## Dry Run

```bash
scripts/bootstrap-jvm-oss-benchmark.sh --dry-run
```

## Bootstrap One Project

```bash
scripts/bootstrap-jvm-oss-benchmark.sh --project spring-framework
```

## Bootstrap All Round-1 Projects

```bash
scripts/bootstrap-jvm-oss-benchmark.sh
```

The script creates shallow, blobless checkouts under:

```text
benchmarks/repos/jvm-oss/
```

It writes smoke notes under:

```text
.sdp-trace-runs/bootstrap/jvm-oss/
```

Both paths are ignored by git. Commit only sanitized research summaries under `docs/research/`.

Round-1 disk footprint after the first bootstrap was about 1.5G. Kotlin is the largest checkout.

## Smoke Detection

The bootstrap script records:

- checked-out commit
- expected build system from the manifest
- top-level Gradle/Maven/Bazel markers
- top-level files

This is not a full assessment. It only proves that the benchmark workspace is ready for Phase A.
