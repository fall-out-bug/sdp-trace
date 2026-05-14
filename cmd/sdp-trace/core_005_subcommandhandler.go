package main

import (
	"io"
)

type subcommandHandler func([]string, io.Writer, io.Writer) int
