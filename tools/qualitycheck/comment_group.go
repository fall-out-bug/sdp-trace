package main

import (
	"go/ast"
	"go/token"
)

// commentGroupInRange rejects package-level and neighboring function comments.
func commentGroupInRange(group *ast.CommentGroup, start token.Pos, end token.Pos) bool {
	return group.Pos() >= start && group.End() <= end
}
