package main

import (
	"go/scanner"
	"go/token"
	"math"
)

func halsteadVolume(source string) float64 {
	// Halstead is computed from scanner tokens so it stays parser-independent
	// for both file slices and function declaration slices.
	counts := halsteadCounts{
		operators: map[string]int{},
		operands:  map[string]int{},
	}
	scanHalsteadTokens(source, &counts)
	vocabulary := len(counts.operators) + len(counts.operands)
	if vocabulary == 0 || counts.length == 0 {
		return 0
	}
	return float64(counts.length) * math.Log2(float64(vocabulary))
}

type halsteadCounts struct {
	operators map[string]int
	operands  map[string]int
	length    int
}

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
