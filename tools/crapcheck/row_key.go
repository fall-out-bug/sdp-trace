package main

// Coverage and complexity rows join on the same file:line subject. Function
// names stay out of the key because coverage profiles report source locations,
// not gocyclo's parsed function labels.
func (row coverageRow) key() string {
	return row.file + ":" + row.line
}

func (row complexityRow) key() string {
	return row.file + ":" + row.line
}
