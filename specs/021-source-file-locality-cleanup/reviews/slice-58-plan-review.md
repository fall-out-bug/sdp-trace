# Slice 58 Plan Review

Status: pass

Date: 2026-06-02T17:35:12+03:00

## Scope

Slice 58 is bounded to numbered `cmd/sdp-trace` packet build-pr flow and PR
packet artifact publication shards:

- `packet_040_run.go`
- `packet_041_runbuildpr.go`
- `packet_042_parsebuildproptions.go`
- `packet_043_buildprresult.go`
- `packet_044_writeprartifacts.go`
- `packet_045_renderprmarkdown.go`
- `packet_046_writeprfiles.go`
- `packet_047_writeprartifactfiles.go`
- `packet_048_prartifactfile.go`
- `packet_049_prartifactfiles.go`
- `packet_050_writeprfile.go`

Planned target after MI checks: `cmd/sdp-trace/packet_command_dispatch.go`,
`cmd/sdp-trace/packet_build_pr_run.go`,
`cmd/sdp-trace/packet_build_pr_options.go`,
`cmd/sdp-trace/packet_build_pr_result.go`,
`cmd/sdp-trace/packet_build_pr_artifact_render.go`,
`cmd/sdp-trace/packet_build_pr_artifact_write.go`, and
`cmd/sdp-trace/packet_build_pr_artifact_files.go`.

Excluded from this slice:

- PR input reconstruction and route/error classification (`packet_051`
  through `packet_066`)
- GitHub Actions artifact hydration and API access (`packet_067` onward)
- fixture loading and packet exit constants (`packet_093` onward)

## Behavior Preservation Claims

- packet dispatch missing-subcommand diagnostic stays unchanged
- `packet build-pr` flag defaults stay unchanged
- required flag validation stays unchanged
- JSON `cannot_verify` output is preserved for input reconstruction and
  rendering failures
- validation and live-gate failures preserve cannot_verify result state and
  error aggregation
- output path names stay `bundle.json`, `change-evidence-packet.md`, and
  `build-pr-result.json`
- validation errors and live-gate errors remain aggregated into the build-pr
  result
- output-directory creation keeps the current behavior, with mode subject to
  process umask
- artifact write labels and ordering stay unchanged
- artifact writes still stop on the first failed file write
- package boundary, dependency direction, and MI baselines stay unchanged

## Review Lanes

- scope/correctness: LGTM after focused evidence names and failure-path
  coverage were added
- trust/evidence: LGTM after exact focused commands and trust-sensitive
  behavior coverage were added
- maintainability/DX: LGTM after exact focused evidence was added and the
  target file was renamed to reflect root packet dispatch plus build-pr command
  flow

## Planned Focused Evidence

- `TestPacketDispatchRequiresKnownPacketSubcommand`
- `TestParsePacketBuildPROptionsKeepsFlagDefaultsAndRequiredOut`
- `TestBuildPacketPRResultKeepsPathsAndAggregatesGateErrors`
- `TestRunPacketBuildPRWritesCannotVerifyJSONForInputFailure`
- `TestRenderPacketPRMarkdownDowngradesResultOnFailure`
- `TestWritePacketPRArtifactsWritesCannotVerifyJSONOnRenderFailure`
- `TestWritePacketPRFilesCreatesOutputDirectoryAndArtifacts`
- `TestWritePacketPRArtifactFilesStopsAfterFirstFailure`
- existing `TestPacketBuildPRFixtureCLIWritesArtifacts`

These tests must cover packet dispatch missing-subcommand behavior,
build-pr flag defaults, required `--out` validation, output path names,
validation/live-gate error aggregation, JSON `cannot_verify` behavior for input
reconstruction and rendering failures, cannot_verify result state for
validation/live-gate failures, output-directory creation with mode subject to
process umask, artifact labels and ordering, and first-write-failure
short-circuiting.

Focused commands:

```text
go test ./cmd/sdp-trace -list 'Test(PacketDispatchRequiresKnownPacketSubcommand|ParsePacketBuildPROptionsKeepsFlagDefaultsAndRequiredOut|BuildPacketPRResultKeepsPathsAndAggregatesGateErrors|RunPacketBuildPRWritesCannotVerifyJSONForInputFailure|RenderPacketPRMarkdownDowngradesResultOnFailure|WritePacketPRArtifactsWritesCannotVerifyJSONOnRenderFailure|WritePacketPRFilesCreatesOutputDirectoryAndArtifacts|WritePacketPRArtifactFilesStopsAfterFirstFailure|PacketBuildPRFixtureCLIWritesArtifacts)$'
go test ./cmd/sdp-trace -run 'Test(PacketDispatchRequiresKnownPacketSubcommand|ParsePacketBuildPROptionsKeepsFlagDefaultsAndRequiredOut|BuildPacketPRResultKeepsPathsAndAggregatesGateErrors|RunPacketBuildPRWritesCannotVerifyJSONForInputFailure|RenderPacketPRMarkdownDowngradesResultOnFailure|WritePacketPRArtifactsWritesCannotVerifyJSONOnRenderFailure|WritePacketPRFilesCreatesOutputDirectoryAndArtifacts|WritePacketPRArtifactFilesStopsAfterFirstFailure|PacketBuildPRFixtureCLIWritesArtifacts)$'
```

## Reviewer Metadata

- harness: multi_agent_v1
- model/provider: not_assessed; the harness does not expose exact
  provider-qualified model IDs in this session
- prompt class: Slice 58 SpecKit plan/task review
- timeout: 600000ms
- retries: 0
- fallback: none
