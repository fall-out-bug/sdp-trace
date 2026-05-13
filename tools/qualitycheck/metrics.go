package main

// functionMetric is the stable row shape for function-level quality evidence.
type functionMetric struct {
	// Identity fields form the report row and function MI baseline key.
	file        string
	line        int
	column      int
	packageName string
	name        string
	// Measurement fields stay with the derived MI so drift can be diagnosed.
	cyclo                int
	cognitive            int
	lines                int
	commentLines         int
	halsteadVolume       float64
	maintainabilityIndex float64
}

// fileMetric is the stable row shape for whole-file quality evidence.
type fileMetric struct {
	// File identity is normalized before reports or baselines consume it.
	file string
	// Whole-file measurements are kept beside the derived MI verdict.
	lines                int
	commentLines         int
	cyclo                int
	halsteadVolume       float64
	maintainabilityIndex float64
}

// qualityReport carries measured evidence before any rendering or gate verdict.
type qualityReport struct {
	functions []functionMetric
	files     []fileMetric
}
