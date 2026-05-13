package main

func validPRInputSource(source string) bool {
	return source == "github-actions" || source == "github-fixture"
}
