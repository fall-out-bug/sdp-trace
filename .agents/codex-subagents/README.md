# Codex Subagent Configuration for sdp-trace

## Roles

Custom role cards for sdp-trace work:

- `roles/sdp-trace-slice-worker.json` — Implementation agent for renaming and merging files within a command family slice.
- `roles/sdp-trace-reviewer.json` — Adversarial review agent (model policy: non-OpenAI/Anthropic/Google preferred).
- `roles/sdp-trace-verifier.json` — Verification agent that runs the full quality gate suite.

## Contexts

- `contexts/sdp-trace-repo-rules.json` — Repository trust rules, quality gates, verification commands, and model policy.

## Usage Examples

### Run a slice worker
```bash
codex-subagent run pi \
  --role-card .agents/codex-subagents/roles/sdp-trace-slice-worker.json \
  --context-pack .agents/codex-subagents/contexts/sdp-trace-repo-rules.json \
  --task "Implement slice for WITNESS family: rename main_374*.go → witness_*.go, merge tiny files, verify" \
  --cwd /home/fall_out_bug/projects/vibe_coding/sdp-trace/.worktrees/010-command-package-organization \
  --model zai/glm-5.1
```

### Run adversarial review
```bash
codex-subagent run pi \
  --role-card .agents/codex-subagents/roles/sdp-trace-reviewer.json \
  --context-pack .agents/codex-subagents/contexts/sdp-trace-repo-rules.json \
  --task-file reviews/slice-review-prompt.md \
  --cwd /home/fall_out_bug/projects/vibe_coding/sdp-trace/.worktrees/010-command-package-organization \
  --model minimax/MiniMax-M2.7
```

### Run verification
```bash
codex-subagent run pi \
  --role-card .agents/codex-subagents/roles/sdp-trace-verifier.json \
  --context-pack .agents/codex-subagents/contexts/sdp-trace-repo-rules.json \
  --task "Run full verification suite after slice N and report pass/fail for every gate" \
  --cwd /home/fall_out_bug/projects/vibe_coding/sdp-trace/.worktrees/010-command-package-organization
```

## Model Policy

Per `AGENTS.md`, adversarial review prefers:
- `zai/glm-5.1` (GLM)
- `minimax/MiniMax-M2.7` (MiniMax)

Avoid OpenAI, Anthropic, and Google models unless explicitly permitted.
