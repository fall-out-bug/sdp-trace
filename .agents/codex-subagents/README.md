# Codex Subagent Configuration for sdp-trace

## Roles

Custom role cards for sdp-trace work (used as reference; `codex-subagent run` cannot combine `--role-card` with `--task`/`--task-file`):

- `roles/sdp-trace-slice-worker.json` — Implementation agent for renaming and merging files within a command family slice.
- `roles/sdp-trace-reviewer.json` — Adversarial review agent (model policy: non-OpenAI/Anthropic/Google preferred).
- `roles/sdp-trace-verifier.json` — Verification agent that runs the full quality gate suite.

## Contexts

Context packs **must** be built via `codex-subagent context build` to satisfy the schema validator.

### Build a context pack

```bash
codex-subagent context build \
  --subject "sdp-trace command package organization" \
  --mode dev \
  --goal "Implement family-prefixed file reorganization while preserving behavior and quality gates" \
  --non-goal "changing CLI behavior or introducing subpackages" \
  --non-goal "using non-Go tooling" \
  --file AGENTS.md \
  --file design-note.md \
  --file specs/010-command-package-organization/spec.md \
  --cwd /home/fall_out_bug/projects/vibe_coding/sdp-trace/.worktrees/010-command-package-organization \
  --out .agents/codex-subagents/contexts/sdp-trace-dev.json
```

Available modes: `review`, `council`, `dev`, `research`.

## Usage

### Single subagent with task-file (no role-card/context-pack)

```bash
codex-subagent run pi \
  --task-file .agents/codex-subagents/tasks/verify-current-slice.md \
  --cwd /home/fall_out_bug/projects/vibe_coding/sdp-trace/.worktrees/010-command-package-organization
```

### Multi-role panel (context-pack + built-in roles)

Panels run multiple built-in roles in parallel against one context pack.

```bash
codex-subagent panel run pi \
  --context-pack .agents/codex-subagents/contexts/sdp-trace-dev.json \
  --role reviewer \
  --role evidence-reviewer \
  --cwd /home/fall_out_bug/projects/vibe_coding/sdp-trace/.worktrees/010-command-package-organization \
  --background
```

Check panel status:
```bash
codex-subagent panel status panel_APdUfKPh3S
codex-subagent panel results panel_APdUfKPh3S
```

### Run adversarial review with custom model

```bash
codex-subagent run pi \
  --model minimax/MiniMax-M2.7 \
  --task-file reviews/slice-review-prompt.md \
  --cwd /home/fall_out_bug/projects/vibe_coding/sdp-trace/.worktrees/010-command-package-organization
```

## Model Policy

Per `AGENTS.md`, adversarial review prefers:
- `zai/glm-5.1` (GLM)
- `minimax/MiniMax-M2.7` (MiniMax)

Avoid OpenAI, Anthropic, and Google models unless explicitly permitted.
