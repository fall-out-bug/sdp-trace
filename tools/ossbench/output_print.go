package main

import "io"

func printResults(w io.Writer, results []benchmarkResult, asJSON bool) error {
	if asJSON {
		return printJSONResults(w, results)
	}
	return printTextResults(w, results)
}
