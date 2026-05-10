package recorder

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
	"github.com/fall_out_bug/sdp-trace/internal/verifier"
)

func TestRunCreatesEventChainAndPreservesExitCode(t *testing.T) {
	sh := mustFindShell(t)
	tempDir := t.TempDir()
	runDir := filepath.Join(tempDir, "run-wrap-positive")
	ctx := context.Background()

	result, err := Run(ctx, RecorderOptions{
		WrapperName:        "test",
		UseDefaultContract: true,
		OutputDir:          runDir,
		Command:            []string{sh, "-c", "printf 'hello world\\n'"},
	})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("exit code = %d, expected 0", result.ExitCode)
	}
	verify, _, _, err := verifier.VerifyRun(result.RunDir)
	if err != nil {
		t.Fatalf("verify returned error: %v", err)
	}
	if verify.Result != trace.VerdictObserved {
		t.Fatalf("unexpected result: %+v", verify.Result)
	}
	eventsDir := filepath.Join(result.RunDir, "events")
	if _, err := os.Stat(filepath.Join(eventsDir, "000000-recorder_attached.json")); err != nil {
		t.Fatalf("missing first event: %v", err)
	}
}

func TestRunPreservesNonZeroExitCode(t *testing.T) {
	sh := mustFindShell(t)
	tempDir := t.TempDir()
	runDir := filepath.Join(tempDir, "run-wrap-nonzero")
	result, err := Run(context.Background(), RecorderOptions{
		WrapperName:        "test",
		UseDefaultContract: true,
		OutputDir:          runDir,
		Command:            []string{sh, "-c", "exit 11"},
	})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if result.ExitCode != 11 {
		t.Fatalf("exit code = %d, expected 11", result.ExitCode)
	}
}

func TestRunCommandHashesOutputAndPreservesExitState(t *testing.T) {
	sh := mustFindShell(t)
	writer := newTestRunWriter(t)

	exitCode, signal := runCommand(context.Background(), []string{sh, "-c", "printf stdout; printf stderr >&2; exit 7"}, os.Environ(), writer)

	if exitCode != 7 {
		t.Fatalf("exit code = %d, expected 7", exitCode)
	}
	if signal != "" {
		t.Fatalf("signal = %q, expected empty", signal)
	}
	if writer.stdoutDigest() != sha256Hex("stdout") {
		t.Fatalf("stdout digest = %q, expected hash of stdout", writer.stdoutDigest())
	}
	if writer.stderrDigest() != sha256Hex("stderr") {
		t.Fatalf("stderr digest = %q, expected hash of stderr", writer.stderrDigest())
	}
}

func TestRunCommandReportsStartFailure(t *testing.T) {
	writer := newTestRunWriter(t)

	exitCode, signal := runCommand(context.Background(), []string{"/path/that/does/not/exist"}, os.Environ(), writer)

	if exitCode != 1 {
		t.Fatalf("exit code = %d, expected 1", exitCode)
	}
	if signal != "start_failed" {
		t.Fatalf("signal = %q, expected start_failed", signal)
	}
}

func newTestRunWriter(t *testing.T) *runWriter {
	t.Helper()
	writer, err := newRunWriter(t.TempDir(), trace.DefaultContract, "test")
	if err != nil {
		t.Fatalf("newRunWriter error: %v", err)
	}
	return writer
}

func sha256Hex(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

func mustFindShell(t *testing.T) string {
	t.Helper()
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("sh shell not available for recorder tests")
	}
	return sh
}
