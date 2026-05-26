package main

import "strings"

func setContainsPrefix(set map[string]bool, prefix string) bool {
	for s := range set {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}
