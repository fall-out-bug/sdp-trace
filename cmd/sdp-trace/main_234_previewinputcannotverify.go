package main

func previewInputCannotVerify(state string) bool {
	return state == "present_unreadable" || state == "present_malformed"
}
