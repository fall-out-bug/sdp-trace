package main

import (
	"context"
	"io"
)

type commandHandler func(context.Context, []string, io.Writer, io.Writer) int
