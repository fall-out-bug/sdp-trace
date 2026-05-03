# JVM And Bazel Guide

Go and JVM are planned assessment targets. Observed behavior is row-specific and remains `not_assessed` until committed evidence exists.

## Minimum JVM Matrix

- Java + Maven
- Java + Gradle
- Kotlin + Gradle
- Kotlin + Bazel
- mixed JVM monorepo
- scoped service assessment inside a monorepo

## Required Behavior

- Do not apply Java-only heuristics to Kotlin without an explicit caveat.
- Infer Bazel ownership only from scope-specific `BUILD` or `BUILD.bazel` target evidence, `MODULE.bazel`, `WORKSPACE`, or `WORKSPACE.bazel` context tied to the assessed target.
- Treat `.bazelrc` as supporting configuration evidence only.
- Treat Maven or Gradle files as dependency metadata only when scoped Bazel evidence proves Bazel ownership for the assessed target.
- Treat Kotlin dependencies as ecosystem context only; Kotlin service-language evidence requires `.kt`, `.kts`, `kt_jvm_*`, or Kotlin compiler/toolchain rules tied to the assessed scope.
- Prefer scoped service assessment over root-level monorepo assessment.
- Emit `not_assessed` when build/test semantics cannot be determined.

## Good Scope Input

```text
Assess service //payments/api.
Language: Kotlin.
Build: Bazel.
Do not infer Maven or Gradle conventions unless files prove them.
```
