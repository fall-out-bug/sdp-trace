package prreview

func RunReview(packet Packet, profile ReviewProfile, opts RunOptions) (RunSet, *RunPreview, error) {
	// RunReview keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if err := validateProfile(profile); err != nil {
		return RunSet{}, nil, err
	}
	opts = normalizeRunOptions(opts)
	if opts.Preview {
		return RunSet{}, preview(packet, profile), nil
	}
	return runReview(packet, profile, opts)
}
