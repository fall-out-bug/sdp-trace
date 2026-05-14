package main

// Index is the machine-readable schema documentation source of truth.
type Index struct {
	Version string        `json:"version"`
	Schemas []SchemaEntry `json:"schemas"`
}

// SchemaEntry describes one schema file.
type SchemaEntry struct {
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	Purpose         string   `json:"purpose"`
	ExampleCoverage string   `json:"example_coverage,omitempty"`
	Examples        []string `json:"examples,omitempty"`
}
