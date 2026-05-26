package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// opaPreflight checks whether the installed OPA supports `import rego.v1`.
// If not, it returns `cannot_verify` so that version-mismatch failures are
// not misreported as conformance `fail`.
func opaPreflight(ctx context.Context) (verifierState, string) {
	tmpDir, err := os.MkdirTemp("", "osscompat-opa-preflight-*")
	if err != nil {
		return stateCannotVerify, fmt.Sprintf("mkdir temp for opa preflight: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	// Probe syntax support before running fixtures that require rego.v1.
	preflightPath := filepath.Join(tmpDir, "preflight.rego")
	if err := os.WriteFile(preflightPath, []byte(opaPreflightRego), 0644); err != nil {
		return stateCannotVerify, fmt.Sprintf("write opa preflight: %v", err)
	}
	out, err := runExternalTool(ctx, "opa", "eval", "--data", preflightPath, "--format", "raw", "data.sdp_trace.preflight.allow")
	if err != nil {
		return stateCannotVerify, fmt.Sprintf("opa does not support rego.v1 syntax (version too old?): %v\n%s", err, strings.TrimSpace(string(out)))
	}
	return "", ""
}

const opaPreflightRego = `package sdp_trace.preflight
import rego.v1
allow := true
`
