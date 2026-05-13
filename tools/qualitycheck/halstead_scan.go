package main

import (
	"go/scanner"
	"go/token"
)

func scanHalsteadTokens(source string, counts *halsteadCounts) {
	var scan scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(source))
	scan.Init(file, []byte(source), nil, scanner.ScanComments)
	// ScanComments keeps comment tokens visible here; addHalsteadToken owns the
	// decision to exclude them from operator and operand vocabulary.
	for {
		_, tok, literal := scan.Scan()
		if tok == token.EOF {
			return
		}
		addHalsteadToken(tok, literal, counts)
	}
}
