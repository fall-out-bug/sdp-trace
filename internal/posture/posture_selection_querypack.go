package posture

import (
	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func readQueryPack(path string) (query.QueryPackResult, error) {
	return readJSONFile[query.QueryPackResult](path)
}
