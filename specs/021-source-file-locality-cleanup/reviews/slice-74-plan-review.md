# Slice 74 Plan Review

Status: pass

Date: 2026-06-04

Scope:
- `cmd/sdp-trace/pr_review_126_runcheck.go`
- `cmd/sdp-trace/pr_review_127_parsecheckargs.go`
- `cmd/sdp-trace/pr_review_128_registercheckflags.go`
- `cmd/sdp-trace/pr_review_129_requirecheckinputs.go`
- `cmd/sdp-trace/pr_review_130_preparecheck.go`
- `cmd/sdp-trace/pr_review_131_executecheck.go`
- `cmd/sdp-trace/pr_review_132_finishcheck.go`
- `cmd/sdp-trace/pr_review_133_writecheckpreview.go`
- `cmd/sdp-trace/pr_review_134_writecheckartifacts.go`
- `cmd/sdp-trace/pr_review_135_writejson.go`
- `cmd/sdp-trace/pr_review_136_validationexit.go`

Planned boundary:
- Move `pr-review check` command orchestration into a cohesive command file.
- Move flag parsing and required input validation into a cohesive args file.
- Move packet/profile preparation and runner execution into a workflow file.
- Move preview/summary publication and validation exit mapping into a
  publication file.
- Move durable artifact writes into an artifacts file.
- Keep shared JSON pretty printing, shared file helpers, packet/profile shared
  readers, repeated flag helpers, runner sets, packet-dir helpers, and exit-code
  helpers out of this slice.

Behavior to preserve:
- Flag-only parsing and positional-argument rejection.
- Required `--out` and packet anchors.
- Packet/profile/readiness failures map to `cannot_verify`.
- Work-dir directory validation.
- Repeated allowed-runner reconstruction from raw args.
- Preview output is non-persisted planning data.
- Run-set persistence happens before ledger and validation publication.
- Summary text is printed only after durable artifacts are written.
- Validation verdict exits through `exitCannotVerify`.
- No package boundary, dependency direction, or MI baseline changes.

Review lanes:
- Lane 1: LGTM (`019e933e-71a0-7331-b52d-59af083707ed`, Hypatia)
- Lane 2: LGTM after fix (`019e933e-790c-74a3-b6f8-284794b56c28`, Turing)
- Lane 3: LGTM (`019e933e-7cc6-7341-a7f6-eb317d718e9d`, Ohm)

Findings:
- Lane 2 reported that T021-5101 did not require focused evidence for all
  behavior the plan says must be preserved: packet/profile/readiness failures
  mapping to `cannot_verify`, repeated allowed-runner reconstruction from raw
  args, and validation verdict exit mapping through `exitCannotVerify`.

Fix:
- Expanded T021-5101 to require packet/profile/readiness `cannot_verify`
  evidence, repeated allowed-runner reconstruction evidence, and validation
  exit mapping evidence.
- Split the initial command/publication files after the MI gate reported
  below-threshold file MI for `pr_review_check_command.go` and
  `pr_review_check_publication.go`.
