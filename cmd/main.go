package main

import (
	"os"
	"runtime"

	"github.com/rs/zerolog/log"

	"github.com/gigiozzz/driver-scanner/internal/command"
	"github.com/gigiozzz/driver-scanner/internal/device"
	"github.com/gigiozzz/driver-scanner/internal/provider"
	"github.com/gigiozzz/driver-scanner/internal/service"
)

// version is set at build time via ldflags.
var version = "dev"

func main() {
	// Initialize logging from env vars before cobra runs.
	provider.InitLogging()

	log.Debug().Str("version", version).Msg("starting driver-scanner")

	deviceProvider := device.NewLsblkProvider()
	mountProvider := device.NewSystemMountInfoProvider()
	scanner := service.NewDeviceScanner(deviceProvider, mountProvider)

	rootCmd := command.NewRootCommand(scanner, version, runtime.Version())
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("command failed")
		os.Exit(1)
	}
}
