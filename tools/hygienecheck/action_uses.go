package main

import "strings"

func usesValue(line string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "- ")
	key, value, ok := strings.Cut(line, ":")
	if !ok || strings.TrimSpace(key) != "uses" {
		return ""
	}
	return cleanUseValue(value)
}

func cleanUseValue(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Split(value, " #")[0]
	return strings.Trim(value, `"'`)
}
