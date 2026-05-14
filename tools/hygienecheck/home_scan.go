package main

import (
	"bytes"
	"strings"
)

// bytesContainHomePath scans data for "/home/" followed by a valid path
// segment start character. Prose patterns such as "/home/..." are skipped
// because the dot after the slash is not a valid segment start.
func bytesContainHomePath(data []byte) bool {
	needle := []byte("/home/")
	for {
		i := bytes.Index(data, needle)
		if i == -1 {
			return false
		}
		after := i + len(needle)
		if isValidHomePrefix(data, after) {
			return true
		}
		data = data[after:]
	}
}

// isValidHomePrefix checks whether the byte after a "/home/" match begins a
// real path segment rather than prose punctuation.
func isValidHomePrefix(data []byte, after int) bool {
	return after < len(data) && isPathSegmentStart(data[after])
}

// isPathSegmentStart returns true for characters that can start a Linux path
// segment: alphanumeric, underscore, and hyphen.
func isPathSegmentStart(b byte) bool {
	return strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-", rune(b))
}
