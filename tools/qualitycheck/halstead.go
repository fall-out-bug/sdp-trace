package main

import (
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
