package posture

import (
	"github.com/fall_out_bug/sdp-trace/internal/query"
)

type metricDef struct {
	id      string
	version string
	source  string
}

type aggregateGroup struct {
	dimensions   map[string]string
	dimensionKey string
	window       string
	rows         []query.QueryPackRow
	inputRefs    []string
	digests      []string
	trustStates  map[string]int
	signals      map[string]PostureSignal
}
