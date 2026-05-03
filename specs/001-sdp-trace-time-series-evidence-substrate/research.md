# Research: sdp-trace Time-Series Evidence Substrate

## Decision 1: SpecKit Artifacts Are Canonical

**Decision**: Use `specs/001-sdp-trace-time-series-evidence-substrate/` as the repository-visible source of truth. Beads issues mirror the work graph but do not replace SpecKit artifacts.

**Rationale**: A repository observer must be able to inspect scope, plan, tasks, and evidence without loading a local Beads database.

**Alternatives Rejected**:

- Beads-only planning: rejected because it hides the plan from plain repository review.
- Loose docs under `docs/research/`: rejected because it does not preserve the SpecKit `spec -> plan -> tasks` flow.

## Decision 2: Import Evidence Concepts, Not Gate Policy, From sdp_lab

**Decision**: Reuse the best portable ideas from `sdp_lab` as source material, but translate them into `sdp-trace` contracts.

Useful source ideas:

| Source | Portable idea |
|---|---|
| `/Users/fall_out_bug/projects/vibe_coding/sdp_lab/schema/evidence-envelope.schema.json` | Strict evidence sections, provenance, boundary, trace fields |
| `/Users/fall_out_bug/projects/vibe_coding/sdp_lab/schema/sdp-pr-gate/evidence-event.schema.json` | Event-level evidence ingestion with source, external ref, timestamps, actor, status, artifact URI/hash |
| `/Users/fall_out_bug/projects/vibe_coding/sdp_lab/docs/PROCESS_HYGIENE_TELEMETRY_SPEC.md` | Process observations such as stale work, blocked debt, missing trace/evidence links |
| `/Users/fall_out_bug/projects/vibe_coding/sdp_lab/internal/metrics/types.go` | Process metric categories: hygiene, waste, git flow, release quality, stabilization, knowledge risk, decay |
| `/Users/fall_out_bug/projects/vibe_coding/sdp_lab/docs/REVIEW_AND_DELIVERY_TRACE_WORKING_MODEL.md` | Trace continuation through review, delivery, rollback, and follow-up |
| `/Users/fall_out_bug/projects/vibe_coding/sdp_lab/schema/telemetry/sdp-trace-events.schema.json` | Attribute allowlist, consent levels, and metadata-first telemetry discipline |
| `/Users/fall_out_bug/projects/vibe_coding/sdp_lab/docs/plans/2026-04-27-f150-product-layering-release-readiness-design.md` | Product layering: substrate vs policy product vs operator runtime |

Rejected source ideas:

| Rejected idea | Reason |
|---|---|
| Gate pass/fail/readiness policies | Belongs in `sdp-gate`, not `sdp-trace`. |
| Traffic-light thresholds from `sdp metrics` | Thresholds are policy. `sdp-trace` records samples and evidence. |
| Beads IDs as required trace fields | Beads is not portable and not a product dependency. |
| Operator Mode phases as required public contract | Runtime-specific and too coupled to `sdp_lab`. |
| PR-only passport assumptions | `sdp-trace` must also support local branches, commits, files, commands, and non-GitHub systems. |

## Decision 3: Use Moving Windows Instead of Fixed Baseline Ownership

**Decision**: `sdp-trace` records ordered metric streams and trace snapshots. Consumers can compare current windows with previous windows, but `sdp-trace` does not own threshold policy or degradation verdicts.

**Rationale**: The user clarified that the baseline source is not important. The process moves over time; we need evidence-backed observations of movement.

Movement is represented structurally:

- previous window reference
- current window reference
- previous value
- current value
- delta
- unit
- dimensions
- evidence coverage and `not_assessed` gaps

The words degrading, improving, pass, fail, ready, and blocked are not native movement labels in `sdp-trace`.

**Alternatives Rejected**:

- Fixed baseline entity: rejected because it overfits to one evaluation method.
- Built-in degradation verdict: rejected because `sdp-gate` owns policy.

## Decision 4: Close Kotlin+Bazel Evidence Gap

**Decision**: Treat Kotlin+Bazel as a first-class pilot gap. Existing Kotlin+Gradle and Java+Bazel examples are not enough.

**Rationale**: The pilot requirement names JVM stack, Kotlin, and Bazel explicitly. A weak agent can misclassify Bazel repos when Maven metadata exists.

**Evidence from current repo**:

- `docs/jvm-bazel-guide.md` states that Bazel markers must prevent Maven/Gradle inference.
- `docs/research/phase-a-bazel-codex-summary.md` records a Java-heavy Bazel baseline, not Kotlin+Bazel.
- `docs/research/bootstrap-smoke-summary.md` records Kotlin as Gradle and Bazel as Java/Starlark/Python.

## Decision 5: Compatibility Claims Require Committed Evidence

**Decision**: Harness/model/stack matrices must remain `TBD` or `not_assessed` until backed by committed examples or sanitized run summaries.

**Rationale**: The repository quality bar forbids unsupported claims. This is especially important for OpenCode+MiniMax/Kimi/GLM and less-known harnesses.

## Decision 6: Pin JSON Schema Draft and Validator Before Schema Work

**Decision**: New schemas target JSON Schema Draft 2020-12. The full-validation command uses pinned local `ajv@8.20.0` through `scripts/validate-json-schema.mjs`.

**Rationale**: Existing schemas already declare Draft 2020-12. Selecting the draft and validator before authoring new schemas prevents incompatible schema features, invalid examples, and rework across `sdp-gate` consumers.

**Alternatives Rejected**:

- Keep validator selection late: rejected because schema authors would not know which draft and features are allowed.
- Use `jq` only: rejected because `jq empty` checks JSON syntax, not schema conformance.
- Choose a harness-specific runtime validator: rejected because validation tooling must not imply dependency on any AI harness runtime.

## Decision 7: Record External Assertions Instead of Native Strength

**Decision**: `sdp-trace` does not assign evidence strength, quality scores, readiness, degradation, pass/fail, or override semantics. If another producer emits such a value, `sdp-trace` records it as an external verdict input or external assertion with producer, origin, policy reference when available, artifact reference, and provenance.

**Rationale**: Evidence strength is policy-adjacent. Recording it as a native trace judgment would erode the `sdp-trace` / `sdp-gate` boundary.

## Decision 8: Use Digest-And-Redaction Integrity for Committed Evidence

**Decision**: Committed evidence references use sanitized summaries, SHA-256 digests when artifacts are committed, redaction notes when raw content is withheld, and `integrity_status` for unverified external references.

**Rationale**: A self-trace or pilot example must be inspectable without leaking secrets, credentials, customer data, or raw prompts. SHA-256 digests provide continuity for committed artifacts but are not authentication signatures. Signing and write authorization are external policy unless a future schema version adds them.

## Decision 9: Version Schemas Before sdp-gate Inherits Them

**Decision**: Every new schema gets a stable `$id` and semver schema version. `sdp-gate` or any other consumer must declare supported versions.

**Rationale**: `sdp-gate` builds on `sdp-trace`; silent breaking changes would break inherited policy inputs. Before v1.0, breaking changes are allowed only when examples and compatibility notes are updated with the same change.
