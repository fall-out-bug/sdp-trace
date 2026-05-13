package main

import (
	"go/ast"
	"go/token"
)

// commentLinesInRange counts comments owned by one function for the MI formula.
func commentLinesInRange(groups []*ast.CommentGroup, start token.Pos, end token.Pos) int {
	lines := 0
	for _, group := range groups {
		// Function MI counts only comments wholly owned by the function body;
		// surrounding package or sibling comments must not inflate the result.
		if !commentGroupInRange(group, start, end) {
			continue
		}
		for _, comment := range group.List {
			// sourceLines preserves the same physical-line counting used for
			// function bodies so block comments and line comments stay comparable.
			lines += sourceLines(comment.Text)
		}
	}
	return lines
}
