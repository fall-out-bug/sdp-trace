package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/releaseproof"
)

var releaseProofExitCodes = map[string]int{
	releaseproof.StatePass: 0,
	releaseproof.StateFail: 1,
}
