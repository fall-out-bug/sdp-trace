package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
)

var ciArtifactExitCodes = map[string]int{
	ciartifact.StatePass: 0,
	ciartifact.StateFail: 1,
}
