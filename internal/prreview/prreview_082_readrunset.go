package prreview

func ReadRunSet(path string) (RunSet, error) {
	// ReadRunSet keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	var runs RunSet

	path = runSetPath(path)
	if err := readJSON(path, &runs); err != nil {
		return runs, err
	}
	if err := validateRunSet(runs); err != nil {
		return runs, err
	}
	return runs, nil
}
