package main

import "go/ast"

// receiverName returns the named receiver type used in function report labels.
func receiverName(expr ast.Expr) (string, bool) {
	// Receiver names are only report labels; unsupported receiver shapes remain
	// anonymous instead of affecting analysis.
	inner, ok := receiverInnerExpr(expr)
	if !ok {
		return "", false
	}
	ident, ok := inner.(*ast.Ident)
	if !ok {
		return "", false
	}
	return ident.Name, true
}

// receiverInnerExpr unwraps receiver syntax until a named type can be checked.
func receiverInnerExpr(expr ast.Expr) (ast.Expr, bool) {
	// Pointer receivers are transparent for naming; generic receivers are
	// unwrapped by the next parser boundary.
	if ident, ok := expr.(*ast.Ident); ok {
		return ident, true
	}
	if star, ok := expr.(*ast.StarExpr); ok {
		return receiverInnerExpr(star.X)
	}
	return receiverGenericInnerExpr(expr)
}

// receiverGenericInnerExpr handles generic receiver syntax added by newer Go
// AST forms.
func receiverGenericInnerExpr(expr ast.Expr) (ast.Expr, bool) {
	// Generic receiver syntax stores the named receiver type under X for both
	// single-index and index-list forms.
	if indexed, ok := expr.(*ast.IndexExpr); ok {
		return indexedReceiverInnerExpr(indexed)
	}
	if indexedList, ok := expr.(*ast.IndexListExpr); ok {
		return indexListReceiverInnerExpr(indexedList)
	}
	return nil, false
}

// indexedReceiverInnerExpr unwraps a single-parameter generic receiver.
func indexedReceiverInnerExpr(expr *ast.IndexExpr) (ast.Expr, bool) {
	return receiverInnerExpr(expr.X)
}

// indexListReceiverInnerExpr unwraps a multi-parameter generic receiver.
func indexListReceiverInnerExpr(expr *ast.IndexListExpr) (ast.Expr, bool) {
	return receiverInnerExpr(expr.X)
}
