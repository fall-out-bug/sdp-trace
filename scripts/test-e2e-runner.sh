#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT_DIR"

runner="scripts/run-opencode-minimax-kotlin-bazel-proof.sh"
if [[ ! -x "$runner" ]]; then
  echo "Missing executable runner: $runner" >&2
  exit 1
fi

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT

"$runner" --help >"$tmp/help.out"
for token in --repo --scope --bazel-target --bazel-command --model --out --timeout-seconds; do
  if ! grep -q -- "$token" "$tmp/help.out"; then
    echo "Runner help is missing option: $token" >&2
    exit 1
  fi
done

if "$runner" \
  --repo "$tmp/missing-repo" \
  --scope services/example \
  --bazel-target //services/example:unit_test \
  --bazel-command "bazel test //services/example:unit_test" \
  --model minimax-coding-plan/MiniMax-M2.5 \
  --out "$tmp/out" >"$tmp/missing.out" 2>&1; then
  echo "Runner should reject missing repository" >&2
  exit 1
fi
if ! grep -q "Repository path not found" "$tmp/missing.out"; then
  echo "Runner missing-repo error was not specific" >&2
  cat "$tmp/missing.out" >&2
  exit 1
fi

fixture="$tmp/repo"
mkdir -p "$fixture/services/example"
git -C "$tmp" init -q
if "$runner" \
  --repo "$fixture" \
  --scope services/example \
  --bazel-target //services/example:unit_test \
  --bazel-command "bazel test //services/example:unit_test" \
  --model minimax-coding-plan/MiniMax-M2.5 \
  --out "$tmp/tracked-output" >"$tmp/tracked.out" 2>&1; then
  echo "Runner should reject tracked or non-ignored output path" >&2
  exit 1
fi
if ! grep -q "Output path must be under .sdp-trace-runs" "$tmp/tracked.out"; then
  echo "Runner tracked-output error was not specific" >&2
  cat "$tmp/tracked.out" >&2
  exit 1
fi

bad_command_out=".sdp-trace-runs/test-e2e-runner-bad-command"
rm -rf "$bad_command_out"
if "$runner" \
  --repo "$fixture" \
  --scope services/example \
  --bazel-target //services/example:unit_test \
  --bazel-command "bazel test //services/example:unit_test; echo unsafe" \
  --model minimax-coding-plan/MiniMax-M2.5 \
  --out "$bad_command_out" >"$tmp/bad-command.out" 2>&1; then
  echo "Runner should reject shell-shaped bazel commands" >&2
  exit 1
fi
if ! grep -q "Bazel command must not contain shell metacharacters" "$tmp/bad-command.out"; then
  echo "Runner bad-command error was not specific" >&2
  cat "$tmp/bad-command.out" >&2
  exit 1
fi
rm -rf "$bad_command_out"

fakebin="$tmp/fakebin"
mkdir -p "$fakebin"
cat >"$fakebin/opencode" <<'FAKE_OPENCODE'
#!/usr/bin/env bash
set -euo pipefail
case "${1:-}" in
  --version)
    echo "fake-opencode-1.0"
    ;;
  models)
    echo "minimax-coding-plan/MiniMax-M2.5"
    ;;
  run)
    echo '{"type":"message","content":"fake minimax run"}'
    ;;
  *)
    echo "unexpected fake opencode args: $*" >&2
    exit 2
    ;;
esac
FAKE_OPENCODE
cat >"$fakebin/bazel" <<'FAKE_BAZEL'
#!/usr/bin/env bash
set -euo pipefail
case "${1:-}" in
  query)
    if [[ "${2:-}" == "--output=build" ]]; then
      echo 'genrule(name = "compile_hello_jar", srcs = ["//services/example:Hello.kt"], outs = ["hello.jar"], cmd = "touch $@")'
    else
      echo "$2"
    fi
    ;;
  build)
    echo "fake bazel build complete"
    ;;
  *)
    echo "unexpected fake bazel args: $*" >&2
    exit 2
    ;;
esac
FAKE_BAZEL
chmod +x "$fakebin/opencode" "$fakebin/bazel"

fake_repo="$tmp/fake-repo"
mkdir -p "$fake_repo/services/example"
git -C "$fake_repo" init -q
printf 'module(name = "fake")\n' >"$fake_repo/MODULE.bazel"
printf 'genrule(name = "compile_hello_jar", outs = ["hello.jar"], cmd = "touch $@")\n' >"$fake_repo/services/example/BUILD.bazel"
printf 'fun main() = println("hello")\n' >"$fake_repo/services/example/Hello.kt"
git -C "$fake_repo" add .
git -C "$fake_repo" -c user.email=sdp-trace@example.invalid -c user.name="SDP Trace Test" commit -q -m "fixture"

fake_out=".sdp-trace-runs/test-e2e-runner-fake"
rm -rf "$fake_out"
if ! PATH="$fakebin:$PATH" "$runner" \
  --repo "$fake_repo" \
  --scope services/example \
  --bazel-target //services/example:compile_hello_jar \
  --bazel-command "bazel build //services/example:compile_hello_jar" \
  --model minimax-coding-plan/MiniMax-M2.5 \
  --out "$fake_out" \
  --timeout-seconds 5 >"$tmp/fake-run.out" 2>&1; then
  if [[ ! -f "$fake_out/evidence/proof-states.json" ]]; then
    echo "Runner should emit proof states with fake OpenCode/Bazel tools" >&2
    cat "$tmp/fake-run.out" >&2
    exit 1
  fi
fi
if [[ "$(jq -r '.states[] | select(.name == "kotlin_bazel_target_identified") | .state' "$fake_out/evidence/proof-states.json")" != "observed" ]]; then
  echo "Fake runner package should observe the Kotlin+Bazel target" >&2
  cat "$fake_out/evidence/proof-states.json" >&2
  exit 1
fi
rm -rf "$fake_out"

echo "e2e runner tests passed"
