package main

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func runSchemaValidation(ctx context.Context, root, tmpDir, runDir string) (verifierState, string) {
	checkCtx, checkCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer checkCancel()
	// The checked fixture exercises the external validator before the generated
	// live manifest is treated as conformance evidence.
	schemaPath := filepath.Join(root, "schema/run-manifest.schema.json")
	positiveFixture := filepath.Join(root, "examples/agentic-sdlc/local-wrap-positive/run.json")
	if _, err := runCheckJSONSchema(checkCtx, schemaPath, positiveFixture); err != nil {
		return stateCannotVerify, fmt.Sprintf("check-jsonschema preflight failed on known-positive fixture (harness/tool error, not conformance): %v", err)
	}
	runJSONPath := filepath.Join(tmpDir, filepath.Clean(runDir), "run.json")
	out, err := runCheckJSONSchema(checkCtx, schemaPath, runJSONPath)
	return interpretSchemaCheckResult(out, err)
}

func interpretSchemaCheckResult(out []byte, err error) (verifierState, string) {
	if err == nil {
		return statePass, "live wrap run.json validates against run-manifest schema"
	}
	// Exit code 1 is a schema verdict; other exits are validator/runtime errors.
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return stateFail, fmt.Sprintf("live wrap run.json fails run-manifest schema: %s", strings.TrimSpace(string(out)))
	}
	return stateCannotVerify, fmt.Sprintf("check-jsonschema exited abnormally on wrap run.json (harness/tool error, not conformance): %v\n%s", err, strings.TrimSpace(string(out)))
}
