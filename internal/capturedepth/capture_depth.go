package capturedepth

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
)

const QueryName = "capture-depth"

// Capture-depth output is investigative evidence, not a gate result.
// The query deliberately preserves cannot_verify/not_assessed details instead
// of collapsing them into pass/fail authority.

func CaptureDepth(runDir string) ([]byte, error) {
	// The evaluator supplies condition semantics, but this query never applies
	// assessment exit semantics. Callers that need a verdict must use assess.
	run, err := loadRun(runDir)
	if err != nil {
		return nil, err
	}
	result := adaptercapture.Evaluate(adaptercapture.Input{Run: run})
	summary := newCaptureDepthSummary(run, result)
	return json.MarshalIndent(summary, "", "  ")
}

func loadRun(runDir string) (adaptercapture.RunEvidence, error) {
	// Capture-depth reads the run manifest only; same-chain or bundle evidence
	// must already be represented in the adapter-capture run evidence.
	// runDir is the caller-selected local evidence root, not a remote input.
	var run adaptercapture.RunEvidence
	path := filepath.Join(runDir, "run.json")
	// A missing or malformed manifest is a query error. It is not converted into
	// a gate state because query mode is outside the verifier verdict surface.
	data, err := os.ReadFile(path) // #nosec G304 -- runDir is a caller-selected local evidence root.
	if err != nil {
		return run, err
	}
	err = json.Unmarshal(data, &run)
	return run, err
}
