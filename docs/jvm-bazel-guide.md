# JVM And Bazel Guide

Go and JVM are first-class stack targets.

## Minimum JVM Matrix

- Java + Maven
- Java + Gradle
- Kotlin + Gradle
- Kotlin + Bazel
- mixed JVM monorepo
- scoped service assessment inside a monorepo

## Required Behavior

- Do not apply Java-only heuristics to Kotlin without an explicit caveat.
- Do not infer Maven or Gradle when Bazel files are present.
- Prefer scoped service assessment over root-level monorepo assessment.
- Emit `not_assessed` when build/test semantics cannot be determined.

## Good Scope Input

```text
Assess service //payments/api.
Language: Kotlin.
Build: Bazel.
Do not infer Maven or Gradle conventions unless files prove them.
```
