package main

import (
	"context"
)

func runWrapAndParse(ctx context.Context, bin, tmpDir string) (string, verifierState, string) {
	stdout, err := runWrapCommand(ctx, bin, wrapArgs(), tmpDir)
	if err != nil {
		return "", stateFail, err.Error()
	}
	return validateWrapOutput(stdout, tmpDir)
}
