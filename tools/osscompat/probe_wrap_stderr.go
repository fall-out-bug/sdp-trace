package main

import "os/exec"

func wrapStderr(err error) []byte {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.Stderr
	}
	return nil
}
