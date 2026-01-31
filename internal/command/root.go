package command

import (
	"github.com/spf13/cobra"

	"github.com/gigiozzz/driver-scanner/internal/service"
)

// Debug is set by the --debug flag. Available globally for future logging setup.
var Debug bool

// NewRootCommand creates the root cobra command for driver-scanner.
func NewRootCommand(scanner service.Scanner, version string, goVersion string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "driver-scanner",
		Short: "Scan and list block devices with mount and filesystem information",
	}

	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "enable debug output")

	rootCmd.AddCommand(newScanCommand(scanner))
	rootCmd.AddCommand(newVersionCommand(version, goVersion))

	return rootCmd
}
