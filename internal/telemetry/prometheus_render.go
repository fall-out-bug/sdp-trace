package telemetry

import (
	"strconv"
	"strings"
)

func renderPrometheusSeries(series []Series) string {
	var out strings.Builder
	currentFamily := ""
	for _, item := range series {
		if item.Name != currentFamily {
			// HELP/TYPE headers are emitted once per metric family.
			writePrometheusFamilyHeader(&out, item)
			currentFamily = item.Name
		}
		// Samples are rendered only after all series have passed label safety and
		// duplicate checks.
		writePrometheusSample(&out, item)
	}
	out.WriteString("# EOF\n")
	return out.String()
}

func writePrometheusFamilyHeader(out *strings.Builder, item Series) {
	// Help and type strings are fixed by the exporter, not sourced from labels.
	// Prometheus metadata therefore cannot carry posture input content.
	out.WriteString("# HELP ")
	out.WriteString(item.Name)
	out.WriteByte(' ')
	out.WriteString(item.Help)
	out.WriteByte('\n')
	out.WriteString("# TYPE ")
	out.WriteString(item.Name)
	out.WriteByte(' ')
	out.WriteString(item.Type)
	out.WriteByte('\n')
}

func writePrometheusSample(out *strings.Builder, item Series) {
	// Values are rendered with Go's shortest decimal form to avoid synthetic
	// precision changes in diffs.
	out.WriteString(item.Name)
	out.WriteString(renderLabels(item.Labels))
	out.WriteByte(' ')
	out.WriteString(strconv.FormatFloat(item.Value, 'f', -1, 64))
	out.WriteByte('\n')
}
