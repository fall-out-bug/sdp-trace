package packet

import "strings"

// Passing rows must cite retained evidence; otherwise a green packet row would
// be a prose claim rather than a replayable trust claim.
func (v *bundleValidator) validatePassRowEvidence(row Row) {
	if row.State == StatePass && len(row.EvidenceRefs) == 0 {
		v.add("%s pass requires retained evidence refs", row.ID)
	}
}

func (v *bundleValidator) validateRowEvidenceRefs(row Row) {
	for _, ref := range row.EvidenceRefs {
		v.validateEvidenceRef(row.ID, row.State, ref)
	}
}

// Evidence-ref checks first prove the manifest can resolve the citation, then
// apply stricter expiry/retention checks only to rows that claim pass.
func (v *bundleValidator) validateEvidenceRef(rowID, state, ref string) {
	entry, ok := v.entryByRef[ref]
	if !ok {
		v.add("%s evidence ref %q is absent from manifest", rowID, ref)
		return
	}
	if strings.TrimSpace(v.resolverByRef[ref]) == "" {
		v.add("%s evidence ref %q has no resolver entry", rowID, ref)
	}
	if state != StatePass {
		return
	}
	v.validatePassEvidenceRef(rowID, ref, entry)
}

// Pass evidence must remain usable at validation time; expired or non-retained
// artifacts cannot support a pass verdict.
func (v *bundleValidator) validatePassEvidenceRef(rowID, ref string, entry BundleEntry) {
	if entryExpired(entry, v.now) {
		v.add("%s pass cites expired artifact ref %q", rowID, ref)
	}
	if passRefUnverifiable(entry) {
		v.add("%s pass cites unverifiable artifact ref %q", rowID, ref)
	}
}
