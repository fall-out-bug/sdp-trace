package main

import "os/exec"

func gitFileExistsAtRef(ref, path string) (bool, error) {
	// cat-file is the authority for source-bound baseline existence because it
	// checks the immutable base ref instead of the mutable working tree.
	cmd := exec.Command("git", "cat-file", "-e", ref+":"+path)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}
	// ExitError means git answered "not present"; other errors are execution
	// failures that policy callers must surface instead of treating as absence.
	if _, ok := err.(*exec.ExitError); ok {
		return false, nil
	}
	return false, err
}
