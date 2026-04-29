# JVM Benchmark Bootstrap Smoke Summary

Status: completed
Date: 2026-04-29

The round-1 JVM OSS benchmark workspace was bootstrapped with shallow, blobless checkouts.

Local checkout path:

```text
benchmarks/repos/jvm-oss/
```

Local smoke notes path:

```text
.sdp-trace-runs/bootstrap/jvm-oss/
```

Both paths are ignored by git.

## Results

| Project | Commit | Expected build | Detected Gradle | Detected Maven | Detected Bazel | Checkout size |
|---|---|---|---:|---:|---:|---:|
| Spring Framework | `3184eb3acc8e2e6e95623c43e979ef6c256887ed` | Gradle | yes | no | no | 102M |
| Apache Kafka | `cb7e3ab375747bfafb97c746a2049443c2ddd9a1` | Gradle | yes | no | no | 111M |
| JetBrains Kotlin | `64dff70c04700f4748a701cf49f7114e8645c57f` | Gradle | yes | no | no | 700M |
| Apache Flink | `f325b4a8d5914d218b394c8a7e2f7e4f0c27a358` | Maven | no | yes | no | 309M |
| Bazel | `5347ca870104f5c6aced1c3a63671c1492ea7bbd` | Bazel | no | no | yes | 258M |

## Observations

- Build-system smoke detection matched the expected build system for all five projects.
- Kotlin is much larger than the other round-1 projects even with shallow blobless clone; expect slower Phase A runs and higher context pressure.
- Bazel has `MODULE.bazel`, top-level `BUILD`, and Bazel config markers. It is a good first target for testing whether agents avoid Maven/Gradle assumptions.
- Flink provides the Maven contrast for the mostly Gradle-heavy set.

## Next Step

Run Phase A read-only intake on one project first:

1. Bazel, because it targets the discovery failure mode directly.
2. Spring Framework, because it is the enterprise Java baseline.

Do not run the full 50-run matrix until these two smoke Phase A runs prove the prompt and artifact contract are usable.
