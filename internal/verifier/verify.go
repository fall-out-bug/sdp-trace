package verifier

import "github.com/fall_out_bug/sdp-trace/internal/trace"

// VerifyRun reads a run directory and returns a verifiable result object.
func VerifyRun(runDir string) (trace.VerifierResult, trace.MissingEvidenceTable, *trace.IntegrityAudit, error) {
	// Verification is a live replay boundary: every returned verdict below is
	// produced from the manifest, retained events, and contract available in
	// this run directory now.
	// Checked-in proof JSON never upgrades the outcome; missing or contradictory
	// source evidence keeps the result at cannot_verify, fail, or not_assessed.
	manifest, events, result, audit, err, ok := verifiedRunInput(runDir)
	if !ok {
		// Input failures return the verifier result and optional integrity audit
		// that explain why replay could not continue.
		return result, trace.MissingEvidenceTable{}, audit, err
	}

	// Observation is provisional until chain, contract, and missing-evidence
	// checks complete.
	result = observedResult(runDir, manifest.RunID)
	chainResult, chainAudit, ok := verifiedChain(runDir, manifest, events)
	if !ok {
		return chainResult, trace.MissingEvidenceTable{}, chainAudit, nil
	}

	contract, table, result, audit, err, ok := verifiedContract(manifest, result)
	if !ok {
		return result, table, audit, err
	}
	result, missingEvidence := resultWithMissingEvidence(result, contract, events)
	return result, missingEvidence, nil, nil
}
