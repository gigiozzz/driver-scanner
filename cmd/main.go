package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gigiozzz/driver-scanner/internal/command"
	"github.com/gigiozzz/driver-scanner/internal/device"
	"github.com/gigiozzz/driver-scanner/internal/service"
)

// version is set at build time via ldflags.
var version = "dev"

func main() {
	deviceProvider := device.NewLsblkProvider()
	mountProvider := device.NewSystemMountInfoProvider()
	scanner := service.NewDeviceScanner(deviceProvider, mountProvider)

	rootCmd := command.NewRootCommand(scanner, version, runtime.Version())
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
