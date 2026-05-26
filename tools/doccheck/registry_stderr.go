package main

import (
	"bytes"
	"strings"
)

func commandSurfaceStderr(stderr *bytes.Buffer) string {
	return strings.TrimSpace(stderr.String())
}
