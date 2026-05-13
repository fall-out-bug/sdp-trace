package main

import "go/ast"

// buildFunctionMetric combines measured values with stable report identity.
func buildFunctionMetric(path string, parsed *ast.File, fn *ast.FuncDecl, measurement functionMeasurement) functionMetric {
	// Position is captured once during measurement so line and column stay tied
	// to the same token lookup.
	// The normalized file path is part of the stable baseline key.
	return functionMetric{
		file:   normalizePath(path),
		line:   measurement.position.Line,
		column: measurement.position.Column,

		// Package and function names make report rows readable without reparsing
		// source at the printing boundary.
		packageName: parsed.Name.Name,
		name:        functionName(fn),

		// Keep raw components beside derived MI so baseline reports can explain
		// why a function changed rather than only showing the verdict.
		cyclo:                measurement.cyclo,
		cognitive:            measurement.cognitive,
		lines:                measurement.lines,
		commentLines:         measurement.commentLines,
		halsteadVolume:       measurement.volume,
		maintainabilityIndex: measurement.mi,
	}
}
