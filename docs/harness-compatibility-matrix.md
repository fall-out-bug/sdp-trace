# Harness Compatibility Matrix

This matrix tracks how `sdp-trace` can integrate with agent harnesses and workflow frameworks.

It is a research artifact until each row is backed by a fixture or real integration.

| Harness | Prompt/rules location | Tool log access | Hook support | Evidence export | Risk |
|---|---|---|---|---|---|
| Superpowers | TBD | TBD | TBD | TBD | Needs evaluation |
| Hyperpowers | TBD | TBD | TBD | TBD | Needs evaluation |
| gsd / gsd2 | TBD | TBD | TBD | TBD | Needs evaluation |
| Oh My OpenAgent | TBD | TBD | TBD | TBD | Needs evaluation |
| Paperclip | TBD | TBD | TBD | TBD | Needs evaluation |
| Claude Code | TBD | TBD | TBD | TBD | Needs model/harness split |
| Codex | TBD | TBD | TBD | TBD | Needs model/harness split |
| OpenCode | TBD | TBD | TBD | TBD | Needs per-model evaluation |
| Kilo | TBD | TBD | TBD | TBD | Needs evaluation |
| Pi | TBD | TBD | TBD | TBD | Needs evaluation |

## Compatibility Levels

| Level | Meaning |
|---|---|
| L0 | Manual export only. |
| L1 | Harness can emit evidence files. |
| L2 | Harness can run gates locally. |
| L3 | Harness can publish decision records to PR/MR/CI. |
| L4 | Harness can enforce policy with human override. |

Do not claim compatibility above L0 without a committed example.
