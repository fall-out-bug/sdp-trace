# Block 26 CI Artifact Observation Fixtures

These fixtures document the first implementation matrix for the CI artifact
observation profile.

The profile records facts about uploaded CI artifact coverage. It does not
enforce merge, release, readiness, audit, or risk decisions.

Fixture scenarios are listed in `fixture-matrix.json`. Implementation tests
exercise the same state and reason-code matrix through generated manifests.

Run the complete-coverage fixture:

```bash
go run ./cmd/sdp-trace assess \
  --profile ci-artifact-observation \
  --artifact-manifest examples/block26-ci-artifact-observation/input/ci-uploaded-bundle-complete-coverage/artifact-manifest.json \
  --out /tmp/block26-observation.json
```

Explain the result:

```bash
go run ./cmd/sdp-trace assess explain \
  --assessment-result /tmp/block26-observation.json
```

`fixture-matrix.json` is the authoritative scenario matrix for the first
implementation. The committed matrix names expected top-level states and primary
reason codes for the negative cases. The Go fixture-matrix test evaluates those
scenarios from the same matrix to keep replayability tied to product behavior.
