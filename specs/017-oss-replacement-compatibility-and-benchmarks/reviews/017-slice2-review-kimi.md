# Slice 2 Review: tools/ossbench

Reviewer: Kimi (wide-context code/spec review)
Date: 2026-05-20

## Summary

The benchmark harness is small, readable, and testable. The CLI is straightforward, tests pass, and `go vet` is clean. However, there are several issues across quality, DX, UX, and operational safety that should be addressed before this slice is considered complete.

---

## Findings

### 1. Important | Quality + DX — Dead/incorrect code in `main.go`
**Location:** `main.go:101–118`

The `envInfo` struct, `getEnv()` function, and `init()` function are never referenced and serve no purpose.

- `getEnv()` incorrectly queries `os.Getenv("GOOS")` / `os.Getenv("GOARCH")`. At runtime these environment variables are almost always empty; runtime platform detection should use `runtime.GOOS` and `runtime.GOARCH`.
- `init()` is a no-op: it checks if `GOOS` is empty and then sets `GOOS` to the same empty value.

**Evidence:**
```
$ grep -rn 'getEnv\|envInfo' tools/ossbench/
main.go:101:// envInfo captures the benchmark environment.
main.go:102:type envInfo struct { ... }
main.go:108:func getEnv() envInfo { ... }
main.go:117:func init() { if getEnv().GoOS == "" { ... } }
```

No other file references these symbols. They should be removed to eliminate misleading cruft.

---

### 2. Important | Quality — `benchmarkResult.Iterations` omitted on error paths
**Location:** `bench.go:44–46`, `bench.go:54–58`

When a benchmark fails early (empty command or iteration error), the returned `benchmarkResult` leaves `Iterations` at its zero value (`0`).

- Empty command path: `Iterations` is not populated even though the caller passed a specific `iterations` value.
- Iteration failure path: the error message says `"iteration %d failed"` but `Iterations` in the struct remains `0`.

This produces misleading output. For example, with `-n 20` and a failure on iteration 5, JSON shows `"iterations": 0` and text shows `n=0`.

**Fix:** Populate `Iterations` with the requested count (or the number of completed iterations) on all return paths.

---

### 3. Important | Quality + Security — No execution timeout for benchmarked commands
**Location:** `bench.go:52`

`exec.Command(def.Cmd, def.Args...)` is created without a `context.Context`. If the target command hangs (e.g., `sdp-trace verify` deadlocks on a bug, or a user-supplied command waits on I/O), the benchmark process will hang indefinitely with no recovery mechanism.

A benchmark harness should bound each invocation. Consider `exec.CommandContext` with a configurable or sensible default timeout.

---

### 4. Important | Quality — Missing test for `-bench` success path
**Location:** `main_test.go`

`TestRun_UnknownBench` verifies the error case for an invalid `-bench` name, but there is no test exercising a *valid* built-in name (e.g., `-bench sdp-trace-verify`). This is a primary CLI feature and should have coverage.

---

### 5. Advisory | UX — Text output misaligns on long benchmark names
**Location:** `main.go:78`

```go
fmt.Fprintf(w, "%-24s median=%6.2f ms  ...", r.Name, ...)
```

The hard-coded `%-24s` width breaks alignment when custom command names exceed 24 characters. Since positional arguments become the benchmark name (`strings.Join(remaining, " ")`), users can easily produce long names.

**Example:**
```
very-long-custom-command-name-that-overflows median=  1.23 ms  ...
```

Consider computing the maximum name width dynamically, or using tab-separated values for the text format.

---

### 6. Advisory | Quality — Missing test for `iterations <= 0` default
**Location:** `bench.go:33–35`

`runBenchmark` defaults `iterations` to `20` when the input is `<= 0`. `bench_test.go` does not exercise this branch. A focused test should verify the default behavior.

---

## Correct

- `run()` is well-factored for testability (injected `stdout`/`stderr`).
- `flag.NewFlagSet` with `flag.ContinueOnError` gives clean CLI error handling.
- `bench.go` correctly avoids shell interpretation by using `exec.Command` with separate command and argument slices (no shell injection vector).
- `stats()` handles empty slices safely and has direct coverage.
- All existing tests pass and `go vet` is clean.
- File organization (CLI in `main.go`, logic in `bench.go`) is clean and within complexity limits.

## Verdict

**Not LGTM.** Findings #1 through #4 should be resolved before approval. Findings #5 and #6 are advisory and may be addressed at the author's discretion.
