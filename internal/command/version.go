package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newVersionCommand creates the "version" subcommand.
func newVersionCommand(version string, goVersion string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of driver-scanner",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("driver-scanner %s (built with %s)\n", version, goVersion)
		},
	}
}
