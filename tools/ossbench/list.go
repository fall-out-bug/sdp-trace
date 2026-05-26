package main

import (
	"fmt"
	"io"
	"strings"
)

func handleList(stdout, stderr io.Writer, args []string) int {
	if len(args) > 0 {
		fmt.Fprintf(stderr, "unexpected positional args with -list: %s\n", strings.Join(args, " "))
		return 2
	}
	for _, b := range builtIns {
		fmt.Fprintf(stdout, "%s\t%s\n", b.Name, b.Description)
	}
	return 0
}
