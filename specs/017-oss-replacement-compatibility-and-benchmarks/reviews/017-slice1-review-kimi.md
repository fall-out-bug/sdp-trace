# Slice 1 Review — Kimi

## tools/osscompat

### Build & Test
- `go test -count=1 ./tools/osscompat` — **PASS**
- `go vet ./tools/osscompat` — **PASS**
- `gofmt -d` — **clean**
- Statement coverage — 54.9%. Uncovered functions are the external-tool integration probes (expected) and `main`/`repoRoot`.

---

## Findings

### 1. Important / Quality / main.go:43
**Silent write-error in single-probe mode.**

```go
r := runProbe(p)
_ = printResults(stdout, []probeResult{r}, *asJSON)
return exitCode([]probeResult{r})
```

The error returned by `printResults` is discarded (`_ = …`). If JSON encoding fails or `stdout` is broken, the user sees no output but receives an exit code based on the probe state. The all-probes path correctly handles this error (lines 49–52). Single-probe mode should do the same.

### 2. Important / Quality+Trust / probe.go:87–91
**`runJSONSchemaFixtures` overclaims verification.**

```go
func runJSONSchemaFixtures() (verifierState, string) {
    return statePass, "fixture validation documented in compatibility doc"
}
```

The probe always returns `statePass` without executing `check-jsonschema` or inspecting any fixture. Per repository trust rules, every `pass` verdict must be evidence-backed or marked `not_assessed`/`cannot_verify`. Tool presence is already checked generically by `runProbe`; this function adds no further evidence, yet claims success. Users seeing `pass` will assume fixtures were validated.

**Recommended fix:** return `stateNotAssessed` or `stateCannotVerify` with a concrete reason, or implement the actual `check-jsonschema` invocation against known fixture paths.

### 3. Important / Quality+Security / probe.go:114–124
**`runCUEImport` mutates the working tree.**

```go
args := []string{
    "import",
    "--package", "sdptrace",
    "schema/flight-recorder-run.schema.json",
}
if out, err := exec.Command("cue", args...).CombinedOutput(); err != nil {
```

`cue import` writes output to a `.cue` file by default (e.g., `schema/flight-recorder-run.schema.cue`). Probes should be read-only; unexpected file creation/overwriting is a side effect. If the tool is run from the repository root, it will create or overwrite an untracked file in the working tree.

**Recommended fix:** redirect `cue import` output to stdout (`-o -` or equivalent) so no files are written, or use a temporary directory.

### 4. Advisory / UX / main.go:22–24
**Duplicate error output on flag parse failure.**

```go
if err := fs.Parse(args); err != nil {
    fmt.Fprintf(stderr, "parse flags: %v\n", err)
    return 2
}
```

`flag.ContinueOnError` already prints the parse error (and usage) to `fs.SetOutput(stderr)`. The explicit `fmt.Fprintf` causes the error message to appear twice.

**Recommended fix:** remove the explicit `fmt.Fprintf`, or set `fs.SetOutput(io.Discard)` and handle all output manually.

### 5. Advisory / DX / probe.go:158–160
**Dead code: `repoRoot()` is never called.**

```go
func repoRoot() string {
    return "."
}
```

No references exist in the package. Remove it or wire it up where needed.

### 6. Advisory / UX / runner.go:44
**Text alignment is not future-proof.**

```go
line := fmt.Sprintf("%-24s %s", r.Name, r.State)
```

Probe names longer than 24 characters will break column alignment. Current longest name is 21 chars (`jsonschema-wrap-drift`), so it works today, but there is no safeguard for future probes.

---

## Summary

| # | Severity | Axis | File:Line | Description |
|---|----------|------|-----------|-------------|
| 1 | Important | Quality | main.go:43 | `printResults` error discarded in single-probe mode |
| 2 | Important | Quality+Trust | probe.go:87 | `runJSONSchemaFixtures` returns `pass` without evidence |
| 3 | Important | Quality+Security | probe.go:114 | `runCUEImport` may create/overwrite `.cue` files |
| 4 | Advisory | UX | main.go:22 | Flag parse errors printed twice |
| 5 | Advisory | DX | probe.go:158 | `repoRoot()` is dead code |
| 6 | Advisory | UX | runner.go:44 | Hardcoded `%-24s` width |

Not `LGTM`. Three Important and three Advisory findings.
