package service

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog/log"

	"github.com/gigiozzz/driver-scanner/internal/device"
)

// ScanFilter holds the filter criteria for scanning devices.
type ScanFilter struct {
	// FSType filters by filesystem type (e.g. "ext4").
	FSType string
	// MinSize filters by minimum device size (e.g. "1G", "500M"). Parsed via go-humanize.
	MinSize string
	// MountPoint filters by mount point substring match.
	MountPoint string
}

// Scanner abstracts the device scanning logic.
type Scanner interface {
	// Scan returns block devices, optionally filtered by the given criteria.
	Scan(filter ScanFilter) ([]device.BlockDevice, error)
}

// DeviceScanner implements Scanner combining lsblk and mount information.
type DeviceScanner struct {
	deviceProvider device.BlockDeviceProvider
	mountProvider  device.MountInfoProvider
}

// NewDeviceScanner creates a new DeviceScanner with the given providers.
func NewDeviceScanner(deviceProvider device.BlockDeviceProvider, mountProvider device.MountInfoProvider) *DeviceScanner {
	return &DeviceScanner{
		deviceProvider: deviceProvider,
		mountProvider:  mountProvider,
	}
}

// Scan retrieves block devices from lsblk, enriches them with mount info,
// and applies filters.
func (s *DeviceScanner) Scan(filter ScanFilter) ([]device.BlockDevice, error) {
	log.Info().Msg("starting device scan")

	devices, err := s.deviceProvider.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	log.Info().Int("deviceCount", len(devices)).Msg("block devices discovered")

	mountEntries, err := s.mountProvider.GetMounts()
	if err != nil {
		return nil, fmt.Errorf("failed to get mount info: %w", err)
	}
	log.Debug().Int("count", len(mountEntries)).Msg("mount entries retrieved")

	mountsBySource := buildMountsBySource(mountEntries)
	enrichDevicesWithMountInfo(devices, mountsBySource)

	log.Info().
		Int("total", len(devices)).
		Str("fstype", filter.FSType).
		Str("minSize", filter.MinSize).
		Str("mountPoint", filter.MountPoint).
		Msg("applying filters")

	filtered, err := applyFilters(devices, filter)
	if err != nil {
		return nil, err
	}

	log.Info().
		Int("before", len(devices)).
		Int("after", len(filtered)).
		Msg("filtering complete")

	return filtered, nil
}

// buildMountsBySource creates a lookup map from device source path to MountEntry.
func buildMountsBySource(mountEntries []device.MountEntry) map[string]device.MountEntry {
	mountsBySource := make(map[string]device.MountEntry, len(mountEntries))
	for _, entry := range mountEntries {
		mountsBySource[entry.Source] = entry
	}
	return mountsBySource
}

// enrichDevicesWithMountInfo fills in mount point information from mountinfo
// for devices where lsblk did not report a mount point.
func enrichDevicesWithMountInfo(devices []device.BlockDevice, mountsBySource map[string]device.MountEntry) {
	for i := range devices {
		if devices[i].MountPoint != "" {
			continue
		}
		mountEntry, found := mountsBySource[devices[i].Path]
		if found {
			log.Debug().
				Str("device", devices[i].Path).
				Str("mountpoint", mountEntry.MountPoint).
				Msg("enriched device with mount info")
			devices[i].MountPoint = mountEntry.MountPoint
		}
	}
}

// applyFilters filters the device list based on the given ScanFilter criteria.
func applyFilters(devices []device.BlockDevice, filter ScanFilter) ([]device.BlockDevice, error) {
	var minSizeBytes uint64
	if filter.MinSize != "" {
		parsed, err := humanize.ParseBytes(filter.MinSize)
		if err != nil {
			return nil, fmt.Errorf("invalid min-size value %q: %w", filter.MinSize, err)
		}
		minSizeBytes = parsed
	}

	result := make([]device.BlockDevice, 0, len(devices))
	for _, dev := range devices {
		if filter.FSType != "" && !strings.EqualFold(dev.FSType, filter.FSType) {
			log.Debug().Str("device", dev.Path).Str("fstype", dev.FSType).Msg("filtered out by fstype")
			continue
		}
		if minSizeBytes > 0 && dev.DeviceSizeBytes < minSizeBytes {
			log.Debug().Str("device", dev.Path).Uint64("size", dev.DeviceSizeBytes).Msg("filtered out by min-size")
			continue
		}
		if filter.MountPoint != "" && !strings.Contains(dev.MountPoint, filter.MountPoint) {
			log.Debug().Str("device", dev.Path).Str("mountpoint", dev.MountPoint).Msg("filtered out by mount-point")
			continue
		}
		result = append(result, dev)
	}
	return result, nil
}
