package query

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/recorder"
)

func TestMissingEvidenceQueryMatchesVerifier(t *testing.T) {
	tempDir := t.TempDir()
	contractPath := filepath.Join(tempDir, "contract.json")
	contract := map[string]any{
		"contract_id":     "query-missing",
		"version":         "v1",
		"required_events": []string{"recorder_attached", "run_started", "command_started", "command_finished", "run_closed", "file_mutation_observed"},
	}
	contractJSON, _ := json.Marshal(contract)
	if err := osWrite(contractPath, contractJSON); err != nil {
		t.Fatalf("write contract: %v", err)
	}

	runDir := runAndCaptureQuery(t, []string{"q"}, contractPath)
	payload, err := MissingEvidence(runDir)
	if err != nil {
		t.Fatalf("query missing evidence: %v", err)
	}
	var table map[string]any
	if err := json.Unmarshal(payload, &table); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if _, ok := table["rows"]; !ok {
		t.Fatalf("expected rows in payload")
	}
}

func runAndCaptureQuery(t *testing.T, command []string, contract string) string {
	t.Helper()
	echo := mustFindCommandQuery(t, "echo")
	ctx := context.Background()
	result, err := recorder.Run(ctx, recorder.RecorderOptions{
		WrapperName:        "query",
		UseDefaultContract: false,
		ContractPath:       contract,
		OutputDir:          t.TempDir(),
		Command:            append([]string{echo}, command...),
	})
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	return result.RunDir
}

func mustFindCommandQuery(t *testing.T, name string) string {
	t.Helper()
	path, err := execLookPath(name)
	if err != nil {
		t.Skipf("%s not available", name)
	}
	return path
}

func execLookPath(file string) (string, error) {
	return exec.LookPath(file)
}

func osWrite(path string, value []byte) error {
	return os.WriteFile(path, value, 0o644)
}
