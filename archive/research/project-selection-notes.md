# Project Selection Notes

Status: source-backed notes
Created: 2026-04-29

## Selection Criteria

Each round-1 project must satisfy at least four conditions:

- large enough to stress context management
- active enough to reflect current OSS practice
- JVM-centered
- build system can be detected from repo files
- useful for enterprise/client adoption conversations
- exposes a distinct failure mode for AI agents

## Chosen Projects

### Spring Framework

Source: https://github.com/spring-projects/spring-framework

Why selected:

- major enterprise Java framework
- large modular repo
- Gradle-based build
- Java with some Kotlin
- useful for client CTO discussions

### Apache Kafka

Source: https://github.com/apache/kafka

Why selected:

- Java and Scala distributed system
- Gradle build
- large, mature Apache project
- good stress test for module and protocol comprehension

### JetBrains Kotlin

Source: https://github.com/JetBrains/kotlin

Why selected:

- Kotlin-first codebase
- compiler and tooling complexity
- Gradle build
- tests whether agents understand Kotlin without treating it as Java

### Apache Flink

Source: https://github.com/apache/flink

Why selected:

- Java and Scala streaming platform
- Maven build
- very large multi-module repository
- useful Maven contrast to Gradle-heavy projects

### Bazel

Source: https://github.com/bazelbuild/bazel

Why selected:

- Bazel-native Java project
- validates Java+Bazel and monorepo-style assumptions
- build-system project with Starlark and Java interaction
- directly targets the discovery feedback about Bazel support

## Deferred Candidates

| Project | Reason deferred |
|---|---|
| Gradle | Strong candidate, but Spring/Kafka/Kotlin already cover Gradle in round 1. Use in round 2 for build-tool dogfood. |
| Elasticsearch | Strong candidate, but licensing/product-boundary concerns make Apache/Spring/JetBrains/Bazel cleaner for round 1. |
| Cassandra | Strong Java distributed system, but overlaps with Kafka/Flink for round 1. |
| Hadoop | Important ecosystem, but slower feedback loop and older build conventions. |
