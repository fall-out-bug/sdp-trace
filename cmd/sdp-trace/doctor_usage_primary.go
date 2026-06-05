package main

const usagePrimaryText = `sdp-trace local recorder and verifier commands.

Usage:
  sdp-trace wrap --name <name> [--contract <file>] [--output-dir <dir>] -- <command...>
  sdp-trace version
  sdp-trace command-surface
  sdp-trace run --task <task-ref> [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace dry-run [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace preview [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace doctor [--contract <file>]
  sdp-trace doctor --profile github-actions-git-hooks-v1 [--out <file>]
  sdp-trace override request --out <file> --id <id> --by <actor> --reason <reason> --source-ref <ref> --scope <scope> [--external-reference <ref>]
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
`
