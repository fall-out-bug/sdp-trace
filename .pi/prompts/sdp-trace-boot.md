---
description: Mandatory sdp-trace process boot reminder
---
When working in the `sdp-trace` repository:

1. **AGENTS.md is loaded** as a context file. The Skills Router section is authoritative.
2. **On "берем блок в работу" or any block-intake request:**
   - Immediately invoke `/skill:sdp-trace-trust-workflow`
   - Follow `intake_protocol` step by step
   - Do NOT ask exploratory questions before intake completes
3. **On review requests:**
   - Immediately invoke `/skill:pi-review`
   - Use `claim-doubt-cycle.md` for trust claims
   - Prefer non-OpenAI, non-Anthropic, non-Google models:
     - `/review-glm` for architecture doubt
     - `/review-qwen` for wide-context code review
     - `/review-deepseek` for reasoning review
4. **Do NOT explore pi infrastructure** (`--list-models`, RPC docs, etc.) unless explicitly asked.
5. **Do NOT improvise** process steps outside the loaded skills.
