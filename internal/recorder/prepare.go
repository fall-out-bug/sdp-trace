package recorder

import (
	"os"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

// Preparation code owns every recorder decision made before the child process
// starts: where evidence is written, which contract names the run, and which
// first events anchor the append-only chain.
//
// The ordering in this file is intentionally strict:
// - reject unsafe output before reading trust inputs,
// - resolve and validate the contract before command execution,
// - persist the initial manifest before appending command events,
// - copy the environment before handing options to the process runner.
//
// Keeping those steps explicit avoids a common recorder failure mode where a
// partially prepared run can look like valid evidence after the command exits.

func prepareRun(options RecorderOptions) (preparedRun, error) {
	// File layout, contract binding, and writer initialization are completed
	// before the user command starts so failed setup cannot emit partial trust
	// evidence.
	runDir, contract, writer, err := prepareRunFiles(options)
	if err != nil {
		return preparedRun{}, err
	}

	options.Env = ifNilCopy(os.Environ())
	commandID, err := startRecordedCommand(writer, options, contract)
	if err != nil {
		return preparedRun{}, err
	}
	return preparedRun{options: options, runDir: runDir, contract: contract, writer: writer, commandID: commandID}, nil
}

func prepareRunFiles(options RecorderOptions) (string, trace.Contract, *runWriter, error) {
	// Output validation happens before contract loading so unsafe destinations
	// fail without touching trust inputs.
	runDir, err := prepareRunDir(options.OutputDir)
	if err != nil {
		return "", trace.Contract{}, nil, err
	}

	contract, err := resolveRecorderContract(options.ContractPath, options.UseDefaultContract)
	if err != nil {
		return "", trace.Contract{}, nil, err
	}
	// Writer setup is the last preparation step because it persists the initial
	// manifest with the resolved contract digest.
	writer, err := initializedRunWriter(runDir, contract, options)
	if err != nil {
		return "", trace.Contract{}, nil, err
	}
	return runDir, contract, writer, nil
}

func initializedRunWriter(runDir string, contract trace.Contract, options RecorderOptions) (*runWriter, error) {
	// Manifest initialization is separated from writer construction so callers
	// never receive a writer whose first run manifest was not persisted.
	writer, err := newRunWriter(runDir, contract, options.Task)
	if err != nil {
		return nil, err
	}

	if err := initializeRunWriter(writer, options, contract); err != nil {
		return nil, err
	}
	return writer, nil
}
