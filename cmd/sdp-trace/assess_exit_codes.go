package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
	"github.com/fall_out_bug/sdp-trace/internal/managed"
)

var adapterCaptureExitCodes = map[string]int{
	adaptercapture.StatePass: 0,
	adaptercapture.StateFail: 1,
}

var managedExitCodes = map[string]int{
	managed.StatePass: 0,
	managed.StateFail: 1,
}

var forensicExitCodes = map[string]int{
	forensic.StatePass: 0,
	forensic.StateFail: 1,
}
