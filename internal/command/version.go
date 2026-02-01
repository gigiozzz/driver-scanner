package command

import (
	"fmt"
	"io"
	"runtime"

	"github.com/spf13/cobra"
)

// Build information set by ldflags.
var (
	Version    = "1.0.0-SNAPSHOT"
	CommitHash = "unknown"
	BuildDate  = "unknown"
)

// VersionOptions holds the configuration for the version command.
type VersionOptions struct {
	Short bool
	Out   io.Writer
}

// Run executes the version command logic.
func (o *VersionOptions) Run() error {
	if o.Short {
		fmt.Fprintf(o.Out, "%s\n", Version)
		return nil
	}

	fmt.Fprintf(o.Out, "driver-scanner version: %s\n", Version)
	fmt.Fprintf(o.Out, "Commit: %s\n", CommitHash)
	fmt.Fprintf(o.Out, "Built: %s\n", BuildDate)
	fmt.Fprintf(o.Out, "Go version: %s\n", runtime.Version())
	fmt.Fprintf(o.Out, "OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	return nil
}

// newVersionCommand creates the "version" subcommand.
func newVersionCommand() *cobra.Command {
	o := &VersionOptions{}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Example: `  # Show full version information
  driver-scanner version

  # Show only version number
  driver-scanner version --short`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Out = cmd.OutOrStdout()
			return o.Run()
		},
	}

	cmd.Flags().BoolVar(&o.Short, "short", false, "Print only the version number")

	return cmd
}
