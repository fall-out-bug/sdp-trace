package posture

import (
	"fmt"
)

func metricForGroup(counter int, def metricDef, group *aggregateGroup) MetricRow {
	// metricForGroup keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	counts := metricCounts(def, group)
	sourceRefs, digestHash := metricEvidenceRefs(group)

	row := metricRowHeader(counter, def, group)
	row.Numerator = counts.numerator
	row.Denominator = len(group.rows)
	row.SourceInputRefs = sourceRefs
	row.SourceArtifactDigestSet = digestHash
	row.SourceFieldState = counts.sourceFieldState
	row.NotAssessedCount = counts.notAssessed
	row.InputTrustStateSummary = copyTrust(group.trustStates)
	return row
}

func metricRowHeader(counter int, def metricDef, group *aggregateGroup) MetricRow {
	// metricRowHeader keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	return MetricRow{
		ID:            fmt.Sprintf("metric.%04d", counter),
		MetricID:      def.id,
		MetricVersion: def.version,
		Unit:          "rows",
		TimeWindow:    group.window,
		Dimensions:    group.dimensions,
		DimensionKey:  group.dimensionKey,
	}
}
