package main

import "strings"

func baseCommand(usage string) string {
	fields := strings.Fields(usage)
	if len(fields) >= 2 && fields[0] == "sdp-trace" {
		return fields[1]
	}
	return ""
}
