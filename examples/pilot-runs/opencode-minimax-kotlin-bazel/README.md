# OpenCode + MiniMax + Kotlin+Bazel Proof Package

Completion state: complete

Inspect package JSON syntax and current repository verifier baseline:

```bash
jq empty examples/pilot-runs/opencode-minimax-kotlin-bazel/evidence/*.json
jq empty examples/pilot-runs/opencode-minimax-kotlin-bazel/handoff/assessment-input.json
go test ./...
jq empty schema/*.json
go run ./cmd/sdp-trace validate-fixtures examples/agentic-sdlc
git diff --check
```

The retired `scripts/validate-e2e-pilot-package.sh` command is not part of the
current product path. The current Go command surface does not provide a dedicated
validator for this historical pilot package shape; package-level schema
validation is therefore `not_assessed` unless a current Go verifier profile or
fixture validator is added. Do not use this package to claim broad OpenCode,
MiniMax, Kotlin, or Bazel support beyond the exact observed slice.
