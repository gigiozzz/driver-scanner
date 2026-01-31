package device

// BlockDevice represents the parsed output of lsblk combined with mount information.
// It is used as the domain DTO to carry block device data across layers.
type BlockDevice struct {
	// Name is the kernel device name (e.g. "sda", "sda1").
	Name string `json:"name"`
	// Path is the full path to the device node (e.g. "/dev/sda").
	Path string `json:"path"`
	// UUID is the filesystem UUID assigned to the device.
	UUID string `json:"uuid"`
	// Serial is the disk serial number.
	Serial string `json:"serial"`
	// FSType is the filesystem type (e.g. "ext4", "xfs", "ntfs"). Empty if unformatted.
	FSType string `json:"fstype"`
	// Type is the device type (e.g. "disk", "part", "loop").
	Type string `json:"type"`
	// Label is the filesystem label, if set.
	Label string `json:"label"`
	// MountPoint is the path where the device is mounted. Empty if not mounted.
	MountPoint string `json:"mountpoint"`
	// DeviceSize is the total physical size of the block device in human-readable format (e.g. "1.0 TB").
	DeviceSize string `json:"deviceSize"`
	// DeviceSizeBytes is the total physical size of the block device in bytes.
	DeviceSizeBytes uint64 `json:"deviceSizeBytes"`
	// FileSystemSize is the total size of the filesystem in human-readable format. Empty if not mounted.
	FileSystemSize string `json:"fileSystemSize"`
	// FileSystemSizeBytes is the total size of the filesystem in bytes. Zero if not mounted.
	FileSystemSizeBytes uint64 `json:"fileSystemSizeBytes"`
	// FileSystemAvail is the available free space in human-readable format. Empty if not mounted.
	FileSystemAvail string `json:"fileSystemAvail"`
	// FileSystemAvailBytes is the available free space in bytes. Zero if not mounted.
	FileSystemAvailBytes uint64 `json:"fileSystemAvailBytes"`
}
