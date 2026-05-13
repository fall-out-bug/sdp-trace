package prreview

import (
	"sort"
)

func SynthesizeLedger(packet Packet, runs RunSet, existing *Ledger) Ledger {
	// SynthesizeLedger keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	findings := synthesizeLedgerFindings(runs, existingFindings(existing))
	sort.Slice(findings, func(i, j int) bool { return findings[i].ID < findings[j].ID })
	return Ledger{SchemaVersion: SchemaVersionLedger, PacketDigest: packet.PacketDigest, Findings: findings}
}
