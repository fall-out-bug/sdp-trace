package main

// benchmarkResult holds the measured statistics for one benchmark.
type benchmarkResult struct {
	Name                string    `json:"name"`
	Description         string    `json:"description,omitempty"`
	Command             string    `json:"command"`
	Argv                []string  `json:"argv,omitempty"`
	WorkingDirectory    string    `json:"working_directory,omitempty"`
	BinaryPath          string    `json:"binary_path,omitempty"`
	BinarySource        string    `json:"binary_source,omitempty"`
	Environment         envInfo   `json:"environment,omitempty"`
	AttemptedIterations int       `json:"attempted_iterations"`
	SucceededIterations int       `json:"succeeded_iterations"`
	MinMs               float64   `json:"min_ms"`
	MaxMs               float64   `json:"max_ms"`
	MedianMs            float64   `json:"median_ms"`
	AllMs               []float64 `json:"all_ms,omitempty"`
	Error               string    `json:"error,omitempty"`
}
