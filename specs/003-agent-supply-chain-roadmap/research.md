# Research: Agent Supply Chain Roadmap

## Decision: Reframe Around Agent Supply Chain

**Decision**: Position the roadmap around the agent supply chain of software
delivery, not generic "AI observability" or "agent governance."

**Rationale**: The CTO-level question is not whether an agent existed. The
question is whether a software change has a traceable route from upstream intent
to repo mutation, review, CI, artifact, and human decision owner.

**Alternatives considered**:

- LLM observability framing: too broad and already served by tracing platforms.
- AI governance/GRC framing: too high-level and risks becoming policy instead
  of evidence.
- Coding-agent replacement framing: wrong boundary; `sdp-trace` should observe
  tools, not replace them.

## Decision: GitHub-First, Not GitHub-Bound

**Decision**: Use GitHub as the first change-host adapter because it is the
fastest route to a C-level evidence packet, but keep product fields
provider-neutral.

**Rationale**: GitHub now has first-party agent surfaces, third-party Claude and
Codex agents, custom Copilot agents, and agentic workflows. That makes GitHub a
good first proof surface and a competitive risk.

**Risk**: A GitHub-only model would miss upstream general-agent intent and would
be painful to port to GitLab, GitFlic, Gitea/Forgejo, or Jenkins-centered flows.

**Implication**: Use `change_host`, `change_ref`, `review_ref`, `ci_ref`, and
`artifact_ref` concepts instead of hard-coding GitHub terms in core semantics.

## Decision: Start With Import And Wrapper Modes

**Decision**: P0 evidence modes are post-hoc import and wrapper/sidecar
observation. Native plugins/hooks are P1+.

**Rationale**: Native plugins produce richer evidence but are costly, fragile,
and tool-specific. Post-hoc import and wrapper modes let us validate product
value without coupling `sdp-trace` to every tool.

**Implication**:

- OpenCode/GSD: wrapper plus raw JSONL import.
- `pi`: session import discovery first.
- GSD2: session import discovery plus wrapper feasibility.
- Superpowers: artifact/intent mapping.
- Hermes/OpenClaw: boundary spike, not full plugin.

## Decision: Treat Harnesses As Intent Sources First

**Decision**: GSD, Superpowers, and similar methodology layers should be
recorded as intent, phase, role, task, and checkpoint evidence unless separate
observations prove compliance.

**Rationale**: A plan or skill invocation does not prove the agent followed the
plan. Treating methodology presence as compliance would recreate evidence
theater.

**Implication**: A packet may say "GSD phase declared" or "Superpowers
verification checkpoint requested." It must not say "methodology complied"
without additional evidence.

## Decision: Evaluate GSD2 Separately

**Decision**: GSD2 gets a separate discovery row from GSD v1.

**Rationale**: Public docs describe GSD2 as a standalone CLI coding agent built
on the Pi SDK with direct runtime control, git isolation, context management,
cost/token state, and crash recovery. That is a different evidence surface from
GSD v1 as a harness/methodology layer.

**Implication**: GSD2 discovery should inspect runtime-owned state, not only
generated plans or phase docs.

## Decision: Limit General-Purpose Agent Scope

**Decision**: General-purpose agents are in scope only at software-delivery
boundaries.

**Rationale**: OpenClaw/Hermes-style agents are increasingly used through chat,
gateway, cron, skills, and tools. They can be used by non-technical staff and
can cross into code, repos, CI, and infrastructure. That creates a real C-level
risk. But observing every personal-agent action would turn `sdp-trace` into a
general monitoring product.

**Implication**: `sdp-trace` records only crossings into repositories, change
hosts, CI, artifacts, infra config, release claims, or secret-bearing software
automation.

## Decision: Signed Attestation Is The Top Profile

**Decision**: Signed attestation caps the trust ladder after packet semantics
stabilize.

**Rationale**: Signing weak or incomplete evidence can make theater look
official. The product should first make gaps explicit, then bind mature packets
to DSSE/in-toto/Sigstore or customer private equivalents.

**Implication**: P0 packets may be local or CI-witnessed. P2 signed packets must
record identity, policy, source refs, witness refs, packet digest, and freshness
evidence.

## Integration Notes

These notes are discovery pointers, not support claims. Re-verify before using
them in external-facing materials.

- GitHub third-party agents currently document Claude and Codex as supported
  agents, with agent sessions creating PRs and requesting review.
- GitHub custom agents are Copilot agent profiles, not arbitrary OSS-agent
  runtime registration.
- GitHub Agentic Workflows run agent-driven repository automation through
  GitHub Actions and emphasize read-only-by-default workflow execution.
- OpenCode exposes terminal-native AI coding behavior plus MCP, LSP, GitHub
  Copilot experimental support, and self-hosted provider configuration.
- `pi` is an agent harness mono repo with a coding-agent CLI, agent runtime,
  unified LLM API, and explicit public session-sharing motivation.
- GSD v1 is a meta-prompting/context/spec-driven system for multiple coding
  harnesses.
- GSD2 is a standalone coding agent built on the Pi SDK.
- Superpowers is a multi-host skills/workflow methodology; useful as checkpoint
  and intent evidence.
- Oh My OpenAgent is a high-autonomy OpenCode harness layer with multi-agent
  orchestration; useful as a future harness-intent source but too broad for P0.
- OpenClaw and Hermes are general-purpose agents with gateways, channels,
  tools, memory, skills, or scheduled automation; useful only for boundary
  spikes into software delivery.

## External Sources

Re-check these before turning discovery notes into product or sales claims.

- GitHub custom agents:
  https://docs.github.com/en/copilot/how-tos/copilot-sdk/use-copilot-sdk/custom-agents
- GitHub Copilot third-party agents in VS Code:
  https://code.visualstudio.com/docs/copilot/agents/third-party-agents
- GitHub Agentic Workflows technical preview:
  https://github.blog/changelog/2026-02-13-github-agentic-workflows-are-now-in-technical-preview
- GitHub Agentic Workflows safe outputs:
  https://github.github.com/gh-aw/reference/safe-outputs/
- OpenCode:
  https://github.com/opencode-ai/opencode
- Pi:
  https://github.com/badlogic/pi-mono
- GSD:
  https://github.com/gsd-build/get-shit-done
- GSD2:
  https://github.com/gsd-build/gsd-2
- Superpowers:
  https://github.com/obra/superpowers
- Oh My OpenAgent:
  https://ohmyopenagent.com/docs
- OpenClaw:
  https://github.com/openclaw/openclaw
- Hermes Agent:
  https://github.com/NousResearch/hermes-agent

## Research Gaps

- Exact `pi` local session storage/export format.
- Exact GSD2 state database/session export format and redaction safety.
- Whether Superpowers emits stable artifacts across Codex, OpenCode, Copilot
  CLI, and Claude Code hosts.
- Whether Hermes or OpenClaw has a stable event/session API suitable for a
  safe boundary spike.
- Minimal GitHub evidence packet shape that produces a CTO wow without needing
  a dashboard.
- Minimal signed private-equivalent profile for customers that cannot use public
  Sigstore/Rekor.
