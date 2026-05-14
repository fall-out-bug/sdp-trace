package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func readOptionalPRRoute(path string) (packet.GitHubPREvidenceInput, error) {
	var route packet.GitHubPREvidenceInput
	return route, readOptionalJSON(path, &route)
}
