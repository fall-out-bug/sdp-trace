package main

func optionalStringMatches(expected, actual string) bool {
	return expected == "" || actual == expected
}
