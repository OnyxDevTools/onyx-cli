package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// These variables can be overridden at build time via -ldflags.
var (
	Version = "0.2.0"
	Commit  = "unknown"
	Date    = "unknown"
)

func newVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show onyx CLI version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "onyx version %s (commit %s, built %s)\n", Version, Commit, Date)
			return nil
		},
	}
}
