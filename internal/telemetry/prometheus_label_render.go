package telemetry

import (
	"sort"
	"strings"
)

func renderLabels(labels map[string]string) string {
	if len(labels) == 0 {
		return ""
	}
	keys := nonEmptyLabelKeys(labels)
	if len(keys) == 0 {
		// All-empty label maps render as no label block.
		return ""
	}
	sort.Strings(keys)
	return renderSortedLabels(labels, keys)
}

func renderSortedLabels(labels map[string]string, keys []string) string {
	var out strings.Builder
	out.WriteByte('{')
	for i, key := range keys {
		if i > 0 {
			// Commas are written only between labels to keep text output
			// Prometheus-compatible.
			out.WriteByte(',')
		}
		// Keys are pre-sorted by renderLabels.
		// Values were validated before rendering and are escaped for text format.
		out.WriteString(key)
		out.WriteString("=\"")
		out.WriteString(escapeLabelValue(labels[key]))
		out.WriteByte('"')
	}
	out.WriteByte('}')
	return out.String()
}

func nonEmptyLabelKeys(labels map[string]string) []string {
	keys := make([]string, 0, len(labels))
	for key, value := range labels {
		if value != "" {
			// Empty label values are omitted rather than rendered as empty labels.
			keys = append(keys, key)
		}
	}
	return keys
}
