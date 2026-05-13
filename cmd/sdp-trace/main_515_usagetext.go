package main

const usageText = `sdp-trace local recorder and verifier commands.

Usage:
  sdp-trace wrap --name <name> [--contract <file>] [--output-dir <dir>] -- <command...>
  sdp-trace version
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
`
