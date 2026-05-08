# Block 25: Compiled CI Demo Pilot

Status: spec draft, pending Socratic review.

Parent artifacts:

- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/23-mvp-closure-drift-and-readiness.md`
- `specs/001-sdp-trace-time-series-evidence-substrate/blocks/24-demo-repo-ci-trace-pilot.md`
- retired research artifact
- retired research artifact

## Goal

Raise the demo repository from a narrow CI trace mechanics pilot to a
demonstrable technical pilot scoped to a selected compiled JVM/Bazel behavior
path and a CI-produced artifact evidence loop.

The intended product signal is not broader production trust. The signal is that
an existing repository can keep its CI and harness shape, run one selected
compiled target, attach `sdp-trace`, publish sanitized artifacts from CI, and
make incomplete or dishonest evidence visible without relying on checked-in
proof JSON. One passing selected target does not imply that other targets,
other BUILD packages, or the full monorepo surface have been assessed.

## Problem

Block 24 proved that a separate repository can run CI commands, attach
`sdp-trace`, upload artifacts, and expose `observed`, `ci_witnessed`,
`cannot_verify`, and `fail` style outcomes. It deliberately stayed narrow.

The demo repository review found two gaps that prevent a stronger buyer pilot:

- the primary Bazel tests are shell `grep` checks over Kotlin source, not
  compiled Kotlin/JVM behavior tests;
- the generated artifact index was written inside the directory it indexed,
  creating a self-referential hash entry that cannot verify the final index
  file.

Those gaps do not invalidate Block 24's narrow authority scope, but they mean
the demo is not yet strong enough as a pilot demonstration for a
Kotlin/Java+Bazel organization that needs behavior evidence.

## Non-Goals

- No claim of customer production trust.
- No claim of broad JVM, Kotlin, Bazel, monorepo, or non-GitHub CI
  compatibility.
- No native policy decision, release approval, risk acceptance, audit verdict,
  readiness score, health score, badge, grade, or rank.
- No dependency from `sdp-trace` product code to the demo repository.
- No raw OIDC token, CI secret, private artifact URL, raw log body, customer
  data, private filesystem path, or model prompt/response body in committed
  docs, committed artifacts, or uploaded CI artifacts.
- No checked-in `.sdp-trace-report/`, `.sdp-trace-runs/`, or built
  `.sdp-trace-tools/` output as authority.
- No use of retired Node/npm/script validators as current proof.
- No questionnaire, sales deck, or coverage-matrix framing.

## Product Boundary

Block 25 may change two repositories:

- `sdp-trace-demo-ci-pilot`: source, Bazel config, workflow, thin helper
  scripts, README, and generated CI artifacts.
- `sdp-trace`: SpecKit docs and sanitized research evidence that point to the
  successful demo run after it exists.

The demo repository remains a separate pilot repository. `sdp-trace` records the
sanitized evidence and interpretation, but does not vendor the demo app, demo CI
runtime, OpenCode, GSD, MiniMax, Bazel rules, or GitHub Actions workflow as
product dependencies.

## Required Demo Behavior

The demo repository must contain a selected compiled Kotlin/JVM behavior path:

- a small Kotlin service surface with behavior that can be tested semantically;
- one `kt_jvm_library` target and one `kt_jvm_test` target using a pinned
  `rules_kotlin` version or commit recorded in `MODULE.bazel`;
- JVM toolchain pinned through Bazel toolchain configuration or an explicit JDK
  setup step, with the version recorded in the demo repository README and CI
  logs;
- test framework recorded in the demo repository README;
- at least one CI command wrapped by `sdp-trace` that runs the compiled target;
- optional source/scope checks may remain only as secondary metadata checks.

The compiled test must assert behavior, not just source-string presence. For the
Feature Flag / Entitlements demo, the domain is a synthetic entitlements service
surface. Acceptable assertions include plan-specific feature enablement,
audit-log entitlement behavior, and seat-overage warning behavior.

One passing selected target does not imply other targets, other BUILD packages,
or the full monorepo surface. Multi-target assessment is outside this demo scope
and requires separate evidence.

## CI Artifact Contract

Artifacts must be produced by GitHub Actions from a clean checkout. Local
generated artifacts are allowed only for debugging and must remain ignored.

The clean CI path must produce:

- wrapped run directories for the selected compiled target and any secondary
  scope checks;
- `verify`, `explain`, `report`, `gate`, and `witness` outputs;
- an artifact index whose listed digests can be recomputed after download;
- redaction scan output;
- a CI run id, source commit, selected `sdp-trace` source commit or immutable
  ref, and artifact retention window.

The artifact index must not index itself. The preferred implementation is:

1. write the index to a temporary file outside the indexed tree;
2. hash all indexed files except the final index file using SHA-256 lowercase
   hex;
3. move the completed index into the artifact tree;
4. run a verification step that recomputes every listed digest from the final
   artifact directory and asserts that the index path is absent from its own
   listed entries.

The index format is JSON:

```json
{
  "schema_version": "demo-artifact-index-v1",
  "authority_scope": "demo_pilot_only",
  "entries": [
    {
      "path": "relative/path/from/artifact/root",
      "sha256": "lowercase-hex-sha256",
      "size_bytes": 123
    }
  ]
}
```

Entries must be sorted by relative path in deterministic lexicographic order.
Each `path` must appear exactly once in `entries`. Paths must be relative to the
artifact root, use `/`, and must not be absolute, empty, parent-traversing,
URL-like, or equal to the final index path. The index enumerates all regular
files in the artifact root recursively; `.git`, `.github`, and other VCS or host
metadata directories are excluded from index scope. The verification command
must enumerate the final artifact directory recursively, exclude the index file,
compute `sha256` and `size_bytes` for every remaining file, compare the set
exactly against the index entries, and exit `0` only when there are no extra,
missing, duplicate, or mismatched files.

The index file's own integrity is not established by self-indexing. It is
established by CI artifact service metadata when available, downloaded archive
digest when available, and the fact that every other file in the extracted
artifact tree matches the final index.

If the index verifier is not implemented yet, Block 25 remains open with
`cannot_verify`.

Local `bazel test //...` is required for local build confidence. Local
`sdp-trace` runs may produce ignored `.sdp-trace-*` directories for developer
debugging, but local artifact digests must not be cited in the Block 25 report,
sanitized report, review ledger, or evidence claim. Only artifacts generated by
an explicitly numbered GitHub Actions run id on the selected source commit serve
as Block 25 proof.

