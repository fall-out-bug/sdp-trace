package main

// coverage.go validates the example_coverage field and its consistency with
// the examples array for a schema index entry.

import "fmt"

// validateExampleCoverage checks the coverage value and cross-validates it
// against the examples list.
func validateExampleCoverage(entry *SchemaEntry) []string {
	var issues []string
	if !validExampleCoverages[entry.ExampleCoverage] {
		issues = append(issues, fmt.Sprintf("%s: invalid example_coverage %q", entry.Name, entry.ExampleCoverage))
	}
	issues = append(issues, validateExamplePresence(entry)...)
	issues = append(issues, validateExampleAbsence(entry)...)
	return issues
}

// validateExamplePresence requires a non-empty examples list when coverage is present.
func validateExamplePresence(entry *SchemaEntry) []string {
	if entry.ExampleCoverage == examplePresent && len(entry.Examples) == 0 {
		return []string{fmt.Sprintf("%s: example_coverage is present but examples list is empty", entry.Name)}
	}
	return nil
}

// validateExampleAbsence rejects a non-empty examples list when coverage is not present.
func validateExampleAbsence(entry *SchemaEntry) []string {
	if entry.ExampleCoverage != examplePresent && len(entry.Examples) > 0 {
		return []string{fmt.Sprintf("%s: example_coverage is %q but examples list is non-empty", entry.Name, entry.ExampleCoverage)}
	}
	return nil
}
