# SpecKit Compatibility

`sdp-trace` should use SpecKit-aligned terms in public docs and examples.

This is conceptual compatibility first. Runtime adapters can come later.

| SpecKit-aligned term | sdp-trace meaning | SDP internal equivalent |
|---|---|---|
| Spec | Intended behavior and acceptance criteria | feature / scope contract |
| Plan | How the change will be made and checked | workstream plan |
| Task | Executable unit of work | beads issue / work item |
| Evidence | Proof collected from execution and review | evidence bundle |
| Gate | External policy check applied to evidence | guard / review / QA check |
| Decision | Human or external automated readiness outcome | verdict / decision record |
| Trace | Link graph across spec, task, evidence, and external decision input | workstream trace |
| Provenance | Actor, model, tool, command, and source chain | attribution chain |
| Beads | Secondary implementation tracker | issue memory / work item mirror |

## Rule

Use SpecKit terms in external docs. Map SDP internals only when needed.

SpecKit artifacts are the planning source of truth in this repository. Beads may mirror execution state and dependencies, but repository observers must be able to understand scope, acceptance, and evidence from committed SpecKit files without reading Beads.
