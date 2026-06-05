# Slice 6 Evidence: Interaction Command Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Removed numbered files:
  - `cmd/sdp-trace/interaction_158_run.go`
  - `cmd/sdp-trace/interaction_159_runrelay.go`
  - `cmd/sdp-trace/interaction_160_relayoptions.go`
  - `cmd/sdp-trace/interaction_161_parserelayargs.go`
  - `cmd/sdp-trace/interaction_162_newrelayflagset.go`
  - `cmd/sdp-trace/interaction_163_relaystringflags.go`
  - `cmd/sdp-trace/interaction_164_requirerest.go`
  - `cmd/sdp-trace/interaction_165_requireonlyflagscode.go`
  - `cmd/sdp-trace/interaction_166_requiredflags.go`
  - `cmd/sdp-trace/interaction_167_runimporttranscript.go`
  - `cmd/sdp-trace/interaction_168_writeimportedtranscript.go`
  - `cmd/sdp-trace/interaction_169_importtranscriptfromoptions.go`
  - `cmd/sdp-trace/interaction_170_parseimporttranscriptargs.go`
  - `cmd/sdp-trace/interaction_171_runsummarize.go`
  - `cmd/sdp-trace/interaction_172_parsesummarizeargs.go`
- Added behavior-named files:
  - `cmd/sdp-trace/interaction_command.go`
  - `cmd/sdp-trace/interaction_relay.go`
  - `cmd/sdp-trace/interaction_relay_args.go`
  - `cmd/sdp-trace/interaction_transcript_import.go`
  - `cmd/sdp-trace/interaction_transcript_import_args.go`
  - `cmd/sdp-trace/interaction_summary.go`
  - `cmd/sdp-trace/cli_flag_requirements.go`

## Local Verification

- `gofmt -w cmd/sdp-trace/interaction_command.go cmd/sdp-trace/interaction_relay.go cmd/sdp-trace/interaction_relay_args.go cmd/sdp-trace/interaction_transcript_import.go cmd/sdp-trace/interaction_transcript_import_args.go cmd/sdp-trace/interaction_summary.go cmd/sdp-trace/cli_flag_requirements.go`: pass
- `go test ./cmd/sdp-trace`: pass
- `go run ./tools/qualitycheck -mi-under 70 cmd/sdp-trace/interaction_command.go cmd/sdp-trace/interaction_relay.go cmd/sdp-trace/interaction_relay_args.go cmd/sdp-trace/interaction_transcript_import.go cmd/sdp-trace/interaction_transcript_import_args.go cmd/sdp-trace/interaction_summary.go cmd/sdp-trace/cli_flag_requirements.go`: pass
  - `cmd/sdp-trace/interaction_command.go`: MI `80.5`
  - `cmd/sdp-trace/interaction_relay.go`: MI `76.2`
  - `cmd/sdp-trace/interaction_relay_args.go`: MI `70.9`
  - `cmd/sdp-trace/interaction_transcript_import.go`: MI `74.1`
  - `cmd/sdp-trace/interaction_transcript_import_args.go`: MI `80.7`
  - `cmd/sdp-trace/interaction_summary.go`: MI `73.1`
  - `cmd/sdp-trace/cli_flag_requirements.go`: MI `73.4`
- Remaining active numbered Go files: `1136`
- `go test ./...`: pass
- `go vet ./...`: pass
- `go run ./tools/doccheck`: pass
- `go run ./tools/hygienecheck`: pass
- `jq empty schema/*.json`: pass
- `git diff --check`: pass
- coverage-backed CRAP and MI bundle: pass

## Targeted Reviews

- `opencode-go/glm-5.1`, Opencode, staged-diff review,
  2026-06-01: `LGTM`
- `opencode-go/qwen3.7-max`, Opencode, staged-diff review,
  2026-06-01: `LGTM`
- `opencode-go/deepseek-v4-flash`, Opencode, staged-diff review,
  2026-06-01: `LGTM`

## Trust States

- behavior preservation: pass
- CRAP < 5: pass
- MI > 70: pass
- spec drift: pass
- constitution drift: pass
- product drift: pass
- CleanArch hex: pass
- CleanCode: pass
- SOLID: pass
- DRY: pass
- YAGNI: pass
- production trust: not_assessed
- release approval: not_assessed
- merge approval: not_assessed
