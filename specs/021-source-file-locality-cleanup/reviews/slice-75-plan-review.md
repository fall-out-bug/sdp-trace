# Slice 75 Plan Review

Status: pass

Date: 2026-06-04

Scope:
- `cmd/sdp-trace/pr_review_137_writeindentedpayload.go`
- `cmd/sdp-trace/pr_review_139_filehelpers.go`
- `cmd/sdp-trace/pr_review_142_exitcodes.go`
- `cmd/sdp-trace/pr_review_144_readpacketandprofile.go`
- `cmd/sdp-trace/pr_review_145_readpacketandprofilevalues.go`
- `cmd/sdp-trace/pr_review_146_repeatedflagvalues.go`
- `cmd/sdp-trace/pr_review_147_appendrepeatedflagvalue.go`
- `cmd/sdp-trace/pr_review_148_runnerset.go`
- `cmd/sdp-trace/pr_review_149_packetdir.go`

Planned boundary:
- Move stdout JSON rendering into a generic CLI shared output locality file,
  because protected gate output also calls `writeIndentedPayload`.
- Move write-once output-file refusal and work-dir checks into a filesystem
  safety locality file.
- Move review validation exit mapping into a validation-exit locality file.
- Move packet/profile shared loading into a packet/profile input locality file.
- Move repeated flag reconstruction into a repeated-flags locality file.
- Move runner allow-list normalization into a runner allow-list locality file.
- Move packet-dir derivation into a packet-dir locality file.
- Keep command-specific `pr-review` packet, run, synthesize, validate,
  summarize, and check files unchanged.
- Keep numbered gate family cleanup out of this slice; only preserve the
  protected gate call site that already depends on the generic output helper.

Behavior to preserve:
- JSON payload output remains terminal-only stdout mirroring with a trailing
  newline for both `pr-review` and protected gate callers.
- Output-file helpers require explicit paths and refuse existing files and
  directories.
- Work-dir diagnostics keep `work-dir:` context.
- Review validation `cannot_verify` and `coverage_unresolved` map to
  `exitCannotVerify`; satisfied validation maps to zero.
- Packet/profile loading reads packet first and returns empty values on profile
  failure to avoid mixing partial inputs.
- Raw repeated flags preserve ordered `--key value` and `--key=value` values,
  with parsed fallback only when raw values are absent.
- Runner allow-list parsing trims comma-separated values and ignores empty
  entries.
- Packet-dir derivation accepts either a directory or a packet file path, and
  uses `filepath.Dir` when the path is missing.
- No package boundary, dependency direction, or MI baseline changes.

Review lanes:
- Lane 1: LGTM (`019e934e-5709-77a3-a153-159b14eaff7a`, Boole)
- Lane 2: LGTM (`019e934e-5b02-7ce1-af12-fc3ce7a511c5`, Descartes)
- Lane 3: LGTM after fix (`019e934e-5e3c-78e3-b87b-3ae590191eef`, Planck)

Findings:
- Lane 3 reported that `writeIndentedPayload` is not `pr-review`-only because
  `protected_gate_core.go` also calls it, so a `pr_review`-only boundary and
  focused evidence would overclaim coverage.

Fix:
- Reframed the helper as generic CLI JSON output.
- Added protected gate output preservation to the planned behavior and focused
  evidence.
