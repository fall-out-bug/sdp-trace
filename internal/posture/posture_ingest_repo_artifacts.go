package posture

import (
	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func ingestRepositoryArtifacts(repo RepositoryWindow) repositoryIngest {
	// ingestRepositoryArtifacts keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	digest, err := verifyDigestManifest(repo.ArtifactDigestManifest, repo.QueryPackResult)
	if err != nil {
		return refusedInput(reasonForDigestErr(err), trustForDigestErr(err), "", true)
	}
	result, err := readQueryPack(repo.QueryPackResult)
	if err != nil {
		return refusedInput("malformed_input", "cannot_verify_input", digest, true)
	}
	return ingestSupportedRepositoryArtifacts(repo, digest, result)
}

func ingestSupportedRepositoryArtifacts(repo RepositoryWindow, digest string, result query.QueryPackResult) repositoryIngest {
	// ingestSupportedRepositoryArtifacts keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	if !isSupportedQueryPack(result) {
		return refusedInput("malformed_input", "cannot_verify_input", digest, true)
	}

	signals, err := readSignals(repo.PostureSignalManifest)
	if err != nil {
		return refusedInput("malformed_input", "cannot_verify_input", digest, true)
	}
	return repositoryIngest{trusted: true, digest: digest, result: result, signals: signals}
}

func isSupportedQueryPack(result query.QueryPackResult) bool {
	return result.SchemaVersion == query.QueryPackSchemaVersion && result.QueryPackID == query.QueryPackForensicsBasic
}

// refusedInput crosses the trust boundary from ingestion to refusal record.
// Refused inputs retain a trust state but do not contribute to metric evidence.
func refusedInput(reason, trustState, digest string, recordSelection bool) repositoryIngest {
	return repositoryIngest{refusalReason: reason, inputTrustState: trustState, digest: digest, recordSelection: recordSelection}
}
