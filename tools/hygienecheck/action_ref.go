package main

import "strings"

func useRef(uses string) string {
	idx := strings.LastIndex(uses, "@")
	if idx < 0 {
		return ""
	}
	return uses[idx+1:]
}
