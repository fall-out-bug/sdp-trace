package prreview

func invalidPacketCIState(opts PacketOptions) bool {
	return opts.CIState != "" && !validCIState(opts.CIState)
}
