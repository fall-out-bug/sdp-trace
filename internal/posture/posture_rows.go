package posture

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func inputSelection(repo RepositoryWindow, digest, trust string) InputSelection {
	// inputSelection keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	return InputSelection{
		InputID:         repo.InputID,
		Repository:      repo.Repo,
		TimeWindow:      repo.TimeWindow,
		PathRedactedID:  "artifact:query_pack_result:" + shortDigest(digest),
		SHA256:          digest,
		InputTrustState: trust,
	}
}

func refusal(counter int, repo RepositoryWindow, reason, trust string) RefusalRow {
	// refusal keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	return RefusalRow{
		ID:              fmt.Sprintf("refusal.%04d", counter),
		InputID:         repo.InputID,
		TimeWindow:      repo.TimeWindow,
		RefusalReason:   reason,
		InputTrustState: trust,
	}
}

func digestSetHash(digests []string) string {
	sum := sha256.Sum256([]byte(strings.Join(digests, "\n")))
	return hex.EncodeToString(sum[:])
}

// deterministicExportID produces identifiers stable across runs. These are identifiers,
// not proof digests; they do not carry evidence authority.
func deterministicExportID(selectionPath string, metrics []MetricRow, refusals []RefusalRow) string {
	// deterministicExportID keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	sum := sha256.Sum256([]byte(fmt.Sprintf("%s:%d:%d", selectionPath, len(metrics), len(refusals))))
	return "export:" + hex.EncodeToString(sum[:8])
}

// shortDigest produces a stable redacted identifier. When crossing from verified digest
// to missing/empty, the evidence boundary requires a placeholder distinct from a real hash.
func shortDigest(digest string) string {
	// shortDigest keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if len(digest) >= 16 {
		return digest[:16]
	}
	return "not_assessed0000"
}

func copyTrust(in map[string]int) map[string]int {
	// copyTrust keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	out := map[string]int{}
	for key, value := range in {
		out[key] = value
	}
	return out
}
