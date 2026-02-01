package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/gigiozzz/driver-scanner/internal/device"
	"github.com/gigiozzz/driver-scanner/internal/service"
)

// newScanCommand creates the "scan" subcommand.
func newScanCommand(scanner service.Scanner) *cobra.Command {
	var filter service.ScanFilter

	cmd := &cobra.Command{
		Use:   "scan",
		Short: "Scan block devices and display their information",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().
				Str("fstype", filter.FSType).
				Str("minSize", filter.MinSize).
				Str("mountPoint", filter.MountPoint).
				Msg("scan command invoked")

			processedFilter, err := buildScanFilter(filter)
			if err != nil {
				return err
			}

			if err := validateScanFilter(processedFilter); err != nil {
				return err
			}

			return runScan(scanner, processedFilter)
		},
	}

	cmd.Flags().StringVar(&filter.FSType, "fstype", "", "filter by filesystem type (e.g. ext4)")
	cmd.Flags().StringVar(&filter.MinSize, "min-size", "", "filter by minimum device size (e.g. 1G, 500M)")
	cmd.Flags().StringVar(&filter.MountPoint, "mount-point", "", "filter by mount point (substring match)")

	return cmd
}

// buildScanFilter processes and normalizes filter input from CLI flags.
// Currently returns the filter as-is; extend here for input transformations.
func buildScanFilter(raw service.ScanFilter) (service.ScanFilter, error) {
	return raw, nil
}

// validateScanFilter validates the filter values before executing the scan.
func validateScanFilter(filter service.ScanFilter) error {
	if filter.FSType != "" {
		log.Debug().Str("fstype", filter.FSType).Msg("validating filesystem type")
		supportedTypes, err := readSupportedFileSystems()
		if err != nil {
			log.Debug().Err(err).Msg("cannot read supported filesystems")
			return fmt.Errorf("cannot validate fstype: %w", err)
		}
		if !supportedTypes[strings.ToLower(filter.FSType)] {
			log.Debug().Str("fstype", filter.FSType).Msg("unsupported filesystem type")
			return fmt.Errorf("unsupported filesystem type %q, supported: %s",
				filter.FSType, joinMapKeys(supportedTypes))
		}
		log.Debug().Str("fstype", filter.FSType).Msg("filesystem type is valid")
	}

	if filter.MinSize != "" {
		log.Debug().Str("minSize", filter.MinSize).Msg("validating min-size")
		if _, err := humanize.ParseBytes(filter.MinSize); err != nil {
			log.Debug().Str("minSize", filter.MinSize).Err(err).Msg("invalid min-size value")
			return fmt.Errorf("invalid min-size value %q: %w", filter.MinSize, err)
		}
	}

	return nil
}

// readSupportedFileSystems reads /proc/filesystems and returns a set of supported types.
func readSupportedFileSystems() (map[string]bool, error) {
	log.Debug().Msg("reading /proc/filesystems")

	file, err := os.Open("/proc/filesystems")
	if err != nil {
		return nil, fmt.Errorf("failed to open /proc/filesystems: %w", err)
	}
	defer file.Close()

	supported := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Each line is either "nodev\t<fstype>" or "\t<fstype>"
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		fsType := fields[len(fields)-1]
		supported[strings.ToLower(fsType)] = true
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read /proc/filesystems: %w", err)
	}

	log.Debug().Int("count", len(supported)).Msg("supported filesystems loaded")
	return supported, nil
}

// joinMapKeys returns a comma-separated string of map keys.
func joinMapKeys(m map[string]bool) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

// runScan executes the scan and prints the results as a table.
func runScan(scanner service.Scanner, filter service.ScanFilter) error {
	devices, err := scanner.Scan(filter)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	log.Info().Int("deviceCount", len(devices)).Msg("scan complete")
	if len(devices) == 0 {
		log.Warn().Msg("no devices matched the filter criteria")
	}

	printDeviceTable(devices)
	return nil
}

// printDeviceTable prints the device list in a formatted table.
func printDeviceTable(devices []device.BlockDevice) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "UUID\tSERIAL\tDEVICE\tFSTYPE\tTYPE\tMOUNTPOINT\tSIZE\tFS SIZE\tFS AVAIL")
	fmt.Fprintln(w, "----\t------\t------\t------\t----\t----------\t----\t-------\t--------")

	for _, dev := range devices {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			valueOrDash(dev.UUID),
			valueOrDash(dev.Serial),
			dev.Path,
			valueOrDash(dev.FSType),
			dev.Type,
			valueOrDash(dev.MountPoint),
			valueOrDash(dev.DeviceSize),
			valueOrDash(dev.FileSystemSize),
			valueOrDash(dev.FileSystemAvail),
		)
	}

	w.Flush()
}

// valueOrDash returns the value if non-empty, otherwise "-".
func valueOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}
