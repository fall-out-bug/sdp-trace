# Agent Instructions

`sdp-trace` is the portable trust substrate for AI-assisted delivery.

## Purpose

Define traceability, provenance, evidence, gate verdicts, and decision records that work across coding agents and harnesses.

## Boundary

This repo must stay independent from `sdp_lab`, Beads, Operator Mode, and any specific harness runtime.

Allowed:
- JSON schemas
- Markdown docs
- portable examples
- tiny validation/rendering tools when needed

Not allowed:
- dependency on SDP Operator Mode
- dependency on Beads
- dependency on agentloop
- hidden assumptions about Claude, Codex, OpenCode, or GitHub

## Product Language

Use SpecKit-aligned terms first: spec, plan, task, evidence, gate, decision, trace, provenance.

Avoid internal SDP terms unless a doc explicitly maps them.

## Quality Bar

Every claim about a gate or verdict must be evidence-backed or marked `not_assessed`.

No opaque health scores.

## Commands

No build is required yet.

Run schema checks with:

```bash
jq empty schema/*.json
```
