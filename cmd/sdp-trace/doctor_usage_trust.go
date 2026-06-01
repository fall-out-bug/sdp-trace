package main

const usageTrustText = `  sdp-trace assess --profile adapter-capture --out <file> --run <run-dir>
  sdp-trace assess --profile managed-harness --out <file> --contract <file> --run <run-dir> --adapter-registry <file> --managed-policy <file> --managed-witness <file>
  sdp-trace assess --profile forensic-retention --out <file> --run <run-dir> --redaction-policy <file>
  sdp-trace assess --profile ci-artifact-observation --out <file> --artifact-manifest <file>
  sdp-trace assess --profile authority-envelope --out <file> --authority-package <file>
  sdp-trace assess preview --profile <adapter-capture|managed-harness|forensic-retention|ci-artifact-observation|authority-envelope> [profile inputs]
  sdp-trace assess explain --assessment-result <file>
  sdp-trace report --out <dir> <runs-root-or-run-dir>
  sdp-trace gate --out <file> <runs-root-or-run-dir>
  sdp-trace checkpoint create --run <run-dir> --out <file> --private-key <file> [--signer-id <id>] [--id <id>]
  sdp-trace checkpoint verify --run <run-dir> --checkpoint <file> [--policy <file>]
  sdp-trace witness --kind <github-actions|gitlab-ci|buildkite|customer-pki> --out <file> [--report-dir <dir>] [--witness-envelope <file>] [--customer-pki-authority-policy <file>] [--customer-pki-public-cert <file> | --customer-pki-public-key <file>] [--customer-pki-payload-digest <sha256>] [--customer-pki-freshness-evidence <file>] <runs-root-or-run-dir>
  sdp-trace release-proof --manifest <file> --out <file>
  sdp-trace pr-review packet --out <dir> --repo-id <safe-id> --change-ref <pr|mr|change-id> --base <sha> --head <sha> --diff <file> [--metadata <file>] [--context <file>]... [--verification <file>]... [--ci-state <state>] [--created-by <actor>]
  sdp-trace pr-review run --packet <dir> --profile <file> --out <dir> [--preview] [--work-dir <dir>] [--allow-external-runner <runner>]... [--not-assessed-reason <reason>]
  sdp-trace pr-review synthesize --packet <dir> --runs <dir> --out <file>
  sdp-trace pr-review validate --packet <dir> --profile <file> --runs <dir> --ledger <file> --out <file>
  sdp-trace pr-review summarize --validation <file> --ledger <file> [--out <file>]
  sdp-trace pr-review check --out <dir> --repo-id <safe-id> --change-ref <pr|mr|change-id> --base <sha> --head <sha> --diff <file> --profile <file> [--metadata <file>] [--context <file>]... [--verification <file>]... [--work-dir <dir>] [--allow-external-runner <runner>]... [--not-assessed-reason <reason>]
`
