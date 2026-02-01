package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/gigiozzz/driver-scanner/internal/command"
	"github.com/gigiozzz/driver-scanner/internal/device"
	"github.com/gigiozzz/driver-scanner/internal/provider"
	"github.com/gigiozzz/driver-scanner/internal/service"
)

func main() {
	// Initialize logging from env vars before cobra runs.
	provider.InitLogging()

	log.Debug().Msg("starting driver-scanner")

	deviceProvider := device.NewLsblkProvider()
	mountProvider := device.NewSystemMountInfoProvider()
	scanner := service.NewDeviceScanner(deviceProvider, mountProvider)

	rootCmd := command.NewRootCommand(scanner)
	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("command failed")
		os.Exit(1)
	}
}
