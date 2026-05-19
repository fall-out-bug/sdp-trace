# Install

`sdp-trace` is delivered as a single binary. Go is required to build or release
the CLI, but it is not required to run observation around an existing harness.

## Binary Path

Download a release artifact for your platform from the repository's GitHub
releases page, verify the digest from `SHA256SUMS`, put it on `PATH`, and
verify:

```text
sdp-trace version
sdp-trace wrap --name smoke --output-dir .sdp-trace-runs/smoke -- /bin/echo ok
sdp-trace verify .sdp-trace-runs/smoke
sdp-trace explain .sdp-trace-runs/smoke
sdp-trace report --out .sdp-trace-report .sdp-trace-runs/smoke
sdp-trace query --query missing-evidence .sdp-trace-runs/smoke
```

On Windows, replace `/bin/echo ok` with a local command available in your
shell, for example `cmd /c echo ok` in Command Prompt.

The wrapped command is the existing harness command. This core path records,
verifies, explains, reports, and queries missing evidence. `sdp-trace` records
command provenance and retained artifacts outside the prompt surface; it does
not inject instructions into the harness or model context.

Assessment profiles, gate facts, witness artifacts, release proof, and PR
packet proof are extension surfaces. Add them only after the core run/report
path is working and an external policy consumer needs those facts.
For PR packet proof, run `packet build-pr --source github-actions` inside
GitHub Actions with `GITHUB_EVENT_PATH`, `GITHUB_RUN_ID`, repository identity,
and retained artifact evidence available. Curated `packet build-github
--github-input` files are fixture/backfill inputs only.

## Build Path

Maintainers can build release artifacts with:

```text
./scripts/build-release-binaries.sh
```

This writes platform binaries and SHA-256 files to `dist/`.

For a source checkout without installing a binary, follow the
[Contributor Quick Start](contributor-quickstart.md) for the canonical `go run`
smoke path.

To install from the checkout onto `PATH`:

```text
go install ./cmd/sdp-trace
```

The current release matrix is:

- `linux/amd64`
- `linux/arm64`
- `darwin/amd64`
- `darwin/arm64`
- `windows/amd64`
- `windows/arm64`
