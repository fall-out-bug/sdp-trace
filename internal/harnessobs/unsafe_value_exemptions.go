package harnessobs

import "strings"

func safeEncodedTokenExemption(path, value string, rawEvent bool) bool {
	return digestField(path) || sha256Pattern.MatchString(value) || rawEventPathLikeField(path, rawEvent)
}

func rawEventPathLikeField(path string, rawEvent bool) bool {
	return rawEvent && rawPathLikeField(path)
}

func digestField(path string) bool {
	last := path
	if idx := strings.LastIndex(last, "."); idx >= 0 {
		last = last[idx+1:]
	}
	return digestFieldNames[last]
}
