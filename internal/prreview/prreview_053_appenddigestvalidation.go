package prreview

func appendDigestValidation(packet Packet, runs RunSet, ledger Ledger, reasons, nextActions *[]string) bool {
	// appendDigestValidation keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	cannotVerify := false
	if packetDigestMismatch(packet, runs, ledger) {
		appendValidationAction(reasons, nextActions, "packet_digest_mismatch", "Create a new packet and rerun review for the current head.")
		cannotVerify = true
	}
	for _, result := range runs.Results {
		if result.PacketDigest != packet.PacketDigest {

			appendValidationAction(reasons, nextActions, "result_packet_digest_mismatch:"+safeID(result.ReviewRunID), "Discard stale reviewer results and rerun review for the current packet.")
			cannotVerify = true
		}
	}
	return cannotVerify
}
