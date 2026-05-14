package main

import (
	"io"
)

type assessHandler func(*flagSet, io.Writer, io.Writer) int
