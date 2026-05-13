package recorder

import "github.com/fall_out_bug/sdp-trace/internal/trace"

type preparedRun struct {
	options   RecorderOptions
	runDir    string
	contract  trace.Contract
	writer    *runWriter
	commandID string
}
