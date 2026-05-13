package main

import (
	"go/ast"
	"go/token"
)

// cyclomaticIncrement maps one AST node to its independent path contribution.
func cyclomaticIncrement(node ast.Node) int {
	// Count only Go control-flow nodes and boolean operator breaks; leaf
	// expressions do not create independent execution paths.
	switch current := node.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
		// Structured control-flow nodes each add one independent path.
		return 1
	case *ast.BinaryExpr:
		// Chained && and || clauses add paths even when they sit inside one if.
		switch current.Op {
		case token.LAND, token.LOR:
			return 1
		}
	}
	return 0
}
