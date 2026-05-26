package main

import (
	"context"
	"time"
)

func runSingleCommand(cmd string, args []string, dir string) (time.Duration, error, bool) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	c := benchmarkCommand(ctx, cmd, args, dir)
	err := c.Run()
	if err != nil {
		return 0, err, ctx.Err() == context.DeadlineExceeded
	}
	return time.Since(start), nil, false
}
