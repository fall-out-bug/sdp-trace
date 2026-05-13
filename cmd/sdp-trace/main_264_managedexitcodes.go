package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/managed"
)

var managedExitCodes = map[string]int{
	managed.StatePass: 0,
	managed.StateFail: 1,
}
