# Install

`sdp-trace` is delivered as a single binary. Go is required to build or release
the CLI, but it is not required to run observation around an existing harness.

## Binary Path

Download a release artifact for your platform, verify the digest from
`SHA256SUMS`, put it on `PATH`, and verify:

```text
sdp-trace version
sdp-trace wrap --name smoke --output-dir .sdp-trace-runs/smoke -- /bin/echo ok
sdp-trace report --out .sdp-trace-report .sdp-trace-runs/smoke
sdp-trace packet build-pr --source github-actions --out packet-artifacts
```

The wrapped command is the existing harness command. `sdp-trace` records command
provenance and retained artifacts outside the prompt surface; it does not inject
instructions into the harness or model context.
For PR packet proof, `packet build-pr` is the live GitHub Actions path; curated
`packet build-github --github-input` files are fixture/backfill inputs only.

## Build Path

Maintainers can build release artifacts with:

```text
./scripts/build-release-binaries.sh
```

This writes platform binaries and SHA-256 files to `dist/`.

The current release matrix is:

- `linux/amd64`
- `linux/arm64`
- `darwin/amd64`
- `darwin/arm64`
- `windows/amd64`
- `windows/arm64`
