package packet

func RenderMarkdown(bundle Bundle) (string, error) {
	// RenderMarkdown keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.

	packet, err := renderablePacket(bundle)
	if err != nil {
		return "", err
	}
	return renderPacketMarkdown(packet, bundle.Manifest), nil
}
