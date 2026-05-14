package main

import (
	"io"
)

type assessPreviewHandler func(*flagSet, io.Writer) int
