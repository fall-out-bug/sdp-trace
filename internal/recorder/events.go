package recorder

import (
	"fmt"
	"os"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

// Event helpers keep each recorder event payload close to the trace schema
// type it persists. That makes evidence fields auditable without requiring the
// caller to know the writer's hash-chain mechanics.
//
// Each helper has one evidence responsibility:
// - recorder/run start binds observer, nonce, task, command, and contract,
// - command start records wrapper context before process execution,
// - command finish binds exit state to output digests,
// - run closed records the terminal chain head and observed event count.
//
// Those payload boundaries match the portable trace contract and keep replay
// failures attributable to a concrete lifecycle phase.
//
// The functions here contain no filesystem writes; persistence is handled by
// runWriter after each payload is normalized and hashed.

func appendRunStartEvents(writer *runWriter, options RecorderOptions, contract trace.Contract) error {
	// Recorder attachment comes first so the run-start event has a verifiable
	// local observer and nonce already anchored in the chain.
	if err := writer.appendEvent(trace.EventRecorderAttached, trace.RecorderAttachedPayload{
		RecorderVersion: trace.RecorderVersion,
		RunNonce:        fmt.Sprintf("run-nonce-%s", writer.manifest.RunID),
		Pid:             os.Getpid(),
	}); err != nil {
		return err
	}

	// The run-start payload binds user intent, command identity, and contract
	// digest into the first task-scoped event.
	return writer.appendEvent(trace.EventRunStarted, trace.RunStartedPayload{
		TaskRef:        options.Task,
		Command:        commandName(options.Command),
		ContractID:     contract.ContractID,
		ContractDigest: writer.manifest.ContractDigest,
		ContractPath:   writer.manifest.ContractPath,
	})
}

func appendCommandStarted(writer *runWriter, options RecorderOptions, commandID string) error {
	// The started event records the wrapper and cwd before process execution so
	// replay can distinguish command behavior from recorder placement.
	return writer.appendEvent(trace.EventCommandStarted, trace.CommandStartedPayload{
		RecordedCommandPayload: trace.RecordedCommandPayload{
			CommandID: commandID,
			Command:   commandName(options.Command),
		},
		CmdPID:         0,
		WrapperName:    options.WrapperName,
		Cwd:            currentWorkingDir(),
		ContractDigest: writer.manifest.ContractDigest,
	})
}

func appendCommandFinished(writer *runWriter, options RecorderOptions, commandID string, exitCode int, signal string) error {
	// Digests are read after the command stream writers finish, keeping command
	// output evidence bound to the terminal process state.
	return writer.appendEvent(trace.EventCommandFinished, trace.CommandFinishedPayload{
		RecordedCommandPayload: trace.RecordedCommandPayload{
			CommandID: commandID,
			Command:   commandName(options.Command),
		},
		ExitCode:     exitCode,
		Signal:       signal,
		StdoutDigest: writer.stdoutDigest(),
		StderrDigest: writer.stderrDigest(),
	})
}

func appendRunClosed(writer *runWriter, commandID, closureState string) error {
	// Closure records the current chain head and event count so later replay can
	// detect both missing events and altered terminal state.
	return writer.appendEvent(trace.EventRunClosed, trace.RunClosedPayload{
		ChainHead:     writer.lastEventHash(),
		ClosureState:  closureState,
		CommandID:     commandID,
		MissingCount:  0,
		ObservedCount: writer.eventCount(),
	})
}
