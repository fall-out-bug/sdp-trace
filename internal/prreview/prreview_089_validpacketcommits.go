package prreview

func validPacketCommits(opts PacketOptions) bool {
	return sha40Pattern.MatchString(opts.BaseCommit) && sha40Pattern.MatchString(opts.HeadCommit)
}
