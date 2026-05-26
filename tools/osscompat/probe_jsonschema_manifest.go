package main

import (
	"context"
	"os"
	"time"
)

// runJSONSchemaWrapManifest builds sdp-trace from source, runs wrap in an isolated
// temp directory, and validates the generated run.json against the schema for
// the current live recorder manifest.
func runJSONSchemaWrapManifest() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	root := repoRoot()
	// Build from the current checkout so probe output tracks local source.
	// This intentionally avoids checked-in run.json fixtures for the final
	// verdict; fixtures are only a validator/schema preflight below.
	bin, tmpDir, err := buildSDPTraceInTemp(ctx, root)
	if tmpDir != "" {
		defer os.RemoveAll(tmpDir)
	}
	if err != nil {
		return stateCannotVerify, err.Error()
	}
	return checkWrapManifest(ctx, root, bin, tmpDir)
}

func checkWrapManifest(ctx context.Context, root, bin, tmpDir string) (verifierState, string) {
	// The wrap step is separated from schema validation so command/runtime
	// failures remain distinct from manifest conformance failures. Missing
	// check-jsonschema is not a conformance result; it is explicitly
	// cannot_verify/not_assessed depending on the probe runner path.
	runDir, s, r := runWrapAndParse(ctx, bin, tmpDir)
	if s != "" {
		return s, r
	}
	if !hasTool("check-jsonschema") {
		return stateCannotVerify, "check-jsonschema not in PATH; cannot validate live wrap output"
	}
	return runSchemaValidation(ctx, root, tmpDir, runDir)
}
