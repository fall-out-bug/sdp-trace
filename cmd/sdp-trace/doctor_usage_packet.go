package main

const usagePacketText = `  sdp-trace packet build-pr --source <github-actions|github-fixture> --out <dir> [--github-event <file>] [--checks-json <file>] [--artifacts-json <file>] [--route-manifest <file>] [--github-api-url <url>]
  sdp-trace packet build-github --github-input <file> --out <file>
  sdp-trace packet validate --bundle <file>
  sdp-trace packet check-demo --bundle <file>
  sdp-trace packet render --bundle <file> --out <file>
  sdp-trace validate-fixtures [root-dir]
`
