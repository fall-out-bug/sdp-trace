No blocking findings in the diff. Here are the key observations:

1. **Packet Build-PR**:
   - Adds `build-pr` subcommand to handle PR packet generation from GitHub Actions/fixtures
   - Correctly handles required flags (`--out`, `--source`) and optional inputs
   - Properly validates input JSON and writes bundle/packet/result artifacts

2. **Prompt Boundary Classification**:
   - Implements `ClassifyPromptBoundary` with clear verdicts (clean, contaminated, missing, malformed)
   - Handles text/digest metadata correctly
   - Enforces `RequirePromptBoundary` flag throughout verification logic

3. **Authority Metadata**:
   - Adds `actor` and `write_authority` fields to bundle entries
   - Properly classifies entry origins (ci_generated, recorder_owned, etc.)
   - Maintains backward compatibility by making new fields optional

4. **Schema Changes**:
   - Extends `evidence-bundle-manifest.v0.schema.json` with new optional fields
   - Enum values are sound and cover all expected cases

The implementation appears correct and backward-compatible. No P0/P1 issues found.
