package main

import (
	"fmt"
	"os"

	"github.com/tienanhnguyen999/gobp/internal/cli"
)

// Version is set at build time via -ldflags.
var Version = "dev"

func main() {
	if err := cli.Execute(Version); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
