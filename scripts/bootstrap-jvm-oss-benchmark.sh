#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MANIFEST="$ROOT/benchmarks/jvm-oss/projects.json"
REPO_DIR="$ROOT/benchmarks/repos/jvm-oss"
RUN_DIR="$ROOT/.sdp-trace-runs/bootstrap/jvm-oss"

usage() {
  cat <<'USAGE'
Usage: scripts/bootstrap-jvm-oss-benchmark.sh [--project <id>] [--dry-run]

Creates shallow, blobless local checkouts for the JVM OSS benchmark and writes
smoke-detection notes under .sdp-trace-runs/bootstrap/jvm-oss/.

Options:
  --project <id>  Bootstrap only one project from benchmarks/jvm-oss/projects.json.
  --dry-run       Print planned clone/update actions without changing repos.
USAGE
}

project_filter=""
dry_run=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --project)
      project_filter="${2:-}"
      if [[ -z "$project_filter" ]]; then
        echo "missing value for --project" >&2
        exit 2
      fi
      shift 2
      ;;
    --dry-run)
      dry_run=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage >&2
      exit 2
      ;;
  esac
done

require() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "required command not found: $1" >&2
    exit 1
  fi
}

require git
require jq

mkdir -p "$REPO_DIR" "$RUN_DIR"

project_query='.projects[]'
if [[ -n "$project_filter" ]]; then
  project_query=".projects[] | select(.id == \"$project_filter\")"
fi

count="$(jq -r "[$project_query] | length" "$MANIFEST")"
if [[ "$count" == "0" ]]; then
  echo "no projects matched filter: ${project_filter:-<none>}" >&2
  exit 1
fi

jq -c "$project_query" "$MANIFEST" | while IFS= read -r project; do
  id="$(jq -r '.id' <<<"$project")"
  repo="$(jq -r '.repo' <<<"$project")"
  expected_build="$(jq -r '.build_system' <<<"$project")"
  target="$REPO_DIR/$id"
  notes="$RUN_DIR/$id.md"

  if [[ "$dry_run" == "1" ]]; then
    if [[ -d "$target/.git" ]]; then
      echo "would update $id at $target"
    else
      echo "would clone $repo -> $target"
    fi
    continue
  fi

  if [[ -d "$target/.git" ]]; then
    git -C "$target" fetch --depth=1 origin
    git -C "$target" checkout --detach FETCH_HEAD >/dev/null
  else
    git clone --depth=1 --filter=blob:none "$repo" "$target"
  fi

  commit="$(git -C "$target" rev-parse HEAD)"
  branch="$(git -C "$target" branch --show-current || true)"
  [[ -n "$branch" ]] || branch="detached"

  gradle="no"
  maven="no"
  bazel="no"
  [[ -f "$target/settings.gradle" || -f "$target/settings.gradle.kts" || -f "$target/build.gradle" || -f "$target/build.gradle.kts" ]] && gradle="yes"
  [[ -f "$target/pom.xml" ]] && maven="yes"
  [[ -f "$target/WORKSPACE" || -f "$target/WORKSPACE.bazel" || -f "$target/MODULE.bazel" ]] && bazel="yes"

  {
    echo "# Bootstrap Smoke: $id"
    echo
    echo "| Field | Value |"
    echo "|---|---|"
    echo "| Repo | $repo |"
    echo "| Commit | $commit |"
    echo "| Branch | $branch |"
    echo "| Expected build | $expected_build |"
    echo "| Gradle detected | $gradle |"
    echo "| Maven detected | $maven |"
    echo "| Bazel detected | $bazel |"
    echo
    echo "## Top-Level Files"
    echo
    find "$target" -maxdepth 1 -type f -print | sed "s#^$target/##" | sort | head -80
  } > "$notes"

  echo "bootstrapped $id -> $target"
  echo "wrote $notes"
done
