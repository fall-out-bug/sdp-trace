package main

import "runtime"

func wrapArgs() []string {
	if runtime.GOOS == "windows" {
		return []string{"wrap", "cmd", "/c", "exit", "0"}
	}
	return []string{"wrap", "true"}
}
