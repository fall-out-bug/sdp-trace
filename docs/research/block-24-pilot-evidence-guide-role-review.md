# Block 24 Pilot Evidence Guide Role Review

Authority scope: `review_disposition_only`

This records the role review used to harden the demo-repository evidence path.
It is not an external handoff or readiness claim.

## Review Packet

Reviewed files:

- `README.md`
- `docs/agent-entrypoint.md`
- `docs/reviewer-entrypoint.md`
- `docs/research/block-24-demo-repo-ci-evidence-guide.md`
- `docs/research/block-24-demo-repo-ci-trace-pilot-report.md`
- `docs/research/block-24-demo-repo-ci-artifact-index.md`

Review constraints:

- keep the documentation evidence-led and repository-native;
- avoid prompt-list framing and grid-based coverage framing;
- keep sidecar-first integration explicit;
- do not claim broad production trust or broad compatibility;
- keep `observed`, `ci_witnessed`, `cannot_verify`, `fail`, and
  `not_assessed` distinct.

## Role Results

CTO buyer review returned `APPROVED`. It found no critical or major issues and
one minor README accuracy risk around JVM/Bazel examples. The README now states
that JVM/Bazel coverage is scoped to design fixtures and pilot evidence.

Head of InfoSec review returned `APPROVED`. It found no critical, major, or
minor issues. The review accepted the raw-log, OIDC, private-artifact,
sanitized-evidence, production-trust, compatibility, and dishonest-case
boundaries as explicit.

Head of Engineering review first returned `REVISE`. Valid findings were accepted:

- `ci_witnessed`, `local_observed`, and `local_dirty_structural_only` needed a
  documented distinction between result state, trust scope, and authority scope.
- `demo_pilot_only` needed a definition.
- `not_assessed` with exit `0` needed an explicit pipeline interpretation.
- reviewer and agent entrypoints duplicated the state contract.
- the Kotlin/Bazel pilot needed a compiled-target upgrade path.

The focused Head of Engineering re-review returned `APPROVED`. It found no
remaining critical or major issues. Minor observations were advisory only and did
not require further changes for this slice.

## Disposition

Accepted and fixed:

- `docs/agent-entrypoint.md` now separates result state, trust scope, and
  authority scope.
- `docs/agent-entrypoint.md` now defines `demo_pilot_only` and
  `local_dirty_structural_only`.
- `docs/agent-entrypoint.md` now defines when `not_assessed` may return exit
  `0`, and warns pipeline authors to inspect emitted state fields.
- `docs/reviewer-entrypoint.md` now points to `docs/agent-entrypoint.md` as the
  canonical state, scope, and exit-code contract.
- `docs/research/block-24-demo-repo-ci-evidence-guide.md` now documents a
  compiled Kotlin/JVM target upgrade path without broad compatibility claims.
- `docs/research/block-24-demo-repo-ci-trace-pilot-report.md` no longer carries
  the old customer-question table; it uses evidence-coverage prose.

No role approved production trust, owner independence, non-GitHub portability,
release binary acquisition, or broad JVM/Bazel compatibility.
