package main

import "io"

// options is the immutable command configuration used after flag parsing.
type options struct {
	// Threshold fields are disabled by zero values.
	cycloOver     int
	cognitiveOver int
	miUnder       float64
	// Baseline paths keep file and function MI ratchets independent.
	fileMIBaseline          string
	writeFileMIBaseline     string
	functionMIUnder         float64
	functionMIBaseline      string
	writeFunctionMIBaseline string
	// Rendering fields affect output, not measured evidence.
	gocyclo  bool
	failOnly bool
	err      io.Writer
}
