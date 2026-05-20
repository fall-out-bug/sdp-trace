# Slice 2 Review: `tools/ossbench` Benchmark Harness

**Reviewer:** GLM (adversarial review)  
**Date:** 2026-05-20  
**Files reviewed:** `tools/ossbench/main.go`, `tools/ossbench/main_test.go`, `tools/ossbench/bench.go`, `tools/ossbench/bench_test.go`

## Verification Summary

- `go vet ./tools/ossbench/` — clean (zero findings)
- `go test -count=1 -v ./tools/ossbench/` — all 11 tests pass (0.095s)
- No TODO/FIXME markers present
- No Node.js/npm/JS/TS violations
- Total source: ~160 lines across 2 files; well within decomposition limits

---

## Findings

### F1. Advisory — Quality — Dead code: `envInfo`, `getEnv()`, and `init()`

**Location:** `main.go:98-117`

`envInfo` struct, `getEnv()` function, and the `init()` function are defined but their output is never consumed anywhere — not in text output, not in JSON output, not in any computation. The `init()` function at lines 115-119 is a particular concern:

```go
func init() {
	if getEnv().GoOS == "" {
		_ = os.Setenv("GOOS", os.Getenv("GOOS"))
	}
}
```

This reads `GOOS` via `os.Getenv("GOOS")`, checks if it's empty, and if so, sets `GOOS` to the result of `os.Getenv("GOOS")` — which is already empty. It is a no-op regardless of the environment state. The `_ =` discard of the error return also hides any potential failure silently.

**Recommendation:** Either remove the dead code or wire `getEnv()` into the JSON output (e.g., as a top-level `environment` field in the JSON results). Given the repo's quality bar, dead code should be removed.

---

### F2. Advisory — Quality — First failed iteration aborts entire benchmark

**Location:** `bench.go:47-52`

```go
if err := cmd.Run(); err != nil {
    return benchmarkResult{
        Name:  def.Name,
        Error: fmt.Sprintf("iteration %d failed: %v", i, err),
    }
}
```

On the first failure of any iteration, `runBenchmark` returns immediately with a single error string. All timing data from successful iterations before the failure is discarded. For benchmark harnesses, a more robust approach would be to either (a) continue running remaining iterations and report partial results + the error, or (b) allow the caller to configure a failure threshold. The current behavior means a single flaky run at iteration 19 of 20 produces zero timing data.

This is a design trade-off; the current behavior is documented by the code and is not a bug per se, but it reduces the tool's value under real-world flaky-command conditions.

---

### F3. Advisory — Security — `exec.Command` uses PATH lookup without validation

**Location:** `bench.go:45`

```go
cmd := exec.Command(def.Cmd, def.Args...)
```

When the user passes a bare command name (e.g., `true`), `exec.Command` performs PATH lookup. This is standard Go behavior and not a vulnerability in itself — arguments are passed separately (no shell injection), and the tool is a local benchmark harness, not a network-facing service. However, since the built-in benchmarks hardcode `"sdp-trace"` as the command, a compromised PATH could cause the tool to execute a malicious binary. The built-in commands use no absolute paths.

**Recommendation:** For the built-in benchmarks, consider using the absolute path to the `sdp-trace` binary (resolved at build time or relative to `os.Args[0]`). For custom user commands, PATH lookup is expected behavior and no change is needed.

---

### F4. Advisory — DX — No `-help` flag; `-help` exits with code 2 and no usage guidance

**Location:** `main.go:26-30`

The `flag` package's built-in `-help` output prints usage text to stderr but exits with code 2. The usage text lists flags but does not explain that positional arguments are treated as custom commands, nor does it describe the built-in benchmarks. Running `ossbench` with no arguments runs all built-in benchmarks silently — this could surprise users who expect help output.

**Recommendation:** Consider adding a short usage preamble to the flag output explaining the command pattern: `ossbench [flags] [command args...]` and pointing to `-list` for built-in benchmarks.

---

### F5. Advisory — Quality — `AllMs` always populated in JSON but never in text mode

**Location:** `bench.go:68` (population), `main.go:73-81` (text output)

`AllMs` is always populated with raw timing data in `runBenchmark`, but `printResults` only emits it in JSON mode. In JSON output with high iteration counts, this can produce very large payloads. The `omitempty` tag means nil slices would be omitted, but `AllMs` is never nil from `runBenchmark`.

**Recommendation:** Consider either (a) omitting `AllMs` from the default JSON output and adding a `-raw` flag to include it, or (b) using a fixed-size summary (p50/p95/p99) instead. This is a minor UX concern, not a bug.

---

## What Is Already Good

1. **Clean separation of concerns.** `bench.go` handles benchmark execution and statistics; `main.go` handles CLI parsing and output formatting. Easy to read and reason about.

2. **Proper testability.** `run()` accepts `io.Writer` parameters for stdout/stderr, making it fully testable without OS-level output capture. All test cases use `bytes.Buffer`.

3. **Good test coverage of edge cases.** Tests cover: list mode, custom commands, JSON output, unknown benchmark name, missing commands, exit codes, empty stats, odd/even/single-element stats, no-command-specified error.

4. **Safe command execution.** `exec.Command` with separate args array prevents shell injection. Arguments are not interpolated into a shell string.

5. **Defensive input handling.** Negative/zero iteration counts fall back to sensible defaults (bench.go:30-32). Empty `Cmd` field returns an explicit error (bench.go:34-38).

6. **Correct median calculation.** Both odd and even-length inputs are handled correctly, verified by tests.

7. **No linter issues.** `go vet` produces zero findings.

---

## Summary

| # | Severity | Axis | Finding |
|---|----------|------|---------|
| F1 | Advisory | Quality | Dead code: `envInfo`, `getEnv()`, `init()` never consumed |
| F2 | Advisory | Quality | First iteration failure discards all prior timing data |
| F3 | Advisory | Security | Built-in benchmarks use PATH-lookup without absolute path |
| F4 | Advisory | DX | No usage preamble; bare invocation silently runs built-ins |
| F5 | Advisory | Quality | `AllMs` always included in JSON; can be large for high `-n` |

**Overall:** No Critical or Important findings. The code is clean, well-tested, and well-structured. The five Advisory findings are minor improvements that would polish the tool but are not blockers.
