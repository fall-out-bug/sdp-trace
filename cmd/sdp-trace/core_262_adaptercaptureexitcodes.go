package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
)

var adapterCaptureExitCodes = map[string]int{
	adaptercapture.StatePass: 0,
	adaptercapture.StateFail: 1,
}
