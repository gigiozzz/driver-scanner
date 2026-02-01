package command

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/gigiozzz/driver-scanner/internal/provider"
	"github.com/gigiozzz/driver-scanner/internal/service"
)

var (
	// debug is set by the --debug flag.
	debug bool
	// verbose is set by the -v flag.
	verbose bool
)

// NewRootCommand creates the root cobra command for driver-scanner.
func NewRootCommand(scanner service.Scanner, version string, goVersion string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "driver-scanner",
		Short: "Scan and list block devices with mount and filesystem information",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Phase 2: Adjust log level based on CLI flags.
			provider.SetLevelFromFlags(debug, verbose)
			log.Debug().
				Bool("debug", debug).
				Bool("verbose", verbose).
				Msg("log level configured from flags")
		},
	}

	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")

	rootCmd.AddCommand(newScanCommand(scanner))
	rootCmd.AddCommand(newVersionCommand(version, goVersion))

	return rootCmd
}
