#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VALIDATOR="${ROOT_DIR}/scripts/validate-discovery-entrypoints.sh"

if [[ ! -x "$VALIDATOR" ]]; then
  echo "Missing executable validator: $VALIDATOR" >&2
  exit 1
fi

source_agent="${ROOT_DIR}/docs/agent-entrypoint.md"
source_reviewer="${ROOT_DIR}/docs/reviewer-entrypoint.md"
source_readme="${ROOT_DIR}/README.md"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

copy_fixture() {
  local fixture_dir="$1"
  mkdir -p "$fixture_dir/docs"
  cp "$source_agent" "$fixture_dir/docs/agent-entrypoint.md"
  cp "$source_reviewer" "$fixture_dir/docs/reviewer-entrypoint.md"
  cp "$source_readme" "$fixture_dir/README.md"
}

run_validator() {
  local fixture_dir="$1"
  local expected="$2"
  local -a args=(
    --agent-doc "${fixture_dir}/docs/agent-entrypoint.md"
    --reviewer-doc "${fixture_dir}/docs/reviewer-entrypoint.md"
    --readme "${fixture_dir}/README.md"
  )
  local out_file="${fixture_dir}/validator.out"

  set +e
  "$VALIDATOR" "${args[@]}" >"$out_file" 2>&1
  local status=$?
  set -e

  if [[ "$expected" == "success" && "$status" -ne 0 ]]; then
    echo "Expected discovery validation to pass" >&2
    cat "$out_file" >&2
    exit 1
  fi

  if [[ "$expected" == "failure" && "$status" -eq 0 ]]; then
    echo "Expected discovery validation to fail" >&2
    echo "--- output ---" >&2
    cat "$out_file" >&2
    exit 1
  fi
}

positive_dir="$tmp_dir/positive"
copy_fixture "$positive_dir"
run_validator "$positive_dir" success

negative_missing_command="$tmp_dir/missing-command"
copy_fixture "$negative_missing_command"
awk '
  {
    if (!removed && $0 ~ /npm run verify:source-bound/) {
      removed=1
      next
    }
    print
  }
' "$negative_missing_command/docs/agent-entrypoint.md" >"${negative_missing_command}/docs/agent-entrypoint.md.new"
mv "${negative_missing_command}/docs/agent-entrypoint.md.new" \
  "${negative_missing_command}/docs/agent-entrypoint.md"
run_validator "$negative_missing_command" failure

negative_forbidden_phrase="$tmp_dir/forbidden-phrase"
copy_fixture "$negative_forbidden_phrase"
printf '\nThis path is production-ready for external trust.\n' >>"${negative_forbidden_phrase}/docs/reviewer-entrypoint.md"
run_validator "$negative_forbidden_phrase" failure

negative_noncanonical_profile="$tmp_dir/noncanonical-profile"
copy_fixture "$negative_noncanonical_profile"
sed 's/source_bound_local_release/source-bound-local-release/g' \
  "$negative_noncanonical_profile/docs/agent-entrypoint.md" >"${negative_noncanonical_profile}/docs/agent-entrypoint.md.new"
mv "${negative_noncanonical_profile}/docs/agent-entrypoint.md.new" \
  "${negative_noncanonical_profile}/docs/agent-entrypoint.md"
run_validator "$negative_noncanonical_profile" failure

negative_missing_section="$tmp_dir/missing-section-marker"
copy_fixture "$negative_missing_section"
grep -v '^## Commands and `--json`' \
  "$negative_missing_section/docs/reviewer-entrypoint.md" >"${negative_missing_section}/docs/reviewer-entrypoint.md.new"
mv "${negative_missing_section}/docs/reviewer-entrypoint.md.new" \
  "${negative_missing_section}/docs/reviewer-entrypoint.md"
run_validator "$negative_missing_section" failure

negative_script_surface="$tmp_dir/noncanonical-script-surface"
copy_fixture "$negative_script_surface"
sed 's#scripts/verify.sh --profile baseline|source-bound|external-trust \[--json\] \[--allow-dirty\] \[--version\]#scripts/verify.sh --profile baseline --json#' \
  "$negative_script_surface/docs/agent-entrypoint.md" >"${negative_script_surface}/docs/agent-entrypoint.md.new"
mv "${negative_script_surface}/docs/agent-entrypoint.md.new" \
  "${negative_script_surface}/docs/agent-entrypoint.md"
run_validator "$negative_script_surface" failure

echo "discovery entrypoint tests passed"
