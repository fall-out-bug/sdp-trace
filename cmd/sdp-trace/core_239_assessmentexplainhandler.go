package main

import (
	"io"
)

type assessmentExplainHandler func(string, io.Writer, io.Writer) int
