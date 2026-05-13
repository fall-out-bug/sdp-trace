package posture

import (
	"github.com/fall_out_bug/sdp-trace/internal/query"
)

type repositoryIngest struct {
	trusted         bool
	recordSelection bool
	digest          string
	refusalReason   string
	inputTrustState string
	result          query.QueryPackResult
	signals         map[string]PostureSignal
}
