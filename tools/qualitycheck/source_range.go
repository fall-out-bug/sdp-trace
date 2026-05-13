package main

// validSourceRange keeps AST offsets from slicing outside the original file.
func validSourceRange(start int, end int, length int) bool {
	return start >= 0 && end <= length && start <= end
}
