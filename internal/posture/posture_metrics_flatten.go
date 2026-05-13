package posture

import (
	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func flattenRows(result query.QueryPackResult) []query.QueryPackRow {
	// flattenRows keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	var rows []query.QueryPackRow

	for _, key := range sortedMapKeys(result.QueryRows) {
		rows = append(rows, result.QueryRows[key]...)
	}
	return rows
}
