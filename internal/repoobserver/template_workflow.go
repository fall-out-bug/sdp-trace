package repoobserver

func githubWorkflow() string {
	// The generated workflow uploads observations; proof is established only
	// after inspecting a real CI artifact.
	// It captures repository metadata and safe status output, not raw secrets or
	// local hook output.
	// The workflow is intentionally small and uses shell commands only as a thin
	// CI metadata launcher.
	// It uploads `.sdp-trace/ci` as an artifact source for later inspection.
	// Repository proof remains not_assessed until that uploaded artifact is
	// inspected by a separate evidence path.
	// The checkout step is the only third-party action used by this template.
	return `name: sdp-trace observe

on:
  pull_request:
  push:

jobs:
  observe:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      actions: read
    steps:
      - uses: actions/checkout@v4
      - name: Capture safe repository metadata
        shell: bash
        run: |
          set -euo pipefail
          mkdir -p .sdp-trace/ci
          {
            printf 'github_repository=%s\n' "$GITHUB_REPOSITORY"
            printf 'github_run_id=%s\n' "$GITHUB_RUN_ID"
            printf 'github_sha=%s\n' "$GITHUB_SHA"
            printf 'github_ref=%s\n' "$GITHUB_REF"
          } > .sdp-trace/ci/metadata.env
          git status --short > .sdp-trace/ci/status.txt
          git diff --check > .sdp-trace/ci/diff-check.txt 2>&1 || true
      - name: Optional Bazel test smoke
        shell: bash
        run: |
          set -euo pipefail
          if [ -f MODULE.bazel ] || [ -f WORKSPACE ] || [ -f WORKSPACE.bazel ]; then
            printf 'bazel_config_present\n' > .sdp-trace/ci/bazel-test.txt
          else
            printf 'bazel_not_configured\n' > .sdp-trace/ci/bazel-test.txt
          fi
      - uses: actions/upload-artifact@v4
        with:
          name: sdp-trace-observer
          path: .sdp-trace/ci/
`
}
