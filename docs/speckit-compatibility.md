# SpecKit Compatibility

`sdp-trace` should use SpecKit-aligned terms in public docs and examples.

This is conceptual compatibility first. Runtime adapters can come later.

| SpecKit-aligned term | sdp-trace meaning | SDP internal equivalent |
|---|---|---|
| Spec | Intended behavior and acceptance criteria | feature / scope contract |
| Plan | How the change will be made and checked | workstream plan |
| Task | Executable unit of work | beads issue / work item |
| Evidence | Proof collected from execution and review | evidence bundle |
| Gate | Check applied to evidence | guard / review / QA check |
| Decision | Human or automated readiness outcome | verdict / decision record |
| Trace | Link graph across spec, task, evidence, and decision | workstream trace |
| Provenance | Actor, model, tool, command, and source chain | attribution chain |

## Rule

Use SpecKit terms in external docs. Map SDP internals only when needed.
