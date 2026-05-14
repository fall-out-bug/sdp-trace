package main

import (
	"errors"
	"os"
)

func protectedInputErrorStatus(err error) string {
	if os.IsNotExist(err) || errors.Is(err, os.ErrPermission) {
		// Protected preview distinguishes unavailable inputs from malformed JSON
		// so users know whether to fix access or content.
		return "present_unreadable"
	}
	return "present_malformed"
}
