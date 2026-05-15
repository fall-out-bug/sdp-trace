I'll conduct adversarial PI review of the 014 spec and current docs. Let me examine the baseline docs first.
→ Read docs/agent-entrypoint.md
→ Read docs/reviewer-entrypoint.md
→ Read docs/concepts.md
✱ Grep "missing_telemetry" 22 matches
✱ Grep "missing_telemetry" 75 matches
$ go run ./cmd/sdp-trace --help 2>&1 | head -60
sdp-trace local recorder and verifier commands.
Usage:
  sdp-trace wrap --name <name> [--contract <file>] [--output-dir <dir>] -- <command...>
  sdp-trace version
  sdp-trace command-surface
  sdp-trace run --task <task-ref> [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace dry-run [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace preview [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace doctor [--contract <file>]
  sdp-trace doctor --profile github-actions-git-hooks-v1 [--out <file>]
  sdp-trace install repo-observer --profile github-actions-git-hooks-v1 [--repository-id <safe-id>] [--write] [--force] [--out <file>]
  sdp-trace interaction relay --task-id <safe-id> --event-type <type> --out <file> -- <forward-command...>
  sdp-trace interaction import-transcript --source preclassified-transcript-import --task-id <safe-id> --events-jsonl <file> --out <file>
  sdp-trace interaction summarize --trace <file> [--out <file>]
  sdp-trace observe setup --profile <session-profile.json> --out <run-dir> [--command <harness-command-preview>]
  sdp-trace observe collect --profile <session-profile.json> --run <run-dir>
  sdp-trace observe session --profile <session-profile.json> --out <run-dir> -- <harness-command...>
  sdp-trace harness observe --profile <harness-profile.json> --source <harness-events.jsonl> --out <run-dir>
  sdp-trace harness validate --profile <harness-profile.json> --run <run-dir> --out <validation.json>
  sdp-trace harness summarize --validation <validation.json>
  sdp-trace envelope summarize --envelope <file> [--out <file>]
  sdp-trace verify <run-dir>
  sdp-trace explain <run-dir>
  sdp-trace query --query <missing-evidence|capture-depth> <run-dir>
  sdp-trace query-pack --pack forensics-basic-v1 --run <run-dir> --out <file>
  sdp-trace query-pack explain --result <file>
  sdp-trace export cross-repo-posture --profile cross-repo-evidence-posture-v1 --selection <file> --out <file>
  sdp-trace export cross-repo-posture explain --result <file>
  sdp-trace export telemetry --profile prometheus-text-v1 --cross-repo-posture <file> --out <file|->
  sdp-trace assess --profile adapter-capture --out <file> --run <run-dir>
  sdp-trace assess --profile managed-harness --out <file> --contract <file> --run <run-dir> --adapter-registry <file> --managed-policy <file> --managed-witness <file>
  sdp-trace assess --profile forensic-retention --out <file> --run <run-dir> --redaction-policy <file>
  sdp-trace assess --profile ci-artifact-observation --out <file> --artifact-manifest <file>
  sdp-trace assess --profile authority-envelope --out <file> --authority-package <file>
  sdp-trace assess preview --profile <adapter-capture|managed-harness|forensic-retention|ci-artifact-observation|authority-envelope> [profile inputs]
  sdp-trace assess explain --assessment-result <file>
  sdp-trace report --out <dir> <runs-root-or-run-dir>
  sdp-trace gate --out <file> <runs-root-or-run-dir>
  sdp-trace witness --kind <github-actions|gitlab-ci|buildkite|customer-pki> --out <file> [--report-dir <dir>] [--witness-envelope <file>] [--customer-pki-authority-policy <file>] [--customer-pki-public-cert <file> | --customer-pki-public-key <file>] [--customer-pki-payload-digest <sha256>] [--customer-pki-freshness-evidence <file>] <runs-root-or-run-dir>
  sdp-trace release-proof --manifest <file> --out <file>
  sdp-trace pr-review packet --out <dir> --repo-id <safe-id> --change-ref <pr|mr|change-id> --base <sha> --head <sha> --diff <file> [--ci-state <state>] [--created-by <actor>]
  sdp-trace pr-review run --packet <dir> --profile <file> --out <dir> [--preview] [--work-dir <dir>] [--allow-external-runner <runner>]...
  sdp-trace pr-review synthesize --packet <dir> --runs <dir> --out <file>
  sdp-trace pr-review validate --packet <dir> --profile <file> --runs <dir> --ledger <file> --out <file>
  sdp-trace pr-review summarize --validation <file> --ledger <file> [--out <file>]
  sdp-trace pr-review check --out <dir> --repo-id <safe-id> --change-ref <pr|mr|change-id> --base <sha> --head <sha> --diff <file> --profile <file> [--work-dir <dir>] [--allow-external-runner <runner>]...
  sdp-trace packet build-pr --source <github-actions|github-fixture> --out <dir> [--github-event <file>] [--checks-json <file>] [--artifacts-json <file>] [--route-manifest <file>] [--github-api-url <url>]
  sdp-trace packet build-github --github-input <file> --out <file>
  sdp-trace packet validate --bundle <file>
  sdp-trace packet check-demo --bundle <file>
  sdp-trace packet render --bundle <file> --out <file>
  sdp-trace validate-fixtures [root-dir]
→ Read docs/claim-authoring.md
→ Read docs/README.md
→ Read docs/concepts.md [limit=15, offset=83]
→ Read docs/agent-entrypoint.md [limit=20, offset=273]
**Issues found:**
**ISSUE 1 — `missing_telemetry` in spec US-002 but absent from canonical exit-code contract**
- `spec.md:34` (US-002) lists `missing_telemetry` as a state a reviewer must distinguish
- `docs/agent-entrypoint.md:273-284` (canonical State And Exit Code Contract) defines only: `observed`, `pass`, `fail`, `not_assessed`, `cannot_verify` — no `missing_telemetry`, no exit-code mapping for it
- `docs/reviewer-entrypoint.md:38-41` repeats the same 5-state contract without `missing_telemetry`
- Implementation has it (`internal/trace/status_types.go:53`, `internal/trace/safety.go:11`, `cmd/sdp-trace/assaptercapturepreviewexpectedevidence.go:6`, etc.) but the user-facing canonical contract omits it
- **A cold user reading the entrypoint docs has no way to know `missing_telemetry` exists or what exit code it produces. The spec's own US-002 acceptance criterion ("distinguish missing_telemetry... without guessing from exit codes") cannot be met by the canonical docs as written.**
**ISSUE 2 — `warn` in `docs/concepts.md:88` is a third undocumented state**
- `docs/concepts.md:88` lists `warn` as an External Verdict alongside `pass`, `fail`, `not_assessed`, `cannot_verify`
- `warn` does not appear in `docs/agent-entrypoint.md:273-284` or in `docs/reviewer-entrypoint.md:38-41`
- No exit code is defined for `warn`
- A cold user reading the entrypoint first would be unable to interpret `warn` correctly
**ISSUE 3 — Full expanded state vocabulary from specs is absent from canonical contract**
- Spec block references show additional states: `unsupported`, `not_integrated`, `retention_limited`, `capped_to_retention_mode`, `suppressed`, `completeness` (`complete`/`partial`/`missing_telemetry`/`unknown`)
- None of these appear in the canonical State And Exit Code Contract in `agent-entrypoint.md:273-284`
- The acceptance criterion "State vocabulary is consistent across README, concepts, agent entrypoint, reviewer entrypoint, and adoption guide" cannot be satisfied when the canonical contract lists only 5 states while implementation and spec use ~10+
**ISSUE 4 — `concepts.md` uses "External Verdict" heading for evidence states; entrypoint uses "State And Exit Code Contract"**
- These are the same concept (verifier outcomes) with different section names and different state sets
- A cold user cannot tell these labels refer to the same vocabulary
**ISSUE 5 — Overclaim checklist in `reviewer-entrypoint.md:135-146` is not identified as the canonical checklist US-004 requires**
- US-004 says "Reviewer entrypoint contains the canonical checklist; README and agent entrypoint link to it"
- `docs/agent-entrypoint.md` has no explicit link to `reviewer-entrypoint.md#what-you-may-state-from-output` as the canonical overclaim checklist
- `docs/README.md` links to `reviewer-entrypoint.md` for "overclaim rules" but does not specifically reference the checklist section
**ISSUE 6 — `docs/agent-entrypoint.md:161` lists `pr-review` but the command surface shows subcommands `pr-review packet`, `pr-review run`, `pr-review synthesize`, `pr-review validate`, `pr-review summarize`, `pr-review check`**
- The English Command Reference table groups by top-level command with a purpose column
- A cold user reading the reviewer entrypoint command list (`docs/reviewer-entrypoint.md:64`) sees only `pr-review` without subcommand guidance
- This matches the entrypoint's stated intent (reviewer subset), but creates discoverability gap for cold users who need to understand `pr-review` has distinct subcommands with different contracts
**ISSUE 7 — Spec FR-006 explicitly prohibits implying interactive CLI support from docs-only work, but no `sdp-trace-claim:future_interactive_guide` marker exists in any doc**
- The spec's own non-goal could serve as a claim tag to prevent overclaim, but no such tag is defined in `docs/claim-authoring.md` or used in the spec docs
