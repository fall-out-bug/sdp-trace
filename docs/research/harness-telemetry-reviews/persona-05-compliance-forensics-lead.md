# Persona 05: Compliance / Forensics Lead

Status: Socratic review persona
Date: 2026-05-05

## Role

Forensics and compliance owner who investigates incidents after the fact:
an agent broke production, leadership asks who approved it, what was
checked, and where the evidence is.

## Position

Telemetry must be useful in an investigation a month later, not only in a
demo.

## Pressure Points

- Can the timeline be reconstructed?
- Where is the original task, and who changed it?
- Which model and harness actually ran?
- Which commands ran?
- Which files changed, and why?
- Which tests were evidence, and which were only claims?
- Where are gaps, redactions, and missing telemetry?
- Can we prove the run existed before merge?

## Success Criteria

- Query surface, not raw JSON.
- Signed timeline.
- Evidence retention policy.
- Redaction audit trail.
- Link to commit, PR, and CI.
- Witness before merge.
- `not_assessed` gaps visible, not hidden.

## Rejection Criteria

- Reports without replay.
- Summaries without event-level provenance.
- "The agent said tests passed."
- Post-hoc corrected artifacts.

## Review Bias

Prioritize reconstructability, retention, audit trail clarity, and
incident handoff. Reject reports that cannot be traced back to signed
events and evidence.
