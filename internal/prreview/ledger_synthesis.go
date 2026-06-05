package prreview

import "sort"

func SynthesizeLedger(packet Packet, runs RunSet, existing *Ledger) Ledger {
	// SynthesizeLedger projects reviewer findings into a durable disposition record.
	// It preserves reviewer evidence and prior human dispositions; it does not
	// approve, validate, or close the review.
	findings := synthesizeLedgerFindings(runs, existingFindings(existing))
	sort.Slice(findings, func(i, j int) bool { return findings[i].ID < findings[j].ID })
	return Ledger{SchemaVersion: SchemaVersionLedger, PacketDigest: packet.PacketDigest, Findings: findings}
}

func existingFindings(existing *Ledger) map[string]LedgerFinding {
	byFinding := map[string]LedgerFinding{}
	if existing != nil {
		for _, finding := range existing.Findings {
			byFinding[finding.ID] = finding
		}
	}
	return byFinding
}

func synthesizeLedgerFindings(runs RunSet, byFinding map[string]LedgerFinding) []LedgerFinding {
	findings := []LedgerFinding{}
	for _, result := range runs.Results {
		for _, finding := range result.Findings {
			findings = append(findings, ledgerFindingFromReviewFinding(result, finding, byFinding))
		}
	}
	return findings
}
