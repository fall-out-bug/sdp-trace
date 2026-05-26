package main

func quickstartHasCommand(qsSet map[string]bool, req string) bool {
	if req == "go run ./cmd/sdp-trace --help" {
		_, ok := qsSet["go run ./cmd/sdp-trace --help"]
		return ok
	}
	return setContainsPrefix(qsSet, req)
}
