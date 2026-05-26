package main

// benchmarkDef defines a single benchmark target.
type benchmarkDef struct {
	Name        string
	Description string
	Cmd         string
	Args        []string
	Dir         string // working directory; empty means current directory
	Cleanup     func() // optional cleanup after benchmark completes
	Source      string // "repo-root", "temp-build", or "PATH"
}
