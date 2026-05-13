package witness

import (
	"encoding/json"
)

func finalizeRecordForWrite(record Record) Record {
	// Scan a copy without OutputSafety so the safety attestation cannot satisfy
	// itself or hide unsafe material in the rest of the record.
	// The original record is only returned when the serialized form is free of
	// known unsafe output classes.
	// Serialization failures are treated as unsafe output because no reviewable
	// witness artifact can be produced from that record.
	// A replacement failure record avoids partial redaction of untrusted payload
	// content.
	raw, ok := outputSafetyScanBytes(record)
	if !ok {
		applyProfileState(&record, StatusFail, stateFail, ReasonUnsafeOutput)
		return record
	}
	if !forbiddenOutputPresent(raw) {
		// Passing output safety states only that known unsafe classes were absent
		// from the serialized record; it does not upgrade the profile verdict.
		applyOutputSafetyPass(&record)
		return record
	}
	return unsafeOutputRecord(record.Kind)
}

func outputSafetyScanBytes(record Record) ([]byte, bool) {
	// Scan bytes exclude OutputSafety so the attestation cannot recursively
	// satisfy the safety check that is about to be recorded.
	scanRecord := record
	scanRecord.OutputSafety = nil
	raw, err := json.Marshal(scanRecord)
	return raw, err == nil
}
