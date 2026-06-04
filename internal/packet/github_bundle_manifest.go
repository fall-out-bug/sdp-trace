package packet

func githubBundleManifest(bundleID string, packet Packet, entries []BundleEntry) BundleManifest {
	// githubBundleManifest keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	return BundleManifest{
		SchemaVersion: BundleSchemaVersion,
		BundleID:      bundleID,
		PacketDigest:  PacketDigest(packet),
		Entries:       entries,
	}
}
