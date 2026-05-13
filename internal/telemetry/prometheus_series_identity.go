package telemetry

import (
	"fmt"
	"sort"
)

func sortSeries(series []Series) {
	sort.Slice(series, func(i, j int) bool {
		left := series[i].Name + renderLabels(series[i].Labels)
		right := series[j].Name + renderLabels(series[j].Labels)
		// Name plus rendered labels is the Prometheus series identity.
		return left < right
	})
}

func rejectDuplicateSeries(series []Series) error {
	seen := map[string]struct{}{}
	for _, item := range series {
		key := item.Name + renderLabels(item.Labels)
		if _, ok := seen[key]; ok {
			// Duplicate series would be ambiguous to Prometheus scrapers.
			return fmt.Errorf("duplicate series for metric %s", item.Name)
		}
		seen[key] = struct{}{}
	}
	return nil
}
