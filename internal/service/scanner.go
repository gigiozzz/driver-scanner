package service

import (
	"fmt"
	"strings"

	"github.com/dustin/go-humanize"

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
	devices, err := s.deviceProvider.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	mountEntries, err := s.mountProvider.GetMounts()
	if err != nil {
		return nil, fmt.Errorf("failed to get mount info: %w", err)
	}

	mountsBySource := buildMountsBySource(mountEntries)
	enrichDevicesWithMountInfo(devices, mountsBySource)

	filtered, err := applyFilters(devices, filter)
	if err != nil {
		return nil, err
	}

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
			continue
		}
		if minSizeBytes > 0 && dev.DeviceSizeBytes < minSizeBytes {
			continue
		}
		if filter.MountPoint != "" && !strings.Contains(dev.MountPoint, filter.MountPoint) {
			continue
		}
		result = append(result, dev)
	}
	return result, nil
}
