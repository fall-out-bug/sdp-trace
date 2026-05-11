# Change Evidence Packet v0

`Change Evidence Packet v0` is the first product artifact for `sdp-trace`.
It organizes retained evidence for one change without turning missing evidence
into trust.

The canonical packet is Markdown rendered from a structured bundle:

```text
sdp-trace packet build-github --github-input github-pr-evidence-input.json --out packet-bundle.json
sdp-trace packet validate --bundle examples/change-evidence-packet/happy-path.bundle.json
sdp-trace packet render --bundle examples/change-evidence-packet/happy-path.bundle.json --out change-evidence-packet.md
```

Every packet has these rows:

- `PC-CHANGE`
- `PC-INITIATOR`
- `PC-AGENT-ROUTE`
- `PC-MUTATION`
- `PC-VERIFICATION`
- `PC-REVIEW`
- `PC-AUTHORITY`
- `PC-THEATER`
- `PC-ATTESTATION`
- `PC-DECISION`
- `PC-RESIDUAL-GAPS`

Allowed row states are `pass`, `partial`, `fail`, `cannot_verify`,
`not_assessed`, and `not_in_scope`. A `pass` row needs retained evidence refs.
Rows with `partial`, `fail`, `cannot_verify`, `not_assessed`, or
`not_in_scope` need a concrete reason.

Evidence refs must resolve through the evidence bundle manifest. Expired or
unresolvable artifacts cannot support `pass`. Contradictory evidence keeps the
affected row `partial` and requires a residual gap.

The packet does not approve merge, release, compliance, production trust,
semantic correctness, or signed external trust. `PC-DECISION` names the next
human owner; it is not an approval verdict.
