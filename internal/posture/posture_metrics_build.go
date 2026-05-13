package posture

func buildMetrics(groups map[string]*aggregateGroup) []MetricRow {
	// Metric rows are emitted in stable group/catalog order because row ids are
	// part of downstream movement and explanation references.
	var rows []MetricRow
	counter := 0
	for _, groupKey := range sortedMapKeys(groups) {
		group := groups[groupKey]
		for _, def := range metricCatalog {
			counter++
			row := metricForGroup(counter, def, group)
			rows = append(rows, row)
		}
	}
	return rows
}
