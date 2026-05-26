package main

import (
	"bytes"
	"os/exec"
)

func commandSurfaceCmd() (*exec.Cmd, *bytes.Buffer) {
	cmd := exec.Command("go", "run", "./cmd/sdp-trace", "command-surface")
	cmd.Dir = repoRoot()
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	return cmd, stderr
}
