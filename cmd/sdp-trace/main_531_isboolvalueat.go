package main

func isBoolValueAt(args []string, idx int) bool {
	return idx < len(args) && isBoolLiteral(args[idx])
}
