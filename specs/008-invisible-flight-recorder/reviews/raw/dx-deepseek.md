### DX/UX Review Findings

1. **Command Clarity (Severity: Medium)**
   - The command `sdp-trace recorder run --profile <profile> -- <developer-command...>` is described but lacks examples in the documentation. This could lead to confusion for operators unfamiliar with the syntax, especially the use of `--` to separate flags from the developer command.

2. **Manual Packet Generation (Severity: Low)**
   - The manual packet generation command (`packet build-github --github-input`) is marked as lower-authority, but the documentation does not clearly explain why this is the case or when to use it. Operators might be confused about when to prefer CI-owned packet generation over manual input.

3. **Defaults and Profiles (Severity: Medium)**
   - The `--profile` flag in `sdp-trace recorder run` is required but does not have a default value specified. This could lead to operator friction, especially in automated environments where specifying a profile might be cumbersome.

4. **Command Help Text (Severity: Medium)**
   - The CLI help text (`sdp-trace --help`) is mentioned in the functional requirements (FR-011), but it’s unclear if it will adequately describe the invisible flow or differentiate between authoritative and backfill commands. This could lead to operators using the wrong commands for their needs.

5. **Integration Action Attribution (Severity: Low)**
   - The separation of integration actions from developer route evidence is mentioned (FR-010), but it’s unclear how this separation is presented to operators in the CLI or documentation. This could lead to confusion when interpreting packet metadata.

6. **Documentation Updates (Severity: High)**
   - The documentation updates (`docs/agent-entrypoint.md`, `docs/change-evidence-packet.md`, and install/release docs) are listed as tasks (T020-T022), but it’s unclear how comprehensive these updates will be in explaining the new invisible flow. Incomplete documentation could lead to operator confusion and friction.

7. **Negative Feedback Handling (Severity: Medium)**
   - The builder fails closed with `cannot_verify` when evidence is missing or contradictory (FR-009), but it’s unclear how this feedback is presented to operators. Poor error messaging could lead to frustration and difficulty in diagnosing issues.

### Recommendations
- Add examples for `sdp-trace recorder run` to clarify the use of `--` for separating flags from the developer command.
- Clearly explain the difference between CI-owned packet generation and manual packet generation in the documentation, including when each should be used.
- Consider providing a default profile for `sdp-trace recorder run` to reduce operator friction.
- Ensure that the CLI help text (`sdp-trace --help`) clearly differentiates between authoritative and backfill commands.
- Explicitly document how integration actions are separated from developer route evidence in the packet metadata.
- Conduct a thorough review of the updated documentation to ensure it comprehensively explains the new invisible flow.
- Improve error messaging for `cannot_verify` failures to help operators diagnose missing or contradictory evidence efficiently.
