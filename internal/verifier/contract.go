package verifier

import (
	"encoding/json"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func verifiedContract(manifest trace.RunManifest, result trace.VerifierResult) (trace.Contract, trace.MissingEvidenceTable, trace.VerifierResult, *trace.IntegrityAudit, error, bool) {
	contract := trace.DefaultContract
	if manifest.ContractPath != "" {
		// Explicit contract paths take precedence because they bind run-specific
		// requirements.
		return verifiedManifestContract(manifest, result, contract)
	} else if manifest.ContractDigest != "" {
		if manifest.ContractDigest != contractDigest(trace.DefaultContract) {
			// A default-contract digest mismatch means the verifier's default has
			// drifted since run capture.
			return contract, trace.MissingEvidenceTable{ContractID: manifest.ContractID}, cannotVerifyReplayResult(result, "default contract digest mismatch"), audit(manifest.RunID, "contract_digest_mismatch", "default contract digest changed after run", "", ""), nil, false
		}
	}
	return contract, trace.MissingEvidenceTable{}, result, nil, nil, true
}

func verifiedManifestContract(manifest trace.RunManifest, result trace.VerifierResult, fallback trace.Contract) (trace.Contract, trace.MissingEvidenceTable, trace.VerifierResult, *trace.IntegrityAudit, error, bool) {
	resolvedContract, err := trace.LoadContract(filepath.Clean(manifest.ContractPath))
	if err != nil {
		// Keep the manifest contract id in the missing-evidence table even when
		// the contract file cannot be loaded.
		return fallback, trace.MissingEvidenceTable{ContractID: manifest.ContractID}, cannotVerifyReplayResult(result, err.Error()), audit(manifest.RunID, "contract_unreadable", err.Error(), "contract_path", manifest.ContractPath), err, false
	}
	if manifest.ContractDigest != "" && manifest.ContractDigest != contractDigest(resolvedContract) {
		// Contract digest mismatch blocks replay because requirements may have
		// changed after the run.
		return fallback, trace.MissingEvidenceTable{ContractID: manifest.ContractID}, cannotVerifyReplayResult(result, "contract digest mismatch"), audit(manifest.RunID, "contract_digest_mismatch", "contract digest changed after run", "contract_path", manifest.ContractPath), nil, false
	}
	return resolvedContract, trace.MissingEvidenceTable{}, result, nil, nil, true
}

func contractDigest(contract trace.Contract) string {
	return trace.SHA256Hex(string(mustMarshalJSON(contract)))
}

func mustMarshalJSON(value any) []byte {
	data, err := json.Marshal(value)
	if err != nil {
		// Contract digesting should stay total for impossible marshal failures.
		return []byte("{}")
	}
	return data
}
