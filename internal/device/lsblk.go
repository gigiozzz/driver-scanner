package device

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/dustin/go-humanize"
)

// lsblkOutput maps the top-level JSON structure returned by lsblk --json.
// lsblk wraps the device array inside a "blockdevices" key.
type lsblkOutput struct {
	BlockDevices []lsblkDevice `json:"blockdevices"`
}

// lsblkDevice maps a single device entry from lsblk JSON output.
// With -b flag, size fields are returned as numeric bytes in JSON.
// JSON null values are unmarshalled to Go zero values (0 for uint64, "" for string).
type lsblkDevice struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	UUID       string `json:"uuid"`
	Serial     string `json:"serial"`
	FSType     string `json:"fstype"`
	Type       string `json:"type"`
	Label      string `json:"label"`
	MountPoint string `json:"mountpoint"`
	Size       uint64 `json:"size"`
	FSSize     uint64 `json:"fssize"`
	FSAvail    uint64 `json:"fsavail"`
}

// BlockDeviceProvider abstracts the retrieval of block device information.
type BlockDeviceProvider interface {
	// List returns all block devices detected by the system.
	List() ([]BlockDevice, error)
}

// LsblkProvider implements BlockDeviceProvider by executing the lsblk command.
type LsblkProvider struct{}

// NewLsblkProvider creates a new LsblkProvider.
func NewLsblkProvider() *LsblkProvider {
	return &LsblkProvider{}
}

// List executes lsblk with -b (bytes) and returns the parsed block devices.
func (l *LsblkProvider) List() ([]BlockDevice, error) {
	out, err := exec.Command(
		"lsblk", "--json", "-b",
		"-o", "NAME,PATH,UUID,SERIAL,FSTYPE,TYPE,LABEL,MOUNTPOINT,SIZE,FSSIZE,FSAVAIL",
	).Output()
	if err != nil {
		return nil, fmt.Errorf("lsblk execution failed: %w", err)
	}

	var raw lsblkOutput
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, fmt.Errorf("lsblk JSON parsing failed: %w", err)
	}

	devices := make([]BlockDevice, 0, len(raw.BlockDevices))
	for _, entry := range raw.BlockDevices {
		devices = append(devices, BlockDevice{
			Name:                 entry.Name,
			Path:                 entry.Path,
			UUID:                 entry.UUID,
			Serial:               entry.Serial,
			FSType:               entry.FSType,
			Type:                 entry.Type,
			Label:                entry.Label,
			MountPoint:           entry.MountPoint,
			DeviceSizeBytes:      entry.Size,
			DeviceSize:           humanizeBytes(entry.Size),
			FileSystemSizeBytes:  entry.FSSize,
			FileSystemSize:       humanizeBytes(entry.FSSize),
			FileSystemAvailBytes: entry.FSAvail,
			FileSystemAvail:      humanizeBytes(entry.FSAvail),
		})
	}
	return devices, nil
}

// humanizeBytes converts a byte count to a human-readable string.
// Returns empty string for 0 (typically means the value was not available).
func humanizeBytes(bytes uint64) string {
	if bytes == 0 {
		return ""
	}
	return humanize.IBytes(bytes)
}
