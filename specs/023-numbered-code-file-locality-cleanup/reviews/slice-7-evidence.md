# Slice 7 Evidence: Wrap Command Locality

Status: pass

## Scope

- Package: `cmd/sdp-trace`
- Removed numbered files:
  - `cmd/sdp-trace/wrap_399_run.go`
  - `cmd/sdp-trace/wrap_400_runlegacyrecorder.go`
  - `cmd/sdp-trace/wrap_401_parseargs.go`
  - `cmd/sdp-trace/wrap_402_command.go`
  - `cmd/sdp-trace/wrap_403_writerunresult.go`
  - `cmd/sdp-trace/wrap_404_runwrappedcommand.go`
  - `cmd/sdp-trace/wrap_405_runtaskrecorder.go`
  - `cmd/sdp-trace/wrap_406_parsewrappedcommandargs.go`
  - `cmd/sdp-trace/wrap_407_requirewrappedcommandargs.go`
  - `cmd/sdp-trace/wrap_408_missingrequiredcontract.go`
  - `cmd/sdp-trace/wrap_409_rundryrun.go`
  - `cmd/sdp-trace/wrap_410_runpreview.go`
  - `cmd/sdp-trace/wrap_411_runpreviewcommand.go`
  - `cmd/sdp-trace/wrap_412_writepreviewcommandpayload.go`
  - `cmd/sdp-trace/wrap_413_previewcommandpayload.go`
  - `cmd/sdp-trace/wrap_414_parsepreviewcommandargs.go`
  - `cmd/sdp-trace/wrap_415_loadpreviewcontract.go`
- Added behavior-named files:
  - `cmd/sdp-trace/wrap_legacy.go`
  - `cmd/sdp-trace/wrap_recorder.go`
  - `cmd/sdp-trace/wrap_run.go`
  - `cmd/sdp-trace/wrap_run_args.go`
  - `cmd/sdp-trace/wrap_preview.go`
  - `cmd/sdp-trace/wrap_preview_args.go`
  - `cmd/sdp-trace/wrap_preview_payload.go`

## Local Verification

- `gofmt -w cmd/sdp-trace/wrap_legacy.go cmd/sdp-trace/wrap_recorder.go cmd/sdp-trace/wrap_run.go cmd/sdp-trace/wrap_run_args.go cmd/sdp-trace/wrap_preview.go cmd/sdp-trace/wrap_preview_args.go cmd/sdp-trace/wrap_preview_payload.go`: pass
- `go test ./cmd/sdp-trace`: pass
- `go run ./tools/qualitycheck -mi-under 70 cmd/sdp-trace/wrap_legacy.go cmd/sdp-trace/wrap_recorder.go cmd/sdp-trace/wrap_run.go cmd/sdp-trace/wrap_run_args.go cmd/sdp-trace/wrap_preview.go cmd/sdp-trace/wrap_preview_args.go cmd/sdp-trace/wrap_preview_payload.go`: pass
  - `cmd/sdp-trace/wrap_legacy.go`: MI `71.0`
  - `cmd/sdp-trace/wrap_recorder.go`: MI `71.6`
  - `cmd/sdp-trace/wrap_run.go`: MI `81.7`
  - `cmd/sdp-trace/wrap_run_args.go`: MI `72.9`
  - `cmd/sdp-trace/wrap_preview.go`: MI `75.1`
  - `cmd/sdp-trace/wrap_preview_args.go`: MI `74.0`
  - `cmd/sdp-trace/wrap_preview_payload.go`: MI `76.4`
- Remaining active numbered Go files: `1119`
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
