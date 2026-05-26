package main

import "fmt"

func runCommandSurface() ([]byte, error) {
	cmd, stderr := commandSurfaceCmd()
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("run sdp-trace command-surface: %w: %s", err, commandSurfaceStderr(stderr))
	}
	return output, nil
}
