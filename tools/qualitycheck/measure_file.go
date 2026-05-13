package main

import "go/ast"

// measureFile records whole-file metrics used by reports and baseline files.
// The normalized path is part of the persisted file identity.
func measureFile(path string, src []byte, parsed *ast.File) fileMetric {
	// File metrics use the complete source text, unlike function metrics that
	// measure exact declarations.
	source := string(src)
	lines := sourceLines(source)
	commentLines := commentLinesInRange(parsed.Comments, parsed.Pos(), parsed.End())
	// The AST-derived scores keep comments out of branch counting while the raw
	// source remains available for Halstead token analysis.
	cyclo := cyclomaticComplexity(parsed)
	volume := halsteadVolume(source)
	// File MI uses whole-file source and AST metrics, matching the subject that
	// file baselines persist and report.
	// Keeping those inputs adjacent prevents path-scoped baselines from being
	// compared against metrics computed from a narrower source range.
	return fileMetric{
		file:                 normalizePath(path),
		lines:                lines,
		commentLines:         commentLines,
		cyclo:                cyclo,
		halsteadVolume:       volume,
		maintainabilityIndex: maintainabilityIndex(volume, cyclo, lines, commentLines),
	}
}
