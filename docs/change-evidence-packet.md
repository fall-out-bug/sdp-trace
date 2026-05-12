# Change Evidence Packet v0

`Change Evidence Packet v0` is the first product artifact for `sdp-trace`.
It organizes retained evidence for one change without turning missing evidence
into trust.

The canonical packet is Markdown rendered from a structured bundle:

```text
sdp-trace packet build-pr --source github-actions --out packet-artifacts
sdp-trace packet build-github --github-input github-pr-evidence-input.json --out packet-bundle.json
sdp-trace packet validate --bundle examples/change-evidence-packet/happy-path.bundle.json
sdp-trace packet check-demo --bundle examples/change-evidence-packet/happy-path.bundle.json
sdp-trace packet render --bundle examples/change-evidence-packet/happy-path.bundle.json --out change-evidence-packet.md
```

`packet build-pr` is the live-demo path. It derives packet artifacts from
GitHub Actions context or supplied GitHub fixture JSON plus recorder route
metadata, then writes `bundle.json`, `change-evidence-packet.md`, and
`build-pr-result.json` into the output directory. Checked-in packet Markdown or
checked-in bundle JSON cannot override those generated artifacts.

`packet build-github --github-input` remains a fixture/backfill command. It is
not sufficient live-demo proof by itself because the input is curated before the
builder runs.

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

`packet check-demo` is a narrow 007 demo gate. It checks the first-packet
minimum bar for the GitHub OSS demo, including `PC-AGENT-ROUTE: pass|partial`
with retained structured harness route observation evidence
(`evidence_kind: harness_route_observation` and `observed_components` covering
OpenCode, GSD, and MiniMax). Route evidence must be digest-bound with
`agent_route_digest`; self-declared route metadata with only a generated
placeholder digest does not satisfy the demo gate. It is not a general merge,
release, or production-trust gate.
