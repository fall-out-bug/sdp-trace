VERDICT: REVISE

Remaining major finding:

**M1 — Pre-write redaction must be a cross-cutting constraint from Block 13B, not a Block 18 deliverable.**

The recommended order places Block 18 (redaction/retention profile) at position 7, after Blocks 13B, 14, 15, 16, and 17. Block 13B explicitly grants the process wrapper observation of "lifecycle and command-level metadata," and Blocks 14–17 will write gate results, checkpoints, and witness records. The roadmap never states that Blocks 13B–17 must implement pre-write redaction from day one.

From a staff-engineer perspective, a single leaked secret, token in argv, or source snippet in stdout during early adoption is enough to kill trust and cause the tool to be silently uninstalled. The convergence notes moved redaction "before broad capture-depth expansion" (Block 19), but staff engineers do not need adapter expansion to leak secrets; the process wrapper and CI observer surfaces are sufficient.

Required correction:

- Elevate pre-write redaction to a cross-cutting constraint for **all** blocks, stated before Block 13B or as a Block 13B acceptance criterion.
- Block 18's acceptance criterion ("default profile stores no raw prompts, model responses, source snippets, stdout, stderr, tokens, or secrets") must be enforced as a **hard floor** for every earlier block that writes telemetry.
- Block 18 should own the *configurable retention menu* and *redaction audit trail*, but the *safe-by-default pre-write floor* must exist before any artifact is written.
