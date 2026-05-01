package postgen

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/tienanhnguyen999/gobp/internal/selection"
)

func Run(dir string, sel selection.Selection) error {
	switch sel.Git {
	case "init":
		if err := runIn(dir, "git", "init"); err != nil {
			fmt.Fprintf(os.Stderr, "warn: git init failed: %v\n", err)
		}
	case "commit":
		_ = runIn(dir, "git", "init")
		_ = runIn(dir, "git", "add", ".")
		if err := runIn(dir, "git", "commit", "-m", "initial commit from gobp"); err != nil {
			fmt.Fprintf(os.Stderr, "warn: git commit failed: %v\n", err)
		}
	case "none":
		// nothing
	}

	if !sel.NoTidy {
		if err := runIn(dir, "go", "mod", "tidy"); err != nil {
			fmt.Fprintf(os.Stderr, "warn: go mod tidy failed: %v\n", err)
		}
	}

	return nil
}

func runIn(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run()
}
