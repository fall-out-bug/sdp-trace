# Pi And Kimi Readiness Note

**Checked**: 2026-05-17
**Pi version**: `0.74.1`
**Local Pi settings**: `.pi/settings.json`

## Current Local State

- `.pi/settings.json` loads project prompts and extensions and uses high thinking by default.
- Project prompts currently include GLM, Qwen, DeepSeek, and boot reminders.
- Pi lists Kimi-capable models:
  - `kimi-coding/kimi-for-coding`
  - `kimi-coding/kimi-k2-thinking`
  - `openrouter/moonshotai/kimi-k2.5`
  - `openrouter/moonshotai/kimi-k2.6`
- `KIMI_API_KEY`, `MOONSHOT_API_KEY`, and `OPENROUTER_API_KEY` were not present in the shell environment at check time.
- A non-interactive smoke prompt to `kimi-coding/kimi-for-coding` returned `OK`, so this local machine may have Pi-managed credentials outside the shell environment. Do not treat that as portable repo evidence.

## Required Configuration For Repeatable Kimi Review

At least one credential path must be configured in the runner environment:

- `KIMI_API_KEY` for Pi's `kimi-coding` provider;
- `MOONSHOT_API_KEY` if using Moonshot's direct OpenAI-compatible API provider;
- `OPENROUTER_API_KEY` if using OpenRouter-hosted Moonshot/Kimi models.

Recommended non-interactive review command for this repo:

```bash
pi --provider kimi-coding --model kimi-for-coding --thinking high --tools read,grep,find,ls -p @specs/004-mvp-readiness-hardening/followup-hardening-spec.md '/review-kimi follow-up readiness hardening spec contract: AGENTS.md trust rules, followup-hardening-spec.md PI Review Prompt'
```

Fallback when the direct Kimi provider is unavailable:

```bash
pi --model openrouter/moonshotai/kimi-k2.5 --thinking high --tools read,grep,find,ls -p @specs/004-mvp-readiness-hardening/followup-hardening-spec.md '/review-kimi follow-up readiness hardening spec contract: AGENTS.md trust rules, followup-hardening-spec.md PI Review Prompt'
```

## External References Checked

- Pi CLI help on this machine lists `MOONSHOT_API_KEY` and `KIMI_API_KEY` environment variables and supports provider-qualified model IDs.
- Pi model listing on this machine exposes `kimi-coding/kimi-for-coding` and `kimi-coding/kimi-k2-thinking`.
- Kimi API Platform docs: https://platform.kimi.ai/docs/api/overview
- Kimi Code docs: https://www.kimi.com/code/docs/en/
- Pi Kimi coding package docs: https://pi.dev/packages/pi-kimi-coder
- Moonshot/Kimi public docs describe OpenAI-compatible API usage and current Kimi models; durable repo policy should say "current official Moonshot coding model" instead of pinning unreplayed future model names.

## Evidence Rules

- A Kimi response is advisory until checked against full files.
- Missing, expired, or hidden credentials make Kimi review `cannot_verify`, not `pass`.
- OpenRouter model availability is not the same evidence as Moonshot direct model availability.
- Benchmark claims are not evidence that a review finding is correct.
