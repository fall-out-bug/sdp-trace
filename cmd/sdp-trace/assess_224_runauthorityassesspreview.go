package main

import (
	"io"
)

func runAuthorityAssessPreview(opts *flagSet, stdout io.Writer) int {
	// Authority preview treats the package path as an input pointer only; it does
	// not evaluate authority rules.
	inputs := map[string]string{
		"authority_package": managedInputStatus(opts.stringValue("authority-package")),
	}
	writeJSONPayloadUnchecked(stdout, newAuthorityPreviewReport(inputs))
	return previewInputExitCode(inputs)
}
