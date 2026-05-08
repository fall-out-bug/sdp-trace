# Block 24 Demo Repository CI Evidence Guide

Authority scope: `demo_pilot_only`

This guide explains how to read the Block 24 demo repository evidence. The demo
repository is `fall-out-bug/sdp-trace-demo-ci-pilot`. It is private, so the
durable record in this repository is the sanitized Block 24 report and artifact
index:

- `docs/research/block-24-demo-repo-ci-trace-pilot-report.md`
- `docs/research/block-24-demo-repo-ci-artifact-index.md`

## Integration Shape

The demo uses a small Kotlin/Bazel service surface. Its GitHub Actions workflow
keeps the existing repository command shape and adds `sdp-trace` around selected
Bazel commands. The workflow then produces trace, report, gate, and witness
artifacts.

The clean CI run captured three scoped commands:

- `bazel test //app:feature_flag_test`
- `bazel test //app:entitlement_matrix_test`
- `bazel test //app:audit_scope_test`

The same run also exercised intentionally incomplete or dishonest paths:

- a CI job without OIDC witness permissions;
- a copied trace/report artifact with source and run binding mismatch;
- a stale digest index.

This mix is deliberate. A useful pilot must show both observed evidence and the
places where the substrate refuses to upgrade trust.

## Evidence Path

Start with the artifact index. It records the demo repository, source commit,
workflow run, uploaded artifact ids, artifact retention, selected digests, and
the safe witness fields copied from GitHub Actions.

Then read the pilot report. It explains how the same CI run produced:

- captured run directories for the three clean Bazel commands;
- `summary.json` with three observed run summaries;
- `verify.txt` and `explain.txt` for reviewer-readable local trace inspection;
- `gate-result.json` with verifier facts rather than readiness decisions;
- `ci-witness.json` showing CI identity evidence for the exact demo topology;
- `ci-witness-no-oidc.json` showing `cannot_verify` when identity evidence is
  missing.

Raw logs, raw OIDC material, provider credentials, and authenticated artifact
URLs are not copied into this repository. The committed record keeps sanitized
summaries, selected digests, reason codes, and access-neutral references.

## State Semantics

`observed` means the selected run or profile had enough evidence for the scoped
local verifier fact.

`ci_witnessed` means the witness profile could bind the selected report or run
artifacts to GitHub Actions identity for that workflow run. It does not create
production trust by itself.

`cannot_verify` means the verifier attempted the check but lacked required
evidence or consistent binding. The no-OIDC case demonstrates this: CI executed
the command, but the job did not expose the identity material needed to raise the
witness state.

`fail` means the evidence conflicts with the selected profile. The stale digest
case demonstrates this: the referenced bytes and digest index do not match.

The demo keeps these states separate. It does not collapse missing telemetry,
stale artifacts, and local observation into one success label.

## Harness Boundary

The pilot path is sidecar-first:

1. Keep the team's existing harness and CI command.
2. Wrap selected commands or emit adapter events where the harness can do so.
3. Publish sanitized run/report artifacts from CI.
4. Run witness or assessment profiles only when the required evidence exists.
5. Feed the resulting facts to an external policy layer when the organization is
   ready to define enforcement.

This lets a team inspect process evidence before replacing its harness. If the
team already has a default harness, `sdp-trace` can start as the evidence and
gate-fact layer around that harness instead of competing with it.

## Kotlin/Bazel Upgrade Path

The Block 24 demo used shell-based Bazel tests over Kotlin source and repository
metadata. That is enough for CI trace mechanics, but it does not assess compiled
Kotlin/JVM behavior.

To assess a compiled JVM target in a later pilot, keep the same sidecar pattern
and replace the scoped command with the repository's real compiled target, for
example:

```bash
go run ./cmd/sdp-trace wrap --name compiled-jvm-test -- bazel test //app:service_jvm_test
```

The target can be a `kt_jvm_test`, `java_test`, or the organization's existing
Bazel target shape. The evidence package then needs to retain the target label,
BUILD rule or sanitized rule excerpt, Bazel version, command result, stdout and
stderr digests, source commit, and CI witness material where available. If
compiled artifacts are part of the claim, retain their digests or an accepted
external artifact reference.

Compiled target evidence can raise the compatibility claim for that selected
target. It does not raise `demo_pilot_only` to production trust, and it does not
prove broad JVM/Bazel compatibility without a wider assessed target selection.

## Demonstrated Behavior

The demo demonstrates that a separate repository can run ordinary CI commands,
attach `sdp-trace`, retain sanitized trace artifacts, show verifier states, and
make missing or inconsistent evidence visible.

It also demonstrates that a passing CI job is not automatically trusted evidence.
The OIDC, source-binding, and digest cases show the difference between command
execution, local observation, CI-witnessed evidence, and rejected or incomplete
trust.

## Limits

The demo scope is intentionally narrow:

- it uses GitHub Actions, not every CI system;
- it runs under the same owner namespace, not a customer-owned organization;
- it builds `sdp-trace` from source, not from a released binary;
- its Kotlin/Bazel tests inspect scoped source and metadata, but do not prove
  full compiled JVM compatibility;
- it does not establish external production trust.

These limits are the pilot boundary. A later pilot can replace the demo
repository with a target repository and retain the same evidence discipline:
observed states remain scoped, missing evidence stays visible, and production
trust is not claimed until the required external trust profile is actually
assessed.
