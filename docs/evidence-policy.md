# Evidence Policy

Evidence policy defines what must be present before a gate may return `pass`.

## Rules

1. A gate cannot pass on model confidence alone.
2. A missing required source produces `not_assessed` or `fail`, never `pass`.
3. Manual evidence is allowed, but it must name the human actor and reason.
4. Generated claims must link to inspected files, commands, tests, or review artifacts.
5. Opaque health scores are not evidence.

## Evidence Strength

| Strength | Meaning | Example |
|---|---|---|
| Strong | Machine or source artifact can be inspected. | CI URL, test log, diff, file path |
| Medium | Human sign-off with clear scope. | reviewer approval with comment |
| Weak | Model summary or inferred claim. | "looks good" from an agent |
| None | Missing or unavailable. | no tests found and no explanation |

Weak evidence can support a `warn`. It should not support `pass` for blocking gates.
