# Phase A Prompt: Read-Only Repo Intake

Use this prompt for every project and harness/model combination.

```text
You are assessing a large JVM OSS repository for traceability and evidence readiness.

Rules:
- Read repository files before making claims.
- Identify the build system from files, not assumptions.
- Do not apply Java-only heuristics to Kotlin or Scala without caveat.
- Do not edit files.
- If evidence is missing, write not_assessed.
- Separate inspected facts from model judgment.

Produce:
1. repo-map.md
2. evidence-bundle.json compatible with sdp-trace
3. trace.json compatible with sdp-trace
4. notes.md with unsupported claims and uncertainty

Assessment questions:
- What languages are present?
- What build systems are present?
- What are the likely scoped modules/services for future experiments?
- Which commands appear to be canonical for build/test?
- What evidence can be collected without running the full build?
- What cannot be assessed yet?
```
