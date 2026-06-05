# Slice 118 Plan Review

Date: 2026-06-05T01:53:46Z

Scope:
- Consolidate `internal/adaptercapture/adaptercapture_valid_event_spec_type.go`
  into `internal/adaptercapture/adaptercapture_valid_event_specs.go`.
- Preserve unexported `validEventSpec` fields `id`, `eventType`, and
  `sequence`.
- Preserve generated valid adapter events and required-event behavior.
- Exclude valid event generation, event ordering, required event type lists,
  adapter-capture assessment behavior, schemas, examples, dependencies, package
  boundary, dependency direction, CRAP/MI baselines.

Decision gate:
- Simpler/Faster: Move the tiny unexported type next to its only catalog; no new
  abstraction or dependency.
- Blocking Edge Cases: Focused tests must preserve valid generated events and
  required-event behavior without claiming the separate required event type list
  is derived from `validEventSpecs`.
- Existing Open Source: Not applicable; this is local fixture type ownership
  cleanup.

Initial plan review:

| Reviewer | Agent ID | Harness | Model/provider | Prompt class | Result |
|---|---|---|---|---|---|
| Beauvoir | `019e9406-f078-7fd2-b8d0-e22ac17a1e3a` | Codex subagent | not_assessed | plan review | LGTM |
| Peirce | `019e9406-f40c-79f1-904e-54d0f0b73866` | Codex subagent | not_assessed | plan review | finding |
| Halley | `019e9406-f7c2-7f80-80d9-86f7cf7e0c22` | Codex subagent | not_assessed | plan review | LGTM |

Finding:
- major: T021-8210 claimed the required event type list is derived from
  `validEventSpecs`, but the code has a separate `validRequiredEventTypes`
  list. That contradicted the explicit exclusion of required event type list
  changes.

Fix:
- Updated T021-8210 to require preserving generated valid adapter events from
  `validEventSpecs` and preserving required-event behavior separately, without
  changing the separate required event type list.

Re-review:

| Reviewer | Agent ID | Harness | Model/provider | Prompt class | Result |
|---|---|---|---|---|---|
| Beauvoir | `019e9406-f078-7fd2-b8d0-e22ac17a1e3a` | Codex subagent | not_assessed | plan re-review | LGTM |
| Peirce | `019e9406-f40c-79f1-904e-54d0f0b73866` | Codex subagent | not_assessed | plan re-review | LGTM |
| Halley | `019e9406-f7c2-7f80-80d9-86f7cf7e0c22` | Codex subagent | not_assessed | plan re-review | LGTM |

Review state: pass.
