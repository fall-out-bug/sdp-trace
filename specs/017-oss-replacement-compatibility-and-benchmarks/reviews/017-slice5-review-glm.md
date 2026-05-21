# Slice 5 Review — GLM Plane

Reviewer: GLM adversarial review
Date: 2026-05-20
Scope: Integration and Final Polish for spec 017
Files reviewed: `tasks.md`, `spec.md`, `plan.md`, `docs/roadmap.md`,
`docs/oss-replacement-compatibility.md`, `docs/oss-benchmark-results.md`,
`docs/oss-policy-prototype.md`, `docs/oss-supply-chain-prototype.md`,
`tools/osscompat/*`, `tools/ossbench/*`

## Verification Baseline

```
go test -count=1 ./...       — PASS (all packages)
go vet ./...                 — PASS (no findings)
go run ./tools/doccheck      — PASS (exit 0)
go run ./tools/hygienecheck  — PASS (exit 0)
git diff --check             — PASS (no whitespace errors)
```

Tool spot-checks:

```
go run ./tools/osscompat -list   — 10 probes listed; in current environment all 10 are not_assessed (no optional tools installed). With check-jsonschema present: 1 pass (fixture validation), 1 fail (live wrap drift), 8 not_assessed (remaining tools absent). Manual-only probes report cannot_verify when their tools are present.
go run ./tools/osscompat -json   — structured JSON output with state/reason
go run ./tools/ossbench -list    — 2 built-in benchmarks (sdp-trace-version, sdp-trace-wrap)
go run ./tools/ossbench -n 3 -json — produces min_ms, max_ms, median_ms, iterations, command, working_directory, binary_path, and argv
```

---

## Findings

### F1. Important — DX: Roadmap Capability Index contradicts Active Specs table for spec 017

- `docs/roadmap.md:35` — Active Specs table: spec 017 status `in_progress`
- `docs/roadmap.md:96` — Capability Index: spec 017 status `Draft`

The same document assigns two different statuses to the same spec. A reader
scanning the Capability Index sees `Draft`; a reader scanning the Active table
sees `in_progress`. The Capability Index was not updated when spec 017 moved
to the Active Specs table.

**Fix:** Change line 96 from `Draft` to `In progress` to match line 35.

---

### F2. Important — Quality: Benchmark doc disclaimer is stale; ossbench tool now satisfies FR-017-004 partially

- `docs/oss-benchmark-results.md:31-34` — "does not satisfy FR-017-004 until `tools/ossbench` produces reproducible structured output with full statistics"
- `specs/.../spec.md:59` — FR-017-004: "benchmark output with median, min, max, iterations, command, and environment"
- `tools/ossbench` verified to emit structured JSON with `min_ms`, `max_ms`, `median_ms`, `iterations`, and error details

The disclaimer was accurate when written but is now stale: the ossbench tool
exists, passes tests, and produces the required structured output. The doc
table still shows `—` for all min/max cells. The disclaimer should be updated
to reflect the current state: the tool can produce full stats for its 2
built-in probes, and the doc table should be regenerated from the tool or
the disclaimer rewritten to state what subset is covered.

Note: the ossbench tool only has 2 built-in probes (sdp-trace-version,
sdp-trace-wrap). The other 5 probes in the doc table (shell, OPA, Cosign,
in-toto, check-jsonschema) require external tools and are not automatable
by the Go harness without those tools installed. This limitation is
acceptable but should be documented explicitly.

**Fix:** Regenerate the doc table rows for the 2 built-in benchmarks from
`ossbench -json` output, or update the disclaimer to state that the tool
covers built-in probes only and the remaining rows are manual one-shot data.

---

### F3. Advisory — DX: Roadmap "Last updated" date is stale

- `docs/roadmap.md` final line before claim tag: `Last updated: 2026-05-16`
- Spec 017 entry was added/updated in the roadmap as part of Slice 5 work on 2026-05-20

The date does not reflect the current content changes.

**Fix:** Update to `2026-05-20` or the actual last-modified date.

---

### F4. Advisory — DX: Prototype docs have `Status: draft` while parent spec is `in_progress`

- `docs/oss-policy-prototype.md:3` — `Status: draft`
- `docs/oss-supply-chain-prototype.md:3` — `Status: draft`
- Parent spec `specs/.../spec.md:3` — `Status: in_progress`

The prototype docs are completed deliverables of an in-progress spec. The
`draft` status is defensible (they describe experimental prototypes) but may
confuse readers into thinking the work hasn't been done. Consider `in_progress`
or adding a note that `draft` means "experimental, not production-ready."

**Suggestion:** Either align status with parent spec or add a one-liner
clarifying that `draft` reflects the experimental nature of the prototype,
not incomplete work.

---

### F5. Advisory — Quality/Security: in-toto reproduction command uses `--key /dev/null`

- `docs/oss-replacement-compatibility.md` (in-toto section) — `in-toto-run ... --key /dev/null -- /bin/true || true`
- Probe result row claims: "Wrap command, sign link metadata, record material/product hashes — `pass`"

The reproduction command passes `/dev/null` as the signing key, which likely
does not produce signed link metadata. The `|| true` masks any error. A reader
copy-pasting this command will not reproduce the claimed `pass` result.

The actual probe was likely run with a real key in a different environment.
The reproduction command should either use a generated key (like the Cosign
example does) or note that the key must be provided.

**Suggestion:** Add `in-toto-keygen` or a pre-step that generates a throwaway
key, similar to the Cosign reproduction command's `cosign generate-key-pair`.

---

## What Looks Good

**Honest status management.** All `not_assessed` and `cannot_verify` states
are correctly applied. The osscompat tool reports `not_assessed` when external
tools are missing and `cannot_verify` for the wrap-schema drift. No local-only
success is inflated to production trust. This is exemplary trust discipline.

**Clean tool implementation.** `tools/osscompat` and `tools/ossbench` are
small, readable Go packages with focused tests. Both pass `go vet` and
`go test`. Structured JSON output is available from both tools.

**Non-authoritative benchmark framing.** The benchmark doc is scrupulous
about not creating health scores or readiness claims from local numbers.
FR-017-005 ("benchmarks non-authoritative") is well satisfied.

**Correct probe state taxonomy.** The compatibility doc uses `pass`, `fail`,
`cannot_verify`, and `not_assessed` correctly and consistently. Failed probes
are marked as expected failures where appropriate (e.g., SLSA negative path).

**No Node.js/npm in product path.** Confirmed: all new tooling is Go. No
JS/TS/MJS files introduced.

**Substitution boundary documentation.** Each OSS tool has a clear breakdown
of what it can replace, what needs adapter glue, and what remains
sdp-trace-specific. This is well-structured decision documentation.

**Reproduction commands use subshell isolation.** All commands in the
compatibility doc use `(cd ... && ...)` pattern per AGENTS.md scanner
verification requirements.

---

## Summary

| # | Severity | Axis | Summary |
|---|---|---|---|
| F1 | Important | DX | Roadmap Capability Index says spec 017 `Draft` but Active table says `in_progress` |
| F2 | Important | Quality | Benchmark doc disclaimer stale; ossbench now produces structured min/max output |
| F3 | Advisory | DX | Roadmap "Last updated" date is 2026-05-16; content changed 2026-05-20 |
| F4 | Advisory | DX | Prototype docs `Status: draft` vs parent spec `in_progress` |
| F5 | Advisory | Quality | in-toto reproduction command uses `/dev/null` key; won't reproduce claimed `pass` |

**Verdict:** 2 Important, 3 Advisory. Not LGTM. Fix F1 through F5 before PR merge.
