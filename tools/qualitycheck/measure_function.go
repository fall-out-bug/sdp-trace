package main

import (
	"go/ast"
	"go/token"
)

// measureFunction returns the report-ready metric for one function declaration.
func measureFunction(path string, src []byte, fset *token.FileSet, parsed *ast.File, fn *ast.FuncDecl) functionMetric {
	// Keep measurement and row construction separate so the metric formula can
	// change without changing report identity fields.
	measurement := measureFunctionSource(src, fset, parsed, fn)
	return buildFunctionMetric(path, parsed, fn, measurement)
}

// measureFunctionSource derives formula inputs from one function source span.
func measureFunctionSource(src []byte, fset *token.FileSet, parsed *ast.File, fn *ast.FuncDecl) functionMeasurement {
	// All inputs come from the same declaration span; package comments and
	// sibling declarations are intentionally outside this measurement.
	source := sourceForNode(src, fset, fn)
	lines := sourceLines(source)
	// Comment credit uses the AST comment ranges so only comments owned by this
	// function can affect the function-level maintainability ratchet.
	commentLines := commentLinesInRange(parsed.Comments, fn.Pos(), fn.End())
	cyclo := cyclomaticComplexity(fn)
	cognitive := cognitiveComplexity(fn)
	// Halstead volume is source-based while complexity is AST-based; deriving
	// them here keeps the final MI row internally consistent.
	volume := halsteadVolume(source)
	mi := maintainabilityIndex(volume, cyclo, lines, commentLines)
	if functionMIOmitted(lines, cyclo, cognitive) {
		mi = 100
	}
	return functionMeasurement{
		lines:        lines,
		commentLines: commentLines,
		cyclo:        cyclo,
		cognitive:    cognitive,
		volume:       volume,
		mi:           mi,
		position:     fset.Position(fn.Pos()),
	}
}

func functionMIOmitted(lines int, cyclo int, cognitive int) bool {
	// Tiny/simple Go helpers are better covered by cyclomatic, cognitive, and
	// CRAP gates; raw function MI overweights their Halstead token noise.
	return lines < 18 && cyclo <= 4 && cognitive <= 3
}
