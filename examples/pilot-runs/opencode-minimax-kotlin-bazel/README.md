# OpenCode + MiniMax + Kotlin+Bazel Proof Package

Completion state: incomplete

Validate package shape and honest incomplete states:

```bash
scripts/validate-e2e-pilot-package.sh --mode package examples/pilot-runs/opencode-minimax-kotlin-bazel
```

Validate completed proof only after every required proof state is observed:

```bash
scripts/validate-e2e-pilot-package.sh --mode complete examples/pilot-runs/opencode-minimax-kotlin-bazel
```
