package recorder

import "github.com/fall_out_bug/sdp-trace/internal/trace"

func finishRun(writer *runWriter, commandID string, exitCode int, signal string, options RecorderOptions) error {
	// Command-finished is emitted before run-closed so the closure event can
	// include the final command state in the observed event count.
	if err := appendCommandFinished(writer, options, commandID, exitCode, signal); err != nil {
		return err
	}
	closureState := trace.ClosureStateCompleted
	if exitCode != 0 {
		closureState = trace.ClosureStateCommandFailure
	}
	if err := appendRunClosed(writer, commandID, closureState); err != nil {
		return err
	}
	return writer.finalize(closureState)
}
