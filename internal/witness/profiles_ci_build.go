package witness

func BuildCIEnvelopeProfile(kind, runsRoot, reportDir, envelopePath string) (Record, error) {
	record := baseRecord(kind)
	// Current artifact hashes are computed before reading the envelope so the
	// envelope can be compared with live local evidence, not trusted by itself.
	if err := populateCIEnvelopeArtifacts(&record, runsRoot, reportDir); err != nil {
		return Record{}, err
	}
	if !applyCIEnvelopeInputState(&record, kind, runsRoot, envelopePath) {
		return record, nil
	}
	return record, nil
}
