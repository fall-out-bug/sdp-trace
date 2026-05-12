No blocking findings. Here's the summary:

1. **CLI Help**: Updated correctly in `cmd/sdp-trace/main.go` to include `build-pr` usage.

2. **Docs**: 
   - `docs/change-evidence-packet.md` clearly explains `build-pr` vs `build-github` authority levels
   - `docs/evidence-bundle-manifest.v0.schema.json` added new fields with appropriate enums
   - `docs/install.md` mentions `build-pr` as the live GitHub Actions path
   - `docs/evidence-bundle-manifest.md` documents new fields

3. **Diagnostics**: Diagnostic messages and error handling look appropriate, with clear error states and reasons.

4. **Authority Levels**: Properly documents `build-pr` as live authority and `build-github` as backfill/fixture only.

The diff appears to maintain consistency across CLI help, docs, and implementation with appropriate guidance on command usage and authority levels.
