package prreview

func appendDigestValidation(packet Packet, runs RunSet, ledger Ledger, reasons, nextActions *[]string) bool {
	cannotVerify := false
	// Packet/run/ledger digest drift invalidates the whole review bundle.
	if packetDigestMismatch(packet, runs, ledger) {
		appendValidationAction(reasons, nextActions, "packet_digest_mismatch", "Create a new packet and rerun review for the current head.")
		cannotVerify = true
	}
	// Individual stale reviewer rows remain visible with their run IDs so the
	// caller can discard only the affected evidence when possible.
	for _, result := range runs.Results {
		if appendResultDigestValidation(packet, result, reasons, nextActions) {
			cannotVerify = true
		}
	}
	return cannotVerify
}

func appendResultDigestValidation(packet Packet, result ReviewerResult, reasons, nextActions *[]string) bool {
	if result.PacketDigest == packet.PacketDigest {
		return false
	}
	appendValidationAction(reasons, nextActions, "result_packet_digest_mismatch:"+safeID(result.ReviewRunID), "Discard stale reviewer results and rerun review for the current packet.")
	return true
}

func packetDigestMismatch(packet Packet, runs RunSet, ledger Ledger) bool {
	return runs.PacketDigest != packet.PacketDigest || ledger.PacketDigest != packet.PacketDigest
}

func appendValidationAction(reasons, nextActions *[]string, reason, nextAction string) {
	*reasons = append(*reasons, reason)
	*nextActions = append(*nextActions, nextAction)
}
