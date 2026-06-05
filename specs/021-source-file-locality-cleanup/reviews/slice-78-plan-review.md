# Slice 78 Plan Review

Status: pass

Date: 2026-06-04

Scope:
- `internal/packet/packet_047_rendercleantheater.go`
- `internal/packet/packet_048_rendertheaterfinding.go`
- `internal/packet/packet_049_rowbyid.go`
- `internal/packet/packet_050_renderdecisions.go`
- `internal/packet/packet_051_renderevidence.go`
- `internal/packet/packet_052_renderresidualgaps.go`
- `internal/packet/packet_053_rendernoresidualgaps.go`
- `internal/packet/packet_054_renderresidualgaprows.go`
- `internal/packet/packet_055_rendernonproof.go`
- `internal/packet/packet_056_requiredrowindex.go`
- `internal/packet/packet_057_resolverfromlist.go`
- `internal/packet/packet_058_md.go`
- `internal/packet/packet_059_packetdigest.go`

Planned boundary:
- Move packet rendering helpers into named rendering locality files.
- Move row/resolver lookup helpers into named lookup locality files.
- Move markdown cell escaping into a markdown escaping locality file.
- Move required-row ordering helper into a row ordering locality file.
- Move packet digest generation into a digest locality file.
- Keep top-level markdown orchestration, validation execution, and downstream
  GitHub row/entry helpers out of this slice.

Behavior to preserve:
- Rendered section headers and table column order.
- Markdown escaping for pipes and newlines.
- `none` rendering for blank cells.
- Clean-theater fallback row shape.
- Resolver fallback behavior.
- Required-row ordering semantics.
- Non-approval fallback text.
- Digest prefix and determinism.
- Package boundary, dependency direction, and MI baselines.

Review lanes:
- Lane 1: `LGTM` (`019e9379-2574-7942-a633-ecfdb571bb2e`, Euclid)
- Lane 2: `LGTM` (`019e9379-28e2-7933-97ec-5629ceff14ad`, Kuhn)
- Lane 3: `LGTM` after fix (`019e9379-2c90-7cd2-80f4-68390ada5fe8`,
  Archimedes)

Findings:
- minor: T021-5381 required helper-level rendering/digest evidence but did not
  require evidence that excluded top-level markdown orchestration stayed
  unchanged or that full `RenderMarkdown` section order and table headers were
  covered.

Fix:
- Updated T021-5381 to require top-level `RenderMarkdown` section order and
  section/table header evidence or source-shape evidence that excluded
  `packet_193` onward orchestration files were untouched.
