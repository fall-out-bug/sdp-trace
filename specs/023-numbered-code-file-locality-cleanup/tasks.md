# Tasks: Numbered Code File Locality Cleanup

Status: in_progress

## Active Slice 1 Tasks

- [x] T023-001 Inventory remaining numbered active Go files by package.
- [x] T023-002 Confirm Slice 1 is bounded to `cmd/sdp-trace` release-proof
  command shards.
- [x] T023-003 Reject broader release-proof grouping when local pre-change MI
  analysis measures candidate files below the absolute file-MI threshold.
- [x] T023-010 Move selected release-proof declarations into
  `cmd/sdp-trace/release_proof_run.go`,
  `cmd/sdp-trace/release_proof_args.go`, and
  `cmd/sdp-trace/release_proof_policy.go`.
- [x] T023-020 Run `gofmt` on changed Go files.
- [x] T023-030 Run focused Go verification for `cmd/sdp-trace`.
- [x] T023-040 Run repository verification: `go test ./...`, `go vet ./...`,
  `go run ./tools/doccheck`, `go run ./tools/hygienecheck`, `jq empty
  schema/*.json`, and `git diff --check`.
- [x] T023-050 Run CRAP and MI quality gates without changing MI baselines.
- [x] T023-060 Run three independent reviewer lanes and record Slice 1 evidence
  in `specs/023-numbered-code-file-locality-cleanup/reviews/slice-1-evidence.md`.
