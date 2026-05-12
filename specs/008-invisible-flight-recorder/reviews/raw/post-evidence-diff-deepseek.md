No blocking findings.

The diff:
1. Adds a new `packet build-pr` command as the live path for PR packet proof, avoiding manual packet edits.
2. Requires `--out` flag for artifact generation (cannot override with checked-in Markdown/JSON).
3. Classifies prompt boundaries to avoid manual prompt contamination.
4. Adds actor/authority metadata without overclaiming:
   - Uses `ci_packet_builder`, `ci_generated` for CI-generated artifacts.
   - Uses `recorder`, `recorder_owned` for recorded evidence.
   - Uses `operator`, `operator_authored` for default decision owners.
5. Fail-closes with `cannot_verify` if critical evidence is missing (e.g., prompt boundary, workflow run ID).

The live path avoids manual prompt mentions, manual packet edits, and overclaiming authority. Cannot_verify fail-closes are correctly implemented.