Artifact retention must be explicit. The demo README and sanitized report must
state the configured artifact retention duration, what the buyer must preserve
for longer audit windows, and that `sdp-trace` does not own CI artifact
retention.

The demo must state whether CI fetches `sdp-trace` source or uses a released
binary. If it fetches source, it must record the source checkout mechanism and
selected source commit or immutable ref. Released binary acquisition and
installation UX remain `not_assessed` unless separately evaluated.

Redaction scan scope covers every downloaded clean and negative artifact set.
The first implementation may reuse the Block 24 denylist pattern file and
portable shell script, but the report must record the exact command, pattern
file digest, scanned artifact roots, exit code, and pass/fail/cannot-verify
state. Pattern classes must include token-like strings, private key material,
provider credentials, private artifact URLs, raw OIDC material, and local
private paths.

The demo CI job configuration should avoid echoing tokens, private URLs, OIDC
material, and private filesystem paths. CI job logs are an external retention
surface and are not part of the uploaded CI artifact contract.

## Dishonest And Incomplete Cases

The demo must keep negative cases because a credible pilot needs to show refusal
behavior, not only green paths.

Required cases:

- no-OIDC witness gap: a separate CI job omits `id-token: write` or runs the
  witness command outside an OIDC-capable context. Expected state:
  `cannot_verify` with an OIDC-unavailable or identity-missing reason.
