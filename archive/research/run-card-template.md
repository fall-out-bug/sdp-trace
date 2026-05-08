# Run Card Template

Use one run card per project, phase, harness, and model.

## Metadata

| Field | Value |
|---|---|
| Project |  |
| Project URL |  |
| Commit |  |
| Phase | A / B / C |
| Harness |  |
| Model |  |
| Operator |  |
| Date |  |
| Prompt SHA-256 |  |
| Prompt release approval | not_assessed |

## Prompt

```text
<redacted prompt, prompt summary, or prompt SHA-256>
```

Record the exact prompt only when the accountable evidence owner explicitly approves it for release. Otherwise commit a prompt hash, redaction note, and access-neutral reference.

## Expected Artifacts

- `evidence-bundle.json`
- `provenance-records.json`
- `trace-snapshot.json`
- `export-limitations.md`
- `redaction-note.md`
- optional external-verdict input
- optional `assessment-input.json`
- optional `repo-map.md`
- optional `notes.md`

## Observations

### Build Detection

Record separately:

- primary build system
- dependency metadata systems
- bootstrap scripts
- wrapper scripts

Example: `maven_install.json` inside a Bazel repo is dependency metadata, not proof that Maven owns the build.

### Language Detection

### Evidence Discipline

### Scope Discipline

### Structured Output

### Uncertainty Handling

### Runtime / Token Practicality

## Evidence State

Use `observed` only when a committed sanitized run artifact or evidence summary supports the row. Use `not_assessed` for missing runs, missing exports, unsafe commands, or design fixtures.

| Dimension | Evidence state | Reason code | Artifact reference | External verdict reference | Notes |
|---|---|---|---|---|---|
| Build detection | observed / not_assessed |  |  | none |  |
| Language detection | observed / not_assessed |  |  | none |  |
| Evidence discipline | observed / not_assessed |  |  | none |  |
| Scope discipline | observed / not_assessed |  |  | none |  |
| Structured output | observed / not_assessed |  |  | none |  |
| Uncertainty handling | observed / not_assessed |  |  | none |  |
| Practicality | observed / not_assessed |  |  | none |  |

External tools may emit their own verdict values. Record those only in an external-verdict input with producer, origin, policy reference when available, and artifact reference.

## Follow-Ups

- 
