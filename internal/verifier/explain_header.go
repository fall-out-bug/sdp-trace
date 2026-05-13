package verifier

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func explainHeaderLines(runDir string, manifest trace.RunManifest, verification trace.VerifierResult) []string {
	// Header rows identify the live replay subject before optional evidence
	// sections add detail.
	return []string{
		// run_dir is local context; run_id and contract_id come from the
		// manifest, while result comes from fresh verifier replay.
		fmt.Sprintf("run_dir: %s", runDir),
		fmt.Sprintf("run_id: %s", manifest.RunID),
		fmt.Sprintf("contract_id: %s", manifest.ContractID),
		fmt.Sprintf("result: %s", verification.Result),
	}
}

func appendClosureState(lines []string, manifest trace.RunManifest) []string {
	if manifest.ClosureState == "" {
		// Older manifests may not carry closure state; omit the row instead of
		// inventing a verdict.
		return lines
	}
	return append(lines, fmt.Sprintf("closure_state: %s", manifest.ClosureState))
}
