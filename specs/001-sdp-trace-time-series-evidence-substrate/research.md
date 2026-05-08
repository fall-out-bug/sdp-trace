# Research: sdp-trace Time-Series Evidence Substrate

## Decision 1: Repository-Visible Planning Artifacts Are Canonical

**Decision**: Use `specs/001-sdp-trace-time-series-evidence-substrate/` as the repository-visible source of truth for this feature. External teams may use SpecKit, gsd, Superpowers, Oh My OpenAgent, ticket trackers, or another planning flow as long as the retained trace exposes spec, plan, task, evidence, provenance, and decision inputs.

**Rationale**: A repository observer must be able to inspect scope, plan, tasks, and evidence without loading a private planning runtime.

**Alternatives Rejected**:

- Private-runtime-only planning: rejected because it hides the plan from plain repository review.
- Loose docs under retired research artifacts: rejected because they do not preserve a stable `spec -> plan -> task -> evidence` flow.

## Decision 2: Preserve Origin, Not Dependency

**Decision**: Record that `sdp-trace` originated from delivery evidence work in
`sdp_lab`, but keep all public contracts independent from that repository.

Rejected source ideas:

| Rejected idea | Reason |
|---|---|
| Gate pass/fail/readiness policies | Belongs in external policy consumer, not `sdp-trace`. |
| Traffic-light thresholds from `sdp metrics` | Thresholds are policy. `sdp-trace` records samples and evidence. |
| Beads IDs as required trace fields | Beads is not portable and not a product dependency. |
| Operator Mode phases as required public contract | Runtime-specific and not portable. |
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
- Built-in degradation verdict: rejected because external policy consumer owns policy.

## Decision 4: Close Kotlin+Bazel Evidence Gap

**Decision**: Treat Kotlin+Bazel as a first-class pilot gap. Existing Kotlin+Gradle and Java+Bazel examples are not enough.

**Rationale**: The pilot requirement names JVM stack, Kotlin, and Bazel explicitly. A weak agent can misclassify Bazel repos when Maven metadata exists.

**Evidence from current repo**:

- `docs/jvm-bazel-guide.md` states that Bazel markers must prevent Maven/Gradle inference.
- retired research artifact records a Java-heavy Bazel baseline, not Kotlin+Bazel.
- retired research artifact records Kotlin as Gradle and Bazel as Java/Starlark/Python.

## Decision 5: Compatibility Claims Require Committed Evidence

**Decision**: Harness/model/stack matrices must remain `TBD` or `not_assessed` until backed by committed examples or sanitized run summaries.

**Rationale**: The repository quality bar forbids unsupported claims. This is especially important for OpenCode+MiniMax/Kimi/GLM and less-known harnesses.

## Decision 6: Pin JSON Schema Draft and Validator Before Schema Work

**Decision**: New schemas target JSON Schema Draft 2020-12. Current active validation is Go-first: `go test ./...`, `jq empty schema/*.json`, `go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc`, and `git diff --check`.

**Rationale**: Existing schemas already declare Draft 2020-12. Selecting the draft and validator before authoring new schemas prevents incompatible schema features, invalid examples, and rework across external policy consumers.

**Alternatives Rejected**:

- Keep validator selection late: rejected because schema authors would not know which draft and features are allowed.
- Use `jq` only: rejected because `jq empty` checks JSON syntax, not verifier behavior or fixture conformance.
- Choose a harness-specific runtime validator: rejected because validation tooling must not imply dependency on any AI harness runtime.

## Decision 7: Record External Assertions Instead of Native Strength

**Decision**: `sdp-trace` does not assign evidence strength, quality scores, readiness, degradation, pass/fail, or override semantics. If another producer emits such a value, `sdp-trace` records it as an external verdict input or external assertion with producer, origin, policy reference when available, artifact reference, and provenance.

**Rationale**: Evidence strength is policy-adjacent. Recording it as a native trace judgment would erode the `sdp-trace` / external policy consumer boundary.

## Decision 8: Use Digest-And-Redaction Integrity for Committed Evidence

**Decision**: Committed evidence references use sanitized summaries, SHA-256 digests when artifacts are committed, redaction notes when raw content is withheld, and `integrity_status` for unverified external references.

**Rationale**: A self-trace or pilot example must be inspectable without leaking secrets, credentials, customer data, or raw prompts. SHA-256 digests provide continuity for committed artifacts but are not authentication signatures. Signing and write authorization are external policy unless a future schema version adds them.

## Decision 9: Version Schemas Before External Consumers Adopt Them

**Decision**: Every new schema gets a stable `$id` and semver schema version. External consumers must declare supported versions.

**Rationale**: silent breaking changes would break policy inputs. Before v1.0, breaking changes are allowed only when examples and compatibility notes are updated with the same change.
