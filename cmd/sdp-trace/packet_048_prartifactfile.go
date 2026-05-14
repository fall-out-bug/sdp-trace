package main

type packetPRArtifactFile struct {
	label string
	write func() error
}
