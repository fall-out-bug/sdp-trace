package main

import "go/token"

// functionMeasurement holds raw measurements before report identity is added.
type functionMeasurement struct {
	// Source-span counts feed the MI formula.
	lines        int
	commentLines int
	// Complexity and Halstead values are measured from the same function body.
	cyclo     int
	cognitive int
	volume    float64
	mi        float64
	// Position pins the report row to the declaration start.
	position token.Position
}
