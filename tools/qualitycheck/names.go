package main

import (
	"go/ast"
)

func functionName(fn *ast.FuncDecl) string {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		// Free functions are keyed by their declared name.
		return fn.Name.Name
	}
	receiver, ok := receiverName(fn.Recv.List[0].Type)
	if !ok {
		// Unsupported receiver syntax falls back to the stable function name
		// rather than inventing a partial key.
		return fn.Name.Name
	}
	// Method names include the normalized receiver so baselines can distinguish
	// identical method names across types.
	return "(" + receiver + ")." + fn.Name.Name
}
