package harnessobs

import (
	"encoding/hex"
	"strings"
)

func sourceCommitHash(commit string) bool {
	if len(commit) != 40 || commit != strings.ToLower(commit) {
		return false
	}
	_, err := hex.DecodeString(commit)
	return err == nil
}
