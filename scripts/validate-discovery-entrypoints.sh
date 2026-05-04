#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
AGENT_ENTRYPOINT_PATH="${ROOT_DIR}/docs/agent-entrypoint.md"
REVIEWER_ENTRYPOINT_PATH="${ROOT_DIR}/docs/reviewer-entrypoint.md"
README_PATH="${ROOT_DIR}/README.md"

usage() {
  cat >&2 <<'USAGE'
Usage: scripts/validate-discovery-entrypoints.sh \
  [--agent-doc <path>] \
  [--reviewer-doc <path>] \
  [--readme <path>]
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --agent-doc)
      AGENT_ENTRYPOINT_PATH="$2"
      shift 2
      ;;
    --reviewer-doc)
      REVIEWER_ENTRYPOINT_PATH="$2"
      shift 2
      ;;
    --readme)
      README_PATH="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown argument: $1" >&2
      usage
      exit 1
      ;;
  esac
done

require_file() {
  local file="$1"
  local label="$2"
  if [[ ! -f "$file" ]]; then
    echo "Missing required file for ${label}: $file" >&2
    exit 1
  fi
}

require_pattern() {
  local file="$1"
  local pattern="$2"
  local label="$3"
  if ! grep -qF -- "$pattern" "$file"; then
    echo "Expected '${pattern}' in ${label}: $file" >&2
    exit 1
  fi
}

has_heading() {
  local file="$1"
  local expected="$2"

  awk -v expected="$expected" '
    $0 ~ /^[[:space:]]*#{1,6}[[:space:]]*/ {
      heading=$0
      sub(/^[[:space:]]*#{1,6}[[:space:]]*/, "", heading)
      sub(/^[0-9]+\)[[:space:]]*/, "", heading)
      sub(/[[:space:]]+$/, "", heading)
      if (heading == expected) {
        found=1
        exit
      }
    }
    END {
      exit !found
    }
  ' "$file"
}

require_section_marker() {
  local file="$1"
  local label="$2"
  local section="$3"
  if ! has_heading "$file" "$section"; then
    echo "Missing required section '${section}' in ${label}: $file" >&2
    exit 1
  fi
}

require_section_markers() {
  local file="$1"
  local label="$2"
  shift 2
  local marker

  for marker in "$@"; do
    require_section_marker "$file" "$label" "$marker"
  done
}

require_not_found() {
  local file="$1"
  local pattern="$2"
  local label="$3"
  if grep -qiF -- "$pattern" "$file"; then
    echo "Forbidden phrase found in ${label}: '${pattern}'" >&2
    echo "  file: $file" >&2
    exit 1
  fi
}

require_dirty_guidance() {
  local file="$1"
  local label="$2"
  if ! grep -qiE -- '(dirty[[:space:]-]*checkout|--allow-dirty)' "$file"; then
    echo "Expected dirty-checkout guidance in ${label}: $file" >&2
    exit 1
  fi
}

require_readme_link() {
  local doc_key="$1"
  local label="$2"
  local required_text="$3"
  local required_alias="$4"
  local link_line
  local line_text

  link_line="$(grep -nF -- "($doc_key)" "$README_PATH" | head -n 1 || true)"
  if [[ -z "$link_line" ]]; then
    link_line="$(grep -nF -- "$doc_key" "$README_PATH" | head -n 1 || true)"
  fi
  if [[ -z "$link_line" ]]; then
    echo "Missing README link to ${label} doc: $doc_key" >&2
    exit 1
  fi
  if ! echo "$link_line" | grep -qiF -- "$required_text" || ! echo "$link_line" | grep -qiF -- "$required_alias"; then
    echo "README link text for ${label} doc must include '${required_text}' and '${required_alias}': $doc_key" >&2
    echo "  line: $link_line" >&2
    exit 1
  fi

  # Only validate forbidden phrases on link text lines to avoid scanning unrelated README prose.
  line_text="${link_line#*: }"
  for phrase in "${FORBIDDEN_PHRASES[@]}"; do
    if echo "$line_text" | grep -qiF -- "$phrase"; then
      echo "Forbidden phrase found in README link text for ${label}: '${phrase}'" >&2
      echo "  line: $link_line" >&2
      exit 1
    fi
  done
}

require_profile_vocabulary() {
  local file="$1"
  local label="$2"

  if ! grep -qF -- "local_dirty_structural_only" "$file"; then
    echo "Expected local_dirty_structural_only vocabulary in ${label}: $file" >&2
    exit 1
  fi
}

