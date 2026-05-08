# SpecKit Compatibility

`sdp-trace` can map SpecKit artifacts into its portable evidence vocabulary.
SpecKit is one supported planning shape, not a required runtime or the only
valid workflow source.

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

Use the portable `sdp-trace` terms in external docs. Map workflow-specific terms
only when needed.

The current repository's working records are SpecKit-shaped. Other teams may use
gsd, Superpowers, Oh My OpenAgent, ticket trackers, or custom workflows as long
as the resulting trace exposes scope, acceptance, evidence, provenance, and
decision inputs without requiring private runtime context.
