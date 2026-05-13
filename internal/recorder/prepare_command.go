package recorder

import "github.com/fall_out_bug/sdp-trace/internal/trace"

// Command preparation appends the pre-execution event prefix and returns the
// command identifier used by every later command lifecycle event.

func startRecordedCommand(writer *runWriter, options RecorderOptions, contract trace.Contract) (string, error) {
	// The command start event depends on the recorder/run start records so the
	// event chain has a stable provenance prefix before command execution.
	if err := appendRunStartEvents(writer, options, contract); err != nil {
		return "", err
	}
	// Command IDs are generated after the run prefix is written so a failed
	// prefix never leaks an identifier for a command that was not recorded.
	commandID := randomHex(12)
	if err := appendCommandStarted(writer, options, commandID); err != nil {
		return "", err
	}
	return commandID, nil
}