require_external_trust_block_vocab() {
  local file="$1"
  local label="$2"

  if ! grep -qF -- "external_trust_profile_selected: fail" "$file"; then
    echo "Expected explicit external trust blocker wording in ${label}: $file" >&2
    echo "  Add external_trust_profile_selected: fail" >&2
    exit 1
  fi
}

AGENT_SECTION_MARKERS=(
  "Profile Selection"
  "Command Contract"
  "Evidence Emission Rules"
  "Forbidden Claims"
)

REVIEWER_SECTION_MARKERS=(
  "Verification Path"
  "Dirty Checkout Behavior"
  "Not-Assessed Interpretation"
  "External Trust Gap"
  "What You May State from Output"
  'Commands and `--json`'
  "Quick Reference"
)

DISCOVERY_COMMANDS=(
  "npm run verify:baseline"
  "npm run verify:source-bound"
  "npm run verify:external-trust"
)

SCRIPT_FORM_COMMANDS=(
  "scripts/verify.sh --profile baseline|source-bound|external-trust [--json] [--allow-dirty] [--version]"
)

PROFILE_MARKERS=(
  "repo_baseline_structural"
  "source_bound_local_release"
  "external_production_trust"
)

FORBIDDEN_PHRASES=(
  "ready for production"
  "production-ready"
  "trusted release"
  "ready for customer"
)

require_file "$AGENT_ENTRYPOINT_PATH" "agent entrypoint doc"
require_file "$REVIEWER_ENTRYPOINT_PATH" "reviewer entrypoint doc"
require_file "$README_PATH" "README"

require_section_markers "$AGENT_ENTRYPOINT_PATH" "agent entrypoint doc" "${AGENT_SECTION_MARKERS[@]}"
require_section_markers "$REVIEWER_ENTRYPOINT_PATH" "reviewer entrypoint doc" "${REVIEWER_SECTION_MARKERS[@]}"

for command in "${DISCOVERY_COMMANDS[@]}"; do
  require_pattern "$AGENT_ENTRYPOINT_PATH" "$command" "agent entrypoint doc"
  require_pattern "$REVIEWER_ENTRYPOINT_PATH" "$command" "reviewer entrypoint doc"
done

for command in "${SCRIPT_FORM_COMMANDS[@]}"; do
  require_pattern "$AGENT_ENTRYPOINT_PATH" "$command" "agent entrypoint doc"
  require_pattern "$REVIEWER_ENTRYPOINT_PATH" "$command" "reviewer entrypoint doc"
done

for profile in "${PROFILE_MARKERS[@]}"; do
  require_pattern "$AGENT_ENTRYPOINT_PATH" "$profile" "agent entrypoint doc"
  require_pattern "$REVIEWER_ENTRYPOINT_PATH" "$profile" "reviewer entrypoint doc"
done

require_profile_vocabulary "$AGENT_ENTRYPOINT_PATH" "agent entrypoint doc"
require_profile_vocabulary "$REVIEWER_ENTRYPOINT_PATH" "reviewer entrypoint doc"
require_external_trust_block_vocab "$REVIEWER_ENTRYPOINT_PATH" "reviewer entrypoint doc"

require_dirty_guidance "$AGENT_ENTRYPOINT_PATH" "agent entrypoint doc"
require_pattern "$AGENT_ENTRYPOINT_PATH" "not_assessed" "agent entrypoint doc"
require_pattern "$REVIEWER_ENTRYPOINT_PATH" "--allow-dirty" "reviewer entrypoint doc"
require_dirty_guidance "$REVIEWER_ENTRYPOINT_PATH" "reviewer entrypoint doc"
require_pattern "$REVIEWER_ENTRYPOINT_PATH" "not_assessed" "reviewer entrypoint doc"

for phrase in "${FORBIDDEN_PHRASES[@]}"; do
  require_not_found "$AGENT_ENTRYPOINT_PATH" "$phrase" "agent entrypoint doc"
  require_not_found "$REVIEWER_ENTRYPOINT_PATH" "$phrase" "reviewer entrypoint doc"
done

require_readme_link "docs/agent-entrypoint.md" "agent" "Agent entrypoint" "current verifier contract"
require_readme_link "docs/reviewer-entrypoint.md" "reviewer" "Reviewer entrypoint" "current proof scope"
