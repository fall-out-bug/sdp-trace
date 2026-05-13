package prreview

type citationResolver struct {
	matches    func(Packet, Citation) bool
	resolvable func(Citation) bool
}
