package main

// summarize counts results by state.
func summarize(results []probeResult) (pass, fail, cant, na int) {
	for _, r := range results {
		pass += boolToInt(r.State == statePass)
		fail += boolToInt(r.State == stateFail)
		cant += boolToInt(r.State == stateCannotVerify)
		na += boolToInt(r.State == stateNotAssessed)
	}
	return
}
