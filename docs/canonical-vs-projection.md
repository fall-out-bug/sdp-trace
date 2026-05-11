# Canonical Packet vs Projection

For Change Evidence Packet v0, the canonical artifact is the generated Markdown
packet plus its evidence bundle manifest.

A GitHub PR comment, PR body section, dashboard card, or CLI summary is a
projection. A projection is useful for navigation, but it must not be marked as
canonical over the uploaded packet artifact.

Validator rule:

- `projection.kind: canonical_markdown_artifact` with `canonical: true` is
  canonical.
- Any other projection with `canonical: true` is rejected.
- A non-canonical projection must point at the canonical packet artifact through
  `artifact_ref`.

This preserves the product boundary: prose shown in a change host is not proof
unless it is bound back to retained packet and bundle evidence.
