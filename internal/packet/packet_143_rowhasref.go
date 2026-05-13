package packet

func rowHasRef(row Row, ref string) bool {
	// Evidence refs are exact identifiers; contradiction attribution must not
	// infer aliases or prefixes when selecting a row.
	for _, rowRef := range row.EvidenceRefs {
		if rowRef == ref {
			return true
		}
	}
	return false
}
