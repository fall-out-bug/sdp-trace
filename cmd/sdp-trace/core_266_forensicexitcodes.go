package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
)

var forensicExitCodes = map[string]int{
	forensic.StatePass: 0,
	forensic.StateFail: 1,
}
