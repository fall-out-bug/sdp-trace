package main

import "strings"

func localUse(uses string) bool {
	return strings.HasPrefix(uses, "./") || strings.HasPrefix(uses, ".github/")
}

func pinnedUse(uses string) bool {
	ref := useRef(uses)
	return len(ref) == 40 && isHex(ref)
}
