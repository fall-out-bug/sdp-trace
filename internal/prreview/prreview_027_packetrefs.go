package prreview

type packetRefs struct {
	diff         SafeRef
	metadata     *SafeRef
	context      []SafeRef
	verification []SafeRef
}
