# Slice 77 Plan Review

Status: pass

Date: 2026-06-04

Scope:
- `internal/packet/packet_033_demofirstpacketchecker.go`
- `internal/packet/packet_034_bundlevalidator.go`
- `internal/packet/packet_035_loadbundle.go`
- `internal/packet/packet_036_loadgithubinput.go`
- `internal/packet/packet_037_buildfromgithubinput.go`
- `internal/packet/packet_038_githubpacket.go`
- `internal/packet/packet_039_appendpromptboundaryfinding.go`
- `internal/packet/packet_040_githubbundlemanifest.go`
- `internal/packet/packet_041_classifypromptboundary.go`
- `internal/packet/packet_042_classifypromptmetadata.go`
- `internal/packet/packet_043_classifyprompttext.go`
- `internal/packet/packet_044_promptboundarymetadatamissing.go`
- `internal/packet/packet_045_promptboundarymetadatacomplete.go`
- `internal/packet/packet_046_forbiddenrecorderdutyphrases.go`

Planned boundary:
- Move validator context structs into named validator context locality files.
- Move JSON loading entrypoints into packet input loading locality files.
- Move GitHub bundle assembly entrypoint, packet shell construction, prompt
  boundary finding append, and bundle manifest construction into named GitHub
  build locality files.
- Move prompt-boundary classification, metadata completeness checks, text
  contamination checks, and forbidden recorder-duty phrase catalog into named
  prompt-boundary locality files.
- Keep rendering, packet digesting, validation execution, downstream GitHub row
  helpers, and downstream GitHub entry helpers out of this slice.

Behavior to preserve:
- Loader error behavior.
- Generated packet and bundle IDs.
- Generated-at UTC formatting.
- Default packet profile fields.
- Integration-action extension behavior.
- Prompt-boundary verdicts, reasons, and route proof effects.
- Contamination theater finding shape.
- Bundle manifest schema, bundle ID, digest, and entries.
- Package boundary, dependency direction, and MI baselines.

Review lanes:
- Lane 1: `LGTM` (`019e936a-f1fd-7e51-ba62-55150e632407`, Confucius)
- Lane 2: `LGTM` (`019e936a-f541-7c71-b178-b630fbbb648e`, Herschel)
- Lane 3: `LGTM` after fix (`019e936a-fce1-70c3-be81-18499da97328`,
  Plato)

Findings:
- minor: T021-5290 did not require source-shape evidence proving numbered
  `packet_033` through `packet_046` files are gone and replacement locality
  filenames exist.
- minor: Slice 77 did not explicitly state `cmd/sdp-trace/FAMILY_INDEX.md` is
  not applicable for non-command `internal/packet` package locality.

Fix:
- Updated T021-5290 to require source-shape evidence and explicit
  `FAMILY_INDEX.md` not-applicable handling.
- Added T021-5312 for focused source-shape evidence.
