package harnessobs

import "sort"

func evaluationDimensions(profile Profile, events []Event) []Dimension {
	// Evaluation dimensions are derived only from observed event families; they
	// do not infer proof beyond the replayed event stream.
	counts := eventFamilyCounts(events)
	dimensions, required := requiredDimensions(profile.RequiredEventFamilies, counts)
	dimensions = append(dimensions, optionalDimensions(profile.OptionalEventFamilies, required, counts)...)

	sortDimensions(dimensions)
	return dimensions
}

func eventFamilyCounts(events []Event) map[string]int {
	// Counts preserve the raw event-family labels so profile validation owns
	// family allow-list decisions.
	counts := map[string]int{}
	for _, event := range events {
		counts[event.EventFamily]++
	}
	return counts
}

func requiredDimensions(families []string, counts map[string]int) ([]Dimension, map[string]bool) {
	// Required dimensions must remain visible as not_assessed when absent; they
	// cannot disappear from the validation output.
	required := map[string]bool{}
	dimensions := make([]Dimension, 0, len(families))
	for _, family := range families {
		required[family] = true
		dimensions = append(dimensions, dimension(family, true, counts[family]))
	}
	return dimensions, required
}

func optionalDimensions(families []string, required map[string]bool, counts map[string]int) []Dimension {
	// Optional dimensions are included only when they are not already required,
	// keeping duplicated profile configuration from duplicating output rows.
	dimensions := []Dimension{}
	for _, family := range families {
		if !required[family] {
			dimensions = append(dimensions, dimension(family, false, counts[family]))
		}
	}
	return dimensions
}

func sortDimensions(dimensions []Dimension) {
	// Required dimensions are shown first, then families are stable-sorted so
	// repeated validations produce deterministic artifacts.
	// The comparator intentionally does not inspect state; state ordering is
	// handled separately by compose.
	sort.Slice(dimensions, func(i, j int) bool {
		if dimensions[i].Required != dimensions[j].Required {
			return dimensions[i].Required
		}
		return dimensions[i].Family < dimensions[j].Family
	})
}
