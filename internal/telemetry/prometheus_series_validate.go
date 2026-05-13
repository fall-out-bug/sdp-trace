package telemetry

import "fmt"

const MaxSeries = 10000

func finalizeSeries(series []Series) ([]Series, error) {
	for _, item := range series {
		if err := validateLabels(item.Labels); err != nil {
			// Reject unsafe labels before rendering any partial output.
			return nil, err
		}
	}
	return checkedSeries(series)
}

func checkedSeries(series []Series) ([]Series, error) {
	if err := rejectDuplicateSeries(series); err != nil {
		return nil, err
	}
	if len(series) > MaxSeries {
		// Bound exported cardinality so one posture file cannot create an
		// unbounded scrape surface.
		return nil, fmt.Errorf("series limit exceeded")
	}
	return series, nil
}
