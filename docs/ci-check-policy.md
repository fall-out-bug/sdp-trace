# CI Check Policy

Pull requests are expected to report GitHub Actions checks from
`.github/workflows/ci.yml`.

Required CI evidence:

- `go test ./...`
- `jq empty schema/*.json examples/block19-adapter-capture/*.json`
- `git diff --check`

If GitHub does not report checks for a PR, record CI as `not_assessed`; do not
treat local verification as a substitute for remote CI evidence. Local
verification may support implementation review, but CI-backed closure requires
the workflow result or an explicit repo-tracked replacement policy.
