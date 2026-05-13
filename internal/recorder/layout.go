package recorder

import (
	"errors"
	"fmt"
	"os"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func prepareRunDir(outputDir string) (string, error) {
	// Caller-provided directories must already be empty; generated run
	// directories are isolated under the default local recorder root.
	if outputDir != "" {
		return outputDir, ensureFreshOutputDir(outputDir)
	}
	if err := os.MkdirAll(defaultBaseOutputDir, 0o755); err != nil {
		return "", err
	}
	return os.MkdirTemp(defaultBaseOutputDir, "run-")
}

func resolveContract(contractPath string, useDefault bool) (trace.Contract, error) {
	// Default contracts are explicit opt-ins so local recordings do not silently
	// claim conformance to an unspecified verifier contract.
	if contractPath == "" {
		if useDefault {
			return trace.DefaultContract, nil
		}
		return trace.Contract{}, errors.New("contract required unless --use-default-contract is set")
	}
	return trace.LoadContract(contractPath)
}

func ensureFreshOutputDir(runDir string) error {
	// Existing empty directories are allowed for deterministic tests and callers
	// that allocate paths before invoking the recorder.
	entries, err := os.ReadDir(runDir)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return fmt.Errorf("output directory must be empty: %s", runDir)
	}
	return nil
}