- stale digest: generate a clean index, then intentionally mutate one
  non-index artifact file while leaving the index entry unchanged. Expected
  state: `fail` with a digest-mismatch tamper-fixture reason.
- source/run binding absent or unresolvable: expected state `cannot_verify`.
- source/run binding present but contradictory, such as commit, digest, or run
  id mismatch: expected state `fail`.

The stale digest case must be independent of the artifact-index implementation
bug fixed in Block 25. It must demonstrate an intentional tamper/staleness
fixture, not a real bug in the clean artifact index.
The fixture must not modify the index file itself.

## Documentation Requirements

The demo repository README must state:

- authority scope is `demo_pilot_only`;
- primary proof path is compiled JVM behavior for selected targets only;
- any grep/source checks are secondary scope checks and do not prove runtime
  behavior;
- one passing selected target does not imply other targets, BUILD packages, or
  the full monorepo surface;
- CI prerequisites, including any read token needed to fetch `sdp-trace` source;
- whether the CI run fetches `sdp-trace` source or uses a released binary, plus
  the exact source commit/ref or artifact version used by the successful run;
- artifacts are generated by CI and expire after the configured retention
  duration unless separately archived;
- production trust, owner independence, non-GitHub portability, release binary
  acquisition, and broad JVM/Bazel compatibility remain outside the demo unless
  separately assessed.

The `sdp-trace` repository may add a sanitized Block 25 report only after the
demo CI run exists and downloaded artifacts have been verified. That report must
be Markdown under retired research artifacts; it must not be a schema-validated fixture
under `examples/`, and it must not contain a buyer questionnaire or sales
narrative.

## Acceptance Criteria

Block 25 is not complete until all of these are true:

1. Demo repository `bazel test //...` passes locally or the local blocker is
   recorded as `cannot_verify` with exact reason.
2. Demo repository GitHub Actions run passes on the selected source commit.
3. Downloaded CI artifacts include a compiled JVM wrapped run and replayable
   `sdp-trace` report outputs.
4. Artifact index verification recomputes every listed digest successfully and
   explicitly asserts that the index does not self-index.
5. No-OIDC, stale digest, and source/run mismatch cases are present and carry
   distinct `cannot_verify` or `fail` states.
6. Redaction scan passes over downloaded artifact sets.
7. Demo README and `sdp-trace` sanitized report keep all residual trust states
   explicit as `not_assessed` or `cannot_verify`.
8. Independent role reviews for CTO buyer, Head of Engineering, and Head of
   InfoSec find no remaining critical or major issues after fixes.
9. `sdp-trace` PR-level review and CI pass after sanitized Block 25 docs are
   added.

## Review Plan

Socratic spec review must run before implementation approval with at least these
planes:

- product/buyer credibility: whether the compiled demo answers the pilot need
  without becoming sales material;
- engineering/replayability: Bazel target shape, CI-only artifact generation,
  artifact-index verification, and local/CI reproduction;
- tracing/evidence: state/scope semantics, negative cases, and report-to-demo
  consistency;
- security/privacy: OIDC/token handling, artifact exposure, redaction, and
  unsafe output boundaries.

Implementation review must repeat across at least:

- demo code and CI correctness;
- tracing/evidence and artifact integrity;
- requirements-vs-implementation for buyer pilot readiness;
- security/privacy for artifact, token, and redaction boundaries.

Any hung, empty, or off-task review is `not_assessed` and must be replaced.
Closure requires usable CTO buyer, Head of Engineering, and Head of InfoSec
review outputs; a hung role is not enough for closure.

## Residual States

Even if Block 25 passes, these remain outside scope unless separately assessed:

- customer production trust;
- customer-owned organization or private network topology;
- non-GitHub CI portability;
- released binary acquisition and installation UX;
- broad JVM/Kotlin/Bazel compatibility;
- broad monorepo scalability;
- CI credential shape, source-fetch token policy, OIDC token access scope, and
  artifact upload/download authentication;
- OpenCode, GSD, or MiniMax as product dependencies.
