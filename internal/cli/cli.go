package cli

import (
	"github.com/spf13/cobra"
)

func Execute(version string) error {
	root := newRoot(version)
	return root.Execute()
}

func newRoot(version string) *cobra.Command {
	root := &cobra.Command{
		Use:           "gobp",
		Short:         "Opinionated Go project boilerplate generator",
		Long:          "gobp generates a buildable Go service skeleton with config, logging, errors, and DB stubs wired in.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
	}
	root.AddCommand(newNewCmd())
	root.AddCommand(newListCmd())
	return root
}
