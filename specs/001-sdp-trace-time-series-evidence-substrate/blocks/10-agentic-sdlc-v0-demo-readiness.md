# Block 10 Demo Readiness: OpenCode + GSD + Bazel + Kotlin

Status: planning draft; demo not yet proven
Date: 2026-05-05
Parent: `10-agentic-sdlc-v0-design.md`

This document defines what must be true before claiming demo readiness.
It is not a record of a successful demo run.

## Demo Objective

Show that `sdp-trace` can be placed around a real agentic development
workflow without forcing harness replacement:

```text
developer command -> GSD/OpenCode -> Kotlin+Bazel workspace
```

The demo must answer:

- What did `sdp-trace` observe?
- What did it not observe?
- What trust scope does the trace have?
- Can a reviewer replay the evidence path?
- What changes when CI witnesses the verifier result?

## Required Demo Scenes

### Scene 1: Local Wrapper

Command shape:

```bash
sdp-trace wrap --name opencode-gsd -- gsd opencode <task>
```

Acceptable fallback if GSD is unavailable:

```bash
sdp-trace wrap --name opencode -- opencode run <args>
```

If either tool is unavailable, the script must record `not_assessed` and
show a fixture mode instead of pretending the live integration ran.

Fixture replay mode is required for reviewers who do not have OpenCode,
GSD, Bazel, and Kotlin tooling installed. Fixture replay demonstrates
verifier/query/export behavior only; it is not evidence of a fresh live
agent run.

Expected output:

```text
Trust scope: local_observed
Verdict: observed
Completeness: partial
Gate usable: false
```

### Scene 2: Missing Telemetry

Run with no harness adapter and no LLM gateway observer.

Expected MissingEvidenceTable rows:

- harness adapter unsupported or absent;
- model identity not observed;
- tool calls not observed;
- gateway not observed;
- CI witness missing.

The point is not to hide gaps. The point is to make gaps impossible to
miss.

### Scene 3: Bazel Test Evidence

Run an operator-provided test command against the Kotlin fixture:

```bash
bazel test <target>
```

The command must be explicit and scoped. The demo must not execute a
model-suggested shell command automatically.

Expected event:

```text
test_observed
```

If Bazel or Bazelisk is unavailable:

```text
test_observed: not_assessed
reason: bazel_unavailable
```

### Scene 4: Tamper

Mutate, delete, or reorder an event in the captured run.

Expected output:

```text
Verdict: fail
Reason: event_hash mismatch or chain link mismatch
Integrity audit: written outside corrupted chain
```

### Scene 5: CI Witness Boundary

Run verifier in a CI-like context and sign or fixture-witness the tuple:

- source digest;
- contract digest;
- run id;
- chain head;
- verifier version;
- verifier result;
- witness identity;
- timestamp;
- witness independence.

Expected output with real independent CI/OIDC:

```text
Trust scope: ci_witnessed
Gate usable: true
```

Expected output with demo local signing only:

```text
Trust scope: local_observed
Gate usable: false
Reason: witness profile is demo/local or not independent
```

### Scene 6: Forensic Query

Run:

```bash
sdp-trace query <run-dir> --query timeline
sdp-trace query <run-dir> --query missing-evidence
sdp-trace explain <run-dir>
sdp-trace export <run-dir> --format audit-bundle
```

The audience must be able to inspect the run without reading raw event
JSON.

## Demo Assets

Required before live demo claim:

- committed Kotlin+Bazel fixture or external demo repo reference;
- demo expected evidence contract;
- generated local run artifact;
- generated tamper artifact;
- generated audit bundle;
- version report for available external tools;
- explicit `not_assessed` list for unavailable tools.

Candidate existing assets:

- `examples/pilot-fixtures/kotlin-bazel-service/`
- `examples/pilot-runs/opencode-minimax-kotlin-bazel/`
- `scripts/run-opencode-minimax-kotlin-bazel-proof.sh`

These assets may be reused, but Block 10 must not inherit Block 06
claims without re-verifying the exact current demo flow.

## Pre-Demo Checklist

- `sdp-trace wrap` preserves child exit code.
- Default redaction stores no raw secrets.
- Local trace verifies as local-only.
- Missing telemetry table is visible.
- Tamper fixture fails.
- Bazel target is scoped and operator-provided.
- CI witness path is demonstrated or explicitly marked `not_assessed`.
- Audit bundle exports successfully.
- Demo script prints exact external tool versions when observed.
- Demo report states what is proven and what remains unassessed.

## Demo Failure Rules

The demo must fail or downgrade, not improvise, when:

- OpenCode is unavailable;
- GSD is unavailable;
- Bazel/Bazelisk is unavailable;
- the Kotlin target cannot be resolved;
- the expected evidence contract is missing;
- event chain verification fails;
- witness tuple does not match source or contract;
- redaction cannot be applied before write.

Every such state should be shown as product behavior. Missing evidence is
part of the product, not a demo embarrassment.

Tamper mode must mutate a copy of a run artifact. It must not destroy the
canonical local demo output needed by other scenes.
