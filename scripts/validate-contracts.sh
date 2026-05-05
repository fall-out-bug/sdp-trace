#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

tmp_validation_dir="$(mktemp -d)"
tmp_release_repo=""
trap 'rm -rf "$tmp_validation_dir" "$tmp_release_repo"' EXIT

jq empty schema/*.json
find examples -name '*.json' \
  -not -path '*/.git/*' \
  -not -path '*/.beads/*' \
  -not -path '*/.sdp-trace-runs/*' \
  -not -path '*/benchmarks/repos/*' \
  -not -path '*/node_modules/*' \
  -print0 | xargs -0 jq empty

validate() {
  local schema="$1"
  local data="$2"
  node scripts/validate-json-schema.mjs "$schema" "$data"
}

expect_fail() {
  local schema="$1"
  local data="$2"
  local negative_out="$tmp_validation_dir/negative.out"
  if validate "$schema" "$data" >"$negative_out" 2>&1; then
    echo "Expected validation failure, but passed: $data" >&2
    cat "$negative_out" >&2
    exit 1
  fi
  rm -f "$negative_out"
}

validate schema/trace.schema.json examples/github-speckit/trace.json
validate schema/evidence-bundle.schema.json examples/go-service/evidence-bundle.json
validate schema/gate-verdict.schema.json examples/go-service/gate-verdict.json
validate schema/evidence-bundle.schema.json examples/jvm-bazel/evidence-bundle.json
validate schema/assessment-input.schema.json examples/contract-foundation/positive-assessment-input.json
validate schema/assessment-input.schema.json examples/contract-foundation/not-assessed-assessment-input.json
validate schema/contract-manifest.schema.json examples/contract-foundation/contract-manifest.example.json
validate schema/contract-release-verification.schema.json examples/contract-foundation/contract-release-verification.example.json
validate schema/review-ledger.schema.json specs/001-sdp-trace-time-series-evidence-substrate/blocks/07-minimum-trust-kernel-review-ledger.json
validate schema/trusted-identity-policy.schema.json examples/contract-foundation/trusted-identity-policy.example.json
validate schema/trusted-identity-policy.schema.json examples/contract-foundation/trusted-identity-policy-wrong-profile.example.json
validate schema/consumer-schema-version-declaration.schema.json examples/contract-foundation/sdp-gate-consumer-declaration.example.json
validate schema/flight-recorder-event.schema.json examples/flight-recorder/local-positive/events/000-run-started.json
validate schema/flight-recorder-event.schema.json examples/flight-recorder/witnessed-positive/events/000-run-started.json
validate schema/flight-recorder-run.schema.json examples/flight-recorder/local-positive/run.json
validate schema/flight-recorder-run.schema.json examples/flight-recorder/witnessed-positive/run.json
validate schema/flight-recorder-witness.schema.json examples/flight-recorder/witnessed-positive/witness.json
scripts/test-flight-recorder.sh
node scripts/validate-pilot-matrices.mjs
scripts/test-e2e-pilot-package.sh
scripts/test-e2e-runner.sh
scripts/validate-e2e-pilot-package.sh --mode package examples/pilot-runs/opencode-minimax-kotlin-bazel

for generated_proof_artifact in \
  examples/contract-foundation/contract-release-verification.example.json \
  examples/contract-foundation/contract-release.dsse.json \
  examples/contract-foundation/local-dev-signing-public.pem \
  examples/self-trace/self-attestation-verification-result.json; do
  if jq -e --arg path "$generated_proof_artifact" '.artifacts[] | select(.path == $path)' examples/contract-foundation/contract-manifest.example.json >/dev/null; then
    echo "Generated release proof artifact must not be a contract manifest subject: $generated_proof_artifact" >&2
    exit 1
  fi
done

expect_fail schema/assessment-input.schema.json examples/contract-foundation/negative-native-policy-field.json
expect_fail schema/trace.schema.json examples/contract-foundation/negative-native-policy-trace.json
expect_fail schema/metric-stream.schema.json examples/contract-foundation/negative-assessed-metric-without-evidence.json
expect_fail schema/accountability.schema.json examples/contract-foundation/negative-ai-sole-accountable-owner.json
expect_fail schema/contract-release-verification.schema.json examples/contract-foundation/negative-unauthorized-signer.json
expect_fail schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-missing-external-evidence.json
expect_fail schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-oidc-issuer-mismatch.json
expect_fail schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-source-uri-mismatch.json
expect_fail schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-protected-ref-mismatch.json
expect_fail schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-workflow-identity-mismatch.json
expect_fail schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-approval-mismatch.json
expect_fail schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-expired-freshness.json
expect_fail schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-local-profile.json
expect_fail schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-mutable-source-commit.json
expect_fail schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-bad-source-counts.json
expect_fail schema/contract-release-verification.schema.json examples/contract-foundation/negative-trusted-release-sigstore-not-required-transparency.json

tmp_release_repo="$(mktemp -d)"
git -C "$tmp_release_repo" init -q
git -C "$tmp_release_repo" config user.email "sdp-trace@example.invalid"
git -C "$tmp_release_repo" config user.name "SDP Trace Test"
printf 'clean artifact\n' >"$tmp_release_repo/artifact.txt"
artifact_sha="$(shasum -a 256 "$tmp_release_repo/artifact.txt" | awk '{print $1}')"
jq -n --arg sha "$artifact_sha" '{
  artifacts: [
    {
      path: "artifact.txt",
      sha256: $sha
    }
  ]
}' >"$tmp_release_repo/manifest.json"
git -C "$tmp_release_repo" add artifact.txt manifest.json
git -C "$tmp_release_repo" commit -q -m "fixture"
tmp_finalization_result="$tmp_release_repo.finalization-result.json"
scripts/finalize-source-bound-release.sh --repo "$tmp_release_repo" --manifest manifest.json --source-ref HEAD >"$tmp_finalization_result"
if [[ "$(jq -r '.proof_states.source_commit_artifacts_verified.value' "$tmp_finalization_result")" != "true" ]]; then
  echo "Expected source-bound finalization to verify clean committed artifact set" >&2
  cat "$tmp_finalization_result" >&2
  exit 1
fi
if [[ "$(jq -r '.trusted_contract_release' "$tmp_finalization_result")" != "false" ]]; then
  echo "Source-bound local finalization must not claim trusted contract release" >&2
  cat "$tmp_finalization_result" >&2
  exit 1
fi
rm -f "$tmp_finalization_result"
printf 'dirty artifact\n' >"$tmp_release_repo/artifact.txt"
dirty_finalization_out="$tmp_validation_dir/dirty-finalization.out"
if scripts/finalize-source-bound-release.sh --repo "$tmp_release_repo" --manifest manifest.json --source-ref HEAD >"$dirty_finalization_out" 2>&1; then
  echo "Expected source-bound finalization to refuse dirty working tree, but passed" >&2
  cat "$dirty_finalization_out" >&2
  exit 1
fi
if ! grep -q "Refusing source-bound finalization from dirty working tree" "$dirty_finalization_out"; then
  echo "Dirty source-bound finalization failed for the wrong reason" >&2
  cat "$dirty_finalization_out" >&2
  exit 1
fi
rm -f "$dirty_finalization_out"

scripts/verify-contract-manifest.sh examples/contract-foundation/contract-manifest.example.json
manifest_negative_out="$tmp_validation_dir/manifest-negative.out"
if scripts/verify-contract-manifest.sh examples/contract-foundation/negative-modified-contract-manifest.json >"$manifest_negative_out" 2>&1; then
  echo "Expected manifest digest verification failure, but passed" >&2
  cat "$manifest_negative_out" >&2
  exit 1
fi
rm -f "$manifest_negative_out"

scripts/verify-release-signature.sh examples/contract-foundation/contract-release.dsse.json examples/contract-foundation/local-dev-signing-public.pem
scripts/validate-self-trace.sh
scripts/verify-artifact-hashes.sh
scripts/verify-self-attestation.sh --all
scripts/check-artifact-safety.sh
scripts/validate-discovery-entrypoints.sh
