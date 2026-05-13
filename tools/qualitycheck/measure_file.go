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
	if fileMIOmitted(parsed, lines) {
		return fileMetric{
			file:                 normalizePath(path),
			lines:                lines,
			commentLines:         commentLines,
			maintainabilityIndex: 100,
		}
	}
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

func fileMIOmitted(parsed *ast.File, lines int) bool {
	return lines < 20 || !hasFunctionDecl(parsed)
}

func hasFunctionDecl(parsed *ast.File) bool {
	// File-level MI uses executable file bodies as the quality subject.
	// Declaration-only files are schema carriers, constants, or static tables.
	// Their generated Halstead volume is noisy because there are no branches to
	// simplify and no smaller behavioral units to extract.
	for _, decl := range parsed.Decls {
		if _, ok := decl.(*ast.FuncDecl); ok {
			return true
		}
	}
	return false
}
