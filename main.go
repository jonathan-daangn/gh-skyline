// Package main provides the entry point for the GitHub CLI gh-skyline extension.
package main

import (
	"context"
	"os"

	"github.com/github/gh-skyline/cmd"
)

type exitCode int

const (
	exitOK    exitCode = 0
	exitError exitCode = 1
)

func main() {
	code := start()
	os.Exit(int(code))
}

func start() exitCode {
	exitCode := exitOK
	ctx := context.Background()

	if err := cmd.Execute(ctx); err != nil {
		exitCode = exitError
	}

	return exitCode
}
