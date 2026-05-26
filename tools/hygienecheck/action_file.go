package main

import "strings"

func actionFile(f string) bool {
	if !strings.HasPrefix(f, ".github/workflows/") {
		return false
	}
	return strings.HasSuffix(f, ".yml") || strings.HasSuffix(f, ".yaml")
}
