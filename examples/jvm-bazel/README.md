# JVM + Bazel Example

This example is a design fixture for scoped JVM+Bazel evidence. It does not prove real Kotlin+Bazel behavior.

Current state:

- Evidence state: `not_assessed`
- Reason code: `design_fixture_only`
- Real run artifact: none

The first real Kotlin+Bazel run must commit sanitized evidence for a scoped target, including Bazel target evidence, Kotlin source or rule evidence, command summary or explicit unavailable reason, provenance, and redaction notes.
