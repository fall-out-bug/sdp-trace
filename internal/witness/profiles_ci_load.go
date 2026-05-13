package witness

func loadSafeCIEnvelopeRecord(record *Record, envelopePath string) (EnvelopeInput, bool) {
	// Envelope loading is the first point where external CI profile facts can
	// affect the record, so malformed or unsafe inputs lower the record directly.
	// The boolean return forces callers to stop before copying untrusted envelope
	// values into the witness record.
	// The output record records the failure class, not the unsafe envelope body.
	var envelope EnvelopeInput
	if err := readSafeJSON(envelopePath, &envelope); err != nil {
		// Malformed or unsafe envelope input is cannot_verify because no
		// trustworthy envelope facts can be extracted.
		applyMalformedEnvelopeState(record)
		return envelope, false
	}
	if unsafeEnvelopeFields(envelope) {
		// Unsafe fields in a parsed envelope are an output-safety failure, not
		// just missing evidence, because persisting them could leak credentials.
		applyUnsafeEnvelopeState(record)
		return envelope, false
	}
	return envelope, true
}

func applyMalformedEnvelopeState(record *Record) {
	applyProfileState(record, StatusCannotVerify, stateCannotVerify, ReasonMalformedInput)
	record.TrustScope = TrustScopeLocalObserved
}

func applyUnsafeEnvelopeState(record *Record) {
	applyProfileState(record, StatusFail, stateFail, ReasonUnsafeOutput)
	record.ProfileStates = defaultProfileStates(stateFail, "cannot_verify")
}
