package prreview

func ReadLedger(path string) (Ledger, error) {
	var ledger Ledger
	return ledger, readJSON(path, &ledger)
}
